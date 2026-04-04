//go:build ignore
// +build ignore

// Quick script to seed comments into the correct MongoDB database
// Run with: go run seed_comments_fix.go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var commentTemplates = []string{
	"这本书太好看了，作者的文笔太棒了！",
	"剧情很精彩，期待后续发展",
	"人物塑造很到位，很喜欢主角的性格",
	"故事节奏很好，不会让人觉得拖沓",
	"题材新颖，看得出来作者很用心",
	"非常推荐，值得一读再读",
	"结尾意犹未尽，希望能有续作",
	"每个章节都有惊喜，完全停不下来",
	"文字优美，描写细腻，很有画面感",
	"配角也很出彩，整个故事很完整",
	"感人至深，看哭了好几次",
	"悬念迭起，让人欲罢不能",
	"有深度有内涵，不是单纯的爽文",
	"看得出作者下了很大功夫研究",
	"强烈推荐给喜欢这类题材的朋友",
	"满分好评，期待作者更多作品",
	"这本书改变了我的看法，很有启发",
	"逻辑清晰，设定合理，很专业",
	"每次更新都追着看，太精彩了",
	"绝对的神作，不接受反驳",
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	// Get all book IDs
	books, err := db.Collection("books").Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	var bookIDs []primitive.ObjectID
	for books.Next(ctx) {
		var book bson.M
		if err := books.Decode(&book); err != nil {
			continue
		}
		if id, ok := book["_id"].(primitive.ObjectID); ok {
			bookIDs = append(bookIDs, id)
		}
	}
	books.Close(ctx)
	fmt.Printf("Found %d books\n", len(bookIDs))

	// Get all user IDs
	users, err := db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	var userIDs []primitive.ObjectID
	for users.Next(ctx) {
		var user bson.M
		if err := users.Decode(&user); err != nil {
			continue
		}
		if id, ok := user["_id"].(primitive.ObjectID); ok {
			userIDs = append(userIDs, id)
		}
	}
	users.Close(ctx)
	fmt.Printf("Found %d users\n", len(userIDs))

	if len(bookIDs) == 0 || len(userIDs) == 0 {
		fmt.Println("No books or users found, skipping comment seeding")
		return
	}

	// Generate comments
	var comments []interface{}
	rand.Seed(time.Now().UnixNano())

	for _, bookID := range bookIDs {
		numComments := 3 + rand.Intn(8) // 3-10 comments per book
		for i := 0; i < numComments; i++ {
			userID := userIDs[rand.Intn(len(userIDs))]
			rating := 3 + rand.Intn(3) // 3-5 stars

			comment := bson.M{
				"_id":            primitive.NewObjectID(),
				"target_id":      bookID,
				"target_type":   "book",
				"author_id":     userID,
				"content":       commentTemplates[rand.Intn(len(commentTemplates))],
				"rating":        rating,
				"like_count":    int64(rand.Intn(50)),
				"reply_count":   int64(rand.Intn(10)),
				"state":         "normal",
				"is_pinned":     false,
				"is_featured":   false,
				"is_author_reply": false,
				"parent_id":     nil,
				"created_at":     time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour),
				"updated_at":    time.Now(),
			}
			comments = append(comments, comment)
		}
	}

	fmt.Printf("Generated %d comments\n", len(comments))

	// Insert comments
	result, err := db.Collection("comments").InsertMany(ctx, comments)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted %d comments\n", len(result.InsertedIDs))

	// Update book ratings
	bookRatings := make(map[primitive.ObjectID][]float64)
	for _, c := range comments {
		cm := c.(bson.M)
		if bookID, ok := cm["target_id"].(primitive.ObjectID); ok {
			if rating, ok := cm["rating"].(int); ok {
				bookRatings[bookID] = append(bookRatings[bookID], float64(rating))
			}
		}
	}

	for bookID, ratings := range bookRatings {
		var sum float64
		for _, r := range ratings {
			sum += r
		}
		avg := sum / float64(len(ratings))
		_, err := db.Collection("books").UpdateOne(ctx, bson.M{"_id": bookID}, bson.M{"$set": bson.M{
			"rating":        avg,
			"rating_count":  len(ratings),
		}})
		if err != nil {
			fmt.Printf("Failed to update book %s: %v\n", bookID.Hex(), err)
		}
	}

	fmt.Println("Done!")
}
