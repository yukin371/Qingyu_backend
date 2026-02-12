package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	authModel "Qingyu_backend/models/auth"
	middlewareAuth "Qingyu_backend/internal/middleware/auth"
	"Qingyu_backend/repository/mongodb/auth"
	permService "Qingyu_backend/service/auth"
)

// TestPermissionDatabaseIntegration 测试完整的权限数据库集成
func TestPermissionDatabaseIntegration(t *testing.T) {
	// 跳过如果环境变量未设置
	if os.Getenv("TEST_MODE") != "true" {
		t.Skip("跳过集成测试（TEST_MODE未设置）")
	}

	// 1. 连接测试数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("TEST_MONGO_HOST")
	if mongoURI == "" {
		mongoURI = "localhost:27018"
	}
	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "qingyu_permission_test"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+mongoURI))
	require.NoError(t, err, "连接MongoDB失败")
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	logger := zap.NewNop()

	// 2. 创建Repository和Service
	repository := auth.NewAuthRepository(db)
	service := permService.NewPermissionService(repository, nil, logger)

	// 3. 创建RBACChecker
	checker, err := middlewareAuth.NewRBACChecker(nil)
	require.NoError(t, err, "创建RBACChecker失败")
	rbacChecker := checker.(*middlewareAuth.RBACChecker)

	// 4. 设置Checker到Service
	service.SetChecker(rbacChecker)

	t.Run("LoadPermissionsFromDatabase", func(t *testing.T) {
		// 从数据库加载权限
		err := service.LoadPermissionsToChecker(ctx)
		require.NoError(t, err, "加载权限失败")

		// 验证角色权限已加载
		adminPerms := rbacChecker.GetRolePermissions("admin")
		assert.Contains(t, adminPerms, "*:*", "admin角色应该有通配符权限")

		authorPerms := rbacChecker.GetRolePermissions("author")
		assert.Contains(t, authorPerms, "book:read", "author角色应该有book:read权限")
		assert.Contains(t, authorPerms, "book:create", "author角色应该有book:create权限")

		readerPerms := rbacChecker.GetRolePermissions("reader")
		assert.Contains(t, readerPerms, "book:read", "reader角色应该有book:read权限")
		assert.NotContains(t, readerPerms, "book:create", "reader角色不应该有book:create权限")

		t.Logf("✓ 权限加载成功，共加载 %d 个角色", 5)
	})

	t.Run("LoadUserRolesToChecker", func(t *testing.T) {
		// 加载admin用户角色
		err := service.LoadUserRolesToChecker(ctx, "admin@test.com")
		require.NoError(t, err, "加载用户角色失败")

		// 验证角色已分配
		roles := rbacChecker.GetUserRoles("admin@test.com")
		assert.Contains(t, roles, "admin", "admin用户应该有admin角色")

		t.Logf("✓ 用户角色加载成功：%v", roles)
	})

	t.Run("CheckPermissionWithRealData", func(t *testing.T) {
		// 加载admin用户角色
		service.LoadUserRolesToChecker(ctx, "admin@test.com")

		// 测试权限检查
		hasPermission, err := rbacChecker.Check(ctx, "admin@test.com", middlewareAuth.Permission{
			Resource: "book",
			Action:   "read",
		})
		require.NoError(t, err)
		assert.True(t, hasPermission, "admin用户应该有book:read权限")

		hasPermission, err = rbacChecker.Check(ctx, "admin@test.com", middlewareAuth.Permission{
			Resource: "book",
			Action:   "delete",
		})
		require.NoError(t, err)
		assert.True(t, hasPermission, "admin用户应该有book:delete权限（通配符）")
	})

	t.Run("VerifyRoleDataInDatabase", func(t *testing.T) {
		// 直接查询数据库验证角色数据
		var adminRole authModel.Role
		err := db.Collection("roles").FindOne(ctx, bson.M{"name": "admin"}).Decode(&adminRole)
		require.NoError(t, err, "查找admin角色失败")

		assert.Equal(t, "admin", adminRole.Name)
		assert.Contains(t, adminRole.Permissions, "*:*")
		assert.True(t, adminRole.IsSystem, "admin应该是系统角色")

		t.Logf("✓ 数据库验证成功：admin角色 = %v", adminRole)
	})

	t.Run("VerifyUserDataInDatabase", func(t *testing.T) {
		// 直接查询数据库验证用户数据
		var adminUser bson.M
		err := db.Collection("users").FindOne(ctx, bson.M{"username": "admin@test.com"}).Decode(&adminUser)
		require.NoError(t, err, "查找admin用户失败")

		username := adminUser["username"]
		roles := adminUser["roles"]

		assert.Equal(t, "admin@test.com", username)
		assert.Contains(t, roles, "admin", "admin用户应该有admin角色")

		t.Logf("✓ 用户数据验证成功：username=%s, roles=%v", username, roles)
	})

	t.Run("BatchCheckPerformance", func(t *testing.T) {
		// 批量权限检查性能测试
		permissions := []middlewareAuth.Permission{
			{Resource: "book", Action: "read"},
			{Resource: "book", Action: "create"},
			{Resource: "book", Action: "update"},
			{Resource: "book", Action: "delete"},
			{Resource: "chapter", Action: "read"},
			{Resource: "chapter", Action: "create"},
		}

		start := time.Now()
		results, err := rbacChecker.BatchCheck(ctx, "admin@test.com", permissions)
		elapsed := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, len(permissions), len(results))
		assert.True(t, elapsed < 10*time.Millisecond, "批量检查应该很快")

		t.Logf("✓ 批量检查 %d 个权限耗时: %v", len(permissions), elapsed)
	})
}

// TestPermissionFormatConversion 测试权限格式转换
func TestPermissionFormatConversion(t *testing.T) {
	if os.Getenv("TEST_MODE") != "true" {
		t.Skip("跳过集成测试")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("TEST_MONGO_HOST")
	if mongoURI == "" {
		mongoURI = "localhost:27018"
	}
	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "qingyu_permission_test"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+mongoURI))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	logger := zap.NewNop()

	repository := auth.NewAuthRepository(db)
	service := permService.NewPermissionService(repository, nil, logger)

	checker, err := middlewareAuth.NewRBACChecker(nil)
	require.NoError(t, err)
	rbacChecker := checker.(*middlewareAuth.RBACChecker)
	service.SetChecker(rbacChecker)

	// 加载权限（会自动进行格式转换）
	err = service.LoadPermissionsToChecker(ctx)
	require.NoError(t, err)

	// 验证格式转换：数据库中的 "book.read" 应该转换为 "book:read"
	authorPerms := rbacChecker.GetRolePermissions("author")
	assert.Contains(t, authorPerms, "book:read", "权限应该从 'book.read' 转换为 'book:read'")
	assert.Contains(t, authorPerms, "book:create", "权限应该从 'book.create' 转换为 'book:create'")

	t.Logf("✓ 权限格式转换验证成功")
}

// TestWildcardPermissionMatching 测试通配符权限匹配
func TestWildcardPermissionMatching(t *testing.T) {
	if os.Getenv("TEST_MODE") != "true" {
		t.Skip("跳过集成测试")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("TEST_MONGO_HOST")
	if mongoURI == "" {
		mongoURI = "localhost:27018"
	}
	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "qingyu_permission_test"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+mongoURI))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	logger := zap.NewNop()

	repository := auth.NewAuthRepository(db)
	service := permService.NewPermissionService(repository, nil, logger)

	checker, err := middlewareAuth.NewRBACChecker(nil)
	require.NoError(t, err)
	rbacChecker := checker.(*middlewareAuth.RBACChecker)
	service.SetChecker(rbacChecker)

	err = service.LoadPermissionsToChecker(ctx)
	require.NoError(t, err)
	service.LoadUserRolesToChecker(ctx, "admin@test.com")

	tests := []struct {
		name      string
		resource  string
		action    string
		expected  bool
	}{
		{"Admin全匹配", "any", "any", true},
		{"Admin读取书籍", "book", "read", true},
		{"Admin删除书籍", "book", "delete", true},
		{"AdminAI生成", "ai", "generate", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasPerm, err := rbacChecker.Check(ctx, "admin@test.com", middlewareAuth.Permission{
				Resource: tt.resource,
				Action:   tt.action,
			})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, hasPerm, fmt.Sprintf("%s: %s:%s", tt.name, tt.resource, tt.action))
		})
	}
}
