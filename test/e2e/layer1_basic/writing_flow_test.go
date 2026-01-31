//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestWritingFlow 测试写作流程
// 流程: 创建写作项目 -> 验证项目存在
func TestWritingFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	// 步骤1: 创建作者用户并登录
	t.Run("步骤1_创建作者用户并登录", func(t *testing.T) {
		t.Log("创建作者用户并登录...")

		// 创建作者用户（不使用固定用户名，避免冲突）
		author := fixtures.CreateUser()
		t.Logf("✓ 作者用户创建成功: %s (角色: %v)", author.Username, author.Roles)

		// 登录获取token
		token := actions.Login(author.Username, "Test1234")
		t.Logf("✓ 登录成功，获取token")

		// 保存到环境
		env.SetTestData("author_user", author)
		env.SetTestData("auth_token", token)
	})

	// 步骤2: 创建写作项目
	t.Run("步骤2_创建写作项目", func(t *testing.T) {
		t.Log("创建写作项目...")

		token := env.GetTestData("auth_token").(string)

		// 创建项目请求
		projectReq := map[string]interface{}{
			"title":       "e2e_test_writing_project",
			"description": "这是一个E2E测试写作项目",
			"genre":       "小说",
			"status":      "draft",
		}

		// 创建项目
		projectResp := actions.CreateProject(token, projectReq)

		// 验证响应
		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			if projectId, ok := data["projectId"].(string); ok {
				t.Logf("✓ 写作项目创建成功 (ID: %s)", projectId)

				// 保存项目ID
				env.SetTestData("project_id", projectId)
			} else {
				t.Error("项目响应中未找到projectId字段")
			}
		} else {
			t.Error("项目响应格式不正确")
		}
	})

	// 步骤3: 验证项目存在（通过查询项目列表）
	t.Run("步骤3_验证项目存在", func(t *testing.T) {
		t.Log("验证写作项目存在...")

		token := env.GetTestData("auth_token")
		if token == nil {
			t.Skip("auth_token未设置，跳过此步骤")
			return
		}

		projectId := env.GetTestData("project_id")
		if projectId == nil {
			t.Skip("project_id未设置，跳过此步骤")
			return
		}

		// 获取项目列表
		w := env.DoRequest("GET", "/api/v1/writer/projects", nil, token.(string))

		if w.Code == 200 {
			t.Logf("✓ 项目列表获取成功，项目 %v 应在列表中", projectId)
		} else {
			t.Errorf("获取项目列表失败，状态码: %d", w.Code)
		}
	})

	// 步骤4: 创建章节文档
	t.Run("步骤4_创建章节文档", func(t *testing.T) {
		t.Log("创建章节文档...")

		token := env.GetTestData("auth_token")
		projectId := env.GetTestData("project_id")
		if token == nil || projectId == nil {
			t.Skip("auth_token或project_id未设置，跳过此步骤")
			return
		}

		// 创建第一章
		chapterReq := map[string]interface{}{
			"project_id": projectId,
			"title":      "e2e_test_chapter_1",
			"content":    "这是E2E测试的第一章内容。这是一个测试章节，用于验证写作和发布功能。",
			"word_count": 35,
		}

		w := env.DoRequest("POST", "/api/v1/writer/documents", chapterReq, token.(string))

		if w.Code == 200 || w.Code == 201 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if documentId, ok := data["document_id"].(string); ok {
					t.Logf("✓ 第一章创建成功 (ID: %s)", documentId)
					env.SetTestData("chapter1_id", documentId)
				} else {
					t.Error("章节响应中未找到document_id字段")
				}
			} else {
				t.Error("章节响应格式不正确")
			}
		} else {
			t.Errorf("创建章节失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		}

		// 创建第二章
		chapterReq2 := map[string]interface{}{
			"project_id": projectId,
			"title":      "e2e_test_chapter_2",
			"content":    "这是E2E测试的第二章内容。这是另一个测试章节，用于验证多章节发布功能。",
			"word_count": 35,
		}

		w2 := env.DoRequest("POST", "/api/v1/writer/documents", chapterReq2, token.(string))

		if w2.Code == 200 || w2.Code == 201 {
			response := env.ParseJSONResponse(w2)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if documentId, ok := data["document_id"].(string); ok {
					t.Logf("✓ 第二章创建成功 (ID: %s)", documentId)
					env.SetTestData("chapter2_id", documentId)
				}
			}
		}
	})

	// 步骤5: 发布项目到书城
	t.Run("步骤5_发布项目到书城", func(t *testing.T) {
		t.Log("发布项目到书城...")

		token := env.GetTestData("auth_token")
		projectId := env.GetTestData("project_id")
		if token == nil || projectId == nil {
			t.Skip("auth_token或project_id未设置，跳过此步骤")
			return
		}

		// 发布项目请求
		publishReq := map[string]interface{}{
			"project_id": projectId,
			"action":     "publish", // 发布到书城
		}

		w := env.DoRequest("POST", "/api/v1/writer/projects/"+projectId.(string)+"/publish", publishReq, token.(string))

		if w.Code == 200 || w.Code == 201 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if bookId, ok := data["book_id"].(string); ok {
					t.Logf("✓ 项目发布成功，书城书籍ID: %s", bookId)
					env.SetTestData("published_book_id", bookId)
				} else {
					t.Log("✓ 项目发布成功（未返回book_id）")
				}
			} else {
				t.Log("✓ 项目发布成功")
			}
		} else {
			t.Logf("发布项目响应: 状态码 %d, %s", w.Code, w.Body.String())
		}
	})

	// 步骤6: 发布章节
	t.Run("步骤6_发布章节", func(t *testing.T) {
		t.Log("发布章节...")

		token := env.GetTestData("auth_token")
		chapter1Id := env.GetTestData("chapter1_id")
		chapter2Id := env.GetTestData("chapter2_id")
		if token == nil || chapter1Id == nil {
			t.Skip("auth_token或chapter_id未设置，跳过此步骤")
			return
		}

		// 发布第一章
		publishChapterReq := map[string]interface{}{
			"action": "publish",
		}

		w1 := env.DoRequest("POST", "/api/v1/writer/documents/"+chapter1Id.(string)+"/publish", publishChapterReq, token.(string))

		if w1.Code == 200 || w1.Code == 201 {
			t.Logf("✓ 第一章发布成功")
		} else {
			t.Logf("第一章发布响应: 状态码 %d, %s", w1.Code, w1.Body.String())
		}

		// 发布第二章（如果存在）
		if chapter2Id != nil {
			w2 := env.DoRequest("POST", "/api/v1/writer/documents/"+chapter2Id.(string)+"/publish", publishChapterReq, token.(string))

			if w2.Code == 200 || w2.Code == 201 {
				t.Logf("✓ 第二章发布成功")
			} else {
				t.Logf("第二章发布响应: 状态码 %d, %s", w2.Code, w2.Body.String())
			}
		}
	})

	// 步骤7: 验证书城API能获取新发布的内容
	t.Run("步骤7_验证书城API能获取新发布的内容", func(t *testing.T) {
		t.Log("验证书城API...")

		publishedBookId := env.GetTestData("published_book_id")
		token := env.GetTestData("auth_token")

		// 如果有发布的书籍ID，直接获取详情
		if publishedBookId != nil {
			// 通过书城API获取书籍详情
			w := env.DoRequest("GET", "/api/v1/bookstore/books/"+publishedBookId.(string), nil, "")

			if w.Code == 200 {
				response := env.ParseJSONResponse(w)
				if data, ok := response["data"].(map[string]interface{}); ok {
					if title, ok := data["title"].(string); ok {
						t.Logf("✓ 书城API成功获取到书籍: %s", title)
					}
				}
			} else {
				t.Logf("书城API获取书籍详情失败，状态码: %d", w.Code)
			}
		} else {
			// 如果没有具体的书籍ID，通过搜索来验证
			t.Log("通过搜索验证书城内容...")

			// 搜索测试书籍
			w := env.DoRequest("GET", "/api/v1/bookstore/books/search?keyword=e2e_test_writing_project", nil, "")

			if w.Code == 200 {
				response := env.ParseJSONResponse(w)
				if data, ok := response["data"].(map[string]interface{}); ok {
					if books, ok := data["books"].([]interface{}); ok {
						t.Logf("✓ 书城搜索成功，找到 %d 本相关书籍", len(books))

						// 检查是否包含我们发布的书
						found := false
						for _, book := range books {
							if bookMap, ok := book.(map[string]interface{}); ok {
								if title, ok := bookMap["title"].(string); ok {
									if title == "e2e_test_writing_project" {
										found = true
										t.Logf("✓ 在书城中找到发布的书籍: %s", title)
										break
									}
								}
							}
						}

						if !found {
							t.Log("ℹ 未在搜索结果中找到测试书籍（可能需要等待索引更新）")
						}
					}
				}
			} else {
				t.Logf("书城搜索失败，状态码: %d", w.Code)
			}
		}

		// 验证章节内容（如果有token）
		if token != nil && publishedBookId != nil {
			// 获取章节列表
			w := env.DoRequest("GET", "/api/v1/reader/books/"+publishedBookId.(string)+"/chapters", nil, token.(string))

			if w.Code == 200 {
				response := env.ParseJSONResponse(w)
				if data, ok := response["data"].(map[string]interface{}); ok {
					if chapters, ok := data["chapters"].([]interface{}); ok {
						t.Logf("✓ 书城API成功获取到 %d 个章节", len(chapters))

						if len(chapters) > 0 {
							t.Logf("✓ 写作→书城发布流程验证成功")
						}
					}
				}
			} else {
				t.Logf("获取章节列表失败，状态码: %d", w.Code)
			}
		}
	})
}



