package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// GenerateHash 生成哈希值
func GenerateHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
