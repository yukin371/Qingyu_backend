//go:build e2e
// +build e2e

package layer1_basic

import (
	"encoding/json"
	"fmt"
	"testing"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestOutlineManagement 测试大纲管理功能
// @P1 重要功能测试 - 大纲管理
// 测试场景：
// 1. 创建大纲
// 2. 编辑大纲
// 3. 大纲排序
// 4. 章节关联
//
// TDD原则：先写测试，看测试失败，再写实现代码
func TestOutlineManagement(t *testing.T) {
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

	// 步骤2: 创建项目
	t.Run("步骤2_创建项目", func(t *testing.T) {
		t.Log("创建项目...")

		token := env.GetTestData("auth_token").(string)

		projectReq := map[string]interface{}{
			"title":       "e2e_test_outline_project",
			"description": "大纲管理测试项目",
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
	})

	// 步骤3: 创建根级大纲节点（卷）
	t.Run("步骤3_创建根级大纲节点", func(t *testing.T) {
		t.Log("创建根级大纲节点...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)

		// 注意：大纲管理可能使用专门的API，这里先用文档API测试
		// 实际实现可能需要调整
		outlineReq := map[string]interface{}{
			"project_id": projectId,
			"title":      "第一卷：起源",
			"type":       "volume",
			"summary":    "这是故事的开端",
			"tension":    3,
			"order":      0,
		}

		w := env.DoRequest("POST", "/api/v1/writer/documents", outlineReq, token)

		if w.Code == 200 || w.Code == 201 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				var outlineId string
				if id, ok := data["documentId"].(string); ok {
					outlineId = id
				} else if id, ok := data["document_id"].(string); ok {
					outlineId = id
				}

				if outlineId != "" {
					t.Logf("✓ 根级大纲节点创建成功 (ID: %s)", outlineId)
					env.SetTestData("outline_root_id", outlineId)
				}
			}
		} else {
			t.Logf("创建根级大纲节点响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤4: 创建子级大纲节点（章）
	t.Run("步骤4_创建子级大纲节点", func(t *testing.T) {
		t.Log("创建子级大纲节点...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)
		rootId := env.GetTestData("outline_root_id")

		if rootId == nil {
			t.Skip("根大纲节点未创建，跳过此步骤")
			return
		}

		chapterOutlineReq := map[string]interface{}{
			"project_id": projectId,
			"parent_id":  rootId,
			"title":      "第一章：觉醒",
			"type":       "chapter",
			"summary":    "主角发现特殊能力",
			"tension":    5,
			"order":      0,
		}

		w := env.DoRequest("POST", "/api/v1/writer/documents", chapterOutlineReq, token)

		if w.Code == 200 || w.Code == 201 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				var chapterOutlineId string
				if id, ok := data["documentId"].(string); ok {
					chapterOutlineId = id
				} else if id, ok := data["document_id"].(string); ok {
					chapterOutlineId = id
				}

				if chapterOutlineId != "" {
					t.Logf("✓ 子级大纲节点创建成功 (ID: %s)", chapterOutlineId)
					env.SetTestData("outline_chapter_id", chapterOutlineId)
				}
			}
		} else {
			t.Logf("创建子级大纲节点响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤5: 测试大纲排序功能
	t.Run("步骤5_测试大纲排序功能", func(t *testing.T) {
		t.Log("测试大纲排序功能...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)

		// 创建多个大纲节点用于测试排序
		for i := 2; i <= 3; i++ {
			outlineReq := map[string]interface{}{
				"project_id": projectId,
				"title":      fmt.Sprintf("第%d章：新的冒险", i),
				"type":       "chapter",
				"order":      i - 1,
			}

			w := env.DoRequest("POST", "/api/v1/writer/documents", outlineReq, token)
			t.Logf("创建第%d个大纲节点，状态码: %d", i, w.Code)
		}

		// 测试批量排序API
		reorderReq := map[string]interface{}{
			"documents": []map[string]interface{}{
				{"id": env.GetTestData("outline_root_id"), "order": 0},
			},
		}

		path := fmt.Sprintf("/api/v1/writer/project/%s/documents/reorder", projectId)
		w := env.DoRequest("PUT", path, reorderReq, token)

		if w.Code == 200 {
			t.Logf("✓ 大纲排序成功")
		} else {
			t.Logf("大纲排序响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤6: 测试编辑大纲功能
	t.Run("步骤6_测试编辑大纲功能", func(t *testing.T) {
		t.Log("测试编辑大纲功能...")

		token := env.GetTestData("auth_token").(string)
		outlineId := env.GetTestData("outline_root_id")

		if outlineId == nil {
			t.Skip("大纲节点未创建，跳过此步骤")
			return
		}

		updateReq := map[string]interface{}{
			"title":      "第一卷：起源（修订版）",
			"summary":    "这是故事的开端，英雄的旅程从这里开始",
			"tension":    4,
			"characters": []string{"character_1", "character_2"},
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s", outlineId.(string))
		w := env.DoRequest("PUT", path, updateReq, token)

		if w.Code == 200 {
			t.Logf("✓ 大纲编辑成功")
		} else {
			t.Logf("大纲编辑响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤7: 测试大纲关联章节功能
	t.Run("步骤7_测试大纲关联章节功能", func(t *testing.T) {
		t.Log("测试大纲关联章节功能...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)
		outlineChapterId := env.GetTestData("outline_chapter_id")

		if outlineChapterId == nil {
			t.Skip("大纲节点未创建，跳过此步骤")
			return
		}

		// 创建实际章节
		chapterReq := map[string]interface{}{
			"project_id": projectId,
			"title":      "第一章：觉醒",
			"content":    "主角在一天清晨醒来，发现世界变得不同...",
			"word_count": 25,
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
					t.Logf("✓ 实际章节创建成功 (ID: %s)", chapterId)
					env.SetTestData("chapter_id", chapterId)
				}
			}
		}

		// 将大纲节点与章节关联
		if chapterId != "" {
			updateReq := map[string]interface{}{
				"chapter_id": chapterId,
			}

			path := fmt.Sprintf("/api/v1/writer/documents/%s", outlineChapterId.(string))
			w = env.DoRequest("PUT", path, updateReq, token)

			if w.Code == 200 {
				t.Logf("✓ 大纲与章节关联成功")
			} else {
				t.Logf("大纲与章节关联响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
			}
		}
	})

	// 步骤8: 测试获取文档树（大纲结构）
	t.Run("步骤8_测试获取文档树", func(t *testing.T) {
		t.Log("测试获取文档树...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)

		path := fmt.Sprintf("/api/v1/writer/project/%s/documents/tree", projectId)
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 获取文档树成功")
			if data, ok := response["data"].(map[string]interface{}); ok {
				treeJSON, _ := json.MarshalIndent(data, "", "  ")
				t.Logf("文档树结构: %s", string(treeJSON))
			}
		} else {
			t.Logf("获取文档树响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤9: 测试大纲属性（紧张度、类型等）
	t.Run("步骤9_测试大纲属性", func(t *testing.T) {
		t.Log("测试大纲属性...")

		token := env.GetTestData("auth_token").(string)
		outlineId := env.GetTestData("outline_root_id")

		if outlineId == nil {
			t.Skip("大纲节点未创建，跳过此步骤")
			return
		}

		// 测试设置大纲特有属性
		updateReq := map[string]interface{}{
			"summary":    "修订后的摘要：英雄之旅的起点",
			"tension":    7,         // 提高紧张度
			"type":       "英雄之旅-召唤", // 结构类型
			"characters": []string{"主角", "导师"},
			"items":      []string{"神秘物品"},
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s", outlineId.(string))
		w := env.DoRequest("PUT", path, updateReq, token)

		if w.Code == 200 {
			t.Logf("✓ 大纲属性设置成功")
		} else {
			t.Logf("大纲属性设置响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})
}

// TestOutlineStructure 测试大纲结构化功能
// @P1 重要功能测试 - 大纲结构
func TestOutlineStructure(t *testing.T) {
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
		"title":       "structure_test_project",
		"description": "大纲结构测试",
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

	// 测试创建层级大纲结构
	t.Run("创建层级大纲结构", func(t *testing.T) {
		// 卷 -> 幕 -> 章 -> 节
		structure := []struct {
			title    string
			typeVal  string
			parentId string
			order    int
		}{
			{"第一卷", "volume", "", 0},
			{"序幕", "prologue", "", 1},
			{"第一章", "chapter", "", 2},
			{"第一节", "section", "", 3},
		}

		var lastId string
		for i, item := range structure {
			req := map[string]interface{}{
				"project_id": projectId,
				"title":      item.title,
				"type":       item.typeVal,
				"order":      item.order,
			}

			w := env.DoRequest("POST", "/api/v1/writer/documents", req, token)

			if w.Code == 200 || w.Code == 201 {
				response := env.ParseJSONResponse(w)
				if data, ok := response["data"].(map[string]interface{}); ok {
					if id, ok := data["documentId"].(string); ok {
						lastId = id
					} else if id, ok := data["document_id"].(string); ok {
						lastId = id
					}
					t.Logf("创建第%d层: %s (ID: %s)", i+1, item.title, lastId)
				}
			}
		}
	})
}
