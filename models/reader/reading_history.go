package reader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/shared/types"
)

// ReadingHistory 阅读历史记录
type ReadingHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId" validate:"required"`
	BookID    primitive.ObjectID `bson:"book_id" json:"bookId" validate:"required"`
	ChapterID primitive.ObjectID `bson:"chapter_id" json:"chapterId" validate:"required"`

	// 阅读信息
	ReadDuration int            `bson:"read_duration" json:"readDuration"` // 阅读时长（秒）
	Progress     types.Progress `bson:"progress" json:"progress"`           // 阅读进度（0-1）

	// 设备信息
	DeviceType string `bson:"device_type,omitempty" json:"deviceType,omitempty"` // web, ios, android
	DeviceID   string `bson:"device_id,omitempty" json:"deviceId,omitempty"`     // 设备唯一标识

	// 时间戳
	StartTime time.Time `bson:"start_time" json:"startTime"` // 阅读开始时间
	EndTime   time.Time `bson:"end_time" json:"endTime"`     // 阅读结束时间
	CreatedAt time.Time `bson:"created_at" json:"createdAt"` // 记录创建时间
}

// TableName 返回集合名称
func (ReadingHistory) TableName() string {
	return "reading_histories"
}
