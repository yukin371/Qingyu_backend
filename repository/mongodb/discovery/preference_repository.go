package discovery

import (
	discoveryModel "Qingyu_backend/models/discovery"
	discoveryRepo "Qingyu_backend/repository/interfaces/discovery"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PreferenceRepository struct {
	collection *mongo.Collection
}

func NewPreferenceRepository(db *mongo.Database) discoveryRepo.PreferenceRepository {
	repo := &PreferenceRepository{
		collection: db.Collection("discovery_preferences"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("idx_discovery_preferences_user_id_unique"),
	})
	if err != nil {
		fmt.Printf("Warning: Failed to create discovery preference indexes: %v\n", err)
	}

	return repo
}

func (r *PreferenceRepository) GetByUserID(ctx context.Context, userID string) (*discoveryModel.PreferenceProfile, error) {
	var profile discoveryModel.PreferenceProfile
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&profile)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r *PreferenceRepository) UpsertByUserID(ctx context.Context, userID string, profile *discoveryModel.PreferenceProfile) error {
	if profile == nil {
		return errors.New("discovery 偏好不能为空")
	}

	now := time.Now()
	update := bson.M{
		"user_id":          userID,
		"favorite_genres":  append([]string{}, profile.FavoriteGenres...),
		"favorite_authors": append([]string{}, profile.FavoriteAuthors...),
		"favorite_tags":    append([]string{}, profile.FavoriteTags...),
		"viewed_books":     append([]string{}, profile.ViewedBooks...),
		"viewed_lists":     append([]string{}, profile.ViewedLists...),
		"updated_at":       now,
		"created_at":       now,
	}

	if !profile.CreatedAt.IsZero() {
		update["created_at"] = profile.CreatedAt
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{
			"$set": update,
			"$setOnInsert": bson.M{
				"created_at": now,
			},
		},
		options.Update().SetUpsert(true),
	)

	return err
}
