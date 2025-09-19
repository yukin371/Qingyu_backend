package document

import (
	"net/http"

	svc "Qingyu_backend/service/document"

	"github.com/gin-gonic/gin"
)

// VersionApi 提供版本相关的 HTTP API
type VersionApi struct {
	service *svc.VersionService
}

func NewVersionApi() *VersionApi {
	return &VersionApi{service: &svc.VersionService{}}
}

// CreateVersion 请求体
type createVersionReq struct {
	AuthorID        string `json:"authorId" binding:"required"`
	Message         string `json:"message"`
	Content         string `json:"content" binding:"required"`
	ExpectedVersion int    `json:"expectedVersion"`
}

// POST /api/v1/document/:nodeId/version
func (a *VersionApi) CreateVersion(c *gin.Context) {
	nodeId := c.Param("nodeId")
	var req createVersionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	rev, err := a.service.UpdateContentWithVersion(c.Param("projectId"), nodeId, req.AuthorID, req.Message, req.Content, req.ExpectedVersion)
	if err != nil {
		if err.Error() == "version_conflict" {
			c.JSON(http.StatusConflict, gin.H{"error": "version_conflict"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rev)
}

// rollback request body
type rollbackReq struct {
	AuthorID string `json:"authorId" binding:"required"`
	Version  int    `json:"version" binding:"required"`
	Message  string `json:"message"`
}

// Rollback 回滚到指定版本
// POST /api/v1/document/:nodeId/rollback
func (a *VersionApi) Rollback(c *gin.Context) {
	nodeId := c.Param("nodeId")
	var req rollbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	rev, err := a.service.RollbackToVersion(c.Param("projectId"), nodeId, req.Version, req.AuthorID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rev)
}

// create patch request
type createPatchReq struct {
	CreatedBy   string `json:"createdBy" binding:"required"`
	DiffFormat  string `json:"diffFormat" binding:"required"`
	DiffPayload string `json:"diffPayload" binding:"required"`
	BaseVersion int    `json:"baseVersion"`
	Message     string `json:"message"`
}

// CreatePatch 创建补丁
// POST /api/v1/document/:nodeId/patch
func (a *VersionApi) CreatePatch(c *gin.Context) {
	nodeId := c.Param("nodeId")
	var req createPatchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	p, err := a.service.CreatePatch(c.Param("projectId"), nodeId, req.BaseVersion, req.DiffFormat, req.DiffPayload, req.CreatedBy, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// ApplyPatch 补丁应用
// POST /api/v1/document/:nodeId/patch/:patchId/apply
func (a *VersionApi) ApplyPatch(c *gin.Context) {
	patchId := c.Param("patchId")
	applier := c.Query("applier")
	if applier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "applier required"})
		return
	}
	rev, err := a.service.ApplyPatch(c.Param("projectId"), patchId, applier)
	if err != nil {
		if err.Error() == "version_conflict" {
			c.JSON(http.StatusConflict, gin.H{"error": "version_conflict"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rev)
}

// ListPatches 补丁列表
// GET /api/v1/document/:nodeId/versions
func (a *VersionApi) ListVersions(c *gin.Context) {
	nodeId := c.Param("nodeId")
	// 简化：返回按时间倒序的最近 50 条
	ctx := c.Request.Context()
	cur, err := a.service.ListRevisions(ctx, c.Param("projectId"), nodeId, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cur)
}
