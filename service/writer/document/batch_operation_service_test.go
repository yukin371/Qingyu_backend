package document

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/writer"
	mongodbwriter "Qingyu_backend/repository/mongodb/writer"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestBatchOperationService_Submit 测试提交批量操作
func TestBatchOperationService_Submit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// 创建repositories
	batchOpRepo := mongodbwriter.NewBatchOperationRepository(testDB)
	opLogRepo := mongodbwriter.NewOperationLogRepository(testDB)
	docRepo := mongodbwriter.NewMongoDocumentRepository(testDB)

	// 创建service
	service := NewBatchOperationService(batchOpRepo, opLogRepo, docRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建测试文档
	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "Test Doc 1",
		StableRef: primitive.NewObjectID().Hex(),
		OrderKey:  writer.DefaultOrderKey,
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.TouchForCreate()

	err := docRepo.Create(ctx, doc1)
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	// 测试提交批量操作
	req := &SubmitBatchOperationRequest{
		ProjectID:       projectID,
		Type:            writer.BatchOpTypeDelete,
		TargetIDs:       []string{doc1.ID.Hex()},
		Atomic:          true,
		ConflictPolicy:  writer.ConflictPolicyAbort,
		UserID:          userID,
		ClientRequestID: "test-client-request-1",
	}

	batchOp, err := service.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	if batchOp == nil {
		t.Fatal("batchOp is nil")
	}

	if batchOp.Status != writer.BatchOpStatusPending {
		t.Errorf("Expected status pending, got %s", batchOp.Status)
	}

	if batchOp.ExecutionMode != writer.ExecutionModeStandardAtomic {
		t.Errorf("Expected execution_mode standard_atomic, got %s", batchOp.ExecutionMode)
	}

	if batchOp.PreflightSummary == nil {
		t.Error("PreflightSummary is nil")
	} else if batchOp.PreflightSummary.ValidCount != 1 {
		t.Errorf("Expected 1 valid document, got %d", batchOp.PreflightSummary.ValidCount)
	}
}

// TestBatchOperationService_Idempotency 测试幂等性
func TestBatchOperationService_Idempotency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	batchOpRepo := mongodbwriter.NewBatchOperationRepository(testDB)
	opLogRepo := mongodbwriter.NewOperationLogRepository(testDB)
	docRepo := mongodbwriter.NewMongoDocumentRepository(testDB)

	service := NewBatchOperationService(batchOpRepo, opLogRepo, docRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建测试文档
	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "Test Doc 1",
		StableRef: primitive.NewObjectID().Hex(),
		OrderKey:  writer.DefaultOrderKey,
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.TouchForCreate()

	err := docRepo.Create(ctx, doc1)
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	clientRequestID := "test-client-request-2"

	// 第一次提交
	req1 := &SubmitBatchOperationRequest{
		ProjectID:       projectID,
		Type:            writer.BatchOpTypeDelete,
		TargetIDs:       []string{doc1.ID.Hex()},
		Atomic:          true,
		ConflictPolicy:  writer.ConflictPolicyAbort,
		UserID:          userID,
		ClientRequestID: clientRequestID,
	}

	batchOp1, err := service.Submit(ctx, req1)
	if err != nil {
		t.Fatalf("First Submit failed: %v", err)
	}

	// 第二次提交（应该返回相同的操作）
	req2 := &SubmitBatchOperationRequest{
		ProjectID:       projectID,
		Type:            writer.BatchOpTypeDelete,
		TargetIDs:       []string{doc1.ID.Hex()},
		Atomic:          true,
		ConflictPolicy:  writer.ConflictPolicyAbort,
		UserID:          userID,
		ClientRequestID: clientRequestID,
	}

	batchOp2, err := service.Submit(ctx, req2)
	if err != nil {
		t.Fatalf("Second Submit failed: %v", err)
	}

	if batchOp1.ID != batchOp2.ID {
		t.Errorf("Expected same batch operation ID, got %s and %s", batchOp1.ID.Hex(), batchOp2.ID.Hex())
	}
}

// TestBatchOperationService_ExecuteStandardAtomic 测试标准原子执行
func TestBatchOperationService_ExecuteStandardAtomic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	batchOpRepo := mongodbwriter.NewBatchOperationRepository(testDB)
	opLogRepo := mongodbwriter.NewOperationLogRepository(testDB)
	docRepo := mongodbwriter.NewMongoDocumentRepository(testDB)

	service := NewBatchOperationService(batchOpRepo, opLogRepo, docRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建测试文档
	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "Test Doc 1",
		StableRef: primitive.NewObjectID().Hex(),
		OrderKey:  writer.DefaultOrderKey,
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.TouchForCreate()

	err := docRepo.Create(ctx, doc1)
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	// 提交批量操作
	req := &SubmitBatchOperationRequest{
		ProjectID:      projectID,
		Type:           writer.BatchOpTypeDelete,
		TargetIDs:      []string{doc1.ID.Hex()},
		Atomic:         true,
		ConflictPolicy: writer.ConflictPolicyAbort,
		UserID:         userID,
	}

	batchOp, err := service.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// 执行批量操作
	err = service.Execute(ctx, batchOp.ID)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// 验证文档已软删除
	deletedDoc, err := docRepo.GetByID(ctx, doc1.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get deleted document: %v", err)
	}

	// DeletedAt是time.Time类型，检查是否为零值来判断是否软删除
	if deletedDoc.DeletedAt.IsZero() {
		t.Error("Document should be soft deleted (DeletedAt should not be zero)")
	}

	// 验证操作日志
	logs, err := opLogRepo.GetByChainID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get operation logs: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 operation log, got %d", len(logs))
	}

	if logs[0].CommandType != writer.CommandDelete {
		t.Errorf("Expected command type delete, got %s", logs[0].CommandType)
	}
}

// TestBatchOperationService_GetProgress 测试获取进度
func TestBatchOperationService_GetProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	batchOpRepo := mongodbwriter.NewBatchOperationRepository(testDB)
	opLogRepo := mongodbwriter.NewOperationLogRepository(testDB)
	docRepo := mongodbwriter.NewMongoDocumentRepository(testDB)

	service := NewBatchOperationService(batchOpRepo, opLogRepo, docRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建测试文档
	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "Test Doc 1",
		StableRef: primitive.NewObjectID().Hex(),
		OrderKey:  writer.DefaultOrderKey,
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.TouchForCreate()

	err := docRepo.Create(ctx, doc1)
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	// 提交并执行批量操作
	req := &SubmitBatchOperationRequest{
		ProjectID:      projectID,
		Type:           writer.BatchOpTypeDelete,
		TargetIDs:      []string{doc1.ID.Hex()},
		Atomic:         true,
		ConflictPolicy: writer.ConflictPolicyAbort,
		UserID:         userID,
	}

	batchOp, err := service.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// 获取进度（pending状态）
	progress, err := service.GetProgress(ctx, batchOp.ID)
	if err != nil {
		t.Fatalf("GetProgress failed: %v", err)
	}

	if progress.TotalItems != 1 {
		t.Errorf("Expected 1 total item, got %d", progress.TotalItems)
	}

	if progress.Status != writer.BatchOpStatusPending {
		t.Errorf("Expected status pending, got %s", progress.Status)
	}

	// 执行操作
	err = service.Execute(ctx, batchOp.ID)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// 获取进度（completed状态）
	progress, err = service.GetProgress(ctx, batchOp.ID)
	if err != nil {
		t.Fatalf("GetProgress after execute failed: %v", err)
	}

	if progress.Status != writer.BatchOpStatusCompleted {
		t.Errorf("Expected status completed, got %s", progress.Status)
	}

	if progress.CompletedItems != 1 {
		t.Errorf("Expected 1 completed item, got %d", progress.CompletedItems)
	}

	if progress.FinishedAt == nil {
		t.Error("FinishedAt should not be nil")
	}
}

// TestBatchOperationService_Cancel 测试取消操作
func TestBatchOperationService_Cancel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	batchOpRepo := mongodbwriter.NewBatchOperationRepository(testDB)
	opLogRepo := mongodbwriter.NewOperationLogRepository(testDB)
	docRepo := mongodbwriter.NewMongoDocumentRepository(testDB)

	service := NewBatchOperationService(batchOpRepo, opLogRepo, docRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建测试文档
	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "Test Doc 1",
		StableRef: primitive.NewObjectID().Hex(),
		OrderKey:  writer.DefaultOrderKey,
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.TouchForCreate()

	err := docRepo.Create(ctx, doc1)
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	// 提交批量操作
	req := &SubmitBatchOperationRequest{
		ProjectID:      projectID,
		Type:           writer.BatchOpTypeDelete,
		TargetIDs:      []string{doc1.ID.Hex()},
		Atomic:         true,
		ConflictPolicy: writer.ConflictPolicyAbort,
		UserID:         userID,
	}

	batchOp, err := service.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// 尝试取消pending状态的操作（应该失败，因为不是running状态）
	err = service.Cancel(ctx, batchOp.ID, userID)
	if err == nil {
		t.Error("Expected error when cancelling non-running operation")
	}

	if err != ErrBatchOperationNotRunning {
		t.Errorf("Expected ErrBatchOperationNotRunning, got %v", err)
	}
}

// TestBatchOperationService_ExecutionModeSelection 测试执行模式选择
func TestBatchOperationService_ExecutionModeSelection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	batchOpRepo := mongodbwriter.NewBatchOperationRepository(testDB)
	opLogRepo := mongodbwriter.NewOperationLogRepository(testDB)
	docRepo := mongodbwriter.NewMongoDocumentRepository(testDB)

	service := NewBatchOperationService(batchOpRepo, opLogRepo, docRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建少量文档（<200）- 应该选择standard_atomic
	var smallTargetIDs []string
	for i := 0; i < 5; i++ {
		doc := &writer.Document{
			ProjectID: projectID,
			Title:     "Small Doc",
			StableRef: primitive.NewObjectID().Hex(),
			OrderKey:  writer.DefaultOrderKey,
			Type:      writer.TypeChapter,
			Level:     0,
		}
		doc.TouchForCreate()
		err := docRepo.Create(ctx, doc)
		if err != nil {
			t.Fatalf("Failed to create test document: %v", err)
		}
		smallTargetIDs = append(smallTargetIDs, doc.ID.Hex())
	}

	req := &SubmitBatchOperationRequest{
		ProjectID:      projectID,
		Type:           writer.BatchOpTypeDelete,
		TargetIDs:      smallTargetIDs,
		Atomic:         true,
		ConflictPolicy: writer.ConflictPolicyAbort,
		UserID:         userID,
	}

	batchOp, err := service.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	if batchOp.ExecutionMode != writer.ExecutionModeStandardAtomic {
		t.Errorf("Expected execution_mode standard_atomic for small batch, got %s", batchOp.ExecutionMode)
	}
}

// TestBatchOperationService_InvalidDocuments 测试处理无效文档
func TestBatchOperationService_InvalidDocuments(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	batchOpRepo := mongodbwriter.NewBatchOperationRepository(testDB)
	opLogRepo := mongodbwriter.NewOperationLogRepository(testDB)
	docRepo := mongodbwriter.NewMongoDocumentRepository(testDB)

	service := NewBatchOperationService(batchOpRepo, opLogRepo, docRepo)

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建一个有效文档
	doc1 := &writer.Document{
		ProjectID: projectID,
		Title:     "Valid Doc",
		StableRef: primitive.NewObjectID().Hex(),
		OrderKey:  writer.DefaultOrderKey,
		Type:      writer.TypeChapter,
		Level:     0,
	}
	doc1.TouchForCreate()
	err := docRepo.Create(ctx, doc1)
	if err != nil {
		t.Fatalf("Failed to create test document: %v", err)
	}

	// 包含有效和无效的文档ID
	req := &SubmitBatchOperationRequest{
		ProjectID:      projectID,
		Type:           writer.BatchOpTypeDelete,
		TargetIDs:      []string{doc1.ID.Hex(), "invalid-id", primitive.NewObjectID().Hex()},
		Atomic:         false, // 非原子模式，允许部分失败
		ConflictPolicy: writer.ConflictPolicySkip,
		UserID:         userID,
	}

	batchOp, err := service.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// 应该只有1个有效文档
	if len(batchOp.TargetIDs) != 1 {
		t.Errorf("Expected 1 valid target ID, got %d", len(batchOp.TargetIDs))
	}

	if batchOp.PreflightSummary.ValidCount != 1 {
		t.Errorf("Expected 1 valid document in summary, got %d", batchOp.PreflightSummary.ValidCount)
	}

	if batchOp.PreflightSummary.InvalidCount != 2 {
		t.Errorf("Expected 2 invalid documents in summary, got %d", batchOp.PreflightSummary.InvalidCount)
	}
}

// setupTestDB 是测试辅助函数（需要在测试文件中定义）
var testDB *mongo.Database // 这里需要根据实际的测试设置进行初始化

// TestMain 测试主函数
func TestMain(m *testing.M) {
	// TODO: 初始化测试数据库
	// 例如：
	// client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	// testDB = client.Database("qingyu_test")
	//
	// m.Run()
	//
	// client.Disconnect(context.Background())
}

// sleep 辅助函数
func sleep(duration time.Duration) {
	time.Sleep(duration)
}
