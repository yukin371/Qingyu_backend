package user

import (
	serviceInterfaces "Qingyu_backend/service/interfaces/user"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
// ⚠️ 注意：用户管理功能已迁移到 user-management 模块
// 新路由位置：/api/v1/user-management/*
// 此函数保留以向后兼容，但不再注册任何路由
func RegisterUserRoutes(r *gin.RouterGroup, userService serviceInterfaces.UserService) {
	// 用户管理功能已迁移到新的 user-management 模块
	// 新路由在 router/enter.go 中通过 usermanagementRouter.RegisterRoutes() 注册
	// 请使用新的路由：
	// - POST   /api/v1/user-management/auth/register
	// - POST   /api/v1/user-management/auth/login
	// - GET    /api/v1/user-management/profile
	// - PUT    /api/v1/user-management/profile
	// - PUT    /api/v1/user-management/password
	// - GET    /api/v1/user-management/users/:id/profile
	// - GET    /api/v1/user-management/users/:id/books
}

// RegisterUserRoutesWithBookstore 注册用户相关路由（包含书店服务）
// ⚠️ 注意：用户管理功能已迁移到 user-management 模块
// 新路由位置：/api/v1/user-management/*
// 此函数保留以向后兼容，但不再注册任何路由
func RegisterUserRoutesWithBookstore(r *gin.RouterGroup, userService serviceInterfaces.UserService, bookstoreService interface{}) {
	// 用户管理功能已迁移到新的 user-management 模块
	// 新路由在 router/enter.go 中通过 usermanagementRouter.RegisterRoutes() 注册
	// 请使用新的路由：
	// - POST   /api/v1/user-management/auth/register
	// - POST   /api/v1/user-management/auth/login
	// - GET    /api/v1/user-management/profile
	// - PUT    /api/v1/user-management/profile
	// - PUT    /api/v1/user-management/password
	// - GET    /api/v1/user-management/users/:id/profile
	// - GET    /api/v1/user-management/users/:id/books
}
