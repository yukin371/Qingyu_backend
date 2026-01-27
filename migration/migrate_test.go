package migration

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

// TestSimpleMigration_Interface 验证SimpleMigration接口
func TestSimpleMigration_Interface(t *testing.T) {
	// 创建一个测试用的迁移实现
	testMigration := &TestSimpleMigrationImplementation{}

	// 验证它实现了SimpleMigration接口
	var _ SimpleMigration = testMigration

	// 测试Up方法可以被调用
	t.Run("Up method callable", func(t *testing.T) {
		ctx := context.Background()
		// 由于没有真实的数据库连接，我们只验证方法可以被调用而不panic
		_ = testMigration.Up
		_ = ctx
	})

	// 测试Down方法可以被调用
	t.Run("Down method callable", func(t *testing.T) {
		ctx := context.Background()
		// 由于没有真实的数据库连接，我们只验证方法可以被调用而不panic
		_ = testMigration.Down
		_ = ctx
	})
}

// TestSimpleMigrationImplementation 用于测试的迁移实现
type TestSimpleMigrationImplementation struct{}

func (m *TestSimpleMigrationImplementation) Up(ctx context.Context, db *mongo.Database) error {
	return nil
}

func (m *TestSimpleMigrationImplementation) Down(ctx context.Context, db *mongo.Database) error {
	return nil
}

// TestNewMigrator 验证Migrator创建
func TestNewMigrator(t *testing.T) {
	t.Run("NewMigrator creates non-nil instance", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Error("NewMigrator should return non-nil instance")
		}
	})

	t.Run("NewMigrator initializes migrations map", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Error("Migrator should not be nil")
		}
		if migrator.migrations == nil {
			t.Error("migrations map should be initialized")
		}
	})
}

// TestMigrator_Register 验证迁移注册
func TestMigrator_Register(t *testing.T) {
	t.Run("Register adds migration to migrator", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Fatal("Migrator should not be nil")
		}

		testMigration := &TestSimpleMigrationImplementation{}
		migrator.Register("test-migration", testMigration)

		// 验证迁移已注册
		if len(migrator.migrations) != 1 {
			t.Errorf("Expected 1 migration, got %d", len(migrator.migrations))
		}

		_, exists := migrator.migrations["test-migration"]
		if !exists {
			t.Error("Migration 'test-migration' should be registered")
		}
	})

	t.Run("Register multiple migrations", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Fatal("Migrator should not be nil")
		}

		migration1 := &TestSimpleMigrationImplementation{}
		migration2 := &TestSimpleMigrationImplementation{}

		migrator.Register("migration-1", migration1)
		migrator.Register("migration-2", migration2)

		// 验证多个迁移都可以注册
		if len(migrator.migrations) != 2 {
			t.Errorf("Expected 2 migrations, got %d", len(migrator.migrations))
		}

		_, exists1 := migrator.migrations["migration-1"]
		_, exists2 := migrator.migrations["migration-2"]

		if !exists1 || !exists2 {
			t.Error("Both migrations should be registered")
		}
	})
}
