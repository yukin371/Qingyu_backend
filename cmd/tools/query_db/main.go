package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	if err != nil {
		log.Fatalf("连接MongoDB失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	// 查询最新的项目相关的outlines
	fmt.Println("=== 查询 outlines 集合（按时间排序） ===")
	coll := db.Collection("outlines")

	// 查找所有outlines并按时间降序排列
	opts := options.Find().SetSort(bson.M{"$natural": -1}).SetLimit(20)
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var outlines []bson.M
	if err = cursor.All(ctx, &outlines); err != nil {
		log.Printf("解码失败: %v", err)
		return
	}

	for _, doc := range outlines {
		fmt.Printf("ID: %v, Title: %v, Type: %v, ParentID: %v, DocumentID: %v, ProjectID: %v\n",
			doc["_id"],
			doc["title"],
			doc["type"],
			doc["parent_id"],
			doc["document_id"],
			doc["project_id"],
		)
	}

	// 查询documents集合
	fmt.Println("\n=== 查询 documents 集合（按时间排序） ===")
	docColl := db.Collection("documents")
	cursor2, err := docColl.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"$natural": -1}).SetLimit(20))
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}
	defer cursor2.Close(ctx)

	var docs []bson.M
	if err = cursor2.All(ctx, &docs); err != nil {
		log.Printf("解码失败: %v", err)
		return
	}

	for _, doc := range docs {
		fmt.Printf("ID: %v, Title: %v, Type: %v, OutlineNodeID: %v, ProjectID: %v\n",
			doc["_id"],
			doc["title"],
			doc["type"],
			doc["outline_node_id"],
			doc["project_id"],
		)
	}
}