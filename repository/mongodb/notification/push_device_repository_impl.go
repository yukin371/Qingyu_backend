package notification

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/notification"
	repo "Qingyu_backend/repository/interfaces/notification"
)

// PushDeviceRepositoryImpl 推送设备仓储实现
type PushDeviceRepositoryImpl struct {
	db               *mongo.Database
	deviceCollection *mongo.Collection
}

// NewPushDeviceRepository 创建推送设备仓储实例
func NewPushDeviceRepository(db *mongo.Database) repo.PushDeviceRepository {
	return &PushDeviceRepositoryImpl{
		db:               db,
		deviceCollection: db.Collection("push_devices"),
	}
}

// Create 创建推送设备
func (r *PushDeviceRepositoryImpl) Create(ctx context.Context, device *notification.PushDevice) error {
	objectID, err := primitive.ObjectIDFromHex(device.ID)
	if err != nil {
		// 如果ID无效，生成新的
		objectID = primitive.NewObjectID()
		device.ID = objectID.Hex()
	}

	doc := bson.M{
		"_id":          objectID,
		"user_id":      device.UserID,
		"device_type":  device.DeviceType,
		"device_token": device.DeviceToken,
		"device_id":    device.DeviceID,
		"is_active":    device.IsActive,
		"last_used_at": device.LastUsedAt,
		"created_at":   device.CreatedAt,
	}

	_, err = r.deviceCollection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("创建推送设备失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取推送设备
func (r *PushDeviceRepositoryImpl) GetByID(ctx context.Context, id string) (*notification.PushDevice, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的设备ID: %w", err)
	}

	var device notification.PushDevice
	err = r.deviceCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询推送设备失败: %w", err)
	}

	return &device, nil
}

// Update 更新推送设备
func (r *PushDeviceRepositoryImpl) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的设备ID: %w", err)
	}

	result, err := r.deviceCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新推送设备失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("推送设备不存在")
	}

	return nil
}

// Delete 删除推送设备
func (r *PushDeviceRepositoryImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的设备ID: %w", err)
	}

	result, err := r.deviceCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除推送设备失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("推送设备不存在")
	}

	return nil
}

// Exists 检查推送设备是否存在
func (r *PushDeviceRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("无效的设备ID: %w", err)
	}

	count, err := r.deviceCollection.CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, fmt.Errorf("检查推送设备存在性失败: %w", err)
	}

	return count > 0, nil
}

// GetByUserID 根据用户ID获取推送设备列表
func (r *PushDeviceRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*notification.PushDevice, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.M{"last_used_at": -1})

	cursor, err := r.deviceCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询推送设备列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var devices []*notification.PushDevice
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, fmt.Errorf("解析推送设备列表失败: %w", err)
	}

	return devices, nil
}

// GetByDeviceID 根据设备ID获取推送设备
func (r *PushDeviceRepositoryImpl) GetByDeviceID(ctx context.Context, deviceID string) (*notification.PushDevice, error) {
	var device notification.PushDevice
	err := r.deviceCollection.FindOne(ctx, bson.M{"device_id": deviceID}).Decode(&device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询推送设备失败: %w", err)
	}

	return &device, nil
}

// GetActiveByUserID 获取用户的活跃推送设备列表
func (r *PushDeviceRepositoryImpl) GetActiveByUserID(ctx context.Context, userID string) ([]*notification.PushDevice, error) {
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
	}
	opts := options.Find().SetSort(bson.M{"last_used_at": -1})

	cursor, err := r.deviceCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询活跃推送设备列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var devices []*notification.PushDevice
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, fmt.Errorf("解析活跃推送设备列表失败: %w", err)
	}

	return devices, nil
}

// BatchDelete 批量删除推送设备
func (r *PushDeviceRepositoryImpl) BatchDelete(ctx context.Context, ids []string) error {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // 跳过无效ID
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return fmt.Errorf("没有有效的设备ID")
	}

	_, err := r.deviceCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return fmt.Errorf("批量删除推送设备失败: %w", err)
	}

	return nil
}

// DeactivateAllForUser 停用用户所有推送设备
func (r *PushDeviceRepositoryImpl) DeactivateAllForUser(ctx context.Context, userID string) error {
	_, err := r.deviceCollection.UpdateMany(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"is_active": false}})
	if err != nil {
		return fmt.Errorf("停用所有推送设备失败: %w", err)
	}

	return nil
}

// UpdateLastUsed 更新最后使用时间
func (r *PushDeviceRepositoryImpl) UpdateLastUsed(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的设备ID: %w", err)
	}

	updates := bson.M{"last_used_at": time.Now()}
	result, err := r.deviceCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新最后使用时间失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("推送设备不存在")
	}

	return nil
}

// DeleteInactiveDevices 删除不活跃的设备
func (r *PushDeviceRepositoryImpl) DeleteInactiveDevices(ctx context.Context, beforeDate time.Time) (int64, error) {
	result, err := r.deviceCollection.DeleteMany(ctx, bson.M{
		"last_used_at": bson.M{"$lt": beforeDate},
	})
	if err != nil {
		return 0, fmt.Errorf("删除不活跃设备失败: %w", err)
	}

	return result.DeletedCount, nil
}
