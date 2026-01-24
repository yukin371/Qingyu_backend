package dto

// ===========================
// 审计 DTO（符合分层架构规范）
// ===========================

// AuditRecordDTO 审计记录数据传输对象
// 用于：Service 层和 API 层数据传输，ID 和时间字段使用字符串类型
type AuditRecordDTO struct {
	ID         string `json:"id" validate:"required"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	ReviewedAt string `json:"reviewedAt,omitempty"` // 复核时间（ISO8601）

	// 审核对象
	TargetType string `json:"targetType" validate:"required,oneof=book chapter comment user"` // 审核对象类型
	TargetID   string `json:"targetId" validate:"required"`                                     // 审核对象ID
	AuthorID   string `json:"authorId" validate:"required"`                                     // 作者ID

	// 审核内容
	Content string `json:"content,omitempty"` // 审核内容
	Status  string `json:"status" validate:"required,oneof=pending approved rejected"` // 审核状态
	Result  string `json:"result,omitempty"`  // 审核结果

	// 风险评估
	RiskLevel  int     `json:"riskLevel" validate:"min=0,max=5"`    // 风险等级 0-5
	RiskScore  float32 `json:"riskScore" validate:"min=0,max=10"`   // 风险分数 0-10
	Violations []ViolationDetailDTO                                   // 违规详情

	// 复核信息
	ReviewerID  string `json:"reviewerId,omitempty"`  // 复核人ID
	ReviewNote  string `json:"reviewNote,omitempty"`  // 复核备注

	// 申诉
	AppealStatus string `json:"appealStatus,omitempty" validate:"oneof=none pending approved rejected"` // 申诉状态

	// 元数据
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 额外元数据
}

// ViolationDetailDTO 违规详情数据传输对象
type ViolationDetailDTO struct {
	Type        string   `json:"type" validate:"required"`       // 违规类型
	Category    string   `json:"category" validate:"required"`   // 违规类别
	Level       int      `json:"level" validate:"required"`      // 违规级别
	Description string   `json:"description,omitempty"`         // 违规描述
	Position    int      `json:"position,omitempty"`            // 违规位置
	Context     string   `json:"context,omitempty"`             // 违规上下文
	Keywords    []string `json:"keywords,omitempty"`            // 违规关键词
}
