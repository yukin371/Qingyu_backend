package shared

import (
	adminModel "Qingyu_backend/models/users"
	"context"
	"fmt"
	"time"

	"Qingyu_backend/repository/interfaces/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAdminRepository MongoDB管理后台Repository实现
type MongoAdminRepository struct {
	db                  *mongo.Database
	auditCollection     *mongo.Collection
	adminLogsCollection *mongo.Collection
}

// NewMongoAdminRepository 创建MongoDB Admin Repository
func NewMongoAdminRepository(db *mongo.Database) shared.AdminRepository {
	return &MongoAdminRepository{
		db:                  db,
		auditCollection:     db.Collection("audit_records"),
		adminLogsCollection: db.Collection("admin_logs"),
	}
}

// ============ 审核记录 ============

// CreateAuditRecord 创建审核记录
func (r *MongoAdminRepository) CreateAuditRecord(ctx context.Context, record *adminModel.AuditRecord) error {
	if record == nil {
		return fmt.Errorf("audit record cannot be nil")
	}

	now := time.Now()
	record.CreatedAt = now
	record.UpdatedAt = now

	if record.ID == "" {
		record.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.auditCollection.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to create audit record: %w", err)
	}

	return nil
}

// GetAuditRecord 获取审核记录
func (r *MongoAdminRepository) GetAuditRecord(ctx context.Context, contentID, contentType string) (*adminModel.AuditRecord, error) {
	var record adminModel.AuditRecord
	query := bson.M{
		"content_id":   contentID,
		"content_type": contentType,
	}

	err := r.auditCollection.FindOne(ctx, query).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("audit record not found: %s/%s", contentType, contentID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get audit record: %w", err)
	}

	return &record, nil
}

// UpdateAuditRecord 更新审核记录
func (r *MongoAdminRepository) UpdateAuditRecord(ctx context.Context, recordID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.auditCollection.UpdateOne(
		ctx,
		bson.M{"_id": recordID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update audit record: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("audit record not found: %s", recordID)
	}

	return nil
}

// ListAuditRecords 获取审核记录列表
func (r *MongoAdminRepository) ListAuditRecords(ctx context.Context, filter *shared.AuditFilter) ([]*adminModel.AuditRecord, error) {
	query := bson.M{}

	if filter != nil {
		if filter.ContentType != "" {
			query["content_type"] = filter.ContentType
		}
		if filter.Status != "" {
			query["status"] = filter.Status
		}
		if filter.ReviewerID != "" {
			query["reviewer_id"] = filter.ReviewerID
		}
		if !filter.StartDate.IsZero() {
			query["created_at"] = bson.M{"$gte": filter.StartDate}
		}
		if !filter.EndDate.IsZero() {
			if query["created_at"] != nil {
				query["created_at"].(bson.M)["$lte"] = filter.EndDate
			} else {
				query["created_at"] = bson.M{"$lte": filter.EndDate}
			}
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if filter != nil && filter.Limit > 0 {
		opts.SetLimit(filter.Limit)
		if filter.Offset > 0 {
			opts.SetSkip(filter.Offset)
		}
	}

	cursor, err := r.auditCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit records: %w", err)
	}
	defer cursor.Close(ctx)

	var records []*adminModel.AuditRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode audit records: %w", err)
	}

	return records, nil
}

// ============ 操作日志 ============

// CreateAdminLog 创建管理员操作日志
func (r *MongoAdminRepository) CreateAdminLog(ctx context.Context, log *adminModel.AdminLog) error {
	if log == nil {
		return fmt.Errorf("admin log cannot be nil")
	}

	now := time.Now()
	log.CreatedAt = now

	if log.ID == "" {
		log.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.adminLogsCollection.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create admin log: %w", err)
	}

	return nil
}

// GetAdminLog 获取管理员日志
func (r *MongoAdminRepository) GetAdminLog(ctx context.Context, logID string) (*adminModel.AdminLog, error) {
	var log adminModel.AdminLog
	err := r.adminLogsCollection.FindOne(ctx, bson.M{"_id": logID}).Decode(&log)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("admin log not found: %s", logID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get admin log: %w", err)
	}

	return &log, nil
}

// ListAdminLogs 获取管理员日志列表
func (r *MongoAdminRepository) ListAdminLogs(ctx context.Context, filter *shared.AdminLogFilter) ([]*adminModel.AdminLog, error) {
	query := bson.M{}

	if filter != nil {
		if filter.AdminID != "" {
			query["admin_id"] = filter.AdminID
		}
		if filter.Operation != "" {
			query["operation"] = filter.Operation
		}
		if !filter.StartDate.IsZero() {
			query["created_at"] = bson.M{"$gte": filter.StartDate}
		}
		if !filter.EndDate.IsZero() {
			if query["created_at"] != nil {
				query["created_at"].(bson.M)["$lte"] = filter.EndDate
			} else {
				query["created_at"] = bson.M{"$lte": filter.EndDate}
			}
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if filter != nil && filter.Limit > 0 {
		opts.SetLimit(filter.Limit)
		if filter.Offset > 0 {
			opts.SetSkip(filter.Offset)
		}
	}

	cursor, err := r.adminLogsCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list admin logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*adminModel.AdminLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("failed to decode admin logs: %w", err)
	}

	return logs, nil
}

// CountAdminLogs 统计管理员日志数量
func (r *MongoAdminRepository) CountAdminLogs(ctx context.Context, filter *shared.AdminLogFilter) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.AdminID != "" {
			query["admin_id"] = filter.AdminID
		}
		if filter.Operation != "" {
			query["operation"] = filter.Operation
		}
		if !filter.StartDate.IsZero() {
			query["created_at"] = bson.M{"$gte": filter.StartDate}
		}
		if !filter.EndDate.IsZero() {
			if query["created_at"] != nil {
				query["created_at"].(bson.M)["$lte"] = filter.EndDate
			} else {
				query["created_at"] = bson.M{"$lte": filter.EndDate}
			}
		}
	}

	count, err := r.adminLogsCollection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin logs: %w", err)
	}

	return count, nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *MongoAdminRepository) Health(ctx context.Context) error {
	if err := r.db.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	return nil
}
