//go:build e2e
// +build e2e

package layer1_basic

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	e2e "Qingyu_backend/test/e2e/framework"
)

func requireStatsResponseData(t *testing.T, env *e2e.TestEnvironment, path string, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	if w.Code != 200 {
		t.Fatalf("请求 %s 失败: 状态码 %d, 响应: %s", path, w.Code, w.Body.String())
	}

	response := env.ParseJSONResponse(w)
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("请求 %s 返回的 data 不是对象: %s", path, w.Body.String())
	}

	return data
}

// TestStatisticsView 测试数据分析查看功能
// @P1 重要功能测试 - 数据分析查看
// 测试场景：
// 1. 阅读数据统计
// 2. 收入统计
// 3. 读者画像
// 4. 数据导出
//
// TDD原则：先写测试，看测试失败，再写实现代码
func TestStatisticsView(t *testing.T) {
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

		author := fixtures.CreateUser()
		t.Logf("✓ 作者用户创建成功: %s", author.Username)

		token := actions.Login(author.Username, "Test1234")
		t.Logf("✓ 登录成功")

		env.SetTestData("auth_token", token)
	})

	// 步骤2: 创建项目并发布（需要发布后才能有统计数据）
	t.Run("步骤2_创建项目并发布", func(t *testing.T) {
		t.Log("创建项目并发布...")

		token := env.GetTestData("auth_token").(string)

		projectReq := map[string]interface{}{
			"title":       "e2e_test_statistics_project",
			"description": "数据分析测试项目",
			"genre":       "小说",
			"status":      "draft",
		}

		projectResp := actions.CreateProject(token, projectReq)

		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			var projectId string
			if id, ok := data["projectId"].(string); ok {
				projectId = id
			} else if id, ok := data["id"].(string); ok {
				projectId = id
			}

			if projectId != "" {
				t.Logf("✓ 项目创建成功 (ID: %s)", projectId)
				env.SetTestData("project_id", projectId)
			}
		}

		// 创建章节
		projectId, ok := env.GetTestData("project_id").(string)
		if !ok || projectId == "" {
			t.Log("⚠ 项目创建失败，跳过章节创建")
			return
		}
		chapterReq := map[string]interface{}{
			"project_id": projectId,
			"title":      "第一章：数据分析测试",
			"content":    "这是用于测试数据分析功能的章节内容",
			"word_count": 22,
		}

		w := env.DoRequest("POST", "/api/v1/writer/documents", chapterReq, token)

		var chapterId string
		if w.Code == 200 || w.Code == 201 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if id, ok := data["documentId"].(string); ok {
					chapterId = id
				} else if id, ok := data["document_id"].(string); ok {
					chapterId = id
				}

				if chapterId != "" {
					t.Logf("✓ 章节创建成功 (ID: %s)", chapterId)
					env.SetTestData("chapter_id", chapterId)
				}
			}
		}

		// 发布项目
		publishReq := map[string]interface{}{
			"action": "publish",
		}

		publishPath := fmt.Sprintf("/api/v1/writer/projects/%s/publish", projectId)
		w = env.DoRequest("POST", publishPath, publishReq, token)

		if w.Code == 200 || w.Code == 201 || w.Code == 202 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if bookId, ok := data["book_id"].(string); ok {
					t.Logf("✓ 项目发布成功，书城书籍ID: %s", bookId)
					env.SetTestData("book_id", bookId)
					env.SetTestData("published_book_id", bookId)
				} else {
					t.Logf("✓ 项目发布成功")
					// 使用projectId作为book_id（某些实现可能相同）
					env.SetTestData("book_id", projectId)
				}
			}
		}
	})

	// 步骤3: 获取作者统计概览
	t.Run("步骤3_获取作者统计概览", func(t *testing.T) {
		t.Log("获取作者统计概览...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)
		path := fmt.Sprintf("/api/v1/writer/stats/overview?projectId=%s", projectID)
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)

		if bookID, ok := data["bookId"].(string); ok {
			t.Logf("  作品ID: %s", bookID)
		}
		if totalViews, ok := data["totalViews"].(float64); ok {
			t.Logf("  总阅读量: %d", int(totalViews))
		}
		if subscribers, ok := data["subscribers"].(float64); ok {
			t.Logf("  订阅数: %d", int(subscribers))
		}
		if totalRevenue, ok := data["totalRevenue"].(float64); ok {
			t.Logf("  总收入: %.2f", totalRevenue)
		}
		if retentionRate, ok := data["retentionRate"].(float64); ok {
			t.Logf("  留存率: %.2f%%", retentionRate*100)
		}

		statsJSON, _ := json.MarshalIndent(data, "  ", "  ")
		t.Logf("作者统计概览数据: %s", string(statsJSON))
	})

	// 步骤4: 获取作者今日统计
	t.Run("步骤4_获取作者今日统计", func(t *testing.T) {
		t.Log("获取作者今日统计...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)
		path := fmt.Sprintf("/api/v1/writer/stats/today?projectId=%s", projectID)
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)

		if todayViews, ok := data["todayViews"].(float64); ok {
			t.Logf("  今日阅读量: %d", int(todayViews))
		}
		if todaySubscribers, ok := data["todaySubscribers"].(float64); ok {
			t.Logf("  今日订阅数: %d", int(todaySubscribers))
		}
		if todayWords, ok := data["todayWords"].(float64); ok {
			t.Logf("  今日字数: %d", int(todayWords))
		}

		statsJSON, _ := json.MarshalIndent(data, "  ", "  ")
		t.Logf("作者今日统计数据: %s", string(statsJSON))
	})

	// 步骤5: 获取阅读热力图
	t.Run("步骤5_获取阅读热力图", func(t *testing.T) {
		t.Log("获取阅读热力图...")

		token := env.GetTestData("auth_token").(string)
		bookId := env.GetTestData("book_id")

		if bookId == nil {
			t.Skip("书籍ID未设置，跳过此步骤")
			return
		}

		path := fmt.Sprintf("/api/v1/writer/books/%s/heatmap", bookId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 获取阅读热力图成功")

			if data, ok := response["data"].(map[string]interface{}); ok {
				if chapters, ok := data["chapters"].([]interface{}); ok {
					t.Logf("  热力图章节数: %d", len(chapters))

					// 显示前3个章节的热度数据
					for i, chapter := range chapters {
						if i >= 3 {
							break
						}
						if chapterMap, ok := chapter.(map[string]interface{}); ok {
							if title, ok := chapterMap["title"].(string); ok {
								if heat, ok := chapterMap["heat"].(float64); ok {
									t.Logf("  %s: 热度 %.2f", title, heat)
								}
							}
						}
					}
				}

				heatmapJSON, _ := json.MarshalIndent(data, "  ", "  ")
				t.Logf("阅读热力图数据: %s", string(heatmapJSON))
			}
		} else {
			t.Logf("获取阅读热力图响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤6: 获取收入统计
	t.Run("步骤6_获取收入统计", func(t *testing.T) {
		t.Log("获取收入统计...")

		token := env.GetTestData("auth_token").(string)
		bookId := env.GetTestData("book_id")

		if bookId == nil {
			t.Skip("书籍ID未设置，跳过此步骤")
			return
		}

		// 获取最近30天的收入数据
		endDate := time.Now().Format("2006-01-02")
		startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

		path := fmt.Sprintf("/api/v1/writer/books/%s/revenue?start_date=%s&end_date=%s", bookId.(string), startDate, endDate)
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 获取收入统计成功")

			if data, ok := response["data"].(map[string]interface{}); ok {
				if totalRevenue, ok := data["total_revenue"].(float64); ok {
					t.Logf("  总收入: %.2f", totalRevenue)
				}
				if chapterRevenue, ok := data["chapter_revenue"].(float64); ok {
					t.Logf("  章节收入: %.2f", chapterRevenue)
				}
				if giftRevenue, ok := data["gift_revenue"].(float64); ok {
					t.Logf("  礼物收入: %.2f", giftRevenue)
				}
				if tipRevenue, ok := data["tip_revenue"].(float64); ok {
					t.Logf("  打赏收入: %.2f", tipRevenue)
				}

				// 检查是否有细分数据
				if breakdown, ok := data["breakdown"].([]interface{}); ok {
					t.Logf("  收入明细记录数: %d", len(breakdown))
				}

				revenueJSON, _ := json.MarshalIndent(data, "  ", "  ")
				t.Logf("收入统计数据: %s", string(revenueJSON))
			}
		} else {
			t.Logf("获取收入统计响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤7: 获取章节表现统计
	t.Run("步骤7_获取章节表现统计", func(t *testing.T) {
		t.Log("获取章节表现统计...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)
		path := fmt.Sprintf("/api/v1/writer/stats/chapters?projectId=%s&page=1&size=10", projectID)
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)

		if total, ok := data["total"].(float64); ok {
			t.Logf("  章节统计总数: %d", int(total))
		}
		if items, ok := data["items"].([]interface{}); ok {
			t.Logf("  当前页章节数: %d", len(items))
		}

		topChaptersJSON, _ := json.MarshalIndent(data, "  ", "  ")
		t.Logf("章节表现数据: %s", string(topChaptersJSON))
	})

	// 步骤8: 获取阅读量趋势
	t.Run("步骤8_获取阅读量趋势", func(t *testing.T) {
		t.Log("获取阅读量趋势...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)
		path := fmt.Sprintf("/api/v1/writer/stats/views?projectId=%s&days=7", projectID)
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)

		if items, ok := data["items"].([]interface{}); ok {
			t.Logf("  趋势记录数: %d", len(items))
		}
		if totalViews, ok := data["total"].(float64); ok {
			t.Logf("  期间总阅读量: %d", int(totalViews))
		}

		dailyStatsJSON, _ := json.MarshalIndent(data, "  ", "  ")
		t.Logf("阅读量趋势数据: %s", string(dailyStatsJSON))
	})

	// 步骤9: 获取跳出点分析
	t.Run("步骤9_获取跳出点分析", func(t *testing.T) {
		t.Log("获取跳出点分析...")

		token := env.GetTestData("auth_token").(string)
		bookId := env.GetTestData("book_id")

		if bookId == nil {
			t.Skip("书籍ID未设置，跳过此步骤")
			return
		}

		path := fmt.Sprintf("/api/v1/writer/books/%s/drop-off-points", bookId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 获取跳出点分析成功")

			if data, ok := response["data"].(map[string]interface{}); ok {
				if dropOffPoints, ok := data["drop_off_points"].([]interface{}); ok {
					t.Logf("  跳出点数量: %d", len(dropOffPoints))

					// 显示前3个跳出点
					for i, point := range dropOffPoints {
						if i >= 3 {
							break
						}
						if pointMap, ok := point.(map[string]interface{}); ok {
							if chapter, ok := pointMap["chapter"].(string); ok {
								if rate, ok := pointMap["drop_off_rate"].(float64); ok {
									t.Logf("  %s: 跳出率 %.2f%%", chapter, rate*100)
								}
							}
						}
					}
				}

				dropOffJSON, _ := json.MarshalIndent(data, "  ", "  ")
				t.Logf("跳出点数据: %s", string(dropOffJSON))
			}
		} else {
			t.Logf("获取跳出点分析响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤10: 获取订阅趋势
	t.Run("步骤10_获取订阅趋势", func(t *testing.T) {
		t.Log("获取订阅趋势...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)
		path := fmt.Sprintf("/api/v1/writer/stats/subscribers?projectId=%s&days=7", projectID)
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)

		if total, ok := data["total"].(float64); ok {
			t.Logf("  新增订阅总数: %d", int(total))
		}
		if items, ok := data["items"].([]interface{}); ok {
			t.Logf("  订阅趋势记录数: %d", len(items))
		}

		retentionJSON, _ := json.MarshalIndent(data, "  ", "  ")
		t.Logf("订阅趋势数据: %s", string(retentionJSON))
	})
}

// TestReaderStatisticsAggregates 测试读者聚合统计接口
func TestReaderStatisticsAggregates(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	reader := fixtures.CreateUser()
	token := actions.Login(reader.Username, "Test1234")

	testCases := []struct {
		name string
		path string
	}{
		{name: "读者概览", path: "/api/v1/reader/statistics"},
		{name: "读者概览别名", path: "/api/v1/reader/statistics/overview"},
		{name: "阅读时长", path: "/api/v1/reader/statistics/reading-time?period=all"},
		{name: "阅读热力图", path: "/api/v1/reader/statistics/heatmap?days=7"},
		{name: "阅读趋势", path: "/api/v1/reader/statistics/trends?days=7"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := env.DoRequest("GET", tc.path, nil, token)
			data := requireStatsResponseData(t, env, tc.path, w)

			payloadJSON, _ := json.MarshalIndent(data, "  ", "  ")
			t.Logf("%s响应: %s", tc.name, string(payloadJSON))
		})
	}
}

// TestStatisticsReaderProfile 测试读者画像功能
// @P1 重要功能测试 - 读者画像
func TestStatisticsReaderProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	// 准备测试数据
	author := fixtures.CreateUser()
	token := actions.Login(author.Username, "Test1234")

	projectReq := map[string]interface{}{
		"title":       "reader_profile_test_project",
		"description": "读者画像测试项目",
		"genre":       "小说",
		"status":      "draft",
	}

	projectResp := actions.CreateProject(token, projectReq)
	var projectId string
	if data, ok := projectResp["data"].(map[string]interface{}); ok {
		if id, ok := data["projectId"].(string); ok {
			projectId = id
		} else if id, ok := data["id"].(string); ok {
			projectId = id
		}
	}

	// 发布项目
	publishReq := map[string]interface{}{
		"action": "publish",
	}

	publishPath := fmt.Sprintf("/api/v1/writer/projects/%s/publish", projectId)
	w := env.DoRequest("POST", publishPath, publishReq, token)

	var bookId string
	if w.Code == 200 || w.Code == 201 || w.Code == 202 {
		response := env.ParseJSONResponse(w)
		if data, ok := response["data"].(map[string]interface{}); ok {
			if id, ok := data["book_id"].(string); ok {
				bookId = id
			} else {
				bookId = projectId
			}
		}
	}

	// 测试读者画像相关API
	t.Run("测试读者年龄分布", func(t *testing.T) {
		// 注意：这个API可能不存在，需要根据实际实现调整
		if bookId == "" {
			t.Skip("书籍ID未设置")
			return
		}

		// 模拟读者画像数据查询
		t.Logf("⚠ 读者画像API需要实现，书籍ID: %s", bookId)
	})

	t.Run("测试读者地域分布", func(t *testing.T) {
		if bookId == "" {
			t.Skip("书籍ID未设置")
			return
		}

		t.Logf("⚠ 读者地域分布API需要实现，书籍ID: %s", bookId)
	})

	t.Run("测试读者阅读偏好", func(t *testing.T) {
		if bookId == "" {
			t.Skip("书籍ID未设置")
			return
		}

		t.Logf("⚠ 读者阅读偏好API需要实现，书籍ID: %s", bookId)
	})
}

// TestStatisticsExport 测试数据导出功能
// @P1 重要功能测试 - 数据导出
func TestStatisticsExport(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过 E2E 测试")
	}

	env, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	fixtures := env.Fixtures()
	actions := env.Actions()

	// 准备测试数据
	author := fixtures.CreateUser()
	token := actions.Login(author.Username, "Test1234")

	projectReq := map[string]interface{}{
		"title":       "export_test_project",
		"description": "数据导出测试项目",
		"genre":       "小说",
		"status":      "draft",
	}

	projectResp := actions.CreateProject(token, projectReq)
	var projectId string
	if data, ok := projectResp["data"].(map[string]interface{}); ok {
		if id, ok := data["projectId"].(string); ok {
			projectId = id
		} else if id, ok := data["id"].(string); ok {
			projectId = id
		}
	}

	// 测试导出功能
	t.Run("测试导出统计数据", func(t *testing.T) {
		// 注意：导出API可能不存在，需要根据实际实现调整
		t.Logf("⚠ 统计数据导出API需要实现，项目ID: %s", projectId)

		// 可能的API路径：
		// GET /api/v1/writer/books/:bookId/stats/export
		// GET /api/v1/writer/projects/:projectId/stats/export
		// GET /api/v1/writer/export/stats?book_id=:bookId
	})
}
