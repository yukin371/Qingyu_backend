package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"
	pb "Qingyu_backend/pkg/grpc/pb"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type quotaAdminAPITestRepo struct {
	listItems []*aiModels.UserQuotaListItem
	listTotal int64

	lastListRole   string
	lastListStatus string
	lastListSearch string
	lastListPage   int
	lastListLimit  int

	quotasByKey   map[string]*aiModels.UserQuota
	quotasByUser  map[string][]*aiModels.UserQuota
	getAllErr     error
	transactions  []*aiModels.QuotaTransaction
	createdQuotas []*aiModels.UserQuota
	updatedQuotas []*aiModels.UserQuota
	lastTxUserID  string
	lastTxStart   time.Time
	lastTxEnd     time.Time
}

func newQuotaAdminAPITestRepo() *quotaAdminAPITestRepo {
	return &quotaAdminAPITestRepo{
		quotasByKey:  make(map[string]*aiModels.UserQuota),
		quotasByUser: make(map[string][]*aiModels.UserQuota),
	}
}

func adminQuotaKey(userID string, quotaType aiModels.QuotaType) string {
	return userID + "|" + string(quotaType)
}

func (s *quotaAdminAPITestRepo) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	s.createdQuotas = append(s.createdQuotas, quota)
	s.quotasByKey[adminQuotaKey(quota.UserID, quota.QuotaType)] = quota
	return nil
}

func (s *quotaAdminAPITestRepo) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	if quota, ok := s.quotasByKey[adminQuotaKey(userID, quotaType)]; ok {
		return quota, nil
	}
	return nil, aiModels.ErrQuotaNotFound
}

func (s *quotaAdminAPITestRepo) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	s.updatedQuotas = append(s.updatedQuotas, quota)
	s.quotasByKey[adminQuotaKey(quota.UserID, quota.QuotaType)] = quota
	return nil
}

func (s *quotaAdminAPITestRepo) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	delete(s.quotasByKey, adminQuotaKey(userID, quotaType))
	return nil
}

func (s *quotaAdminAPITestRepo) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	if s.getAllErr != nil {
		return nil, s.getAllErr
	}
	return s.quotasByUser[userID], nil
}

func (s *quotaAdminAPITestRepo) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	return nil
}

func (s *quotaAdminAPITestRepo) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	s.transactions = append(s.transactions, transaction)
	return nil
}

func (s *quotaAdminAPITestRepo) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	return nil, nil
}

func (s *quotaAdminAPITestRepo) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	s.lastTxUserID = userID
	s.lastTxStart = startTime
	s.lastTxEnd = endTime

	result := make([]*aiModels.QuotaTransaction, 0, len(s.transactions))
	for _, tx := range s.transactions {
		if tx == nil || tx.UserID != userID {
			continue
		}
		if tx.Timestamp.Before(startTime) || tx.Timestamp.After(endTime) {
			continue
		}
		result = append(result, tx)
	}
	return result, nil
}

func (s *quotaAdminAPITestRepo) GetQuotaStatistics(ctx context.Context, userID string) (*aiRepo.QuotaStatistics, error) {
	return nil, nil
}

func (s *quotaAdminAPITestRepo) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	return 0, nil
}

func (s *quotaAdminAPITestRepo) GetDashboardSummary(ctx context.Context) (*aiModels.DashboardSummary, error) {
	return nil, nil
}

func (s *quotaAdminAPITestRepo) GetQuotaDistribution(ctx context.Context) (*aiModels.QuotaDistribution, error) {
	return nil, nil
}

func (s *quotaAdminAPITestRepo) GetTopConsumers(ctx context.Context, limit int) ([]aiModels.UserQuotaRanking, error) {
	return nil, nil
}

func (s *quotaAdminAPITestRepo) GetConsumptionTrend(ctx context.Context, days int) ([]aiModels.TrendPoint, error) {
	return nil, nil
}

func (s *quotaAdminAPITestRepo) GetConsumptionSummary(ctx context.Context, startTime, endTime time.Time, workflowType, groupBy string, page, pageSize int) (*aiModels.QuotaConsumptionSummary, error) {
	return nil, nil
}

func (s *quotaAdminAPITestRepo) ListUserQuotas(ctx context.Context, role, status, search string, page, limit int) ([]*aiModels.UserQuotaListItem, int64, error) {
	s.lastListRole = role
	s.lastListStatus = status
	s.lastListSearch = search
	s.lastListPage = page
	s.lastListLimit = limit
	return s.listItems, s.listTotal, nil
}

func (s *quotaAdminAPITestRepo) Health(ctx context.Context) error {
	return nil
}

var _ aiRepo.QuotaRepository = (*quotaAdminAPITestRepo)(nil)

type noopQuotaAlertRepo struct{}

func (s *noopQuotaAlertRepo) Create(ctx context.Context, alert *aiModels.QuotaAlert) error {
	return nil
}
func (s *noopQuotaAlertRepo) GetByID(ctx context.Context, id string) (*aiModels.QuotaAlert, error) {
	return nil, errors.New("not found")
}
func (s *noopQuotaAlertRepo) List(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	return nil, 0, nil
}
func (s *noopQuotaAlertRepo) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	return nil, 0, nil
}
func (s *noopQuotaAlertRepo) Update(ctx context.Context, alert *aiModels.QuotaAlert) error {
	return nil
}
func (s *noopQuotaAlertRepo) GetRecentGlobal(ctx context.Context, limit int) ([]*aiModels.QuotaAlert, error) {
	return nil, nil
}
func (s *noopQuotaAlertRepo) CountByStatus(ctx context.Context) (map[aiModels.QuotaAlertStatus]int64, error) {
	return map[aiModels.QuotaAlertStatus]int64{}, nil
}
func (s *noopQuotaAlertRepo) Health(ctx context.Context) error { return nil }

var _ aiRepo.QuotaAlertRepository = (*noopQuotaAlertRepo)(nil)

type quotaConsumptionReaderStub struct {
	response         *pb.QuotaConsumptionResponse
	err              error
	lastUserID       string
	lastTimeRange    string
	lastWorkflowType string
}

func (s *quotaConsumptionReaderStub) GetQuotaConsumption(ctx context.Context, userID string, timeRange string, workflowType string) (*pb.QuotaConsumptionResponse, error) {
	s.lastUserID = userID
	s.lastTimeRange = timeRange
	s.lastWorkflowType = workflowType
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}

func setupQuotaAdminAPITestHarness(repo *quotaAdminAPITestRepo, reader *quotaConsumptionReaderStub) (*QuotaAdminAPI, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	quotaService := aiService.NewQuotaService(repo)
	adminService := aiService.NewQuotaAdminService(repo, &noopQuotaAlertRepo{})
	if reader != nil {
		adminService.SetConsumptionReader(reader)
	}
	api := NewQuotaAdminAPI(quotaService, adminService)

	router := gin.New()
	router.GET("/users", api.ListUserQuotas)
	router.GET("/users/:userId", api.GetUserQuotaDetails)
	router.GET("/users/:userId/reconciliation", api.GetUserQuotaReconciliation)
	router.PUT("/users/:userId", api.UpdateUserQuota)
	router.POST("/users/:userId/recharge", api.RechargeUserQuota)
	router.POST("/users/:userId/suspend", api.SuspendUserQuota)
	router.POST("/users/:userId/activate", api.ActivateUserQuota)
	router.POST("/batch-recharge", api.BatchRecharge)
	router.POST("/batch-update", api.BatchUpdateQuota)
	router.POST("/batch-suspend", api.BatchSuspend)
	router.POST("/batch-activate", api.BatchActivate)
	return api, router
}

func TestQuotaAdminAPIListUserQuotasNormalizesPagingAndForwardsFilters(t *testing.T) {
	repo := newQuotaAdminAPITestRepo()
	repo.listItems = []*aiModels.UserQuotaListItem{
		{UserID: "user-1", Username: "user-1", Role: "writer", MemberLevel: "normal", DailyQuota: 120, DailyUsed: 40, UsagePercent: 33.3, Status: "active"},
	}
	repo.listTotal = 1
	_, router := setupQuotaAdminAPITestHarness(repo, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users?role=writer&status=active&search=needle&page=0&limit=999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "writer", repo.lastListRole)
	assert.Equal(t, "active", repo.lastListStatus)
	assert.Equal(t, "needle", repo.lastListSearch)
	assert.Equal(t, 1, repo.lastListPage)
	assert.Equal(t, 20, repo.lastListLimit)

	var resp struct {
		Page  int              `json:"page"`
		Size  int              `json:"size"`
		Total int64            `json:"total"`
		Data  []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.Size)
	assert.Equal(t, int64(1), resp.Total)
	require.Len(t, resp.Data, 1)
}

func TestQuotaAdminAPIGetUserQuotaDetailsSurfacesServiceError(t *testing.T) {
	t.Run("rejects missing user ID for details", func(t *testing.T) {
		repo := newQuotaAdminAPITestRepo()
		api, _ := setupQuotaAdminAPITestHarness(repo, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/users", nil)

		api.GetUserQuotaDetails(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "用户ID不能为空")
	})

	t.Run("surfaces service error when quota details load fails", func(t *testing.T) {
		repo := newQuotaAdminAPITestRepo()
		repo.getAllErr = errors.New("quota service unavailable")
		_, router := setupQuotaAdminAPITestHarness(repo, nil)

		req, _ := http.NewRequest(http.MethodGet, "/users/user-1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "获取配额详情失败")
		assert.Contains(t, w.Body.String(), "quota service unavailable")
	})

	t.Run("returns details and reconciliation summary", func(t *testing.T) {
		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		elapsed := now.Sub(dayStart)
		firstConsumeAt := dayStart.Add(elapsed / 4)
		secondConsumeAt := dayStart.Add(elapsed / 2)
		otherWorkflowAt := dayStart.Add((elapsed * 3) / 4)

		repo := newQuotaAdminAPITestRepo()
		repo.quotasByUser["user-1"] = []*aiModels.UserQuota{
			{
				UserID:         "user-1",
				QuotaType:      aiModels.QuotaTypeDaily,
				TotalQuota:     120,
				RemainingQuota: 80,
			},
			{
				UserID:         "user-1",
				QuotaType:      aiModels.QuotaTypeMonthly,
				TotalQuota:     300,
				RemainingQuota: 280,
			},
		}
		repo.transactions = []*aiModels.QuotaTransaction{
			{UserID: "user-1", Type: "consume", Service: "story_write", Amount: 120, Timestamp: firstConsumeAt},
			{UserID: "user-1", Type: "consume", Service: "story_write", Amount: 80, Timestamp: secondConsumeAt},
			{UserID: "user-1", Type: "consume", Service: "outline", Amount: 50, Timestamp: otherWorkflowAt},
			{UserID: "user-1", Type: "restore", Service: "story_write", Amount: -20, Timestamp: now},
		}
		reader := &quotaConsumptionReaderStub{
			response: &pb.QuotaConsumptionResponse{
				Success:      true,
				TotalTokens:  150,
				TotalRecords: 2,
				Records: []*pb.QuotaRecord{
					{Id: "r1", WorkflowType: "story_write", TokensUsed: 70, ConsumedAt: firstConsumeAt.Format(time.RFC3339)},
					{Id: "r2", WorkflowType: "story_write", TokensUsed: 80, ConsumedAt: secondConsumeAt.Format(time.RFC3339)},
				},
			},
		}
		_, router := setupQuotaAdminAPITestHarness(repo, reader)

		detailsReq, _ := http.NewRequest(http.MethodGet, "/users/user-1", nil)
		detailsResp := httptest.NewRecorder()
		router.ServeHTTP(detailsResp, detailsReq)

		require.Equal(t, http.StatusOK, detailsResp.Code)
		var detailsResult struct {
			Data []struct {
				QuotaType      string `json:"quotaType"`
				TotalQuota     int    `json:"totalQuota"`
				RemainingQuota int    `json:"remainingQuota"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(detailsResp.Body.Bytes(), &detailsResult))
		require.Len(t, detailsResult.Data, 2)
		assert.Equal(t, string(aiModels.QuotaTypeDaily), detailsResult.Data[0].QuotaType)
		assert.Equal(t, 120, detailsResult.Data[0].TotalQuota)
		assert.Equal(t, string(aiModels.QuotaTypeMonthly), detailsResult.Data[1].QuotaType)
		assert.Equal(t, 280, detailsResult.Data[1].RemainingQuota)

		reconciliationReq, _ := http.NewRequest(http.MethodGet, "/users/user-1/reconciliation?workflowType=story_write", nil)
		reconciliationResp := httptest.NewRecorder()
		router.ServeHTTP(reconciliationResp, reconciliationReq)

		require.Equal(t, http.StatusOK, reconciliationResp.Code)
		assert.Equal(t, "user-1", repo.lastTxUserID)
		assert.Equal(t, "day", reader.lastTimeRange)
		assert.Equal(t, "story_write", reader.lastWorkflowType)

		var reconciliationResult struct {
			Data struct {
				TimeRange            string `json:"timeRange"`
				WorkflowType         string `json:"workflowType"`
				BackendTotalTokens   int    `json:"backendTotalTokens"`
				BackendRecordCount   int    `json:"backendRecordCount"`
				AIServiceTotalTokens int    `json:"aiServiceTotalTokens"`
				AIServiceRecordCount int    `json:"aiServiceRecordCount"`
				DifferenceTokens     int    `json:"differenceTokens"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(reconciliationResp.Body.Bytes(), &reconciliationResult))
		assert.Equal(t, "day", reconciliationResult.Data.TimeRange)
		assert.Equal(t, "story_write", reconciliationResult.Data.WorkflowType)
		assert.Equal(t, 200, reconciliationResult.Data.BackendTotalTokens)
		assert.Equal(t, 2, reconciliationResult.Data.BackendRecordCount)
		assert.Equal(t, 150, reconciliationResult.Data.AIServiceTotalTokens)
		assert.Equal(t, 2, reconciliationResult.Data.AIServiceRecordCount)
		assert.Equal(t, 50, reconciliationResult.Data.DifferenceTokens)
	})

	t.Run("returns empty reconciliation records when ai service has no data", func(t *testing.T) {
		repo := newQuotaAdminAPITestRepo()
		repo.transactions = []*aiModels.QuotaTransaction{}
		reader := &quotaConsumptionReaderStub{
			response: &pb.QuotaConsumptionResponse{Success: true},
		}
		_, router := setupQuotaAdminAPITestHarness(repo, reader)

		req, _ := http.NewRequest(http.MethodGet, "/users/user-1/reconciliation", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Data struct {
				BackendTotalTokens   int `json:"backendTotalTokens"`
				BackendRecordCount   int `json:"backendRecordCount"`
				AIServiceTotalTokens int `json:"aiServiceTotalTokens"`
				AIServiceRecordCount int `json:"aiServiceRecordCount"`
				DifferenceTokens     int `json:"differenceTokens"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, 0, resp.Data.BackendTotalTokens)
		assert.Equal(t, 0, resp.Data.BackendRecordCount)
		assert.Equal(t, 0, resp.Data.AIServiceTotalTokens)
		assert.Equal(t, 0, resp.Data.AIServiceRecordCount)
		assert.Equal(t, 0, resp.Data.DifferenceTokens)
	})

	t.Run("maps unsuccessful ai response without message to internal error", func(t *testing.T) {
		repo := newQuotaAdminAPITestRepo()
		reader := &quotaConsumptionReaderStub{
			response: &pb.QuotaConsumptionResponse{Success: false},
		}
		_, router := setupQuotaAdminAPITestHarness(repo, reader)

		req, _ := http.NewRequest(http.MethodGet, "/users/user-1/reconciliation", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "获取配额对账结果失败")
		assert.Contains(t, w.Body.String(), "AI服务返回未成功状态")
	})
}

func TestQuotaAdminAPIGetUserQuotaReconciliationRejectsMissingUserID(t *testing.T) {
	repo := newQuotaAdminAPITestRepo()
	_, _ = setupQuotaAdminAPITestHarness(repo, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users/reconciliation", nil)

	api := NewQuotaAdminAPI(aiService.NewQuotaService(repo), aiService.NewQuotaAdminService(repo, &noopQuotaAlertRepo{}))
	api.GetUserQuotaReconciliation(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "用户ID不能为空")
}

func TestQuotaAdminAPIRechargeUserQuotaRejectsMissingUserID(t *testing.T) {
	repo := newQuotaAdminAPITestRepo()
	_, _ = setupQuotaAdminAPITestHarness(repo, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users/recharge", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	api := NewQuotaAdminAPI(aiService.NewQuotaService(repo), aiService.NewQuotaAdminService(repo, &noopQuotaAlertRepo{}))
	api.RechargeUserQuota(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "用户ID不能为空")
}

func TestQuotaAdminAPIUpdateSuspendAndActivateHandlers(t *testing.T) {
	t.Run("rejects missing user ID for update", func(t *testing.T) {
		repo := newQuotaAdminAPITestRepo()
		api, _ := setupQuotaAdminAPITestHarness(repo, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(`{"quotaType":"daily","totalQuota":100}`))
		c.Request.Header.Set("Content-Type", "application/json")

		api.UpdateUserQuota(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "用户ID不能为空")
	})

	t.Run("maps supported quota types to the quota service", func(t *testing.T) {
		cases := []struct {
			name         string
			userID       string
			quotaType    string
			expectedType aiModels.QuotaType
		}{
			{name: "daily", userID: "user-daily", quotaType: "daily", expectedType: aiModels.QuotaTypeDaily},
			{name: "monthly", userID: "user-monthly", quotaType: "monthly", expectedType: aiModels.QuotaTypeMonthly},
			{name: "total", userID: "user-total", quotaType: "total", expectedType: aiModels.QuotaTypeTotal},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				repo := newQuotaAdminAPITestRepo()
				api, router := setupQuotaAdminAPITestHarness(repo, nil)

				req, _ := http.NewRequest(
					http.MethodPut,
					"/users/"+tt.userID,
					strings.NewReader(`{"quotaType":"`+tt.quotaType+`","totalQuota":100}`),
				)
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, req)

				require.Equal(t, http.StatusOK, resp.Code)
				require.Len(t, repo.createdQuotas, 1)
				created := repo.createdQuotas[0]
				assert.Equal(t, tt.expectedType, created.QuotaType)
				assert.Equal(t, 100, created.TotalQuota)
				switch tt.expectedType {
				case aiModels.QuotaTypeDaily, aiModels.QuotaTypeMonthly:
					assert.False(t, created.ResetAt.IsZero())
				case aiModels.QuotaTypeTotal:
					assert.True(t, created.ResetAt.IsZero())
				}

				_ = api
			})
		}
	})

	t.Run("suspends and activates existing quotas", func(t *testing.T) {
		repo := newQuotaAdminAPITestRepo()
		repo.quotasByUser["user-toggle"] = []*aiModels.UserQuota{
			{
				UserID:    "user-toggle",
				QuotaType: aiModels.QuotaTypeDaily,
				Status:    aiModels.QuotaStatusActive,
			},
			{
				UserID:    "user-toggle",
				QuotaType: aiModels.QuotaTypeMonthly,
				Status:    aiModels.QuotaStatusActive,
			},
		}
		_, router := setupQuotaAdminAPITestHarness(repo, nil)

		suspendReq, _ := http.NewRequest(http.MethodPost, "/users/user-toggle/suspend", nil)
		suspendResp := httptest.NewRecorder()
		router.ServeHTTP(suspendResp, suspendReq)
		require.Equal(t, http.StatusOK, suspendResp.Code)
		require.Len(t, repo.updatedQuotas, 2)
		assert.Equal(t, aiModels.QuotaStatusSuspended, repo.updatedQuotas[0].Status)
		assert.Equal(t, aiModels.QuotaStatusSuspended, repo.updatedQuotas[1].Status)

		repo.updatedQuotas = nil
		for _, quota := range repo.quotasByUser["user-toggle"] {
			quota.Status = aiModels.QuotaStatusSuspended
		}

		activateReq, _ := http.NewRequest(http.MethodPost, "/users/user-toggle/activate", nil)
		activateResp := httptest.NewRecorder()
		router.ServeHTTP(activateResp, activateReq)
		require.Equal(t, http.StatusOK, activateResp.Code)
		require.Len(t, repo.updatedQuotas, 2)
		assert.Equal(t, aiModels.QuotaStatusActive, repo.updatedQuotas[0].Status)
		assert.Equal(t, aiModels.QuotaStatusActive, repo.updatedQuotas[1].Status)
	})
}

func TestQuotaAdminAPINegativeHandlersRejectInvalidPayloads(t *testing.T) {
	repo := newQuotaAdminAPITestRepo()
	_, router := setupQuotaAdminAPITestHarness(repo, nil)

	t.Run("rejects invalid quota type", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/users/user-1", strings.NewReader(`{"quotaType":"yearly","totalQuota":100}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "oneof")
	})

	t.Run("rejects malformed update payload", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/users/user-1", strings.NewReader("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "参数错误")
	})

	t.Run("rejects malformed recharge and batch payloads", func(t *testing.T) {
		cases := []struct {
			name   string
			method string
			path   string
		}{
			{name: "recharge", method: http.MethodPost, path: "/users/user-1/recharge"},
			{name: "batch-recharge", method: http.MethodPost, path: "/batch-recharge"},
			{name: "batch-update", method: http.MethodPost, path: "/batch-update"},
			{name: "batch-suspend", method: http.MethodPost, path: "/batch-suspend"},
			{name: "batch-activate", method: http.MethodPost, path: "/batch-activate"},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest(tt.method, tt.path, strings.NewReader("{"))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				require.Equal(t, http.StatusBadRequest, w.Code)
				assert.Contains(t, w.Body.String(), "参数错误")
			})
		}
	})
}
