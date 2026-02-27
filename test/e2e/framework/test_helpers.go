//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/global"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/users"
	userRepo "Qingyu_backend/repository/mongodb/user"
)

// TestTestHelpersCompilation 测试辅助函数编译验证
func TestTestHelpersCompilation(t *testing.T) {
	// 这个测试仅用于验证代码能够编译
	// 实际的E2E测试在 examples 包中
	t.Skip("E2E框架测试在examples包中")
}

// TestHelpers E2E 测试辅助工具集
type TestHelpers struct {
	env *TestEnvironment
}

// Helpers 获取测试辅助工具
func (env *TestEnvironment) Helpers() *TestHelpers {
	return &TestHelpers{env: env}
}

// ========================================
// 认证辅助函数
// ========================================

// RegisterAndLogin 注册并登录用户，返回 token
func (th *TestHelpers) RegisterAndLogin(username, email, password string) string {
	// 1. 注册用户
	registerData := map[string]interface{}{
		"username": username,
		"email":    email,
		"password": password,
	}
	w := th.env.DoRequest("POST", "/api/v1/register", registerData, "")

	// 如果用户已存在，直接登录
	if w.Code != 200 && w.Code != 201 {
		th.env.T.Logf("用户 %s 可能已存在，尝试直接登录", username)
		return th.Login(username, password)
	}

	// 2. 登录获取 token
	return th.Login(username, password)
}

// Login 用户登录并返回 token
func (th *TestHelpers) Login(username, password string) string {
	loginData := map[string]interface{}{
		"username": username,
		"password": password,
	}
	w := th.env.DoRequest("POST", "/api/v1/login", loginData, "")

	if w.Code != 200 {
		th.env.T.Logf("登录失败: username=%s, status=%d, response=%s",
			username, w.Code, w.Body.String())
		return ""
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(th.env.T, err, "解析登录响应失败")

	data, ok := response["data"].(map[string]interface{})
	require.True(th.env.T, ok, "响应数据格式错误")

	token, ok := data["token"].(string)
	require.True(th.env.T, ok, "获取 token 失败")

	th.env.LogSuccess("用户登录: %s", username)
	return token
}

// LoginAsAdmin 使用管理员账号登录
func (th *TestHelpers) LoginAsAdmin() string {
	return th.Login("admin", "Admin@123456")
}

// LoginAsTestUser 使用测试用户登录
func (th *TestHelpers) LoginAsTestUser() string {
	return th.Login("test_user01", "Test@123456")
}

// ========================================
// 用户数据辅助函数
// ========================================

// CreateTestUser 创建测试用户（直接操作数据库）
func (th *TestHelpers) CreateTestUser(username, email, password string, roles []string) *users.User {
	userRepository := userRepo.NewMongoUserRepository(global.DB)

	user := &users.User{
		Username: username,
		Email:    email,
		Password: password, // 应该是哈希后的密码
		Roles:    roles,
		BaseEntity: shared.BaseEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	err := userRepository.Create(context.Background(), user)
	require.NoError(th.env.T, err, "创建测试用户失败")

	th.env.LogSuccess("创建测试用户: %s (roles: %v)", username, roles)

	// 注册清理函数
	th.env.RegisterCleanup(func() {
		userRepository.Delete(context.Background(), user.ID.Hex())
		th.env.T.Logf("清理测试用户: %s", username)
	})

	return user
}

// GetUserIDByUsername 通过用户名获取用户ID
func (th *TestHelpers) GetUserIDByUsername(username string) string {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	user, err := userRepository.GetByUsername(context.Background(), username)
	require.NoError(th.env.T, err, "获取用户失败")
	return user.ID.Hex()
}

// SetUserVIP 设置用户为VIP（直接操作数据库）
func (th *TestHelpers) SetUserVIP(userID string, vipLevel int) {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	updates := map[string]interface{}{
		"vip_level": vipLevel,
	}
	err := userRepository.Update(context.Background(), userID, updates)
	require.NoError(th.env.T, err, "设置用户VIP失败")

	th.env.LogSuccess("设置用户VIP: %s -> level %d", userID, vipLevel)
}

// ========================================
// 测试数据准备函数
// ========================================

// PrepareTestBookData 准备测试书籍数据
func (th *TestHelpers) PrepareTestBookData() map[string]interface{} {
	bookID := fmt.Sprintf("test-book-%d", time.Now().Unix())
	return map[string]interface{}{
		"book_id":    bookID,
		"title":      "测试书籍",
		"author":     "测试作者",
		"category":   "小说",
		"status":     "published",
		"is_vip":     false,
		"created_at": time.Now(),
	}
}

// PrepareTestChapterData 准备测试章节数据
func (th *TestHelpers) PrepareTestChapterData(bookID string, chapterNum int) map[string]interface{} {
	chapterID := fmt.Sprintf("test-chapter-%s-%d", bookID, chapterNum)
	return map[string]interface{}{
		"chapter_id": chapterID,
		"book_id":    bookID,
		"chapter_num": chapterNum,
		"title":      fmt.Sprintf("第%d章", chapterNum),
		"content":    "这是测试章节内容",
		"is_vip":     false,
		"created_at": time.Now(),
	}
}

// ========================================
// HTTP请求辅助函数
// ========================================

// DoRequestWithRetry 带重试的HTTP请求
func (th *TestHelpers) DoRequestWithRetry(method, path string, body interface{}, token string, maxRetries int) *httptest.ResponseRecorder {
	var w *httptest.ResponseRecorder

	for i := 0; i < maxRetries; i++ {
		w = th.env.DoRequest(method, path, body, token)

		// 如果成功或不是服务器错误，直接返回
		if w.Code < 500 {
			return w
		}

		th.env.T.Logf("请求失败（尝试 %d/%d）: %s %s -> %d",
			i+1, maxRetries, method, path, w.Code)

		// 等待后重试
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}

	return w
}

// DoBatchRequest 批量执行请求
func (th *TestHelpers) DoBatchRequest(requests []map[string]interface{}) []*httptest.ResponseRecorder {
	responses := make([]*httptest.ResponseRecorder, len(requests))

	for i, req := range requests {
		method, _ := req["method"].(string)
		path, _ := req["path"].(string)
		body, _ := req["body"]
		token, _ := req["token"].(string)

		responses[i] = th.env.DoRequest(method, path, body, token)
	}

	return responses
}

// ========================================
// 断言辅助函数
// ========================================

// AssertSuccess 断言请求成功
func (th *TestHelpers) AssertSuccess(w *httptest.ResponseRecorder, expectedCode int) map[string]interface{} {
	require.Equal(th.env.T, expectedCode, w.Code,
		"期望状态码 %d，实际 %d。响应: %s", expectedCode, w.Code, w.Body.String())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(th.env.T, err, "解析响应失败")

	return response
}

// AssertError 断言请求失败
func (th *TestHelpers) AssertError(w *httptest.ResponseRecorder, expectedCode int, expectedMessage string) {
	require.Equal(th.env.T, expectedCode, w.Code,
		"期望错误状态码 %d，实际 %d", expectedCode, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(th.env.T, err, "解析错误响应失败")

	if expectedMessage != "" {
		message, _ := response["message"].(string)
		assert.Contains(th.env.T, message, expectedMessage,
			"错误消息应包含: %s", expectedMessage)
	}
}

// AssertPagination 断言分页响应
func (th *TestHelpers) AssertPagination(w *httptest.ResponseRecorder, expectedMinItems int) map[string]interface{} {
	response := th.AssertSuccess(w, 200)

	data, ok := response["data"].(map[string]interface{})
	require.True(th.env.T, ok, "响应数据应该是对象")

	items, ok := data["items"].([]interface{})
	require.True(th.env.T, ok, "响应应包含items字段")
	require.GreaterOrEqual(th.env.T, len(items), expectedMinItems,
		"至少应有 %d 个项目", expectedMinItems)

	total, ok := data["total"].(float64)
	require.True(th.env.T, ok, "响应应包含total字段")
	th.env.T.Logf("分页结果: %d/%d", len(items), int(total))

	return response
}

// ========================================
// 时间辅助函数
// ========================================

// WaitForCondition 等待条件满足
func (th *TestHelpers) WaitForCondition(condition func() bool, timeout time.Duration, checkInterval time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(checkInterval)
	}

	return false
}

// WaitForStatusCode 等待特定状态码
func (th *TestHelpers) WaitForStatusCode(method, path string, body interface{}, token string, expectedCode int, timeout time.Duration) bool {
	return th.WaitForCondition(func() bool {
		w := th.env.DoRequest(method, path, body, token)
		return w.Code == expectedCode
	}, timeout, 500*time.Millisecond)
}

// ========================================
// 性能测试辅助函数
// ========================================

// MeasureRequestTime 测量请求时间
func (th *TestHelpers) MeasureRequestTime(method, path string, body interface{}, token string) time.Duration {
	start := time.Now()
	th.env.DoRequest(method, path, body, token)
	return time.Since(start)
}

// BenchmarkRequest 性能测试请求
func (th *TestHelpers) BenchmarkRequest(method, path string, body interface{}, token string, iterations int) BenchmarkResult {
	var totalTime time.Duration
	var successCount int
	var minDuration, maxDuration time.Duration

	for i := 0; i < iterations; i++ {
		duration := th.MeasureRequestTime(method, path, body, token)
		totalTime += duration

		if i == 0 {
			minDuration = duration
			maxDuration = duration
		} else {
			if duration < minDuration {
				minDuration = duration
			}
			if duration > maxDuration {
				maxDuration = duration
			}
		}

		successCount++
	}

	avgDuration := totalTime / time.Duration(successCount)

	return BenchmarkResult{
		Iterations:    iterations,
		SuccessCount:  successCount,
		TotalTime:     totalTime,
		AvgDuration:   avgDuration,
		MinDuration:   minDuration,
		MaxDuration:   maxDuration,
	}
}

// BenchmarkResult 性能测试结果
type BenchmarkResult struct {
	Iterations   int
	SuccessCount int
	TotalTime    time.Duration
	AvgDuration  time.Duration
	MinDuration  time.Duration
	MaxDuration  time.Duration
}

// LogBenchmarkResult 记录性能测试结果
func (br *BenchmarkResult) LogBenchmarkResult(t *testing.T) {
	t.Logf("========== 性能测试结果 ==========")
	t.Logf("迭代次数: %d", br.Iterations)
	t.Logf("成功次数: %d", br.SuccessCount)
	t.Logf("总耗时: %v", br.TotalTime)
	t.Logf("平均耗时: %v", br.AvgDuration)
	t.Logf("最小耗时: %v", br.MinDuration)
	t.Logf("最大耗时: %v", br.MaxDuration)
	t.Logf("================================")
}

// ========================================
// 数据清理辅助函数
// ========================================

// CleanupTestUser 清理测试用户
func (th *TestHelpers) CleanupTestUser(username string) {
	userRepository := userRepo.NewMongoUserRepository(global.DB)
	user, err := userRepository.GetByUsername(context.Background(), username)
	if err == nil {
		userRepository.Delete(context.Background(), user.ID.Hex())
		th.env.T.Logf("清理测试用户: %s", username)
	}
}

// CleanupTestBook 清理测试书籍
func (th *TestHelpers) CleanupTestBook(bookID string) {
	// 清理书籍相关数据
	collections := []string{
		"books",
		"chapters",
		"reading_progress",
		"comments",
		"collections",
	}

	for _, coll := range collections {
		_, err := global.DB.Collection(coll).DeleteOne(
			context.Background(),
			map[string]interface{}{"book_id": bookID},
		)
		if err == nil {
			th.env.T.Logf("清理 %s 中的书籍 %s", coll, bookID)
		}
	}
}

// CleanupByPrefix 根据前缀清理测试数据
func (th *TestHelpers) CleanupByPrefix(prefix string) {
	collections := []string{
		"users", "books", "chapters",
	}

	for _, coll := range collections {
		filter := map[string]interface{}{
			"$or": []map[string]interface{}{
				{"username": map[string]interface{}{"$regex": "^" + prefix}},
				{"title": map[string]interface{}{"$regex": "^" + prefix}},
			},
		}

		result, err := global.DB.Collection(coll).DeleteMany(context.Background(), filter)
		if err == nil && result.DeletedCount > 0 {
			th.env.T.Logf("清理 %s: %d 条记录（前缀: %s）", coll, result.DeletedCount, prefix)
		}
	}
}

// ========================================
// 流程辅助函数
// ========================================

// CompleteUserFlow 完整用户流程测试
func (th *TestHelpers) CompleteUserFlow(username, email, password string) UserFlowResult {
	result := UserFlowResult{}

	// 1. 注册并登录
	result.Token = th.RegisterAndLogin(username, email, password)
	require.NotEmpty(th.env.T, result.Token, "注册登录失败")

	// 2. 获取用户信息
	w := th.env.DoRequest("GET", "/api/v1/users/profile", nil, result.Token)
	if w.Code == 200 {
		response := th.AssertSuccess(w, 200)
		data, _ := response["data"].(map[string]interface{})
		if userID, ok := data["user_id"].(string); ok {
			result.UserID = userID
		}
	}

	// 3. 浏览书城
	w = th.env.DoRequest("GET", "/api/v1/bookstore/homepage", nil, "")
	result.BookstoreAccess = w.Code == 200

	return result
}

// UserFlowResult 用户流程结果
type UserFlowResult struct {
	Token          string
	UserID         string
	BookstoreAccess bool
}

// CompleteReaderFlow 完整读者流程测试
func (th *TestHelpers) CompleteReaderFlow(token, bookID string) ReaderFlowResult {
	result := ReaderFlowResult{Token: token}

	// 1. 获取书籍详情
	w := th.env.DoRequest("GET", fmt.Sprintf("/api/v1/bookstore/books/%s", bookID), nil, "")
	result.BookAccess = w.Code == 200

	// 2. 获取章节列表
	w = th.env.DoRequest("GET", fmt.Sprintf("/api/v1/reader/books/%s/chapters", bookID), nil, token)
	result.ChapterListAccess = w.Code == 200

	// 3. 阅读章节（假设第一个章节）
	if result.ChapterListAccess {
		// 这里需要根据实际API调整
	}

	return result
}

// ReaderFlowResult 读者流程结果
type ReaderFlowResult struct {
	Token              string
	BookAccess         bool
	ChapterListAccess  bool
}

// ========================================
// WebSocket测试辅助函数
// ========================================

// CreateWebSocketClient 创建WebSocket客户端（需要实际实现）
func (th *TestHelpers) CreateWebSocketClient(path string, token string) interface{} {
	// TODO: 实现WebSocket客户端创建
	th.env.T.Skip("WebSocket支持待实现")
	return nil
}

// ========================================
// 并发测试辅助函数
// ========================================

// ConcurrentRequest 并发请求
func (th *TestHelpers) ConcurrentRequest(method, path string, body interface{}, token string, concurrency int, iterations int) ConcurrentResult {
	result := ConcurrentResult{
		TotalRequests: concurrency * iterations,
	}

	// 创建结果通道
	resultChan := make(chan *httptest.ResponseRecorder, result.TotalRequests)

	// 执行并发请求
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < iterations; j++ {
				w := th.env.DoRequest(method, path, body, token)
				resultChan <- w
			}
		}()
	}

	// 收集结果
	for i := 0; i < result.TotalRequests; i++ {
		w := <-resultChan

		switch w.Code {
		case 200, 201:
			result.SuccessCount++
		case 429:
			result.RateLimitedCount++
		default:
			result.ErrorCount++
		}
	}

	result.SuccessRate = float64(result.SuccessCount) / float64(result.TotalRequests) * 100

	return result
}

// ConcurrentResult 并发测试结果
type ConcurrentResult struct {
	TotalRequests    int
	SuccessCount     int
	ErrorCount       int
	RateLimitedCount int
	SuccessRate      float64
}

// LogConcurrentResult 记录并发测试结果
func (cr *ConcurrentResult) LogConcurrentResult(t *testing.T) {
	t.Logf("========== 并发测试结果 ==========")
	t.Logf("总请求数: %d", cr.TotalRequests)
	t.Logf("成功数: %d", cr.SuccessCount)
	t.Logf("错误数: %d", cr.ErrorCount)
	t.Logf("限流数: %d", cr.RateLimitedCount)
	t.Logf("成功率: %.2f%%", cr.SuccessRate)
	t.Logf("================================")
}
