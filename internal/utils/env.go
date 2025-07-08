package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv 加载环境变量
func LoadEnv() error {
	// 尝试多个可能的位置
	envPaths := []string{
		".env",                   // 当前工作目录
		"../rename-by-tmdb/.env", // 相对于工作目录的项目目录
	}

	// 获取可执行文件所在目录
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		// 添加可执行文件目录下的 .env
		envPaths = append(envPaths, filepath.Join(execDir, ".env"))
		// 添加可执行文件上级目录下的 rename-by-tmdb/.env
		envPaths = append(envPaths, filepath.Join(execDir, "..", "rename-by-tmdb", ".env"))
	}

	// 获取源代码文件所在的目录
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		sourceDir := filepath.Dir(filename)
		// 添加源代码目录下的 .env
		envPaths = append(envPaths, filepath.Join(sourceDir, ".env"))
	}

	// 尝试每个可能的路径
	var lastErr error
	for _, path := range envPaths {
		err := godotenv.Load(path)
		if err == nil {
			fmt.Printf("成功加载配置文件: %s\n", path)
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("未能找到或加载 .env 文件，尝试过以下路径：\n%s\n最后一个错误: %v",
		strings.Join(envPaths, "\n"), lastErr)
}

// CheckRequiredEnvVars 检查必需的环境变量
func CheckRequiredEnvVars() error {
	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	if tmdbAPIKey == "" {
		return fmt.Errorf("环境变量 TMDB_API_KEY 未设置")
	}

	// 检查上传相关的环境变量
	if IsUploadEnabled() {
		apiBaseURL := os.Getenv("API_BASE_URL")
		if apiBaseURL == "" {
			return fmt.Errorf("启用上传功能时，环境变量 API_BASE_URL 必须设置")
		}

		authToken := os.Getenv("AUTH_TOKEN")
		if authToken == "" {
			return fmt.Errorf("启用上传功能时，环境变量 AUTH_TOKEN 必须设置")
		}
	}
	return nil
}

// IsUploadEnabled 检查是否启用上传
func IsUploadEnabled() bool {
	upload := strings.ToLower(os.Getenv("UPLOAD_MS"))
	return upload == "true"
}
