package document

import (
	"net/http"

	model "Qingyu_backend/models/document"
	svc "Qingyu_backend/service/document"

	"github.com/gin-gonic/gin"
)

// DocumentApi 文档相关API
type DocumentApi struct {
	service *svc.DocumentService
}

// NewDocumentApi 创建文档相关API实例
func NewDocumentApi() *DocumentApi {
	return &DocumentApi{service: &svc.DocumentService{}}
}

// Create 创建文档
func (a *DocumentApi) Create(c *gin.Context) {
	var req model.Document
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	created, err := a.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// List 获取文档列表
func (a *DocumentApi) List(c *gin.Context) {
	// 简化：从查询参数取 userId、limit、offset
	userID := c.Query("userId")
	// 省略参数校验与转换细节，默认分页
	docs, err := a.service.List(userID, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

// Get 获取单个文档
func (a *DocumentApi) Get(c *gin.Context) {
	id := c.Param("id")
	doc, err := a.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	c.JSON(http.StatusOK, doc)
}

// Update 更新文档
func (a *DocumentApi) Update(c *gin.Context) {
	id := c.Param("id")
	var req model.Document
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	updated, err := a.service.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete 删除文档
func (a *DocumentApi) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
