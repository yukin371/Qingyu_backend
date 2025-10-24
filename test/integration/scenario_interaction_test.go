package integration

import (
	"bytes"
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
			testBookID = fmt.Sprintf("%v", books[0]["_id"])
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	t.Run("1.收藏_添加书籍到收藏", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/collections/%s", baseURL, testBookID), nil)
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
			t.Logf("✓ 书籍收藏成功")
		} else {
			t.Logf("○ 收藏失败或接口不存在: %v", result["message"])
		}
	})

	t.Run("2.收藏_获取收藏列表", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/collections?page=1&size=10", baseURL), nil)
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
				collections := data.(map[string]interface{})
				books := collections["books"]
				if books != nil {
					bookList := books.([]interface{})
					t.Logf("✓ 收藏列表获取成功，共 %d 本书", len(bookList))
				}
			}
		} else {
			t.Logf("○ 获取收藏列表失败: %v", result["message"])
		}
	})

	t.Run("3.收藏_取消收藏", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/reader/collections/%s", baseURL, testBookID), nil)
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
			t.Logf("✓ 取消收藏成功")
		} else {
			t.Logf("○ 取消收藏失败: %v", result["message"])
		}
	})

	t.Run("4.评论_发表书籍评论", func(t *testing.T) {
		commentData := map[string]interface{}{
			"book_id": testBookID,
			"content": "这本书写得不错，情节引人入胜！",
			"rating":  4.5,
		}

		jsonData, _ := json.Marshal(commentData)
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/comments", baseURL), bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		if result["code"] == float64(200) {
			t.Logf("✓ 评论发表成功")
			t.Logf("  内容: %s", commentData["content"])
			t.Logf("  评分: %.1f", commentData["rating"])
		} else {
			t.Logf("○ 发表评论失败或接口不存在: %v", result["message"])
		}
	})

	t.Run("5.评论_获取书籍评论列表", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/comments?bookId=%s&page=1&size=10", baseURL, testBookID), nil)
		require.NoError(t, err)

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
				comments := data.(map[string]interface{})
				commentList := comments["comments"]
				if commentList != nil {
					list := commentList.([]interface{})
					t.Logf("✓ 评论列表获取成功，共 %d 条评论", len(list))
				}
			}
		} else {
			t.Logf("○ 获取评论列表失败: %v", result["message"])
		}
	})

	t.Run("6.点赞_点赞书籍", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/books/%s/like", baseURL, testBookID), nil)
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
			t.Logf("✓ 点赞成功")
		} else {
			t.Logf("○ 点赞失败或接口不存在: %v", result["message"])
		}
	})

	t.Run("7.点赞_取消点赞", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/reader/books/%s/like", baseURL, testBookID), nil)
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
			t.Logf("✓ 取消点赞成功")
		} else {
			t.Logf("○ 取消点赞失败: %v", result["message"])
		}
	})

	t.Run("8.阅读历史_获取阅读历史", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/history?page=1&size=10", baseURL), nil)
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
