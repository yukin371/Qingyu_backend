package checkduplicates

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

// Run reports duplicate arc outline titles without modifying data.
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

	cursor, err := client.Database(mongodbconfig.Database()).Collection("outlines").Find(ctx, bson.M{
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

	fmt.Printf("=== Arc 大纲节点列表 ===\n总数: %d\n\n", len(outlines))
	titleCount := make(map[string]int)
	titleToOutlines := make(map[string][]bson.M)
	for _, outline := range outlines {
		title := stringValue(outline, "title")
		titleCount[title]++
		titleToOutlines[title] = append(titleToOutlines[title], outline)
	}

	fmt.Println("标题重复统计:")
	hasDuplicates := false
	for title, count := range titleCount {
		if count > 1 {
			hasDuplicates = true
			fmt.Printf("  \"%s\": %d 个\n", title, count)
		}
	}
	if !hasDuplicates {
		fmt.Println("  无重复")
	}

	fmt.Println("\n详细列表:")
	for i, outline := range outlines {
		id := objectIDValue(outline, "_id")
		title := stringValue(outline, "title")
		dupMark := ""
		if titleCount[title] > 1 {
			dupMark = " [重复!]"
		}
		fmt.Printf("%d. [%s] \"%s\"%s doc=%s parent=%s order=%d\n", i+1, shortID(id), title, dupMark,
			stringValue(outline, "document_id"), stringValue(outline, "parent_id"), orderValue(outline))
	}

	fmt.Println("\n=== 重复项详细对比 ===")
	for title, items := range titleToOutlines {
		if len(items) < 2 {
			continue
		}
		fmt.Printf("\n标题: \"%s\"\n", title)
		for i, outline := range items {
			fmt.Printf("  [%d] ID=%s\n      document_id=%s\n      parent_id=%s\n      order=%d\n      created_at=%v\n      updated_at=%v\n",
				i+1, objectIDValue(outline, "_id").Hex(), stringValue(outline, "document_id"),
				stringValue(outline, "parent_id"), orderValue(outline), dateValue(outline, "created_at"), dateValue(outline, "updated_at"))
		}
	}
	return nil
}

func parseOptions(args []string) (optionsValue, error) {
	fs := flag.NewFlagSet("check_duplicates", flag.ContinueOnError)
	projectID := fs.String("project-id", defaultProjectID, "目标项目 ObjectID")
	if err := fs.Parse(args); err != nil {
		return optionsValue{}, err
	}
	return optionsValue{projectID: *projectID}, nil
}

func stringValue(document bson.M, key string) string {
	value, _ := document[key].(string)
	return value
}

func objectIDValue(document bson.M, key string) primitive.ObjectID {
	value, _ := document[key].(primitive.ObjectID)
	return value
}

func shortID(id primitive.ObjectID) string {
	hex := id.Hex()
	if len(hex) > 8 {
		return hex[:8]
	}
	return hex
}

func orderValue(document bson.M) int64 {
	switch value := document["order"].(type) {
	case int:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case float64:
		return int64(value)
	default:
		return 0
	}
}

func dateValue(document bson.M, key string) time.Time {
	switch value := document[key].(type) {
	case time.Time:
		return value
	case primitive.DateTime:
		return value.Time()
	default:
		return time.Time{}
	}
}
