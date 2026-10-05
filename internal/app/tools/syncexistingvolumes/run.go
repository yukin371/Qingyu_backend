package syncexistingvolumes

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultProjectID = "69c790d7be2b6d81b914bfed"

type optionsValue struct {
	projectID string
	apply     bool
}

// Run previews or explicitly applies volume-to-outline synchronization.
func Run(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if opts.apply && os.Getenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA") != "true" {
		return fmt.Errorf("写入同步需要设置 QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true")
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
	cursor, err := db.Collection("documents").Find(ctx, bson.M{
		"project_id": projectOID,
		"type":       "volume",
		"$or": []bson.M{
			{"outline_node_id": bson.M{"$exists": false}},
			{"outline_node_id": ""},
		},
	})
	if err != nil {
		return fmt.Errorf("查询卷文档失败: %w", err)
	}
	defer cursor.Close(ctx)
	var documents []bson.M
	if err := cursor.All(ctx, &documents); err != nil {
		return fmt.Errorf("解析卷文档失败: %w", err)
	}

	fmt.Printf("找到 %d 个需要同步的卷文档\n", len(documents))
	if !opts.apply {
		fmt.Println("预览模式：未写入数据；使用 --apply 并设置安全环境变量后执行")
		return nil
	}

	synced := 0
	for _, document := range documents {
		documentID, ok := document["_id"].(primitive.ObjectID)
		if !ok {
			return fmt.Errorf("卷文档 _id 类型无效: %v", document["_id"])
		}
		title, ok := document["title"].(string)
		if !ok || title == "" {
			return fmt.Errorf("卷文档 title 无效: %v", document["title"])
		}
		order, err := numericValue(document["order"])
		if err != nil {
			return fmt.Errorf("卷文档 %s order 无效: %w", documentID.Hex(), err)
		}
		outlineID := primitive.NewObjectID()
		outline := bson.M{
			"_id": outlineID, "project_id": projectOID, "title": title,
			"parent_id": "", "order": order, "type": "arc", "summary": "", "tension": 5,
			"document_id": documentID.Hex(), "characters": bson.A{}, "items": bson.A{}, "tags": bson.A{},
			"created_at": time.Now(), "updated_at": time.Now(),
		}
		if _, err := db.Collection("outlines").InsertOne(ctx, outline); err != nil {
			return fmt.Errorf("创建大纲失败 [%s]: %w", title, err)
		}
		if _, err := db.Collection("documents").UpdateOne(ctx, bson.M{"_id": documentID}, bson.M{"$set": bson.M{"outline_node_id": outlineID.Hex()}}); err != nil {
			_, _ = db.Collection("outlines").DeleteOne(context.Background(), bson.M{"_id": outlineID})
			return fmt.Errorf("更新文档 outline_node_id 失败 [%s]: %w", title, err)
		}
		fmt.Printf("✓ 同步完成: %s -> outline=%s\n", title, outlineID.Hex())
		synced++
	}
	fmt.Printf("\n同步完成: %d/%d\n", synced, len(documents))
	return nil
}

func parseOptions(args []string) (optionsValue, error) {
	fs := flag.NewFlagSet("sync_existing_volumes", flag.ContinueOnError)
	projectID := fs.String("project-id", defaultProjectID, "目标项目 ObjectID")
	apply := fs.Bool("apply", false, "执行写入")
	if err := fs.Parse(args); err != nil {
		return optionsValue{}, err
	}
	return optionsValue{projectID: *projectID, apply: *apply}, nil
}

func numericValue(value interface{}) (int64, error) {
	switch value := value.(type) {
	case int:
		return int64(value), nil
	case int32:
		return int64(value), nil
	case int64:
		return value, nil
	case float64:
		return int64(value), nil
	default:
		return 0, fmt.Errorf("期望数字，实际为 %T", value)
	}
}
