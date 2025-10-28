package interfaces

import (
	"Qingyu_backend/models/writer"
	"context"
)

// TimelineService 时间线服务接口
type TimelineService interface {
	// Timeline CRUD
	CreateTimeline(ctx context.Context, projectID string, req *CreateTimelineRequest) (*writer.Timeline, error)
	GetTimeline(ctx context.Context, timelineID, projectID string) (*writer.Timeline, error)
	ListTimelines(ctx context.Context, projectID string) ([]*writer.Timeline, error)
	DeleteTimeline(ctx context.Context, timelineID, projectID string) error

	// Timeline Event CRUD
	CreateEvent(ctx context.Context, projectID string, req *CreateTimelineEventRequest) (*writer.TimelineEvent, error)
	GetEvent(ctx context.Context, eventID, projectID string) (*writer.TimelineEvent, error)
	ListEvents(ctx context.Context, timelineID string) ([]*writer.TimelineEvent, error)
	UpdateEvent(ctx context.Context, eventID, projectID string, req *UpdateTimelineEventRequest) (*writer.TimelineEvent, error)
	DeleteEvent(ctx context.Context, eventID, projectID string) error

	// Visualization
	GetTimelineVisualization(ctx context.Context, timelineID string) (*TimelineVisualization, error)
}

// CreateTimelineRequest 创建时间线请求
type CreateTimelineRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	StartTime   *writer.StoryTime `json:"startTime"`
	EndTime     *writer.StoryTime `json:"endTime"`
}

// CreateTimelineEventRequest 创建时间线事件请求
type CreateTimelineEventRequest struct {
	TimelineID   string            `json:"timelineId" validate:"required"`
	Title        string            `json:"title" validate:"required"`
	Description  string            `json:"description"`
	EventType    string            `json:"eventType"`
	Importance   int               `json:"importance"`
	Participants []string          `json:"participants"`
	LocationIDs  []string          `json:"locationIds"`
	ChapterIDs   []string          `json:"chapterIds"`
	StoryTime    *writer.StoryTime `json:"storyTime"`
	Duration     string            `json:"duration"`
	Impact       string            `json:"impact"`
}

// UpdateTimelineEventRequest 更新时间线事件请求
type UpdateTimelineEventRequest struct {
	Title        *string           `json:"title"`
	Description  *string           `json:"description"`
	EventType    *string           `json:"eventType"`
	Importance   *int              `json:"importance"`
	Participants *[]string         `json:"participants"`
	LocationIDs  *[]string         `json:"locationIds"`
	ChapterIDs   *[]string         `json:"chapterIds"`
	StoryTime    *writer.StoryTime `json:"storyTime"`
	Duration     *string           `json:"duration"`
	Impact       *string           `json:"impact"`
}

// TimelineVisualization 时间线可视化数据
type TimelineVisualization struct {
	Events      []*TimelineEventNode `json:"events"`
	Connections []*EventConnection   `json:"connections"`
}

// TimelineEventNode 时间线事件节点
type TimelineEventNode struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	StoryTime  string   `json:"storyTime"`
	EventType  string   `json:"eventType"`
	Importance int      `json:"importance"`
	Characters []string `json:"characters"`
}

// EventConnection 事件连接
type EventConnection struct {
	FromEventID string `json:"fromEventId"`
	ToEventID   string `json:"toEventId"`
	Type        string `json:"type"` // sequential, causal, parallel
}
