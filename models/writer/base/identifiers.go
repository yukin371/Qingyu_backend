package base

import "go.mongodb.org/mongo-driver/bson/primitive"

// IdentifiedEntity ID实体混入
type IdentifiedEntity struct {
	ID string `bson:"_id,omitempty" json:"id"`
}

// GenerateID 生成新ID
func (i *IdentifiedEntity) GenerateID() {
	if i.ID == "" {
		i.ID = primitive.NewObjectID().Hex()
	}
}

// GetID 获取ID
func (i *IdentifiedEntity) GetID() string {
	return i.ID
}

// SetID 设置ID
func (i *IdentifiedEntity) SetID(id string) {
	i.ID = id
}

// OwnedEntity 所属实体混入
type OwnedEntity struct {
	AuthorID string `bson:"author_id" json:"authorId" validate:"required"`
}

// ProjectScopedEntity 项目范围实体混入
type ProjectScopedEntity struct {
	ProjectID string `bson:"project_id" json:"projectId" validate:"required"`
}

// UserScopedEntity 用户范围实体混入
type UserScopedEntity struct {
	UserID string `bson:"user_id" json:"userId" validate:"required"`
}
