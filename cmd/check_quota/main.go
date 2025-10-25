package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/config"
)

func main() {
	fmt.Println("检查AI配额数据...")
	
	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 连接MongoDB
	ctx := context.Background()
	mongoConfig := config.GlobalConfig.Database.Primary.MongoDB
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URI))
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(mongoConfig.Database)
	users := db.Collection("users")
	quotas := db.Collection("ai_quotas")

	// 获取vip_user01
	var user bson.M
	err = users.FindOne(ctx, bson.M{"username": "vip_user01"}).Decode(&user)
	if err != nil {
		log.Fatal("查找用户失败:", err)
	}

	fmt.Printf("\n用户信息:\n")
	fmt.Printf("  username: %v\n", user["username"])
	fmt.Printf("  _id: %v (type: %T)\n", user["_id"], user["_id"])
	
	userID := user["_id"]

	// 查询该用户的所有配额
	cursor, err := quotas.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Fatal("查询配额失败:", err)
	}
	defer cursor.Close(ctx)

	fmt.Printf("\n该用户的配额记录:\n")
	count := 0
	for cursor.Next(ctx) {
		var quota bson.M
		cursor.Decode(&quota)
		count++
		fmt.Printf("\n配额 #%d:\n", count)
		fmt.Printf("  user_id: %v (type: %T)\n", quota["user_id"], quota["user_id"])
		fmt.Printf("  quota_type: %v\n", quota["quota_type"])
		fmt.Printf("  total_quota: %v\n", quota["total_quota"])
		fmt.Printf("  used_quota: %v\n", quota["used_quota"])
		fmt.Printf("  remaining_quota: %v\n", quota["remaining_quota"])
		fmt.Printf("  status: %v\n", quota["status"])
	}

	if count == 0 {
		fmt.Println("  ❌ 未找到配额记录！")
	}
	fmt.Println()
}

