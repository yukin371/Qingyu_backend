package reader

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection 收藏记录
type Collection struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id" binding:"required"`
	BookID    string             `bson:"book_id" json:"book_id" binding:"required"`
	FolderID  string             `bson:"folder_id,omitempty" json:"folder_id,omitempty"`
	Tags      []string           `bson:"tags" json:"tags"`
	Note      string             `bson:"note" json:"note"`
	IsPublic  bool               `bson:"is_public" json:"is_public"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CollectionFolder 收藏夹
type CollectionFolder struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id" binding:"required"`
	Name        string             `bson:"name" json:"name" binding:"required,min=1,max=50"`
	Description string             `bson:"description" json:"description" binding:"max=200"`
	BookCount   int                `bson:"book_count" json:"book_count"`
	IsPublic    bool               `bson:"is_public" json:"is_public"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// TableName 返回Collection的MongoDB集合名称
func (Collection) TableName() string {
	return "collections"
}

// TableName 返回CollectionFolder的MongoDB集合名称
func (CollectionFolder) TableName() string {
	return "collection_folders"
}
