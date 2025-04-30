package tools

import "fmt"

// BuildQueryString 构建查询字符串
func BuildQueryString(conditions map[string]interface{}) string {
	query := "{ \"selector\": { "
	first := true
	for key, value := range conditions {
		if !first {
			query += ", "
		}
		query += "\"" + key + "\": " + FormatValue(value)
		first = false
	}
	query += " } }"
	return query
}

// FormatValue 格式化条件值
func FormatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return "\"" + v + "\""
	case float64:
		return fmt.Sprintf("%f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
