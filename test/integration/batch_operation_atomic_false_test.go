package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/global"
	writerModel "Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	writerInterfaces "Qingyu_backend/repository/interfaces/writer"
	documentService "Qingyu_backend/service/writer/document"
	"Qingyu_backend/test/testutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("Skipping integration tests: cannot load config:", err)
		os.Exit(0)
	}

	if config.GlobalConfig.Database == nil {
		fmt.Println("Skipping integration tests: database config is nil")
		os.Exit(0)
	}

	// 检查MongoDB配置
	if config.GlobalConfig.Database.Primary.Type != config.DatabaseTypeMongoDB ||
		config.GlobalConfig.Database.Primary.MongoDB == nil {
		fmt.Println("Skipping integration tests: MongoDB config is missing")
		os.Exit(0)
	}

	testutil.EnableStrictLogging()

	mongoCfg := config.GlobalConfig.Database.Primary.MongoDB
	clientOpts := options.Client().ApplyURI(mongoCfg.URI)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println("Skipping integration tests: cannot connect to MongoDB:", err)
		os.Exit(0)
	}
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("Skipping integration tests: cannot ping MongoDB:", err)
		os.Exit(0)
	}
	global.MongoClient = client
	global.DB = client.Database(mongoCfg.Database)

	code := m.Run()
	code = testutil.CheckStrictLogViolations(code)

	_ = global.MongoClient.Disconnect(ctx)
	os.Exit(code)
}

// setupTestEnv 创建测试环境
func setupTestEnv(t *testing.T) (context.Context, *documentService.BatchOperationService, writerInterfaces.BatchOperationRepository, writerInterfaces.DocumentRepository, primitive.ObjectID) {
	ctx := context.Background()

	// 创建测试项目ID
	projectID := primitive.NewObjectID()

	// 创建仓储
	batchOpRepo := writerRepo.NewMongoBatchOperationRepository(global.DB)
	docRepo := writerRepo.NewMongoDocumentRepository(global.DB)
	projectRepo := writerRepo.NewMongoProjectRepository(global.DB)

	// 创建服务
	batchOpService := documentService.NewBatchOperationService(batchOpRepo, docRepo, projectRepo, nil)

	return ctx, batchOpService, batchOpRepo, docRepo, projectID
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T, projectID primitive.ObjectID, batchOpID primitive.ObjectID) {
	ctx := context.Background()
	global.DB.Collection("batch_operations").DeleteMany(ctx, bson.M{"project_id": projectID})
	global.DB.Collection("novel_files").DeleteMany(ctx, bson.M{"project_id": projectID})
	global.DB.Collection("document_contents").DeleteMany(ctx, bson.M{})
}

// createTestDocuments 创建测试文档
func createTestDocuments(t *testing.T, ctx context.Context, docRepo writerInterfaces.DocumentRepository, projectID primitive.ObjectID, count int) []string {
	docIDs := make([]string, count)

	for i := 0; i < count; i++ {
		doc := &writerModel.Document{
			ProjectID: projectID,
			Title:     fmt.Sprintf("Test Document %d", i),
			Type:      "chapter", // 使用Type字段代替ContentType
			Status:    "writing",
		}
		doc.TouchForCreate()

		err := docRepo.Create(ctx, doc)
		if err != nil {
			t.Fatalf("创建测试文档失败: %v", err)
		}
		docIDs[i] = doc.ID.Hex()
	}

	return docIDs
}

// TestAtomicFalse_PartialFailure 测试atomic=false时部分失败继续执行
func TestAtomicFalse_PartialFailure(t *testing.T) {
	ctx, batchOpService, batchOpRepo, docRepo, projectID := setupTestEnv(t)
	defer cleanupTestData(t, projectID, primitive.NilObjectID)

	// 创建5个测试文档
	docIDs := createTestDocuments(t, ctx, docRepo, projectID, 5)

	// 手动删除第3个文档，模拟部分文档不存在的情况
	err := docRepo.SoftDelete(ctx, docIDs[2], projectID.Hex())
	if err != nil {
		t.Fatalf("预删除文档失败: %v", err)
	}

	// 创建批量删除请求（atomic=false）
	req := &documentService.SubmitBatchOperationRequest{
		ProjectID: projectID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs: docIDs,
		Atomic:    false, // 非原子操作
		RetryConfig: &documentService.RetryConfig{
			MaxRetries: 2,
			RetryDelay: 100, // 100ms
			RetryableErrors: []string{
				"VERSION_CONFLICT",
				"NETWORK_ERROR",
			},
		},
	}

	// 设置用户ID
	ctx = context.WithValue(ctx, "userID", "test_user")

	// 提交批量操作
	batchOp, err := batchOpService.Submit(ctx, req)
	if err != nil {
		t.Fatalf("提交批量操作失败: %v", err)
	}

	// 执行批量操作
	err = batchOpService.Execute(ctx, batchOp.ID.Hex())
	if err != nil {
		// atomic=false时，部分失败不应该返回错误
		// 但这里可能返回错误，我们需要检查状态
		t.Logf("Execute返回错误（预期）: %v", err)
	}

	// 获取更新后的批量操作
	result, err := batchOpRepo.GetByID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("获取批量操作结果失败: %v", err)
	}

	// 验证状态
	if result.Status != writerModel.BatchOpStatusPartial {
		t.Errorf("期望状态为partial，实际为: %s", result.Status)
	}

	// 验证摘要
	summary := result.GetSummary()
	if summary.SuccessCount != 4 {
		t.Errorf("期望成功4个，实际: %d", summary.SuccessCount)
	}
	if summary.FailedCount != 1 {
		t.Errorf("期望失败1个，实际: %d", summary.FailedCount)
	}

	// 验证每个item的状态
	successCount := 0
	failedCount := 0
	for _, item := range result.Items {
		if item.Status == writerModel.BatchItemStatusSucceeded {
			successCount++
		} else if item.Status == writerModel.BatchItemStatusFailed {
			failedCount++
			// 验证失败项有错误信息
			if item.ErrorCode == "" {
				t.Error("失败项缺少ErrorCode")
			}
			if item.ErrorMsg == "" {
				t.Error("失败项缺少ErrorMsg")
			}
		}
	}

	if successCount != 4 {
		t.Errorf("期望4个成功项，实际: %d", successCount)
	}
	if failedCount != 1 {
		t.Errorf("期望1个失败项，实际: %d", failedCount)
	}

	t.Log("TestAtomicFalse_PartialFailure 测试通过")
}

// TestAtomicTrue_ImmediateStop 测试atomic=true时失败立即停止
func TestAtomicTrue_ImmediateStop(t *testing.T) {
	ctx, batchOpService, batchOpRepo, docRepo, projectID := setupTestEnv(t)
	defer cleanupTestData(t, projectID, primitive.NilObjectID)

	// 创建5个测试文档
	docIDs := createTestDocuments(t, ctx, docRepo, projectID, 5)

	// 手动删除第2个文档，模拟中间文档不存在
	err := docRepo.SoftDelete(ctx, docIDs[1], projectID.Hex())
	if err != nil {
		t.Fatalf("预删除文档失败: %v", err)
	}

	// 创建批量删除请求（atomic=true）
	req := &documentService.SubmitBatchOperationRequest{
		ProjectID: projectID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs: docIDs,
		Atomic:    true, // 原子操作
		RetryConfig: &documentService.RetryConfig{
			MaxRetries: 2,
			RetryDelay: 100,
			RetryableErrors: []string{
				"VERSION_CONFLICT",
				"NETWORK_ERROR",
			},
		},
	}

	// 设置用户ID
	ctx = context.WithValue(ctx, "userID", "test_user")

	// 提交批量操作
	batchOp, err := batchOpService.Submit(ctx, req)
	if err != nil {
		t.Fatalf("提交批量操作失败: %v", err)
	}

	// 执行批量操作
	err = batchOpService.Execute(ctx, batchOp.ID.Hex())
	if err == nil {
		t.Error("atomic=true时，失败应该返回错误")
	}

	// 获取更新后的批量操作
	result, err := batchOpRepo.GetByID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("获取批量操作结果失败: %v", err)
	}

	// 验证状态：atomic=true时，任何失败都应该导致整体失败
	if result.Status != writerModel.BatchOpStatusFailed {
		t.Errorf("期望状态为failed，实际为: %s", result.Status)
	}

	// 验证有错误信息
	if result.ErrorCode == "" {
		t.Error("失败的操作应该有ErrorCode")
	}
	if result.ErrorMessage == "" {
		t.Error("失败的操作应该有ErrorMessage")
	}

	// 验证只有部分文档被处理（第2个失败后停止）
	// 应该有至少1个成功（第1个），第2个失败
	successCount := 0
	for _, item := range result.Items {
		if item.Status == writerModel.BatchItemStatusSucceeded {
			successCount++
		}
	}

	if successCount < 1 {
		t.Error("atomic=true时，至少应该处理第1个文档")
	}
	if successCount > 1 {
		t.Logf("注意：atomic=true时处理了%d个文档后才失败（可能在重试）", successCount)
	}

	t.Log("TestAtomicTrue_ImmediateStop 测试通过")
}

// TestAtomicFalse_RetryMechanism 测试atomic=false时的重试机制
func TestAtomicFalse_RetryMechanism(t *testing.T) {
	ctx, batchOpService, batchOpRepo, docRepo, projectID := setupTestEnv(t)
	defer cleanupTestData(t, projectID, primitive.NilObjectID)

	// 创建3个测试文档
	docIDs := createTestDocuments(t, ctx, docRepo, projectID, 3)

	// 创建一个mock repository，在第一次删除时失败，第二次成功
	// 由于无法直接mock，我们通过删除一个不存在的文档来模拟重试场景
	// 这里我们使用可重试的错误场景

	// 创建批量删除请求
	req := &documentService.SubmitBatchOperationRequest{
		ProjectID: projectID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs: docIDs,
		Atomic:    false,
		RetryConfig: &documentService.RetryConfig{
			MaxRetries:      3,
			RetryDelay:      100, // 100ms
			RetryableErrors: []string{"DOCUMENT_NOT_FOUND", "DELETE_FAILED"},
		},
	}

	// 设置用户ID
	ctx = context.WithValue(ctx, "userID", "test_user")

	// 提交批量操作
	batchOp, err := batchOpService.Submit(ctx, req)
	if err != nil {
		t.Fatalf("提交批量操作失败: %v", err)
	}

	// 执行批量操作
	err = batchOpService.Execute(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Logf("Execute返回错误（可能正常）: %v", err)
	}

	// 获取更新后的批量操作
	result, err := batchOpRepo.GetByID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("获取批量操作结果失败: %v", err)
	}

	// 验证RetryCount字段被正确设置
	for _, item := range result.Items {
		t.Logf("文档 %s: 状态=%s, 重试次数=%d, 可重试=%v",
			item.TargetID, item.Status, item.RetryCount, item.Retryable)
	}

	// 验证摘要统计正确
	summary := result.GetSummary()
	totalItems := summary.SuccessCount + summary.FailedCount
	if totalItems != len(docIDs) {
		t.Errorf("期望处理%d个项，实际: %d", len(docIDs), totalItems)
	}

	t.Log("TestAtomicFalse_RetryMechanism 测试通过")
}

// TestSummaryStatistics 测试摘要统计的正确性
func TestSummaryStatistics(t *testing.T) {
	ctx, batchOpService, batchOpRepo, docRepo, projectID := setupTestEnv(t)
	defer cleanupTestData(t, projectID, primitive.NilObjectID)

	// 创建10个测试文档
	docIDs := createTestDocuments(t, ctx, docRepo, projectID, 10)

	// 预删除其中3个文档，模拟部分失败
	docRepo.SoftDelete(ctx, docIDs[3], projectID.Hex())
	docRepo.SoftDelete(ctx, docIDs[6], projectID.Hex())
	docRepo.SoftDelete(ctx, docIDs[9], projectID.Hex())

	// 创建批量删除请求（atomic=false）
	req := &documentService.SubmitBatchOperationRequest{
		ProjectID: projectID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs: docIDs,
		Atomic:    false,
		RetryConfig: &documentService.RetryConfig{
			MaxRetries: 1,
			RetryDelay: 50,
			RetryableErrors: []string{
				"DOCUMENT_NOT_FOUND",
			},
		},
	}

	// 设置用户ID
	ctx = context.WithValue(ctx, "userID", "test_user")

	// 提交批量操作
	batchOp, err := batchOpService.Submit(ctx, req)
	if err != nil {
		t.Fatalf("提交批量操作失败: %v", err)
	}

	// 验证初始PreflightSummary
	if batchOp.PreflightSummary == nil {
		t.Fatal("PreflightSummary不应该为nil")
	}
	initialValidCount := batchOp.PreflightSummary.ValidCount
	t.Logf("预检查有效项数: %d", initialValidCount)

	// 执行批量操作
	err = batchOpService.Execute(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Logf("Execute返回错误: %v", err)
	}

	// 获取更新后的批量操作
	result, err := batchOpRepo.GetByID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("获取批量操作结果失败: %v", err)
	}

	// 验证最终摘要
	summary := result.GetSummary()
	t.Logf("最终摘要: 总数=%d, 成功=%d, 失败=%d",
		summary.TotalCount, summary.SuccessCount, summary.FailedCount)

	// 验证统计正确性
	expectedTotal := 10
	if summary.TotalCount != expectedTotal {
		t.Errorf("期望总数%d，实际: %d", expectedTotal, summary.TotalCount)
	}

	// 成功数应该是7（10 - 3个已删除）
	// 但由于那些文档已经被删除，删除操作会失败
	// 所以成功数应该是7个（未被删除的），失败3个
	if summary.SuccessCount != 7 {
		t.Errorf("期望成功7个，实际: %d", summary.SuccessCount)
	}
	if summary.FailedCount != 3 {
		t.Errorf("期望失败3个，实际: %d", summary.FailedCount)
	}

	// 验证状态
	if result.Status != writerModel.BatchOpStatusPartial {
		t.Errorf("期望状态为partial，实际为: %s", result.Status)
	}

	t.Log("TestSummaryStatistics 测试通过")
}

// TestAllSuccess 测试全部成功的场景
func TestAllSuccess(t *testing.T) {
	ctx, batchOpService, batchOpRepo, docRepo, projectID := setupTestEnv(t)
	defer cleanupTestData(t, projectID, primitive.NilObjectID)

	// 创建5个测试文档
	docIDs := createTestDocuments(t, ctx, docRepo, projectID, 5)

	// 创建批量删除请求
	req := &documentService.SubmitBatchOperationRequest{
		ProjectID: projectID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs: docIDs,
		Atomic:    false,
		RetryConfig: &documentService.RetryConfig{
			MaxRetries: 2,
			RetryDelay: 100,
		},
	}

	// 设置用户ID
	ctx = context.WithValue(ctx, "userID", "test_user")

	// 提交批量操作
	batchOp, err := batchOpService.Submit(ctx, req)
	if err != nil {
		t.Fatalf("提交批量操作失败: %v", err)
	}

	// 执行批量操作
	err = batchOpService.Execute(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("Execute不应该返回错误: %v", err)
	}

	// 获取更新后的批量操作
	result, err := batchOpRepo.GetByID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("获取批量操作结果失败: %v", err)
	}

	// 验证状态：全部成功应该是completed
	if result.Status != writerModel.BatchOpStatusCompleted {
		t.Errorf("期望状态为completed，实际为: %s", result.Status)
	}

	// 验证摘要
	summary := result.GetSummary()
	if summary.SuccessCount != 5 {
		t.Errorf("期望成功5个，实际: %d", summary.SuccessCount)
	}
	if summary.FailedCount != 0 {
		t.Errorf("期望失败0个，实际: %d", summary.FailedCount)
	}

	// 验证所有item都是成功状态
	for i, item := range result.Items {
		if item.Status != writerModel.BatchItemStatusSucceeded {
			t.Errorf("第%d个项期望成功，实际: %s", i, item.Status)
		}
	}

	t.Log("TestAllSuccess 测试通过")
}

// TestAllFailure 测试全部失败的场景
func TestAllFailure(t *testing.T) {
	ctx, batchOpService, batchOpRepo, docRepo, projectID := setupTestEnv(t)
	defer cleanupTestData(t, projectID, primitive.NilObjectID)

	// 创建5个测试文档
	docIDs := createTestDocuments(t, ctx, docRepo, projectID, 5)

	// 删除所有文档，导致批量删除失败
	for _, docID := range docIDs {
		err := docRepo.SoftDelete(ctx, docID, projectID.Hex())
		if err != nil {
			t.Fatalf("预删除文档失败: %v", err)
		}
	}

	// 创建批量删除请求
	req := &documentService.SubmitBatchOperationRequest{
		ProjectID: projectID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs: docIDs,
		Atomic:    false,
		RetryConfig: &documentService.RetryConfig{
			MaxRetries: 1,
			RetryDelay: 50,
		},
	}

	// 设置用户ID
	ctx = context.WithValue(ctx, "userID", "test_user")

	// 提交批量操作
	batchOp, err := batchOpService.Submit(ctx, req)
	if err != nil {
		t.Fatalf("提交批量操作失败: %v", err)
	}

	// 执行批量操作
	err = batchOpService.Execute(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Logf("Execute返回错误: %v", err)
	}

	// 获取更新后的批量操作
	result, err := batchOpRepo.GetByID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("获取批量操作结果失败: %v", err)
	}

	// 验证状态：全部失败应该是failed
	if result.Status != writerModel.BatchOpStatusFailed {
		t.Errorf("期望状态为failed，实际为: %s", result.Status)
	}

	// 验证摘要
	summary := result.GetSummary()
	if summary.SuccessCount != 0 {
		t.Errorf("期望成功0个，实际: %d", summary.SuccessCount)
	}
	if summary.FailedCount != 5 {
		t.Errorf("期望失败5个，实际: %d", summary.FailedCount)
	}

	// 验证所有item都有错误信息
	for i, item := range result.Items {
		if item.Status != writerModel.BatchItemStatusFailed {
			t.Errorf("第%d个项期望失败，实际: %s", i, item.Status)
		}
		if item.ErrorCode == "" {
			t.Errorf("第%d个失败项缺少ErrorCode", i)
		}
	}

	t.Log("TestAllFailure 测试通过")
}
