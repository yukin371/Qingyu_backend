package users

import "time"

// UserFilter 用户查询过滤条件
type UserFilter struct {
	// 基础字段
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`

	// 角色和状态
	Role   string     `json:"role,omitempty"`
	Status UserStatus `json:"status,omitempty"`

	// 验证状态
	EmailVerified *bool `json:"emailVerified,omitempty"`
	PhoneVerified *bool `json:"phoneVerified,omitempty"`

	// 时间范围
	CreatedAfter  *time.Time `json:"createdAfter,omitempty"`
	CreatedBefore *time.Time `json:"createdBefore,omitempty"`

	// 搜索
	SearchKeyword string `json:"searchKeyword,omitempty"` // 搜索用户名、昵称、邮箱

	// 分页
	Page     int `json:"page,omitempty"`
	PageSize int `json:"pageSize,omitempty"`

	// 排序
	SortBy    string `json:"sortBy,omitempty"`    // created_at, updated_at, username
	SortOrder string `json:"sortOrder,omitempty"` // asc, desc
}

// SetDefaults 设置默认值
func (f *UserFilter) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	if f.SortBy == "" {
		f.SortBy = "created_at"
	}
	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}
}

// GetSkip 计算跳过的记录数
func (f *UserFilter) GetSkip() int {
	return (f.Page - 1) * f.PageSize
}

// GetLimit 获取限制数量
func (f *UserFilter) GetLimit() int {
	return f.PageSize
}
