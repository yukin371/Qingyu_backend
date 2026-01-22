package main
import (
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type MongoChapterRepository struct {
	collection *mongo.Collection
}
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(ctx)
	repo := &MongoChapterRepository{
		collection: client.Database("qingyu_read").Collection("chapters"),
	}
	bookID, _ := primitive.ObjectIDFromHex("696f35c4cee9d6ed15e66935")
	chapters, err := repo.GetByBookID(ctx, bookID, 20, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d chapters\n", len(chapters))
	for _, c := range chapters {
		fmt.Printf("Chapter: %s - %s\n", c.ID.Hex(), c.Title)
	}
}
