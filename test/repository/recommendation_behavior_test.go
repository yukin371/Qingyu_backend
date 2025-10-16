package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	reco "Qingyu_backend/models/recommendation/reco"
	mongoReco "Qingyu_backend/repository/mongodb/recommendation"
	"Qingyu_backend/test/testutil"
)

// TestBehaviorRepository_Create 测试创建用户行为
func TestBehaviorRepository_Create(t *testing.T) {
	// 跳过短测试
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	// 初始化测试数据库
	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	// 创建Repository
	repo := mongoReco.NewMongoBehaviorRepository(db)

	// 准备测试数据
	behavior := &reco.Behavior{
		UserID:       "test_user_001",
		ItemID:       "test_book_001",
		ChapterID:    "test_chapter_001",
		BehaviorType: "read",
		Value:        120.5,
		Metadata: map[string]interface{}{
			"readTime": 120,
			"progress": 0.3,
		},
	}

	// 执行创建
	ctx := context.Background()
	err := repo.Create(ctx, behavior)

	// 验证结果
	require.NoError(t, err)
	assert.NotEmpty(t, behavior.ID)
	assert.False(t, behavior.CreatedAt.IsZero())
	assert.False(t, behavior.OccurredAt.IsZero())
}

// TestBehaviorRepository_Create_NilBehavior 测试创建nil行为
func TestBehaviorRepository_Create_NilBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoBehaviorRepository(db)

	ctx := context.Background()
	err := repo.Create(ctx, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "behavior cannot be nil")
}

// TestBehaviorRepository_BatchCreate 测试批量创建用户行为
func TestBehaviorRepository_BatchCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoBehaviorRepository(db)

	// 准备批量测试数据
	behaviors := []*reco.Behavior{
		{
			UserID:       "test_user_001",
			ItemID:       "test_book_001",
			BehaviorType: "view",
			Value:        1.0,
		},
		{
			UserID:       "test_user_001",
			ItemID:       "test_book_002",
			BehaviorType: "click",
			Value:        1.0,
		},
		{
			UserID:       "test_user_001",
			ItemID:       "test_book_003",
			BehaviorType: "collect",
			Value:        1.0,
		},
	}

	// 执行批量创建
	ctx := context.Background()
	err := repo.BatchCreate(ctx, behaviors)

	// 验证结果
	require.NoError(t, err)
	for _, behavior := range behaviors {
		assert.False(t, behavior.CreatedAt.IsZero())
		assert.False(t, behavior.OccurredAt.IsZero())
	}
}

// TestBehaviorRepository_BatchCreate_EmptySlice 测试批量创建空切片
func TestBehaviorRepository_BatchCreate_EmptySlice(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoBehaviorRepository(db)

	ctx := context.Background()
	err := repo.BatchCreate(ctx, []*reco.Behavior{})

	assert.NoError(t, err) // 空切片不应该报错
}

// TestBehaviorRepository_GetByUser 测试获取用户行为记录
func TestBehaviorRepository_GetByUser(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoBehaviorRepository(db)
	ctx := context.Background()

	// 准备测试数据
	userID := "test_user_001"
	behaviors := []*reco.Behavior{
		{
			UserID:       userID,
			ItemID:       "book_001",
			BehaviorType: "read",
			Value:        1.0,
			OccurredAt:   time.Now().Add(-3 * time.Hour),
		},
		{
			UserID:       userID,
			ItemID:       "book_002",
			BehaviorType: "collect",
			Value:        1.0,
			OccurredAt:   time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:       userID,
			ItemID:       "book_003",
			BehaviorType: "like",
			Value:        1.0,
			OccurredAt:   time.Now().Add(-1 * time.Hour),
		},
	}

	// 插入测试数据
	err := repo.BatchCreate(ctx, behaviors)
	require.NoError(t, err)

	// 查询用户行为（限制2条）
	result, err := repo.GetByUser(ctx, userID, 2)

	// 验证结果
	require.NoError(t, err)
	assert.Len(t, result, 2)
	// 应该按时间倒序返回（最近的在前）
	assert.Equal(t, "book_003", result[0].ItemID)
	assert.Equal(t, "book_002", result[1].ItemID)
}

// TestBehaviorRepository_GetByUser_EmptyUserID 测试空用户ID
func TestBehaviorRepository_GetByUser_EmptyUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoBehaviorRepository(db)

	ctx := context.Background()
	_, err := repo.GetByUser(ctx, "", 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// TestBehaviorRepository_GetByUser_NoData 测试用户无行为记录
func TestBehaviorRepository_GetByUser_NoData(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoBehaviorRepository(db)

	ctx := context.Background()
	result, err := repo.GetByUser(ctx, "non_existent_user", 10)

	require.NoError(t, err)
	assert.Empty(t, result) // 应该返回空切片
}

// TestBehaviorRepository_CompleteFlow 测试完整流程
func TestBehaviorRepository_CompleteFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoBehaviorRepository(db)
	ctx := context.Background()

	userID := "test_user_flow"

	// 1. 记录浏览行为
	viewBehavior := &reco.Behavior{
		UserID:       userID,
		ItemID:       "book_001",
		BehaviorType: "view",
		Value:        1.0,
	}
	err := repo.Create(ctx, viewBehavior)
	require.NoError(t, err)

	// 2. 记录阅读行为
	readBehavior := &reco.Behavior{
		UserID:       userID,
		ItemID:       "book_001",
		ChapterID:    "chapter_001",
		BehaviorType: "read",
		Value:        300.0, // 阅读时长
		Metadata: map[string]interface{}{
			"readTime": 300,
			"progress": 0.5,
		},
	}
	err = repo.Create(ctx, readBehavior)
	require.NoError(t, err)

	// 3. 记录收藏行为
	collectBehavior := &reco.Behavior{
		UserID:       userID,
		ItemID:       "book_001",
		BehaviorType: "collect",
		Value:        1.0,
	}
	err = repo.Create(ctx, collectBehavior)
	require.NoError(t, err)

	// 4. 查询用户所有行为
	behaviors, err := repo.GetByUser(ctx, userID, 10)
	require.NoError(t, err)
	assert.Len(t, behaviors, 3)

	// 验证行为顺序（最近的在前）
	assert.Equal(t, "collect", behaviors[0].BehaviorType)
	assert.Equal(t, "read", behaviors[1].BehaviorType)
	assert.Equal(t, "view", behaviors[2].BehaviorType)
}
