package base

import "go.mongodb.org/mongo-driver/bson/primitive"

// IdentifiedEntity ID实体混入
type IdentifiedEntity struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
}

// GenerateID 生成新ID
func (i *IdentifiedEntity) GenerateID() {
	if i.ID.IsZero() {
		i.ID = primitive.NewObjectID()
	}
}

// GetID 获取ID
func (i *IdentifiedEntity) GetID() primitive.ObjectID {
	return i.ID
}

// SetID 设置ID
func (i *IdentifiedEntity) SetID(id primitive.ObjectID) {
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
