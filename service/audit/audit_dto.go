package audit

import "time"

// CheckContentRequest 实时检测请求
type CheckContentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=100000"`
}

// AuditDocumentRequest 全文审核请求
type AuditDocumentRequest struct {
	DocumentID string `json:"documentId" validate:"required"`
	Content    string `json:"content" validate:"required,min=1,max=100000"`
}

// ReviewAuditRequest 复核请求
type ReviewAuditRequest struct {
	Approved bool   `json:"approved"`
	Note     string `json:"note" validate:"max=500"`
}

// SubmitAppealRequest 申诉请求
type SubmitAppealRequest struct {
	Reason string `json:"reason" validate:"required,min=10,max=500"`
}

// ReviewAppealRequest 复核申诉请求
type ReviewAppealRequest struct {
	Approved bool   `json:"approved"`
	Note     string `json:"note" validate:"max=500"`
}

// GetAuditRecordsRequest 查询审核记录请求
type GetAuditRecordsRequest struct {
	TargetType string `form:"targetType"`
	Status     string `form:"status"`
	AuthorID   string `form:"authorId"`
	Page       int    `form:"page" validate:"min=1"`
	PageSize   int    `form:"pageSize" validate:"min=1,max=100"`
}

// AuditRecordResponse 审核记录响应
type AuditRecordResponse struct {
	ID           string      `json:"id"`
	TargetType   string      `json:"targetType"`
	TargetID     string      `json:"targetId"`
	AuthorID     string      `json:"authorId"`
	Status       string      `json:"status"`
	Result       string      `json:"result"`
	RiskLevel    int         `json:"riskLevel"`
	RiskScore    float64     `json:"riskScore"`
	Violations   interface{} `json:"violations"`
	ReviewerID   string      `json:"reviewerId,omitempty"`
	ReviewNote   string      `json:"reviewNote,omitempty"`
	AppealStatus string      `json:"appealStatus,omitempty"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	ReviewedAt   *time.Time  `json:"reviewedAt,omitempty"`
	CanAppeal    bool        `json:"canAppeal"`
}

// ViolationRecordResponse 违规记录响应
type ViolationRecordResponse struct {
	ID              string     `json:"id"`
	UserID          string     `json:"userId"`
	TargetType      string     `json:"targetType"`
	TargetID        string     `json:"targetId"`
	ViolationType   string     `json:"violationType"`
	ViolationLevel  int        `json:"violationLevel"`
	ViolationCount  int        `json:"violationCount"`
	PenaltyType     string     `json:"penaltyType,omitempty"`
	PenaltyDuration int        `json:"penaltyDuration,omitempty"`
	IsPenalized     bool       `json:"isPenalized"`
	Description     string     `json:"description"`
	CreatedAt       time.Time  `json:"createdAt"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty"`
	IsActive        bool       `json:"isActive"`
}

// UserViolationSummaryResponse 用户违规统计响应
type UserViolationSummaryResponse struct {
	UserID              string    `json:"userId"`
	TotalViolations     int       `json:"totalViolations"`
	WarningCount        int       `json:"warningCount"`
	RejectCount         int       `json:"rejectCount"`
	HighRiskCount       int       `json:"highRiskCount"`
	LastViolationAt     time.Time `json:"lastViolationAt"`
	ActivePenalties     int       `json:"activePenalties"`
	IsBanned            bool      `json:"isBanned"`
	IsPermanentlyBanned bool      `json:"isPermanentlyBanned"`
	IsHighRiskUser      bool      `json:"isHighRiskUser"`
	ShouldBan           bool      `json:"shouldBan"`
}

// PaginatedAuditRecordsResponse 分页审核记录响应
type PaginatedAuditRecordsResponse struct {
	Records    []AuditRecordResponse `json:"records"`
	TotalCount int64                 `json:"totalCount"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	TotalPages int                   `json:"totalPages"`
}
