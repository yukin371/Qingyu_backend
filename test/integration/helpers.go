package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
)

// ========================================
// API 路径常量
// ========================================

const (
	// 基础路径
	APIBasePath = "/api/v1"

	// 认证相关
	LoginPath    = APIBasePath + "/login"
	RegisterPath = APIBasePath + "/register"

	// 用户相关
	UserProfilePath  = APIBasePath + "/user/profile"
	UserPasswordPath = APIBasePath + "/user/password"

	// 阅读器相关
	ReaderBooksPath          = APIBasePath + "/reader/books"
	ReaderChaptersPath       = APIBasePath + "/reader/chapters"
	ReaderProgressPath       = APIBasePath + "/reader/progress"
	ReaderAnnotationsPath    = APIBasePath + "/reader/annotations"
	ReaderCommentsPath       = APIBasePath + "/reader/comments"
	ReaderCollectionsPath    = APIBasePath + "/reader/collections"
	ReaderLikesPath          = APIBasePath + "/reader/likes"
	ReaderReadingHistoryPath = APIBasePath + "/reader/reading-history"

	// 书城相关
	BookstoreHomePath    = APIBasePath + "/bookstore/homepage"
	BookstoreBooksPath   = APIBasePath + "/bookstore/books"
	BookstoreSearchPath  = APIBasePath + "/bookstore/books/search"
	BookstoreRankingPath = APIBasePath + "/bookstore/rankings"
)

// ========================================
// 测试辅助结构
// ========================================

// TestHelper 测试辅助工具
type TestHelper struct {
	t      *testing.T
	router *gin.Engine
	ctx    context.Context
}

// NewTestHelper 创建测试辅助工具
func NewTestHelper(t *testing.T, router *gin.Engine) *TestHelper {
	return &TestHelper{
		t:      t,
		router: router,
		ctx:    context.Background(),
	}
}

// ========================================
// 认证相关辅助函数
// ========================================

// LoginUser 用户登录并返回token
func (h *TestHelper) LoginUser(username, password string) string {
	loginData := map[string]interface{}{
		"username": username,
		"password": password,
	}

	body, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", LoginPath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		h.t.Logf("❌ 登录失败\n"+
			"  用户名: %s\n"+
			"  状态码: %d (期望: 200)\n"+
			"  响应: %s",
			username, w.Code, w.Body.String())
		return ""
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		h.t.Logf("❌ 解析登录响应失败: %v", err)
		return ""
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		h.t.Logf("❌ 响应数据格式错误: %+v", response)
		return ""
	}

	token, ok := data["token"].(string)
	if !ok {
		h.t.Logf("❌ 获取token失败: %+v", data)
		return ""
	}

	h.t.Logf("✓ 登录成功: %s (token: %s...)", username, token[:20])
	return token
}

// LoginTestUser 登录默认测试用户
func (h *TestHelper) LoginTestUser() string {
	return h.LoginUser("test_user01", "Test@123456")
}

// ========================================
// HTTP 请求辅助函数
// ========================================

// DoRequest 执行HTTP请求
func (h *TestHelper) DoRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)
	return w
}

// DoAuthRequest 执行需要认证的请求
func (h *TestHelper) DoAuthRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	require.NotEmpty(h.t, token, "Token不能为空，请先登录")
	return h.DoRequest(method, path, body, token)
}

// ========================================
// 响应断言辅助函数
// ========================================

// AssertSuccess 断言请求成功
func (h *TestHelper) AssertSuccess(w *httptest.ResponseRecorder, expectedStatus int, msgAndArgs ...interface{}) map[string]interface{} {
	// 构建详细的错误信息
	msg := ""
	if len(msgAndArgs) > 0 {
		if format, ok := msgAndArgs[0].(string); ok {
			msg = fmt.Sprintf(format, msgAndArgs[1:]...)
		}
	}

	detailedMsg := fmt.Sprintf("%s\n"+
		"期望状态码: %d\n"+
		"实际状态码: %d\n"+
		"响应内容: %s",
		msg, expectedStatus, w.Code, h.formatResponse(w.Body.String()))

	assert.Equal(h.t, expectedStatus, w.Code, detailedMsg)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(h.t, err, "解析响应失败: %s", w.Body.String())

	return response
}

// AssertError 断言请求失败并包含特定错误信息
func (h *TestHelper) AssertError(w *httptest.ResponseRecorder, expectedStatus int, expectedMsg string, msgAndArgs ...interface{}) {
	msg := ""
	if len(msgAndArgs) > 0 {
		if format, ok := msgAndArgs[0].(string); ok {
			msg = fmt.Sprintf(format, msgAndArgs[1:]...)
		}
	}

	detailedMsg := fmt.Sprintf("%s\n"+
		"期望状态码: %d\n"+
		"实际状态码: %d\n"+
		"期望错误信息包含: %s\n"+
		"响应内容: %s",
		msg, expectedStatus, w.Code, expectedMsg, h.formatResponse(w.Body.String()))

	assert.Equal(h.t, expectedStatus, w.Code, detailedMsg)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if message, ok := response["message"].(string); ok {
		assert.Contains(h.t, message, expectedMsg, "错误信息不匹配")
	} else if msg, ok := response["msg"].(string); ok {
		assert.Contains(h.t, msg, expectedMsg, "错误信息不匹配")
	}
}

// formatResponse 格式化响应内容（限制长度）
func (h *TestHelper) formatResponse(body string) string {
	if len(body) > 500 {
		return body[:500] + "...(省略)"
	}
	return body
}

// ========================================
// 数据库辅助函数
// ========================================

// GetTestBook 获取测试书籍
func (h *TestHelper) GetTestBook() string {
	if global.DB == nil {
		h.t.Logf("⚠ global.DB 不可用，跳过获取测试书籍")
		return ""
	}

	var book struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	err := global.DB.Collection("books").FindOne(h.ctx, bson.M{}).Decode(&book)
	if err != nil {
		h.t.Logf("⚠ 数据库中没有测试书籍")
		return ""
	}

	return book.ID.Hex()
}

// GetTestBooks 获取多本测试书籍
func (h *TestHelper) GetTestBooks(limit int) []string {
	if global.DB == nil {
		h.t.Logf("⚠ global.DB 不可用，跳过获取测试书籍列表")
		return nil
	}

	cursor, err := global.DB.Collection("books").Find(h.ctx, bson.M{})
	if err != nil {
		h.t.Logf("⚠ 查询测试书籍失败: %v", err)
		return nil
	}
	defer cursor.Close(h.ctx)

	var books []struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	for cursor.Next(h.ctx) && len(books) < limit {
		var book struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&book); err == nil {
			books = append(books, book)
		}
	}

	bookIDs := make([]string, len(books))
	for i, book := range books {
		bookIDs[i] = book.ID.Hex()
	}

	h.t.Logf("✓ 获取%d本测试书籍", len(bookIDs))
	return bookIDs
}

// CleanupTestData 清理测试数据
func (h *TestHelper) CleanupTestData(collections ...string) {
	if global.DB == nil {
		h.t.Logf("⚠ global.DB 不可用，跳过清理测试数据")
		return
	}

	// 获取当前测试用户的ID
	testUsers := []string{"test_user01", "test_user02", "test_user03"}

	for _, coll := range collections {
		// 先获取测试用户的ObjectID
		var userIDs []string
		cursor, _ := global.DB.Collection("users").Find(h.ctx, bson.M{
			"username": bson.M{"$in": testUsers},
		})
		if cursor != nil {
			var users []bson.M
			cursor.All(h.ctx, &users)
			for _, user := range users {
				if id, ok := user["_id"].(primitive.ObjectID); ok {
					userIDs = append(userIDs, id.Hex())
				}
			}
			cursor.Close(h.ctx)
		}

		// 使用获取到的user_id清理
		if len(userIDs) > 0 {
			_, err := global.DB.Collection(coll).DeleteMany(h.ctx, bson.M{
				"user_id": bson.M{"$in": userIDs},
			})
			if err != nil {
				h.t.Logf("⚠ 清理集合 %s 失败: %v", coll, err)
			} else {
				h.t.Logf("✓ 已清理集合 %s 的测试数据", coll)
			}
		}
	}
}

// RemoveCollectionByBookID 通过API删除指定书籍的收藏（如果存在）
// 这是一个辅助方法，用于测试前清理数据
func (h *TestHelper) RemoveCollectionByBookID(bookID, token string) {
	h.t.Logf("检查并清理书籍收藏: %s", bookID)

	// 1. 获取收藏列表
	w := h.DoAuthRequest("GET", ReaderCollectionsPath, nil, token)

	if w.Code != 200 {
		h.t.Logf("  - 获取收藏列表失败，无法清理")
		return
	}

	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)

	// 2. 在收藏列表中查找指定书籍
	if data, ok := listResp["data"].(map[string]interface{}); ok {
		// 注意：API返回的是list字段，不是collections
		if collections, ok := data["list"].([]interface{}); ok {
			for _, item := range collections {
				if collection, ok := item.(map[string]interface{}); ok {
					// 检查是否是目标书籍
					if collection["book_id"] == bookID {
						// 获取collection的_id或id字段
						var collectionID string
						if id, ok := collection["id"].(string); ok {
							collectionID = id
						} else if id, ok := collection["_id"].(string); ok {
							collectionID = id
						}

						if collectionID != "" {
							// 删除这个收藏
							deleteURL := fmt.Sprintf("%s/%s", ReaderCollectionsPath, collectionID)
							w := h.DoAuthRequest("DELETE", deleteURL, nil, token)

							if w.Code == 200 {
								h.t.Logf("✓ 已删除旧收藏记录: %s (book_id: %s)", collectionID, bookID)
								return
							} else {
								h.t.Logf("⚠ 删除收藏失败 (状态码: %d)", w.Code)
							}
						}
					}
				}
			}
			h.t.Logf("  - 书籍未被收藏，无需清理")
		}
	}
}

// CleanupTestCollections 清理测试用户的收藏数据（针对特定书籍）
// 注意：这个方法直接操作数据库，主要用于defer清理
func (h *TestHelper) CleanupTestCollections(bookID string) {
	if global.DB == nil {
		h.t.Logf("⚠ global.DB 不可用，跳过清理收藏测试数据")
		return
	}

	// 直接删除测试用户的收藏
	testUsers := []string{"test_user01", "test_user02", "test_user03"}

	// 获取测试用户的ID
	var userIDs []string
	cursor, err := global.DB.Collection("users").Find(h.ctx, bson.M{
		"username": bson.M{"$in": testUsers},
	})
	if err != nil {
		return
	}
	if cursor != nil {
		var users []bson.M
		cursor.All(h.ctx, &users)
		for _, user := range users {
			if id, ok := user["_id"].(primitive.ObjectID); ok {
				userIDs = append(userIDs, id.Hex())
			}
		}
		cursor.Close(h.ctx)
	}

	// 删除指定书籍的收藏
	if len(userIDs) > 0 {
		result, _ := global.DB.Collection("collections").DeleteMany(h.ctx, bson.M{
			"user_id": bson.M{"$in": userIDs},
			"book_id": bookID,
		})
		if result.DeletedCount > 0 {
			h.t.Logf("✓ 测试结束：已清理 %d 条收藏数据", result.DeletedCount)
		}
	}
}

// ========================================
// 数据验证辅助函数
// ========================================

// VerifyBookExists 验证书籍是否存在
func (h *TestHelper) VerifyBookExists(bookID string) bool {
	if global.DB == nil {
		return false
	}

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return false
	}

	count, err := global.DB.Collection("books").CountDocuments(h.ctx, bson.M{"_id": objectID})
	return err == nil && count > 0
}

// VerifyUserExists 验证用户是否存在
func (h *TestHelper) VerifyUserExists(userID string) bool {
	if global.DB == nil {
		return false
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false
	}

	count, err := global.DB.Collection("users").CountDocuments(h.ctx, bson.M{"_id": objectID})
	return err == nil && count > 0
}

// ========================================
// 日志辅助函数
// ========================================

// LogSuccess 记录成功日志
func (h *TestHelper) LogSuccess(format string, args ...interface{}) {
	h.t.Logf("✓ "+format, args...)
}

// LogInfo 记录信息日志
func (h *TestHelper) LogInfo(format string, args ...interface{}) {
	h.t.Logf("ℹ "+format, args...)
}

// LogWarning 记录警告日志
func (h *TestHelper) LogWarning(format string, args ...interface{}) {
	h.t.Logf("⚠ "+format, args...)
}

// LogError 记录错误日志
func (h *TestHelper) LogError(format string, args ...interface{}) {
	h.t.Logf("❌ "+format, args...)
}

// ========================================
// 增强的错误诊断功能
// ========================================

// LogRequest 记录详细的请求信息
func (h *TestHelper) LogRequest(method, path string, body interface{}, token string) {
	h.t.Logf("→ 请求: %s %s", method, path)

	if body != nil {
		bodyJSON, err := json.MarshalIndent(body, "", "  ")
		if err == nil {
			h.t.Logf("  Body: %s", bodyJSON)
		} else {
			h.t.Logf("  Body: %v", body)
		}
	}

	if token != "" {
		tokenLen := len(token)
		if tokenLen > 20 {
			h.t.Logf("  Token: %s...", token[:20])
		} else {
			h.t.Logf("  Token: %s", token)
		}
	}
}

// LogResponse 记录详细的响应信息
func (h *TestHelper) LogResponse(w *httptest.ResponseRecorder) {
	h.t.Logf("← 响应: %d", w.Code)

	// 尝试格式化JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, w.Body.Bytes(), "", "  "); err == nil {
		responseStr := prettyJSON.String()
		if len(responseStr) > 500 {
			h.t.Logf("  Body (前500字符): %s...", responseStr[:500])
		} else {
			h.t.Logf("  Body: %s", responseStr)
		}
	} else {
		bodyStr := w.Body.String()
		if len(bodyStr) > 500 {
			h.t.Logf("  Body (前500字符): %s...", bodyStr[:500])
		} else {
			h.t.Logf("  Body: %s", bodyStr)
		}
	}
}

// LogTestContext 记录测试执行上下文
func (h *TestHelper) LogTestContext(step string, details ...interface{}) {
	h.t.Logf("┌─ %s ─┐", step)
	for i, detail := range details {
		h.t.Logf("│ [%d] %v", i+1, detail)
	}
	h.t.Logf("└%s┘", strings.Repeat("─", len(step)+4))
}

// GetResponseString 返回格式化的响应字符串
func (h *TestHelper) GetResponseString(w *httptest.ResponseRecorder) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, w.Body.Bytes(), "", "  "); err == nil {
		return prettyJSON.String()
	}
	return w.Body.String()
}

// AssertJSONResponse 断言响应是有效的JSON并返回解析后的数据
func (h *TestHelper) AssertJSONResponse(w *httptest.ResponseRecorder, msgAndArgs ...interface{}) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	if err != nil {
		msg := "响应不是有效的JSON"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}

		h.t.Logf("❌ %s", msg)
		h.t.Logf("   错误: %v", err)
		h.t.Logf("   响应内容: %s", w.Body.String())
		h.t.FailNow()
	}

	return response
}

// ========================================
// 全局辅助函数（保持向后兼容）
// ========================================

// LoginAsUser 登录指定用户（全局函数）
func LoginAsUser(t *testing.T, router *gin.Engine, username, password string) string {
	helper := NewTestHelper(t, router)
	return helper.LoginUser(username, password)
}

// LoginAsTestUser 登录默认测试用户（全局函数）
func LoginAsTestUser(t *testing.T, router *gin.Engine) string {
	helper := NewTestHelper(t, router)
	return helper.LoginTestUser()
}

// setupTestEnvironment 设置测试环境
func setupTestEnvironment(t *testing.T) (*gin.Engine, func()) {
	// 加载配置
	_, err := config.LoadConfig("../..")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	err = core.InitDB()
	if err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化服务器（会自动初始化服务和路由）
	r, err := core.InitServer()
	if err != nil {
		t.Fatalf("初始化服务器失败: %v", err)
	}

	// 清理函数
	cleanup := func() {
		// 关闭数据库连接
		if global.DB != nil {
			global.DB.Client().Disconnect(context.Background())
			global.DB = nil // 重要：将global.DB设为nil，避免后续测试使用断开的连接
		}
	}

	return r, cleanup
}

// loginTestUser 兼容旧测试的登录函数（待迁移的测试使用）
// 注意：baseURL参数被忽略，因为我们使用TestHelper
func loginTestUser(t *testing.T, baseURL, username, password string) string {
	router, _ := setupTestEnvironment(t)
	helper := NewTestHelper(t, router)
	return helper.LoginUser(username, password)
}
