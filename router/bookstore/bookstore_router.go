package bookstore

import (
	bookstoreApi "Qingyu_backend/api/v1/bookstore"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/bookstore"

	"github.com/gin-gonic/gin"
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
) {
	// 创建API实例
	bookstoreApiHandler := bookstoreApi.NewBookstoreAPI(bookstoreService)

	// 初始化其他服务的API处理器
	var bookDetailApiHandler *bookstoreApi.BookDetailAPI
	if bookDetailService != nil {
		bookDetailApiHandler = bookstoreApi.NewBookDetailAPI(bookDetailService)
	}

	// 初始化Rating API处理器
	var ratingApiHandler *bookstoreApi.BookRatingAPI
	if ratingService != nil {
		ratingApiHandler = bookstoreApi.NewBookRatingAPI(ratingService)
	}

	// 初始化Chapter API处理器（暂时跳过，需要ChapterService）
	// var chapterApiHandler *bookstoreApi.ChapterAPI
	// if chapterService != nil {
	// 	chapterApiHandler = bookstoreApi.NewChapterAPI(chapterService)
	// }

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
			public.GET("/books/search", bookstoreApiHandler.SearchBooks)
			public.GET("/books/recommended", bookstoreApiHandler.GetRecommendedBooks)
			public.GET("/books/featured", bookstoreApiHandler.GetFeaturedBooks)
			public.GET("/books/:id", bookstoreApiHandler.GetBookByID)

			// 分类 - 注意：具体路由必须放在参数化路由之前
			public.GET("/categories/tree", bookstoreApiHandler.GetCategoryTree)
			public.GET("/categories/:id/books", bookstoreApiHandler.GetBooksByCategory)
			public.GET("/categories/:id", bookstoreApiHandler.GetCategoryByID)

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

			// 书籍详情接口（当BookDetailAPI可用时）
			if bookDetailApiHandler != nil {
				public.GET("/books/:id/detail", bookDetailApiHandler.GetBookDetail)
				public.GET("/books/:id/similar", bookDetailApiHandler.GetSimilarBooks)
				public.GET("/books/:id/statistics", bookDetailApiHandler.GetBookStatistics)
			}

			// ✅ 统计API（当StatisticsService可用时）
			// 注意：BookDetailAPI中已包含GetBookStatistics
			// Statistics API在BookDetail中已实现，这里备注即可

			// ℹ️ Chapter API路由需要ChapterService支持
			// 当ChapterService实现后，可以启用以下路由:
			// if chapterApiHandler != nil {
			// 	public.GET("/chapters/:id", chapterApiHandler.GetChapter)
			// 	public.GET("/chapters/book/:id", chapterApiHandler.GetChaptersByBookID)
			// }
		}

		// 需要认证的接口
		authenticated := bookstoreGroup.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// ✅ 书籍点击记录（认证接口 - 关联到用户）
			authenticated.POST("/books/:id/view", bookstoreApiHandler.IncrementBookView)

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
}
