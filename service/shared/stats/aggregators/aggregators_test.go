package aggregators

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/service/shared/stats"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAggregator_ContractAndDefaultRange(t *testing.T) {
	var port stats.UserAggregatorPort = NewUserAggregator()

	res, err := port.AggregateUserStats(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, int64(0), res.TotalUsers)

	a := NewUserAggregator()
	start, end := a.BuildRange(nil)
	assert.True(t, end.After(start))
	assert.WithinDuration(t, time.Now(), end, 2*time.Second)
}

func TestContentAggregator_Contract(t *testing.T) {
	var port stats.ContentAggregatorPort = NewContentAggregator()

	res, err := port.AggregateContentStats(context.Background(), &stats.StatsFilter{})
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, int64(0), res.TotalBooks)
	assert.NotNil(t, res.PopularCategories)
}
