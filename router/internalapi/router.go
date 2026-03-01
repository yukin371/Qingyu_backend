package internalapi

import (
	"github.com/gin-gonic/gin"

	internalService "Qingyu_backend/service/internalapi"
	"Qingyu_backend/api/v1/internalapi/ai"
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
