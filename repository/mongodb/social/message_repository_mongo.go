package social

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/models/social"
	"Qingyu_backend/repository/mongodb/base"
)

// MongoMessageRepository MongoDB消息仓储实现
type MongoMessageRepository struct {
	*base.BaseMongoRepository  // 嵌入基类，管理主collection (conversations)
	db               *mongo.Database  // 保留db字段用于事务和健康检查
	messagesCol      *mongo.Collection  // 消息collection独立管理
	mentionsCol      *mongo.Collection  // @提醒collection独立管理
}

// NewMongoMessageRepository 创建MongoDB消息仓储实例
func NewMongoMessageRepository(db *mongo.Database) *MongoMessageRepository {
	messagesCol := db.Collection("messages")
	mentionsCol := db.Collection("mentions")

	repo := &MongoMessageRepository{
		BaseMongoRepository: base.NewBaseMongoRepository(db, "conversations"),
		db:                  db,
		messagesCol:         messagesCol,
		mentionsCol:         mentionsCol,
	}

	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	repo.createIndexes(ctx)

	return repo
}

// createIndexes 创建索引
func (r *MongoMessageRepository) createIndexes(ctx context.Context) {
	// 会话索引
	conversationIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "participants", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "updated_at", Value: -1}},
		},
	}

	_, err := r.GetCollection().Indexes().CreateMany(ctx, conversationIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create conversation indexes: %v\n", err)
	}

	// 消息索引
	messageIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "conversation_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "conversation_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "sender_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "receiver_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_deleted", Value: 1}},
		},
	}

	_, err = r.messagesCol.Indexes().CreateMany(ctx, messageIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create message indexes: %v\n", err)
	}

	// @提醒索引
	mentionIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "is_read", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "content_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "content_id", Value: 1}},
		},
	}

	_, err = r.mentionsCol.Indexes().CreateMany(ctx, mentionIndexes)
	if err != nil {
		fmt.Printf("Warning: Failed to create mention indexes: %v\n", err)
	}
}

// ========== 会话管理 ==========

// CreateConversation 创建会话
func (r *MongoMessageRepository) CreateConversation(ctx context.Context, conversation *social.Conversation) error {
	if conversation.ID.IsZero() {
		conversation.ID = primitive.NewObjectID()
	}

	if conversation.CreatedAt.IsZero() {
		conversation.CreatedAt = time.Now()
	}
	conversation.UpdatedAt = time.Now()

	if conversation.UnreadCount == nil {
		conversation.UnreadCount = make(map[string]int)
	}

	_, err := r.GetCollection().InsertOne(ctx, conversation)
	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

// GetConversationByID 根据ID获取会话
func (r *MongoMessageRepository) GetConversationByID(ctx context.Context, conversationID string) (*social.Conversation, error) {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	var conversation social.Conversation
	err = r.GetCollection().FindOne(ctx, bson.M{"_id": objectID}).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conversation, nil
}

// GetConversationByParticipants 根据参与者获取会话（私聊）
func (r *MongoMessageRepository) GetConversationByParticipants(ctx context.Context, participantIDs []string) (*social.Conversation, error) {
	// 私聊会话的参与者数量相等
	filter := bson.M{
		"participants": bson.M{
			"$all": participantIDs,
			"$size": len(participantIDs),
		},
	}

	var conversation social.Conversation
	err := r.GetCollection().FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation by participants: %w", err)
	}

	return &conversation, nil
}

// GetUserConversations 获取用户的会话列表
func (r *MongoMessageRepository) GetUserConversations(ctx context.Context, userID string, page, size int) ([]*social.Conversation, int64, error) {
	filter := bson.M{
		"participants": userID,
	}

	// 计算总数
	total, err := r.GetCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "updated_at", Value: -1}})

	cursor, err := r.GetCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find conversations: %w", err)
	}
	defer cursor.Close(ctx)

	var conversations []*social.Conversation
	if err = cursor.All(ctx, &conversations); err != nil {
		return nil, 0, fmt.Errorf("failed to decode conversations: %w", err)
	}

	return conversations, total, nil
}

// UpdateConversation 更新会话
func (r *MongoMessageRepository) UpdateConversation(ctx context.Context, conversationID string, updates map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	updates["updated_at"] = time.Now()

	_, err = r.GetCollection().UpdateByID(ctx, objectID, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	return nil
}

// DeleteConversation 删除会话
func (r *MongoMessageRepository) DeleteConversation(ctx context.Context, conversationID string) error {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	_, err = r.GetCollection().DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	return nil
}

// UpdateLastMessage 更新会话最后一条消息
func (r *MongoMessageRepository) UpdateLastMessage(ctx context.Context, conversationID string, message *social.Message) error {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"last_message": message,
			"updated_at":   time.Now(),
		},
	}

	_, err = r.GetCollection().UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to update last message: %w", err)
	}

	return nil
}

// IncrementUnreadCount 增加未读数
func (r *MongoMessageRepository) IncrementUnreadCount(ctx context.Context, conversationID, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	update := bson.M{
		"$inc": bson.M{
			fmt.Sprintf("unread_count.%s", userID): 1,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = r.GetCollection().UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to increment unread count: %w", err)
	}

	return nil
}

// ClearUnreadCount 清空未读数
func (r *MongoMessageRepository) ClearUnreadCount(ctx context.Context, conversationID, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("unread_count.%s", userID): 0,
			"updated_at":                           time.Now(),
		},
	}

	_, err = r.GetCollection().UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to clear unread count: %w", err)
	}

	return nil
}

// ========== 消息管理 ==========

// CreateMessage 创建消息
func (r *MongoMessageRepository) CreateMessage(ctx context.Context, message *social.Message) error {
	if message.ID.IsZero() {
		message.ID = primitive.NewObjectID()
	}

	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}
	message.UpdatedAt = time.Now()

	_, err := r.messagesCol.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

// GetMessageByID 根据ID获取消息
func (r *MongoMessageRepository) GetMessageByID(ctx context.Context, messageID string) (*social.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, fmt.Errorf("invalid message ID: %w", err)
	}

	var message social.Message
	err = r.messagesCol.FindOne(ctx, bson.M{"_id": objectID}).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

// GetMessagesByConversation 获取会话的消息列表
func (r *MongoMessageRepository) GetMessagesByConversation(ctx context.Context, conversationID string, page, size int) ([]*social.Message, int64, error) {
	filter := bson.M{
		"conversation_id": conversationID,
		"is_deleted":      bson.M{"$ne": true},
	}

	// 计算总数
	total, err := r.messagesCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.messagesCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*social.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, 0, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, total, nil
}

// GetMessagesBetweenUsers 获取两个用户之间的消息
func (r *MongoMessageRepository) GetMessagesBetweenUsers(ctx context.Context, userID1, userID2 string, page, size int) ([]*social.Message, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": userID1, "receiver_id": userID2},
			{"sender_id": userID2, "receiver_id": userID1},
		},
		"is_deleted": bson.M{"$ne": true},
	}

	// 计算总数
	total, err := r.messagesCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.messagesCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*social.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, 0, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, total, nil
}

// MarkMessageAsRead 标记消息为已读
func (r *MongoMessageRepository) MarkMessageAsRead(ctx context.Context, messageID string) error {
	objectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	// 检查消息是否存在
	var existing social.Message
	err = r.messagesCol.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existing)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"readstatus.is_read": true,
			"readstatus.read_at": now,
			"updated_at":         now,
		},
	}

	result, err := r.messagesCol.UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no message matched ID %s", messageID)
	}

	return nil
}

// MarkConversationMessagesAsRead 标记会话中所有消息为已读
func (r *MongoMessageRepository) MarkConversationMessagesAsRead(ctx context.Context, conversationID, userID string) error {
	filter := bson.M{
		"conversation_id": conversationID,
		"receiver_id":     userID,
		"is_read":         false,
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"readstatus.is_read": true,
			"readstatus.read_at": now,
			"updated_at":         now,
		},
	}

	_, err := r.messagesCol.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark conversation messages as read: %w", err)
	}

	return nil
}

// DeleteMessage 删除消息（软删除）
func (r *MongoMessageRepository) DeleteMessage(ctx context.Context, messageID string) error {
	objectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": now,
			"updated_at": now,
		},
	}

	_, err = r.messagesCol.UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// CountUnreadMessages 统计未读消息数
func (r *MongoMessageRepository) CountUnreadMessages(ctx context.Context, userID string) (int, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"participants": userID,
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total": bson.M{
					"$sum": fmt.Sprintf("$unread_count.%s", userID),
				},
			},
		},
	}

	cursor, err := r.GetCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		Total int `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("failed to decode count result: %w", err)
		}
		return result.Total, nil
	}

	return 0, nil
}

// ========== @提醒 ==========

// CreateMention 创建@提醒
func (r *MongoMessageRepository) CreateMention(ctx context.Context, mention *social.Mention) error {
	if mention.ID.IsZero() {
		mention.ID = primitive.NewObjectID()
	}

	if mention.CreatedAt.IsZero() {
		mention.CreatedAt = time.Now()
	}
	mention.UpdatedAt = time.Now()

	// 默认未读
	if !mention.IsRead {
		mention.IsRead = false
	}

	_, err := r.mentionsCol.InsertOne(ctx, mention)
	if err != nil {
		return fmt.Errorf("failed to create mention: %w", err)
	}

	return nil
}

// GetMentionByID 根据ID获取@提醒
func (r *MongoMessageRepository) GetMentionByID(ctx context.Context, mentionID string) (*social.Mention, error) {
	objectID, err := primitive.ObjectIDFromHex(mentionID)
	if err != nil {
		return nil, fmt.Errorf("invalid mention ID: %w", err)
	}

	var mention social.Mention
	err = r.mentionsCol.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mention)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("mention not found")
		}
		return nil, fmt.Errorf("failed to get mention: %w", err)
	}

	return &mention, nil
}

// GetUserMentions 获取用户的@提醒列表
func (r *MongoMessageRepository) GetUserMentions(ctx context.Context, userID string, page, size int) ([]*social.Mention, int64, error) {
	filter := bson.M{
		"user_id": userID,
	}

	// 计算总数
	total, err := r.mentionsCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mentions: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.mentionsCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find mentions: %w", err)
	}
	defer cursor.Close(ctx)

	var mentions []*social.Mention
	if err = cursor.All(ctx, &mentions); err != nil {
		return nil, 0, fmt.Errorf("failed to decode mentions: %w", err)
	}

	return mentions, total, nil
}

// GetUnreadMentions 获取未读的@提醒
func (r *MongoMessageRepository) GetUnreadMentions(ctx context.Context, userID string, page, size int) ([]*social.Mention, int64, error) {
	filter := bson.M{
		"user_id":             userID,
		"readstatus.is_read": false,
	}

	// 计算总数
	total, err := r.mentionsCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count unread mentions: %w", err)
	}

	// 分页查询
	skip := (page - 1) * size
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(size)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.mentionsCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find unread mentions: %w", err)
	}
	defer cursor.Close(ctx)

	var mentions []*social.Mention
	if err = cursor.All(ctx, &mentions); err != nil {
		return nil, 0, fmt.Errorf("failed to decode unread mentions: %w", err)
	}

	return mentions, total, nil
}

// MarkMentionAsRead 标记@提醒为已读
func (r *MongoMessageRepository) MarkMentionAsRead(ctx context.Context, mentionID string) error {
	objectID, err := primitive.ObjectIDFromHex(mentionID)
	if err != nil {
		return fmt.Errorf("invalid mention ID: %w", err)
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"readstatus.is_read": true,
			"readstatus.read_at": now,
			"updated_at":         now,
		},
	}

	_, err = r.mentionsCol.UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to mark mention as read: %w", err)
	}

	return nil
}

// MarkAllMentionsAsRead 标记所有@提醒为已读
func (r *MongoMessageRepository) MarkAllMentionsAsRead(ctx context.Context, userID string) error {
	filter := bson.M{
		"user_id":             userID,
		"readstatus.is_read": false,
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"readstatus.is_read": true,
			"readstatus.read_at": now,
			"updated_at":         now,
		},
	}

	_, err := r.mentionsCol.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark all mentions as read: %w", err)
	}

	return nil
}

// CountUnreadMentions 统计未读@提醒数
func (r *MongoMessageRepository) CountUnreadMentions(ctx context.Context, userID string) (int, error) {
	count, err := r.mentionsCol.CountDocuments(ctx, bson.M{
		"user_id":             userID,
		"readstatus.is_read": false,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count unread mentions: %w", err)
	}

	return int(count), nil
}

// DeleteMention 删除@提醒
func (r *MongoMessageRepository) DeleteMention(ctx context.Context, mentionID string) error {
	objectID, err := primitive.ObjectIDFromHex(mentionID)
	if err != nil {
		return fmt.Errorf("invalid mention ID: %w", err)
	}

	_, err = r.mentionsCol.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete mention: %w", err)
	}

	return nil
}

// Health 健康检查
func (r *MongoMessageRepository) Health(ctx context.Context) error {
	// 检查数据库连接
	return r.db.Client().Ping(ctx, nil)
}
