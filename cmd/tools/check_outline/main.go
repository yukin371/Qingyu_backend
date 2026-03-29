package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	db := client.Database("qingyu")

	cursor, _ := db.Collection("outlines").Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc bson.M
		cursor.Decode(&doc)
		fmt.Printf("id=%v project_id=%v (type=%T) title=%s parent_id=%v doc_id=%v\n",
			doc["_id"], doc["project_id"], doc["project_id"], doc["title"], doc["parent_id"], doc["document_id"])
	}
}
