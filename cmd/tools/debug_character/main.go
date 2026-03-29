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

	// Raw MongoDB document
	var raw bson.M
	db.Collection("characters").FindOne(ctx, bson.M{}).Decode(&raw)
	for k, v := range raw {
		fmt.Printf("  %s: %v (type: %T)\n", k, v, v)
	}
}
