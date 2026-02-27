package admin

import (
	adminInterface "Qingyu_backend/repository/interfaces/admin"
	"context"
	"fmt"
	"regexp"
	"time"

	adminModel "Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AdminLogRepositoryImpl 管理员操作日志Repository的MongoDB实现
type AdminLogRepositoryImpl struct {
	collection *mongo.Collection
}

// NewAdminLogRepository 创建管理员操作日志Repository实例
func NewAdminLogRepository(db *mongo.Database) adminInterface.AdminLogRepository {
	return &AdminLogRepositoryImpl{
		collection: db.Collection("admin_logs"),
	}
}

// CreateAdminLog 创建管理员操作日志
func (r *AdminLogRepositoryImpl) CreateAdminLog(ctx context.Context, log *adminModel.AdminLog) error {
	if log.ID == "" {
		log.ID = primitive.NewObjectID().Hex()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("创建管理员日志失败: %w", err)
	}

	return nil
}

// GetAdminLog 根据ID获取管理员日志
func (r *AdminLogRepositoryImpl) GetAdminLog(ctx context.Context, logID string) (*adminModel.AdminLog, error) {
	var log adminModel.AdminLog
	err := r.collection.FindOne(ctx, bson.M{"_id": logID}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("管理员日志不存在")
		}
		return nil, fmt.Errorf("查询管理员日志失败: %w", err)
	}
	return &log, nil
}

// ListAdminLogs 列出管理员操作日志
func (r *AdminLogRepositoryImpl) ListAdminLogs(ctx context.Context, filter *adminInterface.AdminLogFilter) ([]*adminModel.AdminLog, error) {
	bsonFilter := bson.M{}

	if filter.AdminID != "" {
		bsonFilter["admin_id"] = exactMatchRegex(filter.AdminID)
	}
	if filter.Operation != "" {
		bsonFilter["operation"] = exactMatchRegex(filter.Operation)
	}
	// 新增：支持资源类型过滤
	if filter.ResourceType != "" {
		bsonFilter["resource_type"] = exactMatchRegex(filter.ResourceType)
	}
	// 新增：支持资源ID过滤
	if filter.ResourceID != "" {
		bsonFilter["resource_id"] = exactMatchRegex(filter.ResourceID)
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		createdAtFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			createdAtFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			createdAtFilter["$lte"] = filter.EndDate
		}
		bsonFilter["created_at"] = createdAtFilter
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(filter.Limit)

	if filter.Offset > 0 {
		opts.SetSkip(filter.Offset)
	}

	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询管理员日志失败: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*adminModel.AdminLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("解析管理员日志失败: %w", err)
	}

	return logs, nil
}

// CountAdminLogs 统计管理员日志数量
func (r *AdminLogRepositoryImpl) CountAdminLogs(ctx context.Context, filter *adminInterface.AdminLogFilter) (int64, error) {
	bsonFilter := bson.M{}

	if filter.AdminID != "" {
		bsonFilter["admin_id"] = exactMatchRegex(filter.AdminID)
	}
	if filter.Operation != "" {
		bsonFilter["operation"] = exactMatchRegex(filter.Operation)
	}
	// 新增：支持资源类型过滤
	if filter.ResourceType != "" {
		bsonFilter["resource_type"] = exactMatchRegex(filter.ResourceType)
	}
	// 新增：支持资源ID过滤
	if filter.ResourceID != "" {
		bsonFilter["resource_id"] = exactMatchRegex(filter.ResourceID)
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		createdAtFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			createdAtFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			createdAtFilter["$lte"] = filter.EndDate
		}
		bsonFilter["created_at"] = createdAtFilter
	}

	count, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, fmt.Errorf("统计管理员日志数量失败: %w", err)
	}

	return count, nil
}

// Health 健康检查
func (r *AdminLogRepositoryImpl) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}

// ============ 审计追踪新增方法 ============

// GetByResource 根据资源获取日志
func (r *AdminLogRepositoryImpl) GetByResource(ctx context.Context, resourceType, resourceID string) ([]*adminModel.AdminLog, error) {
	bsonFilter := bson.M{
		"resource_type": exactMatchRegex(resourceType),
		"resource_id":   exactMatchRegex(resourceID),
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询资源日志失败: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*adminModel.AdminLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("解析资源日志失败: %w", err)
	}

	return logs, nil
}

// GetByDateRange 根据日期范围获取日志
func (r *AdminLogRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*adminModel.AdminLog, error) {
	bsonFilter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询日期范围日志失败: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*adminModel.AdminLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("解析日期范围日志失败: %w", err)
	}

	return logs, nil
}

// CleanOldLogs 清理旧日志
func (r *AdminLogRepositoryImpl) CleanOldLogs(ctx context.Context, beforeDate time.Time) error {
	filter := bson.M{
		"created_at": bson.M{
			"$lt": beforeDate,
		},
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("清理旧日志失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("没有需要清理的旧日志")
	}

	return nil
}

func exactMatchRegex(value string) primitive.Regex {
	return primitive.Regex{
		Pattern: "^" + regexp.QuoteMeta(value) + "$",
		Options: "",
	}
}
