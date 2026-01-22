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

	// 1. 查询所有书籍
	fmt.Println("=== 查询所有书籍 ===")
	booksCollection := db.Collection("books")
	cursor, err := booksCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("查询书籍失败: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		fmt.Printf("解析书籍数据失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 本书:\n", len(books))
	for i, book := range books {
		title := book["title"]
		author := book["author"]
		id := book["_id"]
		fmt.Printf("%d. _id=%v, title=%s, author=%s\n", i+1, id, title, author)
	}

	// 2. 查询所有用户
	fmt.Println("\n=== 查询所有用户 ===")
	usersCollection := db.Collection("users")
	cursor2, err := usersCollection.Find(ctx, bson.M{"username": "testuser"})
	if err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
		return
	}
	defer cursor2.Close(ctx)

	var users []bson.M
	if err = cursor2.All(ctx, &users); err != nil {
		fmt.Printf("解析用户数据失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 个测试用户:\n", len(users))
	for i, user := range users {
		username := user["username"]
		id := user["_id"]
		fmt.Printf("%d. _id=%v, username=%s\n", i+1, id, username)
	}

	// 3. 测试搜索查询
	fmt.Println("\n=== 测试搜索查询 ===")
	var searchResults []bson.M
	searchCursor, err := booksCollection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": "修仙"}},
			{"author": bson.M{"$regex": "修仙"}},
			{"introduction": bson.M{"$regex": "修仙"}},
		},
	})
	if err != nil {
		fmt.Printf("搜索查询失败: %v\n", err)
		return
	}
	defer searchCursor.Close(ctx)

	if err = searchCursor.All(ctx, &searchResults); err != nil {
		fmt.Printf("解析搜索结果失败: %v\n", err)
		return
	}

	fmt.Printf("搜索'修仙'找到 %d 本书\n", len(searchResults))
}
