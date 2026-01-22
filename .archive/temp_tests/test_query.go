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
	
	// 测试查询
	var chapters []bson.M
	cursor, _ := db.Collection("chapters").Find(ctx, bson.M{"book_id": bookID})
	cursor.All(ctx, &chapters)
	
	fmt.Printf("Found %d chapters\n", len(chapters))
	for _, c := range chapters {
		fmt.Printf("Chapter: %+v\n", c["_id"])
	}
}
