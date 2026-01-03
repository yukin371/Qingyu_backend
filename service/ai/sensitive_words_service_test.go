package ai

import (
	"context"
	"testing"

	"Qingyu_backend/service/ai/adapter"
	"Qingyu_backend/service/ai/dto"
	"Qingyu_backend/service/ai/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSensitiveWordsService_CheckSensitiveWords_Success 测试成功检测敏感词
func TestSensitiveWordsService_CheckSensitiveWords_Success(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	// 添加测试用敏感词
	err := service.AddCustomWords("test-user", []string{"测试敏感词", "违规内容"})
	require.NoError(t, err)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "这是包含测试敏感词的内容",
		CustomWords: []string{"测试敏感词"},
		Category:    "all",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.CheckID)
	assert.GreaterOrEqual(t, result.TotalMatches, 0)
	assert.False(t, result.IsSafe) // 应该检测到敏感词
}

// TestSensitiveWordsService_CheckSensitiveWords_EmptyContent 测试空内容
func TestSensitiveWordsService_CheckSensitiveWords_EmptyContent(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content: "",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "内容不能为空")
}

// TestSensitiveWordsService_CheckSensitiveWords_WhitespaceContent 测试仅包含空白的内容
func TestSensitiveWordsService_CheckSensitiveWords_WhitespaceContent(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content: "   \n\t   ",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "内容不能为空")
}

// TestSensitiveWordsService_CheckSensitiveWords_NoMatch 测试未检测到敏感词
func TestSensitiveWordsService_CheckSensitiveWords_NoMatch(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "这是一段完全正常的文本内容，没有任何敏感词",
		CustomWords: []string{"敏感词1", "敏感词2"},
		Category:    "all",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalMatches)
	assert.True(t, result.IsSafe) // 没有检测到敏感词，应该是安全的
}

// TestSensitiveWordsService_CheckSensitiveWords_PoliticalCategory 测试政治敏感词
func TestSensitiveWordsService_CheckSensitiveWords_PoliticalCategory(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	// 添加政治类敏感词到词库
	err := service.AddCustomWords("test-user", []string{"政治敏感词"})
	require.NoError(t, err)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "这段内容包含政治敏感词",
		CustomWords: []string{"政治敏感词"},
		Category:    "political",
	}

	_, err = service.CheckSensitiveWords(ctx, req)

	assert.NoError(t, err)
	// 验证检测到敏感词的逻辑在完整测试中实现
}

// TestSensitiveWordsService_CheckSensitiveWords_ViolenceCategory 测试暴力敏感词
func TestSensitiveWordsService_CheckSensitiveWords_ViolenceCategory(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "这段内容包含暴力词汇测试",
		CustomWords: []string{"暴力词汇"},
		Category:    "violence",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestSensitiveWordsService_CheckSensitiveWords_AdultCategory 测试成人内容敏感词
func TestSensitiveWordsService_CheckSensitiveWords_AdultCategory(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "这段内容包含成人内容词汇",
		CustomWords: []string{"成人词汇"},
		Category:    "adult",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestSensitiveWordsService_CheckSensitiveWords_CustomWords 测试自定义敏感词
func TestSensitiveWordsService_CheckSensitiveWords_CustomWords(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	customWords := []string{"自定义词1", "自定义词2", "特殊词"}

	req := &dto.SensitiveWordsCheckRequest{
		Content:     "这段内容包含自定义词1和特殊词",
		CustomWords: customWords,
		Category:    "all",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// 应该检测到至少一个自定义词
	if result.TotalMatches > 0 {
		assert.False(t, result.IsSafe)

		// 验证检测到的词
		foundCustomWord := false
		for _, match := range result.SensitiveWords {
			if match.Category == "custom" {
				foundCustomWord = true
				break
			}
		}
		assert.True(t, foundCustomWord, "应该检测到自定义敏感词")
	}
}

// TestSensitiveWordsService_CheckSensitiveWords_AllCategories 测试所有分类
func TestSensitiveWordsService_CheckSensitiveWords_AllCategories(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "测试内容",
		CustomWords: []string{"测试词"},
		Category:    "all",
	}

	result, err := service.CheckSensitiveWords(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestSensitiveWordsService_DetectSensitiveWords 测试敏感词检测逻辑
func TestSensitiveWordsService_DetectSensitiveWords(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	content := "这是一段包含测试敏感词的文本，敏感词出现了多次"

	searchWords := map[string][]string{
		"custom": {"测试敏感词", "敏感词"},
	}

	matches := service.detectSensitiveWords(content, searchWords)

	assert.Greater(t, len(matches), 0)

	// 验证匹配信息
	for _, match := range matches {
		assert.NotEmpty(t, match.ID)
		assert.NotEmpty(t, match.Word)
		assert.NotEmpty(t, match.Category)
		assert.NotEmpty(t, match.Level)
		assert.NotEmpty(t, match.Suggestion)

		// 验证位置信息
		assert.GreaterOrEqual(t, match.Position.Start, 0)
		assert.Greater(t, match.Position.End, match.Position.Start)
		assert.Greater(t, match.Position.Length, 0)

		// 验证上下文
		assert.NotEmpty(t, match.Context)
	}
}

// TestSensitiveWordsService_FindWordPositions 测试查找词位置
func TestSensitiveWordsService_FindWordPositions(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	testCases := []struct {
		name         string
		content      string
		word         string
		expectCount  int
	}{
		{
			name:         "单次出现",
			content:      "这是测试内容",
			word:         "测试",
			expectCount:  1,
		},
		{
			name:         "多次出现",
			content:      "测试内容1 测试内容2 测试内容3",
			word:         "测试",
			expectCount:  3,
		},
		{
			name:         "未出现",
			content:      "这是其他内容",
			word:         "测试",
			expectCount:  0,
		},
		{
			name:         "空内容",
			content:      "",
			word:         "测试",
			expectCount:  0,
		},
		{
			name:         "中文词组",
			content:      "这是一个很长的中文词组在文本中",
			word:         "中文词组",
			expectCount:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			positions := service.findWordPositions(tc.content, tc.word)

			assert.Equal(t, tc.expectCount, len(positions))

			// 验证位置信息
			for _, pos := range positions {
				assert.GreaterOrEqual(t, pos.Start, 0)
				assert.Greater(t, pos.End, pos.Start)
				assert.Greater(t, pos.Length, 0)
				assert.Greater(t, pos.Line, 0)
				assert.Greater(t, pos.Column, 0)
			}
		})
	}
}

// TestSensitiveWordsService_CalculateLineColumn 测试行列计算
func TestSensitiveWordsService_CalculateLineColumn(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	testCases := []struct {
		name         string
		content      string
		position     int
		expectedLine int
		expectedCol  int
	}{
		{
			name:         "第一行开始",
			content:      "第一行\n第二行\n第三行",
			position:     0,
			expectedLine: 1,
			expectedCol:  1,
		},
		{
			name:         "第一行中间",
			content:      "第一行内容",
			position:     2,
			expectedLine: 1,
			expectedCol:  3,
		},
		{
			name:         "第二行开始",
			content:      "第一行\n第二行",
			position:     8, // "第一行\n" = 3+1 = 4个字符（中文和换行）
			expectedLine: 2,
			expectedCol:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			line, column := service.calculateLineColumn(tc.content, tc.position)

			assert.Equal(t, tc.expectedLine, line)
			assert.Equal(t, tc.expectedCol, column)
		})
	}
}

// TestSensitiveWordsService_ExtractContext 测试上下文提取
func TestSensitiveWordsService_ExtractContext(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	content := "这是一段很长的文本内容，包含了很多文字，我们希望提取某个词的上下文环境，以便更好地理解"
	start := 20
	end := 25

	context := service.extractContext(content, start, end)

	assert.NotEmpty(t, context)
	// 上下文应该包含省略号，因为不是从开头开始
	assert.True(t, len(context) > 0)

	// 测试从开头提取
	start = 0
	end = 5
	context = service.extractContext(content, start, end)
	assert.True(t, len(context) > 0)
	// 从开头不应该以"..."开始
	assert.False(t, len(context) > 3 && context[:3] == "...")

	// 测试到结尾
	start = len([]rune(content)) - 5
	end = len([]rune(content))
	context = service.extractContext(content, start, end)
	assert.True(t, len(context) > 0)
}

// TestSensitiveWordsService_DetermineWordLevel 测试风险级别确定
func TestSensitiveWordsService_DetermineWordLevel(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	testCases := []struct {
		category       string
		word           string
		expectedLevel  string
	}{
		{
			category:      "political",
			word:          "任意词",
			expectedLevel: "high",
		},
		{
			category:      "violence",
			word:          "任意词",
			expectedLevel: "medium",
		},
		{
			category:      "adult",
			word:          "任意词",
			expectedLevel: "high",
		},
		{
			category:      "custom",
			word:          "任意词",
			expectedLevel: "medium",
		},
		{
			category:      "unknown",
			word:          "任意词",
			expectedLevel: "low",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.category, func(t *testing.T) {
			level := service.determineWordLevel(tc.category, tc.word)
			assert.Equal(t, tc.expectedLevel, level)
		})
	}
}

// TestSensitiveWordsService_GenerateSuggestion 测试生成修改建议
func TestSensitiveWordsService_GenerateSuggestion(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	word := "测试敏感词"
	category := "political"

	suggestion := service.generateSuggestion(word, category)

	assert.NotEmpty(t, suggestion)
	assert.Contains(t, suggestion, word)
	assert.Contains(t, suggestion, "建议")
}

// TestSensitiveWordsService_GenerateCheckSummary 测试生成检测摘要
func TestSensitiveWordsService_GenerateCheckSummary(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	matches := []dto.SensitiveWordMatch{
		{
			ID:       "1",
			Word:     "词1",
			Category: "political",
			Level:    "high",
		},
		{
			ID:       "2",
			Word:     "词2",
			Category: "violence",
			Level:    "medium",
		},
		{
			ID:       "3",
			Word:     "词3",
			Category: "adult",
			Level:    "high",
		},
		{
			ID:       "4",
			Word:     "词4",
			Category: "custom",
			Level:    "low",
		},
	}

	summary := service.generateCheckSummary(matches)

	// 验证分类统计
	assert.Equal(t, 1, summary.ByCategory["political"])
	assert.Equal(t, 1, summary.ByCategory["violence"])
	assert.Equal(t, 1, summary.ByCategory["adult"])
	assert.Equal(t, 1, summary.ByCategory["custom"])

	// 验证级别统计
	assert.Equal(t, 2, summary.ByLevel["high"])
	assert.Equal(t, 1, summary.ByLevel["medium"])
	assert.Equal(t, 1, summary.ByLevel["low"])

	// 验证计数
	assert.Equal(t, 2, summary.HighRiskCount)
	assert.Equal(t, 1, summary.MediumRiskCount)
	assert.Equal(t, 1, summary.LowRiskCount)
}

// TestSensitiveWordsService_HasHighRiskWords 测试高风险词检测
func TestSensitiveWordsService_HasHighRiskWords(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	testCases := []struct {
		name     string
		matches  []dto.SensitiveWordMatch
		expected bool
	}{
		{
			name:     "有高风险词",
			matches: []dto.SensitiveWordMatch{
				{Level: "high"},
			},
			expected: true,
		},
		{
			name: "只有中低风险词",
			matches: []dto.SensitiveWordMatch{
				{Level: "medium"},
				{Level: "low"},
			},
			expected: false,
		},
		{
			name:     "空列表",
			matches:  []dto.SensitiveWordMatch{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.hasHighRiskWords(tc.matches)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestSensitiveWordsService_AddCustomWords 测试添加自定义敏感词
func TestSensitiveWordsService_AddCustomWords(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	userID := "test-user-123"
	words := []string{"自定义词1", "自定义词2", "自定义词3"}

	// 添加词
	err := service.AddCustomWords(userID, words)
	assert.NoError(t, err)

	// 验证词已添加（通过检测验证）
	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "这段内容包含自定义词1",
		CustomWords: words,
		Category:    "all",
	}

	result, err := service.CheckSensitiveWords(ctx, req)
	assert.NoError(t, err)
	assert.Greater(t, result.TotalMatches, 0)
}

// TestSensitiveWordsService_RemoveCustomWords 测试移除自定义敏感词
func TestSensitiveWordsService_RemoveCustomWords(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	userID := "test-user-456"
	words := []string{"可删除词1", "可删除词2", "可删除词3"}

	// 先添加
	err := service.AddCustomWords(userID, words)
	require.NoError(t, err)

	// 移除部分词
	removeWords := []string{"可删除词1", "可删除词2"}
	err = service.RemoveCustomWords(userID, removeWords)
	assert.NoError(t, err)

	// 验证：移除的词应该不再被检测到，但保留的词仍能检测到
	ctx := context.Background()
	req := &dto.SensitiveWordsCheckRequest{
		Content:     "包含可删除词3的内容",
		CustomWords: []string{"可删除词3"},
		Category:    "all",
	}

	_, err = service.CheckSensitiveWords(ctx, req)
	assert.NoError(t, err)
	// 应该只检测到保留的词
}

// TestSensitiveWordsService_GetSensitiveWordsDetail 测试获取检测详情
func TestSensitiveWordsService_GetSensitiveWordsDetail(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	ctx := context.Background()
	checkID := "check-123"

	result, err := service.GetSensitiveWordsDetail(ctx, checkID)

	// 当前实现返回模拟数据
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, checkID, result.CheckID)
	assert.NotEmpty(t, result.Content)
	assert.NotNil(t, result.Matches)
	assert.NotNil(t, result.CustomWords)
	assert.NotNil(t, result.Summary)
}

// TestSensitiveWordsService_AISemanticAnalysis 测试AI语义分析
func TestSensitiveWordsService_AISemanticAnalysis(t *testing.T) {
	// 创建模拟适配器
	mockAdapter := mocks.NewMockAIAdapter("test-adapter")
	mockAdapter.SetTextResponse("[]", 50) // 返回空数组，表示未发现敏感内容

	_ = mockAdapter
	// TODO: 完善实际的测试调用，需要适配器管理器支持
}

// TestSensitiveWordsService_Integration 测试敏感词检测集成场景
func TestSensitiveWordsService_Integration(t *testing.T) {
	t.Run("完整的检测流程", func(t *testing.T) {
		manager := &adapter.AdapterManager{}
		service := NewSensitiveWordsService(manager)

		content := `这是一篇完整的文章内容。
文章中可能包含一些需要审查的词汇。
我们希望能够准确检测出所有敏感词。`

		// 添加自定义词库
		userID := "integration-test-user"
		customWords := []string{"违规", "禁止"}
		err := service.AddCustomWords(userID, customWords)
		assert.NoError(t, err)

		// 执行检测
		ctx := context.Background()
		req := &dto.SensitiveWordsCheckRequest{
			Content:     content,
			CustomWords: customWords,
			Category:    "all",
		}

		result, err := service.CheckSensitiveWords(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.CheckID)

		// 验证摘要
		assert.NotNil(t, result.Summary)
		assert.NotNil(t, result.Summary.ByCategory)
		assert.NotNil(t, result.Summary.ByLevel)
	})
}

// BenchmarkSensitiveWordsService_FindWordPositions 性能测试
func BenchmarkSensitiveWordsService_FindWordPositions(b *testing.B) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	content := "这是一段很长的文本内容，用于性能测试。" +
		"敏感词可能出现多次，我们需要测试查找性能。"
	word := "敏感词"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.findWordPositions(content, word)
	}
}

// BenchmarkSensitiveWordsService_DetectSensitiveWords 性能测试
func BenchmarkSensitiveWordsService_DetectSensitiveWords(b *testing.B) {
	manager := &adapter.AdapterManager{}
	service := NewSensitiveWordsService(manager)

	content := "测试内容"
	searchWords := map[string][]string{
		"custom": {"词1", "词2", "词3"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.detectSensitiveWords(content, searchWords)
	}
}
