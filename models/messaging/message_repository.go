package messaging

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MessageRepository 消息Repository接口
type MessageRepository interface {
	Create(ctx context.Context, message *DirectMessage) error
	FindByConversationID(ctx context.Context, conversationID, userID string, page, pageSize int, before, after *string) ([]*DirectMessage, int, error)
	MarkConversationRead(ctx context.Context, conversationID, userID string, readAt time.Time) (int, error)
	CountUnreadInConversation(ctx context.Context, conversationID, userID string) (int, error)
}

// MongoMessageRepository MongoDB实现的消息Repository
type MongoMessageRepository struct {
	collection *mongo.Collection
}

// NewMongoMessageRepository 创建MongoDB消息Repository
func NewMongoMessageRepository(db *mongo.Database) *MongoMessageRepository {
	return &MongoMessageRepository{
		collection: db.Collection("direct_messages"),
	}
}

// Create 创建消息
func (r *MongoMessageRepository) Create(ctx context.Context, msg *DirectMessage) error {
	now := time.Now()
	msg.CreatedAt = now
	msg.UpdatedAt = now

	// 生成新的ObjectID
	objID := primitive.NewObjectID()
	msg.ID = objID

	// 转换为BSON
	doc, err := bson.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = r.collection.InsertOne(ctx, doc)
	return err
}

// FindByConversationID 根据会话ID查找消息
func (r *MongoMessageRepository) FindByConversationID(
	ctx context.Context,
	conversationID string,
	userID string,
	page, pageSize int,
	before, after *string,
) ([]*DirectMessage, int, error) {
	// 基础查询条件
	filter := bson.M{
		"conversation_id": conversationID,
		"status": bson.M{"$ne": MessageStatusDeleted},
	}

	// 如果有before参数，获取该消息之前的消息
	if before != nil {
		beforeObjID, err := primitive.ObjectIDFromHex(*before)
		if err == nil {
			filter["_id"] = bson.M{"$lt": beforeObjID}
		}
	}

	// 如果有after参数，获取该消息之后的消息
	if after != nil {
		afterObjID, err := primitive.ObjectIDFromHex(*after)
		if err == nil {
			filter["_id"] = bson.M{"$gt": afterObjID}
		}
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 计算跳过的数量
	skip := (page - 1) * pageSize

	// 查询选项
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 执行查询
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var messages []*DirectMessage
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	return messages, int(total), nil
}

// MarkConversationRead 标记会话中的所有未读消息为已读
func (r *MongoMessageRepository) MarkConversationRead(
	ctx context.Context,
	conversationID string,
	userID string,
	readAt time.Time,
) (int, error) {
	// 查询条件：会话中未读的、发给该用户的消息
	filter := bson.M{
		"conversation_id": conversationID,
		"receiver_id":     userID,
		"is_read":         false,
		"status":          bson.M{"$ne": MessageStatusDeleted},
	}

	// 更新操作
	update := bson.M{
		"$set": bson.M{
			"is_read":   true,
			"read_at":   readAt,
			"updated_at": time.Now(),
		},
	}

	// 执行更新
	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return int(result.ModifiedCount), nil
}

// CountUnreadInConversation 统计会话中用户的未读消息数
func (r *MongoMessageRepository) CountUnreadInConversation(
	ctx context.Context,
	conversationID string,
	userID string,
) (int, error) {
	filter := bson.M{
		"conversation_id": conversationID,
		"receiver_id":     userID,
		"is_read":         false,
		"status":          bson.M{"$ne": MessageStatusDeleted},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
