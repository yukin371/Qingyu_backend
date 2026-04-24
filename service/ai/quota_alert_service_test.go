package ai

import (
	"context"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type quotaAlertRepoStub struct {
	alerts map[string]*aiModels.QuotaAlert
}

func newQuotaAlertRepoStub(alerts ...*aiModels.QuotaAlert) *quotaAlertRepoStub {
	items := make(map[string]*aiModels.QuotaAlert, len(alerts))
	for _, alert := range alerts {
		if alert == nil {
			continue
		}
		if alert.ID.IsZero() {
			alert.ID = primitive.NewObjectID()
		}
		items[alert.ID.Hex()] = alert
	}
	return &quotaAlertRepoStub{alerts: items}
}

func (s *quotaAlertRepoStub) Create(ctx context.Context, alert *aiModels.QuotaAlert) error {
	if alert.ID.IsZero() {
		alert.ID = primitive.NewObjectID()
	}
	s.alerts[alert.ID.Hex()] = alert
	return nil
}

func (s *quotaAlertRepoStub) GetByID(ctx context.Context, id string) (*aiModels.QuotaAlert, error) {
	return s.alerts[id], nil
}

func (s *quotaAlertRepoStub) List(ctx context.Context, alertType, level, status string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	items := make([]*aiModels.QuotaAlert, 0, len(s.alerts))
	for _, alert := range s.alerts {
		if alertType != "" && string(alert.Type) != alertType {
			continue
		}
		if level != "" && string(alert.Level) != level {
			continue
		}
		if status != "" && string(alert.Status) != status {
			continue
		}
		items = append(items, alert)
	}
	return items, int64(len(items)), nil
}

func (s *quotaAlertRepoStub) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*aiModels.QuotaAlert, int64, error) {
	items := make([]*aiModels.QuotaAlert, 0)
	for _, alert := range s.alerts {
		if alert.UserID == userID {
			items = append(items, alert)
		}
	}
	return items, int64(len(items)), nil
}

func (s *quotaAlertRepoStub) Update(ctx context.Context, alert *aiModels.QuotaAlert) error {
	s.alerts[alert.ID.Hex()] = alert
	return nil
}

func (s *quotaAlertRepoStub) GetRecentGlobal(ctx context.Context, limit int) ([]*aiModels.QuotaAlert, error) {
	return nil, nil
}

func (s *quotaAlertRepoStub) CountByStatus(ctx context.Context) (map[aiModels.QuotaAlertStatus]int64, error) {
	return map[aiModels.QuotaAlertStatus]int64{}, nil
}

func (s *quotaAlertRepoStub) Health(ctx context.Context) error {
	return nil
}

func TestResolveRecoveredConsistencyAlertsResolvesOnlyCheckedRecoveredOpenAlerts(t *testing.T) {
	now := time.Now()
	workflowAlert := &aiModels.QuotaAlert{
		ID:      primitive.NewObjectID(),
		Type:    aiModels.QuotaAlertTypeConsistency,
		Level:   aiModels.QuotaAlertLevelCritical,
		Title:   "workflow mismatch",
		Status:  aiModels.QuotaAlertStatusPending,
		Message: "workflow mismatch",
		Data: map[string]interface{}{
			"scope":     quotaConsistencyScopeWorkflow,
			"timeRange": "day",
			"groupBy":   "workflow",
			"groupKey":  "story_write",
		},
		CreatedAt: now,
	}
	globalAlert := &aiModels.QuotaAlert{
		ID:      primitive.NewObjectID(),
		Type:    aiModels.QuotaAlertTypeConsistency,
		Level:   aiModels.QuotaAlertLevelWarning,
		Title:   "global mismatch",
		Status:  aiModels.QuotaAlertStatusAcknowledged,
		Message: "global mismatch",
		Data: map[string]interface{}{
			"scope":     quotaConsistencyScopeGlobal,
			"timeRange": "day",
			"groupBy":   "workflow",
		},
		CreatedAt: now,
	}
	activeUserAlert := &aiModels.QuotaAlert{
		ID:      primitive.NewObjectID(),
		Type:    aiModels.QuotaAlertTypeConsistency,
		UserID:  "user-1",
		Level:   aiModels.QuotaAlertLevelWarning,
		Title:   "user mismatch",
		Status:  aiModels.QuotaAlertStatusPending,
		Message: "user mismatch",
		Data: map[string]interface{}{
			"scope":     quotaConsistencyScopeUser,
			"timeRange": "day",
		},
		CreatedAt: now,
	}
	ignoredAlert := &aiModels.QuotaAlert{
		ID:      primitive.NewObjectID(),
		Type:    aiModels.QuotaAlertTypeConsistency,
		Level:   aiModels.QuotaAlertLevelCritical,
		Title:   "ignored mismatch",
		Status:  aiModels.QuotaAlertStatusIgnored,
		Message: "ignored mismatch",
		Data: map[string]interface{}{
			"scope":     quotaConsistencyScopeWorkflow,
			"timeRange": "day",
			"groupBy":   "workflow",
			"groupKey":  "reader_chat",
		},
		CreatedAt: now,
	}
	otherTypeAlert := &aiModels.QuotaAlert{
		ID:        primitive.NewObjectID(),
		Type:      aiModels.QuotaAlertTypeThreshold,
		Level:     aiModels.QuotaAlertLevelCritical,
		Title:     "threshold",
		Status:    aiModels.QuotaAlertStatusPending,
		Message:   "threshold",
		CreatedAt: now,
	}

	repo := newQuotaAlertRepoStub(workflowAlert, globalAlert, activeUserAlert, ignoredAlert, otherTypeAlert)
	service := NewQuotaAlertService(repo)

	checkedKeys := map[string]struct{}{
		buildConsistencyAlertKey("", workflowAlert.Data):                       {},
		buildConsistencyAlertKey("", globalAlert.Data):                         {},
		buildConsistencyAlertKey(activeUserAlert.UserID, activeUserAlert.Data): {},
	}
	activeKeys := map[string]struct{}{
		buildConsistencyAlertKey(activeUserAlert.UserID, activeUserAlert.Data): {},
	}

	err := service.ResolveRecoveredConsistencyAlerts(context.Background(), checkedKeys, activeKeys, "system-consistency-check")
	assert.NoError(t, err)

	assert.Equal(t, aiModels.QuotaAlertStatusResolved, workflowAlert.Status)
	assert.Equal(t, "system-consistency-check", workflowAlert.ResolvedBy)
	assert.NotNil(t, workflowAlert.ResolvedAt)

	assert.Equal(t, aiModels.QuotaAlertStatusResolved, globalAlert.Status)
	assert.Equal(t, "system-consistency-check", globalAlert.ResolvedBy)
	assert.NotNil(t, globalAlert.ResolvedAt)

	assert.Equal(t, aiModels.QuotaAlertStatusPending, activeUserAlert.Status)
	assert.Nil(t, activeUserAlert.ResolvedAt)

	assert.Equal(t, aiModels.QuotaAlertStatusIgnored, ignoredAlert.Status)
	assert.Equal(t, aiModels.QuotaAlertStatusPending, otherTypeAlert.Status)
}
