package ai

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	internalService "Qingyu_backend/service/internalapi"
)

// 全局服务实例（用于向后兼容的函数式handler）
var documentService *internalService.WriterDraftService
var conceptService *internalService.ConceptService

// 确保 response 包被导入以支持 swagger 文档生成
var _ = response.APIResponse{}

// InitDocumentHandlers 初始化document handlers（用于向后兼容）
func InitDocumentHandlers(service *internalService.WriterDraftService) {
	documentService = service
}

// InitConceptHandlers 初始化concept handlers（用于向后兼容）
func InitConceptHandlers(service *internalService.ConceptService) {
	conceptService = service
}

// DocumentAPI 文档API处理器
type DocumentAPI struct {
	service *internalService.WriterDraftService
}

// NewDocumentAPI 创建DocumentAPI实例
func NewDocumentAPI(service *internalService.WriterDraftService) *DocumentAPI {
	return &DocumentAPI{service: service}
}

// CreateOrUpdateDocument 创建或更新文档
// @Summary 创建或更新文档
// @Description 根据action参数执行创建、更新、创建或更新、追加内容操作
// @Tags Internal-AI-Documents
// @Accept json
// @Produce json
// @Param request body CreateOrUpdateDocumentRequest true "文档信息"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/internal/ai/documents [post]
func (api *DocumentAPI) CreateOrUpdateDocument(c *gin.Context) {
	var req internalService.CreateOrUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Document.ChapterNum == 0 && req.Document.Title == "" && req.Document.Content == "" && req.Document.Format == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document is required"})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	doc, err := api.service.CreateOrUpdate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "操作成功",
		"data":    doc,
	})
}

// GetDocument 获取文档
// @Summary 获取文档
// @Description 根据ID获取文档详情
// @Tags Internal-AI-Documents
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param user_id query string true "用户ID"
// @Param project_id query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/internal/ai/documents/{id} [get]
func (api *DocumentAPI) GetDocument(c *gin.Context) {
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	docID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "user_id和project_id参数必填",
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	doc, err := api.service.GetDocument(c.Request.Context(), userID, projectID, docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文档不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    doc,
	})
}

// ListDocuments 列出文档
// @Summary 列出文档
// @Description 获取项目下的文档列表
// @Tags Internal-AI-Documents
// @Accept json
// @Produce json
// @Param user_id query string true "用户ID"
// @Param project_id query string true "项目ID"
// @Param limit query int false "每页数量" default(50)
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/internal/ai/documents [get]
func (api *DocumentAPI) ListDocuments(c *gin.Context) {
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "user_id和project_id参数必填",
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	docs, total, err := api.service.ListDocuments(c.Request.Context(), userID, projectID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文档列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"documents": docs,
			"total":     total,
		},
	})
}

// DeleteDocument 删除文档
// @Summary 删除文档
// @Description 根据ID删除文档
// @Tags Internal-AI-Documents
// @Accept json
// @Produce json
// @Param id path string true "文档ID"
// @Param user_id query string true "用户ID"
// @Param project_id query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/internal/ai/documents/{id} [delete]
func (api *DocumentAPI) DeleteDocument(c *gin.Context) {
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	docID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "user_id和project_id参数必填",
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	if err := api.service.DeleteDocument(c.Request.Context(), userID, projectID, docID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除文档失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data": gin.H{
			"success": true,
		},
	})
}

// BatchGetDocuments 批量获取文档
// @Summary 批量获取文档
// @Description 根据ID列表批量获取文档
// @Tags Internal-AI-Documents
// @Accept json
// @Produce json
// @Param request body object{user_id:string,project_id:string,document_ids:[]string} true "批量请求参数"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/internal/ai/documents/batch [post]
func (api *DocumentAPI) BatchGetDocuments(c *gin.Context) {
	var req struct {
		UserID      string   `json:"user_id" binding:"required"`
		ProjectID   string   `json:"project_id" binding:"required"`
		DocumentIDs []string `json:"document_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	docs, err := api.service.BatchGetDocuments(c.Request.Context(), req.UserID, req.ProjectID, req.DocumentIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量获取文档失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"documents": docs,
			"total":     len(docs),
		},
	})
}

// ============ 向后兼容的函数式Handler ============
// 这些函数用于与现有的enter.go路由注册兼容

// CreateOrUpdateDocument 创建或更新文档（函数式handler）
func CreateOrUpdateDocument(c *gin.Context) {
	if documentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	var req internalService.CreateOrUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := documentService.CreateOrUpdate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// GetDocument 获取文档（函数式handler）
func GetDocument(c *gin.Context) {
	if documentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	docID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and project_id required"})
		return
	}

	doc, err := documentService.GetDocument(c.Request.Context(), userID, projectID, docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ListDocuments 列出文档（函数式handler）
func ListDocuments(c *gin.Context) {
	if documentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and project_id required"})
		return
	}

	docs, total, err := documentService.ListDocuments(c.Request.Context(), userID, projectID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": docs,
		"total":     total,
	})
}

// DeleteDocument 删除文档（函数式handler）
func DeleteDocument(c *gin.Context) {
	if documentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	docID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and project_id required"})
		return
	}

	if err := documentService.DeleteDocument(c.Request.Context(), userID, projectID, docID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// BatchGetDocuments 批量获取文档（函数式handler）
func BatchGetDocuments(c *gin.Context) {
	if documentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	var req struct {
		UserID      string   `json:"user_id" binding:"required"`
		ProjectID   string   `json:"project_id" binding:"required"`
		DocumentIDs []string `json:"document_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	docs, err := documentService.BatchGetDocuments(c.Request.Context(), req.UserID, req.ProjectID, req.DocumentIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": docs,
		"total":     len(docs),
	})
}

// ============ Concept Handlers ============

// ConceptAPI 设定概念API处理器
type ConceptAPI struct {
	service *internalService.ConceptService
}

// NewConceptAPI 创建ConceptAPI实例
func NewConceptAPI(service *internalService.ConceptService) *ConceptAPI {
	return &ConceptAPI{service: service}
}

// CreateConcept 创建概念
// @Summary 创建设定概念
// @Description AI助手创建新的设定概念（角色、地点、魔法等）
// @Tags Internal-AI-Concepts
// @Accept json
// @Produce json
// @Param request body CreateConceptRequest true "创建请求"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/internal/ai/concepts [post]
func (api *ConceptAPI) CreateConcept(c *gin.Context) {
	var req internalService.CreateConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	concept, err := api.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建概念失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    concept,
	})
}

// GetConcept 获取概念
// @Summary 获取设定概念
// @Description 根据ID获取设定概念详情
// @Tags Internal-AI-Concepts
// @Accept json
// @Produce json
// @Param id path string true "概念ID"
// @Param user_id query string true "用户ID"
// @Param project_id query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /api/v1/internal/ai/concepts/{id} [get]
func (api *ConceptAPI) GetConcept(c *gin.Context) {
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	conceptID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "user_id和project_id参数必填",
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	concept, err := api.service.GetConcept(c.Request.Context(), userID, projectID, conceptID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "概念不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    concept,
	})
}

// UpdateConcept 更新概念
// @Summary 更新设定概念
// @Description 更新现有的设定概念
// @Tags Internal-AI-Concepts
// @Accept json
// @Produce json
// @Param id path string true "概念ID"
// @Param request body UpdateConceptRequest true "更新请求"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/internal/ai/concepts/{id} [put]
func (api *ConceptAPI) UpdateConcept(c *gin.Context) {
	conceptID := c.Param("id")
	var req internalService.UpdateConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	concept, err := api.service.Update(c.Request.Context(), conceptID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新概念失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    concept,
	})
}

// DeleteConcept 删除概念
// @Summary 删除设定概念
// @Description 根据ID删除设定概念
// @Tags Internal-AI-Concepts
// @Accept json
// @Produce json
// @Param id path string true "概念ID"
// @Param user_id query string true "用户ID"
// @Param project_id query string true "项目ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/internal/ai/concepts/{id} [delete]
func (api *ConceptAPI) DeleteConcept(c *gin.Context) {
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	conceptID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "user_id和project_id参数必填",
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	if err := api.service.Delete(c.Request.Context(), userID, projectID, conceptID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除概念失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
		"data": gin.H{
			"success": true,
		},
	})
}

// SearchConcepts 搜索概念
// @Summary 搜索设定概念
// @Description 按分类和关键词搜索设定概念
// @Tags Internal-AI-Concepts
// @Accept json
// @Produce json
// @Param user_id query string true "用户ID"
// @Param project_id query string true "项目ID"
// @Param category query string false "分类"
// @Param keyword query string false "关键词"
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/internal/ai/concepts [get]
func (api *ConceptAPI) SearchConcepts(c *gin.Context) {
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	category := c.Query("category")
	keyword := c.Query("keyword")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "user_id和project_id参数必填",
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	concepts, total, err := api.service.Search(c.Request.Context(), userID, projectID, category, keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索概念失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索成功",
		"data": gin.H{
			"concepts": concepts,
			"total":    total,
		},
	})
}

// BatchGetConcepts 批量获取概念
// @Summary 批量获取设定概念
// @Description 根据ID列表批量获取设定概念
// @Tags Internal-AI-Concepts
// @Accept json
// @Produce json
// @Param request body object{user_id:string,project_id:string,concept_ids:[]string} true "批量请求参数"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/internal/ai/concepts/batch [post]
func (api *ConceptAPI) BatchGetConcepts(c *gin.Context) {
	var req struct {
		UserID     string   `json:"user_id" binding:"required"`
		ProjectID  string   `json:"project_id" binding:"required"`
		ConceptIDs []string `json:"concept_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}
	if api.service == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}

	concepts, err := api.service.BatchGet(c.Request.Context(), req.UserID, req.ProjectID, req.ConceptIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量获取概念失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"concepts": concepts,
			"total":    len(concepts),
		},
	})
}

// ============ 向后兼容的函数式Handler ============
// 用于与现有的enter.go路由注册兼容

// CreateConcept 创建概念（函数式handler）
func CreateConcept(c *gin.Context) {
	if conceptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	var req internalService.CreateConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	concept, err := conceptService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, concept)
}

// GetConcept 获取概念（函数式handler）
func GetConcept(c *gin.Context) {
	if conceptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	conceptID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and project_id required"})
		return
	}

	concept, err := conceptService.GetConcept(c.Request.Context(), userID, projectID, conceptID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "concept not found"})
		return
	}

	c.JSON(http.StatusOK, concept)
}

// UpdateConcept 更新概念（函数式handler）
func UpdateConcept(c *gin.Context) {
	if conceptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	conceptID := c.Param("id")
	var req internalService.UpdateConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	concept, err := conceptService.Update(c.Request.Context(), conceptID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, concept)
}

// DeleteConcept 删除概念（函数式handler）
func DeleteConcept(c *gin.Context) {
	if conceptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	conceptID := c.Param("id")

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and project_id required"})
		return
	}

	if err := conceptService.Delete(c.Request.Context(), userID, projectID, conceptID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SearchConcepts 搜索概念（函数式handler）
func SearchConcepts(c *gin.Context) {
	if conceptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	userID := c.Query("user_id")
	projectID := c.Query("project_id")
	category := c.Query("category")
	keyword := c.Query("keyword")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	if userID == "" || projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and project_id required"})
		return
	}

	concepts, total, err := conceptService.Search(c.Request.Context(), userID, projectID, category, keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"concepts": concepts,
		"total":    total,
	})
}

// BatchGetConcepts 批量获取概念（函数式handler）
func BatchGetConcepts(c *gin.Context) {
	if conceptService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务未初始化",
		})
		return
	}
	var req struct {
		UserID     string   `json:"user_id" binding:"required"`
		ProjectID  string   `json:"project_id" binding:"required"`
		ConceptIDs []string `json:"concept_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	concepts, err := conceptService.BatchGet(c.Request.Context(), req.UserID, req.ProjectID, req.ConceptIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"concepts": concepts,
		"total":    len(concepts),
	})
}

// ========== Swagger 请求结构体定义 ==========

// CreateOrUpdateDocumentRequest 创建或更新文档请求（用于swagger文档）
type CreateOrUpdateDocumentRequest struct {
	UserID    string              `json:"user_id" binding:"required" example:"user123"`
	ProjectID string              `json:"project_id" binding:"required" example:"project456"`
	Action    string              `json:"action" binding:"required" example:"create"` // create, update, create_or_update, append
	Document  WriterDraftDocument `json:"document" binding:"required"`
}

// WriterDraftDocument 文档数据（用于swagger文档）
type WriterDraftDocument struct {
	ChapterNum int    `json:"chapter_num" example:"1"`
	Title      string `json:"title" example:"第一章 开始"`
	Content    string `json:"content" example:"文档内容..."`
	Format     string `json:"format" example:"markdown"`
}

// CreateConceptRequest 创建概念请求（用于swagger文档）
type CreateConceptRequest struct {
	UserID    string   `json:"user_id" binding:"required" example:"user123"`
	ProjectID string   `json:"project_id" binding:"required" example:"project456"`
	Name      string   `json:"name" binding:"required" example:"主角"`
	Category  string   `json:"category" example:"character"`
	Content   string   `json:"content" example:"主角的详细描述..."`
	Tags      []string `json:"tags" example:"主角,重要"`
}

// UpdateConceptRequest 更新概念请求（用于swagger文档）
type UpdateConceptRequest struct {
	UserID    string   `json:"user_id" binding:"required" example:"user123"`
	ProjectID string   `json:"project_id" binding:"required" example:"project456"`
	Name      string   `json:"name" example:"主角"`
	Content   string   `json:"content" example:"更新后的描述..."`
	Tags      []string `json:"tags" example:"主角,重要,更新"`
}
