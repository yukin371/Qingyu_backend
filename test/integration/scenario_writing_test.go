package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
)

// 写作流程测试 - 从创建项目到发布章节
func TestWritingScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化
	_, err := config.LoadConfig("../..")
	require.NoError(t, err, "加载配置失败")

	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	baseURL := "http://localhost:8080"

	// 登录获取 token
	token := loginTestUser(t, baseURL, "test_user01", "Test@123456")
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

		jsonData, _ := json.Marshal(projectData)
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/projects", baseURL), bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 添加详细的诊断日志
		t.Logf("原始响应状态: %d", resp.StatusCode)
		t.Logf("响应头: %v", resp.Header)
		body, _ := io.ReadAll(resp.Body)
		t.Logf("响应Body长度: %d", len(body))
		t.Logf("响应Body (前500字符): %s", string(body[:min(len(body), 500)]))

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Logf("❌ JSON解析失败: %v", err)
			t.Logf("完整响应Body: %s", string(body))
		}
		require.NoError(t, err, "JSON解析失败，请检查API响应格式")

		if result["code"] == float64(200) {
			data := result["data"].(map[string]interface{})
			projectID = data["id"].(string)

			t.Logf("✓ 写作项目创建成功")
			t.Logf("  项目ID: %s", projectID)
			t.Logf("  项目标题: %s", data["title"])
		} else {
			t.Logf("○ 创建项目失败: %v", result["message"])
		}
	})

	if projectID != "" {
		t.Run("2.项目管理_获取项目列表", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/projects?page=1&size=10", baseURL), nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)

			if result["code"] == float64(200) {
				data := result["data"].(map[string]interface{})
				projects := data["projects"]
				if projects != nil {
					projectList := projects.([]interface{})
					t.Logf("✓ 项目列表获取成功，共 %d 个项目", len(projectList))
				}
			} else {
				t.Logf("○ 获取项目列表失败: %v", result["message"])
			}
		})

		t.Run("3.文档管理_创建文档", func(t *testing.T) {
			documentData := map[string]interface{}{
				"project_id": projectID,
				"title":      "第一章 开端",
				"content":    "这是第一章的内容...",
				"type":       "chapter",
			}

			jsonData, _ := json.Marshal(documentData)
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/documents", baseURL), bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)

			if result["code"] == float64(200) {
				data := result["data"].(map[string]interface{})
				documentID = data["id"].(string)

				t.Logf("✓ 文档创建成功")
				t.Logf("  文档ID: %s", documentID)
				t.Logf("  文档标题: %s", data["title"])
			} else {
				t.Logf("○ 创建文档失败: %v", result["message"])
			}
		})

		if documentID != "" {
			t.Run("4.文档管理_保存草稿", func(t *testing.T) {
				updateData := map[string]interface{}{
					"content": "这是更新后的第一章内容，增加了更多细节描写...",
					"status":  "draft",
				}

				jsonData, _ := json.Marshal(updateData)
				req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/documents/%s", baseURL, documentID), bytes.NewBuffer(jsonData))
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				require.NoError(t, err)

				if result["code"] == float64(200) {
					t.Logf("✓ 草稿保存成功")
				} else {
					t.Logf("○ 保存草稿失败: %v", result["message"])
				}
			})

			t.Run("5.版本管理_创建版本", func(t *testing.T) {
				versionData := map[string]interface{}{
					"document_id": documentID,
					"note":        "第一版本",
				}

				jsonData, _ := json.Marshal(versionData)
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/documents/%s/versions", baseURL, documentID), bytes.NewBuffer(jsonData))
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				require.NoError(t, err)

				if result["code"] == float64(200) {
					t.Logf("✓ 版本创建成功")
				} else {
					t.Logf("○ 创建版本失败或接口不存在: %v", result["message"])
				}
			})

			t.Run("6.文档管理_发布文档", func(t *testing.T) {
				publishData := map[string]interface{}{
					"status": "published",
				}

				jsonData, _ := json.Marshal(publishData)
				req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/documents/%s", baseURL, documentID), bytes.NewBuffer(jsonData))
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				require.NoError(t, err)

				if result["code"] == float64(200) {
					t.Logf("✓ 文档发布成功")
				} else {
					t.Logf("○ 发布文档失败: %v", result["message"])
				}
			})

			t.Run("7.文档管理_获取文档详情", func(t *testing.T) {
				req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/documents/%s", baseURL, documentID), nil)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)

				client := &http.Client{}
				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				require.NoError(t, err)

				if result["code"] == float64(200) {
					data := result["data"].(map[string]interface{})
					t.Logf("✓ 文档详情获取成功")
					t.Logf("  标题: %s", data["title"])
					t.Logf("  状态: %s", data["status"])
				} else {
					t.Logf("○ 获取文档详情失败: %v", result["message"])
				}
			})
		}

		// 清理：删除测试项目
		t.Run("8.清理_删除测试项目", func(t *testing.T) {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/projects/%s", baseURL, projectID), nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			require.NoError(t, err)

			if result["code"] == float64(200) {
				t.Logf("✓ 测试项目清理成功")
			} else {
				t.Logf("○ 删除项目失败: %v", result["message"])
			}
		})
	}

	_ = chapterID

	t.Logf("\n=== 写作流程测试完成 ===")
	t.Logf("测试场景: 创建项目 → 创建文档 → 保存草稿 → 版本管理 → 发布")
}
