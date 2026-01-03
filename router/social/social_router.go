package social

import (
	"github.com/gin-gonic/gin"

	socialApi "Qingyu_backend/api/v1/social"
	"Qingyu_backend/middleware"
)

// RegisterSocialRoutes 注册所有社交相关路由
func RegisterSocialRoutes(r *gin.RouterGroup,
	relationAPI *socialApi.UserRelationAPI,
	commentAPI *socialApi.CommentAPI,
	likeAPI *socialApi.LikeAPI,
	collectionAPI *socialApi.CollectionAPI,
	followAPI *socialApi.FollowAPI,
	messageAPI *socialApi.MessageAPI,
	reviewAPI *socialApi.ReviewAPI,
	bookListAPI *socialApi.BookListAPI) {

	// 社交路由需要认证
	socialGroup := r.Group("/social")
	socialGroup.Use(middleware.JWTAuth())
	{
		// ========== 用户关注相关（原有） ==========
		if relationAPI != nil {
			socialGroup.POST("/follow/:userId", relationAPI.FollowUser)
			socialGroup.DELETE("/follow/:userId", relationAPI.UnfollowUser)
			socialGroup.GET("/follow/:userId/status", relationAPI.CheckIsFollowing)

			// ========== 关注列表相关 ==========
			socialGroup.GET("/users/:userId/followers", relationAPI.GetFollowers)
			socialGroup.GET("/users/:userId/following", relationAPI.GetFollowing)
			socialGroup.GET("/users/:userId/follow-stats", relationAPI.GetFollowStats)
		}

		// ========== 评论相关 ==========
		if commentAPI != nil {
			socialGroup.POST("/comments", commentAPI.CreateComment)
			socialGroup.GET("/comments", commentAPI.GetCommentList)
			socialGroup.GET("/comments/:id", commentAPI.GetCommentDetail)
			socialGroup.PUT("/comments/:id", commentAPI.UpdateComment)
			socialGroup.DELETE("/comments/:id", commentAPI.DeleteComment)
			socialGroup.POST("/comments/:id/reply", commentAPI.ReplyComment)
			socialGroup.GET("/comments/:id/thread", commentAPI.GetCommentThread)
			socialGroup.POST("/comments/:id/like", commentAPI.LikeComment)
			socialGroup.DELETE("/comments/:id/like", commentAPI.UnlikeComment)
		}

		// ========== 点赞相关 ==========
		if likeAPI != nil {
			socialGroup.POST("/books/:bookId/like", likeAPI.LikeBook)
			socialGroup.DELETE("/books/:bookId/like", likeAPI.UnlikeBook)
			socialGroup.GET("/books/:bookId/like/info", likeAPI.GetBookLikeInfo)
		}

		// ========== 收藏相关 ==========
		if collectionAPI != nil {
			socialGroup.POST("/collections", collectionAPI.AddCollection)
			socialGroup.GET("/collections", collectionAPI.GetCollections)
			socialGroup.PUT("/collections/:id", collectionAPI.UpdateCollection)
			socialGroup.DELETE("/collections/:id", collectionAPI.DeleteCollection)
			socialGroup.GET("/collections/check", collectionAPI.CheckCollected)
			socialGroup.GET("/collections/by-tag", collectionAPI.GetCollectionsByTag)
			socialGroup.POST("/collections/folders", collectionAPI.CreateFolder)
			socialGroup.GET("/collections/folders", collectionAPI.GetFolders)
			socialGroup.PUT("/collections/folders/:id", collectionAPI.UpdateFolder)
			socialGroup.DELETE("/collections/folders/:id", collectionAPI.DeleteFolder)
			socialGroup.POST("/collections/:id/share", collectionAPI.ShareCollection)
			socialGroup.DELETE("/collections/:id/share", collectionAPI.UnshareCollection)
			socialGroup.GET("/collections/public", collectionAPI.GetPublicCollections)
			socialGroup.GET("/collections/stats", collectionAPI.GetCollectionStats)
		}

		// ========== 新增：关注系统 ==========
		if followAPI != nil {
			// 用户关注
			socialGroup.POST("/users/:userId/follow", followAPI.FollowUser)
			socialGroup.DELETE("/users/:userId/unfollow", followAPI.UnfollowUser)
			socialGroup.GET("/users/:userId/followers", followAPI.GetFollowers)
			socialGroup.GET("/users/:userId/following", followAPI.GetFollowing)
			socialGroup.GET("/users/:userId/follow-status", followAPI.CheckFollowStatus)

			// 作者关注
			socialGroup.POST("/authors/:authorId/follow", followAPI.FollowAuthor)
			socialGroup.DELETE("/authors/:authorId/unfollow", followAPI.UnfollowAuthor)
			socialGroup.GET("/following/authors", followAPI.GetFollowingAuthors)
		}

		// ========== 新增：私信系统 ==========
		if messageAPI != nil {
			// 会话管理
			socialGroup.GET("/messages/conversations", messageAPI.GetConversations)
			socialGroup.GET("/messages/:conversationId", messageAPI.GetConversationMessages)

			// 消息管理
			socialGroup.POST("/messages", messageAPI.SendMessage)
			socialGroup.PUT("/messages/:id/read", messageAPI.MarkMessageAsRead)
			socialGroup.DELETE("/messages/:id", messageAPI.DeleteMessage)

			// @提醒
			socialGroup.POST("/mentions", messageAPI.CreateMention)
			socialGroup.GET("/mentions", messageAPI.GetMentions)
			socialGroup.PUT("/mentions/:id/read", messageAPI.MarkMentionAsRead)
		}

		// ========== 新增：书评系统 ==========
		if reviewAPI != nil {
			socialGroup.GET("/reviews", reviewAPI.GetReviews)
			socialGroup.POST("/reviews", reviewAPI.CreateReview)
			socialGroup.GET("/reviews/:id", reviewAPI.GetReviewDetail)
			socialGroup.PUT("/reviews/:id", reviewAPI.UpdateReview)
			socialGroup.DELETE("/reviews/:id", reviewAPI.DeleteReview)
			socialGroup.POST("/reviews/:id/like", reviewAPI.LikeReview)
		}

		// ========== 新增：书单系统 ==========
		if bookListAPI != nil {
			socialGroup.GET("/booklists", bookListAPI.GetBookLists)
			socialGroup.POST("/booklists", bookListAPI.CreateBookList)
			socialGroup.GET("/booklists/:id", bookListAPI.GetBookListDetail)
			socialGroup.PUT("/booklists/:id", bookListAPI.UpdateBookList)
			socialGroup.DELETE("/booklists/:id", bookListAPI.DeleteBookList)
			socialGroup.POST("/booklists/:id/like", bookListAPI.LikeBookList)
			socialGroup.POST("/booklists/:id/fork", bookListAPI.ForkBookList)
			socialGroup.GET("/booklists/:id/books", bookListAPI.GetBooksInList)
		}
	}
}
