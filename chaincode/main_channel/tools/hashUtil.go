package tools

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateHash 生成哈希值
func GenerateHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateRandomHash 生成随机哈希值
func GenerateRandomHash() string {
	// 创建一个包含16个字节的随机数
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return GenerateHash(fmt.Sprintf("%d", time.Now().UnixNano()))
	}

	// 使用SHA256对随机字节进行哈希处理
	return GenerateHash(hex.EncodeToString(randomBytes))
}

// GenerateRandomHashWithPrefix 生成带前缀的随机哈希值
func GenerateRandomHashWithPrefix(prefix string) string {
	return prefix + "-" + GenerateRandomHash()
}
