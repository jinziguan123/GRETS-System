package db

import (
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

// 全局数据库实例
var GlobalDB *BoltDB

// BoltDB结构体，封装bolt.DB
type BoltDB struct {
	db *bolt.DB
}

// 初始化BoltDB，如果数据库不存在，会自动创建
func InitBoltDB(dbPath string) (*BoltDB, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("打开BoltDB失败: %v", err)
	}

	// 初始化全局实例
	GlobalDB = &BoltDB{db: db}
	return GlobalDB, nil
}

// 关闭数据库
func (b *BoltDB) Close() error {
	return b.db.Close()
}

// 创建桶(表)
func (b *BoltDB) CreateBucketIfNotExists(bucketName string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("创建桶失败: %v", err)
		}
		return nil
	})
}

// 插入或更新记录
func (b *BoltDB) Put(bucketName, key string, value interface{}) error {
	// 将value转换为JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("桶[%s]不存在", bucketName)
		}
		return bucket.Put([]byte(key), jsonData)
	})
}

// 获取记录
func (b *BoltDB) Get(bucketName, key string, result interface{}) error {
	return b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("桶[%s]不存在", bucketName)
		}

		data := bucket.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("记录不存在")
		}

		return json.Unmarshal(data, result)
	})
}

// 删除记录
func (b *BoltDB) Delete(bucketName, key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("桶[%s]不存在", bucketName)
		}
		return bucket.Delete([]byte(key))
	})
}

// 获取桶中的所有记录
func (b *BoltDB) GetAll(bucketName string) (map[string][]byte, error) {
	result := make(map[string][]byte)

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("桶[%s]不存在", bucketName)
		}

		return bucket.ForEach(func(k, v []byte) error {
			// 创建k的副本，因为k在迭代外部可能不再有效
			key := make([]byte, len(k))
			copy(key, k)

			// 创建v的副本
			value := make([]byte, len(v))
			copy(value, v)

			result[string(key)] = value
			return nil
		})
	})

	return result, err
}

// 批量查询记录
func (b *BoltDB) Query(bucketName string, filter func(k, v []byte) bool, result interface{}) error {
	var matches []json.RawMessage

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("桶[%s]不存在", bucketName)
		}

		return bucket.ForEach(func(k, v []byte) error {
			if filter == nil || filter(k, v) {
				matches = append(matches, json.RawMessage(v))
			}
			return nil
		})
	})

	if err != nil {
		return err
	}

	// 将匹配的数据编码为JSON数组
	data, err := json.Marshal(matches)
	if err != nil {
		return err
	}

	// 将JSON数组解码到result
	return json.Unmarshal(data, result)
}

// 检查键是否存在
func (b *BoltDB) Exists(bucketName, key string) (bool, error) {
	var exists bool

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("桶[%s]不存在", bucketName)
		}

		value := bucket.Get([]byte(key))
		exists = value != nil
		return nil
	})

	return exists, err
}
