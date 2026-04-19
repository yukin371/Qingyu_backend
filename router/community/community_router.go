package community

import (
	"github.com/gin-gonic/gin"

	communityApi "Qingyu_backend/api/v1/community"
	"Qingyu_backend/internal/middleware/auth"
)

// RegisterCommunityRoutes 注册社区动态相关路由
func RegisterCommunityRoutes(r *gin.RouterGroup, postAPI *communityApi.PostAPI) {
	// 社区路由组
	communityGroup := r.Group("/community")
	{
		// 动态列表（公开访问）
		communityGroup.GET("/posts", postAPI.GetPosts)

		// 动态详情（公开访问）
		communityGroup.GET("/posts/:id", postAPI.GetPostDetail)

		// 以下接口需要认证
		authGroup := communityGroup.Group("")
		authGroup.Use(auth.JWTAuth())
		{
			// 创建动态
			authGroup.POST("/posts", postAPI.CreatePost)

			// 更新动态
			authGroup.PUT("/posts/:id", postAPI.UpdatePost)
			authGroup.PATCH("/posts/:id", postAPI.UpdatePost)

			// 删除动态
			authGroup.DELETE("/posts/:id", postAPI.DeletePost)

			// 点赞动态
			authGroup.POST("/posts/:id/like", postAPI.LikePost)

			// 取消点赞
			authGroup.DELETE("/posts/:id/like", postAPI.UnlikePost)

			// 切换点赞状态
			authGroup.PATCH("/posts/:id/like", postAPI.ToggleLikePost)
		}
	}
}
