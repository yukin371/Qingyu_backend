package users

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/system"
	serviceInterfaces "Qingyu_backend/service/interfaces"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.RouterGroup, userService serviceInterfaces.UserService) {
	// 创建UserAPI
	userAPI := system.NewUserAPI(userService)

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
	authenticated := r.Group("")
	// TODO: 添加JWT认证中间件
	// authenticated.Use(middleware.JWTAuth())
	{
		// 获取当前用户信息
		authenticated.GET("/users/profile", userAPI.GetProfile)

		// 更新当前用户信息
		authenticated.PUT("/users/profile", userAPI.UpdateProfile)

		// 修改密码
		authenticated.PUT("/users/password", userAPI.ChangePassword)
	}

	// =========================
	// 管理员路由
	// =========================
	admin := r.Group("/admin/users")
	// TODO: 添加JWT认证中间件和管理员权限中间件
	// admin.Use(middleware.JWTAuth())
	// admin.Use(middleware.AdminPermission())
	{
		// 获取用户列表
		admin.GET("", userAPI.ListUsers)

		// 获取指定用户信息
		admin.GET("/:id", userAPI.GetUser)

		// 更新指定用户信息
		admin.PUT("/:id", userAPI.UpdateUser)

		// 删除用户
		admin.DELETE("/:id", userAPI.DeleteUser)
	}
}
