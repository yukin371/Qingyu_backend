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
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type quotaAlertRepoForAPITest struct {
	alerts        []*aiModels.QuotaAlert
	lastAlertType string
	lastLevel     string
	lastStatus    string
	lastPage      int
	lastLimit     int
	getByIDErr    error
	listErr       error
}

func (s *quotaAlertRepoForAPITest) Create(ctx context.Context, alert *aiModels.QuotaAlert) error {
	return nil
}

func (s *quotaAlertRepoForAPITest) GetByID(ctx context.Context, id string) (*aiModels.QuotaAlert, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	for _, alert := range s.alerts {
		if alert != nil && alert.ID.Hex() == id {
			return alert, nil
		}
	}
	return nil, nil
}

func (s *quotaAlertRepoForAPITest) List(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	s.lastAlertType = alertType
	s.lastLevel = level
	s.lastStatus = status
	s.lastPage = page
	s.lastLimit = limit
	if s.listErr != nil {
		return nil, 0, s.listErr
	}

	items := make([]*aiModels.QuotaAlert, 0, len(s.alerts))
	for _, alert := range s.alerts {
		if alert == nil {
			continue
		}
		if alertType != "" && string(alert.Type) != alertType {
			continue
		}
		if level != "" && string(alert.Level) != level {
			continue
		}
		switch status {
		case "", "all":
		case "open":
			if alert.Status != aiModels.QuotaAlertStatusPending && alert.Status != aiModels.QuotaAlertStatusAcknowledged {
				continue
			}
		default:
			if string(alert.Status) != status {
				continue
			}
		}
		items = append(items, alert)
	}
	return items, int64(len(items)), nil
}

func (s *quotaAlertRepoForAPITest) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	return nil, 0, nil
}

func (s *quotaAlertRepoForAPITest) Update(ctx context.Context, alert *aiModels.QuotaAlert) error {
	return nil
}

func (s *quotaAlertRepoForAPITest) GetRecentGlobal(ctx context.Context, limit int) ([]*aiModels.QuotaAlert, error) {
	return nil, nil
}

func (s *quotaAlertRepoForAPITest) CountByStatus(ctx context.Context) (map[aiModels.QuotaAlertStatus]int64, error) {
	return map[aiModels.QuotaAlertStatus]int64{}, nil
}

func (s *quotaAlertRepoForAPITest) Health(ctx context.Context) error {
	return nil
}

var _ aiRepo.QuotaAlertRepository = (*quotaAlertRepoForAPITest)(nil)

func setupQuotaAlertAPITestRouter(alerts ...*aiModels.QuotaAlert) *gin.Engine {
	gin.SetMode(gin.TestMode)
	repo := &quotaAlertRepoForAPITest{alerts: alerts}
	service := aiService.NewQuotaAlertService(repo)
	api := NewQuotaAlertAPI(service)

	router := gin.New()
	router.GET("/alerts", api.ListAlerts)
	router.GET("/alerts/:id", api.GetAlert)
	router.PUT("/alerts/:id/acknowledge", api.AcknowledgeAlert)
	router.PUT("/alerts/:id/resolve", api.ResolveAlert)
	router.PUT("/alerts/:id/ignore", api.IgnoreAlert)
	return router
}

func setupQuotaAlertAPITestRouterWithRepo(repo *quotaAlertRepoForAPITest) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := aiService.NewQuotaAlertService(repo)
	api := NewQuotaAlertAPI(service)

	router := gin.New()
	router.GET("/alerts", api.ListAlerts)
	router.GET("/alerts/:id", api.GetAlert)
	router.PUT("/alerts/:id/acknowledge", api.AcknowledgeAlert)
	router.PUT("/alerts/:id/resolve", api.ResolveAlert)
	router.PUT("/alerts/:id/ignore", api.IgnoreAlert)
	return router
}

func TestQuotaAlertAPIListAlertsWithoutStatusKeepsHistoryView(t *testing.T) {
	router := setupQuotaAlertAPITestRouter(
		&aiModels.QuotaAlert{
			ID:        primitive.NewObjectID(),
			Type:      aiModels.QuotaAlertTypeConsistency,
			Status:    aiModels.QuotaAlertStatusPending,
			Title:     "pending",
			Message:   "pending",
			Level:     aiModels.QuotaAlertLevelWarning,
			CreatedAt: primitive.NewObjectID().Timestamp(),
		},
		&aiModels.QuotaAlert{
			ID:        primitive.NewObjectID(),
			Type:      aiModels.QuotaAlertTypeConsistency,
			Status:    aiModels.QuotaAlertStatusResolved,
			Title:     "resolved",
			Message:   "resolved",
			Level:     aiModels.QuotaAlertLevelInfo,
			CreatedAt: primitive.NewObjectID().Timestamp(),
		},
	)

	req, _ := http.NewRequest("GET", "/alerts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Code  int              `json:"code"`
		Total int              `json:"total"`
		Data  []map[string]any `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 200, response.Code)
	assert.Equal(t, 2, response.Total)
	require.Len(t, response.Data, 2)
}

func TestQuotaAlertAPIListAlertsSupportsAllStatus(t *testing.T) {
	router := setupQuotaAlertAPITestRouter(
		&aiModels.QuotaAlert{
			ID:        primitive.NewObjectID(),
			Type:      aiModels.QuotaAlertTypeConsistency,
			Status:    aiModels.QuotaAlertStatusPending,
			Title:     "pending",
			Message:   "pending",
			Level:     aiModels.QuotaAlertLevelWarning,
			CreatedAt: primitive.NewObjectID().Timestamp(),
		},
		&aiModels.QuotaAlert{
			ID:        primitive.NewObjectID(),
			Type:      aiModels.QuotaAlertTypeConsistency,
			Status:    aiModels.QuotaAlertStatusResolved,
			Title:     "resolved",
			Message:   "resolved",
			Level:     aiModels.QuotaAlertLevelInfo,
			CreatedAt: primitive.NewObjectID().Timestamp(),
		},
	)

	req, _ := http.NewRequest("GET", "/alerts?status=all", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Code  int              `json:"code"`
		Total int              `json:"total"`
		Data  []map[string]any `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 200, response.Code)
	assert.Equal(t, 2, response.Total)
	require.Len(t, response.Data, 2)
}

func TestQuotaAlertAPIListAlertsReturnsEmptyPageWhenNoAlerts(t *testing.T) {
	router := setupQuotaAlertAPITestRouterWithRepo(&quotaAlertRepoForAPITest{})

	req, _ := http.NewRequest("GET", "/alerts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Code  int              `json:"code"`
		Total int              `json:"total"`
		Data  []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, 0, response.Total)
	require.Len(t, response.Data, 0)
}

func TestQuotaAlertAPIListAlertsNormalizesPagingAndFiltersOpenStatus(t *testing.T) {
	repo := &quotaAlertRepoForAPITest{
		alerts: []*aiModels.QuotaAlert{
			{
				ID:        primitive.NewObjectID(),
				Type:      aiModels.QuotaAlertTypeConsistency,
				Status:    aiModels.QuotaAlertStatusPending,
				Title:     "pending",
				Message:   "pending",
				Level:     aiModels.QuotaAlertLevelWarning,
				CreatedAt: primitive.NewObjectID().Timestamp(),
			},
			{
				ID:        primitive.NewObjectID(),
				Type:      aiModels.QuotaAlertTypeConsistency,
				Status:    aiModels.QuotaAlertStatusAcknowledged,
				Title:     "ack",
				Message:   "ack",
				Level:     aiModels.QuotaAlertLevelWarning,
				CreatedAt: primitive.NewObjectID().Timestamp(),
			},
			{
				ID:        primitive.NewObjectID(),
				Type:      aiModels.QuotaAlertTypeConsistency,
				Status:    aiModels.QuotaAlertStatusResolved,
				Title:     "resolved",
				Message:   "resolved",
				Level:     aiModels.QuotaAlertLevelWarning,
				CreatedAt: primitive.NewObjectID().Timestamp(),
			},
		},
	}
	router := setupQuotaAlertAPITestRouterWithRepo(repo)

	req, _ := http.NewRequest("GET", "/alerts?type=consistency&level=warning&status=open&page=0&limit=101", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "consistency", repo.lastAlertType)
	assert.Equal(t, "warning", repo.lastLevel)
	assert.Equal(t, "open", repo.lastStatus)
	assert.Equal(t, 1, repo.lastPage)
	assert.Equal(t, 20, repo.lastLimit)

	var response struct {
		Code  int              `json:"code"`
		Total int              `json:"total"`
		Data  []map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, 2, response.Total)
	require.Len(t, response.Data, 2)
}

func TestQuotaAlertAPIListAlertsSurfacesRepositoryFailure(t *testing.T) {
	repo := &quotaAlertRepoForAPITest{
		listErr: errors.New("list failed"),
	}
	router := setupQuotaAlertAPITestRouterWithRepo(repo)

	req, _ := http.NewRequest("GET", "/alerts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取告警列表失败")
}

func TestQuotaAlertAPIGetAlertReturnsDetails(t *testing.T) {
	alert := &aiModels.QuotaAlert{
		ID:        primitive.NewObjectID(),
		Type:      aiModels.QuotaAlertTypeConsistency,
		Status:    aiModels.QuotaAlertStatusPending,
		Title:     "pending",
		Message:   "pending",
		Level:     aiModels.QuotaAlertLevelWarning,
		CreatedAt: primitive.NewObjectID().Timestamp(),
	}
	router := setupQuotaAlertAPITestRouter(alert)

	req, _ := http.NewRequest("GET", "/alerts/"+alert.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Title  string `json:"title"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, alert.ID.Hex(), response.Data.ID)
	assert.Equal(t, string(aiModels.QuotaAlertStatusPending), response.Data.Status)
	assert.Equal(t, "pending", response.Data.Title)
}

func TestQuotaAlertAPIGetAlertSurfacesRepositoryFailure(t *testing.T) {
	repo := &quotaAlertRepoForAPITest{
		getByIDErr: errors.New("lookup failed"),
	}
	router := setupQuotaAlertAPITestRouterWithRepo(repo)

	req, _ := http.NewRequest("GET", "/alerts/507f1f77bcf86cd799439011", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取告警失败")
	assert.Contains(t, w.Body.String(), "lookup failed")
}

func TestQuotaAlertAPIResolveAndIgnoreEndpointsUpdateAlertStatus(t *testing.T) {
	resolvedAlert := &aiModels.QuotaAlert{
		ID:        primitive.NewObjectID(),
		Type:      aiModels.QuotaAlertTypeConsistency,
		Status:    aiModels.QuotaAlertStatusPending,
		Title:     "resolve",
		Message:   "resolve",
		Level:     aiModels.QuotaAlertLevelWarning,
		CreatedAt: time.Now(),
	}
	ignoredAlert := &aiModels.QuotaAlert{
		ID:        primitive.NewObjectID(),
		Type:      aiModels.QuotaAlertTypeConsistency,
		Status:    aiModels.QuotaAlertStatusPending,
		Title:     "ignore",
		Message:   "ignore",
		Level:     aiModels.QuotaAlertLevelWarning,
		CreatedAt: time.Now(),
	}
	router := setupQuotaAlertAPITestRouter(resolvedAlert, ignoredAlert)

	resolveReq, _ := http.NewRequest("PUT", "/alerts/"+resolvedAlert.ID.Hex()+"/resolve", strings.NewReader(`{"operatorId":"admin-1"}`))
	resolveReq.Header.Set("Content-Type", "application/json")
	resolveResp := httptest.NewRecorder()
	router.ServeHTTP(resolveResp, resolveReq)
	require.Equal(t, http.StatusOK, resolveResp.Code)
	assert.Equal(t, aiModels.QuotaAlertStatusResolved, resolvedAlert.Status)
	assert.Equal(t, "admin-1", resolvedAlert.ResolvedBy)
	assert.NotNil(t, resolvedAlert.ResolvedAt)

	ignoreReq, _ := http.NewRequest("PUT", "/alerts/"+ignoredAlert.ID.Hex()+"/ignore", strings.NewReader(`{"operatorId":"admin-2"}`))
	ignoreReq.Header.Set("Content-Type", "application/json")
	ignoreResp := httptest.NewRecorder()
	router.ServeHTTP(ignoreResp, ignoreReq)
	require.Equal(t, http.StatusOK, ignoreResp.Code)
	assert.Equal(t, aiModels.QuotaAlertStatusIgnored, ignoredAlert.Status)
	assert.Equal(t, "admin-2", ignoredAlert.ResolvedBy)
}

func TestQuotaAlertAPIActionHandlersRejectMissingIDAndMalformedJSON(t *testing.T) {
	t.Run("rejects missing alert id", func(t *testing.T) {
		api := NewQuotaAlertAPI(aiService.NewQuotaAlertService(&quotaAlertRepoForAPITest{}))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/alerts", nil)

		api.GetAlert(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "告警ID不能为空")
	})

	t.Run("rejects malformed action payloads", func(t *testing.T) {
		router := setupQuotaAlertAPITestRouter()
		cases := []struct {
			name   string
			method string
			path   string
		}{
			{name: "acknowledge", method: http.MethodPut, path: "/alerts/507f1f77bcf86cd799439011/acknowledge"},
			{name: "resolve", method: http.MethodPut, path: "/alerts/507f1f77bcf86cd799439011/resolve"},
			{name: "ignore", method: http.MethodPut, path: "/alerts/507f1f77bcf86cd799439011/ignore"},
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

func TestQuotaAlertAPIListAlertsSurfacesRepositoryErrors(t *testing.T) {
	repo := &quotaAlertRepoForAPITest{
		listErr: errors.New("query failed"),
	}
	router := setupQuotaAlertAPITestRouterWithRepo(repo)

	req, _ := http.NewRequest("GET", "/alerts?status=open", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "获取告警列表失败")
	assert.Contains(t, w.Body.String(), "query failed")
}
