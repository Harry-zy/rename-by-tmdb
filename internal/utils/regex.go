package utils

import (
	"fmt"
	"strings"
)

// GenerateRangePattern 生成指定范围的正则表达式模式
func GenerateRangePattern(start, end int, digits int) string {
	// 确保范围有效，不生成负数集数
	if start < 1 {
		start = 1
	}
	if end < 1 {
		// 如果结束集数小于1，返回空模式
		return "()"
	}
	if start > end {
		// 如果开始集数大于结束集数，返回空模式
		return "()"
	}

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
