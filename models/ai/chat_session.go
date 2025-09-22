package ai

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatSession 聊天会话模型
type ChatSession struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SessionID   string             `json:"sessionId" bson:"session_id"`
	ProjectID   string             `json:"projectId" bson:"project_id"`
	UserID      string             `json:"userId" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status"`
	Settings    *ChatSettings      `json:"settings" bson:"settings"`
	Metadata    *ChatMetadata      `json:"metadata" bson:"metadata"`
	Messages    []ChatMessage      `json:"messages" bson:"messages"`
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SessionID string             `json:"sessionId" bson:"session_id"`
	Role      string             `json:"role" bson:"role"`
	Content   string             `json:"content" bson:"content"`
	TokenUsed int                `json:"tokenUsed" bson:"token_used"`
	Metadata  *MessageMeta       `json:"metadata" bson:"metadata"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deletedAt,omitempty" bson:"deleted_at,omitempty"`
}

// ChatSettings 聊天设置
type ChatSettings struct {
	Model           string  `json:"model"`           // AI模型
	Temperature     float32 `json:"temperature"`     // 温度参数
	MaxTokens       int     `json:"maxTokens"`       // 最大token数
	TopP            float32 `json:"topP"`            // Top-p参数
	FrequencyPenalty float32 `json:"frequencyPenalty"` // 频率惩罚
	PresencePenalty  float32 `json:"presencePenalty"`  // 存在惩罚
	SystemPrompt     string  `json:"systemPrompt"`     // 系统提示词
	ContextLength    int     `json:"contextLength"`    // 上下文长度
	EnableMemory     bool    `json:"enableMemory"`     // 是否启用记忆
	EnableContext    bool    `json:"enableContext"`    // 是否启用上下文
}

// ChatMetadata 聊天元数据
type ChatMetadata struct {
	TotalMessages    int                    `json:"totalMessages"`    // 总消息数
	TotalTokens      int                    `json:"totalTokens"`      // 总token数
	AverageResponse  float64                `json:"averageResponse"`  // 平均响应时间
	LastActiveAt     time.Time              `json:"lastActiveAt"`     // 最后活跃时间
	Tags             []string               `json:"tags"`             // 标签
	Category         string                 `json:"category"`         // 分类
	Priority         int                    `json:"priority"`         // 优先级
	CustomFields     map[string]interface{} `json:"customFields"`     // 自定义字段
}

// MessageMeta 消息元数据
type MessageMeta struct {
	ResponseTime    float64                `json:"responseTime"`    // 响应时间(毫秒)
	ModelUsed       string                 `json:"modelUsed"`       // 使用的模型
	ContextUsed     bool                   `json:"contextUsed"`     // 是否使用了上下文
	MemoryUsed      bool                   `json:"memoryUsed"`      // 是否使用了记忆
	Sources         []string               `json:"sources"`         // 信息来源
	Confidence      float64                `json:"confidence"`      // 置信度
	Sentiment       string                 `json:"sentiment"`       // 情感分析
	Intent          string                 `json:"intent"`          // 意图识别
	Entities        []Entity               `json:"entities"`        // 实体识别
	CustomData      map[string]interface{} `json:"customData"`      // 自定义数据
}

// Entity 实体信息
type Entity struct {
	Type       string  `json:"type"`       // 实体类型
	Value      string  `json:"value"`      // 实体值
	Confidence float64 `json:"confidence"` // 置信度
	StartPos   int     `json:"startPos"`   // 开始位置
	EndPos     int     `json:"endPos"`     // 结束位置
}

// BeforeCreate MongoDB钩子 - 创建前
func (cs *ChatSession) BeforeCreate() {
	if cs.ID.IsZero() {
		cs.ID = primitive.NewObjectID()
	}
	if cs.SessionID == "" {
		cs.SessionID = generateSessionID()
	}
	if cs.Status == "" {
		cs.Status = "active"
	}
	cs.CreatedAt = time.Now()
	cs.UpdatedAt = time.Now()
}

// BeforeUpdate MongoDB钩子 - 更新前
func (cs *ChatSession) BeforeUpdate() {
	cs.UpdatedAt = time.Now()
}

// TableName 指定集合名
func (ChatSession) CollectionName() string {
	return "ai_chat_sessions"
}

// TableName 指定集合名
func (ChatMessage) CollectionName() string {
	return "ai_chat_messages"
}

// generateSessionID 生成会话ID
func generateSessionID() string {
	// 这里应该使用更安全的ID生成方法，比如UUID
	return "chat_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// GetActiveMessages 获取活跃消息（未删除）
func (cs *ChatSession) GetActiveMessages() []ChatMessage {
	var activeMessages []ChatMessage
	for _, msg := range cs.Messages {
		if msg.DeletedAt == nil {
			activeMessages = append(activeMessages, msg)
		}
	}
	return activeMessages
}

// GetTotalTokens 获取总token数
func (cs *ChatSession) GetTotalTokens() int {
	total := 0
	for _, msg := range cs.Messages {
		total += msg.TokenUsed
	}
	return total
}

// IsActive 检查会话是否活跃
func (cs *ChatSession) IsActive() bool {
	return cs.Status == "active"
}

// Archive 归档会话
func (cs *ChatSession) Archive() {
	cs.Status = "archived"
	cs.UpdatedAt = time.Now()
}

// Activate 激活会话
func (cs *ChatSession) Activate() {
	cs.Status = "active"
	cs.UpdatedAt = time.Now()
}

// AddMessage 添加消息
func (cs *ChatSession) AddMessage(role, content string, tokenUsed int, metadata *MessageMeta) *ChatMessage {
	message := &ChatMessage{
		SessionID: cs.SessionID,
		Role:      role,
		Content:   content,
		TokenUsed: tokenUsed,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
	
	cs.Messages = append(cs.Messages, *message)
	cs.UpdatedAt = time.Now()
	
	// 更新元数据
	if cs.Metadata == nil {
		cs.Metadata = &ChatMetadata{}
	}
	cs.Metadata.TotalMessages = len(cs.Messages)
	cs.Metadata.TotalTokens = cs.GetTotalTokens()
	cs.Metadata.LastActiveAt = time.Now()
	
	return message
}

// GetLastMessage 获取最后一条消息
func (cs *ChatSession) GetLastMessage() *ChatMessage {
	if len(cs.Messages) == 0 {
		return nil
	}
	return &cs.Messages[len(cs.Messages)-1]
}

// GetMessagesByRole 根据角色获取消息
func (cs *ChatSession) GetMessagesByRole(role string) []ChatMessage {
	var messages []ChatMessage
	for _, msg := range cs.Messages {
		if msg.Role == role {
			messages = append(messages, msg)
		}
	}
	return messages
}

// UpdateSettings 更新设置
func (cs *ChatSession) UpdateSettings(settings *ChatSettings) {
	cs.Settings = settings
	cs.UpdatedAt = time.Now()
}

// UpdateMetadata 更新元数据
func (cs *ChatSession) UpdateMetadata(metadata *ChatMetadata) {
	cs.Metadata = metadata
	cs.UpdatedAt = time.Now()
}