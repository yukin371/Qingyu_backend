package writer

import (
	"Qingyu_backend/models/writer/base"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchOperation 批量操作任务
type BatchOperation struct {
	base.IdentifiedEntity `bson:",inline"`
	base.Timestamps       `bson:",inline"`

	// 基础信息
	ProjectID         primitive.ObjectID   `bson:"project_id" json:"projectId" validate:"required"`
	Type              BatchOperationType  `bson:"type" json:"type" validate:"required"`
	TargetIDs         []string            `bson:"target_ids" json:"targetIds" validate:"required,min=1"`
	OriginalTargetIDs []string            `bson:"original_target_ids,omitempty" json:"originalTargetIds,omitempty"`

	// 执行配置
	ExecutionMode  ExecutionMode          `bson:"execution_mode,omitempty" json:"executionMode,omitempty"`
	Atomic         bool                   `bson:"atomic" json:"atomic"`
	Payload        map[string]interface{} `bson:"payload,omitempty" json:"payload,omitempty"`

	// 冲突处理
	ConflictPolicy   ConflictPolicy    `bson:"conflict_policy,omitempty" json:"conflictPolicy,omitempty"`
	ExpectedVersions map[string]int    `bson:"expected_versions,omitempty" json:"expectedVersions,omitempty"`
	ClientRequestID  string            `bson:"client_request_id,omitempty" json:"clientRequestId,omitempty"`

	// 状态
	Status     BatchOperationStatus `bson:"status" json:"status"`
	Cancelable bool                 `bson:"cancelable" json:"cancelable"`
	CreatedBy  primitive.ObjectID    `bson:"created_by" json:"createdBy" validate:"required"`
	StartedAt  *time.Time           `bson:"started_at,omitempty" json:"startedAt,omitempty"`
	FinishedAt *time.Time           `bson:"finished_at,omitempty" json:"finishedAt,omitempty"`

	// Preflight摘要
	PreflightSummary *PreflightSummary `bson:"preflight_summary,omitempty" json:"preflightSummary,omitempty"`
}

// BatchOperationType 批量操作类型
type BatchOperationType string

const (
	BatchOpTypeDelete        BatchOperationType = "delete"
	BatchOpTypeMove          BatchOperationType = "move"
	BatchOpTypeCopy          BatchOperationType = "copy"
	BatchOpTypeExport        BatchOperationType = "export"
	BatchOpTypeApplyTemplate BatchOperationType = "apply_template"
)

// ExecutionMode 执行模式
type ExecutionMode string

const (
	ExecutionModeStandardAtomic ExecutionMode = "standard_atomic" // <=200节点，单事务
	ExecutionModeSagaAtomic      ExecutionMode = "saga_atomic"      // >200节点，补偿事务
	ExecutionModeNonAtomic       ExecutionMode = "non_atomic"       // 允许部分成功
)

// ConflictPolicy 冲突策略
type ConflictPolicy string

const (
	ConflictPolicyAbort    ConflictPolicy = "abort"    // 中止整个操作
	ConflictPolicyOverwrite ConflictPolicy = "overwrite" // 覆盖
	ConflictPolicySkip     ConflictPolicy = "skip"     // 跳过冲突项
)

// BatchOperationStatus 批量操作状态
type BatchOperationStatus string

const (
	BatchOpStatusPending       BatchOperationStatus = "pending"
	BatchOpStatusRunning       BatchOperationStatus = "running"
	BatchOpStatusCompleted     BatchOperationStatus = "completed"
	BatchOpStatusFailed        BatchOperationStatus = "failed"
	BatchOpStatusCancelled     BatchOperationStatus = "cancelled"
	BatchOpStatusPartiallyFailed BatchOperationStatus = "partially_failed"
)

// PreflightSummary Preflight预检查摘要
type PreflightSummary struct {
	TotalCount   int `bson:"total_count" json:"totalCount"`
	ValidCount   int `bson:"valid_count" json:"validCount"`
	InvalidCount int `bson:"invalid_count" json:"invalidCount"`
	SkippedCount int `bson:"skipped_count" json:"skippedCount"`
}

// TouchForCreate 创建时设置默认值
func (b *BatchOperation) TouchForCreate() {
	b.IdentifiedEntity.GenerateID()
	b.Timestamps.TouchForCreate()
	if b.Status == "" {
		b.Status = BatchOpStatusPending
	}
	b.Cancelable = true
}

// IsRunning 判断是否正在运行
func (b *BatchOperation) IsRunning() bool {
	return b.Status == BatchOpStatusRunning
}

// IsTerminal 判断是否已终止
func (b *BatchOperation) IsTerminal() bool {
	return b.Status == BatchOpStatusCompleted ||
		b.Status == BatchOpStatusFailed ||
		b.Status == BatchOpStatusCancelled ||
		b.Status == BatchOpStatusPartiallyFailed
}

// CanCancel 判断是否可取消
func (b *BatchOperation) CanCancel() bool {
	return b.Cancelable && b.IsRunning()
}

// BatchOperationItem 批量操作子项
type BatchOperationItem struct {
	base.IdentifiedEntity `bson:",inline"`
	base.Timestamps       `bson:",inline"`

	// 关联
	BatchID        primitive.ObjectID `bson:"batch_id" json:"batchId" validate:"required"`
	TargetID       string             `bson:"target_id" json:"targetId" validate:"required"`
	TargetStableRef string            `bson:"target_stable_ref" json:"targetStableRef" validate:"required"`

	// 状态
	Status      BatchItemStatus `bson:"status" json:"status"`
	ErrorCode   string          `bson:"error_code,omitempty" json:"errorCode,omitempty"`
	ErrorMessage string         `bson:"error_message,omitempty" json:"errorMessage,omitempty"`
	SkipReason  string          `bson:"skip_reason,omitempty" json:"skipReason,omitempty"`

	// 版本控制
	ExpectedVersion *int `bson:"expected_version,omitempty" json:"expectedVersion,omitempty"`
	ActualVersion   *int `bson:"actual_version,omitempty" json:"actualVersion,omitempty"`

	// 撤销信息
	InverseCommand map[string]interface{} `bson:"inverse_command,omitempty" json:"inverseCommand,omitempty"`
	InverseLogID   primitive.ObjectID      `bson:"inverse_log_id,omitempty" json:"inverseLogId,omitempty"`

	// 时间
	StartedAt  *time.Time `bson:"started_at,omitempty" json:"startedAt,omitempty"`
	FinishedAt *time.Time `bson:"finished_at,omitempty" json:"finishedAt,omitempty"`
}

// BatchItemStatus 子项状态
type BatchItemStatus string

const (
	BatchItemStatusPending    BatchItemStatus = "pending"
	BatchItemStatusProcessing BatchItemStatus = "processing"
	BatchItemStatusSucceeded  BatchItemStatus = "succeeded"
	BatchItemStatusFailed     BatchItemStatus = "failed"
	BatchItemStatusSkipped    BatchItemStatus = "skipped"
	BatchItemStatusCancelled  BatchItemStatus = "cancelled"
)

// TouchForCreate 创建时设置默认值
func (bi *BatchOperationItem) TouchForCreate() {
	bi.IdentifiedEntity.GenerateID()
	bi.Timestamps.TouchForCreate()
	if bi.Status == "" {
		bi.Status = BatchItemStatusPending
	}
}
