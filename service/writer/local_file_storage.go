package writer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalFileStorage 本地文件存储实现
type LocalFileStorage struct {
	baseDir string
	baseURL string
}

// NewLocalFileStorage 创建本地文件存储
func NewLocalFileStorage(baseDir, baseURL string) *LocalFileStorage {
	_ = os.MkdirAll(baseDir, 0755)
	return &LocalFileStorage{baseDir: baseDir, baseURL: baseURL}
}

// Upload 上传文件
func (s *LocalFileStorage) Upload(ctx context.Context, filename string, content io.Reader, mimeType string) (string, error) {
	dateDir := time.Now().Format("2006/01/02")
	dir := filepath.Join(s.baseDir, dateDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create directory failed: %w", err)
	}

	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	filePath := filepath.Join(dir, uniqueName)

	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create file failed: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, content); err != nil {
		return "", fmt.Errorf("write file failed: %w", err)
	}

	relPath := filepath.Join(dateDir, uniqueName)
	return fmt.Sprintf("%s/exports/%s", s.baseURL, relPath), nil
}

// Download 下载文件
func (s *LocalFileStorage) Download(ctx context.Context, url string) (io.ReadCloser, error) {
	relPath := ""
	fmt.Sscanf(url, s.baseURL+"/exports/%s", &relPath)
	if relPath == "" {
		return nil, fmt.Errorf("invalid file URL")
	}

	f, err := os.Open(filepath.Join(s.baseDir, relPath))
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	return f, nil
}

// Delete 删除文件
func (s *LocalFileStorage) Delete(ctx context.Context, url string) error {
	relPath := ""
	fmt.Sscanf(url, s.baseURL+"/exports/%s", &relPath)
	if relPath == "" {
		return fmt.Errorf("invalid file URL")
	}
	return os.Remove(filepath.Join(s.baseDir, relPath))
}

// GetSignedURL 获取签名URL（本地存储直接返回原URL）
func (s *LocalFileStorage) GetSignedURL(ctx context.Context, url string, expiration time.Duration) (string, error) {
	return url, nil
}

// Verify interface
var _ FileStorage = (*LocalFileStorage)(nil)
