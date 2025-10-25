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
	commentService *reading.CommentService,
	likeService *reading.LikeService,
	collectionService *reading.CollectionService,
	readingHistoryService *reading.ReadingHistoryService,
) {
	// 创建API实例
	progressApiHandler := readerApi.NewProgressAPI(readerService)
	chaptersApiHandler := readerApi.NewChaptersAPI(readerService)
	annotationsApiHandler := readerApi.NewAnnotationsAPI(readerService)
	settingApiHandler := readerApi.NewSettingAPI(readerService)
	booksApiHandler := readerApi.NewBooksAPI(readerService)

	// 评论API（如果commentService可用）
	var commentApiHandler *readerApi.CommentAPI
	if commentService != nil {
		commentApiHandler = readerApi.NewCommentAPI(commentService)
	}

	// 点赞API（如果likeService可用）
	var likeApiHandler *readerApi.LikeAPI
	if likeService != nil {
		likeApiHandler = readerApi.NewLikeAPI(likeService)
	}

	// 收藏API（如果collectionService可用）
	var collectionApiHandler *readerApi.CollectionAPI
	if collectionService != nil {
		collectionApiHandler = readerApi.NewCollectionAPI(collectionService)
	}

	// 阅读历史API（如果readingHistoryService可用）
	var historyApiHandler *readerApi.ReadingHistoryAPI
	if readingHistoryService != nil {
		historyApiHandler = readerApi.NewReadingHistoryAPI(readingHistoryService)
	}

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

			// 书籍点赞（如果likeApiHandler可用）
			if likeApiHandler != nil {
				books.POST("/:id/like", likeApiHandler.LikeBook)            // 点赞书籍
				books.DELETE("/:id/like", likeApiHandler.UnlikeBook)        // 取消点赞书籍
				books.GET("/:id/like/info", likeApiHandler.GetBookLikeInfo) // 获取点赞信息
			}
		}

		// 章节内容（阅读）
		chapters := readerGroup.Group("/chapters")
		{
			// 注意：查询参数形式的路由必须在参数化路由之前
			chapters.GET("", chaptersApiHandler.GetBookChapters)                      // 获取书籍章节列表（查询参数：bookId）
			chapters.GET("/:id", chaptersApiHandler.GetChapterByID)                   // 获取章节信息
			chapters.GET("/:id/content", chaptersApiHandler.GetChapterContent)        // 获取章节内容
			chapters.GET("/:id/navigation", chaptersApiHandler.GetNavigationChapters) // 获取导航章节
			chapters.GET("/book/:bookId", chaptersApiHandler.GetBookChapters)         // 获取书籍章节列表（路径参数）
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

		// 评论模块（如果commentService可用）
		if commentApiHandler != nil {
			comments := readerGroup.Group("/comments")
			{
				comments.POST("", commentApiHandler.CreateComment)            // 发表评论
				comments.GET("", commentApiHandler.GetCommentList)            // 获取评论列表
				comments.GET("/:id", commentApiHandler.GetCommentDetail)      // 获取评论详情
				comments.PUT("/:id", commentApiHandler.UpdateComment)         // 更新评论
				comments.DELETE("/:id", commentApiHandler.DeleteComment)      // 删除评论
				comments.POST("/:id/reply", commentApiHandler.ReplyComment)   // 回复评论
				comments.POST("/:id/like", commentApiHandler.LikeComment)     // 点赞评论
				comments.DELETE("/:id/like", commentApiHandler.UnlikeComment) // 取消点赞
			}
		}

		// 点赞模块（如果likeApiHandler可用）
		if likeApiHandler != nil {
			likes := readerGroup.Group("/likes")
			{
				likes.GET("/books", likeApiHandler.GetUserLikedBooks) // 获取用户点赞的书籍列表
				likes.GET("/stats", likeApiHandler.GetUserLikeStats)  // 获取用户点赞统计
			}
		}

		// 收藏模块（如果collectionApiHandler可用）
		if collectionApiHandler != nil {
			collections := readerGroup.Group("/collections")
			{
				// 收藏管理
				collections.POST("", collectionApiHandler.AddCollection)                 // 添加收藏
				collections.GET("", collectionApiHandler.GetCollections)                 // 获取收藏列表
				collections.PUT("/:id", collectionApiHandler.UpdateCollection)           // 更新收藏
				collections.DELETE("/:id", collectionApiHandler.DeleteCollection)        // 删除收藏
				collections.GET("/check/:book_id", collectionApiHandler.CheckCollected)  // 检查是否已收藏
				collections.GET("/tags/:tag", collectionApiHandler.GetCollectionsByTag)  // 根据标签获取收藏
				collections.GET("/stats", collectionApiHandler.GetCollectionStats)       // 获取收藏统计
				collections.POST("/:id/share", collectionApiHandler.ShareCollection)     // 分享收藏
				collections.DELETE("/:id/share", collectionApiHandler.UnshareCollection) // 取消分享收藏
				collections.GET("/public", collectionApiHandler.GetPublicCollections)    // 获取公开收藏列表

				// 收藏夹管理
				collections.POST("/folders", collectionApiHandler.CreateFolder)       // 创建收藏夹
				collections.GET("/folders", collectionApiHandler.GetFolders)          // 获取收藏夹列表
				collections.PUT("/folders/:id", collectionApiHandler.UpdateFolder)    // 更新收藏夹
				collections.DELETE("/folders/:id", collectionApiHandler.DeleteFolder) // 删除收藏夹
			}
		}

		// 阅读历史模块（如果historyApiHandler可用）
		if historyApiHandler != nil {
			history := readerGroup.Group("/reading-history")
			{
				history.POST("", historyApiHandler.RecordReading)        // 记录阅读历史
				history.GET("", historyApiHandler.GetReadingHistories)   // 获取阅读历史列表
				history.GET("/stats", historyApiHandler.GetReadingStats) // 获取阅读统计
				history.DELETE("/:id", historyApiHandler.DeleteHistory)  // 删除单条历史记录
				history.DELETE("", historyApiHandler.ClearHistories)     // 清空历史记录
			}
		}
	}
}
