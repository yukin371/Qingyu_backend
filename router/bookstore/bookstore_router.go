package bookstore

import (
	bookstoreApi "Qingyu_backend/api/v1/bookstore"
	"Qingyu_backend/internal/middleware/auth"
	"Qingyu_backend/pkg/logger"
	"Qingyu_backend/service/bookstore"
	searchService "Qingyu_backend/service/search"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// =====================================================
// 书店路由配置文档
// =====================================================
//
//  路由设计原则：
//
// 1️ 公开路由 (public) - 无需认证
//    - 适用于首页数据、浏览、搜索、排行榜等内容消费场景
//    - 可被任何客户端（已登录或未登录）访问
//
// 2️ 认证路由 (authenticated) - 需要JWT Token
//    - 适用于用户个人数据、行为追踪、点赞评论等需关联用户身份的场景
//    - 需要用户提供有效的JWT Token
//
//  具体划分：
//
//  公开 (无需登录)
//  - GET /homepage               : 获取首页数据
//  - GET /books/*                : 书籍信息查询
//  - GET /categories/*           : 分类信息查询
//  - GET /rankings/*             : 榜单查询
//  - GET /banners                : 获取可用Banner
//  - POST /banners/:id/click     : 👈 Banner点击记录（匿名可用）
//
//  认证 (需要登录)
//  - POST /books/:id/view        : 书籍点击记录（关联用户）
//  - POST /ratings/*             : 评分、评论等（关联用户）
//
//  为什么这样设计：
//  - Banner点击是**广告统计**，不需要关联用户身份
//  - 书籍点击是**用户行为数据**，用于个性化推荐
//  - 这种设计让前端在登录前就能完全使用首页和浏览功能

// InitBookstoreRouter 初始化书店路由
func InitBookstoreRouter(
	r *gin.RouterGroup,
	bookstoreService bookstore.BookstoreService,
	bookDetailService bookstore.BookDetailService,
	ratingService bookstore.BookRatingService,
	statisticsService bookstore.BookStatisticsService,
	chapterService bookstore.ChapterService,
	purchaseService bookstore.ChapterPurchaseService,
	searchSvc *searchService.SearchService,
	zapLogger *zap.Logger,
) {
	// 创建API实例
	bookstoreApiHandler := bookstoreApi.NewBookstoreAPI(bookstoreService, searchSvc, logger.Get())

	// 初始化其他服务的API处理器
	var bookDetailApiHandler *bookstoreApi.BookDetailAPI
	if bookstoreService != nil {
		bookDetailApiHandler = bookstoreApi.NewBookDetailAPI(bookstoreService)
	}

	// 初始化Rating API处理器
	var ratingApiHandler *bookstoreApi.BookRatingAPI
	if ratingService != nil {
		ratingApiHandler = bookstoreApi.NewBookRatingAPI(ratingService)
	}

	// 初始化Chapter API处理器
	var chapterApiHandler *bookstoreApi.ChapterAPI
	if chapterService != nil {
		chapterApiHandler = bookstoreApi.NewChapterAPI(chapterService)
	}

	// 初始化Chapter Catalog API处理器（章节目录和购买）
	var chapterCatalogApiHandler *bookstoreApi.ChapterCatalogAPI
	if chapterService != nil && purchaseService != nil {
		chapterCatalogApiHandler = bookstoreApi.NewChapterCatalogAPI(chapterService, purchaseService)
	}

	// ℹ️ Statistics API已通过BookDetailAPI实现
	// 如需单独的Statistics API处理器，可在这里初始化
	// var statisticsApiHandler *bookstoreApi.BookStatisticsAPI
	// if statisticsService != nil {
	// 	statisticsApiHandler = bookstoreApi.NewBookStatisticsAPI(statisticsService)
	// }
	// chapterApiHandler := bookstoreApi.NewChapterAPI(...)

	// 书店主路由组
	bookstoreGroup := r.Group("/bookstore")
	{
		// 公开接口（不需要认证）
		public := bookstoreGroup.Group("")
		{
			// 书城首页
			public.GET("/homepage", bookstoreApiHandler.GetHomepage)

			// 书籍列表和搜索 - 注意：具体路由必须放在参数化路由之前
			public.GET("/books", bookstoreApiHandler.GetBooks) // 必须放在 /books/:id 之前
			public.GET("/books/search", bookstoreApiHandler.SearchBooks)
			public.GET("/books/stream", bookstoreApiHandler.StreamSearchBooks)            // 新增：流式搜索
			public.GET("/books/stream-batch", bookstoreApiHandler.StreamSearchBooksBatch) // 新增：批量流式搜索
			public.GET("/books/search/title", bookstoreApiHandler.SearchByTitle)          // 新增：按标题搜索
			public.GET("/books/search/author", bookstoreApiHandler.SearchByAuthor)        // 新增：按作者搜索
			public.GET("/books/recommended", bookstoreApiHandler.GetRecommendedBooks)
			public.GET("/books/featured", bookstoreApiHandler.GetFeaturedBooks)
			public.GET("/books/tags", bookstoreApiHandler.GetBooksByTags)     // 新增：按标签筛选
			public.GET("/books/status", bookstoreApiHandler.GetBooksByStatus) // 新增：按状态筛选
			public.GET("/books/years", bookstoreApiHandler.GetYears)          // 新增：获取年份列表（必须在 /books/:id 之前）

			// 书籍详情接口（当BookDetailAPI可用时）- 必须在 /books/:id 之前注册
			if bookDetailApiHandler != nil {
				public.GET("/books/:id/detail", bookDetailApiHandler.GetBookDetail)
				// ✅ 浏览量记录API（公开接口，不需要认证）
				public.POST("/books/:id/view", bookDetailApiHandler.IncrementViewCount)
			}

			// /books/:id 路由必须放在所有更具体的 /books/:id/* 路由之后
			public.GET("/books/:id", bookstoreApiHandler.GetBookByID)
			public.GET("/books/:id/similar", bookstoreApiHandler.GetSimilarBooks) // 新增：相似书籍推荐（四层降级策略）

			// 分类 - 注意：具体路由必须放在参数化路由之前
			public.GET("/categories/tree", bookstoreApiHandler.GetCategoryTree)
			public.GET("/categories/:id/books", bookstoreApiHandler.GetBooksByCategory)
			public.GET("/categories/:id", bookstoreApiHandler.GetCategoryByID)

			// 元数据查询 - 用于筛选
			public.GET("/tags", bookstoreApiHandler.GetTags) // 新增：获取标签列表

			// Banner - 公开API
			public.GET("/banners", bookstoreApiHandler.GetActiveBanners)
			// ✅ Banner 点击记录（公开，不需要认证）
			public.POST("/banners/:id/click", bookstoreApiHandler.IncrementBannerClick)

			// 排行榜
			public.GET("/rankings/realtime", bookstoreApiHandler.GetRealtimeRanking)
			public.GET("/rankings/weekly", bookstoreApiHandler.GetWeeklyRanking)
			public.GET("/rankings/monthly", bookstoreApiHandler.GetMonthlyRanking)
			public.GET("/rankings/newbie", bookstoreApiHandler.GetNewbieRanking)
			public.GET("/rankings/:type", bookstoreApiHandler.GetRankingByType)

			// ✅ 统计API（当StatisticsService可用时）
			// 注意：BookDetailAPI中已包含GetBookStatistics
			// Statistics API在BookDetail中已实现，这里备注即可

			// ✅ Chapter API（章节管理）- 公开接口（优先注册，不需要PurchaseService）
			if chapterApiHandler != nil {
				public.GET("/books/:id/chapters", chapterApiHandler.GetChaptersByBookID)      // 获取章节列表（使用ChapterAPI）
				public.GET("/books/:id/chapters/list", chapterApiHandler.GetChaptersByBookID) // 别名路由
				public.GET("/chapters/:id", chapterApiHandler.GetChapter)                     // 获取章节详情
				public.GET("/chapters/:id/content", chapterApiHandler.GetChapterContent)      // 获取章节内容（公开接口）
			}

			// ✅ Chapter Catalog API（章节目录和购买）- 需要ChapterService和PurchaseService
			if chapterCatalogApiHandler != nil {
				// 这些路由需要购买服务支持（VIP章节、价格、购买记录等）
				public.GET("/books/:id/trial-chapters", chapterCatalogApiHandler.GetTrialChapters)     // 获取试读章节
				public.GET("/books/:id/vip-chapters", chapterCatalogApiHandler.GetVIPChapters)         // 获取VIP章节列表
				public.GET("/chapters/:chapterId/price", chapterCatalogApiHandler.GetChapterPrice)     // 获取章节价格
				public.GET("/chapters/:chapterId/access", chapterCatalogApiHandler.CheckChapterAccess) // 检查章节访问权限
			}
		}
	}

	// 需要认证的接口
	authenticated := bookstoreGroup.Group("")
	authenticated.Use(auth.JWTAuth())
	{
		// ✅ 评分API（当RatingAPI可用时）
		if ratingApiHandler != nil {
			authenticated.GET("/books/:id/rating", ratingApiHandler.GetBookRating)
			authenticated.POST("/books/:id/rating", ratingApiHandler.CreateRating)
			authenticated.PUT("/books/:id/rating", ratingApiHandler.UpdateRating)
			authenticated.DELETE("/books/:id/rating", ratingApiHandler.DeleteRating)
			authenticated.GET("/ratings/user/:id", ratingApiHandler.GetRatingsByUserID)
		}
	}
}

// InitReaderPurchaseRouter 初始化读者购买路由（用于章节购买相关接口）
// 这些接口放在 /api/v1/reader 路径下，因为它们与读者的个人购买记录相关
func InitReaderPurchaseRouter(
	r *gin.RouterGroup,
	purchaseService bookstore.ChapterPurchaseService,
) {
	// 如果没有提供购买服务，直接返回
	if purchaseService == nil {
		return
	}

	readerGroup := r.Group("/reader")
	readerGroup.Use(auth.JWTAuth())
	{
		// 创建章节目录API处理器
		chapterCatalogApiHandler := bookstoreApi.NewChapterCatalogAPI(nil, purchaseService)

		// ✅ 购买相关接口（需要认证）
		readerGroup.POST("/chapters/:chapterId/purchase", chapterCatalogApiHandler.PurchaseChapter) // 购买单个章节
		readerGroup.POST("/books/:id/buy-all", chapterCatalogApiHandler.PurchaseBook)               // 购买全书

		// ✅ 购买记录查询（需要认证）
		readerGroup.GET("/purchases", chapterCatalogApiHandler.GetPurchases)         // 获取所有购买记录
		readerGroup.GET("/purchases/:id", chapterCatalogApiHandler.GetBookPurchases) // 获取某本书的购买记录
	}
}
