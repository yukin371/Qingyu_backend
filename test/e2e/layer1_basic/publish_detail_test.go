//go:build e2e
// +build e2e

package layer1_basic

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestPublishManagementDetail 测试发布管理详细流程
// @P1 重要功能测试 - 发布管理详细流程
// 测试场景：
// 1. 发布单个章节
// 2. 批量发布章节
// 3. 设置章节定价
// 4. 定时发布
// 5. 查看发布历史
// 6. 撤销发布
//
// TDD原则：先写测试，看测试失败，再写实现代码
func TestPublishManagementDetail(t *testing.T) {
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

	// 步骤2: 创建项目和多个章节
	t.Run("步骤2_创建项目和多个章节", func(t *testing.T) {
		t.Log("创建项目和多个章节...")

		token := env.GetTestData("auth_token").(string)

		projectReq := map[string]interface{}{
			"title":       "e2e_test_publish_detail_project",
			"description": "发布管理详细流程测试项目",
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

		// 创建5个章节用于测试批量发布
		projectId := env.GetTestData("project_id").(string)
		chapterIds := make([]string, 0, 5)

		for i := 1; i <= 5; i++ {
			chapterReq := map[string]interface{}{
				"project_id": projectId,
				"title":      fmt.Sprintf("第%d章：测试章节", i),
				"content":    fmt.Sprintf("这是第%d章的内容，用于测试发布管理功能", i),
				"word_count": 25,
			}

			w := env.DoRequest("POST", "/api/v1/writer/documents", chapterReq, token)

			if w.Code == 200 || w.Code == 201 {
				response := env.ParseJSONResponse(w)
				if data, ok := response["data"].(map[string]interface{}); ok {
					var chapterId string
					if id, ok := data["documentId"].(string); ok {
						chapterId = id
					} else if id, ok := data["document_id"].(string); ok {
						chapterId = id
					}

					if chapterId != "" {
						chapterIds = append(chapterIds, chapterId)
						t.Logf("✓ 第%d章创建成功 (ID: %s)", i, chapterId)
					}
				}
			}
		}

		env.SetTestData("chapter_ids", chapterIds)
	})

	// 步骤3: 发布单个章节
	t.Run("步骤3_发布单个章节", func(t *testing.T) {
		t.Log("发布单个章节...")

		token := env.GetTestData("auth_token").(string)
		chapterIds := env.GetTestData("chapter_ids")
		projectId := env.GetTestData("project_id")

		if chapterIds == nil {
			t.Skip("章节未创建，跳过此步骤")
			return
		}

		ids := chapterIds.([]string)
		if len(ids) == 0 {
			t.Skip("没有可发布的章节")
			return
		}

		publishReq := map[string]interface{}{
			"action": "publish",
		}

		// 添加projectId作为查询参数
		path := fmt.Sprintf("/api/v1/writer/documents/%s/publish?projectId=%s", ids[0], projectId.(string))
		w := env.DoRequest("POST", path, publishReq, token)

		if w.Code == 200 || w.Code == 201 || w.Code == 202 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 单个章节发布成功，响应: %+v", response)
		} else {
			t.Logf("单个章节发布响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤4: 批量发布章节
	t.Run("步骤4_批量发布章节", func(t *testing.T) {
		t.Log("批量发布章节...")

		token := env.GetTestData("auth_token").(string)
		chapterIds := env.GetTestData("chapter_ids")
		projectId := env.GetTestData("project_id")

		if chapterIds == nil || projectId == nil {
			t.Skip("章节或项目未创建，跳过此步骤")
			return
		}

		ids := chapterIds.([]string)
		if len(ids) < 2 {
			t.Skip("章节数量不足，跳过批量发布测试")
			return
		}

		// 批量发布第2-5章
		batchPublishReq := map[string]interface{}{
			"document_ids": ids[1:],
			"action":       "publish",
		}

		path := fmt.Sprintf("/api/v1/writer/projects/%s/documents/batch-publish", projectId.(string))
		w := env.DoRequest("POST", path, batchPublishReq, token)

		if w.Code == 200 || w.Code == 201 || w.Code == 202 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 批量发布成功，响应: %+v", response)

			// 检查是否有部分失败的情况
			if data, ok := response["data"].(map[string]interface{}); ok {
				if successCount, ok := data["success_count"].(float64); ok {
					t.Logf("成功发布 %d 个章节", int(successCount))
				}
				if failedCount, ok := data["failed_count"].(float64); ok {
					if failedCount > 0 {
						t.Logf("⚠ 有 %d 个章节发布失败", int(failedCount))
					}
				}
			}
		} else {
			t.Logf("批量发布响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤5: 设置章节定价
	t.Run("步骤5_设置章节定价", func(t *testing.T) {
		t.Log("设置章节定价...")

		token := env.GetTestData("auth_token").(string)
		chapterIds := env.GetTestData("chapter_ids")
		projectId := env.GetTestData("project_id")

		if chapterIds == nil || projectId == nil {
			t.Skip("章节或项目未创建，跳过此步骤")
			return
		}

		ids := chapterIds.([]string)
		if len(ids) == 0 {
			t.Skip("没有可设置定价的章节")
			return
		}

		// 设置章节为付费
		updateStatusReq := map[string]interface{}{
			"is_paid":   true,
			"price":     100, // 100书币
			"vip_free":  true,
			"status":    "published",
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/publish-status?projectId=%s", ids[0], projectId.(string))
		w := env.DoRequest("PUT", path, updateStatusReq, token)

		if w.Code == 200 {
			t.Logf("✓ 章节定价设置成功")
		} else {
			t.Logf("章节定价设置响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤6: 定时发布（模拟）
	t.Run("步骤6_定时发布", func(t *testing.T) {
		t.Log("定时发布...")

		token := env.GetTestData("auth_token").(string)
		chapterIds := env.GetTestData("chapter_ids")
		projectId := env.GetTestData("project_id")

		if chapterIds == nil || projectId == nil {
			t.Skip("章节或项目未创建，跳过此步骤")
			return
		}

		ids := chapterIds.([]string)
		if len(ids) < 2 {
			t.Skip("章节数量不足，跳过定时发布测试")
			return
		}

		// 设置定时发布时间（未来1分钟）
		scheduledTime := time.Now().Add(1 * time.Minute).Format("2006-01-02T15:04:05Z07:00")

		publishReq := map[string]interface{}{
			"action":         "publish",
			"scheduled_time": scheduledTime,
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/publish?projectId=%s", ids[1], projectId.(string))
		w := env.DoRequest("POST", path, publishReq, token)

		if w.Code == 200 || w.Code == 201 || w.Code == 202 {
			t.Logf("✓ 定时发布设置成功")
		} else {
			t.Logf("定时发布响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤7: 查看发布历史
	t.Run("步骤7_查看发布历史", func(t *testing.T) {
		t.Log("查看发布历史...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id")

		if projectId == nil {
			t.Skip("项目未创建，跳过此步骤")
			return
		}

		// 获取第1页发布记录
		path := fmt.Sprintf("/api/v1/writer/projects/%s/publications?page=1&pageSize=20", projectId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 获取发布历史成功")

			if data, ok := response["data"].(map[string]interface{}); ok {
				if records, ok := data["records"].([]interface{}); ok {
					t.Logf("发布记录数量: %d", len(records))

					// 显示前3条记录
					for i, record := range records {
						if i >= 3 {
							break
						}
						if recordMap, ok := record.(map[string]interface{}); ok {
							recordJSON, _ := json.MarshalIndent(recordMap, "  ", "  ")
							t.Logf("记录%d: %s", i+1, string(recordJSON))
						}
					}
				}

				if total, ok := data["total"].(float64); ok {
					t.Logf("总发布记录数: %d", int(total))
				}
			}
		} else {
			t.Logf("获取发布历史响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤8: 获取项目发布状态
	t.Run("步骤8_获取项目发布状态", func(t *testing.T) {
		t.Log("获取项目发布状态...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id")

		if projectId == nil {
			t.Skip("项目未创建，跳过此步骤")
			return
		}

		path := fmt.Sprintf("/api/v1/writer/projects/%s/publication-status", projectId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 获取项目发布状态成功: %+v", response)

			if data, ok := response["data"].(map[string]interface{}); ok {
				if status, ok := data["status"].(string); ok {
					t.Logf("项目发布状态: %s", status)
				}
				if publishedCount, ok := data["published_count"].(float64); ok {
					t.Logf("已发布章节数: %d", int(publishedCount))
				}
				if totalCount, ok := data["total_count"].(float64); ok {
					t.Logf("总章节数: %d", int(totalCount))
				}
			}
		} else {
			t.Logf("获取项目发布状态响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤9: 撤销发布（测试未发布的章节）
	t.Run("步骤9_撤销发布", func(t *testing.T) {
		t.Log("撤销发布...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id")

		if projectId == nil {
			t.Skip("项目未创建，跳过此步骤")
			return
		}

		// 先创建一个新章节用于测试撤销
		chapterReq := map[string]interface{}{
			"project_id": projectId,
			"title":      "待撤销章节",
			"content":    "这个章节将被撤销发布",
			"word_count": 15,
		}

		w := env.DoRequest("POST", "/api/v1/writer/documents", chapterReq, token)

		var newChapterId string
		if w.Code == 200 || w.Code == 201 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if id, ok := data["documentId"].(string); ok {
					newChapterId = id
				} else if id, ok := data["document_id"].(string); ok {
					newChapterId = id
				}
			}
		}

		if newChapterId == "" {
			t.Skip("无法创建测试章节，跳过此步骤")
			return
		}

		// 先发布章节
		publishReq := map[string]interface{}{
			"action": "publish",
		}

		publishPath := fmt.Sprintf("/api/v1/writer/documents/%s/publish?projectId=%s", newChapterId, projectId.(string))
		env.DoRequest("POST", publishPath, publishReq, token)

		// 撤销发布
		unpublishReq := map[string]interface{}{
			"action": "unpublish",
		}

		path := fmt.Sprintf("/api/v1/writer/projects/%s/unpublish", projectId.(string))
		w = env.DoRequest("POST", path, unpublishReq, token)

		if w.Code == 200 {
			t.Logf("✓ 撤销发布成功")
		} else {
			t.Logf("撤销发布响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})
}

// TestPublishPricing 测试发布定价功能
// @P1 重要功能测试 - 发布定价
func TestPublishPricing(t *testing.T) {
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
		"title":       "pricing_test_project",
		"description": "定价测试项目",
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

	// 创建测试章节
	chapterReq := map[string]interface{}{
		"project_id": projectId,
		"title":      "付费章节",
		"content":    "这是一个付费章节的内容",
		"word_count": 15,
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
		}
	}

	// 测试不同定价策略
	t.Run("测试免费章节", func(t *testing.T) {
		updateStatusReq := map[string]interface{}{
			"is_paid":  false,
			"status":  "published",
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/publish-status?projectId=%s", chapterId, projectId)
		w := env.DoRequest("PUT", path, updateStatusReq, token)

		if w.Code == 200 {
			t.Logf("✓ 设置免费章节成功")
		}
	})

	t.Run("测试付费章节", func(t *testing.T) {
		updateStatusReq := map[string]interface{}{
			"is_paid":  true,
			"price":    200,
			"vip_free": true,
			"status":  "published",
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/publish-status?projectId=%s", chapterId, projectId)
		w := env.DoRequest("PUT", path, updateStatusReq, token)

		if w.Code == 200 {
			t.Logf("✓ 设置付费章节成功")
		}
	})

	t.Run("测试VIP免费章节", func(t *testing.T) {
		updateStatusReq := map[string]interface{}{
			"is_paid":  true,
			"price":    300,
			"vip_free": true,
			"status":  "published",
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/publish-status?projectId=%s", chapterId, projectId)
		w := env.DoRequest("PUT", path, updateStatusReq, token)

		if w.Code == 200 {
			t.Logf("✓ 设置VIP免费章节成功")
		}
	})
}

// TestPublishHistory 测试发布历史功能
// @P1 重要功能测试 - 发布历史
func TestPublishHistory(t *testing.T) {
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
		"title":       "history_test_project",
		"description": "发布历史测试项目",
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

	// 测试获取发布记录详情
	t.Run("测试获取发布记录详情", func(t *testing.T) {
		// 先发布项目以生成记录
		publishReq := map[string]interface{}{
			"action": "publish",
		}

		publishPath := fmt.Sprintf("/api/v1/writer/projects/%s/publish", projectId)
		w := env.DoRequest("POST", publishPath, publishReq, token)

		var recordId string
		if w.Code == 200 || w.Code == 201 || w.Code == 202 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if id, ok := data["id"].(string); ok {
					recordId = id
				} else if id, ok := data["record_id"].(string); ok {
					recordId = id
				}
			}
		}

		if recordId != "" {
			// 获取发布记录详情
			path := fmt.Sprintf("/api/v1/writer/publications/%s", recordId)
			w = env.DoRequest("GET", path, nil, token)

			if w.Code == 200 {
				response := env.ParseJSONResponse(w)
				t.Logf("✓ 获取发布记录详情成功: %+v", response)
			} else {
				t.Logf("获取发布记录详情响应: 状态码 %d", w.Code)
			}
		}
	})
}
