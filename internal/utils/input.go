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
