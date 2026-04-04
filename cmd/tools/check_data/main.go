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

	// Check one character
	var ch bson.M
	db.Collection("characters").FindOne(ctx, bson.M{}).Decode(&ch)
	fmt.Printf("Character project_id: %v (type: %T)\n", ch["project_id"], ch["project_id"])

	// Check one location
	var loc bson.M
	db.Collection("locations").FindOne(ctx, bson.M{}).Decode(&loc)
	fmt.Printf("Location project_id: %v (type: %T)\n", loc["project_id"], loc["project_id"])

	// Check one outline
	var ol bson.M
	db.Collection("outlines").FindOne(ctx, bson.M{}).Decode(&ol)
	fmt.Printf("Outline project_id: %v (type: %T)\n", ol["project_id"], ol["project_id"])

	// Check one timeline
	var tl bson.M
	db.Collection("timelines").FindOne(ctx, bson.M{}).Decode(&tl)
	fmt.Printf("Timeline project_id: %v (type: %T)\n", tl["project_id"], tl["project_id"])

	// Check one document
	var doc bson.M
	db.Collection("documents").FindOne(ctx, bson.M{}).Decode(&doc)
	fmt.Printf("Document project_id: %v (type: %T)\n", doc["project_id"], doc["project_id"])
}
