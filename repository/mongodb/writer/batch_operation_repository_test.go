package writer_test

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	"Qingyu_backend/test/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 测试辅助函数
func setupBatchOperationRepo(t *testing.T) (writerRepo.BatchOperationRepository, *mongo.Database, context.Context, func()) {
	t.Helper()
	db, cleanup := testutil.SetupTestDB(t)
	repo := writerRepo.NewBatchOperationRepository(db)
	ctx := context.Background()
	return repo, db, ctx, func() {
		// 清理batch_operations和batch_operation_items集合
		_ = db.Collection("batch_operations").Drop(ctx)
		_ = db.Collection("batch_operation_items").Drop(ctx)
		cleanup()
	}
}

func setupOperationLogRepo(t *testing.T) (writerRepo.OperationLogRepository, *mongo.Database, context.Context, func()) {
	t.Helper()
	db, cleanup := testutil.SetupTestDB(t)
	repo := writerRepo.NewOperationLogRepository(db)
	ctx := context.Background()
	return repo, db, ctx, func() {
		// 清理operation_logs集合
		_ = db.Collection("operation_logs").Drop(ctx)
		cleanup()
	}
}

// ==================== BatchOperationRepository 测试 ====================

// TestBatchOperationRepository_Create 测试创建批量操作
func TestBatchOperationRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	op := &writer.BatchOperation{
		ProjectID:   projectID,
		Type:        writer.BatchOpTypeDelete,
		TargetIDs:   []string{"doc-1", "doc-2"},
		Atomic:      true,
		CreatedBy:   userID.Hex(),
	}

	err := repo.Create(ctx, op)
	require.NoError(t, err)
	assert.False(t, op.ID.IsZero())
	assert.Equal(t, writer.BatchOpStatusPending, op.Status)
}

// TestBatchOperationRepository_GetByID 测试根据ID获取
func TestBatchOperationRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	op := &writer.BatchOperation{
		ProjectID:   projectID,
		Type:        writer.BatchOpTypeDelete,
		TargetIDs:   []string{"doc-1"},
		CreatedBy:   userID.Hex(),
	}
	err := repo.Create(ctx, op)
	require.NoError(t, err)

	// 获取操作
	retrieved, err := repo.GetByID(ctx, op.ID)
	require.NoError(t, err)
	assert.Equal(t, op.ID, retrieved.ID)
	assert.Equal(t, op.Type, retrieved.Type)
}

// TestBatchOperationRepository_GetByID_NotFound 测试获取不存在的操作
func TestBatchOperationRepository_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	_, err := repo.GetByID(ctx, primitive.NewObjectID())
	assert.Error(t, err)
	assert.Equal(t, writerRepo.ErrBatchOperationNotFound, err)
}

// TestBatchOperationRepository_GetByClientRequestID 测试幂等性检查
func TestBatchOperationRepository_GetByClientRequestID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	clientRequestID := "idempotent-test-123"
	userID := primitive.NewObjectID()

	// 第一次创建
	op1 := &writer.BatchOperation{
		ProjectID:       projectID,
		Type:            writer.BatchOpTypeDelete,
		TargetIDs:       []string{"doc-1"},
		ClientRequestID: clientRequestID,
		CreatedBy:       userID.Hex(),
	}
	err := repo.Create(ctx, op1)
	require.NoError(t, err)

	// 通过clientRequestID查询
	op2, err := repo.GetByClientRequestID(ctx, projectID, clientRequestID)
	require.NoError(t, err)
	assert.Equal(t, op1.ID, op2.ID)
	assert.Equal(t, op1.ClientRequestID, op2.ClientRequestID)
}

// TestBatchOperationRepository_GetByClientRequestID_NotFound 测试幂等性检查-未找到
func TestBatchOperationRepository_GetByClientRequestID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	_, err := repo.GetByClientRequestID(ctx, projectID, "non-existent-request-id")
	assert.Error(t, err)
	assert.Equal(t, writerRepo.ErrBatchOperationNotFound, err)
}

// TestBatchOperationRepository_UpdateStatus 测试更新状态
func TestBatchOperationRepository_UpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	op := &writer.BatchOperation{
		ProjectID:   projectID,
		Type:        writer.BatchOpTypeDelete,
		TargetIDs:   []string{"doc-1"},
		CreatedBy:   userID.Hex(),
		Status:      writer.BatchOpStatusPending,
	}
	err := repo.Create(ctx, op)
	require.NoError(t, err)

	// 更新状态为running
	err = repo.UpdateStatus(ctx, op.ID, writer.BatchOpStatusRunning)
	require.NoError(t, err)

	// 验证状态已更新
	retrieved, err := repo.GetByID(ctx, op.ID)
	require.NoError(t, err)
	assert.Equal(t, writer.BatchOpStatusRunning, retrieved.Status)
}

// TestBatchOperationRepository_UpdateStatus_NotFound 测试更新状态-未找到
func TestBatchOperationRepository_UpdateStatus_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	err := repo.UpdateStatus(ctx, primitive.NewObjectID(), writer.BatchOpStatusRunning)
	assert.Error(t, err)
	assert.Equal(t, writerRepo.ErrBatchOperationNotFound, err)
}

// TestBatchOperationRepository_Update 测试更新操作
func TestBatchOperationRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	op := &writer.BatchOperation{
		ProjectID:   projectID,
		Type:        writer.BatchOpTypeDelete,
		TargetIDs:   []string{"doc-1"},
		CreatedBy:   userID.Hex(),
		Status:      writer.BatchOpStatusPending,
	}
	err := repo.Create(ctx, op)
	require.NoError(t, err)

	// 修改状态
	op.Status = writer.BatchOpStatusRunning
	now := time.Now()
		// Convert time.Time to primitive.DateTime
		dt := primitive.NewDateTimeFromTime(now)
		op.StartedAt = &dt

	// 更新
	err = repo.Update(ctx, op)
	require.NoError(t, err)

	// 验证已更新
	retrieved, err := repo.GetByID(ctx, op.ID)
	require.NoError(t, err)
	assert.Equal(t, writer.BatchOpStatusRunning, retrieved.Status)
	assert.NotNil(t, retrieved.StartedAt)
}

// TestBatchOperationRepository_ListByProject 测试查询项目操作列表
func TestBatchOperationRepository_ListByProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建多个操作
	for i := 0; i < 3; i++ {
		op := &writer.BatchOperation{
			ProjectID:   projectID,
			Type:        writer.BatchOpTypeDelete,
			TargetIDs:   []string{primitive.NewObjectID().Hex()},
			CreatedBy:   userID.Hex(),
			Status:      writer.BatchOpStatusCompleted,
		}
		err := repo.Create(ctx, op)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // 确保时间差异
	}

	// 查询所有
	ops, err := repo.ListByProject(ctx, projectID, &writerRepo.ListOptions{
		Limit: 10,
		Skip:  0,
	})
	require.NoError(t, err)
	assert.Len(t, ops, 3)
}

// TestBatchOperationRepository_ListByProject_WithStatus 测试按状态查询
func TestBatchOperationRepository_ListByProject_WithStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建不同状态的操作
	pendingOp := &writer.BatchOperation{
		ProjectID:   projectID,
		Type:        writer.BatchOpTypeDelete,
		TargetIDs:   []string{"doc-1"},
		CreatedBy:   userID.Hex(),
		Status:      writer.BatchOpStatusPending,
	}
	err := repo.Create(ctx, pendingOp)
	require.NoError(t, err)

	completedOp := &writer.BatchOperation{
		ProjectID:   projectID,
		Type:        writer.BatchOpTypeDelete,
		TargetIDs:   []string{"doc-2"},
		CreatedBy:   userID.Hex(),
		Status:      writer.BatchOpStatusCompleted,
	}
	err = repo.Create(ctx, completedOp)
	require.NoError(t, err)

	// 查询pending状态
	ops, err := repo.ListByProject(ctx, projectID, &writerRepo.ListOptions{
		Status: writer.BatchOpStatusPending,
	})
	require.NoError(t, err)
	assert.Len(t, ops, 1)
	assert.Equal(t, writer.BatchOpStatusPending, ops[0].Status)
}

// TestBatchOperationRepository_GetRunningCount 测试获取运行中操作数量
func TestBatchOperationRepository_GetRunningCount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupBatchOperationRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建2个运行中操作
	for i := 0; i < 2; i++ {
		op := &writer.BatchOperation{
			ProjectID:   projectID,
			Type:        writer.BatchOpTypeDelete,
			TargetIDs:   []string{"doc-1"},
			CreatedBy:   userID.Hex(),
			Status:      writer.BatchOpStatusRunning,
		}
		err := repo.Create(ctx, op)
		require.NoError(t, err)
	}

	// 创建1个已完成操作
	completedOp := &writer.BatchOperation{
		ProjectID:   projectID,
		Type:        writer.BatchOpTypeDelete,
		TargetIDs:   []string{"doc-2"},
		CreatedBy:   userID.Hex(),
		Status:      writer.BatchOpStatusCompleted,
	}
	err := repo.Create(ctx, completedOp)
	require.NoError(t, err)

	// 统计运行中数量
	count, err := repo.GetRunningCount(ctx, projectID)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// ==================== OperationLogRepository 测试 ====================

// TestOperationLogRepository_Create 测试创建操作日志
func TestOperationLogRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	log := &writer.OperationLog{
		ProjectID:    projectID,
		UserID:       userID,
		CommandType:  writer.CommandDelete,
		TargetIDs:    []string{"doc-1"},
		CommandPayload: map[string]interface{}{
			"reason": "test",
		},
	}

	err := repo.Create(ctx, log)
	require.NoError(t, err)
	assert.False(t, log.ID.IsZero())
	assert.NotEmpty(t, log.ChainID)
	assert.Equal(t, writer.OpLogStatusExecuted, log.Status)
}

// TestOperationLogRepository_GetByID 测试根据ID获取
func TestOperationLogRepository_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	log := &writer.OperationLog{
		ProjectID:   projectID,
		UserID:      userID,
		CommandType: writer.CommandDelete,
		TargetIDs:   []string{"doc-1"},
	}
	err := repo.Create(ctx, log)
	require.NoError(t, err)

	// 获取日志
	retrieved, err := repo.GetByID(ctx, log.ID)
	require.NoError(t, err)
	assert.Equal(t, log.ID, retrieved.ID)
	assert.Equal(t, log.CommandType, retrieved.CommandType)
}

// TestOperationLogRepository_GetByID_NotFound 测试获取不存在的日志
func TestOperationLogRepository_GetByID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	_, err := repo.GetByID(ctx, primitive.NewObjectID())
	assert.Error(t, err)
	assert.Equal(t, writerRepo.ErrOperationLogNotFound, err)
}

// TestOperationLogRepository_GetByChainID 测试根据链ID获取
func TestOperationLogRepository_GetByChainID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	chainID := "test-chain-123"

	// 创建多个同链的日志
	for i := 0; i < 3; i++ {
		log := &writer.OperationLog{
			ProjectID:    projectID,
			UserID:       userID,
			ChainID:      chainID,
			CommandType:  writer.CommandDelete,
			TargetIDs:    []string{primitive.NewObjectID().Hex()},
		}
		err := repo.Create(ctx, log)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// 查询链
	logs, err := repo.GetByChainID(ctx, chainID)
	require.NoError(t, err)
	assert.Len(t, logs, 3)
}

// TestOperationLogRepository_GetLatestByProject 测试获取项目最新操作（Undo栈）
func TestOperationLogRepository_GetLatestByProject(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	// 创建5个操作日志
	for i := 0; i < 5; i++ {
		log := &writer.OperationLog{
			ProjectID:   projectID,
			UserID:      userID,
			CommandType: writer.CommandDelete,
			TargetIDs:   []string{primitive.NewObjectID().Hex()},
			IsCommitted: true,
		}
		err := repo.Create(ctx, log)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// 获取最新3个（用于Undo栈）
	logs, err := repo.GetLatestByProject(ctx, projectID, 3)
	require.NoError(t, err)
	assert.Len(t, logs, 3)
	// 应该按created_at降序排列
}

// TestOperationLogRepository_UpdateStatus 测试更新状态
func TestOperationLogRepository_UpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	log := &writer.OperationLog{
		ProjectID:    projectID,
		UserID:       userID,
		CommandType:  writer.CommandDelete,
		TargetIDs:    []string{"doc-1"},
		Status:       writer.OpLogStatusExecuted,
		IsCommitted:  true,
	}
	err := repo.Create(ctx, log)
	require.NoError(t, err)

	// 更新为undone
	err = repo.UpdateStatus(ctx, log.ID, writer.OpLogStatusUndone)
	require.NoError(t, err)

	// 验证状态已更新
	retrieved, err := repo.GetByID(ctx, log.ID)
	require.NoError(t, err)
	assert.Equal(t, writer.OpLogStatusUndone, retrieved.Status)
	assert.NotNil(t, retrieved.UndoneAt)
}

// TestOperationLogRepository_UpdateStatus_NotFound 测试更新状态-未找到
func TestOperationLogRepository_UpdateStatus_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	err := repo.UpdateStatus(ctx, primitive.NewObjectID(), writer.OpLogStatusUndone)
	assert.Error(t, err)
	assert.Equal(t, writerRepo.ErrOperationLogNotFound, err)
}

// TestOperationLogRepository_MarkAsCommitted 测试标记为已提交
func TestOperationLogRepository_MarkAsCommitted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	projectID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	log := &writer.OperationLog{
		ProjectID:    projectID,
		UserID:       userID,
		CommandType:  writer.CommandDelete,
		TargetIDs:    []string{"doc-1"},
		IsCommitted:  false,
	}
	err := repo.Create(ctx, log)
	require.NoError(t, err)

	// 标记为已提交
	err = repo.MarkAsCommitted(ctx, log.ID)
	require.NoError(t, err)

	// 验证已标记
	retrieved, err := repo.GetByID(ctx, log.ID)
	require.NoError(t, err)
	assert.True(t, retrieved.IsCommitted)
}

// TestOperationLogRepository_MarkAsCommitted_NotFound 测试标记为已提交-未找到
func TestOperationLogRepository_MarkAsCommitted_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	repo, _, ctx, cleanup := setupOperationLogRepo(t)
	defer cleanup()

	err := repo.MarkAsCommitted(ctx, primitive.NewObjectID())
	assert.Error(t, err)
	assert.Equal(t, writerRepo.ErrOperationLogNotFound, err)
}
