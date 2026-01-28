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

// NotificationRepositoryImpl 通知仓储实现
type NotificationRepositoryImpl struct {
	db                     *mongo.Database
	notificationCollection *mongo.Collection
}

// NewNotificationRepository 创建通知仓储实例
func NewNotificationRepository(db *mongo.Database) repo.NotificationRepository {
	return &NotificationRepositoryImpl{
		db:                     db,
		notificationCollection: db.Collection("notifications"),
	}
}

// Create 创建通知
func (r *NotificationRepositoryImpl) Create(ctx context.Context, notif *notification.Notification) error {
	objectID, err := primitive.ObjectIDFromHex(notif.ID)
	if err != nil {
		// 如果ID无效，生成新的
		objectID = primitive.NewObjectID()
		notif.ID = objectID.Hex()
	}

	doc := bson.M{
		"_id":        objectID,
		"user_id":    notif.UserID,
		"type":       notif.Type,
		"priority":   notif.Priority,
		"title":      notif.Title,
		"content":    notif.Content,
		"data":       notif.Data,
		"read":       notif.Read,
		"read_at":    notif.ReadAt,
		"created_at": notif.CreatedAt,
		"expires_at": notif.ExpiresAt,
	}

	_, err = r.notificationCollection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("创建通知失败: %w", err)
	}

	return nil
}

// GetByID 根据ID获取通知
func (r *NotificationRepositoryImpl) GetByID(ctx context.Context, id string) (*notification.Notification, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的通知ID: %w", err)
	}

	var notif notification.Notification
	err = r.notificationCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&notif)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("查询通知失败: %w", err)
	}

	return &notif, nil
}

// Update 更新通知
func (r *NotificationRepositoryImpl) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	updates["updated_at"] = time.Now()

	result, err := r.notificationCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新通知失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("通知不存在")
	}

	return nil
}

// Delete 删除通知
func (r *NotificationRepositoryImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	result, err := r.notificationCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除通知失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("通知不存在")
	}

	return nil
}

// Exists 检查通知是否存在
func (r *NotificationRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("无效的通知ID: %w", err)
	}

	count, err := r.notificationCollection.CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, fmt.Errorf("检查通知存在性失败: %w", err)
	}

	return count > 0, nil
}

// List 获取通知列表
func (r *NotificationRepositoryImpl) List(ctx context.Context, filter *notification.NotificationFilter) ([]*notification.Notification, error) {
	mongoFilter := r.buildFilter(filter)

	// 构建排序
	sortOrder := 1 // 默认升序
	if filter.SortOrder == "desc" {
		sortOrder = -1
	}

	sortField := "created_at"
	if filter.SortBy != "" {
		sortField = filter.SortBy
	}

	opts := options.Find()
	opts.SetSort(bson.M{sortField: sortOrder})
	opts.SetSkip(int64(filter.Offset))
	opts.SetLimit(int64(filter.Limit))

	cursor, err := r.notificationCollection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询通知列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []*notification.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("解析通知列表失败: %w", err)
	}

	return notifications, nil
}

// Count 统计通知数量
func (r *NotificationRepositoryImpl) Count(ctx context.Context, filter *notification.NotificationFilter) (int64, error) {
	mongoFilter := r.buildFilter(filter)

	count, err := r.notificationCollection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, fmt.Errorf("统计通知数量失败: %w", err)
	}

	return count, nil
}

// BatchMarkAsRead 批量标记为已读
func (r *NotificationRepositoryImpl) BatchMarkAsRead(ctx context.Context, ids []string) error {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // 跳过无效ID
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return fmt.Errorf("没有有效的通知ID")
	}

	now := time.Now()
	updates := bson.M{
		"read":    true,
		"read_at": now,
	}

	_, err := r.notificationCollection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("批量标记已读失败: %w", err)
	}

	return nil
}

// MarkAllAsReadForUser 标记用户所有通知为已读
func (r *NotificationRepositoryImpl) MarkAllAsReadForUser(ctx context.Context, userID string) error {
	now := time.Now()
	updates := bson.M{
		"read":    true,
		"read_at": now,
	}

	_, err := r.notificationCollection.UpdateMany(ctx, bson.M{"user_id": userID, "read": false}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("标记全部已读失败: %w", err)
	}

	return nil
}

// BatchDelete 批量删除通知
func (r *NotificationRepositoryImpl) BatchDelete(ctx context.Context, ids []string) error {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // 跳过无效ID
		}
		objectIDs = append(objectIDs, objectID)
	}

	if len(objectIDs) == 0 {
		return fmt.Errorf("没有有效的通知ID")
	}

	_, err := r.notificationCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return fmt.Errorf("批量删除失败: %w", err)
	}

	return nil
}

// DeleteAllForUser 删除用户所有通知
func (r *NotificationRepositoryImpl) DeleteAllForUser(ctx context.Context, userID string) error {
	_, err := r.notificationCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return fmt.Errorf("删除所有通知失败: %w", err)
	}

	return nil
}

// CountUnread 统计未读通知数量
func (r *NotificationRepositoryImpl) CountUnread(ctx context.Context, userID string) (int64, error) {
	count, err := r.notificationCollection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"read":    false,
	})
	if err != nil {
		return 0, fmt.Errorf("统计未读通知失败: %w", err)
	}

	return count, nil
}

// GetStats 获取通知统计
func (r *NotificationRepositoryImpl) GetStats(ctx context.Context, userID string) (*notification.NotificationStats, error) {
	// 统计总数
	totalCount, err := r.notificationCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("统计总数失败: %w", err)
	}

	// 统计未读数
	unreadCount, err := r.CountUnread(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 按类型统计
	typePipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":   "$type",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := r.notificationCollection.Aggregate(ctx, typePipeline)
	if err != nil {
		return nil, fmt.Errorf("按类型统计失败: %w", err)
	}
	defer cursor.Close(ctx)

	typeCounts := make(map[notification.NotificationType]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    notification.NotificationType `bson:"_id"`
			Count int64                         `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		typeCounts[result.ID] = result.Count
	}

	// 按优先级统计
	priorityPipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":   "$priority",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor2, err := r.notificationCollection.Aggregate(ctx, priorityPipeline)
	if err != nil {
		return nil, fmt.Errorf("按优先级统计失败: %w", err)
	}
	defer cursor2.Close(ctx)

	priorityCounts := make(map[notification.NotificationPriority]int64)
	for cursor2.Next(ctx) {
		var result struct {
			ID    notification.NotificationPriority `bson:"_id"`
			Count int64                             `bson:"count"`
		}
		if err := cursor2.Decode(&result); err != nil {
			continue
		}
		priorityCounts[result.ID] = result.Count
	}

	return &notification.NotificationStats{
		TotalCount:     totalCount,
		UnreadCount:    unreadCount,
		TypeCounts:     typeCounts,
		PriorityCounts: priorityCounts,
	}, nil
}

// GetUnreadByType 获取指定类型的未读通知
func (r *NotificationRepositoryImpl) GetUnreadByType(ctx context.Context, userID string, notificationType notification.NotificationType) ([]*notification.Notification, error) {
	filter := bson.M{
		"user_id": userID,
		"type":    notificationType,
		"read":    false,
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(50)

	cursor, err := r.notificationCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询未读通知失败: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []*notification.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("解析未读通知失败: %w", err)
	}

	return notifications, nil
}

// DeleteExpired 删除过期通知
func (r *NotificationRepositoryImpl) DeleteExpired(ctx context.Context) (int64, error) {
	now := time.Now()
	result, err := r.notificationCollection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lt": now},
	})
	if err != nil {
		return 0, fmt.Errorf("删除过期通知失败: %w", err)
	}

	return result.DeletedCount, nil
}

// DeleteOldNotifications 删除旧通知
func (r *NotificationRepositoryImpl) DeleteOldNotifications(ctx context.Context, beforeDate time.Time) (int64, error) {
	result, err := r.notificationCollection.DeleteMany(ctx, bson.M{
		"created_at": bson.M{"$lt": beforeDate},
	})
	if err != nil {
		return 0, fmt.Errorf("删除旧通知失败: %w", err)
	}

	return result.DeletedCount, nil
}

// buildFilter 构建MongoDB查询条件
func (r *NotificationRepositoryImpl) buildFilter(filter *notification.NotificationFilter) bson.M {
	mongoFilter := bson.M{}

	if filter.UserID != nil {
		mongoFilter["user_id"] = *filter.UserID
	}

	if filter.Type != nil {
		mongoFilter["type"] = *filter.Type
	}

	if filter.Read != nil {
		mongoFilter["read"] = *filter.Read
	}

	if filter.Priority != nil {
		mongoFilter["priority"] = *filter.Priority
	}

	if filter.StartDate != nil {
		mongoFilter["created_at"] = bson.M{"$gte": *filter.StartDate}
	}

	if filter.EndDate != nil {
		if _, exists := mongoFilter["created_at"]; exists {
			mongoFilter["created_at"].(bson.M)["$lte"] = *filter.EndDate
		} else {
			mongoFilter["created_at"] = bson.M{"$lte": *filter.EndDate}
		}
	}

	if filter.Keyword != nil {
		mongoFilter["$or"] = []bson.M{
			{"title": bson.M{"$regex": *filter.Keyword, "$options": "i"}},
			{"content": bson.M{"$regex": *filter.Keyword, "$options": "i"}},
		}
	}

	return mongoFilter
}

// DeleteReadForUser 删除用户所有已读通知
func (r *NotificationRepositoryImpl) DeleteReadForUser(ctx context.Context, userID string) (int64, error) {
	filter := bson.M{
		"user_id": userID,
		"read":    true,
	}

	result, err := r.notificationCollection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}


