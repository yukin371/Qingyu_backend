package writer

import (
	"Qingyu_backend/models/writer/base"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchOperation 批量操作模型
// 用于管理文档的批量操作，包括移动、删除、导出、复制等
type BatchOperation struct {
	base.IdentifiedEntity `bson:",inline"` // ID
	base.Timestamps       `bson:",inline"` // CreatedAt, UpdatedAt, DeletedAt

	// 项目信息
	ProjectID primitive.ObjectID `bson:"project_id" json:"projectId" validate:"required"`

	// 操作类型
	Type BatchOperationType `bson:"type" json:"type" validate:"required"`

	// 目标对象ID列表（文档ID、章节ID等）
	TargetIDs []string `bson:"target_ids" json:"targetIds" validate:"required,min=1"`

	// 操作状态
	Status BatchOperationStatus `bson:"status" json:"status"`

	// 操作载荷（根据type不同，payload内容不同）
	// 例如：move操作包含目标父节点、位置等；export操作包含格式等
	Payload map[string]interface{} `bson:"payload,omitempty" json:"payload,omitempty"`

	// 操作项列表（详细记录每个目标项的操作结果）
	Items []BatchOperationItem `bson:"items,omitempty" json:"items,omitempty"`

	// 是否原子操作（true：任一失败则全部回滚；false：部分失败继续）
	Atomic bool `bson:"atomic" json:"atomic"`

	// 冲突策略（当遇到冲突时的处理方式）
	ConflictPolicy ConflictPolicy `bson:"conflict_policy,omitempty" json:"conflictPolicy,omitempty"`

	// 预检查摘要（操作前的验证结果）
	PreflightSummary *PreflightSummary `bson:"preflight_summary,omitempty" json:"preflightSummary,omitempty"`

	// 客户端请求ID（用于幂等性）
	ClientRequestID string `bson:"client_request_id,omitempty" json:"clientRequestId,omitempty"`

	// 重试配置（存储为原始数据，避免循环依赖）
	RetryConfig map[string]interface{} `bson:"retry_config,omitempty" json:"retryConfig,omitempty"`

	// 错误信息（操作级别的错误）
	ErrorCode    string `bson:"error_code,omitempty" json:"errorCode,omitempty"`
	ErrorMessage string `bson:"error_message,omitempty" json:"errorMessage,omitempty"`

	// 开始时间
	StartedAt *primitive.DateTime `bson:"started_at,omitempty" json:"startedAt,omitempty"`

	// 完成时间
	CompletedAt *primitive.DateTime `bson:"completed_at,omitempty" json:"completedAt,omitempty"`

	// 创建者
	CreatedBy string `bson:"created_by" json:"createdBy" validate:"required"`
}

// BatchOperationType 批量操作类型常量
type BatchOperationType string

const (
	BatchOpTypeMove   BatchOperationType = "move"   // 批量移动
	BatchOpTypeDelete BatchOperationType = "delete" // 批量删除
	BatchOpTypeExport BatchOperationType = "export" // 批量导出
	BatchOpTypeCopy   BatchOperationType = "copy"   // 批量复制
	BatchOpTypeApply  BatchOperationType = "apply"  // 批量应用模板
)

// BatchOperationStatus 批量操作状态常量
type BatchOperationStatus string

const (
	BatchOpStatusPending        BatchOperationStatus = "pending"         // 待处理
	BatchOpStatusPreflight      BatchOperationStatus = "preflight"       // 预检查中
	BatchOpStatusRunning        BatchOperationStatus = "running"         // 执行中（别名，向后兼容）
	BatchOpStatusProcessing     BatchOperationStatus = "processing"      // 执行中
	BatchOpStatusCompleted      BatchOperationStatus = "completed"       // 已完成
	BatchOpStatusFailed         BatchOperationStatus = "failed"          // 失败
	BatchOpStatusCancelled      BatchOperationStatus = "cancelled"       // 已取消
	BatchOpStatusPartial        BatchOperationStatus = "partial"         // 部分成功（atomic=false时）
	BatchOpStatusPartiallyFailed BatchOperationStatus = "partially_failed" // 部分失败（别名，向后兼容）
)

// BatchOperationItem 批量操作项
// 记录单个目标项的操作结果
type BatchOperationItem struct {
	// 目标ID
	TargetID string `bson:"target_id" json:"targetId"`

	// 目标类型（document, chapter等）
	TargetType string `bson:"target_type,omitempty" json:"targetType,omitempty"`

	// 操作状态
	Status BatchItemStatus `bson:"status" json:"status"`

	// 错误信息（失败时记录）
	ErrorCode   string `bson:"error_code,omitempty" json:"errorCode,omitempty"`
	ErrorMsg    string `bson:"error_msg,omitempty" json:"errorMsg,omitempty"`
	ErrorMessage string `bson:"error_message,omitempty" json:"errorMessage,omitempty"` // 兼容字段

	// 是否可重试
	Retryable bool `bson:"retryable" json:"retryable"`

	// 重试次数
	RetryCount int `bson:"retry_count,omitempty" json:"retryCount,omitempty"`

	// 处理开始时间
	StartedAt *primitive.DateTime `bson:"started_at,omitempty" json:"startedAt,omitempty"`

	// 处理完成时间
	CompletedAt *primitive.DateTime `bson:"completed_at,omitempty" json:"completedAt,omitempty"`

	// 结果数据（成功时记录，如新文档ID等）
	Result map[string]interface{} `bson:"result,omitempty" json:"result,omitempty"`
}

// BatchItemStatus 批量操作项状态常量
type BatchItemStatus string

const (
	BatchItemStatusPending    BatchItemStatus = "pending"    // 待处理
	BatchItemStatusProcessing BatchItemStatus = "processing" // 处理中
	BatchItemStatusSucceeded  BatchItemStatus = "succeeded"  // 成功
	BatchItemStatusFailed     BatchItemStatus = "failed"     // 失败
	BatchItemStatusSkipped    BatchItemStatus = "skipped"    // 跳过
	BatchItemStatusCancelled  BatchItemStatus = "cancelled"  // 已取消
)

// ConflictPolicy 冲突策略常量
type ConflictPolicy string

const (
	ConflictPolicySkip     ConflictPolicy = "skip"     // 跳过冲突项
	ConflictPolicyOverwrite ConflictPolicy = "overwrite" // 覆盖
	ConflictPolicyRename   ConflictPolicy = "rename"    // 重命名（自动添加后缀）
	ConflictPolicyAbort    ConflictPolicy = "abort"     // 中止操作
)

// PreflightSummary 预检查摘要
// 记录操作前的验证结果
type PreflightSummary struct {
	// 总数
	TotalCount int `bson:"total_count" json:"totalCount"`

	// 有效项数
	ValidCount int `bson:"valid_count" json:"validCount"`

	// 无效项数
	InvalidCount int `bson:"invalid_count" json:"invalidCount"`

	// 跳过项数
	SkippedCount int `bson:"skipped_count" json:"skippedCount"`

	// 成功项数（执行后更新）
	SuccessCount int `bson:"success_count,omitempty" json:"successCount,omitempty"`

	// 失败项数（执行后更新）
	FailedCount int `bson:"failed_count,omitempty" json:"failedCount,omitempty"`

	// 警告信息
	Warnings []string `bson:"warnings,omitempty" json:"warnings,omitempty"`

	// 错误信息
	Errors []string `bson:"errors,omitempty" json:"errors,omitempty"`
}

// TouchForCreate 创建时设置默认值
func (b *BatchOperation) TouchForCreate() {
	b.IdentifiedEntity.GenerateID()
	b.Timestamps.TouchForCreate()

	if b.Status == "" {
		b.Status = BatchOpStatusPending
	}

	// 初始化items
	if b.Items == nil {
		b.Items = make([]BatchOperationItem, 0)
	}

	// 为每个targetID创建对应的item
	if len(b.Items) == 0 && len(b.TargetIDs) > 0 {
		b.Items = make([]BatchOperationItem, len(b.TargetIDs))
		for i, targetID := range b.TargetIDs {
			b.Items[i] = BatchOperationItem{
				TargetID: targetID,
				Status:   BatchItemStatusPending,
				Retryable: true,
			}
		}
	}
}

// TouchForUpdate 更新时设置默认值
func (b *BatchOperation) TouchForUpdate() {
	b.Timestamps.Touch()
}

// IsCompleted 判断操作是否完成
func (b *BatchOperation) IsCompleted() bool {
	return b.Status == BatchOpStatusCompleted ||
		b.Status == BatchOpStatusFailed ||
		b.Status == BatchOpStatusCancelled ||
		b.Status == BatchOpStatusPartial
}

// GetSummary 获取操作摘要
func (b *BatchOperation) GetSummary() *PreflightSummary {
	if b.PreflightSummary == nil {
		b.PreflightSummary = &PreflightSummary{}
	}

	// 统计items状态
	successCount := 0
	failedCount := 0
	skippedCount := 0

	for _, item := range b.Items {
		switch item.Status {
		case BatchItemStatusSucceeded:
			successCount++
		case BatchItemStatusFailed:
			failedCount++
		case BatchItemStatusSkipped:
			skippedCount++
		}
	}

	b.PreflightSummary.SuccessCount = successCount
	b.PreflightSummary.FailedCount = failedCount
	if skippedCount > 0 {
		b.PreflightSummary.SkippedCount = skippedCount
	}

	return b.PreflightSummary
}

// UpdateItemStatus 更新操作项状态
func (b *BatchOperation) UpdateItemStatus(targetID string, status BatchItemStatus, errCode, errMsg string) {
	for i := range b.Items {
		if b.Items[i].TargetID == targetID {
			b.Items[i].Status = status
			if errCode != "" {
				b.Items[i].ErrorCode = errCode
			}
			if errMsg != "" {
				b.Items[i].ErrorMsg = errMsg
				b.Items[i].ErrorMessage = errMsg
			}

			now := primitive.NewDateTimeFromTime(time.Now())
			if status == BatchItemStatusProcessing && b.Items[i].StartedAt == nil {
				b.Items[i].StartedAt = &now
			}
			if status == BatchItemStatusSucceeded || status == BatchItemStatusFailed {
				b.Items[i].CompletedAt = &now
			}

			break
		}
	}
}
