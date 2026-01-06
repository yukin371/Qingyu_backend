package writer

import (
	"Qingyu_backend/models/writer"
	"context"
)

// CharacterRepository 角色卡Repository接口
type CharacterRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, character *writer.Character) error
	FindByID(ctx context.Context, characterID string) (*writer.Character, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*writer.Character, error)
	Update(ctx context.Context, character *writer.Character) error
	Delete(ctx context.Context, characterID string) error

	// 关系管理
	CreateRelation(ctx context.Context, relation *writer.CharacterRelation) error
	FindRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error)
	FindRelationByID(ctx context.Context, relationID string) (*writer.CharacterRelation, error)
	DeleteRelation(ctx context.Context, relationID string) error

	// 辅助方法
	ExistsByID(ctx context.Context, characterID string) (bool, error)
	CountByProjectID(ctx context.Context, projectID string) (int64, error)
}
