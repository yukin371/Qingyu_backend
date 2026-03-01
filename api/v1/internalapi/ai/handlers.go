package ai

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	internalService "Qingyu_backend/service/internal"
)

// 全局服务实例（用于向后兼容的函数式handler）
var documentService *internalService.WriterDraftService

// InitDocumentHandlers 初始化document handlers（用于向后兼容）
func InitDocumentHandlers(service *internalService.WriterDraftService) {
	documentService = service
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
// @Param request body internalService.CreateOrUpdateRequest true "文档信息"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/internal/ai/documents [post]
func (api *DocumentAPI) CreateOrUpdateDocument(c *gin.Context) {
	var req internalService.CreateOrUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		UserID     string   `json:"user_id" binding:"required"`
		ProjectID  string   `json:"project_id" binding:"required"`
		DocumentIDs []string `json:"document_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
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
		UserID     string   `json:"user_id" binding:"required"`
		ProjectID  string   `json:"project_id" binding:"required"`
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

// Concept Handlers (待实现)
func CreateConcept(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func GetConcept(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func UpdateConcept(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func DeleteConcept(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func SearchConcepts(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func BatchGetConcepts(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
