package mongodb

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	"Qingyu_backend/repository/interfaces/writing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CharacterRepositoryMongo Character Repository的MongoDB实现
type CharacterRepositoryMongo struct {
	db                  *mongo.Database
	characterCollection *mongo.Collection
	relationCollection  *mongo.Collection
}

// NewCharacterRepository 创建CharacterRepository实例
func NewCharacterRepository(db *mongo.Database) writing.CharacterRepository {
	return &CharacterRepositoryMongo{
		db:                  db,
		characterCollection: db.Collection("characters"),
		relationCollection:  db.Collection("character_relations"),
	}
}

// Create 创建角色
func (r *CharacterRepositoryMongo) Create(ctx context.Context, character *writer.Character) error {
	if character.ID == "" {
		character.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now()
	character.CreatedAt = now
	character.UpdatedAt = now

	_, err := r.characterCollection.InsertOne(ctx, character)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create character failed", err)
	}

	return nil
}

// FindByID 根据ID查询角色
func (r *CharacterRepositoryMongo) FindByID(ctx context.Context, characterID string) (*writer.Character, error) {
	var character writer.Character
	filter := bson.M{"_id": characterID}

	err := r.characterCollection.FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "character not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find character failed", err)
	}

	return &character, nil
}

// FindByProjectID 查询项目下的所有角色
func (r *CharacterRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) ([]*writer.Character, error) {
	filter := bson.M{"project_id": projectID}

	cursor, err := r.characterCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find characters failed", err)
	}
	defer cursor.Close(ctx)

	var characters []*writer.Character
	if err = cursor.All(ctx, &characters); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode characters failed", err)
	}

	return characters, nil
}

// Update 更新角色
func (r *CharacterRepositoryMongo) Update(ctx context.Context, character *writer.Character) error {
	character.UpdatedAt = time.Now()

	filter := bson.M{"_id": character.ID}
	update := bson.M{"$set": character}

	result, err := r.characterCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update character failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "character not found", nil)
	}

	return nil
}

// Delete 删除角色
func (r *CharacterRepositoryMongo) Delete(ctx context.Context, characterID string) error {
	filter := bson.M{"_id": characterID}

	result, err := r.characterCollection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete character failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "character not found", nil)
	}

	return nil
}

// CreateRelation 创建角色关系
func (r *CharacterRepositoryMongo) CreateRelation(ctx context.Context, relation *writer.CharacterRelation) error {
	if relation.ID == "" {
		relation.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now()
	relation.CreatedAt = now
	relation.UpdatedAt = now

	_, err := r.relationCollection.InsertOne(ctx, relation)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create character relation failed", err)
	}

	return nil
}

// FindRelations 查询角色关系
// 如果characterID为nil，返回项目下所有关系
// 如果characterID不为nil，返回该角色相关的所有关系
func (r *CharacterRepositoryMongo) FindRelations(ctx context.Context, projectID string, characterID *string) ([]*writer.CharacterRelation, error) {
	filter := bson.M{"project_id": projectID}

	if characterID != nil {
		filter["$or"] = []bson.M{
			{"from_id": *characterID},
			{"to_id": *characterID},
		}
	}

	cursor, err := r.relationCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find character relations failed", err)
	}
	defer cursor.Close(ctx)

	var relations []*writer.CharacterRelation
	if err = cursor.All(ctx, &relations); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode character relations failed", err)
	}

	return relations, nil
}

// FindRelationByID 根据ID查询关系
func (r *CharacterRepositoryMongo) FindRelationByID(ctx context.Context, relationID string) (*writer.CharacterRelation, error) {
	var relation writer.CharacterRelation
	filter := bson.M{"_id": relationID}

	err := r.relationCollection.FindOne(ctx, filter).Decode(&relation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "character relation not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find character relation failed", err)
	}

	return &relation, nil
}

// DeleteRelation 删除角色关系
func (r *CharacterRepositoryMongo) DeleteRelation(ctx context.Context, relationID string) error {
	filter := bson.M{"_id": relationID}

	result, err := r.relationCollection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete character relation failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "character relation not found", nil)
	}

	return nil
}

// ExistsByID 检查角色是否存在
func (r *CharacterRepositoryMongo) ExistsByID(ctx context.Context, characterID string) (bool, error) {
	filter := bson.M{"_id": characterID}
	count, err := r.characterCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.NewRepositoryError(errors.RepositoryErrorInternal, fmt.Sprintf("check character exists failed: %s", characterID), err)
	}
	return count > 0, nil
}

// CountByProjectID 统计项目下的角色数量
func (r *CharacterRepositoryMongo) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	filter := bson.M{"project_id": projectID}
	count, err := r.characterCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorInternal, "count characters failed", err)
	}
	return count, nil
}
