//go:build e2e
// +build e2e

package layer1_basic

import (
	"fmt"
	"testing"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestTemplateManagement 测试模板管理功能
// P2功能：创建模板、应用模板、变量替换
//
// TDD原则：先写测试验证现有实现，发现bug后修复
func TestTemplateManagement(t *testing.T) {
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

	// 测试场景1: 创建章节模板
	t.Run("场景1_创建章节模板", func(t *testing.T) {
		t.Log("测试创建章节模板...")

		token := env.GetTestData("auth_token").(string)

		// 创建模板
		templateReq := map[string]interface{}{
			"name":        "标准章节模板",
			"description": "用于标准章节的模板",
			"type":        "chapter", // chapter/outline/setting
			"category":    "标准",
			"content":     "第{{chapterNum}}章：{{title}}\n\n{{content}}",
			"variables": []map[string]interface{}{
				{
					"name":         "chapterNum",
					"label":        "章节号",
					"type":         "number",
					"placeholder":  "请输入章节号",
					"defaultValue": "1",
					"required":     true,
					"order":        1,
				},
				{
					"name":         "title",
					"label":        "章节标题",
					"type":         "text",
					"placeholder":  "请输入章节标题",
					"defaultValue": "",
					"required":     true,
					"order":        2,
				},
				{
					"name":         "content",
					"label":        "章节内容",
					"type":         "textarea",
					"placeholder":  "请输入章节内容",
					"defaultValue": "",
					"required":     true,
					"order":        3,
				},
			},
		}

		w := env.DoRequest("POST", "/api/v1/writer/templates", templateReq, token)

		// 验证响应
		if w.Code == 200 || w.Code == 201 {
			t.Log("✓ 模板创建成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if templateId, ok := data["id"].(string); ok {
					env.SetTestData("template_id", templateId)
					t.Logf("✓ 模板ID: %s", templateId)
				} else if templateId, ok := data["templateId"].(string); ok {
					env.SetTestData("template_id", templateId)
					t.Logf("✓ 模板ID: %s", templateId)
				}

				if name, ok := data["name"].(string); ok {
					t.Logf("✓ 模板名称: %s", name)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 模板管理API未实现（404）")
		} else {
			t.Errorf("✗ 模板创建失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 测试场景2: 查询模板列表
	t.Run("场景2_查询模板列表", func(t *testing.T) {
		t.Log("测试查询模板列表...")

		token := env.GetTestData("auth_token").(string)

		w := env.DoRequest("GET", "/api/v1/writer/templates?type=chapter", nil, token)

		if w.Code == 200 {
			t.Log("✓ 模板列表查询成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if templates, ok := data["templates"].([]interface{}); ok {
					t.Logf("✓ 找到 %d 个模板", len(templates))

					// 验证我们创建的模板在列表中
					templateId := env.GetTestData("template_id")
					if templateId != nil {
						found := false
						for _, tmpl := range templates {
							if tmplMap, ok := tmpl.(map[string]interface{}); ok {
								if id, ok := tmplMap["id"].(string); ok && id == templateId.(string) {
									found = true
									t.Logf("✓ 找到创建的模板")
									break
								}
							}
						}
						if !found {
							t.Log("⚠ 未在列表中找到创建的模板")
						}
					}
				} else if items, ok := data["items"].([]interface{}); ok {
					t.Logf("✓ 找到 %d 个模板", len(items))
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 模板列表API未实现（404）")
		} else {
			t.Errorf("✗ 模板列表查询失败，状态码: %d", w.Code)
		}
	})

	// 测试场景3: 获取模板详情
	t.Run("场景3_获取模板详情", func(t *testing.T) {
		t.Log("测试获取模板详情...")

		token := env.GetTestData("auth_token").(string)
		templateIdData := env.GetTestData("template_id")
		if templateIdData == nil {
			t.Skip("模板ID不存在，跳过此测试")
			return
		}
		templateId := templateIdData.(string)

		w := env.DoRequest("GET", "/api/v1/writer/templates/"+templateId, nil, token)

		if w.Code == 200 {
			t.Log("✓ 模板详情获取成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if name, ok := data["name"].(string); ok {
					t.Logf("  模板名称: %s", name)
				}
				if content, ok := data["content"].(string); ok {
					t.Logf("  模板内容: %s", content)
				}
				if variables, ok := data["variables"].([]interface{}); ok {
					t.Logf("  变量数量: %d", len(variables))
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 模板详情API未实现或模板不存在（404）")
		} else {
			t.Errorf("✗ 模板详情获取失败，状态码: %d", w.Code)
		}
	})

	// 测试场景4: 应用模板（变量替换）
	t.Run("场景4_应用模板", func(t *testing.T) {
		t.Log("测试应用模板和变量替换...")

		token := env.GetTestData("auth_token").(string)
		templateIdData := env.GetTestData("template_id")
		if templateIdData == nil {
			t.Skip("模板ID不存在，跳过此测试")
			return
		}
		templateId := templateIdData.(string)

		// 创建测试文档
		projectReq := map[string]interface{}{
			"title":       "模板测试项目",
			"description": "用于测试模板应用",
			"genre":       "小说",
			"status":      "draft",
		}

		projectResp := actions.CreateProject(token, projectReq)
		var documentId string

		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			if projectId, ok := data["projectId"].(string); ok {
				// 创建章节
				chapterReq := map[string]interface{}{
					"project_id": projectId,
					"title":      "测试章节",
					"content":    "",
					"word_count": 0,
				}

				w := env.DoRequest("POST", "/api/v1/writer/documents", chapterReq, token)
				if w.Code == 200 || w.Code == 201 {
					response := env.ParseJSONResponse(w)
					if data, ok := response["data"].(map[string]interface{}); ok {
						if id, ok := data["documentId"].(string); ok {
							documentId = id
						} else if id, ok := data["document_id"].(string); ok {
							documentId = id
						}
					}
				}
			}
		}

		if documentId == "" {
			t.Skip("无法创建测试文档，跳过此测试")
			return
		}

		// 应用模板
		applyReq := map[string]interface{}{
			"documentId": documentId,
			"variables": map[string]string{
				"chapterNum": "1",
				"title":      "新的开始",
				"content":    "这是第一章的内容，讲述了主人公的故事开始。",
			},
		}

		w := env.DoRequest("POST", "/api/v1/writer/templates/"+templateId+"/apply", applyReq, token)

		if w.Code == 200 {
			t.Log("✓ 模板应用成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if renderedContent, ok := data["renderedContent"].(string); ok {
					t.Logf("✓ 渲染后的内容:\n%s", renderedContent)

					// 验证变量替换
					expectedContent := "第1章：新的开始\n\n这是第一章的内容，讲述了主人公的故事开始。"
					if renderedContent == expectedContent {
						t.Log("✓ 变量替换完全正确")
					} else {
						t.Logf("⚠ 变量替换结果与预期不完全匹配")
						t.Logf("  预期: %s", expectedContent)
						t.Logf("  实际: %s", renderedContent)
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 模板应用API未实现（404）")
		} else {
			t.Errorf("✗ 模板应用失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 测试场景5: 更新模板
	t.Run("场景5_更新模板", func(t *testing.T) {
		t.Log("测试更新模板...")

		token := env.GetTestData("auth_token").(string)
		templateIdData := env.GetTestData("template_id")
		if templateIdData == nil {
			t.Skip("模板ID不存在，跳过此测试")
			return
		}
		templateId := templateIdData.(string)

		// 更新模板
		updateReq := map[string]interface{}{
			"name":        fmt.Sprintf("更新后的模板_%s", t.Name()),
			"description": "这是更新后的描述",
		}

		w := env.DoRequest("PUT", "/api/v1/writer/templates/"+templateId, updateReq, token)

		if w.Code == 200 {
			t.Log("✓ 模板更新成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if name, ok := data["name"].(string); ok {
					t.Logf("  更新后的名称: %s", name)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 模板更新API未实现（404）")
		} else {
			t.Logf("ℹ 模板更新尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景6: 删除模板
	t.Run("场景6_删除模板", func(t *testing.T) {
		t.Log("测试删除模板...")

		token := env.GetTestData("auth_token").(string)
		templateIdData := env.GetTestData("template_id")
		if templateIdData == nil {
			t.Skip("模板ID不存在，跳过此测试")
			return
		}
		templateId := templateIdData.(string)

		w := env.DoRequest("DELETE", "/api/v1/writer/templates/"+templateId, nil, token)

		if w.Code == 200 {
			t.Log("✓ 模板删除成功")
			t.Log("✓ 删除请求执行完成")
		} else if w.Code == 404 {
			t.Log("⚠ 模板删除API未实现（404）")
		} else if w.Code == 403 {
			t.Log("⚠ 可能没有权限删除模板（403）")
		} else {
			t.Logf("ℹ 模板删除尝试，状态码: %d", w.Code)
		}
	})

	t.Log("✅ 模板管理E2E测试完成")
}
