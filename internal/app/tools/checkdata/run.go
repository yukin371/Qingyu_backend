package checkdata

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run prints project_id values from representative collections.
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

	collections := []struct {
		label      string
		collection string
	}{
		{"Character", "characters"},
		{"Location", "locations"},
		{"Outline", "outlines"},
		{"Timeline", "timelines"},
		{"Document", "documents"},
	}
	for _, item := range collections {
		var document bson.M
		if err := db.Collection(item.collection).FindOne(ctx, bson.M{}).Decode(&document); err != nil {
			return fmt.Errorf("读取 %s 失败: %w", item.collection, err)
		}
		value := document["project_id"]
		fmt.Printf("%s project_id: %v (type: %T)\n", item.label, value, value)
	}
	return nil
}
