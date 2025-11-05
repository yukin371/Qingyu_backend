package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"Qingyu_backend/config"
	aiModels "Qingyu_backend/models/ai"
)

func main() {
	fmt.Println("========== AI配额调试工具 ==========\n")

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

	userID := fmt.Sprintf("%v", user["_id"])
	fmt.Printf("1. 用户信息:\n")
	fmt.Printf("   username: %v\n", user["username"])
	fmt.Printf("   _id: %v (原始类型: %T)\n", user["_id"], user["_id"])
	fmt.Printf("   userID字符串: %s\n\n", userID)

	// 模拟QuotaService.GetQuotaByUserID的查询
	quotaType := "daily"
	filter := bson.M{
		"user_id":    userID,
		"quota_type": quotaType,
	}

	fmt.Printf("2. 模拟配额查询:\n")
	fmt.Printf("   查询条件: %+v\n\n", filter)

	var quota aiModels.UserQuota
	err = quotas.FindOne(ctx, filter).Decode(&quota)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("   ❌ 未找到配额记录！\n")

			// 尝试查找所有相关配额
			fmt.Println("3. 查找该用户的所有配额记录:\n")
			cursor, _ := quotas.Find(ctx, bson.M{"user_id": userID})
			defer cursor.Close(ctx)

			count := 0
			for cursor.Next(ctx) {
				var q bson.M
				cursor.Decode(&q)
				count++
				fmt.Printf("   配额 #%d:\n", count)
				fmt.Printf("     user_id: %v\n", q["user_id"])
				fmt.Printf("     quota_type: %v\n", q["quota_type"])
				fmt.Printf("     remaining_quota: %v\n\n", q["remaining_quota"])
			}

			if count == 0 {
				fmt.Println("   ❌ 该用户没有任何配额记录！")
			}
		} else {
			log.Fatal("查询配额失败:", err)
		}
		return
	}

	fmt.Printf("3. 找到配额记录:\n")
	fmt.Printf("   UserID: %s\n", quota.UserID)
	fmt.Printf("   QuotaType: %s\n", quota.QuotaType)
	fmt.Printf("   TotalQuota: %d\n", quota.TotalQuota)
	fmt.Printf("   UsedQuota: %d\n", quota.UsedQuota)
	fmt.Printf("   RemainingQuota: %d\n", quota.RemainingQuota)
	fmt.Printf("   Status: %s\n\n", quota.Status)

	// 测试CanConsume
	testAmount := 100
	canConsume := quota.CanConsume(testAmount)
	fmt.Printf("4. 配额检查:\n")
	fmt.Printf("   请求消费: %d tokens\n", testAmount)
	fmt.Printf("   CanConsume(): %v\n", canConsume)
	fmt.Printf("   IsAvailable(): %v\n", quota.IsAvailable())

	if !canConsume {
		fmt.Println("\n   ❌ 配额不足！这就是为什么返回429错误")
		fmt.Printf("   原因分析:\n")
		if !quota.IsAvailable() {
			fmt.Printf("     - 配额不可用 (Status: %s)\n", quota.Status)
		}
		if quota.RemainingQuota < testAmount {
			fmt.Printf("     - 剩余配额不足 (%d < %d)\n", quota.RemainingQuota, testAmount)
		}
	} else {
		fmt.Println("\n   ✅ 配额充足！应该可以通过检查")
	}

	fmt.Println("\n========================================")
}
