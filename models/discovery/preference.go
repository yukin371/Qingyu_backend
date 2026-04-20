package discovery

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PreferenceProfile 表示 discovery 页自身的轻量偏好配置。
// 它只承载页面展示偏好与轻量历史，不作为推荐算法画像真相。
type PreferenceProfile struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"userId"`
	FavoriteGenres  []string           `bson:"favorite_genres" json:"favoriteGenres"`
	FavoriteAuthors []string           `bson:"favorite_authors" json:"favoriteAuthors"`
	FavoriteTags    []string           `bson:"favorite_tags" json:"favoriteTags"`
	ViewedBooks     []string           `bson:"viewed_books" json:"viewedBooks"`
	ViewedLists     []string           `bson:"viewed_lists" json:"viewedLists"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}
