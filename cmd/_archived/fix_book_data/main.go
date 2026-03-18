package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	booksColl := db.Collection("books")

	// 查找没有author或categoryIds的书籍
	filter := bson.M{
		"$or": []bson.M{
			{"author": bson.M{"$exists": true, "$in": []interface{}{nil, ""}}},
			{"$or": []bson.M{
				{"category_ids": bson.M{"$exists": false}},
				{"category_ids": bson.M{"$size": 0}},
				{"category_ids": nil},
			}},
		},
	}

	cursor, err := booksColl.Find(ctx, filter)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Printf("找到 %d 本需要修复的书籍\n\n", len(results))

	// 使用"玄幻"分类ID（这是一个有效的ObjectId）
	defaultCategoryID := "696f35c4cee9d6ed15e66933"

	for _, book := range results {
		bookID := book["_id"]
		title := book["title"]

		// 准备更新数据
		update := bson.M{}

		// 如果author为空，设置默认值
		if book["author"] == nil || book["author"] == "" {
			update["author"] = "匿名作者"
		}

		// 如果category_ids为空或无效，设置默认分类
		categoryIDs := book["category_ids"]
		if categoryIDs == nil {
			update["category_ids"] = []string{defaultCategoryID}
		} else if arr, ok := categoryIDs.(primitive.A); ok && len(arr) == 0 {
			update["category_ids"] = []string{defaultCategoryID}
		} else if str, ok := categoryIDs.(string); ok && str == "fiction" {
			// 修复之前错误的 "fiction" 字符串
			update["category_ids"] = []string{defaultCategoryID}
		}

		// 执行更新
		if len(update) > 0 {
			_, err := booksColl.UpdateOne(
				ctx,
				bson.M{"_id": bookID},
				bson.M{"$set": update},
			)
			if err != nil {
				fmt.Printf("❌ 更新失败 [%s]: %v\n", title, err)
			} else {
				fmt.Printf("✓ 已修复 [%s]\n", title)
			}
		}
	}

	fmt.Println("\n修复完成！")
}
