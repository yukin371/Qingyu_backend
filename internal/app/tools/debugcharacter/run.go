package debugcharacter

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run prints the BSON shape of one character document without modifying data.
func Run(args []string) error {
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

	var raw bson.M
	err = client.Database(mongodbconfig.Database()).Collection("characters").FindOne(ctx, bson.M{}).Decode(&raw)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("characters 集合没有文档")
		}
		return fmt.Errorf("读取 character 文档失败: %w", err)
	}
	for key, value := range raw {
		fmt.Printf("  %s: %v (type: %T)\n", key, value, value)
	}
	return nil
}
