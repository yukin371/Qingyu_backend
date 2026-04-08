package writer

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apiWriter "Qingyu_backend/api/v1/writer"
	"Qingyu_backend/service/interfaces"
)

// InitEntityRoutes 初始化统一实体路由
func InitEntityRoutes(router *gin.RouterGroup, entityService interfaces.EntityService) {
	zap.L().Info("InitEntityRoutes: 开始注册实体路由")

	api := apiWriter.NewEntityApi(entityService)

	// 项目级别的实体路由
	projectGroup := router.Group("/projects/:projectId")
	{
		projectGroup.GET("/entities", api.ListEntities)
		projectGroup.GET("/entities/graph", api.GetEntityGraph)

		zap.L().Info("InitEntityRoutes: 项目级实体路由已注册到 /projects/:projectId/entities")
	}

	// 实体级别的路由
	entityGroup := router.Group("/entities")
	{
		entityGroup.PUT("/:entityId/state-fields", api.UpdateEntityStateFields)

		zap.L().Info("InitEntityRoutes: 实体级路由已注册到 /entities")
	}

	zap.L().Info("InitEntityRoutes: 实体路由注册完成")
}
