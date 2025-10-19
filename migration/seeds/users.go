package seeds

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers 用户种子数据
func SeedUsers(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// 检查是否已有数据
	count, err := collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		fmt.Printf("Users collection already has %d documents, skipping seed\n", count)
		return nil
	}

	// 准备测试用户
	testUsers := []users.User{
		{
			Username:  "admin",
			Email:     "admin@qingyu.com",
			Phone:     "13800138000",
			Role:      "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "author1",
			Email:     "author1@qingyu.com",
			Phone:     "13800138001",
			Role:      "author",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "reader1",
			Email:     "reader1@qingyu.com",
			Phone:     "13800138002",
			Role:      "reader",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:  "reader2",
			Email:     "reader2@qingyu.com",
			Phone:     "13800138003",
			Role:      "reader",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 设置密码（所有测试用户密码都是 "password123"）
	for i := range testUsers {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		testUsers[i].Password = string(hashedPassword)
	}

	// 插入用户
	docs := make([]interface{}, len(testUsers))
	for i, user := range testUsers {
		docs[i] = user
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}

	fmt.Printf("✓ Seeded %d users\n", len(result.InsertedIDs))
	fmt.Println("  Test accounts:")
	fmt.Println("    - admin:admin@qingyu.com (password: password123)")
	fmt.Println("    - author1:author1@qingyu.com (password: password123)")
	fmt.Println("    - reader1:reader1@qingyu.com (password: password123)")
	fmt.Println("    - reader2:reader2@qingyu.com (password: password123)")

	return nil
}
