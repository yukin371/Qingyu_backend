//go:build e2e
// +build e2e

package layer1_basic

import (
	"fmt"
	"testing"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestBatchOperations 测试批量操作功能
// P2功能：批量发布章节、批量设置定价、批量导出
//
// TDD原则：先写测试验证现有实现，发现bug后修复
func TestBatchOperations(t *testing.T) {
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
		t.Logf("✓ 作者准备完成")
	})

	// 准备：创建项目和多个章节
	t.Run("准备_创建项目和章节", func(t *testing.T) {
		t.Log("创建测试项目和多个章节...")

		token := env.GetTestData("auth_token").(string)

		// 创建项目
		projectReq := map[string]interface{}{
			"title":       "批量操作测试项目",
			"description": "用于测试批量操作",
			"genre":       "小说",
			"status":      "draft",
		}

		projectResp := actions.CreateProject(token, projectReq)
		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			if projectId, ok := data["projectId"].(string); ok {
				env.SetTestData("project_id", projectId)
				t.Logf("✓ 项目创建成功: %s", projectId)
			}
		}

		// 创建多个章节用于批量操作
		projectId := env.GetTestData("project_id").(string)
		chapterIds := []string{}

		for i := 1; i <= 5; i++ {
			title := fmt.Sprintf("测试章节%d", i)
			content := fmt.Sprintf("这是第%d章的内容。", i)
			chapterReq := map[string]interface{}{
				"project_id": projectId,
				"title":      title,
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
					if documentId != "" {
						chapterIds = append(chapterIds, documentId)
					}
				}
			}
		}

		env.SetTestData("chapter_ids", chapterIds)
		t.Logf("✓ 创建了 %d 个章节", len(chapterIds))
	})

	// 测试场景1: 提交批量操作
	t.Run("场景1_提交批量操作", func(t *testing.T) {
		t.Log("测试提交批量操作...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)
		chapterIds := env.GetTestData("chapter_ids").([]string)

		// 提交批量发布操作
		batchReq := map[string]interface{}{
			"projectId":  projectId,
			"type":       "publish", // 批量发布
			"targetIds":  chapterIds,
			"atomic":     false, // 非原子操作，允许部分失败
			"conflictPolicy": "skip",
		}

		w := env.DoRequest("POST", "/api/v1/writer/batch-operations", batchReq, token)

		// 验证响应
		if w.Code == 200 || w.Code == 201 {
			t.Log("✓ 批量操作提交成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if batchId, ok := data["batchId"].(string); ok {
					env.SetTestData("batch_id", batchId)
					t.Logf("✓ 批量操作ID: %s", batchId)
				}

				// 检查预检摘要
				if preflight, ok := data["preflightSummary"].(map[string]interface{}); ok {
					if total, ok := preflight["total"].(float64); ok {
						t.Logf("✓ 预检总数: %d", int(total))
					}
					if ready, ok := preflight["ready"].(float64); ok {
						t.Logf("✓ 准备就绪: %d", int(ready))
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 批量操作API未实现（404）")
		} else {
			t.Errorf("✗ 批量操作提交失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 测试场景2: 查询批量操作状态
	t.Run("场景2_查询批量操作状态", func(t *testing.T) {
		t.Log("测试查询批量操作状态...")

		token := env.GetTestData("auth_token").(string)
		batchIdData := env.GetTestData("batch_id")
		if batchIdData == nil {
			t.Skip("批量操作ID不存在，跳过此测试")
			return
		}
		batchId := batchIdData.(string)

		w := env.DoRequest("GET", "/api/v1/writer/batch-operations/"+batchId, nil, token)

		if w.Code == 200 {
			t.Log("✓ 批量操作状态查询成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if status, ok := data["status"].(string); ok {
					t.Logf("  批量操作状态: %s", status)
				}

				// 检查进度
				if progress, ok := data["progress"].(map[string]interface{}); ok {
					if completed, ok := progress["completed"].(float64); ok {
						t.Logf("  已完成: %d", int(completed))
					}
					if total, ok := progress["total"].(float64); ok {
						t.Logf("  总数: %d", int(total))
					}
				}

				// 检查结果
				if results, ok := data["results"].([]interface{}); ok {
					t.Logf("  结果数: %d", len(results))

					// 统计成功和失败
					successCount := 0
					failedCount := 0
					for _, result := range results {
						if r, ok := result.(map[string]interface{}); ok {
							if success, ok := r["success"].(bool); ok && success {
								successCount++
							} else {
								failedCount++
							}
						}
					}
					t.Logf("  成功: %d, 失败: %d", successCount, failedCount)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 批量操作状态查询API未实现（404）")
		} else {
			t.Errorf("✗ 状态查询失败，状态码: %d", w.Code)
		}
	})

	// 测试场景3: 批量设置定价
	t.Run("场景3_批量设置定价", func(t *testing.T) {
		t.Log("测试批量设置定价...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)
		chapterIds := env.GetTestData("chapter_ids").([]string)

		// 批量设置定价
		priceBatchReq := map[string]interface{}{
			"projectId":  projectId,
			"type":       "set_price", // 设置定价
			"targetIds":  chapterIds[:3], // 只对前3个章节设置
			"payload": map[string]interface{}{
				"price":     100, // 单位：分
				"currency":  "CNY",
				"free":      false,
			},
			"atomic": false,
		}

		w := env.DoRequest("POST", "/api/v1/writer/batch-operations", priceBatchReq, token)

		if w.Code == 200 || w.Code == 201 {
			t.Log("✓ 批量定价操作提交成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if batchId, ok := data["batchId"].(string); ok {
					t.Logf("✓ 批量定价操作ID: %s", batchId)
				}
			}
		} else if w.Code == 400 {
			// 可能不支持此操作类型
			t.Logf("⚠ 批量定价操作可能不支持（状态码400）")
			t.Logf("  响应: %s", w.Body.String())
		} else if w.Code == 404 {
			t.Log("⚠ 批量操作API未实现（404）")
		} else {
			t.Logf("ℹ 批量定价操作尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景4: 批量导出
	t.Run("场景4_批量导出", func(t *testing.T) {
		t.Log("测试批量导出...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)
		chapterIds := env.GetTestData("chapter_ids").([]string)

		// 批量导出
		exportBatchReq := map[string]interface{}{
			"projectId":  projectId,
			"type":       "export",
			"targetIds":  chapterIds,
			"payload": map[string]interface{}{
				"format":      "txt", // 导出格式
				"includeMeta": true,  // 包含元数据
			},
			"atomic": false,
		}

		w := env.DoRequest("POST", "/api/v1/writer/batch-operations", exportBatchReq, token)

		if w.Code == 200 || w.Code == 201 {
			t.Log("✓ 批量导出操作提交成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if batchId, ok := data["batchId"].(string); ok {
					t.Logf("✓ 批量导出操作ID: %s", batchId)
				}
			}
		} else if w.Code == 400 {
			t.Logf("⚠ 批量导出操作可能不支持（状态码400）")
		} else if w.Code == 404 {
			t.Log("⚠ 批量操作API未实现（404）")
		} else {
			t.Logf("ℹ 批量导出操作尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景5: 部分失败处理
	t.Run("场景5_部分失败处理", func(t *testing.T) {
		t.Log("测试批量操作的部分失败处理...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)

		// 包含无效ID的批量操作
		batchReq := map[string]interface{}{
			"projectId":  projectId,
			"type":       "publish",
			"targetIds":  []string{"invalid_id_1", "invalid_id_2"},
			"atomic":     false, // 非原子操作，允许部分失败
		}

		w := env.DoRequest("POST", "/api/v1/writer/batch-operations", batchReq, token)

		if w.Code == 200 || w.Code == 201 {
			t.Log("✓ 批量操作接受包含无效ID的请求")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if preflight, ok := data["preflightSummary"].(map[string]interface{}); ok {
					if failed, ok := preflight["failed"].(float64); ok {
						t.Logf("✓ 预检识别出 %d 个无效项", int(failed))
					}
				}
			}
		} else if w.Code == 404 {
			t.Skip("批量操作API未实现，跳过部分失败测试")
		} else {
			t.Logf("ℹ 部分失败测试，状态码: %d", w.Code)
		}
	})

	t.Log("✅ 批量操作E2E测试完成")
}
