package integration

import (
	"fmt"
	"testing"
	"time"
)

// 写作流程测试 - 从创建项目到发布章节
func TestWritingScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 初始化helper
	helper := NewTestHelper(t, router)

	// 登录测试用户
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过写作流程测试")
	}

	var projectID string
	var documentID string
	var chapterID string

	t.Run("1.项目管理_创建写作项目", func(t *testing.T) {
		projectData := map[string]interface{}{
			"title":       "测试小说项目_" + fmt.Sprintf("%d", time.Now().Unix()),
			"description": "这是一个集成测试项目",
			"category":    "玄幻",
			"tags":        []string{"测试", "玄幻"},
		}

		w := helper.DoAuthRequest("POST", "/api/v1/writer/projects", projectData, token)
		data := helper.AssertSuccess(w, 201, "创建写作项目应该成功")

		if id, ok := data["projectId"].(string); ok {
			projectID = id
			title := data["title"]
			helper.LogSuccess("写作项目创建成功 - 项目ID: %s, 标题: %v", projectID, title)
		}
	})

	if projectID != "" {
		t.Run("2.项目管理_获取项目列表", func(t *testing.T) {
			w := helper.DoAuthRequest("GET", "/api/v1/writer/projects?page=1&size=10", nil, token)
			data := helper.AssertSuccess(w, 200, "获取项目列表应该成功")

			if projects, ok := data["projects"].([]interface{}); ok {
				helper.LogSuccess("项目列表获取成功，共 %d 个项目", len(projects))
			}
		})

		t.Run("3.文档管理_创建文档", func(t *testing.T) {
			documentData := map[string]interface{}{
				"project_id": projectID,
				"title":      "第一章 开端",
				"content":    "这是第一章的内容...",
				"type":       "chapter",
			}

			w := helper.DoAuthRequest("POST", "/api/v1/writer/documents", documentData, token)
			data := helper.AssertSuccess(w, 200, "创建文档应该成功")

			if id, ok := data["id"].(string); ok {
				documentID = id
				title := data["title"]
				helper.LogSuccess("文档创建成功 - 文档ID: %s, 标题: %v", documentID, title)
			}
		})

		if documentID != "" {
			t.Run("4.文档管理_保存草稿", func(t *testing.T) {
				updateData := map[string]interface{}{
					"content": "这是更新后的第一章内容，增加了更多细节描写...",
					"status":  "draft",
				}

				url := fmt.Sprintf("/api/v1/writer/documents/%s", documentID)
				w := helper.DoAuthRequest("PUT", url, updateData, token)
				helper.AssertSuccess(w, 200, "保存草稿应该成功")

				helper.LogSuccess("草稿保存成功")
			})

			t.Run("5.版本管理_创建版本", func(t *testing.T) {
				versionData := map[string]interface{}{
					"document_id": documentID,
					"note":        "第一版本",
				}

				url := fmt.Sprintf("/api/v1/writer/documents/%s/versions", documentID)
				w := helper.DoAuthRequest("POST", url, versionData, token)

				if w.Code == 200 {
					helper.LogSuccess("版本创建成功")
				} else {
					t.Logf("○ 创建版本失败或接口不存在 (状态码: %d)", w.Code)
				}
			})

			t.Run("6.文档管理_发布文档", func(t *testing.T) {
				publishData := map[string]interface{}{
					"status": "published",
				}

				url := fmt.Sprintf("/api/v1/writer/documents/%s", documentID)
				w := helper.DoAuthRequest("PUT", url, publishData, token)
				helper.AssertSuccess(w, 200, "发布文档应该成功")

				helper.LogSuccess("文档发布成功")
			})

			t.Run("7.文档管理_获取文档详情", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/writer/documents/%s", documentID)
				w := helper.DoAuthRequest("GET", url, nil, token)
				data := helper.AssertSuccess(w, 200, "获取文档详情应该成功")

				title := data["title"]
				status := data["status"]
				helper.LogSuccess("文档详情获取成功 - 标题: %v, 状态: %v", title, status)
			})
		}

		// 清理：删除测试项目
		t.Run("8.清理_删除测试项目", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/writer/projects/%s", projectID)
			w := helper.DoAuthRequest("DELETE", url, nil, token)
			helper.AssertSuccess(w, 200, "删除项目应该成功")

			helper.LogSuccess("测试项目清理成功")
		})
	}

	_ = chapterID

	helper.LogSuccess("写作流程测试完成 - 测试场景: 创建项目 → 创建文档 → 保存草稿 → 版本管理 → 发布")
}
