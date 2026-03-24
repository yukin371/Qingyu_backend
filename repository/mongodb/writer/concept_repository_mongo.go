package writer

import (
	"context"
	"regexp"
	"strings"

	"Qingyu_backend/models/writer"
	writerInterface "Qingyu_backend/repository/interfaces/writer"
	"Qingyu_backend/repository/mongodb/base"

	"Qingyu_backend/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conceptCategoryPattern = regexp.MustCompile(`^[\p{L}\p{N}_\-\s]{1,64}$`)

func normalizeAndValidateConceptCategory(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return "", nil
	}
	if !conceptCategoryPattern.MatchString(normalized) {
		return "", errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid category format", nil)
	}
	return normalized, nil
}

// ConceptRepositoryMongo Concept Repository的MongoDB实现
type ConceptRepositoryMongo struct {
	*base.BaseMongoRepository
}

// NewConceptRepository 创建ConceptRepository实例
//
// 参数：
//   - db: MongoDB数据库实例
//
// 返回：
//   - writerInterface.ConceptRepository: Repository接口实例
func NewConceptRepository(db *mongo.Database) writerInterface.ConceptRepository {
	return &ConceptRepositoryMongo{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "concepts"),
	}
}

// Create 创建设定
//
// 在创建前会自动设置默认值（时间戳等）喵~
func (r *ConceptRepositoryMongo) Create(ctx context.Context, concept *writer.Concept) error {
	// 设置默认值
	concept.BeforeCreate()

	_, err := r.GetCollection().InsertOne(ctx, concept)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "create concept failed", err)
	}

	return nil
}

// GetByID 根据ID获取设定
func (r *ConceptRepositoryMongo) GetByID(ctx context.Context, id string) (*writer.Concept, error) {
	var concept writer.Concept

	err := r.FindByID(ctx, id, &concept)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewRepositoryError(errors.RepositoryErrorNotFound, "concept not found", err)
		}
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find concept failed", err)
	}

	return &concept, nil
}

// ListByProject 获取项目的设定列表
//
// 按分类和名称排序返回喵~
func (r *ConceptRepositoryMongo) ListByProject(ctx context.Context, projectID string) ([]*writer.Concept, error) {
	safeProjectID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}

	filter := bson.M{
		"project_id": safeProjectID,
	}

	// 按分类和名称排序
	opts := options.Find().SetSort(bson.M{"category": 1, "name": 1})

	// codeql[go/nosql-injection]: query fields are validated (project/category) and keyword matching is done in-memory
	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find concepts failed", err)
	}
	defer cursor.Close(ctx)

	var concepts []*writer.Concept
	if err = cursor.All(ctx, &concepts); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode concepts failed", err)
	}

	return concepts, nil
}

// Search 搜索设定
//
// 支持按分类和关键词搜索喵~
func (r *ConceptRepositoryMongo) Search(ctx context.Context, projectID, category, keyword string) ([]*writer.Concept, error) {
	safeProjectID, err := r.ParseID(projectID)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid project ID", err)
	}
	safeCategory, err := normalizeAndValidateConceptCategory(category)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"project_id": safeProjectID,
	}

	// 添加分类筛选
	if safeCategory != "" {
		filter["category"] = safeCategory
	}

	// 按分类和名称排序
	opts := options.Find().SetSort(bson.M{"category": 1, "name": 1})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "search concepts failed", err)
	}
	defer cursor.Close(ctx)

	var concepts []*writer.Concept
	if err = cursor.All(ctx, &concepts); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode concepts failed", err)
	}

	normalizedKeyword := strings.ToLower(strings.TrimSpace(keyword))
	if normalizedKeyword == "" {
		return concepts, nil
	}

	filtered := make([]*writer.Concept, 0, len(concepts))
	for _, concept := range concepts {
		name := strings.ToLower(concept.Name)
		content := strings.ToLower(concept.Content)
		if strings.Contains(name, normalizedKeyword) || strings.Contains(content, normalizedKeyword) {
			filtered = append(filtered, concept)
		}
	}

	return filtered, nil

}

// Update 更新设定
//
// 更新前会自动刷新时间戳喵~
func (r *ConceptRepositoryMongo) Update(ctx context.Context, concept *writer.Concept) error {
	// 设置更新时间
	concept.BeforeUpdate()

	filter := bson.M{"_id": concept.ID}
	update := bson.M{"$set": concept}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "update concept failed", err)
	}

	if result.MatchedCount == 0 {
		return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "concept not found", nil)
	}

	return nil
}

// Delete 删除设定
func (r *ConceptRepositoryMongo) Delete(ctx context.Context, id string) error {
	err := r.DeleteByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.NewRepositoryError(errors.RepositoryErrorNotFound, "concept not found", err)
		}
		return errors.NewRepositoryError(errors.RepositoryErrorInternal, "delete concept failed", err)
	}

	return nil
}

// BatchGetByIDs 批量获取设定
//
// 支持通过ID列表批量获取多个设定喵~
func (r *ConceptRepositoryMongo) BatchGetByIDs(ctx context.Context, ids []string) ([]*writer.Concept, error) {
	if len(ids) == 0 {
		return []*writer.Concept{}, nil
	}

	objectIDs, err := r.ParseIDs(ids)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorValidation, "invalid concept id format", err)
	}

	filter := bson.M{
		"_id": bson.M{"$in": objectIDs},
	}

	cursor, err := r.GetCollection().Find(ctx, filter)
	if err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "find concepts failed", err)
	}
	defer cursor.Close(ctx)

	var concepts []*writer.Concept
	if err = cursor.All(ctx, &concepts); err != nil {
		return nil, errors.NewRepositoryError(errors.RepositoryErrorInternal, "decode concepts failed", err)
	}

	return concepts, nil
}
