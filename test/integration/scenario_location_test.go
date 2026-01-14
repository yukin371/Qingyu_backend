package integration

import (
	"fmt"
	"testing"
	"time"
)

// TestLocationScenario 地点管理完整流程测试
func TestLocationScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	helper := NewTestHelper(t, router)

	// 登录测试用户
	token := helper.LoginTestUser()
	if token == "" {
		t.Skip("无法登录测试用户，跳过地点流程测试")
	}

	var projectID string
	var continentID string
	var regionID string
	var cityID string
	var relationID string

	t.Run("1.创建测试项目", func(t *testing.T) {
		projectData := map[string]interface{}{
			"title":       "地点测试项目_" + fmt.Sprintf("%d", time.Now().Unix()),
			"description": "用于测试地点功能的项目",
			"category":    "仙侠",
		}

		w := helper.DoAuthRequest("POST", "/api/v1/writer/projects", projectData, token)
		data := helper.AssertSuccess(w, 201, "创建项目应该成功")

		if id, ok := data["projectId"].(string); ok {
			projectID = id
			helper.LogSuccess("项目创建成功 - ID: %s", projectID)
		}
	})

	if projectID != "" {
		t.Run("2.创建顶层地点-大陆", func(t *testing.T) {
			locationData := map[string]interface{}{
				"name":        "修真大陆",
				"description": "仙侠世界的主要大陆",
				"climate":     "多样化气候",
				"culture":     "修真文化盛行",
				"geography":   "地域辽阔，灵气充沛",
			}

			url := fmt.Sprintf("/api/v1/writer/projects/%s/locations", projectID)
			w := helper.DoAuthRequest("POST", url, locationData, token)
			data := helper.AssertSuccess(w, 201, "创建大陆应该成功")

			if id, ok := data["id"].(string); ok {
				continentID = id
				name := data["name"]
				helper.LogSuccess("大陆创建成功 - ID: %s, 名字: %v", continentID, name)
			}
		})

		if continentID != "" {
			t.Run("3.创建子地点-区域", func(t *testing.T) {
				locationData := map[string]interface{}{
					"name":        "东部仙域",
					"parentId":    continentID,
					"description": "修真大陆东部区域，宗门林立",
					"climate":     "四季分明，灵气浓郁",
				}

				url := fmt.Sprintf("/api/v1/writer/projects/%s/locations", projectID)
				w := helper.DoAuthRequest("POST", url, locationData, token)
				data := helper.AssertSuccess(w, 201, "创建区域应该成功")

				if id, ok := data["id"].(string); ok {
					regionID = id
					helper.LogSuccess("区域创建成功 - ID: %s", regionID)
				}
			})

			if regionID != "" {
				t.Run("4.创建三级地点-城市", func(t *testing.T) {
					locationData := map[string]interface{}{
						"name":        "天剑宗",
						"parentId":    regionID,
						"description": "东部第一剑修宗门",
						"atmosphere":  "严肃庄重，剑气森然",
					}

					url := fmt.Sprintf("/api/v1/writer/projects/%s/locations", projectID)
					w := helper.DoAuthRequest("POST", url, locationData, token)
					data := helper.AssertSuccess(w, 201, "创建城市应该成功")

					if id, ok := data["id"].(string); ok {
						cityID = id
						helper.LogSuccess("城市创建成功 - ID: %s", cityID)
					}
				})
			}
		}

		t.Run("5.获取地点列表", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/writer/projects/%s/locations", projectID)
			w := helper.DoAuthRequest("GET", url, nil, token)
			_ = helper.AssertSuccess(w, 200, "获取地点列表应该成功")
			helper.LogSuccess("地点列表获取成功")
		})

		t.Run("6.获取地点层级树", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/writer/projects/%s/locations/tree", projectID)
			w := helper.DoAuthRequest("GET", url, nil, token)
			_ = helper.AssertSuccess(w, 200, "获取层级树应该成功")
			helper.LogSuccess("层级树获取成功")
		})

		if continentID != "" {
			t.Run("7.获取地点详情", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/locations/%s?projectId=%s", continentID, projectID)
				w := helper.DoAuthRequest("GET", url, nil, token)
				data := helper.AssertSuccess(w, 200, "获取地点详情应该成功")

				if name, ok := data["name"].(string); ok && name == "修真大陆" {
					helper.LogSuccess("地点详情获取成功 - 名字: %s", name)
				}
			})

			t.Run("8.更新地点信息", func(t *testing.T) {
				updateData := map[string]interface{}{
					"description": "仙侠世界的主要大陆，灵气充沛，宗门林立",
				}

				url := fmt.Sprintf("/api/v1/locations/%s?projectId=%s", continentID, projectID)
				w := helper.DoAuthRequest("PUT", url, updateData, token)
				helper.AssertSuccess(w, 200, "更新地点应该成功")
				helper.LogSuccess("地点更新成功")
			})
		}

		if continentID != "" && regionID != "" {
			t.Run("9.创建地点关系", func(t *testing.T) {
				relationData := map[string]interface{}{
					"fromId":   continentID,
					"toId":     regionID,
					"type":     "contains",
					"distance": "无",
					"notes":    "大陆包含区域",
				}

				url := fmt.Sprintf("/api/v1/locations/relations?projectId=%s", projectID)
				w := helper.DoAuthRequest("POST", url, relationData, token)
				data := helper.AssertSuccess(w, 201, "创建关系应该成功")

				if id, ok := data["id"].(string); ok {
					relationID = id
					helper.LogSuccess("关系创建成功 - ID: %s", relationID)
				}
			})
		}

		if relationID != "" {
			t.Run("10.删除地点关系", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/locations/relations/%s?projectId=%s", relationID, projectID)
				w := helper.DoAuthRequest("DELETE", url, nil, token)
				helper.AssertSuccess(w, 200, "删除关系应该成功")
				helper.LogSuccess("关系删除成功")
			})
		}

		if cityID != "" {
			t.Run("11.删除地点", func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/locations/%s?projectId=%s", cityID, projectID)
				w := helper.DoAuthRequest("DELETE", url, nil, token)
				helper.AssertSuccess(w, 200, "删除地点应该成功")
				helper.LogSuccess("地点删除成功")
			})
		}
	}
}
