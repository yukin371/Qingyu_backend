package impl

import (
	"context"

	readerif "Qingyu_backend/service/interfaces/reader"
	readerService "Qingyu_backend/service/reader"
)

// ChapterContentImpl 章节内容获取端口实现
type ChapterContentImpl struct {
	readerService  *readerService.ReaderService
	chapterService readerService.ChapterService
	serviceName    string
	version        string
}

// NewChapterContentImpl 创建章节内容获取端口实现
func NewChapterContentImpl(
	readerService *readerService.ReaderService,
	chapterService readerService.ChapterService,
) readerif.ChapterContentPort {
	return &ChapterContentImpl{
		readerService:  readerService,
		chapterService: chapterService,
		serviceName:    "ChapterContentPort",
		version:        "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (c *ChapterContentImpl) Initialize(ctx context.Context) error {
	return c.readerService.Initialize(ctx)
}

func (c *ChapterContentImpl) Health(ctx context.Context) error {
	return c.readerService.Health(ctx)
}

func (c *ChapterContentImpl) Close(ctx context.Context) error {
	return c.readerService.Close(ctx)
}

func (c *ChapterContentImpl) GetServiceName() string {
	return c.serviceName
}

func (c *ChapterContentImpl) GetVersion() string {
	return c.version
}

// ============================================================================
// ChapterContentPort 方法实现
// ============================================================================

// GetChapterContent 获取章节内容（基本）
func (c *ChapterContentImpl) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	return c.readerService.GetChapterContent(ctx, userID, chapterID)
}

// GetChapterByID 获取章节信息（不含内容）
func (c *ChapterContentImpl) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	return c.readerService.GetChapterByID(ctx, chapterID)
}

// GetBookChapters 获取书籍的章节列表
func (c *ChapterContentImpl) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	return c.readerService.GetBookChapters(ctx, bookID, page, size)
}

// GetChapterContentWithProgress 获取章节内容（带权限检查、进度保存）
func (c *ChapterContentImpl) GetChapterContentWithProgress(ctx context.Context, req *readerif.GetChapterContentRequest) (*readerif.ChapterContentResponse, error) {
	// 委托给 ChapterService，并进行类型转换
	resp, err := c.chapterService.GetChapterContent(ctx, req.UserID, req.BookID, req.ChapterID)
	if err != nil {
		return nil, err
	}

	// 转换响应类型
	return &readerif.ChapterContentResponse{
		ChapterID:   resp.ChapterID,
		BookID:      resp.BookID,
		ChapterNum:  resp.ChapterNum,
		Title:       resp.Title,
		Content:     resp.Content,
		WordCount:   resp.WordCount,
		IsVIP:       !resp.CanAccess, // 如果不能访问则为VIP章节
		PublishedAt: resp.LastReadAt.Unix(), // 使用 LastReadAt 作为近似值
	}, nil
}

// GetChapterByNumber 根据章节号获取内容
func (c *ChapterContentImpl) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (*readerif.ChapterContentResponse, error) {
	// 委托给 ChapterService，并进行类型转换
	resp, err := c.chapterService.GetChapterByNumber(ctx, userID, bookID, chapterNum)
	if err != nil {
		return nil, err
	}

	// 转换响应类型
	return &readerif.ChapterContentResponse{
		ChapterID:   resp.ChapterID,
		BookID:      resp.BookID,
		ChapterNum:  resp.ChapterNum,
		Title:       resp.Title,
		Content:     resp.Content,
		WordCount:   resp.WordCount,
		IsVIP:       !resp.CanAccess,
		PublishedAt: resp.LastReadAt.Unix(),
	}, nil
}

// GetNextChapter 获取下一章
func (c *ChapterContentImpl) GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (*readerif.ChapterInfo, error) {
	// 委托给 ChapterService，并进行类型转换
	info, err := c.chapterService.GetNextChapter(ctx, userID, bookID, chapterID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	// 转换响应类型
	return &readerif.ChapterInfo{
		ChapterID:   info.ChapterID,
		BookID:      info.BookID,
		ChapterNum:  info.ChapterNum,
		Title:       info.Title,
		WordCount:   info.WordCount,
		IsVIP:       !info.CanAccess,
		PublishedAt: info.PublishTime.Unix(),
	}, nil
}

// GetPreviousChapter 获取上一章
func (c *ChapterContentImpl) GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (*readerif.ChapterInfo, error) {
	// 委托给 ChapterService，并进行类型转换
	info, err := c.chapterService.GetPreviousChapter(ctx, userID, bookID, chapterID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	// 转换响应类型
	return &readerif.ChapterInfo{
		ChapterID:   info.ChapterID,
		BookID:      info.BookID,
		ChapterNum:  info.ChapterNum,
		Title:       info.Title,
		WordCount:   info.WordCount,
		IsVIP:       !info.CanAccess,
		PublishedAt: info.PublishTime.Unix(),
	}, nil
}

// GetChapterList 获取章节目录
func (c *ChapterContentImpl) GetChapterList(ctx context.Context, userID, bookID string, page, size int) (*readerif.ChapterListResponse, error) {
	// 委托给 ChapterService，并进行类型转换
	resp, err := c.chapterService.GetChapterList(ctx, userID, bookID, page, size)
	if err != nil {
		return nil, err
	}

	// 转换章节列表
	chapters := make([]*readerif.ChapterInfo, 0, len(resp.Chapters))
	for _, ch := range resp.Chapters {
		chapters = append(chapters, &readerif.ChapterInfo{
			ChapterID:   ch.ChapterID,
			BookID:      ch.BookID,
			ChapterNum:  ch.ChapterNum,
			Title:       ch.Title,
			WordCount:   ch.WordCount,
			IsVIP:       !ch.CanAccess,
			PublishedAt: ch.PublishTime.Unix(),
		})
	}

	// 转换响应类型
	return &readerif.ChapterListResponse{
		Chapters: chapters,
		Total:    resp.Total,
		Page:     resp.Page,
		Size:     resp.Size,
	}, nil
}

// GetChapterInfo 获取章节信息（不含内容）
func (c *ChapterContentImpl) GetChapterInfo(ctx context.Context, userID, chapterID string) (*readerif.ChapterInfo, error) {
	// 委托给 ChapterService，并进行类型转换
	info, err := c.chapterService.GetChapterInfo(ctx, userID, chapterID)
	if err != nil {
		return nil, err
	}

	// 转换响应类型
	return &readerif.ChapterInfo{
		ChapterID:   info.ChapterID,
		BookID:      info.BookID,
		ChapterNum:  info.ChapterNum,
		Title:       info.Title,
		WordCount:   info.WordCount,
		IsVIP:       !info.CanAccess,
		PublishedAt: info.PublishTime.Unix(),
	}, nil
}
