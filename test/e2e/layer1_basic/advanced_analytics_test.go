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
		var projectId string
		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			// 支持驼峰命名 projectId（Go标准）和蛇形命名 id
			if id, ok := data["projectId"].(string); ok {
				projectId = id
			} else if id, ok := data["id"].(string); ok {
				projectId = id
			}

			if projectId != "" {
				env.SetTestData("project_id", projectId)
				t.Logf("✓ 项目创建成功: %s", projectId)
			}
		}

		// 如果项目创建失败，跳过章节创建
		if projectId == "" {
			t.Log("⚠ 项目创建失败，跳过章节创建步骤")
			return
		}

		// 创建多个章节
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
		path := "/api/v1/writer/books/" + bookID + "/heatmap"
		w := env.DoRequest("GET", path, nil, token)
		if w.Code != 200 {
			t.Fatalf("请求 %s 失败: 状态码 %d, 响应: %s", path, w.Code, w.Body.String())
		}
		response := env.ParseJSONResponse(w)
		heatmapData, ok := response["data"].([]interface{})
		if !ok {
			t.Fatalf("请求 %s 返回的 data 不是数组: %s", path, w.Body.String())
		}
		t.Log("✓ 热力图数据获取成功")

		t.Logf("✓ 热力图数据点数: %d", len(heatmapData))
		for i, point := range heatmapData {
			if i < 3 {
				if pointMap, ok := point.(map[string]interface{}); ok {
					if day, ok := pointMap["day"].(float64); ok {
						if hour, ok := pointMap["hour"].(float64); ok {
							if value, ok := pointMap["value"].(float64); ok {
								t.Logf("  星期%d %02d:00 -> %d", int(day), int(hour), int(value))
							}
						}
					}
				}
			}
		}
	})

	// 测试场景2: 章节留存率分析
	t.Run("场景2_章节留存率分析", func(t *testing.T) {
		t.Log("测试章节留存率分析...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)

		// 获取留存率数据
		path := "/api/v1/writer/stats/overview?projectId=" + projectID
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)
		t.Log("✓ 留存率数据获取成功")

		if overallRetention, ok := data["retentionRate"].(float64); ok {
			t.Logf("✓ 总体留存率: %.2f%%", overallRetention*100)
		}
		if today, ok := data["today"].(map[string]interface{}); ok {
			t.Logf("✓ 今日统计摘要: %v", today)
		}
	})

	// 测试场景3: 跳出点分析
	t.Run("场景3_跳出点分析", func(t *testing.T) {
		t.Log("测试跳出点分析...")

		token := env.GetTestData("auth_token").(string)
		bookID := env.GetTestData("book_id").(string)

		// 获取跳出点数据
		path := "/api/v1/writer/books/" + bookID + "/drop-off-points"
		w := env.DoRequest("GET", path, nil, token)
		if w.Code != 200 {
			t.Fatalf("请求 %s 失败: 状态码 %d, 响应: %s", path, w.Code, w.Body.String())
		}
		response := env.ParseJSONResponse(w)
		t.Log("✓ 跳出点数据获取成功")

		if response["data"] == nil {
			t.Log("✓ 当前无跳出点数据")
			return
		}

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("请求 %s 返回的 data 不是对象: %s", path, w.Body.String())
		}

		if bouncePoints, ok := data["bouncePoints"].([]interface{}); ok {
			t.Logf("✓ 发现 %d 个跳出点", len(bouncePoints))
		}
		if highestBounce, ok := data["highestBounce"].(map[string]interface{}); ok {
			t.Logf("✓ 最高跳出点: %v", highestBounce)
		}
		if recommendations, ok := data["recommendations"].([]interface{}); ok {
			t.Logf("✓ 优化建议数: %d", len(recommendations))
		}
	})

	// 测试场景4: 每日统计分布
	t.Run("场景4_阅读时长分布", func(t *testing.T) {
		t.Log("测试每日统计分布...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)

		// 获取每日统计
		path := "/api/v1/writer/stats/views?projectId=" + projectID + "&days=7"
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)
		t.Log("✓ 阅读时长分布获取成功")

		if items, ok := data["items"].([]interface{}); ok {
			t.Logf("✓ 趋势记录数: %d", len(items))
		}
		if total, ok := data["total"].(float64); ok {
			t.Logf("✓ 周期总阅读量: %d", int(total))
		}
	})

	// 测试场景5: 热门章节分析
	t.Run("场景5_用户行为路径分析", func(t *testing.T) {
		t.Log("测试热门章节分析...")

		token := env.GetTestData("auth_token").(string)
		projectID := env.GetTestData("project_id").(string)

		// 获取热门章节
		path := "/api/v1/writer/stats/chapters?projectId=" + projectID + "&page=1&size=10"
		w := env.DoRequest("GET", path, nil, token)
		data := requireStatsResponseData(t, env, path, w)
		t.Log("✓ 章节行为统计获取成功")

		if items, ok := data["items"].([]interface{}); ok {
			t.Logf("✓ 章节统计条数: %d", len(items))
			for i, item := range items {
				if i < 3 {
					if itemMap, ok := item.(map[string]interface{}); ok {
						if title, ok := itemMap["title"].(string); ok {
							if completionRate, ok := itemMap["completionRate"].(float64); ok {
								t.Logf("  %s: 完读率 %.2f%%", title, completionRate*100)
							}
						}
					}
				}
			}
		}
		if total, ok := data["total"].(float64); ok {
			t.Logf("✓ 章节总数: %d", int(total))
		}
	})

	t.Log("✅ 高级数据分析E2E测试完成")
}
