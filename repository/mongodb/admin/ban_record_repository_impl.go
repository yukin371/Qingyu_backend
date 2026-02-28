package admin

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/admin"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
	"Qingyu_backend/repository/mongodb/base"
)

// banRecordRepositoryImpl 封禁记录仓储实现
type banRecordRepositoryImpl struct {
	*base.BaseMongoRepository
}

// NewBanRecordRepository 创建封禁记录仓储实例
func NewBanRecordRepository(db *mongo.Database) adminrepo.BanRecordRepository {
	return &banRecordRepositoryImpl{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "ban_records"),
	}
}

// Create 创建封禁记录
func (r *banRecordRepositoryImpl) Create(ctx context.Context, record *admin.BanRecord) error {
	record.ID = primitive.NewObjectID()
	record.TouchForCreate()
	return r.Create(ctx, record)
}

// GetByUserID 获取用户的封禁历史
func (r *banRecordRepositoryImpl) GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*admin.BanRecord, int64, error) {
	// 解析用户ID
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, 0, err
	}

	// 构造查询条件
	filter := bson.M{"user_id": userOID}

	// 计算跳过数量
	skip := int64((page - 1) * pageSize)
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(skip).
		SetLimit(int64(pageSize))

	// 查询记录
	var records []*admin.BanRecord
	if err := r.Find(ctx, filter, &records, opts); err != nil {
		return nil, 0, err
	}

	// 统计总数
	total, err := r.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetActiveBan 获取用户当前生效的封禁记录
func (r *banRecordRepositoryImpl) GetActiveBan(ctx context.Context, userID string) (*admin.BanRecord, error) {
	userOID, err := r.ParseID(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"user_id": userOID,
		"action":  "ban",
	}

	opts := options.FindOne().SetSort(bson.M{"created_at": -1})

	var record admin.BanRecord
	if err := r.FindOne(ctx, filter, &record, opts); err != nil {
		return nil, err
	}

	return &record, nil
}
