package messaging

import (
	messagingModel "Qingyu_backend/models/messaging"
	messagingInterface "Qingyu_backend/repository/interfaces/messaging"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MessageRepositoryImpl MongoDB消息Repository实现
type MessageRepositoryImpl struct {
	db                      *mongo.Database
	messagesCollection      *mongo.Collection
	notificationsCollection *mongo.Collection
	templatesCollection     *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) messagingInterface.MessageRepository {
	return &MessageRepositoryImpl{
		db:                      db,
		messagesCollection:      db.Collection("messages"),
		notificationsCollection: db.Collection("notifications"),
		templatesCollection:     db.Collection("message_templates"),
	}
}

func toObjectID(id string, resource string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid %s ID: %w", resource, err)
	}
	return objectID, nil
}

// ============ 消息队列管理 ============

func (r *MessageRepositoryImpl) CreateMessage(ctx context.Context, message *messagingModel.Message) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	now := time.Now()
	message.CreatedAt = now
	if message.UpdatedAt.IsZero() {
		message.UpdatedAt = now
	}

	if message.ID.IsZero() {
		message.ID = primitive.NewObjectID()
	}

	_, err := r.messagesCollection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *MessageRepositoryImpl) GetMessage(ctx context.Context, messageID string) (*messagingModel.Message, error) {
	objectID, err := toObjectID(messageID, "message")
	if err != nil {
		return nil, err
	}

	var message messagingModel.Message
	err = r.messagesCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&message)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("message not found: %s", messageID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

func (r *MessageRepositoryImpl) UpdateMessage(ctx context.Context, messageID string, updates map[string]interface{}) error {
	objectID, err := toObjectID(messageID, "message")
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()

	result, err := r.messagesCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("message not found: %s", messageID)
	}
	return nil
}

func (r *MessageRepositoryImpl) DeleteMessage(ctx context.Context, messageID string) error {
	objectID, err := toObjectID(messageID, "message")
	if err != nil {
		return err
	}

	result, err := r.messagesCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("message not found: %s", messageID)
	}
	return nil
}

func (r *MessageRepositoryImpl) ListMessages(ctx context.Context, filter *messagingInterface.MessageFilter) ([]*messagingModel.Message, error) {
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

func (r *MessageRepositoryImpl) CountMessages(ctx context.Context, filter *messagingInterface.MessageFilter) (int64, error) {
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

func (r *MessageRepositoryImpl) CreateNotification(ctx context.Context, notification *messagingModel.NotificationDelivery) error {
	if notification == nil {
		return fmt.Errorf("notification cannot be nil")
	}

	now := time.Now()
	notification.CreatedAt = now
	if notification.ID.IsZero() {
		notification.ID = primitive.NewObjectID()
	}

	_, err := r.notificationsCollection.InsertOne(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

func (r *MessageRepositoryImpl) GetNotification(ctx context.Context, notificationID string) (*messagingModel.NotificationDelivery, error) {
	objectID, err := toObjectID(notificationID, "notification")
	if err != nil {
		return nil, err
	}

	var notification messagingModel.NotificationDelivery
	err = r.notificationsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&notification)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("notification not found: %s", notificationID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return &notification, nil
}

func (r *MessageRepositoryImpl) UpdateNotification(ctx context.Context, notificationID string, updates map[string]interface{}) error {
	objectID, err := toObjectID(notificationID, "notification")
	if err != nil {
		return err
	}

	result, err := r.notificationsCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found: %s", notificationID)
	}
	return nil
}

func (r *MessageRepositoryImpl) DeleteNotification(ctx context.Context, notificationID string) error {
	objectID, err := toObjectID(notificationID, "notification")
	if err != nil {
		return err
	}

	result, err := r.notificationsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("notification not found: %s", notificationID)
	}
	return nil
}

func (r *MessageRepositoryImpl) ListNotifications(ctx context.Context, filter *messagingInterface.NotificationFilter) ([]*messagingModel.NotificationDelivery, int64, error) {
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

	total, err := r.notificationsCollection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

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

	var notifications []*messagingModel.NotificationDelivery
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, 0, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return notifications, total, nil
}

func (r *MessageRepositoryImpl) CountNotifications(ctx context.Context, filter *messagingInterface.NotificationFilter) (int64, error) {
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

func (r *MessageRepositoryImpl) MarkAsRead(ctx context.Context, notificationID string) error {
	objectID, err := toObjectID(notificationID, "notification")
	if err != nil {
		return err
	}

	now := time.Now()
	updates := bson.M{"is_read": true, "read_at": now}
	result, err := r.notificationsCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found: %s", notificationID)
	}
	return nil
}

func (r *MessageRepositoryImpl) MarkMultipleAsRead(ctx context.Context, userID string, notificationIDs []string) error {
	objectIDs := make([]primitive.ObjectID, 0, len(notificationIDs))
	for _, id := range notificationIDs {
		objectID, err := toObjectID(id, "notification")
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	now := time.Now()
	updates := bson.M{"is_read": true, "read_at": now}
	query := bson.M{"user_id": userID, "_id": bson.M{"$in": objectIDs}}

	_, err := r.notificationsCollection.UpdateMany(ctx, query, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}

	return nil
}

func (r *MessageRepositoryImpl) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	query := bson.M{"user_id": userID, "is_read": false}
	count, err := r.notificationsCollection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

// ============ 批量操作 ============

func (r *MessageRepositoryImpl) BatchDeleteNotifications(ctx context.Context, userID string, notificationIDs []string) error {
	objectIDs := make([]primitive.ObjectID, 0, len(notificationIDs))
	for _, id := range notificationIDs {
		objectID, err := toObjectID(id, "notification")
		if err != nil {
			return err
		}
		objectIDs = append(objectIDs, objectID)
	}

	query := bson.M{"user_id": userID, "_id": bson.M{"$in": objectIDs}}
	_, err := r.notificationsCollection.DeleteMany(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to batch delete notifications: %w", err)
	}
	return nil
}

func (r *MessageRepositoryImpl) ClearAllNotifications(ctx context.Context, userID string) error {
	query := bson.M{"user_id": userID}
	_, err := r.notificationsCollection.DeleteMany(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to clear all notifications: %w", err)
	}
	return nil
}

// ============ 消息模板管理 ============

func (r *MessageRepositoryImpl) CreateTemplate(ctx context.Context, template *messagingModel.MessageTemplate) error {
	if template == nil {
		return fmt.Errorf("template cannot be nil")
	}

	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now
	if template.ID.IsZero() {
		template.ID = primitive.NewObjectID()
	}

	_, err := r.templatesCollection.InsertOne(ctx, template)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}
	return nil
}

func (r *MessageRepositoryImpl) GetTemplate(ctx context.Context, templateID string) (*messagingModel.MessageTemplate, error) {
	objectID, err := toObjectID(templateID, "template")
	if err != nil {
		return nil, err
	}

	var template messagingModel.MessageTemplate
	err = r.templatesCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&template)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return &template, nil
}

func (r *MessageRepositoryImpl) GetTemplateByName(ctx context.Context, name string) (*messagingModel.MessageTemplate, error) {
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

func (r *MessageRepositoryImpl) UpdateTemplate(ctx context.Context, templateID string, updates map[string]interface{}) error {
	objectID, err := toObjectID(templateID, "template")
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()
	result, err := r.templatesCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("template not found: %s", templateID)
	}
	return nil
}

func (r *MessageRepositoryImpl) DeleteTemplate(ctx context.Context, templateID string) error {
	objectID, err := toObjectID(templateID, "template")
	if err != nil {
		return err
	}

	result, err := r.templatesCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("template not found: %s", templateID)
	}
	return nil
}

func (r *MessageRepositoryImpl) ListTemplates(ctx context.Context, templateType string, isActive *bool) ([]*messagingModel.MessageTemplate, error) {
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

func (r *MessageRepositoryImpl) Health(ctx context.Context) error {
	if err := r.db.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	return nil
}
