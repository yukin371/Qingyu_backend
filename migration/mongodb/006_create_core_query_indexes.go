package migration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateCoreQueryIndexes struct{}

func (m *CreateCoreQueryIndexes) Up(ctx context.Context, db *mongo.Database) error {
	if err := createUserCoreIndexes(ctx, db); err != nil {
		return err
	}
	if err := createCommentCoreIndexes(ctx, db); err != nil {
		return err
	}
	if err := createNotificationCoreIndexes(ctx, db); err != nil {
		return err
	}

	return nil
}

func (m *CreateCoreQueryIndexes) Down(ctx context.Context, db *mongo.Database) error {
	indexGroups := map[string][]string{
		"users": {
			"username_1_unique",
			"email_1_unique_non_empty",
			"phone_1_unique_non_empty",
		},
		"comments": {
			"target_id_1_target_type_1_state_1_created_at_-1",
			"author_id_1_state_1_created_at_-1",
			"parent_id_1_state_1_created_at_1",
		},
		"notifications": {
			"user_id_1_created_at_-1",
			"user_id_1_is_read_1_created_at_-1",
			"user_id_1_type_1_created_at_-1",
		},
	}

	for collectionName, indexNames := range indexGroups {
		col := db.Collection(collectionName)
		for _, indexName := range indexNames {
			_, err := col.Indexes().DropOne(ctx, indexName)
			if err != nil {
				log.Printf("删除索引失败 %s.%s: %v", collectionName, indexName, err)
			} else {
				log.Printf("✅ 删除索引: %s.%s", collectionName, indexName)
			}
		}
	}

	return nil
}

func createUserCoreIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("users")
	models := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().
				SetName("username_1_unique").
				SetUnique(true).
				SetBackground(true),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetName("email_1_unique_non_empty").
				SetUnique(true).
				SetBackground(true).
				SetPartialFilterExpression(bson.M{"email": bson.M{"$exists": true, "$gt": ""}}),
		},
		{
			Keys: bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().
				SetName("phone_1_unique_non_empty").
				SetUnique(true).
				SetBackground(true).
				SetPartialFilterExpression(bson.M{"phone": bson.M{"$exists": true, "$gt": ""}}),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create users core indexes: %w", err)
	}

	log.Printf("✅ Users核心索引创建成功: %v", names)
	return nil
}

func createCommentCoreIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("comments")
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "target_id", Value: 1},
				{Key: "target_type", Value: 1},
				{Key: "state", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("target_id_1_target_type_1_state_1_created_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "state", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("author_id_1_state_1_created_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "parent_id", Value: 1},
				{Key: "state", Value: 1},
				{Key: "created_at", Value: 1},
			},
			Options: options.Index().
				SetName("parent_id_1_state_1_created_at_1").
				SetBackground(true),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create comments core indexes: %w", err)
	}

	log.Printf("✅ Comments核心索引创建成功: %v", names)
	return nil
}

func createNotificationCoreIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("notifications")
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("user_id_1_created_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "is_read", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("user_id_1_is_read_1_created_at_-1").
				SetBackground(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "type", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("user_id_1_type_1_created_at_-1").
				SetBackground(true),
		},
	}

	names, err := col.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create notifications core indexes: %w", err)
	}

	log.Printf("✅ Notifications核心索引创建成功: %v", names)
	return nil
}
