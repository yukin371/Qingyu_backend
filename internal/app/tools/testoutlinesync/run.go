package testoutlinesync

import (
	"context"
	"flag"
	"fmt"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultProjectID = "69c790d7be2b6d81b914bfed"

type optionsValue struct{ projectID string }

// Run checks bidirectional outline/document associations without modifying data.
func Run(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	projectOID, err := primitive.ObjectIDFromHex(opts.projectID)
	if err != nil {
		return fmt.Errorf("项目 ID 无效: %w", err)
	}
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

	outlines, err := findDocuments(ctx, db.Collection("outlines"), bson.M{
		"project_id": projectOID,
		"$or":        []bson.M{{"document_id": bson.M{"$exists": false}}, {"document_id": ""}},
	})
	if err != nil {
		return fmt.Errorf("查询未关联大纲失败: %w", err)
	}
	documents, err := findDocuments(ctx, db.Collection("documents"), bson.M{
		"project_id": projectOID,
		"$or":        []bson.M{{"outline_node_id": bson.M{"$exists": false}}, {"outline_node_id": ""}},
	})
	if err != nil {
		return fmt.Errorf("查询未关联文档失败: %w", err)
	}
	fmt.Printf("找到 %d 个没有关联文档的大纲节点\n", len(outlines))
	fmt.Printf("找到 %d 个没有关联大纲的文档\n", len(documents))

	fmt.Println("\n=== 双向绑定检查 ===")
	counts := []struct {
		label  string
		coll   string
		filter bson.M
	}{
		{"文档有 outline_node_id", "documents", bson.M{"project_id": projectOID, "outline_node_id": bson.M{"$ne": ""}}},
		{"大纲有 document_id", "outlines", bson.M{"project_id": projectOID, "document_id": bson.M{"$ne": ""}}},
		{"卷文档数量", "documents", bson.M{"project_id": projectOID, "type": "volume"}},
		{"章节文档数量", "documents", bson.M{"project_id": projectOID, "type": "chapter"}},
		{"arc 大纲数量", "outlines", bson.M{"project_id": projectOID, "type": "arc"}},
		{"scene 大纲数量", "outlines", bson.M{"project_id": projectOID, "type": "scene"}},
	}
	for _, item := range counts {
		count, err := db.Collection(item.coll).CountDocuments(ctx, item.filter)
		if err != nil {
			return fmt.Errorf("统计 %s 失败: %w", item.label, err)
		}
		fmt.Printf("%s: %d\n", item.label, count)
	}

	fmt.Println("\n建议:")
	if len(outlines) > 0 {
		fmt.Printf("- 需要为 %d 个大纲创建对应文档\n", len(outlines))
	}
	if len(documents) > 0 {
		fmt.Printf("- 需要为 %d 个文档创建对应大纲\n", len(documents))
	}
	return nil
}

func findDocuments(ctx context.Context, collection *mongo.Collection, filter bson.M) ([]bson.M, error) {
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var documents []bson.M
	if err := cursor.All(ctx, &documents); err != nil {
		return nil, err
	}
	return documents, nil
}

func parseOptions(args []string) (optionsValue, error) {
	fs := flag.NewFlagSet("test_outline_sync", flag.ContinueOnError)
	projectID := fs.String("project-id", defaultProjectID, "目标项目 ObjectID")
	if err := fs.Parse(args); err != nil {
		return optionsValue{}, err
	}
	return optionsValue{projectID: *projectID}, nil
}
