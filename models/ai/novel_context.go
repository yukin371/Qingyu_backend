package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NovelContext 小说上下文模型
type NovelContext struct {
	ID         string                 `bson:"_id,omitempty" json:"id"`
	ProjectID  string                 `bson:"project_id" json:"projectId"`
	Type       string                 `bson:"type" json:"type"` // character, plot, setting, chapter
	Title      string                 `bson:"title" json:"title"`
	Content    string                 `bson:"content" json:"content"`
	Summary    string                 `bson:"summary" json:"summary"`
	Embedding  []float32              `bson:"embedding" json:"-"`
	Metadata   map[string]interface{} `bson:"metadata" json:"metadata"`
	Importance int                    `bson:"importance" json:"importance"` // 1-10
	CreatedAt  time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updatedAt"`
}

// MemoryRelation 记忆关联模型
type MemoryRelation struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	SourceID     string    `bson:"source_id" json:"sourceId"`
	TargetID     string    `bson:"target_id" json:"targetId"`
	RelationType string    `bson:"relation_type" json:"relationType"`
	Strength     float32   `bson:"strength" json:"strength"`
	CreatedAt    time.Time `bson:"created_at" json:"createdAt"`
}

// ContextMemory 上下文记忆层次
type ContextMemory struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	ProjectID   string    `bson:"project_id" json:"projectId"`
	MemoryType  string    `bson:"memory_type" json:"memoryType"` // short_term, medium_term, long_term
	Content     string    `bson:"content" json:"content"`
	Summary     string    `bson:"summary" json:"summary"`
	Embedding   []float32 `bson:"embedding" json:"-"`
	Importance  int       `bson:"importance" json:"importance"`
	AccessCount int       `bson:"access_count" json:"accessCount"`
	LastAccess  time.Time `bson:"last_access" json:"lastAccess"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

// RetrievalResult 检索结果
type RetrievalResult struct {
	Context    *NovelContext `json:"context"`
	Score      float32       `json:"score"`
	Relevance  string        `json:"relevance"`
	SourceType string        `json:"sourceType"` // vector, keyword, metadata
}

// ContextBuildRequest 上下文构建请求
type ContextBuildRequest struct {
	ProjectID       string `json:"projectId"`
	CurrentPosition string `json:"currentPosition"`
	MaxTokens       int    `json:"maxTokens"`
	FocusType       string `json:"focusType,omitempty"` // character, plot, setting
	FocusID         string `json:"focusId,omitempty"`
}

// ContextBuildResponse 上下文构建响应
type ContextBuildResponse struct {
	Context       *AIContext         `json:"context"`
	TokenCount    int                `json:"tokenCount"`
	Sources       []*RetrievalResult `json:"sources"`
	BuildStrategy string             `json:"buildStrategy"`
	BuildTime     time.Duration      `json:"buildTime"`
}

// BeforeCreate 在创建前设置时间戳
func (nc *NovelContext) BeforeCreate() {
	now := time.Now()
	nc.CreatedAt = now
	nc.UpdatedAt = now
	if nc.ID == "" {
		nc.ID = primitive.NewObjectID().Hex()
	}
}

// BeforeUpdate 在更新前刷新更新时间戳
func (nc *NovelContext) BeforeUpdate() {
	nc.UpdatedAt = time.Now()
}

// BeforeCreate 在创建前设置时间戳
func (mr *MemoryRelation) BeforeCreate() {
	mr.CreatedAt = time.Now()
	if mr.ID == "" {
		mr.ID = primitive.NewObjectID().Hex()
	}
}

// BeforeCreate 在创建前设置时间戳
func (cm *ContextMemory) BeforeCreate() {
	now := time.Now()
	cm.CreatedAt = now
	cm.UpdatedAt = now
	cm.LastAccess = now
	if cm.ID == "" {
		cm.ID = primitive.NewObjectID().Hex()
	}
}

// BeforeUpdate 在更新前刷新更新时间戳
func (cm *ContextMemory) BeforeUpdate() {
	cm.UpdatedAt = time.Now()
}

// UpdateAccess 更新访问信息
func (cm *ContextMemory) UpdateAccess() {
	cm.AccessCount++
	cm.LastAccess = time.Now()
}