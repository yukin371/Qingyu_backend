package ai

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/models/ai"
	"Qingyu_backend/service/ai/adapter"
	documentService "Qingyu_backend/service/writer/project"

	pb "Qingyu_backend/pkg/grpc/pb" // 假设proto路径

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service AI服务
type Service struct {
	contextService *ContextService
	adapterManager *adapter.AdapterManager
	PythonConfig   *config.PythonAIServiceConfig // 新增：Python配置
}

// NewService 创建AI服务（向后兼容，但不推荐使用）
// 废弃：请使用 NewServiceWithDependencies
func NewService() *Service {
	// 使用全局配置
	cfg := config.GlobalConfig
	if cfg == nil {
		panic("GlobalConfig is not initialized")
	}

	// 创建document服务实例（空的，没有依赖注入）
	// 警告：这会导致nil pointer错误！
	docService := &documentService.DocumentService{}
	projService := &documentService.ProjectService{}
	nodeService := &documentService.NodeService{}
	versionService := &documentService.VersionService{}

	// 创建上下文服务
	// TODO: 这里使用nil是因为旧架构没有依赖注入，需要重构为使用RepositoryFactory
	contextService := NewContextService(docService, projService, nodeService, versionService, nil)

	// 创建适配器管理器 - 使用简化的配置
	var adapterManager *adapter.AdapterManager
	if cfg.AI != nil {
		// 创建一个简化的ExternalAPIConfig
		externalAPIConfig := &config.ExternalAPIConfig{
			DefaultProvider: "openai",
			Providers: map[string]*config.ProviderConfig{
				"openai": {
					Name:            "openai",
					APIKey:          cfg.AI.APIKey,
					BaseURL:         cfg.AI.BaseURL,
					Priority:        1,
					Enabled:         true,
					SupportedModels: []string{"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"},
				},
			},
		}
		adapterManager = adapter.NewAdapterManager(externalAPIConfig)
	}

	return &Service{
		contextService: contextService,
		adapterManager: adapterManager,
		PythonConfig:   cfg.AI.PythonService, // 注入Python配置
	}
}

// NewServiceWithDependencies 创建AI服务（使用依赖注入，推荐）
func NewServiceWithDependencies(projectService *documentService.ProjectService) *Service {
	// 使用全局配置
	cfg := config.GlobalConfig
	if cfg == nil {
		panic("GlobalConfig is not initialized")
	}

	// 创建其他document服务实例（暂时仍为空，待后续迁移）
	docService := &documentService.DocumentService{}
	nodeService := &documentService.NodeService{}
	versionService := &documentService.VersionService{}

	// 创建上下文服务，使用注入的ProjectService
	contextService := NewContextService(docService, projectService, nodeService, versionService, nil)

	// 创建适配器管理器 - 优先使用External配置（支持多提供商）
	var adapterManager *adapter.AdapterManager
	if cfg.External != nil {
		adapterManager = adapter.NewAdapterManager(cfg.External)
	} else {
		// 默认空配置（应该不会走到这里，因为config.test.yaml有External配置）
		fmt.Println("警告: 未找到External API配置，AI功能可能无法使用")
	}

	return &Service{
		contextService: contextService,
		adapterManager: adapterManager,
		PythonConfig:   cfg.AI.PythonService, // 注入Python配置
	}
}

// GenerateContentRequest 生成内容请求
type GenerateContentRequest struct {
	ProjectID string              `json:"projectId"`
	ChapterID string              `json:"chapterId,omitempty"`
	Prompt    string              `json:"prompt"`
	Options   *ai.GenerateOptions `json:"options,omitempty"`
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
	_, err := s.contextService.BuildContext(ctx, req.ProjectID, req.ChapterID)
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

	// 检查适配器管理器是否已初始化
	if s.adapterManager == nil {
		return nil, fmt.Errorf("AI适配器管理器未初始化，请检查配置文件中的External API配置")
	}

	// 新增：如果Python服务配置可用，优先使用gRPC调用Python AI服务
	if s.PythonConfig != nil && s.PythonConfig.GrpcPort > 0 {
		pythonAddr := fmt.Sprintf("%s:%d", s.PythonConfig.Host, s.PythonConfig.GrpcPort)
		conn, err := grpc.NewClient(pythonAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("警告: 无法连接Python AI服务 (%s)，回退到本地适配器: %v\n", pythonAddr, err)
		} else {
			defer conn.Close()
			client := pb.NewAIServiceClient(conn)

			// 示例gRPC调用 (假设proto定义了GenerateContent方法)
			grpcReq := &pb.GenerateContentRequest{
				Prompt:    req.Prompt,
				ProjectId: req.ProjectID,
				ChapterId: req.ChapterID,
				Options: &pb.GenerateOptions{
					Temperature: float32(options.Temperature),
					MaxTokens:   int32(options.MaxTokens),
					Model:       options.Model,
				},
			}
			grpcResp, err := client.GenerateContent(ctx, grpcReq)
			if err != nil {
				fmt.Printf("gRPC调用Python AI失败，回退到本地: %v\n", err)
			} else {
				// 使用gRPC响应
				return &GenerateContentResponse{
					Content:     grpcResp.Content,
					TokensUsed:  int(grpcResp.TokensUsed),
					Model:       grpcResp.Model,
					GeneratedAt: time.Now(),
				}, nil
			}
		}
	}

	// 回退：使用原有适配器逻辑
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      req.Prompt,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Model:       options.Model,
	}

	result, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("生成内容失败: %w", err)
	}

	response := &GenerateContentResponse{
		Content:     result.Text,
		TokensUsed:  result.Usage.TotalTokens,
		Model:       result.Model,
		GeneratedAt: result.CreatedAt,
	}

	return response, nil
}

// GenerateContentStream 流式生成内容
func (s *Service) GenerateContentStream(ctx context.Context, req *GenerateContentRequest) (<-chan *GenerateContentResponse, error) {
	// 构建AI上下文
	_, err := s.contextService.BuildContext(ctx, req.ProjectID, req.ChapterID)
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

	// 使用适配器管理器进行流式生成
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      req.Prompt,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Model:       options.Model,
		Stream:      true, // 启用流式响应
	}

	// 获取流式响应通道
	streamChan, err := s.adapterManager.AutoTextGenerationStream(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("流式生成内容失败: %w", err)
	}

	// 创建响应通道
	responseChan := make(chan *GenerateContentResponse, 10)

	// 转换流式响应格式
	go func() {
		defer close(responseChan)

		for result := range streamChan {
			if result == nil {
				continue
			}

			response := &GenerateContentResponse{
				Content:     result.Text,
				TokensUsed:  result.Usage.TotalTokens,
				Model:       result.Model,
				GeneratedAt: result.CreatedAt,
			}

			select {
			case responseChan <- response:
			case <-ctx.Done():
				return
			}
		}
	}()

	return responseChan, nil
}

// AnalyzeContentRequest 分析内容请求
type AnalyzeContentRequest struct {
	Content      string `json:"content"`
	AnalysisType string `json:"analysisType"` // plot, character, style, general
}

// AnalyzeContentResponse 分析内容响应
type AnalyzeContentResponse struct {
	Type       string    `json:"type"`
	Analysis   string    `json:"analysis"`
	TokensUsed int       `json:"tokensUsed"`
	Model      string    `json:"model"`
	AnalyzedAt time.Time `json:"analyzedAt"`
}

// AnalyzeContent 分析内容
func (s *Service) AnalyzeContent(ctx context.Context, req *AnalyzeContentRequest) (*AnalyzeContentResponse, error) {
	// 构建分析提示词
	prompt := fmt.Sprintf("请对以下内容进行%s分析：\n\n%s", req.AnalysisType, req.Content)

	// 使用适配器管理器生成分析
	adapterReq := &adapter.TextGenerationRequest{
		Prompt:      prompt,
		Temperature: 0.3, // 分析任务使用较低的温度
		MaxTokens:   1500,
	}

	result, err := s.adapterManager.AutoTextGeneration(ctx, adapterReq)
	if err != nil {
		return nil, fmt.Errorf("分析内容失败: %w", err)
	}

	response := &AnalyzeContentResponse{
		Type:       req.AnalysisType,
		Analysis:   result.Text,
		TokensUsed: result.Usage.TotalTokens,
		Model:      result.Model,
		AnalyzedAt: result.CreatedAt,
	}

	return response, nil
}

// ContinueWritingRequest 续写请求
type ContinueWritingRequest struct {
	ProjectID      string              `json:"projectId"`
	ChapterID      string              `json:"chapterId"`
	CurrentText    string              `json:"currentText"`
	ContinueLength int                 `json:"continueLength,omitempty"` // 续写长度（字数）
	Options        *ai.GenerateOptions `json:"options,omitempty"`
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
	ProjectID    string              `json:"projectId"`
	ChapterID    string              `json:"chapterId,omitempty"`
	OriginalText string              `json:"originalText"`
	OptimizeType string              `json:"optimizeType"`           // grammar, style, flow, dialogue
	Instructions string              `json:"instructions,omitempty"` // 具体优化指示
	Options      *ai.GenerateOptions `json:"options,omitempty"`
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

// GetAdapterManager 获取适配器管理器
func (s *Service) GetAdapterManager() *adapter.AdapterManager {
	return s.adapterManager
}
