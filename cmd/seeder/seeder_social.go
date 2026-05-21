// Package main 提供社交数据填充功能
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"
	socialModel "Qingyu_backend/models/social"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SocialSeeder 社交数据填充器
type SocialSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

type socialSeedUser struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Nickname string             `bson:"nickname"`
	Avatar   string             `bson:"avatar"`
}

type socialSeedBook struct {
	ID           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	Author       string             `bson:"author"`
	Introduction string             `bson:"introduction"`
	Cover        string             `bson:"cover"`
	Categories   []string           `bson:"categories"`
	Tags         []string           `bson:"tags"`
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
	users, err := s.getUsers(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	books, err := s.getBooks(ctx)
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

	userIDs := socialSeedUserIDs(users)
	bookIDs := socialSeedBookIDs(books)

	if err := s.seedPosts(ctx, users, books); err != nil {
		return err
	}

	// 用于收集每本书的评分
	bookRatings := make(map[string][]float64)

	// 创建评论
	if err := s.seedComments(ctx, userIDs, bookIDs, bookRatings); err != nil {
		return err
	}

	// 更新书籍评分
	if err := s.updateBookRatings(ctx, bookRatings); err != nil {
		return err
	}

	// 创建点赞
	if err := s.seedLikes(ctx, userIDs, bookIDs); err != nil {
		return err
	}

	// 创建收藏
	if err := s.seedCollections(ctx, userIDs, bookIDs); err != nil {
		return err
	}

	// 创建关注
	if err := s.seedFollows(ctx, userIDs); err != nil {
		return err
	}

	return nil
}

// getUsers 获取用户列表
func (s *SocialSeeder) getUsers(ctx context.Context) ([]socialSeedUser, error) {
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []socialSeedUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// getBooks 获取书籍列表
func (s *SocialSeeder) getBooks(ctx context.Context) ([]socialSeedBook, error) {
	cursor, err := s.db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []socialSeedBook
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}
	return books, nil
}

func socialSeedUserIDs(users []socialSeedUser) []string {
	userIDs := make([]string, 0, len(users))
	for _, user := range users {
		if user.ID.IsZero() {
			continue
		}
		userIDs = append(userIDs, user.ID.Hex())
	}
	return userIDs
}

func socialSeedBookIDs(books []socialSeedBook) []string {
	bookIDs := make([]string, 0, len(books))
	for _, book := range books {
		if book.ID.IsZero() {
			continue
		}
		bookIDs = append(bookIDs, book.ID.Hex())
	}
	return bookIDs
}

func (s *SocialSeeder) seedPosts(ctx context.Context, users []socialSeedUser, books []socialSeedBook) error {
	collection := s.db.Collection("posts")

	existing, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计动态数量失败: %w", err)
	}
	if existing > 0 {
		fmt.Printf("  动态集合已有 %d 条数据，跳过帖子播种\n", existing)
		return nil
	}

	posts := buildSeedPosts(users, books)
	if len(posts) == 0 {
		fmt.Println("  没有可用于生成动态的用户或书籍，跳过帖子播种")
		return nil
	}

	if _, err := collection.InsertMany(ctx, posts); err != nil {
		return fmt.Errorf("插入社区动态失败: %w", err)
	}

	fmt.Printf("  创建了 %d 条社区动态\n", len(posts))
	return nil
}

func buildSeedPosts(users []socialSeedUser, books []socialSeedBook) []interface{} {
	if len(users) == 0 || len(books) == 0 {
		return nil
	}

	limit := len(books)
	if limit > 12 {
		limit = 12
	}

	posts := make([]interface{}, 0, limit)
	for i := 0; i < limit; i++ {
		user := users[i%len(users)]
		book := books[i]
		topics := buildPostTopics(book)
		createdAt := time.Now().Add(-time.Duration(i*6+rand.Intn(6)) * time.Hour)

		post := bson.M{
			"_id":           primitive.NewObjectID(),
			"user_id":       user.ID.Hex(),
			"user_name":     firstNonEmpty(strings.TrimSpace(user.Nickname), strings.TrimSpace(user.Username), "测试读者"),
			"user_avatar":   strings.TrimSpace(user.Avatar),
			"user_level":    1 + (i % 20),
			"type":          socialModel.PostTypeBookRecommendation,
			"content":       buildPostContent(book, topics),
			"images":        []string{},
			"book_id":       book.ID.Hex(),
			"book_title":    book.Title,
			"book_cover":    book.Cover,
			"book_author":   book.Author,
			"chapter_id":    "",
			"chapter_title": "",
			"progress":      0,
			"topics":        topics,
			"like_count":    24 + rand.Intn(180),
			"comment_count": 3 + rand.Intn(36),
			"share_count":   rand.Intn(18),
			"created_at":    createdAt,
			"updated_at":    createdAt,
		}

		if i%3 == 1 {
			post["type"] = socialModel.PostTypeReadingProgress
			post["content"] = fmt.Sprintf("追读《%s》到最新章节了，节奏比我预期还稳，准备继续蹲更新。", book.Title)
			post["chapter_id"] = fmt.Sprintf("seed-chapter-%02d", i+1)
			post["chapter_title"] = fmt.Sprintf("第%d章 精彩延续", 30+i)
			post["progress"] = 40 + (i*7)%55
		}
		if i%4 == 3 {
			post["type"] = socialModel.PostTypeText
			post["content"] = fmt.Sprintf("最近在补《%s》，%s 这个话题真的越聊越上头，有同好吗？", book.Title, firstNonEmpty(firstString(topics), "推荐"))
		}

		posts = append(posts, post)
	}

	return posts
}

func buildPostTopics(book socialSeedBook) []string {
	topics := make([]string, 0, 4)
	topics = appendIfMissing(topics, "推荐")
	for _, category := range book.Categories {
		topics = appendIfMissing(topics, category)
		if len(topics) >= 3 {
			break
		}
	}
	for _, tag := range book.Tags {
		topics = appendIfMissing(topics, tag)
		if len(topics) >= 4 {
			break
		}
	}
	if len(topics) == 1 {
		topics = append(topics, "书评")
	}
	return topics
}

func buildPostContent(book socialSeedBook, topics []string) string {
	lead := firstNonEmpty(book.Introduction, fmt.Sprintf("这本 %s 题材作品最近讨论度很高。", firstNonEmpty(firstString(topics), "小说")))
	lead = strings.TrimSpace(lead)
	runes := []rune(lead)
	if len(runes) > 40 {
		lead = string(runes[:40]) + "..."
	}
	return fmt.Sprintf("刚刷到《%s》，%s 写得很有记忆点。%s 这个方向值得继续聊。%s", book.Title, book.Author, firstNonEmpty(firstString(topics), "推荐", "书评"), lead)
}

func appendIfMissing(values []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func firstString(values []string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// seedComments 创建评论
func (s *SocialSeeder) seedComments(ctx context.Context, users, books []string, bookRatings map[string][]float64) error {
	collection := s.db.Collection("comments")
	existing, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计评论数量失败: %w", err)
	}
	if existing > 0 {
		fmt.Printf("  评论集合已有 %d 条数据，跳过评论播种\n", existing)
		return nil
	}

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
	existing, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计点赞数量失败: %w", err)
	}
	if existing > 0 {
		fmt.Printf("  点赞集合已有 %d 条数据，跳过点赞播种\n", existing)
		return nil
	}

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
	existing, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计收藏数量失败: %w", err)
	}
	if existing > 0 {
		fmt.Printf("  收藏集合已有 %d 条数据，跳过收藏播种\n", existing)
		return nil
	}

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
	existing, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计关注数量失败: %w", err)
	}
	if existing > 0 {
		fmt.Printf("  关注集合已有 %d 条数据，跳过关注播种\n", existing)
		return nil
	}

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
		bookObjectID, err := primitive.ObjectIDFromHex(bookID)
		if err != nil {
			return fmt.Errorf("书籍ID格式无效 (书籍ID: %s): %w", bookID, err)
		}

		// 计算平均评分
		sum := 0.0
		for _, r := range ratings {
			sum += r
		}
		avgRating := sum / float64(len(ratings))

		// 更新书籍
		_, err = booksCollection.UpdateOne(
			ctx,
			bson.M{"_id": bookObjectID},
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

	collections := []string{"post_likes", "posts", "comments", "likes", "collections", "follows"}

	for _, collName := range collections {
		_, err := s.db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空 %s 集合失败: %w", collName, err)
		}
	}

	fmt.Println("  已清空社交数据集合")
	return nil
}
