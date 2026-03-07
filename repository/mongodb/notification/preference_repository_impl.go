package notification

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/models/notification"
	repo "Qingyu_backend/repository/interfaces/notification"
)

// NotificationPreferenceRepositoryImpl 通知偏好设置仓储实现
type NotificationPreferenceRepositoryImpl struct {
	db                   *mongo.Database
	preferenceCollection *mongo.Collection
}

// NewNotificationPreferenceRepository 创建通知偏好设置仓储实例
func NewNotificationPreferenceRepository(db *mongo.Database) repo.NotificationPreferenceRepository {
	return &NotificationPreferenceRepositoryImpl{
		db:                   db,
		preferenceCollection: db.Collection("notification_preferences"),
	}
}

// Create 创建通知偏好设置
func (r *NotificationPreferenceRepositoryImpl) Create(ctx context.Context, preference *notification.NotificationPreference) error {
	if preference.ID.IsZero() {
		preference.ID = primitive.NewObjectID()
	}
	_, err := r.preferenceCollection.InsertOne(ctx, preference)
	if err != nil {
		return fmt.Errorf("创建通知偏好设置失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取通知偏好设置
func (r *NotificationPreferenceRepositoryImpl) GetByID(ctx context.Context, id string) (*notification.NotificationPreference, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的通知偏好设置ID: %w", err)
	}

	var preference notification.NotificationPreference
	err = r.preferenceCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&preference)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询通知偏好设置失败: %w", err)
	}

	return &preference, nil
}

// GetByUserID 根据用户ID获取通知偏好设置
func (r *NotificationPreferenceRepositoryImpl) GetByUserID(ctx context.Context, userID string) (*notification.NotificationPreference, error) {
	var preference notification.NotificationPreference
	err := r.preferenceCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&preference)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询通知偏好设置失败: %w", err)
	}

	return &preference, nil
}

// Update 更新通知偏好设置
func (r *NotificationPreferenceRepositoryImpl) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知偏好设置ID: %w", err)
	}

	result, err := r.preferenceCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新通知偏好设置失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("通知偏好设置不存在")
	}

	return nil
}

// Delete 删除通知偏好设置
func (r *NotificationPreferenceRepositoryImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知偏好设置ID: %w", err)
	}

	result, err := r.preferenceCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除通知偏好设置失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("通知偏好设置不存在")
	}

	return nil
}

// Exists 检查通知偏好设置是否存在
func (r *NotificationPreferenceRepositoryImpl) Exists(ctx context.Context, userID string) (bool, error) {
	count, err := r.preferenceCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, fmt.Errorf("检查通知偏好设置存在性失败: %w", err)
	}

	return count > 0, nil
}

// BatchUpdate 批量更新通知偏好设置
func (r *NotificationPreferenceRepositoryImpl) BatchUpdate(ctx context.Context, ids []string, updates map[string]interface{}) error {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // 跳过无效ID
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return fmt.Errorf("没有有效的通知偏好设置ID")
	}

	_, err := r.preferenceCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("批量更新通知偏好设置失败: %w", err)
	}

	return nil
}
