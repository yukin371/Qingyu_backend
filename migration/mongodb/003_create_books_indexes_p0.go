package migration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateBooksIndexesP0 struct{}

func (m *CreateBooksIndexesP0) Up(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("books")

	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("status_1_created_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "rating", Value: -1},
			},
			Options: options.Index().
				SetName("status_1_rating_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("author_id_1_status_1_created_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "category_ids", Value: 1},
				{Key: "rating", Value: -1},
			},
			Options: options.Index().
				SetName("category_ids_1_rating_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "is_completed", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().
				SetName("is_completed_1_status_1").
				SetBackground(true),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create books p0 indexes: %w", err)
	}

	log.Printf("✅ Books P0索引创建成功: %v", names)
	return nil
}

func (m *CreateBooksIndexesP0) Down(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("books")
	indexNames := []string{
		"status_1_created_at_-1",
		"status_1_rating_-1",
		"author_id_1_status_1_created_at_-1",
		"category_ids_1_rating_-1",
		"is_completed_1_status_1",
	}

	for _, idxName := range indexNames {
		_, err := col.Indexes().DropOne(ctx, idxName)
		if err != nil {
			log.Printf("删除索引失败 books.%s: %v", idxName, err)
		} else {
			log.Printf("✅ 删除索引: books.%s", idxName)
		}
	}

	return nil
}
