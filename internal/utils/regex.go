package utils

import (
	"fmt"
	"strings"
)

// GenerateRangePattern 生成指定范围的正则表达式模式
func GenerateRangePattern(start, end int, digits int) string {
	// 生成一个简单的范围模式
	var patterns []string
	for i := start; i <= end; i++ {
		if digits > 1 {
			// 使用固定位数（补0）
			patterns = append(patterns, fmt.Sprintf("%0*d", digits, i))
		} else {
			// 不补0
			patterns = append(patterns, fmt.Sprintf("%d", i))
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(patterns, "|"))
}
