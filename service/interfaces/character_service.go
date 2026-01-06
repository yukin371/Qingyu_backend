package interfaces

import (
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
type CreateCharacterRequest struct {
	Name              string   `json:"name" validate:"required"`
	Alias             []string `json:"alias"`
	Summary           string   `json:"summary"`
	Traits            []string `json:"traits"`
	Background        string   `json:"background"`
	AvatarURL         string   `json:"avatarUrl"`
	PersonalityPrompt string   `json:"personalityPrompt"`
	SpeechPattern     string   `json:"speechPattern"`
	CurrentState      string   `json:"currentState"`
}

// UpdateCharacterRequest 更新角色请求
type UpdateCharacterRequest struct {
	Name              *string   `json:"name"`
	Alias             *[]string `json:"alias"`
	Summary           *string   `json:"summary"`
	Traits            *[]string `json:"traits"`
	Background        *string   `json:"background"`
	AvatarURL         *string   `json:"avatarUrl"`
	PersonalityPrompt *string   `json:"personalityPrompt"`
	SpeechPattern     *string   `json:"speechPattern"`
	CurrentState      *string   `json:"currentState"`
}

// CreateRelationRequest 创建关系请求
type CreateRelationRequest struct {
	FromID   string `json:"fromId" validate:"required"`
	ToID     string `json:"toId" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Strength int    `json:"strength"`
	Notes    string `json:"notes"`
}

// CharacterGraph 角色关系图
type CharacterGraph struct {
	Nodes []*writer.Character         `json:"nodes"`
	Edges []*writer.CharacterRelation `json:"edges"`
}
