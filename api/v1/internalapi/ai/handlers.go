package ai

import (
	"github.com/gin-gonic/gin"
)

// TODO: 后续任务中实现这些handler函数
// 这里先定义函数签名，用于路由注册

// Document Handlers
func CreateOrUpdateDocument(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func GetDocument(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func ListDocuments(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func DeleteDocument(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

func BatchGetDocuments(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// Concept Handlers
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
