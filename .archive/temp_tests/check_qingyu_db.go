package main
import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func main() {
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(ctx)
	bookID, _ := primitive.ObjectIDFromHex("696f35c4cee9d6ed15e66935")
	
	// 检查 qingyu 数据库
	count, _ := client.Database("qingyu").Collection("chapters").CountDocuments(ctx, bson.M{"book_id": bookID})
	fmt.Printf("qingyu.chapters: %d\n", count)
}
