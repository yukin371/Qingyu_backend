package shared

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Announcement 公告模型
type Announcement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" validate:"required,min=1,max=200"`                              // 标题
	Content     string             `bson:"content" json:"content" validate:"required,min=1"`                                  // 内容（支持HTML）
	Type        string             `bson:"type" json:"type" validate:"required,oneof=info warning notice"`                    // 类型：info信息 warning警告 notice通知
	Priority    int                `bson:"priority" json:"priority"`                                                          // 优先级（数字越大优先级越高）
	IsActive    bool               `bson:"is_active" json:"isActive"`                                                         // 是否启用
	StartTime   *time.Time         `bson:"start_time,omitempty" json:"startTime,omitempty"`                                   // 开始时间
	EndTime     *time.Time         `bson:"end_time,omitempty" json:"endTime,omitempty"`                                       // 结束时间
	TargetUsers string             `bson:"target_users" json:"targetUsers" validate:"required,oneof=all reader writer admin"` // 目标用户：all所有 reader读者 writer作者 admin管理员
	ViewCount   int64              `bson:"view_count" json:"viewCount"`                                                       // 查看次数
	CreatedBy   primitive.ObjectID `bson:"created_by,omitempty" json:"createdBy,omitempty"`                                   // 创建者ID
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`                                                       // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`                                                       // 更新时间
}

// AnnouncementFilter 公告查询过滤器
type AnnouncementFilter struct {
	IsActive    *bool   `json:"isActive,omitempty"`
	Type        *string `json:"type,omitempty"`
	TargetUsers *string `json:"targetUsers,omitempty"`
	SortBy      string  `json:"sortBy,omitempty"`    // priority, created_at, view_count
	SortOrder   string  `json:"sortOrder,omitempty"` // asc, desc
	Limit       int     `json:"limit,omitempty"`
	Offset      int     `json:"offset,omitempty"`
}

// GetConditions 实现Filter接口
func (f *AnnouncementFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	if f.IsActive != nil {
		conditions["is_active"] = *f.IsActive
	}
	if f.Type != nil {
		conditions["type"] = *f.Type
	}
	if f.TargetUsers != nil {
		// 目标用户可以是特定类型或all
		conditions["$or"] = []map[string]interface{}{
			{"target_users": *f.TargetUsers},
			{"target_users": "all"},
		}
	}

	return conditions
}

// GetSort 实现Filter接口
func (f *AnnouncementFilter) GetSort() map[string]int {
	sort := make(map[string]int)

	if f.SortBy != "" {
		sortValue := 1 // 默认升序
		if f.SortOrder == "desc" {
			sortValue = -1
		}

		switch f.SortBy {
		case "priority":
			sort["priority"] = sortValue
		case "view_count":
			sort["view_count"] = sortValue
		case "created_at":
			sort["created_at"] = sortValue
		default:
			sort["priority"] = -1 // 默认按优先级降序
		}
	} else {
		// 默认按优先级降序，创建时间降序
		sort["priority"] = -1
		sort["created_at"] = -1
	}

	return sort
}

// GetLimit 实现Filter接口
func (f *AnnouncementFilter) GetLimit() int {
	return f.Limit
}

// GetOffset 实现Filter接口
func (f *AnnouncementFilter) GetOffset() int {
	return f.Offset
}

// GetFields 实现Filter接口
func (f *AnnouncementFilter) GetFields() []string {
	return []string{} // Announcement查询返回所有字段
}

// Validate 实现Filter接口
func (f *AnnouncementFilter) Validate() error {
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

// IsEffective 检查公告是否在有效期内
func (a *Announcement) IsEffective() bool {
	if !a.IsActive {
		return false
	}

	now := time.Now()

	// 检查开始时间
	if a.StartTime != nil && now.Before(*a.StartTime) {
		return false
	}

	// 检查结束时间
	if a.EndTime != nil && now.After(*a.EndTime) {
		return false
	}

	return true
}
