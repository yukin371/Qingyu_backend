package ai

import "Qingyu_backend/models/writer"

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

// NewCharacterInfoFromEntity 摘要
// TODO
func NewCharacterInfoFromEntity(c *writer.Character) *CharacterInfo {
	return &CharacterInfo{
		ID:      c.ID,
		Name:    c.Name,
		Summary: c.ShortDescription, // 取摘要，节省 Token
	}
}
