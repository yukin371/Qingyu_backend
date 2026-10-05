package cleartestdata

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run clears the known test-data collections only after explicit confirmation.
func Run(args []string) error {
	fs := flag.NewFlagSet("clear_test_data", flag.ContinueOnError)
	apply := fs.Bool("apply", false, "确认清空测试数据")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if !*apply {
		fmt.Println("未执行清理：请显式传入 --apply。")
		return nil
	}
	if os.Getenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA") != "true" {
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

	// 清空相关集合
	collections := []string{"outlines", "documents", "projects", "chapters", "characters", "locations"}
	totalDeleted := int64(0)

	for _, collName := range collections {
		result, err := db.Collection(collName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清空 %s 失败: %w", collName, err)
		}
		totalDeleted += result.DeletedCount
		fmt.Printf("✓ 清空 %s: 删除 %d 条记录\n", collName, result.DeletedCount)
	}

	fmt.Printf("\n总计删除 %d 条记录，数据库已清空\n", totalDeleted)
	return nil
}
