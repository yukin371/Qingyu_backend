package project

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service"
	projectService "Qingyu_backend/service/project"
)

// RegisterRoutes 注册项目管理路由到 /api/v1/projects
func RegisterRoutes(r *gin.RouterGroup) {
	// 从服务容器获取依赖
	serviceContainer := service.GetServiceContainer()
	if serviceContainer == nil {
		// 如果服务容器未初始化，跳过路由注册
		return
	}

	// 获取Repository工厂和EventBus
	repositoryFactory := serviceContainer.GetRepositoryFactory()
	if repositoryFactory == nil {
		return
	}

	eventBus := serviceContainer.GetEventBus()

	// 创建ProjectRepository
	projectRepo := repositoryFactory.CreateProjectRepository()

	// 创建ProjectService
	projectSvc := projectService.NewProjectService(projectRepo, eventBus)

	// 创建ProjectApi实例
	projectApi := writer.NewProjectApi(projectSvc)

	// 项目管理路由组 - 需要JWT认证
	projectGroup := r.Group("/projects")
	projectGroup.Use(middleware.JWTAuth())
	{
		// 项目CRUD
		projectGroup.POST("", projectApi.CreateProject)       // 创建项目
		projectGroup.GET("", projectApi.ListProjects)         // 获取项目列表
		projectGroup.GET("/:id", projectApi.GetProject)       // 获取项目详情
		projectGroup.PUT("/:id", projectApi.UpdateProject)    // 更新项目
		projectGroup.DELETE("/:id", projectApi.DeleteProject) // 删除项目

		// 项目统计
		projectGroup.PUT("/:id/statistics", projectApi.UpdateProjectStatistics) // 更新项目统计
	}
}
