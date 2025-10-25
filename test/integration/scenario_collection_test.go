package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCollectionScenario 收藏系统场景测试
func TestCollectionScenario(t *testing.T) {
	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 使用测试用户登录
	token := loginAsTestUser(t, router)

	// 测试用书籍ID（使用书城中的书籍）
	testBookID := "507f1f77bcf86cd799439011"

	t.Run("1.收藏管理_添加收藏", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"book_id":   testBookID,
			"note":      "这本书真不错！",
			"tags":      []string{"玄幻", "推荐"},
			"is_public": true,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/reader/collections", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "应该成功添加收藏")

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(201), response["code"])
		assert.NotNil(t, response["data"])

		t.Logf("✓ 添加收藏成功")
	})

	t.Run("2.收藏管理_重复收藏检测", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"book_id": testBookID,
			"note":    "重复收藏",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/reader/collections", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "不应该允许重复收藏")

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["message"], "已经收藏")

		t.Logf("✓ 重复收藏检测通过")
	})

	t.Run("3.收藏管理_检查收藏状态", func(t *testing.T) {
		url := fmt.Sprintf("/api/v1/reader/collections/check/%s", testBookID)
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.True(t, data["is_collected"].(bool), "应该显示已收藏")

		t.Logf("✓ 收藏状态检查通过")
	})

	t.Run("4.收藏管理_获取收藏列表", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/collections?page=1&size=20", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		list := data["list"].([]interface{})
		assert.Greater(t, len(list), 0, "应该有至少一条收藏记录")

		t.Logf("✓ 获取收藏列表成功，共%d条", len(list))
	})

	t.Run("5.收藏夹管理_创建收藏夹", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "我的最爱",
			"description": "收藏的经典作品",
			"is_public":   true,
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/reader/collections/folders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "应该成功创建收藏夹")

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(201), response["code"])
		assert.NotNil(t, response["data"])

		t.Logf("✓ 创建收藏夹成功")
	})

	t.Run("6.收藏夹管理_获取收藏夹列表", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/collections/folders", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		list := data["list"].([]interface{})
		assert.Greater(t, len(list), 0, "应该有至少一个收藏夹")

		t.Logf("✓ 获取收藏夹列表成功，共%d个", len(list))
	})

	t.Run("7.收藏分享_获取公开收藏", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/collections/public?page=1&size=10", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(200), response["code"])

		t.Logf("✓ 获取公开收藏列表成功")
	})

	t.Run("8.收藏统计_获取统计数据", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/reader/collections/stats", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.NotNil(t, data["total_collections"])
		assert.NotNil(t, data["total_folders"])

		t.Logf("✓ 收藏统计: %v条收藏, %v个收藏夹",
			data["total_collections"], data["total_folders"])
	})
}
