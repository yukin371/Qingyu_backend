package integration

import (
	"fmt"
	"testing"
	"time"
)

// TestTimelineScenario 时间线管理完整流程测试
func TestTimelineScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	// 登录测试用户
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过时间线流程测试")
	}

	var projectID string
	var timelineID string
	var event1ID string
	var event2ID string

	t.Run("1.创建测试项目", func(t *testing.T) {
		projectData := map[string]interface{}{
			"title":       "时间线测试项目_" + fmt.Sprintf("%d", time.Now().Unix()),
			"description": "用于测试时间线功能的项目",
			"category":    "玄幻",
		}

		w := helper.DoAuthRequest("POST", "/api/v1/writer/projects", projectData, token)
		data := helper.AssertSuccess(w, 201, "创建项目应该成功")

		if id, ok := data["projectId"].(string); ok {
			projectID = id
			helper.LogSuccess("项目创建成功 - ID: %s", projectID)
		}
	})

	if projectID != "" {
		t.Run("2.创建时间线", func(t *testing.T) {
			timelineData := map[string]interface{}{
				"name":        "主线剧情",
				"description": "故事主线发展时间线",
				"startTime": map[string]interface{}{
					"year": 1,
					"era":  "天宝元年",
				},
			}

			url := fmt.Sprintf("/api/v1/writer/projects/%s/timelines", projectID)
			w := helper.DoAuthRequest("POST", url, timelineData, token)
			data := helper.AssertSuccess(w, 201, "创建时间线应该成功")

			if id, ok := data["id"].(string); ok {
				timelineID = id
				name := data["name"]
				helper.LogSuccess("时间线创建成功 - ID: %s, 名字: %v", timelineID, name)
			}
		})

		t.Run("3.获取时间线列表", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/writer/projects/%s/timelines", projectID)
			w := helper.DoAuthRequest("GET", url, nil, token)
			_ = helper.AssertSuccess(w, 200, "获取时间线列表应该成功")
			helper.LogSuccess("时间线列表获取成功")
		})

		if timelineID != "" {
			t.Run("4.获取时间线详情", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/writer/timelines/%s?projectId=%s", timelineID, projectID)
				w := helper.DoAuthRequest("GET", url, nil, token)
				data := helper.AssertSuccess(w, 200, "获取时间线详情应该成功")

				if name, ok := data["name"].(string); ok && name == "主线剧情" {
					helper.LogSuccess("时间线详情获取成功 - 名字: %s", name)
				}
			})

			t.Run("5.创建第一个时间线事件", func(t *testing.T) {
				eventData := map[string]interface{}{
					"title":       "主角出生",
					"description": "在一个小村庄，主角降生",
					"eventType":   "character",
					"importance":  8,
					"storyTime": map[string]interface{}{
						"year":  1,
						"month": 1,
						"day":   1,
					},
					"impact": "主角的诞生，命运的开始",
				}

				url := fmt.Sprintf("/api/v1/writer/timelines/%s/events?projectId=%s", timelineID, projectID)
				w := helper.DoAuthRequest("POST", url, eventData, token)
				data := helper.AssertSuccess(w, 201, "创建事件应该成功")

				if id, ok := data["id"].(string); ok {
					event1ID = id
					title := data["title"]
					helper.LogSuccess("事件创建成功 - ID: %s, 标题: %v", event1ID, title)
				}
			})

			t.Run("6.创建第二个时间线事件", func(t *testing.T) {
				eventData := map[string]interface{}{
					"title":       "拜师学艺",
					"description": "主角遇到神秘高人，开始修炼之路",
					"eventType":   "milestone",
					"importance":  9,
					"storyTime": map[string]interface{}{
						"year":  16,
						"month": 3,
					},
					"impact": "踏上修炼之路，人生转折点",
				}

				url := fmt.Sprintf("/api/v1/writer/timelines/%s/events?projectId=%s", timelineID, projectID)
				w := helper.DoAuthRequest("POST", url, eventData, token)
				data := helper.AssertSuccess(w, 201, "创建事件应该成功")

				if id, ok := data["id"].(string); ok {
					event2ID = id
					helper.LogSuccess("事件创建成功 - ID: %s", event2ID)
				}
			})

			t.Run("7.获取时间线事件列表", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/writer/timelines/%s/events", timelineID)
				w := helper.DoAuthRequest("GET", url, nil, token)
				_ = helper.AssertSuccess(w, 200, "获取事件列表应该成功")
				helper.LogSuccess("事件列表获取成功")
			})

			if event1ID != "" {
				t.Run("8.获取事件详情", func(t *testing.T) {
					url := fmt.Sprintf("/api/v1/timeline-events/%s?projectId=%s", event1ID, projectID)
					w := helper.DoAuthRequest("GET", url, nil, token)
					data := helper.AssertSuccess(w, 200, "获取事件详情应该成功")

					if title, ok := data["title"].(string); ok && title == "主角出生" {
						helper.LogSuccess("事件详情获取成功 - 标题: %s", title)
					}
				})

				t.Run("9.更新事件信息", func(t *testing.T) {
					updateData := map[string]interface{}{
						"description": "在一个小村庄，主角降生，天生异象",
						"importance":  9,
					}

					url := fmt.Sprintf("/api/v1/timeline-events/%s?projectId=%s", event1ID, projectID)
					w := helper.DoAuthRequest("PUT", url, updateData, token)
					helper.AssertSuccess(w, 200, "更新事件应该成功")
					helper.LogSuccess("事件更新成功")
				})
			}

			t.Run("10.获取时间线可视化数据", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/writer/timelines/%s/visualization", timelineID)
				w := helper.DoAuthRequest("GET", url, nil, token)
				_ = helper.AssertSuccess(w, 200, "获取可视化数据应该成功")
				helper.LogSuccess("可视化数据获取成功")
			})

			if event2ID != "" {
				t.Run("11.删除事件", func(t *testing.T) {
					url := fmt.Sprintf("/api/v1/timeline-events/%s?projectId=%s", event2ID, projectID)
					w := helper.DoAuthRequest("DELETE", url, nil, token)
					helper.AssertSuccess(w, 200, "删除事件应该成功")
					helper.LogSuccess("事件删除成功")
				})
			}

			t.Run("12.删除时间线", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/writer/timelines/%s?projectId=%s", timelineID, projectID)
				w := helper.DoAuthRequest("DELETE", url, nil, token)
				helper.AssertSuccess(w, 200, "删除时间线应该成功")
				helper.LogSuccess("时间线删除成功")
			})
		}
	}
}
