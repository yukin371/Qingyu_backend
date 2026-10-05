package checkoutline

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run prints outline hierarchy fields for diagnostic inspection.
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

	cursor, err := db.Collection("outlines").Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("查询 outlines 失败: %w", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("解析 outline 失败: %w", err)
		}
		fmt.Printf("id=%v project_id=%v (type=%T) title=%s parent_id=%v doc_id=%v\n",
			doc["_id"], doc["project_id"], doc["project_id"], doc["title"], doc["parent_id"], doc["document_id"])
	}
	if err := cursor.Err(); err != nil {
		return fmt.Errorf("读取 outlines 游标失败: %w", err)
	}
	return nil
}
