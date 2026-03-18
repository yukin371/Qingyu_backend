package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var knownTextReplacements = map[string]string{
	"ϵͳ��������":     "系统维护通知",
	"����һ�����Թ���": "这是一条测试公告",
}

type collectionPlan struct {
	Name          string
	IDField       string
	Fields        []string
	SuspiciousDoc func(bson.M) bool
}

func main() {
	var (
		mongoURI = flag.String("uri", "mongodb://127.0.0.1:27017", "MongoDB 连接 URI")
		dbName   = flag.String("db", "qingyu", "数据库名称")
		apply    = flag.Bool("apply", false, "应用修复；默认仅预览")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		log.Fatalf("连接 MongoDB 失败: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(*dbName)
	fmt.Printf("连接数据库成功: %s/%s\n", *mongoURI, *dbName)
	if *apply {
		fmt.Println("模式: APPLY")
	} else {
		fmt.Println("模式: DRY RUN")
	}
	fmt.Println()

	plans := []collectionPlan{
		{
			Name:    "announcements",
			IDField: "_id",
			Fields:  []string{"title", "content"},
			SuspiciousDoc: func(doc bson.M) bool {
				return containsSuspiciousText(asString(doc["title"])) || containsSuspiciousText(asString(doc["content"]))
			},
		},
		{
			Name:    "banners",
			IDField: "_id",
			Fields:  []string{"title", "description"},
			SuspiciousDoc: func(doc bson.M) bool {
				return containsSuspiciousText(asString(doc["title"])) || containsSuspiciousText(asString(doc["description"]))
			},
		},
	}

	totalUpdated := int64(0)
	totalSuspicious := int64(0)
	for _, plan := range plans {
		updated, suspicious, err := processCollection(ctx, db.Collection(plan.Name), plan, *apply)
		if err != nil {
			log.Fatalf("处理集合 %s 失败: %v", plan.Name, err)
		}
		totalUpdated += updated
		totalSuspicious += suspicious
	}

	fmt.Println()
	fmt.Printf("处理完成: 更新 %d 条记录，仍有 %d 条可疑记录需要人工确认\n", totalUpdated, totalSuspicious)
}

func processCollection(ctx context.Context, coll *mongo.Collection, plan collectionPlan, apply bool) (int64, int64, error) {
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var updatedCount int64
	var suspiciousCount int64

	fmt.Printf("== 集合: %s ==\n", coll.Name())
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return updatedCount, suspiciousCount, err
		}

		updates := bson.M{}
		for _, field := range plan.Fields {
			original := asString(doc[field])
			fixed := applyKnownReplacements(original)
			if fixed != original {
				updates[field] = fixed
			}
		}

		// 已确认这是一条前后端联调阶段遗留的测试 Banner，需要连同失效图片地址一起修复。
		if coll.Name() == "banners" &&
			asString(doc["target"]) == "/test" &&
			asString(doc["image"]) == "https://example.com/test.jpg" {
			updates["image"] = "/images/banners/showcase-yunhai.jpg"
			if containsSuspiciousText(asString(doc["title"]) + asString(doc["description"])) {
				updates["title"] = "前后端联调测试 Banner"
				updates["description"] = "前后端联调用测试资源占位"
			}
		}

		id := formatID(doc[plan.IDField])
		if len(updates) > 0 {
			fmt.Printf("[FIX] %s %s\n", coll.Name(), id)
			for field, value := range updates {
				fmt.Printf("  - %s: %q\n", field, value)
			}

			if apply {
				updateDoc := bson.M{"$set": updates}
				if _, err := coll.UpdateByID(ctx, doc[plan.IDField], updateDoc); err != nil {
					return updatedCount, suspiciousCount, err
				}
			}
			updatedCount++
			continue
		}

		if plan.SuspiciousDoc != nil && plan.SuspiciousDoc(doc) {
			suspiciousCount++
			fmt.Printf("[WARN] %s %s 仍包含可疑文本\n", coll.Name(), id)
			for _, field := range plan.Fields {
				value := asString(doc[field])
				if containsSuspiciousText(value) {
					fmt.Printf("  - %s: %q\n", field, value)
				}
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return updatedCount, suspiciousCount, err
	}

	fmt.Printf("汇总: 更新 %d 条，待确认 %d 条\n\n", updatedCount, suspiciousCount)
	return updatedCount, suspiciousCount, nil
}

func applyKnownReplacements(text string) string {
	fixed := text
	for from, to := range knownTextReplacements {
		fixed = strings.ReplaceAll(fixed, from, to)
	}
	return fixed
}

func containsSuspiciousText(text string) bool {
	if text == "" {
		return false
	}

	suspiciousFragments := []string{
		"����",
		"���",
		"ϵͳ",
		"֪ͨ",
		"�",
	}

	for _, fragment := range suspiciousFragments {
		if strings.Contains(text, fragment) {
			return true
		}
	}
	return false
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		return ""
	}
}

func formatID(v any) string {
	switch id := v.(type) {
	case primitive.ObjectID:
		return id.Hex()
	case string:
		return id
	default:
		return fmt.Sprintf("%v", v)
	}
}
