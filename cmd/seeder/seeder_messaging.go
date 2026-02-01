// Package main 提供消息数据填充功能
package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/models"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessagingSeeder 消息数据填充器
type MessagingSeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewMessagingSeeder 创建消息数据填充器
func NewMessagingSeeder(db *utils.Database, cfg *config.Config) *MessagingSeeder {
	return &MessagingSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedMessagingData 填充所有消息数据
func (s *MessagingSeeder) SeedMessagingData() error {
	ctx := context.Background()

	// 获取用户列表
	users, err := s.getUserIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	if len(users) < 2 {
		fmt.Println("  用户数量不足，至少需要2个用户")
		return nil
	}

	// 填充对话和消息
	if err := s.seedConversationsAndMessages(ctx, users); err != nil {
		return err
	}

	// 填充公告
	if err := s.seedAnnouncements(ctx); err != nil {
		return err
	}

	return nil
}

// seedConversationsAndMessages 创建对话和消息
func (s *MessagingSeeder) seedConversationsAndMessages(ctx context.Context, users []string) error {
	conversationCollection := s.db.Collection("conversations")
	messageCollection := s.db.Collection("messages")

	// 30%的用户有对话
	userCount := len(users) * 3 / 10
	activeUsers := users[:userCount]

	var allMessages []interface{}

	// 创建对话
	for i := 0; i < userCount; i++ {
		user1 := activeUsers[i]
		user2 := users[(i+1)%len(users)]

		conversationID := primitive.NewObjectID().Hex()
		now := time.Now()

		// 创建对话
		conversation := models.Conversation{
			ID:            conversationID,
			Participants:  []string{user1, user2},
			Type:          "private",
			Title:         "",
			LastMessageAt: now.Add(-time.Duration(i) * time.Hour),
			CreatedAt:     now.Add(-time.Duration(i*24) * time.Hour),
			UpdatedAt:     now,
		}

		_, err := conversationCollection.InsertOne(ctx, conversation)
		if err != nil {
			return fmt.Errorf("插入对话失败: %w", err)
		}

		// 为每个对话创建5-20条消息
		messageCount := 5 + i%16
		for j := 0; j < messageCount; j++ {
			isFromUser1 := j%2 == 0
			senderID := user1
			receiverID := user2
			if !isFromUser1 {
				senderID, receiverID = receiverID, senderID
			}

			isRead := j < messageCount-2 // 最后2条未读
			var readAt time.Time
			if isRead {
				readAt = now.Add(-time.Duration(j) * time.Hour)
			}

			message := models.Message{
				ID:             primitive.NewObjectID().Hex(),
				ConversationID: conversationID,
				SenderID:       senderID,
				ReceiverID:     receiverID,
				Content:        s.getRandomMessage(j),
				MessageType:    "text",
				IsRead:         isRead,
				ReadAt:         readAt,
				CreatedAt:      now.Add(-time.Duration((messageCount-j)*2) * time.Minute),
			}

			allMessages = append(allMessages, message)
		}
	}

	// 批量插入消息
	if len(allMessages) > 0 {
		batchSize := 100
		for i := 0; i < len(allMessages); i += batchSize {
			end := i + batchSize
			if end > len(allMessages) {
				end = len(allMessages)
			}

			_, err := messageCollection.InsertMany(ctx, allMessages[i:end])
			if err != nil {
				return fmt.Errorf("插入消息失败（批次 %d）: %w", i/batchSize, err)
			}
		}

		fmt.Printf("  创建了 %d 个对话\n", userCount)
		fmt.Printf("  创建了 %d 条消息\n", len(allMessages))
	}

	return nil
}

// seedAnnouncements 创建公告
func (s *MessagingSeeder) seedAnnouncements(ctx context.Context) error {
	collection := s.db.Collection("announcements")

	now := time.Now()

	announcements := []interface{}{
		models.Announcement{
			ID:          primitive.NewObjectID().Hex(),
			Title:       "欢迎使用青羽写作平台",
			Content:     "感谢您注册青羽写作平台，祝您创作愉快！",
			Type:        "system",
			Priority:    5,
			IsPublished: true,
			PublishedAt: now.Add(-30 * 24 * time.Hour),
			ExpireAt:    now.Add(365 * 24 * time.Hour),
			CreatedAt:   now.Add(-30 * 24 * time.Hour),
		},
		models.Announcement{
			ID:          primitive.NewObjectID().Hex(),
			Title:       "新功能上线通知",
			Content:     "我们上线了AI辅助写作功能，欢迎体验！",
			Type:        "event",
			Priority:    3,
			IsPublished: true,
			PublishedAt: now.Add(-7 * 24 * time.Hour),
			ExpireAt:    now.Add(30 * 24 * time.Hour),
			CreatedAt:   now.Add(-7 * 24 * time.Hour),
		},
		models.Announcement{
			ID:          primitive.NewObjectID().Hex(),
			Title:       "系统维护通知",
			Content:     "系统将于本周六凌晨2:00-4:00进行维护，届时将暂停服务。",
			Type:        "maintenance",
			Priority:    4,
			IsPublished: true,
			PublishedAt: now.Add(-1 * 24 * time.Hour),
			ExpireAt:    now.Add(7 * 24 * time.Hour),
			CreatedAt:   now.Add(-1 * 24 * time.Hour),
		},
	}

	_, err := collection.InsertMany(ctx, announcements)
	if err != nil {
		return fmt.Errorf("插入公告失败: %w", err)
	}

	fmt.Printf("  创建了 %d 条公告\n", len(announcements))
	return nil
}

// getRandomMessage 获取随机消息内容
func (s *MessagingSeeder) getRandomMessage(index int) string {
	messages := []string{
		"你好，在吗？",
		"最近在写什么作品？",
		"能不能帮我看一下我的作品？",
		"谢谢你的建议！",
		"好的，我知道了",
		"期待你的更新",
		"加油！",
		"一起努力",
		"有空多交流",
		"收到",
	}
	return messages[index%len(messages)]
}

// getUserIDs 获取用户ID列表
func (s *MessagingSeeder) getUserIDs(ctx context.Context) ([]string, error) {
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}
	return userIDs, nil
}

// Clean 清空消息数据
func (s *MessagingSeeder) Clean() error {
	ctx := context.Background()

	collections := []string{"conversations", "messages", "announcements"}

	for _, collName := range collections {
		_, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}

	fmt.Println("  已清空消息数据集合")
	return nil
}
