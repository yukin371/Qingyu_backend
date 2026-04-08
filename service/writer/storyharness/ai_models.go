package storyharness

// ChapterAnalysisRequest 对应 Python AI 服务的章节分析请求。
type ChapterAnalysisRequest struct {
	ProjectID        string               `json:"project_id"`
	ChapterID        string               `json:"chapter_id"`
	Text             string               `json:"text"`
	ExistingEntities []ChapterKnownEntity `json:"existing_entities,omitempty"`
	Model            string               `json:"model,omitempty"`
}

// ChapterKnownEntity 提供给 AI 的已知实体上下文。
type ChapterKnownEntity struct {
	EntityID     string   `json:"entity_id,omitempty"`
	EntityType   string   `json:"entity_type,omitempty"`
	Name         string   `json:"name,omitempty"`
	Aliases      []string `json:"aliases,omitempty"`
	Summary      string   `json:"summary,omitempty"`
	CurrentState string   `json:"current_state,omitempty"`
}

// ChapterAnalysisResponse 对应 Python AI 服务的章节分析响应。
type ChapterAnalysisResponse struct {
	StateChanges    []AIStateChange    `json:"state_changes"`
	RelationChanges []AIRelationChange `json:"relation_changes"`
	NewEntities     []AINewEntity      `json:"new_entities"`
	Events          []AIEventFact      `json:"events"`
	Usage           map[string]int     `json:"usage"`
}

// AIStateChange AI 返回的实体状态变化。
type AIStateChange struct {
	EntityName string `json:"entity_name"`
	FieldKey   string `json:"field_key"`
	OldValue   any    `json:"old_value"`
	NewValue   any    `json:"new_value"`
	Evidence   string `json:"evidence"`
}

// AIRelationChange AI 返回的关系变化。
type AIRelationChange struct {
	FromEntity   string `json:"from_entity"`
	ToEntity     string `json:"to_entity"`
	RelationType string `json:"relation_type"`
	ChangeType   string `json:"change_type"`
	Evidence     string `json:"evidence"`
}

// AINewEntity AI 返回的新实体出场。
type AINewEntity struct {
	Name         string `json:"name"`
	EntityType   string `json:"entity_type"`
	FirstMention string `json:"first_mention"`
	Description  string `json:"description"`
}

// AIEventFact AI 返回的事件事实。
type AIEventFact struct {
	Description      string   `json:"description"`
	ChapterPosition  string   `json:"chapter_position"`
	InvolvedEntities []string `json:"involved_entities"`
	Evidence         string   `json:"evidence"`
}
