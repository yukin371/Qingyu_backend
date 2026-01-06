package ai

import (
	"context"

	"Qingyu_backend/service/ai/dto"
)

// WritingAssistantService 写作辅助服务接口
// 提供内容总结、校对、敏感词检测等功能
type WritingAssistantService interface {
	// ===========================
	// 内容总结功能
	// ===========================

	// SummarizeContent 总结文档内容
	SummarizeContent(ctx context.Context, req *dto.SummarizeRequest) (*dto.SummarizeResponse, error)

	// SummarizeChapter 总结章节内容
	SummarizeChapter(ctx context.Context, req *dto.ChapterSummaryRequest) (*dto.ChapterSummaryResponse, error)

	// ===========================
	// 文本校对功能
	// ===========================

	// ProofreadContent 校对文本内容
	ProofreadContent(ctx context.Context, req *dto.ProofreadRequest) (*dto.ProofreadResponse, error)

	// GetProofreadSuggestion 获取校对建议详情
	GetProofreadSuggestion(ctx context.Context, suggestionID string) (*dto.ProofreadSuggestion, error)

	// ===========================
	// 敏感词检测功能
	// ===========================

	// CheckSensitiveWords 检测敏感词
	CheckSensitiveWords(ctx context.Context, req *dto.SensitiveWordsCheckRequest) (*dto.SensitiveWordsCheckResponse, error)

	// GetSensitiveWordsDetail 获取敏感词检测结果
	GetSensitiveWordsDetail(ctx context.Context, checkID string) (*dto.SensitiveWordsDetail, error)
}
