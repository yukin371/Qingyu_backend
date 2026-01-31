//go:build e2e
// +build e2e

package layer1_basic

import (
	"testing"

	"Qingyu_backend/models/users"
	e2e "Qingyu_backend/test/e2e/framework"
)

// TestReadingFlow 测试阅读流程
// 流程: 获取书籍列表 -> 查看书籍详情 -> 获取章节列表 -> 阅读章节内容 -> 导航到下一章 -> 保存阅读进度
func TestReadingFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	// 步骤1: 创建测试用户并登录
	t.Run("步骤1_创建测试用户并登录", func(t *testing.T) {
		t.Log("创建测试用户并登录...")

		// 创建读者用户
		reader := fixtures.CreateUser()
		t.Logf("✓ 读者用户创建成功: %s (角色: %v)", reader.Username, reader.Roles)

		// 登录获取token
		token := actions.Login(reader.Username, "Test1234")
		t.Logf("✓ 登录成功，获取token")

		// 保存到环境
		env.SetTestData("reader_user", reader)
		env.SetTestData("auth_token", token)
	})

	// 步骤2: 获取书籍列表
	t.Run("步骤2_获取书籍列表", func(t *testing.T) {
		t.Log("获取书籍列表...")

		token := env.GetTestData("auth_token")
		if token == nil {
			t.Skip("auth_token未设置，跳过此步骤")
			return
		}

		// 获取书城首页（包含书籍列表）
		homepage := actions.GetBookstoreHomepage()

		if data, ok := homepage["data"].(map[string]interface{}); ok {
			t.Logf("✓ 书城首页获取成功")

			// 检查是否有书籍数据
			if featuredBooks, ok := data["featured_books"].([]interface{}); ok && len(featuredBooks) > 0 {
				t.Logf("✓ 推荐书籍数量: %d", len(featuredBooks))

				// 保存第一本书的ID用于后续测试
				if firstBook, ok := featuredBooks[0].(map[string]interface{}); ok {
					if bookID, ok := firstBook["id"].(string); ok {
						env.SetTestData("test_book_id", bookID)
						t.Logf("✓ 选择测试书籍: %s", bookID)
					}
				}
			} else {
				// 如果首页没有书籍，创建一本测试书籍
				t.Log("首页没有书籍，创建测试书籍...")
				reader := env.GetTestData("reader_user").(*users.User)
				book := fixtures.CreateBook(reader.ID.Hex(),
					e2e.WithBookTitle("e2e_test_reading_book"),
					e2e.WithBookPrice(0),
				)
				env.SetTestData("test_book_id", book.ID.Hex())
				t.Logf("✓ 创建测试书籍: %s", book.ID.Hex())

				// 为书籍创建章节
				fixtures.CreateChapter(book.ID.Hex(),
					e2e.WithChapterTitle("e2e_test_chapter_1"),
					e2e.WithChapterFree(true),
				)
				fixtures.CreateChapter(book.ID.Hex(),
					e2e.WithChapterTitle("e2e_test_chapter_2"),
					e2e.WithChapterFree(true),
				)
				t.Logf("✓ 创建2个测试章节")
			}
		} else {
			t.Error("首页响应格式不正确")
		}
	})

	// 步骤3: 获取书籍详情
	t.Run("步骤3_获取书籍详情", func(t *testing.T) {
		t.Log("获取书籍详情...")

		bookID := env.GetTestData("test_book_id")
		if bookID == nil {
			t.Skip("test_book_id未设置，跳过此步骤")
			return
		}

		// 获取书籍详情
		bookDetail := actions.GetBookDetail(bookID.(string))

		if data, ok := bookDetail["data"].(map[string]interface{}); ok {
			if title, ok := data["title"].(string); ok {
				t.Logf("✓ 书籍详情获取成功: %s", title)

				// 验证书籍详情完整性
				requiredFields := []string{"id", "title", "author", "introduction", "status", "wordCount"}
				missingFields := []string{}
				for _, field := range requiredFields {
					if _, exists := data[field]; !exists {
						missingFields = append(missingFields, field)
					}
				}

				if len(missingFields) > 0 {
					t.Errorf("书籍详情缺少字段: %v", missingFields)
				} else {
					t.Logf("✓ 书籍详情完整性验证通过")
				}
			} else {
				t.Error("书籍详情中未找到title字段")
			}
		} else {
			t.Error("书籍详情响应格式不正确")
		}
	})

	// 步骤4: 获取章节目录
	t.Run("步骤4_获取章节目录", func(t *testing.T) {
		t.Log("获取章节目录...")

		bookID := env.GetTestData("test_book_id")
		token := env.GetTestData("auth_token")
		if bookID == nil || token == nil {
			t.Skip("test_book_id或auth_token未设置，跳过此步骤")
			return
		}

		// 获取章节列表
		chapterList := actions.GetChapterList(bookID.(string), token.(string))

		if data, ok := chapterList["data"].(map[string]interface{}); ok {
			if chapters, ok := data["chapters"].([]interface{}); ok {
				t.Logf("✓ 章节列表获取成功，章节数量: %d", len(chapters))

				if len(chapters) > 0 {
					// 保存第一章的ID用于后续测试
					if firstChapter, ok := chapters[0].(map[string]interface{}); ok {
						if chapterID, ok := firstChapter["id"].(string); ok {
							env.SetTestData("test_chapter_id", chapterID)
							t.Logf("✓ 选择测试章节: %s", chapterID)
						}
					}
				} else {
					t.Error("章节列表为空")
				}
			} else {
				t.Error("章节列表格式不正确")
			}
		} else {
			t.Error("章节列表响应格式不正确")
		}
	})

	// 步骤5: 阅读第一章 - 获取章节内容
	t.Run("步骤5_阅读第一章_获取章节内容", func(t *testing.T) {
		t.Log("阅读第一章...")

		chapterID := env.GetTestData("test_chapter_id")
		token := env.GetTestData("auth_token")
		if chapterID == nil || token == nil {
			t.Skip("test_chapter_id或auth_token未设置，跳过此步骤")
			return
		}

		// 获取章节内容
		chapterContent := actions.GetChapter(chapterID.(string), token.(string))

		if data, ok := chapterContent["data"].(map[string]interface{}); ok {
			if title, ok := data["title"].(string); ok {
				t.Logf("✓ 章节内容获取成功: %s", title)

				// 验证章节内容完整性
				if content, ok := data["content"].(string); ok && len(content) > 0 {
					t.Logf("✓ 章节内容完整，内容长度: %d", len(content))
				} else {
					t.Error("章节内容为空")
				}

				// 保存章节信息用于导航测试
				if bookID, ok := data["book_id"].(string); ok {
					env.SetTestData("current_book_id", bookID)
				}
			} else {
				t.Error("章节内容中未找到title字段")
			}
		} else {
			t.Error("章节内容响应格式不正确")
		}
	})

	// 步骤6: 阅读下一章 - 验证导航正确
	t.Run("步骤6_阅读下一章_验证导航正确", func(t *testing.T) {
		t.Log("导航到下一章...")

		bookID := env.GetTestData("test_book_id")
		chapterID := env.GetTestData("test_chapter_id")
		token := env.GetTestData("auth_token")
		if bookID == nil || chapterID == nil || token == nil {
			t.Skip("必要参数未设置，跳过此步骤")
			return
		}

		// 获取下一章
		path := "/api/v1/reader/books/" + bookID.(string) + "/chapters/" + chapterID.(string) + "/next"
		w := env.DoRequest("GET", path, nil, token.(string))

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if nextChapterID, ok := data["id"].(string); ok {
					t.Logf("✓ 下一章获取成功: %s", nextChapterID)

					// 验证下一章ID与当前章节不同
					if nextChapterID != chapterID.(string) {
						t.Logf("✓ 导航正确：下一章ID与当前章节不同")
					} else {
						t.Error("导航错误：下一章ID与当前章节相同")
					}
				} else {
					t.Log("✓ 已是最后一章（没有下一章）")
				}
			} else {
				t.Error("下一章响应格式不正确")
			}
		} else {
			t.Logf("✓ 已是最后一章（状态码: %d）", w.Code)
		}
	})

	// 步骤7: 记录阅读进度 - 验证进度保存
	t.Run("步骤7_记录阅读进度_验证进度保存", func(t *testing.T) {
		t.Log("记录阅读进度...")

		reader := env.GetTestData("reader_user")
		bookID := env.GetTestData("test_book_id")
		chapterID := env.GetTestData("test_chapter_id")
		token := env.GetTestData("auth_token")
		if reader == nil || bookID == nil || chapterID == nil || token == nil {
			t.Skip("必要参数未设置，跳过此步骤")
			return
		}

		readerUser := reader.(*users.User)

		// 保存阅读进度
		progressResp := actions.StartReading(readerUser.ID.Hex(), bookID.(string), chapterID.(string), token.(string))

		if data, ok := progressResp["data"].(map[string]interface{}); ok {
			if _, ok := data["progress"]; ok {
				t.Logf("✓ 阅读进度保存成功")
			} else {
				t.Error("阅读进度响应中未找到progress字段")
			}
		} else {
			t.Error("阅读进度响应格式不正确")
		}

		// 验证进度可以正确获取
		savedProgress := actions.GetReadingProgress(readerUser.ID.Hex(), bookID.(string), token.(string))

		if data, ok := savedProgress["data"].(map[string]interface{}); ok {
			if progressList, ok := data["progress"].([]interface{}); ok && len(progressList) > 0 {
				t.Logf("✓ 阅读进度获取成功，进度记录数量: %d", len(progressList))
			} else {
				t.Error("未找到阅读进度记录")
			}
		} else {
			t.Error("阅读进度获取响应格式不正确")
		}
	})
}



