package ai

// Phase3创作相关的请求和响应模型

// ============ 大纲生成 ============

// GenerateOutlineRequest 生成大纲请求
type GenerateOutlineRequest struct {
	Task             string            `json:"task" binding:"required"`     // 创作任务描述
	ProjectID        string            `json:"project_id"`                  // 项目ID（可选）
	WorkspaceContext map[string]string `json:"workspace_context,omitempty"` // 工作区上下文
	CorrectionPrompt string            `json:"correction_prompt,omitempty"` // 修正提示
}

// GenerateOutlineResponse 生成大纲响应
type GenerateOutlineResponse struct {
	Outline       *OutlineData `json:"outline"`        // 大纲数据
	ExecutionTime float32      `json:"execution_time"` // 执行时间（秒）
}

// OutlineData 大纲数据
type OutlineData struct {
	Title               string        `json:"title"`                 // 标题
	Genre               string        `json:"genre"`                 // 类型
	CoreTheme           string        `json:"core_theme"`            // 核心主题
	TargetAudience      string        `json:"target_audience"`       // 目标受众
	EstimatedTotalWords int32         `json:"estimated_total_words"` // 预估总字数
	Chapters            []ChapterData `json:"chapters"`              // 章节列表
	StoryArc            *StoryArc     `json:"story_arc"`             // 故事结构
}

// ChapterData 章节数据
type ChapterData struct {
	ChapterID          int32    `json:"chapter_id"`           // 章节ID
	Title              string   `json:"title"`                // 标题
	Summary            string   `json:"summary"`              // 概要
	KeyEvents          []string `json:"key_events"`           // 关键事件
	CharactersInvolved []string `json:"characters_involved"`  // 参与角色
	ConflictType       string   `json:"conflict_type"`        // 冲突类型
	EmotionalTone      string   `json:"emotional_tone"`       // 情感基调
	EstimatedWordCount int32    `json:"estimated_word_count"` // 预估字数
	ChapterGoal        string   `json:"chapter_goal"`         // 章节目标
	Cliffhanger        string   `json:"cliffhanger"`          // 悬念
}

// StoryArc 故事结构
type StoryArc struct {
	Setup         []int32 `json:"setup"`          // 起（章节ID列表）
	RisingAction  []int32 `json:"rising_action"`  // 承
	Climax        []int32 `json:"climax"`         // 转
	FallingAction []int32 `json:"falling_action"` // 合
	Resolution    []int32 `json:"resolution"`     // 结局
}

// ============ 角色生成 ============

// GenerateCharactersRequest 生成角色请求
type GenerateCharactersRequest struct {
	Task             string            `json:"task" binding:"required"` // 创作任务描述
	ProjectID        string            `json:"project_id"`              // 项目ID
	Outline          *OutlineData      `json:"outline"`                 // 大纲数据（可选，建议提供）
	WorkspaceContext map[string]string `json:"workspace_context,omitempty"`
	CorrectionPrompt string            `json:"correction_prompt,omitempty"`
}

// GenerateCharactersResponse 生成角色响应
type GenerateCharactersResponse struct {
	Characters    *CharactersData `json:"characters"`     // 角色数据
	ExecutionTime float32         `json:"execution_time"` // 执行时间
}

// CharactersData 角色数据
type CharactersData struct {
	Characters          []CharacterData      `json:"characters"`           // 角色列表
	RelationshipNetwork *RelationshipNetwork `json:"relationship_network"` // 关系网络
}

// CharacterData 角色数据
type CharacterData struct {
	CharacterID      string             `json:"character_id"`      // 角色ID
	Name             string             `json:"name"`              // 姓名
	RoleType         string             `json:"role_type"`         // 角色类型
	Importance       string             `json:"importance"`        // 重要性
	Age              int32              `json:"age"`               // 年龄
	Gender           string             `json:"gender"`            // 性别
	Appearance       string             `json:"appearance"`        // 外貌
	Personality      *PersonalityData   `json:"personality"`       // 性格
	Background       *BackgroundData    `json:"background"`        // 背景
	Motivation       string             `json:"motivation"`        // 动机
	Relationships    []RelationshipData `json:"relationships"`     // 关系
	DevelopmentArc   *DevelopmentArc    `json:"development_arc"`   // 发展弧线
	RoleInStory      string             `json:"role_in_story"`     // 在故事中的作用
	FirstAppearance  int32              `json:"first_appearance"`  // 首次出场章节
	ChaptersInvolved []int32            `json:"chapters_involved"` // 参与章节
}

// PersonalityData 性格数据
type PersonalityData struct {
	Traits     []string `json:"traits"`      // 性格特质
	Strengths  []string `json:"strengths"`   // 优点
	Weaknesses []string `json:"weaknesses"`  // 缺点
	CoreValues string   `json:"core_values"` // 核心价值观
	Fears      string   `json:"fears"`       // 恐惧
}

// BackgroundData 背景数据
type BackgroundData struct {
	Summary        string   `json:"summary"`         // 背景概要
	Family         string   `json:"family"`          // 家庭
	Education      string   `json:"education"`       // 教育
	KeyExperiences []string `json:"key_experiences"` // 关键经历
}

// RelationshipData 关系数据
type RelationshipData struct {
	Character    string `json:"character"`     // 关联角色
	RelationType string `json:"relation_type"` // 关系类型
	Description  string `json:"description"`   // 描述
	Dynamics     string `json:"dynamics"`      // 互动动态
}

// DevelopmentArc 发展弧线
type DevelopmentArc struct {
	StartingPoint string   `json:"starting_point"` // 起点
	TurningPoints []string `json:"turning_points"` // 转折点
	EndingPoint   string   `json:"ending_point"`   // 终点
	GrowthTheme   string   `json:"growth_theme"`   // 成长主题
}

// RelationshipNetwork 关系网络
type RelationshipNetwork struct {
	Alliances   [][]string       `json:"alliances"`   // 联盟（成员列表）
	Conflicts   [][]string       `json:"conflicts"`   // 冲突（对立方列表）
	Mentorships []MentorshipData `json:"mentorships"` // 师徒关系
}

// MentorshipData 师徒关系
type MentorshipData struct {
	Mentor  string `json:"mentor"`  // 导师
	Student string `json:"student"` // 学生
}

// ============ 情节生成 ============

// GeneratePlotRequest 生成情节请求
type GeneratePlotRequest struct {
	Task             string            `json:"task" binding:"required"` // 创作任务描述
	ProjectID        string            `json:"project_id"`              // 项目ID
	Outline          *OutlineData      `json:"outline"`                 // 大纲数据（建议提供）
	Characters       *CharactersData   `json:"characters"`              // 角色数据（建议提供）
	WorkspaceContext map[string]string `json:"workspace_context,omitempty"`
	CorrectionPrompt string            `json:"correction_prompt,omitempty"`
}

// GeneratePlotResponse 生成情节响应
type GeneratePlotResponse struct {
	Plot          *PlotData `json:"plot"`           // 情节数据
	ExecutionTime float32   `json:"execution_time"` // 执行时间
}

// PlotData 情节数据
type PlotData struct {
	TimelineEvents []TimelineEventData `json:"timeline_events"` // 时间线事件
	PlotThreads    []PlotThreadData    `json:"plot_threads"`    // 情节线索
	Conflicts      []ConflictData      `json:"conflicts"`       // 冲突
	KeyPlotPoints  *KeyPlotPoints      `json:"key_plot_points"` // 关键情节点
}

// TimelineEventData 时间线事件数据
type TimelineEventData struct {
	EventID      string       `json:"event_id"`     // 事件ID
	Timestamp    string       `json:"timestamp"`    // 时间戳
	Location     string       `json:"location"`     // 地点
	Title        string       `json:"title"`        // 标题
	Description  string       `json:"description"`  // 描述
	Participants []string     `json:"participants"` // 参与者
	EventType    string       `json:"event_type"`   // 事件类型
	Impact       *EventImpact `json:"impact"`       // 影响
	Causes       []string     `json:"causes"`       // 原因
	Consequences []string     `json:"consequences"` // 后果
	ChapterID    int32        `json:"chapter_id"`   // 所属章节
}

// EventImpact 事件影响
type EventImpact struct {
	OnPlot          string            `json:"on_plot"`          // 对情节的影响
	OnCharacters    map[string]string `json:"on_characters"`    // 对角色的影响
	EmotionalImpact string            `json:"emotional_impact"` // 情感影响
}

// PlotThreadData 情节线索数据
type PlotThreadData struct {
	ThreadID           string   `json:"thread_id"`           // 线索ID
	Title              string   `json:"title"`               // 标题
	Description        string   `json:"description"`         // 描述
	Type               string   `json:"type"`                // 类型（主线/支线）
	Events             []string `json:"events"`              // 相关事件ID列表
	StartingEvent      string   `json:"starting_event"`      // 起始事件
	ClimaxEvent        string   `json:"climax_event"`        // 高潮事件
	ResolutionEvent    string   `json:"resolution_event"`    // 解决事件
	CharactersInvolved []string `json:"characters_involved"` // 参与角色
}

// ConflictData 冲突数据
type ConflictData struct {
	ConflictID       string   `json:"conflict_id"`       // 冲突ID
	Type             string   `json:"type"`              // 类型
	Parties          []string `json:"parties"`           // 对立方
	Description      string   `json:"description"`       // 描述
	EscalationEvents []string `json:"escalation_events"` // 升级事件
	ResolutionEvent  string   `json:"resolution_event"`  // 解决事件
}

// KeyPlotPoints 关键情节点
type KeyPlotPoints struct {
	IncitingIncident string `json:"inciting_incident"` // 触发事件
	PlotPoint1       string `json:"plot_point_1"`      // 第一转折点
	Midpoint         string `json:"midpoint"`          // 中点
	PlotPoint2       string `json:"plot_point_2"`      // 第二转折点
	Climax           string `json:"climax"`            // 高潮
	Resolution       string `json:"resolution"`        // 结局
}

// ============ 完整工作流 ============

// ExecuteCreativeWorkflowRequest 执行创作工作流请求
type ExecuteCreativeWorkflowRequest struct {
	Task              string            `json:"task" binding:"required"` // 创作任务描述
	ProjectID         string            `json:"project_id"`              // 项目ID
	MaxReflections    int32             `json:"max_reflections"`         // 最大反思次数（默认3）
	EnableHumanReview bool              `json:"enable_human_review"`     // 是否启用人工审核
	WorkspaceContext  map[string]string `json:"workspace_context,omitempty"`
}

// ExecuteCreativeWorkflowResponse 执行创作工作流响应
type ExecuteCreativeWorkflowResponse struct {
	ExecutionID      string                `json:"execution_id"`      // 执行ID
	ReviewPassed     bool                  `json:"review_passed"`     // 审核是否通过
	ReflectionCount  int32                 `json:"reflection_count"`  // 反思次数
	Outline          *OutlineData          `json:"outline"`           // 大纲
	Characters       *CharactersData       `json:"characters"`        // 角色
	Plot             *PlotData             `json:"plot"`              // 情节
	DiagnosticReport *DiagnosticReportData `json:"diagnostic_report"` // 诊断报告
	Reasoning        []string              `json:"reasoning"`         // 推理链
	ExecutionTimes   map[string]float32    `json:"execution_times"`   // 执行时间
	TokensUsed       int32                 `json:"tokens_used"`       // Token使用量
}

// DiagnosticReportData 诊断报告数据
type DiagnosticReportData struct {
	Passed             bool              `json:"passed"`              // 是否通过
	QualityScore       int32             `json:"quality_score"`       // 质量分数
	Issues             []DiagnosticIssue `json:"issues"`              // 问题列表
	CorrectionStrategy string            `json:"correction_strategy"` // 修正策略
	AffectedAgents     []string          `json:"affected_agents"`     // 受影响的Agent
	ReasoningChain     []string          `json:"reasoning_chain"`     // 推理链
}

// DiagnosticIssue 诊断问题
type DiagnosticIssue struct {
	ID               string   `json:"id"`                // 问题ID
	Severity         string   `json:"severity"`          // 严重程度
	Category         string   `json:"category"`          // 类别
	SubCategory      string   `json:"sub_category"`      // 子类别
	Title            string   `json:"title"`             // 标题
	Description      string   `json:"description"`       // 描述
	RootCause        string   `json:"root_cause"`        // 根本原因
	AffectedEntities []string `json:"affected_entities"` // 受影响的实体
	Impact           string   `json:"impact"`            // 影响
}
