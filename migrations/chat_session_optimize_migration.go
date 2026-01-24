package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatSessionMigration 聊天会话迁移工具
// 用于将 ChatSession 内嵌的 Messages 数组拆分到独立的 ai_chat_messages 集合
type ChatSessionMigration struct {
	db         *mongo.Database
	sessionCol *mongo.Collection
	messageCol *mongo.Collection
}

// NewChatSessionMigration 创建迁移实例
func NewChatSessionMigration(db *mongo.Database) *ChatSessionMigration {
	return &ChatSessionMigration{
		db:         db,
		sessionCol: db.Collection("ai_chat_sessions"),
		messageCol: db.Collection("ai_chat_messages"),
	}
}

// BackupSessionMessages 备份会话中的 messages 数组
// 将包含 messages 的会话备份到 ai_chat_sessions_backup 集合
func (m *ChatSessionMigration) BackupSessionMessages(ctx context.Context) error {
	log.Println("开始备份会话消息...")

	// 查找所有包含 messages 字段的会话
	filter := bson.M{"messages": bson.M{"$exists": true, "$not": bson.M{"$size": 0}}}
	cursor, err := m.sessionCol.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("查询会话失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 创建备份集合
	backupCol := m.db.Collection("ai_chat_sessions_backup")

	var backedUpCount int
	for cursor.Next(ctx) {
		var session bson.M
		if err := cursor.Decode(&session); err != nil {
			log.Printf("解码会话失败: %v", err)
			continue
		}

		// 插入到备份集合
		_, err := backupCol.InsertOne(ctx, session)
		if err != nil {
			log.Printf("备份会话 %s 失败: %v", session["_id"], err)
			continue
		}

		backedUpCount++
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("遍历会话失败: %w", err)
	}

	log.Printf("备份完成，共备份 %d 个会话", backedUpCount)
	return nil
}

// MigrateSessionMessages 执行迁移
// 将会话中的 messages 数组拆分到独立集合
func (m *ChatSessionMigration) MigrateSessionMessages(ctx context.Context) error {
	log.Println("开始迁移会话消息...")

	// 查找所有包含 messages 字段的会话
	filter := bson.M{"messages": bson.M{"$exists": true, "$not": bson.M{"$size": 0}}}
	cursor, err := m.sessionCol.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("查询会话失败: %w", err)
	}
	defer cursor.Close(ctx)

	var migratedSessions int
	var migratedMessages int

	for cursor.Next(ctx) {
		var session struct {
			ID       primitive.ObjectID `bson:"_id"`
			SessionID string             `bson:"session_id"`
			Messages  []ChatMessage      `bson:"messages"`
		}

		if err := cursor.Decode(&session); err != nil {
			log.Printf("解码会话失败: %v", err)
			continue
		}

		// 迁移每条消息
		for _, msg := range session.Messages {
			// 检查消息是否已存在
			existingFilter := bson.M{"_id": msg.ID}
			count, err := m.messageCol.CountDocuments(ctx, existingFilter)
			if err != nil {
				log.Printf("检查消息 %s 失败: %v", msg.ID, err)
				continue
			}

			// 如果消息不存在，则插入
			if count == 0 {
				// 确保 session_id 正确
				if msg.SessionID == "" {
					msg.SessionID = session.SessionID
				}

				// 设置创建时间（如果不存在）
				if msg.CreatedAt.IsZero() {
					msg.CreatedAt = time.Now()
					msg.UpdatedAt = time.Now()
				}

				_, err := m.messageCol.InsertOne(ctx, msg)
				if err != nil {
					log.Printf("插入消息 %s 失败: %v", msg.ID, err)
					continue
				}

				migratedMessages++
			}
		}

		// 从会话文档中移除 messages 字段
		updateFilter := bson.M{"_id": session.ID}
		update := bson.M{"$unset": bson.M{"messages": ""}}
		_, err := m.sessionCol.UpdateOne(ctx, updateFilter, update)
		if err != nil {
			log.Printf("移除会话 %s 的 messages 字段失败: %v", session.ID, err)
			continue
		}

		migratedSessions++
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("遍历会话失败: %w", err)
	}

	log.Printf("迁移完成，共迁移 %d 个会话，%d 条消息", migratedSessions, migratedMessages)
	return nil
}

// VerifyMigration 验证迁移结果
func (m *ChatSessionMigration) VerifyMigration(ctx context.Context) error {
	log.Println("开始验证迁移结果...")

	// 1. 检查是否还有会话包含 messages 字段
	filter := bson.M{"messages": bson.M{"$exists": true, "$not": bson.M{"$size": 0}}}
	count, err := m.sessionCol.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("统计未迁移会话失败: %w", err)
	}

	if count > 0 {
		log.Printf("警告：仍有 %d 个会话包含 messages 字段", count)
	} else {
		log.Println("✓ 所有会话的 messages 字段已成功移除")
	}

	// 2. 统计消息总数
	messageCount, err := m.messageCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计消息数失败: %w", err)
	}
	log.Printf("✓ ai_chat_messages 集合共有 %d 条消息", messageCount)

	// 3. 统计会话总数
	sessionCount, err := m.sessionCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计会话数失败: %w", err)
	}
	log.Printf("✓ ai_chat_sessions 集合共有 %d 个会话", sessionCount)

	// 4. 检查孤立消息（没有对应会话的消息）
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
				"session": bson.M{"$size": 0},
			},
		},
		{
			"$count": "orphan_messages",
		},
	}

	cursor, err := m.messageCol.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("查找孤立消息失败: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			OrphanMessages int64 `bson:"orphan_messages"`
		}
		if err := cursor.Decode(&result); err == nil {
			if result.OrphanMessages > 0 {
				log.Printf("警告：发现 %d 条孤立消息", result.OrphanMessages)
			} else {
				log.Println("✓ 没有发现孤立消息")
			}
		}
	}

	log.Println("验证完成")
	return nil
}

// RollbackMigration 回滚迁移
// 从备份集合恢复会话的 messages 字段
func (m *ChatSessionMigration) RollbackMigration(ctx context.Context) error {
	log.Println("开始回滚迁移...")

	backupCol := m.db.Collection("ai_chat_sessions_backup")

	// 检查备份集合是否存在
	count, err := backupCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("检查备份集合失败: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("备份集合为空，无法回滚")
	}

	log.Printf("找到 %d 个备份会话", count)

	cursor, err := backupCol.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("查询备份会话失败: %w", err)
	}
	defer cursor.Close(ctx)

	var restoredCount int
	for cursor.Next(ctx) {
		var session bson.M
		if err := cursor.Decode(&session); err != nil {
			log.Printf("解码备份会话失败: %v", err)
			continue
		}

		sessionID := session["_id"]
		messages, ok := session["messages"]
		if !ok {
			continue
		}

		// 恢复 messages 字段
		filter := bson.M{"_id": sessionID}
		update := bson.M{"$set": bson.M{"messages": messages}}
		_, err := m.sessionCol.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("恢复会话 %s 失败: %v", sessionID, err)
			continue
		}

		restoredCount++
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("遍历备份会话失败: %w", err)
	}

	log.Printf("回滚完成，共恢复 %d 个会话", restoredCount)
	return nil
}

// ChatMessage 旧的消息模型（用于迁移）
type ChatMessage struct {
	ID        primitive.ObjectID `bson:"_id"`
	SessionID string             `bson:"session_id"`
	Role      string             `bson:"role"`
	Content   string             `bson:"content"`
	TokenUsed int                `bson:"token_used"`
	Metadata  *MessageMeta       `bson:"metadata"`
	Timestamp time.Time          `bson:"timestamp"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

// MessageMeta 消息元数据
type MessageMeta struct {
	ResponseTime float64                `bson:"response_time"`
	ModelUsed    string                 `bson:"model_used"`
	ContextUsed  bool                   `bson:"context_used"`
	MemoryUsed   bool                   `bson:"memory_used"`
	Sources      []string               `bson:"sources"`
	Confidence   float64                `bson:"confidence"`
	Sentiment    string                 `bson:"sentiment"`
	Intent       string                 `bson:"intent"`
	Entities     []Entity               `bson:"entities"`
	CustomData   map[string]interface{} `bson:"custom_data"`
}

// Entity 实体信息
type Entity struct {
	Type       string  `bson:"type"`
	Value      string  `bson:"value"`
	Confidence float64 `bson:"confidence"`
	StartPos   int     `bson:"start_pos"`
	EndPos     int     `bson:"end_pos"`
}

func main() {
	// 示例使用代码
	// 实际使用时需要根据项目配置调整

	// 连接 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接 MongoDB 失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	migration := NewChatSessionMigration(db)

	// 执行迁移流程
	fmt.Println("=== 聊天会话存储优化迁移 ===")
	fmt.Println("步骤 1/4: 备份原始数据")
	if err := migration.BackupSessionMessages(ctx); err != nil {
		log.Fatalf("备份失败: %v", err)
	}

	fmt.Println("\n步骤 2/4: 执行迁移")
	if err := migration.MigrateSessionMessages(ctx); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	fmt.Println("\n步骤 3/4: 验证迁移结果")
	if err := migration.VerifyMigration(ctx); err != nil {
		log.Fatalf("验证失败: %v", err)
	}

	fmt.Println("\n步骤 4/4: 清理备份（可选）")
	fmt.Println("如需回滚，请运行 RollbackMigration() 方法")
	fmt.Println("确认迁移成功后，可手动删除备份集合: ai_chat_sessions_backup")

	fmt.Println("\n=== 迁移完成 ===")
}
