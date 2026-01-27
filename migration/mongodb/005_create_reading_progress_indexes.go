package migration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateReadingProgressIndexes struct{}

func (m *CreateReadingProgressIndexes) Up(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("reading_progress")

	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "updated_at", Value: -1},
			},
			Options: options.Index().
				SetName("user_id_1_updated_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "book_id", Value: 1}},
			Options: options.Index().
				SetName("book_id_1").
				SetBackground(true),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create reading_progress indexes: %w", err)
	}

	log.Printf("✅ ReadingProgress索引创建成功: %v", names)
	return nil
}

func (m *CreateReadingProgressIndexes) Down(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("reading_progress")
	indexNames := []string{
		"user_id_1_updated_at_-1",
		"book_id_1",
	}

	for _, idxName := range indexNames {
		_, err := col.Indexes().DropOne(ctx, idxName)
		if err != nil {
			log.Printf("删除索引失败 reading_progress.%s: %v", idxName, err)
		} else {
			log.Printf("✅ 删除索引: reading_progress.%s", idxName)
		}
	}

	return nil
}
