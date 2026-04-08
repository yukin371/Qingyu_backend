package interfaces

import (
	"context"

	"Qingyu_backend/models/writer"
)

// EntitySummary 统一实体摘要
type EntitySummary struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	EntityType  writer.EntityType            `json:"entityType"`
	Summary     string                       `json:"summary,omitempty"`
	StateFields map[string]writer.StateValue `json:"stateFields,omitempty"`
}

// EntityGraph 统一实体图谱
type EntityGraph struct {
	Nodes []EntitySummary `json:"nodes"`
	Edges []RelationEdge  `json:"edges"`
}

// RelationEdge 统一关系边
type RelationEdge struct {
	FromID   string            `json:"fromId"`
	ToID     string            `json:"toId"`
	FromType writer.EntityType `json:"fromType"`
	ToType   writer.EntityType `json:"toType"`
	Type     string            `json:"type"`
	Strength int               `json:"strength,omitempty"`
	Notes    string            `json:"notes,omitempty"`
}

// EntityService 统一实体服务接口
type EntityService interface {
	// ListEntities 查询项目下所有实体（支持按 type 筛选）
	ListEntities(ctx context.Context, projectID string, entityType *string) ([]EntitySummary, error)
	// GetEntityGraph 统一实体图谱
	GetEntityGraph(ctx context.Context, projectID string) (*EntityGraph, error)
	// UpdateEntityStateFields 更新实体状态字段
	UpdateEntityStateFields(ctx context.Context, entityID string, stateFields map[string]writer.StateValue) error
}
