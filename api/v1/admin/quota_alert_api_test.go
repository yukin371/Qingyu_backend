package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	aiModels "Qingyu_backend/models/ai"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	aiService "Qingyu_backend/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type quotaAlertRepoForAPITest struct {
	alerts []*aiModels.QuotaAlert
}

func (s *quotaAlertRepoForAPITest) Create(ctx context.Context, alert *aiModels.QuotaAlert) error {
	return nil
}

func (s *quotaAlertRepoForAPITest) GetByID(ctx context.Context, id string) (*aiModels.QuotaAlert, error) {
	for _, alert := range s.alerts {
		if alert != nil && alert.ID.Hex() == id {
			return alert, nil
		}
	}
	return nil, nil
}

func (s *quotaAlertRepoForAPITest) List(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
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
