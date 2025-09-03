package system

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册系统相关路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 用户路由组
	userRouter := r.Group("/user")
	{
		// 用户注册
		userRouter.POST("/register", nil) // TODO: 添加用户注册处理函数

		// 用户登录
		userRouter.POST("/login", nil) // TODO: 添加用户登录处理函数
	}
}