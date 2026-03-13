package main

import (
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/models"
	bookstoreModel "Qingyu_backend/models/bookstore"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const showcaseTag = "seed:showcase"

type showcaseBookSpec struct {
	Title         string
	Author        string
	Category      string
	Introduction  string
	Cover         string
	Tags          []string
	Status        string
	Rating        float64
	RatingCount   int64
	ViewCount     int64
	WordCount     int64
	ChapterCount  int
	Price         float64
	IsFree        bool
	IsRecommended bool
	IsFeatured    bool
	IsHot         bool
	PublishedDays int
	LastUpdateAgo int
}

var showcaseBookSpecs = []showcaseBookSpec{
	{
		Title:         "云海问剑录",
		Author:        "闻舟",
		Category:      "仙侠",
		Introduction:  "少年从废弃山门拾得残卷，自此被卷入旧宗门、海上剑墟与天外秘境的连锁风暴。节奏稳定，适合做首页推荐和追更演示。",
		Cover:         "/images/covers/showcase-yunhai.jpg",
		Tags:          []string{"修仙", "剑修", "宗门", "成长", "演示精选", showcaseTag},
		Status:        "ongoing",
		Rating:        9.3,
		RatingCount:   12860,
		ViewCount:     1680000,
		WordCount:     1280000,
		ChapterCount:  386,
		Price:         29.9,
		IsFree:        false,
		IsRecommended: true,
		IsFeatured:    true,
		IsHot:         true,
		PublishedDays: 240,
		LastUpdateAgo: 1,
	},
	{
		Title:         "长安雪尽时",
		Author:        "阿迟",
		Category:      "历史",
		Introduction:  "以边城小吏的抉择切入朝局更替，兼具战争调度、人物群像与长线伏笔，适合作为完结高分作品演示。",
		Cover:         "/images/covers/showcase-changan.jpg",
		Tags:          []string{"权谋", "战争", "群像", "完结佳作", "演示精选", showcaseTag},
		Status:        "completed",
		Rating:        9.6,
		RatingCount:   20480,
		ViewCount:     2260000,
		WordCount:     1560000,
		ChapterCount:  428,
		Price:         36.0,
		IsFree:        false,
		IsRecommended: true,
		IsFeatured:    true,
		IsHot:         true,
		PublishedDays: 420,
		LastUpdateAgo: 8,
	},
	{
		Title:         "霓虹停机坪",
		Author:        "林见山",
		Category:      "都市",
		Introduction:  "退役试飞员回到海港城经营民间机库，在资本、旧案和家族债务之间重新起飞。都市职业线清晰，适合详情页演示。",
		Cover:         "/images/covers/showcase-nihong.jpg",
		Tags:          []string{"都市", "职业", "悬疑", "逆袭", "演示精选", showcaseTag},
		Status:        "ongoing",
		Rating:        8.9,
		RatingCount:   7640,
		ViewCount:     980000,
		WordCount:     820000,
		ChapterCount:  214,
		Price:         19.9,
		IsFree:        false,
		IsRecommended: true,
		IsFeatured:    false,
		IsHot:         true,
		PublishedDays: 160,
		LastUpdateAgo: 0,
	},
	{
		Title:         "夜航星尘档案",
		Author:        "纪遥",
		Category:      "科幻",
		Introduction:  "在失重港口负责打捞黑匣子的调查员，逐步拼出殖民舰队失踪真相。题材辨识度高，适合作为榜单头部演示。",
		Cover:         "/images/covers/showcase-yehang.jpg",
		Tags:          []string{"星际", "调查", "赛博", "悬疑", "演示精选", showcaseTag},
		Status:        "ongoing",
		Rating:        9.1,
		RatingCount:   9320,
		ViewCount:     1320000,
		WordCount:     940000,
		ChapterCount:  267,
		Price:         24.0,
		IsFree:        false,
		IsRecommended: true,
		IsFeatured:    true,
		IsHot:         true,
		PublishedDays: 95,
		LastUpdateAgo: 1,
	},
	{
		Title:         "旧城游戏策展人",
		Author:        "孟白",
		Category:      "游戏",
		Introduction:  "博物馆策展人与独立游戏工作室联手，把历史叙事做成沉浸式副本。题材轻新，适合新人榜和发现页演示。",
		Cover:         "/images/covers/showcase-youxi.jpg",
		Tags:          []string{"游戏", "副本", "策展", "轻小说感", "演示精选", showcaseTag},
		Status:        "ongoing",
		Rating:        8.7,
		RatingCount:   2860,
		ViewCount:     420000,
		WordCount:     286000,
		ChapterCount:  72,
		Price:         0,
		IsFree:        true,
		IsRecommended: true,
		IsFeatured:    false,
		IsHot:         false,
		PublishedDays: 18,
		LastUpdateAgo: 0,
	},
}

func buildShowcaseBooks(
	categories map[string]*bookstoreModel.Category,
	authorIDs []primitive.ObjectID,
	limit int,
) ([]models.Book, error) {
	if limit <= 0 {
		return nil, nil
	}
	count := limit
	if count > len(showcaseBookSpecs) {
		count = len(showcaseBookSpecs)
	}
	if len(authorIDs) == 0 {
		return nil, fmt.Errorf("没有可用作者用于生成精选书籍")
	}

	now := time.Now()
	books := make([]models.Book, 0, count)
	for i := 0; i < count; i++ {
		spec := showcaseBookSpecs[i]
		category, ok := categories[spec.Category]
		if !ok {
			return nil, fmt.Errorf("精选书籍缺少分类: %s", spec.Category)
		}
		publishedAt := now.Add(-time.Duration(spec.PublishedDays) * 24 * time.Hour)
		lastUpdateAt := now.Add(-time.Duration(spec.LastUpdateAgo) * 24 * time.Hour)
		if lastUpdateAt.Before(publishedAt) {
			lastUpdateAt = publishedAt
		}

		books = append(books, models.Book{
			ID:            primitive.NewObjectID(),
			Title:         spec.Title,
			Author:        spec.Author,
			AuthorID:      authorIDs[i%len(authorIDs)],
			Introduction:  spec.Introduction,
			Cover:         spec.Cover,
			CategoryIDs:   []primitive.ObjectID{category.ID},
			Categories:    []string{category.Name},
			Tags:          spec.Tags,
			Status:        spec.Status,
			Rating:        spec.Rating,
			RatingCount:   spec.RatingCount,
			ViewCount:     spec.ViewCount,
			WordCount:     spec.WordCount,
			ChapterCount:  spec.ChapterCount,
			Price:         spec.Price,
			IsFree:        spec.IsFree,
			IsRecommended: spec.IsRecommended,
			IsFeatured:    spec.IsFeatured,
			IsHot:         spec.IsHot,
			PublishedAt:   publishedAt,
			LastUpdateAt:  lastUpdateAt,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}

	return books, nil
}
