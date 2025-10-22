package reader

import (
	readerApi "Qingyu_backend/api/v1/reader"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/reading"

	"github.com/gin-gonic/gin"
)

// InitReaderRouter 初始化阅读器路由
func InitReaderRouter(
	r *gin.RouterGroup,
	readerService *reading.ReaderService,
) {
	// 创建API实例
	progressApiHandler := readerApi.NewProgressAPI(readerService)
	chaptersApiHandler := readerApi.NewChaptersAPI(readerService)
	annotationsApiHandler := readerApi.NewAnnotationsAPI(readerService)
	settingApiHandler := readerApi.NewSettingAPI(readerService)
	booksApiHandler := readerApi.NewBooksAPI(readerService)

	// 阅读器主路由组（需要认证）
	readerGroup := r.Group("/reader")
	readerGroup.Use(middleware.JWTAuth())
	{
		// 书架管理
		books := readerGroup.Group("/books")
		{
			books.GET("", booksApiHandler.GetBookshelf)                   // 获取书架
			books.GET("/recent", booksApiHandler.GetRecentReading)        // 获取最近阅读
			books.GET("/unfinished", booksApiHandler.GetUnfinishedBooks)  // 获取未读完
			books.GET("/finished", booksApiHandler.GetFinishedBooks)      // 获取已读完
			books.POST("/:bookId", booksApiHandler.AddToBookshelf)        // 添加到书架
			books.DELETE("/:bookId", booksApiHandler.RemoveFromBookshelf) // 从书架移除
		}

		// 章节内容（阅读）
		chapters := readerGroup.Group("/chapters")
		{
			chapters.GET("/:id", chaptersApiHandler.GetChapterByID)                   // 获取章节信息
			chapters.GET("/:id/content", chaptersApiHandler.GetChapterContent)        // 获取章节内容
			chapters.GET("/book/:bookId", chaptersApiHandler.GetBookChapters)         // 获取书籍章节列表
			chapters.GET("/:id/navigation", chaptersApiHandler.GetNavigationChapters) // 获取导航章节
			chapters.GET("/book/:bookId/first", chaptersApiHandler.GetFirstChapter)   // 获取第一章
			chapters.GET("/book/:bookId/last", chaptersApiHandler.GetLastChapter)     // 获取最后一章
		}

		// 阅读进度
		progress := readerGroup.Group("/progress")
		{
			progress.GET("/:bookId", progressApiHandler.GetReadingProgress)    // 获取阅读进度
			progress.POST("", progressApiHandler.SaveReadingProgress)          // 保存阅读进度
			progress.POST("/time", progressApiHandler.UpdateReadingTime)       // 更新阅读时长
			progress.GET("/recent", progressApiHandler.GetRecentReading)       // 获取最近阅读
			progress.GET("/history", progressApiHandler.GetReadingHistory)     // 获取阅读历史
			progress.GET("/stats", progressApiHandler.GetReadingStats)         // 获取阅读统计
			progress.GET("/unfinished", progressApiHandler.GetUnfinishedBooks) // 获取未读完的书
			progress.GET("/finished", progressApiHandler.GetFinishedBooks)     // 获取已读完的书
		}

		// 标注管理
		annotations := readerGroup.Group("/annotations")
		{
			// 基础CRUD
			annotations.POST("", annotationsApiHandler.CreateAnnotation)       // 创建标注
			annotations.PUT("/:id", annotationsApiHandler.UpdateAnnotation)    // 更新标注
			annotations.DELETE("/:id", annotationsApiHandler.DeleteAnnotation) // 删除标注

			// 批量操作
			annotations.POST("/batch", annotationsApiHandler.BatchCreateAnnotations)   // 批量创建
			annotations.PUT("/batch", annotationsApiHandler.BatchUpdateAnnotations)    // 批量更新
			annotations.DELETE("/batch", annotationsApiHandler.BatchDeleteAnnotations) // 批量删除

			// 按类型查询
			annotations.GET("/notes", annotationsApiHandler.GetNotes)           // 获取笔记
			annotations.GET("/bookmarks", annotationsApiHandler.GetBookmarks)   // 获取书签
			annotations.GET("/highlights", annotationsApiHandler.GetHighlights) // 获取高亮

			// 按范围查询
			annotations.GET("/book/:bookId", annotationsApiHandler.GetAnnotationsByBook)          // 获取书籍标注
			annotations.GET("/chapter/:chapterId", annotationsApiHandler.GetAnnotationsByChapter) // 获取章节标注
			annotations.GET("/recent", annotationsApiHandler.GetRecentAnnotations)                // 获取最近标注
			annotations.GET("/public", annotationsApiHandler.GetPublicAnnotations)                // 获取公开标注

			// 搜索和统计
			annotations.GET("/search", annotationsApiHandler.SearchNotes)       // 搜索笔记
			annotations.GET("/stats", annotationsApiHandler.GetAnnotationStats) // 标注统计

			// 同步和导出
			annotations.POST("/sync", annotationsApiHandler.SyncAnnotations)             // 同步标注
			annotations.GET("/export", annotationsApiHandler.ExportAnnotations)          // 导出标注
			annotations.GET("/bookmark/latest", annotationsApiHandler.GetLatestBookmark) // 获取最新书签
		}

		// 阅读设置
		settings := readerGroup.Group("/settings")
		{
			settings.GET("", settingApiHandler.GetReadingSettings)    // 获取阅读设置
			settings.POST("", settingApiHandler.SaveReadingSettings)  // 保存阅读设置
			settings.PUT("", settingApiHandler.UpdateReadingSettings) // 更新阅读设置
		}
	}
}
