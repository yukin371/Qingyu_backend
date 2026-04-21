package ai

import (
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"

	"github.com/stretchr/testify/assert"
)

func TestSelectLatestQuota_PrefersNewerUpdatedAt(t *testing.T) {
	older := &aiModels.UserQuota{
		UserID:    "user-1",
		QuotaType: aiModels.QuotaTypeDaily,
		UpdatedAt: time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC),
	}
	newer := &aiModels.UserQuota{
		UserID:    "user-1",
		QuotaType: aiModels.QuotaTypeDaily,
		UpdatedAt: time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC),
	}

	assert.Same(t, newer, selectLatestQuota(older, newer))
	assert.Same(t, older, selectLatestQuota(older, nil))
	assert.Same(t, newer, selectLatestQuota(nil, newer))
}

func TestMergeLatestQuotasByKey_DeduplicatesAndKeepsLatest(t *testing.T) {
	dailyOlder := &aiModels.UserQuota{
		UserID:    "user-1",
		QuotaType: aiModels.QuotaTypeDaily,
		UpdatedAt: time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC),
	}
	dailyNewer := &aiModels.UserQuota{
		UserID:    "user-1",
		QuotaType: aiModels.QuotaTypeDaily,
		UpdatedAt: time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC),
	}
	monthly := &aiModels.UserQuota{
		UserID:    "user-1",
		QuotaType: aiModels.QuotaTypeMonthly,
		UpdatedAt: time.Date(2026, 4, 19, 10, 0, 0, 0, time.UTC),
	}

	merged := map[string]*aiModels.UserQuota{}
	mergeLatestQuotasByKey(merged, []*aiModels.UserQuota{dailyOlder, nil, monthly})
	mergeLatestQuotasByKey(merged, []*aiModels.UserQuota{dailyNewer})

	assert.Len(t, merged, 2)
	assert.Same(t, dailyNewer, merged[quotaDedupKey("user-1", aiModels.QuotaTypeDaily)])
	assert.Same(t, monthly, merged[quotaDedupKey("user-1", aiModels.QuotaTypeMonthly)])
}
