package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AIContext AI上下文数据结构
type AIContext struct {
	ProjectID        string            `json:"projectId"`
	CurrentChapter   *ChapterInfo      `json:"currentChapter"`
	ActiveCharacters []*CharacterInfo  `json:"activeCharacters"`
	CurrentLocations []*LocationInfo   `json:"currentLocations"`
	RelevantEvents   []*TimelineEvent  `json:"relevantEvents"`
	PreviousChapters []*ChapterSummary `json:"previousChapters"`
	NextChapters     []*ChapterOutline `json:"nextChapters"`
	WorldSettings    *WorldSettings    `json:"worldSettings"`
	PlotThreads      []*PlotThread     `json:"plotThreads"`
	TokenCount       int               `json:"tokenCount"`
}

// StoryTime 故事时间
type StoryTime struct {
	Year   int    `json:"year,omitempty"`
	Month  int    `json:"month,omitempty"`
	Day    int    `json:"day,omitempty"`
	Hour   int    `json:"hour,omitempty"`
	Minute int    `json:"minute,omitempty"`
	Era    string `json:"era,omitempty"`
}

// ChapterSummary 章节摘要
type ChapterSummary struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Order     int      `json:"order"`
	KeyEvents []string `json:"keyEvents,omitempty"`
}

// ChapterOutline 章节大纲
type ChapterOutline struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Outline string   `json:"outline"`
	Order   int      `json:"order"`
	Goals   []string `json:"goals,omitempty"`
}

// WorldSettings 世界观设定
type WorldSettings struct {
	ID          string                 `json:"id"`
	ProjectID   string                 `json:"projectId"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Rules       []string               `json:"rules,omitempty"`
	History     string                 `json:"history,omitempty"`
	Geography   string                 `json:"geography,omitempty"`
	Culture     string                 `json:"culture,omitempty"`
	Magic       string                 `json:"magic,omitempty"`
	Technology  string                 `json:"technology,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// PlotThread 情节线索
type PlotThread struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"projectId"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`   // active, resolved, pending, suspended
	Priority     int       `json:"priority"` // 1-10, 10为最高优先级
	ChapterIDs   []string  `json:"chapterIds"`
	Characters   []string  `json:"characters,omitempty"`
	StartChapter string    `json:"startChapter,omitempty"`
	EndChapter   string    `json:"endChapter,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// BeforeCreate 在创建前设置时间戳
func (w *WorldSettings) BeforeCreate() {
	now := time.Now()
	w.CreatedAt = now
	w.UpdatedAt = now
	if w.ID == "" {
		w.ID = primitive.NewObjectID().Hex()
	}
}

// BeforeUpdate 在更新前刷新更新时间戳
func (w *WorldSettings) BeforeUpdate() {
	w.UpdatedAt = time.Now()
}

// BeforeCreate 在创建前设置时间戳
func (p *PlotThread) BeforeCreate() {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.ID == "" {
		p.ID = primitive.NewObjectID().Hex()
	}
}

// BeforeUpdate 在更新前刷新更新时间戳
func (p *PlotThread) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}
