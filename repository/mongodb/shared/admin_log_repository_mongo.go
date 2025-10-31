package shared

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/service/shared/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AdminLogRepositoryImpl AdminLog Repository 的 MongoDB 实现
type AdminLogRepositoryImpl struct {
	collection *mongo.Collection
}

// NewAdminLogRepository 创建 AdminLog Repository 实例
func NewAdminLogRepository(db *mongo.Database) *AdminLogRepositoryImpl {
	return &AdminLogRepositoryImpl{
		collection: db.Collection("admin_logs"),
	}
}

// Create 创建管理员操作日志
func (r *AdminLogRepositoryImpl) Create(ctx context.Context, log *admin.AdminLog) error {
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

// List 查询管理员操作日志
func (r *AdminLogRepositoryImpl) List(ctx context.Context, filter *admin.LogFilter) ([]*admin.AdminLog, error) {
	// 构建过滤条件
	bsonFilter := bson.M{}

	if filter.AdminID != "" {
		bsonFilter["admin_id"] = filter.AdminID
	}
	if filter.Operation != "" {
		bsonFilter["operation"] = filter.Operation
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

	// 设置分页
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	// 设置查询选项
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // 按创建时间倒序
	opts.SetSkip(skip)
	opts.SetLimit(limit)

	// 执行查询
	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询管理员日志失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var logs []*admin.AdminLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("解析管理员日志失败: %w", err)
	}

	return logs, nil
}

// Health 健康检查
func (r *AdminLogRepositoryImpl) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
