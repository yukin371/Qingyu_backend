package project

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"Qingyu_backend/api/v1/writer"
	"Qingyu_backend/middleware"
	"Qingyu_backend/service"
	projectService "Qingyu_backend/service/project"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

// RegisterRoutes 注册项目管理路由到 /api/v1/projects
func RegisterRoutes(r *gin.RouterGroup) {
	logger.Info("开始注册项目路由...")

	// 从服务容器获取依赖
	serviceContainer := service.GetServiceContainer()
	if serviceContainer == nil {
		logger.Warn("服务容器未初始化，跳过项目路由注册")
		return
	}
	logger.Info("服务容器已获取")

	// 获取Repository工厂和EventBus
	repositoryFactory := serviceContainer.GetRepositoryFactory()
	if repositoryFactory == nil {
		logger.Warn("Repository工厂未初始化，跳过项目路由注册")
		return
	}
	logger.Info("Repository工厂已获取")

	eventBus := serviceContainer.GetEventBus()
	logger.Info("EventBus已获取")

	// 创建ProjectRepository
	projectRepo := repositoryFactory.CreateProjectRepository()
	logger.Info("ProjectRepository已创建")

	// 创建ProjectService
	projectSvc := projectService.NewProjectService(projectRepo, eventBus)
	logger.Info("ProjectService已创建")

	// 创建ProjectApi实例
	projectApi := writer.NewProjectApi(projectSvc)
	logger.Info("ProjectApi已创建")

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

	logger.Info("✓ 项目路由已成功注册到: /api/v1/projects/")
}
