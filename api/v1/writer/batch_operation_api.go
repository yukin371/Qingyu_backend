package writer

import (
	"context"
	"net/http"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/service/writer/document"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchOperationAPI 批量操作API
type BatchOperationAPI struct {
	batchOpSvc document.BatchOperationService
}

// NewBatchOperationAPI 创建批量操作API
func NewBatchOperationAPI(batchOpSvc document.BatchOperationService) *BatchOperationAPI {
	return &BatchOperationAPI{
		batchOpSvc: batchOpSvc,
	}
}

// SubmitBatchOperationRequest 提交批量操作请求
type SubmitBatchOperationRequest struct {
	ProjectID          string                 `json:"projectId" binding:"required"`
	Type               writer.BatchOperationType `json:"type" binding:"required"`
	TargetIDs          []string               `json:"targetIds" binding:"required,min=1"`
	Atomic             bool                   `json:"atomic"`
	Payload            map[string]interface{} `json:"payload"`
	ConflictPolicy     writer.ConflictPolicy  `json:"conflictPolicy"`
	ExpectedVersions   map[string]int         `json:"expectedVersions"`
	ClientRequestID    string                 `json:"clientRequestId"`
	IncludeDescendants bool                   `json:"includeDescendants"`
}

// SubmitBatchOperationResponse 提交批量操作响应
type SubmitBatchOperationResponse struct {
	BatchID          string                    `json:"batchId"`
	Status           writer.BatchOperationStatus `json:"status"`
	PreflightSummary *writer.PreflightSummary  `json:"preflightSummary"`
}

// SubmitBatchOperation 提交批量操作
// @Summary 提交批量操作
// @Description 提交批量删除/移动/复制等操作，包含Preflight预检查
// @Tags batch-operations
// @Accept json
// @Produce json
// @Param request body SubmitBatchOperationRequest true "批量操作请求"
// @Success 200 {object} SubmitBatchOperationResponse
// @Router /api/v1/writer/batch-operations [post]
func (api *BatchOperationAPI) SubmitBatchOperation(c *gin.Context) {
	var req SubmitBatchOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := primitive.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project ID"})
		return
	}

	submitReq := &document.SubmitBatchOperationRequest{
		ProjectID:          projectID,
		Type:               req.Type,
		TargetIDs:          req.TargetIDs,
		Atomic:             req.Atomic,
		Payload:            req.Payload,
		ConflictPolicy:     req.ConflictPolicy,
		ExpectedVersions:   req.ExpectedVersions,
		ClientRequestID:    req.ClientRequestID,
		UserID:             userID.(primitive.ObjectID),
		IncludeDescendants: req.IncludeDescendants,
	}

	batchOp, err := api.batchOpSvc.Submit(c.Request.Context(), submitReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 异步执行
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		api.batchOpSvc.Execute(ctx, batchOp.ID)
	}()

	c.JSON(http.StatusOK, SubmitBatchOperationResponse{
		BatchID:          batchOp.ID.Hex(),
		Status:           batchOp.Status,
		PreflightSummary: batchOp.PreflightSummary,
	})
}

// GetBatchOperationStatus 获取批量操作状态
// @Summary 获取批量操作状态
// @Description 查询批量操作的执行状态和进度
// @Tags batch-operations
// @Accept json
// @Produce json
// @Param id path string true "批量操作ID"
// @Success 200 {object} document.BatchOperationProgress
// @Router /api/v1/writer/batch-operations/{id} [get]
func (api *BatchOperationAPI) GetBatchOperationStatus(c *gin.Context) {
	id := c.Param("id")
	batchID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch operation ID"})
		return
	}

	progress, err := api.batchOpSvc.GetProgress(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch operation not found"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// CancelBatchOperation 取消批量操作
// @Summary 取消批量操作
// @Description 取消正在运行的批量操作
// @Tags batch-operations
// @Accept json
// @Produce json
// @Param id path string true "批量操作ID"
// @Success 200
// @Router /api/v1/writer/batch-operations/{id}/cancel [post]
func (api *BatchOperationAPI) CancelBatchOperation(c *gin.Context) {
	id := c.Param("id")
	batchID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch operation ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = api.batchOpSvc.Cancel(c.Request.Context(), batchID, userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "batch operation cancelled"})
}

// UndoBatchOperation 撤销批量操作
// @Summary 撤销批量操作
// @Description 撤销已完成的批量操作
// @Tags batch-operations
// @Accept json
// @Produce json
// @Param id path string true "批量操作ID"
// @Success 200
// @Router /api/v1/writer/batch-operations/{id}/undo [post]
func (api *BatchOperationAPI) UndoBatchOperation(c *gin.Context) {
	id := c.Param("id")
	batchID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch operation ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err = api.batchOpSvc.Undo(c.Request.Context(), batchID, userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "batch operation undone"})
}

// RegisterRoutes 注册路由
func (api *BatchOperationAPI) RegisterRoutes(r *gin.RouterGroup) {
	batchOps := r.Group("/batch-operations")
	{
		batchOps.POST("", api.SubmitBatchOperation)
		batchOps.GET("/:id", api.GetBatchOperationStatus)
		batchOps.POST("/:id/cancel", api.CancelBatchOperation)
		batchOps.POST("/:id/undo", api.UndoBatchOperation)
	}
}
