package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringSlice []string

// Scan 解析 JSON 字符串到数组
func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("无法解析为字节数组")
	}
	return json.Unmarshal(bytes, s)
}

// Value 将数组序列化为 JSON 字符串
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}
