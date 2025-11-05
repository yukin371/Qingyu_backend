package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"Qingyu_backend/config"
	"Qingyu_backend/global"
	"Qingyu_backend/models/writer"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("========== 创建测试项目数据 ==========\n")

	// 1. 加载配置
	fmt.Println("1. 加载配置...")
	_, err := config.LoadConfig("./config/config.test.yaml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}
	fmt.Println("✓ 配置加载成功")

	// 2. 连接MongoDB
	fmt.Println("\n2. 连接MongoDB...")
	ctx := context.Background()

	// 获取MongoDB配置
	mongoURI := config.GlobalConfig.Database.Primary.MongoDB.URI
	dbName := config.GlobalConfig.Database.Primary.MongoDB.Database

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("连接MongoDB失败:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	global.DB = db
	fmt.Println("✓ MongoDB连接成功")

	// 3. 获取用户ID
	fmt.Println("\n3. 查找测试用户...")
	users := db.Collection("users")

	testUsers := []string{"test_user01", "vip_user01"}
	userIDs := make(map[string]string)

	for _, username := range testUsers {
		var user bson.M
		err := users.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			fmt.Printf("  ! 用户 %s 不存在\n", username)
			continue
		}

		var userID string
		switch v := user["_id"].(type) {
		case string:
			userID = v
		default:
			fmt.Printf("  ! 用户 %s 的ID类型不支持\n", username)
			continue
		}

		userIDs[username] = userID
		fmt.Printf("  ✓ %s: %s\n", username, userID)
	}

	if len(userIDs) == 0 {
		log.Fatal("没有找到测试用户")
	}

	// 4. 创建测试项目
	fmt.Println("\n4. 创建测试项目...")
	projects := db.Collection("projects")

	now := time.Now()

	// 为每个用户创建一个测试项目
	for username, userID := range userIDs {
		projectID := "test_project_001_" + username

		// 检查是否已存在
		var existing bson.M
		err := projects.FindOne(ctx, bson.M{"_id": projectID}).Decode(&existing)
		if err == nil {
			fmt.Printf("  项目 %s 已存在，跳过\n", projectID)
			continue
		}

		project := writer.Project{
			ID:         projectID,
			AuthorID:   userID,
			Title:      "测试项目",
			Summary:    "AI功能测试专用项目",
			Category:   "fantasy",
			Status:     "draft",
			Visibility: "private",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// 插入项目
		_, err = projects.InsertOne(ctx, project)
		if err != nil {
			fmt.Printf("  ! 创建项目 %s 失败: %v\n", projectID, err)
			continue
		}

		fmt.Printf("  ✓ 创建项目: %s (用户: %s)\n", projectID, username)
	}

	// 5. 创建通用测试项目（projectId固定为 test_project_001）
	fmt.Println("\n5. 创建通用测试项目...")

	// 使用test_user01的ID
	defaultUserID := userIDs["test_user01"]
	if defaultUserID == "" {
		defaultUserID = userIDs["vip_user01"]
	}

	var existing bson.M
	err = projects.FindOne(ctx, bson.M{"_id": "test_project_001"}).Decode(&existing)
	if err == nil {
		fmt.Println("  通用测试项目已存在")
	} else {
		project := writer.Project{
			ID:         "test_project_001",
			AuthorID:   defaultUserID,
			Title:      "通用测试项目",
			Summary:    "AI功能测试专用项目（通用）",
			Category:   "fantasy",
			Status:     "draft",
			Visibility: "private",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		_, err = projects.InsertOne(ctx, project)
		if err != nil {
			fmt.Printf("  ! 创建通用测试项目失败: %v\n", err)
		} else {
			fmt.Println("  ✓ 创建通用测试项目: test_project_001")
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("✓ 测试项目数据创建完成！")
	fmt.Println("========================================")
}
