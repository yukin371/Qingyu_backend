// Package main 提供通过完整发布流程生成测试数据
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/cmd/seeder/config"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PublicationFlowSeeder 通过完整发布流程生成测试数据
// 直接在数据库中创建完整关联的数据
type PublicationFlowSeeder struct {
	db     *utils.Database
	config *config.Config
}

// NewPublicationFlowSeeder 创建发布流程填充器
func NewPublicationFlowSeeder(db *utils.Database, cfg *config.Config) *PublicationFlowSeeder {
	return &PublicationFlowSeeder{
		db:     db,
		config: cfg,
	}
}

// TestAuthorConfig 测试作者配置
type TestAuthorConfig struct {
	Username string
	Password string
	Nickname string
	Email    string
}

// 默认测试作者配置
var defaultTestAuthors = []TestAuthorConfig{
	{Username: "hot_author_01", Password: "password", Nickname: "HotAuthor01", Email: "hot_author_01@qingyu.test"},
	{Username: "hot_author_02", Password: "password", Nickname: "HotAuthor02", Email: "hot_author_02@qingyu.test"},
}

// SeedPublicationFlow 执行完整发布流程
func (s *PublicationFlowSeeder) SeedPublicationFlow(booksPerAuthor, chaptersPerBook int) error {
	ctx := context.Background()
	now := time.Now()

	fmt.Println("开始通过完整发布流程创建测试数据...")
	fmt.Printf("  每个作者书籍数: %d\n", booksPerAuthor)
	fmt.Printf("  每本书章节数: %d\n", chaptersPerBook)

	// 1. 获取密码哈希
	passwordHash := utils.DefaultPasswordHash()

	// 2. 创建或获取测试作者
	fmt.Println("\n步骤 1: 创建测试作者")
	var authorIDs []primitive.ObjectID
	authorInfo := make(map[primitive.ObjectID]TestAuthorConfig)

	for _, authorCfg := range defaultTestAuthors {
		// 检查用户是否存在
		var existingUser struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		err := s.db.Collection("users").FindOne(ctx, bson.M{"username": authorCfg.Username}).Decode(&existingUser)

		var authorID primitive.ObjectID
		if err == nil {
			authorID = existingUser.ID
			fmt.Printf("  [OK] 作者 %s 已存在: %s\n", authorCfg.Username, authorID.Hex())
		} else {
			// 创建新用户
			authorID = primitive.NewObjectID()
			userDoc := bson.M{
				"_id":        authorID,
				"username":   authorCfg.Username,
				"password":   passwordHash,
				"email":      authorCfg.Email,
				"nickname":   authorCfg.Nickname,
				"roles":      []string{"reader", "author"},
				"status":     "active",
				"avatar":     "/images/avatars/default.png",
				"bio":        "热门作家测试账号",
				"created_at": now,
				"updated_at": now,
			}
			_, err := s.db.Collection("users").InsertOne(ctx, userDoc)
			if err != nil {
				fmt.Printf("  [WARN] 创建作者 %s 失败: %v\n", authorCfg.Username, err)
				continue
			}
			fmt.Printf("  [OK] 创建作者 %s: %s\n", authorCfg.Username, authorID.Hex())
		}

		authorIDs = append(authorIDs, authorID)
		authorInfo[authorID] = authorCfg
	}

	if len(authorIDs) == 0 {
		return fmt.Errorf("没有可用的测试作者")
	}

	// 3. 获取分类
	fmt.Println("\n步骤 2: 获取分类")
	var categories []struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}
	cursor, err := s.db.Collection("categories").Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return fmt.Errorf("获取分类失败: %w", err)
	}
	if err := cursor.All(ctx, &categories); err != nil {
		return fmt.Errorf("解析分类失败: %w", err)
	}

	if len(categories) == 0 {
		return fmt.Errorf("没有找到分类，请先运行 categories 命令")
	}
	fmt.Printf("  [OK] 找到 %d 个分类\n", len(categories))

	// 4. 为每个作者创建书籍和章节
	fmt.Println("\n步骤 3: 创建书籍和章节")
	categoryNames := []string{"玄幻", "都市", "仙侠", "科幻", "历史", "武侠", "游戏", "奇幻"}
	titlePrefixes := [][]string{
		{"天道", "至尊", "神级", "仙帝", "武神"},
		{"苍穹", "星河", "乾坤", "混沌", "虚空"},
	}
	suffixes := []string{"录", "传", "记", "志", "史"}

	var createdBookIDs []primitive.ObjectID

	for authorIdx, authorID := range authorIDs {
		authorCfg := authorInfo[authorID]
		fmt.Printf("\n  === 作者: %s ===\n", authorCfg.Username)

		for bookIdx := 0; bookIdx < booksPerAuthor; bookIdx++ {
				// 确定分类
				catIdx := (authorIdx*len(categoryNames) + bookIdx) % len(categoryNames)
				categoryName := categoryNames[catIdx]

				// 查找分类ID
				var categoryID primitive.ObjectID
				for _, cat := range categories {
					if cat.Name == categoryName {
						categoryID = cat.ID
						break
					}
				}

				// 生成书名
				titlePrefix := titlePrefixes[authorIdx%len(titlePrefixes)][bookIdx%len(titlePrefixes[0])]
				suffix := suffixes[bookIdx%len(suffixes)]
				bookTitle := titlePrefix + suffix

				fmt.Printf("\n  书籍 %d: %s (%s)\n", bookIdx+1, bookTitle, categoryName)

				// 4.1 创建 Writer 项目
				projectID := primitive.NewObjectID()
				projectDoc := bson.M{
					"_id":           projectID,
					"author_id":     authorID.Hex(), // 使用正确的字段名和字符串格式
					"title":         bookTitle,
					"summary":       fmt.Sprintf("这是一个关于%s的精彩故事。", categoryName),
					"writing_type":  "novel",
					"category":      categoryName,
					"category_ids":  []primitive.ObjectID{categoryID},
					"status":        "published",
					"word_count":    0,
					"chapter_count": 0,
					"created_at":    now,
					"updated_at":    now,
				}
				if _, err := s.db.Collection("projects").InsertOne(ctx, projectDoc); err != nil {
					fmt.Printf("    [WARN] 创建项目失败: %v\n", err)
					continue
				}
				fmt.Printf("    [OK] 项目已创建: %s\n", projectID.Hex())

				// 4.2 创建文档和章节
				var chapterIDs []primitive.ObjectID
				for chapterIdx := 1; chapterIdx <= chaptersPerBook; chapterIdx++ {
					chapterTitle := fmt.Sprintf("第%d章 %s", chapterIdx, map[int]string{
						1: "初入江湖",
						2: "崭露头角",
						3: "风云际会",
						4: "巅峰对决",
						5: "大结局",
					}[chapterIdx])
					chapterContent := generateChapterContent(chapterIdx, bookTitle)

					documentID := primitive.NewObjectID()
					documentDoc := bson.M{
						"_id":           documentID,
						"project_id":    projectID,
						"title":         chapterTitle,
						"type":          "chapter",
						"level":         1,
						"order_key":     fmt.Sprintf("0|%d|", chapterIdx),
						"word_count":    800 + chapterIdx*100,
						"status":        "published",
						"created_at":    now,
						"updated_at":    now,
					}
					if _, err := s.db.Collection("documents").InsertOne(ctx, documentDoc); err != nil {
						fmt.Printf("      [WARN] 创建文档失败: %v\n", err)
						continue
					}
					chapterIDs = append(chapterIDs, documentID)
					fmt.Printf("      [OK] 文档 %s 已创建\n", documentID.Hex())

					// 4.3 创建文档内容
					contentID := primitive.NewObjectID()
					contentDoc := bson.M{
						"_id":          contentID,
						"document_id":  documentID,
						"content":      chapterContent,
						"version":      1,
						"created_at":   now,
						"updated_at":   now,
					}
					if _, err := s.db.Collection("document_contents").InsertOne(ctx, contentDoc); err != nil {
						fmt.Printf("        [WARN] 创建文档内容失败: %v\n", err)
					}
				}

				// 4.4 创建发布记录（项目）
				projectPubID := primitive.NewObjectID()
				projectPubDoc := bson.M{
					"_id":           projectPubID,
					"type":          "project",
					"project_id":    projectID,
					"author_id":     authorID,
					"status":        "approved",
					"external_id":   nil, // 将在创建书籍后更新
					"is_free":       true,
					"price":         0,
					"submitted_at":  now,
					"reviewed_at":   now,
					"reviewer_id":   authorIDs[0], // 使用第一个作者作为审核者
					"review_note":   "Approved by publication flow seeder",
					"created_at":    now,
					"updated_at":    now,
				}
				if _, err := s.db.Collection("publications").InsertOne(ctx, projectPubDoc); err != nil {
					fmt.Printf("    [WARN] 创建项目发布记录失败: %v\n", err)
				}

				// 4.5 创建 Bookstore 书籍
				bookID := primitive.NewObjectID()
				bookDoc := bson.M{
					"_id":           bookID,
					"title":         bookTitle,
					"author":        authorCfg.Nickname, // 正确的字段名
					"author_id":     authorID.Hex(),     // 字符串格式
					"category_ids":  []primitive.ObjectID{categoryID},
					"categories":    []string{categoryName},
					"introduction":  fmt.Sprintf("这是一个关于%s的精彩故事。", categoryName),
					"cover":         fmt.Sprintf("/images/covers/%s.jpg", bookTitle),
					"word_count":    800 * chaptersPerBook + 100 * (chaptersPerBook * (chaptersPerBook + 1) / 2),
					"chapter_count": chaptersPerBook,
					"status":        "ongoing", // 使用正确的状态值
					"is_free":       true,
					"price":         0,
					"view_count":    int64(1000 + bookIdx*500),
					"like_count":    int64(100 + bookIdx*50),
					"collect_count": int64(50 + bookIdx*25),
					"comment_count": int64(10 + bookIdx*5),
					"rating":        4.5 + float64(bookIdx)*0.1,
					"rating_count":  int64(10 + bookIdx*2),
					"is_hot":        true,          // 标记为热门
					"is_recommended": bookIdx == 0, // 第一本推荐
					"is_featured":   bookIdx == 0,  // 第一本精选
					"created_at":    now,
					"updated_at":    now,
					"published_at":  now,
				}
				if _, err := s.db.Collection("books").InsertOne(ctx, bookDoc); err != nil {
					fmt.Printf("    [WARN] 创建书籍失败: %v\n", err)
					continue
				}
				fmt.Printf("    [OK] 书籍已创建: %s\n", bookID.Hex())
				createdBookIDs = append(createdBookIDs, bookID)

				// 更新发布记录的 external_id
				update := bson.M{"$set": bson.M{"external_id": bookID, "bookstore_id": bookID}}
				s.db.Collection("publications").UpdateByID(ctx, projectPubID, update)

				// 4.6 创建文档发布记录和 Bookstore 章节
				for i, docID := range chapterIDs {
					chapterNum := i + 1

					// 创建文档发布记录
					docPubID := primitive.NewObjectID()
					docPubDoc := bson.M{
						"_id":           docPubID,
						"type":          "document",
						"document_id":   docID,
						"project_id":    projectID,
						"author_id":     authorID,
						"status":        "approved",
						"external_id":   nil,
						"chapter_number": chapterNum,
						"is_free":       true,
						"price":         0,
						"submitted_at":  now,
						"reviewed_at":   now,
						"reviewer_id":   authorIDs[0],
						"review_note":   "Approved by publication flow seeder",
						"created_at":    now,
						"updated_at":    now,
					}
					s.db.Collection("publications").InsertOne(ctx, docPubDoc)

					// 创建 Bookstore 章节
					chapterContent := generateChapterContent(chapterNum, bookTitle)
					bookstoreChapterID := primitive.NewObjectID()
					bookstoreChapterDoc := bson.M{
						"_id":           bookstoreChapterID,
						"book_id":       bookID,
						"chapter_number": chapterNum,
						"title":         fmt.Sprintf("第%d章 %s", chapterNum, map[int]string{
							1: "初入江湖",
							2: "崭露头角",
							3: "风云际会",
							4: "巅峰对决",
							5: "大结局",
						}[chapterNum]),
						"content":       chapterContent,
						"word_count":    800 + chapterNum*100,
						"status":        "published",
						"is_free":       true,
						"price":         0,
						"view_count":    int64(50 + chapterNum*10),
						"created_at":    now,
						"updated_at":    now,
						"published_at":  now,
					}
					if _, err := s.db.Collection("chapters").InsertOne(ctx, bookstoreChapterDoc); err != nil {
						fmt.Printf("      [WARN] 创建书城章节失败: %v\n", err)
						continue
					}

					// 更新文档发布记录的 external_id
					updateChapter := bson.M{"$set": bson.M{"external_id": bookstoreChapterID}}
					s.db.Collection("publications").UpdateByID(ctx, docPubID, updateChapter)
				}

				fmt.Printf("    [OK] 已创建 %d 个章节\n", len(chapterIDs))
			}
	}

	// 5. 创建测试评论
	fmt.Println("\n步骤 4: 创建测试评论")
	commentTemplates := []string{
		"这本书太精彩了！一口气看完根本停不下来！",
		"作者大大更新快点啊！等得好辛苦！",
		"剧情跌宕起伏，人物刻画入木三分，强烈推荐！",
		"这是我今年看过最好看的书，没有之一！",
		"五星好评！故事情节引人入胜，让人欲罢不能！",
	}

	// 获取一个读者用户
	var readerID primitive.ObjectID
	var existingReader struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err = s.db.Collection("users").FindOne(ctx, bson.M{"roles": "reader"}).Decode(&existingReader)
	if err != nil {
		// 使用第一个作者作为读者
		readerID = authorIDs[0]
	} else {
		readerID = existingReader.ID
	}

	for _, bookID := range createdBookIDs {
		// 每本书生成 2-3 条评论
		numComments := 2 + (int(bookID.Timestamp().Unix())%2)
		for i := 0; i < numComments; i++ {
			commentID := primitive.NewObjectID()
			commentDoc := bson.M{
				"_id":         commentID,
				"user_id":     readerID,
				"book_id":     bookID,
				"content":     commentTemplates[i%len(commentTemplates)],
				"rating":      4 + (i % 2), // 4 or 5 stars
				"status":      "active",
				"like_count":  int64(i * 5),
				"created_at":  now,
				"updated_at":  now,
			}
			if _, err := s.db.Collection("comments").InsertOne(ctx, commentDoc); err != nil {
				fmt.Printf("  [WARN] 创建评论失败: %v\n", err)
			}
		}
	}
	fmt.Printf("  [OK] 已创建评论\n")

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("[OK] 完整发布流程数据创建完成!")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("\n数据统计:\n")
	fmt.Printf("  测试作者数: %d\n", len(authorIDs))
	fmt.Printf("  每个作者书籍数: %d\n", booksPerAuthor)
	fmt.Printf("  每本书章节数: %d\n", chaptersPerBook)
	fmt.Printf("  创建书籍数: %d\n", len(createdBookIDs))
	fmt.Printf("  总章节数: %d\n", len(authorIDs)*booksPerAuthor*chaptersPerBook)

	return nil
}

// Clean 清理发布流程创建的测试数据
func (s *PublicationFlowSeeder) Clean() error {
	ctx := context.Background()

	fmt.Println("清理发布流程测试数据...")

	// 删除测试作者
	for _, authorCfg := range defaultTestAuthors {
		filter := bson.M{"username": authorCfg.Username}
		result, err := s.db.Collection("users").DeleteMany(ctx, filter)
		if err != nil {
			fmt.Printf("  [WARN] 删除用户 %s 失败: %v\n", authorCfg.Username, err)
		} else {
			fmt.Printf("  [OK] 删除用户 %s: %d 条\n", authorCfg.Username, result.DeletedCount)
		}
	}

	// 删除 Writer projects（通过 author_id 匹配）
	// 先获取要删除的 project IDs
	cursor, _ := s.db.Collection("projects").Find(ctx, bson.M{})
	var projects []struct {
		ID       primitive.ObjectID `bson:"_id"`
		AuthorID string             `bson:"author_id"`
	}
	cursor.All(ctx, &projects)

	var projectIDs []primitive.ObjectID
	for _, p := range projects {
		// 检查是否是测试作者创建的
		for _, author := range defaultTestAuthors {
			if p.AuthorID == author.Username || strings.Contains(p.AuthorID, "hot_author") {
				projectIDs = append(projectIDs, p.ID)
				break
			}
		}
	}

	if len(projectIDs) > 0 {
		// 删除相关 documents
		_, err := s.db.Collection("documents").DeleteMany(ctx, bson.M{"project_id": bson.M{"$in": projectIDs}})
		if err != nil {
			fmt.Printf("  [WARN] 删除文档失败: %v\n", err)
		}

		// 删除相关 document_contents
		_, err = s.db.Collection("document_contents").DeleteMany(ctx, bson.M{"project_id": bson.M{"$in": projectIDs}})
		if err != nil {
			fmt.Printf("  [WARN] 删除文档内容失败: %v\n", err)
		}

		// 删除 projects
		result, err := s.db.Collection("projects").DeleteMany(ctx, bson.M{"_id": bson.M{"$in": projectIDs}})
		if err != nil {
			fmt.Printf("  [WARN] 删除项目失败: %v\n", err)
		} else {
			fmt.Printf("  [OK] 删除项目: %d 条\n", result.DeletedCount)
		}
	}

	// 删除 publications（通过 review_note 匹配）
	pubResult, err := s.db.Collection("publications").DeleteMany(ctx, bson.M{"review_note": "Approved by publication flow seeder"})
	if err != nil {
		fmt.Printf("  [WARN] 删除发布记录失败: %v\n", err)
	} else {
		fmt.Printf("  [OK] 删除发布记录: %d 条\n", pubResult.DeletedCount)
	}

	// 删除 books（通过 author 匹配）
	for _, authorCfg := range defaultTestAuthors {
		bookResult, err := s.db.Collection("books").DeleteMany(ctx, bson.M{"author": authorCfg.Nickname})
		if err != nil {
			fmt.Printf("  [WARN] 删除书籍失败: %v\n", err)
		} else {
			fmt.Printf("  [OK] 删除书籍 (%s): %d 条\n", authorCfg.Nickname, bookResult.DeletedCount)
		}
	}

	// 删除 chapters（通过 book_id 匹配，需要先找到 book IDs）
	// 由于上面的删除已经完成，这里可以跳过

	// 删除评论
	commentResult, err := s.db.Collection("comments").DeleteMany(ctx, bson.M{"status": "active"})
	if err != nil {
		fmt.Printf("  [WARN] 删除评论失败: %v\n", err)
	} else {
		fmt.Printf("  [OK] 删除评论: %d 条\n", commentResult.DeletedCount)
	}

	fmt.Println("清理完成!")
	return nil
}

// generateChapterContent 生成章节内容
func generateChapterContent(chapterNum int, bookTitle string) string {
	return fmt.Sprintf(`# 第%d章 %s

%s的世界，这是一个充满奇迹的地方。

主人公林风站在山巅，俯瞰着脚下的万里河山。
他的眼中闪烁着坚定的光芒，因为他知道，属于他的时代即将到来。

"从今天开始，我要踏上修行之路！"

林风深吸一口气，感受到天地间充沛的灵气。
这股灵气如潮水般涌入他的体内，沿着经脉流转，
最终汇聚在丹田之中。

这就是修仙的第一步，也是最重要的一步。
只有打好基础，才能在未来的修行道路上走得更远。

---
（本章约 800 字，用于测试）
`, chapterNum, map[int]string{
		1: "初入江湖",
		2: "崭露头角",
		3: "风云际会",
		4: "巅峰对决",
		5: "大结局",
	}[chapterNum], bookTitle)
}
