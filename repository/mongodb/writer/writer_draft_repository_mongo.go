package writer

import (
	"context"
	"fmt"
	"strings"

	"Qingyu_backend/models/writer"
	writerInterface "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/repository/mongodb/base"

	"Qingyu_backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// normalizeAndValidateDraftQueryID 规范化并验证草稿查询ID
//
// 确保ID格式正确，避免不同大小写/格式带来的查询歧义喵~
func normalizeAndValidateDraftQueryID(field, value string, allowEmpty bool) (string, error) {
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
	// 返回标准化的hex字符串
	return objectID.Hex(), nil
}

// WriterDraftRepositoryMongo WriterDraft Repository的MongoDB实现
type WriterDraftRepositoryMongo struct {
	*base.BaseMongoRepository
}

// NewWriterDraftRepository 创建WriterDraftRepository实例
//
// 参数：
//   - db: MongoDB数据库实例
//
// 返回：
//   - writerInterface.WriterDraftRepository: Repository接口实例
func NewWriterDraftRepository(db *mongo.Database) writerInterface.WriterDraftRepository {
	return &WriterDraftRepositoryMongo{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "writer_drafts"),
	}
}

// Create 创建草稿
//
// 在创建前会自动设置默认值（时间戳、格式、状态等）喵~
func (r *WriterDraftRepositoryMongo) Create(ctx context.Context, doc *writer.WriterDraft) error {
	// 设置默认值
	doc.BeforeCreate()

	_, err := r.GetCollection().InsertOne(ctx, doc)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create writer draft failed", err)
	}

	return nil
}

// GetByID 根据ID获取草稿
func (r *WriterDraftRepositoryMongo) GetByID(ctx context.Context, id string) (*writer.WriterDraft, error) {
	var doc writer.WriterDraft

	err := r.FindByID(ctx, id, &doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "writer draft not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find writer draft failed", err)
	}

	return &doc, nil
}

// GetByProjectAndChapter 根据项目ID和章节号获取草稿
//
// 用于获取特定项目的特定章节草稿，确保每个章节只有一个草稿喵~
func (r *WriterDraftRepositoryMongo) GetByProjectAndChapter(ctx context.Context, projectID string, chapterNum int) (*writer.WriterDraft, error) {
	safeProjectID, err := normalizeAndValidateDraftQueryID("project id", projectID, false)
	if err != nil {
		return nil, err
	}

	var doc writer.WriterDraft
	filter := bson.M{
		"project_id":  safeProjectID,
		"chapter_num": chapterNum,
	}

	err = r.GetCollection().FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "writer draft not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find writer draft failed", err)
	}

	return &doc, nil
}

// ListByProject 获取项目的草稿列表
//
// 按章节号排序返回，限制返回数量喵~
func (r *WriterDraftRepositoryMongo) ListByProject(ctx context.Context, projectID string, limit int) ([]*writer.WriterDraft, error) {
	safeProjectID, err := normalizeAndValidateDraftQueryID("project id", projectID, false)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"project_id": safeProjectID,
	}

	// 按章节号排序
	opts := options.Find().SetSort(bson.M{"chapter_num": 1})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find writer drafts failed", err)
	}
	defer cursor.Close(ctx)

	var docs []*writer.WriterDraft
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode writer drafts failed", err)
	}

	return docs, nil
}

// Update 更新草稿
//
// 更新前会自动刷新时间戳并增加版本号喵~
func (r *WriterDraftRepositoryMongo) Update(ctx context.Context, doc *writer.WriterDraft) error {
	// 设置更新时间
	doc.BeforeUpdate()

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update writer draft failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "writer draft not found", nil)
	}

	return nil
}

// Delete 删除草稿
func (r *WriterDraftRepositoryMongo) Delete(ctx context.Context, id string) error {
	err := r.DeleteByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "writer draft not found", err)
		}
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete writer draft failed", err)
	}

	return nil
}

// BatchGetByIDs 批量获取草稿
//
// 支持通过ID列表批量获取多个草稿喵~
func (r *WriterDraftRepositoryMongo) BatchGetByIDs(ctx context.Context, ids []string) ([]*writer.WriterDraft, error) {
	if len(ids) == 0 {
		return []*writer.WriterDraft{}, nil
	}

	objectIDs, err := r.ParseIDs(ids)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid draft id format", err)
	}

	filter := bson.M{
		"_id": bson.M{"$in": objectIDs},
	}

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find writer drafts failed", err)
	}
	defer cursor.Close(ctx)

	var docs []*writer.WriterDraft
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode writer drafts failed", err)
	}

	return docs, nil
}
