package checkcontent

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run reports the number of stored document contents.
func Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbconfig.URI()))
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping 失败: %w", err)
	}

	db := client.Database(mongodbconfig.Database())

	count, err := db.Collection("document_contents").CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("统计 document_contents 失败: %w", err)
	}
	fmt.Printf("document_contents total: %d\n", count)
	return nil
}
