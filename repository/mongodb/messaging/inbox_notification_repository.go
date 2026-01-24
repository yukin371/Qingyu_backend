package messaging

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/messaging"
)

// InboxNotificationRepository 站内通知仓储
type InboxNotificationRepository struct {
	db                        *mongo.Database
	inboxNotificationCollection *mongo.Collection
}

// NewInboxNotificationRepository 创建站内通知仓储实例
func NewInboxNotificationRepository(db *mongo.Database) *InboxNotificationRepository {
	return &InboxNotificationRepository{
		db:                        db,
		inboxNotificationCollection: db.Collection("inbox_notifications"),
	}
}

// InboxNotificationFilter 站内通知过滤器
type InboxNotificationFilter struct {
	ReceiverID string                      // 接收者ID
	Type       messaging.InboxNotificationType // 通知类型
	IsRead     *bool                       // 已读状态
	IsPinned   *bool                       // 是否置顶
	Priority   messaging.InboxNotificationPriority // 优先级
	StartDate  *time.Time                  // 开始日期
	EndDate    *time.Time                  // 结束日期
	IsExpired  bool                        // 是否包含过期通知
	Page       int                         // 页码
	PageSize   int                         // 每页大小
	Limit      int64                       // 限制数量
	Offset     int64                       // 偏移量
}

// Create 创建站内通知
func (r *InboxNotificationRepository) Create(ctx context.Context, notification *messaging.InboxNotification) (string, error) {
	// 设置创建时间
	now := time.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now

	// 如果ID为空，生成新的ObjectID
	if notification.ID.IsZero() {
		notification.ID = primitive.NewObjectID()
	}

	// 插入数据库
	_, err := r.inboxNotificationCollection.InsertOne(ctx, notification)
	if err != nil {
		return "", fmt.Errorf("创建站内通知失败: %w", err)
	}

	return notification.ID.Hex(), nil
}

// GetByID 根据ID获取站内通知
func (r *InboxNotificationRepository) GetByID(ctx context.Context, id string) (*messaging.InboxNotification, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的通知ID: %w", err)
	}

	var notification messaging.InboxNotification
	err = r.inboxNotificationCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 未找到
		}
		return nil, fmt.Errorf("查询站内通知失败: %w", err)
	}

	return &notification, nil
}

// Update 更新站内通知
func (r *InboxNotificationRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	// 自动更新 updated_at
	updates["updated_at"] = time.Now()

	result, err := r.inboxNotificationCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("更新站内通知失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("站内通知不存在")
	}

	return nil
}

// Delete 删除站内通知
func (r *InboxNotificationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	result, err := r.inboxNotificationCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("删除站内通知失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("站内通知不存在")
	}

	return nil
}

// List 获取站内通知列表
func (r *InboxNotificationRepository) List(ctx context.Context, filter *InboxNotificationFilter) ([]*messaging.InboxNotification, int64, error) {
	mongoFilter := r.buildFilter(filter)

	// 统计总数
	total, err := r.inboxNotificationCollection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("统计站内通知失败: %w", err)
	}

	// 构建查询选项
	opts := r.buildFindOptions(filter)

	// 执行查询
	cursor, err := r.inboxNotificationCollection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("查询站内通知列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 解析结果
	var notifications []*messaging.InboxNotification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, 0, fmt.Errorf("解析站内通知失败: %w", err)
	}

	return notifications, total, nil
}

// GetByReceiverID 获取用户的站内通知列表
func (r *InboxNotificationRepository) GetByReceiverID(ctx context.Context, receiverID string, filter *InboxNotificationFilter) ([]*messaging.InboxNotification, int64, error) {
	if filter == nil {
		filter = &InboxNotificationFilter{}
	}
	filter.ReceiverID = receiverID

	return r.List(ctx, filter)
}

// GetUnreadCount 获取用户未读通知数量
func (r *InboxNotificationRepository) GetUnreadCount(ctx context.Context, receiverID string) (int64, error) {
	filter := bson.M{
		"receiver_id": receiverID,
		"is_read":     false,
		"is_deleted":  bson.M{"$ne": true}, // 排除已删除的
	}

	count, err := r.inboxNotificationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("统计未读通知失败: %w", err)
	}

	return count, nil
}

// MarkAsRead 标记通知为已读
func (r *InboxNotificationRepository) MarkAsRead(ctx context.Context, id, receiverID string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	now := time.Now()
	updates := bson.M{
		"is_read":     true,
		"read_at":     &now,
		"updated_at":  now,
	}

	// 确保只能标记自己的通知
	filter := bson.M{
		"_id":         objectID,
		"receiver_id": receiverID,
	}

	result, err := r.inboxNotificationCollection.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("标记已读失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("站内通知不存在或无权操作")
	}

	return nil
}

// MarkAllAsRead 标记用户所有通知为已读
func (r *InboxNotificationRepository) MarkAllAsRead(ctx context.Context, receiverID string) (int64, error) {
	now := time.Now()
	updates := bson.M{
		"is_read":     true,
		"read_at":     &now,
		"updated_at":  now,
	}

	filter := bson.M{
		"receiver_id": receiverID,
		"is_read":     false,
		"is_deleted":  bson.M{"$ne": true},
	}

	result, err := r.inboxNotificationCollection.UpdateMany(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return 0, fmt.Errorf("标记全部已读失败: %w", err)
	}

	return result.ModifiedCount, nil
}

// SoftDelete 软删除站内通知
func (r *InboxNotificationRepository) SoftDelete(ctx context.Context, id, receiverID string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	now := time.Now()
	updates := bson.M{
		"is_deleted": true,
		"deleted_at": &now,
		"updated_at": now,
	}

	// 确保只能删除自己的通知
	filter := bson.M{
		"_id":         objectID,
		"receiver_id": receiverID,
	}

	result, err := r.inboxNotificationCollection.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("软删除失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("站内通知不存在或无权操作")
	}

	return nil
}

// DeleteExpired 删除过期的站内通知
func (r *InboxNotificationRepository) DeleteExpired(ctx context.Context) (int64, error) {
	now := time.Now()
	filter := bson.M{
		"expires_at": bson.M{"$lt": now},
		"is_deleted": bson.M{"$ne": true},
	}

	// 软删除过期的通知
	updates := bson.M{
		"is_deleted": true,
		"deleted_at": &now,
		"updated_at": now,
	}

	result, err := r.inboxNotificationCollection.UpdateMany(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return 0, fmt.Errorf("删除过期通知失败: %w", err)
	}

	return result.ModifiedCount, nil
}

// PinNotification 置顶通知
func (r *InboxNotificationRepository) PinNotification(ctx context.Context, id, receiverID string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	updates := bson.M{
		"is_pinned":   true,
		"updated_at": time.Now(),
	}

	filter := bson.M{
		"_id":         objectID,
		"receiver_id": receiverID,
	}

	result, err := r.inboxNotificationCollection.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("置顶通知失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("站内通知不存在或无权操作")
	}

	return nil
}

// UnpinNotification 取消置顶通知
func (r *InboxNotificationRepository) UnpinNotification(ctx context.Context, id, receiverID string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("无效的通知ID: %w", err)
	}

	updates := bson.M{
		"is_pinned":   false,
		"updated_at": time.Now(),
	}

	filter := bson.M{
		"_id":         objectID,
		"receiver_id": receiverID,
	}

	result, err := r.inboxNotificationCollection.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("取消置顶失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("站内通知不存在或无权操作")
	}

	return nil
}

// Health 健康检查
func (r *InboxNotificationRepository) Health(ctx context.Context) error {
	return r.db.Client().Ping(ctx, nil)
}

// buildFilter 构建查询过滤器
func (r *InboxNotificationRepository) buildFilter(filter *InboxNotificationFilter) bson.M {
	mongoFilter := bson.M{}

	if filter == nil {
		// 默认只查询未删除的通知
		mongoFilter["is_deleted"] = bson.M{"$ne": true}
		return mongoFilter
	}

	// 接收者ID
	if filter.ReceiverID != "" {
		mongoFilter["receiver_id"] = filter.ReceiverID
	}

	// 通知类型
	if filter.Type != "" {
		mongoFilter["type"] = filter.Type
	}

	// 已读状态
	if filter.IsRead != nil {
		mongoFilter["is_read"] = *filter.IsRead
	}

	// 置顶状态
	if filter.IsPinned != nil {
		mongoFilter["is_pinned"] = *filter.IsPinned
	}

	// 优先级
	if filter.Priority != "" {
		mongoFilter["priority"] = filter.Priority
	}

	// 时间范围
	if filter.StartDate != nil || filter.EndDate != nil {
		timeFilter := bson.M{}
		if filter.StartDate != nil {
			timeFilter["$gte"] = filter.StartDate
		}
		if filter.EndDate != nil {
			timeFilter["$lte"] = filter.EndDate
		}
		mongoFilter["created_at"] = timeFilter
	}

	// 是否包含过期通知
	if !filter.IsExpired {
		mongoFilter["$or"] = []bson.M{
			{"expires_at": bson.M{"$exists": false}},
			{"expires_at": bson.M{"$gt": time.Now()}},
		}
	}

	// 默认只查询未删除的通知
	mongoFilter["is_deleted"] = bson.M{"$ne": true}

	return mongoFilter
}

// buildFindOptions 构建查询选项
func (r *InboxNotificationRepository) buildFindOptions(filter *InboxNotificationFilter) *options.FindOptions {
	opts := options.Find()

	// 默认排序：置顶优先，然后按创建时间倒序
	sort := bson.D{
		{Key: "is_pinned", Value: -1},
		{Key: "priority", Value: -1},
		{Key: "created_at", Value: -1},
	}
	opts.SetSort(sort)

	// 分页
	if filter != nil {
		if filter.Page > 0 && filter.PageSize > 0 {
			opts.SetSkip(int64((filter.Page - 1) * filter.PageSize))
			opts.SetLimit(int64(filter.PageSize))
		} else if filter.Limit > 0 {
			opts.SetLimit(filter.Limit)
			if filter.Offset > 0 {
				opts.SetSkip(filter.Offset)
			}
		}
	}

	return opts
}
