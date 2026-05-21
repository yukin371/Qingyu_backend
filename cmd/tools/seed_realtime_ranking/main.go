package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	// 获取书籍列表，按 viewCount 降序
	cursor, err := db.Collection("books").Find(ctx, bson.M{
		"status": bson.M{"$in": []string{"ongoing", "completed"}},
	}, options.Find().SetSort(bson.D{{"view_count", -1}}).SetLimit(20))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	type Book struct {
		ID        primitive.ObjectID `bson:"_id"`
		Title     string             `bson:"title"`
		ViewCount int64              `bson:"view_count"`
	}
	var books []Book
	if err := cursor.All(ctx, &books); err != nil {
		log.Fatal(err)
	}
	if len(books) == 0 {
		log.Fatal("没有找到书籍数据")
	}

	// 清除旧 realtime 数据（集合名是 ranking_items）
	rankingColl := db.Collection("ranking_items")
	delResult, _ := rankingColl.DeleteMany(ctx, bson.M{"type": "realtime"})
	fmt.Printf("清除旧数据: %d 条\n", delResult.DeletedCount)

	// 生成飙升榜数据
	period := time.Now().Format("2006-01-02")
	var docs []interface{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i, book := range books {
		// 飙升指数：排名越前，飙升越高
		surgeScore := math.Round((100.0 - float64(i)*4.5 + r.Float64()*3 - 1.5) * 100) / 100
		if surgeScore < 0 {
			surgeScore = r.Float64() * 5
		}

		viewCount := int64(float64(book.ViewCount) * (0.15 - float64(i)*0.006) * (0.9 + r.Float64()*0.2))
		likeCount := int64(float64(viewCount) * (0.08 - float64(i)*0.003) * (0.8 + r.Float64()*0.4))

		now := time.Now()
		docs = append(docs, bson.M{
			"book_id":    book.ID,
			"type":       "realtime",
			"rank":       i + 1,
			"score":      surgeScore,
			"view_count": viewCount,
			"like_count": likeCount,
			"period":     period,
			"created_at": now,
			"updated_at": now,
		})

		fmt.Printf("  #%2d  %-25s  surge=%.1f  views=%d  likes=%d\n",
			i+1, book.Title, surgeScore, viewCount, likeCount)
	}

	result, err := rankingColl.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n成功插入 %d 条飙升榜数据\n", len(result.InsertedIDs))
}
