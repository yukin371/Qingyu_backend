package writer

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/service/base"
	eventservice "Qingyu_backend/service/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type captureEventStore struct {
	lastEvent base.Event
}

func (s *captureEventStore) Store(_ context.Context, event base.Event) error {
	s.lastEvent = event
	return nil
}

func (s *captureEventStore) StoreBatch(_ context.Context, events []base.Event) error {
	if len(events) > 0 {
		s.lastEvent = events[len(events)-1]
	}
	return nil
}

func (s *captureEventStore) GetByID(_ context.Context, _ string) (*eventservice.StoredEvent, error) {
	return nil, nil
}

func (s *captureEventStore) GetByType(_ context.Context, _ string, _, _ int64) ([]*eventservice.StoredEvent, error) {
	return nil, nil
}

func (s *captureEventStore) GetBySource(_ context.Context, _ string, _, _ int64) ([]*eventservice.StoredEvent, error) {
	return nil, nil
}

func (s *captureEventStore) GetByTimeRange(_ context.Context, _, _ time.Time, _, _ int64) ([]*eventservice.StoredEvent, error) {
	return nil, nil
}

func (s *captureEventStore) GetByTypeAndTimeRange(_ context.Context, _ string, _, _ time.Time, _, _ int64) ([]*eventservice.StoredEvent, error) {
	return nil, nil
}

func (s *captureEventStore) Replay(_ context.Context, _ base.EventHandler, _ eventservice.EventFilter) (*eventservice.ReplayResult, error) {
	return &eventservice.ReplayResult{}, nil
}

func (s *captureEventStore) Cleanup(_ context.Context, _ time.Time) (int64, error) {
	return 0, nil
}

func (s *captureEventStore) Count(_ context.Context, _ eventservice.EventFilter) (int64, error) {
	return 0, nil
}

func (s *captureEventStore) Health(_ context.Context) error {
	return nil
}

func TestPublishEventBusAdapterPreservesPublishMetadata(t *testing.T) {
	store := &captureEventStore{}
	bus := eventservice.NewPersistedEventBus(base.NewSimpleEventBus(), store, false)
	adapter := NewPublishEventBusAdapter(bus)

	err := adapter.PublishAsync(context.Background(), map[string]interface{}{
		"eventType": "project.published",
		"source":    "writer.publish_service",
		"projectId": "project-1",
	})

	require.NoError(t, err)
	require.NotNil(t, store.lastEvent)
	assert.Equal(t, "project.published", store.lastEvent.GetEventType())
	assert.Equal(t, "writer.publish_service", store.lastEvent.GetSource())

	data, ok := store.lastEvent.GetEventData().(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "project-1", data["projectId"])
}
