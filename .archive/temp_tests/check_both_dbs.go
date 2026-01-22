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
	
	// 检查所有数据库
	dbs := []string{"Qingyu_backend", "qingyu_read", "qingyu"}
	for _, dbName := range dbs {
		count, _ := client.Database(dbName).Collection("chapters").CountDocuments(ctx, bson.M{"book_id": bookID})
		fmt.Printf("%s.chapters: %d\n", dbName, count)
	}
}
