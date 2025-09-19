package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Name string

const (
	OpenAI Name = "openai"
	Google Name = "google"
	Baidu  Name = "baidu"
)

// AIProvider AI服务提供商配置
type AIProvider struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`        // 服务提供商名称：openai, baidu等
	APIKey    string    `bson:"api_key" json:"-"`        // API密钥（不返回给客户端）
	BaseURL   string    `bson:"base_url" json:"baseUrl"` // API基础URL
	Enabled   bool      `bson:"enabled" json:"enabled"`  // 是否启用
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// BeforeCreate 在创建前设置时间戳
func (p *AIProvider) BeforeCreate() {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.ID == "" {
		// 生成唯一ID
		p.ID = primitive.NewObjectID().Hex()
	}
}

// BeforeUpdate 在更新前刷新更新时间戳
func (p *AIProvider) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}
