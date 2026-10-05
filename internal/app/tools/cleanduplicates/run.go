package cleanduplicates

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
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

// Run reports duplicate arc outlines and deletes them only with explicit opt-in.
func Run(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	projectOID, err := primitive.ObjectIDFromHex(opts.projectID)
	if err != nil {
		return fmt.Errorf("项目 ID 无效: %w", err)
	}
	if opts.apply && os.Getenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA") != "true" {
		return fmt.Errorf("拒绝执行破坏性清理：请设置 QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true")
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
	cursor, err := db.Collection("outlines").Find(ctx, bson.M{
		"project_id": projectOID,
		"type":       "arc",
	}, options.Find().SetSort(bson.M{"order": 1}))
	if err != nil {
		return fmt.Errorf("查询 outlines 失败: %w", err)
	}
	var outlines []bson.M
	if err := cursor.All(ctx, &outlines); err != nil {
		return fmt.Errorf("解析 outlines 失败: %w", err)
	}

	docToOutlines := make(map[string][]bson.M)
	for _, outline := range outlines {
		docID, ok := stringField(outline, "document_id")
		if ok && docID != "" {
			docToOutlines[docID] = append(docToOutlines[docID], outline)
		}
	}

	fmt.Printf("=== 清理重复大纲节点 ===\n\n")
	totalDeleted := 0
	for docID, outlineList := range docToOutlines {
		if len(outlineList) < 2 {
			continue
		}
		fmt.Printf("文档 %s 有 %d 个关联大纲:\n", docID, len(outlineList))
		sort.SliceStable(outlineList, func(i, j int) bool {
			return outlineCreatedAt(outlineList[i]).Before(outlineCreatedAt(outlineList[j]))
		})
		keepID, ok := objectIDField(outlineList[0], "_id")
		if !ok {
			return fmt.Errorf("文档 %s 的保留大纲缺少有效 _id", docID)
		}
		fmt.Printf("  保留: %s\n", keepID.Hex())
		for _, outline := range outlineList {
			id, ok := objectIDField(outline, "_id")
			if !ok {
				return fmt.Errorf("文档 %s 存在无效 _id", docID)
			}
			fmt.Printf("  - %s (created_at: %v)\n", id.Hex(), outlineCreatedAt(outline))
			if id == keepID || !opts.apply {
				continue
			}
			if _, err := db.Collection("outlines").DeleteOne(ctx, bson.M{"_id": id}); err != nil {
				return fmt.Errorf("删除重复大纲 %s 失败: %w", id.Hex(), err)
			}
			totalDeleted++
		}
	}

	if !opts.apply {
		fmt.Println("\n预览模式：未删除数据，请显式传入 --apply 并设置 QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true。")
		return nil
	}
	count, err := db.Collection("outlines").CountDocuments(ctx, bson.M{
		"project_id": projectOID,
		"type":       "arc",
	})
	if err != nil {
		return fmt.Errorf("清理后统计 outlines 失败: %w", err)
	}
	fmt.Printf("清理完成: 共删除 %d 个重复大纲节点\n剩余 arc 大纲数量: %d\n", totalDeleted, count)
	return nil
}

func parseOptions(args []string) (optionsValue, error) {
	fs := flag.NewFlagSet("clean_duplicates", flag.ContinueOnError)
	projectID := fs.String("project-id", defaultProjectID, "目标项目 ObjectID")
	apply := fs.Bool("apply", false, "确认删除重复大纲")
	if err := fs.Parse(args); err != nil {
		return optionsValue{}, err
	}
	return optionsValue{projectID: *projectID, apply: *apply}, nil
}

func stringField(document bson.M, key string) (string, bool) {
	value, ok := document[key].(string)
	return value, ok
}

func objectIDField(document bson.M, key string) (primitive.ObjectID, bool) {
	value, ok := document[key].(primitive.ObjectID)
	return value, ok
}

func outlineCreatedAt(document bson.M) time.Time {
	switch value := document["created_at"].(type) {
	case time.Time:
		return value
	case primitive.DateTime:
		return value.Time()
	default:
		return time.Time{}
	}
}
