package writer

import (
	"context"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/service/base"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// TimelineService 时间线服务实现
type TimelineService struct {
	timelineRepo      writerRepo.TimelineRepository
	timelineEventRepo writerRepo.TimelineEventRepository
	eventBus          base.EventBus
}

// NewTimelineService 创建TimelineService实例
func NewTimelineService(
	timelineRepo writerRepo.TimelineRepository,
	timelineEventRepo writerRepo.TimelineEventRepository,
	eventBus base.EventBus,
) serviceInterfaces.TimelineService {
	return &TimelineService{
		timelineRepo:      timelineRepo,
		timelineEventRepo: timelineEventRepo,
		eventBus:          eventBus,
	}
}

// CreateTimeline 创建时间线
func (s *TimelineService) CreateTimeline(
	ctx context.Context,
	projectID string,
	req *serviceInterfaces.CreateTimelineRequest,
) (*writer.Timeline, error) {
	timeline := &writer.Timeline{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	if err := s.timelineRepo.Create(ctx, timeline); err != nil {
		return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorInternal, "create timeline failed", "", err)
	}

	return timeline, nil
}

// GetTimeline 根据ID获取时间线
func (s *TimelineService) GetTimeline(
	ctx context.Context,
	timelineID, projectID string,
) (*writer.Timeline, error) {
	timeline, err := s.timelineRepo.FindByID(ctx, timelineID)
	if err != nil {
		return nil, err
	}

	if timeline.ProjectID != projectID {
		return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorForbidden, "no permission to access this timeline", "", nil)
	}

	return timeline, nil
}

// ListTimelines 获取项目下的所有时间线
func (s *TimelineService) ListTimelines(
	ctx context.Context,
	projectID string,
) ([]*writer.Timeline, error) {
	return s.timelineRepo.FindByProjectID(ctx, projectID)
}

// DeleteTimeline 删除时间线
func (s *TimelineService) DeleteTimeline(
	ctx context.Context,
	timelineID, projectID string,
) error {
	if _, err := s.GetTimeline(ctx, timelineID, projectID); err != nil {
		return err
	}

	return s.timelineRepo.Delete(ctx, timelineID)
}

// CreateEvent 创建时间线事件
func (s *TimelineService) CreateEvent(
	ctx context.Context,
	projectID string,
	req *serviceInterfaces.CreateTimelineEventRequest,
) (*writer.TimelineEvent, error) {
	// 验证时间线存在
	timeline, err := s.timelineRepo.FindByID(ctx, req.TimelineID)
	if err != nil {
		return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorNotFound, "timeline not found", "", err)
	}

	if timeline.ProjectID != projectID {
		return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorForbidden, "no permission", "", nil)
	}

	// 验证事件类型
	if req.EventType != "" && !writer.IsValidEventType(req.EventType) {
		return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorValidation, "invalid event type", "", nil)
	}

	event := &writer.TimelineEvent{
		ProjectID:    projectID,
		TimelineID:   req.TimelineID,
		Title:        req.Title,
		Description:  req.Description,
		EventType:    writer.EventType(req.EventType),
		Importance:   req.Importance,
		Participants: req.Participants,
		LocationIDs:  req.LocationIDs,
		ChapterIDs:   req.ChapterIDs,
		StoryTime:    req.StoryTime,
		Duration:     req.Duration,
		Impact:       req.Impact,
	}

	if err := s.timelineEventRepo.Create(ctx, event); err != nil {
		return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorInternal, "create timeline event failed", "", err)
	}

	// 发布事件
	if s.eventBus != nil {
		busEvent := &base.BaseEvent{
			EventType: "timeline_event.created",
			EventData: map[string]interface{}{
				"event_id":    event.ID,
				"timeline_id": req.TimelineID,
				"project_id":  projectID,
			},
			Timestamp: time.Now(),
			Source:    "TimelineService",
		}
		s.eventBus.PublishAsync(ctx, busEvent)
	}

	return event, nil
}

// GetEvent 根据ID获取事件
func (s *TimelineService) GetEvent(
	ctx context.Context,
	eventID, projectID string,
) (*writer.TimelineEvent, error) {
	event, err := s.timelineEventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.ProjectID != projectID {
		return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorForbidden, "no permission to access this event", "", nil)
	}

	return event, nil
}

// ListEvents 获取时间线下的所有事件
func (s *TimelineService) ListEvents(
	ctx context.Context,
	timelineID string,
) ([]*writer.TimelineEvent, error) {
	return s.timelineEventRepo.FindByTimelineID(ctx, timelineID)
}

// UpdateEvent 更新时间线事件
func (s *TimelineService) UpdateEvent(
	ctx context.Context,
	eventID, projectID string,
	req *serviceInterfaces.UpdateTimelineEventRequest,
) (*writer.TimelineEvent, error) {
	event, err := s.GetEvent(ctx, eventID, projectID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.EventType != nil {
		if !writer.IsValidEventType(*req.EventType) {
			return nil, errors.NewServiceError("TimelineService", errors.ServiceErrorValidation, "invalid event type", "", nil)
		}
		event.EventType = writer.EventType(*req.EventType)
	}
	if req.Importance != nil {
		event.Importance = *req.Importance
	}
	if req.Participants != nil {
		event.Participants = *req.Participants
	}
	if req.LocationIDs != nil {
		event.LocationIDs = *req.LocationIDs
	}
	if req.ChapterIDs != nil {
		event.ChapterIDs = *req.ChapterIDs
	}
	if req.StoryTime != nil {
		event.StoryTime = req.StoryTime
	}
	if req.Duration != nil {
		event.Duration = *req.Duration
	}
	if req.Impact != nil {
		event.Impact = *req.Impact
	}

	if err := s.timelineEventRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	// 发布事件
	if s.eventBus != nil {
		busEvent := &base.BaseEvent{
			EventType: "timeline_event.updated",
			EventData: map[string]interface{}{
				"event_id":   eventID,
				"project_id": projectID,
			},
			Timestamp: time.Now(),
			Source:    "TimelineService",
		}
		s.eventBus.PublishAsync(ctx, busEvent)
	}

	return event, nil
}

// DeleteEvent 删除时间线事件
func (s *TimelineService) DeleteEvent(
	ctx context.Context,
	eventID, projectID string,
) error {
	if _, err := s.GetEvent(ctx, eventID, projectID); err != nil {
		return err
	}

	return s.timelineEventRepo.Delete(ctx, eventID)
}

// GetTimelineVisualization 获取时间线可视化数据
func (s *TimelineService) GetTimelineVisualization(
	ctx context.Context,
	timelineID string,
) (*serviceInterfaces.TimelineVisualization, error) {
	events, err := s.timelineEventRepo.FindByTimelineID(ctx, timelineID)
	if err != nil {
		return nil, err
	}

	var nodes []*serviceInterfaces.TimelineEventNode
	for _, event := range events {
		storyTimeStr := ""
		if event.StoryTime != nil {
			storyTimeStr = event.StoryTime.GetTimeString()
		}

		nodes = append(nodes, &serviceInterfaces.TimelineEventNode{
			ID:         event.ID,
			Title:      event.Title,
			StoryTime:  storyTimeStr,
			EventType:  string(event.EventType),
			Importance: event.Importance,
			Characters: event.Participants,
		})
	}

	// 分析事件连接（简化版，基于时间顺序）
	connections := analyzeEventConnections(events)

	return &serviceInterfaces.TimelineVisualization{
		Events:      nodes,
		Connections: connections,
	}, nil
}

// analyzeEventConnections 分析事件连接关系
func analyzeEventConnections(events []*writer.TimelineEvent) []*serviceInterfaces.EventConnection {
	var connections []*serviceInterfaces.EventConnection

	// 简单实现：相邻事件建立sequential连接
	for i := 0; i < len(events)-1; i++ {
		connections = append(connections, &serviceInterfaces.EventConnection{
			FromEventID: events[i].ID,
			ToEventID:   events[i+1].ID,
			Type:        "sequential",
		})
	}

	return connections
}
