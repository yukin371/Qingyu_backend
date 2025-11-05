package writer

import (
	"fmt"
	"time"
)

// Timeline 时间线
type Timeline struct {
	ID          string     `bson:"_id,omitempty" json:"id"`
	ProjectID   string     `bson:"project_id" json:"projectId"`
	Name        string     `bson:"name" json:"name"`
	Description string     `bson:"description,omitempty" json:"description,omitempty"`
	StartTime   *StoryTime `bson:"start_time,omitempty" json:"startTime,omitempty"`
	EndTime     *StoryTime `bson:"end_time,omitempty" json:"endTime,omitempty"`
	CreatedAt   time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updatedAt"`
}

// TimelineEvent 时间线事件
type TimelineEvent struct {
	ID           string     `bson:"_id,omitempty" json:"id"`
	ProjectID    string     `bson:"project_id" json:"projectId"`
	TimelineID   string     `bson:"timeline_id" json:"timelineId"`
	Title        string     `bson:"title" json:"title"`
	Description  string     `bson:"description,omitempty" json:"description,omitempty"`
	StoryTime    *StoryTime `bson:"story_time,omitempty" json:"storyTime,omitempty"`
	Duration     string     `bson:"duration,omitempty" json:"duration,omitempty"`         // 事件持续时间描述
	Impact       string     `bson:"impact,omitempty" json:"impact,omitempty"`             // 事件影响
	Participants []string   `bson:"participants,omitempty" json:"participants,omitempty"` // 参与角色ID
	LocationIDs  []string   `bson:"location_ids,omitempty" json:"locationIds,omitempty"`  // 相关地点ID
	ChapterIDs   []string   `bson:"chapter_ids,omitempty" json:"chapterIds,omitempty"`    // 相关章节ID
	EventType    EventType  `bson:"event_type" json:"eventType"`
	Importance   int        `bson:"importance" json:"importance"` // 重要性等级 1-10
	CreatedAt    time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `bson:"updated_at" json:"updatedAt"`
}

// StoryTime 故事时间
type StoryTime struct {
	Year        int    `bson:"year,omitempty" json:"year,omitempty"`
	Month       int    `bson:"month,omitempty" json:"month,omitempty"`
	Day         int    `bson:"day,omitempty" json:"day,omitempty"`
	Hour        int    `bson:"hour,omitempty" json:"hour,omitempty"`
	Minute      int    `bson:"minute,omitempty" json:"minute,omitempty"`
	Era         string `bson:"era,omitempty" json:"era,omitempty"`
	Season      string `bson:"season,omitempty" json:"season,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"` // 时间描述
}

// EventType 事件类型
type EventType string

const (
	EventTypePlot       EventType = "plot"       // 情节事件
	EventTypeCharacter  EventType = "character"  // 角色事件
	EventTypeWorld      EventType = "world"      // 世界事件
	EventTypeBackground EventType = "background" // 背景事件
	EventTypeMilestone  EventType = "milestone"  // 里程碑事件
)

// IsValidEventType 验证事件类型是否合法
func IsValidEventType(t string) bool {
	switch EventType(t) {
	case EventTypePlot, EventTypeCharacter, EventTypeWorld, EventTypeBackground, EventTypeMilestone:
		return true
	default:
		return false
	}
}

// GetTimeString 获取时间字符串表示
func (st *StoryTime) GetTimeString() string {
	if st == nil {
		return ""
	}

	if st.Description != "" {
		return st.Description
	}

	timeStr := ""
	if st.Era != "" {
		timeStr += st.Era + " "
	}
	if st.Year > 0 {
		timeStr += fmt.Sprintf("%d年", st.Year)
	}
	if st.Month > 0 {
		timeStr += fmt.Sprintf("%d月", st.Month)
	}
	if st.Day > 0 {
		timeStr += fmt.Sprintf("%d日", st.Day)
	}
	if st.Season != "" {
		timeStr += " " + st.Season
	}

	return timeStr
}
