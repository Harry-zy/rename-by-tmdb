package utils

import (
	"fmt"
	"strings"
)

// GenerateRangePattern 生成指定范围的正则表达式模式
func GenerateRangePattern(start, end int, digits int) string {
	// 生成一个简单的范围模式，确保使用固定位数
	var patterns []string
	for i := start; i <= end; i++ {
		patterns = append(patterns, fmt.Sprintf("%0*d", digits, i))
	}
	return fmt.Sprintf("(%s)", strings.Join(patterns, "|"))
}

// Join 连接字符串切片
func Join(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
