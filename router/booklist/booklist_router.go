package booklist

import (
	"github.com/gin-gonic/gin"

	booklistAPI "Qingyu_backend/api/v1/booklist"
	"Qingyu_backend/middleware"
)

// RegisterBooklistRoutes 注册书单路由
func RegisterBooklistRoutes(
	r *gin.RouterGroup,
	booklistAPI *booklistAPI.BookListAPI,
) {
	// 应用中间件
	booklistGroup := r.Group("/booklists")
	booklistGroup.Use(middleware.ResponseFormatterMiddleware())
	booklistGroup.Use(middleware.ResponseTimingMiddleware())
	booklistGroup.Use(middleware.CORSMiddleware())
	booklistGroup.Use(middleware.Recovery())

	// 公开路由（无需认证）
	public := booklistGroup.Group("")
	{
		public.GET("", booklistAPI.GetBookLists)           // 获取书单列表
		public.GET("/:id", booklistAPI.GetBookListDetail)  // 获取书单详情
		public.GET("/:id/books", booklistAPI.GetBooksInList) // 获取书单中的书籍
	}

	// 需要认证的路由
	authenticated := booklistGroup.Group("")
	authenticated.Use(middleware.JWTAuth())
	{
		authenticated.POST("", booklistAPI.CreateBookList)        // 创建书单
		authenticated.PUT("/:id", booklistAPI.UpdateBookList)     // 更新书单
		authenticated.DELETE("/:id", booklistAPI.DeleteBookList)  // 删除书单
		authenticated.POST("/:id/like", booklistAPI.LikeBookList) // 点赞书单
		authenticated.POST("/:id/fork", booklistAPI.ForkBookList) // 复制书单
	}
}
