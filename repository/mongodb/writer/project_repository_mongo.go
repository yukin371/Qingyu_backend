package writer

import (
	"Qingyu_backend/models/writer"
	"Qingyu_backend/repository/mongodb/base"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/repository/interfaces/infrastructure"
	writingInterface "Qingyu_backend/repository/interfaces/writer"
)

// MongoProjectRepository MongoDB项目仓储实现
type MongoProjectRepository struct {
	*base.BaseMongoRepository
	db *mongo.Database
}

// NewMongoProjectRepository 创建MongoDB项目仓储
func NewMongoProjectRepository(db *mongo.Database) writingInterface.ProjectRepository {
	return &MongoProjectRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "projects"),
		db:                 db,
	}
}

// Create 创建项目
func (r *MongoProjectRepository) Create(ctx context.Context, project *writer.Project) error {
	if project == nil {
		return fmt.Errorf("项目对象不能为空")
	}

	// 生成ID
	if project.ID.IsZero() {
		project.ID = primitive.NewObjectID()
	}

	// 设置时间戳
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	// 设置默认状态
	if project.Status == "" {
		project.Status = writer.StatusDraft
	}

	// 设置默认可见性
	if project.Visibility == "" {
		project.Visibility = writer.VisibilityPrivate
	}

	// 初始化统计信息
	project.Statistics = writer.ProjectStats{
		TotalWords:    0,
		ChapterCount:  0,
		DocumentCount: 0,
		LastUpdateAt:  now,
	}

	// 初始化设置
	project.Settings = writer.ProjectSettings{
		AutoBackup:     true,
		BackupInterval: 24,
	}

	// 验证数据
	if err := project.Validate(); err != nil {
		return fmt.Errorf("项目数据验证失败: %w", err)
	}

	// 插入数据库
	_, err := r.GetCollection().InsertOne(ctx, project)
	if err != nil {
		return fmt.Errorf("创建项目失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取项目
func (r *MongoProjectRepository) GetByID(ctx context.Context, id string) (*writer.Project, error) {
	var project writer.Project

	objectID, err := r.ParseID(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	filter := bson.M{
		"_id":        objectID,
		"deleted_at": nil, // 排除已删除的项目
	}

	err = r.GetCollection().FindOne(ctx, filter).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 项目不存在
		}
		return nil, fmt.Errorf("查询项目失败: %w", err)
	}

	return &project, nil
}

// Update 更新项目
func (r *MongoProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	// 自动更新updated_at
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": objectID, "deleted_at": nil}
	update := bson.M{"$set": updates}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新项目失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("项目不存在或已删除")
	}

	return nil
}

// Delete 物理删除项目
func (r *MongoProjectRepository) Delete(ctx context.Context, id string) error {
	objectID, err := r.ParseID(id)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	filter := bson.M{"_id": objectID}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("删除项目失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("项目不存在")
	}

	return nil
}

// List 查询项目列表
func (r *MongoProjectRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.Project, error) {
	mongoFilter := bson.M{"deleted_at": nil}

	// 如果有筛选条件，合并条件
	if filter != nil {
		conditions := filter.GetConditions()
		for key, value := range conditions {
			mongoFilter[key] = value
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})

	// 如果Filter提供了排序，使用Filter的排序
	if filter != nil && filter.GetSort() != nil {
		opts.SetSort(filter.GetSort())
	}

	cursor, err := r.GetCollection().Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询项目列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var projects []*writer.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, fmt.Errorf("解析项目数据失败: %w", err)
	}

	return projects, nil
}

// Exists 检查项目是否存在
func (r *MongoProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := r.ParseID(id)
	if err != nil {
		return false, fmt.Errorf("无效的ID格式: %w", err)
	}

	filter := bson.M{"_id": objectID, "deleted_at": nil}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查项目存在失败: %w", err)
	}

	return count > 0, nil
}

// Count 统计项目总数
func (r *MongoProjectRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	mongoFilter := bson.M{"deleted_at": nil}

	// 如果有筛选条件，合并条件
	if filter != nil {
		conditions := filter.GetConditions()
		for key, value := range conditions {
			mongoFilter[key] = value
		}
	}

	count, err := r.GetCollection().CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, fmt.Errorf("统计项目数失败: %w", err)
	}

	return count, nil
}

// GetListByOwnerID 获取作者的项目列表
func (r *MongoProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error) {
	filter := bson.M{
		"author_id":  ownerID,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询项目列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var projects []*writer.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, fmt.Errorf("解析项目数据失败: %w", err)
	}

	return projects, nil
}

// GetByOwnerAndStatus 根据作者和状态查询项目
func (r *MongoProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error) {
	filter := bson.M{
		"author_id":  ownerID,
		"status":     status,
		"deleted_at": nil,
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询项目列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var projects []*writer.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, fmt.Errorf("解析项目数据失败: %w", err)
	}

	return projects, nil
}

// UpdateByOwner 根据所有者更新项目
func (r *MongoProjectRepository) UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return fmt.Errorf("无效的项目ID格式: %w", err)
	}

	// 自动更新updated_at
	updates["updated_at"] = time.Now()

	filter := bson.M{
		"_id":        projectObjectID,
		"author_id":  ownerID,
		"deleted_at": nil,
	}
	update := bson.M{"$set": updates}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新项目失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("项目不存在或无权限")
	}

	return nil
}

// IsOwner 检查用户是否为项目所有者
func (r *MongoProjectRepository) IsOwner(ctx context.Context, projectID, ownerID string) (bool, error) {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return false, fmt.Errorf("无效的项目ID格式: %w", err)
	}

	filter := bson.M{
		"_id":        projectObjectID,
		"author_id":  ownerID,
		"deleted_at": nil,
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("检查所有者失败: %w", err)
	}

	return count > 0, nil
}

// SoftDelete 软删除项目
func (r *MongoProjectRepository) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return fmt.Errorf("无效的项目ID格式: %w", err)
	}

	now := time.Now()

	filter := bson.M{
		"_id":        projectObjectID,
		"author_id":  ownerID,
		"deleted_at": nil,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("软删除项目失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("项目不存在或无权限")
	}

	return nil
}

// HardDelete 物理删除项目
func (r *MongoProjectRepository) HardDelete(ctx context.Context, projectID string) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return fmt.Errorf("无效的项目ID格式: %w", err)
	}

	filter := bson.M{"_id": projectObjectID}

	result, err := r.GetCollection().DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("物理删除项目失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("项目不存在")
	}

	return nil
}

// Restore 恢复已删除的项目
func (r *MongoProjectRepository) Restore(ctx context.Context, projectID, ownerID string) error {
	projectObjectID, err := r.ParseID(projectID)
	if err != nil {
		return fmt.Errorf("无效的项目ID格式: %w", err)
	}

	filter := bson.M{
		"_id":       projectObjectID,
		"author_id": ownerID,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted_at": nil,
			"updated_at": time.Now(),
		},
	}

	result, err := r.GetCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("恢复项目失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("项目不存在或无权限")
	}

	return nil
}

// CountByOwner 统计作者的项目数
func (r *MongoProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	filter := bson.M{
		"author_id":  ownerID,
		"deleted_at": nil,
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计项目数失败: %w", err)
	}

	return count, nil
}

// CountByStatus 统计指定状态的项目数
func (r *MongoProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	filter := bson.M{
		"status":     status,
		"deleted_at": nil,
	}

	count, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计项目数失败: %w", err)
	}

	return count, nil
}

// CreateWithTransaction 在事务中创建项目
func (r *MongoProjectRepository) CreateWithTransaction(ctx context.Context, project *writer.Project, callback func(ctx context.Context) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return fmt.Errorf("启动事务失败: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 创建项目
		if err := r.Create(sessCtx, project); err != nil {
			return nil, err
		}

		// 执行回调
		if callback != nil {
			if err := callback(sessCtx); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("事务执行失败: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoProjectRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// EnsureIndexes 创建索引
func (r *MongoProjectRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "author_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "updated_at", Value: -1},
			},
		},
	}

	_, err := r.GetCollection().Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	return nil
}
