package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/service/ai/adapter"
	"Qingyu_backend/service/ai/dto"
)

// SummarizeService 内容总结服务
type SummarizeService struct {
	adapterManager *adapter.AdapterManager
}

// NewSummarizeService 创建内容总结服务
func NewSummarizeService(adapterManager *adapter.AdapterManager) *SummarizeService {
	return &SummarizeService{
		adapterManager: adapterManager,
	}
}

// SummarizeContent 总结文档内容
func (s *SummarizeService) SummarizeContent(ctx context.Context, req *dto.SummarizeRequest) (*dto.SummarizeResponse, error) {
	// 参数验证
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("内容不能为空")
	}

	// 构建总结提示词
	prompt := s.buildSummarizePrompt(req)

	// 设置默认参数
	maxTokens := 1000
	if req.MaxLength > 0 {
		maxTokens = req.MaxLength
	}

	// 调用AI生成摘要
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      prompt,
		Temperature: 0.5, // 总结任务使用中等温度
		MaxTokens:   maxTokens,
	}

	result, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("生成摘要失败: %w", err)
	}

	// 解析关键点
	keyPoints := s.extractKeyPoints(result.Text)

	// 计算统计信息
	originalLength := len([]rune(req.Content))
	summaryLength := len([]rune(result.Text))
	compressionRate := 0.0
	if originalLength > 0 {
		compressionRate = float64(summaryLength) / float64(originalLength)
	}

	return &dto.SummarizeResponse{
		Summary:         result.Text,
		KeyPoints:       keyPoints,
		OriginalLength:  originalLength,
		SummaryLength:   summaryLength,
		CompressionRate: compressionRate,
		TokensUsed:      result.Usage.TotalTokens,
		Model:           result.Model,
		ProcessedAt:     time.Now(),
	}, nil
}

// SummarizeChapter 总结章节内容
func (s *SummarizeService) SummarizeChapter(ctx context.Context, req *dto.ChapterSummaryRequest) (*dto.ChapterSummaryResponse, error) {
	// TODO: 从数据库获取章节内容
	// 这里暂时使用模拟数据
	chapterContent := "这是章节内容的占位符。实际实现需要从数据库获取。"
	chapterTitle := "章节标题占位符"

	// 构建章节总结提示词
	prompt := fmt.Sprintf(`请对以下小说章节进行总结分析：

章节标题：%s
章节内容：%s

请提供：
1. 章节摘要（100-200字）
2. 关键情节点（3-5个）
3. 情节发展大纲（分层结构）
4. 涉及的主要角色及其作用

请以JSON格式返回，包含summary, keyPoints, plotOutline, characters字段。`, chapterTitle, chapterContent)

	// 调用AI生成总结
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      prompt,
		Temperature: 0.5,
		MaxTokens:   1500,
	}

	result, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("生成章节总结失败: %w", err)
	}

	// TODO: 解析AI返回的JSON结构
	// 这里暂时返回简化版本
	return &dto.ChapterSummaryResponse{
		ChapterID:    req.ChapterID,
		ChapterTitle: chapterTitle,
		Summary:      result.Text,
		KeyPoints:    []string{"关键点1", "关键点2", "关键点3"},
		PlotOutline:  []dto.OutlineItem{},
		Characters:   []dto.CharacterMention{},
		TokensUsed:   result.Usage.TotalTokens,
		ProcessedAt:  time.Now(),
	}, nil
}

// buildSummarizePrompt 构建总结提示词
func (s *SummarizeService) buildSummarizePrompt(req *dto.SummarizeRequest) string {
	var promptBuilder strings.Builder

	switch req.SummaryType {
	case "brief":
		promptBuilder.WriteString("请为以下内容生成简洁摘要（50-100字）：\n\n")
	case "detailed":
		promptBuilder.WriteString("请为以下内容生成详细摘要（200-300字）：\n\n")
	case "keypoints":
		promptBuilder.WriteString("请提取以下内容的关键要点，并以列表形式总结：\n\n")
	default:
		promptBuilder.WriteString("请为以下内容生成摘要：\n\n")
	}

	promptBuilder.WriteString(req.Content)

	if req.IncludeQuotes {
		promptBuilder.WriteString("\n\n请在摘要中包含1-2句关键引用。")
	}

	return promptBuilder.String()
}

// extractKeyPoints 从摘要中提取关键点
func (s *SummarizeService) extractKeyPoints(summary string) []string {
	// 简单实现：按句号分割并选择较长的句子
	sentences := strings.Split(summary, "。")
	var keyPoints []string

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len([]rune(sentence)) > 10 { // 只保留长度大于10的句子
			keyPoints = append(keyPoints, sentence)
		}
		if len(keyPoints) >= 5 { // 最多返回5个关键点
			break
		}
	}

	return keyPoints
}
