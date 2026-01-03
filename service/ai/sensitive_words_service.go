package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"Qingyu_backend/service/ai/adapter"
	"Qingyu_backend/service/ai/dto"
	"github.com/google/uuid"
)

// SensitiveWordsService 敏感词检测服务
type SensitiveWordsService struct {
	adapterManager *adapter.AdapterManager
	wordLibrary    *SensitiveWordLibrary
	// 可以添加存储层来保存检测结果
	// repository SensitiveWordsRepository
}

// SensitiveWordLibrary 敏感词库
type SensitiveWordLibrary struct {
	mu             sync.RWMutex
	politicalWords []string
	violenceWords  []string
	adultWords     []string
	customWords    map[string][]string // 按用户ID组织的自定义词库
}

// NewSensitiveWordLibrary 创建敏感词库
func NewSensitiveWordLibrary() *SensitiveWordLibrary {
	return &SensitiveWordLibrary{
		politicalWords: []string{
			// 示例政治敏感词（实际应用中应该从配置或数据库加载）
			// "敏感词1", "敏感词2",
		},
		violenceWords: []string{
			// 暴力词汇示例
			// "暴力词1", "暴力词2",
		},
		adultWords: []string{
			// 成人内容词汇示例
			// "成人词1", "成人词2",
		},
		customWords: make(map[string][]string),
	}
}

// NewSensitiveWordsService 创建敏感词检测服务
func NewSensitiveWordsService(adapterManager *adapter.AdapterManager) *SensitiveWordsService {
	return &SensitiveWordsService{
		adapterManager: adapterManager,
		wordLibrary:    NewSensitiveWordLibrary(),
	}
}

// CheckSensitiveWords 检测敏感词
func (s *SensitiveWordsService) CheckSensitiveWords(ctx context.Context, req *dto.SensitiveWordsCheckRequest) (*dto.SensitiveWordsCheckResponse, error) {
	// 参数验证
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("内容不能为空")
	}

	// 生成检查ID
	checkID := uuid.New().String()

	// 收集所有需要检测的敏感词
	searchWords := s.collectSearchWords(req)

	// 执行敏感词检测
	matches := s.detectSensitiveWords(req.Content, searchWords)

	// 生成摘要
	summary := s.generateCheckSummary(matches)

	// 判断是否安全
	isSafe := len(matches) == 0 || !s.hasHighRiskWords(matches)

	// 如果需要，可以使用AI进行语义分析
	if len(matches) > 0 {
		aiMatches, err := s.aiSemanticAnalysis(ctx, req.Content, req)
		if err == nil {
			matches = append(matches, aiMatches...)
		}
	}

	response := &dto.SensitiveWordsCheckResponse{
		CheckID:        checkID,
		IsSafe:         isSafe,
		TotalMatches:   len(matches),
		SensitiveWords: matches,
		Summary:        summary,
		TokensUsed:     0,
		ProcessedAt:    time.Now(),
	}

	// TODO: 将检测结果保存到存储层
	// go s.saveCheckResult(checkID, req, response)

	return response, nil
}

// GetSensitiveWordsDetail 获取敏感词检测结果
func (s *SensitiveWordsService) GetSensitiveWordsDetail(ctx context.Context, checkID string) (*dto.SensitiveWordsDetail, error) {
	// TODO: 从存储层获取检测结果
	// 这里返回模拟数据
	return &dto.SensitiveWordsDetail{
		CheckID:    checkID,
		Content:    "示例内容",
		IsSafe:     true,
		Matches:    []dto.SensitiveWordMatch{},
		CustomWords: []string{},
		Summary: dto.CheckSummary{
			ByCategory:      map[string]int{},
			ByLevel:         map[string]int{},
			HighRiskCount:   0,
			MediumRiskCount: 0,
			LowRiskCount:    0,
		},
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}, nil
}

// collectSearchWords 收集所有需要检测的敏感词
func (s *SensitiveWordsService) collectSearchWords(req *dto.SensitiveWordsCheckRequest) map[string][]string {
	s.wordLibrary.mu.RLock()
	defer s.wordLibrary.mu.RUnlock()

	searchWords := make(map[string][]string)

	// 根据分类选择词库
	category := req.Category
	if category == "" || category == "all" {
		// 添加所有分类的词库
		searchWords["political"] = s.wordLibrary.politicalWords
		searchWords["violence"] = s.wordLibrary.violenceWords
		searchWords["adult"] = s.wordLibrary.adultWords
	} else {
		// 根据指定分类添加词库
		switch category {
		case "political":
			searchWords["political"] = s.wordLibrary.politicalWords
		case "violence":
			searchWords["violence"] = s.wordLibrary.violenceWords
		case "adult":
			searchWords["adult"] = s.wordLibrary.adultWords
		}
	}

	// 添加自定义敏感词
	if len(req.CustomWords) > 0 {
		searchWords["custom"] = req.CustomWords
	}

	return searchWords
}

// detectSensitiveWords 检测敏感词
func (s *SensitiveWordsService) detectSensitiveWords(content string, searchWords map[string][]string) []dto.SensitiveWordMatch {
	var matches []dto.SensitiveWordMatch

	for category, words := range searchWords {
		for _, word := range words {
			if word == "" {
				continue
			}

			// 查找所有匹配位置
			positions := s.findWordPositions(content, word)
			for _, pos := range positions {
				match := dto.SensitiveWordMatch{
					ID:       uuid.New().String(),
					Word:     word,
					Category: category,
					Level:    s.determineWordLevel(category, word),
					Position: pos,
					Context:  s.extractContext(content, pos.Start, pos.End),
				}

				// 添加建议修改
				match.Suggestion = s.generateSuggestion(word, category)

				matches = append(matches, match)
			}
		}
	}

	return matches
}

// findWordPositions 查找词的所有出现位置
func (s *SensitiveWordsService) findWordPositions(content, word string) []dto.TextPosition {
	var positions []dto.TextPosition

	runes := []rune(content)
	wordRunes := []rune(word)

	start := 0
	for {
		// 查找词的位置
		idx := strings.Index(string(runes[start:]), word)
		if idx == -1 {
			break
		}

		absStart := start + idx
		absEnd := absStart + len(wordRunes)

		// 计算行列信息
		line, column := s.calculateLineColumn(content, absStart)

		positions = append(positions, dto.TextPosition{
			Start:  absStart,
			End:    absEnd,
			Length: len(wordRunes),
			Line:   line,
			Column: column,
		})

		start = absEnd
	}

	return positions
}

// calculateLineColumn 计算行列位置
func (s *SensitiveWordsService) calculateLineColumn(content string, pos int) (line, column int) {
	runes := []rune(content)
	if pos >= len(runes) {
		pos = len(runes) - 1
	}

	line = 1
	column = 1

	for i := 0; i < pos; i++ {
		if runes[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}

	return line, column
}

// extractContext 提取上下文
func (s *SensitiveWordsService) extractContext(content string, start, end int) string {
	runes := []rune(content)
	contextStart := start - 50
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := end + 50
	if contextEnd > len(runes) {
		contextEnd = len(runes)
	}

	context := string(runes[contextStart:contextEnd])
	if contextStart > 0 {
		context = "..." + context
	}
	if contextEnd < len(runes) {
		context = context + "..."
	}

	return context
}

// determineWordLevel 确定词汇风险级别
func (s *SensitiveWordsService) determineWordLevel(category, word string) string {
	// 根据分类确定默认级别
	switch category {
	case "political":
		return "high"
	case "violence":
		return "medium"
	case "adult":
		return "high"
	case "custom":
		return "medium"
	default:
		return "low"
	}
}

// generateSuggestion 生成修改建议
func (s *SensitiveWordsService) generateSuggestion(word, category string) string {
	return fmt.Sprintf("建议修改或删除敏感词「%s」", word)
}

// generateCheckSummary 生成检测摘要
func (s *SensitiveWordsService) generateCheckSummary(matches []dto.SensitiveWordMatch) dto.CheckSummary {
	summary := dto.CheckSummary{
		ByCategory: make(map[string]int),
		ByLevel:    make(map[string]int),
	}

	for _, match := range matches {
		summary.ByCategory[match.Category]++
		summary.ByLevel[match.Level]++

		switch match.Level {
		case "high":
			summary.HighRiskCount++
		case "medium":
			summary.MediumRiskCount++
		case "low":
			summary.LowRiskCount++
		}
	}

	return summary
}

// hasHighRiskWords 检查是否有高风险词
func (s *SensitiveWordsService) hasHighRiskWords(matches []dto.SensitiveWordMatch) bool {
	for _, match := range matches {
		if match.Level == "high" {
			return true
		}
	}
	return false
}

// aiSemanticAnalysis 使用AI进行语义分析（可选）
func (s *SensitiveWordsService) aiSemanticAnalysis(ctx context.Context, content string, req *dto.SensitiveWordsCheckRequest) ([]dto.SensitiveWordMatch, error) {
	// 构建提示词
	prompt := fmt.Sprintf(`请分析以下内容是否存在潜在的敏感内容（包括政治、暴力、成人内容等）：

%s

如果发现敏感内容，请以JSON格式返回，包含以下字段：
- word: 敏感词
- category: 分类
- level: 风险级别
- reason: 原因说明
- suggestion: 修改建议

如果未发现敏感内容，请返回空数组。`, content)

	// 调用AI
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      prompt,
		Temperature: 0.3,
		MaxTokens:   1000,
	}

	_, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, err
	}

	// TODO: 解析AI返回的结果
	// 这里暂时返回空数组
	return []dto.SensitiveWordMatch{}, nil
}

// AddCustomWords 添加自定义敏感词（可选方法）
func (s *SensitiveWordsService) AddCustomWords(userID string, words []string) error {
	s.wordLibrary.mu.Lock()
	defer s.wordLibrary.mu.Unlock()

	if s.wordLibrary.customWords[userID] == nil {
		s.wordLibrary.customWords[userID] = []string{}
	}

	s.wordLibrary.customWords[userID] = append(s.wordLibrary.customWords[userID], words...)

	return nil
}

// RemoveCustomWords 移除自定义敏感词（可选方法）
func (s *SensitiveWordsService) RemoveCustomWords(userID string, words []string) error {
	s.wordLibrary.mu.Lock()
	defer s.wordLibrary.mu.Unlock()

	if s.wordLibrary.customWords[userID] == nil {
		return nil
	}

	// 创建需要删除的词的map
	wordMap := make(map[string]bool)
	for _, word := range words {
		wordMap[word] = true
	}

	// 过滤掉需要删除的词
	filtered := []string{}
	for _, word := range s.wordLibrary.customWords[userID] {
		if !wordMap[word] {
			filtered = append(filtered, word)
		}
	}

	s.wordLibrary.customWords[userID] = filtered

	return nil
}
