package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Printf("连接MongoDB失败: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")
	booksCollection := db.Collection("books")

	fmt.Println("=== 测试不同的查询条件 ===")

	// 1. 空查询（应该返回所有书籍）
	fmt.Println("\n1. 空查询条件: {}")
	count1, _ := booksCollection.CountDocuments(ctx, bson.M{})
	fmt.Printf("   结果: %d 本书\n", count1)

	// 2. 查询特定status
	fmt.Println("\n2. 查询 status='ongoing'")
	count2, _ := booksCollection.CountDocuments(ctx, bson.M{"status": "ongoing"})
	fmt.Printf("   结果: %d 本书\n", count2)

	// 3. 查询status为nil（测试）
	fmt.Println("\n3. 查询 status=null")
	count3, _ := booksCollection.CountDocuments(ctx, bson.M{"status": nil})
	fmt.Printf("   结果: %d 本书\n", count3)

	// 4. 测试模拟API的filter构建
	fmt.Println("\n4. 模拟API层的filter")
	filter := struct {
		Keyword *string
		Status  *string
		Limit   int
		Offset  int
	}{
		Limit:  20,
		Offset: 0,
	}

	keyword := "修仙"
	filter.Keyword = &keyword
	// 注意：filter.Status 未设置（nil）

	fmt.Printf("   filter.Keyword: %v\n", *filter.Keyword)
	fmt.Printf("   filter.Status: %v\n", filter.Status)
	fmt.Printf("   filter.Limit: %d\n", filter.Limit)
	fmt.Printf("   filter.Offset: %d\n", filter.Offset)

	// 构建MongoDB查询
	query := bson.M{}
	// if filter.Status != nil {  // 这个条件不会执行，因为filter.Status是nil
	//     query["status"] = *filter.Status
	// }
	fmt.Printf("   构建的MongoDB查询: %+v\n", query)

	count4, _ := booksCollection.CountDocuments(ctx, query)
	fmt.Printf("   匹配的书籍数量: %d\n", count4)

	// 5. 测试获取书籍并检查status
	fmt.Println("\n5. 获取前5本书并检查status")
	cursor, _ := booksCollection.Find(ctx, bson.M{}, options.Find().SetLimit(5))
	defer cursor.Close(ctx)

	var results []bson.M
	cursor.All(ctx, &results)
	for i, book := range results {
		fmt.Printf("   书%d: %s, status=%s\n", i+1, book["title"], book["status"])
	}
}
