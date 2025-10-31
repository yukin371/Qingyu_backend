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

// AdminAuditRepositoryImpl 管理员审核记录 Repository 的 MongoDB 实现
type AdminAuditRepositoryImpl struct {
	collection *mongo.Collection
}

// NewAdminAuditRepository 创建管理员审核记录 Repository 实例
func NewAdminAuditRepository(db *mongo.Database) *AdminAuditRepositoryImpl {
	return &AdminAuditRepositoryImpl{
		collection: db.Collection("admin_audit_records"),
	}
}

// Create 创建审核记录
func (r *AdminAuditRepositoryImpl) Create(ctx context.Context, record *admin.AuditRecord) error {
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

// Get 根据ID获取审核记录
func (r *AdminAuditRepositoryImpl) Get(ctx context.Context, recordID string) (*admin.AuditRecord, error) {
	var record admin.AuditRecord
	err := r.collection.FindOne(ctx, bson.M{"_id": recordID}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("审核记录不存在")
		}
		return nil, fmt.Errorf("查询审核记录失败: %w", err)
	}
	return &record, nil
}

// Update 更新审核记录
func (r *AdminAuditRepositoryImpl) Update(ctx context.Context, recordID string, updates map[string]interface{}) error {
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

// ListByStatus 根据状态查询审核记录
func (r *AdminAuditRepositoryImpl) ListByStatus(ctx context.Context, contentType, status string) ([]*admin.AuditRecord, error) {
	filter := bson.M{}
	if contentType != "" {
		filter["content_type"] = contentType
	}
	if status != "" {
		filter["status"] = status
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(100)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询审核记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var records []*admin.AuditRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("解析审核记录失败: %w", err)
	}

	return records, nil
}

// ListByContent 根据内容ID查询审核记录
func (r *AdminAuditRepositoryImpl) ListByContent(ctx context.Context, contentID string) ([]*admin.AuditRecord, error) {
	filter := bson.M{"content_id": contentID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询审核记录失败: %w", err)
	}
	defer cursor.Close(ctx)

	var records []*admin.AuditRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("解析审核记录失败: %w", err)
	}

	return records, nil
}

// Health 健康检查
func (r *AdminAuditRepositoryImpl) Health(ctx context.Context) error {
	return r.collection.Database().Client().Ping(ctx, nil)
}
