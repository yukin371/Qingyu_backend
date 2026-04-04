package ai

// StoryContext 故事上下文（三层漏斗压缩结果）
type StoryContext struct {
	// Layer 1: 结构化舞台
	Stage *StageContext `json:"stage"`
	// Layer 2: 大纲即摘要
	OutlineContext string `json:"outlineContext"`
	// Layer 3: RAG 检索结果
	RAGExcerpts string `json:"ragExcerpts"`
	// 最近已写文本
	RecentText string `json:"recentText"`
	// 写作指令
	Instruction string `json:"instruction"`
	// 模式: continue | rewrite | suggest
	Mode string `json:"mode"`
	// 选中文本（改写模式用）
	SelectedText string `json:"selectedText,omitempty"`
	// Token 统计
	StageTokens   int `json:"stageTokens"`
	OutlineTokens int `json:"outlineTokens"`
	RAGTokens     int `json:"ragTokens"`
	TotalTokens   int `json:"totalTokens"`
}

// StageContext 结构化舞台（Layer 1）
type StageContext struct {
	SceneGoal      string             `json:"sceneGoal,omitempty"`
	ActiveConflict string             `json:"activeConflict,omitempty"`
	Location       string             `json:"location,omitempty"`
	PlotThreads    []string           `json:"plotThreads,omitempty"`
	Characters     []*StageCharacter  `json:"characters,omitempty"`
	Relations      []*StageRelation   `json:"relations,omitempty"`
}

// StageCharacter 场景中的角色
type StageCharacter struct {
	Name          string `json:"name"`
	ShortDesc     string `json:"shortDesc,omitempty"`
	CurrentState  string `json:"currentState,omitempty"`
	Personality   string `json:"personality,omitempty"`
	SpeechPattern string `json:"speechPattern,omitempty"`
}

// StageRelation 场景中的角色关系
type StageRelation struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Type    string `json:"type"`
	Tension string `json:"tension,omitempty"` // high | medium | low
}

// StoryGenerateRequest 故事生成请求
type StoryGenerateRequest struct {
	ProjectID    string `json:"projectId" binding:"required"`
	DocumentID   string `json:"documentId" binding:"required"`
	Mode         string `json:"mode" binding:"required,oneof=continue rewrite suggest"`
	Instruction  string `json:"instruction,omitempty"`
	SelectedText string `json:"selectedText,omitempty"`
}

// StoryGenerateResponse 故事生成响应
type StoryGenerateResponse struct {
	Content      string        `json:"content"`
	ContextStats *ContextStats `json:"contextStats,omitempty"`
	TokensUsed   int           `json:"tokensUsed"`
	Model        string        `json:"model"`
}

// ContextStats 上下文统计
type ContextStats struct {
	StageTokens   int `json:"stageTokens"`
	OutlineTokens int `json:"outlineTokens"`
	RAGTokens     int `json:"ragTokens"`
	TotalTokens   int `json:"totalTokens"`
}
