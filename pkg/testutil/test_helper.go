package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHelper 测试助手
type TestHelper struct {
	router    *gin.Engine
	t         *testing.T
	requestID string
}

// NewTestHelper 创建测试助手
func NewTestHelper(t *testing.T, router *gin.Engine) *TestHelper {
	return &TestHelper{
		router:    router,
		t:         t,
		requestID: "test_req_123",
	}
}

// PerformRequest 执行测试请求
func (h *TestHelper) PerformRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var bodyBuf *bytes.Buffer
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		assert.NoError(h.t, err, "请求体序列化失败")
		bodyBuf = bytes.NewBuffer(bodyBytes)
	} else {
		bodyBuf = bytes.NewBuffer([]byte("{}"))
	}

	req, err := http.NewRequest(method, path, bodyBuf)
	assert.NoError(h.t, err, "创建HTTP请求失败")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", h.requestID)

	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)
	return w
}

// ParseResponse 解析响应
func (h *TestHelper) ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(h.t, err, "响应体反序列化失败")
	return resp
}
