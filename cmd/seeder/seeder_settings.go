// Package main 提供用户设置数据填充功能
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SettingsSeeder 用户设置数据填充器
type SettingsSeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewSettingsSeeder 创建用户设置填充器
func NewSettingsSeeder(db *utils.Database, cfg *config.Config) *SettingsSeeder {
	return &SettingsSeeder{
		db:     db,
		config: cfg,
	}
}

// ReadingPreferences 阅读偏好设置
type ReadingPreferences struct {
	FontSize     float64 `bson:"font_size" json:"fontSize"`
	Theme        string  `bson:"theme" json:"theme"`
	LineHeight   float64 `bson:"line_height" json:"lineHeight"`
	AutoPageTurn bool    `bson:"auto_page_turn" json:"autoPageTurn"`
}

// NotificationSettings 通知设置
type NotificationSettings struct {
	CommentNotification bool `bson:"comment_notification" json:"commentNotification"`
	UpdateNotification  bool `bson:"update_notification" json:"updateNotification"`
	SystemNotification  bool `bson:"system_notification" json:"systemNotification"`
	EmailNotification   bool `bson:"email_notification" json:"emailNotification"`
}

// PrivacySettings 隐私设置
type PrivacySettings struct {
	ReadingHistoryVisibility string `bson:"reading_history_visibility" json:"readingHistoryVisibility"`
	BookshelfVisibility      string `bson:"bookshelf_visibility" json:"bookshelfVisibility"`
}

// UserSettings 用户设置
type UserSettings struct {
	ID                  string              `bson:"_id" json:"id"`
	UserID              string              `bson:"user_id" json:"userId"`
	ReadingPreferences  ReadingPreferences  `bson:"reading_preferences" json:"readingPreferences"`
	NotificationSettings NotificationSettings `bson:"notification_settings" json:"notificationSettings"`
	PrivacySettings     PrivacySettings     `bson:"privacy_settings" json:"privacySettings"`
	CreatedAt           time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt           time.Time           `bson:"updated_at" json:"updatedAt"`
}

// SeedUserSettings 填充用户设置数据
func (s *SettingsSeeder) SeedUserSettings() error {
	ctx := context.Background()

	// 获取所有用户
	users, err := s.getUserIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("  没有找到用户，请先运行 users 命令创建用户")
		return nil
	}

	// 为每个用户生成设置
	var settings []interface{}
	for _, userID := range users {
		setting := s.generateUserSettings(userID)
		settings = append(settings, setting)
	}

	// 批量插入
	collection := s.db.Collection("user_settings")
	if len(settings) > 0 {
		// 先清空现有数据
		_, err := collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空现有设置数据失败: %w", err)
		}

		// 批量插入新数据
		_, err = collection.InsertMany(ctx, settings)
		if err != nil {
			return fmt.Errorf("插入用户设置失败: %w", err)
		}

		fmt.Printf("  成功为 %d 个用户创建了设置\n", len(settings))
	}

	return nil
}

// getUserIDs 获取用户ID列表
func (s *SettingsSeeder) getUserIDs(ctx context.Context) ([]string, error) {
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

// generateUserSettings 生成用户设置
func (s *SettingsSeeder) generateUserSettings(userID string) bson.M {
	// 随机主题：70% light，30% dark
	theme := "light"
	if rand.Intn(10) < 3 {
		theme = "dark"
	}

	// 随机字体大小：14-20
	fontSize := 14 + rand.Intn(7)

	// 随机行间距：1.4-2.0
	lineHeight := 1.4 + rand.Float64()*0.6

	// 随机可见性：60% public，30% friends，10% private
	var readingVisibility, bookshelfVisibility string
	randVal := rand.Intn(10)
	if randVal < 6 {
		readingVisibility = "public"
	} else if randVal < 9 {
		readingVisibility = "friends"
	} else {
		readingVisibility = "private"
	}

	randVal = rand.Intn(10)
	if randVal < 6 {
		bookshelfVisibility = "public"
	} else if randVal < 9 {
		bookshelfVisibility = "friends"
	} else {
		bookshelfVisibility = "private"
	}

	now := time.Now()

	return bson.M{
		"_id":         primitive.NewObjectID(),
		"user_id":     userID,
		"reading_preferences": bson.M{
			"font_size":      fontSize,
			"theme":          theme,
			"line_height":    lineHeight,
			"auto_page_turn": rand.Intn(2) == 1, // 50%概率
		},
		"notification_settings": bson.M{
			"comment_notification": rand.Intn(10) < 8, // 80%开启
			"update_notification":  rand.Intn(10) < 8, // 80%开启
			"system_notification":  rand.Intn(10) < 9, // 90%开启
			"email_notification":   rand.Intn(10) < 3, // 30%开启
		},
		"privacy_settings": bson.M{
			"reading_history_visibility": readingVisibility,
			"bookshelf_visibility":       bookshelfVisibility,
		},
		"created_at": now,
		"updated_at": now,
	}
}

// Clean 清空设置数据
func (s *SettingsSeeder) Clean() error {
	ctx := context.Background()
	_, err := s.db.Collection("user_settings").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 user_settings 集合失败: %w", err)
	}
	fmt.Println("  已清空 user_settings 集合")
	return nil
}

// Count 统计设置数量
func (s *SettingsSeeder) Count() (int64, error) {
	ctx := context.Background()
	return s.db.Collection("user_settings").CountDocuments(ctx, bson.M{})
}
