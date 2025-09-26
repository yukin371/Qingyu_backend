package reading

import (
	"github.com/gin-gonic/gin"
	
	readingAPI "Qingyu_backend/api/v1/reading"
)

// BookstoreRouter 书城路由组
type BookstoreRouter struct {
	api *readingAPI.BookstoreAPI
}

// NewBookstoreRouter 创建书城路由实例
func NewBookstoreRouter(api *readingAPI.BookstoreAPI) *BookstoreRouter {
	return &BookstoreRouter{
		api: api,
	}
}

// RegisterRoutes 注册书城相关路由
func (r *BookstoreRouter) RegisterRoutes(rg *gin.RouterGroup) {
	// 书城路由组
	bookstore := rg.Group("/bookstore")
	{
		// 首页数据
		bookstore.GET("/homepage", r.api.GetHomepage)
		
		// 书籍相关路由
		books := bookstore.Group("/books")
		{
			// 获取书籍详情
			books.GET("/:id", r.api.GetBookByID)
			
			// 增加书籍浏览量
			books.POST("/:id/view", r.api.IncrementBookView)
			
			// 获取推荐书籍
			books.GET("/recommended", r.api.GetRecommendedBooks)
			
			// 获取精选书籍
			books.GET("/featured", r.api.GetFeaturedBooks)
			
			// 搜索书籍
			books.GET("/search", r.api.SearchBooks)
		}
		
		// 分类相关路由
		categories := bookstore.Group("/categories")
		{
			// 获取分类树
			categories.GET("/tree", r.api.GetCategoryTree)
			
			// 获取分类详情
			categories.GET("/:id", r.api.GetCategoryByID)
			
			// 根据分类获取书籍列表
			categories.GET("/:categoryId/books", r.api.GetBooksByCategory)
		}
		
		// 榜单相关路由
		rankings := bookstore.Group("/rankings")
		{
			// 获取实时榜
			rankings.GET("/realtime", r.api.GetRealtimeRanking)
			
			// 获取周榜
			rankings.GET("/weekly", r.api.GetWeeklyRanking)
			
			// 获取月榜
			rankings.GET("/monthly", r.api.GetMonthlyRanking)
			
			// 获取新人榜
			rankings.GET("/newbie", r.api.GetNewbieRanking)
			
			// 根据类型获取榜单
			rankings.GET("/:type", r.api.GetRankingByType)
		}
		
		// Banner相关路由
		banners := bookstore.Group("/banners")
		{
			// 获取激活的Banner列表
			banners.GET("", r.api.GetActiveBanners)
			
			// 增加Banner点击次数
			banners.POST("/:id/click", r.api.IncrementBannerClick)
		}
	}
}

// RegisterPublicRoutes 注册公开路由（不需要认证）
func (r *BookstoreRouter) RegisterPublicRoutes(rg *gin.RouterGroup) {
	// 公开的书城路由
	public := rg.Group("/public/bookstore")
	{
		// 首页数据（公开访问）
		public.GET("/homepage", r.api.GetHomepage)
		
		// 书籍相关公开路由
		books := public.Group("/books")
		{
			// 获取书籍详情（公开访问）
			books.GET("/:id", r.api.GetBookByID)
			
			// 获取推荐书籍（公开访问）
			books.GET("/recommended", r.api.GetRecommendedBooks)
			
			// 获取精选书籍（公开访问）
			books.GET("/featured", r.api.GetFeaturedBooks)
			
			// 搜索书籍（公开访问）
			books.GET("/search", r.api.SearchBooks)
		}
		
		// 分类相关公开路由
		categories := public.Group("/categories")
		{
			// 获取分类树（公开访问）
			categories.GET("/tree", r.api.GetCategoryTree)
			
			// 获取分类详情（公开访问）
			categories.GET("/:id", r.api.GetCategoryByID)
			
			// 根据分类获取书籍列表（公开访问）
			categories.GET("/:categoryId/books", r.api.GetBooksByCategory)
		}
		
		// Banner相关公开路由
		banners := public.Group("/banners")
		{
			// 获取激活的Banner列表（公开访问）
			banners.GET("", r.api.GetActiveBanners)
		}
	}
}

// RegisterPrivateRoutes 注册需要认证的路由
func (r *BookstoreRouter) RegisterPrivateRoutes(rg *gin.RouterGroup) {
	// 需要认证的书城路由
	private := rg.Group("/bookstore")
	// 这里可以添加认证中间件
	// private.Use(middleware.AuthRequired())
	{
		// 书籍相关私有路由
		books := private.Group("/books")
		{
			// 增加书籍浏览量（需要认证以防刷量）
			books.POST("/:id/view", r.api.IncrementBookView)
		}
		
		// Banner相关私有路由
		banners := private.Group("/banners")
		{
			// 增加Banner点击次数（需要认证以防刷量）
			banners.POST("/:id/click", r.api.IncrementBannerClick)
		}
	}
}