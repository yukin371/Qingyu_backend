package seedrealtimeranking

import (
	"context"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type optionsValue struct{ apply bool }

type book struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	ViewCount int64              `bson:"view_count"`
}

// Run previews or explicitly regenerates realtime ranking data.
func Run(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if opts.apply && os.Getenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA") != "true" {
		return fmt.Errorf("写入排名数据需要设置 QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbconfig.URI()))
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping 失败: %w", err)
	}
	db := client.Database(mongodbconfig.Database())
	cursor, err := db.Collection("books").Find(ctx, bson.M{"status": bson.M{"$in": []string{"ongoing", "completed"}}}, options.Find().SetSort(bson.D{{Key: "view_count", Value: -1}}).SetLimit(20))
	if err != nil {
		return fmt.Errorf("查询书籍失败: %w", err)
	}
	defer cursor.Close(ctx)
	var books []book
	if err := cursor.All(ctx, &books); err != nil {
		return fmt.Errorf("解析书籍失败: %w", err)
	}
	if len(books) == 0 {
		return fmt.Errorf("没有找到书籍数据")
	}

	rankingColl := db.Collection("ranking_items")
	oldCount, err := rankingColl.CountDocuments(ctx, bson.M{"type": "realtime"})
	if err != nil {
		return fmt.Errorf("统计旧排名失败: %w", err)
	}
	period := time.Now().Format("2006-01-02")
	docs := make([]interface{}, 0, len(books))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, currentBook := range books {
		surgeScore := math.Round((100.0-float64(i)*4.5+r.Float64()*3-1.5)*100) / 100
		if surgeScore < 0 {
			surgeScore = r.Float64() * 5
		}
		viewCount := int64(float64(currentBook.ViewCount) * (0.15 - float64(i)*0.006) * (0.9 + r.Float64()*0.2))
		likeCount := int64(float64(viewCount) * (0.08 - float64(i)*0.003) * (0.8 + r.Float64()*0.4))
		now := time.Now()
		docs = append(docs, bson.M{"book_id": currentBook.ID, "type": "realtime", "rank": i + 1, "score": surgeScore, "view_count": viewCount, "like_count": likeCount, "period": period, "created_at": now, "updated_at": now})
		fmt.Printf("  #%2d  %-25s  surge=%.1f  views=%d  likes=%d\n", i+1, currentBook.Title, surgeScore, viewCount, likeCount)
	}
	if !opts.apply {
		fmt.Printf("预览模式：将删除 %d 条旧 realtime 数据并插入 %d 条新数据；未写入。\n", oldCount, len(docs))
		return nil
	}
	delResult, err := rankingColl.DeleteMany(ctx, bson.M{"type": "realtime"})
	if err != nil {
		return fmt.Errorf("清理旧排名失败: %w", err)
	}
	fmt.Printf("清除旧数据: %d 条\n", delResult.DeletedCount)
	result, err := rankingColl.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("插入新排名失败: %w", err)
	}
	fmt.Printf("\n成功插入 %d 条飙升榜数据\n", len(result.InsertedIDs))
	return nil
}

func parseOptions(args []string) (optionsValue, error) {
	fs := flag.NewFlagSet("seed_realtime_ranking", flag.ContinueOnError)
	apply := fs.Bool("apply", false, "执行删除旧数据并插入新数据")
	if err := fs.Parse(args); err != nil {
		return optionsValue{}, err
	}
	return optionsValue{apply: *apply}, nil
}
