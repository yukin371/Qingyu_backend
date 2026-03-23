// Package main 提供收藏夹数据填充功能
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

// CollectionFolderSeeder 收藏夹数据填充器
type CollectionFolderSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

// NewCollectionFolderSeeder 创建收藏夹数据填充器
func NewCollectionFolderSeeder(db *utils.Database, cfg *config.Config) *CollectionFolderSeeder {
	collection := db.Collection("collection_folders")
	return &CollectionFolderSeeder{
		db:       db,
		config:   cfg,
		inserter: utils.NewBulkInserter(collection, cfg.BatchSize),
	}
}

// SeedCollectionFolders 填充收藏夹数据
func (s *CollectionFolderSeeder) SeedCollectionFolders() error {
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

	// 创建收藏夹
	if err := s.seedFolders(ctx, users); err != nil {
		return err
	}

	return nil
}

// getUserIDs 获取用户ID列表
func (s *CollectionFolderSeeder) getUserIDs(ctx context.Context) ([]string, error) {
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

// seedFolders 创建收藏夹
func (s *CollectionFolderSeeder) seedFolders(ctx context.Context, userIDs []string) error {
	collection := s.db.Collection("collection_folders")

	var folders []interface{}
	now := time.Now()

	// 预定义收藏夹模板
	folderTemplates := []struct {
		name        string
		description string
		isPublic    bool
	}{
		{"我的书架", "默认收藏夹", true},
		{"想读的书", "想读但还没开始的书", false},
		{"已读完", "已经读完的书籍", false},
		{"推荐收藏", "值得推荐的优质作品", true},
		{"追更中", "正在追更的连载作品", false},
		{"经典必读", "经典好书收藏", true},
	}

	// 为每个用户创建收藏夹
	for _, userID := range userIDs {
		// 每个用户 2-4 个收藏夹
		folderCount := 2 + rand.Intn(3)

		for i := 0; i < folderCount && i < len(folderTemplates); i++ {
			template := folderTemplates[i]

			folder := bson.M{
				"_id":         primitive.NewObjectID(),
				"user_id":     userID,
				"name":        template.name,
				"description": template.description,
				"book_count":  rand.Intn(20),
				"is_public":   template.isPublic,
				"created_at":  now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
				"updated_at":  now,
			}
			folders = append(folders, folder)
		}
	}

	// 批量插入
	if len(folders) > 0 {
		if _, err := collection.InsertMany(ctx, folders); err != nil {
			return fmt.Errorf("插入收藏夹失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个收藏夹\n", len(folders))
	}

	return nil
}

// Clean 清空收藏夹数据
func (s *CollectionFolderSeeder) Clean() error {
	ctx := context.Background()

	_, err := s.db.Collection("collection_folders").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 collection_folders 集合失败: %w", err)
	}

	fmt.Println("  已清空 collection_folders 集合")
	return nil
}
