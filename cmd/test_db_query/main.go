package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	// 测试不同的数据库名称
	databases := []string{"Qingyu_backend", "qingyu", "qingyu_read"}
	testID := "6980b950550350d562e6a231"

	for _, dbName := range databases {
		fmt.Printf("\n========== 测试数据库: %s ==========\n", dbName)
		db := client.Database(dbName)

		// 列出所有集合
		collections, err := db.ListCollectionNames(ctx, bson.M{})
		if err != nil {
			fmt.Printf("获取集合列表失败: %v\n", err)
			continue
		}
		fmt.Printf("集合列表: %v\n", collections)

		// 检查 books 集合
		if !contains(collections, "books") {
			fmt.Println("books 集合不存在，跳过")
			continue
		}

		booksCollection := db.Collection("books")

		// 统计文档数量
		total, err := booksCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Printf("统计文档失败: %v\n", err)
			continue
		}
		fmt.Printf("books 集合文档总数: %d\n", total)

		if total == 0 {
			fmt.Println("集合为空，跳过")
			continue
		}

		// 尝试通过 ID 查询
		objectID, err := primitive.ObjectIDFromHex(testID)
		if err != nil {
			fmt.Printf("解析 ObjectID 失败: %v\n", err)
			continue
		}

		var book bson.M
		err = booksCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
		if err != nil {
			fmt.Printf("通过 ID 查询失败: %v\n", err)
		} else {
			fmt.Printf("✓ 通过 ID 查询成功: %v\n", book["title"])
		}

		// 查询前 5 个文档
		cursor, err := booksCollection.Find(ctx, bson.M{}, options.Find().SetLimit(5))
		if err != nil {
			fmt.Printf("查询前 5 个文档失败: %v\n", err)
			continue
		}
		defer cursor.Close(ctx)

		// 使用 Next 方法逐个读取文档
		fmt.Printf("前 5 个文档:\n")
		count := 0
		for cursor.Next(ctx) {
			var book bson.M
			if err := cursor.Decode(&book); err != nil {
				fmt.Printf("  解码文档失败: %v\n", err)
				continue
			}
			count++

			// 检查 _id 字段的类型和值
			idField := book["_id"]
			fmt.Printf("  [%d] _id 类型: %T, 值: %v, 标题: %v\n", count, idField, idField, book["title"])

			if idObj, ok := idField.(primitive.ObjectID); ok {
				fmt.Printf("      这是 ObjectID, Hex: %s\n", idObj.Hex())
			}

			if count >= 5 {
				break
			}
		}

		if err := cursor.Err(); err != nil {
			fmt.Printf("游标错误: %v\n", err)
		}

		if count == 0 {
			fmt.Println("  没有找到任何文档！")

			// 尝试使用 Aggregate 来检查
			pipeline := mongo.Pipeline{
				bson.D{{"$limit", 5}},
			}
			cursor2, err := booksCollection.Aggregate(ctx, pipeline)
			if err != nil {
				fmt.Printf("  Aggregate 查询失败: %v\n", err)
			} else {
				defer cursor2.Close(ctx)
				var results []bson.M
				if err = cursor2.All(ctx, &results); err != nil {
					fmt.Printf("  Aggregate 解码失败: %v\n", err)
				} else {
					fmt.Printf("  Aggregate 查询返回 %d 个文档\n", len(results))
				}
			}
		}

		// 测试：使用字符串 ID 查询
		fmt.Printf("\n  测试使用不同方式查询 ID %s:\n", testID)

		// 方式1: ObjectID (已定义)
		var book1 bson.M
		err1 := booksCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book1)
		fmt.Printf("    使用 ObjectID: %v\n", err1)

		// 方式2: 字符串
		var book2 bson.M
		err2 := booksCollection.FindOne(ctx, bson.M{"_id": testID}).Decode(&book2)
		fmt.Printf("    使用字符串: %v\n", err2)

		// 方式3: $in 操作符
		var book3 bson.M
		err3 := booksCollection.FindOne(ctx, bson.M{"_id": bson.M{"$in": []interface{}{objectID, testID}}}).Decode(&book3)
		fmt.Printf("    使用 $in: %v\n", err3)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
