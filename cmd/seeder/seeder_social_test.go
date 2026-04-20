package main

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBuildSeedPosts_ProducesBookBoundPosts(t *testing.T) {
	users := []socialSeedUser{
		{ID: primitive.NewObjectID(), Username: "reader_a", Nickname: "读者A", Avatar: "/images/a.png"},
		{ID: primitive.NewObjectID(), Username: "reader_b", Nickname: "读者B", Avatar: "/images/b.png"},
	}
	books := []socialSeedBook{
		{
			ID:           primitive.NewObjectID(),
			Title:        "星海回声",
			Author:       "作者甲",
			Introduction: "关于远航与回响的故事。",
			Cover:        "/images/covers/a.jpg",
			Categories:   []string{"科幻"},
			Tags:         []string{"太空", "冒险"},
		},
		{
			ID:           primitive.NewObjectID(),
			Title:        "雾隐都市",
			Author:       "作者乙",
			Introduction: "悬疑氛围拉满。",
			Cover:        "/images/covers/b.jpg",
			Categories:   []string{"悬疑"},
			Tags:         []string{"烧脑"},
		},
	}

	posts := buildSeedPosts(users, books)
	if len(posts) != len(books) {
		t.Fatalf("expected %d seeded posts, got %d", len(books), len(posts))
	}

	first, ok := posts[0].(bson.M)
	if !ok {
		t.Fatalf("expected seeded post to be bson.M, got %T", posts[0])
	}
	if first["book_id"] == "" || first["book_title"] == "" {
		t.Fatalf("expected seeded post to keep book binding, got %+v", first)
	}
	topics, ok := first["topics"].([]string)
	if !ok {
		t.Fatalf("expected topics to be []string, got %T", first["topics"])
	}
	if len(topics) < 2 || topics[0] != "推荐" {
		t.Fatalf("expected topics to start with 推荐 and include book metadata, got %+v", topics)
	}
}

func TestBuildPostTopics_DeduplicatesAndKeepsCategories(t *testing.T) {
	topics := buildPostTopics(socialSeedBook{
		Categories: []string{"玄幻", "玄幻"},
		Tags:       []string{"热血", "热血", "冒险"},
	})

	if len(topics) < 3 {
		t.Fatalf("expected multiple topics, got %+v", topics)
	}
	if topics[0] != "推荐" {
		t.Fatalf("expected 推荐 as first topic, got %+v", topics)
	}
	if topics[1] != "玄幻" {
		t.Fatalf("expected category topic to be preserved, got %+v", topics)
	}
}
