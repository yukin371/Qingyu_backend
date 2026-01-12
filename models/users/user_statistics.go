package users

import "time"

// UserStatistics 用户统计信息
type UserStatistics struct {
	UserID             string    `json:"userId"`
	TotalBooks         int       `json:"totalBooks"`         // 发布书籍总数
	TotalWords         int64     `json:"totalWords"`         // 总字数
	TotalChapters      int       `json:"totalChapters"`      // 总章节数
	TotalReaders       int       `json:"totalReaders"`       // 读者总数
	TotalComments      int       `json:"totalComments"`      // 评论总数
	TotalLikes         int       `json:"totalLikes"`         // 点赞总数
	TotalCollections   int       `json:"totalCollections"`   // 收藏总数
	ReadingTimeMinutes int64     `json:"readingTimeMinutes"` // 阅读时长（分钟）
	LastActiveAt       time.Time `json:"lastActiveAt"`       // 最后活跃时间
	RegisterDays       int       `json:"registerDays"`       // 注册天数
	AvgDailyWords      float64   `json:"avgDailyWords"`      // 平均日字数
	MostActiveMonth    string    `json:"mostActiveMonth"`    // 最活跃月份
}

// UserActivity 用户活动记录
type UserActivity struct {
	ID          string            `json:"id"`
	UserID      string            `json:"userId"`
	Action      string            `json:"action"` // 登录、注册、发布、评论等
	Description string            `json:"description"`
	IPAddress   string            `json:"ipAddress,omitempty"`
	UserAgent   string            `json:"userAgent,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
}

// UserManagementInfo 用户管理信息（管理员视图）
type UserManagementInfo struct {
	User             *User           `json:"user"`
	Statistics       *UserStatistics `json:"statistics,omitempty"`
	RecentActivities []*UserActivity `json:"recentActivities,omitempty"`
	CanModify        bool            `json:"canModify"`   // 是否可以修改
	CanDelete        bool            `json:"canDelete"`   // 是否可以删除
	CanResetPwd      bool            `json:"canResetPwd"` // 是否可以重置密码
}
