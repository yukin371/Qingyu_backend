package messaging

import (
	"time"

	"Qingyu_backend/models/messaging/base"
)

// Announcement 公告模型
type Announcement struct {
	base.IdentifiedEntity `bson:",inline"`
	base.Timestamps       `bson:",inline"`
	base.TitledEntity     `bson:",inline"`
	base.Pinned           `bson:",inline"`
	base.Expirable        `bson:",inline"`

	// 公告内容
	Content string `bson:"content" json:"content" validate:"required,min=1"` // 内容（支持HTML）

	// 分类和优先级
	Type     AnnouncementType `bson:"type" json:"type" validate:"required,oneof=info warning notice"` // 类型：info信息 warning警告 notice通知
	Priority int              `bson:"priority" json:"priority"`                                       // 优先级（数字越大优先级越高）

	// 状态和目标
	IsActive   bool   `bson:"is_active" json:"isActive"`                                                       // 是否启用
	TargetRole string `bson:"target_role" json:"targetRole" validate:"required,oneof=all reader writer admin"` // 目标用户：all所有 reader读者 writer作者 admin管理员

	// 统计数据
	ViewCount int64 `bson:"view_count" json:"viewCount"` // 查看次数

	// 创建者
	CreatedBy string `bson:"created_by" json:"createdBy"` // 创建者ID

	// 有效期（使用Expirable混入的ExpiresAt字段，同时保留StartTime以支持延迟发布）
	StartTime *time.Time `bson:"start_time,omitempty" json:"startTime,omitempty"` // 开始时间（延迟发布）
	EndTime   *time.Time `bson:"end_time,omitempty" json:"endTime,omitempty"`     // 结束时间
}

// AnnouncementType 公告类型
type AnnouncementType string

const (
	AnnouncementTypeInfo    AnnouncementType = "info"    // 信息
	AnnouncementTypeWarning AnnouncementType = "warning" // 警告
	AnnouncementTypeNotice  AnnouncementType = "notice"  // 通知
)

// AnnouncementFilter 公告查询过滤器
type AnnouncementFilter struct {
	IsActive   *bool             `json:"isActive,omitempty"`
	Type       *AnnouncementType `json:"type,omitempty"`
	TargetRole *string           `json:"targetRole,omitempty"`
	IsPinned   *bool             `json:"isPinned,omitempty"`
	SortBy     string            `json:"sortBy,omitempty"`    // priority, created_at, view_count
	SortOrder  string            `json:"sortOrder,omitempty"` // asc, desc
	Limit      int               `json:"limit,omitempty"`
	Offset     int               `json:"offset,omitempty"`
}

// GetConditions 获取查询条件
func (f *AnnouncementFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	if f.IsActive != nil {
		conditions["is_active"] = *f.IsActive
	}
	if f.Type != nil {
		conditions["type"] = *f.Type
	}
	if f.TargetRole != nil {
		// 目标用户可以是特定类型或all
		conditions["$or"] = []map[string]interface{}{
			{"target_role": *f.TargetRole},
			{"target_role": "all"},
		}
	}
	if f.IsPinned != nil {
		conditions["is_pinned"] = *f.IsPinned
	}

	// 排除过期公告
	conditions["$or"] = []map[string]interface{}{
		{"expires_at": nil},
		{"expires_at": map[string]interface{}{"$gt": time.Now()}},
	}

	return conditions
}

// GetSort 获取排序
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
		// 默认：置顶优先，然后按优先级降序，创建时间降序
		sort["is_pinned"] = -1
		sort["priority"] = -1
		sort["created_at"] = -1
	}

	return sort
}

// GetLimit 获取限制
func (f *AnnouncementFilter) GetLimit() int {
	return f.Limit
}

// GetOffset 获取偏移
func (f *AnnouncementFilter) GetOffset() int {
	return f.Offset
}

// GetFields 获取字段
func (f *AnnouncementFilter) GetFields() []string {
	return []string{} // Announcement查询返回所有字段
}

// Validate 验证
func (f *AnnouncementFilter) Validate() error {
	if f.Limit < 0 {
		return &ValidationError{Message: "limit不能为负数"}
	}
	if f.Offset < 0 {
		return &ValidationError{Message: "offset不能为负数"}
	}
	return nil
}

// IsEffective 检查公告是否有效（激活且在有效期内）
func (a *Announcement) IsEffective() bool {
	if !a.IsActive {
		return false
	}

	// 检查是否已过期
	if a.IsExpired() {
		return false
	}

	now := time.Now()

	// 检查开始时间（延迟发布）
	if a.StartTime != nil && now.Before(*a.StartTime) {
		return false
	}

	// 检查结束时间
	if a.EndTime != nil && now.After(*a.EndTime) {
		return false
	}

	return true
}

// ShouldShow 判断公告是否应该显示给指定角色
func (a *Announcement) ShouldShow(role string) bool {
	if !a.IsEffective() {
		return false
	}

	// all表示显示给所有用户
	if a.TargetRole == "all" {
		return true
	}

	// 检查角色匹配
	return a.TargetRole == role
}

// Publish 发布公告
func (a *Announcement) Publish() {
	a.IsActive = true
	a.Touch()
}

// Unpublish 取消发布
func (a *Announcement) Unpublish() {
	a.IsActive = false
	a.Touch()
}

// IncrementView 增加查看次数
func (a *Announcement) IncrementView() {
	a.ViewCount++
	a.Touch()
}

// SetExpirationByEndTime 根据EndTime设置过期时间
func (a *Announcement) SetExpirationByEndTime() {
	if a.EndTime != nil {
		a.ExpiresAt = a.EndTime
	}
}
