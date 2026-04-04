// Package main 提供书单数据填充功能
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

// BooklistSeeder 书单数据填充器
type BooklistSeeder struct {
	db       *utils.Database
	config   *config.Config
	inserter *utils.BulkInserter
}

// NewBooklistSeeder 创建书单数据填充器
func NewBooklistSeeder(db *utils.Database, cfg *config.Config) *BooklistSeeder {
	collection := db.Collection("booklists")
	return &BooklistSeeder{
		db:       db,
		config:   cfg,
		inserter: utils.NewBulkInserter(collection, cfg.BatchSize),
	}
}

// SeedBooklists 填充书单数据
func (s *BooklistSeeder) SeedBooklists() error {
	ctx := context.Background()

	// 获取用户和书籍
	users, err := s.getUserIDs(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	books, err := s.getBooksWithAuthors(ctx)
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

	// 创建书单
	if err := s.seedBooklistsData(ctx, users, books); err != nil {
		return err
	}

	return nil
}

// getUserIDs 获取用户ID列表
func (s *BooklistSeeder) getUserIDs(ctx context.Context) ([]map[string]string, error) {
	cursor, err := s.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []struct {
		ID       string `bson:"_id"`
		Username string `bson:"username"`
		Avatar   string `bson:"avatar"`
	}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	result := make([]map[string]string, len(users))
	for i, u := range users {
		result[i] = map[string]string{
			"id":       u.ID,
			"username": u.Username,
			"avatar":   u.Avatar,
		}
	}
	return result, nil
}

// getBooksWithAuthors 获取书籍ID和作者信息
func (s *BooklistSeeder) getBooksWithAuthors(ctx context.Context) ([]map[string]interface{}, error) {
	cursor, err := s.db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []struct {
		ID          string  `bson:"_id"`
		Title       string  `bson:"title"`
		Cover       string  `bson:"cover"`
		Description string  `bson:"description"`
		AuthorID    string  `bson:"author_id"`
	}
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(books))
	for i, b := range books {
		result[i] = map[string]interface{}{
			"id":          b.ID,
			"title":       b.Title,
			"cover":       b.Cover,
			"description": b.Description,
			"author_id":   b.AuthorID,
		}
	}
	return result, nil
}

// seedBooklistsData 创建书单数据
func (s *BooklistSeeder) seedBooklistsData(ctx context.Context, users []map[string]string, books []map[string]interface{}) error {
	booklistCollection := s.db.Collection("booklists")
	likeCollection := s.db.Collection("booklist_likes")

	var booklists []interface{}
	var likes []interface{}
	now := time.Now()

	// 书单标题模板
	booklistTemplates := []struct {
		title       string
		description string
		category    string
	}{
		{"年度必读好书", "精选年度最受欢迎的优质作品，涵盖各类题材", "推荐"},
		{"适合深夜阅读的治愈系", "温暖治愈的故事，适合在宁静的夜晚慢慢品味", "治愈"},
		{"高能反转神作推荐", "剧情跌宕起伏、让人欲罢不能的精彩作品", "悬疑"},
		{"新人入坑指南", "适合新读者的入门佳作，轻松易懂", "入门"},
		{"经典文学名著", "经过时间考验的经典作品，值得一读再读", "经典"},
		{"甜宠文合集", "甜甜的恋爱故事，让你相信爱情", "言情"},
		{"热血爽文推荐", "节奏明快，让人看得停不下来", "爽文"},
		{"完结好文精选", "已完结的优质作品，一次看过瘾", "完结"},
	}

	// 为 30% 的用户创建书单
	userCount := len(users) * 3 / 10
	if userCount < 1 {
		userCount = 1
	}

	createdCount := 0
	for i := 0; i < userCount && i < len(users); i++ {
		user := users[i]

		// 每个用户创建 1-3 个书单
		booklistCount := 1 + rand.Intn(3)
		for j := 0; j < booklistCount; j++ {
			template := booklistTemplates[(i+j)%len(booklistTemplates)]

			// 随机选择 3-10 本书
			bookCount := 3 + rand.Intn(8)
			if bookCount > len(books) {
				bookCount = len(books)
			}

			// 洗牌选择书籍
			shuffledBooks := make([]map[string]interface{}, len(books))
			copy(shuffledBooks, books)
			rand.Shuffle(len(shuffledBooks), func(a, b int) {
				shuffledBooks[a], shuffledBooks[b] = shuffledBooks[b], shuffledBooks[a]
			})

			// 构建书单中的书籍列表
			var bookItems []interface{}
			selectedBooks := shuffledBooks[:bookCount]
			for k := 0; k < len(selectedBooks); k++ {
				book := selectedBooks[k]
				bookItems = append(bookItems, bson.M{
					"book_id":     book["id"],
					"book_title":  book["title"],
					"book_cover":  book["cover"],
					"author_name":  "",
					"description": "",
					"comment":     s.getRandomComment(),
					"order":       k + 1,
					"add_time":    now.Add(-time.Duration(rand.Intn(30)) * 24 * time.Hour),
				})
			}

			booklistID := primitive.NewObjectID()
			isPublic := rand.Intn(10) > 2 // 70% 公开

			viewCount := rand.Intn(500) + 50

			likeCount := rand.Intn(100) + 10
			forkCount := rand.Intn(20)

			booklist := bson.M{
				"_id":         booklistID,
				"user_id":     user["id"],
				"user_name":   user["username"],
				"user_avatar": user["avatar"],
				"title":       template.title,
				"description": template.description,
				"cover":       "",
				"books":       bookItems,
				"book_count":  len(bookItems),
				"like_count":  likeCount,
				"fork_count":  forkCount,
				"view_count":  viewCount,
				"is_public":   isPublic,
				"tags":         []string{template.category, "书单"},
				"category":    template.category,
				"created_at": now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
				"updated_at": now,
			}
			booklists = append(booklists, booklist)
			createdCount++

			// 为部分公开书单创建点赞
			if isPublic && likeCount > 0 {
				for k := 0; k < likeCount && k < len(users); k++ {
					if users[k]["id"] == user["id"] {
						continue // 不能给自己点赞
					}
					likes = append(likes, bson.M{
						"_id":         primitive.NewObjectID(),
						"booklist_id": booklistID.Hex(),
						"user_id":     users[k]["id"],
						"created_at":  now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
					})
				}
			}
		}
	}

	// 批量插入书单
	if len(booklists) > 0 {
		if _, err := booklistCollection.InsertMany(ctx, booklists); err != nil {
			return fmt.Errorf("插入书单失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个书单\n", len(booklists))
	}

	// 批量插入点赞
	if len(likes) > 0 {
		if _, err := likeCollection.InsertMany(ctx, likes); err != nil {
			return fmt.Errorf("插入书单点赞失败: %w", err)
		}
		fmt.Printf("  创建了 %d 个书单点赞\n", len(likes))
	}

	return nil
}

// getRandomComment 获取随机推荐语
func (s *BooklistSeeder) getRandomComment() string {
	comments := []string{
		"强烈推荐！",
		"非常好看，值得一看",
		"这本书太精彩了",
		"五星好评！",
		"入坑不后悔系列",
		"看完忍不住推荐给朋友",
		"熬夜也要看完",
		"已经二刷三遍了",
		"看完心情很好",
		"期待作者更新",
		"必须收藏",
	}
	return comments[rand.Intn(len(comments))]
}

// Clean 清空书单数据
func (s *BooklistSeeder) Clean() error {
	ctx := context.Background()

	_, err := s.db.Collection("booklists").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 booklists 集合失败: %w", err)
	 }

	_, err = s.db.Collection("booklist_likes").DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清空 booklist_likes 集合失败: %w", err)
        }

	fmt.Println("  已清空 booklists 和 booklist_likes 集合")
	return nil
}
