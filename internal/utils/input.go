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
