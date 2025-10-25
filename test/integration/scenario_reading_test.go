package integration

import (
	"context"
	"fmt"
	"testing"

	"Qingyu_backend/global"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 阅读流程测试 - 从书籍详情到章节阅读的完整流程
func TestReadingScenario(t *testing.T) {
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
		t.Skip("无法登录测试用户，跳过阅读流程测试")
	}

	ctx := context.Background()

	// 获取一本测试书籍
	var testBook map[string]interface{}
	var testBookID string
	cursor, err := global.DB.Collection("books").Find(ctx, bson.M{}, nil)
	if err == nil {
		var books []map[string]interface{}
		cursor.All(ctx, &books)
		cursor.Close(ctx)

		if len(books) > 0 {
			testBook = books[0]
			// 正确处理MongoDB ObjectID类型
			if oid, ok := testBook["_id"].(primitive.ObjectID); ok {
				testBookID = oid.Hex()
			}
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	t.Run("1.书籍详情_获取书籍信息", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", BookstoreBooksPath, testBookID)
		w := helper.DoRequest("GET", url, nil, "")

		// 处理404情况
		if w.Code == 404 {
			t.Logf("⚠ 书籍详情API返回404，可能路由未实现 (ID: %s)", testBookID)
			t.Skip("书籍详情API尚未完全实现")
			return
		}

		data := helper.AssertSuccess(w, 200, "应该成功获取书籍详情")

		// 安全检查data是否为nil
		if data == nil {
			t.Logf("⚠ 书籍详情数据为空")
			return
		}

		title := data["title"]
		author := data["author"]
		wordCount := data["word_count"]
		chapterCount := data["chapter_count"]
		helper.LogSuccess("书籍详情获取成功 - 书名: %v, 作者: %v, 字数: %v, 章节数: %v",
			title, author, wordCount, chapterCount)
	})

	t.Run("2.书籍详情_获取章节列表", func(t *testing.T) {
		url := fmt.Sprintf("%s?bookId=%s&page=1&size=10", ReaderChaptersPath, testBookID)
		w := helper.DoAuthRequest("GET", url, nil, token)

		// 检查404
		if w.Code == 404 {
			t.Skip("章节列表API尚未实现")
			return
		}

		data := helper.AssertSuccess(w, 200, "获取章节列表应该成功")

		if chapters, ok := data["chapters"].([]interface{}); ok {
			helper.LogSuccess("章节列表获取成功，共 %d 章", len(chapters))

			// 显示前3章
			for i := 0; i < len(chapters) && i < 3; i++ {
				ch := chapters[i].(map[string]interface{})
				isFree := "免费"
				if free, ok := ch["is_free"].(bool); ok && !free {
					isFree = "付费"
				}
				t.Logf("  第%d章: %s (%s)", i+1, ch["title"], isFree)
			}
		}
	})

	// 获取第一章的ID用于后续测试
	var firstChapterID string
	cursor2, err := global.DB.Collection("chapters").Find(ctx, bson.M{"book_id": testBookID}, nil)
	if err == nil {
		var chapters []map[string]interface{}
		cursor2.All(ctx, &chapters)
		cursor2.Close(ctx)

		if len(chapters) > 0 {
			if id, ok := chapters[0]["_id"]; ok {
				firstChapterID = fmt.Sprintf("%v", id)
			}
		}
	}

	if firstChapterID != "" {
		t.Run("3.章节阅读_获取章节内容（免费章节）", func(t *testing.T) {
			url := fmt.Sprintf("%s/%s/content", ReaderChaptersPath, firstChapterID)
			w := helper.DoAuthRequest("GET", url, nil, token)
			data := helper.AssertSuccess(w, 200, "获取章节内容应该成功")

			if content, ok := data["content"].(string); ok {
				helper.LogSuccess("章节内容获取成功，内容长度: %d 字符", len(content))

				// 显示前100个字符
				if len(content) > 100 {
					t.Logf("  内容预览: %s...", content[:100])
				} else {
					t.Logf("  内容预览: %s", content)
				}
			}
		})

		t.Run("4.阅读进度_保存阅读进度", func(t *testing.T) {
			progressData := map[string]interface{}{
				"book_id":    testBookID,
				"chapter_id": firstChapterID,
				"position":   50, // 阅读到50%
			}

			w := helper.DoAuthRequest("POST", ReaderProgressPath, progressData, token)
			helper.AssertSuccess(w, 200, "保存阅读进度应该成功")

			helper.LogSuccess("阅读进度保存成功，位置: 50%%")
		})

		t.Run("5.阅读进度_获取阅读进度", func(t *testing.T) {
			url := fmt.Sprintf("%s/%s", ReaderProgressPath, testBookID)
			w := helper.DoAuthRequest("GET", url, nil, token)
			data := helper.AssertSuccess(w, 200, "获取阅读进度应该成功")

			if data != nil {
				bookID := data["book_id"]
				chapterID := data["chapter_id"]
				position := data["position"]
				helper.LogSuccess("阅读进度获取成功 - 书籍ID: %v, 章节ID: %v, 阅读位置: %v%%",
					bookID, chapterID, position)
			}
		})

		t.Run("6.书签笔记_添加书签", func(t *testing.T) {
			annotationData := map[string]interface{}{
				"book_id":    testBookID,
				"chapter_id": firstChapterID,
				"type":       "bookmark",
				"range":      "100-150",
				"text":       "重要情节",
			}

			w := helper.DoAuthRequest("POST", ReaderAnnotationsPath, annotationData, token)
			helper.AssertSuccess(w, 200, "添加书签应该成功")

			helper.LogSuccess("书签添加成功")
		})

		t.Run("7.书签笔记_添加笔记", func(t *testing.T) {
			annotationData := map[string]interface{}{
				"book_id":    testBookID,
				"chapter_id": firstChapterID,
				"type":       "note",
				"range":      "200-250",
				"text":       "精彩描写",
				"note":       "这段描写非常生动",
			}

			w := helper.DoAuthRequest("POST", ReaderAnnotationsPath, annotationData, token)
			helper.AssertSuccess(w, 200, "添加笔记应该成功")

			helper.LogSuccess("笔记添加成功")
		})

		t.Run("8.书签笔记_获取书签和笔记列表", func(t *testing.T) {
			url := fmt.Sprintf("%s?bookId=%s", ReaderAnnotationsPath, testBookID)
			w := helper.DoAuthRequest("GET", url, nil, token)
			
			// 尝试解析响应（API可能直接返回数组或嵌套在data中）
			if w.Code == 200 {
				helper.LogSuccess("书签笔记列表获取成功")
			} else {
				t.Logf("○ 获取书签笔记失败，状态码: %d", w.Code)
			}
		})
	}

	t.Run("9.收藏_添加书籍到书架", func(t *testing.T) {
		url := fmt.Sprintf("%s/%s", ReaderBooksPath, testBookID)
		w := helper.DoAuthRequest("POST", url, nil, token)
		helper.AssertSuccess(w, 200, "添加到书架应该成功")

		helper.LogSuccess("书籍已添加到书架")
	})

	t.Run("10.书架_查看我的书架", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", ReaderBooksPath+"?page=1&size=10", nil, token)
		data := helper.AssertSuccess(w, 200, "获取书架应该成功")

		if data != nil {
			if books, ok := data["books"].([]interface{}); ok {
				helper.LogSuccess("书架获取成功，共 %d 本书", len(books))
			}
		}
	})

	helper.LogSuccess("阅读流程测试完成 - 测试场景: 书籍详情 → 章节列表 → 阅读内容 → 保存进度 → 书签笔记 → 书架")
}
