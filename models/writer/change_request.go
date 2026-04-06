package writer

import (
	"Qingyu_backend/models/writer/base"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- 枚举 ---

// ChangeRequestStatus 建议状态
type ChangeRequestStatus string

const (
	CRStatusPending  ChangeRequestStatus = "pending"
	CRStatusAccepted ChangeRequestStatus = "accepted"
	CRStatusIgnored  ChangeRequestStatus = "ignored"
	CRStatusDeferred ChangeRequestStatus = "deferred"
)

// ChangeRequestPriority 建议优先级
type ChangeRequestPriority string

const (
	CRPriorityHigh   ChangeRequestPriority = "high"
	CRPriorityMedium ChangeRequestPriority = "medium"
	CRPriorityLow    ChangeRequestPriority = "low"
)

// ChangeRequestCategory 建议类别
type ChangeRequestCategory string

const (
	CRCategoryCharacterState  ChangeRequestCategory = "character_state"
	CRCategoryRelationChange  ChangeRequestCategory = "relation_change"
	CRCategoryInlineDirective ChangeRequestCategory = "inline_directive"
	CRCategoryScopeDrift      ChangeRequestCategory = "scope_drift"
)

// --- 证据引用 ---

// EvidenceRef 证据引用
type EvidenceRef struct {
	DocumentID   string `bson:"document_id" json:"documentId"`
	ParagraphIdx int    `bson:"paragraph_idx" json:"paragraphIdx"`
	QuoteText    string `bson:"quote_text" json:"quoteText" validate:"max=500"`
}

// --- ChangeRequest 单条变更建议 ---

// ChangeRequest 单条变更建议
type ChangeRequest struct {
	base.IdentifiedEntity    `bson:",inline"`
	base.Timestamps          `bson:",inline"`
	base.ProjectScopedEntity `bson:",inline"`

	BatchID   primitive.ObjectID    `bson:"batch_id" json:"batchId"`
	ChapterID primitive.ObjectID    `bson:"chapter_id" json:"chapterId"`

	Category   ChangeRequestCategory `bson:"category" json:"category"`
	Priority   ChangeRequestPriority `bson:"priority" json:"priority"`
	Status     ChangeRequestStatus   `bson:"status" json:"status"`
	Title      string                `bson:"title" json:"title" validate:"required,max=200"`
	Description string               `bson:"description,omitempty" json:"description,omitempty" validate:"max=1000"`

	// 建议的具体变更内容（JSON 自由结构，后续按 category 细化）
	SuggestedChange map[string]interface{} `bson:"suggested_change,omitempty" json:"suggestedChange,omitempty"`

	// 证据链
	Evidence []EvidenceRef `bson:"evidence,omitempty" json:"evidence,omitempty" validate:"max=10"`

	// 处理信息
	ProcessedAt *primitive.DateTime `bson:"processed_at,omitempty" json:"processedAt,omitempty"`
	ProcessedBy *string             `bson:"processed_by,omitempty" json:"processedBy,omitempty"`

	// 来源
	Source string `bson:"source" json:"source"` // "rule" | "ai"
}

// TouchForCreate 创建时设置默认值
func (cr *ChangeRequest) TouchForCreate() {
	cr.IdentifiedEntity.GenerateID()
	cr.Timestamps.TouchForCreate()
	if cr.Status == "" {
		cr.Status = CRStatusPending
	}
	if cr.Priority == "" {
		cr.Priority = CRPriorityMedium
	}
}

// Validate 验证数据
func (cr *ChangeRequest) Validate() error {
	if cr.ProjectID.IsZero() {
		return base.ErrProjectIDRequired
	}
	if cr.Title == "" {
		return base.ErrTitleRequired
	}
	return nil
}

// --- ChangeRequestBatch 建议批次 ---

// ChangeRequestBatch 一次保存触发的建议批次
type ChangeRequestBatch struct {
	base.IdentifiedEntity    `bson:",inline"`
	base.Timestamps          `bson:",inline"`
	base.ProjectScopedEntity `bson:",inline"`

	ChapterID primitive.ObjectID `bson:"chapter_id" json:"chapterId"`
	Source    string             `bson:"source" json:"source"` // "rule" | "ai" | "save"

	TotalCount   int `bson:"total_count" json:"totalCount"`
	PendingCount int `bson:"pending_count" json:"pendingCount"`

	SaveTriggerID *string `bson:"save_trigger_id,omitempty" json:"saveTriggerId,omitempty"`
}

// TouchForCreate 创建时设置默认值
func (b *ChangeRequestBatch) TouchForCreate() {
	b.IdentifiedEntity.GenerateID()
	b.Timestamps.TouchForCreate()
}
