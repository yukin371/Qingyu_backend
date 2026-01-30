package writer

import (
	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository/mongodb/base"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	writerInterface "Qingyu_backend/repository/interfaces/writer"
)

// MongoTemplateRepository MongoDB模板仓储实现
type MongoTemplateRepository struct {
	*base.BaseMongoRepository
	db *mongo.Database
}

// NewMongoTemplateRepository 创建MongoDB模板仓储
func NewMongoTemplateRepository(db *mongo.Database) writerInterface.TemplateRepository {
	return &MongoTemplateRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "templates"),
		db:                 db,
	}
}

// Create 创建模板
func (r *MongoTemplateRepository) Create(ctx context.Context, template *writer.Template) error {
	if template == nil {
		return fmt.Errorf("模板对象不能为空")
	}

	// 设置默认值
	template.TouchForCreate()

	// 插入数据库
	_, err := r.GetCollection().InsertOne(ctx, template)
	if err != nil {
		return fmt.Errorf("创建模板失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取模板
func (r *MongoTemplateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*writer.Template, error) {
	var template writer.Template

	filter := bson.M{
		"_id": id,
		// 同时排除软删除的模板
		"deleted_at": nil,
		"status": bson.M{"$ne": writer.TemplateStatusDeleted},
	}

	err := r.GetCollection().FindOne(ctx, filter).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}

	return &template, nil
}

// Update 更新模板
func (r *MongoTemplateRepository) Update(ctx context.Context, template *writer.Template) error {
	if template == nil || template.ID.IsZero() {
		return fmt.Errorf("模板对象或ID不能为空")
	}

	// 更新版本和时间戳
	template.TouchForUpdate()

	filter := bson.M{
		"_id":        template.ID,
		"deleted_at": nil,
		"status":     bson.M{"$ne": writer.TemplateStatusDeleted},
	}

	update := bson.M{"$set": template}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新模板失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("模板不存在或已删除")
	}

	return nil
}

// Delete 删除模板（硬删除）
func (r *MongoTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("模板不存在")
	}

	return nil
}

// buildFilter 构建查询条件
func (r *MongoTemplateRepository) buildFilter(projectID *primitive.ObjectID, workspaceID string, isGlobal bool, filter *writerInterface.TemplateFilter) bson.M {
	mongoFilter := bson.M{
		// 排除软删除的模板
		"deleted_at": nil,
		"status":     bson.M{"$ne": writer.TemplateStatusDeleted},
	}

	// 项目过滤
	if isGlobal {
		// 全局模板：project_id为nil
		mongoFilter["project_id"] = nil
	} else if projectID != nil {
		// 项目模板：指定项目ID
		mongoFilter["project_id"] = *projectID
	}

	// 工作区过滤（如果有工作区ID，用于权限验证）
	// 注意：这里假设工作区ID在应用层处理，repository层不直接处理workspaceID
	// 如果需要工作区级别的过滤，可以在应用层添加额外逻辑

	// 模板类型过滤
	if filter != nil && filter.Type != nil {
		mongoFilter["type"] = *filter.Type
	}

	// 状态过滤
	if filter != nil && filter.Status != nil {
		mongoFilter["status"] = *filter.Status
	}

	// 分类过滤
	if filter != nil && filter.Category != "" {
		mongoFilter["category"] = filter.Category
	}

	// 系统模板过滤
	if filter != nil && filter.IsSystem != nil {
		mongoFilter["is_system"] = *filter.IsSystem
	}

	// 关键词搜索（名称或描述）
	if filter != nil && filter.Keyword != "" {
		keywordPattern := primitive.Regex{Pattern: filter.Keyword, Options: "i"}
		mongoFilter["$or"] = []bson.M{
			{"name": keywordPattern},
			{"description": keywordPattern},
		}
	}

	return mongoFilter
}

// buildSortOptions 构建排序选项
func (r *MongoTemplateRepository) buildSortOptions(filter *writerInterface.TemplateFilter) *options.FindOptions {
	opts := options.Find()

	// 默认按更新时间降序
	sortField := "updated_at"
	sortOrder := -1

	if filter != nil {
		// 排序字段
		if filter.SortBy != "" {
			switch filter.SortBy {
			case "name", "created_at", "updated_at":
				sortField = filter.SortBy
			}
		}

		// 排序顺序
		if filter.SortOrder != 0 {
			sortOrder = filter.SortOrder
		}
	}

	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// 分页
	if filter != nil {
		page := filter.Page
		if page < 1 {
			page = 1
		}

		pageSize := filter.PageSize
		if pageSize < 1 {
			pageSize = 20 // 默认每页20条
		}
		if pageSize > 100 {
			pageSize = 100 // 最大每页100条
		}

		skip := int64((page - 1) * pageSize)
		limit := int64(pageSize)
		opts.SetSkip(skip)
		opts.SetLimit(limit)
	} else {
		// 默认分页
		opts.SetSkip(0)
		opts.SetLimit(20)
	}

	return opts
}

// ListByProject 列出项目的所有模板
func (r *MongoTemplateRepository) ListByProject(ctx context.Context, projectID primitive.ObjectID, filter *writerInterface.TemplateFilter) ([]*writer.Template, error) {
	mongoFilter := r.buildFilter(&projectID, "", false, filter)
	opts := r.buildSortOptions(filter)

	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询项目模板列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*writer.Template
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析模板数据失败: %w", err)
	}

	return templates, nil
}

// ListByWorkspace 列出工作区的所有模板（包括全局模板）
func (r *MongoTemplateRepository) ListByWorkspace(ctx context.Context, workspaceID string, filter *writerInterface.TemplateFilter) ([]*writer.Template, error) {
	// 注意：workspaceID的使用方式取决于具体业务逻辑
	// 这里假设模板模型中可能需要workspace_id字段，或者在应用层处理工作区权限
	// 当前实现：返回所有模板（包括全局和项目模板），由应用层根据workspaceID过滤

	mongoFilter := bson.M{
		"deleted_at": nil,
		"status":     bson.M{"$ne": writer.TemplateStatusDeleted},
	}

	// 应用过滤器
	if filter != nil {
		if filter.Type != nil {
			mongoFilter["type"] = *filter.Type
		}
		if filter.Status != nil {
			mongoFilter["status"] = *filter.Status
		}
		if filter.Category != "" {
			mongoFilter["category"] = filter.Category
		}
		if filter.IsSystem != nil {
			mongoFilter["is_system"] = *filter.IsSystem
		}
		if filter.Keyword != "" {
			keywordPattern := primitive.Regex{Pattern: filter.Keyword, Options: "i"}
			mongoFilter["$or"] = []bson.M{
				{"name": keywordPattern},
				{"description": keywordPattern},
			}
		}
	}

	opts := r.buildSortOptions(filter)

	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询工作区模板列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*writer.Template
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析模板数据失败: %w", err)
	}

	return templates, nil
}

// ListGlobal 列出所有全局模板
func (r *MongoTemplateRepository) ListGlobal(ctx context.Context, filter *writerInterface.TemplateFilter) ([]*writer.Template, error) {
	mongoFilter := r.buildFilter(nil, "", true, filter)
	opts := r.buildSortOptions(filter)

	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询全局模板列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*writer.Template
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("解析模板数据失败: %w", err)
	}

	return templates, nil
}

// CountByProject 统计项目的模板数量
func (r *MongoTemplateRepository) CountByProject(ctx context.Context, projectID primitive.ObjectID) (int64, error) {
	filter := bson.M{
		"project_id": projectID,
		"deleted_at": nil,
		"status":     bson.M{"$ne": writer.TemplateStatusDeleted},
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计项目模板数量失败: %w", err)
	}

	return count, nil
}

// ExistsByName 检查同名模板是否存在
func (r *MongoTemplateRepository) ExistsByName(ctx context.Context, projectID *primitive.ObjectID, name string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"name":       name,
		"deleted_at": nil,
		"status":     bson.M{"$ne": writer.TemplateStatusDeleted},
	}

	// 项目ID过滤
	if projectID != nil {
		filter["project_id"] = *projectID
	} else {
		// 全局模板
		filter["project_id"] = nil
	}

	// 排除指定ID（用于更新时检查）
	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查模板名称是否存在失败: %w", err)
	}

	return count > 0, nil
}

// EnsureIndexes 创建索引
//
// 建议的索引配置：
// 1. project_id: 用于快速查询项目模板
// 2. {project_id: 1, type: 1}: 用于按项目+类型查询
// 3. {project_id: 1, status: 1}: 用于按项目+状态查询
// 4. {project_id: 1, category: 1}: 用于按项目+分类查询
// 5. {project_id: 1, name: 1}: 用于按项目+名称查询（唯一索引）
// 6. is_system: 用于查询系统模板
// 7. {deleted_at: 1, status: 1}: 用于过滤已删除的模板
// 8. created_at: 用于按创建时间排序
// 9. updated_at: 用于按更新时间排序
// 10. {name: "text", description: "text"}: 用于全文搜索
func (r *MongoTemplateRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "project_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "type", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "category", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "project_id", Value: 1},
				{Key: "name", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
				"deleted_at": nil,
				"status":     bson.M{"$ne": writer.TemplateStatusDeleted},
			}),
		},
		{
			Keys: bson.D{{Key: "is_system", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "deleted_at", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "description", Value: "text"},
			},
		},
	}

	_, err := r.GetCollection().Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoTemplateRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}
