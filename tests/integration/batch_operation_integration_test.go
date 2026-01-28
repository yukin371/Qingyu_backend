package integration

import (
	"context"
	"testing"

	writerModel "Qingyu_backend/models/writer"
	"Qingyu_backend/repository/mongodb/writer"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	serviceBase "Qingyu_backend/service/base"
	"Qingyu_backend/service/writer/document"
	"Qingyu_backend/test/testutil"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestBatchOperation_DeleteDocuments 测试批量删除文档
func TestBatchOperation_DeleteDocuments(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	// 初始化repositories
	opLogRepo := writerRepo.NewOperationLogRepository(db)
	docRepo := writerRepo.NewMongoDocumentRepository(db)

	// 1. 创建测试项目
	project := testutil.CreateTestProject(primitive.NewObjectID().Hex())
	project.TouchForCreate()
	projectRepo := writerRepo.NewMongoProjectRepository(db)
	err := projectRepo.Create(ctx, project)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	userID := primitive.NewObjectID()

	// 2. 创建测试文档
	docs := make([]*writerModel.Document, 5)
	for i := 0; i < 5; i++ {
		doc := testutil.CreateTestDocument(
			project.ID.Hex(),
			testutil.WithDocumentTitle("测试文档"+string(rune('A'+i))),
		)
		doc.StableRef = "doc-" + primitive.NewObjectID().Hex()
		doc.OrderKey = "0|" + doc.StableRef
		doc.TouchForCreate()

		err := docRepo.Create(ctx, doc)
		if err != nil {
			t.Fatalf("Failed to create test document %d: %v", i, err)
		}
		docs[i] = doc
	}

	// 3. 创建EventBus和BatchOperationService
	eventBus := serviceBase.NewSimpleEventBus()
	batchOpRepo := writerRepo.NewBatchOperationRepository(db)
	batchOpSvc := document.NewBatchOperationService(
		batchOpRepo.(*writerRepo.BatchOperationRepositoryImpl),
		opLogRepo.(*writerRepo.OperationLogRepositoryImpl),
		docRepo,
		eventBus,
	)

	req := &document.SubmitBatchOperationRequest{
		ProjectID: project.ID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs:  []string{docs[0].ID.Hex(), docs[1].ID.Hex()},
		Atomic:     true,
		UserID:     userID.Hex(),
	}

	batchOp, err := batchOpSvc.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Failed to submit batch operation: %v", err)
	}

	// 4. 验证Preflight结果
	if batchOp.PreflightSummary == nil {
		t.Error("PreflightSummary should not be nil")
	}
	if batchOp.PreflightSummary.ValidCount != 2 {
		t.Errorf("Expected 2 valid documents, got %d", batchOp.PreflightSummary.ValidCount)
	}

	// 5. 执行批量操作
	err = batchOpSvc.Execute(ctx, batchOp.ID)
	if err != nil {
		t.Fatalf("Failed to execute batch operation: %v", err)
	}

	// 6. 验证文档已软删除
	for i := 0; i < 2; i++ {
		doc, err := docRepo.GetByID(ctx, docs[i].ID.Hex())
		if err != nil {
			t.Errorf("Document %s should be retrievable (soft delete)", docs[i].ID.Hex())
		}
		if !doc.IsDeleted() {
			t.Errorf("Document %s should be deleted", docs[i].ID.Hex())
		}
	}

	// 7. 验证其他文档未被删除
	for i := 2; i < 5; i++ {
		doc, err := docRepo.GetByID(ctx, docs[i].ID.Hex())
		if err != nil {
			t.Errorf("Document %s should not be deleted", docs[i].ID.Hex())
		}
		if doc.IsDeleted() {
			t.Errorf("Document %s should not be deleted", docs[i].ID.Hex())
		}
	}

	// 8. 验证OperationLog已创建
	logs, err := opLogRepo.GetByChainID(ctx, batchOp.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get operation logs: %v", err)
	}
	if len(logs) == 0 {
		t.Errorf("Expected at least 1 operation log, got %d", len(logs))
	}
}

// TestBatchOperation_Undo 测试撤销批量操作
func TestBatchOperation_Undo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	// 初始化repositories
	opLogRepo := writerRepo.NewOperationLogRepository(db)
	docRepo := writerRepo.NewMongoDocumentRepository(db)

	// 创建测试项目
	project := testutil.CreateTestProject(primitive.NewObjectID().Hex())
	project.TouchForCreate()
	projectRepo := writerRepo.NewMongoProjectRepository(db)
	err := projectRepo.Create(ctx, project)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	userID := primitive.NewObjectID()

	// 创建测试文档
	docs := make([]*writerModel.Document, 3)
	for i := 0; i < 3; i++ {
		doc := testutil.CreateTestDocument(
			project.ID.Hex(),
			testutil.WithDocumentTitle("测试文档"+string(rune('A'+i))),
		)
		doc.StableRef = "doc-" + primitive.NewObjectID().Hex()
		doc.OrderKey = "0|" + doc.StableRef
		doc.TouchForCreate()

		err := docRepo.Create(ctx, doc)
		if err != nil {
			t.Fatalf("Failed to create test document %d: %v", i, err)
		}
		docs[i] = doc
	}

	// 创建EventBus和BatchOperationService
	eventBus := serviceBase.NewSimpleEventBus()
	batchOpRepo := writerRepo.NewBatchOperationRepository(db)
	batchOpSvc := document.NewBatchOperationService(
		batchOpRepo.(*writerRepo.BatchOperationRepositoryImpl),
		opLogRepo.(*writerRepo.OperationLogRepositoryImpl),
		docRepo,
		eventBus,
	)

	// 执行批量删除
	req := &document.SubmitBatchOperationRequest{
		ProjectID: project.ID.Hex(),
		Type:      writerModel.BatchOpTypeDelete,
		TargetIDs:  []string{docs[0].ID.Hex(), docs[1].ID.Hex()},
		Atomic:     true,
		UserID:     userID.Hex(),
	}

	batchOp, err := batchOpSvc.Submit(ctx, req)
	if err != nil {
		t.Fatalf("Failed to submit: %v", err)
	}

	err = batchOpSvc.Execute(ctx, batchOp.ID)
	if err != nil {
		t.Fatalf("Failed to execute: %v", err)
	}

	// 验证文档已删除
	for i := 0; i < 2; i++ {
		doc, _ := docRepo.GetByID(ctx, docs[i].ID.Hex())
		if !doc.IsDeleted() {
			t.Errorf("Document %s should be deleted before undo", docs[i].ID.Hex())
		}
	}

	// 撤销操作
	err = batchOpSvc.Undo(ctx, batchOp.ID, userID)
	if err != nil {
		t.Fatalf("Failed to undo: %v", err)
	}

	// 验证文档已恢复
	for i := 0; i < 2; i++ {
		doc, err := docRepo.GetByID(ctx, docs[i].ID.Hex())
		if err != nil {
			t.Errorf("Document %s should be restored after undo", docs[i].ID.Hex())
		}
		if doc.IsDeleted() {
			t.Errorf("Document %s should not be deleted after undo", docs[i].ID.Hex())
		}
	}

	// 验证第三个文档未被影响
	doc3, _ := docRepo.GetByID(ctx, docs[2].ID.Hex())
	if doc3.IsDeleted() {
		t.Errorf("Document %s should not be affected by undo", docs[2].ID.Hex())
	}
}
