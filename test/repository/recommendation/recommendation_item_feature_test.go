package recommendation_test

import (
	reco "Qingyu_backend/models/recommendation"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mongoReco "Qingyu_backend/repository/mongodb/recommendation"
	"Qingyu_backend/test/testutil"
)

// TestItemFeatureRepository_Create 测试创建物品特征
func TestItemFeatureRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)

	// 准备测试数据
	feature := &reco.ItemFeature{
		ItemID: "book_001",
		Tags: map[string]float64{
			"玄幻": 0.9,
			"修真": 0.8,
			"热血": 0.7,
		},
		Authors:    []string{"作者A", "作者B"},
		Categories: []string{"玄幻", "修真"},
		Vector:     []float64{0.1, 0.2, 0.3, 0.4},
	}

	// 执行创建
	ctx := context.Background()
	err := repo.Create(ctx, feature)

	// 验证结果
	require.NoError(t, err)
	assert.False(t, feature.CreatedAt.IsZero())
	assert.False(t, feature.UpdatedAt.IsZero())
}

// TestItemFeatureRepository_Upsert 测试Upsert功能
func TestItemFeatureRepository_Upsert(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)
	ctx := context.Background()

	itemID := "book_upsert"

	// 1. 第一次upsert（创建）
	feature1 := &reco.ItemFeature{
		ItemID: itemID,
		Tags: map[string]float64{
			"玄幻": 0.5,
		},
		Categories: []string{"玄幻"},
	}
	err := repo.Upsert(ctx, feature1)
	require.NoError(t, err)

	// 2. 第二次upsert（更新）
	feature2 := &reco.ItemFeature{
		ItemID: itemID,
		Tags: map[string]float64{
			"玄幻": 0.9, // 更新
			"修真": 0.8, // 新增
		},
		Authors:    []string{"作者A"},      // 新增
		Categories: []string{"玄幻", "修真"}, // 更新
	}
	err = repo.Upsert(ctx, feature2)
	require.NoError(t, err)

	// 3. 验证更新结果
	result, err := repo.GetByItemID(ctx, itemID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, len(result.Tags))
	assert.Equal(t, 0.9, result.Tags["玄幻"])
	assert.Equal(t, 0.8, result.Tags["修真"])
	assert.Equal(t, 1, len(result.Authors))
	assert.Equal(t, 2, len(result.Categories))
}

// TestItemFeatureRepository_GetByItemID 测试根据物品ID获取特征
func TestItemFeatureRepository_GetByItemID(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)
	ctx := context.Background()

	itemID := "book_get"

	// 先创建特征
	feature := &reco.ItemFeature{
		ItemID: itemID,
		Tags: map[string]float64{
			"玄幻": 0.9,
		},
		Authors:    []string{"作者A"},
		Categories: []string{"玄幻"},
	}
	err := repo.Create(ctx, feature)
	require.NoError(t, err)

	// 查询特征
	result, err := repo.GetByItemID(ctx, itemID)

	// 验证结果
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, itemID, result.ItemID)
	assert.Equal(t, 1, len(result.Tags))
	assert.Equal(t, 0.9, result.Tags["玄幻"])
}

// TestItemFeatureRepository_GetByItemID_NotFound 测试获取不存在的物品特征
func TestItemFeatureRepository_GetByItemID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)

	ctx := context.Background()
	result, err := repo.GetByItemID(ctx, "non_existent_item")

	require.NoError(t, err)
	assert.Nil(t, result)
}

// TestItemFeatureRepository_BatchGetByItemIDs 测试批量获取物品特征
func TestItemFeatureRepository_BatchGetByItemIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)
	ctx := context.Background()

	// 创建多个物品特征
	features := []*reco.ItemFeature{
		{
			ItemID:     "book_001",
			Tags:       map[string]float64{"玄幻": 0.9},
			Categories: []string{"玄幻"},
		},
		{
			ItemID:     "book_002",
			Tags:       map[string]float64{"都市": 0.8},
			Categories: []string{"都市"},
		},
		{
			ItemID:     "book_003",
			Tags:       map[string]float64{"科幻": 0.7},
			Categories: []string{"科幻"},
		},
	}

	for _, feature := range features {
		err := repo.Create(ctx, feature)
		require.NoError(t, err)
	}

	// 批量查询
	itemIDs := []string{"book_001", "book_002", "book_003"}
	results, err := repo.BatchGetByItemIDs(ctx, itemIDs)

	// 验证结果
	require.NoError(t, err)
	assert.Len(t, results, 3)
}

// TestItemFeatureRepository_BatchGetByItemIDs_EmptySlice 测试空切片
func TestItemFeatureRepository_BatchGetByItemIDs_EmptySlice(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)

	ctx := context.Background()
	results, err := repo.BatchGetByItemIDs(ctx, []string{})

	require.NoError(t, err)
	assert.Empty(t, results)
}

// TestItemFeatureRepository_GetByCategory 测试根据分类获取物品特征
func TestItemFeatureRepository_GetByCategory(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)
	ctx := context.Background()

	// 创建不同分类的物品特征
	features := []*reco.ItemFeature{
		{
			ItemID:     "book_xuanhuan_001",
			Tags:       map[string]float64{"玄幻": 0.9},
			Categories: []string{"玄幻"},
		},
		{
			ItemID:     "book_xuanhuan_002",
			Tags:       map[string]float64{"玄幻": 0.8},
			Categories: []string{"玄幻"},
		},
		{
			ItemID:     "book_dushi_001",
			Tags:       map[string]float64{"都市": 0.7},
			Categories: []string{"都市"},
		},
	}

	for _, feature := range features {
		err := repo.Create(ctx, feature)
		require.NoError(t, err)
	}

	// 查询玄幻分类
	results, err := repo.GetByCategory(ctx, "玄幻", 10)

	// 验证结果
	require.NoError(t, err)
	assert.Len(t, results, 2)
	for _, result := range results {
		assert.Contains(t, result.Categories, "玄幻")
	}
}

// TestItemFeatureRepository_GetByTags 测试根据标签获取物品特征
func TestItemFeatureRepository_GetByTags(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)
	ctx := context.Background()

	// 创建具有不同标签的物品特征
	features := []*reco.ItemFeature{
		{
			ItemID: "book_001",
			Tags: map[string]float64{
				"玄幻": 0.9,
				"修真": 0.8,
			},
			Categories: []string{"玄幻"},
		},
		{
			ItemID: "book_002",
			Tags: map[string]float64{
				"玄幻": 0.7,
				"热血": 0.6,
			},
			Categories: []string{"玄幻"},
		},
		{
			ItemID: "book_003",
			Tags: map[string]float64{
				"都市": 0.8,
			},
			Categories: []string{"都市"},
		},
	}

	for _, feature := range features {
		err := repo.Create(ctx, feature)
		require.NoError(t, err)
	}

	// 查询包含"玄幻"标签的物品
	searchTags := map[string]float64{
		"玄幻": 0.5,
	}
	results, err := repo.GetByTags(ctx, searchTags, 10)

	// 验证结果
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 2) // 至少包含book_001和book_002
}

// TestItemFeatureRepository_Delete 测试删除物品特征
func TestItemFeatureRepository_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)
	ctx := context.Background()

	itemID := "book_delete"

	// 先创建特征
	feature := &reco.ItemFeature{
		ItemID:     itemID,
		Tags:       map[string]float64{"玄幻": 0.9},
		Categories: []string{"玄幻"},
	}
	err := repo.Create(ctx, feature)
	require.NoError(t, err)

	// 删除特征
	err = repo.Delete(ctx, itemID)
	require.NoError(t, err)

	// 验证已删除
	result, err := repo.GetByItemID(ctx, itemID)
	require.NoError(t, err)
	assert.Nil(t, result)
}

// TestItemFeatureRepository_Health 测试健康检查
func TestItemFeatureRepository_Health(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoItemFeatureRepository(db)

	ctx := context.Background()
	err := repo.Health(ctx)

	assert.NoError(t, err)
}
