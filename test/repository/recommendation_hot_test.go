package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	mongoReco "Qingyu_backend/repository/mongodb/recommendation"
	"Qingyu_backend/test/testutil"
)

// TestHotRecommendationRepository_GetHotBooks 测试获取热门书籍
func TestHotRecommendationRepository_GetHotBooks(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoHotRecommendationRepository(db)
	ctx := context.Background()

	// 准备测试数据（book_statistics集合）
	statsCollection := db.Collection("book_statistics")
	now := time.Now()

	stats := []interface{}{
		bson.M{
			"book_id":        "book_001",
			"views":          1000,
			"favorites":      500,
			"average_rating": 4.5,
			"updated_at":     now,
		},
		bson.M{
			"book_id":        "book_002",
			"views":          800,
			"favorites":      600,
			"average_rating": 4.8,
			"updated_at":     now,
		},
		bson.M{
			"book_id":        "book_003",
			"views":          500,
			"favorites":      200,
			"average_rating": 4.0,
			"updated_at":     now,
		},
	}

	_, err := statsCollection.InsertMany(ctx, stats)
	require.NoError(t, err)

	// 获取热门书籍（最近7天）
	results, err := repo.GetHotBooks(ctx, 10, 7)

	// 验证结果
	require.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.GreaterOrEqual(t, len(results), 3)

	// 第一个应该是热度最高的（book_002：收藏量最高）
	assert.Equal(t, "book_002", results[0])
}

// TestHotRecommendationRepository_GetHotBooks_NoData 测试无数据情况
func TestHotRecommendationRepository_GetHotBooks_NoData(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoHotRecommendationRepository(db)

	ctx := context.Background()
	results, err := repo.GetHotBooks(ctx, 10, 7)

	require.NoError(t, err)
	assert.Empty(t, results)
}

// TestHotRecommendationRepository_GetHotBooksByCategory 测试获取分类热门书籍
func TestHotRecommendationRepository_GetHotBooksByCategory(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoHotRecommendationRepository(db)
	ctx := context.Background()
	now := time.Now()

	// 准备测试数据（books集合）
	booksCollection := db.Collection("books")
	books := []interface{}{
		bson.M{
			"_id":        "book_xuanhuan_001",
			"category":   "玄幻",
			"status":     "published",
			"created_at": now,
		},
		bson.M{
			"_id":        "book_xuanhuan_002",
			"category":   "玄幻",
			"status":     "published",
			"created_at": now,
		},
		bson.M{
			"_id":        "book_dushi_001",
			"category":   "都市",
			"status":     "published",
			"created_at": now,
		},
	}
	_, err := booksCollection.InsertMany(ctx, books)
	require.NoError(t, err)

	// 准备统计数据
	statsCollection := db.Collection("book_statistics")
	stats := []interface{}{
		bson.M{
			"book_id":        "book_xuanhuan_001",
			"views":          1000,
			"favorites":      500,
			"average_rating": 4.5,
			"updated_at":     now,
		},
		bson.M{
			"book_id":        "book_xuanhuan_002",
			"views":          800,
			"favorites":      400,
			"average_rating": 4.3,
			"updated_at":     now,
		},
		bson.M{
			"book_id":        "book_dushi_001",
			"views":          600,
			"favorites":      300,
			"average_rating": 4.0,
			"updated_at":     now,
		},
	}
	_, err = statsCollection.InsertMany(ctx, stats)
	require.NoError(t, err)

	// 获取玄幻分类的热门书籍
	results, err := repo.GetHotBooksByCategory(ctx, "玄幻", 10, 7)

	// 验证结果
	require.NoError(t, err)
	assert.Len(t, results, 2) // 只有2本玄幻书
	// 验证都是玄幻分类的书
	for _, bookID := range results {
		assert.Contains(t, bookID, "xuanhuan")
	}
}

// TestHotRecommendationRepository_GetTrendingBooks 测试获取飙升书籍
func TestHotRecommendationRepository_GetTrendingBooks(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoHotRecommendationRepository(db)
	ctx := context.Background()

	// 准备最近3天的数据
	statsCollection := db.Collection("book_statistics")
	recentTime := time.Now().Add(-1 * time.Hour) // 1小时前（在3天内）

	stats := []interface{}{
		bson.M{
			"book_id":    "book_trending_001",
			"views":      500,
			"favorites":  300, // 收藏权重更高
			"updated_at": recentTime,
		},
		bson.M{
			"book_id":    "book_trending_002",
			"views":      800,
			"favorites":  100,
			"updated_at": recentTime,
		},
		bson.M{
			"book_id":    "book_old",
			"views":      1000,
			"favorites":  500,
			"updated_at": time.Now().AddDate(0, 0, -5), // 5天前（超过3天）
		},
	}
	_, err := statsCollection.InsertMany(ctx, stats)
	require.NoError(t, err)

	// 获取飙升书籍
	results, err := repo.GetTrendingBooks(ctx, 10)

	// 验证结果
	require.NoError(t, err)
	assert.NotEmpty(t, results)
	// book_old应该不在结果中（太旧）
	assert.NotContains(t, results, "book_old")
}

// TestHotRecommendationRepository_GetNewPopularBooks 测试获取新书热门
func TestHotRecommendationRepository_GetNewPopularBooks(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoHotRecommendationRepository(db)
	ctx := context.Background()

	now := time.Now()
	recentTime := now.AddDate(0, 0, -10) // 10天前（在30天内）
	oldTime := now.AddDate(0, 0, -40)    // 40天前（超过30天）

	// 准备书籍数据
	booksCollection := db.Collection("books")
	books := []interface{}{
		bson.M{
			"_id":        "book_new_001",
			"status":     "published",
			"created_at": recentTime,
		},
		bson.M{
			"_id":        "book_new_002",
			"status":     "published",
			"created_at": recentTime,
		},
		bson.M{
			"_id":        "book_old_001",
			"status":     "published",
			"created_at": oldTime,
		},
	}
	_, err := booksCollection.InsertMany(ctx, books)
	require.NoError(t, err)

	// 准备统计数据
	statsCollection := db.Collection("book_statistics")
	stats := []interface{}{
		bson.M{
			"book_id":   "book_new_001",
			"views":     500,
			"favorites": 300,
		},
		bson.M{
			"book_id":   "book_new_002",
			"views":     400,
			"favorites": 200,
		},
		bson.M{
			"book_id":   "book_old_001",
			"views":     1000,
			"favorites": 800,
		},
	}
	_, err = statsCollection.InsertMany(ctx, stats)
	require.NoError(t, err)

	// 获取新书热门（最近30天）
	results, err := repo.GetNewPopularBooks(ctx, 10, 30)

	// 验证结果
	require.NoError(t, err)
	assert.Len(t, results, 2) // 只有2本新书
	// 验证都是新书
	for _, bookID := range results {
		assert.Contains(t, bookID, "new")
	}
	// 老书不应该在结果中
	assert.NotContains(t, results, "book_old_001")
}

// TestHotRecommendationRepository_Health 测试健康检查
func TestHotRecommendationRepository_Health(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoHotRecommendationRepository(db)

	ctx := context.Background()
	err := repo.Health(ctx)

	assert.NoError(t, err)
}

// TestHotRecommendationRepository_GetHotBooksByCategory_EmptyCategory 测试空分类
func TestHotRecommendationRepository_GetHotBooksByCategory_EmptyCategory(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := mongoReco.NewMongoHotRecommendationRepository(db)

	ctx := context.Background()
	_, err := repo.GetHotBooksByCategory(ctx, "", 10, 7)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category cannot be empty")
}
