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

type Book struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Title    string `bson:"title" json:"title"`
	Author   string `bson:"author" json:"author"`
	Status   string `bson:"status" json:"status"`
	ViewCount int64  `bson:"view_count" json:"viewCount"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	booksCollection := db.Collection("books")
	testID := "6980b950550350d562e6a231"

	fmt.Println("========== 测试 FindOne vs Find ==========")

	// 测试1: FindOne 使用 string _id
	fmt.Println("\n1. FindOne 查询 (string _id):")
	var book1 Book
	err = booksCollection.FindOne(ctx, bson.M{"_id": testID}).Decode(&book1)
	if err != nil {
		fmt.Printf("   失败: %v\n", err)
	} else {
		fmt.Printf("   成功! 书名: %s, 状态: %s\n", book1.Title, book1.Status)
	}

	// 测试2: FindOne 使用 ObjectID
	fmt.Println("\n2. FindOne 查询 (ObjectID _id):")
	objectID, _ := primitive.ObjectIDFromHex(testID)
	var book2 Book
	err = booksCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book2)
	if err != nil {
		fmt.Printf("   失败: %v\n", err)
	} else {
		fmt.Printf("   成功! 书名: %s, 状态: %s\n", book2.Title, book2.Status)
	}

	// 测试3: Find 查询并过滤
	fmt.Println("\n3. Find 查询 (status 过滤):")
	cursor, err := booksCollection.Find(ctx, bson.M{
		"status": bson.M{"$in": []string{"ongoing", "completed"}},
	}, options.Find().SetLimit(3))
	if err != nil {
		fmt.Printf("   失败: %v\n", err)
	} else {
		var books []Book
		if err = cursor.All(ctx, &books); err != nil {
			fmt.Printf("   解码失败: %v\n", err)
		} else {
			fmt.Printf("   成功! 找到 %d 本书\n", len(books))
			for i, b := range books {
				fmt.Printf("   [%d] ID: %s, 书名: %s\n", i+1, b.ID, b.Title)
			}
		}
		cursor.Close(ctx)
	}

	// 测试4: Find 查询特定 ID (使用字符串)
	fmt.Println("\n4. Find 查询特定 ID (string):")
	cursor2, err := booksCollection.Find(ctx, bson.M{"_id": testID}, options.Find().SetLimit(1))
	if err != nil {
		fmt.Printf("   失败: %v\n", err)
	} else {
		var books []Book
		if err = cursor2.All(ctx, &books); err != nil {
			fmt.Printf("   解码失败: %v\n", err)
		} else if len(books) > 0 {
			fmt.Printf("   成功! 书名: %s, 状态: %s\n", books[0].Title, books[0].Status)
		} else {
			fmt.Printf("   没有找到结果\n")
		}
		cursor2.Close(ctx)
	}

	// 测试5: Find 查询特定 ID (使用 ObjectID)
	fmt.Println("\n5. Find 查询特定 ID (ObjectID):")
	cursor3, err := booksCollection.Find(ctx, bson.M{"_id": objectID}, options.Find().SetLimit(1))
	if err != nil {
		fmt.Printf("   失败: %v\n", err)
	} else {
		var books []Book
		if err = cursor3.All(ctx, &books); err != nil {
			fmt.Printf("   解码失败: %v\n", err)
		} else if len(books) > 0 {
			fmt.Printf("   成功! 书名: %s, 状态: %s\n", books[0].Title, books[0].Status)
		} else {
			fmt.Printf("   没有找到结果\n")
		}
		cursor3.Close(ctx)
	}
}
