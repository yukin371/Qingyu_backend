package recommendation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	reco "Qingyu_backend/models/recommendation/reco"
	mongoReco "Qingyu_backend/repository/mongodb/recommendation"
	"Qingyu_backend/test/testutil"
)

// TestProfileRepository_Upsert_Create 测试创建用户画像
func TestProfileRepository_Upsert_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)

	// 准备测试数据
	profile := &reco.UserProfile{
		UserID: "test_user_001",
		Tags: map[string]float64{
			"玄幻": 0.8,
			"修真": 0.6,
			"热血": 0.7,
		},
		Authors: map[string]float64{
			"作者A": 0.9,
			"作者B": 0.7,
		},
		Categories: map[string]float64{
			"玄幻": 0.8,
			"都市": 0.5,
		},
	}

	// 执行upsert（创建）
	ctx := context.Background()
	err := repo.Upsert(ctx, profile)

	// 验证结果
	require.NoError(t, err)
	assert.False(t, profile.UpdatedAt.IsZero())
	assert.False(t, profile.CreatedAt.IsZero())
}

// TestProfileRepository_Upsert_Update 测试更新用户画像
func TestProfileRepository_Upsert_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)
	ctx := context.Background()

	userID := "test_user_update"

	// 1. 先创建画像
	profile1 := &reco.UserProfile{
		UserID: userID,
		Tags: map[string]float64{
			"玄幻": 0.5,
		},
		Categories: map[string]float64{
			"玄幻": 0.5,
		},
	}
	err := repo.Upsert(ctx, profile1)
	require.NoError(t, err)

	// 2. 更新画像（增加新标签）
	profile2 := &reco.UserProfile{
		UserID: userID,
		Tags: map[string]float64{
			"玄幻": 0.8, // 更新权重
			"修真": 0.6, // 新增标签
		},
		Categories: map[string]float64{
			"玄幻": 0.8,
			"都市": 0.4, // 新增分类
		},
	}
	err = repo.Upsert(ctx, profile2)
	require.NoError(t, err)

	// 3. 查询验证
	result, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, result)

	// 验证标签已更新
	assert.Equal(t, 0.8, result.Tags["玄幻"])
	assert.Equal(t, 0.6, result.Tags["修真"])
	assert.Equal(t, 0.8, result.Categories["玄幻"])
	assert.Equal(t, 0.4, result.Categories["都市"])
}

// TestProfileRepository_Upsert_NilProfile 测试nil画像
func TestProfileRepository_Upsert_NilProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)

	ctx := context.Background()
	err := repo.Upsert(ctx, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "profile cannot be nil")
}

// TestProfileRepository_Upsert_EmptyUserID 测试空用户ID
func TestProfileRepository_Upsert_EmptyUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)

	profile := &reco.UserProfile{
		UserID: "",
		Tags: map[string]float64{
			"玄幻": 0.5,
		},
	}

	ctx := context.Background()
	err := repo.Upsert(ctx, profile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// TestProfileRepository_GetByUserID 测试获取用户画像
func TestProfileRepository_GetByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)
	ctx := context.Background()

	userID := "test_user_get"

	// 先创建画像
	profile := &reco.UserProfile{
		UserID: userID,
		Tags: map[string]float64{
			"玄幻": 0.8,
			"修真": 0.6,
		},
		Authors: map[string]float64{
			"作者A": 0.9,
		},
		Categories: map[string]float64{
			"玄幻": 0.8,
		},
	}
	err := repo.Upsert(ctx, profile)
	require.NoError(t, err)

	// 查询画像
	result, err := repo.GetByUserID(ctx, userID)

	// 验证结果
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, 2, len(result.Tags))
	assert.Equal(t, 0.8, result.Tags["玄幻"])
	assert.Equal(t, 0.6, result.Tags["修真"])
	assert.Equal(t, 1, len(result.Authors))
	assert.Equal(t, 0.9, result.Authors["作者A"])
	assert.Equal(t, 1, len(result.Categories))
	assert.Equal(t, 0.8, result.Categories["玄幻"])
}

// TestProfileRepository_GetByUserID_NotFound 测试获取不存在的用户画像
func TestProfileRepository_GetByUserID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)

	ctx := context.Background()
	result, err := repo.GetByUserID(ctx, "non_existent_user")

	require.NoError(t, err)
	assert.Nil(t, result) // 不存在应该返回nil
}

// TestProfileRepository_GetByUserID_EmptyUserID 测试空用户ID
func TestProfileRepository_GetByUserID_EmptyUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)

	ctx := context.Background()
	_, err := repo.GetByUserID(ctx, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// TestProfileRepository_CompleteFlow 测试完整流程
func TestProfileRepository_CompleteFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoProfileRepository(db)
	ctx := context.Background()

	userID := "test_user_flow"

	// 1. 初始画像（新用户）
	initialProfile := &reco.UserProfile{
		UserID: userID,
		Tags: map[string]float64{
			"玄幻": 0.3,
		},
		Categories: map[string]float64{
			"玄幻": 0.3,
		},
	}
	err := repo.Upsert(ctx, initialProfile)
	require.NoError(t, err)

	// 2. 用户阅读行为后，更新画像（增加权重）
	updatedProfile := &reco.UserProfile{
		UserID: userID,
		Tags: map[string]float64{
			"玄幻": 0.6, // 权重增加
			"修真": 0.4, // 新增兴趣
		},
		Categories: map[string]float64{
			"玄幻": 0.6,
		},
	}
	err = repo.Upsert(ctx, updatedProfile)
	require.NoError(t, err)

	// 3. 继续阅读，画像持续更新
	finalProfile := &reco.UserProfile{
		UserID: userID,
		Tags: map[string]float64{
			"玄幻": 0.8, // 权重继续增加
			"修真": 0.7,
			"热血": 0.5, // 发现新兴趣
		},
		Authors: map[string]float64{
			"作者A": 0.9, // 喜欢的作者
		},
		Categories: map[string]float64{
			"玄幻": 0.8,
			"都市": 0.3, // 新增分类
		},
	}
	err = repo.Upsert(ctx, finalProfile)
	require.NoError(t, err)

	// 4. 获取最终画像
	result, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, result)

	// 验证最终画像
	assert.Equal(t, 3, len(result.Tags))
	assert.Equal(t, 0.8, result.Tags["玄幻"])
	assert.Equal(t, 0.7, result.Tags["修真"])
	assert.Equal(t, 0.5, result.Tags["热血"])
	assert.Equal(t, 1, len(result.Authors))
	assert.Equal(t, 0.9, result.Authors["作者A"])
	assert.Equal(t, 2, len(result.Categories))
	assert.Equal(t, 0.8, result.Categories["玄幻"])
	assert.Equal(t, 0.3, result.Categories["都市"])
}
