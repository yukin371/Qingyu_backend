package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookList 书单
type BookList struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID      string              `bson:"user_id" json:"userId"`
	UserName    string              `bson:"user_name" json:"userName"`
	UserAvatar  string              `bson:"user_avatar,omitempty" json:"user_avatar,omitempty"`
	Title       string              `bson:"title" json:"title"`
	Description string              `bson:"description" json:"description"`
	Cover       string              `bson:"cover,omitempty" json:"cover,omitempty"`
	Books       []BookListItem      `bson:"books" json:"books"`
	BookCount   int                 `bson:"book_count" json:"bookCount"`
	LikeCount   int                 `bson:"like_count" json:"likeCount"`
	ForkCount   int                 `bson:"fork_count" json:"forkCount"` // 被复制次数
	ViewCount   int                 `bson:"view_count" json:"viewCount"`
	IsPublic    bool                `bson:"is_public" json:"isPublic"`
	Tags        []string            `bson:"tags" json:"tags"`
	Category    string              `bson:"category" json:"category"`                           // 书单分类
	OriginalID  *primitive.ObjectID `bson:"original_id,omitempty" json:"original_id,omitempty"` // 原始书单ID（用于复制）
	CreatedAt   time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updatedAt"`
}

// BookListItem 书单中的书籍项
type BookListItem struct {
	BookID      string    `bson:"book_id" json:"book_id"`
	BookTitle   string    `bson:"book_title" json:"bookTitle"`
	BookCover   string    `bson:"book_cover,omitempty" json:"book_cover,omitempty"`
	AuthorName  string    `bson:"author_name,omitempty" json:"author_name,omitempty"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Comment     string    `bson:"comment,omitempty" json:"comment,omitempty"` // 个人推荐语
	Order       int       `bson:"order" json:"order"`                         // 排序
	AddTime     time.Time `bson:"add_time" json:"addTime"`
}

// BookListLike 书单点赞
type BookListLike struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookListID string             `bson:"booklist_id" json:"booklistId"`
	UserID     string             `bson:"user_id" json:"userId"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
}

// BookListInfo 书单信息
type BookListInfo struct {
	ID          string         `json:"id"`
	UserID      string         `json:"userId"`
	UserName    string         `json:"userName"`
	UserAvatar  string         `json:"user_avatar,omitempty"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Cover       string         `json:"cover,omitempty"`
	Books       []BookListItem `json:"books"`
	BookCount   int            `json:"bookCount"`
	LikeCount   int            `json:"likeCount"`
	ForkCount   int            `json:"forkCount"`
	ViewCount   int            `json:"viewCount"`
	IsLiked     bool           `json:"isLiked"`
	IsPublic    bool           `json:"isPublic"`
	Tags        []string       `json:"tags"`
	Category    string         `json:"category"`
	CreatedAt   time.Time      `json:"createdAt"`
}
