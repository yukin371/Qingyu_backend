package debugoutlinequery

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

const defaultProjectID = "69c787f0883ddcfe3b7ff47f"

type optionsValue struct{ projectID string }

// Run compares ObjectID and string project filters and reports decoded outline data.
func Run(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	projectOID, err := primitive.ObjectIDFromHex(opts.projectID)
	if err != nil {
		return fmt.Errorf("项目 ID 无效: %w", err)
	}

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
	coll := client.Database(mongodbconfig.Database()).Collection("outlines")

	count, err := coll.CountDocuments(ctx, bson.M{"project_id": projectOID})
	if err != nil {
		return fmt.Errorf("ObjectID 查询失败: %w", err)
	}
	fmt.Printf("Query with ObjectID: count=%d\n", count)
	count, err = coll.CountDocuments(ctx, bson.M{"project_id": opts.projectID})
	if err != nil {
		return fmt.Errorf("string 查询失败: %w", err)
	}
	fmt.Printf("Query with string: count=%d\n", count)

	var raw bson.M
	err = coll.FindOne(ctx, bson.M{"project_id": projectOID}).Decode(&raw)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("读取原始大纲失败: %w", err)
	}
	if err == mongo.ErrNoDocuments {
		fmt.Println("Raw project_id: <无匹配文档>")
	} else {
		fmt.Printf("Raw project_id: %v (type: %T)\n", raw["project_id"], raw["project_id"])
	}

	type outlineNode struct {
		ID        primitive.ObjectID `bson:"_id"`
		ProjectID primitive.ObjectID `bson:"project_id"`
		Title     string             `bson:"title"`
		ParentID  string             `bson:"parent_id,omitempty"`
	}
	cursor, err := coll.Find(ctx, bson.M{"project_id": projectOID})
	if err != nil {
		return fmt.Errorf("读取大纲列表失败: %w", err)
	}
	defer cursor.Close(ctx)
	var nodes []outlineNode
	if err := cursor.All(ctx, &nodes); err != nil {
		return fmt.Errorf("解析大纲列表失败: %w", err)
	}
	fmt.Printf("Decoded nodes: %d\n", len(nodes))
	for i, node := range nodes {
		if i > 2 {
			break
		}
		fmt.Printf("  - %s (project_id: %s, parent_id: %s)\n", node.Title, node.ProjectID.Hex(), node.ParentID)
	}
	return nil
}

func parseOptions(args []string) (optionsValue, error) {
	fs := flag.NewFlagSet("debug_outline_query", flag.ContinueOnError)
	projectID := fs.String("project-id", defaultProjectID, "目标项目 ObjectID")
	if err := fs.Parse(args); err != nil {
		return optionsValue{}, err
	}
	return optionsValue{projectID: *projectID}, nil
}
