// Package models 定义数据模型
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book 书籍模型
type Book struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Author        string             `bson:"author" json:"author"`
	AuthorID      primitive.ObjectID  `bson:"author_id,omitempty" json:"authorId,omitempty"` // 作者用户ID
	Introduction  string             `bson:"introduction" json:"introduction"`
	Cover         string             `bson:"cover" json:"cover"`
	Categories    []string           `bson:"categories" json:"categories"`
	Tags          []string           `bson:"tags" json:"tags"`
	Status        string             `bson:"status" json:"status"`
	Rating        float64            `bson:"rating" json:"rating"`
	RatingCount   int64              `bson:"rating_count" json:"rating_count"`
	ViewCount     int64              `bson:"view_count" json:"view_count"`
	WordCount     int64              `bson:"word_count" json:"word_count"`
	ChapterCount  int                `bson:"chapter_count" json:"chapter_count"`
	Price         float64            `bson:"price" json:"price"`
	IsFree        bool               `bson:"is_free" json:"is_free"`
	IsRecommended bool               `bson:"is_recommended" json:"is_recommended"`
	IsFeatured    bool               `bson:"is_featured" json:"is_featured"`
	IsHot         bool               `bson:"is_hot" json:"is_hot"`
	PublishedAt   time.Time          `bson:"published_at" json:"published_at"`
	LastUpdateAt  time.Time          `bson:"last_update_at" json:"last_update_at"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// Chapter 章节模型
type Chapter struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID      primitive.ObjectID `bson:"book_id" json:"book_id"`
	ChapterNum  int                `bson:"chapter_num" json:"chapter_num"`
	Title       string             `bson:"title" json:"title"`
	WordCount   int                `bson:"word_count" json:"word_count"`
	Price       float64            `bson:"price" json:"price"`
	IsFree      bool               `bson:"is_free" json:"is_free"`
	Status      string             `bson:"status" json:"status"`
	PublishedAt time.Time          `bson:"published_at" json:"published_at"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ChapterContent 章节内容模型
type ChapterContent struct {
	ChapterID primitive.ObjectID `bson:"chapter_id" json:"chapter_id"`
	Content   string             `bson:"content" json:"content"`
	WordCount int                `bson:"word_count" json:"word_count"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
