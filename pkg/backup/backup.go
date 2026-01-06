package backup

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BackupConfig 备份配置
type BackupConfig struct {
	// 备份目录
	BackupDir string `json:"backup_dir" yaml:"backup_dir"`

	// 备份保留天数
	RetentionDays int `json:"retention_days" yaml:"retention_days"`

	// Cron表达式（例如："0 2 * * *" 表示每天凌晨2点）
	CronSchedule string `json:"cron_schedule" yaml:"cron_schedule"`

	// 是否启用压缩
	Compress bool `json:"compress" yaml:"compress"`

	// 备份失败通知回调
	NotificationCallback func(error) `json:"-" yaml:"-"`
}

// BackupManager 备份管理器
type BackupManager struct {
	config     *BackupConfig
	dbClient   *mongo.Client
	database   string
	cron       *cron.Cron
	logger     BackupLogger
}

// BackupLogger 备份日志接口
type BackupLogger interface {
	Info(msg string)
	Error(msg string, err error)
	Debug(msg string)
}

// DefaultBackupLogger 默认日志实现
type DefaultBackupLogger struct{}

func (l *DefaultBackupLogger) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (l *DefaultBackupLogger) Error(msg string, err error) {
	fmt.Printf("[ERROR] %s: %v\n", msg, err)
}

func (l *DefaultBackupLogger) Debug(msg string) {
	fmt.Printf("[DEBUG] %s\n", msg)
}

// NewBackupManager 创建备份管理器
func NewBackupManager(config *BackupConfig, dbClient *mongo.Client, database string) *BackupManager {
	if config.BackupDir == "" {
		config.BackupDir = "./backups"
	}

	if config.RetentionDays == 0 {
		config.RetentionDays = 7 // 默认保留7天
	}

	if config.CronSchedule == "" {
		config.CronSchedule = "0 2 * * *" // 默认每天凌晨2点
	}

	return &BackupManager{
		config:   config,
		dbClient: dbClient,
		database: database,
		cron:     cron.New(),
		logger:   &DefaultBackupLogger{},
	}
}

// SetLogger 设置日志记录器
func (m *BackupManager) SetLogger(logger BackupLogger) {
	m.logger = logger
}

// Start 启动定时备份任务
func (m *BackupManager) Start() error {
	// 确保备份目录存在
	if err := os.MkdirAll(m.config.BackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// 添加定时任务
	_, err := m.cron.AddFunc(m.config.CronSchedule, func() {
		if err := m.PerformBackup(context.Background()); err != nil {
			m.logger.Error("Backup failed", err)
			if m.config.NotificationCallback != nil {
				m.config.NotificationCallback(err)
			}
		} else {
			m.logger.Info("Backup completed successfully")
		}
	})

	if err != nil {
		return fmt.Errorf("failed to schedule backup: %w", err)
	}

	// 启动cron调度器
	m.cron.Start()
	m.logger.Info(fmt.Sprintf("Backup scheduler started: %s", m.config.CronSchedule))

	// 清理过期备份
	go m.CleanupOldBackups()

	return nil
}

// Stop 停止定时备份任务
func (m *BackupManager) Stop() {
	m.cron.Stop()
	m.logger.Info("Backup scheduler stopped")
}

// PerformBackup 执行备份
func (m *BackupManager) PerformBackup(ctx context.Context) error {
	startTime := time.Now()
	timestamp := startTime.Format("20060102_150405")
	backupFileName := fmt.Sprintf("%s_backup_%s", m.database, timestamp)
	backupPath := filepath.Join(m.config.BackupDir, backupFileName)

	m.logger.Info(fmt.Sprintf("Starting backup: %s", backupFileName))

	// 确保备份目录存在
	if err := os.MkdirAll(m.config.BackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// 1. MongoDB数据备份
	if err := m.backupMongoDB(ctx, backupPath); err != nil {
		return fmt.Errorf("failed to backup MongoDB: %w", err)
	}

	// 2. 压缩备份（如果启用）
	if m.config.Compress {
		if err := m.compressBackup(backupPath); err != nil {
			return fmt.Errorf("failed to compress backup: %w", err)
		}
	}

	duration := time.Since(startTime)
	m.logger.Info(fmt.Sprintf("Backup completed in %s", duration))

	return nil
}

// backupMongoDB 备份MongoDB数据
func (m *BackupManager) backupMongoDB(ctx context.Context, basePath string) error {
	// 如果dbClient为nil，跳过MongoDB备份
	if m.dbClient == nil {
		m.logger.Info("MongoDB client is nil, skipping database backup")
		return nil
	}

	// 创建备份目录
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return err
	}

	// 获取所有集合
	collections, err := m.dbClient.Database(m.database).ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	// 导出每个集合
	for _, collection := range collections {
		m.logger.Debug(fmt.Sprintf("Backing up collection: %s", collection))

		// 获取集合中的所有文档
		cursor, err := m.dbClient.Database(m.database).Collection(collection).Find(ctx, map[string]interface{}{})
		if err != nil {
			m.logger.Error(fmt.Sprintf("Failed to query collection %s", collection), err)
			continue
		}

		// 创建输出文件
		outputFile := filepath.Join(basePath, fmt.Sprintf("%s.json", collection))
		file, err := os.Create(outputFile)
		if err != nil {
			m.logger.Error(fmt.Sprintf("Failed to create output file for %s", collection), err)
			cursor.Close(ctx)
			continue
		}

		// 写入文档
		if err := m.dumpCollection(ctx, cursor, file); err != nil {
			file.Close()
			m.logger.Error(fmt.Sprintf("Failed to dump collection %s", collection), err)
			continue
		}

		file.Close()
		cursor.Close(ctx)
	}

	return nil
}

// dumpCollection 导出集合数据到文件
func (m *BackupManager) dumpCollection(ctx context.Context, cursor *mongo.Cursor, file *os.File) error {
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	defer cursor.Close(ctx)

	// 写入JSON数组开始
	file.WriteString("[\n")

	first := true
	for cursor.Next(ctx) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			m.logger.Error("Failed to decode document", err)
			continue
		}

		// 处理特殊类型
		processedDoc := m.processDocument(document)

		// 添加逗号分隔符
		if !first {
			file.WriteString(",\n")
		}
		first = false

		// 编码文档
		if err := encoder.Encode(processedDoc); err != nil {
			m.logger.Error("Failed to encode document", err)
			continue
		}
	}

	// 写入JSON数组结束
	file.WriteString("\n]")

	return cursor.Err()
}

// processDocument 处理文档中的特殊类型（如ObjectID、Date等）
func (m *BackupManager) processDocument(doc bson.M) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range doc {
		switch v := value.(type) {
		case primitive.ObjectID:
			result[key] = map[string]interface{}{
				"$oid": v.Hex(),
			}
		case primitive.DateTime:
			result[key] = map[string]interface{}{
				"$date": v.Time().Format(time.RFC3339),
			}
		case time.Time:
			result[key] = map[string]interface{}{
				"$date": v.Format(time.RFC3339),
			}
		case primitive.Timestamp:
			result[key] = map[string]interface{}{
				"$timestamp": map[string]interface{}{
					"t": v.T,
					"i": v.I,
				},
			}
		case primitive.Binary:
			result[key] = map[string]interface{}{
				"$binary": map[string]interface{}{
					"base64": v.Data,
					"subtype": v.Subtype,
				},
			}
		case primitive.A:
			// 处理数组中的嵌套文档
			arr := make([]interface{}, len(v))
			for i, item := range v {
				if subDoc, ok := item.(bson.M); ok {
					arr[i] = m.processDocument(subDoc)
				} else {
					arr[i] = item
				}
			}
			result[key] = arr
		case bson.M:
			// 递归处理嵌套文档
			result[key] = m.processDocument(v)
		case []interface{}:
			// 处理普通数组
			arr := make([]interface{}, len(v))
			for i, item := range v {
				if subDoc, ok := item.(bson.M); ok {
					arr[i] = m.processDocument(subDoc)
				} else {
					arr[i] = item
				}
			}
			result[key] = arr
		default:
			result[key] = v
		}
	}

	return result
}

// compressBackup 压缩备份文件
func (m *BackupManager) compressBackup(backupPath string) error {
	m.logger.Info(fmt.Sprintf("Compressing backup: %s", backupPath))

	// 创建gzip文件
	gzPath := backupPath + ".tar.gz"
	gzFile, err := os.Create(gzPath)
	if err != nil {
		return fmt.Errorf("failed to create gzip file: %w", err)
	}
	defer gzFile.Close()

	// 创建gzip写入器
	gzWriter := gzip.NewWriter(gzFile)
	defer gzWriter.Close()

	// 遍历备份目录中的所有文件
	err = filepath.Walk(backupPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if info.IsDir() && filePath == backupPath {
			return nil
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 读取文件
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		// 获取相对路径
		relPath, err := filepath.Rel(backupPath, filePath)
		if err != nil {
			return err
		}

		// 写入路径和分隔符
		if _, err := gzWriter.Write([]byte(fmt.Sprintf(">>>FILE:%s<<<", relPath))); err != nil {
			return fmt.Errorf("failed to write file header: %w", err)
		}

		// 写入文件内容
		if _, err := gzWriter.Write(data); err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}

		m.logger.Debug(fmt.Sprintf("Compressed file: %s", relPath))
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to compress backup: %w", err)
	}

	// 确保所有数据都写入
	if err := gzWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// 压缩成功后，删除原始备份目录
	if err := os.RemoveAll(backupPath); err != nil {
		m.logger.Error("Failed to remove original backup directory", err)
	}

	m.logger.Info(fmt.Sprintf("Backup compressed successfully: %s", gzPath))
	return nil
}

// CleanupOldBackups 清理过期备份
func (m *BackupManager) CleanupOldBackups() {
	ticker := time.NewTicker(24 * time.Hour) // 每天检查一次
	defer ticker.Stop()

	for range ticker.C {
		m.logger.Info("Checking for old backups to clean up...")

		if err := m.performCleanup(); err != nil {
			m.logger.Error("Failed to cleanup old backups", err)
		}
	}
}

// performCleanup 执行清理操作
func (m *BackupManager) performCleanup() error {
	entries, err := os.ReadDir(m.config.BackupDir)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -m.config.RetentionDays)
	cleanedCount := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			m.logger.Error(fmt.Sprintf("Failed to get info for %s", entry.Name()), err)
			continue
		}

		// 删除超过保留期的备份
		if info.ModTime().Before(cutoffTime) {
			backupPath := filepath.Join(m.config.BackupDir, entry.Name())
			if err := os.RemoveAll(backupPath); err != nil {
				m.logger.Error(fmt.Sprintf("Failed to remove %s", entry.Name()), err)
				continue
			}
			cleanedCount++
			m.logger.Info(fmt.Sprintf("Removed old backup: %s", entry.Name()))
		}
	}

	if cleanedCount > 0 {
		m.logger.Info(fmt.Sprintf("Cleaned up %d old backup(s)", cleanedCount))
	}

	return nil
}

// ListBackups 列出所有备份
func (m *BackupManager) ListBackups() ([]BackupInfo, error) {
	entries, err := os.ReadDir(m.config.BackupDir)
	if err != nil {
		return nil, err
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			Name:    entry.Name(),
			Size:    calculateSize(filepath.Join(m.config.BackupDir, entry.Name())),
			Created: info.ModTime(),
		})
	}

	return backups, nil
}

// BackupInfo 备份信息
type BackupInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Created time.Time `json:"created"`
}

// calculateSize 计算目录大小
func calculateSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// RestoreBackup 恢复备份
func (m *BackupManager) RestoreBackup(ctx context.Context, backupName string) error {
	backupPath := filepath.Join(m.config.BackupDir, backupName)

	// 检查备份是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupName)
	}

	m.logger.Info(fmt.Sprintf("Starting restore: %s", backupName))

	// 检查是否是压缩备份
	if filepath.Ext(backupPath) == ".gz" {
		return m.restoreCompressedBackup(ctx, backupPath)
	}

	// 恢复未压缩的备份
	return m.restoreUncompressedBackup(ctx, backupPath)
}

// restoreCompressedBackup 恢复压缩备份
func (m *BackupManager) restoreCompressedBackup(ctx context.Context, gzPath string) error {
	m.logger.Info(fmt.Sprintf("Decompressing backup: %s", gzPath))

	// 打开gzip文件
	gzFile, err := os.Open(gzPath)
	if err != nil {
		return fmt.Errorf("failed to open gzip file: %w", err)
	}
	defer gzFile.Close()

	// 创建gzip读取器
	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// 读取所有内容
	data, err := io.ReadAll(gzReader)
	if err != nil {
		return fmt.Errorf("failed to read gzip data: %w", err)
	}

	// 解析文件并恢复
	currentFile := ""
	var currentContent []byte
	content := string(data)

	for len(content) > 0 {
		// 查找文件标记
		startIdx := indexOf(content, ">>>FILE:")
		if startIdx == -1 {
			// 没有更多文件标记，结束
			break
		}

		// 如果有之前的内容，先写入
		if currentFile != "" {
			if err := m.writeRestoredFile(currentFile, currentContent); err != nil {
				m.logger.Error(fmt.Sprintf("Failed to write file: %s", currentFile), err)
			}
		}

		// 提取文件名
		endIdx := indexOf(content[startIdx:], "<<<")
		if endIdx == -1 {
			return fmt.Errorf("invalid backup format: missing file marker end")
		}
		endIdx += startIdx

		currentFile = content[startIdx+8 : endIdx] // 跳过 ">>>FILE:"
		content = content[endIdx+3:]              // 跳过 "<<<"

		// 查找下一个文件标记或结束
		nextFileIdx := indexOf(content, ">>>FILE:")
		if nextFileIdx == -1 {
			// 最后一个文件，读取剩余所有内容
			currentContent = []byte(content)
			content = ""
		} else {
			currentContent = []byte(content[:nextFileIdx])
			content = content[nextFileIdx:]
		}
	}

	// 写入最后一个文件
	if currentFile != "" {
		if err := m.writeRestoredFile(currentFile, currentContent); err != nil {
			return fmt.Errorf("failed to write final file: %w", err)
		}
	}

	// 解压后恢复数据
	basePath := filepath.Dir(gzPath)
	baseName := filepath.Base(gzPath)
	restoreDir := filepath.Join(basePath, strings.TrimSuffix(baseName, ".tar.gz")+"_restored")

	return m.restoreUncompressedBackup(ctx, restoreDir)
}

// restoreUncompressedBackup 恢复未压缩的备份
func (m *BackupManager) restoreUncompressedBackup(ctx context.Context, backupPath string) error {
	if m.dbClient == nil {
		return fmt.Errorf("database client is nil, cannot restore")
	}

	database := m.dbClient.Database(m.database)

	// 遍历备份目录中的所有JSON文件
	err := filepath.Walk(backupPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理JSON文件
		if filepath.Ext(filePath) != ".json" {
			return nil
		}

		// 获取集合名称（文件名去掉.json扩展名）
		collectionName := strings.TrimSuffix(filepath.Base(filePath), ".json")

		m.logger.Info(fmt.Sprintf("Restoring collection: %s", collectionName))

		// 读取JSON文件
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read backup file: %w", err)
		}

		// 解析JSON数组
		var documents []bson.M
		if err := json.Unmarshal(data, &documents); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}

		// 获取或创建集合
		collection := database.Collection(collectionName)

		// 清空现有数据（可选）
		// collection.DeleteMany(ctx, bson.M{})

		// 批量插入文档
		if len(documents) > 0 {
			// 将处理后的文档转换回bson.M
			var interfaceDocs []interface{}
			for _, doc := range documents {
				interfaceDocs = append(interfaceDocs, doc)
			}

			// 使用InsertMany插入文档
			_, err := collection.InsertMany(ctx, interfaceDocs)
			if err != nil {
				return fmt.Errorf("failed to insert documents: %w", err)
			}

			m.logger.Info(fmt.Sprintf("Restored %d documents to collection: %s", len(documents), collectionName))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	m.logger.Info("Backup restored successfully")
	return nil
}

// writeRestoredFile 写入恢复的文件
func (m *BackupManager) writeRestoredFile(relPath string, content []byte) error {
	// 创建完整的文件路径
	fullPath := filepath.Join(m.config.BackupDir, "restore_temp", relPath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(fullPath, content, 0644)
}

// indexOf 查找子字符串位置
func indexOf(s, substr string) int {
	idx := strings.Index(s, substr)
	return idx
}

// BackupStatus 备份状态
type BackupStatus struct {
	LastBackup    time.Time `json:"last_backup"`
	NextBackup    time.Time `json:"next_backup"`
	TotalBackups  int       `json:"total_backups"`
	TotalSize     int64     `json:"total_size"`
	Schedule      string    `json:"schedule"`
	IsRunning     bool      `json:"is_running"`
}

// GetStatus 获取备份状态
func (m *BackupManager) GetStatus() (*BackupStatus, error) {
	backups, err := m.ListBackups()
	if err != nil {
		return nil, err
	}

	var lastBackup time.Time
	if len(backups) > 0 {
		// 找到最新的备份
		for _, backup := range backups {
			if backup.Created.After(lastBackup) {
				lastBackup = backup.Created
			}
		}
	}

	status := &BackupStatus{
		LastBackup:   lastBackup,
		TotalBackups: len(backups),
		Schedule:     m.config.CronSchedule,
		// 注意：cron.Entries()在cron启动后可能仍为空，
		// 这是cron库的预期行为，我们不能依赖它来判断是否运行
		IsRunning: false, // 默认为false，可以通过外部状态管理
	}

	// 计算总大小
	for _, backup := range backups {
		status.TotalSize += backup.Size
	}

	return status, nil
}
