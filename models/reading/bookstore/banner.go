package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Banner 轮播图模型
type Banner struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" validate:"required,min=1,max=100"`                      // 标题
	Description string             `bson:"description" json:"description" validate:"max=200"`                         // 描述
	Image       string             `bson:"image" json:"image" validate:"required,url"`                                // 图片URL
	Target      string             `bson:"target" json:"target" validate:"required"`                                  // 跳转目标
	TargetType  string             `bson:"target_type" json:"targetType" validate:"required,oneof=book category url"` // 目标类型：book书籍 category分类 url外链
	SortOrder   int                `bson:"sort_order" json:"sortOrder"`                                               // 排序权重
	IsActive    bool               `bson:"is_active" json:"isActive"`                                                 // 是否启用
	StartTime   *time.Time         `bson:"start_time,omitempty" json:"startTime,omitempty"`                           // 开始时间
	EndTime     *time.Time         `bson:"end_time,omitempty" json:"endTime,omitempty"`                               // 结束时间
	ClickCount  int64              `bson:"click_count" json:"clickCount"`                                             // 点击次数
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`                                               // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`                                               // 更新时间
}

// BannerFilter Banner查询过滤器
type BannerFilter struct {
	IsActive   *bool   `json:"isActive,omitempty"`
	TargetType *string `json:"targetType,omitempty"`
	SortBy     string  `json:"sortBy,omitempty"`    // sort_order, click_count, created_at
	SortOrder  string  `json:"sortOrder,omitempty"` // asc, desc
	Limit      int     `json:"limit,omitempty"`
	Offset     int     `json:"offset,omitempty"`
}
