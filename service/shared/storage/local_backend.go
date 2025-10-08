package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalBackend 本地文件系统存储后端
type LocalBackend struct {
	basePath string // 存储根目录
	baseURL  string // 访问基础URL
}

// NewLocalBackend 创建本地存储后端
func NewLocalBackend(basePath, baseURL string) *LocalBackend {
	return &LocalBackend{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// Save 保存文件到本地
func (b *LocalBackend) Save(ctx context.Context, path string, reader io.Reader) error {
	// 1. 构建完整路径
	fullPath := filepath.Join(b.basePath, path)

	// 2. 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 3. 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 4. 写入内容
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// Load 从本地加载文件
func (b *LocalBackend) Load(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(b.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("文件不存在")
		}
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}

	return file, nil
}

// Delete 删除本地文件
func (b *LocalBackend) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(b.basePath, path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件已不存在，视为成功
		}
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (b *LocalBackend) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(b.basePath, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetURL 生成访问URL
func (b *LocalBackend) GetURL(ctx context.Context, path string, expiresIn time.Duration) (string, error) {
	// 本地存储简单拼接URL
	// 实际生产环境可能需要生成带签名的临时URL
	url := fmt.Sprintf("%s/%s", b.baseURL, path)
	return url, nil
}
