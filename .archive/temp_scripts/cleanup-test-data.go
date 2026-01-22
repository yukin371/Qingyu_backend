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

	fmt.Println("清理旧的测试数据...")

	// 删除测试书籍
	_, err = db.Collection("books").DeleteMany(ctx, bson.M{"title": bson.M{"$regex": "修仙"}})
	if err != nil {
		fmt.Printf("删除测试书籍失败: %v\n", err)
	} else {
		fmt.Println("✓ 已删除测试书籍")
	}

	// 删除测试分类
	_, err = db.Collection("categories").DeleteMany(ctx, bson.M{"name": bson.M{"$in": bson.A{"玄幻", "修仙"}}})
	if err != nil {
		fmt.Printf("删除测试分类失败: %v\n", err)
	} else {
		fmt.Println("✓ 已删除测试分类")
	}

	// 删除测试用户
	_, err = db.Collection("users").DeleteMany(ctx, bson.M{"username": "testuser"})
	if err != nil {
		fmt.Printf("删除测试用户失败: %v\n", err)
	} else {
		fmt.Println("✓ 已删除测试用户")
	}

	fmt.Println("✓ 旧测试数据清理完成!")
}
