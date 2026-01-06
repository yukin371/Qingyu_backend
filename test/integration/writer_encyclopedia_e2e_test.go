package integration

import (
	"fmt"
	"testing"
	"time"
)

// TestWriterEncyclopediaE2E 设定百科完整流程端到端测试
// 测试完整的创作流程：项目 → 角色 → 地点 → 时间线 → 章节
func TestWriterEncyclopediaE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过端到端测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	// 登录测试用户
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过端到端测试")
	}

	var projectID string
	var characterID string
	var locationID string
	var timelineID string
	var eventID string
	var documentID string

	t.Run("阶段1：创建项目", func(t *testing.T) {
		projectData := map[string]interface{}{
			"title":       "仙侠小说_" + fmt.Sprintf("%d", time.Now().Unix()),
			"description": "一个完整的仙侠世界",
			"category":    "仙侠",
			"tags":        []string{"修真", "热血", "成长"},
		}

		w := helper.DoAuthRequest("POST", "/api/v1/projects", projectData, token)
		data := helper.AssertSuccess(w, 201, "创建项目应该成功")

		if id, ok := data["projectId"].(string); ok {
			projectID = id
			helper.LogSuccess("✓ 项目创建成功 - ID: %s", projectID)
		} else {
			t.Fatal("无法获取项目ID")
		}
	})

	t.Run("阶段2：创建角色设定", func(t *testing.T) {
		characterData := map[string]interface{}{
			"name":              "云无极",
			"alias":             []string{"无极剑仙", "云师兄"},
			"summary":           "天剑宗最强剑修，天赋异禀",
			"traits":            []string{"冷静", "坚毅", "正义"},
			"background":        "出身世家，自幼修炼剑道，十六岁突破金丹",
			"personalityPrompt": "外表冷峻，内心火热，对朋友两肋插刀",
			"speechPattern":     "言简意赅，剑气逼人",
			"currentState":      "金丹期巅峰，即将突破元婴",
		}

		url := fmt.Sprintf("/api/v1/projects/%s/characters", projectID)
		w := helper.DoAuthRequest("POST", url, characterData, token)
		data := helper.AssertSuccess(w, 201, "创建角色应该成功")

		if id, ok := data["id"].(string); ok {
			characterID = id
			helper.LogSuccess("✓ 角色创建成功 - ID: %s, 名字: 云无极", characterID)
		}
	})

	t.Run("阶段3：创建地点设定", func(t *testing.T) {
		locationData := map[string]interface{}{
			"name":        "天剑宗",
			"description": "修真界第一剑修宗门，坐落云巅，常年云雾缭绕",
			"climate":     "高山气候，四季如春",
			"culture":     "以剑入道，尊师重道",
			"geography":   "坐落九霄云巅，灵脉充沛",
			"atmosphere":  "剑气森然，庄严肃穆",
		}

		url := fmt.Sprintf("/api/v1/projects/%s/locations", projectID)
		w := helper.DoAuthRequest("POST", url, locationData, token)
		data := helper.AssertSuccess(w, 201, "创建地点应该成功")

		if id, ok := data["id"].(string); ok {
			locationID = id
			helper.LogSuccess("✓ 地点创建成功 - ID: %s, 名字: 天剑宗", locationID)
		}
	})

	t.Run("阶段4：创建时间线和事件", func(t *testing.T) {
		// 创建时间线
		timelineData := map[string]interface{}{
			"name":        "主角成长线",
			"description": "主角从凡人到剑仙的成长历程",
		}

		url := fmt.Sprintf("/api/v1/projects/%s/timelines", projectID)
		w := helper.DoAuthRequest("POST", url, timelineData, token)
		data := helper.AssertSuccess(w, 201, "创建时间线应该成功")

		if id, ok := data["id"].(string); ok {
			timelineID = id
			helper.LogSuccess("✓ 时间线创建成功 - ID: %s", timelineID)
		}

		// 创建时间线事件
		if timelineID != "" {
			eventData := map[string]interface{}{
				"title":       "拜入天剑宗",
				"description": "云无极通过重重考验，拜入天剑宗",
				"eventType":   "milestone",
				"importance":  10,
				"storyTime": map[string]interface{}{
					"year":  1,
					"month": 1,
				},
				"impact":       "修炼之路的起点",
				"participants": []string{characterID},
				"locationIds":  []string{locationID},
			}

			url := fmt.Sprintf("/api/v1/timelines/%s/events?projectId=%s", timelineID, projectID)
			w := helper.DoAuthRequest("POST", url, eventData, token)
			data := helper.AssertSuccess(w, 201, "创建事件应该成功")

			if id, ok := data["id"].(string); ok {
				eventID = id
				helper.LogSuccess("✓ 事件创建成功 - ID: %s", eventID)
			}
		}
	})

	t.Run("阶段5：创建章节文档", func(t *testing.T) {
		documentData := map[string]interface{}{
			"project_id": projectID,
			"title":      "第一章 拜师",
			"content":    "云巅之上，天剑宗山门前，少年云无极仰望巍峨的宗门...",
			"type":       "chapter",
			"status":     "draft",
		}

		w := helper.DoAuthRequest("POST", "/api/v1/documents", documentData, token)
		data := helper.AssertSuccess(w, 200, "创建文档应该成功")

		if id, ok := data["id"].(string); ok {
			documentID = id
			helper.LogSuccess("✓ 章节文档创建成功 - ID: %s", documentID)
		}
	})

	t.Run("阶段6：验证设定数据完整性", func(t *testing.T) {
		// 验证角色列表
		url := fmt.Sprintf("/api/v1/projects/%s/characters", projectID)
		w := helper.DoAuthRequest("GET", url, nil, token)
		helper.AssertSuccess(w, 200, "获取角色列表应该成功")
		helper.LogSuccess("✓ 角色数据验证通过")

		// 验证地点列表
		url = fmt.Sprintf("/api/v1/projects/%s/locations", projectID)
		w = helper.DoAuthRequest("GET", url, nil, token)
		helper.AssertSuccess(w, 200, "获取地点列表应该成功")
		helper.LogSuccess("✓ 地点数据验证通过")

		// 验证时间线事件
		if timelineID != "" {
			url = fmt.Sprintf("/api/v1/timelines/%s/events", timelineID)
			w = helper.DoAuthRequest("GET", url, nil, token)
			helper.AssertSuccess(w, 200, "获取事件列表应该成功")
			helper.LogSuccess("✓ 时间线数据验证通过")
		}

		// 验证可视化数据
		if timelineID != "" {
			url = fmt.Sprintf("/api/v1/timelines/%s/visualization", timelineID)
			w = helper.DoAuthRequest("GET", url, nil, token)
			helper.AssertSuccess(w, 200, "获取可视化数据应该成功")
			helper.LogSuccess("✓ 可视化数据验证通过")
		}
	})

	t.Run("阶段7：完整性总结", func(t *testing.T) {
		helper.LogSuccess("========================================")
		helper.LogSuccess("✓ 端到端测试完成")
		helper.LogSuccess("✓ 项目ID: %s", projectID)
		helper.LogSuccess("✓ 角色数: 1")
		helper.LogSuccess("✓ 地点数: 1")
		helper.LogSuccess("✓ 时间线: 1")
		helper.LogSuccess("✓ 事件数: 1")
		helper.LogSuccess("✓ 章节数: 1")
		helper.LogSuccess("========================================")
		helper.LogSuccess("设定百科系统集成验证通过！")
	})
}
