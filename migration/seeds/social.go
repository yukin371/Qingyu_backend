package seeds

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Comment 评论
type Comment struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty"`
	TargetID   string              `bson:"target_id"`   // 书籍ID或章节ID
	TargetType string              `bson:"target_type"` // book, chapter
	UserID     string              `bson:"user_id"`
	Content    string              `bson:"content"`
	Rating     int                 `bson:"rating"` // 1-5星评分
	LikeCount  int                 `bson:"like_count"`
	ReplyCount int                 `bson:"reply_count"`
	ParentID   *primitive.ObjectID `bson:"parent_id,omitempty"` // 父评论ID（回复时）
	Status     string              `bson:"status"`              // normal, hidden, deleted
	CreatedAt  time.Time           `bson:"created_at"`
	UpdatedAt  time.Time           `bson:"updated_at"`
}

// Like 点赞
type Like struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	TargetID   string             `bson:"target_id"`
	TargetType string             `bson:"target_type"` // book, chapter, comment
	CreatedAt  time.Time          `bson:"created_at"`
}

// Collection 收藏
type Collection struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	BookID     string             `bson:"book_id"`
	FolderName string             `bson:"folder_name"`
	Note       string             `bson:"note"`
	IsPublic   bool               `bson:"is_public"`
	CreatedAt  time.Time          `bson:"created_at"`
}

// Follow 关注
type Follow struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FollowerID  string             `bson:"follower_id"`  // 关注者ID
	FollowingID string             `bson:"following_id"` // 被关注者ID
	CreatedAt   time.Time          `bson:"created_at"`
}

// SeedSocialData 社交数据种子
func SeedSocialData(ctx context.Context, db *mongo.Database) error {
	fmt.Println("========================================")
	fmt.Println("开始创建社交测试数据...")
	fmt.Println("========================================")

	// 获取用户和书籍
	userCollection := db.Collection("users")
	bookCollection := db.Collection("books")

	users, err := getUserIDs(ctx, userCollection)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	books, err := getBookIDs(ctx, bookCollection)
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}

	if len(users) == 0 || len(books) == 0 {
		fmt.Println("警告：没有足够的用户或书籍，跳过社交数据创建")
		return nil
	}

	// 创建评论
	err = seedComments(ctx, db, users, books)
	if err != nil {
		return fmt.Errorf("创建评论失败: %w", err)
	}

	// 创建点赞
	err = seedLikes(ctx, db, users, books)
	if err != nil {
		return fmt.Errorf("创建点赞失败: %w", err)
	}

	// 创建收藏
	err = seedCollections(ctx, db, users, books)
	if err != nil {
		return fmt.Errorf("创建收藏失败: %w", err)
	}

	// 创建关注
	err = seedFollows(ctx, db, users)
	if err != nil {
		return fmt.Errorf("创建关注失败: %w", err)
	}

	fmt.Println("========================================")
	fmt.Println("社交数据创建完成")
	fmt.Println("========================================")

	return nil
}

func seedComments(ctx context.Context, db *mongo.Database, users, books []string) error {
	collection := db.Collection("comments")

	// 检查是否已有评论
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		fmt.Printf("评论数据已存在 (%d条)，跳过创建\n", count)
		return nil
	}

	commentTemplates := []string{
		"这本书太好看了！强烈推荐！",
		"作者大大更新快点啊，等不及了",
		"情节非常吸引人，一口气读完了",
		"人物塑造很饱满，对话很精彩",
		"结局有点仓促，不过整体还是不错的",
		"期待第二部！",
		"这个伏笔埋得好，后面肯定有反转",
		"主角的性格我很喜欢，很真实",
		"文笔流畅，故事紧凑",
		"已加入收藏，慢慢看",
		"这本书简直是神作！",
		"有点虐心，不过很感人",
		"配角也很出彩",
		"世界观设定很宏大",
		"希望作者加油，会更的",
	}

	now := time.Now()
	commentCount := 0

	for _, bookID := range books {
		// 每本书随机5-20条评论
		numComments := 5 + rand.Intn(16)

		for i := 0; i < numComments; i++ {
			userID := users[rand.Intn(len(users))]
			content := commentTemplates[rand.Intn(len(commentTemplates))]

			comment := Comment{
				ID:         primitive.NewObjectID(),
				TargetID:   bookID,
				TargetType: "book",
				UserID:     userID,
				Content:    content,
				Rating:     3 + rand.Intn(3), // 3-5星
				LikeCount:  rand.Intn(50),
				ReplyCount: rand.Intn(5),
				Status:     "normal",
				CreatedAt:  now.Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
				UpdatedAt:  now,
			}

			_, err := collection.InsertOne(ctx, comment)
			if err == nil {
				commentCount++
			}

			// 随机添加回复
			if rand.Float32() < 0.3 { // 30%概率有回复
				replyCount := 1 + rand.Intn(3)
				for j := 0; j < replyCount; j++ {
					replyUserID := users[rand.Intn(len(users))]
					reply := Comment{
						ID:         primitive.NewObjectID(),
						TargetID:   bookID,
						TargetType: "book",
						UserID:     replyUserID,
						Content:    randomReply(),
						ParentID:   &comment.ID,
						Status:     "normal",
						CreatedAt:  comment.CreatedAt.Add(time.Duration(j+1) * time.Hour),
						UpdatedAt:  now,
					}
					if _, err := collection.InsertOne(ctx, reply); err != nil {
						log.Printf("Warning: failed to insert reply: %v", err)
					}
				}
			}
		}
	}

	fmt.Printf("✓ 创建了 %d 条评论\n", commentCount)
	return nil
}

func seedLikes(ctx context.Context, db *mongo.Database, users, books []string) error {
	collection := db.Collection("likes")

	// 检查是否已有点赞
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		fmt.Printf("点赞数据已存在 (%d条)，跳过创建\n", count)
		return nil
	}

	now := time.Now()
	likeCount := 0

	for _, bookID := range books {
		// 每本书随机20-80个点赞
		numLikes := 20 + rand.Intn(61)

		for i := 0; i < numLikes; i++ {
			userID := users[rand.Intn(len(users))]

			like := Like{
				ID:         primitive.NewObjectID(),
				UserID:     userID,
				TargetID:   bookID,
				TargetType: "book",
				CreatedAt:  now.Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
			}

			_, err := collection.InsertOne(ctx, like)
			if err == nil {
				likeCount++
			}
		}
	}

	fmt.Printf("✓ 创建了 %d 条点赞记录\n", likeCount)
	return nil
}

func seedCollections(ctx context.Context, db *mongo.Database, users, books []string) error {
	collection := db.Collection("collections")

	// 检查是否已有收藏
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		fmt.Printf("收藏数据已存在 (%d条)，跳过创建\n", count)
		return nil
	}

	folderNames := []string{
		"我的收藏",
		"玄幻小说",
		"仙侠修真",
		"都市言情",
		"科幻奇幻",
		"历史军事",
		"必读榜单",
		"精品推荐",
		"睡前读物",
		"追更中",
	}

	now := time.Now()
	collectionCount := 0

	for _, userID := range users {
		// 每个用户随机5-30个收藏
		numCollections := 5 + rand.Intn(26)

		for i := 0; i < numCollections; i++ {
			bookID := books[rand.Intn(len(books))]

			collectionItem := Collection{
				ID:         primitive.NewObjectID(),
				UserID:     userID,
				BookID:     bookID,
				FolderName: folderNames[rand.Intn(len(folderNames))],
				Note:       randomNote(),
				IsPublic:   rand.Float32() < 0.2, // 20%公开
				CreatedAt:  now.Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
			}

			_, err := collection.InsertOne(ctx, collectionItem)
			if err == nil {
				collectionCount++
			}
		}
	}

	fmt.Printf("✓ 创建了 %d 条收藏记录\n", collectionCount)
	return nil
}

func seedFollows(ctx context.Context, db *mongo.Database, users []string) error {
	collection := db.Collection("follows")

	// 检查是否已有关注
	count, _ := collection.CountDocuments(ctx, bson.M{})
	if count > 0 {
		fmt.Printf("关注数据已存在 (%d条)，跳过创建\n", count)
		return nil
	}

	now := time.Now()
	followCount := 0

	// 为作者创建粉丝
	for i, authorID := range users {
		// 每个用户随机关注5-20个其他用户
		numFollows := 5 + rand.Intn(16)

		for j := 0; j < numFollows; j++ {
			// 随机选择一个不同的用户来关注
			targetIdx := (i + j + 1) % len(users)
			if targetIdx == i {
				continue
			}
			targetID := users[targetIdx]

			follow := Follow{
				ID:          primitive.NewObjectID(),
				FollowerID:  authorID,
				FollowingID: targetID,
				CreatedAt:   now.Add(-time.Duration(rand.Intn(90)) * 24 * time.Hour),
			}

			_, err := collection.InsertOne(ctx, follow)
			if err == nil {
				followCount++
			}
		}
	}

	fmt.Printf("✓ 创建了 %d 条关注记录\n", followCount)
	return nil
}

func getUserIDs(ctx context.Context, collection *mongo.Collection) ([]string, error) {
	cursor, err := collection.Find(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	userIDs := make([]string, len(results))
	for i, r := range results {
		userIDs[i] = r.ID
	}
	return userIDs, nil
}

func getBookIDs(ctx context.Context, collection *mongo.Collection) ([]string, error) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID string `bson:"_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	bookIDs := make([]string, len(results))
	for i, r := range results {
		bookIDs[i] = r.ID
	}
	return bookIDs, nil
}

func randomReply() string {
	replies := []string{
		"同意！",
		"说得对",
		"我也是这么觉得的",
		"哈哈，确实",
		"同感",
		"赞同",
	}
	return replies[rand.Intn(len(replies))]
}

func randomNote() string {
	notes := []string{
		"",
		"很好看",
		"值得推荐",
		"经典之作",
		"必须追",
		"收藏备用",
		"周末阅读",
	}
	return notes[rand.Intn(len(notes))]
}
