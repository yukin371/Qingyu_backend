package content

import (
	"context"

	"Qingyu_backend/models/dto"
)

// ============================================================================
// 章节服务接口详细定义
// ============================================================================

// ChapterService 章节内容服务接口
// 提供章节内容的获取和导航功能
type ChapterService interface {
	// === 章节内容获取 ===

	// GetChapter 获取章节内容
	// 返回指定章节的完整内容，包括标题、正文、字数等信息
	GetChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error)

	// GetChapterByNumber 根据章节号获取章节
	// 通过章节序号获取章节内容，支持快速定位
	GetChapterByNumber(ctx context.Context, bookID string, chapterNum int) (*dto.ChapterResponse, error)

	// GetChapterInfo 获取章节信息
	// 返回章节的基本信息（标题、序号、字数等），不包含正文内容
	GetChapterInfo(ctx context.Context, chapterID string) (*dto.ChapterInfo, error)

	// === 章节导航 ===

	// GetNextChapter 获取下一章
	// 返回当前章节的下一章内容，用于"下一章"导航
	GetNextChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error)

	// GetPreviousChapter 获取上一章
	// 返回当前章节的上一章内容，用于"上一章"导航
	GetPreviousChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error)

	// ListChapters 获取章节列表
	// 返回书籍的所有章节列表，用于目录展示
	ListChapters(ctx context.Context, bookID string) (*dto.ChapterListResponse, error)

	// === 章节搜索与筛选 ===

	// SearchChapters 搜索章节
	// 根据关键词搜索章节标题和内容
	SearchChapters(ctx context.Context, bookID string, keyword string, page, pageSize int) (*dto.ChapterListResponse, error)

	// GetChaptersByType 根据类型获取章节
	// 按章节类型（正传、番外等）筛选章节
	GetChaptersByType(ctx context.Context, bookID, chapterType string) (*dto.ChapterListResponse, error)

	// === 章节状态管理 ===

	// GetChapterPublishStatus 获取章节发布状态
	// 返回章节的发布状态、发布时间等信息
	GetChapterPublishStatus(ctx context.Context, chapterID string) (*dto.ChapterPublishStatus, error)

	// UpdateChapterPublishStatus 更新章节发布状态
	// 更新章节的发布状态（已发布/草稿/下架）
	UpdateChapterPublishStatus(ctx context.Context, chapterID string, req *dto.UpdateChapterPublishStatusRequest) error

	// === 批量操作 ===

	// BatchGetChapters 批量获取章节
	// 批量获取多个章节的内容，用于离线阅读等场景
	BatchGetChapters(ctx context.Context, chapterIDs []string) ([]*dto.ChapterResponse, error)

	// GetChapterRange 获取章节范围
	// 获取指定章节号范围内的所有章节
	GetChapterRange(ctx context.Context, bookID string, startNum, endNum int) (*dto.ChapterListResponse, error)
}
