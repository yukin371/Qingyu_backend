package handlers

import (
	"context"
	"net/http"
	"time"

	"Qingyu_backend/database"
	"Qingyu_backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 文档集合名称
const documentCollection = "documents"

// DocumentRequest 文档请求结构
type DocumentRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreateDocumentHandler 创建新文档
func CreateDocumentHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析请求
	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 创建文档
	now := time.Now()
	document := models.Document{
		UserID:    userID.(string),
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
		Versions:  []models.DocumentVersion{},
	}

	// 保存到数据库
	coll := database.GetCollection(documentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文档创建失败"})
		return
	}

	// 获取插入的ID
	document.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "文档创建成功",
		"document": document,
	})
}

// GetDocumentsHandler 获取用户的所有文档
func GetDocumentsHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询数据库
	coll := database.GetCollection(documentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建查询选项，按更新时间降序排序
	opts := options.Find().SetSort(bson.D{{"updated_at", -1}})

	// 执行查询
	cursor, err := coll.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文档失败"})
		return
	}
	defer cursor.Close(ctx)

	// 解码结果
	var documents []models.Document
	if err := cursor.All(ctx, &documents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析文档失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
	})
}

// GetDocumentHandler 获取单个文档
func GetDocumentHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文档ID
	docID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文档ID"})
		return
	}

	// 查询数据库
	coll := database.GetCollection(documentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查找文档
	var document models.Document
	err = coll.FindOne(ctx, bson.M{"_id": docID, "user_id": userID}).Decode(&document)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文档失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"document": document,
	})
}

// UpdateDocumentHandler 更新文档
func UpdateDocumentHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文档ID
	docID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文档ID"})
		return
	}

	// 解析请求
	var req DocumentRequest
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 查询数据库
	coll := database.GetCollection(documentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取当前文档
	var currentDoc models.Document
	err = coll.FindOne(ctx, bson.M{"_id": docID, "user_id": userID}).Decode(&currentDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文档失败"})
		}
		return
	}

	// 创建版本记录
	version := models.DocumentVersion{
		Title:     currentDoc.Title,
		Content:   currentDoc.Content,
		Time:      currentDoc.UpdatedAt,
		VersionID: primitive.NewObjectID(),
	}

	// 更新文档
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":      req.Title,
			"content":    req.Content,
			"updated_at": now,
		},
		"$push": bson.M{
			"versions": version,
		},
	}

	// 执行更新
	_, err = coll.UpdateOne(ctx, bson.M{"_id": docID, "user_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文档失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文档更新成功",
	})
}

// DeleteDocumentHandler 删除文档
func DeleteDocumentHandler(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取文档ID
	docID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文档ID"})
		return
	}

	// 查询数据库
	coll := database.GetCollection(documentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行删除
	result, err := coll.DeleteOne(ctx, bson.M{"_id": docID, "user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文档失败"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在或无权限删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文档删除成功",
	})
}
