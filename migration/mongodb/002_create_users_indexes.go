package migration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateUsersIndexes struct{}

func (m *CreateUsersIndexes) Up(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("users")

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
			Keys: bson.D{{Key: "roles", Value: 1}},
			Options: options.Index().
				SetName("roles_1").
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "last_login_at", Value: -1}},
			Options: options.Index().
				SetName("last_login_at_-1").
				SetBackground(true),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create users indexes: %w", err)
	}

	log.Printf("✅ Users索引创建成功: %v", names)
	return nil
}

func (m *CreateUsersIndexes) Down(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("users")
	indexNames := []string{
		"status_1_created_at_-1",
		"roles_1",
		"last_login_at_-1",
	}

	for _, idxName := range indexNames {
		_, err := col.Indexes().DropOne(ctx, idxName)
		if err != nil {
			log.Printf("删除索引失败 users.%s: %v", idxName, err)
		} else {
			log.Printf("✅ 删除索引: users.%s", idxName)
		}
	}

	return nil
}
