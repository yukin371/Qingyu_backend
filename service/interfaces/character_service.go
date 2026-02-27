package interfaces

import (
	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/writer"
	"context"
)

// CharacterService 角色服务接口
type CharacterService interface {
	// 基础CRUD
	Create(ctx context.Context, projectID, userID string, req *CreateCharacterRequest) (*writer.Character, error)
	GetByID(ctx context.Context, characterID, projectID string) (*writer.Character, error)
	List(ctx context.Context, projectID string) ([]*writer.Character, error)
	Update(ctx context.Context, characterID, projectID string, req *UpdateCharacterRequest) (*writer.Character, error)
	Delete(ctx context.Context, characterID, projectID string) error

	// 关系管理
	CreateRelation(ctx context.Context, projectID string, req *CreateRelationRequest) (*writer.CharacterRelation, error)
	ListRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error)
	DeleteRelation(ctx context.Context, relationID, projectID string) error

	// 关系图
	GetCharacterGraph(ctx context.Context, projectID string) (*CharacterGraph, error)
}

// CreateCharacterRequest 创建角色请求
// Deprecated: 使用 dto.CreateCharacterRequest 替代
type CreateCharacterRequest = dto.CreateCharacterRequest

// UpdateCharacterRequest 更新角色请求
// Deprecated: 使用 dto.UpdateCharacterRequest 替代
type UpdateCharacterRequest = dto.UpdateCharacterRequest

// CreateRelationRequest 创建关系请求
// Deprecated: 使用 dto.CreateRelationRequest 替代
type CreateRelationRequest = dto.CreateRelationRequest

// CharacterGraph 角色关系图
// Deprecated: 保留角色模块特定类型
type CharacterGraph struct {
	Nodes []*writer.Character         `json:"nodes"`
	Edges []*writer.CharacterRelation `json:"edges"`
}
