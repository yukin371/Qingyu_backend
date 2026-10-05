package checktypes

import (
	"context"
	"flag"
	"fmt"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run inspects projects and optionally removes non-demo projects.
// Destructive cleanup is opt-in through --apply.
func Run(args []string) error {
	apply, err := parseArgs(args)
	if err != nil {
		return err
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
	db := client.Database(mongodbconfig.Database())

	// Get all projects
	fmt.Println("=== Projects ===")
	cursor, err := db.Collection("projects").Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("查询 projects 失败: %w", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var p bson.M
		if err := cursor.Decode(&p); err != nil {
			return fmt.Errorf("解析 project 失败: %w", err)
		}
		fmt.Printf("  _id=%v (type=%T), title=%v\n", p["_id"], p["_id"], p["title"])
	}
	if err := cursor.Err(); err != nil {
		return fmt.Errorf("读取 projects 游标失败: %w", err)
	}
	if !apply {
		fmt.Println("\n未执行删除：如需清理非演示项目，请显式传入 --apply。")
		return nil
	}

	// Delete non-demo projects (keep only the one created by demo seeder)
	fmt.Println("\n=== Deleting non-demo projects ===")
	// Find the demo project (title=星际觉醒)
	var demoProj bson.M
	if err := db.Collection("projects").FindOne(ctx, bson.M{"title": "星际觉醒"}).Decode(&demoProj); err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("查询演示项目失败: %w", err)
	}
	if demoProj != nil {
		demoID := demoProj["_id"]
		fmt.Printf("Demo project _id=%v\n", demoID)
		// Delete all other projects
		res, err := db.Collection("projects").DeleteMany(ctx, bson.M{"_id": bson.M{"$ne": demoID}})
		if err != nil {
			return fmt.Errorf("删除非演示项目失败: %w", err)
		}
		fmt.Printf("Deleted %d non-demo projects\n", res.DeletedCount)
	} else {
		return fmt.Errorf("未找到演示项目，拒绝执行全量删除")
	}

	// Verify after cleanup
	fmt.Println("\n=== Projects after cleanup ===")
	cursor5, err := db.Collection("projects").Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("清理后查询 projects 失败: %w", err)
	}
	defer cursor5.Close(ctx)
	for cursor5.Next(ctx) {
		var p bson.M
		if err := cursor5.Decode(&p); err != nil {
			return fmt.Errorf("清理后解析 project 失败: %w", err)
		}
		fmt.Printf("  _id=%v, title=%v\n", p["_id"], p["title"])
	}
	if err := cursor5.Err(); err != nil {
		return fmt.Errorf("清理后读取 projects 游标失败: %w", err)
	}
	return nil
}

func parseArgs(args []string) (bool, error) {
	fs := flag.NewFlagSet("check_types", flag.ContinueOnError)
	apply := fs.Bool("apply", false, "确认删除非演示项目")
	if err := fs.Parse(args); err != nil {
		return false, err
	}
	return *apply, nil
}
