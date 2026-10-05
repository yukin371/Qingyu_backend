package fixduplicateglobaloutline

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

type optionsValue struct{ apply bool }

type duplicateGroup struct {
	ID struct {
		ProjectID primitive.ObjectID `bson:"project_id"`
	} `bson:"_id"`
	Count   int                  `bson:"count"`
	IDs     []primitive.ObjectID `bson:"ids"`
	FirstID primitive.ObjectID   `bson:"first_id"`
}

// Run previews or explicitly removes duplicate global outlines.
func Run(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if opts.apply && os.Getenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA") != "true" {
		return fmt.Errorf("删除操作需要设置 QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	duplicates, err := findDuplicates(ctx, coll)
	if err != nil {
		return err
	}
	fmt.Printf("找到 %d 个项目有重复的 global outline\n", len(duplicates))
	for _, duplicate := range duplicates {
		fmt.Printf("\n项目 %s:\n  重复数量: %d\n  保留 ID: %s\n", duplicate.ID.ProjectID.Hex(), duplicate.Count, duplicate.FirstID.Hex())
		for i, id := range duplicate.IDs[1:] {
			if opts.apply {
				fmt.Printf("  删除 ID[%d]: %s\n", i+1, id.Hex())
			} else {
				fmt.Printf("  预览删除 ID[%d]: %s\n", i+1, id.Hex())
			}
		}
	}
	if !opts.apply {
		fmt.Println("预览模式：未删除数据；使用 --apply 并设置安全环境变量后执行")
		return nil
	}
	for _, duplicate := range duplicates {
		for _, id := range duplicate.IDs[1:] {
			result, err := coll.DeleteOne(ctx, bson.M{"_id": id})
			if err != nil {
				return fmt.Errorf("删除重复大纲 %s 失败: %w", id.Hex(), err)
			}
			if result.DeletedCount != 1 {
				return fmt.Errorf("删除重复大纲 %s 未命中记录", id.Hex())
			}
		}
	}

	remaining, err := findDuplicates(ctx, coll)
	if err != nil {
		return fmt.Errorf("复核重复大纲失败: %w", err)
	}
	if len(remaining) == 0 {
		fmt.Println("✓ 所有重复的 global outline 已清理完成！")
	} else {
		fmt.Printf("⚠ 仍有 %d 个项目存在重复\n", len(remaining))
	}
	return printGlobalOutlines(ctx, coll)
}

func findDuplicates(ctx context.Context, coll *mongo.Collection) ([]duplicateGroup, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"type": "global"}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
		{{Key: "$group", Value: bson.M{
			"_id":      bson.M{"project_id": "$project_id"},
			"count":    bson.M{"$sum": 1},
			"ids":      bson.M{"$push": "$_id"},
			"first_id": bson.M{"$first": "$_id"},
		}}},
		{{Key: "$match", Value: bson.M{"count": bson.M{"$gt": 1}}}},
	}
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("聚合查询失败: %w", err)
	}
	defer cursor.Close(ctx)
	var duplicates []duplicateGroup
	if err := cursor.All(ctx, &duplicates); err != nil {
		return nil, fmt.Errorf("解码重复分组失败: %w", err)
	}
	return duplicates, nil
}

func printGlobalOutlines(ctx context.Context, coll *mongo.Collection) error {
	cursor, err := coll.Find(ctx, bson.M{"type": "global"})
	if err != nil {
		return fmt.Errorf("查询 global outline 失败: %w", err)
	}
	defer cursor.Close(ctx)
	var outlines []bson.M
	if err := cursor.All(ctx, &outlines); err != nil {
		return fmt.Errorf("解析 global outline 失败: %w", err)
	}
	fmt.Printf("\n总共 %d 个 global outline:\n", len(outlines))
	for _, outline := range outlines {
		projectID := "unknown"
		if pid, ok := outline["project_id"].(primitive.ObjectID); ok {
			projectID = pid.Hex()
		}
		if len(projectID) > 8 {
			projectID = projectID[:8] + "..."
		}
		fmt.Printf("  - ID: %v, Project: %s\n", outline["_id"], projectID)
	}
	return nil
}

func parseOptions(args []string) (optionsValue, error) {
	fs := flag.NewFlagSet("fix_duplicate_global_outline", flag.ContinueOnError)
	apply := fs.Bool("apply", false, "执行删除")
	if err := fs.Parse(args); err != nil {
		return optionsValue{}, err
	}
	return optionsValue{apply: *apply}, nil
}
