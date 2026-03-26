package writer

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Timeline 时间线
type Timeline struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID   string             `bson:"project_id" json:"projectId"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	StartTime   *StoryTime         `bson:"start_time,omitempty" json:"startTime,omitempty"`
	EndTime     *StoryTime         `bson:"end_time,omitempty" json:"endTime,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// TimelineEvent 时间线事件
type TimelineEvent struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID    string             `bson:"project_id" json:"projectId"`
	TimelineID   string             `bson:"timeline_id" json:"timelineId"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	StoryTime    *StoryTime         `bson:"story_time,omitempty" json:"storyTime,omitempty"`
	Duration     string             `bson:"duration,omitempty" json:"duration,omitempty"`
	Impact       string             `bson:"impact,omitempty" json:"impact,omitempty"`
	Participants []string           `bson:"participants,omitempty" json:"participants,omitempty"`
	LocationIDs  []string           `bson:"location_ids,omitempty" json:"locationIds,omitempty"`
	ChapterIDs   []string           `bson:"chapter_ids,omitempty" json:"chapterIds,omitempty"`
	EventType    EventType          `bson:"event_type" json:"eventType"`
	Importance   int                `bson:"importance" json:"importance"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

type StoryTime struct {
	Year        int    `bson:"year,omitempty" json:"year,omitempty"`
	Month       int    `bson:"month,omitempty" json:"month,omitempty"`
	Day         int    `bson:"day,omitempty" json:"day,omitempty"`
	Hour        int    `bson:"hour,omitempty" json:"hour,omitempty"`
	Minute      int    `bson:"minute,omitempty" json:"minute,omitempty"`
	Era         string `bson:"era,omitempty" json:"era,omitempty"`
	Season      string `bson:"season,omitempty" json:"season,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
}

type EventType string

const (
	EventTypePlot       EventType = "plot"
	EventTypeCharacter  EventType = "character"
	EventTypeWorld      EventType = "world"
	EventTypeBackground EventType = "background"
	EventTypeMilestone  EventType = "milestone"
)

func IsValidEventType(t string) bool {
	switch EventType(t) {
	case EventTypePlot, EventTypeCharacter, EventTypeWorld, EventTypeBackground, EventTypeMilestone:
		return true
	default:
		return false
	}
}

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

// RelationTimelineEvent 角色关系时序变化事件
type RelationTimelineEvent struct {
	ChapterID     string       `bson:"chapter_id" json:"chapterId"`             // 变化发生的章节ID
	ChapterTitle  string       `bson:"chapter_title" json:"chapterTitle"`       // 章节标题
	OldType       RelationType `bson:"old_type,omitempty" json:"oldType,omitempty"` // 变化前的关系类型
	NewType       RelationType `bson:"new_type" json:"newType"`                 // 变化后的关系类型
	Strength      int          `bson:"strength" json:"strength"`               // 变化后的强度
	Notes         string       `bson:"notes" json:"notes"`                     // 变化原因/剧情描述
	Timestamp     time.Time    `bson:"timestamp" json:"timestamp"`             // 记录时间
}
