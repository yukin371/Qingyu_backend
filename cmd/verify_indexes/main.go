package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

type indexSpecFile struct {
	Collections map[string]collectionSpec `yaml:"collections"`
}

type collectionSpec struct {
	Indexes []indexSpec `yaml:"indexes"`
}

type indexSpec struct {
	Name string `yaml:"name"`
}

func main() {
	specPath := flag.String("spec", "docs/database/indexes.yaml", "索引规范文件路径")
	showExtra := flag.Bool("show-extra", true, "是否输出规范外索引")
	flag.Parse()

	mongoURI := firstNonEmpty(
		os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_URI"),
		os.Getenv("MONGODB_URI"),
		os.Getenv("MONGO_URI"),
		"mongodb://localhost:27017",
	)
	dbName := firstNonEmpty(
		os.Getenv("QINGYU_DATABASE_PRIMARY_MONGODB_DATABASE"),
		os.Getenv("MONGODB_DATABASE"),
		"qingyu_dev",
	)

	spec, err := loadSpec(*specPath)
	if err != nil {
		exitf("加载索引规范失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		exitf("连接 MongoDB 失败: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		exitf("MongoDB ping 失败: %v", err)
	}

	db := client.Database(dbName)
	fmt.Printf("Verifying indexes against %s on database %s\n\n", *specPath, dbName)

	hasMissing := false
	for collectionName, collection := range spec.Collections {
		expected := make(map[string]struct{}, len(collection.Indexes))
		for _, idx := range collection.Indexes {
			expected[idx.Name] = struct{}{}
		}

		actual, err := listIndexNames(ctx, db.Collection(collectionName))
		if err != nil {
			exitf("读取集合 %s 索引失败: %v", collectionName, err)
		}

		missing := diffMapKeys(expected, actual)
		extra := diffMapKeys(actual, expected)

		sort.Strings(missing)
		sort.Strings(extra)

		status := "OK"
		if len(missing) > 0 {
			status = "MISSING"
			hasMissing = true
		}

		fmt.Printf("[%s] %s\n", status, collectionName)
		if len(missing) == 0 {
			fmt.Println("  missing: none")
		} else {
			for _, name := range missing {
				fmt.Printf("  missing: %s\n", name)
			}
		}

		if *showExtra {
			if len(extra) == 0 {
				fmt.Println("  extra: none")
			} else {
				for _, name := range extra {
					fmt.Printf("  extra: %s\n", name)
				}
			}
		}
		fmt.Println()
	}

	if hasMissing {
		os.Exit(1)
	}
}

func loadSpec(path string) (*indexSpecFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var spec indexSpecFile
	if err := yaml.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}

	if len(spec.Collections) == 0 {
		return nil, fmt.Errorf("索引规范为空")
	}
	return &spec, nil
}

func listIndexNames(ctx context.Context, collection *mongo.Collection) (map[string]struct{}, error) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	names := make(map[string]struct{})
	for cursor.Next(ctx) {
		var indexDoc bson.M
		if err := cursor.Decode(&indexDoc); err != nil {
			return nil, err
		}
		name, _ := indexDoc["name"].(string)
		if name == "" || name == "_id_" {
			continue
		}
		names[name] = struct{}{}
	}
	return names, cursor.Err()
}

func diffMapKeys(left, right map[string]struct{}) []string {
	var diff []string
	for key := range left {
		if _, ok := right[key]; !ok {
			diff = append(diff, key)
		}
	}
	return diff
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func exitf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
