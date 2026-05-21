package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type quotaDashboardRepoForAPITest struct {
	summary          *aiModels.DashboardSummary
	distribution     *aiModels.QuotaDistribution
	topConsumers     []aiModels.UserQuotaRanking
	trend            []aiModels.TrendPoint
	consumption      *aiModels.QuotaConsumptionSummary
	summaryErr       error
	distributionErr  error
	topConsumersErr  error
	trendErr         error
	consumptionErr   error
	lastTrendDays    int
	lastWorkflowType string
	lastGroupBy      string
	lastPage         int
	lastPageSize     int
}

func (s *quotaDashboardRepoForAPITest) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	return nil
}

func (s *quotaDashboardRepoForAPITest) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	return nil, aiModels.ErrQuotaNotFound
}

func (s *quotaDashboardRepoForAPITest) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	return nil
}

func (s *quotaDashboardRepoForAPITest) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaDashboardRepoForAPITest) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	return nil, nil
}

func (s *quotaDashboardRepoForAPITest) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaDashboardRepoForAPITest) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	return nil
}

func (s *quotaDashboardRepoForAPITest) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaDashboardRepoForAPITest) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaDashboardRepoForAPITest) GetQuotaStatistics(ctx context.Context, userID string) (*aiRepo.QuotaStatistics, error) {
	return nil, nil
}

func (s *quotaDashboardRepoForAPITest) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	return 0, nil
}

func (s *quotaDashboardRepoForAPITest) GetDashboardSummary(ctx context.Context) (*aiModels.DashboardSummary, error) {
	if s.summaryErr != nil {
		return nil, s.summaryErr
	}
	if s.summary != nil {
		return s.summary, nil
	}
	return &aiModels.DashboardSummary{}, nil
}

func (s *quotaDashboardRepoForAPITest) GetQuotaDistribution(ctx context.Context) (*aiModels.QuotaDistribution, error) {
	if s.distributionErr != nil {
		return nil, s.distributionErr
	}
	if s.distribution != nil {
		return s.distribution, nil
	}
	return &aiModels.QuotaDistribution{}, nil
}

func (s *quotaDashboardRepoForAPITest) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	if s.topConsumersErr != nil {
		return nil, s.topConsumersErr
	}
	return s.topConsumers, nil
}

func (s *quotaDashboardRepoForAPITest) GetConsumptionTrend(ctx context.Context, days int) ([]aiModels.TrendPoint, error) {
	s.lastTrendDays = days
	if s.trendErr != nil {
		return nil, s.trendErr
	}
	return s.trend, nil
}

func (s *quotaDashboardRepoForAPITest) GetConsumptionSummary(ctx context.Context, startTime, endTime time.Time, workflowType, groupBy string, page, pageSize int) (*aiModels.QuotaConsumptionSummary, error) {
	s.lastWorkflowType = workflowType
	s.lastGroupBy = groupBy
	s.lastPage = page
	s.lastPageSize = pageSize
	if s.consumptionErr != nil {
		return nil, s.consumptionErr
	}
	if s.consumption != nil {
		return s.consumption, nil
	}
	return &aiModels.QuotaConsumptionSummary{}, nil
}

func (s *quotaDashboardRepoForAPITest) ListUserQuotas(ctx context.Context, role, status, search string, page, limit int) ([]*aiModels.UserQuotaListItem, int64, error) {
	return nil, 0, nil
}

func (s *quotaDashboardRepoForAPITest) Health(ctx context.Context) error {
	return nil
}

type quotaDashboardAlertRepoForAPITest struct {
	recent []*aiModels.QuotaAlert
}

func (s *quotaDashboardAlertRepoForAPITest) Create(ctx context.Context, alert *aiModels.QuotaAlert) error {
	return nil
}

func (s *quotaDashboardAlertRepoForAPITest) GetByID(ctx context.Context, id string) (*aiModels.QuotaAlert, error) {
	return nil, nil
}

func (s *quotaDashboardAlertRepoForAPITest) List(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	return nil, 0, nil
}

func (s *quotaDashboardAlertRepoForAPITest) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	return nil, 0, nil
}

func (s *quotaDashboardAlertRepoForAPITest) Update(ctx context.Context, alert *aiModels.QuotaAlert) error {
	return nil
}

func (s *quotaDashboardAlertRepoForAPITest) GetRecentGlobal(ctx context.Context, limit int) ([]*aiModels.QuotaAlert, error) {
	return s.recent, nil
}

func (s *quotaDashboardAlertRepoForAPITest) CountByStatus(ctx context.Context) (map[aiModels.QuotaAlertStatus]int64, error) {
	return map[aiModels.QuotaAlertStatus]int64{}, nil
}

func (s *quotaDashboardAlertRepoForAPITest) Health(ctx context.Context) error {
	return nil
}

type quotaDashboardSummaryReaderForAPITest struct {
	response         *pb.QuotaConsumptionSummaryResponse
	err              error
	lastTimeRange    string
	lastWorkflowType string
	lastGroupBy      string
	lastPage         int32
	lastPageSize     int32
}

func (s *quotaDashboardSummaryReaderForAPITest) GetQuotaConsumptionSummary(
	ctx context.Context,
	timeRange string,
	workflowType string,
	groupBy string,
	page int32,
	pageSize int32,
) (*pb.QuotaConsumptionSummaryResponse, error) {
	s.lastTimeRange = timeRange
	s.lastWorkflowType = workflowType
	s.lastGroupBy = groupBy
	s.lastPage = page
	s.lastPageSize = pageSize
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}

type quotaConsistencyRunnerForAPITest struct {
	called int
	err    error
}

func (s *quotaConsistencyRunnerForAPITest) RunConsistencyCheck(ctx context.Context) error {
	s.called++
	return s.err
}

func setupQuotaDashboardAPITestRouter(
	repo *quotaDashboardRepoForAPITest,
	alertRepo *quotaDashboardAlertRepoForAPITest,
	reader *quotaDashboardSummaryReaderForAPITest,
	runner *quotaConsistencyRunnerForAPITest,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := aiService.NewQuotaDashboardService(repo, alertRepo, nil)
	if reader != nil {
		service.SetConsumptionSummaryReader(reader)
	}
	if runner != nil {
		service.SetConsistencyRunner(runner)
	}
	api := NewQuotaDashboardAPI(service)

	router := gin.New()
	router.GET("/dashboard", api.GetDashboard)
	router.GET("/statistics/global", api.GetStatistics)
	router.GET("/statistics/trend", api.GetTrend)
	router.GET("/statistics/reconciliation", api.GetReconciliationSummary)
	router.POST("/dashboard/refresh", api.RefreshCache)
	router.POST("/statistics/reconciliation/check", api.RunConsistencyCheck)
	return router
}

func TestQuotaDashboardAPIGetTrendDefaultsInvalidDaysToSeven(t *testing.T) {
	repo := &quotaDashboardRepoForAPITest{
		trend: []aiModels.TrendPoint{{Date: "2026-05-20", Consumption: 42, Users: 3}},
	}
	router := setupQuotaDashboardAPITestRouter(repo, &quotaDashboardAlertRepoForAPITest{}, nil, nil)

	req, _ := http.NewRequest("GET", "/statistics/trend?days=-7", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 7, repo.lastTrendDays)

	var resp struct {
		Data []struct {
			Date        string `json:"date"`
			Consumption int    `json:"consumption"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Data, 1)
	assert.Equal(t, "2026-05-20", resp.Data[0].Date)
	assert.Equal(t, 42, resp.Data[0].Consumption)
}

func TestQuotaDashboardAPIGetReconciliationSummaryNormalizesPagingAndGroupBy(t *testing.T) {
	repo := &quotaDashboardRepoForAPITest{
		consumption: &aiModels.QuotaConsumptionSummary{
			GroupBy:      "user",
			TotalGroups:  1,
			TotalTokens:  120,
			TotalRecords: 2,
			Items: []aiModels.QuotaConsumptionSummaryItem{
				{GroupKey: "user-1", TotalTokens: 120, TotalRecords: 2},
			},
		},
	}
	reader := &quotaDashboardSummaryReaderForAPITest{
		response: &pb.QuotaConsumptionSummaryResponse{
			Success:      true,
			TimeRange:    "month",
			GroupBy:      "user",
			Page:         1,
			PageSize:     100,
			TotalGroups:  1,
			TotalTokens:  110,
			TotalRecords: 2,
			Items: []*pb.QuotaConsumptionSummaryItem{
				{GroupKey: "user-1", TotalTokens: 110, TotalRecords: 2},
			},
		},
	}
	router := setupQuotaDashboardAPITestRouter(repo, &quotaDashboardAlertRepoForAPITest{}, reader, nil)

	req, _ := http.NewRequest("GET", "/statistics/reconciliation?timeRange=month&workflowType=story_write&groupBy=unknown&page=0&pageSize=999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "story_write", repo.lastWorkflowType)
	assert.Equal(t, "user", repo.lastGroupBy)
	assert.Equal(t, 1, repo.lastPage)
	assert.Equal(t, 100, repo.lastPageSize)
	assert.Equal(t, "month", reader.lastTimeRange)
	assert.Equal(t, "story_write", reader.lastWorkflowType)
	assert.Equal(t, "user", reader.lastGroupBy)
	assert.EqualValues(t, 1, reader.lastPage)
	assert.EqualValues(t, 100, reader.lastPageSize)

	var resp struct {
		Data struct {
			GroupBy  string `json:"groupBy"`
			PageSize int    `json:"pageSize"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user", resp.Data.GroupBy)
	assert.Equal(t, 100, resp.Data.PageSize)
}

func TestQuotaDashboardAPIGetDashboardReturnsInternalErrorWhenSummaryFails(t *testing.T) {
	repo := &quotaDashboardRepoForAPITest{summaryErr: errors.New("boom")}
	alertRepo := &quotaDashboardAlertRepoForAPITest{
		recent: []*aiModels.QuotaAlert{
			{ID: primitive.NewObjectID(), Title: "unused", CreatedAt: time.Now()},
		},
	}
	router := setupQuotaDashboardAPITestRouter(repo, alertRepo, nil, nil)

	req, _ := http.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取仪表盘数据失败")
}

func TestQuotaDashboardAPIGetStatisticsReturnsInternalErrorWhenSummaryFails(t *testing.T) {
	repo := &quotaDashboardRepoForAPITest{summaryErr: errors.New("boom")}
	router := setupQuotaDashboardAPITestRouter(repo, &quotaDashboardAlertRepoForAPITest{}, nil, nil)

	req, _ := http.NewRequest("GET", "/statistics/global", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取统计数据失败")
}

func TestQuotaDashboardAPIRunConsistencyCheckInvokesRunner(t *testing.T) {
	repo := &quotaDashboardRepoForAPITest{}
	runner := &quotaConsistencyRunnerForAPITest{}
	router := setupQuotaDashboardAPITestRouter(repo, &quotaDashboardAlertRepoForAPITest{}, nil, runner)

	req, _ := http.NewRequest("POST", "/statistics/reconciliation/check", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, runner.called)

	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "一致性检查已执行", resp.Message)
}

func TestQuotaDashboardAPIRefreshCacheAndConsistencyCheckTailErrors(t *testing.T) {
	t.Run("refresh cache surfaces dashboard rebuild failure", func(t *testing.T) {
		repo := &quotaDashboardRepoForAPITest{summaryErr: errors.New("boom")}
		router := setupQuotaDashboardAPITestRouter(repo, &quotaDashboardAlertRepoForAPITest{}, nil, nil)

		req, _ := http.NewRequest("POST", "/dashboard/refresh", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "刷新缓存失败")
	})

	t.Run("run consistency check rejects missing runner", func(t *testing.T) {
		repo := &quotaDashboardRepoForAPITest{}
		router := setupQuotaDashboardAPITestRouter(repo, &quotaDashboardAlertRepoForAPITest{}, nil, nil)

		req, _ := http.NewRequest("POST", "/statistics/reconciliation/check", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "执行一致性检查失败")
	})
}

func TestQuotaDashboardAPIGetReconciliationSummarySurfacesReaderFailure(t *testing.T) {
	repo := &quotaDashboardRepoForAPITest{
		consumption: &aiModels.QuotaConsumptionSummary{
			GroupBy:      "user",
			TotalGroups:  1,
			TotalTokens:  100,
			TotalRecords: 1,
			Items: []aiModels.QuotaConsumptionSummaryItem{
				{GroupKey: "user-1", TotalTokens: 100, TotalRecords: 1},
			},
		},
	}
	reader := &quotaDashboardSummaryReaderForAPITest{
		err: errors.New("ai service down"),
	}
	router := setupQuotaDashboardAPITestRouter(repo, &quotaDashboardAlertRepoForAPITest{}, reader, nil)

	req, _ := http.NewRequest("GET", "/statistics/reconciliation?timeRange=day&groupBy=user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取对账摘要失败")
	assert.Contains(t, w.Body.String(), "获取AI服务消费聚合摘要失败")
}
