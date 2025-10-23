package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HomepageResponse 首页响应
type HomepageResponse struct {
	Banners        []bookstore2.Banner `json:"banners"`
	RecommendBooks []bookstore2.Book   `json:"recommendBooks"`
	FeaturedBooks  []bookstore2.Book   `json:"featuredBooks"`
	HotBooks       []bookstore2.Book   `json:"hotBooks"`
	NewReleases    []bookstore2.Book   `json:"newReleases"`
	FreeBooks      []bookstore2.Book   `json:"freeBooks"`
	Categories     []CategoryNode      `json:"categories"`
}

// CategoryNode 分类树节点
type CategoryNode struct {
	ID          primitive.ObjectID  `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Icon        string              `json:"icon"`
	ParentID    *primitive.ObjectID `json:"parentId,omitempty"`
	Level       int                 `json:"level"`
	BookCount   int64               `json:"bookCount"`
	SortOrder   int                 `json:"sortOrder"`
	Children    []CategoryNode      `json:"children,omitempty"`
}

// BannerResponse Banner响应
type BannerResponse struct {
	ID         primitive.ObjectID `json:"id"`
	Title      string             `json:"title"`
	ImageURL   string             `json:"imageUrl"`
	TargetType string             `json:"targetType"` // book, category, external
	TargetID   string             `json:"targetId"`
	TargetURL  string             `json:"targetUrl"`
	Position   int                `json:"position"`
	StartTime  time.Time          `json:"startTime"`
	EndTime    time.Time          `json:"endTime"`
	IsActive   bool               `json:"isActive"`
}

// RankingResponse 榜单响应
type RankingResponse struct {
	Type   string              `json:"type"`   // realtime, weekly, monthly, newbie
	Period string              `json:"period"` // 时间周期
	Items  []RankingItemDetail `json:"items"`
	Total  int64               `json:"total"`
}

// RankingItemDetail 榜单项详情
type RankingItemDetail struct {
	Rank  int             `json:"rank"`
	Book  bookstore2.Book `json:"book"`
	Score float64         `json:"score"`
	ViewCount int64          `json:"viewCount"`
	LikeCount int64          `json:"likeCount"`
	Change    int            `json:"change"` // 排名变化
}

// BookListResponse 书籍列表响应
type BookListResponse struct {
	Books []bookstore2.Book `json:"books"`
	Total int64             `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	Category  bookstore2.Category `json:"category"`
	BookCount int64               `json:"bookCount"`
	Children  []CategoryNode     `json:"children,omitempty"`
}
