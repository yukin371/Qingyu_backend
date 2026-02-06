// Package main 提供社交数据填充功能
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SocialSeeder 社交数据填充器
type SocialSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

// NewSocialSeeder 创建社交数据填充器
func NewSocialSeeder(db *utils.Database, cfg *config.Config) *SocialSeeder {
	return &SocialSeeder{
		db:     db,
		config: cfg,
	}
}

// SeedSocialData 填充所有社交数据
func (s *SocialSeeder) SeedSocialData() error {
	ctx := context.Background()

	// 获取用户和书籍
	users, err := s.getUserIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	books, err := s.getBookIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("  没有找到用户，请先运行 users 命令创建用户")
		return nil
	}

	if len(books) == 0 {
		fmt.Println("  没有找到书籍，请先运行 bookstore 命令创建书籍")
		return nil
	}

	// 用于收集每本书的评分
	bookRatings := make(map[string][]float64)

	// 创建评论
	if err := s.seedComments(ctx, users, books, bookRatings); err != nil {
		return err
	}

	// 更新书籍评分
	if err := s.updateBookRatings(ctx, bookRatings); err != nil {
		return err
	}

	// 创建点赞
	if err := s.seedLikes(ctx, users, books); err != nil {
		return err
	}

	// 创建收藏
	if err := s.seedCollections(ctx, users, books); err != nil {
		return err
	}

	// 创建关注
	if err := s.seedFollows(ctx, users); err != nil {
		return err
	}

	return nil
}

// getUserIDs 获取用户ID列表
func (s *SocialSeeder) getUserIDs(ctx context.Context) ([]string, error) {
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

// getBookIDs 获取书籍ID列表
func (s *SocialSeeder) getBookIDs(ctx context.Context) ([]string, error) {
	cursor, err := s.db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	bookIDs := make([]string, len(books))
	for i, b := range books {
		bookIDs[i] = b.ID
	}
	return bookIDs, nil
}

// seedComments 创建评论
func (s *SocialSeeder) seedComments(ctx context.Context, users, books []string, bookRatings map[string][]float64) error {
	collection := s.db.Collection("comments")

	var comments []interface{}

	// 为每本书创建评论
	for _, bookID := range books {
		// 每本书5-20条评论
		commentCount := 5 + rand.Intn(16)

		for i := 0; i < commentCount; i++ {
			userID := users[rand.Intn(len(users))]
			rating := 3 + rand.Intn(3) // 3-5星

			// 收集评分用于后续更新书籍
			bookRatings[bookID] = append(bookRatings[bookID], float64(rating))

			comments = append(comments, bson.M{
				"_id":         primitive.NewObjectID(),
				"target_id":   bookID,
				"target_type": "book",
				"user_id":     userID,
				"content":     s.getRandomComment(),
				"rating":      rating,
				"like_count":  rand.Intn(50),
				"reply_count": rand.Intn(10),
				"status":      "normal",
				"created_at":  time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
				"updated_at":  time.Now(),
			})
		}
	}

	if len(comments) > 0 {
		if _, err := collection.InsertMany(ctx, comments); err != nil {
			return fmt.Errorf("插入评论失败: %w", err)
		}
		fmt.Printf("  创建了 %d 条评论\n", len(comments))
	}

	return nil
}

// seedLikes 创建点赞
func (s *SocialSeeder) seedLikes(ctx context.Context, users, books []string) error {
	collection := s.db.Collection("likes")

	var likes []interface{}

	// 为每本书创建点赞
	for _, bookID := range books {
		// 每本书10-100个点赞
		likeCount := 10 + rand.Intn(91)

		// 使用 Fisher-Yates 洗牌算法选择用户
		shuffledUsers := make([]string, len(users))
		copy(shuffledUsers, users)
		rand.Shuffle(len(shuffledUsers), func(i, j int) {
			shuffledUsers[i], shuffledUsers[j] = shuffledUsers[j], shuffledUsers[i]
		})

		for i := 0; i < likeCount && i < len(shuffledUsers); i++ {
			likes = append(likes, bson.M{
				"_id":         primitive.NewObjectID(),
				"user_id":     shuffledUsers[i],
				"target_id":   bookID,
				"target_type": "book",
				"created_at":  time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
			})
		}
	}

	if len(likes) > 0 {
		if _, err := collection.InsertMany(ctx, likes); err != nil {
			return fmt.Errorf("插入点赞失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个点赞\n", len(likes))
	}

	return nil
}

// seedCollections 创建收藏
func (s *SocialSeeder) seedCollections(ctx context.Context, users, books []string) error {
	collection := s.db.Collection("collections")

	var collections []interface{}

	// 为每个用户创建收藏
	for _, userID := range users {
		// 每个用户5-30个收藏
		collectionCount := 5 + rand.Intn(26)

		// 使用 Fisher-Yates 洗牌算法选择书籍
		shuffledBooks := make([]string, len(books))
		copy(shuffledBooks, books)
		rand.Shuffle(len(shuffledBooks), func(i, j int) {
			shuffledBooks[i], shuffledBooks[j] = shuffledBooks[j], shuffledBooks[i]
		})

		for i := 0; i < collectionCount && i < len(shuffledBooks); i++ {
			collections = append(collections, bson.M{
				"_id":         primitive.NewObjectID(),
				"user_id":     userID,
				"book_id":     shuffledBooks[i],
				"folder_name": "我的书架",
				"note":        "",
				"is_public":   rand.Intn(2) == 1,
				"created_at":  time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
			})
		}
	}

	if len(collections) > 0 {
		if _, err := collection.InsertMany(ctx, collections); err != nil {
			return fmt.Errorf("插入收藏失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个收藏\n", len(collections))
	}

	return nil
}

// seedFollows 创建关注
func (s *SocialSeeder) seedFollows(ctx context.Context, users []string) error {
	collection := s.db.Collection("follows")

	var follows []interface{}

	// 为每个用户创建关注
	for _, userID := range users {
		// 每个用户5-50个关注
		followCount := 5 + rand.Intn(46)

		// 使用 Fisher-Yates 洗牌算法选择被关注者
		shuffledUsers := make([]string, len(users))
		copy(shuffledUsers, users)
		rand.Shuffle(len(shuffledUsers), func(i, j int) {
			shuffledUsers[i], shuffledUsers[j] = shuffledUsers[j], shuffledUsers[i]
		})

		for i := 0; i < followCount && i < len(shuffledUsers); i++ {
			// 不能关注自己
			if shuffledUsers[i] == userID {
				continue
			}

			follows = append(follows, bson.M{
				"_id":          primitive.NewObjectID(),
				"follower_id":  userID,
				"following_id": shuffledUsers[i],
				"created_at":   time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
			})
		}
	}

	if len(follows) > 0 {
		if _, err := collection.InsertMany(ctx, follows); err != nil {
			return fmt.Errorf("插入关注失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个关注\n", len(follows))
	}

	return nil
}

// getRandomComment 获取随机评论内容
func (s *SocialSeeder) getRandomComment() string {
	comments := []string{
		"很好看，推荐！",
		"太精彩了，作者加油！",
		"剧情紧凑，人物鲜明。",
		"这本书值得一看。",
		"非常有趣的故事。",
		"期待后续更新！",
		"文笔流畅，情节引人入胜。",
		"题材新颖，想象力丰富。",
		"一口气读完，太爽了！",
		"希望作者能继续努力。",
		"这是一本不可多得的好书。",
		"强烈推荐给大家！",
		"支持作者，继续加油！",
		"非常喜欢这个风格。",
		"故事情节跌宕起伏，很吸引人。",
	}
	return comments[rand.Intn(len(comments))]
}

// updateBookRatings 更新书籍评分
func (s *SocialSeeder) updateBookRatings(ctx context.Context, bookRatings map[string][]float64) error {
	booksCollection := s.db.Collection("books")

	for bookID, ratings := range bookRatings {
		if len(ratings) == 0 {
			continue
		}

		// 计算平均评分
		sum := 0.0
		for _, r := range ratings {
			sum += r
		}
		avgRating := sum / float64(len(ratings))

		// 更新书籍
		_, err := booksCollection.UpdateOne(
			ctx,
			bson.M{"_id": bookID},
			bson.M{
				"$set": bson.M{
					"rating":       roundToOneDecimal(avgRating),
					"rating_count": len(ratings),
					"updated_at":   time.Now(),
				},
			},
		)
		if err != nil {
			return fmt.Errorf("更新书籍评分失败 (书籍ID: %s): %w", bookID, err)
		}
	}

	fmt.Printf("  更新了 %d 本书的评分\n", len(bookRatings))
	return nil
}

// roundToOneDecimal 保留一位小数
func roundToOneDecimal(n float64) float64 {
	return math.Round(n*10) / 10
}

// Clean 清空社交数据
func (s *SocialSeeder) Clean() error {
	ctx := context.Background()

	collections := []string{"comments", "likes", "collections", "follows"}

	for _, collName := range collections {
		_, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}

	fmt.Println("  已清空社交数据集合")
	return nil
}
