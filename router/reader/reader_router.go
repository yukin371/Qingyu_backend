package reader

import (
	readerApi "Qingyu_backend/api/v1/reader"
	socialApi "Qingyu_backend/api/v1/social"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service/bookstore"
	readerservice "Qingyu_backend/service/reader"
	socialService "Qingyu_backend/service/social"
	syncService "Qingyu_backend/pkg/sync"

	"github.com/gin-gonic/gin"
)

// InitReaderRouter 初始化阅读器路由
func InitReaderRouter(
	r *gin.RouterGroup,
	readerService *readerservice.ReaderService,
	chapterService bookstore.ChapterService,
	commentService *socialService.CommentService,
	likeService *socialService.LikeService,
	collectionService *socialService.CollectionService,
	readingHistoryService *readerservice.ReadingHistoryService,
	progressSyncService *syncService.ProgressSyncService,
	bookmarkService readerservice.BookmarkService,
) {
	// 创建API实例
	progressApiHandler := readerApi.NewProgressAPI(readerService)
	annotationsApiHandler := readerApi.NewAnnotationsAPI(readerService)
	settingApiHandler := readerApi.NewSettingAPI(readerService)
	booksApiHandler := readerApi.NewBooksAPI(readerService)
	themeApiHandler := readerApi.NewThemeAPI()
	fontApiHandler := readerApi.NewFontAPI()
	chapterCommentApiHandler := readerApi.NewChapterCommentAPI()

	// 章节API（使用阅读器专属服务）
	chapterServiceForReader := readerservice.NewChapterService(chapterService, readerService, nil)
	chapterApiHandler := readerApi.NewChapterAPI(chapterServiceForReader)

	// 书签API（如果可用）
	var bookmarkApiHandler *readerApi.BookmarkAPI
	if bookmarkService != nil {
		bookmarkApiHandler = readerApi.NewBookmarkAPI(bookmarkService)
	}

	// 进度同步API（如果syncService可用）
	var syncApiHandler *readerApi.SyncAPI
	if progressSyncService != nil {
		syncApiHandler = readerApi.NewSyncAPI(progressSyncService)
		// 启动同步服务
		progressSyncService.Start()
	}

	// 评论API（如果commentService可用，使用 social 包的统一实现）
	var commentApiHandler *socialApi.CommentAPI
	if commentService != nil {
		commentApiHandler = socialApi.NewCommentAPI(commentService)
	}

	// 点赞API（如果likeService可用，使用 social 包的统一实现）
	var likeApiHandler *socialApi.LikeAPI
	if likeService != nil {
		likeApiHandler = socialApi.NewLikeAPI(likeService)
	}

	// 收藏API（如果collectionService可用，使用 social 包的统一实现）
	var collectionApiHandler *socialApi.CollectionAPI
	if collectionService != nil {
		collectionApiHandler = socialApi.NewCollectionAPI(collectionService)
	}

	// 阅读历史API（如果readingHistoryService可用）
	var historyApiHandler *readerApi.ReadingHistoryAPI
	if readingHistoryService != nil {
		historyApiHandler = readerApi.NewReadingHistoryAPI(readingHistoryService)
	}

	// ========================================
	// 公开路由（不需要认证）
	// ========================================
	readerPublic := r.Group("/reader")
	{
		// 公开的收藏列表
		if collectionApiHandler != nil {
			readerPublic.GET("/collections/public", collectionApiHandler.GetPublicCollections)
		}
	}

	// ========================================
	// 需要认证的路由
	// ========================================
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
				books.POST("/:bookId/like", likeApiHandler.LikeBook)            // 点赞书籍
				books.DELETE("/:bookId/like", likeApiHandler.UnlikeBook)        // 取消点赞书籍
				books.GET("/:bookId/like/info", likeApiHandler.GetBookLikeInfo) // 获取点赞信息
			}
		}

		// 章节阅读
		{
			// 章节内容获取（支持按ID和按章节号）
			readerGroup.GET("/books/:bookId/chapters/:chapterId", chapterApiHandler.GetChapterContent)
			readerGroup.GET("/books/:bookId/chapters/by-number/:chapterNum", chapterApiHandler.GetChapterByNumber)

			// 章节导航
			readerGroup.GET("/books/:bookId/chapters/:chapterId/next", chapterApiHandler.GetNextChapter)
			readerGroup.GET("/books/:bookId/chapters/:chapterId/previous", chapterApiHandler.GetPreviousChapter)

			// 章节目录
			readerGroup.GET("/books/:bookId/chapters", chapterApiHandler.GetChapterList)

			// 章节信息（不含内容）
			readerGroup.GET("/chapters/:chapterId/info", chapterApiHandler.GetChapterInfo)
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

			// 进度同步（如果syncApiHandler可用）
			if syncApiHandler != nil {
				progress.GET("/ws", syncApiHandler.SyncWebSocket)              // WebSocket同步
				progress.POST("/sync", syncApiHandler.SyncProgress)            // HTTP同步
				progress.POST("/merge", syncApiHandler.MergeOfflineProgresses) // 合并离线进度
				progress.GET("/sync-status", syncApiHandler.GetSyncStatus)     // 同步状态
			}
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

		// 书签管理（如果bookmarkApiHandler可用）
		if bookmarkApiHandler != nil {
			bookmarks := readerGroup.Group("/bookmarks")
			{
				// 基础CRUD
				bookmarks.GET("", bookmarkApiHandler.GetBookmarks)                // 获取书签列表
				bookmarks.GET("/:id", bookmarkApiHandler.GetBookmark)            // 获取书签详情
				bookmarks.PUT("/:id", bookmarkApiHandler.UpdateBookmark)         // 更新书签
				bookmarks.DELETE("/:id", bookmarkApiHandler.DeleteBookmark)      // 删除书签

				// 按书籍获取
				readerGroup.GET("/books/:bookId/bookmarks", bookmarkApiHandler.GetBookmarks) // 获取某本书的书签
				readerGroup.POST("/books/:bookId/bookmarks", bookmarkApiHandler.CreateBookmark) // 创建书签

				// 搜索和统计
				bookmarks.GET("/search", bookmarkApiHandler.SearchBookmarks) // 搜索书签
				bookmarks.GET("/stats", bookmarkApiHandler.GetBookmarkStats) // 书签统计

				// 导出
				bookmarks.GET("/export", bookmarkApiHandler.ExportBookmarks) // 导出书签
			}
		}

		// 阅读设置
		settings := readerGroup.Group("/settings")
		{
			settings.GET("", settingApiHandler.GetReadingSettings)    // 获取阅读设置
			settings.POST("", settingApiHandler.SaveReadingSettings)  // 保存阅读设置
			settings.PUT("", settingApiHandler.UpdateReadingSettings) // 更新阅读设置
		}

		// 主题管理
		themes := readerGroup.Group("/themes")
		{
			themes.GET("", themeApiHandler.GetThemes)                    // 获取主题列表
			themes.GET("/:name", themeApiHandler.GetThemeByName)         // 获取单个主题
			themes.POST("", themeApiHandler.CreateCustomTheme)           // 创建自定义主题
			themes.PUT("/:id", themeApiHandler.UpdateTheme)              // 更新主题
			themes.DELETE("/:id", themeApiHandler.DeleteTheme)           // 删除主题
			themes.POST("/:name/activate", themeApiHandler.ActivateTheme) // 激活主题
		}

		// 字体管理
		fonts := readerGroup.Group("/fonts")
		{
			fonts.GET("", fontApiHandler.GetFonts)              // 获取字体列表
			fonts.GET("/:name", fontApiHandler.GetFontByName)   // 获取单个字体
			fonts.POST("", fontApiHandler.CreateCustomFont)     // 创建自定义字体
			fonts.PUT("/:id", fontApiHandler.UpdateFont)        // 更新字体
			fonts.DELETE("/:id", fontApiHandler.DeleteFont)     // 删除字体
		}

		// 字体偏好设置
		readerGroup.POST("/settings/font", fontApiHandler.SetFontPreference) // 设置字体偏好

		// 章节评论
		chapters := readerGroup.Group("/chapters")
		{
			// 章节级评论
			chapters.GET("/:chapterId/comments", chapterCommentApiHandler.GetChapterComments)           // 获取章节评论列表
			chapters.POST("/:chapterId/comments", chapterCommentApiHandler.CreateChapterComment)       // 发表章节评论

			// 段落级评论
			chapters.GET("/:chapterId/paragraph-comments", chapterCommentApiHandler.GetChapterParagraphComments) // 获取章节段落评论概览
			chapters.POST("/:chapterId/paragraph-comments", chapterCommentApiHandler.CreateParagraphComment)     // 发表段落评论
			chapters.GET("/:chapterId/paragraphs/:paragraphIndex/comments", chapterCommentApiHandler.GetParagraphComments) // 获取特定段落评论
		}

		// 章节评论管理（单条评论操作）
		chapterComments := readerGroup.Group("/chapter-comments")
		{
			chapterComments.GET("/:commentId", chapterCommentApiHandler.GetChapterComment)        // 获取评论详情
			chapterComments.PUT("/:commentId", chapterCommentApiHandler.UpdateChapterComment)     // 更新评论
			chapterComments.DELETE("/:commentId", chapterCommentApiHandler.DeleteChapterComment)  // 删除评论
			chapterComments.POST("/:commentId/like", chapterCommentApiHandler.LikeChapterComment)     // 点赞评论
			chapterComments.DELETE("/:commentId/like", chapterCommentApiHandler.UnlikeChapterComment) // 取消点赞
		}

		// 通用评论模块（如果commentService可用）
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
				// 注意：GET("/public") 已移至公开路由组，不需要认证

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
