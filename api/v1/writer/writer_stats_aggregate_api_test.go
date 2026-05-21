package writer

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	bookstoreModel "Qingyu_backend/models/bookstore"
	sharedModel "Qingyu_backend/models/shared"
	statsModel "Qingyu_backend/models/stats"
	bookstoreRepo "Qingyu_backend/repository/mongodb/bookstore"
	statsRepo "Qingyu_backend/repository/mongodb/stats"
	readingStats "Qingyu_backend/service/reader/stats"
	"Qingyu_backend/test/testutil"
)

func setupWriterStatsAggregateAPITest(t *testing.T) (*WriterStatsAggregateAPI, string, string) {
	t.Helper()

	db, cleanup := testutil.SetupTestDB(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	bookRepo := bookstoreRepo.NewMongoBookRepository(db.Client(), db.Name())
	bookStatsRepo := statsRepo.NewMongoBookStatsRepository(db)
	readerBehaviorRepo := statsRepo.NewMongoReaderBehaviorRepository(db)
	chapterStatsRepo := statsRepo.NewMongoChapterStatsRepository(db)
	statsService := readingStats.NewReadingStatsService(chapterStatsRepo, readerBehaviorRepo, bookStatsRepo)

	bookOID := primitive.NewObjectID()
	bookID := bookOID.Hex()
	projectID := "project-lookup"
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	twoDaysAgo := today.AddDate(0, 0, -2)
	tenDaysAgo := today.AddDate(0, 0, -10)
	targetDate := now.AddDate(0, 0, -7)
	activeToday := today.Add(12 * time.Hour)

	must := func(err error) {
		if err != nil {
			t.Fatalf("seed writer stats fixtures failed: %v", err)
		}
	}

	must(bookRepo.Create(ctx, &bookstoreModel.Book{
		IdentifiedEntity: sharedModel.IdentifiedEntity{ID: bookOID},
		ProjectID:        &projectID,
		Title:            "测试作品",
		AuthorID:         "author-1",
	}))

	must(bookStatsRepo.Create(ctx, &statsModel.BookStats{
		BookID:             bookID,
		Title:              "测试作品",
		AuthorID:           "author-1",
		TotalViews:         120,
		TotalSubscribers:   33,
		TotalBookmarks:     7,
		TotalComments:      5,
		TotalWords:         45678,
		TotalRevenue:       88.8,
		AvgCompletionRate:  0.66,
		AvgDropOffRate:     0.12,
		AvgReadingDuration: 123.4,
		StatDate:           today,
	}))

	must(bookStatsRepo.CreateDailyStats(ctx, &statsModel.BookStatsDaily{
		BookID:           bookID,
		Date:             today,
		DailyViews:       8,
		DailySubscribers: 2,
		DailyRevenue:     12,
	}))
	must(bookStatsRepo.CreateDailyStats(ctx, &statsModel.BookStatsDaily{
		BookID:           bookID,
		Date:             twoDaysAgo,
		DailyViews:       5,
		DailySubscribers: 1,
		DailyRevenue:     3,
	}))
	must(bookStatsRepo.CreateDailyStats(ctx, &statsModel.BookStatsDaily{
		BookID:           bookID,
		Date:             tenDaysAgo,
		DailyViews:       4,
		DailySubscribers: 1,
		DailyRevenue:     2,
	}))

	must(readerBehaviorRepo.Create(ctx, &statsModel.ReaderBehavior{
		UserID:       "user-1",
		BookID:       bookID,
		ChapterID:    "chapter-1",
		BehaviorType: statsModel.BehaviorTypeView,
		ReadAt:       targetDate.Add(2 * time.Hour),
	}))
	must(readerBehaviorRepo.Create(ctx, &statsModel.ReaderBehavior{
		UserID:       "user-2",
		BookID:       bookID,
		ChapterID:    "chapter-1",
		BehaviorType: statsModel.BehaviorTypeView,
		ReadAt:       targetDate.Add(3 * time.Hour),
	}))
	must(readerBehaviorRepo.Create(ctx, &statsModel.ReaderBehavior{
		UserID:       "user-1",
		BookID:       bookID,
		ChapterID:    "chapter-1",
		BehaviorType: statsModel.BehaviorTypeView,
		ReadAt:       activeToday,
	}))

	must(chapterStatsRepo.Create(ctx, &statsModel.ChapterStats{
		BookID:         bookID,
		ChapterID:      "chapter-1",
		Title:          "第一章",
		ViewCount:      50,
		CompletionRate: 0.82,
		DropOffRate:    0.10,
		Revenue:        18,
		StatDate:       today,
	}))
	must(chapterStatsRepo.Create(ctx, &statsModel.ChapterStats{
		BookID:         bookID,
		ChapterID:      "chapter-2",
		Title:          "第二章",
		ViewCount:      75,
		CompletionRate: 0.64,
		DropOffRate:    0.22,
		Revenue:        22,
		StatDate:       today,
	}))

	return NewWriterStatsAggregateAPI(statsService, bookRepo), bookID, projectID
}

func decodeWriterStatsAggregateResponse(t *testing.T, body []byte) map[string]any {
	t.Helper()

	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

func TestWriterStatsAggregateAPIResolveAggregateBookIDPrefersExplicitParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewWriterStatsAggregateAPI(nil, nil)

	c, _ := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/overview?bookId=book-1&projectId=project-1", "", nil)
	bookID, ok := api.resolveAggregateBookID(c)
	assert.True(t, ok)
	assert.Equal(t, "book-1", bookID)

	c, _ = newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/overview?projectId=project-2", "", nil)
	bookID, ok = api.resolveAggregateBookID(c)
	assert.True(t, ok)
	assert.Equal(t, "project-2", bookID)
}

func TestSumViewMetricsAggregatesViewsAndSubscribers(t *testing.T) {
	views, subscribers := sumViewMetrics([]*statsModel.BookStatsDaily{
		{
			DailyViews:       12,
			DailySubscribers: 3,
		},
		{
			DailyViews:       5,
			DailySubscribers: 2,
		},
	})

	assert.Equal(t, int64(17), views)
	assert.Equal(t, int64(5), subscribers)
}

func TestWriterStatsAggregateAPIGetOverview_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api, bookID, _ := setupWriterStatsAggregateAPITest(t)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/overview?bookId="+bookID, "", nil)

	api.GetOverview(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeWriterStatsAggregateResponse(t, w.Body.Bytes())
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]any)
	assert.Equal(t, bookID, data["bookId"])
	assert.Equal(t, "测试作品", data["title"])
	assert.Equal(t, float64(120), data["totalViews"])
	assert.Equal(t, float64(33), data["subscribers"])
	assert.Equal(t, float64(7), data["favorites"])
	assert.Equal(t, float64(5), data["comments"])
	assert.Equal(t, float64(8), data["todayViews"])
	assert.Equal(t, float64(17), data["monthViews"])
	assert.Equal(t, float64(0), data["todayWords"])
	assert.Equal(t, float64(45678), data["wordCount"])
	assert.Equal(t, float64(88.8), data["totalRevenue"])
	assert.Equal(t, float64(0.5), data["retentionRate"])

	overview := data["overview"].(map[string]any)
	assert.Equal(t, float64(120), overview["views"])
	assert.Equal(t, float64(33), overview["subscribers"])
	assert.Equal(t, float64(7), overview["bookmarks"])
	assert.Equal(t, float64(5), overview["comments"])

	today := data["today"].(map[string]any)
	assert.Equal(t, float64(8), today["views"])
	assert.Equal(t, float64(2), today["subscribers"])
	assert.Equal(t, float64(0), today["words"])
}

func TestWriterStatsAggregateAPIGetViews_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api, bookID, _ := setupWriterStatsAggregateAPITest(t)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/views?bookId="+bookID+"&days=7", "", nil)

	api.GetViews(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeWriterStatsAggregateResponse(t, w.Body.Bytes())
	data := resp["data"].(map[string]any)
	assert.Equal(t, bookID, data["bookId"])
	assert.Equal(t, float64(7), data["days"])
	assert.Equal(t, float64(13), data["total"])

	items := data["items"].([]any)
	assert.Len(t, items, 2)
	first := items[0].(map[string]any)
	second := items[1].(map[string]any)
	assert.Equal(t, float64(5), first["dailyViews"])
	assert.Equal(t, float64(8), second["dailyViews"])
}

func TestWriterStatsAggregateAPIGetSubscribers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api, bookID, _ := setupWriterStatsAggregateAPITest(t)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/subscribers?bookId="+bookID+"&days=7", "", nil)

	api.GetSubscribers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeWriterStatsAggregateResponse(t, w.Body.Bytes())
	data := resp["data"].(map[string]any)
	assert.Equal(t, bookID, data["bookId"])
	assert.Equal(t, float64(3), data["total"])

	items := data["items"].([]any)
	assert.Len(t, items, 2)
	first := items[0].(map[string]any)
	second := items[1].(map[string]any)
	assert.Equal(t, float64(1), first["newSubscribers"])
	assert.Equal(t, float64(2), second["newSubscribers"])
}

func TestWriterStatsAggregateAPIGetChapters_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api, bookID, _ := setupWriterStatsAggregateAPITest(t)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/chapters?bookId="+bookID, "", nil)

	api.GetChapters(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeWriterStatsAggregateResponse(t, w.Body.Bytes())
	data := resp["data"].(map[string]any)
	assert.Equal(t, bookID, data["bookId"])
	assert.Equal(t, float64(2), data["total"])

	items := data["items"].([]any)
	assert.Len(t, items, 2)
	first := items[0].(map[string]any)
	second := items[1].(map[string]any)
	assert.Equal(t, "chapter-2", first["chapterId"])
	assert.Equal(t, "chapter-1", second["chapterId"])
	assert.Equal(t, float64(75), first["views"])
	assert.Equal(t, float64(50), second["views"])
}

func TestWriterStatsAggregateAPIGetToday_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api, bookID, _ := setupWriterStatsAggregateAPITest(t)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/today?bookId="+bookID, "", nil)

	api.GetToday(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeWriterStatsAggregateResponse(t, w.Body.Bytes())
	data := resp["data"].(map[string]any)
	assert.Equal(t, bookID, data["bookId"])
	assert.Equal(t, float64(8), data["todayViews"])
	assert.Equal(t, float64(2), data["todaySubscribers"])
	assert.Equal(t, float64(0), data["todayWords"])
	assert.Equal(t, float64(8), data["views"])
	assert.Equal(t, float64(2), data["subscribers"])
	assert.Equal(t, float64(0), data["words"])
}

func TestWriterStatsAggregateAPIGetOverview_ResolvesProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api, bookID, projectID := setupWriterStatsAggregateAPITest(t)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/overview?projectId="+projectID, "", nil)

	api.GetOverview(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeWriterStatsAggregateResponse(t, w.Body.Bytes())
	data := resp["data"].(map[string]any)
	assert.Equal(t, bookID, data["bookId"])
	assert.Equal(t, "测试作品", data["title"])
	assert.Equal(t, float64(120), data["totalViews"])
	assert.Equal(t, float64(8), data["todayViews"])
	assert.Equal(t, float64(17), data["monthViews"])
	assert.Equal(t, float64(0.5), data["retentionRate"])
}
