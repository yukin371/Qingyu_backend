package testutil

import (
	"Qingyu_backend/models/writer"
	"context"
	"os"
	"testing"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/users"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ 测试配置助手 ============

// InitTestConfig 初始化测试配置
func InitTestConfig() {
	if config.GlobalConfig == nil {
		// 从环境变量获取 MongoDB URI，如果没有则使用默认值（带认证）
		mongoURI := os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://admin:password@localhost:27017"
		}

		// 从环境变量获取数据库名，如果没有则使用默认值
		mongoDB := os.Getenv("MONGODB_DATABASE")
		if mongoDB == "" {
			mongoDB = "qingyu_test"
		}

		config.GlobalConfig = &config.Config{
			JWT: &config.JWTConfig{
				Secret:          "test-secret-key-for-testing-only",
				ExpirationHours: 24,
			},
			Database: &config.DatabaseConfig{
				Type: "mongodb",
				Primary: config.DatabaseConnection{
					Type: config.DatabaseTypeMongoDB,
					MongoDB: &config.MongoDBConfig{
						URI:      mongoURI,
						Database: mongoDB,
					},
				},
			},
			Server: &config.ServerConfig{
				Port: ":8080",
				Mode: "test",
			},
		}
	}
}

// ============ 用户相关测试助手 ============

// UserOption 用户选项函数类型
type UserOption func(*users.User)

// CreateTestUser 创建测试用户
func CreateTestUser(opts ...UserOption) *users.User {
	id := primitive.NewObjectID()
	now := time.Now()
	user := &users.User{
		IdentifiedEntity: shared.IdentifiedEntity{ID: id},
		BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
		Username:         "testuser",
		Email:            "test@example.com",
		Password:         "hashed_password_123",
		Roles:            []string{"reader"},
		Status:           users.UserStatusActive,
	}

	// 应用选项
	for _, opt := range opts {
		opt(user)
	}

	return user
}

// WithUsername 设置用户名
func WithUsername(username string) UserOption {
	return func(u *users.User) {
		u.Username = username
	}
}

// WithEmail 设置邮箱
func WithEmail(email string) UserOption {
	return func(u *users.User) {
		u.Email = email
	}
}

// WithRole 设置角色
func WithRole(role string) UserOption {
	return func(u *users.User) {
		u.Roles = []string{role}
	}
}

// WithStatus 设置状态
func WithStatus(status string) UserOption {
	return func(u *users.User) {
		u.Status = users.UserStatus(status)
	}
}

// CreateTestUsers 批量创建测试用户
func CreateTestUsers(count int) []*users.User {
	result := make([]*users.User, count)
	for i := 0; i < count; i++ {
		result[i] = CreateTestUser(
			WithUsername("testuser"+string(rune(i))),
			WithEmail("test"+string(rune(i))+"@example.com"),
		)
	}
	return result
}

// AssertUserEqual 断言用户相等
func AssertUserEqual(t *testing.T, expected, actual *users.User) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Roles, actual.Roles)
	assert.Equal(t, expected.Status, actual.Status)
}

// ============ 项目相关测试助手 ============

// ProjectOption 项目选项函数类型
type ProjectOption func(*writer.Project)

// CreateTestProject 创建测试项目
func CreateTestProject(userID string, opts ...ProjectOption) *writer.Project {
	now := time.Now()
	project := &writer.Project{
		WritingType: "novel",
		Summary:     "这是一个测试项目",
		Status:      writer.StatusDraft,
		Visibility:  writer.VisibilityPrivate,
		Statistics: writer.ProjectStats{
			TotalWords:    0,
			ChapterCount:  0,
			DocumentCount: 0,
			LastUpdateAt:  now,
		},
		Settings: writer.ProjectSettings{
			AutoBackup:     false,
			BackupInterval: 24,
		},
	}
	project.OwnedEntity.AuthorID = userID
	project.TitledEntity.Title = "测试项目"
	project.Timestamps.CreatedAt = now
	project.Timestamps.UpdatedAt = now

	// 应用选项
	for _, opt := range opts {
		opt(project)
	}

	return project
}

// WithProjectName 设置项目名称（使用Title字段）
func WithProjectName(name string) ProjectOption {
	return func(p *writer.Project) {
		p.Title = name
	}
}

// WithProjectDescription 设置项目描述（使用Summary字段）
func WithProjectDescription(description string) ProjectOption {
	return func(p *writer.Project) {
		p.Summary = description
	}
}

// WithProjectStatus 设置项目状态
func WithProjectStatus(status string) ProjectOption {
	return func(p *writer.Project) {
		p.Status = writer.ProjectStatus(status)
	}
}

// CreateTestProjects 批量创建测试项目
func CreateTestProjects(userID string, count int) []*writer.Project {
	result := make([]*writer.Project, count)
	for i := 0; i < count; i++ {
		result[i] = CreateTestProject(
			userID,
			WithProjectName("测试项目"+string(rune(i))),
		)
	}
	return result
}

// AssertProjectEqual 断言项目相等
func AssertProjectEqual(t *testing.T, expected, actual *writer.Project) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Summary, actual.Summary)
	assert.Equal(t, expected.AuthorID, actual.AuthorID)
	assert.Equal(t, expected.Status, actual.Status)
}

// ============ 文档相关测试助手 ============

// DocumentOption 文档选项函数类型
type DocumentOption func(*writer.Document)

// CreateTestDocument 创建测试文档
// 注意：此方法只创建Document元数据，不包含内容
// 如需创建文档内容，请使用CreateTestDocumentContent
func CreateTestDocument(projectID string, opts ...DocumentOption) *writer.Document {
	now := time.Now()
	objectID := primitive.NewObjectID()
	var projectObjectID primitive.ObjectID
	if projectID != "" {
		projectObjectID, _ = primitive.ObjectIDFromHex(projectID)
	}
	doc := &writer.Document{}
	doc.ID = objectID
	doc.Type = writer.TypeChapter
	doc.Status = "draft"
	doc.Title = "测试文档"
	doc.ProjectID = projectObjectID
	doc.CreatedAt = now
	doc.UpdatedAt = now

	// 应用选项
	for _, opt := range opts {
		opt(doc)
	}

	return doc
}

// CreateTestDocumentContent 创建测试文档内容
func CreateTestDocumentContent(documentID string, content string) *writer.DocumentContent {
	id := primitive.NewObjectID()
	var documentObjectID primitive.ObjectID
	if documentID != "" {
		documentObjectID, _ = primitive.ObjectIDFromHex(documentID)
	}
	now := time.Now()
	return &writer.DocumentContent{
		ID:               id,
		DocumentID:       documentObjectID,
		Content:          content,
		ContentType:      "markdown",
		WordCount:        len([]rune(content)),
		CharCount:        len(content),
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
		LastSavedAt:      now,
	}
}

// WithDocumentTitle 设置文档标题
func WithDocumentTitle(title string) DocumentOption {
	return func(d *writer.Document) {
		d.Title = title
	}
}

// WithDocumentContent 设置文档内容（已废弃）
// 注意：Document模型不再包含Content字段
// 请使用DocumentContent模型来处理文档内容
// 此函数保留用于向后兼容，但实际上不会产生任何效果
func WithDocumentContent(content string) DocumentOption {
	return func(d *writer.Document) {
		// 不再设置Content字段，保留函数签名用于兼容
	}
}

// WithDocumentStatus 设置文档状态
func WithDocumentStatus(status string) DocumentOption {
	return func(d *writer.Document) {
		d.Status = status
	}
}

// ============ 上下文助手 ============

// CreateTestContext 创建测试上下文
func CreateTestContext() context.Context {
	return context.Background()
}

// CreateTestContextWithUser 创建带用户信息的测试上下文
func CreateTestContextWithUser(userID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", userID)
	return ctx
}

// CreateTestContextWithTimeout 创建带超时的测试上下文
func CreateTestContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// ============ 通用断言助手 ============

// AssertNoErrorWithMessage 断言无错误并输出自定义消息
func AssertNoErrorWithMessage(t *testing.T, err error, message string) {
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertErrorContains 断言错误包含指定字符串
func AssertErrorContains(t *testing.T, err error, substr string) {
	assert.Error(t, err)
	assert.Contains(t, err.Error(), substr)
}

// AssertTimeAlmostEqual 断言时间几乎相等（允许1秒误差）
func AssertTimeAlmostEqual(t *testing.T, expected, actual time.Time) {
	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}
	assert.True(t, diff < time.Second, "时间差异超过1秒: %v", diff)
}

// ============ 数据清理助手 ============

// CleanupFunc 清理函数类型
type CleanupFunc func()

// RegisterCleanup 注册清理函数
func RegisterCleanup(t *testing.T, cleanup CleanupFunc) {
	t.Cleanup(cleanup)
}

// ============ 随机数据生成助手 ============

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

// RandomEmail 生成随机邮箱
func RandomEmail() string {
	return RandomString(10) + "@test.com"
}

// RandomInt 生成随机整数
func RandomInt(min, max int) int {
	return min + int(time.Now().UnixNano()%(int64(max-min)))
}

// ============ Repository Filter 助手 ============

// SimpleFilter 简单的Filter实现，用于测试
type SimpleFilter struct {
	Conditions map[string]interface{}
	SortFields map[string]int
	Fields     []string
}

// GetConditions 返回筛选条件
func (f *SimpleFilter) GetConditions() map[string]interface{} {
	if f.Conditions == nil {
		return make(map[string]interface{})
	}
	return f.Conditions
}

// GetSort 返回排序字段
func (f *SimpleFilter) GetSort() map[string]int {
	if f.SortFields == nil {
		return map[string]int{"createdAt": -1}
	}
	return f.SortFields
}

// GetFields 返回字段选择
func (f *SimpleFilter) GetFields() []string {
	return f.Fields
}

// Validate 验证过滤器
func (f *SimpleFilter) Validate() error {
	return nil // 简单实现，测试用
}
