package migration

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"go.mongodb.org/mongo-driver/mongo"
)

// SimpleMigration 简化的迁移接口（只包含Up和Down方法）
type SimpleMigration interface {
	Up(ctx context.Context, db *mongo.Database) error
	Down(ctx context.Context, db *mongo.Database) error
}

// Migrator 迁移执行器
type Migrator struct {
	db         *mongo.Database
	migrations map[string]SimpleMigration
}

// NewMigrator 创建迁移执行器
func NewMigrator(db *mongo.Database) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make(map[string]SimpleMigration),
	}
}

// Register 注册迁移
func (m *Migrator) Register(name string, migration SimpleMigration) error {
	// 验证名称格式：数字_小写字母数字下划线（例如: 001_create_users_indexes, 003_create_books_indexes_p0）
	matched, err := regexp.MatchString(`^\d{3}_[a-z0-9_]+$`, name)
	if err != nil {
		return fmt.Errorf("迁移名称验证失败: %w", err)
	}
	if !matched {
		return fmt.Errorf("无效的迁移名称格式: %s (期望格式: 001_description)", name)
	}

	m.migrations[name] = migration
	log.Printf("✅ 已注册迁移: %s", name)
	return nil
}

// Up 执行迁移
func (m *Migrator) Up(ctx context.Context, name string) error {
	migration, exists := m.migrations[name]
	if !exists {
		return fmt.Errorf("迁移不存在: %s", name)
	}

	log.Printf("🚀 开始执行迁移: %s", name)
	if err := migration.Up(ctx, m.db); err != nil {
		return fmt.Errorf("执行迁移失败 %s: %w", name, err)
	}

	log.Printf("✅ 迁移执行成功: %s", name)
	return nil
}

// Down 回滚迁移
func (m *Migrator) Down(ctx context.Context, name string) error {
	migration, exists := m.migrations[name]
	if !exists {
		return fmt.Errorf("迁移不存在: %s", name)
	}

	log.Printf("🔄 开始回滚迁移: %s", name)
	if err := migration.Down(ctx, m.db); err != nil {
		return fmt.Errorf("回滚迁移失败 %s: %w", name, err)
	}

	log.Printf("✅ 迁移回滚成功: %s", name)
	return nil
}
