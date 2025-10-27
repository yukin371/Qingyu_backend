package ai

import (
	"context"
	"fmt"
	"time"

	aiModels "Qingyu_backend/models/ai"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatRepository 聊天数据访问层
type ChatRepository struct {
	db *mongo.Database
}

// NewChatRepository 创建聊天数据访问层
func NewChatRepository(db *mongo.Database) *ChatRepository {
	return &ChatRepository{db: db}
}

// getSessionCollection 获取会话集合
func (r *ChatRepository) getSessionCollection() *mongo.Collection {
	return r.db.Collection("ai_chat_sessions")
}

// getMessageCollection 获取消息集合
func (r *ChatRepository) getMessageCollection() *mongo.Collection {
	return r.db.Collection("ai_chat_messages")
}

// CreateSession 创建会话
func (r *ChatRepository) CreateSession(ctx context.Context, session *aiModels.ChatSession) error {
	session.BeforeCreate()

	collection := r.getSessionCollection()
	_, err := collection.InsertOne(ctx, session)
	if err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}

	return nil
}

// GetSessionByID 根据ID获取会话
func (r *ChatRepository) GetSessionByID(ctx context.Context, sessionID string) (*aiModels.ChatSession, error) {
	collection := r.getSessionCollection()

	var session aiModels.ChatSession
	filter := bson.M{
		"session_id": sessionID,
		"deleted_at": bson.M{"$exists": false},
	}

	err := collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("会话不存在")
		}
		return nil, fmt.Errorf("获取会话失败: %w", err)
	}

	// 加载消息
	messages, err := r.GetMessagesBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("加载消息失败: %w", err)
	}
	session.Messages = messages

	return &session, nil
}

// GetSessionsByProjectID 根据项目ID获取会话列表
func (r *ChatRepository) GetSessionsByProjectID(ctx context.Context, projectID string, limit, offset int) ([]*aiModels.ChatSession, error) {
	collection := r.getSessionCollection()

	filter := bson.M{
		"project_id": projectID,
		"deleted_at": bson.M{"$exists": false},
	}

	opts := options.Find().
		SetSort(bson.M{"updated_at": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询会话失败: %w", err)
	}
	defer cursor.Close(ctx)

	var sessions []*aiModels.ChatSession
	for cursor.Next(ctx) {
		var session aiModels.ChatSession
		if err := cursor.Decode(&session); err != nil {
			return nil, fmt.Errorf("解码会话失败: %w", err)
		}

		// 加载最新的几条消息
		messages, err := r.GetRecentMessagesBySessionID(ctx, session.SessionID, 10)
		if err == nil {
			session.Messages = messages
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// UpdateSession 更新会话
func (r *ChatRepository) UpdateSession(ctx context.Context, session *aiModels.ChatSession) error {
	session.BeforeUpdate()

	collection := r.getSessionCollection()
	filter := bson.M{"session_id": session.SessionID}

	update := bson.M{
		"$set": bson.M{
			"title":       session.Title,
			"description": session.Description,
			"status":      session.Status,
			"settings":    session.Settings,
			"metadata":    session.Metadata,
			"updated_at":  session.UpdatedAt,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("更新会话失败: %w", err)
	}

	return nil
}

// DeleteSession 删除会话（软删除）
func (r *ChatRepository) DeleteSession(ctx context.Context, sessionID string) error {
	collection := r.getSessionCollection()
	filter := bson.M{"session_id": sessionID}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("删除会话失败: %w", err)
	}

	// 同时软删除相关消息
	return r.DeleteMessagesBySessionID(ctx, sessionID)
}

// CreateMessage 创建消息
func (r *ChatRepository) CreateMessage(ctx context.Context, message *aiModels.ChatMessage) error {
	if message.ID.IsZero() {
		message.ID = primitive.NewObjectID()
	}
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	collection := r.getMessageCollection()
	_, err := collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("创建消息失败: %w", err)
	}

	return nil
}

// GetMessagesBySessionID 根据会话ID获取消息
func (r *ChatRepository) GetMessagesBySessionID(ctx context.Context, sessionID string) ([]aiModels.ChatMessage, error) {
	collection := r.getMessageCollection()

	filter := bson.M{
		"session_id": sessionID,
		"deleted_at": bson.M{"$exists": false},
	}

	opts := options.Find().SetSort(bson.M{"timestamp": 1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []aiModels.ChatMessage
	for cursor.Next(ctx) {
		var message aiModels.ChatMessage
		if err := cursor.Decode(&message); err != nil {
			return nil, fmt.Errorf("解码消息失败: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetRecentMessagesBySessionID 获取会话的最近消息
func (r *ChatRepository) GetRecentMessagesBySessionID(ctx context.Context, sessionID string, limit int) ([]aiModels.ChatMessage, error) {
	collection := r.getMessageCollection()

	filter := bson.M{
		"session_id": sessionID,
		"deleted_at": bson.M{"$exists": false},
	}

	opts := options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []aiModels.ChatMessage
	for cursor.Next(ctx) {
		var message aiModels.ChatMessage
		if err := cursor.Decode(&message); err != nil {
			return nil, fmt.Errorf("解码消息失败: %w", err)
		}
		messages = append(messages, message)
	}

	// 反转顺序，使其按时间正序排列
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// DeleteMessagesBySessionID 删除会话的所有消息（软删除）
func (r *ChatRepository) DeleteMessagesBySessionID(ctx context.Context, sessionID string) error {
	collection := r.getMessageCollection()
	filter := bson.M{"session_id": sessionID}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("删除消息失败: %w", err)
	}

	return nil
}

// GetSessionStatistics 获取会话统计信息
func (r *ChatRepository) GetSessionStatistics(ctx context.Context, projectID string) (*ChatStatistics, error) {
	collection := r.getSessionCollection()

	// 统计总会话数
	totalSessions, err := collection.CountDocuments(ctx, bson.M{
		"project_id": projectID,
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, fmt.Errorf("统计会话数失败: %w", err)
	}

	// 统计活跃会话数
	activeSessions, err := collection.CountDocuments(ctx, bson.M{
		"project_id": projectID,
		"status":     "active",
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, fmt.Errorf("统计活跃会话数失败: %w", err)
	}

	// 统计总消息数
	messageCollection := r.getMessageCollection()
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "ai_chat_sessions",
				"localField":   "session_id",
				"foreignField": "session_id",
				"as":           "session",
			},
		},
		{
			"$match": bson.M{
				"session.project_id": projectID,
				"deleted_at":         bson.M{"$exists": false},
			},
		},
		{
			"$count": "total",
		},
	}

	cursor, err := messageCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("统计消息数失败: %w", err)
	}
	defer cursor.Close(ctx)

	var totalMessages int64
	if cursor.Next(ctx) {
		var result struct {
			Total int64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err == nil {
			totalMessages = result.Total
		}
	}

	return &ChatStatistics{
		TotalSessions:  int(totalSessions),
		ActiveSessions: int(activeSessions),
		TotalMessages:  int(totalMessages),
	}, nil
}

// ChatStatistics 聊天统计信息
type ChatStatistics struct {
	TotalSessions  int `json:"totalSessions"`
	ActiveSessions int `json:"activeSessions"`
	TotalMessages  int `json:"totalMessages"`
}

// EnsureIndexes 确保索引存在
func (r *ChatRepository) EnsureIndexes(ctx context.Context) error {
	// 会话集合索引
	sessionCollection := r.getSessionCollection()
	sessionIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"session_id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"project_id": 1, "updated_at": -1},
		},
		{
			Keys: bson.M{"user_id": 1},
		},
		{
			Keys: bson.M{"status": 1},
		},
		{
			Keys:    bson.M{"deleted_at": 1},
			Options: options.Index().SetSparse(true),
		},
	}

	_, err := sessionCollection.Indexes().CreateMany(ctx, sessionIndexes)
	if err != nil {
		return fmt.Errorf("创建会话索引失败: %w", err)
	}

	// 消息集合索引
	messageCollection := r.getMessageCollection()
	messageIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{"session_id": 1, "timestamp": 1},
		},
		{
			Keys: bson.M{"role": 1},
		},
		{
			Keys:    bson.M{"deleted_at": 1},
			Options: options.Index().SetSparse(true),
		},
	}

	_, err = messageCollection.Indexes().CreateMany(ctx, messageIndexes)
	if err != nil {
		return fmt.Errorf("创建消息索引失败: %w", err)
	}

	return nil
}
