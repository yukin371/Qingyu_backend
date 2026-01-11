package stats_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	statsModel "Qingyu_backend/models/stats"
	infra "Qingyu_backend/repository/interfaces/infrastructure"
	statsRepo "Qingyu_backend/repository/mongodb/stats"
	"Qingyu_backend/test/testutil"
)

// 测试辅助函数
func setupBookStatsRepo(t *testing.T) (interface {
	Create(ctx context.Context, bookStats *statsModel.BookStats) error
	GetByID(ctx context.Context, id string) (*statsModel.BookStats, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetByBookID(ctx context.Context, bookID string) (*statsModel.BookStats, error)
	GetByAuthorID(ctx context.Context, authorID string, limit, offset int64) ([]*statsModel.BookStats, error)
	GetByDateRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*statsModel.BookStats, error)
	CreateDailyStats(ctx context.Context, dailyStats *statsModel.BookStatsDaily) error
	GetDailyStats(ctx context.Context, bookID string, date time.Time) (*statsModel.BookStatsDaily, error)
	GetDailyStatsRange(ctx context.Context, bookID string, startDate, endDate time.Time) ([]*statsModel.BookStatsDaily, error)
	GetRevenueBreakdown(ctx context.Context, bookID string, startDate, endDate time.Time) (*statsModel.RevenueBreakdown, error)
	CalculateTotalRevenue(ctx context.Context, bookID string) (float64, error)
	CalculateRevenueByType(ctx context.Context, bookID string, revenueType string) (float64, error)
	GetTopChapters(ctx context.Context, bookID string) (*statsModel.TopChapters, error)
	AnalyzeViewTrend(ctx context.Context, bookID string, days int) (string, error)
	AnalyzeRevenueTrend(ctx context.Context, bookID string, days int) (string, error)
	CalculateAvgCompletionRate(ctx context.Context, bookID string) (float64, error)
	CalculateAvgDropOffRate(ctx context.Context, bookID string) (float64, error)
	CalculateAvgReadingDuration(ctx context.Context, bookID string) (float64, error)
	GetTopBooksByViews(ctx context.Context, limit int) ([]*statsModel.BookStats, error)
	GetTopBooksByRevenue(ctx context.Context, limit int) ([]*statsModel.BookStats, error)
	GetTopBooksByCompletion(ctx context.Context, limit int) ([]*statsModel.BookStats, error)
	BatchCreate(ctx context.Context, bookStats []*statsModel.BookStats) error
	BatchUpdate(ctx context.Context, updates []map[string]interface{}) error
	Count(ctx context.Context, filter infra.Filter) (int64, error)
	CountByAuthor(ctx context.Context, authorID string) (int64, error)
	Health(ctx context.Context) error
}, context.Context) {
	t.Helper()

	db, cleanup := testutil.SetupTestDB(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	repo := statsRepo.NewMongoBookStatsRepository(db)

	return repo, ctx
}

func createTestBookStats(bookID string) *statsModel.BookStats {
	return &statsModel.BookStats{
		BookID:         bookID,
		AuthorID:       "author_" + primitive.NewObjectID().Hex(),
		Title:          "测试书籍",
		TotalViews:     0,
		TotalWords:     10000,
		TotalChapter:   10,
		StatDate:       time.Now(),
		TotalComments:  0,
		TotalLikes:     0,
		TotalBookmarks: 0,
	}
}

// TestBookStatsRepository_Create 测试创建书籍统计
func TestBookStatsRepository_Create(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	bookStats := createTestBookStats("book_" + primitive.NewObjectID().Hex())

	err := repo.Create(ctx, bookStats)
	assert.NoError(t, err)
	assert.NotEmpty(t, bookStats.ID)
}

// TestBookStatsRepository_GetByID 测试根据ID获取
func TestBookStatsRepository_GetByID(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	bookStats := createTestBookStats("book_" + primitive.NewObjectID().Hex())
	err := repo.Create(ctx, bookStats)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, bookStats.ID)
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, bookStats.ID, found.ID)
}

// TestBookStatsRepository_GetByID_NotFound 测试获取不存在的记录
func TestBookStatsRepository_GetByID_NotFound(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	found, err := repo.GetByID(ctx, "nonexistent_id")
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestBookStatsRepository_Update 测试更新
func TestBookStatsRepository_Update(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	bookStats := createTestBookStats("book_" + primitive.NewObjectID().Hex())
	err := repo.Create(ctx, bookStats)
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	updates := map[string]interface{}{
		"total_views": int64(100),
		"title":       "更新后的标题",
	}

	err = repo.Update(ctx, bookStats.ID, updates)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, bookStats.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(100), found.TotalViews)
}

// TestBookStatsRepository_Delete 测试删除
func TestBookStatsRepository_Delete(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	bookStats := createTestBookStats("book_" + primitive.NewObjectID().Hex())
	err := repo.Create(ctx, bookStats)
	require.NoError(t, err)

	err = repo.Delete(ctx, bookStats.ID)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, bookStats.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}

// TestBookStatsRepository_GetByBookID 测试根据书籍ID获取统计
func TestBookStatsRepository_GetByBookID(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	bookID := "book_" + primitive.NewObjectID().Hex()
	bookStats := createTestBookStats(bookID)
	err := repo.Create(ctx, bookStats)
	require.NoError(t, err)

	found, err := repo.GetByBookID(ctx, bookID)
	require.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, bookID, found.BookID)
}

// TestBookStatsRepository_GetByAuthorID 测试根据作者ID获取作品列表
func TestBookStatsRepository_GetByAuthorID(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	authorID := "author_" + primitive.NewObjectID().Hex()
	for i := 0; i < 3; i++ {
		bookStats := createTestBookStats("book_" + primitive.NewObjectID().Hex())
		bookStats.AuthorID = authorID
		err := repo.Create(ctx, bookStats)
		require.NoError(t, err)
	}

	found, err := repo.GetByAuthorID(ctx, authorID, 10, 0)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(found), 3)
}

// TestBookStatsRepository_Count 测试统计数量
func TestBookStatsRepository_Count(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	for i := 0; i < 3; i++ {
		bookStats := createTestBookStats("book_" + primitive.NewObjectID().Hex())
		err := repo.Create(ctx, bookStats)
		require.NoError(t, err)
	}

	count, err := repo.Count(ctx, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))
}

// TestBookStatsRepository_Health 测试健康检查
func TestBookStatsRepository_Health(t *testing.T) {
	repo, ctx := setupBookStatsRepo(t)

	err := repo.Health(ctx)
	assert.NoError(t, err)
}
