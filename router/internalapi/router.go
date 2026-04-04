package internalapi

import (
	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/internalapi/ai"
	"Qingyu_backend/api/v1/internalapi/context"
	internalService "Qingyu_backend/service/internalapi"
)

// RegisterInternalAPIRoutes 注册内部AI API路由
// 供AI服务内部调用，需要AI服务认证中间件
func RegisterInternalAPIRoutes(routerGroup *gin.RouterGroup, draftService *internalService.WriterDraftService, conceptService *internalService.ConceptService) {
	// 创建DocumentAPI实例
	documentAPI := ai.NewDocumentAPI(draftService)

	// 文档管理路由
	routerGroup.POST("/documents", documentAPI.CreateOrUpdateDocument)
	routerGroup.GET("/documents/:id", documentAPI.GetDocument)
	routerGroup.GET("/documents", documentAPI.ListDocuments)
	routerGroup.DELETE("/documents/:id", documentAPI.DeleteDocument)
	routerGroup.POST("/documents/batch", documentAPI.BatchGetDocuments)

	// 创建ConceptAPI实例
	conceptAPI := ai.NewConceptAPI(conceptService)

	// 概念管理路由
	routerGroup.POST("/concepts", conceptAPI.CreateConcept)
	routerGroup.GET("/concepts/:id", conceptAPI.GetConcept)
	routerGroup.PUT("/concepts/:id", conceptAPI.UpdateConcept)
	routerGroup.DELETE("/concepts/:id", conceptAPI.DeleteConcept)
	routerGroup.GET("/concepts", conceptAPI.SearchConcepts)
	routerGroup.POST("/concepts/batch", conceptAPI.BatchGetConcepts)
}

// RegisterContextRoutes 注册内部上下文API路由
// 供AI服务获取项目上下文数据（角色、大纲、文档内容、角色关系）
func RegisterContextRoutes(v1 *gin.RouterGroup, aggregator *internalService.ContextAggregator) {
	// 使用与现有内部AI API相同的认证中间件
	internalGroup := v1.Group("/internal")
	internalGroup.Use(internalService.AIAuthMiddleware())

	contextAPI := context.NewContextAPI(aggregator)

	// 项目上下文路由
	internalGroup.GET("/projects/:id/context", contextAPI.GetProjectContext)
	internalGroup.GET("/projects/:id/characters", contextAPI.GetCharacters)
	internalGroup.GET("/projects/:id/outline", contextAPI.GetOutline)
	internalGroup.GET("/projects/:id/relations", contextAPI.GetCharacterRelations)

	// 文档内容路由
	internalGroup.GET("/documents/:id/content", contextAPI.GetDocumentContent)
}
