package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/harry/rename-by-tmdb/internal/models"
)

// WordGroupService 处理词组相关的操作
type WordGroupService struct {
	apiBaseURL string
	authToken  string
}

// NewWordGroupService 创建新的词组服务实例
func NewWordGroupService() (*WordGroupService, error) {
	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		return nil, fmt.Errorf("API_BASE_URL 未设置")
	}
	// 确保URL以/api/v1结尾
	if !strings.HasSuffix(apiBaseURL, "/api/v1") {
		apiBaseURL = strings.TrimRight(apiBaseURL, "/") + "/api/v1"
	}

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("AUTH_TOKEN 未设置")
	}

	return &WordGroupService{
		apiBaseURL: apiBaseURL,
		authToken:  authToken,
	}, nil
}

// CreateWordGroup 创建词组
func (s *WordGroupService) CreateWordGroup(title string) (*models.WordGroup, error) {
	url := fmt.Sprintf("%s/wordGroup/add", s.apiBaseURL)

	body := map[string]string{
		"title": title,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("JSON编码失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d，错误信息: %s", resp.StatusCode, apiResp.Message)
	}

	if apiResp.Code != 20000 {
		return nil, fmt.Errorf("API返回错误: [%d] %s", apiResp.Code, apiResp.Message)
	}

	var response models.WordGroupResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("解析响应数据失败: %v", err)
	}

	return &response.Data, nil
}

// AddWordUnit 添加替换规则
func (s *WordGroupService) AddWordUnit(groupID int, beReplaced, replace, front, back string, offset int) error {
	url := fmt.Sprintf("%s/wordUnit/add", s.apiBaseURL)

	// 设置偏移量字符串和规则类型
	var offsetStr string
	ruleType := 200 // 默认类型，用于无偏移的情况

	if offset != 0 {
		// 有偏移量时
		ruleType = 300
		if offset > 0 {
			offsetStr = fmt.Sprintf("EP+%d", offset)
		} else {
			offsetStr = fmt.Sprintf("EP%d", offset) // 负数已经包含负号
		}
	}

	wordUnit := models.WordUnit{
		ID:          0,
		WordGroupID: groupID,
		BeReplaced:  beReplaced,
		Replace:     replace,
		Front:       front,
		Back:        back,
		Offset:      offsetStr,
		Enabled:     true,
		Type:        ruleType,
		Regex:       true,
		Note:        "",
	}

	jsonData, err := json.Marshal(wordUnit)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API请求失败，状态码: %d，错误信息: %s", resp.StatusCode, apiResp.Message)
	}

	if apiResp.Code != 20000 {
		return fmt.Errorf("API返回错误: [%d] %s", apiResp.Code, apiResp.Message)
	}

	return nil
}

// WordGroupList 表示词组列表响应
type WordGroupList struct {
	Total    int                `json:"total"`
	PageNum  int                `json:"pageNum"`
	PageSize int                `json:"pageSize"`
	List     []models.WordGroup `json:"list"`
}

// GetWordGroupList 获取词组列表
func (s *WordGroupService) GetWordGroupList() (*WordGroupList, error) {
	url := fmt.Sprintf("%s/wordGroup/page?pageNum=1&pageSize=9999&keyword=", s.apiBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var apiResp struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    WordGroupList `json:"data"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d，错误信息: %s", resp.StatusCode, apiResp.Message)
	}

	if apiResp.Code != 20000 {
		return nil, fmt.Errorf("API返回错误: [%d] %s", apiResp.Code, apiResp.Message)
	}

	return &apiResp.Data, nil
}
