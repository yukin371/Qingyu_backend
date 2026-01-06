package main

import (
	"context"
	"fmt"
	"log"

	"Qingyu_backend/config"
	"Qingyu_backend/core"
	"Qingyu_backend/service"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   青羽写作平台 - 更新书籍字段工具       ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// 加载配置
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v\n", err)
	}

	// 初始化服务
	if err := core.InitServices(); err != nil {
		log.Fatalf("❌ 初始化服务失败: %v\n", err)
	}

	ctx := context.Background()
	db := service.GetServiceContainer().GetMongoDB()
	collection := db.Collection("books")

	// 查看需要更新的书籍数量
	total, _ := collection.CountDocuments(ctx, bson.M{})
	fmt.Printf("书籍总数: %d\n\n", total)

	// 更新所有缺失字段的书籍
	fmt.Println("📝 开始更新书籍字段...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	updateResult, err := collection.UpdateMany(
		ctx,
		bson.M{
			"$or": []bson.M{
				{"rating": bson.M{"$exists": false}},
				{"rating_count": bson.M{"$exists": false}},
				{"view_count": bson.M{"$exists": false}},
			},
		},
		bson.M{
			"$set": bson.M{
				"rating":       0.0,
				"rating_count": 0,
				"view_count":   0,
			},
		},
	)

	if err != nil {
		log.Fatalf("❌ 更新失败: %v\n", err)
	}

	fmt.Printf("✓ 匹配到 %d 本书\n", updateResult.MatchedCount)
	fmt.Printf("✓ 成功更新 %d 本书\n", updateResult.ModifiedCount)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 验证更新结果
	fmt.Println("\n🔍 验证更新结果...")
	cursor, err := collection.Find(ctx, bson.M{}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	hasAllFields := 0
	for cursor.Next(ctx) {
		var book bson.M
		cursor.Decode(&book)
		if book["rating"] != nil && book["view_count"] != nil {
			hasAllFields++
		}
	}

	fmt.Printf("✓ 已有完整字段的书籍: %d/%d\n", hasAllFields, total)
	fmt.Println("\n✓ 更新完成！")
}
