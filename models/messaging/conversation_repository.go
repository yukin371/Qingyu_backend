package messaging

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConversationRepository 会话Repository接口
type ConversationRepository interface {
	Create(ctx context.Context, conversation *Conversation) error
	FindByID(ctx context.Context, id string) (*Conversation, error)
	FindByParticipants(ctx context.Context, participantIDs []string) (*Conversation, error)
	FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]*Conversation, int64, error)
	Update(ctx context.Context, conversation *Conversation) error
	IncrementUnreadCount(ctx context.Context, conversationID string, userID string) error
}

// MongoConversationRepository MongoDB实现的会话Repository
type MongoConversationRepository struct {
	collection *mongo.Collection
}

// NewMongoConversationRepository 创建MongoDB会话Repository
func NewMongoConversationRepository(db *mongo.Database) *MongoConversationRepository {
	return &MongoConversationRepository{
		collection: db.Collection("conversations"),
	}
}

// Create 创建会话
func (r *MongoConversationRepository) Create(ctx context.Context, conv *Conversation) error {
	now := time.Now()
	conv.CreatedAt = now
	conv.UpdatedAt = now

	// 生成新的ObjectID
	objID := primitive.NewObjectID()
	conv.ID = objID

	// 转换为BSON
	doc, err := bson.Marshal(conv)
	if err != nil {
		return err
	}

	_, err = r.collection.InsertOne(ctx, doc)
	return err
}

// FindByID 根据ID查找会话
func (r *MongoConversationRepository) FindByID(ctx context.Context, id string) (*Conversation, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var conv Conversation
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&conv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}

	return &conv, nil
}

// FindByParticipants 根据参与者列表查找会话
func (r *MongoConversationRepository) FindByParticipants(ctx context.Context, participantIDs []string) (*Conversation, error) {
	// 查找包含所有指定参与者的会话
	// 对于私聊，参与者数量应该相同
	filter := bson.M{
		"participant_ids": bson.M{"$all": participantIDs},
		"type":            ConversationTypeDirect,
		"$expr": bson.M{"$eq": bson.A{"$size", "$participant_ids", len(participantIDs)}},
	}

	var conv Conversation
	err := r.collection.FindOne(ctx, filter).Decode(&conv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}

	return &conv, nil
}

// FindByUserID 查找用户的所有会话
func (r *MongoConversationRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]*Conversation, int64, error) {
	// 计算跳过的数量
	skip := (page - 1) * pageSize

	// 查询条件
	filter := bson.M{
		"participant_ids": userID,
		"is_active":       true,
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 查询选项
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "last_message_at", Value: -1}})

	// 执行查询
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var conversations []*Conversation
	if err = cursor.All(ctx, &conversations); err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// Update 更新会话
func (r *MongoConversationRepository) Update(ctx context.Context, conv *Conversation) error {
	conv.UpdatedAt = time.Now()

	update := bson.M{"$set": conv}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": conv.ID}, update)
	return err
}

// IncrementUnreadCount 增加未读计数
func (r *MongoConversationRepository) IncrementUnreadCount(ctx context.Context, conversationID string, userID string) error {
	objID, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return err
	}

	// 为特定用户增加未读计数
	// 注意：这里假设未读计数是按用户存储的
	update := bson.M{
		"$inc": bson.M{
			"unread_count." + userID: 1,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// ErrConversationNotFound 会话不存在错误
var ErrConversationNotFound = errors.New("会话不存在")
