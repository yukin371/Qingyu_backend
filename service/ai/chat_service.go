package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Qingyu_backend/models/ai"
	aiModels "Qingyu_backend/models/ai"
	"Qingyu_backend/service/ai/adapter"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatRepositoryInterface 聊天仓库接口
type ChatRepositoryInterface interface {
	CreateSession(ctx context.Context, session *aiModels.ChatSession) error
	GetSessionByID(ctx context.Context, sessionID string) (*aiModels.ChatSession, error)
	UpdateSession(ctx context.Context, session *aiModels.ChatSession) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetSessionsByProjectID(ctx context.Context, projectID string, limit, offset int) ([]*aiModels.ChatSession, error)
	CreateMessage(ctx context.Context, message *aiModels.ChatMessage) error
	GetSessionStatistics(ctx context.Context, projectID string) (*ChatStatistics, error)
}

// AIServiceInterface AI服务接口
type AIServiceInterface interface {
	GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)
}

// ChatService AI聊天服务
type ChatService struct {
	aiService           AIServiceInterface
	adapterManager      *adapter.AdapterManager
	novelContextService *NovelContextService
	repository          ChatRepositoryInterface
}

// NewChatService 创建聊天服务
func NewChatService(aiService *Service, repository ChatRepositoryInterface) *ChatService {
	return &ChatService{
		aiService:           aiService,
		adapterManager:      aiService.adapterManager,
		novelContextService: nil, // 暂时设为nil，避免循环依赖
		repository:          repository,
	}
}

// ChatMessage 聊天消息
type ChatMessage struct {
	ID        string                 `json:"id"`
	Role      string                 `json:"role"` // user, assistant, system
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ChatSession 聊天会话
type ChatSession struct {
	ID        string         `json:"id"`
	ProjectID string         `json:"projectId,omitempty"`
	Title     string         `json:"title"`
	Messages  []*ChatMessage `json:"messages"`
	Context   *ai.AIContext  `json:"context,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID   string                 `json:"sessionId,omitempty"`
	ProjectID   string                 `json:"projectId,omitempty"`
	Message     string                 `json:"message" binding:"required"`
	UseContext  bool                   `json:"useContext"`
	ContextType string                 `json:"contextType,omitempty"` // novel, general
	Options     *ai.GenerateOptions    `json:"options,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	SessionID    string        `json:"sessionId"`
	Message      *ChatMessage  `json:"message"`
	TokensUsed   int           `json:"tokensUsed"`
	Model        string        `json:"model"`
	ContextUsed  bool          `json:"contextUsed"`
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

	// 构建聊天上下文 - 转换为服务层类型
	serviceSession := convertToServiceChatSession(session)
	context, err := s.buildChatContext(ctx, serviceSession, req)
	if err != nil {
		return nil, fmt.Errorf("构建上下文失败: %w", err)
	}

	// 添加用户消息到会话
	userMessage := &aiModels.ChatMessage{
		ID:        primitive.NewObjectID(),
		SessionID: session.SessionID, // 使用SessionID字段
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
		SessionID:    session.SessionID, // 使用SessionID字段
		Message:      convertToChatMessage(assistantMessage),
		TokensUsed:   aiResponse.TokensUsed,
		Model:        aiResponse.Model,
		ContextUsed:  len(context) > 1, // 如果有系统消息则表示使用了上下文
		ResponseTime: responseTime,
	}

	return response, nil
}

// buildChatContext 构建对话上下文
func (s *ChatService) buildChatContext(ctx context.Context, session *ChatSession, req *ChatRequest) ([]*ChatMessage, error) {
	messages := make([]*ChatMessage, 0)

	// 添加系统提示
	systemPrompt := s.buildSystemPrompt(req.ContextType, req.ProjectID)
	if systemPrompt != "" {
		messages = append(messages, &ChatMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	// 如果需要使用小说上下文
	if req.UseContext && req.ProjectID != "" && s.novelContextService != nil {
		contextReq := &ai.ContextBuildRequest{
			ProjectID:       req.ProjectID,
			CurrentPosition: req.Message, // 使用用户消息作为当前位置
			MaxTokens:       1000,        // 为上下文预留token
		}

		contextResp, err := s.novelContextService.BuildContext(ctx, contextReq)
		if err == nil && contextResp.Context != nil {
			// 将上下文信息转换为系统消息
			contextMessage := s.convertContextToMessage(contextResp.Context)
			if contextMessage != "" {
				messages = append(messages, &ChatMessage{
					Role:    "system",
					Content: contextMessage,
				})
			}
		}
	}

	// 添加历史对话（保留最近的对话）
	historyLimit := 10 // 保留最近10轮对话
	startIndex := 0
	if len(session.Messages) > historyLimit*2 {
		startIndex = len(session.Messages) - historyLimit*2
	}

	for i := startIndex; i < len(session.Messages); i++ {
		messages = append(messages, session.Messages[i])
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

// getSession 获取会话（占位实现）
func (s *ChatService) getSession(ctx context.Context, sessionID string) (*ChatSession, error) {
	// 这里需要从数据库获取会话
	// 暂时返回错误，表示会话不存在
	return nil, fmt.Errorf("会话不存在")
}

// saveSession 保存会话（占位实现）
func (s *ChatService) saveSession(ctx context.Context, session *ChatSession) error {
	// 这里需要保存会话到数据库
	// 暂时返回nil，表示保存成功
	return nil
}

// convertToChatMessage 转换为服务层ChatMessage
func convertToChatMessage(msg *aiModels.ChatMessage) *ChatMessage {
	return &ChatMessage{
		ID:        msg.ID.Hex(),
		Role:      msg.Role,
		Content:   msg.Content,
		Timestamp: msg.Timestamp,
		Metadata:  make(map[string]interface{}),
	}
}

// convertToServiceChatSession 转换为服务层ChatSession
func convertToServiceChatSession(session *aiModels.ChatSession) *ChatSession {
	messages := make([]*ChatMessage, len(session.Messages))
	for i, msg := range session.Messages {
		messages[i] = &ChatMessage{
			ID:        msg.ID.Hex(),
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
			Metadata:  make(map[string]interface{}),
		}
	}
	
	return &ChatSession{
		ID:        session.SessionID,
		ProjectID: session.ProjectID,
		Title:     session.Title,
		Messages:  messages,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}

// convertToAdapterMessages 转换消息格式
func convertToAdapterMessages(messages []*ChatMessage) []adapter.Message {
	adapterMessages := make([]adapter.Message, len(messages))
	for i, msg := range messages {
		adapterMessages[i] = adapter.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return adapterMessages
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
func (s *ChatService) GetChatHistory(ctx context.Context, sessionID string) (*ChatSession, error) {
	session, err := s.repository.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	
	return convertToServiceChatSession(session), nil
}

// ListChatSessions 获取会话列表
func (s *ChatService) ListChatSessions(ctx context.Context, projectID string, limit, offset int) ([]*ChatSession, error) {
	sessions, err := s.repository.GetSessionsByProjectID(ctx, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	
	result := make([]*ChatSession, len(sessions))
	for i, session := range sessions {
		result[i] = convertToServiceChatSession(session)
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