package events

import (
	"time"

	"Qingyu_backend/service/base"
)

// CharacterCreatedEvent 角色创建事件
type CharacterCreatedEvent struct {
	base.BaseEvent
	CharacterID string
	ProjectID   string
	UserID      string
	Name        string
}

// GetEventType 实现Event接口
func (e *CharacterCreatedEvent) GetEventType() string {
	return "character.created"
}

// GetEventData 实现Event接口
func (e *CharacterCreatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"character_id": e.CharacterID,
		"project_id":   e.ProjectID,
		"user_id":      e.UserID,
		"name":         e.Name,
	}
}

// GetTimestamp 实现Event接口
func (e *CharacterCreatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *CharacterCreatedEvent) GetSource() string {
	return e.Source
}

// CharacterUpdatedEvent 角色更新事件
type CharacterUpdatedEvent struct {
	base.BaseEvent
	CharacterID string
	ProjectID   string
}

// GetEventType 实现Event接口
func (e *CharacterUpdatedEvent) GetEventType() string {
	return "character.updated"
}

// GetEventData 实现Event接口
func (e *CharacterUpdatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"character_id": e.CharacterID,
		"project_id":   e.ProjectID,
	}
}

// GetTimestamp 实现Event接口
func (e *CharacterUpdatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *CharacterUpdatedEvent) GetSource() string {
	return e.Source
}

// CharacterDeletedEvent 角色删除事件
type CharacterDeletedEvent struct {
	base.BaseEvent
	CharacterID string
	ProjectID   string
}

// GetEventType 实现Event接口
func (e *CharacterDeletedEvent) GetEventType() string {
	return "character.deleted"
}

// GetEventData 实现Event接口
func (e *CharacterDeletedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"character_id": e.CharacterID,
		"project_id":   e.ProjectID,
	}
}

// GetTimestamp 实现Event接口
func (e *CharacterDeletedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *CharacterDeletedEvent) GetSource() string {
	return e.Source
}

// LocationCreatedEvent 地点创建事件
type LocationCreatedEvent struct {
	base.BaseEvent
	LocationID string
	ProjectID  string
	UserID     string
}

// GetEventType 实现Event接口
func (e *LocationCreatedEvent) GetEventType() string {
	return "location.created"
}

// GetEventData 实现Event接口
func (e *LocationCreatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"location_id": e.LocationID,
		"project_id":  e.ProjectID,
		"user_id":     e.UserID,
	}
}

// GetTimestamp 实现Event接口
func (e *LocationCreatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *LocationCreatedEvent) GetSource() string {
	return e.Source
}

// LocationUpdatedEvent 地点更新事件
type LocationUpdatedEvent struct {
	base.BaseEvent
	LocationID string
	ProjectID  string
}

// GetEventType 实现Event接口
func (e *LocationUpdatedEvent) GetEventType() string {
	return "location.updated"
}

// GetEventData 实现Event接口
func (e *LocationUpdatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"location_id": e.LocationID,
		"project_id":  e.ProjectID,
	}
}

// GetTimestamp 实现Event接口
func (e *LocationUpdatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *LocationUpdatedEvent) GetSource() string {
	return e.Source
}

// LocationDeletedEvent 地点删除事件
type LocationDeletedEvent struct {
	base.BaseEvent
	LocationID string
	ProjectID  string
}

// GetEventType 实现Event接口
func (e *LocationDeletedEvent) GetEventType() string {
	return "location.deleted"
}

// GetEventData 实现Event接口
func (e *LocationDeletedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"location_id": e.LocationID,
		"project_id":  e.ProjectID,
	}
}

// GetTimestamp 实现Event接口
func (e *LocationDeletedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *LocationDeletedEvent) GetSource() string {
	return e.Source
}

// TimelineEventCreatedEvent 时间线事件创建事件
type TimelineEventCreatedEvent struct {
	base.BaseEvent
	EventID    string
	TimelineID string
	ProjectID  string
}

// GetEventType 实现Event接口
func (e *TimelineEventCreatedEvent) GetEventType() string {
	return "timeline_event.created"
}

// GetEventData 实现Event接口
func (e *TimelineEventCreatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"event_id":    e.EventID,
		"timeline_id": e.TimelineID,
		"project_id":  e.ProjectID,
	}
}

// GetTimestamp 实现Event接口
func (e *TimelineEventCreatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *TimelineEventCreatedEvent) GetSource() string {
	return e.Source
}

// TimelineEventUpdatedEvent 时间线事件更新事件
type TimelineEventUpdatedEvent struct {
	base.BaseEvent
	EventID   string
	ProjectID string
}

// GetEventType 实现Event接口
func (e *TimelineEventUpdatedEvent) GetEventType() string {
	return "timeline_event.updated"
}

// GetEventData 实现Event接口
func (e *TimelineEventUpdatedEvent) GetEventData() interface{} {
	return map[string]interface{}{
		"event_id":   e.EventID,
		"project_id": e.ProjectID,
	}
}

// GetTimestamp 实现Event接口
func (e *TimelineEventUpdatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetSource 实现Event接口
func (e *TimelineEventUpdatedEvent) GetSource() string {
	return e.Source
}
