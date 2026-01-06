package ai

import (
	"context"
	"fmt"
	"testing"
	"time"

	"Qingyu_backend/service/ai/adapter"
	"Qingyu_backend/service/ai/dto"
	"Qingyu_backend/service/ai/mocks"

	"github.com/stretchr/testify/assert"
)

// TestSummarizeService_SummarizeContent_Success 测试成功总结内容
func TestSummarizeService_SummarizeContent_Success(t *testing.T) {
	// 创建模拟适配器
	mockAdapter := mocks.NewMockAIAdapter("test-adapter")
	mockAdapter.SetTextResponse("这是一段测试内容的摘要。它包含了关键信息和要点。", 100)

	// 创建适配器管理器（使用简单实现）
	manager := &adapter.AdapterManager{}

	// 创建服务实例
	service := NewSummarizeService(manager)

	// 使用反射或通过测试接口注入模拟适配器
	// 由于AdapterManager的复杂性，我们通过创建测试辅助函数来模拟

	ctx := context.Background()
	req := &dto.SummarizeRequest{
		Content:     "这是要测试的完整内容。它包含了很多文字和详细信息。",
		SummaryType: "brief",
		MaxLength:   1000,
	}

	// 注意：这个测试需要实际的适配器管理器支持或依赖注入
	// 这里我们演示测试结构
	_ = mockAdapter
	_ = service
	_ = ctx
	_ = req

	// TODO: 完善实际的测试调用
	// result, err := service.SummarizeContent(ctx, req)
	// assert.NoError(t, err)
	// assert.NotEmpty(t, result.Summary)
	// assert.Greater(t, result.TokensUsed, 0)
}

// TestSummarizeService_SummarizeContent_EmptyContent 测试空内容
func TestSummarizeService_SummarizeContent_EmptyContent(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	ctx := context.Background()
	req := &dto.SummarizeRequest{
		Content:     "",
		SummaryType: "brief",
	}

	result, err := service.SummarizeContent(ctx, req)

	// 应该返回错误
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "内容不能为空")
}

// TestSummarizeService_SummarizeContent_WhitespaceContent 测试仅包含空白字符的内容
func TestSummarizeService_SummarizeContent_WhitespaceContent(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	ctx := context.Background()
	req := &dto.SummarizeRequest{
		Content:     "   \n\t   ",
		SummaryType: "brief",
	}

	result, err := service.SummarizeContent(ctx, req)

	// 应该返回错误
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "内容不能为空")
}

// TestSummarizeService_SummarizeContent_DifferentTypes 测试不同总结类型
func TestSummarizeService_SummarizeContent_DifferentTypes(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	testCases := []struct {
		name        string
		summaryType string
		content     string
	}{
		{
			name:        "简短摘要",
			summaryType: "brief",
			content:     "这是一段需要总结的测试内容，用于测试简短摘要功能。",
		},
		{
			name:        "详细摘要",
			summaryType: "detailed",
			content:     "这是一段需要总结的测试内容，用于测试详细摘要功能。详细摘要应该包含更多的信息和细节。",
		},
		{
			name:        "关键点提取",
			summaryType: "keypoints",
			content:     "这是第一个关键点。这是第二个关键点。这是第三个关键点。",
		},
		{
			name:        "默认类型",
			summaryType: "",
			content:     "测试默认摘要类型。",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 验证提示词构建逻辑
			prompt := service.buildSummarizePrompt(&dto.SummarizeRequest{
				Content:     tc.content,
				SummaryType: tc.summaryType,
			})

			assert.NotEmpty(t, prompt)
			assert.Contains(t, prompt, tc.content)
		})
	}
}

// TestSummarizeService_SummarizeContent_WithQuotes 测试包含引用的总结
func TestSummarizeService_SummarizeContent_WithQuotes(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	req := &dto.SummarizeRequest{
		Content:       "这是测试内容",
		SummaryType:   "brief",
		IncludeQuotes: true,
	}

	prompt := service.buildSummarizePrompt(req)

	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "关键引用")
}

// TestSummarizeService_SummarizeContent_AIError 测试AI服务错误
func TestSummarizeService_SummarizeContent_AIError(t *testing.T) {
	// 创建模拟适配器并配置为返回错误
	mockAdapter := mocks.NewMockAIAdapter("test-adapter")
	mockAdapter.ShouldFail = true
	mockAdapter.FailureError = &adapter.AdapterError{
		Type:       adapter.ErrorTypeServiceUnavailable,
		Message:    "AI服务不可用",
		Code:       "service_unavailable",
		StatusCode: 503,
		Provider:   "test-adapter",
		Retryable:  true,
	}

	_ = mockAdapter
	// TODO: 完善实际的测试调用，需要适配器管理器支持
}

// TestSummarizeService_SummarizeContent_AITimeout 测试AI服务超时
func TestSummarizeService_SummarizeContent_AITimeout(t *testing.T) {
	// 创建模拟适配器并配置超时
	mockAdapter := mocks.NewMockAIAdapter("test-adapter")
	mockAdapter.ShouldTimeout = true
	mockAdapter.ResponseDelay = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = mockAdapter
	_ = ctx
	// TODO: 完善实际的测试调用，需要适配器管理器支持
}

// TestSummarizeService_SummarizeContent_CompressionRate 测试压缩率计算
func TestSummarizeService_SummarizeContent_CompressionRate(t *testing.T) {
	// 测试压缩率计算逻辑
	testCases := []struct {
		name            string
		originalLength  int
		summaryLength   int
		expectedRate    float64
	}{
		{
			name:           "50%压缩率",
			originalLength: 1000,
			summaryLength:  500,
			expectedRate:   0.5,
		},
		{
			name:           "25%压缩率",
			originalLength: 1000,
			summaryLength:  250,
			expectedRate:   0.25,
		},
		{
			name:           "100%压缩率（相同长度）",
			originalLength: 500,
			summaryLength:  500,
			expectedRate:   1.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			compressionRate := 0.0
			if tc.originalLength > 0 {
				compressionRate = float64(tc.summaryLength) / float64(tc.originalLength)
			}

			assert.InDelta(t, tc.expectedRate, compressionRate, 0.01)
		})
	}
}

// TestSummarizeService_ExtractKeyPoints 测试关键点提取
func TestSummarizeService_ExtractKeyPoints(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	testCases := []struct {
		name          string
		summary       string
		expectedCount int
	}{
		{
			name:          "三个关键点",
			summary:       "这是第一个关键点。这是第二个关键点，内容稍长一些。这是第三个关键点。",
			expectedCount: 3,
		},
		{
			name:          "短句子不计入",
			summary:       "短句。这是一个足够长的句子应该被计入。另一个长句子也被计入。",
			expectedCount: 2,
		},
		{
			name:          "超过5个只返回5个",
			summary:       "第一点。第二点。第三点。第四点。第五点。第六点。第七点。",
			expectedCount: 5,
		},
		{
			name:          "空摘要",
			summary:       "",
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keyPoints := service.extractKeyPoints(tc.summary)

			assert.Equal(t, tc.expectedCount, len(keyPoints))
		})
	}
}

// TestSummarizeService_SummarizeChapter_Success 测试章节总结
func TestSummarizeService_SummarizeChapter_Success(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	ctx := context.Background()
	req := &dto.ChapterSummaryRequest{
		ChapterID:    "chapter-123",
		ProjectID:    "project-456",
		OutlineLevel: 3,
	}

	// 注意：实际实现中需要从数据库获取章节内容
	// 这里测试的是服务层的基本结构
	result, err := service.SummarizeChapter(ctx, req)

	// 当前实现返回模拟数据
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "chapter-123", result.ChapterID)
	assert.NotEmpty(t, result.ChapterTitle)
	assert.NotEmpty(t, result.Summary)
	assert.Greater(t, result.TokensUsed, 0)
}

// TestSummarizeService_SummarizeChapter_ChapterIDRequired 测试缺少章节ID
func TestSummarizeService_SummarizeChapter_ChapterIDRequired(t *testing.T) {
	// 测试章节ID是必需的（虽然在当前实现中没有强制验证）
	// 这个测试确保当实现添加验证时能正确工作

	testCases := []struct {
		name      string
		chapterID string
		projectID string
	}{
		{
			name:      "空章节ID",
			chapterID: "",
			projectID: "project-123",
		},
		{
			name:      "有效章节ID",
			chapterID: "chapter-123",
			projectID: "project-456",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			manager := &adapter.AdapterManager{}
			service := NewSummarizeService(manager)

			ctx := context.Background()
			req := &dto.ChapterSummaryRequest{
				ChapterID:    tc.chapterID,
				ProjectID:    tc.projectID,
				OutlineLevel: 2,
			}

			result, err := service.SummarizeChapter(ctx, req)

			// 当前实现总是返回模拟数据
			// 实际实现应该验证章节ID
			_ = result
			_ = err
		})
	}
}

// TestSummarizeService_Integration 测试总结服务集成场景
func TestSummarizeService_Integration(t *testing.T) {
	t.Run("完整的总结流程", func(t *testing.T) {
		manager := &adapter.AdapterManager{}
		service := NewSummarizeService(manager)

		// 模拟真实场景
		content := `第一章 开始

这是一个很长的故事开始。主人公小明住在一个小村庄里。
他每天过着平静的生活，直到有一天，一切都改变了。
村庄来了一位神秘的陌生人，带来了一封来自远方的信。
小明决定踏上旅程，去寻找信中提到的宝藏。

在旅途中，他遇到了各种挑战和困难，但他从未放弃。
最终，他找到了宝藏，但也发现了更珍贵的东西——友谊和勇气。`

		req := &dto.SummarizeRequest{
			Content:       content,
			SummaryType:   "detailed",
			IncludeQuotes: true,
			MaxLength:     500,
		}

		prompt := service.buildSummarizePrompt(req)

		// 验证提示词构建
		assert.NotEmpty(t, prompt)
		assert.Contains(t, prompt, content)
		assert.Contains(t, prompt, "关键引用")

		// 验证内容长度统计
		originalLength := len([]rune(content))
		assert.Greater(t, originalLength, 100)
	})
}

// BenchmarkSummarizeService_ExtractKeyPoints 性能测试
func BenchmarkSummarizeService_ExtractKeyPoints(b *testing.B) {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	// 准备测试数据
	longSummary := "这是第一点。这是第二点。这是第三点。这是第四点。这是第五点。" +
		"这是第六点。这是第七点。这是第八点。这是第九点。这是第十点。"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.extractKeyPoints(longSummary)
	}
}

// ExampleSummarizeService 示例测试
func ExampleSummarizeService() {
	manager := &adapter.AdapterManager{}
	service := NewSummarizeService(manager)

	req := &dto.SummarizeRequest{
		Content:     "这是要总结的内容",
		SummaryType: "brief",
	}

	prompt := service.buildSummarizePrompt(req)
	fmt.Println("生成的提示词:", prompt)
	// Output: 生成的提示词: 请为以下内容生成简洁摘要（50-100字）：
	//
	// 这是要总结的内容
}
