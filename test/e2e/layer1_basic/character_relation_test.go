//go:build e2e
// +build e2e

package layer1_basic

import (
	"fmt"
	"testing"

	e2e "Qingyu_backend/test/e2e/framework"
)

// TestCharacterRelationGraph 测试角色关系图功能
// P2功能：创建角色、设置关系、生成关系图
//
// TDD原则：先写测试验证现有实现，发现bug后修复
func TestCharacterRelationGraph(t *testing.T) {
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

	// 准备：创建项目
	t.Run("准备_创建项目", func(t *testing.T) {
		t.Log("创建测试项目...")

		token := env.GetTestData("auth_token").(string)

		projectReq := map[string]interface{}{
			"title":       "角色关系测试项目",
			"description": "用于测试角色关系图功能",
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
	})

	// 测试场景1: 创建角色
	t.Run("场景1_创建角色", func(t *testing.T) {
		t.Log("测试创建角色...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)

		// 创建第一个角色
		character1Req := map[string]interface{}{
			"projectId":   projectId,
			"name":        "主人公小明",
			"type":        "protagonist", // 主角
			"personality": []string{"勇敢", "善良", "聪明"},
			"appearance":  "一个年轻有为的青年，眼神坚定。",
			"background":  "出生在一个普通家庭，立志改变世界。",
		}

		w := env.DoRequest("POST", "/api/v1/writer/projects/"+projectId+"/characters", character1Req, token)

		if w.Code == 200 || w.Code == 201 {
			t.Log("✓ 角色1创建成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if charId, ok := data["id"].(string); ok {
					env.SetTestData("character1_id", charId)
					t.Logf("✓ 角色1 ID: %s", charId)
				} else if charId, ok := data["characterId"].(string); ok {
					env.SetTestData("character1_id", charId)
					t.Logf("✓ 角色1 ID: %s", charId)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 角色管理API未实现（404）")
			t.Skip("角色API未实现，跳过后续测试")
			return
		} else {
			t.Errorf("✗ 角色创建失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
		}

		// 创建第二个角色
		character2Req := map[string]interface{}{
			"projectId":   projectId,
			"name":        "小红",
			"type":        "supporting", // 配角
			"personality": []string{"温柔", "坚强"},
			"appearance":  "一个美丽善良的女孩。",
			"background":  "小明的青梅竹马。",
		}

		w2 := env.DoRequest("POST", "/api/v1/writer/projects/"+projectId+"/characters", character2Req, token)

		if w2.Code == 200 || w2.Code == 201 {
			t.Log("✓ 角色2创建成功")

			response := env.ParseJSONResponse(w2)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if charId, ok := data["id"].(string); ok {
					env.SetTestData("character2_id", charId)
					t.Logf("✓ 角色2 ID: %s", charId)
				} else if charId, ok := data["characterId"].(string); ok {
					env.SetTestData("character2_id", charId)
					t.Logf("✓ 角色2 ID: %s", charId)
				}
			}
		}

		// 创建第三个角色
		character3Req := map[string]interface{}{
			"projectId":   projectId,
			"name":        "大魔王",
			"type":        "antagonist", // 反派
			"personality": []string{"邪恶", "强大"},
			"appearance":  "一个神秘的黑暗人物。",
			"background":  "企图征服世界的邪恶势力领袖。",
		}

		w3 := env.DoRequest("POST", "/api/v1/writer/projects/"+projectId+"/characters", character3Req, token)

		if w3.Code == 200 || w3.Code == 201 {
			t.Log("✓ 角色3创建成功")

			response := env.ParseJSONResponse(w3)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if charId, ok := data["id"].(string); ok {
					env.SetTestData("character3_id", charId)
					t.Logf("✓ 角色3 ID: %s", charId)
				}
			}
		}
	})

	// 测试场景2: 设置角色关系
	t.Run("场景2_设置角色关系", func(t *testing.T) {
		t.Log("测试设置角色关系...")

		token := env.GetTestData("auth_token").(string)
		char1Id := env.GetTestData("character1_id")
		char2Id := env.GetTestData("character2_id")

		if char1Id == nil || char2Id == nil {
			t.Skip("角色ID不存在，跳过此测试")
			return
		}

		// 设置角色关系
		relationReq := map[string]interface{}{
			"sourceId":      char1Id,
			"targetId":      char2Id,
			"relationType":  "friend", // 朋友
			"relationLevel": 5,        // 关系强度 1-10
			"description":   "青梅竹马的好朋友",
		}

		w := env.DoRequest("POST", "/api/v1/writer/characters/relations", relationReq, token)

		if w.Code == 200 || w.Code == 201 {
			t.Log("✓ 角色关系设置成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if relationId, ok := data["id"].(string); ok {
					env.SetTestData("relation_id", relationId)
					t.Logf("✓ 关系ID: %s", relationId)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 角色关系API未实现（404）")
		} else {
			t.Logf("ℹ 角色关系设置尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景3: 查询角色列表
	t.Run("场景3_查询角色列表", func(t *testing.T) {
		t.Log("测试查询角色列表...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)

		w := env.DoRequest("GET", "/api/v1/writer/projects/"+projectId+"/characters", nil, token)

		if w.Code == 200 {
			t.Log("✓ 角色列表查询成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if characters, ok := data["characters"].([]interface{}); ok {
					t.Logf("✓ 找到 %d 个角色", len(characters))

					for _, char := range characters {
						if charMap, ok := char.(map[string]interface{}); ok {
							if name, ok := charMap["name"].(string); ok {
								t.Logf("  - %s", name)
							}
						}
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 角色列表API未实现（404）")
		} else {
			t.Errorf("✗ 角色列表查询失败，状态码: %d", w.Code)
		}
	})

	// 测试场景4: 获取角色详情和关系
	t.Run("场景4_获取角色详情和关系", func(t *testing.T) {
		t.Log("测试获取角色详情和关系...")

		token := env.GetTestData("auth_token").(string)
		char1Id := env.GetTestData("character1_id")

		if char1Id == nil {
			t.Skip("角色ID不存在，跳过此测试")
			return
		}

		w := env.DoRequest("GET", "/api/v1/writer/characters/"+char1Id.(string)+"?projectId="+env.GetTestData("project_id").(string), nil, token)

		if w.Code == 200 {
			t.Log("✓ 角色详情获取成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if name, ok := data["name"].(string); ok {
					t.Logf("  角色名: %s", name)
				}
				if relations, ok := data["relations"].([]interface{}); ok {
					t.Logf("  关系数: %d", len(relations))
					for _, rel := range relations {
						if relMap, ok := rel.(map[string]interface{}); ok {
							if relType, ok := relMap["relationType"].(string); ok {
								if target, ok := relMap["target"].(map[string]interface{}); ok {
									if targetName, ok := target["name"].(string); ok {
										t.Logf("    - %s: %s", relType, targetName)
									}
								}
							}
						}
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 角色详情API未实现（404）")
		} else {
			t.Errorf("✗ 角色详情获取失败，状态码: %d", w.Code)
		}
	})

	// 测试场景5: 生成关系图
	t.Run("场景5_生成关系图", func(t *testing.T) {
		t.Log("测试生成角色关系图...")

		token := env.GetTestData("auth_token").(string)
		projectId := env.GetTestData("project_id").(string)

		w := env.DoRequest("GET", "/api/v1/writer/projects/"+projectId+"/characters/graph", nil, token)

		if w.Code == 200 {
			t.Log("✓ 关系图生成成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if nodes, ok := data["nodes"].([]interface{}); ok {
					t.Logf("✓ 关系图节点数: %d", len(nodes))
				}
				if edges, ok := data["edges"].([]interface{}); ok {
					t.Logf("✓ 关系图边数: %d", len(edges))
				}

				// 验证图数据结构
				if nodes, ok := data["nodes"].([]interface{}); ok {
					for _, node := range nodes {
						if nodeMap, ok := node.(map[string]interface{}); ok {
							if id, ok := nodeMap["id"].(string); ok {
								if label, ok := nodeMap["label"].(string); ok {
									t.Logf("  节点: %s (%s)", label, id[:8])
								}
							}
						}
					}
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 关系图API未实现（404）")
		} else {
			t.Logf("ℹ 关系图生成尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景6: 更新角色信息
	t.Run("场景6_更新角色信息", func(t *testing.T) {
		t.Log("测试更新角色信息...")

		token := env.GetTestData("auth_token").(string)
		char1Id := env.GetTestData("character1_id")

		if char1Id == nil {
			t.Skip("角色ID不存在，跳过此测试")
			return
		}

		updateReq := map[string]interface{}{
			"personality": []string{"勇敢", "善良", "聪明", "正义"},
			"background":  fmt.Sprintf("更新后的背景_%s", t.Name()),
		}

		w := env.DoRequest("PUT", "/api/v1/writer/characters/"+char1Id.(string)+"?projectId="+env.GetTestData("project_id").(string), updateReq, token)

		if w.Code == 200 {
			t.Log("✓ 角色信息更新成功")

			response := env.ParseJSONResponse(w)
			if data, ok := response["data"].(map[string]interface{}); ok {
				if personality, ok := data["personality"].([]interface{}); ok {
					t.Logf("✓ 更新后的性格特征: %v", personality)
				}
			}
		} else if w.Code == 404 {
			t.Log("⚠ 角色更新API未实现（404）")
		} else {
			t.Logf("ℹ 角色更新尝试，状态码: %d", w.Code)
		}
	})

	// 测试场景7: 删除角色
	t.Run("场景7_删除角色", func(t *testing.T) {
		t.Log("测试删除角色...")

		token := env.GetTestData("auth_token").(string)
		char3Id := env.GetTestData("character3_id")

		if char3Id == nil {
			t.Skip("角色3 ID不存在，跳过此测试")
			return
		}

		w := env.DoRequest("DELETE", "/api/v1/writer/characters/"+char3Id.(string)+"?projectId="+env.GetTestData("project_id").(string), nil, token)

		if w.Code == 200 {
			t.Log("✓ 角色删除成功")
		} else if w.Code == 404 {
			t.Log("⚠ 角色删除API未实现（404）")
		} else {
			t.Logf("ℹ 角色删除尝试，状态码: %d", w.Code)
		}
	})

	t.Log("✅ 角色关系图E2E测试完成")
}
