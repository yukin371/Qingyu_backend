package dto

import (
	"time"
)

// ===========================
// 内容总结相关 DTO
// ===========================

// SummarizeRequest 总结请求
type SummarizeRequest struct {
	Content       string `json:"content" binding:"required"` // 要总结的内容
	ProjectID     string `json:"projectId,omitempty"`        // 项目ID（可选，用于上下文）
	ChapterID     string `json:"chapterId,omitempty"`        // 章节ID（可选）
	MaxLength     int    `json:"maxLength,omitempty"`        // 摘要最大长度
	SummaryType   string `json:"summaryType,omitempty"`      // 总结类型: brief, detailed, keypoints
	IncludeQuotes bool   `json:"includeQuotes,omitempty"`    // 是否包含关键引用
}

// SummarizeResponse 总结响应
type SummarizeResponse struct {
	Summary         string    `json:"summary"`             // 摘要内容
	KeyPoints       []string  `json:"keyPoints,omitempty"` // 关键点列表
	OriginalLength  int       `json:"originalLength"`      // 原文长度
	SummaryLength   int       `json:"summaryLength"`       // 摘要长度
	CompressionRate float64   `json:"compressionRate"`     // 压缩率
	TokensUsed      int       `json:"tokensUsed"`          // 使用的Token数
	Model           string    `json:"model"`               // 使用的模型
	ProcessedAt     time.Time `json:"processedAt"`         // 处理时间
}

// ChapterSummaryRequest 章节总结请求
type ChapterSummaryRequest struct {
	ChapterID    string `json:"chapterId" binding:"required"` // 章节ID
	ProjectID    string `json:"projectId" binding:"required"` // 项目ID
	OutlineLevel int    `json:"outlineLevel,omitempty"`       // 大纲级别（可选）
}

// ChapterSummaryResponse 章节总结响应
type ChapterSummaryResponse struct {
	ChapterID    string             `json:"chapterId"`             // 章节ID
	ChapterTitle string             `json:"chapterTitle"`          // 章节标题
	Summary      string             `json:"summary"`               // 章节摘要
	KeyPoints    []string           `json:"keyPoints,omitempty"`   // 关键点列表
	PlotOutline  []OutlineItem      `json:"plotOutline,omitempty"` // 情节大纲
	Characters   []CharacterMention `json:"characters,omitempty"`  // 涉及角色
	TokensUsed   int                `json:"tokensUsed"`            // 使用的Token数
	ProcessedAt  time.Time          `json:"processedAt"`           // 处理时间
}

// OutlineItem 大纲项目
type OutlineItem struct {
	Level       int           `json:"level"`                 // 层级
	Title       string        `json:"title"`                 // 标题
	Description string        `json:"description,omitempty"` // 描述
	Children    []OutlineItem `json:"children,omitempty"`    // 子项
}

// CharacterMention 角色提及
type CharacterMention struct {
	Name       string `json:"name"`                 // 角色名称
	Role       string `json:"role,omitempty"`       // 角色定位
	Appearance string `json:"appearance,omitempty"` // 出现描述
}

// ===========================
// 文本校对相关 DTO
// ===========================

// ProofreadRequest 校对请求
type ProofreadRequest struct {
	Content     string   `json:"content" binding:"required"` // 要校对的内容
	ProjectID   string   `json:"projectId,omitempty"`        // 项目ID（可选）
	ChapterID   string   `json:"chapterId,omitempty"`        // 章节ID（可选）
	CheckTypes  []string `json:"checkTypes,omitempty"`       // 检查类型: spelling, grammar, punctuation, style
	Language    string   `json:"language,omitempty"`         // 语言（自动检测）
	Suggestions bool     `json:"suggestions,omitempty"`      // 是否提供修改建议
}

// ProofreadResponse 校对响应
type ProofreadResponse struct {
	OriginalContent string         `json:"originalContent"` // 原始内容
	Issues          []Issue        `json:"issues"`          // 问题列表
	Score           float64        `json:"score"`           // 整体评分 (0-100)
	Statistics      ProofreadStats `json:"statistics"`      // 统计信息
	TokensUsed      int            `json:"tokensUsed"`      // 使用的Token数
	Model           string         `json:"model"`           // 使用的模型
	ProcessedAt     time.Time      `json:"processedAt"`     // 处理时间
}

// Issue 问题详情
type Issue struct {
	ID           string       `json:"id"`                    // 问题ID
	Type         string       `json:"type"`                  // 问题类型: spelling, grammar, punctuation, style
	Severity     string       `json:"severity"`              // 严重程度: error, warning, suggestion
	Message      string       `json:"message"`               // 问题描述
	Position     TextPosition `json:"position"`              // 位置信息
	OriginalText string       `json:"originalText"`          // 原文
	Suggestions  []string     `json:"suggestions,omitempty"` // 修改建议
	Category     string       `json:"category,omitempty"`    // 问题分类
	Rule         string       `json:"rule,omitempty"`        // 违反的规则
}

// TextPosition 文本位置
type TextPosition struct {
	Line   int `json:"line"`   // 行号（从1开始）
	Column int `json:"column"` // 列号（从1开始）
	Start  int `json:"start"`  // 字符开始位置
	End    int `json:"end"`    // 字符结束位置
	Length int `json:"length"` // 长度
}

// ProofreadStats 校对统计
type ProofreadStats struct {
	TotalIssues     int            `json:"totalIssues"`     // 总问题数
	ErrorCount      int            `json:"errorCount"`      // 错误数
	WarningCount    int            `json:"warningCount"`    // 警告数
	SuggestionCount int            `json:"suggestionCount"` // 建议数
	IssuesByType    map[string]int `json:"issuesByType"`    // 按类型统计
	WordCount       int            `json:"wordCount"`       // 词数
	CharacterCount  int            `json:"characterCount"`  // 字符数
}

// ProofreadSuggestion 校对建议详情
type ProofreadSuggestion struct {
	IssueID      string           `json:"issueId"`               // 问题ID
	Type         string           `json:"type"`                  // 问题类型
	Message      string           `json:"message"`               // 问题描述
	Position     TextPosition     `json:"position"`              // 位置信息
	OriginalText string           `json:"originalText"`          // 原文
	Suggestions  []SuggestionItem `json:"suggestions"`           // 建议列表
	Explanation  string           `json:"explanation,omitempty"` // 解释说明
	Examples     []string         `json:"examples,omitempty"`    // 示例
}

// SuggestionItem 建议项
type SuggestionItem struct {
	Text       string  `json:"text"`             // 建议文本
	Confidence float64 `json:"confidence"`       // 置信度 (0-1)
	Reason     string  `json:"reason,omitempty"` // 原因说明
}

// ===========================
// 敏感词检测相关 DTO
// ===========================

// SensitiveWordsCheckRequest 敏感词检测请求
type SensitiveWordsCheckRequest struct {
	Content     string   `json:"content" binding:"required"` // 要检测的内容
	ProjectID   string   `json:"projectId,omitempty"`        // 项目ID（可选）
	ChapterID   string   `json:"chapterId,omitempty"`        // 章节ID（可选）
	CustomWords []string `json:"customWords,omitempty"`      // 自定义敏感词
	Category    string   `json:"category,omitempty"`         // 检测分类: all, political, violence, adult, other
}

// SensitiveWordsCheckResponse 敏感词检测响应
type SensitiveWordsCheckResponse struct {
	CheckID        string               `json:"checkId"`        // 检查ID
	IsSafe         bool                 `json:"isSafe"`         // 是否安全
	TotalMatches   int                  `json:"totalMatches"`   // 总匹配数
	SensitiveWords []SensitiveWordMatch `json:"sensitiveWords"` // 敏感词列表
	Summary        CheckSummary         `json:"summary"`        // 检测摘要
	TokensUsed     int                  `json:"tokensUsed"`     // 使用的Token数
	ProcessedAt    time.Time            `json:"processedAt"`    // 处理时间
}

// SensitiveWordMatch 敏感词匹配
type SensitiveWordMatch struct {
	ID         string       `json:"id"`                   // 匹配ID
	Word       string       `json:"word"`                 // 敏感词
	Category   string       `json:"category"`             // 分类
	Level      string       `json:"level"`                // 级别: high, medium, low
	Position   TextPosition `json:"position"`             // 位置信息
	Context    string       `json:"context,omitempty"`    // 上下文（前后各50字符）
	Suggestion string       `json:"suggestion,omitempty"` // 建议修改
	Reason     string       `json:"reason,omitempty"`     // 原因说明
}

// CheckSummary 检测摘要
type CheckSummary struct {
	ByCategory      map[string]int `json:"byCategory,omitempty"` // 按分类统计
	ByLevel         map[string]int `json:"byLevel,omitempty"`    // 按级别统计
	HighRiskCount   int            `json:"highRiskCount"`        // 高风险数量
	MediumRiskCount int            `json:"mediumRiskCount"`      // 中风险数量
	LowRiskCount    int            `json:"lowRiskCount"`         // 低风险数量
}

// SensitiveWordsDetail 敏感词检测结果详情
type SensitiveWordsDetail struct {
	CheckID     string               `json:"checkId"`     // 检查ID
	Content     string               `json:"content"`     // 检测内容
	IsSafe      bool                 `json:"isSafe"`      // 是否安全
	Matches     []SensitiveWordMatch `json:"matches"`     // 匹配列表
	CustomWords []string             `json:"customWords"` // 自定义词库
	Summary     CheckSummary         `json:"summary"`     // 检测摘要
	CreatedAt   time.Time            `json:"createdAt"`   // 创建时间
	ExpiresAt   time.Time            `json:"expiresAt"`   // 过期时间（30天后）
}
