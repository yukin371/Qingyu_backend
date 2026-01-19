package migration

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MigrationRecord 迁移记录
type MigrationRecord struct {
	Version      string     `bson:"version" json:"version"`
	Description  string     `bson:"description" json:"description"`
	AppliedAt    time.Time  `bson:"applied_at" json:"appliedAt"`
	RolledBack   bool       `bson:"rolled_back" json:"rolledBack"`
	RolledBackAt *time.Time `bson:"rolled_back_at,omitempty" json:"rolledBackAt,omitempty"`
}

// Migration 迁移接口
type Migration interface {
	Version() string
	Description() string
	Up(ctx context.Context, db *mongo.Database) error
	Down(ctx context.Context, db *mongo.Database) error
}

// Manager 迁移管理器
type Manager struct {
	db         *mongo.Database
	collection *mongo.Collection
	migrations []Migration
}

// NewManager 创建迁移管理器
func NewManager(db *mongo.Database) *Manager {
	return &Manager{
		db:         db,
		collection: db.Collection("migrations"),
		migrations: make([]Migration, 0),
	}
}

// Register 注册迁移
func (m *Manager) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// RegisterMultiple 批量注册迁移
func (m *Manager) RegisterMultiple(migrations ...Migration) {
	m.migrations = append(m.migrations, migrations...)
}

// Up 执行迁移（升级）
func (m *Manager) Up(ctx context.Context) error {
	// 获取已应用的迁移
	appliedVersions, err := m.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// 排序迁移
	m.sortMigrations()

	// 执行未应用的迁移
	for _, migration := range m.migrations {
		version := migration.Version()

		// 检查是否已应用
		if appliedVersions[version] {
			fmt.Printf("Migration %s already applied, skipping\n", version)
			continue
		}

		fmt.Printf("Applying migration %s: %s\n", version, migration.Description())

		// 执行迁移
		if err := migration.Up(ctx, m.db); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", version, err)
		}

		// 记录迁移
		record := MigrationRecord{
			Version:     version,
			Description: migration.Description(),
			AppliedAt:   time.Now(),
			RolledBack:  false,
		}

		if _, err := m.collection.InsertOne(ctx, record); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		fmt.Printf("✓ Migration %s applied successfully\n", version)
	}

	return nil
}

// Down 回滚迁移
func (m *Manager) Down(ctx context.Context, steps int) error {
	// 获取已应用的迁移（按时间倒序）
	appliedRecords, err := m.getAppliedRecords(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedRecords) == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// 确定要回滚的数量
	rollbackCount := steps
	if steps <= 0 || steps > len(appliedRecords) {
		rollbackCount = len(appliedRecords)
	}

	// 创建迁移映射
	migrationMap := make(map[string]Migration)
	for _, migration := range m.migrations {
		migrationMap[migration.Version()] = migration
	}

	// 回滚迁移
	for i := 0; i < rollbackCount; i++ {
		record := appliedRecords[i]
		migration, exists := migrationMap[record.Version]

		if !exists {
			fmt.Printf("Warning: Migration %s not found, skipping\n", record.Version)
			continue
		}

		fmt.Printf("Rolling back migration %s: %s\n", record.Version, migration.Description())

		// 执行回滚
		if err := migration.Down(ctx, m.db); err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", record.Version, err)
		}

		// 更新记录
		now := time.Now()
		update := bson.M{
			"$set": bson.M{
				"rolled_back":    true,
				"rolled_back_at": now,
			},
		}

		filter := bson.M{"version": record.Version}
		if _, err := m.collection.UpdateOne(ctx, filter, update); err != nil {
			return fmt.Errorf("failed to update migration record %s: %w", record.Version, err)
		}

		fmt.Printf("✓ Migration %s rolled back successfully\n", record.Version)
	}

	return nil
}

// Status 获取迁移状态
func (m *Manager) Status(ctx context.Context) error {
	// 获取已应用的迁移
	appliedVersions, err := m.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// 排序迁移
	m.sortMigrations()

	fmt.Println("\n=== Migration Status ===")
	fmt.Printf("%-20s %-10s %-50s\n", "VERSION", "STATUS", "DESCRIPTION")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, migration := range m.migrations {
		version := migration.Version()
		status := "Pending"
		if appliedVersions[version] {
			status = "Applied"
		}

		fmt.Printf("%-20s %-10s %-50s\n", version, status, migration.Description())
	}

	fmt.Printf("\nTotal: %d migrations, %d applied, %d pending\n",
		len(m.migrations), len(appliedVersions), len(m.migrations)-len(appliedVersions))

	return nil
}

// getAppliedVersions 获取已应用的迁移版本
func (m *Manager) getAppliedVersions(ctx context.Context) (map[string]bool, error) {
	filter := bson.M{"rolled_back": false}
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	versions := make(map[string]bool)
	for cursor.Next(ctx) {
		var record MigrationRecord
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		versions[record.Version] = true
	}

	return versions, cursor.Err()
}

// getAppliedRecords 获取已应用的迁移记录（按时间倒序）
func (m *Manager) getAppliedRecords(ctx context.Context) ([]MigrationRecord, error) {
	filter := bson.M{"rolled_back": false}
	opts := options.Find().SetSort(bson.D{{Key: "applied_at", Value: -1}})

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []MigrationRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// sortMigrations 排序迁移（按版本号）
func (m *Manager) sortMigrations() {
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version() < m.migrations[j].Version()
	})
}

// Reset 重置所有迁移（谨慎使用！）
func (m *Manager) Reset(ctx context.Context) error {
	// 回滚所有迁移
	if err := m.Down(ctx, 0); err != nil {
		return err
	}

	// 删除迁移记录
	if _, err := m.collection.DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("failed to delete migration records: %w", err)
	}

	fmt.Println("✓ All migrations reset successfully")
	return nil
}
