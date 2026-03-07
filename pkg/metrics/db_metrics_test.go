package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
)

func TestDbMetricsCollectorMonitorSuccess(t *testing.T) {
	collector := NewDbMetricsCollector("testdb", 100*time.Millisecond, 1)
	monitor := collector.GetMonitorCommand()

	before := testutil.ToFloat64(collector.metrics.dbQueryTotal.WithLabelValues("testdb", "find", "success"))

	monitor.Started(context.Background(), &event.CommandStartedEvent{
		RequestID:    101,
		CommandName:  "find",
		ConnectionID: "localhost:27017[-1]",
		Command: mustMarshalBSON(t, bson.D{
			{Key: "find", Value: "users"},
			{Key: "filter", Value: bson.D{{Key: "username", Value: "author_new"}}},
		}),
	})
	monitor.Succeeded(context.Background(), &event.CommandSucceededEvent{
		CommandFinishedEvent: event.CommandFinishedEvent{
			RequestID:     101,
			CommandName:   "find",
			DurationNanos: int64(120 * time.Millisecond),
		},
	})

	after := testutil.ToFloat64(collector.metrics.dbQueryTotal.WithLabelValues("testdb", "find", "success"))
	assert.Equal(t, before+1, after)

	_, exists := collector.inFlightQueries.Load(int64(101))
	assert.False(t, exists)
}

func TestDbMetricsCollectorMonitorFailure(t *testing.T) {
	collector := NewDbMetricsCollector("testdb", 100*time.Millisecond, 1)
	monitor := collector.GetMonitorCommand()

	beforeTotal := testutil.ToFloat64(collector.metrics.dbQueryTotal.WithLabelValues("testdb", "aggregate", "error"))
	beforeErrors := testutil.ToFloat64(collector.metrics.dbQueryErrors.WithLabelValues("testdb", "aggregate"))

	monitor.Started(context.Background(), &event.CommandStartedEvent{
		RequestID:   202,
		CommandName: "aggregate",
		Command: mustMarshalBSON(t, bson.D{
			{Key: "aggregate", Value: "comments"},
		}),
	})
	monitor.Failed(context.Background(), &event.CommandFailedEvent{
		CommandFinishedEvent: event.CommandFinishedEvent{
			RequestID:     202,
			CommandName:   "aggregate",
			DurationNanos: int64(50 * time.Millisecond),
		},
		Failure: "mock query failure",
	})

	afterTotal := testutil.ToFloat64(collector.metrics.dbQueryTotal.WithLabelValues("testdb", "aggregate", "error"))
	afterErrors := testutil.ToFloat64(collector.metrics.dbQueryErrors.WithLabelValues("testdb", "aggregate"))
	assert.Equal(t, beforeTotal+1, afterTotal)
	assert.Equal(t, beforeErrors+1, afterErrors)

	_, exists := collector.inFlightQueries.Load(int64(202))
	assert.False(t, exists)
}

func TestCollectionNameFromCommand(t *testing.T) {
	command, err := bson.Marshal(bson.D{{Key: "find", Value: "notifications"}})
	assert.NoError(t, err)

	assert.Equal(t, "notifications", collectionNameFromCommand(command, "find"))
	assert.Equal(t, "", collectionNameFromCommand(command, "aggregate"))
}

func TestDbMetricsCollectorShouldLogQuery(t *testing.T) {
	t.Run("profiling off", func(t *testing.T) {
		collector := NewDbMetricsCollector("testdb", 100*time.Millisecond, 0)
		assert.False(t, collector.shouldLogQuery(200*time.Millisecond))
	})

	t.Run("slow only", func(t *testing.T) {
		collector := NewDbMetricsCollector("testdb", 100*time.Millisecond, 1)
		assert.False(t, collector.shouldLogQuery(90*time.Millisecond))
		assert.True(t, collector.shouldLogQuery(100*time.Millisecond))
	})

	t.Run("all queries", func(t *testing.T) {
		collector := NewDbMetricsCollector("testdb", 100*time.Millisecond, 2)
		assert.True(t, collector.shouldLogQuery(10*time.Millisecond))
	})
}

func mustMarshalBSON(t *testing.T, document bson.D) bson.Raw {
	t.Helper()
	raw, err := bson.Marshal(document)
	if err != nil {
		t.Fatalf("marshal bson failed: %v", err)
	}
	return raw
}
