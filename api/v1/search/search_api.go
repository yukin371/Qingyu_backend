package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchAPI 统一搜索 API
type SearchAPI struct {
	// TODO: 注入 SearchService
}

// NewSearchAPI 创建搜索 API 实例
func NewSearchAPI() *SearchAPI {
	return &SearchAPI{}
}

// SearchBooks 书籍搜索
// GET /api/v1/search/books?q={query}&page=1&page_size=20
func (api *SearchAPI) SearchBooks(c *gin.Context) {
	// TODO: 实现书籍搜索
	c.JSON(http.StatusOK, gin.H{
		"message": "Search books endpoint - to be implemented",
	})
}

// SearchProjects 项目搜索
// GET /api/v1/search/projects?q={query}&page=1&page_size=20
func (api *SearchAPI) SearchProjects(c *gin.Context) {
	// TODO: 实现项目搜索
	c.JSON(http.StatusOK, gin.H{
		"message": "Search projects endpoint - to be implemented",
	})
}

// SearchDocuments 文档搜索
// GET /api/v1/search/documents?q={query}&project_id={id}&page=1&page_size=20
func (api *SearchAPI) SearchDocuments(c *gin.Context) {
	// TODO: 实现文档搜索
	c.JSON(http.StatusOK, gin.H{
		"message": "Search documents endpoint - to be implemented",
	})
}

// SearchUsers 用户搜索
// GET /api/v1/search/users?q={query}&page=1&page_size=20
func (api *SearchAPI) SearchUsers(c *gin.Context) {
	// TODO: 实现用户搜索
	c.JSON(http.StatusOK, gin.H{
		"message": "Search users endpoint - to be implemented",
	})
}
