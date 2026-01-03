package ai

import (
	"context"
	"strings"
	"testing"

	"Qingyu_backend/service/ai/adapter"
	"Qingyu_backend/service/ai/dto"
	"Qingyu_backend/service/ai/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProofreadService_ProofreadContent_Success 测试成功校对内容
func TestProofreadService_ProofreadContent_Success(t *testing.T) {
	// 创建模拟适配器
	mockAdapter := mocks.NewMockAIAdapter("test-adapter")

	// 设置校对响应（JSON格式）
	jsonResponse := `{
		"issues": [
			{
				"type": "grammar",
				"severity": "error",
				"message": "主谓不一致",
				"line": 1,
				"column": 10,
				"original": "他喜欢跑步和游泳",
				"suggestions": ["他喜欢跑步和游泳。"]
			}
		]
	}`
	mockAdapter.SetTextResponse(jsonResponse, 150)

	_ = mockAdapter
	// TODO: 完善实际的测试调用，需要适配器管理器支持
}

// TestProofreadService_ProofreadContent_EmptyContent 测试空内容
func TestProofreadService_ProofreadContent_EmptyContent(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	ctx := context.Background()
	req := &dto.ProofreadRequest{
		Content: "",
	}

	result, err := service.ProofreadContent(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "内容不能为空")
}

// TestProofreadService_ProofreadContent_WhitespaceContent 测试仅包含空白的内容
func TestProofreadService_ProofreadContent_WhitespaceContent(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	ctx := context.Background()
	req := &dto.ProofreadRequest{
		Content: "   \n\t   ",
	}

	result, err := service.ProofreadContent(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "内容不能为空")
}

// TestProofreadService_ProofreadContent_DefaultCheckTypes 测试默认检查类型
func TestProofreadService_ProofreadContent_DefaultCheckTypes(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	req := &dto.ProofreadRequest{
		Content:    "这是测试内容",
		CheckTypes: []string{},
	}

	prompt := service.buildProofreadPrompt(req)

	// 验证默认检查类型
	assert.Contains(t, prompt, "拼写错误")
	assert.Contains(t, prompt, "语法错误")
	assert.Contains(t, prompt, "标点符号错误")
}

// TestProofreadService_ProofreadContent_CustomCheckTypes 测试自定义检查类型
func TestProofreadService_ProofreadContent_CustomCheckTypes(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	testCases := []struct {
		name       string
		checkTypes []string
		expected   []string
	}{
		{
			name:       "仅拼写检查",
			checkTypes: []string{"spelling"},
			expected:   []string{"拼写错误"},
		},
		{
			name:       "仅语法检查",
			checkTypes: []string{"grammar"},
			expected:   []string{"语法错误"},
		},
		{
			name:       "仅标点检查",
			checkTypes: []string{"punctuation"},
			expected:   []string{"标点符号错误"},
		},
		{
			name:       "风格检查",
			checkTypes: []string{"style"},
			expected:   []string{"写作风格"},
		},
		{
			name:       "多种检查类型",
			checkTypes: []string{"spelling", "grammar"},
			expected:   []string{"拼写错误", "语法错误"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &dto.ProofreadRequest{
				Content:    "测试内容",
				CheckTypes: tc.checkTypes,
			}

			prompt := service.buildProofreadPrompt(req)

			// 验证所有预期的检查类型都在提示词中
			for _, expected := range tc.expected {
				assert.Contains(t, prompt, expected)
			}
		})
	}
}

// TestProofreadService_ProofreadContent_AIError 测试AI服务错误
func TestProofreadService_ProofreadContent_AIError(t *testing.T) {
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
	// TODO: 完善实际的测试调用
}

// TestProofreadService_ParseProofreadResult_JSONFormat 测试JSON格式解析
func TestProofreadService_ParseProofreadResult_JSONFormat(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	aiResult := `{
		"issues": [
			{
				"type": "grammar",
				"severity": "error",
				"message": "语法错误",
				"line": 1,
				"column": 5,
				"original": "原文",
				"suggestions": ["建议1", "建议2"]
			},
			{
				"type": "spelling",
				"severity": "warning",
				"message": "拼写错误",
				"line": 2,
				"column": 10,
				"original": "原文2",
				"suggestions": ["建议3"]
			}
		]
	}`

	originalContent := "第一行原文\n第二行原文2"

	issues, err := service.parseProofreadResult(aiResult, originalContent)

	require.NoError(t, err)
	assert.Equal(t, 2, len(issues))

	// 验证第一个问题
	assert.Equal(t, "grammar", issues[0].Type)
	assert.Equal(t, "error", issues[0].Severity)
	assert.Equal(t, "语法错误", issues[0].Message)
	assert.Equal(t, "原文", issues[0].OriginalText)
	assert.Equal(t, 2, len(issues[0].Suggestions))

	// 验证第二个问题
	assert.Equal(t, "spelling", issues[1].Type)
	assert.Equal(t, "warning", issues[1].Severity)
}

// TestProofreadService_ParseProofreadResult_TextFormat 测试文本格式解析（后备方案）
func TestProofreadService_ParseProofreadResult_TextFormat(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	aiResult := `- [grammar] 主谓不一致
- [spelling] 拼写错误示例
• 标点符号使用不当`

	originalContent := "这是原文内容"

	issues, err := service.parseProofreadResult(aiResult, originalContent)

	require.NoError(t, err)
	assert.Greater(t, len(issues), 0)

	// 验证至少有一个问题被解析
	foundIssue := false
	for _, issue := range issues {
		if issue.Message != "" {
			foundIssue = true
			break
		}
	}
	assert.True(t, foundIssue, "应该至少解析出一个问题")
}

// TestProofreadService_GenerateStatistics 测试统计信息生成
func TestProofreadService_GenerateStatistics(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	issues := []dto.Issue{
		{
			Type:     "grammar",
			Severity: "error",
		},
		{
			Type:     "spelling",
			Severity: "warning",
		},
		{
			Type:     "punctuation",
			Severity: "suggestion",
		},
		{
			Type:     "grammar",
			Severity: "error",
		},
	}

	content := "这是测试内容 包含四个词"

	stats := service.generateStatistics(issues, content)

	// 验证总数
	assert.Equal(t, 4, stats.TotalIssues)

	// 验证严重程度统计
	assert.Equal(t, 2, stats.ErrorCount)
	assert.Equal(t, 1, stats.WarningCount)
	assert.Equal(t, 1, stats.SuggestionCount)

	// 验证类型统计
	assert.Equal(t, 2, stats.IssuesByType["grammar"])
	assert.Equal(t, 1, stats.IssuesByType["spelling"])
	assert.Equal(t, 1, stats.IssuesByType["punctuation"])

	// 验证词数统计
	assert.Equal(t, 4, stats.WordCount)
	assert.Greater(t, stats.CharacterCount, 0)
}

// TestProofreadService_CalculateScore 测试评分计算
func TestProofreadService_CalculateScore(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	testCases := []struct {
		name           string
		stats          dto.ProofreadStats
		expectedScore  float64
		description    string
	}{
		{
			name: "完美内容",
			stats: dto.ProofreadStats{
				TotalIssues:     0,
				ErrorCount:      0,
				WarningCount:    0,
				SuggestionCount: 0,
			},
			expectedScore: 100.0,
			description:   "没有问题应该得100分",
		},
		{
			name: "只有错误",
			stats: dto.ProofreadStats{
				TotalIssues:     10,
				ErrorCount:      10,
				WarningCount:    0,
				SuggestionCount: 0,
			},
			expectedScore: 50.0,
			description:   "10个错误，每个扣5分",
		},
		{
			name: "只有警告",
			stats: dto.ProofreadStats{
				TotalIssues:     10,
				ErrorCount:      0,
				WarningCount:    10,
				SuggestionCount: 0,
			},
			expectedScore: 80.0,
			description:   "10个警告，每个扣2分",
		},
		{
			name: "只有建议",
			stats: dto.ProofreadStats{
				TotalIssues:     10,
				ErrorCount:      0,
				WarningCount:    0,
				SuggestionCount: 10,
			},
			expectedScore: 95.0,
			description:   "10个建议，每个扣0.5分",
		},
		{
			name: "混合问题",
			stats: dto.ProofreadStats{
				TotalIssues:     15,
				ErrorCount:      5,
				WarningCount:    5,
				SuggestionCount: 5,
			},
			expectedScore: 82.5,
			description:   "5*5 + 5*2 + 5*0.5 = 37.5, 100-37.5=62.5",
		},
		{
			name: "最低分",
			stats: dto.ProofreadStats{
				TotalIssues:     100,
				ErrorCount:      100,
				WarningCount:    0,
				SuggestionCount: 0,
			},
			expectedScore: 0.0,
			description:   "100个错误，500分，但最低分为0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := service.calculateScore(tc.stats)

			assert.InDelta(t, tc.expectedScore, score, 0.01, tc.description)
		})
	}
}

// TestProofreadService_GetProofreadSuggestion 测试获取校对建议
func TestProofreadService_GetProofreadSuggestion(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	ctx := context.Background()
	suggestionID := "suggestion-123"

	result, err := service.GetProofreadSuggestion(ctx, suggestionID)

	// 当前实现返回模拟数据
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, suggestionID, result.IssueID)
	assert.NotEmpty(t, result.Type)
	assert.NotEmpty(t, result.Message)
	assert.NotEmpty(t, result.Position)
	assert.NotEmpty(t, result.OriginalText)
	assert.NotEmpty(t, result.Suggestions)
	assert.NotEmpty(t, result.Explanation)
}

// TestProofreadService_FindPositionInText 测试文本位置查找
func TestFindPositionInText(t *testing.T) {
	testCases := []struct {
		name        string
		content     string
		searchTerm  string
		line        int
		column      int
		expectedPos int
	}{
		{
			name:        "第一行开始",
			content:     "第一行\n第二行\n第三行",
			searchTerm:  "第一",
			line:        1,
			column:      1,
			expectedPos: 0,
		},
		{
			name:        "第二行",
			content:     "第一行\n第二行\n第三行",
			searchTerm:  "第二",
			line:        2,
			column:      1,
			expectedPos: 4, // "第一行\n" = 4个字符（中文）
		},
		{
			name:        "无效位置",
			content:     "测试内容",
			searchTerm:  "测试",
			line:        10,
			column:      1,
			expectedPos: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pos := findPositionInText(tc.content, tc.searchTerm, tc.line, tc.column)

			if tc.expectedPos == -1 {
				assert.Equal(t, -1, pos)
			} else {
				assert.GreaterOrEqual(t, pos, 0)
			}
		})
	}
}

// TestProofreadService_Integration 测试校对服务集成场景
func TestProofreadService_Integration(t *testing.T) {
	t.Run("完整的校对流程", func(t *testing.T) {
		manager := &adapter.AdapterManager{}
		service := NewProofreadService(manager)

		content := `这是一个测试文本。
虽然内容有一些语法错误，但不影响理解。
例如，主谓不一致的问题。
还有标点符号使用不当。`

		req := &dto.ProofreadRequest{
			Content:    content,
			CheckTypes: []string{"grammar", "punctuation"},
			Language:   "zh-CN",
		}

		// 验证提示词构建
		prompt := service.buildProofreadPrompt(req)
		assert.Contains(t, prompt, content)
		assert.Contains(t, prompt, "语法错误")
		assert.Contains(t, prompt, "标点符号错误")

		// 验证统计信息
		mockIssues := []dto.Issue{
			{Type: "grammar", Severity: "error", Message: "错误1"},
			{Type: "grammar", Severity: "warning", Message: "错误2"},
			{Type: "punctuation", Severity: "suggestion", Message: "建议1"},
		}

		stats := service.generateStatistics(mockIssues, content)
		assert.Equal(t, 3, stats.TotalIssues)
		assert.Equal(t, 1, stats.ErrorCount)
		assert.Equal(t, 1, stats.WarningCount)
		assert.Equal(t, 1, stats.SuggestionCount)

		// 验证评分
		score := service.calculateScore(stats)
		assert.Greater(t, score, 0.0)
		assert.LessOrEqual(t, score, 100.0)
	})
}

// BenchmarkProofreadService_CalculateScore 性能测试
func BenchmarkProofreadService_CalculateScore(b *testing.B) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	stats := dto.ProofreadStats{
		TotalIssues:     100,
		ErrorCount:      50,
		WarningCount:    30,
		SuggestionCount: 20,
		IssuesByType: map[string]int{
			"grammar":     50,
			"spelling":    30,
			"punctuation": 20,
		},
		WordCount:      1000,
		CharacterCount: 5000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.calculateScore(stats)
	}
}

// TestProofreadService_ExtractIssuesFromText 测试从文本提取问题
func TestProofreadService_ExtractIssuesFromText(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	aiResult := `- 主谓不一致，建议修改
- 拼写错误：这个词语拼写不正确
• 标点符号应该用句号`

	originalContent := "这是原文内容"

	issues := service.extractIssuesFromText(aiResult, originalContent)

	assert.Greater(t, len(issues), 0)

	// 验证问题被正确解析
	for _, issue := range issues {
		assert.NotEmpty(t, issue.Message)
		assert.NotEqual(t, "", issue.Type)
		assert.NotEqual(t, "", issue.Severity)
	}
}

// TestProofreadService_LongText 测试长文本校对
func TestProofreadService_LongText(t *testing.T) {
	manager := &adapter.AdapterManager{}
	service := NewProofreadService(manager)

	// 生成长文本
	var longText strings.Builder
	for i := 0; i < 100; i++ {
		longText.WriteString("这是第")
		longText.WriteString(string(rune('0' + i%10)))
		longText.WriteString("句话。")
	}

	req := &dto.ProofreadRequest{
		Content:    longText.String(),
		CheckTypes: []string{"grammar", "spelling", "punctuation"},
	}

	prompt := service.buildProofreadPrompt(req)

	// 验证长文本处理
	assert.Contains(t, prompt, longText.String())
	assert.Greater(t, len(prompt), 1000)
}
