//go:build ignore

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	// 获取测试书籍ID
	bookID, _ := primitive.ObjectIDFromHex("696f35c4cee9d6ed15e66935")
	authorID, _ := primitive.ObjectIDFromHex("6768bf1fe3fada58740d95ba")

	// 更新书籍的authorId
	result, err := db.Collection("books").UpdateOne(
		ctx,
		bson.M{"_id": bookID},
		bson.M{"$set": bson.M{"author_id": authorID}},
	)
	if err != nil {
		log.Printf("更新失败: %v", err)
		return
	}

	fmt.Printf("✓ 成功更新书籍authorId，匹配记录数: %d\n", result.MatchedCount)
}
