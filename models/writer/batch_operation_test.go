package writer

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBatchOperation_TouchForCreate(t *testing.T) {
	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1", "doc-2"},
		Atomic:    true,
		CreatedBy: "user123",
	}

	op.TouchForCreate()

	if op.ID.IsZero() {
		t.Error("ID should be generated")
	}
	if op.Status != BatchOpStatusPending {
		t.Errorf("Status should be pending, got %s", op.Status)
	}
	if op.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if op.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
	if len(op.Items) != 2 {
		t.Errorf("Items length should be 2, got %d", len(op.Items))
	}
}

func TestBatchOperation_TouchForCreate_WithExistingStatus(t *testing.T) {
	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1"},
		Status:    BatchOpStatusRunning,
		CreatedBy: "user123",
	}

	op.TouchForCreate()

	// Status should not be overwritten if already set
	if op.Status != BatchOpStatusRunning {
		t.Errorf("Status should remain running, got %s", op.Status)
	}
}

func TestBatchOperation_IsCompleted(t *testing.T) {
	tests := []struct {
		name   string
		status BatchOperationStatus
		want   bool
	}{
		{"pending", BatchOpStatusPending, false},
		{"running", BatchOpStatusRunning, false},
		{"processing", BatchOpStatusProcessing, false},
		{"completed", BatchOpStatusCompleted, true},
		{"failed", BatchOpStatusFailed, true},
		{"cancelled", BatchOpStatusCancelled, true},
		{"partial", BatchOpStatusPartial, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BatchOperation{Status: tt.status}
			if got := op.IsCompleted(); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBatchOperationItem_Structure(t *testing.T) {
	item := &BatchOperationItem{
		TargetID:    "doc-1",
		TargetType:  "document",
		Status:      BatchItemStatusPending,
		Retryable:   true,
		RetryCount:  0,
		ErrorCode:   "",
		ErrorMsg:    "",
		ErrorMessage: "",
	}

	if item.TargetID != "doc-1" {
		t.Errorf("TargetID should be doc-1, got %s", item.TargetID)
	}
	if item.Status != BatchItemStatusPending {
		t.Errorf("Status should be pending, got %s", item.Status)
	}
	if !item.Retryable {
		t.Error("Retryable should be true")
	}
}

func TestBatchOperationType_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value BatchOperationType
	}{
		{"delete", BatchOpTypeDelete},
		{"move", BatchOpTypeMove},
		{"copy", BatchOpTypeCopy},
		{"export", BatchOpTypeExport},
		{"apply", BatchOpTypeApply},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("BatchOperationType %s should not be empty", tt.name)
			}
		})
	}
}

func TestBatchOperationStatus_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value BatchOperationStatus
	}{
		{"pending", BatchOpStatusPending},
		{"preflight", BatchOpStatusPreflight},
		{"running", BatchOpStatusRunning},
		{"processing", BatchOpStatusProcessing},
		{"completed", BatchOpStatusCompleted},
		{"failed", BatchOpStatusFailed},
		{"cancelled", BatchOpStatusCancelled},
		{"partial", BatchOpStatusPartial},
		{"partially_failed", BatchOpStatusPartiallyFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("BatchOperationStatus %s should not be empty", tt.name)
			}
		})
	}
}

func TestBatchItemStatus_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value BatchItemStatus
	}{
		{"pending", BatchItemStatusPending},
		{"processing", BatchItemStatusProcessing},
		{"succeeded", BatchItemStatusSucceeded},
		{"failed", BatchItemStatusFailed},
		{"skipped", BatchItemStatusSkipped},
		{"cancelled", BatchItemStatusCancelled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("BatchItemStatus %s should not be empty", tt.name)
			}
		})
	}
}

func TestConflictPolicy_Constants(t *testing.T) {
	tests := []struct {
		name  string
		value ConflictPolicy
	}{
		{"skip", ConflictPolicySkip},
		{"overwrite", ConflictPolicyOverwrite},
		{"rename", ConflictPolicyRename},
		{"abort", ConflictPolicyAbort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) == "" {
				t.Errorf("ConflictPolicy %s should not be empty", tt.name)
			}
		})
	}
}

func TestPreflightSummary_Structure(t *testing.T) {
	summary := PreflightSummary{
		TotalCount:   100,
		ValidCount:   80,
		InvalidCount: 15,
		SkippedCount: 5,
	}

	if summary.TotalCount != 100 {
		t.Errorf("TotalCount = %d, want 100", summary.TotalCount)
	}
	if summary.ValidCount != 80 {
		t.Errorf("ValidCount = %d, want 80", summary.ValidCount)
	}
	if summary.InvalidCount != 15 {
		t.Errorf("InvalidCount = %d, want 15", summary.InvalidCount)
	}
	if summary.SkippedCount != 5 {
		t.Errorf("SkippedCount = %d, want 5", summary.SkippedCount)
	}
}

func TestBatchOperation_TimeFields(t *testing.T) {
	now := primitive.NewDateTimeFromTime(time.Now())
	op := &BatchOperation{
		Status:     BatchOpStatusRunning,
		StartedAt:  &now,
		CompletedAt: nil,
	}

	if op.CompletedAt != nil {
		t.Error("CompletedAt should be nil initially")
	}

	later := primitive.NewDateTimeFromTime(time.Now().Add(time.Hour))
	op.CompletedAt = &later
	op.Status = BatchOpStatusCompleted

	if op.CompletedAt == nil {
		t.Error("CompletedAt should be set after completion")
	}
}

func TestBatchOperationItem_TimeFields(t *testing.T) {
	now := primitive.NewDateTimeFromTime(time.Now())
	item := &BatchOperationItem{
		Status:      BatchItemStatusProcessing,
		StartedAt:   &now,
		CompletedAt: nil,
	}

	if item.Status != BatchItemStatusProcessing {
		t.Errorf("Status = %s, want processing", item.Status)
	}

	later := primitive.NewDateTimeFromTime(time.Now().Add(time.Minute))
	item.CompletedAt = &later
	item.Status = BatchItemStatusSucceeded

	if item.Status != BatchItemStatusSucceeded {
		t.Errorf("Status = %s, want succeeded", item.Status)
	}
}

func TestBatchOperation_GetSummary(t *testing.T) {
	items := []BatchOperationItem{
		{TargetID: "doc-1", Status: BatchItemStatusSucceeded},
		{TargetID: "doc-2", Status: BatchItemStatusSucceeded},
		{TargetID: "doc-3", Status: BatchItemStatusFailed},
	}

	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1", "doc-2", "doc-3"},
		Items:     items,
	}

	summary := op.GetSummary()

	if summary.SuccessCount != 2 {
		t.Errorf("SuccessCount = %d, want 2", summary.SuccessCount)
	}
	if summary.FailedCount != 1 {
		t.Errorf("FailedCount = %d, want 1", summary.FailedCount)
	}
}

func TestBatchOperation_UpdateItemStatus(t *testing.T) {
	items := []BatchOperationItem{
		{TargetID: "doc-1", Status: BatchItemStatusPending},
		{TargetID: "doc-2", Status: BatchItemStatusPending},
	}

	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1", "doc-2"},
		Items:     items,
	}

	// Update doc-1 to succeeded
	op.UpdateItemStatus("doc-1", BatchItemStatusSucceeded, "", "")

	if op.Items[0].Status != BatchItemStatusSucceeded {
		t.Errorf("Items[0].Status = %s, want succeeded", op.Items[0].Status)
	}

	// Update doc-2 to failed with error
	op.UpdateItemStatus("doc-2", BatchItemStatusFailed, "NOT_FOUND", "Document not found")

	if op.Items[1].Status != BatchItemStatusFailed {
		t.Errorf("Items[1].Status = %s, want failed", op.Items[1].Status)
	}
	if op.Items[1].ErrorCode != "NOT_FOUND" {
		t.Errorf("Items[1].ErrorCode = %s, want NOT_FOUND", op.Items[1].ErrorCode)
	}
}

func TestBatchOperation_WithPayload(t *testing.T) {
	payload := map[string]interface{}{
		"recursive":      true,
		"skip_conflicts": false,
	}

	op := &BatchOperation{
		ProjectID: primitive.NewObjectID(),
		Type:      BatchOpTypeDelete,
		TargetIDs: []string{"doc-1"},
		Payload:   payload,
		Atomic:    true,
	}

	if !op.Atomic {
		t.Error("Operation should be atomic")
	}

	if op.Payload == nil {
		t.Error("Payload should not be nil")
	}

	if op.Payload["recursive"] != true {
		t.Errorf("Payload[recursive] = %v, want true", op.Payload["recursive"])
	}
}

// ===== 错误处理和重试测试 =====

func TestBatchOperationItem_WithRetryable(t *testing.T) {
	item := &BatchOperationItem{
		TargetID:     "doc-1",
		TargetType:   "document",
		Status:       BatchItemStatusFailed,
		ErrorCode:    "CONFLICT",
		ErrorMessage: "Version conflict detected",
		Retryable:    true,
		RetryCount:   3,
	}

	if !item.Retryable {
		t.Error("Retryable should be true")
	}

	if item.ErrorCode != "CONFLICT" {
		t.Errorf("ErrorCode = %s, want CONFLICT", item.ErrorCode)
	}

	if item.ErrorMessage != "Version conflict detected" {
		t.Errorf("ErrorMessage = %s, want 'Version conflict detected'", item.ErrorMessage)
	}

	if item.RetryCount != 3 {
		t.Errorf("RetryCount = %d, want 3", item.RetryCount)
	}
}

func TestBatchOperationItem_ErrorFields_OmitEmpty(t *testing.T) {
	// 验证空值情况下omitempty标签的作用
	item := &BatchOperationItem{
		TargetID:  "doc-1",
		TargetType: "document",
		Status:    BatchItemStatusPending,
		Retryable: false, // 默认值false应该被序列化
	}

	// 空字符串不应该影响JSON序列化
	if item.ErrorCode != "" {
		t.Errorf("ErrorCode should be empty, got %s", item.ErrorCode)
	}

	if item.ErrorMessage != "" {
		t.Errorf("ErrorMessage should be empty, got %s", item.ErrorMessage)
	}

	// Retryable的默认值应该是false
	if item.Retryable {
		t.Error("Retryable should be false by default")
	}
}

func TestBatchOperationItem_WithFullErrorInfo(t *testing.T) {
	tests := []struct {
		name         string
		status       BatchItemStatus
		errorCode    string
		errorMessage string
		retryable    bool
	}{
		{
			name:         "conflict error retryable",
			status:       BatchItemStatusFailed,
			errorCode:    "VERSION_CONFLICT",
			errorMessage: "Document version 5 does not match expected version 3",
			retryable:    true,
		},
		{
			name:         "not found error not retryable",
			status:       BatchItemStatusFailed,
			errorCode:    "NOT_FOUND",
			errorMessage: "Document doc-1 not found",
			retryable:    false,
		},
		{
			name:         "permission error not retryable",
			status:       BatchItemStatusFailed,
			errorCode:    "PERMISSION_DENIED",
			errorMessage: "User does not have permission",
			retryable:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &BatchOperationItem{
				TargetID:     "doc-1",
				TargetType:   "document",
				Status:       tt.status,
				ErrorCode:    tt.errorCode,
				ErrorMessage: tt.errorMessage,
				Retryable:    tt.retryable,
			}

			if item.Status != tt.status {
				t.Errorf("Status = %s, want %s", item.Status, tt.status)
			}

			if item.ErrorCode != tt.errorCode {
				t.Errorf("ErrorCode = %s, want %s", item.ErrorCode, tt.errorCode)
			}

			if item.ErrorMessage != tt.errorMessage {
				t.Errorf("ErrorMessage = %s, want %s", item.ErrorMessage, tt.errorMessage)
			}

			if item.Retryable != tt.retryable {
				t.Errorf("Retryable = %v, want %v", item.Retryable, tt.retryable)
			}
		})
	}
}

func TestPreflightSummary_WithP1Extensions(t *testing.T) {
	summary := &PreflightSummary{
		TotalCount:   100,
		ValidCount:   85,
		InvalidCount: 10,
		SkippedCount: 5,

		// P1扩展字段
		SuccessCount: 70,
		FailedCount:  15,
	}

	// 验证原有字段
	if summary.TotalCount != 100 {
		t.Errorf("TotalCount = %d, want 100", summary.TotalCount)
	}

	if summary.ValidCount != 85 {
		t.Errorf("ValidCount = %d, want 85", summary.ValidCount)
	}

	if summary.InvalidCount != 10 {
		t.Errorf("InvalidCount = %d, want 10", summary.InvalidCount)
	}

	if summary.SkippedCount != 5 {
		t.Errorf("SkippedCount = %d, want 5", summary.SkippedCount)
	}

	// 验证P1扩展字段
	if summary.SuccessCount != 70 {
		t.Errorf("SuccessCount = %d, want 70", summary.SuccessCount)
	}

	if summary.FailedCount != 15 {
		t.Errorf("FailedCount = %d, want 15", summary.FailedCount)
	}

	// 验证一致性：SuccessCount + FailedCount 应该等于 ValidCount
	if summary.SuccessCount+summary.FailedCount != summary.ValidCount {
		t.Errorf("SuccessCount + FailedCount = %d, ValidCount = %d, should be equal",
			summary.SuccessCount+summary.FailedCount, summary.ValidCount)
	}
}

func TestPreflightSummary_NonAtomicModeResults(t *testing.T) {
	// 测试atomic=false模式下的结果统计
	// 场景：100个文档，80个有效，20个无效
	// 执行结果：70个成功，10个失败（部分失败场景）

	summary := &PreflightSummary{
		TotalCount:   100,
		ValidCount:   80,
		InvalidCount: 20,
		SkippedCount: 0,

		// 执行结果
		SuccessCount: 70,
		FailedCount:  10,
	}

	// 验证总数一致
	if summary.TotalCount != summary.ValidCount+summary.InvalidCount+summary.SkippedCount {
		t.Error("TotalCount should equal ValidCount + InvalidCount + SkippedCount")
	}

	// 验证执行结果统计
	executedCount := summary.SuccessCount + summary.FailedCount
	if executedCount != summary.ValidCount {
		t.Errorf("Executed count (%d) should equal ValidCount (%d)", executedCount, summary.ValidCount)
	}

	// 验证有部分失败
	if summary.FailedCount == 0 {
		t.Error("Expected some failures in non-atomic mode")
	}

	// 验证不是全部失败
	if summary.SuccessCount == 0 {
		t.Error("Expected some successes in non-atomic mode")
	}
}

func TestPreflightSummary_AtomicModeResults(t *testing.T) {
	// 测试atomic=true模式下的结果统计
	// 场景：第一个操作失败，整个批次中止

	summary := &PreflightSummary{
		TotalCount:   100,
		ValidCount:   100,
		InvalidCount: 0,
		SkippedCount: 0,

		// 原子模式下，要么全部成功，要么全部失败
		SuccessCount: 0,  // 全部失败
		FailedCount:  100, // 因为第一个失败导致整个批次失败
	}

	// 验证原子模式的一致性
	if summary.SuccessCount > 0 && summary.FailedCount > 0 {
		t.Error("Atomic mode should not have partial success/failure")
	}

	// 验证总执行数
	executedCount := summary.SuccessCount + summary.FailedCount
	if executedCount != summary.ValidCount {
		t.Errorf("Executed count (%d) should equal ValidCount (%d)", executedCount, summary.ValidCount)
	}
}

func TestBatchOperationItem_AllStatusesWithErrors(t *testing.T) {
	// 验证所有状态类型与新错误字段的兼容性
	statuses := []BatchItemStatus{
		BatchItemStatusPending,
		BatchItemStatusProcessing,
		BatchItemStatusSucceeded,
		BatchItemStatusFailed,
		BatchItemStatusSkipped,
		BatchItemStatusCancelled,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			item := &BatchOperationItem{
				TargetID:     "doc-1",
				TargetType:   "document",
				Status:       status,
				ErrorCode:    "TEST_ERROR",
				ErrorMessage: "Test error message",
				Retryable:    true,
			}

			if item.Status != status {
				t.Errorf("Status = %s, want %s", item.Status, status)
			}

			// 所有状态都应该能设置错误信息
			if item.ErrorCode != "TEST_ERROR" {
				t.Errorf("ErrorCode should be set for status %s", status)
			}

			if item.ErrorMessage != "Test error message" {
				t.Errorf("ErrorMessage should be set for status %s", status)
			}

			if !item.Retryable {
				t.Errorf("Retryable should be set for status %s", status)
			}
		})
	}
}
