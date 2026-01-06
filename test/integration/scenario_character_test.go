package integration

import (
	"fmt"
	"testing"
	"time"
)

// TestCharacterScenario 角色管理完整流程测试
func TestCharacterScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	// 登录测试用户
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过角色流程测试")
	}

	var projectID string
	var characterID1 string
	var characterID2 string
	var relationID string

	t.Run("1.创建测试项目", func(t *testing.T) {
		projectData := map[string]interface{}{
			"title":       "角色测试项目_" + fmt.Sprintf("%d", time.Now().Unix()),
			"description": "用于测试角色功能的项目",
			"category":    "玄幻",
		}

		w := helper.DoAuthRequest("POST", "/api/v1/projects", projectData, token)
		data := helper.AssertSuccess(w, 201, "创建项目应该成功")

		if id, ok := data["projectId"].(string); ok {
			projectID = id
			helper.LogSuccess("项目创建成功 - ID: %s", projectID)
		}
	})

	if projectID != "" {
		t.Run("2.创建第一个角色", func(t *testing.T) {
			characterData := map[string]interface{}{
				"name":              "李逍遥",
				"alias":             []string{"逍遥哥哥"},
				"summary":           "本作主人公，天真烂漫的少年",
				"traits":            []string{"勇敢", "善良", "乐观"},
				"background":        "余杭镇客栈小二，意外卷入灵儿被困仙灵岛的事件",
				"personalityPrompt": "性格开朗，行侠仗义，重情重义",
				"speechPattern":     "说话直爽，喜欢开玩笑",
			}

			url := fmt.Sprintf("/api/v1/projects/%s/characters", projectID)
			w := helper.DoAuthRequest("POST", url, characterData, token)
			data := helper.AssertSuccess(w, 201, "创建角色应该成功")

			if id, ok := data["id"].(string); ok {
				characterID1 = id
				name := data["name"]
				helper.LogSuccess("角色创建成功 - ID: %s, 名字: %v", characterID1, name)
			}
		})

		t.Run("3.创建第二个角色", func(t *testing.T) {
			characterData := map[string]interface{}{
				"name":    "赵灵儿",
				"alias":   []string{"灵儿"},
				"summary": "女娲后人，拥有强大的法力",
				"traits":  []string{"温柔", "善良", "坚强"},
			}

			url := fmt.Sprintf("/api/v1/projects/%s/characters", projectID)
			w := helper.DoAuthRequest("POST", url, characterData, token)
			data := helper.AssertSuccess(w, 201, "创建角色应该成功")

			if id, ok := data["id"].(string); ok {
				characterID2 = id
				helper.LogSuccess("角色创建成功 - ID: %s", characterID2)
			}
		})

		t.Run("4.获取角色列表", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/projects/%s/characters", projectID)
			w := helper.DoAuthRequest("GET", url, nil, token)
			_ = helper.AssertSuccess(w, 200, "获取角色列表应该成功")
			helper.LogSuccess("角色列表获取成功")
		})

		if characterID1 != "" {
			t.Run("5.获取角色详情", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/characters/%s?projectId=%s", characterID1, projectID)
				w := helper.DoAuthRequest("GET", url, nil, token)
				data := helper.AssertSuccess(w, 200, "获取角色详情应该成功")

				if name, ok := data["name"].(string); ok && name == "李逍遥" {
					helper.LogSuccess("角色详情获取成功 - 名字: %s", name)
				}
			})

			t.Run("6.更新角色信息", func(t *testing.T) {
				updateData := map[string]interface{}{
					"currentState": "已经踏上拯救灵儿的旅程",
				}

				url := fmt.Sprintf("/api/v1/characters/%s?projectId=%s", characterID1, projectID)
				w := helper.DoAuthRequest("PUT", url, updateData, token)
				helper.AssertSuccess(w, 200, "更新角色应该成功")
				helper.LogSuccess("角色更新成功")
			})
		}

		if characterID1 != "" && characterID2 != "" {
			t.Run("7.创建角色关系", func(t *testing.T) {
				relationData := map[string]interface{}{
					"fromId":   characterID1,
					"toId":     characterID2,
					"type":     "恋人",
					"strength": 90,
					"notes":    "命中注定的一对",
				}

				url := fmt.Sprintf("/api/v1/characters/relations?projectId=%s", projectID)
				w := helper.DoAuthRequest("POST", url, relationData, token)
				data := helper.AssertSuccess(w, 201, "创建关系应该成功")

				if id, ok := data["id"].(string); ok {
					relationID = id
					helper.LogSuccess("关系创建成功 - ID: %s", relationID)
				}
			})

			t.Run("8.获取角色关系图", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/projects/%s/characters/graph", projectID)
				w := helper.DoAuthRequest("GET", url, nil, token)
				_ = helper.AssertSuccess(w, 200, "获取关系图应该成功")
				helper.LogSuccess("关系图获取成功")
			})
		}

		if relationID != "" {
			t.Run("9.删除角色关系", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/characters/relations/%s?projectId=%s", relationID, projectID)
				w := helper.DoAuthRequest("DELETE", url, nil, token)
				helper.AssertSuccess(w, 200, "删除关系应该成功")
				helper.LogSuccess("关系删除成功")
			})
		}

		if characterID2 != "" {
			t.Run("10.删除角色", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/characters/%s?projectId=%s", characterID2, projectID)
				w := helper.DoAuthRequest("DELETE", url, nil, token)
				helper.AssertSuccess(w, 200, "删除角色应该成功")
				helper.LogSuccess("角色删除成功")
			})
		}
	}
}
