package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetUserInput 从用户获取输入
func GetUserInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("读取输入失败: %v", err)
	}
	return strings.TrimSpace(input), nil
}

// GetHasSeasonChoice 从用户获取是否包含季数的选择（直接回车默认为包含）
func GetHasSeasonChoice() (bool, error) {
	input, err := GetUserInput("是否使用原文件名季数？(y/n，直接回车默认为y): ")
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" || input == "y" || input == "yes" {
		return true, nil
	}
	return false, nil
}

// GetEpisodeOffset 从用户获取集数偏移量（直接回车默认为0）
func GetEpisodeOffset() (int, error) {
	input, err := GetUserInput("请输入集数偏移量（如：+1、-1，直接回车表示不偏移）: ")
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return 0, nil
	}

	// 尝试解析偏移量
	offset, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("无效的偏移量: %v", err)
	}

	return offset, nil
}

// GetSpecificSeasons 从用户获取指定的季数
func GetSpecificSeasons() ([]int, bool, error) {
	input, err := GetUserInput("请输入要生成的季数（多季用;分隔，直接回车生成所有季，0表示特别篇）: ")
	if err != nil {
		return nil, false, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return nil, true, nil // 返回 nil, true 表示生成所有季
	}

	// 分割输入的季数
	seasonStrs := strings.Split(input, ";")
	var seasons []int

	// 解析每个季数
	for _, s := range seasonStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		season, err := strconv.Atoi(s)
		if err != nil {
			return nil, false, fmt.Errorf("无效的季数 '%s': %v", s, err)
		}
		if season < 0 {
			return nil, false, fmt.Errorf("季数不能为负数: %d", season)
		}
		seasons = append(seasons, season)
	}

	if len(seasons) == 0 {
		return nil, true, nil // 如果没有有效的季数，则生成所有季
	}

	return seasons, false, nil
}

// GetIncludeSpecialSeason 从用户获取是否包含第0季（特别篇）的选择（直接回车默认为n）
func GetIncludeSpecialSeason() (bool, error) {
	input, err := GetUserInput("是否包含第0季（特别篇）？(y/n，直接回车默认为n): ")
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}

// GetPadZeroChoice 从用户获取是否需要补0站位的选择（直接回车默认为y）
func GetPadZeroChoice() (bool, error) {
	input, err := GetUserInput("集数是否补0站位？(y/n，直接回车默认为y): ")
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" || input == "y" || input == "yes" {
		return true, nil
	}
	return false, nil
}

// GetEpisodeContinuousChoice 从用户获取集数是否连续的选择（直接回车默认为y）
func GetEpisodeContinuousChoice() (bool, error) {
	input, err := GetUserInput("集数是否连续？(y/n，直接回车默认为y): ")
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" || input == "y" || input == "yes" {
		return true, nil
	}
	return false, nil
}

// GetPartEpisodeChoice 从用户获取是否有part剧集的选择（直接回车默认为n）
func GetPartEpisodeChoice() (bool, error) {
	input, err := GetUserInput("是否有part剧集（y/n，直接回车默认为n）: ")
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}

// GetPartEpisodeInfo 从用户获取part剧集信息
func GetPartEpisodeInfo() (map[int][]int, error) {
	input, err := GetUserInput("请输入有part的集数和part数（格式为：集数:part数，多集之间以;间隔）例如：2:2;5:2，代表第二集和第五集都有part1和part2: ")
	if err != nil {
		return nil, err
	}

	// 清理输入，移除所有不可见字符和BOM
	input = strings.TrimSpace(input)

	// 添加调试信息，显示输入的原始字符
	fmt.Printf("调试信息 - 输入长度: %d, 字符码点: ", len(input))
	for i, r := range input {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("'%c'(%d)", r, r)
	}
	fmt.Println()

	// 移除可能的BOM和其他不可见字符，只保留数字、冒号、分号和空格
	input = strings.Map(func(r rune) rune {
		// 只保留数字、冒号、分号、空格和换行符
		if (r >= '0' && r <= '9') || r == ':' || r == ';' || r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return r
		}
		return -1 // 移除其他所有字符
	}, input)

	// 再次清理空格
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, fmt.Errorf("输入不能为空")
	}

	// 分割多集信息
	episodeStrs := strings.Split(input, ";")
	partInfo := make(map[int][]int)

	// 解析每个集数的part信息
	for _, episodeStr := range episodeStrs {
		episodeStr = strings.TrimSpace(episodeStr)
		if episodeStr == "" {
			continue
		}

		// 分割集数和part数
		parts := strings.Split(episodeStr, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("无效的格式 '%s'，应为 '集数:part数'", episodeStr)
		}

		// 添加调试信息
		fmt.Printf("调试信息 - 分割结果: parts[0]='%s'(长度:%d), parts[1]='%s'(长度:%d)\n",
			parts[0], len(parts[0]), parts[1], len(parts[1]))

		// 解析集数，添加调试信息
		episodeStrClean := strings.TrimSpace(parts[0])
		episodeNum, err := strconv.Atoi(episodeStrClean)
		if err != nil {
			return nil, fmt.Errorf("无效的集数 '%s' (长度:%d): %v", episodeStrClean, len(episodeStrClean), err)
		}
		if episodeNum <= 0 {
			return nil, fmt.Errorf("集数必须大于0: %d", episodeNum)
		}

		// 解析part数，添加调试信息
		partStrClean := strings.TrimSpace(parts[1])
		partCount, err := strconv.Atoi(partStrClean)
		if err != nil {
			return nil, fmt.Errorf("无效的part数 '%s' (长度:%d): %v", partStrClean, len(partStrClean), err)
		}
		if partCount <= 0 {
			return nil, fmt.Errorf("part数必须大于0: %d", partCount)
		}

		// 生成part列表（从1开始）
		var partList []int
		for i := 1; i <= partCount; i++ {
			partList = append(partList, i)
		}

		partInfo[episodeNum] = partList
	}

	if len(partInfo) == 0 {
		return nil, fmt.Errorf("没有有效的part剧集信息")
	}

	return partInfo, nil
}
