package blockchain

import (
	"context"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"grets_server/pkg/utils"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
)

const (
	_BlocksBucket = "blocks"
	_LatestBucket = "latestBucket"

	_RetryInterval = 3 * time.Second
)

// BlockData 区块数据
type BlockData struct {
	BlockNumber uint64    `json:"blockNumber"`
	BlockHash   string    `json:"blockHash"`
	DataHash    string    `json:"dataHash"`
	PrevHash    string    `json:"prevHash"`
	TxCount     int       `json:"txCount"`
	SaveTime    time.Time `json:"saveTime"`
}

// LatestBlock 最新区块信息
type LatestBlock struct {
	BlockNum uint64    `json:"blockNum"`
	SaveTime time.Time `json:"saveTime"`
}

// BlockListener 区块监听器
type blockListener struct {
	db *bolt.DB
	sync.RWMutex
	networks map[string]*client.Network
	ctx      context.Context
	cancel   context.CancelFunc
	dataDir  string
}

var (
	listener     *blockListener
	listenerOnce sync.Once
)

func initBlockchainListener(dataDir string) error {
	utils.Log.Info("初始化区块链监听器")
	var initErr error
	listenerOnce.Do(func() {
		// 创建数据目录
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			initErr = fmt.Errorf("创建数据目录失败: %v", err)
			return
		}

		// 打开数据库
		dbPath := filepath.Join(dataDir, "blocks.db")

		if _, err := os.Stat(dbPath); err == nil {
			// 文件存在，删除它
			if err := os.Remove(dbPath); err != nil {
				initErr = fmt.Errorf("删除已存在的数据库文件失败: %v", err)
				return
			}
		}

		db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 10 * time.Second})
		if err != nil {
			initErr = fmt.Errorf("打开数据库失败: %v", err)
			return
		}

		// 创建bucket
		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(_BlocksBucket))
			if err != nil {
				return fmt.Errorf("创建bucket失败: %v", err)
			}
			_, err = tx.CreateBucketIfNotExists([]byte(_LatestBucket))
			if err != nil {
				return fmt.Errorf("创建latest bucket失败: %v", err)
			}
			return nil
		}); err != nil {
			db.Close()
			initErr = fmt.Errorf("初始化数据库失败: %v", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		listener = &blockListener{
			networks: make(map[string]*client.Network),
			db:       db,
			dataDir:  dataDir,
			ctx:      ctx,
			cancel:   cancel,
		}
	})
	utils.Log.Info("初始化区块链监听器完成")

	return initErr
}

// GetBlockListener 获取区块监听器实例
func GetBlockListener() *blockListener {
	return listener
}

// addNetwork 添加网络
func addNetwork(orgName string, network *client.Network) error {
	if listener == nil {
		return fmt.Errorf("区块监听器未初始化")
	}

	listener.Lock()
	defer listener.Unlock()

	listener.networks[orgName] = network
	go listener.startBlockListener(orgName)

	return nil
}

// getLastBlockNum 获取最后保存的区块号
func (l *blockListener) getLastBlockNum(orgName string) (uint64, bool) {
	var lastBlock LatestBlock

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_LatestBucket))
		data := b.Get([]byte(orgName))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &lastBlock)
	})

	if err != nil {
		fmt.Printf("获取最后区块号失败：%v\n", err)
		return 0, false
	}

	// 如果没有数据，返回false表示是首次启动
	if lastBlock.BlockNum == 0 && lastBlock.SaveTime.IsZero() {
		return 0, false
	}

	return lastBlock.BlockNum, true
}

// startBlockListener 开始监听区块
func (l *blockListener) startBlockListener(orgName string) {
	utils.Log.Info(fmt.Sprintf("开始监听组织[%s]的区块", orgName))

	retryCount := 0
	network := l.networks[orgName]
	if network == nil {
		fmt.Printf("组织[%s]的网络未找到\n", orgName)
		return
	}
	for {
		lastBlockNum, exists := l.getLastBlockNum(orgName)
		var startBlock uint64
		if !exists {
			// 首次启动，从0开始
			startBlock = 0
		} else {
			// 已有数据，从下一个开始
			startBlock = lastBlockNum + 1
		}

		events, err := network.BlockEvents(l.ctx, client.WithStartBlock(startBlock))
		if err != nil {
			fmt.Printf("创建区块事件请求失败（已重试%d次）：%v\n", retryCount, err)
			retryCount++
			select {
			case <-l.ctx.Done():
				return
			case <-time.After(_RetryInterval):
				continue
			}
		}

		for {
			select {
			case <-l.ctx.Done():
				return
			case block, ok := <-events:
				if !ok {
					fmt.Printf("组织[%s]的区块事件监听中断（已重试%d次），准备重试...\n", orgName, retryCount)
					retryCount++
					select {
					case <-l.ctx.Done():
						return
					case <-time.After(_RetryInterval):
						break
					}
					goto RETRY
				}
				l.saveBlock(orgName, block)
			}
		}

	RETRY:
		continue
	}
}

// saveBlock 保存区块
func (l *blockListener) saveBlock(orgName string, block *common.Block) {
	if block == nil {
		return
	}

	blockNum := block.GetHeader().GetNumber()

	// 计算区块哈希
	blockHeader := struct {
		Number       *big.Int
		PreviousHash []byte
		DataHash     []byte
	}{
		Number:       new(big.Int).SetUint64(blockNum),
		PreviousHash: block.GetHeader().GetPreviousHash(),
		DataHash:     block.GetHeader().GetDataHash(),
	}
	headerBytes, err := asn1.Marshal(blockHeader)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("区块头序列化失败: %v", err))
		return
	}

	blockHash := sha256.Sum256(headerBytes)

	// 准备区块数据
	blockData := BlockData{
		BlockNumber: blockNum,
		BlockHash:   fmt.Sprintf("%x", blockHash[:]),
		DataHash:    fmt.Sprintf("%x", block.GetHeader().GetDataHash()),
		PrevHash:    fmt.Sprintf("%x", block.GetHeader().GetPreviousHash()),
		TxCount:     len(block.GetData().GetData()),
		SaveTime:    time.Now(),
	}

	// 事务保存区块链
	err = l.db.Update(func(tx *bolt.Tx) error {
		// 保存区块数据
		_BlocksBucket := tx.Bucket([]byte(_BlocksBucket))
		blockKey := fmt.Sprintf("%s_%d", orgName, blockNum)
		blockJson, err := json.Marshal(blockData)
		if err != nil {
			return fmt.Errorf("区块数据序列化失败: %v", err)
		}
		if err := _BlocksBucket.Put([]byte(blockKey), blockJson); err != nil {
			return fmt.Errorf("保存区块数据失败: %v", err)
		}

		// 更新最新区块信息
		_LatestBucket := tx.Bucket([]byte(_LatestBucket))
		latestBlock := LatestBlock{
			BlockNum: blockNum,
			SaveTime: time.Now(),
		}
		latestJson, err := json.Marshal(latestBlock)
		if err != nil {
			return fmt.Errorf("最新区块信息序列化失败: %v", err)
		}
		if err := _LatestBucket.Put([]byte(orgName), latestJson); err != nil {
			return fmt.Errorf("保存最新区块信息失败: %v", err)
		}

		return nil
	})

	if err != nil {
		utils.Log.Error(fmt.Sprintf("保存区块失败: %v", err))
		return
	}

	utils.Log.Info(fmt.Sprintf("组织[%s]区块[%d]保存成功", orgName, blockNum))
}

// GetBlockByNumber 根据组织名和区块号查询区块
func (l *blockListener) GetBlockByNumber(orgName string, blockNum uint64) (*BlockData, error) {
	var blockData BlockData

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_BlocksBucket))
		blockKey := fmt.Sprintf("%s_%d", orgName, blockNum)
		data := b.Get([]byte(blockKey))
		if data == nil {
			return fmt.Errorf("区块不存在")
		}
		return json.Unmarshal(data, &blockData)
	})

	if err != nil {
		return nil, fmt.Errorf("查询区块失败: %v", err)
	}

	return &blockData, nil
}

// BlockQueryResult 区块查询结果
type BlockQueryResult struct {
	Blocks   []*BlockData `json:"blocks"`   // 区块数据列表
	Total    int          `json:"total"`    // 总记录数
	PageSize int          `json:"pageSize"` // 每页大小
	PageNum  int          `json:"pageNum"`  // 当前页码
	HasMore  bool         `json:"hasMore"`  // 是否还有更多数据
}

// GetBlocksByOrg 分页查询组织的区块列表（按区块号降序）
func (l *blockListener) GetBlocksByOrg(orgName string, pageNum int, pageSize int) (*BlockQueryResult, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	var result BlockQueryResult
	result.PageNum = pageNum
	result.PageSize = pageSize

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_BlocksBucket))
		if b == nil {
			return fmt.Errorf("区块bucket不存在")
		}

		// 获取组织的最新区块号
		_LatestBucket := tx.Bucket([]byte(_LatestBucket))
		if _LatestBucket == nil {
			return fmt.Errorf("latest_blocks bucket不存在")
		}

		var latestBlock LatestBlock
		latestData := _LatestBucket.Get([]byte(orgName))
		if latestData == nil {
			return fmt.Errorf("组织[%s]没有区块数据", orgName)
		}
		if err := json.Unmarshal(latestData, &latestBlock); err != nil {
			return fmt.Errorf("最新区块信息反序列化失败: %v", err)
		}

		// 计算总记录数
		result.Total = int(latestBlock.BlockNum) + 1

		// 计算开始和结束区块号
		startIndex := result.Total - (pageNum * pageSize)
		endIndex := startIndex + pageSize
		if startIndex < 0 {
			startIndex = 0
		}
		if endIndex > result.Total {
			endIndex = result.Total
		}

		result.HasMore = startIndex > 0

		// 收集区块数据
		blocks := make([]*BlockData, 0, pageSize)
		for i := endIndex - 1; i >= startIndex; i-- {
			blockKey := fmt.Sprintf("%s_%d", orgName, i)
			data := b.Get([]byte(blockKey))
			if data != nil {
				var block BlockData
				if err := json.Unmarshal(data, &block); err != nil {
					return fmt.Errorf("区块数据反序列化失败: %v", err)
				}
				blocks = append(blocks, &block)
			}
		}

		result.Blocks = blocks
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("查询区块失败: %v", err)
	}

	return &result, nil
}
