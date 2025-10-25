package integration

import (
	"context"
	"fmt"
	"testing"

	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 互动功能测试 - 评论、点赞、收藏、阅读历史等
func TestInteractionScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 初始化helper
	helper := NewTestHelper(t, router)

	// 登录测试用户
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过互动测试")
	}

	ctx := context.Background()

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
		// 先清理可能存在的旧收藏
		helper.RemoveCollectionByBookID(testBookID, token)

		// 添加收藏
		requestBody := map[string]interface{}{
			"book_id": testBookID,
			"notes":   "集成测试收藏",
			"tags":    []string{"测试"},
		}

		w := helper.DoAuthRequest("POST", ReaderCollectionsPath, requestBody, token)
		data := helper.AssertSuccess(w, 201, "添加收藏应该成功")

		// 保存收藏ID - 根据API响应结构获取
		if id, ok := data["_id"].(string); ok {
			collectionID = id
		} else if id, ok := data["id"].(string); ok {
			collectionID = id
		}
		helper.LogSuccess("书籍添加到收藏成功，收藏ID: %s", collectionID)
	})

	t.Run("2.收藏_获取收藏列表", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderCollectionsPath+"?page=1&page_size=10", nil, token)
		data := helper.AssertSuccess(w, 200, "获取收藏列表应该成功")

		if list, ok := data["list"].([]interface{}); ok {
			helper.LogSuccess("收藏列表获取成功，共 %d 条记录", len(list))
		}
	})

	t.Run("3.收藏_获取收藏统计", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderCollectionsPath+"/stats", nil, token)
		data := helper.AssertSuccess(w, 200, "获取收藏统计应该成功")

		totalCount := data["total_count"]
		folderCount := data["folder_count"]
		helper.LogSuccess("收藏统计获取成功 - 总收藏数: %v, 收藏夹数: %v", totalCount, folderCount)
	})

	t.Run("4.评论_发表书籍评论", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"book_id": testBookID,
			"content": "这是一本很棒的书！集成测试评论。",
			"rating":  5,
		}

		w := helper.DoAuthRequest("POST", ReaderCommentsPath, requestBody, token)
		data := helper.AssertSuccess(w, 201, "评论发表应该成功")

		// 保存评论ID
		if id, ok := data["_id"].(string); ok {
			commentID = id
		} else if id, ok := data["id"].(string); ok {
			commentID = id
		}
		status := data["status"]
		helper.LogSuccess("评论发表成功 - 评论ID: %s, 审核状态: %v", commentID, status)
	})

	t.Run("5.评论_获取书籍评论列表", func(t *testing.T) {
		// 注意：API需要book_id作为必需参数
		url := fmt.Sprintf("%s?book_id=%s&page=1&page_size=10", ReaderCommentsPath, testBookID)
		w := helper.DoAuthRequest("GET", url, nil, token)
		data := helper.AssertSuccess(w, 200, "获取评论列表应该成功")

		if comments, ok := data["comments"].([]interface{}); ok {
			helper.LogSuccess("评论列表获取成功，共 %d 条评论", len(comments))
			if len(comments) > 0 {
				comment := comments[0].(map[string]interface{})
				t.Logf("  第一条评论: %v", comment["content"])
			}
		}
	})

	t.Run("6.点赞_点赞书籍", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s/like", ReaderBooksPath, testBookID)
		w := helper.DoAuthRequest("POST", url, nil, token)
		data := helper.AssertSuccess(w, 200, "书籍点赞应该成功")

		isLiked := data["is_liked"]
		helper.LogSuccess("书籍点赞成功，点赞状态: %v", isLiked)
	})

	t.Run("7.点赞_获取点赞信息", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s/like/info", ReaderBooksPath, testBookID)
		w := helper.DoAuthRequest("GET", url, nil, token)
		data := helper.AssertSuccess(w, 200, "获取点赞信息应该成功")

		isLiked := data["is_liked"]
		likeCount := data["like_count"]
		helper.LogSuccess("点赞信息获取成功 - 当前用户是否已点赞: %v, 总点赞数: %v", isLiked, likeCount)
	})

	t.Run("8.阅读历史_记录阅读", func(t *testing.T) {
		// 注意：添加必需的StartTime和EndTime字段
		requestBody := map[string]interface{}{
			"book_id":       testBookID,
			"chapter_id":    "chapter_001",
			"start_time":    "2024-01-01T10:00:00Z",
			"end_time":      "2024-01-01T10:05:00Z",
			"read_duration": 300,
			"progress":      25.5,
			"device_type":   "web",
			"device_id":     "integration_test_device",
		}

		w := helper.DoAuthRequest("POST", ReaderReadingHistoryPath, requestBody, token)
		data := helper.AssertSuccess(w, 201, "阅读历史记录应该成功")

		historyID := data["id"]
		readDuration := data["read_duration"]
		helper.LogSuccess("阅读历史记录成功 - 历史ID: %v, 阅读时长: %v秒", historyID, readDuration)
	})

	t.Run("9.阅读历史_获取历史列表", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderReadingHistoryPath+"?page=1&page_size=10", nil, token)
		data := helper.AssertSuccess(w, 200, "获取阅读历史列表应该成功")

		if histories, ok := data["histories"].([]interface{}); ok {
			helper.LogSuccess("阅读历史列表获取成功，共 %d 条记录", len(histories))
		}
	})

	t.Run("10.阅读历史_获取统计", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderReadingHistoryPath+"/stats", nil, token)
		data := helper.AssertSuccess(w, 200, "获取阅读统计应该成功")

		totalDuration := data["total_duration"]
		totalBooks := data["total_books"]
		totalChapters := data["total_chapters"]
		helper.LogSuccess("阅读统计获取成功 - 总阅读时长: %v秒, 阅读书籍数: %v, 阅读章节数: %v",
			totalDuration, totalBooks, totalChapters)
	})

	t.Run("11.点赞_取消点赞", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s/like", ReaderBooksPath, testBookID)
		w := helper.DoAuthRequest("DELETE", url, nil, token)
		helper.AssertSuccess(w, 200, "取消点赞应该成功")

		helper.LogSuccess("取消点赞成功")
	})

	t.Run("12.收藏_取消收藏", func(t *testing.T) {
		// 如果有收藏ID，尝试删除
		if collectionID != "" {
			url := fmt.Sprintf("%s/%s", ReaderCollectionsPath, collectionID)
			w := helper.DoAuthRequest("DELETE", url, nil, token)
			helper.AssertSuccess(w, 200, "取消收藏应该成功")

			helper.LogSuccess("取消收藏成功")
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

	helper.LogSuccess("互动功能集成测试完成 - 测试场景: 收藏 → 评论 → 点赞 → 阅读历史 | 已测试API: 13个端点")
}
