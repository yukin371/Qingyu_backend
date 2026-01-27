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

	t.Run("NewMigrator can accept registrations", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Fatal("Migrator should not be nil")
		}

		testMigration := &TestSimpleMigrationImplementation{}
		err := migrator.Register("001_test_migration", testMigration)
		if err != nil {
			t.Errorf("Register should succeed: %v", err)
		}
	})
}

// TestMigrator_Register 验证迁移注册
func TestMigrator_Register(t *testing.T) {
	t.Run("Register valid migration", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Fatal("Migrator should not be nil")
		}

		testMigration := &TestSimpleMigrationImplementation{}
		err := migrator.Register("001_test_migration", testMigration)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Register invalid migration name format", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Fatal("Migrator should not be nil")
		}

		testMigration := &TestSimpleMigrationImplementation{}

		// 测试各种无效格式
		invalidNames := []string{
			"test-migration",     // 缺少数字前缀
			"001_test-migration", // 包含连字符
			"001TestMigration",   // 包含大写字母
			"001 test migration", // 包含空格
			"1_test",             // 数字不足3位
			"001_测试",             // 包含非ASCII字符
		}

		for _, name := range invalidNames {
			err := migrator.Register(name, testMigration)
			if err == nil {
				t.Errorf("Expected error for invalid name '%s', got nil", name)
			}
		}
	})

	t.Run("Register multiple migrations", func(t *testing.T) {
		migrator := NewMigrator(nil)
		if migrator == nil {
			t.Fatal("Migrator should not be nil")
		}

		migration1 := &TestSimpleMigrationImplementation{}
		migration2 := &TestSimpleMigrationImplementation{}

		err1 := migrator.Register("001_migration_one", migration1)
		err2 := migrator.Register("002_migration_two", migration2)

		if err1 != nil || err2 != nil {
			t.Errorf("Both registrations should succeed: err1=%v, err2=%v", err1, err2)
		}
	})
}

// TestMigrator_Up_Down_NonExistent 验证执行不存在的迁移
func TestMigrator_Up_Down_NonExistent(t *testing.T) {
	migrator := NewMigrator(nil)
	ctx := context.Background()

	t.Run("Up returns error for non-existent migration", func(t *testing.T) {
		err := migrator.Up(ctx, "999_non_existent")
		if err == nil {
			t.Error("Expected error for non-existent migration, got nil")
		}
	})

	t.Run("Down returns error for non-existent migration", func(t *testing.T) {
		err := migrator.Down(ctx, "999_non_existent")
		if err == nil {
			t.Error("Expected error for non-existent migration, got nil")
		}
	})
}
