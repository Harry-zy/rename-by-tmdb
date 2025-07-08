package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/harry/rename-by-tmdb/internal/models"
)

// MediaType 表示媒体类型
type MediaType string

const (
	// MovieType 电影类型
	MovieType MediaType = "movie"
	// TVType 剧集类型
	TVType MediaType = "tv"
)

// TMDBService 处理TMDB API相关的操作
type TMDBService struct {
	apiKey string
}

// NewTMDBService 创建新的TMDB服务实例
func NewTMDBService() (*TMDBService, error) {
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEY 环境变量为空")
	}
	return &TMDBService{apiKey: apiKey}, nil
}

// checkTMDBResponse 检查TMDB API响应
func (s *TMDBService) checkTMDBResponse(resp *http.Response) error {
	// 如果HTTP状态码是200，直接返回成功
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP状态码: %d，读取响应失败: %v", resp.StatusCode, err)
	}

	// 重新设置响应体，因为ReadAll会消耗它
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	// 尝试解析错误响应
	var tmdbErr models.TMDBError
	if err := json.Unmarshal(body, &tmdbErr); err != nil || tmdbErr.StatusMessage == "" {
		// 如果解析失败或没有错误消息，返回原始响应内容
		return fmt.Errorf("HTTP状态码: %d，响应内容: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("TMDB API错误: [%d] %s", tmdbErr.StatusCode, tmdbErr.StatusMessage)
}

// FetchMovieInfo 获取电影信息
func (s *TMDBService) FetchMovieInfo(movieID string) (*models.TMDBMovie, error) {
	url := fmt.Sprintf("https://api.tmdb.org/3/movie/%s?language=zh-CN", movieID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建TMDB请求失败: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送TMDB请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if err := s.checkTMDBResponse(resp); err != nil {
		return nil, err
	}

	var movie models.TMDBMovie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("解析TMDB响应失败: %v", err)
	}

	return &movie, nil
}

// FetchShowInfo 获取剧集信息
func (s *TMDBService) FetchShowInfo(seriesID string) (*models.TMDBShow, error) {
	url := fmt.Sprintf("https://api.tmdb.org/3/tv/%s?language=zh-CN", seriesID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建TMDB请求失败: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送TMDB请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if err := s.checkTMDBResponse(resp); err != nil {
		return nil, err
	}

	var show models.TMDBShow
	if err := json.NewDecoder(resp.Body).Decode(&show); err != nil {
		return nil, fmt.Errorf("解析TMDB响应失败: %v", err)
	}

	return &show, nil
}

// FetchSeasonDetails 获取季度详细信息
func (s *TMDBService) FetchSeasonDetails(seriesID string, seasonNumber int) (*models.TMDBSeason, error) {
	url := fmt.Sprintf("https://api.tmdb.org/3/tv/%s/season/%d?language=zh-CN", seriesID, seasonNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建TMDB请求失败: %v", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送TMDB请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if err := s.checkTMDBResponse(resp); err != nil {
		return nil, err
	}

	var season models.TMDBSeason
	if err := json.NewDecoder(resp.Body).Decode(&season); err != nil {
		return nil, fmt.Errorf("解析TMDB响应失败: %v", err)
	}

	return &season, nil
}
