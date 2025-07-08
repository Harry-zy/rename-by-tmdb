package utils

import (
	"fmt"
)

// GenerateRangePattern 生成指定范围的正则表达式模式
func GenerateRangePattern(start, end int, digits int) string {
	// 确保开始和结束都是相同位数的字符串
	startStr := fmt.Sprintf("%0*d", digits, start)
	endStr := fmt.Sprintf("%0*d", digits, end)

	// 如果在同一个数量级内（比如都是062-077）
	if len(startStr) == len(endStr) {
		// 找到不同的部分
		var patterns []string

		// 处理每个可能的百位数（根据位数调整）
		startHundreds := start / 100
		endHundreds := end / 100

		for i := startHundreds; i <= endHundreds; i++ {
			// 计算当前组合的最小值和最大值
			currentStart := i * 100
			currentEnd := (i+1)*100 - 1

			// 调整范围边界
			if i == startHundreds {
				currentStart = start
			}
			if i == endHundreds {
				currentEnd = end
			}

			// 如果这个范围有效
			if currentStart <= currentEnd {
				// 提取十位和个位
				startTens := (currentStart % 100) / 10
				startOnes := currentStart % 10
				endTens := (currentEnd % 100) / 10
				endOnes := currentEnd % 10

				if startTens == endTens {
					// 如果十位相同
					patterns = append(patterns, fmt.Sprintf("%d[%d-%d]", i*10+startTens, startOnes, endOnes))
				} else {
					// 处理第一个十位（如果不是从0开始）
					if startOnes > 0 {
						patterns = append(patterns, fmt.Sprintf("%d[%d-9]", i*10+startTens, startOnes))
					}
					// 处理中间的十位
					if endTens-startTens > 1 || (startOnes == 0 && endTens > startTens) {
						nextStartTens := startTens
						if startOnes > 0 {
							nextStartTens = startTens + 1
						}
						for t := nextStartTens; t < endTens; t++ {
							patterns = append(patterns, fmt.Sprintf("%d[0-9]", i*10+t))
						}
					}
					// 处理最后一个十位（如果不是到9结束）
					if endOnes < 9 {
						patterns = append(patterns, fmt.Sprintf("%d[0-%d]", i*10+endTens, endOnes))
					}
					// 如果最后一个十位到9结束
					if endOnes == 9 {
						patterns = append(patterns, fmt.Sprintf("%d[0-9]", i*10+endTens))
					}
				}
			}
		}

		return "(" + Join(patterns, "|") + ")"
	}

	// 如果跨越了数量级，则使用范围表示
	return fmt.Sprintf("(%0*d-%0*d)", digits, start, digits, end)
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
