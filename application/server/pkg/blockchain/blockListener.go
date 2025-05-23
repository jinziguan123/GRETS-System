package blockchain

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"grets_server/config"
	"grets_server/pkg/utils"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/boltdb/bolt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/msp"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

const (
	_BlocksBucket      = "blocks"
	_LatestBucket      = "latestBucket"
	_RetryInterval     = 3 * time.Second
	_BlockQueueSize    = 1000                   // 内部区块处理队列大小
	_BlockBatchSize    = 100                    // 一次数据库事务保存的区块数量
	_BlockBatchTimeout = 200 * time.Millisecond // 批量保存的最大等待时间
)

// 添加区块保存任务结构体
type blockToSave struct {
	channelName string
	orgName     string
	block       *common.Block
}

// BlockHeader 区块头
type BlockHeader struct {
	Number       *big.Int
	PreviousHash []byte
	DataHash     []byte
}

// BlockData 区块数据
type BlockData struct {
	BlockNumber uint64    `json:"blockNumber"`
	BlockHash   string    `json:"blockHash"`
	DataHash    string    `json:"dataHash"`
	PrevHash    string    `json:"prevHash"`
	TxCount     int       `json:"txCount"`
	SaveTime    time.Time `json:"saveTime"`
	Data        [][]byte  `json:"data"`
	ChannelName string    `json:"channelName"`
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
	mainNetworks      map[string]*client.Network
	subNetworks       map[string]map[string]*client.Network
	ctx               context.Context
	cancel            context.CancelFunc
	dataDir           string
	blockProcessQueue chan blockToSave // 区块处理队列
	wg                sync.WaitGroup   // 用于等待保存协程完成
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
			mainNetworks:      make(map[string]*client.Network),
			subNetworks:       make(map[string]map[string]*client.Network),
			db:                db,
			dataDir:           dataDir,
			ctx:               ctx,
			cancel:            cancel,
			blockProcessQueue: make(chan blockToSave, _BlockQueueSize),
		}

		// 启动专门的区块保存协程
		listener.wg.Add(1)
		go listener.runBlockSaver()
	})
	utils.Log.Info("初始化区块链监听器完成")

	return initErr
}

// GetBlockListener 获取区块监听器实例
func GetBlockListener() *blockListener {
	return listener
}

// Stop 优雅关闭区块监听器及其保存协程
func (l *blockListener) Stop() {
	utils.Log.Info("正在停止区块链监听器...")
	l.cancel()                 // 通知所有监听协程和保存器停止
	close(l.blockProcessQueue) // 关闭队列，通知保存器不再有新区块
	l.wg.Wait()                // 等待保存协程处理完成剩余区块
	if l.db != nil {
		l.db.Close()
	}
	utils.Log.Info("区块链监听器已停止")
}

// addMainNetwork 添加主通道网络
func addMainNetwork(orgName string, network *client.Network) error {
	if listener == nil {
		return fmt.Errorf("区块监听器未初始化")
	}

	listener.Lock()
	defer listener.Unlock()

	listener.mainNetworks[orgName] = network
	go listener.startMainBlockListener(orgName)

	return nil
}

// addSubNetwork 添加子通道网络
func addSubNetwork(subChannelName string, orgName string, network *client.Network) error {
	if listener == nil {
		return fmt.Errorf("区块监听器未初始化")
	}

	listener.Lock()
	defer listener.Unlock()

	// 确保subChannelName对应的map已初始化
	if listener.subNetworks[subChannelName] == nil {
		listener.subNetworks[subChannelName] = make(map[string]*client.Network)
	}

	listener.subNetworks[subChannelName][orgName] = network
	go listener.startSubBlockListener(subChannelName, orgName)

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

// startMainNetworkListener 开始监听主通道区块
func (l *blockListener) startMainNetworkListener(mainNetwork *client.Network, orgName string) {
	retryCount := 0

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

		events, err := mainNetwork.BlockEvents(l.ctx, client.WithStartBlock(startBlock))
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
					fmt.Printf("通道[%s]组织[%s]的区块事件监听中断（已重试%d次），准备重试...\n", config.GlobalConfig.Fabric.MainChannelName, orgName, retryCount)
					retryCount++
					select {
					case <-l.ctx.Done():
						return
					case <-time.After(_RetryInterval):
						break
					}
					goto RETRY
				}
				l.enqueueBlockForSaving(config.GlobalConfig.Fabric.MainChannelName, orgName, block)
			}
		}

	RETRY:
		continue
	}
}

// startSubNetworkListener 开始监听子通道区块
func (l *blockListener) startSubNetworkListener(subNetwork *client.Network, orgName string, subChannelName string) {
	utils.Log.Info(fmt.Sprintf("开始监听组织[%s]的子通道[%s]区块", orgName, subChannelName))

	retryCount := 0

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

		events, err := subNetwork.BlockEvents(l.ctx, client.WithStartBlock(startBlock))
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
					fmt.Printf("通道[%s]组织[%s]的区块事件监听中断（已重试%d次），准备重试...\n", subChannelName, orgName, retryCount)
					retryCount++
					select {
					case <-l.ctx.Done():
						return
					case <-time.After(_RetryInterval):
						break
					}
					goto RETRY
				}
				l.enqueueBlockForSaving(subChannelName, orgName, block)
			}
		}

	RETRY:
		continue
	}
}

// startBlockListener 开始监听区块
func (l *blockListener) startMainBlockListener(orgName string) {
	utils.Log.Info(fmt.Sprintf("开始监听组织[%s]的区块", orgName))

	mainNetwork := l.mainNetworks[orgName]
	if mainNetwork == nil {
		fmt.Printf("组织[%s]的主通道网络未找到\n", orgName)
		return
	}

	go l.startMainNetworkListener(mainNetwork, orgName)
}

func (l *blockListener) startSubBlockListener(subChannelName string, orgName string) {
	utils.Log.Info(fmt.Sprintf("开始监听组织[%s]的子通道[%s]区块", orgName, subChannelName))

	subNetwork := l.subNetworks[subChannelName][orgName]
	if subNetwork == nil {
		fmt.Printf("组织[%s]的子通道[%s]网络未找到\n", orgName, subChannelName)
		return
	}

	go l.startSubNetworkListener(subNetwork, orgName, subChannelName)
}

// saveBlock 保存区块 (已废弃，保留以兼容旧代码，应使用enqueueBlockForSaving)
func (l *blockListener) saveBlock(channelName string, orgName string, block *common.Block) {
	// 调用新方法处理，不再直接保存
	l.enqueueBlockForSaving(channelName, orgName, block)
}

// GetBlockByNumber 根据组织名和区块号查询区块
func (l *blockListener) GetBlockByNumber(channelName string, orgName string, blockNum uint64) (*BlockData, error) {
	var blockData BlockData

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(_BlocksBucket))
		blockKey := fmt.Sprintf("%s_%s_%d", channelName, orgName, blockNum)
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

// BlockTransaction 区块交易
type BlockTransactionDetail struct {
	TransactionID         string `json:"transactionID"`         // 交易ID
	Creator               string `json:"creator"`               // 创建者地址
	TransactionTimestamp  string `json:"transactionTimestamp"`  // 交易时间戳
	ChainCodeFunctionName string `json:"chainCodeFunctionName"` // 链码函数名称
}

// GetBlocksByChannelAndOrg 分页查询组织的区块列表（按区块号降序）
func (l *blockListener) GetBlocksByChannelAndOrg(channelName string, orgName string, pageNum int, pageSize int) (*BlockQueryResult, error) {
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
			blockKey := fmt.Sprintf("%s_%s_%d", channelName, orgName, i)
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

// GetAllBlocksByChannelAndOrg 获取所有区块
func (l *blockListener) GetAllBlocksByChannelAndOrg(channelName string, orgName string) ([]*BlockData, error) {
	var result []*BlockData

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
		total := int(latestBlock.BlockNum) + 1

		// 收集区块数据
		blocks := make([]*BlockData, 0, total)
		for i := total - 1; i >= 0; i-- {
			blockKey := fmt.Sprintf("%s_%s_%d", channelName, orgName, i)
			data := b.Get([]byte(blockKey))
			if data != nil {
				var block BlockData
				if err := json.Unmarshal(data, &block); err != nil {
					return fmt.Errorf("区块数据反序列化失败: %v", err)
				}
				blocks = append(blocks, &block)
			}
		}

		result = blocks

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("查询区块失败: %v", err)
	}

	return result, nil
}

func (l *blockListener) GetEnvelopeListFromBlock(block *common.Block) ([]*common.Envelope, error) {
	var envelopes []*common.Envelope
	for _, envBytes := range block.Data.Data {
		envelope := &common.Envelope{}
		if err := proto.Unmarshal(envBytes, envelope); err != nil {
			return nil, fmt.Errorf("getEnvelopeFromBlock error: %v", err)
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes, nil
}

func (l *blockListener) GetEnvelopeListFromBoltBlockData(blockData *BlockData) ([]*common.Envelope, error) {
	var envelopes []*common.Envelope
	for _, envBytes := range blockData.Data {
		envelope := &common.Envelope{}
		if err := proto.Unmarshal(envBytes, envelope); err != nil {
			return nil, fmt.Errorf("getEnvelopeFromBlock error: %v", err)
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes, nil
}

func (l *blockListener) GetChannelHeaderFromEnvelope(env *common.Envelope) (*common.ChannelHeader, error) {
	payload := &common.Payload{}
	if err := proto.Unmarshal(env.Payload, payload); err != nil {
		return nil, fmt.Errorf("getChannelHeaderFromEnvelope error: %v", err)
	}
	channelHeader := &common.ChannelHeader{}
	if err := proto.Unmarshal(payload.Header.ChannelHeader, channelHeader); err != nil {
		return nil, fmt.Errorf("getChannelHeaderFromEnvelope error: %v", err)
	}
	return channelHeader, nil
}

func (l *blockListener) GetSignatureHeaderFromEnvelope(env *common.Envelope) (*common.SignatureHeader, error) {
	payload := &common.Payload{}
	if err := proto.Unmarshal(env.Payload, payload); err != nil {
		return nil, fmt.Errorf("getSignatureHeaderFromEnvelope error: %v", err)
	}
	signatureHeader := &common.SignatureHeader{}
	if err := proto.Unmarshal(payload.Header.SignatureHeader, signatureHeader); err != nil {
		return nil, fmt.Errorf("getSignatureHeaderFromEnvelope error: %v", err)
	}
	return signatureHeader, nil
}

// parseCreatorCertificate 解析 PEM 编码的证书并提取用户信息
func parseCreatorCertificate(certBytes []byte) (string, error) {
	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %v", err)
	}

	// 使用公钥派生地址
	pubKey := cert.PublicKey.(*ecdsa.PublicKey).X.Bytes()
	hash := sha256.Sum256(pubKey)
	return hex.EncodeToString(hash[:]), nil
}

func (l *blockListener) GetTransactionDetailListFromEnvelopeList(envList []*common.Envelope) ([]*BlockTransactionDetail, error) {
	transactionDetailList := make([]*BlockTransactionDetail, 0)
	// 遍历envList
	for _, env := range envList {
		// 首先获取channelHeader来拿到交易时间和txID
		channelHeader, err := l.GetChannelHeaderFromEnvelope(env)
		if err != nil {
			return nil, fmt.Errorf("获取channelHeader失败: %v", err)
		}

		signatureHeader, err := l.GetSignatureHeaderFromEnvelope(env)
		if err != nil {
			return nil, fmt.Errorf("获取signatureHeader失败: %v", err)
		}
		creator := &msp.SerializedIdentity{}
		if err := proto.Unmarshal(signatureHeader.Creator, creator); err != nil {
			return nil, fmt.Errorf("解析creator失败: %v", err)
		}
		creatorInfo, err := parseCreatorCertificate(creator.IdBytes)
		if err != nil {
			return nil, fmt.Errorf("解析creator证书失败: %v", err)
		}

		payload := &common.Payload{}
		if err := proto.Unmarshal(env.Payload, payload); err != nil {
			return nil, fmt.Errorf("解析payload失败: %v", err)
		}
		txPayload := &peer.Transaction{}
		if err := proto.Unmarshal(payload.Data, txPayload); err != nil {
			return nil, fmt.Errorf("解析TransactionAction失败: %v", err)
		}
		for _, action := range txPayload.Actions {
			chaincodeActionPayload := &peer.ChaincodeActionPayload{}
			if err := proto.Unmarshal(action.Payload, chaincodeActionPayload); err != nil {
				return nil, fmt.Errorf("解析ChaincodeActionPayload失败: %v", err)
			}
			chaincodeProposalPayload := &peer.ChaincodeProposalPayload{}
			if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, chaincodeProposalPayload); err != nil {
				return nil, fmt.Errorf("解析ChaincodeProposalPayload失败: %v", err)
			}
			chaincodeInvocationSpec := &peer.ChaincodeInvocationSpec{}
			if err := proto.Unmarshal(chaincodeProposalPayload.Input, chaincodeInvocationSpec); err != nil {
				return nil, fmt.Errorf("解析ChaincodeInvocationSpec失败: %v", err)
			}
			transactionDetail := &BlockTransactionDetail{
				TransactionID:         channelHeader.TxId,
				Creator:               creatorInfo,
				TransactionTimestamp:  channelHeader.Timestamp.AsTime().Format(time.RFC3339),
				ChainCodeFunctionName: string(chaincodeInvocationSpec.ChaincodeSpec.Input.Args[0]),
			}
			transactionDetailList = append(transactionDetailList, transactionDetail)
		}
	}
	return transactionDetailList, nil
}

// runBlockSaver 是一个长期运行的协程，批量保存区块到数据库
func (l *blockListener) runBlockSaver() {
	defer l.wg.Done()
	batch := make([]blockToSave, 0, _BlockBatchSize)
	ticker := time.NewTicker(_BlockBatchTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done(): // 上下文取消，监听器正在关闭
			// 处理批处理队列中的剩余块，然后退出
			if len(batch) > 0 {
				l.persistBatch(batch)
				batch = make([]blockToSave, 0, _BlockBatchSize) // 清空批处理
			}
			// 消耗队列
			for item := range l.blockProcessQueue { // 队列关闭且为空时将自动退出
				batch = append(batch, item)
				if len(batch) >= _BlockBatchSize {
					l.persistBatch(batch)
					batch = make([]blockToSave, 0, _BlockBatchSize)
				}
			}
			if len(batch) > 0 { // 处理最后可能剩余的区块
				l.persistBatch(batch)
			}
			utils.Log.Info("区块保存器已停止")
			return

		case item, ok := <-l.blockProcessQueue:
			if !ok { // 通道关闭，表示Stop()被调用且队列已耗尽
				if len(batch) > 0 {
					l.persistBatch(batch)
				}
				utils.Log.Info("区块保存器处理完剩余任务后停止 (通道关闭)")
				return // 退出协程
			}
			batch = append(batch, item)
			if len(batch) >= _BlockBatchSize {
				l.persistBatch(batch)
				batch = make([]blockToSave, 0, _BlockBatchSize) // 重置批处理
				// 重置定时器，防止批处理快速填充时立即触发
				ticker.Reset(_BlockBatchTimeout)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				l.persistBatch(batch)
				batch = make([]blockToSave, 0, _BlockBatchSize) // 重置批处理
			}
		}
	}
}

// persistBatch 将批量区块保存到数据库中的单个事务中
func (l *blockListener) persistBatch(batch []blockToSave) {
	if len(batch) == 0 {
		return
	}

	latestBlockUpdates := make(map[string]LatestBlock) // orgName -> LatestBlock

	err := l.db.Update(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(_BlocksBucket))
		if blocksBucket == nil {
			return fmt.Errorf("bucket %s not found", _BlocksBucket)
		}
		latestBucket := tx.Bucket([]byte(_LatestBucket))
		if latestBucket == nil {
			return fmt.Errorf("bucket %s not found", _LatestBucket)
		}

		for _, item := range batch {
			if item.block == nil {
				utils.Log.Warn("尝试保存空区块，已跳过")
				continue
			}
			blockNum := item.block.GetHeader().GetNumber()

			// 计算区块哈希
			blockHeader := BlockHeader{
				Number:       new(big.Int).SetUint64(blockNum),
				PreviousHash: item.block.GetHeader().GetPreviousHash(),
				DataHash:     item.block.GetHeader().GetDataHash(),
			}
			headerBytes, err := asn1.Marshal(blockHeader)
			if err != nil {
				utils.Log.Error(fmt.Sprintf("区块头序列化失败 (区块 %d): %v", blockNum, err))
				continue // 跳过这个区块，或者更健壮地处理错误
			}
			blockHash := sha256.Sum256(headerBytes)

			// 从区块中的第一个交易中提取SaveTime
			var saveTime time.Time
			if len(item.block.GetData().GetData()) > 0 {
				envelopes, err := l.GetEnvelopeListFromBlock(item.block) // 重用现有的解析逻辑
				if err == nil && len(envelopes) > 0 {
					channelHeader, err := l.GetChannelHeaderFromEnvelope(envelopes[0])
					if err == nil {
						saveTime = channelHeader.Timestamp.AsTime()
					} else {
						utils.Log.Warn(fmt.Sprintf("无法从区块 %d 的信封中获取ChannelHeader: %v", blockNum, err))
						saveTime = time.Now().UTC() // 备用
					}
				} else {
					utils.Log.Warn(fmt.Sprintf("无法从区块 %d 中获取信封列表: %v", blockNum, err))
					saveTime = time.Now().UTC() // 备用
				}
			} else {
				// 配置区块可能没有相同方式的Data.Data中的交易
				// 假设使用当前时间或尝试从块元数据获取
				utils.Log.Info(fmt.Sprintf("区块 %d 没有数据交易, 使用当前时间作为SaveTime", blockNum))
				saveTime = time.Now().UTC()
			}

			blockData := BlockData{
				BlockNumber: blockNum,
				BlockHash:   fmt.Sprintf("%x", blockHash[:]),
				DataHash:    fmt.Sprintf("%x", item.block.GetHeader().GetDataHash()),
				PrevHash:    fmt.Sprintf("%x", item.block.GetHeader().GetPreviousHash()),
				TxCount:     len(item.block.GetData().GetData()),
				SaveTime:    saveTime, // 使用提取或备用时间
				Data:        item.block.GetData().GetData(),
				ChannelName: item.channelName,
			}

			blockKey := fmt.Sprintf("%s_%s_%d", item.channelName, item.orgName, blockNum)
			blockJson, err := json.Marshal(blockData)
			if err != nil {
				utils.Log.Error(fmt.Sprintf("区块数据序列化失败 (区块 %s): %v", blockKey, err))
				// 可能返回err以回滚事务，或记录并继续处理下一个区块
				return fmt.Errorf("区块 %s 数据序列化失败: %v", blockKey, err)
			}
			if err := blocksBucket.Put([]byte(blockKey), blockJson); err != nil {
				utils.Log.Error(fmt.Sprintf("保存区块数据失败 (区块 %s): %v", blockKey, err))
				return fmt.Errorf("保存区块 %s 数据失败: %v", blockKey, err)
			}

			// 跟踪当前批处理中该orgName的最新区块
			if lb, ok := latestBlockUpdates[item.orgName]; !ok || blockNum > lb.BlockNum {
				latestBlockUpdates[item.orgName] = LatestBlock{
					BlockNum: blockNum,
					SaveTime: time.Now(), // 保存批处理的时间
				}
			}
			utils.Log.Debug(fmt.Sprintf("通道[%s]组织[%s]区块[%d]已准备在批处理中保存", item.channelName, item.orgName, blockNum))
		}

		// 处理完批处理中的所有区块后，更新最新区块信息
		for orgName, latestBlock := range latestBlockUpdates {
			latestJson, err := json.Marshal(latestBlock)
			if err != nil {
				return fmt.Errorf("组织 %s 最新区块信息序列化失败: %v", orgName, err)
			}
			if err := latestBucket.Put([]byte(orgName), latestJson); err != nil {
				return fmt.Errorf("组织 %s 保存最新区块信息失败: %v", orgName, err)
			}
		}
		return nil
	})

	if err != nil {
		utils.Log.Error(fmt.Sprintf("批量保存区块失败 (共 %d 个区块): %v", len(batch), err))
		// 如有必要，为失败的批处理实现重试逻辑或死信队列
	} else {
		var firstBlockNum, lastBlockNum uint64
		if len(batch) > 0 {
			firstBlockNum = batch[0].block.GetHeader().GetNumber()
			lastBlockNum = batch[len(batch)-1].block.GetHeader().GetNumber()
		}
		utils.Log.Info(fmt.Sprintf("批量保存 %d 个区块成功 (范围 %d-%d)", len(batch), firstBlockNum, lastBlockNum))
	}
}

// enqueueBlockForSaving 将区块加入队列以进行保存
func (l *blockListener) enqueueBlockForSaving(channelName string, orgName string, block *common.Block) {
	if block == nil {
		return
	}

	item := blockToSave{
		channelName: channelName,
		orgName:     orgName,
		block:       block,
	}

	select {
	case l.blockProcessQueue <- item:
		// 成功入队
		blockNum := block.GetHeader().GetNumber()
		utils.Log.Debug(fmt.Sprintf("通道[%s]组织[%s]区块[%d]已入队等待保存", channelName, orgName, blockNum))
	case <-l.ctx.Done():
		utils.Log.Warn(fmt.Sprintf("区块监听器已停止，区块 %d (通道 %s, 组织 %s) 未能进入队列", block.GetHeader().GetNumber(), channelName, orgName))
	default:
		// 队列已满，表示保存协程处理速度跟不上
		// 记录错误，区块可能被丢弃或处理延迟
		utils.Log.Error(fmt.Sprintf("区块处理队列已满！通道[%s]组织[%s]区块[%d]可能被丢弃或处理延迟。",
			channelName, orgName, block.GetHeader().GetNumber()))
	}
}
