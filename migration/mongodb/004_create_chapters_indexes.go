package migration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateChaptersIndexes struct{}

func (m *CreateChaptersIndexes) Up(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("chapters")

	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "book_id", Value: 1},
				{Key: "chapter_num", Value: 1},
			},
			Options: options.Index().
				SetName("book_id_1_chapter_num_1").
				SetBackground(true).
				SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "book_id", Value: 1},
				{Key: "status", Value: 1},
				{Key: "chapter_num", Value: 1},
			},
			Options: options.Index().
				SetName("book_id_1_status_1_chapter_num_1").
				SetBackground(true),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create chapters indexes: %w", err)
	}

	log.Printf("✅ Chapters索引创建成功: %v", names)
	return nil
}

func (m *CreateChaptersIndexes) Down(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("chapters")
	indexNames := []string{
		"book_id_1_chapter_num_1",
		"book_id_1_status_1_chapter_num_1",
	}

	for _, idxName := range indexNames {
		_, err := col.Indexes().DropOne(ctx, idxName)
		if err != nil {
			log.Printf("删除索引失败 chapters.%s: %v", idxName, err)
		} else {
			log.Printf("✅ 删除索引: chapters.%s", idxName)
		}
	}

	return nil
}
