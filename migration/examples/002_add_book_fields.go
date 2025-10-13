package examples

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddBookFields 添加书籍字段迁移
type AddBookFields struct{}

func (m *AddBookFields) Version() string {
	return "002"
}

func (m *AddBookFields) Description() string {
	return "Add view_count and like_count fields to books"
}

func (m *AddBookFields) Up(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("books")

	// 更新所有书籍，添加新字段
	filter := bson.M{}
	update := bson.M{
		"$set": bson.M{
			"view_count": 0,
			"like_count": 0,
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add book fields: %w", err)
	}

	fmt.Printf("  ✓ Updated %d books with new fields\n", result.ModifiedCount)

	// 创建索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "view_count", Value: -1},
			{Key: "like_count", Value: -1},
		},
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create book index: %w", err)
	}

	fmt.Println("  ✓ Created view_count and like_count index")

	return nil
}

func (m *AddBookFields) Down(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection("books")

	// 删除字段
	filter := bson.M{}
	update := bson.M{
		"$unset": bson.M{
			"view_count": "",
			"like_count": "",
		},
	}

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove book fields: %w", err)
	}

	fmt.Printf("  ✓ Removed fields from %d books\n", result.ModifiedCount)

	// 删除索引
	if _, err := collection.Indexes().DropOne(ctx, "view_count_-1_like_count_-1"); err != nil {
		fmt.Printf("  Warning: Failed to drop index: %v\n", err)
	} else {
		fmt.Println("  ✓ Dropped view_count and like_count index")
	}

	return nil
}







