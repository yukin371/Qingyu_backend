package examples

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddUserIndexes 添加用户索引迁移
type AddUserIndexes struct{}

func (m *AddUserIndexes) Version() string {
	return "001"
}

func (m *AddUserIndexes) Description() string {
	return "Add indexes to users collection"
}

func (m *AddUserIndexes) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// 创建用户名唯一索引
	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// 创建邮箱唯一索引
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	}

	// 创建手机号唯一索引
	phoneIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true),
	}

	// 创建创建时间索引
	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: -1}},
	}

	indexes := []mongo.IndexModel{
		usernameIndex,
		emailIndex,
		phoneIndex,
		createdAtIndex,
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	fmt.Println("  ✓ Created username unique index")
	fmt.Println("  ✓ Created email unique index")
	fmt.Println("  ✓ Created phone unique index")
	fmt.Println("  ✓ Created created_at index")

	return nil
}

func (m *AddUserIndexes) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("users")

	// 删除索引
	indexNames := []string{
		"username_1",
		"email_1",
		"phone_1",
		"created_at_-1",
	}

	for _, name := range indexNames {
		if _, err := collection.Indexes().DropOne(ctx, name); err != nil {
			fmt.Printf("  Warning: Failed to drop index %s: %v\n", name, err)
		} else {
			fmt.Printf("  ✓ Dropped index %s\n", name)
		}
	}

	return nil
}





