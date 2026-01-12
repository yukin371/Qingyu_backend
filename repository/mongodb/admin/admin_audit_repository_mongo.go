package admin

import (
	adminInterface "Qingyu_backend/repository/interfaces/admin"
	"context"
	"fmt"
	"time"

	adminModel "Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditRepositoryImpl 审核记录Repository的MongoDB实现
type AuditRepositoryImpl struct {
	collection *mongo.Collection
}

// NewAuditRepository 创建审核记录Repository实例
func NewAuditRepository(db *mongo.Database) adminInterface.AuditRepository {
	return &AuditRepositoryImpl{
		collection: db.Collection("admin_audit_records"),
	}
}

// CreateAuditRecord 创建审核记录
func (r *AuditRepositoryImpl) CreateAuditRecord(ctx context.Context, record *adminModel.AuditRecord) error {
	if record.ID == "" {
		record.ID = primitive.NewObjectID().Hex()
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}
	if record.UpdatedAt.IsZero() {
		record.UpdatedAt = time.Now()
	}

	_, err := r.collection.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("创建审核记录失败: %w", err)
	}

	return nil
}

// GetAuditRecord 根据内容ID和类型获取审核记录
func (r *AuditRepositoryImpl) GetAuditRecord(ctx context.Context, contentID, contentType string) (*adminModel.AuditRecord, error) {
	filter := bson.M{"content_id": contentID}
	if contentType != "" {
		filter["content_type"] = contentType
	}

	var record adminModel.AuditRecord
	err := r.collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("审核记录不存在")
		}
		return nil, fmt.Errorf("查询审核记录失败: %w", err)
	}
	return &record, nil
}

// UpdateAuditRecord 更新审核记录
func (r *AuditRepositoryImpl) UpdateAuditRecord(ctx context.Context, recordID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": recordID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return fmt.Errorf("更新审核记录失败: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("审核记录不存在")
	}

	return nil
}

// ListAuditRecords 列出审核记录
func (r *AuditRepositoryImpl) ListAuditRecords(ctx context.Context, filter *adminInterface.AuditFilter) ([]*adminModel.AuditRecord, error) {
	bsonFilter := bson.M{}

	if filter.ContentType != "" {
		bsonFilter["content_type"] = filter.ContentType
	}
	if filter.Status != "" {
		bsonFilter["status"] = filter.Status
	}
	if filter.ReviewerID != "" {
		bsonFilter["reviewer_id"] = filter.ReviewerID
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
		return nil, fmt.Errorf("查询审核记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var records []*adminModel.AuditRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("解析审核记录失败: %w", err)
	}

	return records, nil
}

// Health 健康检查
func (r *AuditRepositoryImpl) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
