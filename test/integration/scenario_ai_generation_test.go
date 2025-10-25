package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// AI文本生成测试 - 测试DeepSeek API集成
func TestAIGenerationScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 设置测试环境
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 初始化helper
	helper := NewTestHelper(t, router)

	// 登录获取 token（优先使用VIP用户以获得AI权限）
	token := helper.LoginUser("vip_user01", "Vip@123456")
	if token == "" {
		token = helper.LoginTestUser() // 降级到普通用户
	}

	if token == "" {
		t.Skip("无法登录测试用户，跳过AI测试")
	}

	t.Logf("ℹ 本测试使用DeepSeek API (deepseek-chat模型)")

	t.Run("1.文本生成_续写功能", func(t *testing.T) {
		requestData := map[string]interface{}{
			"prompt": "在一个风雪交加的夜晚，主角踏上了寻找失散多年亲人的旅程。",
			"type":   "continue",
			"length": 200,
		}

		w := helper.DoAuthRequest("POST", "/api/v1/ai/generate", requestData, token)
		data := helper.AssertSuccess(w, 200, "续写功能应该成功")

		if generatedText, ok := data["generated_text"].(string); ok {
			t.Logf("  原文: %s", requestData["prompt"])
			t.Logf("  续写: %s", generatedText)

			assert.NotEmpty(t, generatedText, "生成的文本不应为空")
			assert.Greater(t, len(generatedText), 50, "生成的文本长度应大于50字符")

			helper.LogSuccess("续写功能测试成功 (DeepSeek)")
		}
	})

	t.Run("2.文本生成_改写功能", func(t *testing.T) {
		originalText := "他快速地跑向门口，打开门冲了出去。"

		requestData := map[string]interface{}{
			"text":  originalText,
			"type":  "rewrite",
			"style": "literary", // 文学化
		}

		w := helper.DoAuthRequest("POST", "/api/v1/ai/rewrite", requestData, token)
		data := helper.AssertSuccess(w, 200, "改写功能应该成功")

		if rewrittenText, ok := data["rewritten_text"].(string); ok {
			t.Logf("  原文: %s", originalText)
			t.Logf("  改写: %s", rewrittenText)

			assert.NotEmpty(t, rewrittenText, "改写的文本不应为空")
			assert.NotEqual(t, originalText, rewrittenText, "改写后的文本应与原文不同")

			helper.LogSuccess("改写功能测试成功 (DeepSeek)")
		}
	})

	t.Run("3.文本生成_扩写功能", func(t *testing.T) {
		outline := "主角在森林中遇到了神秘的老人"

		requestData := map[string]interface{}{
			"outline": outline,
			"type":    "expand",
			"length":  300,
		}

		w := helper.DoAuthRequest("POST", "/api/v1/ai/expand", requestData, token)
		data := helper.AssertSuccess(w, 200, "扩写功能应该成功")

		if expandedText, ok := data["expanded_text"].(string); ok {
			previewLen := 100
			if len(expandedText) < 100 {
				previewLen = len(expandedText)
			}
			t.Logf("  大纲: %s", outline)
			t.Logf("  扩写: %s...", expandedText[:previewLen])

			assert.NotEmpty(t, expandedText, "扩写的文本不应为空")
			assert.Greater(t, len(expandedText), len(outline)*2, "扩写后的文本应显著长于大纲")

			helper.LogSuccess("扩写功能测试成功 (DeepSeek)")
		}
	})

	t.Run("4.文本生成_润色功能", func(t *testing.T) {
		rawText := "天气很热，他很累，想要休息一下。"

		requestData := map[string]interface{}{
			"text": rawText,
			"type": "polish",
		}

		w := helper.DoAuthRequest("POST", "/api/v1/ai/polish", requestData, token)
		data := helper.AssertSuccess(w, 200, "润色功能应该成功")

		if polishedText, ok := data["polished_text"].(string); ok {
			t.Logf("  原文: %s", rawText)
			t.Logf("  润色: %s", polishedText)

			assert.NotEmpty(t, polishedText, "润色的文本不应为空")

			helper.LogSuccess("润色功能测试成功 (DeepSeek)")
		}
	})

	t.Run("5.AI服务_Token使用统计", func(t *testing.T) {
		requestData := map[string]interface{}{
			"prompt": "测试短文本生成",
			"type":   "continue",
			"length": 50,
		}

		w := helper.DoAuthRequest("POST", "/api/v1/ai/generate", requestData, token)
		data := helper.AssertSuccess(w, 200, "AI生成应该成功")

		// 检查usage信息
		if usage, ok := data["usage"].(map[string]interface{}); ok {
			t.Logf("  提示词Tokens: %.0f", usage["prompt_tokens"])
			t.Logf("  生成Tokens: %.0f", usage["completion_tokens"])
			t.Logf("  总计Tokens: %.0f", usage["total_tokens"])

			assert.Greater(t, usage["total_tokens"].(float64), float64(0), "总Token数应大于0")

			helper.LogSuccess("Token使用统计测试成功 (DeepSeek)")
		} else {
			t.Logf("○ 响应中没有usage信息")
		}
	})

	t.Run("6.错误处理_空文本", func(t *testing.T) {
		requestData := map[string]interface{}{
			"prompt": "",
			"type":   "continue",
		}

		w := helper.DoAuthRequest("POST", "/api/v1/ai/generate", requestData, token)
		helper.AssertError(w, 400, "空文本应该返回错误")

		helper.LogSuccess("空文本错误处理正确")
	})

	t.Run("7.错误处理_超长文本", func(t *testing.T) {
		// 生成一个非常长的文本
		longText := ""
		for i := 0; i < 10000; i++ {
			longText += "测试文本"
		}

		requestData := map[string]interface{}{
			"prompt": longText,
			"type":   "continue",
		}

		w := helper.DoAuthRequest("POST", "/api/v1/ai/generate", requestData, token)

		// 可能返回错误或截断
		if w.Code != 200 {
			t.Logf("✓ 超长文本错误处理: 状态码 %d", w.Code)
			helper.LogSuccess("超长文本错误处理正确")
		} else {
			t.Logf("○ 超长文本被接受（可能被截断）")
		}
	})

	helper.LogSuccess("AI文本生成测试完成 (DeepSeek) - 测试场景: 续写 → 改写 → 扩写 → 润色 → Token统计 → 错误处理")
}

// 测试AI服务健康检查
func TestAIServiceHealth(t *testing.T) {
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
		t.Skip("无法登录测试用户")
	}

	t.Run("AI服务_健康检查", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", "/api/v1/ai/health", nil, token)
		helper.AssertSuccess(w, 200, "AI服务健康检查应该成功")

		helper.LogSuccess("AI服务健康检查通过 (DeepSeek)")
	})

	t.Run("AI服务_获取提供商列表", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", "/api/v1/ai/providers", nil, token)
		data := helper.AssertSuccess(w, 200, "获取提供商列表应该成功")

		if providers, ok := data["providers"].([]interface{}); ok {
			helper.LogSuccess("AI提供商列表获取成功，共 %d 个提供商", len(providers))
		}
	})

	t.Run("AI服务_获取模型列表", func(t *testing.T) {
		w := helper.DoAuthRequest("GET", "/api/v1/ai/models", nil, token)
		data := helper.AssertSuccess(w, 200, "获取模型列表应该成功")

		if models, ok := data["models"].([]interface{}); ok {
			helper.LogSuccess("AI模型列表获取成功，共 %d 个模型", len(models))

			// 显示前3个模型
			for i := 0; i < len(models) && i < 3; i++ {
				if model, ok := models[i].(map[string]interface{}); ok {
					t.Logf("  - %v", model["name"])
				}
			}
		}
	})
}
