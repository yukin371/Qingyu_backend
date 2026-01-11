package interfaces

import (
	"context"

	readerservice "Qingyu_backend/service/reader"
)

// ReaderChapterService 阅读器章节服务接口
type ReaderChapterService interface {
	// GetChapterContent 获取章节内容（带权限检查、进度保存）
	GetChapterContent(ctx context.Context, userID, bookID, chapterID string) (*readerservice.ChapterContentResponse, error)

	// GetChapterByNumber 根据章节号获取内容
	GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*readerservice.ChapterContentResponse, error)

	// GetNextChapter 获取下一章
	GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*readerservice.ChapterInfo, error)

	// GetPreviousChapter 获取上一章
	GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*readerservice.ChapterInfo, error)

	// GetChapterList 获取章节目录
	GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*readerservice.ChapterListResponse, error)

	// GetChapterInfo 获取章节信息（不含内容）
	GetChapterInfo(ctx context.Context, userID, chapterID string) (*readerservice.ChapterInfo, error)
}
