package integration

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"Qingyu_backend/service/shared/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorageIntegration 文件存储集成测试
func TestStorageIntegration(t *testing.T) {
	// 创建临时测试目录
	tempDir := "./test_uploads"
	defer os.RemoveAll(tempDir)

	// 创建LocalBackend
	backend := storage.NewLocalBackend(tempDir, "http://localhost:8080/files")
	ctx := context.Background()

	t.Run("Save_Load_Delete_流程测试", func(t *testing.T) {
		// 准备测试数据
		testContent := []byte("Hello, Storage Service!")
		testPath := "test/example.txt"

		// 1. 保存文件
		err := backend.Save(ctx, testPath, bytes.NewReader(testContent))
		require.NoError(t, err, "保存文件应该成功")

		// 验证文件是否存在
		fullPath := filepath.Join(tempDir, testPath)
		_, err = os.Stat(fullPath)
		assert.NoError(t, err, "文件应该存在于文件系统中")

		// 2. 检查文件是否存在
		exists, err := backend.Exists(ctx, testPath)
		require.NoError(t, err)
		assert.True(t, exists, "文件应该存在")

		// 3. 加载文件
		reader, err := backend.Load(ctx, testPath)
		require.NoError(t, err, "加载文件应该成功")

		// 验证文件内容
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		require.NoError(t, err)
		assert.Equal(t, string(testContent), buf.String(), "文件内容应该一致")

		// 立即关闭reader（Windows需要）
		reader.Close()

		// 4. 获取URL
		url, err := backend.GetURL(ctx, testPath, 1*time.Hour)
		require.NoError(t, err)
		assert.Contains(t, url, testPath, "URL应该包含文件路径")

		// 5. 删除文件
		err = backend.Delete(ctx, testPath)
		require.NoError(t, err, "删除文件应该成功")

		// 验证文件已删除
		exists, err = backend.Exists(ctx, testPath)
		require.NoError(t, err)
		assert.False(t, exists, "文件应该已被删除")
	})

	t.Run("文件不存在时的行为", func(t *testing.T) {
		// 加载不存在的文件
		_, err := backend.Load(ctx, "nonexistent/file.txt")
		assert.Error(t, err, "加载不存在的文件应该返回错误")

		// 检查不存在的文件
		exists, err := backend.Exists(ctx, "nonexistent/file.txt")
		require.NoError(t, err)
		assert.False(t, exists, "不存在的文件应该返回false")

		// 删除不存在的文件（应该成功，幂等性）
		err = backend.Delete(ctx, "nonexistent/file.txt")
		assert.NoError(t, err, "删除不存在的文件应该成功（幂等）")
	})

	t.Run("目录自动创建", func(t *testing.T) {
		// 保存到深层嵌套目录
		deepPath := "a/b/c/d/e/file.txt"
		content := []byte("Deep file")

		err := backend.Save(ctx, deepPath, bytes.NewReader(content))
		require.NoError(t, err, "应该自动创建所有父目录")

		// 验证文件存在
		exists, err := backend.Exists(ctx, deepPath)
		require.NoError(t, err)
		assert.True(t, exists)

		// 清理
		backend.Delete(ctx, deepPath)
	})

	t.Run("并发操作测试", func(t *testing.T) {
		const concurrency = 10
		done := make(chan bool, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer func() { done <- true }()

				path := filepath.Join("concurrent", "file"+string(rune('0'+index))+".txt")
				content := []byte("Concurrent test " + string(rune('0'+index)))

				// 保存
				err := backend.Save(ctx, path, bytes.NewReader(content))
				assert.NoError(t, err)

				// 加载
				reader, err := backend.Load(ctx, path)
				if err == nil {
					reader.Close()
				}

				// 删除
				backend.Delete(ctx, path)
			}(i)
		}

		// 等待所有goroutine完成
		for i := 0; i < concurrency; i++ {
			<-done
		}
	})
}

// TestStorageBackendPerformance 性能基准测试
func TestStorageBackendPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}

	tempDir := "./test_perf"
	defer os.RemoveAll(tempDir)

	backend := storage.NewLocalBackend(tempDir, "http://localhost:8080/files")
	ctx := context.Background()

	// 准备不同大小的测试数据
	testSizes := []struct {
		name string
		size int
	}{
		{"1KB", 1024},
		{"10KB", 10 * 1024},
		{"100KB", 100 * 1024},
		{"1MB", 1024 * 1024},
	}

	for _, ts := range testSizes {
		t.Run(ts.name, func(t *testing.T) {
			content := make([]byte, ts.size)
			path := "perf/" + ts.name + ".bin"

			// 测试保存性能
			start := time.Now()
			err := backend.Save(ctx, path, bytes.NewReader(content))
			saveTime := time.Since(start)

			require.NoError(t, err)
			t.Logf("保存 %s: %v", ts.name, saveTime)

			// 测试加载性能
			start = time.Now()
			reader, err := backend.Load(ctx, path)
			require.NoError(t, err)

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(reader)
			reader.Close()
			loadTime := time.Since(start)

			require.NoError(t, err)
			t.Logf("加载 %s: %v", ts.name, loadTime)

			// 清理
			backend.Delete(ctx, path)

			// 性能断言（本地文件系统应该很快）
			assert.Less(t, saveTime.Milliseconds(), int64(1000), "保存应该在1秒内完成")
			assert.Less(t, loadTime.Milliseconds(), int64(1000), "加载应该在1秒内完成")
		})
	}
}
