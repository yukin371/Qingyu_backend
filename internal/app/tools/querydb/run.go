package querydb

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run prints recent outline and document records without modifying data.
func Run(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	if err := printCollection(ctx, db.Collection("outlines"), "outlines", func(doc bson.M) {
		fmt.Printf("ID: %v, Title: %v, Type: %v, ParentID: %v, DocumentID: %v, ProjectID: %v\n", doc["_id"], doc["title"], doc["type"], doc["parent_id"], doc["document_id"], doc["project_id"])
	}); err != nil {
		return err
	}
	return printCollection(ctx, db.Collection("documents"), "documents", func(doc bson.M) {
		fmt.Printf("ID: %v, Title: %v, Type: %v, OutlineNodeID: %v, ProjectID: %v\n", doc["_id"], doc["title"], doc["type"], doc["outline_node_id"], doc["project_id"])
	})
}

func printCollection(ctx context.Context, collection *mongo.Collection, name string, printDocument func(bson.M)) error {
	fmt.Printf("=== 查询 %s 集合（按时间排序） ===\n", name)
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"$natural": -1}).SetLimit(20))
	if err != nil {
		return fmt.Errorf("查询 %s 失败: %w", name, err)
	}
	defer cursor.Close(ctx)
	var documents []bson.M
	if err := cursor.All(ctx, &documents); err != nil {
		return fmt.Errorf("解析 %s 失败: %w", name, err)
	}
	for _, document := range documents {
		printDocument(document)
	}
	return nil
}
