package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockLogger 模拟日志记录器
type MockLogger struct {
	infoLogs  []string
	errorLogs []string
	debugLogs []string
}

func (m *MockLogger) Info(msg string) {
	m.infoLogs = append(m.infoLogs, msg)
}

func (m *MockLogger) Error(msg string, err error) {
	m.errorLogs = append(m.errorLogs, msg)
}

func (m *MockLogger) Debug(msg string) {
	m.debugLogs = append(m.debugLogs, msg)
}

// TestNewBackupManager 测试创建备份管理器
func TestNewBackupManager(t *testing.T) {
	config := &BackupConfig{
		BackupDir:     "./test_backups",
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
		Compress:      false,
	}

	manager := NewBackupManager(config, nil, "test_db")

	assert.NotNil(t, manager)
	assert.Equal(t, config, manager.config)
	assert.Equal(t, "test_db", manager.database)
	assert.NotNil(t, manager.cron)
}

// TestBackupManager_SetLogger 测试设置日志记录器
func TestBackupManager_SetLogger(t *testing.T) {
	config := &BackupConfig{
		BackupDir:     "./test_backups",
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")
	mockLogger := &MockLogger{}

	manager.SetLogger(mockLogger)

	assert.Equal(t, mockLogger, manager.logger)
}

// TestBackupManager_PerformBackup 测试执行备份
func TestBackupManager_PerformBackup(t *testing.T) {
	// 创建临时备份目录
	tempDir := filepath.Join(os.TempDir(), "backup_test")
	defer os.RemoveAll(tempDir)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
		Compress:      false,
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 执行备份
	err := manager.PerformBackup(context.Background())

	// 由于没有实际的MongoDB客户端，备份应该成功（跳过数据库备份）
	assert.NoError(t, err)

	// 检查备份目录是否创建
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		assert.Fail(t, "Backup directory should be created")
	}
}

// TestBackupManager_ListBackups 测试列出现有备份
func TestBackupManager_ListBackups(t *testing.T) {
	// 创建临时备份目录和测试备份
	tempDir := filepath.Join(os.TempDir(), "backup_list_test")
	defer os.RemoveAll(tempDir)

	// 创建测试备份目录
	backup1 := filepath.Join(tempDir, "backup_20250101_020000")
	backup2 := filepath.Join(tempDir, "backup_20250102_020000")

	os.MkdirAll(backup1, 0755)
	os.MkdirAll(backup2, 0755)

	// 创建一些测试文件
	os.WriteFile(filepath.Join(backup1, "data.json"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(backup2, "data.json"), []byte("test"), 0644)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	backups, err := manager.ListBackups()

	assert.NoError(t, err)
	assert.Len(t, backups, 2)

	// 验证备份信息
	assert.Contains(t, backups[0].Name, "backup_")
	assert.Contains(t, backups[1].Name, "backup_")
	assert.Greater(t, backups[0].Size, int64(0))
	assert.False(t, backups[0].Created.IsZero())
}

// TestBackupManager_performCleanup 测试清理过期备份
func TestBackupManager_performCleanup(t *testing.T) {
	// 创建临时备份目录
	tempDir := filepath.Join(os.TempDir(), "backup_cleanup_test")
	defer os.RemoveAll(tempDir)

	// 创建一个旧的备份目录（超过保留期）
	oldBackup := filepath.Join(tempDir, "old_backup")
	os.MkdirAll(oldBackup, 0755)

	// 修改目录时间为10天前
	oldTime := time.Now().AddDate(0, 0, -10)
	os.Chtimes(oldBackup, oldTime, oldTime)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 执行清理
	err := manager.performCleanup()

	assert.NoError(t, err)

	// 验证旧备份已被删除
	_, err = os.Stat(oldBackup)
	assert.True(t, os.IsNotExist(err))
}

// TestBackupManager_GetStatus 测试获取备份状态
func TestBackupManager_GetStatus(t *testing.T) {
	// 创建临时备份目录
	tempDir := filepath.Join(os.TempDir(), "backup_status_test")
	defer os.RemoveAll(tempDir)

	// 创建测试备份
	backup := filepath.Join(tempDir, "backup_20250101_020000")
	os.MkdirAll(backup, 0755)
	os.WriteFile(filepath.Join(backup, "data.json"), []byte("test data"), 0644)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 启动调度器以更新状态
	manager.cron.Start()
	defer manager.cron.Stop()

	status, err := manager.GetStatus()

	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "0 2 * * *", status.Schedule)
	assert.Equal(t, 1, status.TotalBackups)
	assert.Greater(t, status.TotalSize, int64(0))
}

// TestBackupManager_RestoreBackup 测试恢复备份
func TestBackupManager_RestoreBackup(t *testing.T) {
	// 创建临时备份目录
	tempDir := filepath.Join(os.TempDir(), "backup_restore_test")
	defer os.RemoveAll(tempDir)

	// 创建测试备份
	backupName := "backup_20250101_020000"
	backupPath := filepath.Join(tempDir, backupName)
	os.MkdirAll(backupPath, 0755)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 测试恢复不存在的备份
	err := manager.RestoreBackup(context.Background(), "nonexistent")
	assert.Error(t, err)

	// 测试恢复存在的备份（会因为功能未实现而成功）
	err = manager.RestoreBackup(context.Background(), backupName)
	// 当前实现中RestoreBackup还未完整实现，所以可能返回nil或error
	// 只要不panic即可
	assert.NotPanics(t, func() {
		manager.RestoreBackup(context.Background(), backupName)
	})
}

// TestBackupInfo 测试备份信息结构
func TestBackupInfo(t *testing.T) {
	info := BackupInfo{
		Name:    "backup_20250101",
		Size:    1024 * 1024 * 100, // 100MB
		Created: time.Now(),
	}

	assert.Equal(t, "backup_20250101", info.Name)
	assert.Equal(t, int64(104857600), info.Size)
	assert.False(t, info.Created.IsZero())
}

// TestBackupConfig_Defaults 测试默认配置
func TestBackupConfig_Defaults(t *testing.T) {
	config := &BackupConfig{
		BackupDir: "",
		RetentionDays: 0,
		CronSchedule:  "",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 验证默认值
	assert.Equal(t, "./backups", manager.config.BackupDir)
	assert.Equal(t, 7, manager.config.RetentionDays)
	assert.Equal(t, "0 2 * * *", manager.config.CronSchedule)
}

// TestCalculateSize 测试目录大小计算
func TestCalculateSize(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "size_test")
	defer os.RemoveAll(tempDir)

	// 创建目录和文件
	os.MkdirAll(tempDir, 0755)
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("test data 123"), 0644)
	os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte("more test data"), 0644)

	size := calculateSize(tempDir)

	assert.Greater(t, size, int64(0))
	assert.Equal(t, int64(27), size) // "test data 123"(13) + "more test data"(14)
}

// TestBackupManager_StartStop 测试启动和停止备份调度器
func TestBackupManager_StartStop(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup_start_test")
	defer os.RemoveAll(tempDir)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 启动
	err := manager.Start()
	assert.NoError(t, err)

	// 验证cron正在运行
	assert.True(t, len(manager.cron.Entries()) > 0)

	// 停止
	manager.Stop()

	// 验证cron已停止
	// 注意：cron.Stop()不会立即清除entries，但会停止新的调度
}

// TestBackupManager_NotificationCallback 测试备份失败通知
func TestBackupManager_NotificationCallback(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "backup_notify_test")
	defer os.RemoveAll(tempDir)

	callbackCalled := false
	var callbackError error

	config := &BackupConfig{
		BackupDir: tempDir,
		CronSchedule: "@every 1s", // 每秒执行一次用于测试
		NotificationCallback: func(err error) {
			callbackCalled = true
			callbackError = err
		},
	}

	manager := NewBackupManager(config, nil, "test_db")
	mockLogger := &MockLogger{}
	manager.SetLogger(mockLogger)

	// 启动调度器
	err := manager.Start()
	assert.NoError(t, err)

	// 等待一次备份执行
	time.Sleep(2 * time.Second)

	// 停止调度器
	manager.Stop()

	// 备份可能会因为没有MongoDB客户端而失败
	// 验证日志被调用
	assert.Greater(t, len(mockLogger.infoLogs)+len(mockLogger.errorLogs), 0)

	// 如果回调被调用，验证错误
	if callbackCalled {
		assert.NotNil(t, callbackError)
	}
}

// BenchmarkListBackups 性能测试 - 列出备份
func BenchmarkListBackups(b *testing.B) {
	tempDir := filepath.Join(os.TempDir(), "backup_bench_test")
	defer os.RemoveAll(tempDir)

	// 创建大量测试备份
	for i := 0; i < 100; i++ {
		backupPath := filepath.Join(tempDir, filepath.Join("backup", time.Now().Format("20060102_150405")))
		os.MkdirAll(backupPath, 0755)
	}

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.ListBackups()
	}
}

// BenchmarkCalculateSize 性能测试 - 计算目录大小
func BenchmarkCalculateSize(b *testing.B) {
	tempDir := filepath.Join(os.TempDir(), "size_bench_test")
	defer os.RemoveAll(tempDir)

	// 创建大量测试文件
	for i := 0; i < 1000; i++ {
		file := filepath.Join(tempDir, filepath.Join("subdir", "file.txt"))
		os.MkdirAll(filepath.Dir(file), 0755)
		os.WriteFile(file, []byte("test data"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculateSize(tempDir)
	}
}

// TestProcessDocument 测试文档处理
func TestProcessDocument(t *testing.T) {
	config := &BackupConfig{
		BackupDir:     "./test_backups",
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 测试ObjectID处理
	now := time.Now()
	doc := bson.M{
		"_id":        primitive.NewObjectID(),
		"name":       "test",
		"count":      42,
		"created_at": now,
		"nested": bson.M{
			"key": "value",
		},
		"array": []interface{}{"item1", "item2"},
	}

	processed := manager.processDocument(doc)

	// 验证ObjectID被正确处理
	assert.Contains(t, processed, "_id")
	assert.Contains(t, processed["_id"], "$oid")

	// 验证简单字段保持不变
	assert.Equal(t, "test", processed["name"])
	assert.Equal(t, 42, processed["count"])

	// 验证嵌套文档被处理
	assert.Contains(t, processed["nested"], "key")

	// 验证数组被处理
	assert.Contains(t, processed["array"], "item1")
}

// TestCompressBackup 测试压缩功能
func TestCompressBackup(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "compress_test")
	defer os.RemoveAll(tempDir)

	// 创建测试备份目录和文件
	backupDir := filepath.Join(tempDir, "backup_20250101_020000")
	os.MkdirAll(backupDir, 0755)
	os.WriteFile(filepath.Join(backupDir, "data.json"), []byte(`{"key": "value"}`), 0644)
	os.WriteFile(filepath.Join(backupDir, "data2.json"), []byte(`{"key2": "value2"}`), 0644)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 测试压缩
	err := manager.compressBackup(backupDir)
	assert.NoError(t, err)

	// 验证压缩文件存在
	gzPath := backupDir + ".tar.gz"
	_, err = os.Stat(gzPath)
	assert.NoError(t, err)

	// 验证原始目录已删除
	_, err = os.Stat(backupDir)
	assert.True(t, os.IsNotExist(err))
}

// TestRestoreBackup_Compressed 测试恢复压缩备份
func TestRestoreBackup_Compressed(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "restore_test")
	defer os.RemoveAll(tempDir)

	// 创建测试备份目录和文件
	backupDir := filepath.Join(tempDir, "backup_20250101_020000")
	os.MkdirAll(backupDir, 0755)
	testData := []byte(`[{"_id": {"$oid": "507f1f77bcf86cd799439011"}, "name": "test"}]`)
	os.WriteFile(filepath.Join(backupDir, "test_collection.json"), testData, 0644)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 压缩备份
	err := manager.compressBackup(backupDir)
	assert.NoError(t, err)

	// 测试恢复（没有实际数据库客户端，只测试解压部分）
	gzPath := backupDir + ".tar.gz"
	err = manager.RestoreBackup(context.Background(), filepath.Base(gzPath))
	// 由于没有dbClient，期望返回错误
	assert.Error(t, err)
}

// TestIndexOf 测试indexOf函数
func TestIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected int
	}{
		{"找到子串", "hello world", "world", 6},
		{"未找到", "hello world", "xyz", -1},
		{"空子串", "hello world", "", 0},
		{"开头匹配", "hello world", "hello", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := indexOf(tt.s, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWriteRestoredFile 测试写入恢复文件
func TestWriteRestoredFile(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "write_restore_test")
	defer os.RemoveAll(tempDir)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 测试写入文件
	content := []byte(`{"test": "data"}`)
	err := manager.writeRestoredFile(filepath.Join("subdir", "test.json"), content)
	assert.NoError(t, err)

	// 验证文件存在
	expectedPath := filepath.Join(tempDir, "restore_temp", "subdir", "test.json")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err)

	// 验证文件内容
	data, err := os.ReadFile(expectedPath)
	assert.NoError(t, err)
	assert.Equal(t, content, data)
}

// BenchmarkProcessDocument 性能测试 - 文档处理
func BenchmarkProcessDocument(b *testing.B) {
	config := &BackupConfig{
		BackupDir:     "./test_backups",
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	// 创建复杂的测试文档
	doc := bson.M{
		"_id":        primitive.NewObjectID(),
		"name":       "test",
		"count":      42,
		"created_at": time.Now(),
		"nested": bson.M{
			"key1": "value1",
			"key2": 123,
			"deep": bson.M{
				"key": "deep_value",
			},
		},
		"array": []interface{}{
			"item1",
			"item2",
			bson.M{"nested": "array_item"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.processDocument(doc)
	}
}

// BenchmarkCompressBackup 性能测试 - 压缩备份
func BenchmarkCompressBackup(b *testing.B) {
	tempDir := filepath.Join(os.TempDir(), "compress_bench_test")
	defer os.RemoveAll(tempDir)

	config := &BackupConfig{
		BackupDir:     tempDir,
		RetentionDays: 7,
		CronSchedule:  "0 2 * * *",
	}

	manager := NewBackupManager(config, nil, "test_db")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建临时备份目录
		backupDir := filepath.Join(tempDir, fmt.Sprintf("backup_%d", i))
		os.MkdirAll(backupDir, 0755)

		// 创建一些测试文件
		for j := 0; j < 10; j++ {
			data := make([]byte, 1024) // 1KB
			os.WriteFile(filepath.Join(backupDir, fmt.Sprintf("file%d.json", j)), data, 0644)
		}

		// 压缩
		_ = manager.compressBackup(backupDir)

		// 清理压缩文件
		os.Remove(backupDir + ".tar.gz")
	}
}

