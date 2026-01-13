package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"
	"Qingyu_backend/service/ai/adapter"
	"Qingyu_backend/service/ai/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatRepositoryInterface 聊天仓库接口
type ChatRepositoryInterface interface {
	CreateSession(ctx context.Context, session *aiModels.ChatSession) error                                           // 创建会话
	GetSessionByID(ctx context.Context, sessionID string) (*aiModels.ChatSession, error)                              // 获取会话
	UpdateSession(ctx context.Context, session *aiModels.ChatSession) error                                           // 更新会话
	DeleteSession(ctx context.Context, sessionID string) error                                                        // 删除会话
	GetSessionsByProjectID(ctx context.Context, projectID string, limit, offset int) ([]*aiModels.ChatSession, error) // 获取项目会话列表
	CreateMessage(ctx context.Context, message *aiModels.ChatMessage) error                                           // 创建消息
	GetSessionStatistics(ctx context.Context, projectID string) (*ChatStatistics, error)                              // 获取会话统计信息
}

// AIServiceInterface AI服务接口
type AIServiceInterface interface {
	GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)
	GenerateContentStream(ctx context.Context, req *GenerateContentRequest) (<-chan *GenerateContentResponse, error)
}

// ChatService AI聊天服务
type ChatService struct {
	aiService      AIServiceInterface
	adapterManager *adapter.AdapterManager
	repository     ChatRepositoryInterface
}

// NewChatService 创建聊天服务
func NewChatService(aiService *Service, repository ChatRepositoryInterface) *ChatService {
	return &ChatService{
		aiService:      aiService,
		adapterManager: aiService.adapterManager,
		repository:     repository,
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID   string                    `json:"sessionId,omitempty"`
	ProjectID   string                    `json:"projectId,omitempty"`
	Message     string                    `json:"message" binding:"required"`
	UseContext  bool                      `json:"useContext"`
	ContextType string                    `json:"contextType,omitempty"` // novel, general
	Options     *aiModels.GenerateOptions `json:"options,omitempty"`
	Metadata    map[string]interface{}    `json:"metadata,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	SessionID    string              `json:"sessionId"`
	Message      *dto.ChatMessageDTO `json:"message"`
	TokensUsed   int                 `json:"tokensUsed"`
	Model        string              `json:"model"`
	ContextUsed  bool                `json:"contextUsed"`
	ResponseTime time.Duration       `json:"responseTime"`
}

// StreamChatResponse 流式聊天响应
type StreamChatResponse struct {
	SessionID    string        `json:"sessionId"`
	MessageID    string        `json:"messageId"`
	Content      string        `json:"content"`
	Delta        string        `json:"delta"`
	TokensUsed   int           `json:"tokensUsed"`
	Model        string        `json:"model"`
	ContextUsed  bool          `json:"contextUsed"`
	IsComplete   bool          `json:"isComplete"`
	ResponseTime time.Duration `json:"responseTime"`
}

// StartChat 开始聊天
func (s *ChatService) StartChat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 验证请求
	if req.Message == "" {
		return nil, errors.New("消息内容不能为空")
	}

	// 获取或创建会话
	session, err := s.getOrCreateSession(ctx, req.SessionID, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("获取会话失败: %w", err)
	}

	// 构建聊天上下文
	context, err := s.buildChatContext(ctx, session, req)
	if err != nil {
		return nil, fmt.Errorf("构建上下文失败: %w", err)
	}

	// 添加用户消息到会话
	userMessage := &aiModels.ChatMessage{
		ID:        primitive.NewObjectID(),
		SessionID: session.SessionID,
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
	}
	session.Messages = append(session.Messages, *userMessage)

	// 保存用户消息
	if err := s.repository.CreateMessage(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 调用AI生成响应
	startTime := time.Now()

	// 准备AI请求 - 使用GenerateContentRequest
	aiRequest := &GenerateContentRequest{
		ProjectID: req.ProjectID,
		Prompt:    req.Message,
		Options:   req.Options,
	}

	aiResponse, err := s.aiService.GenerateContent(ctx, aiRequest)
	if err != nil {
		return nil, fmt.Errorf("生成AI响应失败: %w", err)
	}

	responseTime := time.Since(startTime)

	// 添加AI响应到会话
	assistantMessage := &aiModels.ChatMessage{
		ID:        primitive.NewObjectID(),
		SessionID: session.SessionID, // 使用SessionID字段
		Role:      "assistant",
		Content:   aiResponse.Content,
		TokenUsed: aiResponse.TokensUsed,
		Timestamp: time.Now(),
	}
	session.Messages = append(session.Messages, *assistantMessage)

	// 保存AI消息
	if err := s.repository.CreateMessage(ctx, assistantMessage); err != nil {
		return nil, fmt.Errorf("保存AI消息失败: %w", err)
	}

	// 更新会话
	session.UpdatedAt = time.Now()
	err = s.repository.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("保存会话失败: %w", err)
	}

	// 构建响应
	response := &ChatResponse{
		SessionID:    session.SessionID,
		Message:      dto.ToMessageDTO(assistantMessage),
		TokensUsed:   aiResponse.TokensUsed,
		Model:        aiResponse.Model,
		ContextUsed:  len(context) > 1, // 如果有系统消息则表示使用了上下文
		ResponseTime: responseTime,
	}

	return response, nil
}

// StartChatStream 开始流式聊天
func (s *ChatService) StartChatStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChatResponse, error) {
	// 验证请求
	if req.Message == "" {
		return nil, errors.New("消息内容不能为空")
	}

	// 获取或创建会话
	session, err := s.getOrCreateSession(ctx, req.SessionID, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("获取会话失败: %w", err)
	}

	// 构建聊天上下文
	context, err := s.buildChatContext(ctx, session, req)
	if err != nil {
		return nil, fmt.Errorf("构建上下文失败: %w", err)
	}

	// 添加用户消息到会话
	userMessage := &aiModels.ChatMessage{
		ID:        primitive.NewObjectID(),
		SessionID: session.SessionID,
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
	}
	session.Messages = append(session.Messages, *userMessage)

	// 保存用户消息
	if err := s.repository.CreateMessage(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("保存用户消息失败: %w", err)
	}

	// 创建响应通道
	responseChan := make(chan *StreamChatResponse, 10)

	// 启动流式生成
	go func() {
		defer close(responseChan)

		startTime := time.Now()

		// 准备AI请求
		aiRequest := &GenerateContentRequest{
			ProjectID: req.ProjectID,
			Prompt:    req.Message,
			Options:   req.Options,
		}

		// 调用流式AI生成
		streamChan, err := s.aiService.GenerateContentStream(ctx, aiRequest)
		if err != nil {
			// 发送错误响应
			select {
			case responseChan <- &StreamChatResponse{
				SessionID:    session.SessionID,
				MessageID:    primitive.NewObjectID().Hex(),
				Content:      "",
				Delta:        "",
				TokensUsed:   0,
				Model:        "",
				ContextUsed:  len(context) > 1,
				IsComplete:   true,
				ResponseTime: time.Since(startTime),
			}:
			case <-ctx.Done():
			}
			return
		}

		var fullContent string
		var totalTokens int
		var model string
		messageID := primitive.NewObjectID().Hex()

		// 处理流式响应
		for aiResponse := range streamChan {
			if aiResponse == nil {
				continue
			}

			fullContent += aiResponse.Content
			totalTokens = aiResponse.TokensUsed
			model = aiResponse.Model

			// 发送流式响应
			select {
			case responseChan <- &StreamChatResponse{
				SessionID:    session.SessionID,
				MessageID:    messageID,
				Content:      fullContent,
				Delta:        aiResponse.Content,
				TokensUsed:   totalTokens,
				Model:        model,
				ContextUsed:  len(context) > 1,
				IsComplete:   false,
				ResponseTime: time.Since(startTime),
			}:
			case <-ctx.Done():
				return
			}
		}

		// 发送完成响应
		select {
		case responseChan <- &StreamChatResponse{
			SessionID:    session.SessionID,
			MessageID:    messageID,
			Content:      fullContent,
			Delta:        "",
			TokensUsed:   totalTokens,
			Model:        model,
			ContextUsed:  len(context) > 1,
			IsComplete:   true,
			ResponseTime: time.Since(startTime),
		}:
		case <-ctx.Done():
			return
		}

		// 保存AI响应消息
		assistantMessage := &aiModels.ChatMessage{
			ID:        primitive.NewObjectID(),
			SessionID: session.SessionID,
			Role:      "assistant",
			Content:   fullContent,
			TokenUsed: totalTokens,
			Timestamp: time.Now(),
		}
		session.Messages = append(session.Messages, *assistantMessage)

		// 保存AI消息
		if err := s.repository.CreateMessage(ctx, assistantMessage); err != nil {
			// 记录错误但不中断流式响应
			fmt.Printf("保存AI消息失败: %v\n", err)
		}

		// 更新会话
		session.UpdatedAt = time.Now()
		if err := s.repository.UpdateSession(ctx, session); err != nil {
			// 记录错误但不中断流式响应
			fmt.Printf("保存会话失败: %v\n", err)
		}
	}()

	return responseChan, nil
}

// buildChatContext 构建对话上下文
func (s *ChatService) buildChatContext(ctx context.Context, session *aiModels.ChatSession, req *ChatRequest) ([]*aiModels.ChatMessage, error) {
	messages := make([]*aiModels.ChatMessage, 0)

	// 添加系统提示
	systemPrompt := s.buildSystemPrompt(req.ContextType, req.ProjectID)
	if systemPrompt != "" {
		messages = append(messages, &aiModels.ChatMessage{
			ID:        primitive.NewObjectID(),
			Role:      "system",
			Content:   systemPrompt,
			Timestamp: time.Now(),
		})
	}

	// 注意: NovelContextService 功能已禁用（未完成）
	// 如需使用小说上下文，请等待该功能完全实现后再启用

	// 添加历史对话（保留最近的对话）
	historyLimit := 10 // 保留最近10轮对话
	startIndex := 0
	if len(session.Messages) > historyLimit*2 {
		startIndex = len(session.Messages) - historyLimit*2
	}

	for i := startIndex; i < len(session.Messages); i++ {
		messages = append(messages, &session.Messages[i])
	}

	return messages, nil
}

// buildSystemPrompt 构建系统提示
func (s *ChatService) buildSystemPrompt(contextType, projectID string) string {
	switch contextType {
	case "novel":
		return `你是一个专业的小说创作助手。你可以帮助用户：
1. 分析小说情节和角色
2. 提供创作建议和灵感
3. 协助完善故事结构
4. 解答创作相关问题

请根据用户提供的上下文信息，给出专业、有建设性的建议。`
	default:
		return `你是一个智能助手，可以回答各种问题并提供帮助。请保持友好、专业的态度。`
	}
}

// convertContextToMessage 将上下文转换为消息
func (s *ChatService) convertContextToMessage(context *ai.AIContext) string {
	if context == nil {
		return ""
	}

	message := "当前小说上下文信息：\n"

	// 添加当前章节信息
	if context.CurrentChapter != nil {
		message += fmt.Sprintf("当前章节：%s\n", context.CurrentChapter.Title)
		if context.CurrentChapter.Summary != "" {
			message += fmt.Sprintf("章节摘要：%s\n", context.CurrentChapter.Summary)
		}
	}

	// 添加活跃角色信息
	if len(context.ActiveCharacters) > 0 {
		message += "主要角色：\n"
		for _, char := range context.ActiveCharacters {
			message += fmt.Sprintf("- %s: %s\n", char.Name, char.Summary)
		}
	}

	// 添加情节线信息
	if len(context.PlotThreads) > 0 {
		message += "当前情节线：\n"
		for _, plot := range context.PlotThreads {
			message += fmt.Sprintf("- %s: %s\n", plot.Name, plot.Description)
		}
	}

	return message
}

// getOrCreateSession 获取或创建会话
func (s *ChatService) getOrCreateSession(ctx context.Context, sessionID, projectID string) (*aiModels.ChatSession, error) {
	if sessionID != "" {
		// 尝试获取现有会话
		session, err := s.repository.GetSessionByID(ctx, sessionID)
		if err == nil {
			return session, nil
		}
	}

	// 创建新会话
	session := &aiModels.ChatSession{
		ID:          primitive.NewObjectID(),
		SessionID:   generateSessionID(),
		ProjectID:   projectID,
		Title:       "新对话",
		Description: "",
		Status:      "active",
		Settings:    &aiModels.ChatSettings{},
		Metadata:    &aiModels.ChatMetadata{},
		Messages:    []aiModels.ChatMessage{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存会话到数据库
	if err := s.repository.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("创建会话失败: %w", err)
	}

	return session, nil
}

// generateSessionID 生成会话ID
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

// GetChatHistory 获取聊天历史
func (s *ChatService) GetChatHistory(ctx context.Context, sessionID string) (*dto.ChatSessionDTO, error) {
	session, err := s.repository.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return dto.ToSessionDTO(session), nil
}

// ListChatSessions 获取会话列表
func (s *ChatService) ListChatSessions(ctx context.Context, projectID string, limit, offset int) ([]*dto.ChatSessionDTO, error) {
	sessions, err := s.repository.GetSessionsByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.ChatSessionDTO, len(sessions))
	for i, session := range sessions {
		result[i] = dto.ToSessionDTO(session)
	}

	return result, nil
}

// DeleteChatSession 删除聊天会话
func (s *ChatService) DeleteChatSession(ctx context.Context, sessionID string) error {
	return s.repository.DeleteSession(ctx, sessionID)
}

// GetSessions 获取会话列表
func (s *ChatService) GetSessions(projectID string, limit, offset int) ([]*aiModels.ChatSession, error) {
	ctx := context.Background()
	return s.repository.GetSessionsByProjectID(ctx, projectID, limit, offset)
}

// GetHistory 获取聊天历史
func (s *ChatService) GetHistory(sessionID string) (*aiModels.ChatSession, error) {
	ctx := context.Background()
	return s.repository.GetSessionByID(ctx, sessionID)
}

// UpdateSession 更新会话
func (s *ChatService) UpdateSession(sessionID, title, description string, metadata map[string]interface{}) error {
	ctx := context.Background()

	// 获取现有会话
	session, err := s.repository.GetSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("获取会话失败: %w", err)
	}

	// 更新字段
	if title != "" {
		session.Title = title
	}
	if description != "" {
		session.Description = description
	}
	if metadata != nil {
		// 将map[string]interface{}转换为*ChatMetadata
		chatMetadata := &aiModels.ChatMetadata{
			CustomFields: metadata,
		}
		session.Metadata = chatMetadata
	}

	// 保存更新
	return s.repository.UpdateSession(ctx, session)
}

// DeleteSession 删除会话
func (s *ChatService) DeleteSession(sessionID string) error {
	ctx := context.Background()
	return s.repository.DeleteSession(ctx, sessionID)
}

// GetStatistics 获取统计信息
func (s *ChatService) GetStatistics(projectID string) (*ChatStatistics, error) {
	ctx := context.Background()
	return s.repository.GetSessionStatistics(ctx, projectID)
}
