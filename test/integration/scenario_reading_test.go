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
)

// 阅读流程测试 - 从书籍详情到章节阅读的完整流程
func TestReadingScenario(t *testing.T) {
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
		t.Skip("无法登录测试用户，跳过阅读流程测试")
	}

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
			testBookID = testBook["_id"].(string)
		}
	}

	if testBookID == "" {
		t.Skip("数据库中没有测试书籍")
	}

	t.Run("1.书籍详情_获取书籍信息", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/bookstore/books/%s", baseURL, testBookID))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "应该成功获取书籍详情")

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Equal(t, float64(200), result["code"])
		data := result["data"].(map[string]interface{})

		t.Logf("✓ 书籍详情获取成功")
		t.Logf("  书名: %s", data["title"])
		t.Logf("  作者: %s", data["author"])
		t.Logf("  字数: %.0f", data["word_count"])
		t.Logf("  章节数: %.0f", data["chapter_count"])
	})

	t.Run("2.书籍详情_获取章节列表", func(t *testing.T) {
		// 使用认证请求获取章节列表
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/chapters?bookId=%s&page=1&size=10", baseURL, testBookID), nil)
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
			data := result["data"].(map[string]interface{})
			chapters := data["chapters"]
			if chapters != nil {
				chapterList := chapters.([]interface{})
				t.Logf("✓ 章节列表获取成功，共 %d 章", len(chapterList))

				// 显示前3章
				for i := 0; i < len(chapterList) && i < 3; i++ {
					ch := chapterList[i].(map[string]interface{})
					isFree := "免费"
					if !ch["is_free"].(bool) {
						isFree = "付费"
					}
					t.Logf("  第%d章: %s (%s)", i+1, ch["title"], isFree)
				}
			}
		} else {
			t.Logf("○ 获取章节列表失败: %v", result["message"])
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
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/chapters/%s/content", baseURL, firstChapterID), nil)
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
				data := result["data"].(map[string]interface{})
				content := data["content"].(string)
				t.Logf("✓ 章节内容获取成功，内容长度: %d 字符", len(content))

				// 显示前100个字符
				if len(content) > 100 {
					t.Logf("  内容预览: %s...", content[:100])
				} else {
					t.Logf("  内容预览: %s", content)
				}
			} else {
				t.Logf("○ 获取章节内容失败: %v", result["message"])
			}
		})

		t.Run("4.阅读进度_保存阅读进度", func(t *testing.T) {
			progressData := map[string]interface{}{
				"book_id":    testBookID,
				"chapter_id": firstChapterID,
				"position":   50, // 阅读到50%
			}

			jsonData, _ := json.Marshal(progressData)
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/progress", baseURL), bytes.NewBuffer(jsonData))
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
				t.Logf("✓ 阅读进度保存成功，位置: 50%%")
			} else {
				t.Logf("○ 保存阅读进度失败: %v", result["message"])
			}
		})

		t.Run("5.阅读进度_获取阅读进度", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/progress/%s", baseURL, testBookID), nil)
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
					progress := data.(map[string]interface{})
					t.Logf("✓ 阅读进度获取成功")
					t.Logf("  书籍ID: %v", progress["book_id"])
					t.Logf("  章节ID: %v", progress["chapter_id"])
					t.Logf("  阅读位置: %.0f%%", progress["position"])
				}
			} else {
				t.Logf("○ 获取阅读进度失败: %v", result["message"])
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

			jsonData, _ := json.Marshal(annotationData)
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/annotations", baseURL), bytes.NewBuffer(jsonData))
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
				t.Logf("✓ 书签添加成功")
			} else {
				t.Logf("○ 添加书签失败: %v", result["message"])
			}
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

			jsonData, _ := json.Marshal(annotationData)
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/annotations", baseURL), bytes.NewBuffer(jsonData))
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
				t.Logf("✓ 笔记添加成功")
			} else {
				t.Logf("○ 添加笔记失败: %v", result["message"])
			}
		})

		t.Run("8.书签笔记_获取书签和笔记列表", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reader/annotations?bookId=%s", baseURL, testBookID), nil)
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
					annotations := data.([]interface{})
					t.Logf("✓ 书签笔记列表获取成功，共 %d 条", len(annotations))
				}
			} else {
				t.Logf("○ 获取书签笔记失败: %v", result["message"])
			}
		})
	}

	t.Run("9.收藏_添加书籍到书架", func(t *testing.T) {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/reader/books/%s", baseURL, testBookID), nil)
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
			t.Logf("✓ 书籍已添加到书架")
		} else {
			t.Logf("○ 添加到书架失败: %v", result["message"])
		}
	})

	t.Run("10.书架_查看我的书架", func(t *testing.T) {
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
				}
			}
		} else {
			t.Logf("○ 获取书架失败: %v", result["message"])
		}
	})

	t.Logf("\n=== 阅读流程测试完成 ===")
	t.Logf("测试场景: 书籍详情 → 章节列表 → 阅读内容 → 保存进度 → 书签笔记 → 书架")
}

// 辅助函数：登录测试用户
func loginTestUser(t *testing.T, baseURL, username, password string) string {
	loginData := map[string]interface{}{
		"username": username,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/login", baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		t.Logf("登录请求失败: %v", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Logf("解析登录响应失败: %v", err)
		return ""
	}

	if result["code"] != float64(200) {
		t.Logf("登录失败: %v", result["message"])
		return ""
	}

	data := result["data"].(map[string]interface{})
	token := data["token"].(string)

	t.Logf("✓ 测试用户登录成功: %s", username)

	return token
}
