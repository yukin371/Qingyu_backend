package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APIRequestLog struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	UserID     string    `bson:"user_id" json:"userId"`         // 用户ID
	Provider   string    `bson:"provider" json:"provider"`      // 服务提供商名称
	Model      string    `bson:"model" json:"model"`            // 模型名称
	Endpoint   string    `bson:"endpoint" json:"endpoint"`      // 请求的API端点
	Request    string    `bson:"request" json:"-"`              // 请求数据（不返回给客户端）
	Response   string    `bson:"response" json:"-"`             // 响应数据（不返回给客户端）
	StatusCode int       `bson:"status_code" json:"statusCode"` // HTTP状态码
	Duration   int64     `bson:"duration" json:"duration"`      // 请求耗时(毫秒)
	IP         string    `bson:"ip" json:"ip"`                  // 请求来源IP
	Error      string    `bson:"error,omitempty" json:"error"`  // 错误信息
	CreatedAt  time.Time `bson:"created_at" json:"createdAt"`
}

// BeforeCreate 在创建前设置时间戳
func (l *APIRequestLog) BeforeCreate() {
	now := time.Now()
	l.CreatedAt = now
	if l.ID == "" {
		// 生成唯一ID
		l.ID = primitive.NewObjectID().Hex()
	}
}
