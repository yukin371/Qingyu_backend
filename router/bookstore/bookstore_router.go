package bookstore

import (
	bookstoreApi "Qingyu_backend/api/v1/bookstore"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/bookstore"

	"github.com/gin-gonic/gin"
)

// InitBookstoreRouter 初始化书店路由
func InitBookstoreRouter(
	r *gin.RouterGroup,
	bookstoreService bookstore.BookstoreService,
	bookDetailService interface{}, // TODO: 改为具体类型
	ratingService interface{}, // TODO: 改为具体类型
	statisticsService interface{}, // TODO: 改为具体类型
) {
	// 创建API实例
	bookstoreApiHandler := bookstoreApi.NewBookstoreAPI(bookstoreService)

	// TODO: 当其他服务实现后，取消注释
	// if bookDetailService != nil {
	// 	bookDetailApiHandler := bookstoreApi.NewBookDetailAPI(bookDetailService.(bookstore.BookDetailService))
	// }
	// if ratingService != nil {
	// 	ratingApiHandler := bookstoreApi.NewBookRatingAPI(ratingService.(bookstore.RatingService))
	// }
	// if statisticsService != nil {
	// 	statisticsApiHandler := bookstoreApi.NewBookStatisticsAPI(statisticsService.(bookstore.StatisticsService))
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

			// 书籍列表和搜索
			public.GET("/books/:id", bookstoreApiHandler.GetBookByID)
			public.GET("/books/search", bookstoreApiHandler.SearchBooks)

			// 推荐和精选
			public.GET("/recommended", bookstoreApiHandler.GetRecommendedBooks)
			public.GET("/featured", bookstoreApiHandler.GetFeaturedBooks)

			// 分类
			public.GET("/categories", bookstoreApiHandler.GetCategoryTree)
			public.GET("/categories/:id", bookstoreApiHandler.GetCategoryByID)
			public.GET("/categories/:id/books", bookstoreApiHandler.GetBooksByCategory)

			// Banner
			public.GET("/banners", bookstoreApiHandler.GetActiveBanners)

			// 排行榜
			public.GET("/rankings/realtime", bookstoreApiHandler.GetRealtimeRanking)
			public.GET("/rankings/weekly", bookstoreApiHandler.GetWeeklyRanking)
			public.GET("/rankings/monthly", bookstoreApiHandler.GetMonthlyRanking)
			public.GET("/rankings/newbie", bookstoreApiHandler.GetNewbieRanking)
			public.GET("/rankings/:type", bookstoreApiHandler.GetRankingByType)

			// TODO: 当BookDetailAPI实现后添加
			// public.GET("/books/:id/detail", bookDetailApiHandler.GetBookDetail)
			// public.GET("/books/:id/similar", bookDetailApiHandler.GetSimilarBooks)
			// public.GET("/books/:id/statistics", bookDetailApiHandler.GetBookStatistics)

			// TODO: 当ChapterAPI实现后添加
			// public.GET("/chapters/:id", chapterApiHandler.GetChapter)
			// public.GET("/chapters/book/:id", chapterApiHandler.GetChaptersByBookID)
		}

		// 需要认证的接口
		authenticated := bookstoreGroup.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// 统计点击
			authenticated.POST("/books/:id/view", bookstoreApiHandler.IncrementBookView)
			authenticated.POST("/banners/:id/click", bookstoreApiHandler.IncrementBannerClick)

			// TODO: 当RatingAPI实现后添加
			// authenticated.GET("/books/:id/rating", ratingApiHandler.GetBookRating)
			// authenticated.POST("/books/:id/rating", ratingApiHandler.CreateRating)
			// authenticated.PUT("/books/:id/rating", ratingApiHandler.UpdateRating)
			// authenticated.DELETE("/books/:id/rating", ratingApiHandler.DeleteRating)
			// authenticated.GET("/ratings/user/:id", ratingApiHandler.GetRatingsByUserID)
		}
	}
}
