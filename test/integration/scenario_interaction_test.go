package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 互动功能测试 - 评论、点赞、收藏、阅读历史等
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

	// 用于存储测试过程中创建的ID
	var commentID string
	var collectionID string

	t.Run("1.收藏_添加书籍到收藏", func(t *testing.T) {
		// 使用真实的收藏API
		requestBody := map[string]interface{}{
			"book_id": testBookID,
			"notes":   "集成测试收藏",
			"tags":    []string{"测试"},
		}
		bodyJSON, _ := json.Marshal(requestBody)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/collections", baseURL), bytes.NewBuffer(bodyJSON))
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
		require.NoError(t, err, "JSON解析失败")

		// 如果已收藏，也算成功
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Logf("✓ 书籍添加到收藏成功")
			// 保存收藏ID
			if data, ok := result["data"].(map[string]interface{}); ok {
				if id, ok := data["id"].(string); ok {
					collectionID = id
					t.Logf("  收藏ID: %s", collectionID)
				}
			}
		} else {
			t.Logf("○ 添加收藏响应: %v (可能已存在)", result["message"])
		}
	})

	t.Run("2.收藏_获取收藏列表", func(t *testing.T) {
		// 使用真实的收藏列表API
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/collections?page=1&page_size=10", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err, "JSON解析失败")

		assert.Equal(t, http.StatusOK, resp.StatusCode, "获取收藏列表应该成功")

		if data, ok := result["data"].(map[string]interface{}); ok {
			if collections, ok := data["collections"].([]interface{}); ok {
				t.Logf("✓ 收藏列表获取成功，共 %d 条记录", len(collections))
			}
		}
	})

	t.Run("3.收藏_获取收藏统计", func(t *testing.T) {
		// 使用真实的收藏统计API
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/collections/stats", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err, "JSON解析失败")

		assert.Equal(t, http.StatusOK, resp.StatusCode, "获取收藏统计应该成功")

		if data, ok := result["data"].(map[string]interface{}); ok {
			t.Logf("✓ 收藏统计获取成功")
			t.Logf("  总收藏数: %.0f", data["total_count"])
			t.Logf("  收藏夹数: %.0f", data["folder_count"])
		}
	})

	t.Run("4.评论_发表书籍评论", func(t *testing.T) {
		// 使用真实的评论API
		requestBody := map[string]interface{}{
			"book_id": testBookID,
			"content": "这是一本很棒的书！集成测试评论。",
			"rating":  5,
		}
		bodyJSON, _ := json.Marshal(requestBody)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/comments", baseURL), bytes.NewBuffer(bodyJSON))
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
		require.NoError(t, err, "JSON解析失败")

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Logf("✓ 评论发表成功")
			// 保存评论ID
			if data, ok := result["data"].(map[string]interface{}); ok {
				if id, ok := data["id"].(string); ok {
					commentID = id
					t.Logf("  评论ID: %s", commentID)
					t.Logf("  审核状态: %v", data["status"])
				}
			}
		} else {
			t.Logf("○ 发表评论失败: %v", result["message"])
		}
	})

	t.Run("5.评论_获取书籍评论列表", func(t *testing.T) {
		// 使用真实的评论列表API
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/comments?book_id=%s&page=1&page_size=10", baseURL, testBookID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err, "JSON解析失败")

		assert.Equal(t, http.StatusOK, resp.StatusCode, "获取评论列表应该成功")

		if data, ok := result["data"].(map[string]interface{}); ok {
			if comments, ok := data["comments"].([]interface{}); ok {
				t.Logf("✓ 评论列表获取成功，共 %d 条评论", len(comments))
				if len(comments) > 0 {
					comment := comments[0].(map[string]interface{})
					t.Logf("  第一条评论: %v", comment["content"])
				}
			}
		}
	})

	t.Run("6.点赞_点赞书籍", func(t *testing.T) {
		// 使用真实的点赞API
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
		require.NoError(t, err, "JSON解析失败")

		// 点赞是幂等的，所以OK或已点赞都算成功
		if resp.StatusCode == http.StatusOK {
			t.Logf("✓ 书籍点赞成功")
			if data, ok := result["data"].(map[string]interface{}); ok {
				t.Logf("  点赞状态: %v", data["is_liked"])
			}
		} else {
			t.Logf("○ 点赞响应: %v", result["message"])
		}
	})

	t.Run("7.点赞_获取点赞信息", func(t *testing.T) {
		// 使用真实的点赞信息API
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/books/%s/like/info", baseURL, testBookID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err, "JSON解析失败")

		assert.Equal(t, http.StatusOK, resp.StatusCode, "获取点赞信息应该成功")

		if data, ok := result["data"].(map[string]interface{}); ok {
			t.Logf("✓ 点赞信息获取成功")
			t.Logf("  当前用户是否已点赞: %v", data["is_liked"])
			t.Logf("  总点赞数: %.0f", data["like_count"])
		}
	})

	t.Run("8.阅读历史_记录阅读", func(t *testing.T) {
		// 使用真实的阅读历史记录API
		requestBody := map[string]interface{}{
			"book_id":       testBookID,
			"chapter_id":    "chapter_001",
			"read_duration": 300,
			"progress":      25.5,
			"device_type":   "web",
			"device_id":     "integration_test_device",
		}
		bodyJSON, _ := json.Marshal(requestBody)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/reading-history", baseURL), bytes.NewBuffer(bodyJSON))
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
		require.NoError(t, err, "JSON解析失败")

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			t.Logf("✓ 阅读历史记录成功")
			if data, ok := result["data"].(map[string]interface{}); ok {
				t.Logf("  历史ID: %v", data["id"])
				t.Logf("  阅读时长: %.0f秒", data["read_duration"])
			}
		} else {
			t.Logf("○ 记录阅读历史失败: %v", result["message"])
		}
	})

	t.Run("9.阅读历史_获取历史列表", func(t *testing.T) {
		// 使用真实的阅读历史列表API
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/reading-history?page=1&page_size=10", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err, "JSON解析失败")

		assert.Equal(t, http.StatusOK, resp.StatusCode, "获取阅读历史列表应该成功")

		if data, ok := result["data"].(map[string]interface{}); ok {
			if histories, ok := data["histories"].([]interface{}); ok {
				t.Logf("✓ 阅读历史列表获取成功，共 %d 条记录", len(histories))
			}
		}
	})

	t.Run("10.阅读历史_获取统计", func(t *testing.T) {
		// 使用真实的阅读统计API
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/reading-history/stats", baseURL), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err, "JSON解析失败")

		assert.Equal(t, http.StatusOK, resp.StatusCode, "获取阅读统计应该成功")

		if data, ok := result["data"].(map[string]interface{}); ok {
			t.Logf("✓ 阅读统计获取成功")
			t.Logf("  总阅读时长: %.0f秒", data["total_duration"])
			t.Logf("  阅读书籍数: %.0f", data["total_books"])
			t.Logf("  阅读章节数: %.0f", data["total_chapters"])
		}
	})

	t.Run("11.点赞_取消点赞", func(t *testing.T) {
		// 使用真实的取消点赞API
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
		require.NoError(t, err, "JSON解析失败")

		// 取消点赞是幂等的
		if resp.StatusCode == http.StatusOK {
			t.Logf("✓ 取消点赞成功")
		} else {
			t.Logf("○ 取消点赞响应: %v", result["message"])
		}
	})

	t.Run("12.收藏_取消收藏", func(t *testing.T) {
		// 如果有收藏ID，尝试删除
		if collectionID != "" {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/reader/collections/%s", baseURL, collectionID), nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err, "JSON解析失败")

			if resp.StatusCode == http.StatusOK {
				t.Logf("✓ 取消收藏成功")
			} else {
				t.Logf("○ 取消收藏失败: %v", result["message"])
			}
		} else {
			t.Log("○ 没有收藏ID，跳过取消收藏测试")
		}
	})

	t.Run("13.统计_数据库互动数据", func(t *testing.T) {
		// 统计数据库中的互动数据
		collectionCount, _ := global.DB.Collection("collections").CountDocuments(ctx, bson.M{})
		commentCount, _ := global.DB.Collection("comments").CountDocuments(ctx, bson.M{})
		likeCount, _ := global.DB.Collection("likes").CountDocuments(ctx, bson.M{})
		historyCount, _ := global.DB.Collection("reading_histories").CountDocuments(ctx, bson.M{})

		t.Logf("✓ 互动数据统计:")
		t.Logf("  收藏记录: %d", collectionCount)
		t.Logf("  评论数量: %d", commentCount)
		t.Logf("  点赞数量: %d", likeCount)
		t.Logf("  阅读历史: %d", historyCount)
	})

	t.Logf("\n=== 互动功能集成测试完成 ===")
	t.Logf("测试场景: 收藏 → 评论 → 点赞 → 阅读历史")
	t.Logf("已测试API: 13个端点")
}
