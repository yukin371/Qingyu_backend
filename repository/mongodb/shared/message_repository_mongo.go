package shared

import (
	messagingModel "Qingyu_backend/models/messaging"
	"context"
	"fmt"
	"time"

	"Qingyu_backend/repository/interfaces/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoMessageRepository MongoDB消息Repository实现
type MongoMessageRepository struct {
	db                      *mongo.Database
	messagesCollection      *mongo.Collection
	notificationsCollection *mongo.Collection
	templatesCollection     *mongo.Collection
}

// NewMongoMessageRepository 创建MongoDB Message Repository
func NewMongoMessageRepository(db *mongo.Database) shared.MessageRepository {
	return &MongoMessageRepository{
		db:                      db,
		messagesCollection:      db.Collection("messages"),
		notificationsCollection: db.Collection("notifications"),
		templatesCollection:     db.Collection("message_templates"),
	}
}

// ============ 消息队列管理 ============

// CreateMessage 创建消息
func (r *MongoMessageRepository) CreateMessage(ctx context.Context, message *messagingModel.Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	now := time.Now()
	message.CreatedAt = now

	if message.ID == "" {
		message.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.messagesCollection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

// GetMessage 获取消息
func (r *MongoMessageRepository) GetMessage(ctx context.Context, messageID string) (*messagingModel.Message, error) {
	var message messagingModel.Message
	err := r.messagesCollection.FindOne(ctx, bson.M{"_id": messageID}).Decode(&message)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("message not found: %s", messageID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

// UpdateMessage 更新消息
func (r *MongoMessageRepository) UpdateMessage(ctx context.Context, messageID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.messagesCollection.UpdateOne(
		ctx,
		bson.M{"_id": messageID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("message not found: %s", messageID)
	}

	return nil
}

// DeleteMessage 删除消息
func (r *MongoMessageRepository) DeleteMessage(ctx context.Context, messageID string) error {
	result, err := r.messagesCollection.DeleteOne(ctx, bson.M{"_id": messageID})

	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("message not found: %s", messageID)
	}

	return nil
}

// ListMessages 获取消息列表
func (r *MongoMessageRepository) ListMessages(ctx context.Context, filter *shared.MessageFilter) ([]*messagingModel.Message, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Topic != "" {
			query["topic"] = filter.Topic
		}
		if filter.Status != "" {
			query["status"] = filter.Status
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

	cursor, err := r.messagesCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*messagingModel.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, nil
}

// CountMessages 统计消息数量
func (r *MongoMessageRepository) CountMessages(ctx context.Context, filter *shared.MessageFilter) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.Topic != "" {
			query["topic"] = filter.Topic
		}
		if filter.Status != "" {
			query["status"] = filter.Status
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

	count, err := r.messagesCollection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

// ============ 通知记录管理 ============

// CreateNotification 创建通知
func (r *MongoMessageRepository) CreateNotification(ctx context.Context, notification *messagingModel.Notification) error {
	if notification == nil {
		return fmt.Errorf("notification cannot be nil")
	}

	now := time.Now()
	notification.CreatedAt = now

	if notification.ID == "" {
		notification.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.notificationsCollection.InsertOne(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

// GetNotification 获取通知
func (r *MongoMessageRepository) GetNotification(ctx context.Context, notificationID string) (*messagingModel.Notification, error) {
	var notification messagingModel.Notification
	err := r.notificationsCollection.FindOne(ctx, bson.M{"_id": notificationID}).Decode(&notification)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("notification not found: %s", notificationID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return &notification, nil
}

// UpdateNotification 更新通知
func (r *MongoMessageRepository) UpdateNotification(ctx context.Context, notificationID string, updates map[string]interface{}) error {
	result, err := r.notificationsCollection.UpdateOne(
		ctx,
		bson.M{"_id": notificationID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found: %s", notificationID)
	}

	return nil
}

// DeleteNotification 删除通知
func (r *MongoMessageRepository) DeleteNotification(ctx context.Context, notificationID string) error {
	result, err := r.notificationsCollection.DeleteOne(ctx, bson.M{"_id": notificationID})

	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("notification not found: %s", notificationID)
	}

	return nil
}

// ListNotifications 获取通知列表
func (r *MongoMessageRepository) ListNotifications(ctx context.Context, filter *shared.NotificationFilter) ([]*messagingModel.Notification, int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.UserID != "" {
			query["user_id"] = filter.UserID
		}
		if filter.Type != "" {
			query["type"] = filter.Type
		}
		if filter.Status != "" {
			query["status"] = filter.Status
		}
		if filter.IsRead != nil {
			query["is_read"] = *filter.IsRead
		}
		if filter.StartDate != nil {
			query["created_at"] = bson.M{"$gte": *filter.StartDate}
		}
		if filter.EndDate != nil {
			if query["created_at"] != nil {
				query["created_at"].(bson.M)["$lte"] = *filter.EndDate
			} else {
				query["created_at"] = bson.M{"$lte": *filter.EndDate}
			}
		}
	}

	// 统计总数
	total, err := r.notificationsCollection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	// 查询列表
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	if filter != nil {
		if filter.PageSize > 0 {
			opts.SetLimit(int64(filter.PageSize))
			if filter.Page > 0 {
				opts.SetSkip(int64((filter.Page - 1) * filter.PageSize))
			}
		} else if filter.Limit > 0 {
			opts.SetLimit(filter.Limit)
			if filter.Offset > 0 {
				opts.SetSkip(filter.Offset)
			}
		}
	}

	cursor, err := r.notificationsCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []*messagingModel.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, 0, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return notifications, total, nil
}

// CountNotifications 统计通知数量
func (r *MongoMessageRepository) CountNotifications(ctx context.Context, filter *shared.NotificationFilter) (int64, error) {
	query := bson.M{}

	if filter != nil {
		if filter.UserID != "" {
			query["user_id"] = filter.UserID
		}
		if filter.Type != "" {
			query["type"] = filter.Type
		}
		if filter.Status != "" {
			query["status"] = filter.Status
		}
		if filter.IsRead != nil {
			query["is_read"] = *filter.IsRead
		}
	}

	count, err := r.notificationsCollection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	return count, nil
}

// ============ 通知已读状态管理 ============

// MarkAsRead 标记为已读
func (r *MongoMessageRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	now := time.Now()
	updates := bson.M{
		"is_read": true,
		"read_at": now,
	}

	result, err := r.notificationsCollection.UpdateOne(
		ctx,
		bson.M{"_id": notificationID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found: %s", notificationID)
	}

	return nil
}

// MarkMultipleAsRead 批量标记为已读
func (r *MongoMessageRepository) MarkMultipleAsRead(ctx context.Context, userID string, notificationIDs []string) error {
	now := time.Now()
	updates := bson.M{
		"is_read": true,
		"read_at": now,
	}

	query := bson.M{
		"user_id": userID,
		"_id":     bson.M{"$in": notificationIDs},
	}

	_, err := r.notificationsCollection.UpdateMany(
		ctx,
		query,
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}

	return nil
}

// GetUnreadCount 获取未读数量
func (r *MongoMessageRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	query := bson.M{
		"user_id": userID,
		"is_read": false,
	}

	count, err := r.notificationsCollection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// ============ 批量操作 ============

// BatchDeleteNotifications 批量删除通知
func (r *MongoMessageRepository) BatchDeleteNotifications(ctx context.Context, userID string, notificationIDs []string) error {
	query := bson.M{
		"user_id": userID,
		"_id":     bson.M{"$in": notificationIDs},
	}

	_, err := r.notificationsCollection.DeleteMany(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to batch delete notifications: %w", err)
	}

	return nil
}

// ClearAllNotifications 清空所有通知
func (r *MongoMessageRepository) ClearAllNotifications(ctx context.Context, userID string) error {
	query := bson.M{"user_id": userID}

	_, err := r.notificationsCollection.DeleteMany(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to clear all notifications: %w", err)
	}

	return nil
}

// ============ 消息模板管理 ============

// CreateTemplate 创建消息模板
func (r *MongoMessageRepository) CreateTemplate(ctx context.Context, template *messagingModel.MessageTemplate) error {
	if template == nil {
		return fmt.Errorf("template cannot be nil")
	}

	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now

	if template.ID == "" {
		template.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.templatesCollection.InsertOne(ctx, template)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

// GetTemplate 获取模板
func (r *MongoMessageRepository) GetTemplate(ctx context.Context, templateID string) (*messagingModel.MessageTemplate, error) {
	var template messagingModel.MessageTemplate
	err := r.templatesCollection.FindOne(ctx, bson.M{"_id": templateID}).Decode(&template)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &template, nil
}

// GetTemplateByName 根据名称获取模板
func (r *MongoMessageRepository) GetTemplateByName(ctx context.Context, name string) (*messagingModel.MessageTemplate, error) {
	var template messagingModel.MessageTemplate
	err := r.templatesCollection.FindOne(ctx, bson.M{"name": name}).Decode(&template)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("template not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get template by name: %w", err)
	}

	return &template, nil
}

// UpdateTemplate 更新模板
func (r *MongoMessageRepository) UpdateTemplate(ctx context.Context, templateID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	result, err := r.templatesCollection.UpdateOne(
		ctx,
		bson.M{"_id": templateID},
		bson.M{"$set": updates},
	)

	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("template not found: %s", templateID)
	}

	return nil
}

// DeleteTemplate 删除模板
func (r *MongoMessageRepository) DeleteTemplate(ctx context.Context, templateID string) error {
	result, err := r.templatesCollection.DeleteOne(ctx, bson.M{"_id": templateID})

	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("template not found: %s", templateID)
	}

	return nil
}

// ListTemplates 获取模板列表
func (r *MongoMessageRepository) ListTemplates(ctx context.Context, templateType string, isActive *bool) ([]*messagingModel.MessageTemplate, error) {
	query := bson.M{}

	if templateType != "" {
		query["type"] = templateType
	}
	if isActive != nil {
		query["is_active"] = *isActive
	}

	cursor, err := r.templatesCollection.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	defer cursor.Close(ctx)

	var templates []*messagingModel.MessageTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, fmt.Errorf("failed to decode templates: %w", err)
	}

	return templates, nil
}

// ============ 健康检查 ============

// Health 健康检查
func (r *MongoMessageRepository) Health(ctx context.Context) error {
	if err := r.db.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	return nil
}
