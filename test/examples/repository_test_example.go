package examples

import (
	"context"
	"testing"

	"Qingyu_backend/models/users"
	"Qingyu_backend/test/fixtures"
	"Qingyu_backend/test/testutil"

	"github.com/stretchr/testify/assert"
)

/*
本文件提供Repository层测试的完整示例
展示数据库操作测试、数据准备和清理等最佳实践

注意: 这是一个示例文件，实际测试需要连接真实的测试数据库
*/

// ============ 测试数据准备 ============

// setupTestDB 准备测试数据库（示例）
func setupTestDB(t *testing.T) (*ExampleUserRepository, func()) {
	// 在实际测试中，这里应该:
	// 1. 连接到测试数据库
	// 2. 创建必要的集合和索引
	// 3. 准备测试数据

	repo := &ExampleUserRepository{
		// 初始化repository
	}

	// 返回repository和清理函数
	cleanup := func() {
		// 清理测试数据
		// 关闭数据库连接
	}

	return repo, cleanup
}

// ============ 示例Repository ============

// ExampleUserRepository 示例用户仓储
type ExampleUserRepository struct {
	// 在实际实现中，这里会有数据库连接
	// collection *mongo.Collection
}

func (r *ExampleUserRepository) Create(ctx context.Context, user *users.User) error {
	// 实际实现会插入到MongoDB
	return nil
}

func (r *ExampleUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	// 实际实现会从MongoDB查询
	return nil, nil
}

func (r *ExampleUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	// 实际实现会从MongoDB查询
	return nil, nil
}

func (r *ExampleUserRepository) List(ctx context.Context, limit, skip int) ([]*users.User, error) {
	// 实际实现会从MongoDB查询
	return nil, nil
}

func (r *ExampleUserRepository) Delete(ctx context.Context, id string) error {
	// 实际实现会从MongoDB删除
	return nil
}

// ============ Table-Driven Repository测试示例 ============

func TestUserRepository_Create_TableDriven(t *testing.T) {
	// 使用工厂创建测试数据
	userFactory := fixtures.NewUserFactory()

	tests := []struct {
		name    string
		user    *users.User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "成功创建用户",
			user:    userFactory.Create(),
			wantErr: false,
		},
		{
			name:    "创建管理员用户",
			user:    userFactory.CreateAdmin(),
			wantErr: false,
		},
		{
			name:    "创建作者用户",
			user:    userFactory.CreateAuthor(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo, cleanup := setupTestDB(t)
			defer cleanup()
			ctx := testutil.CreateTestContext()

			// Act
			err := repo.Create(ctx, tt.user)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				// 验证数据是否真的插入了
				// 在实际测试中，应该查询数据库确认
			}
		})
	}
}

// ============ CRUD操作完整测试示例 ============

func TestUserRepository_CRUD_Complete(t *testing.T) {
	// Arrange
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := testutil.CreateTestContext()
	userFactory := fixtures.NewUserFactory()

	t.Run("创建用户", func(t *testing.T) {
		user := userFactory.Create()
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
	})

	t.Run("根据ID查询用户", func(t *testing.T) {
		// 先创建
		user := userFactory.Create()
		err := repo.Create(ctx, user)
		assert.NoError(t, err)

		// 再查询
		found, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		testutil.AssertUserEqual(t, user, found)
	})

	t.Run("根据邮箱查询用户", func(t *testing.T) {
		// 先创建
		user := userFactory.Create()
		err := repo.Create(ctx, user)
		assert.NoError(t, err)

		// 再查询
		found, err := repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, user.Email, found.Email)
	})

	t.Run("列出用户", func(t *testing.T) {
		// 创建多个用户
		users := userFactory.CreateBatch(5)
		for _, user := range users {
			err := repo.Create(ctx, user)
			assert.NoError(t, err)
		}

		// 查询列表
		result, err := repo.List(ctx, 10, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("删除用户", func(t *testing.T) {
		// 先创建
		user := userFactory.Create()
		err := repo.Create(ctx, user)
		assert.NoError(t, err)

		// 删除
		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err)

		// 验证已删除
		found, err := repo.GetByID(ctx, user.ID)
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

// ============ 数据准备和清理示例 ============

func TestUserRepository_WithDataPreparation(t *testing.T) {
	// Arrange: 准备测试数据
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := testutil.CreateTestContext()
	userFactory := fixtures.NewUserFactory()

	// 准备测试数据
	testUsers := []*users.User{
		userFactory.CreateAdmin(),
		userFactory.CreateAuthor(),
		userFactory.Create(),
	}

	// 插入测试数据
	for _, user := range testUsers {
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
	}

	// 注册清理函数
	testutil.RegisterCleanup(t, func() {
		for _, user := range testUsers {
			_ = repo.Delete(ctx, user.ID)
		}
	})

	// Act & Assert: 执行测试
	t.Run("验证管理员用户存在", func(t *testing.T) {
		admin := testUsers[0]
		found, err := repo.GetByID(ctx, admin.ID)
		assert.NoError(t, err)
		assert.Equal(t, "admin", found.Role)
	})

	t.Run("验证作者用户存在", func(t *testing.T) {
		author := testUsers[1]
		found, err := repo.GetByID(ctx, author.ID)
		assert.NoError(t, err)
		assert.Equal(t, "author", found.Role)
	})
}

// ============ 批量操作测试示例 ============

func TestUserRepository_BatchOperations(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := testutil.CreateTestContext()
	userFactory := fixtures.NewUserFactory()

	t.Run("批量创建用户", func(t *testing.T) {
		users := userFactory.CreateBatch(10)

		for _, user := range users {
			err := repo.Create(ctx, user)
			assert.NoError(t, err)
		}

		// 验证所有用户都创建成功
		result, err := repo.List(ctx, 20, 0)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 10)
	})

	t.Run("批量删除用户", func(t *testing.T) {
		// 创建测试用户
		users := userFactory.CreateBatch(5)
		for _, user := range users {
			err := repo.Create(ctx, user)
			assert.NoError(t, err)
		}

		// 批量删除
		for _, user := range users {
			err := repo.Delete(ctx, user.ID)
			assert.NoError(t, err)
		}

		// 验证都已删除
		for _, user := range users {
			found, err := repo.GetByID(ctx, user.ID)
			assert.Error(t, err)
			assert.Nil(t, found)
		}
	})
}

// ============ 边界条件测试示例 ============

func TestUserRepository_EdgeCases(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := testutil.CreateTestContext()

	t.Run("查询不存在的用户", func(t *testing.T) {
		user, err := repo.GetByID(ctx, "nonexistent_id")
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("删除不存在的用户", func(t *testing.T) {
		err := repo.Delete(ctx, "nonexistent_id")
		// 根据实际实现，可能返回错误或nil
		// assert.Error(t, err) 或 assert.NoError(t, err)
		_ = err // 忽略错误（示例）
	})

	t.Run("空列表查询", func(t *testing.T) {
		// 在清空的数据库中查询
		users, err := repo.List(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("分页边界测试", func(t *testing.T) {
		// 测试大偏移量
		users, err := repo.List(ctx, 10, 1000000)
		assert.NoError(t, err)
		assert.Empty(t, users)
	})
}

// ============ 并发测试示例 ============

func TestUserRepository_Concurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过并发测试")
	}

	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := testutil.CreateTestContext()
	userFactory := fixtures.NewUserFactory()

	t.Run("并发创建用户", func(t *testing.T) {
		const goroutines = 10
		usersToCreate := userFactory.CreateBatch(goroutines)

		// 并发创建
		errCh := make(chan error, goroutines)
		for _, user := range usersToCreate {
			go func(u *users.User) {
				errCh <- repo.Create(ctx, u)
			}(user)
		}

		// 检查结果
		for i := 0; i < goroutines; i++ {
			err := <-errCh
			assert.NoError(t, err)
		}
	})
}

// ============ 性能基准测试示例 ============

func BenchmarkUserRepository_Create(b *testing.B) {
	repo, cleanup := setupTestDB(&testing.T{})
	defer cleanup()

	ctx := testutil.CreateTestContext()
	userFactory := fixtures.NewUserFactory()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := userFactory.Create()
		_ = repo.Create(ctx, user)
	}
}

func BenchmarkUserRepository_GetByID(b *testing.B) {
	repo, cleanup := setupTestDB(&testing.T{})
	defer cleanup()

	ctx := testutil.CreateTestContext()
	userFactory := fixtures.NewUserFactory()

	// 准备测试数据
	user := userFactory.Create()
	_ = repo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetByID(ctx, user.ID)
	}
}

/*
运行本示例测试:
  go test -v ./test/examples/repository_test_example.go

运行基准测试:
  go test -bench=. ./test/examples/repository_test_example.go

跳过并发测试:
  go test -short ./test/examples/repository_test_example.go

实际使用建议:
1. 使用Docker启动测试数据库
2. 在每个测试前清理数据
3. 使用事务或defer清理测试数据
4. 避免测试之间的数据污染
5. 使用testutil助手函数简化测试代码
*/
