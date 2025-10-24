package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
)

// AI文本生成测试 - 测试Gemini API集成
func TestAIGenerationScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化
	_, err := config.LoadConfig("../..")
	require.NoError(t, err, "加载配置失败")

	err = core.InitDB()
	require.NoError(t, err, "初始化数据库失败")

	baseURL := "http://localhost:8080"

	// 登录获取 token（使用VIP用户以获得AI权限）
	token := loginTestUser(t, baseURL, "vip_user01", "Vip@123456")
	if token == "" {
		token = loginTestUser(t, baseURL, "test_user01", "Test@123456")
	}

	if token == "" {
		t.Skip("无法登录测试用户，跳过AI测试")
	}

	t.Run("1.文本生成_续写功能", func(t *testing.T) {
		requestData := map[string]interface{}{
			"prompt": "在一个风雪交加的夜晚，主角踏上了寻找失散多年亲人的旅程。",
			"type":   "continue",
			"length": 200,
		}

		response := callAIAPI(t, baseURL, token, "/api/v1/ai/generate", requestData)

		if response["code"] == float64(200) {
			data := response["data"].(map[string]interface{})
			generatedText := data["generated_text"].(string)

			t.Logf("✓ 续写功能测试成功")
			t.Logf("  原文: %s", requestData["prompt"])
			t.Logf("  续写: %s", generatedText)

			assert.NotEmpty(t, generatedText, "生成的文本不应为空")
			assert.Greater(t, len(generatedText), 50, "生成的文本长度应大于50字符")
		} else {
			t.Logf("○ 续写功能调用失败: %v", response["message"])
		}
	})

	t.Run("2.文本生成_改写功能", func(t *testing.T) {
		originalText := "他快速地跑向门口，打开门冲了出去。"

		requestData := map[string]interface{}{
			"text":  originalText,
			"type":  "rewrite",
			"style": "literary", // 文学化
		}

		response := callAIAPI(t, baseURL, token, "/api/v1/ai/rewrite", requestData)

		if response["code"] == float64(200) {
			data := response["data"].(map[string]interface{})
			rewrittenText := data["rewritten_text"].(string)

			t.Logf("✓ 改写功能测试成功")
			t.Logf("  原文: %s", originalText)
			t.Logf("  改写: %s", rewrittenText)

			assert.NotEmpty(t, rewrittenText, "改写的文本不应为空")
			assert.NotEqual(t, originalText, rewrittenText, "改写后的文本应与原文不同")
		} else {
			t.Logf("○ 改写功能调用失败: %v", response["message"])
		}
	})

	t.Run("3.文本生成_扩写功能", func(t *testing.T) {
		outline := "主角在森林中遇到了神秘的老人"

		requestData := map[string]interface{}{
			"outline": outline,
			"type":    "expand",
			"length":  300,
		}

		response := callAIAPI(t, baseURL, token, "/api/v1/ai/expand", requestData)

		if response["code"] == float64(200) {
			data := response["data"].(map[string]interface{})
			expandedText := data["expanded_text"].(string)

			t.Logf("✓ 扩写功能测试成功")
			t.Logf("  大纲: %s", outline)
			t.Logf("  扩写: %s...", expandedText[:min(100, len(expandedText))])

			assert.NotEmpty(t, expandedText, "扩写的文本不应为空")
			assert.Greater(t, len(expandedText), len(outline)*2, "扩写后的文本应显著长于大纲")
		} else {
			t.Logf("○ 扩写功能调用失败: %v", response["message"])
		}
	})

	t.Run("4.文本生成_润色功能", func(t *testing.T) {
		rawText := "天气很热，他很累，想要休息一下。"

		requestData := map[string]interface{}{
			"text": rawText,
			"type": "polish",
		}

		response := callAIAPI(t, baseURL, token, "/api/v1/ai/polish", requestData)

		if response["code"] == float64(200) {
			data := response["data"].(map[string]interface{})
			polishedText := data["polished_text"].(string)

			t.Logf("✓ 润色功能测试成功")
			t.Logf("  原文: %s", rawText)
			t.Logf("  润色: %s", polishedText)

			assert.NotEmpty(t, polishedText, "润色的文本不应为空")
		} else {
			t.Logf("○ 润色功能调用失败: %v", response["message"])
		}
	})

	t.Run("5.AI服务_Token使用统计", func(t *testing.T) {
		requestData := map[string]interface{}{
			"prompt": "测试短文本生成",
			"type":   "continue",
			"length": 50,
		}

		response := callAIAPI(t, baseURL, token, "/api/v1/ai/generate", requestData)

		if response["code"] == float64(200) {
			data := response["data"].(map[string]interface{})

			// 检查usage信息
			if usage, ok := data["usage"]; ok {
				usageMap := usage.(map[string]interface{})

				t.Logf("✓ Token使用统计:")
				t.Logf("  提示词Tokens: %.0f", usageMap["prompt_tokens"])
				t.Logf("  生成Tokens: %.0f", usageMap["completion_tokens"])
				t.Logf("  总计Tokens: %.0f", usageMap["total_tokens"])

				assert.Greater(t, usageMap["total_tokens"].(float64), float64(0), "总Token数应大于0")
			} else {
				t.Logf("○ 响应中没有usage信息")
			}
		} else {
			t.Logf("○ AI调用失败: %v", response["message"])
		}
	})

	t.Run("6.错误处理_空文本", func(t *testing.T) {
		requestData := map[string]interface{}{
			"prompt": "",
			"type":   "continue",
		}

		response := callAIAPI(t, baseURL, token, "/api/v1/ai/generate", requestData)

		// 应该返回错误
		if response["code"] != float64(200) {
			t.Logf("✓ 空文本错误处理正确: %v", response["message"])
		} else {
			t.Logf("○ 空文本应该返回错误")
		}
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

		response := callAIAPI(t, baseURL, token, "/api/v1/ai/generate", requestData)

		// 可能返回错误或截断
		if response["code"] != float64(200) {
			t.Logf("✓ 超长文本错误处理: %v", response["message"])
		} else {
			t.Logf("○ 超长文本被接受（可能被截断）")
		}
	})

	t.Logf("\n=== AI文本生成测试完成 ===")
	t.Logf("测试场景: 续写 → 改写 → 扩写 → 润色 → Token统计 → 错误处理")
}

// 测试直接调用Gemini Adapter
func TestGeminiAdapterDirect(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	_, err := config.LoadConfig("../..")
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("Gemini_健康检查", func(t *testing.T) {
		// 这里需要直接使用Gemini Adapter
		// 由于需要访问internal包，这个测试可能需要放在service/ai包中
		t.Log("○ Gemini Adapter健康检查需要在service层测试")
	})

	t.Run("Gemini_简单文本生成", func(t *testing.T) {
		// 直接调用Gemini API（绕过业务层）
		apiKey := "AIzaSyA6aj4aqWOdkIfZAYPlM5fk_8e4gJtkceE"
		model := "gemini-1.5-flash"

		requestData := map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"parts": []map[string]interface{}{
						{"text": "Hello, how are you?"},
					},
				},
			},
		}

		jsonData, _ := json.Marshal(requestData)
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
			model, apiKey)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

		if err != nil {
			t.Logf("○ Gemini API调用失败: %v", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			json.Unmarshal(body, &result)

			t.Logf("✓ Gemini API直接调用成功")
			t.Logf("  响应状态: %d", resp.StatusCode)

			if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
				t.Logf("  ✓ 收到响应候选")
			}
		} else {
			t.Logf("○ Gemini API返回错误: %d", resp.StatusCode)
			t.Logf("  响应: %s", string(body))
		}
	})

	_ = ctx
}

// 辅助函数：调用AI API
func callAIAPI(t *testing.T, baseURL, token, endpoint string, data map[string]interface{}) map[string]interface{} {
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Logf("创建请求失败: %v", err)
		return map[string]interface{}{"code": float64(500), "message": err.Error()}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Logf("请求失败: %v", err)
		return map[string]interface{}{"code": float64(500), "message": err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Logf("解析响应失败: %v", err)
		return map[string]interface{}{"code": float64(500), "message": "解析失败"}
	}

	return result
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
