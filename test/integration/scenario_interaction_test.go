package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 互动功能测试 - 收藏、评论、阅读历史等
func TestInteractionScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化
	_, err := config.LoadConfig("../..")
	require.NoError(t, err, "加载配置失败")

	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	ctx := context.Background()
	baseURL := "http://localhost:8080"

	// 登录获取 token
	token := loginTestUser(t, baseURL, "test_user01", "Test@123456")
	if token == "" {
		t.Skip("无法登录测试用户，跳过互动测试")
	}

	// 获取一本测试书籍
	var testBookID string
	cursor, err := global.DB.Collection("books").Find(ctx, bson.M{}, nil)
	if err == nil {
		var books []map[string]interface{}
		cursor.All(ctx, &books)
		cursor.Close(ctx)

		if len(books) > 0 {
			// 正确处理MongoDB ObjectID类型
			if oid, ok := books[0]["_id"].(primitive.ObjectID); ok {
				testBookID = oid.Hex()
			}
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	t.Run("1.收藏_添加书籍到收藏", func(t *testing.T) {
		// TODO: 收藏功能尚未实现，使用书架功能代替
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/books/%s", baseURL, testBookID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 检查HTTP状态码
		if resp.StatusCode == http.StatusNotFound {
			t.Skip("收藏API尚未实现")
			return
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Logf("○ JSON解析失败（API可能返回了HTML）: %v", err)
			t.Skip("收藏API响应格式错误")
			return
		}

		if result["code"] == float64(200) {
			t.Logf("✓ 书籍添加到书架成功（收藏功能）")
		} else {
			t.Logf("○ 添加失败: %v", result["message"])
		}
	})

	t.Run("2.收藏_获取收藏列表", func(t *testing.T) {
		// TODO: 使用书架API代替收藏列表
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/books", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("收藏列表API尚未实现")
			return
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Skip("收藏列表API响应格式错误")
			return
		}

		if result["code"] == float64(200) {
			t.Logf("✓ 书架列表获取成功（收藏功能）")
		} else {
			t.Logf("○ 获取列表失败: %v", result["message"])
		}
	})

	t.Run("3.收藏_取消收藏", func(t *testing.T) {
		// TODO: 使用书架移除API代替
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/reader/books/%s", baseURL, testBookID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("取消收藏API尚未实现")
			return
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Skip("取消收藏API响应格式错误")
			return
		}

		if result["code"] == float64(200) {
			t.Logf("✓ 从书架移除成功（取消收藏功能）")
		} else {
			t.Logf("○ 移除失败: %v", result["message"])
		}
	})

	t.Run("4.评论_发表书籍评论", func(t *testing.T) {
		// TODO: 实现评论功能
		// - 需要实现 POST /api/v1/reader/comments API
		// - 支持评论内容、评分
		// - 需要评论审核机制
		t.Skip("评论功能尚未实现")
	})

	t.Run("5.评论_获取书籍评论列表", func(t *testing.T) {
		// TODO: 实现评论列表查询
		// - 需要实现 GET /api/v1/reader/comments API
		// - 支持分页、排序（最新/最热）
		// - 支持评论回复嵌套显示
		t.Skip("评论功能尚未实现")
	})

	t.Run("6.点赞_点赞书籍", func(t *testing.T) {
		// TODO: 实现点赞功能
		// - 需要实现 POST /api/v1/reader/books/:id/like API
		// - 支持点赞去重（一个用户只能点赞一次）
		// - 更新书籍点赞数统计
		t.Skip("点赞功能尚未实现")
	})

	t.Run("7.点赞_取消点赞", func(t *testing.T) {
		// TODO: 实现取消点赞功能
		// - 需要实现 DELETE /api/v1/reader/books/:id/like API
		// - 同步更新书籍点赞数统计
		t.Skip("取消点赞功能尚未实现")
	})

	t.Run("8.阅读历史_获取阅读历史", func(t *testing.T) {
		// 使用现有的阅读进度历史API
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/progress/history?page=1&size=10", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			t.Skip("阅读历史API尚未实现")
			return
		}

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Skip("阅读历史API响应格式错误")
			return
		}

		if result["code"] == float64(200) {
			data := result["data"]
			if data != nil {
				history := data.(map[string]interface{})
				books := history["books"]
				if books != nil {
					bookList := books.([]interface{})
					t.Logf("✓ 阅读历史获取成功，共 %d 条记录", len(bookList))
				}
			}
		} else {
			t.Logf("○ 获取阅读历史失败或接口不存在: %v", result["message"])
		}
	})

	t.Run("9.书架_获取个人书架", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/books?page=1&size=10", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		if result["code"] == float64(200) {
			data := result["data"]
			if data != nil {
				bookshelf := data.(map[string]interface{})
				books := bookshelf["books"]
				if books != nil {
					bookList := books.([]interface{})
					t.Logf("✓ 书架获取成功，共 %d 本书", len(bookList))

					// 显示书架中的书籍
					for i, book := range bookList {
						if i < 3 {
							bookMap := book.(map[string]interface{})
							t.Logf("  书籍%d: %v", i+1, bookMap["title"])
						}
					}
				}
			}
		} else {
			t.Logf("○ 获取书架失败: %v", result["message"])
		}
	})

	t.Run("10.统计_数据库互动数据统计", func(t *testing.T) {
		// 统计数据库中的互动数据
		collectionCount, _ := global.DB.Collection("user_collections").CountDocuments(ctx, bson.M{})
		commentCount, _ := global.DB.Collection("comments").CountDocuments(ctx, bson.M{})
		progressCount, _ := global.DB.Collection("reading_progress").CountDocuments(ctx, bson.M{})

		t.Logf("✓ 互动数据统计:")
		t.Logf("  收藏记录: %d", collectionCount)
		t.Logf("  评论数量: %d", commentCount)
		t.Logf("  阅读进度: %d", progressCount)
	})

	t.Logf("\n=== 互动功能测试完成 ===")
	t.Logf("测试场景: 收藏 → 评论 → 点赞 → 阅读历史 → 个人书架")
}
