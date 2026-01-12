package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookList 书单
type BookList struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID      string              `bson:"user_id" json:"user_id"`
	UserName    string              `bson:"user_name" json:"user_name"`
	UserAvatar  string              `bson:"user_avatar,omitempty" json:"user_avatar,omitempty"`
	Title       string              `bson:"title" json:"title"`
	Description string              `bson:"description" json:"description"`
	Cover       string              `bson:"cover,omitempty" json:"cover,omitempty"`
	Books       []BookListItem      `bson:"books" json:"books"`
	BookCount   int                 `bson:"book_count" json:"book_count"`
	LikeCount   int                 `bson:"like_count" json:"like_count"`
	ForkCount   int                 `bson:"fork_count" json:"fork_count"` // 被复制次数
	ViewCount   int                 `bson:"view_count" json:"view_count"`
	IsPublic    bool                `bson:"is_public" json:"is_public"`
	Tags        []string            `bson:"tags" json:"tags"`
	Category    string              `bson:"category" json:"category"`                           // 书单分类
	OriginalID  *primitive.ObjectID `bson:"original_id,omitempty" json:"original_id,omitempty"` // 原始书单ID（用于复制）
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updated_at"`
}

// BookListItem 书单中的书籍项
type BookListItem struct {
	BookID      string    `bson:"book_id" json:"book_id"`
	BookTitle   string    `bson:"book_title" json:"book_title"`
	BookCover   string    `bson:"book_cover,omitempty" json:"book_cover,omitempty"`
	AuthorName  string    `bson:"author_name,omitempty" json:"author_name,omitempty"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Comment     string    `bson:"comment,omitempty" json:"comment,omitempty"` // 个人推荐语
	Order       int       `bson:"order" json:"order"`                         // 排序
	AddTime     time.Time `bson:"add_time" json:"add_time"`
}

// BookListLike 书单点赞
type BookListLike struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookListID string             `bson:"booklist_id" json:"booklist_id"`
	UserID     string             `bson:"user_id" json:"user_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

// BookListInfo 书单信息
type BookListInfo struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	UserName    string         `json:"user_name"`
	UserAvatar  string         `json:"user_avatar,omitempty"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Cover       string         `json:"cover,omitempty"`
	Books       []BookListItem `json:"books"`
	BookCount   int            `json:"book_count"`
	LikeCount   int            `json:"like_count"`
	ForkCount   int            `json:"fork_count"`
	ViewCount   int            `json:"view_count"`
	IsLiked     bool           `json:"is_liked"`
	IsPublic    bool           `json:"is_public"`
	Tags        []string       `json:"tags"`
	Category    string         `json:"category"`
	CreatedAt   time.Time      `json:"created_at"`
}
