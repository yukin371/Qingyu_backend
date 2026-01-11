package interfaces

import (
	"context"
	bookstoreModel "Qingyu_backend/models/bookstore"
)

// BookstoreService 书城服务接口 - 专注于书城列表展示和首页聚合
// 用于书城首页、分类页面、搜索结果等列表场景
type BookstoreService interface {
	// 书籍列表相关服务 - 使用Book模型
	GetAllBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	GetBookByID(ctx context.Context, id string) (*bookstoreModel.Book, error)
	GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	GetBooksByAuthorID(ctx context.Context, authorID string, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	GetHotBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	GetNewReleases(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	GetFreeBooks(ctx context.Context, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstoreModel.Book, int64, error)
	SearchBooksWithFilter(ctx context.Context, filter *bookstoreModel.BookFilter) ([]*bookstoreModel.Book, int64, error)

	// 分类相关服务
	GetCategoryTree(ctx context.Context) ([]*bookstoreModel.CategoryTree, error)
	GetCategoryByID(ctx context.Context, id string) (*bookstoreModel.Category, error)
	GetRootCategories(ctx context.Context) ([]*bookstoreModel.Category, error)

	// Banner相关方法
	GetActiveBanners(ctx context.Context, limit int) ([]*bookstoreModel.Banner, error)
	IncrementBannerClick(ctx context.Context, bannerID string) error

	// 榜单相关方法
	GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstoreModel.RankingItem, error)
	GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error)
	GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error)
	GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstoreModel.RankingItem, error)
	GetRankingByType(ctx context.Context, rankingType bookstoreModel.RankingType, period string, limit int) ([]*bookstoreModel.RankingItem, error)
	UpdateRankings(ctx context.Context, rankingType bookstoreModel.RankingType, period string) error

	// 首页数据聚合
	GetHomepageData(ctx context.Context) (interface{}, error)

	// 统计和计数
	GetBookStats(ctx context.Context) (*bookstoreModel.BookStats, error)
	IncrementBookView(ctx context.Context, bookID string) error
}

// ChapterService 章节服务接口
type ChapterService interface {
	// 章节查询
	GetChapterByID(ctx context.Context, chapterID interface{}) (interface{}, error)
	GetChaptersByBookID(ctx context.Context, bookID interface{}, page, size int) (interface{}, int64, error)

	// 章节内容
	GetChapterContent(ctx context.Context, chapterID interface{}, userID interface{}) (string, error)

	// 章节购买相关
	GetChapterPrice(ctx context.Context, chapterID interface{}) (int64, error)
	IsChapterPurchased(ctx context.Context, userID, chapterID string) (bool, error)
	PurchaseChapter(ctx context.Context, userID, chapterID string) error

	// 章节列表
	GetChapterList(ctx context.Context, bookID string, page, size int) (interface{}, int64, error)
}
