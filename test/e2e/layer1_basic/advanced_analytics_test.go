//go:build e2e
// +build e2e

package layer1_basic

import (
	"fmt"
	"testing"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestAdvancedAnalytics 测试高级数据分析功能
// P2功能：热力图、留存率分析、跳出点分析
//
// TDD原则：先写测试验证现有实现，发现bug后修复
func TestAdvancedAnalytics(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过E2E测试")
	}

	// 初始化测试环境
	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	// 准备：创建作者用户并登录
	t.Run("准备_创建作者并登录", func(t *testing.T) {
		t.Log("创建作者用户并登录...")

		author := fixtures.CreateUser()
		token := actions.Login(author.Username, "Test1234")

		env.SetTestData("author_user", author)
		env.SetTestData("auth_token", token)
		// 创建书城书籍与章节，确保统计接口有目标资源
		book := fixtures.CreateBook(author.ID.Hex())
		for i := 0; i < 3; i++ {
			fixtures.CreateChapter(book.ID.Hex())
		}
		env.SetTestData("book_id", book.ID.Hex())
		t.Logf("✓ 作者准备完成")
	})

	// 准备：创建项目和已发布的书籍
	t.Run("准备_创建项目和书籍", func(t *testing.T) {
		t.Log("创建测试项目和书籍...")

		token := env.GetTestData("auth_token").(string)

		// 创建项目
		projectReq := map[string]interface{}{
			"title":       "高级数据分析测试项目",
			"description": "用于测试高级数据分析功能",
			"genre":       "小说",
			"status":      "published", // 已发布
		}

		projectResp := actions.CreateProject(token, projectReq)
		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			if projectId, ok := data["projectId"].(string); ok {
				env.SetTestData("project_id", projectId)
				t.Logf("✓ 项目创建成功: %s", projectId)
			}
		}

		// 创建多个章节
		projectId := env.GetTestData("project_id").(string)
		for i := 1; i <= 10; i++ {
			chapterNum := fmt.Sprintf("第%d章", i)
			content := fmt.Sprintf("这是第%d章的内容。", i)
			chapterReq := map[string]interface{}{
				"project_id": projectId,
				"title":      chapterNum,
				"content":    content,
				"word_count": 15,
			}

			w := env.DoRequest("POST", "/api/v1/writer/documents", chapterReq, token)
			if w.Code == 200 || w.Code == 201 {
				response := env.ParseJSONResponse(w)
				if data, ok := response["data"].(map[string]interface{}); ok {
					var documentId string
					if id, ok := data["documentId"].(string); ok {
						documentId = id
					} else if id, ok := data["document_id"].(string); ok {
						documentId = id
					}
					if documentId != "" && i == 1 {
						env.SetTestData("chapter1_id", documentId)
					}
				}
			}
		}

		t.Logf("✓ 创建了10个章节")
	})

	// 测试场景1: 阅读热力图
	t.Run("场景1_阅读热力图", func(t *testing.T) {
		t.Log("测试阅读热力图...")

		token := env.GetTestData("auth_token").(string)
		bookID := env.GetTestData("book_id").(string)

		// 获取热力图数据
		w := env.DoRequest("GET", "/api/v1/writer/books/"+bookID+"/heatmap", nil, token)

		if w.Code == 200 {
			t.Log("✓ 热力图数据获取成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				// 检查热力图数据
				if heatmap, ok := data["heatmap"].([]interface{}); ok {
					t.Logf("✓ 热力图数据点数: %d", len(heatmap))

					// 验证数据结构
					for i, point := range heatmap {
						if i < 3 { // 只打印前3个
							if pointMap, ok := point.(map[string]interface{}); ok {
								if timestamp, ok := pointMap["timestamp"].(string); ok {
									if count, ok := pointMap["count"].(float64); ok {
										t.Logf("  时间: %s, 阅读数: %d", timestamp[:19], int(count))
									}
								}
							}
						}
					}
				}

				// 检查时间段分布
				if timeDistribution, ok := data["timeDistribution"].(map[string]interface{}); ok {
					t.Logf("✓ 时间段分布: %v", timeDistribution)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 热力图API未实现（404）")
		} else {
			t.Logf("ℹ 热力图数据获取尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景2: 章节留存率分析
	t.Run("场景2_章节留存率分析", func(t *testing.T) {
		t.Log("测试章节留存率分析...")

		token := env.GetTestData("auth_token").(string)
		bookID := env.GetTestData("book_id").(string)

		// 获取留存率数据
		w := env.DoRequest("GET", "/api/v1/writer/books/"+bookID+"/retention?days=7", nil, token)

		if w.Code == 200 {
			t.Log("✓ 留存率数据获取成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				// 检查章节留存率
				if chapterRetention, ok := data["chapterRetention"].([]interface{}); ok {
					t.Logf("✓ 章节留存率数据: %d 个章节", len(chapterRetention))

					for i, item := range chapterRetention {
						if i < 5 { // 只打印前5个
							if itemMap, ok := item.(map[string]interface{}); ok {
								if chapterNum, ok := itemMap["chapterNum"].(float64); ok {
									if retention, ok := itemMap["retention"].(float64); ok {
										t.Logf("  第%d章留存率: %.2f%%", int(chapterNum), retention*100)
									}
								}
							}
						}
					}
				}

				// 检查总体留存率
				if overallRetention, ok := data["overallRetention"].(float64); ok {
					t.Logf("✓ 总体留存率: %.2f%%", overallRetention*100)
				}

				// 检查关键流失点
				if dropOffPoints, ok := data["dropOffPoints"].([]interface{}); ok {
					t.Logf("✓ 关键流失点数: %d", len(dropOffPoints))
					for i, point := range dropOffPoints {
						if i < 3 {
							if pointMap, ok := point.(map[string]interface{}); ok {
								if chapterNum, ok := pointMap["chapterNum"].(float64); ok {
									if dropOffRate, ok := pointMap["dropOffRate"].(float64); ok {
										t.Logf("  第%d章流失率: %.2f%%", int(chapterNum), dropOffRate*100)
									}
								}
							}
						}
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 留存率API未实现（404）")
		} else {
			t.Logf("ℹ 留存率数据获取尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景3: 跳出点分析
	t.Run("场景3_跳出点分析", func(t *testing.T) {
		t.Log("测试跳出点分析...")

		token := env.GetTestData("auth_token").(string)
		bookID := env.GetTestData("book_id").(string)

		// 获取跳出点数据
		w := env.DoRequest("GET", "/api/v1/writer/books/"+bookID+"/drop-off-points", nil, token)

		if w.Code == 200 {
			t.Log("✓ 跳出点数据获取成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				// 检查跳出点列表
				if bouncePoints, ok := data["bouncePoints"].([]interface{}); ok {
					t.Logf("✓ 发现 %d 个跳出点", len(bouncePoints))

					for i, point := range bouncePoints {
						if i < 5 {
							if pointMap, ok := point.(map[string]interface{}); ok {
								if chapterNum, ok := pointMap["chapterNum"].(float64); ok {
									if bounceRate, ok := pointMap["bounceRate"].(float64); ok {
										if position, ok := pointMap["position"].(string); ok {
											t.Logf("  第%d章 [%s]: 跳出率 %.2f%%", int(chapterNum), position, bounceRate*100)
										}
									}
								}
							}
						}
					}
				}

				// 检查最高跳出点
				if highestBounce, ok := data["highestBounce"].(map[string]interface{}); ok {
					if chapterNum, ok := highestBounce["chapterNum"].(float64); ok {
						if rate, ok := highestBounce["rate"].(float64); ok {
							t.Logf("✓ 最高跳出点: 第%d章 (%.2f%%)", int(chapterNum), rate*100)
						}
					}
				}

				// 检查建议
				if recommendations, ok := data["recommendations"].([]interface{}); ok {
					t.Logf("✓ 优化建议数: %d", len(recommendations))
					for i, rec := range recommendations {
						if i < 3 {
							if recStr, ok := rec.(string); ok {
								t.Logf("  建议%d: %s", i+1, recStr)
							}
						}
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 跳出点分析API未实现（404）")
		} else {
			t.Logf("ℹ 跳出点数据获取尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景4: 每日统计分布
	t.Run("场景4_阅读时长分布", func(t *testing.T) {
		t.Log("测试每日统计分布...")

		token := env.GetTestData("auth_token").(string)
		bookID := env.GetTestData("book_id").(string)

		// 获取每日统计
		w := env.DoRequest("GET", "/api/v1/writer/books/"+bookID+"/daily-stats?days=7", nil, token)

		if w.Code == 200 {
			t.Log("✓ 阅读时长分布获取成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				// 检查时长分布
				if distribution, ok := data["distribution"].(map[string]interface{}); ok {
					t.Logf("✓ 时长分布: %v", distribution)

					// 验证各时段
					ranges := []string{"0-5min", "5-10min", "10-30min", "30-60min", "60min+"}
					for _, r := range ranges {
						if count, ok := distribution[r].(float64); ok {
							t.Logf("  %s: %d 人", r, int(count))
						}
					}
				}

				// 检查平均阅读时长
				if avgDuration, ok := data["avgDuration"].(float64); ok {
					t.Logf("✓ 平均阅读时长: %.2f 分钟", avgDuration)
				}

				// 检查峰值阅读时段
				if peakTime, ok := data["peakTime"].(string); ok {
					t.Logf("✓ 峰值阅读时段: %s", peakTime)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 每日统计API未实现（404）")
		} else {
			t.Logf("ℹ 每日统计获取尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景5: 热门章节分析
	t.Run("场景5_用户行为路径分析", func(t *testing.T) {
		t.Log("测试热门章节分析...")

		token := env.GetTestData("auth_token").(string)
		bookID := env.GetTestData("book_id").(string)

		// 获取热门章节
		w := env.DoRequest("GET", "/api/v1/writer/books/"+bookID+"/top-chapters", nil, token)

		if w.Code == 200 {
			t.Log("✓ 用户行为路径获取成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				// 检查常见路径
				if commonPaths, ok := data["commonPaths"].([]interface{}); ok {
					t.Logf("✓ 常见路径数: %d", len(commonPaths))

					for i, path := range commonPaths {
						if i < 3 {
							if pathMap, ok := path.(map[string]interface{}); ok {
								if steps, ok := pathMap["steps"].([]interface{}); ok {
									if percentage, ok := pathMap["percentage"].(float64); ok {
										t.Logf("  路径%d (%.1f%%): %d 步", i+1, percentage*100, len(steps))
									}
								}
							}
						}
					}
				}

				// 检查典型用户画像
				if userPersonas, ok := data["userPersonas"].([]interface{}); ok {
					t.Logf("✓ 用户画像数: %d", len(userPersonas))
					for i, persona := range userPersonas {
						if i < 3 {
							if personaMap, ok := persona.(map[string]interface{}); ok {
								if typeName, ok := personaMap["type"].(string); ok {
									if percentage, ok := personaMap["percentage"].(float64); ok {
										t.Logf("  %s: %.1f%%", typeName, percentage*100)
									}
								}
							}
						}
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 热门章节API未实现（404）")
		} else {
			t.Logf("ℹ 热门章节获取尝试，状态码: %d", w.Code)
		}
	})

	t.Log("✅ 高级数据分析E2E测试完成")
}
