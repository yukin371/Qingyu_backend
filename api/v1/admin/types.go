package admin

import (
	"Qingyu_backend/api/v1/shared"
	"time"
)

// DTO (Data Transfer Object) - Admin API层的请求和响应结构

// ===========================
// 用户管理相关DTO
// ===========================

// ListUsersRequest 获取用户列表请求（查询参数）
type ListUsersRequest struct {
	Page     int    `form:"page" validate:"omitempty,min=1"`
	PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100"`
	Username string `form:"username" validate:"omitempty"`
	Email    string `form:"email" validate:"omitempty,email"`
	Role     string `form:"role" validate:"omitempty"`
	Status   string `form:"status" validate:"omitempty"`
}

// AdminUpdateUserRequest 管理员更新用户请求
type AdminUpdateUserRequest struct {
	Nickname      *string `json:"nickname,omitempty" validate:"omitempty,max=50"`
	Bio           *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	Avatar        *string `json:"avatar,omitempty" validate:"omitempty,url"`
	Phone         *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Role          *string `json:"role,omitempty" validate:"omitempty,oneof=reader author admin"`
	Status        *string `json:"status,omitempty" validate:"omitempty"`
	EmailVerified *bool   `json:"email_verified,omitempty"`
	PhoneVerified *bool   `json:"phone_verified,omitempty"`
}

// UserProfileResponse 用户信息响应（引用共享定义）
// 保留类型别名以保持向后兼容
type UserProfileResponse = shared.UserProfileResponse

// ===========================
// AI配额管理相关DTO
// ===========================

// UpdateQuotaRequest 更新配额请求
type UpdateQuotaRequest struct {
	TotalQuota int    `json:"totalQuota" binding:"required,min=0"`
	QuotaType  string `json:"quotaType" binding:"required,oneof=daily monthly total"`
}

// ===========================
// 审核管理相关DTO
// ===========================

// ReviewContentRequest 审核内容请求
type ReviewContentRequest struct {
	ContentID   string `json:"content_id" binding:"required,min=1"`
	ContentType string `json:"content_type" binding:"required"`
	Action      string `json:"action" binding:"required,oneof=approve reject"`
	Reason      string `json:"reason" validate:"omitempty,max=500"`
}

// ReviewAuditRequest 审核记录处理请求
type ReviewAuditRequest struct {
	Action      string `json:"action" binding:"required,oneof=approve reject"`
	ReviewNote  string `json:"review_note" validate:"omitempty,max=500"`
	PenaltyType string `json:"penalty_type,omitempty" validate:"omitempty,oneof=warning ban"`
}

// ReviewAppealRequest 审核申诉请求
type ReviewAppealRequest struct {
	Action     string `json:"action" binding:"required,oneof=approve reject"`
	ReviewNote string `json:"review_note" validate:"omitempty,max=500"`
}

// ===========================
// 系统管理相关DTO
// ===========================

// ReviewWithdrawRequest 审核提现请求
type ReviewWithdrawRequest struct {
	WithdrawID string `json:"withdraw_id" binding:"required"`
	Approved   bool   `json:"approved"`
	Reason     string `json:"reason"`
}

// GetLogsRequest 获取操作日志请求
type GetLogsRequest struct {
	Page      int    `form:"page" validate:"omitempty,min=1"`
	PageSize  int    `form:"page_size" validate:"omitempty,min=1,max=100"`
	AdminID   string `form:"admin_id" validate:"omitempty"`
	Operation string `form:"operation" validate:"omitempty"`
}

// SystemConfigRequest 系统配置请求
type SystemConfigRequest struct {
	AllowRegistration        *bool  `json:"allowRegistration,omitempty"`
	RequireEmailVerification *bool  `json:"requireEmailVerification,omitempty"`
	MaxUploadSize            *int64 `json:"maxUploadSize,omitempty"`
	EnableAudit              *bool  `json:"enableAudit,omitempty"`
}

// SystemStatsResponse 系统统计响应
type SystemStatsResponse struct {
	TotalUsers    int64   `json:"totalUsers"`
	ActiveUsers   int64   `json:"activeUsers"`
	TotalBooks    int64   `json:"totalBooks"`
	TotalRevenue  float64 `json:"totalRevenue"`
	PendingAudits int64   `json:"pendingAudits"`
}

// BanUserRequest 封禁用户请求
type BanUserRequest struct {
	Reason       string     `json:"reason" binding:"required"`
	Duration     int        `json:"duration" binding:"required,min=1"`
	DurationUnit string     `json:"durationUnit" binding:"required,oneof=hours days weeks months"`
	BanUntil     *time.Time `json:"banUntil,omitempty"`
}
