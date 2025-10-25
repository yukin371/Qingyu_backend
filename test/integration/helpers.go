package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	UserProfilePath  = APIBasePath + "/users/profile"
	UserPasswordPath = APIBasePath + "/users/password"

	// 阅读器相关
	ReaderBooksPath       = APIBasePath + "/reader/books"
	ReaderChaptersPath    = APIBasePath + "/reader/chapters"
	ReaderProgressPath    = APIBasePath + "/reader/progress"
	ReaderAnnotationsPath = APIBasePath + "/reader/annotations"
	ReaderCommentsPath    = APIBasePath + "/reader/comments"
	ReaderCollectionsPath = APIBasePath + "/reader/collections"
	ReaderLikesPath       = APIBasePath + "/reader/likes"

	// 书城相关
	BookstoreHomePath    = APIBasePath + "/bookstore/homepage"
	BookstoreBooksPath   = APIBasePath + "/bookstore/books"
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
	for _, coll := range collections {
		_, err := global.DB.Collection(coll).DeleteMany(h.ctx, bson.M{
			"user_id": bson.M{"$regex": "^test_"},
		})
		if err != nil {
			h.t.Logf("⚠ 清理集合 %s 失败: %v", coll, err)
		}
	}
}

// ========================================
// 数据验证辅助函数
// ========================================

// VerifyBookExists 验证书籍是否存在
func (h *TestHelper) VerifyBookExists(bookID string) bool {
	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return false
	}

	count, err := global.DB.Collection("books").CountDocuments(h.ctx, bson.M{"_id": objectID})
	return err == nil && count > 0
}

// VerifyUserExists 验证用户是否存在
func (h *TestHelper) VerifyUserExists(userID string) bool {
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
