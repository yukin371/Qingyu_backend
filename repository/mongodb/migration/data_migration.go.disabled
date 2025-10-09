package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"Qingyu_backend/models/users"
	repository "Qingyu_backend/repository/mongodb"
)

// MigrationStatus 迁移状态
type MigrationStatus string

const (
	MigrationStatusPending   MigrationStatus = "pending"
	MigrationStatusRunning   MigrationStatus = "running"
	MigrationStatusCompleted MigrationStatus = "completed"
	MigrationStatusFailed    MigrationStatus = "failed"
	MigrationStatusRollback  MigrationStatus = "rollback"
)

// MigrationDirection 迁移方向
type MigrationDirection string

const (
	MigrationDirectionUp   MigrationDirection = "up"
	MigrationDirectionDown MigrationDirection = "down"
)

// MigrationRecord 迁移记录
type MigrationRecord struct {
	ID             string             `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	Direction      MigrationDirection `json:"direction" bson:"direction"`
	Status         MigrationStatus    `json:"status" bson:"status"`
	StartTime      time.Time          `json:"start_time" bson:"start_time"`
	EndTime        *time.Time         `json:"end_time,omitempty" bson:"end_time,omitempty"`
	Error          string             `json:"error,omitempty" bson:"error,omitempty"`
	Progress       int                `json:"progress" bson:"progress"` // 0-100
	TotalItems     int64              `json:"total_items" bson:"total_items"`
	ProcessedItems int64              `json:"processed_items" bson:"processed_items"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

// Migration 迁移接口
type Migration interface {
	// GetName 获取迁移名称
	GetName() string

	// Up 执行向上迁移
	Up(ctx context.Context, source, target repository.RepositoryFactory) error

	// Down 执行向下迁移（回滚）
	Down(ctx context.Context, source, target repository.RepositoryFactory) error

	// Validate 验证迁移前置条件
	Validate(ctx context.Context, source, target repository.RepositoryFactory) error

	// EstimateTime 估算迁移时间
	EstimateTime(ctx context.Context, source repository.RepositoryFactory) (time.Duration, error)

	// GetProgress 获取迁移进度
	GetProgress() (processed, total int64)
}

// UserMigration 用户数据迁移
type UserMigration struct {
	batchSize    int
	processed    int64
	total        int64
	mu           sync.RWMutex
	progressChan chan MigrationProgress
}

// MigrationProgress 迁移进度
type MigrationProgress struct {
	Processed int64  `json:"processed"`
	Total     int64  `json:"total"`
	Progress  int    `json:"progress"`
	Message   string `json:"message"`
}

// NewUserMigration 创建用户迁移实例
func NewUserMigration(batchSize int) *UserMigration {
	return &UserMigration{
		batchSize:    batchSize,
		progressChan: make(chan MigrationProgress, 100),
	}
}

// GetName 获取迁移名称
func (m *UserMigration) GetName() string {
	return "user_migration"
}

// Up 执行向上迁移（从源数据库迁移到目标数据库）
func (m *UserMigration) Up(ctx context.Context, source, target repository.RepositoryFactory) error {
	sourceRepo := source.CreateUserRepository()
	targetRepo := target.CreateUserRepository()

	// 统计总数
	total, err := sourceRepo.Count(ctx, repository.UserFilter{})
	if err != nil {
		return fmt.Errorf("统计源数据库用户数量失败: %w", err)
	}

	m.mu.Lock()
	m.total = total
	m.processed = 0
	m.mu.Unlock()

	// 分批迁移
	var offset int64 = 0
	for offset < total {
		// 获取一批用户
		users, err := sourceRepo.List(ctx, repository.UserFilter{
			Limit:  int64(m.batchSize),
			Offset: offset,
		})
		if err != nil {
			return fmt.Errorf("获取用户数据失败 (offset: %d): %w", offset, err)
		}

		// 迁移这批用户
		for _, user := range users {
			// 检查目标数据库中是否已存在
			exists, err := targetRepo.Exists(ctx, user.ID)
			if err != nil {
				return fmt.Errorf("检查用户存在性失败 (ID: %s): %w", user.ID, err)
			}

			if !exists {
				// 创建用户
				if err := targetRepo.Create(ctx, user); err != nil {
					// 如果是重复错误，跳过
					if !repository.IsDuplicateError(err) {
						return fmt.Errorf("创建用户失败 (ID: %s): %w", user.ID, err)
					}
				}
			}

			// 更新进度
			m.mu.Lock()
			m.processed++
			progress := int((m.processed * 100) / m.total)
			m.mu.Unlock()

			// 发送进度更新
			select {
			case m.progressChan <- MigrationProgress{
				Processed: m.processed,
				Total:     m.total,
				Progress:  progress,
				Message:   fmt.Sprintf("已迁移用户: %s", user.Username),
			}:
			default:
				// 如果通道满了，跳过这次进度更新
			}
		}

		offset += int64(len(users))

		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

// Down 执行向下迁移（回滚）
func (m *UserMigration) Down(ctx context.Context, source, target repository.RepositoryFactory) error {
	// 回滚逻辑：从目标数据库删除迁移的数据
	targetRepo := target.CreateUserRepository()

	// 获取所有用户
	users, err := targetRepo.List(ctx, repository.UserFilter{})
	if err != nil {
		return fmt.Errorf("获取目标数据库用户失败: %w", err)
	}

	// 删除用户
	for _, user := range users {
		if err := targetRepo.HardDelete(ctx, user.ID); err != nil {
			log.Printf("删除用户失败 (ID: %s): %v", user.ID, err)
		}
	}

	return nil
}

// Validate 验证迁移前置条件
func (m *UserMigration) Validate(ctx context.Context, source, target repository.RepositoryFactory) error {
	// 检查源数据库连接
	if err := source.Health(ctx); err != nil {
		return fmt.Errorf("源数据库连接失败: %w", err)
	}

	// 检查目标数据库连接
	if err := target.Health(ctx); err != nil {
		return fmt.Errorf("目标数据库连接失败: %w", err)
	}

	// 检查源数据库是否有数据
	sourceRepo := source.CreateUserRepository()
	count, err := sourceRepo.Count(ctx, repository.UserFilter{})
	if err != nil {
		return fmt.Errorf("检查源数据库数据失败: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("源数据库没有用户数据")
	}

	return nil
}

// EstimateTime 估算迁移时间
func (m *UserMigration) EstimateTime(ctx context.Context, source repository.RepositoryFactory) (time.Duration, error) {
	sourceRepo := source.CreateUserRepository()
	count, err := sourceRepo.Count(ctx, repository.UserFilter{})
	if err != nil {
		return 0, fmt.Errorf("统计用户数量失败: %w", err)
	}

	// 估算每个用户迁移需要10毫秒
	estimatedSeconds := (count * 10) / 1000
	return time.Duration(estimatedSeconds) * time.Second, nil
}

// GetProgress 获取迁移进度
func (m *UserMigration) GetProgress() (processed, total int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.processed, m.total
}

// GetProgressChan 获取进度通道
func (m *UserMigration) GetProgressChan() <-chan MigrationProgress {
	return m.progressChan
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	migrations []Migration
	records    map[string]*MigrationRecord
	mu         sync.RWMutex
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager() *MigrationManager {
	return &MigrationManager{
		migrations: make([]Migration, 0),
		records:    make(map[string]*MigrationRecord),
	}
}

// RegisterMigration 注册迁移
func (m *MigrationManager) RegisterMigration(migration Migration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.migrations = append(m.migrations, migration)
}

// RunMigration 执行迁移
func (m *MigrationManager) RunMigration(ctx context.Context, name string, direction MigrationDirection, source, target repository.RepositoryFactory) error {
	// 查找迁移
	var migration Migration
	for _, mig := range m.migrations {
		if mig.GetName() == name {
			migration = mig
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("迁移 %s 不存在", name)
	}

	// 创建迁移记录
	record := &MigrationRecord{
		ID:        fmt.Sprintf("%s_%d", name, time.Now().Unix()),
		Name:      name,
		Direction: direction,
		Status:    MigrationStatusPending,
		StartTime: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.mu.Lock()
	m.records[record.ID] = record
	m.mu.Unlock()

	// 更新状态为运行中
	m.updateRecordStatus(record.ID, MigrationStatusRunning, "")

	// 验证前置条件
	if err := migration.Validate(ctx, source, target); err != nil {
		m.updateRecordStatus(record.ID, MigrationStatusFailed, err.Error())
		return fmt.Errorf("迁移验证失败: %w", err)
	}

	// 执行迁移
	var err error
	switch direction {
	case MigrationDirectionUp:
		err = migration.Up(ctx, source, target)
	case MigrationDirectionDown:
		err = migration.Down(ctx, source, target)
	default:
		err = fmt.Errorf("不支持的迁移方向: %s", direction)
	}

	// 更新迁移状态
	if err != nil {
		m.updateRecordStatus(record.ID, MigrationStatusFailed, err.Error())
		return fmt.Errorf("迁移执行失败: %w", err)
	}

	m.updateRecordStatus(record.ID, MigrationStatusCompleted, "")
	return nil
}

// updateRecordStatus 更新迁移记录状态
func (m *MigrationManager) updateRecordStatus(recordID string, status MigrationStatus, errorMsg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if record, exists := m.records[recordID]; exists {
		record.Status = status
		record.UpdatedAt = time.Now()
		if errorMsg != "" {
			record.Error = errorMsg
		}
		if status == MigrationStatusCompleted || status == MigrationStatusFailed {
			now := time.Now()
			record.EndTime = &now
		}
	}
}

// GetMigrationRecord 获取迁移记录
func (m *MigrationManager) GetMigrationRecord(recordID string) (*MigrationRecord, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	record, exists := m.records[recordID]
	return record, exists
}

// ListMigrationRecords 列出所有迁移记录
func (m *MigrationManager) ListMigrationRecords() []*MigrationRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	records := make([]*MigrationRecord, 0, len(m.records))
	for _, record := range m.records {
		records = append(records, record)
	}
	return records
}

// GetAvailableMigrations 获取可用的迁移列表
func (m *MigrationManager) GetAvailableMigrations() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.migrations))
	for _, migration := range m.migrations {
		names = append(names, migration.GetName())
	}
	return names
}

// EstimateMigrationTime 估算迁移时间
func (m *MigrationManager) EstimateMigrationTime(ctx context.Context, name string, source repository.RepositoryFactory) (time.Duration, error) {
	for _, migration := range m.migrations {
		if migration.GetName() == name {
			return migration.EstimateTime(ctx, source)
		}
	}
	return 0, fmt.Errorf("迁移 %s 不存在", name)
}

// DataSyncManager 数据同步管理器
// 用于混合架构中的数据同步
type DataSyncManager struct {
	source     repository.RepositoryFactory
	target     repository.RepositoryFactory
	syncConfig *SyncConfig
	mu         sync.RWMutex
}

// SyncConfig 同步配置
type SyncConfig struct {
	Interval    time.Duration `json:"interval"`            // 同步间隔
	BatchSize   int           `json:"batch_size"`          // 批处理大小
	Enabled     bool          `json:"enabled"`             // 是否启用同步
	Direction   string        `json:"direction"`           // 同步方向: "source_to_target", "target_to_source", "bidirectional"
	ConflictRes string        `json:"conflict_resolution"` // 冲突解决策略: "source_wins", "target_wins", "latest_wins"
}

// NewDataSyncManager 创建数据同步管理器
func NewDataSyncManager(source, target repository.RepositoryFactory, config *SyncConfig) *DataSyncManager {
	return &DataSyncManager{
		source:     source,
		target:     target,
		syncConfig: config,
	}
}

// StartSync 启动数据同步
func (s *DataSyncManager) StartSync(ctx context.Context) error {
	if !s.syncConfig.Enabled {
		return fmt.Errorf("数据同步未启用")
	}

	ticker := time.NewTicker(s.syncConfig.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.syncUsers(ctx); err != nil {
				log.Printf("用户数据同步失败: %v", err)
			}
		}
	}
}

// syncUsers 同步用户数据
func (s *DataSyncManager) syncUsers(ctx context.Context) error {
	sourceRepo := s.source.CreateUserRepository()
	targetRepo := s.target.CreateUserRepository()

	// 获取源数据库中的用户
	users, err := sourceRepo.List(ctx, repository.UserFilter{
		Limit: int64(s.syncConfig.BatchSize),
	})
	if err != nil {
		return fmt.Errorf("获取源用户数据失败: %w", err)
	}

	// 同步到目标数据库
	for _, user := range users {
		exists, err := targetRepo.Exists(ctx, user.ID)
		if err != nil {
			log.Printf("检查用户存在性失败 (ID: %s): %v", user.ID, err)
			continue
		}

		if !exists {
			// 用户不存在，创建
			if err := targetRepo.Create(ctx, user); err != nil {
				log.Printf("创建用户失败 (ID: %s): %v", user.ID, err)
			}
		} else {
			// 用户存在，根据冲突解决策略处理
			if err := s.resolveUserConflict(ctx, user, targetRepo); err != nil {
				log.Printf("解决用户冲突失败 (ID: %s): %v", user.ID, err)
			}
		}
	}

	return nil
}

// resolveUserConflict 解决用户冲突
func (s *DataSyncManager) resolveUserConflict(ctx context.Context, sourceUser *system.User, targetRepo repository.UserRepository) error {
	switch s.syncConfig.ConflictRes {
	case "source_wins":
		// 源数据库优先，更新目标数据库
		updates := map[string]interface{}{
			"username":   sourceUser.Username,
			"email":      sourceUser.Email,
			"updated_at": sourceUser.UpdatedAt,
		}
		return targetRepo.Update(ctx, sourceUser.ID, updates)

	case "target_wins":
		// 目标数据库优先，不做任何操作
		return nil

	case "latest_wins":
		// 最新数据优先
		targetUser, err := targetRepo.GetByID(ctx, sourceUser.ID)
		if err != nil {
			return err
		}

		if sourceUser.UpdatedAt.After(targetUser.UpdatedAt) {
			// 源数据更新，更新目标数据库
			updates := map[string]interface{}{
				"username":   sourceUser.Username,
				"email":      sourceUser.Email,
				"updated_at": sourceUser.UpdatedAt,
			}
			return targetRepo.Update(ctx, sourceUser.ID, updates)
		}
		return nil

	default:
		return fmt.Errorf("不支持的冲突解决策略: %s", s.syncConfig.ConflictRes)
	}
}
