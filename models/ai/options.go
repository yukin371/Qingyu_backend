package ai

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
