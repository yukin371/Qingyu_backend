package reader

import (
	"Qingyu_backend/models"
	"Qingyu_backend/repository/mongodb/base"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDeviceRepository 设备仓储MongoDB实现
type MongoDeviceRepository struct {
	*base.BaseMongoRepository
	db *mongo.Database
}

// NewMongoDeviceRepository 创建设备仓储实例
func NewMongoDeviceRepository(db *mongo.Database) *MongoDeviceRepository {
	return &MongoDeviceRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "devices"),
		db:                  db,
	}
}

// UpsertDevice 创建或更新设备记录
func (r *MongoDeviceRepository) UpsertDevice(ctx context.Context, device *models.Device) error {
	if device.ID.IsZero() {
		device.ID = primitive.NewObjectID()
	}
	device.LastSeen = time.Now()

	filter := bson.M{
		"user_id":    device.UserID,
		"user_agent": device.UserAgent,
	}

	update := bson.M{
		"$set": bson.M{
			"name":      device.Name,
			"type":      device.Type,
			"ip":        device.IP,
			"last_seen": device.LastSeen,
		},
		"$setOnInsert": bson.M{
			"_id":        device.ID,
			"user_id":    device.UserID,
			"user_agent": device.UserAgent,
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.GetCollection().UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("upsert device failed: %w", err)
	}

	return nil
}

// GetByUserID 获取用户所有设备
func (r *MongoDeviceRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Device, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid userID: %w", err)
	}

	filter := bson.M{"user_id": objectID}
	opts := options.Find().SetSort(bson.M{"last_seen": -1})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find devices failed: %w", err)
	}
	defer cursor.Close(ctx)

	devices := make([]*models.Device, 0)
	if err := cursor.All(ctx, &devices); err != nil {
		return nil, fmt.Errorf("decode devices failed: %w", err)
	}

	return devices, nil
}

// DeleteByUserID 删除用户所有设备
func (r *MongoDeviceRepository) DeleteByUserID(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid userID: %w", err)
	}

	_, err = r.GetCollection().DeleteMany(ctx, bson.M{"user_id": objectID})
	if err != nil {
		return fmt.Errorf("delete devices failed: %w", err)
	}

	return nil
}
