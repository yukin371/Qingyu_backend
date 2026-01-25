package document

import (
	"errors"
	"testing"
	"time"

	"Qingyu_backend/models/writer"
	pkgErrors "Qingyu_backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestRetryConfig 测试重试配置的转换
func TestRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("期望MaxRetries=3，实际: %d", config.MaxRetries)
	}
	if config.RetryDelay != 1000 {
		t.Errorf("期望RetryDelay=1000，实际: %d", config.RetryDelay)
	}
	if len(config.RetryableErrors) == 0 {
		t.Error("RetryableErrors不应为空")
	}

	t.Logf("默认配置: MaxRetries=%d, RetryDelay=%d, RetryableErrors=%v",
		config.MaxRetries, config.RetryDelay, config.RetryableErrors)
}

// TestRetryService_ShouldRetry 测试ShouldRetry方法
func TestRetryService_ShouldRetry(t *testing.T) {
	service := NewRetryService()
	config := DefaultRetryConfig()

	// 测试可重试的错误
	retryableErr := &pkgErrors.ServiceError{
		Type:    pkgErrors.ServiceErrorType("VERSION_CONFLICT"),
		Message: "version conflict",
	}
	if !service.ShouldRetry(retryableErr, config) {
		t.Error("VERSION_CONFLICT应该可重试")
	}

	// 测试不可重试的错误
	nonRetryableErr := errors.New("invalid argument")
	if service.ShouldRetry(nonRetryableErr, config) {
		t.Error("普通错误不应该可重试")
	}

	t.Log("ShouldRetry测试通过")
}

// TestRetryService_GetRetryDelay 测试GetRetryDelay方法
func TestRetryService_GetRetryDelay(t *testing.T) {
	service := NewRetryService()
	config := &RetryConfig{
		MaxRetries: 3,
		RetryDelay: 100, // 100ms
	}

	// 测试指数退避
	delay0 := service.GetRetryDelay(0, config)
	delay1 := service.GetRetryDelay(1, config)
	delay2 := service.GetRetryDelay(2, config)

	expected0 := 100 * time.Millisecond
	expected1 := 200 * time.Millisecond
	expected2 := 400 * time.Millisecond

	if delay0 != expected0 {
		t.Errorf("attempt 0: 期望%v，实际%v", expected0, delay0)
	}
	if delay1 != expected1 {
		t.Errorf("attempt 1: 期望%v，实际%v", expected1, delay1)
	}
	if delay2 != expected2 {
		t.Errorf("attempt 2: 期望%v，实际%v", expected2, delay2)
	}

	t.Logf("指数退避测试通过: %v, %v, %v", delay0, delay1, delay2)
}

// TestRetryService_CanRetrySimple 测试CanRetry方法
func TestRetryService_CanRetrySimple(t *testing.T) {
	service := NewRetryService()
	config := &RetryConfig{
		MaxRetries: 3,
	}

	if !service.CanRetry(0, config) {
		t.Error("attempt 0应该可以重试")
	}
	if !service.CanRetry(1, config) {
		t.Error("attempt 1应该可以重试")
	}
	if !service.CanRetry(2, config) {
		t.Error("attempt 2应该可以重试")
	}
	if service.CanRetry(3, config) {
		t.Error("attempt 3不应该可以重试")
	}

	t.Log("CanRetry测试通过")
}

// TestBatchOperation_GetSummary 测试GetSummary方法
func TestBatchOperation_GetSummary(t *testing.T) {
	projectID := primitive.NewObjectID()
	batchOp := &writer.BatchOperation{
		ProjectID: projectID,
		Type:      writer.BatchOpTypeDelete,
		TargetIDs: []string{"doc1", "doc2", "doc3", "doc4", "doc5"},
		Status:    writer.BatchOpStatusProcessing,
		Atomic:    false,
		PreflightSummary: &writer.PreflightSummary{
			TotalCount:  5,
			ValidCount:  5,
			InvalidCount: 0,
		},
		Items: []writer.BatchOperationItem{
			{TargetID: "doc1", Status: writer.BatchItemStatusSucceeded},
			{TargetID: "doc2", Status: writer.BatchItemStatusSucceeded},
			{TargetID: "doc3", Status: writer.BatchItemStatusFailed, ErrorCode: "DELETE_FAILED", ErrorMsg: "not found"},
			{TargetID: "doc4", Status: writer.BatchItemStatusSucceeded},
			{TargetID: "doc5", Status: writer.BatchItemStatusFailed, ErrorCode: "DELETE_FAILED", ErrorMsg: "permission denied"},
		},
	}
	batchOp.TouchForCreate()

	summary := batchOp.GetSummary()

	if summary.TotalCount != 5 {
		t.Errorf("期望TotalCount=5，实际: %d", summary.TotalCount)
	}
	if summary.SuccessCount != 3 {
		t.Errorf("期望SuccessCount=3，实际: %d", summary.SuccessCount)
	}
	if summary.FailedCount != 2 {
		t.Errorf("期望FailedCount=2，实际: %d", summary.FailedCount)
	}

	t.Logf("GetSummary测试通过: Total=%d, Success=%d, Failed=%d",
		summary.TotalCount, summary.SuccessCount, summary.FailedCount)
}

// TestBatchOperation_IsCompleted 测试IsCompleted方法
func TestBatchOperation_IsCompleted(t *testing.T) {
	projectID := primitive.NewObjectID()

	tests := []struct {
		name     string
		status   writer.BatchOperationStatus
		expected bool
	}{
		{"已完成", writer.BatchOpStatusCompleted, true},
		{"失败", writer.BatchOpStatusFailed, true},
		{"已取消", writer.BatchOpStatusCancelled, true},
		{"部分成功", writer.BatchOpStatusPartial, true},
		{"待处理", writer.BatchOpStatusPending, false},
		{"处理中", writer.BatchOpStatusProcessing, false},
		{"预检查", writer.BatchOpStatusPreflight, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchOp := &writer.BatchOperation{
				ProjectID: projectID,
				Type:      writer.BatchOpTypeDelete,
				TargetIDs: []string{"doc1"},
				Status:    tt.status,
			}
			batchOp.TouchForCreate()

			result := batchOp.IsCompleted()
			if result != tt.expected {
				t.Errorf("IsCompleted()期望%v，实际%v", tt.expected, result)
			}
		})
	}

	t.Log("IsCompleted测试通过")
}

// TestConvertMapToRetryConfig 测试convertMapToRetryConfig方法
func TestConvertMapToRetryConfig(t *testing.T) {
	service := &BatchOperationService{}

	// 测试nil输入
	config := service.convertMapToRetryConfig(nil)
	if config != nil {
		t.Error("nil输入应该返回nil")
	}

	// 测试正常输入
	data := map[string]interface{}{
		"maxRetries":      5,
		"retryDelay":      2000,
		"retryableErrors": []string{"ERROR1", "ERROR2"},
	}
	config = service.convertMapToRetryConfig(data)

	if config == nil {
		t.Fatal("convertMapToRetryConfig不应该返回nil")
	}
	if config.MaxRetries != 5 {
		t.Errorf("期望MaxRetries=5，实际: %d", config.MaxRetries)
	}
	if config.RetryDelay != 2000 {
		t.Errorf("期望RetryDelay=2000，实际: %d", config.RetryDelay)
	}
	if len(config.RetryableErrors) != 2 {
		t.Errorf("期望RetryableErrors长度=2，实际: %d", len(config.RetryableErrors))
	}

	t.Log("convertMapToRetryConfig测试通过")
}

// TestConvertRetryConfigToMap 测试convertRetryConfigToMap方法
func TestConvertRetryConfigToMap(t *testing.T) {
	service := &BatchOperationService{}

	config := &RetryConfig{
		MaxRetries:      3,
		RetryDelay:      1000,
		RetryableErrors: []string{"ERROR1", "ERROR2"},
	}

	data := service.convertRetryConfigToMap(config)

	if data == nil {
		t.Fatal("convertRetryConfigToMap不应该返回nil")
	}
	if data["maxRetries"].(int) != 3 {
		t.Errorf("期望maxRetries=3，实际: %v", data["maxRetries"])
	}
	if data["retryDelay"].(int) != 1000 {
		t.Errorf("期望retryDelay=1000，实际: %v", data["retryDelay"])
	}
	if len(data["retryableErrors"].([]string)) != 2 {
		t.Errorf("期望retryableErrors长度=2，实际: %v", data["retryableErrors"])
	}

	// 测试nil输入
	data = service.convertRetryConfigToMap(nil)
	if data != nil {
		t.Error("nil输入应该返回nil")
	}

	t.Log("convertRetryConfigToMap测试通过")
}

// TestBatchOperationItem_RetryableFields 测试BatchOperationItem的可重试字段
func TestBatchOperationItem_RetryableFields(t *testing.T) {
	item := writer.BatchOperationItem{
		TargetID:   "doc1",
		Status:     writer.BatchItemStatusPending,
		Retryable:  true,
		RetryCount: 0,
	}

	// 模拟第一次失败
	item.Status = writer.BatchItemStatusFailed
	item.ErrorCode = "VERSION_CONFLICT"
	item.ErrorMsg = "version conflict"
	item.Retryable = true
	item.RetryCount = 1

	if item.Retryable != true {
		t.Error("VERSION_CONFLICT应该可重试")
	}

	// 模拟第二次失败（不可重试）
	item.ErrorCode = "PERMISSION_DENIED"
	item.ErrorMsg = "permission denied"
	item.Retryable = false
	item.RetryCount = 2

	if item.Retryable != false {
		t.Error("PERMISSION_DENIED不应该可重试")
	}

	t.Log("BatchOperationItem可重试字段测试通过")
}
