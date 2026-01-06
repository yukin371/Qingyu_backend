package users

import "time"

// AuditRecord 审核记录
type AuditRecord struct {
	ID          string                 `json:"id" bson:"_id,omitempty"`
	ContentID   string                 `json:"content_id" bson:"content_id"`
	ContentType string                 `json:"content_type" bson:"content_type"` // book, chapter, comment
	Status      string                 `json:"status" bson:"status"`             // pending, approved, rejected
	ReviewerID  string                 `json:"reviewer_id,omitempty" bson:"reviewer_id,omitempty"`
	Reason      string                 `json:"reason,omitempty" bson:"reason,omitempty"`
	ReviewedAt  time.Time              `json:"reviewed_at,omitempty" bson:"reviewed_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"` // 额外信息
	CreatedAt   time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" bson:"updated_at"`
}

// AdminLog 管理员操作日志
type AdminLog struct {
	ID         string                 `json:"id" bson:"_id,omitempty"`
	AdminID    string                 `json:"admin_id" bson:"admin_id"`
	Operation  string                 `json:"operation" bson:"operation"`                         // review_content, ban_user, approve_withdraw
	Target     string                 `json:"target,omitempty" bson:"target,omitempty"`           // 操作对象ID
	TargetType string                 `json:"target_type,omitempty" bson:"target_type,omitempty"` // user, content, withdraw
	Details    map[string]interface{} `json:"details,omitempty" bson:"details,omitempty"`
	IP         string                 `json:"ip" bson:"ip"`
	UserAgent  string                 `json:"user_agent,omitempty" bson:"user_agent,omitempty"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
}

// 审核状态
const (
	AuditStatusPending  = "pending"  // 待审核
	AuditStatusApproved = "approved" // 已通过
	AuditStatusRejected = "rejected" // 已驳回
)

// 内容类型
const (
	ContentTypeBook    = "book"    // 书籍
	ContentTypeChapter = "chapter" // 章节
	ContentTypeComment = "comment" // 评论
	ContentTypeArticle = "article" // 文章
)

// 操作类型
const (
	OperationReviewContent   = "review_content"   // 审核内容
	OperationBanUser         = "ban_user"         // 封禁用户
	OperationUnbanUser       = "unban_user"       // 解封用户
	OperationDeleteUser      = "delete_user"      // 删除用户
	OperationApproveWithdraw = "approve_withdraw" // 批准提现
	OperationRejectWithdraw  = "reject_withdraw"  // 驳回提现
	OperationUpdateRole      = "update_role"      // 更新角色
	OperationModifyContent   = "modify_content"   // 修改内容
)
