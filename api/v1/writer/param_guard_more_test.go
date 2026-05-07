package writer

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimelineApiGetTimeline_RequiresProjectIDQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewTimelineApi(nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/timelines/timeline-1", "", gin.Params{
		{Key: "timelineId", Value: "timeline-1"},
	})

	api.GetTimeline(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestOutlineApiGetOutline_RequiresProjectIDQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewOutlineApi(nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/outlines/outline-1", "", gin.Params{
		{Key: "outlineId", Value: "outline-1"},
	})

	api.GetOutline(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestVersionApiCompareVersions_RequiresFromVersionQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewVersionApi(nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/documents/doc-1/versions/compare?toVersion=v2", "", gin.Params{
		{Key: "documentId", Value: "doc-1"},
	})

	api.CompareVersions(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "源版本ID不能为空")
}

func TestVersionApiCompareVersions_RequiresToVersionQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewVersionApi(nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/documents/doc-1/versions/compare?fromVersion=v1", "", gin.Params{
		{Key: "documentId", Value: "doc-1"},
	})

	api.CompareVersions(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "目标版本ID不能为空")
}

func TestKeywordApiSearchKeywords_RequiresQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewKeywordApi(nil, nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/projects/507f1f77bcf86cd799439011/keywords/search", "", gin.Params{
		{Key: "id", Value: "507f1f77bcf86cd799439011"},
	})

	api.SearchKeywords(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "查询关键词不能为空")
}

func TestStatsApiGetDailyStats_RejectsInvalidDaysQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewStatsApi(nil, nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/books/book-1/daily-stats?days=abc", "", gin.Params{
		{Key: "book_id", Value: "book-1"},
	})

	api.GetDailyStats(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "天数必须在1-365之间")
}

func TestStatsApiGetRetentionRate_RejectsOutOfRangeDaysQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewStatsApi(nil, nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/books/book-1/retention?days=91", "", gin.Params{
		{Key: "book_id", Value: "book-1"},
	})

	api.GetRetentionRate(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "天数必须在1-90之间")
}

func TestStatsApiGetBookRevenue_RejectsReversedDateRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewStatsApi(nil, nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/books/book-1/revenue?start_date=2026-05-08&end_date=2026-05-07", "", gin.Params{
		{Key: "book_id", Value: "book-1"},
	})

	api.GetBookRevenue(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "开始日期不能晚于结束日期")
}

func TestWriterStatsAggregateAPIGetViews_RejectsInvalidDaysQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewWriterStatsAggregateAPI(nil, nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/stats/views?bookId=book-1&days=0", "", nil)

	api.GetViews(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "天数必须在1-365之间")
}

func TestStoryHarnessApiGetChapterContext_RequiresChapterIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewStoryHarnessApi(nil, nil, nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/projects/project-1/chapters/context", "", gin.Params{
		{Key: "id", Value: "project-1"},
	})

	api.GetChapterContext(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "章节ID不能为空")
}

func TestChangeRequestApiListChangeRequests_RequiresChapterIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewChangeRequestApi(nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/projects/project-1/chapters/change-requests", "", gin.Params{
		{Key: "id", Value: "project-1"},
	})

	api.ListChangeRequests(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "章节ID不能为空")
}
