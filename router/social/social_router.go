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
	collectionAPI *socialApi.CollectionAPI) {

	// 社交路由需要认证
	socialGroup := r.Group("/social")
	socialGroup.Use(middleware.JWTAuth())
	{
		// ========== 用户关注相关 ==========
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
	}
}
