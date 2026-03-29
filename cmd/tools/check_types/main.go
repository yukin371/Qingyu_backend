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

	// Get all projects
	fmt.Println("=== Projects ===")
	cursor, _ := db.Collection("projects").Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var p bson.M
		cursor.Decode(&p)
		fmt.Printf("  _id=%v (type=%T), title=%v\n", p["_id"], p["_id"], p["title"])
	}

	// Delete non-demo projects (keep only the one created by demo seeder)
	fmt.Println("\n=== Deleting non-demo projects ===")
	// Find the demo project (title=星际觉醒)
	var demoProj bson.M
	db.Collection("projects").FindOne(ctx, bson.M{"title": "星际觉醒"}).Decode(&demoProj)
	if demoProj != nil {
		demoID := demoProj["_id"]
		fmt.Printf("Demo project _id=%v\n", demoID)
		// Delete all other projects
		res, _ := db.Collection("projects").DeleteMany(ctx, bson.M{"_id": bson.M{"$ne": demoID}})
		fmt.Printf("Deleted %d non-demo projects\n", res.DeletedCount)
	}

	// Verify after cleanup
	fmt.Println("\n=== Projects after cleanup ===")
	cursor5, _ := db.Collection("projects").Find(ctx, bson.M{})
	defer cursor5.Close(ctx)
	for cursor5.Next(ctx) {
		var p bson.M
		cursor5.Decode(&p)
		fmt.Printf("  _id=%v, title=%v\n", p["_id"], p["title"])
	}
}
