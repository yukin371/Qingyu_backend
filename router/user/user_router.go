package user

import (
	"github.com/gin-gonic/gin"

	userApi "Qingyu_backend/api/v1/user"
	"Qingyu_backend/middleware"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.RouterGroup, userService serviceInterfaces.UserService) {
	// 创建UserAPI
	userAPI := userApi.NewUserAPI(userService)

	// =========================
	// 公开路由（无需认证）
	// =========================
	{
		// 用户注册
		r.POST("/register", userAPI.Register)

		// 用户登录
		r.POST("/login", userAPI.Login)
	}

	// =========================
	// 需要认证的路由
	// =========================
	authenticated := r.Group("/users")
	authenticated.Use(middleware.JWTAuth()) // 启用JWT认证中间件
	{
		// 获取当前用户信息
		authenticated.GET("/profile", userAPI.GetProfile)

		// 更新当前用户信息
		authenticated.PUT("/profile", userAPI.UpdateProfile)

		// 修改密码
		authenticated.PUT("/password", userAPI.ChangePassword)
	}
}

