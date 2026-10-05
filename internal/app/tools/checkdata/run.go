package checkdata

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run prints project_id values from representative collections.
func Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI()))
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping 失败: %w", err)
	}
	db := client.Database(databaseName())

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

func mongoURI() string {
	if uri := strings.TrimSpace(os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI")); uri != "" {
		return uri
	}
	if uri := strings.TrimSpace(os.Getenv("MONGODB_URI")); uri != "" {
		return uri
	}
	return "mongodb://localhost:27017/?directConnection=true"
}

func databaseName() string {
	if name := strings.TrimSpace(os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE")); name != "" {
		return name
	}
	if name := strings.TrimSpace(os.Getenv("MONGODB_DATABASE")); name != "" {
		return name
	}
	return "qingyu"
}
