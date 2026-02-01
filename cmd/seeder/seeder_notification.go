// Package main 提供通知数据填充功能
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

// NotificationSeeder 通知数据填充器
type NotificationSeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewNotificationSeeder 创建通知数据填充器
func NewNotificationSeeder(db *utils.Database, cfg *config.Config) *NotificationSeeder {
	return &NotificationSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedNotifications 填充通知数据
func (s *NotificationSeeder) SeedNotifications() error {
	ctx := context.Background()

	// 获取用户列表
	users, err := s.getUserIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("  没有找到用户，请先运行 users 命令创建用户")
		return nil
	}

	collection := s.db.Collection("notifications")

	var notifications []interface{}

	// 通知类型分布
	notificationTypes := []struct {
		typ      string
		ratio    float64
		titles   []string
		contents []string
	}{
		{
			typ:    "comment",
			ratio:  0.4,
			titles: []string{"新评论通知", "有人评论了你的作品"},
			contents: []string{
				"用户{{user}}评论了你的作品《{{book}}》",
				"收到一条新评论：{{content}}",
			},
		},
		{
			typ:    "like",
			ratio:  0.3,
			titles: []string{"点赞通知", "有人赞了你的作品"},
			contents: []string{
				"用户{{user}}赞了你的作品《{{book}}》",
				"你的作品获得了一个点赞",
			},
		},
		{
			typ:    "follow",
			ratio:  0.2,
			titles: []string{"新粉丝通知", "有人关注了你"},
			contents: []string{
				"用户{{user}}关注了你",
				"你有了新的粉丝",
			},
		},
		{
			typ:    "system",
			ratio:  0.1,
			titles: []string{"系统通知", "平台公告"},
			contents: []string{
				"欢迎来到青羽写作平台",
				"你的作品已通过审核",
				"系统维护通知",
			},
		},
	}

	scale := config.GetScaleConfig(s.config.Scale)
	notificationsPerUser := 20 + scale.Users%30 // 每个用户20-50条通知

	for _, userID := range users {
		// 生成该用户的通知
		userNotifs := s.generateUserNotifications(userID, notificationsPerUser, notificationTypes)
		notifications = append(notifications, userNotifs...)
	}

	if len(notifications) > 0 {
		// 批量插入
		batchSize := 100
		for i := 0; i < len(notifications); i += batchSize {
			end := i + batchSize
			if end > len(notifications) {
				end = len(notifications)
			}

			_, err := collection.InsertMany(ctx, notifications[i:end])
			if err != nil {
				return fmt.Errorf("插入通知失败（批次 %d）: %w", i/batchSize, err)
			}
		}

		fmt.Printf("  创建了 %d 条通知\n", len(notifications))
	}

	return nil
}

// generateUserNotifications 为单个用户生成通知
func (s *NotificationSeeder) generateUserNotifications(userID string, count int, types []struct {
	typ      string
	ratio    float64
	titles   []string
	contents []string
}) []interface{} {
	var notifications []interface{}
	now := time.Now()

	for i := 0; i < count; i++ {
		// 根据比例选择通知类型
		rand := float64(i%10) / 10.0
		var typ string
		var titles, contents []string

		cumRatio := 0.0
		for _, t := range types {
			cumRatio += t.ratio
			if rand < cumRatio {
				typ = t.typ
				titles = t.titles
				contents = t.contents
				break
			}
		}

		// 70%已读，30%未读
		isRead := (i % 10) < 7
		var readAt time.Time
		if isRead {
			readAt = now.Add(-time.Duration(i) * time.Hour)
		}

		title := titles[i%len(titles)]
		content := contents[i%len(contents)]

		notifications = append(notifications, models.Notification{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    userID,
			Type:      typ,
			Title:     title,
			Content:   content,
			Link:      fmt.Sprintf("/%s", typ),
			IsRead:    isRead,
			ReadAt:    readAt,
			CreatedAt: now.Add(-time.Duration(i*2) * time.Hour),
		})
	}

	return notifications
}

// getUserIDs 获取用户ID列表
func (s *NotificationSeeder) getUserIDs(ctx context.Context) ([]string, error) {
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

// Clean 清空通知数据
func (s *NotificationSeeder) Clean() error {
	ctx := context.Background()
	_, err := s.db.Collection("notifications").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 notifications 集合失败: %w", err)
	}

	fmt.Println("  已清空 notifications 集合")
	return nil
}
