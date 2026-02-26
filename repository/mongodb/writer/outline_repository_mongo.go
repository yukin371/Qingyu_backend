package writer

import (
	"context"
	"fmt"
	"regexp"
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

var outlineQueryIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

func normalizeAndValidateOutlineQueryID(field, value string, allowEmpty bool) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		if allowEmpty {
			return "", nil
		}
		return "", errors.NewRepositoryError(errors.RepositoryErrorValidation, fmt.Sprintf("%s is required", field), nil)
	}
	if !outlineQueryIDPattern.MatchString(normalized) {
		return "", errors.NewRepositoryError(errors.RepositoryErrorValidation, fmt.Sprintf("invalid %s format", field), nil)
	}
	return normalized, nil
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

	_, err := r.GetCollection().InsertOne(ctx, outline)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create outline failed", err)
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

	filter := bson.M{"project_id": safeProjectID}

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

	filter := bson.M{
		"project_id": safeProjectID,
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

	filter := bson.M{
		"project_id": safeProjectID,
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

	filter := bson.M{"project_id": safeProjectID}
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

	var filter bson.M

	if safeParentID == "" {
		// 查询根节点（parent_id为空或不存在）
		filter = bson.M{
			"project_id": safeProjectID,
			"$or": []bson.M{
				{"parent_id": ""},
				{"parent_id": bson.M{"$exists": false}},
			},
		}
	} else {
		// 查询指定父节点的子节点
		filter = bson.M{
			"project_id": safeProjectID,
			"parent_id":  safeParentID,
		}
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.NewRepositoryError(errors.RepositoryErrorInternal, "count outline children failed", err)
	}
	return count, nil
}
