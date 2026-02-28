package integration

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCollectionScenario 收藏系统场景测试
func TestCollectionScenario(t *testing.T) {
	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 创建TestHelper
	helper := NewTestHelper(t, router)

	// 使用动态测试用户登录，避免依赖固定种子账号
	username := fmt.Sprintf("it_collection_%d", time.Now().UnixNano())
	password := "Test@123456"
	registerData := map[string]interface{}{
		"username": username,
		"email":    fmt.Sprintf("%s@test.com", username),
		"password": password,
	}
	registerResp := helper.DoRequest("POST", RegisterPath, registerData, "")
	if registerResp.Code != 200 && registerResp.Code != 201 {
		t.Fatalf("创建测试用户失败，状态码: %d, 响应: %s", registerResp.Code, registerResp.Body.String())
	}

	loginResp := helper.DoRequest("POST", LoginPath, map[string]interface{}{
		"username": username,
		"password": password,
	}, "")
	if loginResp.Code != 200 {
		t.Fatalf("测试用户登录失败，状态码: %d, 响应: %s", loginResp.Code, loginResp.Body.String())
	}

	var loginBody map[string]interface{}
	_ = json.Unmarshal(loginResp.Body.Bytes(), &loginBody)
	data, _ := loginBody["data"].(map[string]interface{})
	token, _ := data["token"].(string)
	if token == "" {
		t.Fatal("登录失败，无法继续测试")
	}

	// 测试用书籍ID（使用书城中的书籍）
	testBookID := "507f1f77bcf86cd799439011"

	// 清理可能的遗留测试数据（使用API方式）
	helper.RemoveCollectionByBookID(testBookID, token)

	// 测试结束后清理数据
	defer helper.CleanupTestData("collections", "collection_folders")

	t.Run("1.收藏管理_添加收藏", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"book_id":   testBookID,
			"note":      "这本书真不错！",
			"tags":      []string{"玄幻", "推荐"},
			"is_public": true,
		}

		w := helper.DoAuthRequest("POST", ReaderCollectionsPath, reqBody, token)
		response := helper.AssertSuccess(w, 201, "添加收藏失败")

		assert.EqualValues(t, 0, response["code"])
		assert.NotNil(t, response["data"])

		helper.LogSuccess("添加收藏成功")
	})

	t.Run("2.收藏管理_重复收藏检测", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"book_id": testBookID,
			"note":    "重复收藏",
		}

		w := helper.DoAuthRequest("POST", ReaderCollectionsPath, reqBody, token)
		helper.AssertError(w, 500, "已经收藏过该书籍", "重复收藏检测失败")

		helper.LogSuccess("重复收藏检测通过")
	})

	t.Run("3.收藏管理_检查收藏状态", func(t *testing.T) {
		url := fmt.Sprintf("%s/check/%s", ReaderCollectionsPath, testBookID)
		w := helper.DoAuthRequest("GET", url, nil, token)
		response := helper.AssertSuccess(w, 200, "检查收藏状态失败")

		data := response["data"].(map[string]interface{})
		assert.True(t, data["is_collected"].(bool), "应该显示已收藏")

		helper.LogSuccess("收藏状态检查通过")
	})

	t.Run("4.收藏管理_获取收藏列表", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderCollectionsPath+"?page=1&size=20", nil, token)
		response := helper.AssertSuccess(w, 200, "获取收藏列表失败")

		data := response["data"].(map[string]interface{})
		list := data["list"].([]interface{})
		assert.Greater(t, len(list), 0, "应该有至少一条收藏记录")

		helper.LogSuccess("获取收藏列表成功，共%d条", len(list))
	})

	t.Run("5.收藏夹管理_创建收藏夹", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "我的最爱",
			"description": "收藏的经典作品",
			"is_public":   true,
		}

		w := helper.DoAuthRequest("POST", ReaderCollectionsPath+"/folders", reqBody, token)
		response := helper.AssertSuccess(w, 201, "创建收藏夹失败")

		assert.EqualValues(t, 0, response["code"])
		assert.NotNil(t, response["data"])

		helper.LogSuccess("创建收藏夹成功")
	})

	t.Run("6.收藏夹管理_获取收藏夹列表", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderCollectionsPath+"/folders", nil, token)
		response := helper.AssertSuccess(w, 200, "获取收藏夹列表失败")

		data := response["data"].(map[string]interface{})
		list := data["list"].([]interface{})
		assert.Greater(t, len(list), 0, "应该有至少一个收藏夹")

		helper.LogSuccess("获取收藏夹列表成功，共%d个", len(list))
	})

	t.Run("7.收藏分享_获取公开收藏", func(t *testing.T) {
		// 公开接口不需要认证
		w := helper.DoRequest("GET", ReaderCollectionsPath+"/public?page=1&size=10", nil, "")
		response := helper.AssertSuccess(w, 200, "获取公开收藏列表失败")

		assert.EqualValues(t, 0, response["code"])

		helper.LogSuccess("获取公开收藏列表成功")
	})

	t.Run("8.收藏统计_获取统计数据", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderCollectionsPath+"/stats", nil, token)
		response := helper.AssertSuccess(w, 200, "获取收藏统计失败")

		data := response["data"].(map[string]interface{})
		assert.NotNil(t, data["total_collections"])
		assert.NotNil(t, data["total_folders"])

		helper.LogSuccess("收藏统计: %v条收藏, %v个收藏夹",
			data["total_collections"], data["total_folders"])
	})
}
