package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"Qingyu_backend/service/ai/adapter"
	"Qingyu_backend/service/ai/dto"
	"github.com/google/uuid"
)

// ProofreadService 文本校对服务
type ProofreadService struct {
	adapterManager *adapter.AdapterManager
	// 可以添加存储层来保存校对结果
	// repository ProofreadRepository
}

// NewProofreadService 创建文本校对服务
func NewProofreadService(adapterManager *adapter.AdapterManager) *ProofreadService {
	return &ProofreadService{
		adapterManager: adapterManager,
	}
}

// ProofreadContent 校对文本内容
func (s *ProofreadService) ProofreadContent(ctx context.Context, req *dto.ProofreadRequest) (*dto.ProofreadResponse, error) {
	// 参数验证
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("内容不能为空")
	}

	// 设置默认检查类型
	checkTypes := req.CheckTypes
	if len(checkTypes) == 0 {
		checkTypes = []string{"spelling", "grammar", "punctuation"}
	}

	// 构建校对提示词
	prompt := s.buildProofreadPrompt(req)

	// 调用AI进行校对
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      prompt,
		Temperature: 0.3, // 校对任务使用较低温度以确保准确性
		MaxTokens:   2000,
	}

	result, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("校对失败: %w", err)
	}

	// 解析AI返回的问题列表
	issues, err := s.parseProofreadResult(result.Text, req.Content)
	if err != nil {
		return nil, fmt.Errorf("解析校对结果失败: %w", err)
	}

	// 生成统计信息
	stats := s.generateStatistics(issues, req.Content)

	// 计算整体评分
	score := s.calculateScore(stats)

	return &dto.ProofreadResponse{
		OriginalContent: req.Content,
		Issues:          issues,
		Score:           score,
		Statistics:      stats,
		TokensUsed:      result.Usage.TotalTokens,
		Model:           result.Model,
		ProcessedAt:     time.Now(),
	}, nil
}

// GetProofreadSuggestion 获取校对建议详情
func (s *ProofreadService) GetProofreadSuggestion(ctx context.Context, suggestionID string) (*dto.ProofreadSuggestion, error) {
	// TODO: 从存储层获取建议详情
	// 这里返回模拟数据
	return &dto.ProofreadSuggestion{
		IssueID: suggestionID,
		Type:    "grammar",
		Message: "建议修改语法错误",
		Position: dto.TextPosition{
			Line:   1,
			Column: 10,
			Start:  10,
			End:    20,
			Length: 10,
		},
		OriginalText: "原文示例",
		Suggestions: []dto.SuggestionItem{
			{
				Text:       "建议文本",
				Confidence: 0.95,
				Reason:     "语法更通顺",
			},
		},
		Explanation: "这是一个语法错误的示例说明",
		Examples:    []string{"正确示例1", "正确示例2"},
	}, nil
}

// buildProofreadPrompt 构建校对提示词
func (s *ProofreadService) buildProofreadPrompt(req *dto.ProofreadRequest) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("请对以下文本进行校对，找出所有错误并提供修改建议：\n\n")
	promptBuilder.WriteString(req.Content)
	promptBuilder.WriteString("\n\n请检查以下方面：\n")

	if len(req.CheckTypes) == 0 {
		promptBuilder.WriteString("- 拼写错误\n")
		promptBuilder.WriteString("- 语法错误\n")
		promptBuilder.WriteString("- 标点符号错误\n")
	} else {
		for _, checkType := range req.CheckTypes {
			switch checkType {
			case "spelling":
				promptBuilder.WriteString("- 拼写错误\n")
			case "grammar":
				promptBuilder.WriteString("- 语法错误\n")
			case "punctuation":
				promptBuilder.WriteString("- 标点符号错误\n")
			case "style":
				promptBuilder.WriteString("- 写作风格\n")
			}
		}
	}

	promptBuilder.WriteString("\n请以JSON格式返回结果，包含以下字段：\n")
	promptBuilder.WriteString("- issues: 问题列表，每个问题包含type, severity, message, position, suggestions等字段\n")
	promptBuilder.WriteString("- 总体评分（0-100分）\n")

	return promptBuilder.String()
}

// parseProofreadResult 解析AI返回的校对结果
func (s *ProofreadService) parseProofreadResult(aiResult, originalContent string) ([]dto.Issue, error) {
	var issues []dto.Issue

	// 尝试解析JSON格式
	var result struct {
		Issues []struct {
			Type        string   `json:"type"`
			Severity    string   `json:"severity"`
			Message     string   `json:"message"`
			Line        int      `json:"line"`
			Column      int      `json:"column"`
			Original    string   `json:"original"`
			Suggestions []string `json:"suggestions"`
		} `json:"issues"`
	}

	if err := json.Unmarshal([]byte(aiResult), &result); err == nil {
		// 成功解析JSON
		for _, issue := range result.Issues {
			// 查找原文中的位置
			start := findPositionInText(originalContent, issue.Original, issue.Line, issue.Column)
			if start == -1 {
				continue
			}

			issues = append(issues, dto.Issue{
				ID:       uuid.New().String(),
				Type:     issue.Type,
				Severity: issue.Severity,
				Message:  issue.Message,
				Position: dto.TextPosition{
					Line:   issue.Line,
					Column: issue.Column,
					Start:  start,
					End:    start + len([]rune(issue.Original)),
					Length: len([]rune(issue.Original)),
				},
				OriginalText: issue.Original,
				Suggestions:  issue.Suggestions,
			})
		}
	} else {
		// JSON解析失败，使用文本解析作为后备方案
		issues = s.extractIssuesFromText(aiResult, originalContent)
	}

	return issues, nil
}

// extractIssuesFromText 从文本中提取问题（后备方案）
func (s *ProofreadService) extractIssuesFromText(aiResult, originalContent string) []dto.Issue {
	var issues []dto.Issue

	// 简单的文本解析逻辑
	lines := strings.Split(aiResult, "\n")
	currentIssue := &dto.Issue{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") {
			if currentIssue.Message != "" {
				issues = append(issues, *currentIssue)
				currentIssue = &dto.Issue{}
			}

			// 提取问题描述
			re := regexp.MustCompile(`^[-•]\s*\[?(\w+)\]?\s*(.+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 3 {
				currentIssue.Type = "grammar" // 默认类型
				currentIssue.Severity = "warning"
				currentIssue.Message = matches[2]
			}
		}
	}

	if currentIssue.Message != "" {
		issues = append(issues, *currentIssue)
	}

	return issues
}

// findPositionInText 在文本中查找指定位置
func findPositionInText(content, searchTerm string, line, column int) int {
	lines := strings.Split(content, "\n")
	if line > 0 && line <= len(lines) {
		targetLine := lines[line-1]
		if column > 0 && column <= len([]rune(targetLine)) {
			// 计算绝对位置
			pos := 0
			for i := 0; i < line-1; i++ {
				pos += len([]rune(lines[i])) + 1 // +1 for newline
			}
			pos += column - 1
			return pos
		}
	}
	return -1
}

// generateStatistics 生成统计信息
func (s *ProofreadService) generateStatistics(issues []dto.Issue, content string) dto.ProofreadStats {
	stats := dto.ProofreadStats{
		TotalIssues:  len(issues),
		IssuesByType: make(map[string]int),
	}

	for _, issue := range issues {
		stats.IssuesByType[issue.Type]++

		switch issue.Severity {
		case "error":
			stats.ErrorCount++
		case "warning":
			stats.WarningCount++
		case "suggestion":
			stats.SuggestionCount++
		}
	}

	// 统计词数和字符数
	words := strings.Fields(content)
	stats.WordCount = len(words)
	stats.CharacterCount = len([]rune(content))

	return stats
}

// calculateScore 计算整体评分
func (s *ProofreadService) calculateScore(stats dto.ProofreadStats) float64 {
	if stats.TotalIssues == 0 {
		return 100.0
	}

	// 基础分100分，根据问题扣分
	score := 100.0

	// 错误扣分较多
	score -= float64(stats.ErrorCount) * 5

	// 警告扣分中等
	score -= float64(stats.WarningCount) * 2

	// 建议扣分较少
	score -= float64(stats.SuggestionCount) * 0.5

	// 确保分数在0-100之间
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}
