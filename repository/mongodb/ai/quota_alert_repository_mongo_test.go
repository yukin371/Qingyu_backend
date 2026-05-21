package ai

import (
	"testing"

	aiModels "Qingyu_backend/models/ai"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildQuotaAlertListFilterSupportsOpenAlias(t *testing.T) {
	filter := buildQuotaAlertListFilter(
		string(aiModels.QuotaAlertTypeConsistency),
		string(aiModels.QuotaAlertLevelCritical),
		quotaAlertStatusFilterOpen,
	)

	assert.Equal(t, string(aiModels.QuotaAlertTypeConsistency), filter["type"])
	assert.Equal(t, string(aiModels.QuotaAlertLevelCritical), filter["level"])
	assert.Equal(t, bson.M{
		"$in": []string{
			string(aiModels.QuotaAlertStatusPending),
			string(aiModels.QuotaAlertStatusAcknowledged),
		},
	}, filter["status"])
}

func TestBuildQuotaAlertListFilterKeepsExactStatus(t *testing.T) {
	filter := buildQuotaAlertListFilter("", "", string(aiModels.QuotaAlertStatusResolved))

	assert.Equal(t, string(aiModels.QuotaAlertStatusResolved), filter["status"])
}

func TestBuildQuotaAlertListFilterSupportsAllAlias(t *testing.T) {
	filter := buildQuotaAlertListFilter(
		string(aiModels.QuotaAlertTypeConsistency),
		"",
		quotaAlertStatusFilterAll,
	)

	assert.Equal(t, string(aiModels.QuotaAlertTypeConsistency), filter["type"])
	_, hasStatus := filter["status"]
	assert.False(t, hasStatus)
}
