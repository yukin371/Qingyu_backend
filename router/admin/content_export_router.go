package admin

import (
	"github.com/gin-gonic/gin"

	adminAPI "Qingyu_backend/api/v1/admin"
	"Qingyu_backend/internal/middleware/auth"
	serviceinterfaces "Qingyu_backend/service/interfaces"
)

// RegisterContentExportRoutes 注册管理员内容导出路由。
//
// 内容导出依赖组合根提供的最小读取端口。端口缺失时不注册路由，
// 保持无 MongoDB 环境下其余管理员路由的渐进式启动行为。
func RegisterContentExportRoutes(
	r *gin.RouterGroup,
	bookPort serviceinterfaces.AdminContentBookPort,
	chapterPort serviceinterfaces.AdminContentChapterPort,
) {
	if r == nil || bookPort == nil || chapterPort == nil {
		return
	}

	api := adminAPI.NewContentExportAPI(bookPort, chapterPort)
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.JWTAuth())
	adminGroup.Use(auth.RequireRole("admin"))

	contentGroup := adminGroup.Group("/content")
	contentGroup.GET("/books/export", api.ExportBooks)
	contentGroup.GET("/books/export/template", api.GetBookExportTemplate)
	contentGroup.GET("/chapters/export", api.ExportChapters)
	contentGroup.GET("/chapters/export/template", api.GetChapterExportTemplate)
}
