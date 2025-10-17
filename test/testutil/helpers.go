package testutil

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/document"
	"Qingyu_backend/models/users"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ 用户相关测试助手 ============

// UserOption 用户选项函数类型
type UserOption func(*users.User)

// CreateTestUser 创建测试用户
func CreateTestUser(opts ...UserOption) *users.User {
	user := &users.User{
		ID:        primitive.NewObjectID().Hex(),
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashed_password_123",
		Role:      "user",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
		u.Role = role
	}
}

// WithStatus 设置状态
func WithStatus(status string) UserOption {
	return func(u *users.User) {
		u.Status = status
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
	assert.Equal(t, expected.Role, actual.Role)
	assert.Equal(t, expected.Status, actual.Status)
}

// ============ 项目相关测试助手 ============

// ProjectOption 项目选项函数类型
type ProjectOption func(*document.Project)

// CreateTestProject 创建测试项目
func CreateTestProject(userID string, opts ...ProjectOption) *document.Project {
	project := &document.Project{
		ID:          primitive.NewObjectID().Hex(),
		Name:        "测试项目",
		Description: "这是一个测试项目",
		AuthorID:    userID,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 应用选项
	for _, opt := range opts {
		opt(project)
	}

	return project
}

// WithProjectName 设置项目名称
func WithProjectName(name string) ProjectOption {
	return func(p *document.Project) {
		p.Name = name
	}
}

// WithProjectDescription 设置项目描述
func WithProjectDescription(description string) ProjectOption {
	return func(p *document.Project) {
		p.Description = description
	}
}

// WithProjectStatus 设置项目状态
func WithProjectStatus(status string) ProjectOption {
	return func(p *document.Project) {
		p.Status = status
	}
}

// CreateTestProjects 批量创建测试项目
func CreateTestProjects(userID string, count int) []*document.Project {
	result := make([]*document.Project, count)
	for i := 0; i < count; i++ {
		result[i] = CreateTestProject(
			userID,
			WithProjectName("测试项目"+string(rune(i))),
		)
	}
	return result
}

// AssertProjectEqual 断言项目相等
func AssertProjectEqual(t *testing.T, expected, actual *document.Project) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.AuthorID, actual.AuthorID)
	assert.Equal(t, expected.Status, actual.Status)
}

// ============ 文档相关测试助手 ============

// DocumentOption 文档选项函数类型
type DocumentOption func(*document.Document)

// CreateTestDocument 创建测试文档
func CreateTestDocument(projectID string, opts ...DocumentOption) *document.Document {
	doc := &document.Document{
		ID:        primitive.NewObjectID().Hex(),
		ProjectID: projectID,
		Title:     "测试文档",
		Content:   "这是测试内容",
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 应用选项
	for _, opt := range opts {
		opt(doc)
	}

	return doc
}

// WithDocumentTitle 设置文档标题
func WithDocumentTitle(title string) DocumentOption {
	return func(d *document.Document) {
		d.Title = title
	}
}

// WithDocumentContent 设置文档内容
func WithDocumentContent(content string) DocumentOption {
	return func(d *document.Document) {
		d.Content = content
	}
}

// WithDocumentStatus 设置文档状态
func WithDocumentStatus(status string) DocumentOption {
	return func(d *document.Document) {
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
