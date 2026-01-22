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
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(ctx)
	db := client.Database("qingyu_read")
	bookID, _ := primitive.ObjectIDFromHex("696f35c4cee9d6ed15e66935")
	// 检查章节
	count, _ := db.Collection("chapters").CountDocuments(ctx, bson.M{"book_id": bookID})
	fmt.Printf("chapters collection count with book_id %s: %d\n", bookID.Hex(), count)
	// 查看前5个章节
	cursor, _ := db.Collection("chapters").Find(ctx, bson.M{"book_id": bookID}, options.Find().SetLimit(5))
	defer cursor.Close(ctx)
	var results []bson.M
	if err := cursor.All(ctx, &results); err == nil {
		for i, r := range results {
			fmt.Printf("Chapter %d: %+v\n", i+1, r)
		}
	}
}
