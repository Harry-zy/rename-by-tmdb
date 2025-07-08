package models

import "encoding/json"

// WordGroup 表示词组信息
type WordGroup struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	WordGroupType int    `json:"wordGroupType"`
}

// WordUnit 表示替换规则
type WordUnit struct {
	ID          int    `json:"id"`
	WordGroupID int    `json:"wordGroupId"`
	BeReplaced  string `json:"beReplaced"`
	Replace     string `json:"replace"`
	Front       string `json:"front"`
	Back        string `json:"back"`
	Offset      string `json:"offset"`
	Enabled     bool   `json:"enabled"`
	Type        int    `json:"type"`
	Regex       bool   `json:"regex"`
	Note        string `json:"note"`
}

// APIResponse 表示API通用响应格式
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// WordGroupResponse 表示词组API响应
type WordGroupResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    WordGroup `json:"data"`
}
