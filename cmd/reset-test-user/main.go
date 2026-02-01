package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 连接到MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("qingyu")
	collection := db.Collection("users")

	// 生成 password123 的哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("生成哈希失败:", err)
	}

	// 更新 reader1 的密码
	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"username": "reader1"},
		bson.M{"$set": bson.M{"password": string(passwordHash)}},
	)
	if err != nil {
		log.Fatal("更新失败:", err)
	}

	fmt.Printf("更新成功! MatchedCount: %d, ModifiedCount: %d\n", result.MatchedCount, result.ModifiedCount)
	fmt.Printf("新密码哈希: %s\n", string(passwordHash))

	// 测试验证
	var user bson.M
	err = collection.FindOne(context.Background(), bson.M{"username": "reader1"}).Decode(&user)
	if err == nil {
		fmt.Printf("验证后密码: %v\n", user["password"])
	}
}
