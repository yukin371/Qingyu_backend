package reader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReadingHistory 阅读历史记录
type ReadingHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id" validate:"required"`
	BookID    string             `bson:"book_id" json:"book_id" validate:"required"`
	ChapterID string             `bson:"chapter_id" json:"chapter_id" validate:"required"`

	// 阅读信息
	ReadDuration int     `bson:"read_duration" json:"read_duration"` // 阅读时长（秒）
	Progress     float64 `bson:"progress" json:"progress"`           // 阅读进度（0-100）

	// 设备信息
	DeviceType string `bson:"device_type,omitempty" json:"device_type,omitempty"` // web, ios, android
	DeviceID   string `bson:"device_id,omitempty" json:"device_id,omitempty"`     // 设备唯一标识

	// 时间戳
	StartTime time.Time `bson:"start_time" json:"start_time"` // 阅读开始时间
	EndTime   time.Time `bson:"end_time" json:"end_time"`     // 阅读结束时间
	CreatedAt time.Time `bson:"created_at" json:"created_at"` // 记录创建时间
}

// TableName 返回集合名称
func (ReadingHistory) TableName() string {
	return "reading_histories"
}
