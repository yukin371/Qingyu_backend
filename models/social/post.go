package social

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostType 动态类型
type PostType string

const (
	PostTypeText              PostType = "text"               // 纯文本
	PostTypeImage             PostType = "image"             // 图片
	PostTypeBookRecommendation PostType = "book_recommendation" // 书籍推荐
	PostTypeReadingProgress   PostType = "reading_progress"   // 阅读进度
	PostTypePoll             PostType = "poll"               // 投票
)

// Post 动态/帖子
type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"userId"`
	UserName  string             `bson:"user_name" json:"userName"`
	UserAvatar string            `bson:"user_avatar,omitempty" json:"userAvatar,omitempty"`
	UserLevel int               `bson:"user_level" json:"userLevel"`
	Type      PostType           `bson:"type" json:"type"` // text|image|book_recommendation|reading_progress
	Content   string             `bson:"content" json:"content"`
	Images    []string           `bson:"images,omitempty" json:"images,omitempty"`
	BookID    string             `bson:"book_id,omitempty" json:"bookId,omitempty"`
	BookTitle string             `bson:"book_title,omitempty" json:"bookTitle,omitempty"`
	BookCover string             `bson:"book_cover,omitempty" json:"bookCover,omitempty"`
	BookAuthor string            `bson:"book_author,omitempty" json:"bookAuthor,omitempty"`
	ChapterID string             `bson:"chapter_id,omitempty" json:"chapterId,omitempty"`
	ChapterTitle string          `bson:"chapter_title,omitempty" json:"chapterTitle,omitempty"`
	Progress  int               `bson:"progress,omitempty" json:"progress,omitempty"`
	Topics    []string           `bson:"topics,omitempty" json:"topics,omitempty"`
	LikeCount    int            `bson:"like_count" json:"likeCount"`
	CommentCount int            `bson:"comment_count" json:"commentCount"`
	ShareCount   int            `bson:"share_count" json:"shareCount"`
	CreatedAt time.Time        `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updatedAt"`
}

// PostLike 动态点赞
type PostLike struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID    string             `bson:"post_id" json:"postId"`
	UserID    string             `bson:"user_id" json:"userId"`
	CreatedAt time.Time         `bson:"created_at" json:"createdAt"`
}

// PostInfo 返回给前端的动态信息
type PostInfo struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"userId"`
	UserName     string                 `json:"userName"`
	UserAvatar   string                 `json:"userAvatar"`
	UserLevel    int                    `json:"userLevel"`
	Type         PostType               `json:"type"`
	Content      string                 `json:"content"`
	Images       []string               `json:"images,omitempty"`
	Book         *BookInfo              `json:"book,omitempty"`
	ReadingProgress *ReadingProgressInfo  `json:"readingProgress,omitempty"`
	Topics       []string               `json:"topics"`
	LikeCount    int                    `json:"likeCount"`
	CommentCount int                    `json:"commentCount"`
	ShareCount   int                    `json:"shareCount"`
	IsLiked      bool                   `json:"isLiked"`
	IsBookmarked bool                   `json:"isBookmarked"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt,omitempty"`
}

// BookInfo 书籍信息
type BookInfo struct {
	BookID  string `json:"bookId"`
	Title   string `json:"title"`
	Cover   string `json:"cover"`
	Author  string `json:"author"`
}

// ReadingProgressInfo 阅读进度信息
type ReadingProgressInfo struct {
	BookID       string `json:"bookId"`
	ChapterID    string `json:"chapterId"`
	ChapterTitle string `json:"chapterTitle"`
	Progress     int    `json:"progress"`
}

// ToPostInfo 将 Post 转换为 PostInfo
func (p *Post) ToPostInfo(isLiked, isBookmarked bool) *PostInfo {
	info := &PostInfo{
		ID:           p.ID.Hex(),
		UserID:       p.UserID,
		UserName:     p.UserName,
		UserAvatar:   p.UserAvatar,
		UserLevel:    p.UserLevel,
		Type:         p.Type,
		Content:      p.Content,
		Topics:       p.Topics,
		LikeCount:    p.LikeCount,
		CommentCount: p.CommentCount,
		ShareCount:   p.ShareCount,
		IsLiked:      isLiked,
		IsBookmarked: isBookmarked,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	if len(p.Images) > 0 {
		info.Images = p.Images
	}

	if p.BookID != "" {
		info.Book = &BookInfo{
			BookID:  p.BookID,
			Title:   p.BookTitle,
			Cover:   p.BookCover,
			Author:  p.BookAuthor,
		}
	}

	if p.Type == PostTypeReadingProgress && p.ChapterID != "" {
		info.ReadingProgress = &ReadingProgressInfo{
			BookID:       p.BookID,
			ChapterID:    p.ChapterID,
			ChapterTitle: p.ChapterTitle,
			Progress:     p.Progress,
		}
	}

	return info
}
