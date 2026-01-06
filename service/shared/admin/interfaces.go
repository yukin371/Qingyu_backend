package admin

import (
	"context"
	"time"
)

// AdminService 管理后台服务接口（对外暴露）
type AdminService interface {
	// 内容审核
	ReviewContent(ctx context.Context, req *ReviewContentRequest) error
	GetPendingReviews(ctx context.Context, contentType string) ([]*AuditRecord, error)

	// 用户管理
	ManageUser(ctx context.Context, req *ManageUserRequest) error
	BanUser(ctx context.Context, userID, reason string, duration time.Duration) error
	UnbanUser(ctx context.Context, userID string) error
	GetUserStatistics(ctx context.Context, userID string) (*UserStatistics, error)

	// 提现审核
	ReviewWithdraw(ctx context.Context, withdrawID, adminID string, approved bool, reason string) error

	// 操作日志
	LogOperation(ctx context.Context, req *LogOperationRequest) error
	GetOperationLogs(ctx context.Context, req *GetLogsRequest) ([]*AdminLog, error)
	ExportLogs(ctx context.Context, startDate, endDate time.Time) ([]byte, error)

	// 系统管理
	GetSystemStats(ctx context.Context) (interface{}, error)
	GetSystemConfig(ctx context.Context) (interface{}, error)
	UpdateSystemConfig(ctx context.Context, req interface{}) error
	CreateAnnouncement(ctx context.Context, adminID, title, content, announceType, priority string) (interface{}, error)
	GetAnnouncements(ctx context.Context, page, pageSize int) (interface{}, int64, error)

	// 审核管理
	GetAuditStatistics(ctx context.Context) (interface{}, error)

	// Health 健康检查
	Health(ctx context.Context) error
}

// ============ 请求结构 ============

// ReviewContentRequest 审核内容请求
type ReviewContentRequest struct {
	ContentID   string `json:"content_id" binding:"required"`
	ContentType string `json:"content_type" binding:"required"` // book, chapter, comment
	Action      string `json:"action" binding:"required"`       // approve, reject
	Reason      string `json:"reason,omitempty"`
	ReviewerID  string `json:"reviewer_id" binding:"required"`
}

// ManageUserRequest 管理用户请求
type ManageUserRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Action   string `json:"action" binding:"required"` // ban, unban, delete
	Reason   string `json:"reason,omitempty"`
	Duration int64  `json:"duration,omitempty"` // 封禁时长（秒）
	AdminID  string `json:"admin_id" binding:"required"`
}

// LogOperationRequest 记录操作日志请求
type LogOperationRequest struct {
	AdminID   string                 `json:"admin_id" binding:"required"`
	Operation string                 `json:"operation" binding:"required"`
	Target    string                 `json:"target"`
	Details   map[string]interface{} `json:"details,omitempty"`
	IP        string                 `json:"ip"`
}

// GetLogsRequest 查询日志请求
type GetLogsRequest struct {
	AdminID   string    `json:"admin_id"`
	Operation string    `json:"operation"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
}

// ============ 数据结构 ============

// AuditRecord 审核记录
type AuditRecord struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ContentID   string    `json:"content_id" bson:"content_id"`
	ContentType string    `json:"content_type" bson:"content_type"`
	Status      string    `json:"status" bson:"status"` // pending, approved, rejected
	ReviewerID  string    `json:"reviewer_id,omitempty" bson:"reviewer_id,omitempty"`
	Reason      string    `json:"reason,omitempty" bson:"reason,omitempty"`
	ReviewedAt  time.Time `json:"reviewed_at,omitempty" bson:"reviewed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// AdminLog 管理员操作日志
type AdminLog struct {
	ID        string                 `json:"id" bson:"_id,omitempty"`
	AdminID   string                 `json:"admin_id" bson:"admin_id"`
	Operation string                 `json:"operation" bson:"operation"`
	Target    string                 `json:"target,omitempty" bson:"target,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty" bson:"details,omitempty"`
	IP        string                 `json:"ip" bson:"ip"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
}

// UserStatistics 用户统计信息
type UserStatistics struct {
	UserID           string    `json:"user_id"`
	TotalBooks       int       `json:"total_books"`
	TotalChapters    int       `json:"total_chapters"`
	TotalWords       int64     `json:"total_words"`
	TotalReads       int64     `json:"total_reads"`
	TotalIncome      float64   `json:"total_income"`
	RegistrationDate time.Time `json:"registration_date"`
	LastLoginDate    time.Time `json:"last_login_date"`
}
