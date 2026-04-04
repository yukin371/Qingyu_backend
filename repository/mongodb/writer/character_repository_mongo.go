// --- 角色仓库MongoDB实现 ---
package writer

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CharacterRepositoryMongo Character Repository的MongoDB实现
type CharacterRepositoryMongo struct {
	*base.BaseMongoRepository // 嵌入基类，继承ID转换和通用CRUD方法
	db                        *mongo.Database
	relationCollection        *mongo.Collection // 关系collection独立管理
}

// NewCharacterRepository 创建CharacterRepository实例
func NewCharacterRepository(db *mongo.Database) writerRepo.CharacterRepository {
	return &CharacterRepositoryMongo{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "characters"),
		db:                  db,
		relationCollection:  db.Collection("character_relations"),
	}
}

// Create 创建角色
func (r *CharacterRepositoryMongo) Create(ctx context.Context, character *writer.Character) error {
	if character.ID.IsZero() {
		character.ID = primitive.NewObjectID()
	}

	now := time.Now()
	character.CreatedAt = now
	character.UpdatedAt = now

	_, err := r.GetCollection().InsertOne(ctx, character)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create character failed", err)
	}

	return nil
}

// FindByID 根据ID查询角色
func (r *CharacterRepositoryMongo) FindByID(ctx context.Context, characterID string) (*writer.Character, error) {
	var character writer.Character
	oid, err := r.ParseID(characterID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid character ID", err)
	}
	filter := bson.M{"_id": oid}

	err = r.GetCollection().FindOne(ctx, filter).Decode(&character)
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
	log.Printf("[FindByProjectID] 开始查询, projectID=%s", projectID)

	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		log.Printf("[FindByProjectID] ID转换失败: %v", err)
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	filter := bson.M{"project_id": projectObjectID}
	log.Printf("[FindByProjectID] 查询filter: %+v", filter)

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		log.Printf("[FindByProjectID] 查询失败: %v", err)
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find characters failed", err)
	}
	defer cursor.Close(ctx)

	var characters []*writer.Character
	if err = cursor.All(ctx, &characters); err != nil {
		log.Printf("[FindByProjectID] 解码失败: %v", err)
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode characters failed", err)
	}

	log.Printf("[FindByProjectID] 查询成功, 返回角色数量=%d", len(characters))

	return characters, nil
}

// Update 更新角色
func (r *CharacterRepositoryMongo) Update(ctx context.Context, character *writer.Character) error {
	character.UpdatedAt = time.Now()

	filter := bson.M{"_id": character.ID}
	update := bson.M{"$set": character}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
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
	oid, err := r.ParseID(characterID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid character ID", err)
	}
	filter := bson.M{"_id": oid}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
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
	if relation.ID.IsZero() {
		relation.ID = primitive.NewObjectID()
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
	log.Printf("[FindRelations] 开始查询, projectID=%s, characterID=%v", projectID, characterID)

	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		log.Printf("[FindRelations] ID转换失败: %v", err)
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}
	filter := bson.M{"project_id": projectObjectID}

	if characterID != nil {
		filter["$or"] = []bson.M{
			{"from_id": *characterID},
			{"to_id": *characterID},
		}
	}

	log.Printf("[FindRelations] 查询filter: %+v", filter)

	cursor, err := r.relationCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("[FindRelations] 查询失败: %v", err)
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find character relations failed", err)
	}
	defer cursor.Close(ctx)

	var relations []*writer.CharacterRelation
	if err = cursor.All(ctx, &relations); err != nil {
		log.Printf("[FindRelations] 解码失败: %v", err)
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode character relations failed", err)
	}

	log.Printf("[FindRelations] 查询成功, 返回关系数量=%d", len(relations))

	return relations, nil
}

// FindRelationByID 根据ID查询关系
func (r *CharacterRepositoryMongo) FindRelationByID(ctx context.Context, relationID string) (*writer.CharacterRelation, error) {
	var relation writer.CharacterRelation
	oid, err := r.ParseID(relationID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid relation ID", err)
	}
	filter := bson.M{"_id": oid}

	err = r.relationCollection.FindOne(ctx, filter).Decode(&relation)
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
	oid, err := r.ParseID(relationID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid relation ID", err)
	}
	filter := bson.M{"_id": oid}

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
	oid, err := r.ParseID(characterID)
	if err != nil {
		return false, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid character ID", err)
	}
	filter := bson.M{"_id": oid}
	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.NewRepositoryError(errors.RepositoryErrorInternal, fmt.Sprintf("check character exists failed: %s", characterID), err)
	}
	return count > 0, nil
}

// CountByProjectID 统计项目下的角色数量
func (r *CharacterRepositoryMongo) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}
	filter := bson.M{"project_id": projectObjectID}
	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorInternal, "count characters failed", err)
	}
	return count, nil
}

// CreateRelationTimelineEvent 为关系创建时序事件
func (r *CharacterRepositoryMongo) CreateRelationTimelineEvent(ctx context.Context, relationID string, event *writer.RelationTimelineEvent) error {
	objectID, err := primitive.ObjectIDFromHex(relationID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid relation id", err)
	}

	// 设置事件时间戳
	event.Timestamp = time.Now()

	// 使用 $push 添加到 timeline_events 数组
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$push": bson.M{
			"timeline_events": event,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.relationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create timeline event failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "relation not found", nil)
	}

	return nil
}

// GetRelationTimeline 获取关系的时序事件列表
func (r *CharacterRepositoryMongo) GetRelationTimeline(ctx context.Context, relationID string) ([]writer.RelationTimelineEvent, error) {
	objectID, err := primitive.ObjectIDFromHex(relationID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid relation id", err)
	}

	filter := bson.M{"_id": objectID}

	cursor, err := r.relationCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find timeline events failed", err)
	}
	defer cursor.Close(ctx)

	var relations []*writer.CharacterRelation
	if err = cursor.All(ctx, &relations); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode timeline events failed", err)
	}

	if len(relations) == 0 {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "relation not found", nil)
	}

	return relations[0].TimelineEvents, nil
}

// UpdateRelationTimelineEvent 更新关系的指定时序事件
func (r *CharacterRepositoryMongo) UpdateRelationTimelineEvent(ctx context.Context, relationID string, eventIndex int, event *writer.RelationTimelineEvent) error {
	objectID, err := primitive.ObjectIDFromHex(relationID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid relation id", err)
	}

	filter := bson.M{"_id": objectID}

	// 构建更新表达式，使用索引访问数组元素
	fieldName := fmt.Sprintf("timeline_events.%d", eventIndex)
	update := bson.M{
		"$set": bson.M{
			fieldName:     event,
			"updated_at": time.Now(),
		},
	}

	result, err := r.relationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update timeline event failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "relation not found", nil)
	}

	return nil
}

// DeleteRelationTimelineEvent 删除关系的指定时序事件
func (r *CharacterRepositoryMongo) DeleteRelationTimelineEvent(ctx context.Context, relationID string, eventIndex int) error {
	objectID, err := primitive.ObjectIDFromHex(relationID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid relation id", err)
	}

	filter := bson.M{"_id": objectID}

	// 使用 $unset 删除数组元素
	fieldName := fmt.Sprintf("timeline_events.%d", eventIndex)
	update := bson.M{
		"$unset": bson.M{
			fieldName: "",
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.relationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete timeline event failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "relation not found", nil)
	}

	// 清理数组中的 null 值
	cleanupFilter := bson.M{"_id": objectID}
	cleanupUpdate := bson.M{
		"$pull": bson.M{
			"timeline_events": nil,
		},
	}
	_, _ = r.relationCollection.UpdateOne(ctx, cleanupFilter, cleanupUpdate)

	return nil
}
