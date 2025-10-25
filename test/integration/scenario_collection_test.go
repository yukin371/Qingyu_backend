package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"
)

// TestCollectionScenario 收藏系统场景测试
func TestCollectionScenario(t *testing.T) {
	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 创建TestHelper
	helper := NewTestHelper(t, router)

	// 使用测试用户登录
	token := helper.LoginTestUser()
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

		assert.Equal(t, float64(201), response["code"])
		assert.NotNil(t, response["data"])

		helper.LogSuccess("添加收藏成功")
	})

	t.Run("2.收藏管理_重复收藏检测", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"book_id": testBookID,
			"note":    "重复收藏",
		}

		w := helper.DoAuthRequest("POST", ReaderCollectionsPath, reqBody, token)
		helper.AssertError(w, 400, "该书籍已经收藏", "重复收藏检测失败")

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

		assert.Equal(t, float64(201), response["code"])
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

		assert.Equal(t, float64(200), response["code"])

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
		}
	}

	return r, cleanup
}
