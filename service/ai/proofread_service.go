package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"Qingyu_backend/service/ai/dto"

	"github.com/google/uuid"
)

// ProofreadService 文本校对服务
type ProofreadService struct {
	textGenerator TextGenerator
	cache         *proofreadResultCache
}

// NewProofreadService 创建文本校对服务
func NewProofreadService(textGenerator TextGenerator) *ProofreadService {
	return &ProofreadService{
		textGenerator: textGenerator,
		cache:         newProofreadResultCache(10 * time.Minute),
	}
}

// ProofreadContent 校对文本内容
func (s *ProofreadService) ProofreadContent(ctx context.Context, req *dto.ProofreadRequest) (*dto.ProofreadResponse, error) {
	// 参数验证
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("内容不能为空")
	}

	// 设置默认检查类型
	checkTypes := req.CheckTypes //nolint:ineffassign // 可能在后续被重新赋值
	if len(checkTypes) == 0 {
		checkTypes = []string{"spelling", "grammar", "punctuation"}
	}
	req.CheckTypes = normalizeProofreadCheckTypes(checkTypes)
	contentHash := hashProofreadContent(req.Content)
	cacheKey := buildProofreadCacheKey(contentHash, req)
	if cached, ok := s.cache.Get(cacheKey); ok {
		cached.ProcessedAt = time.Now()
		return cached, nil
	}

	// 构建校对提示词
	prompt := s.buildProofreadPrompt(req)

	// 调用AI进行校对
	generateReq := &TextGenerateRequest{
		Prompt:      prompt,
		Temperature: 0.3, // 校对任务使用较低温度以确保准确性
		MaxTokens:   2000,
	}

	result, err := s.textGenerator.GenerateText(ctx, generateReq)
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

	response := &dto.ProofreadResponse{
		ReviewID:        uuid.New().String(),
		ContentHash:     contentHash,
		OriginalContent: req.Content,
		Issues:          issues,
		Score:           score,
		Statistics:      stats,
		PreviewWarnings: s.generateMobilePreviewWarnings(req.Content),
		TokensUsed:      result.TokensUsed,
		Model:           result.Model,
		ProcessedAt:     time.Now(),
	}
	s.cache.Set(cacheKey, response)
	return response, nil
}

type proofreadCacheEntry struct {
	response  *dto.ProofreadResponse
	expiresAt time.Time
}

type proofreadResultCache struct {
	mu      sync.RWMutex
	ttl     time.Duration
	records map[string]proofreadCacheEntry
}

func newProofreadResultCache(ttl time.Duration) *proofreadResultCache {
	return &proofreadResultCache{
		ttl:     ttl,
		records: make(map[string]proofreadCacheEntry),
	}
}

func (c *proofreadResultCache) Get(key string) (*dto.ProofreadResponse, bool) {
	if c == nil {
		return nil, false
	}

	c.mu.RLock()
	entry, ok := c.records[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.records, key)
		c.mu.Unlock()
		return nil, false
	}

	return cloneProofreadResponse(entry.response), true
}

func (c *proofreadResultCache) Set(key string, response *dto.ProofreadResponse) {
	if c == nil || response == nil {
		return
	}

	c.mu.Lock()
	c.records[key] = proofreadCacheEntry{
		response:  cloneProofreadResponse(response),
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

func buildProofreadCacheKey(contentHash string, req *dto.ProofreadRequest) string {
	parts := []string{
		contentHash,
		strings.Join(req.CheckTypes, ","),
		strings.TrimSpace(req.Genre),
		strings.TrimSpace(req.ReaderProfile),
		strings.TrimSpace(req.Mode),
	}
	return strings.Join(parts, "|")
}

func cloneProofreadResponse(response *dto.ProofreadResponse) *dto.ProofreadResponse {
	if response == nil {
		return nil
	}

	clone := *response
	clone.Issues = append([]dto.Issue(nil), response.Issues...)
	for index := range clone.Issues {
		clone.Issues[index].Suggestions = append([]string(nil), response.Issues[index].Suggestions...)
		clone.Issues[index].SuggestionDetails = append([]dto.SuggestionItem(nil), response.Issues[index].SuggestionDetails...)
	}
	clone.PreviewWarnings = append([]dto.MobilePreviewWarning(nil), response.PreviewWarnings...)
	return &clone
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

	promptBuilder.WriteString("你是中文小说章节审校助手。请对以下文本进行校对，找出错误并提供修改建议。\n")
	promptBuilder.WriteString("只返回 JSON，不要返回 Markdown 代码块，不要整段重写正文。\n\n")
	promptBuilder.WriteString(req.Content)
	promptBuilder.WriteString("\n\n请检查以下方面：\n")

	if len(req.CheckTypes) == 0 {
		promptBuilder.WriteString("- 拼写错误\n")
		promptBuilder.WriteString("- 语法错误\n")
		promptBuilder.WriteString("- 标点符号错误\n")
	} else {
		for _, checkType := range req.CheckTypes {
			switch checkType {
			case "spelling", "typo":
				promptBuilder.WriteString("- 错别字、同音误写、形近误写\n")
			case "grammar":
				promptBuilder.WriteString("- 病句、语义不顺、主谓搭配问题\n")
			case "punctuation":
				promptBuilder.WriteString("- 标点符号错误\n")
			case "style":
				promptBuilder.WriteString("- 重复表达、口水词、风格不一致\n")
			case "readability":
				promptBuilder.WriteString("- 移动端阅读体验、句子过长、段落过长\n")
			}
		}
	}

	if req.Genre != "" {
		promptBuilder.WriteString("\n题材参考：")
		promptBuilder.WriteString(req.Genre)
	}
	if req.ReaderProfile != "" {
		promptBuilder.WriteString("\n读者画像：")
		promptBuilder.WriteString(req.ReaderProfile)
	}

	promptBuilder.WriteString("\n\n请以JSON格式返回结果，schema如下：\n")
	promptBuilder.WriteString(`{"score":86,"issues":[{"type":"typo","severity":"error","message":"问题说明","start":0,"end":2,"originalText":"原文","suggestions":[{"text":"建议文本","reason":"原因","confidence":0.9}]}]}`)
	promptBuilder.WriteString("\nseverity 只能使用 error、warning、suggestion。type 只能使用 typo、grammar、punctuation、style、readability。")

	return promptBuilder.String()
}

// parseProofreadResult 解析AI返回的校对结果
func (s *ProofreadService) parseProofreadResult(aiResult, originalContent string) ([]dto.Issue, error) {
	var issues []dto.Issue

	// 尝试解析JSON格式
	var result struct {
		Score  float64 `json:"score"`
		Issues []struct {
			Type         string        `json:"type"`
			Severity     string        `json:"severity"`
			Message      string        `json:"message"`
			Line         int           `json:"line"`
			Column       int           `json:"column"`
			Start        *int          `json:"start"`
			End          *int          `json:"end"`
			Original     string        `json:"original"`
			OriginalText string        `json:"originalText"`
			Suggestions  []interface{} `json:"suggestions"`
		} `json:"issues"`
	}

	if err := json.Unmarshal([]byte(aiResult), &result); err == nil {
		// 成功解析JSON
		for _, issue := range result.Issues {
			originalText := strings.TrimSpace(issue.OriginalText)
			if originalText == "" {
				originalText = strings.TrimSpace(issue.Original)
			}

			// 查找原文中的位置
			start := -1
			if issue.Start != nil && issue.End != nil && *issue.Start >= 0 && *issue.End > *issue.Start {
				start = *issue.Start
			} else {
				start = findPositionInText(originalContent, originalText, issue.Line, issue.Column)
			}
			if start == -1 {
				continue
			}
			end := start + len([]rune(originalText))
			if issue.End != nil && *issue.End > start {
				end = *issue.End
			}
			if start < 0 || end > len([]rune(originalContent)) || end <= start {
				continue
			}

			suggestions := make([]string, 0, len(issue.Suggestions))
			suggestionDetails := make([]dto.SuggestionItem, 0, len(issue.Suggestions))
			for _, suggestion := range issue.Suggestions {
				switch value := suggestion.(type) {
				case string:
					if strings.TrimSpace(value) != "" {
						text := strings.TrimSpace(value)
						suggestions = append(suggestions, text)
						suggestionDetails = append(suggestionDetails, dto.SuggestionItem{Text: text})
					}
				case map[string]interface{}:
					if textValue, ok := value["text"].(string); ok && strings.TrimSpace(textValue) != "" {
						item := dto.SuggestionItem{
							Text:   strings.TrimSpace(textValue),
							Reason: stringFromMap(value, "reason"),
						}
						if confidence, ok := numberFromMap(value, "confidence"); ok {
							item.Confidence = confidence
						}
						suggestions = append(suggestions, item.Text)
						suggestionDetails = append(suggestionDetails, item)
					}
				}
			}
			issues = append(issues, dto.Issue{
				ID:       uuid.New().String(),
				Type:     normalizeProofreadIssueType(issue.Type),
				Severity: normalizeProofreadSeverity(issue.Severity),
				Message:  issue.Message,
				Position: dto.TextPosition{
					Line:   issue.Line,
					Column: issue.Column,
					Start:  start,
					End:    end,
					Length: end - start,
				},
				OriginalText:      originalText,
				Suggestions:       suggestions,
				SuggestionDetails: suggestionDetails,
			})
		}
	} else {
		// JSON解析失败，使用文本解析作为后备方案
		issues = s.extractIssuesFromText(aiResult, originalContent)
	}

	return issues, nil
}

func stringFromMap(record map[string]interface{}, key string) string {
	if value, ok := record[key].(string); ok {
		return strings.TrimSpace(value)
	}
	return ""
}

func numberFromMap(record map[string]interface{}, key string) (float64, bool) {
	switch value := record[key].(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		return float64(value), true
	case json.Number:
		number, err := value.Float64()
		return number, err == nil
	default:
		return 0, false
	}
}

func normalizeProofreadCheckTypes(checkTypes []string) []string {
	normalized := make([]string, 0, len(checkTypes))
	for _, checkType := range checkTypes {
		switch strings.ToLower(strings.TrimSpace(checkType)) {
		case "typo", "spelling":
			normalized = append(normalized, "spelling")
		case "grammar", "punctuation", "style", "readability":
			normalized = append(normalized, strings.ToLower(strings.TrimSpace(checkType)))
		}
	}
	if len(normalized) == 0 {
		return []string{"spelling", "grammar", "punctuation"}
	}
	return normalized
}

func normalizeProofreadIssueType(issueType string) string {
	switch strings.ToLower(strings.TrimSpace(issueType)) {
	case "spelling", "typo":
		return "typo"
	case "grammar", "punctuation", "style", "readability", "continuity":
		return strings.ToLower(strings.TrimSpace(issueType))
	default:
		return "style"
	}
}

func normalizeProofreadSeverity(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "error", "high", "critical":
		return "error"
	case "warning", "warn", "medium":
		return "warning"
	default:
		return "suggestion"
	}
}

func hashProofreadContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (s *ProofreadService) generateMobilePreviewWarnings(content string) []dto.MobilePreviewWarning {
	paragraphs := splitNonEmptyParagraphs(content)
	warnings := make([]dto.MobilePreviewWarning, 0)
	searchByteOffset := 0
	for index, paragraph := range paragraphs {
		paragraphLength := len([]rune(paragraph))
		start := 0
		byteStart := strings.Index(content[searchByteOffset:], paragraph)
		if byteStart >= 0 {
			absoluteByteStart := searchByteOffset + byteStart
			start = len([]rune(content[:absoluteByteStart]))
			searchByteOffset = absoluteByteStart + len(paragraph)
		}
		if paragraphLength >= 220 {
			warnings = append(warnings, dto.MobilePreviewWarning{
				ID:             uuid.New().String(),
				Type:           "paragraph_too_long",
				Severity:       "warning",
				Message:        fmt.Sprintf("第 %d 段在手机上偏长，建议拆分以降低阅读压力。", index+1),
				ParagraphIndex: index,
				Position: dto.TextPosition{
					Start:  start,
					End:    start + paragraphLength,
					Length: paragraphLength,
				},
			})
		}
		if containsLongSentence(paragraph, 90) {
			warnings = append(warnings, dto.MobilePreviewWarning{
				ID:             uuid.New().String(),
				Type:           "sentence_too_long",
				Severity:       "suggestion",
				Message:        fmt.Sprintf("第 %d 段存在较长句，手机端阅读可能需要断句。", index+1),
				ParagraphIndex: index,
			})
		}
	}
	return warnings
}

func splitNonEmptyParagraphs(content string) []string {
	rawParagraphs := regexp.MustCompile(`\n\s*\n`).Split(content, -1)
	paragraphs := make([]string, 0, len(rawParagraphs))
	for _, paragraph := range rawParagraphs {
		trimmed := strings.TrimSpace(paragraph)
		if trimmed != "" {
			paragraphs = append(paragraphs, trimmed)
		}
	}
	return paragraphs
}

func containsLongSentence(paragraph string, threshold int) bool {
	sentences := regexp.MustCompile(`[。！？!?；;]`).Split(paragraph, -1)
	for _, sentence := range sentences {
		if len([]rune(strings.TrimSpace(sentence))) >= threshold {
			return true
		}
	}
	return false
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
