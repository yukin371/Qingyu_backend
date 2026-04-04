package writer

import (
	"context"
	"fmt"
	"strings"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/pkg/errors"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/repository/mongodb/base"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func normalizeAndValidateOutlineQueryID(field, value string, allowEmpty bool) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		if allowEmpty {
			return "", nil
		}
		return "", errors.NewRepositoryError(errors.RepositoryErrorValidation, fmt.Sprintf("%s is required", field), nil)
	}
	objectID, err := primitive.ObjectIDFromHex(normalized)
	if err != nil {
		return "", errors.NewRepositoryError(errors.RepositoryErrorValidation, fmt.Sprintf("invalid %s format", field), nil)
	}
	// 返回标准化的hex字符串，避免不同大小写/格式带来的查询歧义。
	return objectID.Hex(), nil
}

// OutlineRepositoryMongo Outline Repository的MongoDB实现
type OutlineRepositoryMongo struct {
	*base.BaseMongoRepository
}

// NewOutlineRepository 创建OutlineRepository实例
func NewOutlineRepository(db *mongo.Database) writerRepo.OutlineRepository {
	return &OutlineRepositoryMongo{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "outlines"),
	}
}

// Create 创建大纲节点
func (r *OutlineRepositoryMongo) Create(ctx context.Context, outline *writer.OutlineNode) error {
	// 设置默认值
	outline.TouchForCreate()

	result, err := r.GetCollection().InsertOne(ctx, outline)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create outline failed", err)
	}

	// 将插入后的ID设置回outline对象，确保后续操作能正确使用
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		outline.ID = oid
	}

	return nil
}

// FindByID 根据ID查询大纲节点
func (r *OutlineRepositoryMongo) FindByID(ctx context.Context, outlineID string) (*writer.OutlineNode, error) {
	var outline writer.OutlineNode

	objectID, err := primitive.ObjectIDFromHex(outlineID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid outline id format", err)
	}

	filter := bson.M{"_id": objectID}

	err = r.GetCollection().FindOne(ctx, filter).Decode(&outline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "outline not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find outline failed", err)
	}

	return &outline, nil
}

// FindByProjectID 查询项目下的所有大纲节点
func (r *OutlineRepositoryMongo) FindByProjectID(ctx context.Context, projectID string) ([]*writer.OutlineNode, error) {
	safeProjectID, err := normalizeAndValidateOutlineQueryID("project id", projectID, false)
	if err != nil {
		return nil, err
	}
	projectObjectID, _ := primitive.ObjectIDFromHex(safeProjectID)

	filter := bson.M{"project_id": projectObjectID}

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find outlines failed", err)
	}
	defer cursor.Close(ctx)

	var outlines []*writer.OutlineNode
	if err = cursor.All(ctx, &outlines); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode outlines failed", err)
	}

	return outlines, nil
}

// Update 更新大纲节点
func (r *OutlineRepositoryMongo) Update(ctx context.Context, outline *writer.OutlineNode) error {
	// 设置更新时间
	outline.TouchForUpdate()

	filter := bson.M{"_id": outline.ID}
	update := bson.M{"$set": outline}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update outline failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "outline not found", nil)
	}

	return nil
}

// Delete 删除大纲节点
func (r *OutlineRepositoryMongo) Delete(ctx context.Context, outlineID string) error {
	objectID, err := primitive.ObjectIDFromHex(outlineID)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid outline id format", err)
	}

	filter := bson.M{"_id": objectID}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete outline failed", err)
	}

	if result.DeletedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "outline not found", nil)
	}

	return nil
}

// FindByParentID 根据父节点ID查询子节点
func (r *OutlineRepositoryMongo) FindByParentID(ctx context.Context, projectID, parentID string) ([]*writer.OutlineNode, error) {
	safeProjectID, err := normalizeAndValidateOutlineQueryID("project id", projectID, false)
	if err != nil {
		return nil, err
	}
	safeParentID, err := normalizeAndValidateOutlineQueryID("parent id", parentID, true)
	if err != nil {
		return nil, err
	}
	projectObjectID, _ := primitive.ObjectIDFromHex(safeProjectID)

	filter := bson.M{
		"project_id": projectObjectID,
		"parent_id":  safeParentID,
	}

	// 按order排序
	opts := options.Find().SetSort(bson.M{"order": 1})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find outline children failed", err)
	}
	defer cursor.Close(ctx)

	var outlines []*writer.OutlineNode
	if err = cursor.All(ctx, &outlines); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode outline children failed", err)
	}

	return outlines, nil
}

// FindRoots 查询项目的所有根节点（parent_id为空的节点）
func (r *OutlineRepositoryMongo) FindRoots(ctx context.Context, projectID string) ([]*writer.OutlineNode, error) {
	safeProjectID, err := normalizeAndValidateOutlineQueryID("project id", projectID, false)
	if err != nil {
		return nil, err
	}
	projectObjectID, _ := primitive.ObjectIDFromHex(safeProjectID)

	filter := bson.M{
		"project_id": projectObjectID,
		"$or": []bson.M{
			{"parent_id": ""},
			{"parent_id": bson.M{"$exists": false}},
		},
	}

	// 按order排序
	opts := options.Find().SetSort(bson.M{"order": 1})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find outline roots failed", err)
	}
	defer cursor.Close(ctx)

	var outlines []*writer.OutlineNode
	if err = cursor.All(ctx, &outlines); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode outline roots failed", err)
	}

	return outlines, nil
}

// FindByDocumentID 根据关联文档ID查询大纲节点（用于双向同步去重）
func (r *OutlineRepositoryMongo) FindByDocumentID(ctx context.Context, documentID string) (*writer.OutlineNode, error) {
	filter := bson.M{
		"document_id": documentID,
	}

	var outline writer.OutlineNode
	err := r.GetCollection().FindOne(ctx, filter).Decode(&outline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 没有找到，返回 nil 而不是错误
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find outline by document_id failed", err)
	}

	return &outline, nil
}

// ExistsByID 检查大纲节点是否存在
func (r *OutlineRepositoryMongo) ExistsByID(ctx context.Context, outlineID string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(outlineID)
	if err != nil {
		return false, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid outline id format", err)
	}

	filter := bson.M{"_id": objectID}
	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.NewRepositoryError(errors.RepositoryErrorInternal, "check outline exists failed", err)
	}
	return count > 0, nil
}

// CountByProjectID 统计项目下的大纲节点数量
func (r *OutlineRepositoryMongo) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	safeProjectID, err := normalizeAndValidateOutlineQueryID("project id", projectID, false)
	if err != nil {
		return 0, err
	}
	projectObjectID, _ := primitive.ObjectIDFromHex(safeProjectID)

	filter := bson.M{"project_id": projectObjectID}
	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorInternal, "count outlines failed", err)
	}
	return count, nil
}

// CountByParentID 统计指定父节点下的子节点数量
func (r *OutlineRepositoryMongo) CountByParentID(ctx context.Context, projectID, parentID string) (int64, error) {
	safeProjectID, err := normalizeAndValidateOutlineQueryID("project id", projectID, false)
	if err != nil {
		return 0, err
	}
	safeParentID, err := normalizeAndValidateOutlineQueryID("parent id", parentID, true)
	if err != nil {
		return 0, err
	}
	projectObjectID, _ := primitive.ObjectIDFromHex(safeProjectID)

	var filter bson.M

	if safeParentID == "" {
		// 查询根节点（parent_id为空或不存在）
		filter = bson.M{
			"project_id": projectObjectID,
			"$or": []bson.M{
				{"parent_id": ""},
				{"parent_id": bson.M{"$exists": false}},
			},
		}
	} else {
		// 查询指定父节点的子节点
		filter = bson.M{
			"project_id": projectObjectID,
			"parent_id":  safeParentID,
		}
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorInternal, "count outline children failed", err)
	}
	return count, nil
}

// FindByGlobalOutline 查找项目的全局总纲节点（type="global", parent_id为空或不存在, 无 document_id）
// 如果不存在则原子性地创建（使用 findOneAndUpdate + upsert）
func (r *OutlineRepositoryMongo) FindByGlobalOutline(ctx context.Context, projectID primitive.ObjectID) (*writer.OutlineNode, error) {
	// 使用 findOneAndUpdate + upsert 实现原子性的"查找或创建"
	// filter 只使用精确匹配的字段来确保最多只有一个匹配
	filter := bson.M{
		"project_id":  projectID,
		"type":        "global",
		"document_id": "",
	}

	// 使用 $setOnInsert 在插入时设置默认值
	update := bson.M{
		"$setOnInsert": bson.M{
			"title":     "总纲",
			"parent_id": "", // 对于 upsert，必须设置一个具体的值
			"order":     0,
			"tags":      []string{},
			"tension":   0,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var outline writer.OutlineNode
	err := r.GetCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&outline)
	if err == nil {
		// 成功（可能是找到了已有的，或者是创建了新的）
		return &outline, nil
	}

	if err == mongo.ErrNoDocuments {
		// findOneAndUpdate + upsert 模式下理论上不会返回 ErrNoDocuments
		// 但为了安全，再查一次
		err = r.GetCollection().FindOne(ctx, filter).Decode(&outline)
		if err == nil {
			return &outline, nil
		}
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find global outline after upsert failed", err)
	}

	return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "findOrCreate global outline failed", err)
}

// EnsureIndexes 创建索引
func (r *OutlineRepositoryMongo) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		// project_id + parent_id 索引，用于查询子节点
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "parent_id", Value: 1},
			},
			Options: options.Index(),
		},
		// document_id 索引，用于双向引用查询
		{
			Keys: bson.D{
				{Key: "document_id", Value: 1},
			},
			Options: options.Index().SetSparse(true),
		},
	}

	_, err := r.GetCollection().Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create outline indexes failed", err)
	}

	return nil
}
