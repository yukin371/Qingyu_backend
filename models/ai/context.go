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

// ChapterInfo 当前章节信息
type ChapterInfo struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	Content      string   `json:"content"`
	CharacterIDs []string `json:"characterIds"`
	LocationIDs  []string `json:"locationIds"`
	TimelineIDs  []string `json:"timelineIds"`
	PlotThreads  []string `json:"plotThreads"`
	KeyPoints    []string `json:"keyPoints"`
	WritingHints string   `json:"writingHints"`
}

// CharacterInfo 角色信息（用于AI上下文）
type CharacterInfo struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Alias             []string `json:"alias,omitempty"`
	Summary           string   `json:"summary"`
	Traits            []string `json:"traits"`
	Background        string   `json:"background"`
	PersonalityPrompt string   `json:"personalityPrompt,omitempty"`
	SpeechPattern     string   `json:"speechPattern,omitempty"`
	CurrentState      string   `json:"currentState,omitempty"`
	Relationships     []string `json:"relationships,omitempty"`
}

// LocationInfo 地点信息（用于AI上下文）
type LocationInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Climate     string `json:"climate,omitempty"`
	Culture     string `json:"culture,omitempty"`
	Geography   string `json:"geography,omitempty"`
	Atmosphere  string `json:"atmosphere,omitempty"`
}

// TimelineEvent 时间线事件
type TimelineEvent struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	StoryTime    *StoryTime `json:"storyTime,omitempty"`
	Impact       string     `json:"impact,omitempty"`
	Participants []string   `json:"participants,omitempty"`
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

// ContextOptions 上下文构建选项
type ContextOptions struct {
	MaxTokens      int    `json:"maxTokens"`
	IncludeHistory bool   `json:"includeHistory"`
	IncludeOutline bool   `json:"includeOutline"`
	FocusCharacter string `json:"focusCharacter,omitempty"`
	FocusLocation  string `json:"focusLocation,omitempty"`
	HistoryDepth   int    `json:"historyDepth"` // 包含多少个历史章节
	OutlineDepth   int    `json:"outlineDepth"` // 包含多少个未来章节
}

// GenerateOptions AI生成选项
type GenerateOptions struct {
	Model       string  `json:"model,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"maxTokens,omitempty"`
	TopP        float64 `json:"topP,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
}

// GenerateResult AI生成结果
type GenerateResult struct {
	Content      string `json:"content"`
	TokensUsed   int    `json:"tokensUsed"`
	Model        string `json:"model"`
	FinishReason string `json:"finishReason,omitempty"`
}

// AnalysisResult AI分析结果
type AnalysisResult struct {
	Type       string `json:"type"`
	Analysis   string `json:"analysis"`
	TokensUsed int    `json:"tokensUsed"`
	Model      string `json:"model"`
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
