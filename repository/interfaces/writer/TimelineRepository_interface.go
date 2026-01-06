package writing

import (
	"Qingyu_backend/models/writer"
	"context"
)

// TimelineRepository 时间线Repository接口
type TimelineRepository interface {
	Create(ctx context.Context, timeline *writer.Timeline) error
	FindByID(ctx context.Context, timelineID string) (*writer.Timeline, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*writer.Timeline, error)
	Update(ctx context.Context, timeline *writer.Timeline) error
	Delete(ctx context.Context, timelineID string) error
}

// TimelineEventRepository 时间线事件Repository接口
type TimelineEventRepository interface {
	Create(ctx context.Context, event *writer.TimelineEvent) error
	FindByID(ctx context.Context, eventID string) (*writer.TimelineEvent, error)
	FindByTimelineID(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*writer.TimelineEvent, error)
	Update(ctx context.Context, event *writer.TimelineEvent) error
	Delete(ctx context.Context, eventID string) error
}
