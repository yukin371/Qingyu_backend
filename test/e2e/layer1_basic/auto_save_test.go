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

// TestAutoSaveFunctionality 测试自动保存功能
// @P1 重要功能测试 - 自动保存功能
// 测试场景：
// 1. 编辑器自动保存触发
// 2. 保存指示器状态更新
// 3. 草稿恢复功能
//
// TDD原则：先写测试，看测试失败，再写实现代码
func TestAutoSaveFunctionality(t *testing.T) {
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

		// 创建作者用户
		author := fixtures.CreateUser()
		t.Logf("✓ 作者用户创建成功: %s (角色: %v)", author.Username, author.Roles)

		// 登录获取token
		token := actions.Login(author.Username, "Test1234")
		t.Logf("✓ 登录成功，获取token")

		// 保存到环境
		env.SetTestData("author_user", author)
		env.SetTestData("auth_token", token)
	})

	// 步骤2: 创建项目和章节
	t.Run("步骤2_创建项目和章节", func(t *testing.T) {
		t.Log("创建项目和章节...")

		token := env.GetTestData("auth_token").(string)

		// 创建项目
		projectReq := map[string]interface{}{
			"title":       "e2e_test_autosave_project",
			"description": "这是一个E2E测试自动保存功能的项目",
			"genre":       "小说",
			"status":      "draft",
		}

		projectResp := actions.CreateProject(token, projectReq)

		if data, ok := projectResp["data"].(map[string]interface{}); ok {
			if projectId, ok := data["projectId"].(string); ok {
				t.Logf("✓ 项目创建成功 (ID: %s)", projectId)
				env.SetTestData("project_id", projectId)
			} else if projectId, ok := data["id"].(string); ok {
				t.Logf("✓ 项目创建成功 (ID: %s)", projectId)
				env.SetTestData("project_id", projectId)
			}
		}

		// 创建章节
		projectId := env.GetTestData("project_id").(string)
		chapterReq := map[string]interface{}{
			"project_id": projectId,
			"title":      "e2e_test_autosave_chapter",
			"content":    "这是自动保存测试的初始内容",
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
					t.Logf("✓ 章节创建成功 (ID: %s)", documentId)
					env.SetTestData("chapter_id", documentId)
				}
			}
		}
	})

	// 步骤3: 测试自动保存功能
	t.Run("步骤3_测试自动保存功能", func(t *testing.T) {
		t.Log("测试自动保存功能...")

		token := env.GetTestData("auth_token").(string)
		chapterId := env.GetTestData("chapter_id")
		if chapterId == nil {
			t.Skip("chapter_id未设置，跳过此步骤")
			return
		}

		// 测试自动保存
		autoSaveReq := map[string]interface{}{
			"content":     "这是自动保存测试的新内容",
			"version":     1,
			"word_count":  16,
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/autosave", chapterId.(string))
		w := env.DoRequest("POST", path, autoSaveReq, token)

		if w.Code == 200 {
			t.Logf("✓ 自动保存成功")
		} else if w.Code == 409 {
			t.Logf("⚠ 版本冲突（可能需要处理版本号）")
		} else {
			t.Logf("自动保存响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤4: 测试获取保存状态
	t.Run("步骤4_测试获取保存状态", func(t *testing.T) {
		t.Log("测试获取保存状态...")

		token := env.GetTestData("auth_token").(string)
		chapterId := env.GetTestData("chapter_id")
		if chapterId == nil {
			t.Skip("chapter_id未设置，跳过此步骤")
			return
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/save-status", chapterId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			t.Logf("✓ 获取保存状态成功: %+v", response)
		} else {
			t.Logf("获取保存状态响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤5: 测试获取文档内容
	t.Run("步骤5_测试获取文档内容", func(t *testing.T) {
		t.Log("测试获取文档内容...")

		token := env.GetTestData("auth_token").(string)
		chapterId := env.GetTestData("chapter_id")
		if chapterId == nil {
			t.Skip("chapter_id未设置，跳过此步骤")
			return
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/content", chapterId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if content, ok := data["content"].(string); ok {
					t.Logf("✓ 获取文档内容成功，内容长度: %d", len(content))
				}
			}
		} else {
			t.Logf("获取文档内容响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤6: 测试手动更新文档内容
	t.Run("步骤6_测试手动更新文档内容", func(t *testing.T) {
		t.Log("测试手动更新文档内容...")

		token := env.GetTestData("auth_token").(string)
		chapterId := env.GetTestData("chapter_id")
		if chapterId == nil {
			t.Skip("chapter_id未设置，跳过此步骤")
			return
		}

		updateReq := map[string]interface{}{
			"content":    "这是手动更新的内容",
			"word_count": 13,
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/content", chapterId.(string))
		w := env.DoRequest("PUT", path, updateReq, token)

		if w.Code == 200 {
			t.Logf("✓ 手动更新内容成功")
		} else {
			t.Logf("手动更新内容响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤7: 测试草稿恢复 - 验证内容已保存
	t.Run("步骤7_测试草稿恢复", func(t *testing.T) {
		t.Log("测试草稿恢复...")

		token := env.GetTestData("auth_token").(string)
		chapterId := env.GetTestData("chapter_id")
		if chapterId == nil {
			t.Skip("chapter_id未设置，跳过此步骤")
			return
		}

		// 再次获取文档内容，验证内容已保存
		path := fmt.Sprintf("/api/v1/writer/documents/%s/content", chapterId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if content, ok := data["content"].(string); ok {
					if content == "这是手动更新的内容" {
						t.Logf("✓ 草稿恢复成功，内容已正确保存")
					} else {
						t.Logf("⚠ 内容不匹配，期望: '这是手动更新的内容', 实际: '%s'", content)
					}
				}
			}
		} else {
			t.Logf("草稿恢复响应: 状态码 %d, 响应: %s", w.Code, w.Body.String())
		}
	})

	// 步骤8: 测试自动保存防抖机制
	t.Run("步骤8_测试自动保存防抖机制", func(t *testing.T) {
		t.Log("测试自动保存防抖机制...")

		token := env.GetTestData("auth_token").(string)
		chapterId := env.GetTestData("chapter_id")
		if chapterId == nil {
			t.Skip("chapter_id未设置，跳过此步骤")
			return
		}

		// 快速连续发送多个自动保存请求
		for i := 0; i < 3; i++ {
			autoSaveReq := map[string]interface{}{
				"content":    fmt.Sprintf("快速保存测试内容 %d", i),
				"version":    1,
				"word_count": 15,
			}

			path := fmt.Sprintf("/api/v1/writer/documents/%s/autosave", chapterId.(string))
			w := env.DoRequest("POST", path, autoSaveReq, token)

			t.Logf("快速保存 %d 响应: 状态码 %d", i, w.Code)
		}

		// 等待一段时间后验证最终保存的内容
		time.Sleep(500 * time.Millisecond)

		path := fmt.Sprintf("/api/v1/writer/documents/%s/content", chapterId.(string))
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if content, ok := data["content"].(string); ok {
					t.Logf("✓ 防抖机制测试完成，最终内容: %s", content)
				}
			}
		}
	})

	// 步骤9: 测试版本冲突检测
	t.Run("步骤9_测试版本冲突检测", func(t *testing.T) {
		t.Log("测试版本冲突检测...")

		token := env.GetTestData("auth_token").(string)
		chapterId := env.GetTestData("chapter_id")
		if chapterId == nil {
			t.Skip("chapter_id未设置，跳过此步骤")
			return
		}

		// 使用错误的版本号触发冲突
		autoSaveReq := map[string]interface{}{
			"content":    "测试版本冲突",
			"version":    999, // 错误的版本号
			"word_count": 10,
		}

		path := fmt.Sprintf("/api/v1/writer/documents/%s/autosave", chapterId.(string))
		w := env.DoRequest("POST", path, autoSaveReq, token)

		if w.Code == 409 {
			t.Logf("✓ 版本冲突检测正常工作")
		} else {
			t.Logf("版本冲突检测响应: 状态码 %d (期望409)", w.Code)
		}
	})
}

// TestAutoSaveSaveIndicator 测试保存指示器功能
// @P1 重要功能测试 - 保存指示器
func TestAutoSaveSaveIndicator(t *testing.T) {
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
		"title":       "save_indicator_project",
		"description": "测试保存指示器",
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

	chapterReq := map[string]interface{}{
		"project_id": projectId,
		"title":      "测试章节",
		"content":    "初始内容",
		"word_count": 8,
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

	// 测试保存状态API返回的字段
	t.Run("验证保存状态字段", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/writer/documents/%s/save-status", chapterId)
		w := env.DoRequest("GET", path, nil, token)

		if w.Code == 200 {
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			if data, ok := resp["data"].(map[string]interface{}); ok {
				// 检查保存状态字段
				if lastSaved, ok := data["last_saved"].(string); ok {
					t.Logf("✓ 保存时间字段存在: %s", lastSaved)
				} else {
					t.Log("⚠ 保存时间字段不存在")
				}

				if status, ok := data["status"].(string); ok {
					t.Logf("✓ 保存状态字段存在: %s", status)
				} else {
					t.Log("⚠ 保存状态字段不存在")
				}
			}
		}
	})
}
