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

// GetConditions 实现Filter接口
func (f *BannerFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	if f.IsActive != nil {
		conditions["is_active"] = *f.IsActive
	}
	if f.TargetType != nil {
		conditions["target_type"] = *f.TargetType
	}

	return conditions
}

// GetSort 实现Filter接口
func (f *BannerFilter) GetSort() map[string]int {
	sort := make(map[string]int)

	if f.SortBy != "" {
		sortValue := 1 // 默认升序
		if f.SortOrder == "desc" {
			sortValue = -1
		}

		switch f.SortBy {
		case "sort_order":
			sort["sort_order"] = sortValue
		case "click_count":
			sort["click_count"] = sortValue
		case "created_at":
			sort["created_at"] = sortValue
		default:
			sort["sort_order"] = 1
		}
	} else {
		// 默认按排序权重升序
		sort["sort_order"] = 1
	}

	return sort
}

// GetLimit 实现Filter接口
func (f *BannerFilter) GetLimit() int {
	return f.Limit
}

// GetOffset 实现Filter接口
func (f *BannerFilter) GetOffset() int {
	return f.Offset
}

// GetFields 实现Filter接口
func (f *BannerFilter) GetFields() []string {
	return []string{} // Banner查询返回所有字段
}

// Validate 实现Filter接口
func (f *BannerFilter) Validate() error {
	if f.Limit < 0 {
		return &ValidationError{Message: "limit不能为负数"}
	}
	if f.Offset < 0 {
		return &ValidationError{Message: "offset不能为负数"}
	}
	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
