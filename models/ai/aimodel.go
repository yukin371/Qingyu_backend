package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelType string

const (
	ModelTypeChat  ModelType = "chat"
	ModelTypeImage ModelType = "image"
)

type AIModel struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Provider    string    `bson:"provider" json:"provider"`       // 服务提供商名称
	Name        string    `bson:"name" json:"name"`               // 模型名称
	Type        ModelType `bson:"type" json:"type"`               // 模型类型
	MaxTokens   int       `bson:"max_tokens" json:"maxTokens"`    // 最大令牌数
	Enabled     bool      `bson:"enabled" json:"enabled"`         // 是否启用
	Description string    `bson:"description" json:"description"` // 模型描述
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

// BeforeCreate 在创建前设置时间戳
func (m *AIModel) BeforeCreate() {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	if m.ID == "" {
		// 生成唯一ID
		m.ID = primitive.NewObjectID().Hex()
	}
}

// BeforeUpdate 在更新前刷新更新时间戳
func (m *AIModel) BeforeUpdate() {
	m.UpdatedAt = time.Now()
}
