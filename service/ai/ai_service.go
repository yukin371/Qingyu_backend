package ai

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/models/ai"
	documentService "Qingyu_backend/service/document"
)

// Service AI服务
type Service struct {
	contextService     *ContextService
	externalAPIService *ExternalAPIService
}

// NewService 创建AI服务
func NewService() *Service {
	// 加载配置
	cfg := config.LoadConfig()
	
	// 创建document服务实例
	docService := &documentService.DocumentService{}
	projService := &documentService.ProjectService{}
	nodeService := &documentService.NodeService{}
	versionService := &documentService.VersionService{}
	
	// 创建上下文服务
	contextService := NewContextService(docService, projService, nodeService, versionService)
	
	// 创建外部API服务
	externalAPIService := NewExternalAPIService(cfg.AI)
	
	return &Service{
		contextService:     contextService,
		externalAPIService: externalAPIService,
	}
}

// GenerateContentRequest 生成内容请求
type GenerateContentRequest struct {
	ProjectID string                `json:"projectId"`
	ChapterID string                `json:"chapterId,omitempty"`
	Prompt    string                `json:"prompt"`
	Options   *ai.GenerateOptions   `json:"options,omitempty"`
}

// GenerateContentResponse 生成内容响应
type GenerateContentResponse struct {
	Content     string    `json:"content"`
	TokensUsed  int       `json:"tokensUsed"`
	Model       string    `json:"model"`
	GeneratedAt time.Time `json:"generatedAt"`
}

// GenerateContent 生成内容
func (s *Service) GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error) {
	// 构建AI上下文
	aiContext, err := s.contextService.BuildContext(ctx, req.ProjectID, req.ChapterID)
	if err != nil {
		return nil, fmt.Errorf("构建AI上下文失败: %w", err)
	}

	// 设置默认选项
	options := req.Options
	if options == nil {
		options = &ai.GenerateOptions{
			Temperature: 0.7,
			MaxTokens:   2000,
		}
	}

	// 调用外部AI API生成内容
	result, err := s.externalAPIService.GenerateContent(ctx, aiContext, req.Prompt, options)
	if err != nil {
		return nil, fmt.Errorf("生成内容失败: %w", err)
	}

	response := &GenerateContentResponse{
		Content:     result.Content,
		TokensUsed:  result.TokensUsed,
		Model:       result.Model,
		GeneratedAt: time.Now(),
	}

	return response, nil
}

// AnalyzeContentRequest 分析内容请求
type AnalyzeContentRequest struct {
	Content      string `json:"content"`
	AnalysisType string `json:"analysisType"` // plot, character, style, general
}

// AnalyzeContentResponse 分析内容响应
type AnalyzeContentResponse struct {
	Type        string    `json:"type"`
	Analysis    string    `json:"analysis"`
	TokensUsed  int       `json:"tokensUsed"`
	Model       string    `json:"model"`
	AnalyzedAt  time.Time `json:"analyzedAt"`
}

// AnalyzeContent 分析内容
func (s *Service) AnalyzeContent(ctx context.Context, req *AnalyzeContentRequest) (*AnalyzeContentResponse, error) {
	// 调用外部AI API分析内容
	result, err := s.externalAPIService.AnalyzeContent(ctx, req.Content, req.AnalysisType)
	if err != nil {
		return nil, fmt.Errorf("分析内容失败: %w", err)
	}

	response := &AnalyzeContentResponse{
		Type:       result.Type,
		Analysis:   result.Analysis,
		TokensUsed: result.TokensUsed,
		Model:      result.Model,
		AnalyzedAt: time.Now(),
	}

	return response, nil
}

// ContinueWritingRequest 续写请求
type ContinueWritingRequest struct {
	ProjectID     string              `json:"projectId"`
	ChapterID     string              `json:"chapterId"`
	CurrentText   string              `json:"currentText"`
	ContinueLength int                `json:"continueLength,omitempty"` // 续写长度（字数）
	Options       *ai.GenerateOptions `json:"options,omitempty"`
}

// ContinueWriting 续写内容
func (s *Service) ContinueWriting(ctx context.Context, req *ContinueWritingRequest) (*GenerateContentResponse, error) {
	// 构建续写提示词
	prompt := fmt.Sprintf("请基于以下内容进行续写，保持风格和情节的连贯性：\n\n%s", req.CurrentText)
	
	if req.ContinueLength > 0 {
		prompt += fmt.Sprintf("\n\n请续写约%d字的内容。", req.ContinueLength)
	}

	// 调用生成内容方法
	generateReq := &GenerateContentRequest{
		ProjectID: req.ProjectID,
		ChapterID: req.ChapterID,
		Prompt:    prompt,
		Options:   req.Options,
	}

	return s.GenerateContent(ctx, generateReq)
}

// OptimizeTextRequest 文本优化请求
type OptimizeTextRequest struct {
	ProjectID      string              `json:"projectId"`
	ChapterID      string              `json:"chapterId,omitempty"`
	OriginalText   string              `json:"originalText"`
	OptimizeType   string              `json:"optimizeType"`   // grammar, style, flow, dialogue
	Instructions   string              `json:"instructions,omitempty"` // 具体优化指示
	Options        *ai.GenerateOptions `json:"options,omitempty"`
}

// OptimizeText 优化文本
func (s *Service) OptimizeText(ctx context.Context, req *OptimizeTextRequest) (*GenerateContentResponse, error) {
	// 构建优化提示词
	var prompt string
	switch req.OptimizeType {
	case "grammar":
		prompt = "请修正以下文本的语法错误，保持原意不变："
	case "style":
		prompt = "请优化以下文本的写作风格，使其更加流畅自然："
	case "flow":
		prompt = "请优化以下文本的逻辑流程和段落结构："
	case "dialogue":
		prompt = "请优化以下文本中的对话，使其更加生动自然："
	default:
		prompt = "请对以下文本进行综合优化："
	}

	if req.Instructions != "" {
		prompt += fmt.Sprintf("\n\n具体要求：%s", req.Instructions)
	}

	prompt += fmt.Sprintf("\n\n原文：\n%s", req.OriginalText)

	// 调用生成内容方法
	generateReq := &GenerateContentRequest{
		ProjectID: req.ProjectID,
		ChapterID: req.ChapterID,
		Prompt:    prompt,
		Options:   req.Options,
	}

	return s.GenerateContent(ctx, generateReq)
}

// GenerateOutlineRequest 生成大纲请求
type GenerateOutlineRequest struct {
	ProjectID   string              `json:"projectId"`
	Theme       string              `json:"theme"`       // 主题
	Genre       string              `json:"genre"`       // 类型
	Length      string              `json:"length"`      // 长度（短篇、中篇、长篇）
	KeyElements []string            `json:"keyElements"` // 关键元素
	Options     *ai.GenerateOptions `json:"options,omitempty"`
}

// GenerateOutline 生成大纲
func (s *Service) GenerateOutline(ctx context.Context, req *GenerateOutlineRequest) (*GenerateContentResponse, error) {
	// 构建大纲生成提示词
	prompt := fmt.Sprintf("请为以下小说创作一个详细的大纲：\n\n主题：%s\n类型：%s\n长度：%s", 
		req.Theme, req.Genre, req.Length)

	if len(req.KeyElements) > 0 {
		prompt += "\n\n关键元素："
		for _, element := range req.KeyElements {
			prompt += fmt.Sprintf("\n- %s", element)
		}
	}

	prompt += "\n\n请包含以下内容：\n1. 故事概述\n2. 主要角色设定\n3. 章节大纲\n4. 主要情节线\n5. 关键转折点"

	// 调用生成内容方法
	generateReq := &GenerateContentRequest{
		ProjectID: req.ProjectID,
		Prompt:    prompt,
		Options:   req.Options,
	}

	return s.GenerateContent(ctx, generateReq)
}

// GetContextInfo 获取上下文信息
func (s *Service) GetContextInfo(ctx context.Context, projectID, chapterID string) (*ai.AIContext, error) {
	return s.contextService.BuildContext(ctx, projectID, chapterID)
}

// UpdateContextWithFeedback 根据反馈更新上下文
func (s *Service) UpdateContextWithFeedback(ctx context.Context, projectID, chapterID, feedback string) error {
	// 构建上下文
	aiContext, err := s.contextService.BuildContext(ctx, projectID, chapterID)
	if err != nil {
		return fmt.Errorf("构建AI上下文失败: %w", err)
	}

	// 更新上下文
	return s.contextService.UpdateContextWithFeedback(ctx, aiContext, feedback)
}