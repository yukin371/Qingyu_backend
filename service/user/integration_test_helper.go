package user

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/config"
	usersModel "Qingyu_backend/models/users"
	authInterface "Qingyu_backend/repository/interfaces/auth"
	roleRepo "Qingyu_backend/repository/mongodb/auth"
	repoUser "Qingyu_backend/repository/mongodb/user"
)

// IntegrationTestEnvironment 集成测试环境
type IntegrationTestEnvironment struct {
	UserService    *UserServiceImpl
	AuthRepository authInterface.RoleRepository
	DB             *mongo.Database
	DBName         string
	CleanupFunc    func()
	TestConfig     *TestConfig
	JWTTestHelper  *JWTTestHelper
}

// SetupIntegrationTestEnvironment 设置集成测试环境
// 返回测试环境和清理函数
func SetupIntegrationTestEnvironment(t *testing.T) *IntegrationTestEnvironment {
	t.Helper()

	// 跳过短测试
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
	}

	// 获取测试配置
	testCfg := GetTestConfig()

	// 初始化全局配置（如果还没有初始化）
	if config.GlobalConfig == nil {
		config.GlobalConfig = &config.Config{
			JWT: &config.JWTConfig{
				Secret:          testCfg.JWTSecret,
				ExpirationHours: int(testCfg.JWTExpiration.Hours()),
			},
			Database: &config.DatabaseConfig{
				Type: "mongodb",
				Primary: config.DatabaseConnection{
					Type: config.DatabaseTypeMongoDB,
					MongoDB: &config.MongoDBConfig{
						URI:      testCfg.MongoURI,
						Database: testCfg.DatabaseName,
					},
				},
			},
			Server: &config.ServerConfig{
				Port: ":9090",
				Mode: "test",
			},
		}
	} else {
		// 确保测试使用当前测试配置，避免受其他测试污染
		if config.GlobalConfig.JWT == nil {
			config.GlobalConfig.JWT = &config.JWTConfig{}
		}
		config.GlobalConfig.JWT.Secret = testCfg.JWTSecret
		config.GlobalConfig.JWT.ExpirationHours = int(testCfg.JWTExpiration.Hours())

		if config.GlobalConfig.Database == nil {
			config.GlobalConfig.Database = &config.DatabaseConfig{}
		}
		config.GlobalConfig.Database.Type = "mongodb"
		config.GlobalConfig.Database.Primary = config.DatabaseConnection{
			Type: config.DatabaseTypeMongoDB,
			MongoDB: &config.MongoDBConfig{
				URI:      testCfg.MongoURI,
				Database: testCfg.DatabaseName,
			},
		}
	}

	// 创建MongoDB客户端（不使用service container避免循环依赖）
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer connectCancel()
	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(testCfg.MongoURI))
	if err != nil {
		t.Skipf("MongoDB不可用，跳过集成测试: %v", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		t.Skipf("MongoDB不可用，跳过集成测试: %v", err)
	}

	db := client.Database(testCfg.DatabaseName)

	// 创建UserRepository
	userRepo := repoUser.NewMongoUserRepository(db)

	// 创建AuthRepository (使用RoleRepository)
	authRepository := roleRepo.NewRoleRepository(db)

	// 创建UserService
	userService := NewUserService(userRepo, authRepository)
	userServiceImpl, ok := userService.(*UserServiceImpl)
	require.True(t, ok, "UserService类型转换失败")

	// 初始化UserService
	ctx := context.Background()
	err = userServiceImpl.Initialize(ctx)
	require.NoError(t, err, "初始化UserService失败")

	// 创建JWT测试辅助工具
	jwtHelper := NewJWTTestHelper(testCfg.JWTSecret)

	// 清理函数
	cleanup := func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 清理测试集合
		_ = db.Collection("users").Drop(cleanupCtx)
		_ = db.Collection("roles").Drop(cleanupCtx)
		_ = db.Collection("verification_codes").Drop(cleanupCtx)
		_ = db.Collection("password_reset_tokens").Drop(cleanupCtx)
		_ = db.Collection("email_verification_codes").Drop(cleanupCtx)
		_ = db.Collection("sessions").Drop(cleanupCtx)

		// 断开连接
		_ = client.Disconnect(cleanupCtx)
	}

	return &IntegrationTestEnvironment{
		UserService:    userServiceImpl,
		AuthRepository: authRepository,
		DB:             db,
		DBName:         testCfg.DatabaseName,
		CleanupFunc:    cleanup,
		TestConfig:     testCfg,
		JWTTestHelper:  jwtHelper,
	}
}

// CleanupTestData 清理测试数据
func (env *IntegrationTestEnvironment) CleanupTestData(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 清理测试集合
	_ = env.DB.Collection("users").Drop(ctx)
	_ = env.DB.Collection("roles").Drop(ctx)
	_ = env.DB.Collection("verification_codes").Drop(ctx)
	_ = env.DB.Collection("password_reset_tokens").Drop(ctx)
	_ = env.DB.Collection("email_verification_codes").Drop(ctx)
	_ = env.DB.Collection("sessions").Drop(ctx)
}

// GenerateUniqueTestUser 生成唯一测试用户数据
func (env *IntegrationTestEnvironment) GenerateUniqueTestUser(t *testing.T) *usersModel.User {
	t.Helper()

	timestamp := time.Now().UnixNano()
	return &usersModel.User{
		Username:      userStringFormat("testuser_%d", timestamp),
		Email:         userStringFormat("testuser_%d@example.com", timestamp),
		Password:      "TestPassword123!",
		Roles:         []string{"reader"},
		Status:        usersModel.UserStatusActive,
		EmailVerified: false,
		PhoneVerified: false,
	}
}

// GenerateUniqueTestUserWithPrefix 使用前缀生成唯一测试用户
func (env *IntegrationTestEnvironment) GenerateUniqueTestUserWithPrefix(t *testing.T, prefix string) *usersModel.User {
	t.Helper()

	timestamp := time.Now().UnixNano()
	return &usersModel.User{
		Username:      userStringFormat("%s_%d", prefix, timestamp),
		Email:         userStringFormat("%s_%d@example.com", prefix, timestamp),
		Password:      "TestPassword123!",
		Roles:         []string{"reader"},
		Status:        usersModel.UserStatusActive,
		EmailVerified: false,
		PhoneVerified: false,
	}
}

// CreateTestUserInDB 在数据库中创建测试用户
func (env *IntegrationTestEnvironment) CreateTestUserInDB(t *testing.T, testUser *usersModel.User) string {
	t.Helper()

	ctx := context.Background()

	// 设置密码
	err := testUser.SetPassword(testUser.Password)
	require.NoError(t, err, "设置密码失败")

	// 创建UserRepository
	userRepo := repoUser.NewMongoUserRepository(env.DB)

	// 创建用户
	err = userRepo.Create(ctx, testUser)
	require.NoError(t, err, "创建用户失败")

	return testUser.ID.Hex()
}

// CreateDefaultTestUser 创建默认测试用户
func (env *IntegrationTestEnvironment) CreateDefaultTestUser(t *testing.T) (userID, username, email, password string) {
	t.Helper()

	testUser := env.GenerateUniqueTestUser(t)
	password = testUser.Password
	userID = env.CreateTestUserInDB(t, testUser)
	username = testUser.Username
	email = testUser.Email

	return userID, username, email, password
}

// CreateTestUserWithStatus 创建指定状态的测试用户
func (env *IntegrationTestEnvironment) CreateTestUserWithStatus(t *testing.T, status usersModel.UserStatus) string {
	t.Helper()

	testUser := env.GenerateUniqueTestUser(t)
	testUser.Status = status

	return env.CreateTestUserInDB(t, testUser)
}

// CreateVerifiedTestUser 创建已验证邮箱的测试用户
func (env *IntegrationTestEnvironment) CreateVerifiedTestUser(t *testing.T) (userID, username, email, password string) {
	t.Helper()

	testUser := env.GenerateUniqueTestUser(t)
	password = testUser.Password
	testUser.EmailVerified = true

	userID = env.CreateTestUserInDB(t, testUser)
	username = testUser.Username
	email = testUser.Email

	return userID, username, email, password
}

// CreateTestUserWithRoles 创建带指定角色的测试用户
func (env *IntegrationTestEnvironment) CreateTestUserWithRoles(t *testing.T, roles []string) string {
	t.Helper()

	testUser := env.GenerateUniqueTestUser(t)
	testUser.Roles = roles

	return env.CreateTestUserInDB(t, testUser)
}

// AssertUserExists 断言用户存在
func (env *IntegrationTestEnvironment) AssertUserExists(t *testing.T, userID string) *usersModel.User {
	t.Helper()

	ctx := context.Background()
	userRepo := repoUser.NewMongoUserRepository(env.DB)

	user, err := userRepo.GetByID(ctx, userID)
	require.NoError(t, err, "获取用户失败")
	require.NotNil(t, user, "用户不存在")

	return user
}

// AssertUserNotExists 断言用户不存在
func (env *IntegrationTestEnvironment) AssertUserNotExists(t *testing.T, userID string) {
	t.Helper()

	ctx := context.Background()
	userRepo := repoUser.NewMongoUserRepository(env.DB)

	_, err := userRepo.GetByID(ctx, userID)
	require.Error(t, err, "期望用户不存在，但找到了用户")
}

// AssertUserStatus 断言用户状态
func (env *IntegrationTestEnvironment) AssertUserStatus(t *testing.T, userID string, expectedStatus usersModel.UserStatus) {
	t.Helper()

	user := env.AssertUserExists(t, userID)
	require.Equal(t, expectedStatus, user.Status, "用户状态不匹配")
}

// AssertUserEmailVerified 断言用户邮箱已验证
func (env *IntegrationTestEnvironment) AssertUserEmailVerified(t *testing.T, userID string, expectedVerified bool) {
	t.Helper()

	user := env.AssertUserExists(t, userID)
	require.Equal(t, expectedVerified, user.EmailVerified, "邮箱验证状态不匹配")
}

// GetEnvOrDefault 获取环境变量或使用默认值
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDockerEnvironment 检测是否在Docker环境中运行
func IsDockerEnvironment() bool {
	return GetEnvOrDefault("DOCKER_ENV", "false") == "true" ||
		os.Getenv("MONGODB_URI") != ""
}

// GetTestMongoDBURI 获取测试用MongoDB URI
func GetTestMongoDBURI() string {
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}
	return "mongodb://localhost:27017"
}

// GetTestDatabaseName 获取测试数据库名称
func GetTestDatabaseName() string {
	if name := os.Getenv("MONGODB_DATABASE"); name != "" {
		return name
	}
	return "qingyu_test"
}

// userStringFormat 格式化字符串的辅助函数
func userStringFormat(format string, args ...interface{}) string {
	// 简单的字符串格式化
	if len(args) == 0 {
		return format
	}
	// 这里使用fmt.Sprintf，但为了避免循环导入，我们简化处理
	result := format
	for _, arg := range args {
		switch v := arg.(type) {
		case int64:
			// 替换%d
			for i := 0; i < len(result)-1; i++ {
				if result[i] == '%' && result[i+1] == 'd' {
					result = result[:i] + fmt.Sprintf("%d", v) + result[i+2:]
					break
				}
			}
		case string:
			// 替换%s
			for i := 0; i < len(result)-1; i++ {
				if result[i] == '%' && result[i+1] == 's' {
					result = result[:i] + v + result[i+2:]
					break
				}
			}
		}
	}
	return result
}
