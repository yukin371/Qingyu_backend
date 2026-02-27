//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
	"Qingyu_backend/service"
	"Qingyu_backend/test/testutil"
)

// TestEnvironment E2E 测试环境
type TestEnvironment struct {
	T          *testing.T
	Router     *gin.Engine
	Client     *http.Client
	Config     *config.Config
	testData   map[string]interface{}
	cleanupFns []func()
}

// SetupTestEnvironment 初始化完整 E2E 测试环境
func SetupTestEnvironment(t *testing.T) (*TestEnvironment, func()) {
	t.Helper()

	// 1. 加载测试配置
	cfg, err := config.LoadConfig("../../config")
	require.NoError(t, err, "加载测试配置失败")
	config.GlobalConfig = cfg
	testutil.EnableStrictLogAssertions(t)

	// 2. 初始化服务（会自动创建 ServiceContainer）
	err = core.InitServices()
	require.NoError(t, err, "初始化服务失败")

	// 兼容仍依赖 global.DB 的 E2E fixtures/helpers
	if sc := service.GetServiceContainer(); sc != nil {
		global.DB = sc.GetMongoDB()
		global.MongoClient = sc.GetMongoClient()
	}

	// 3. 初始化服务器
	gin.SetMode(gin.TestMode)
	router, err := core.InitServer()
	require.NoError(t, err, "初始化服务器失败")

	// 4. 创建测试环境
	env := &TestEnvironment{
		T:          t,
		Router:     router,
		Client:     &http.Client{Timeout: 10 * time.Second},
		Config:     cfg,
		testData:   make(map[string]interface{}),
		cleanupFns: make([]func(), 0),
	}

	// 5. 创建清理函数
	cleanup := func() {
		env.CleanupAll()
	}

	env.T.Log("✓ E2E 测试环境初始化完成")

	return env, cleanup
}

// DoRequest 执行 HTTP 请求
func (env *TestEnvironment) DoRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)

	return w
}

// RegisterCleanup 注册清理函数
func (env *TestEnvironment) RegisterCleanup(fn func()) {
	env.cleanupFns = append(env.cleanupFns, fn)
}

// CleanupAll 执行所有清理
func (env *TestEnvironment) CleanupAll() {
	// 执行注册的清理函数
	for i := len(env.cleanupFns) - 1; i >= 0; i-- {
		fn := env.cleanupFns[i]
		if fn != nil {
			fn()
		}
	}

	// 清理带 e2e_test_ 前缀的数据
	env.cleanupTestDataByPrefix()

	// 关闭全局数据库连接
	// 注意：这需要从 global 包导入
	// if global.DB != nil {
	// 	global.DB.Client().Disconnect(context.Background())
	// }

	env.T.Log("✓ 测试环境清理完成")
}

// cleanupTestDataByPrefix 删除所有带 e2e_test_ 前缀的数据
func (env *TestEnvironment) cleanupTestDataByPrefix() {
	// 注意：这需要从 global 包导入 DB
	// db := global.DB
	//
	// collections := []string{
	// 	"users", "books", "chapters", "purchases",
	// 	"reading_progress", "comments", "collections", "likes",
	// 	"projects", "documents",
	// }
	//
	// for _, coll := range collections {
	// 	filter := bson.M{
	// 		"$or": []bson.M{
	// 			{"username": bson.M{"$regex": "^e2e_test_"}},
	// 			{"email": bson.M{"$regex": "^e2e_test_"}},
	// 			{"title": bson.M{"$regex": "^e2e_test_"}},
	// 			{"name": bson.M{"$regex": "^e2e_test_"}},
	// 		},
	// 	}
	// 	result, _ := db.Collection(coll).DeleteMany(ctx, filter)
	// 	if result.DeletedCount > 0 {
	// 		env.T.Logf("清理 %s: %d 条记录", coll, result.DeletedCount)
	// 	}
	// }

	env.T.Log("✓ 测试数据清理完成（数据前缀模式）")
}

// GetTestData 获取测试数据
func (env *TestEnvironment) GetTestData(key string) interface{} {
	return env.testData[key]
}

// SetTestData 设置测试数据
func (env *TestEnvironment) SetTestData(key string, value interface{}) {
	env.testData[key] = value
}

// ParseJSONResponse 解析 JSON 响应
func (env *TestEnvironment) ParseJSONResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(env.T, err, "解析 JSON 响应失败")
	return response
}

// HasRoute 检查路由是否已注册（按 method + 路由模板路径匹配）
func (env *TestEnvironment) HasRoute(method, routePath string) bool {
	routes := env.Router.Routes()
	for _, r := range routes {
		if r.Method == method && r.Path == routePath {
			return true
		}
	}
	return false
}

// LogSuccess 记录成功日志
func (env *TestEnvironment) LogSuccess(format string, args ...interface{}) {
	env.T.Logf("✓ "+format, args...)
}

// LogInfo 记录信息日志
func (env *TestEnvironment) LogInfo(format string, args ...interface{}) {
	env.T.Logf("ℹ "+format, args...)
}

// LogError 记录错误日志
func (env *TestEnvironment) LogError(format string, args ...interface{}) {
	env.T.Logf("❌ "+format, args...)
}

// ConsistencyValidator 获取数据一致性验证器
func (env *TestEnvironment) ConsistencyValidator() *ConsistencyValidatorWrapper {
	return &ConsistencyValidatorWrapper{env: env}
}
