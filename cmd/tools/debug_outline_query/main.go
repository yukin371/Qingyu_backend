package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?directConnection=true"))
	db := client.Database("qingyu")

	projectID := "69c787f0883ddcfe3b7ff47f"
	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	// Test 1: Query with ObjectID
	filter1 := bson.M{"project_id": projectOID}
	count1, _ := db.Collection("outlines").CountDocuments(ctx, filter1)
	fmt.Printf("Query with ObjectID: count=%d\n", count1)

	// Test 2: Query with string
	filter2 := bson.M{"project_id": projectID}
	count2, _ := db.Collection("outlines").CountDocuments(ctx, filter2)
	fmt.Printf("Query with string: count=%d\n", count2)

	// Test 3: Get one raw document
	var raw bson.M
	db.Collection("outlines").FindOne(ctx, filter1).Decode(&raw)
	fmt.Printf("Raw project_id: %v (type: %T)\n", raw["project_id"], raw["project_id"])

	// Test 4: Try decoding into OutlineNode struct
	type OutlineNode struct {
		ID        primitive.ObjectID `bson:"_id"`
		ProjectID primitive.ObjectID `bson:"project_id"`
		Title     string             `bson:"title"`
		ParentID  string             `bson:"parent_id,omitempty"`
	}
	cursor, _ := db.Collection("outlines").Find(ctx, filter1)
	defer cursor.Close(ctx)
	var nodes []OutlineNode
	cursor.All(ctx, &nodes)
	fmt.Printf("Decoded nodes: %d\n", len(nodes))
	for i, n := range nodes {
		if i > 2 { break }
		fmt.Printf("  - %s (project_id: %s, parent_id: %s)\n", n.Title, n.ProjectID.Hex(), n.ParentID)
	}
}
