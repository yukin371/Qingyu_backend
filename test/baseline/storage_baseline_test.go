package baseline

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockStorageService Mock存储服务（支持状态管理）
type MockStorageService struct {
	UploadFunc   func(ctx context.Context, filename string, data []byte) (string, error)
	DownloadFunc func(ctx context.Context, fileID string) ([]byte, error)
	DeleteFunc   func(ctx context.Context, fileID string) error

	// 内部状态管理
	files   map[string][]byte
	deleted map[string]bool
}

func NewMockStorageService() *MockStorageService {
	return &MockStorageService{
		files:   make(map[string][]byte),
		deleted: make(map[string]bool),
	}
}

func (m *MockStorageService) Upload(ctx context.Context, filename string, data []byte) (string, error) {
	if m.UploadFunc != nil {
		return m.UploadFunc(ctx, filename, data)
	}
	// 默认实现：生成fileID并存储数据
	fileID := "file-" + filename
	m.files[fileID] = data
	return fileID, nil
}

func (m *MockStorageService) Download(ctx context.Context, fileID string) ([]byte, error) {
	if m.DownloadFunc != nil {
		return m.DownloadFunc(ctx, fileID)
	}
	// 默认实现：检查文件是否存在且未被删除
	if m.deleted[fileID] {
		return nil, io.EOF
	}
	data, exists := m.files[fileID]
	if !exists {
		return nil, io.EOF
	}
	return data, nil
}

func (m *MockStorageService) Delete(ctx context.Context, fileID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, fileID)
	}
	// 默认实现：标记文件为已删除
	m.deleted[fileID] = true
	delete(m.files, fileID)
	return nil
}

// StoragePort 存储服务端口接口（简化版，用于测试）
type StoragePort interface {
	Upload(ctx context.Context, filename string, data []byte) (string, error)
	Download(ctx context.Context, fileID string) ([]byte, error)
	Delete(ctx context.Context, fileID string) error
}

// TestStoragePortInterface 测试存储服务端口接口
func TestStoragePortInterface(t *testing.T) {
	t.Run("接口应该包含核心存储方法", func(t *testing.T) {
		mock := NewMockStorageService()
		var port StoragePort = mock

		// 测试Upload - 会生成file-test.txt
		fileID, err := port.Upload(context.Background(), "test.txt", []byte("data"))
		assert.NoError(t, err)
		assert.Equal(t, "file-test.txt", fileID)

		// 测试Download - 应该返回上传的数据
		data, err := port.Download(context.Background(), "file-test.txt")
		assert.NoError(t, err)
		assert.Equal(t, []byte("data"), data)

		// 测试Delete
		err = port.Delete(context.Background(), "file-test.txt")
		assert.NoError(t, err)
	})

	t.Run("应该支持自定义行为", func(t *testing.T) {
		uploadCalled := false
		mock := &MockStorageService{
			UploadFunc: func(ctx context.Context, filename string, data []byte) (string, error) {
				uploadCalled = true
				assert.Equal(t, "custom.txt", filename)
				return "custom-file", nil
			},
		}

		var port StoragePort = mock
		fileID, _ := port.Upload(context.Background(), "custom.txt", []byte("data"))

		assert.True(t, uploadCalled)
		assert.Equal(t, "custom-file", fileID)
	})

	t.Run("应该能够处理错误", func(t *testing.T) {
		mock := &MockStorageService{
			DownloadFunc: func(ctx context.Context, fileID string) ([]byte, error) {
				return nil, io.EOF
			},
		}

		var port StoragePort = mock
		_, err := port.Download(context.Background(), "nonexistent")

		assert.Error(t, err)
	})
}

// TestFileOperations 测试文件操作基线
func TestFileOperations(t *testing.T) {
	t.Run("完整的文件生命周期", func(t *testing.T) {
		mock := NewMockStorageService()
		var port StoragePort = mock
		ctx := context.Background()

		// 1. 上传文件
		testData := []byte("test file content")
		fileID, err := port.Upload(ctx, "test.txt", testData)
		require.NoError(t, err)
		assert.NotEmpty(t, fileID)

		// 2. 下载文件
		downloadedData, err := port.Download(ctx, fileID)
		require.NoError(t, err)
		assert.Equal(t, testData, downloadedData)

		// 3. 删除文件
		err = port.Delete(ctx, fileID)
		assert.NoError(t, err)

		// 4. 验证文件已删除（下载应该失败）
		_, err = port.Download(ctx, fileID)
		assert.Error(t, err)
	})

	t.Run("大文件处理", func(t *testing.T) {
		mock := NewMockStorageService()
		var port StoragePort = mock
		ctx := context.Background()

		// 创建1MB测试数据
		largeData := make([]byte, 1024*1024)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		fileID, err := port.Upload(ctx, "large.bin", largeData)
		assert.NoError(t, err)
		assert.NotEmpty(t, fileID)

		// 验证可以下载
		downloaded, err := port.Download(ctx, fileID)
		assert.NoError(t, err)
		assert.Len(t, downloaded, len(largeData))
	})
}

// BenchmarkStorageOperations 性能基线测试
func BenchmarkStorageOperations(b *testing.B) {
	mock := &MockStorageService{
		UploadFunc:   func(ctx context.Context, filename string, data []byte) (string, error) { return "file-123", nil },
		DownloadFunc: func(ctx context.Context, fileID string) ([]byte, error) { return []byte("data"), nil },
		DeleteFunc:   func(ctx context.Context, fileID string) error { return nil },
	}

	var port StoragePort = mock
	ctx := context.Background()

	smallData := []byte("small test data")
	largeData := make([]byte, 1024*100) // 100KB

	b.Run("UploadSmallFile", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = port.Upload(ctx, "small.txt", smallData)
		}
	})

	b.Run("UploadLargeFile", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = port.Upload(ctx, "large.bin", largeData)
		}
	})

	b.Run("DownloadFile", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = port.Download(ctx, "file-123")
		}
	})
}

// TestStorageErrorHandling 测试错误处理基线
func TestStorageErrorHandling(t *testing.T) {
	t.Run("应该处理上传失败", func(t *testing.T) {
		expectedErr := io.EOF
		mock := &MockStorageService{
			UploadFunc: func(ctx context.Context, filename string, data []byte) (string, error) {
				return "", expectedErr
			},
		}

		var port StoragePort = mock
		_, err := port.Upload(context.Background(), "test.txt", []byte("data"))

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("应该处理下载失败", func(t *testing.T) {
		expectedErr := io.ErrUnexpectedEOF
		mock := &MockStorageService{
			DownloadFunc: func(ctx context.Context, fileID string) ([]byte, error) {
				return nil, expectedErr
			},
		}

		var port StoragePort = mock
		_, err := port.Download(context.Background(), "file-123")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("应该处理删除失败", func(t *testing.T) {
		expectedErr := assert.AnError
		mock := &MockStorageService{
			DeleteFunc: func(ctx context.Context, fileID string) error {
				return expectedErr
			},
		}

		var port StoragePort = mock
		err := port.Delete(context.Background(), "file-123")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestStorageConcurrency 并发测试基线
func TestStorageConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过并发测试（short模式）")
	}

	t.Run("并发上传应该安全", func(t *testing.T) {
		mock := &MockStorageService{
			UploadFunc: func(ctx context.Context, filename string, data []byte) (string, error) {
				// 模拟延迟
				time.Sleep(10 * time.Millisecond)
				return "file-" + filename, nil
			},
		}

		var port StoragePort = mock
		ctx := context.Background()

		// 并发上传10个文件
		const concurrency = 10
		results := make(chan string, concurrency)
		errors := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(index int) {
				filename := fmt.Sprintf("file%d.txt", index)
				fileID, err := port.Upload(ctx, filename, []byte("data"))
				if err != nil {
					errors <- err
				} else {
					results <- fileID
				}
			}(i)
		}

		// 收集结果
		successCount := 0
		for i := 0; i < concurrency; i++ {
			select {
			case <-results:
				successCount++
			case err := <-errors:
				t.Logf("上传失败: %v", err)
			}
		}

		assert.Equal(t, concurrency, successCount, "所有上传都应该成功")
	})
}
