// Package generators 提供阅读数据生成功能
package generators

import (
	"math/rand"
	"time"

	"Qingyu_backend/cmd/seeder/models"
)

// ReaderGenerator 阅读数据生成器
type ReaderGenerator struct {
	base *BaseGenerator
}

// NewReaderGenerator 创建阅读数据生成器
func NewReaderGenerator() *ReaderGenerator {
	return &ReaderGenerator{
		base: NewBaseGenerator(),
	}
}

// GenerateReadingHistories 生成阅读历史
func (g *ReaderGenerator) GenerateReadingHistories(userIDs, bookIDs, chapterIDs []string, count int) []models.ReadingHistory {
	histories := make([]models.ReadingHistory, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		bookID := bookIDs[rand.Intn(len(bookIDs))]
		chapterID := chapterIDs[rand.Intn(len(chapterIDs))]

		histories[i] = models.ReadingHistory{
			ID:        g.base.ID(),
			UserID:    userID,
			BookID:    bookID,
			ChapterID: chapterID,
			ReadTime:  now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
			Duration:  300 + rand.Intn(3600), // 5分钟到1小时
			Device:    []string{"mobile", "tablet", "desktop"}[rand.Intn(3)],
			CreatedAt: now,
		}
	}

	return histories
}

// GenerateReadingProgresses 生成阅读进度
func (g *ReaderGenerator) GenerateReadingProgresses(userIDs, bookIDs []string, count int) []models.ReadingProgress {
	progresses := make([]models.ReadingProgress, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		bookID := bookIDs[rand.Intn(len(bookIDs))]

		// 进度分布
		progressRand := rand.Float64()
		var progress float64
		var chapterNum int

		switch {
		case progressRand < 0.3:
			progress = float64(5 + rand.Intn(15)) // 5-20%
			chapterNum = 1 + rand.Intn(10)
		case progressRand < 0.7:
			progress = float64(20 + rand.Intn(60)) // 20-80%
			chapterNum = 10 + rand.Intn(50)
		default:
			progress = float64(80 + rand.Intn(20)) // 80-100%
			chapterNum = 50 + rand.Intn(100)
		}

		progresses[i] = models.ReadingProgress{
			ID:            g.base.ID(),
			UserID:        userID,
			BookID:        bookID,
			ChapterNum:    chapterNum,
			Progress:      progress,
			LastReadAt:    now.Add(-time.Duration(rand.Intn(48)) * time.Hour),
			TotalReadTime: 1800 + rand.Intn(72000), // 30分钟到20小时
			CreatedAt:     now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
			UpdatedAt:     now,
		}
	}

	return progresses
}

// GenerateReadingProgress 为指定用户和书籍生成单个阅读进度
func (g *ReaderGenerator) GenerateReadingProgress(userID, bookID string) models.ReadingProgress {
	now := time.Now()

	// 进度分布
	progressRand := rand.Float64()
	var progress float64
	var chapterNum int

	switch {
	case progressRand < 0.3:
		progress = float64(5 + rand.Intn(15)) // 5-20%
		chapterNum = 1 + rand.Intn(10)
	case progressRand < 0.7:
		progress = float64(20 + rand.Intn(60)) // 20-80%
		chapterNum = 10 + rand.Intn(50)
	default:
		progress = float64(80 + rand.Intn(20)) // 80-100%
		chapterNum = 50 + rand.Intn(100)
	}

	return models.ReadingProgress{
		ID:            g.base.ID(),
		UserID:        userID,
		BookID:        bookID,
		ChapterNum:    chapterNum,
		Progress:      progress,
		LastReadAt:    now.Add(-time.Duration(rand.Intn(48)) * time.Hour),
		TotalReadTime: 1800 + rand.Intn(72000), // 30分钟到20小时
		CreatedAt:     now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
		UpdatedAt:     now,
	}
}

// GenerateBookmarks 生成书签
func (g *ReaderGenerator) GenerateBookmarks(userIDs, bookIDs, chapterIDs []string, count int) []models.Bookmark {
	bookmarks := make([]models.Bookmark, count)
	now := time.Now()

	notes := []string{
		"精彩段落", "重要情节", "值得回味", "描写细腻",
		"高潮部分", "转折点", "伏笔", "人物塑造",
	}

	for i := 0; i < count; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		bookID := bookIDs[rand.Intn(len(bookIDs))]
		chapterID := chapterIDs[rand.Intn(len(chapterIDs))]

		bookmarks[i] = models.Bookmark{
			ID:        g.base.ID(),
			UserID:    userID,
			BookID:    bookID,
			ChapterID: chapterID,
			Position:  rand.Intn(5000),
			Note:      notes[rand.Intn(len(notes))],
			CreatedAt: now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
		}
	}

	return bookmarks
}

// GenerateAnnotations 生成批注
func (g *ReaderGenerator) GenerateAnnotations(userIDs, bookIDs, chapterIDs []string, count int) []models.Annotation {
	annotations := make([]models.Annotation, count)
	now := time.Now()

	contents := []string{
		"这里描写得很生动", "人物刻画深刻", "情节设计巧妙",
		"对话自然流畅", "环境描写细腻", "节奏把握很好",
		"伏笔埋得很深", "这个转折很意外",
	}

	colors := []string{"yellow", "green", "blue", "pink", "orange"}

	for i := 0; i < count; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		bookID := bookIDs[rand.Intn(len(bookIDs))]
		chapterID := chapterIDs[rand.Intn(len(chapterIDs))]

		annotations[i] = models.Annotation{
			ID:        g.base.ID(),
			UserID:    userID,
			BookID:    bookID,
			ChapterID: chapterID,
			Position:  rand.Intn(5000),
			Content:   contents[rand.Intn(len(contents))],
			Color:     colors[rand.Intn(len(colors))],
			CreatedAt: now.Add(-time.Duration(rand.Intn(720)) * time.Hour),
		}
	}

	return annotations
}
