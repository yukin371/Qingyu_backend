package document

import (
	"net/http"
	"strconv"

	model "Qingyu_backend/models/document"
	svc "Qingyu_backend/service/project"

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

// BatchCommitRequest 批量提交请求体
type BatchCommitRequest struct {
	AuthorID string `json:"authorId" binding:"required"`
	Message  string `json:"message"`
	Files    []struct {
		NodeID          string `json:"nodeId" binding:"required"`
		Content         string `json:"content" binding:"required"`
		ExpectedVersion int    `json:"expectedVersion"`
	} `json:"files" binding:"required,min=1"`
}

// CreateCommit 创建批量提交
// POST /api/v1/document/:projectId/commit
func (a *VersionApi) CreateCommit(c *gin.Context) {
	projectID := c.Param("projectId")
	var req BatchCommitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	
	// 转换请求格式为model.CommitFile
	var files []model.CommitFile
	for _, file := range req.Files {
		files = append(files, model.CommitFile{
			NodeID:          file.NodeID,
			Content:         file.Content,
			ExpectedVersion: file.ExpectedVersion,
		})
	}

	commit, err := a.service.CreateCommit(ctx, projectID, req.AuthorID, req.Message, files)
	if err != nil {
		if err.Error() == "commit_conflicts_detected" {
			c.JSON(http.StatusConflict, gin.H{"error": "commit_conflicts_detected"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, commit)
}

// ListCommits 查询提交历史
// GET /api/v1/document/:projectId/commits
func (a *VersionApi) ListCommits(c *gin.Context) {
	projectID := c.Param("projectId")
	authorID := c.Query("authorId")
	
	// 简化分页参数处理
	limit := int64(50)
	offset := int64(0)

	ctx := c.Request.Context()
	commits, err := a.service.ListCommits(ctx, projectID, authorID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commits)
}

// GetCommitDetails 获取提交详情
// GET /api/v1/document/:projectId/commit/:commitId
func (a *VersionApi) GetCommitDetails(c *gin.Context) {
	projectID := c.Param("projectId")
	commitID := c.Param("commitId")

	ctx := c.Request.Context()
	commit, revisions, err := a.service.GetCommitDetails(ctx, projectID, commitID)
	if err != nil {
		if err.Error() == "commit_not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "commit_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"commit":    commit,
		"revisions": revisions,
	})
}

// DetectConflicts 检测版本冲突
// POST /api/v1/document/:projectId/conflicts
func (a *VersionApi) DetectConflicts(c *gin.Context) {
	projectID := c.Param("projectId")
	
	var req struct {
		Files []struct {
			NodeID          string `json:"nodeId" binding:"required"`
			ExpectedVersion int    `json:"expectedVersion"`
		} `json:"files" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()
	
	// 转换请求格式
	var files []struct {
		NodeID          string `json:"node_id"`
		ExpectedVersion int    `json:"expected_version"`
	}
	for _, file := range req.Files {
		files = append(files, struct {
			NodeID          string `json:"node_id"`
			ExpectedVersion int    `json:"expected_version"`
		}{
			NodeID:          file.NodeID,
			ExpectedVersion: file.ExpectedVersion,
		})
	}

	result, err := a.service.BatchDetectConflicts(ctx, projectID, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateVersion 创建新版本
func (api *VersionApi) CreateNewVersion(c *gin.Context) {
	var req struct {
		ProjectID string `json:"projectId" binding:"required"`
		NodeID    string `json:"nodeId" binding:"required"`
		AuthorID  string `json:"authorId" binding:"required"`
		Message   string `json:"message"`
		Content   string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	revision, err := api.service.BumpVersionAndCreateRevision(req.ProjectID, req.NodeID, req.AuthorID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": revision})
}

// UpdateVersion 更新版本内容（乐观锁）
func (api *VersionApi) UpdateVersion(c *gin.Context) {
	var req struct {
		ProjectID       string `json:"projectId" binding:"required"`
		NodeID          string `json:"nodeId" binding:"required"`
		AuthorID        string `json:"authorId" binding:"required"`
		Message         string `json:"message"`
		Content         string `json:"content" binding:"required"`
		ExpectedVersion int    `json:"expectedVersion" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	revision, err := api.service.UpdateContentWithVersion(req.ProjectID, req.NodeID, req.AuthorID, req.Message, req.Content, req.ExpectedVersion)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": revision})
}

// GetRevisions 获取文件修订历史
func (api *VersionApi) GetRevisions(c *gin.Context) {
	projectID := c.Query("projectId")
	nodeID := c.Query("nodeId")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	if projectID == "" || nodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "projectId and nodeId are required"})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	revisions, err := api.service.ListRevisions(c.Request.Context(), projectID, nodeID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": revisions})
}

// ResolveBatchConflicts 批量解决冲突
func (api *VersionApi) ResolveBatchConflicts(c *gin.Context) {
	var req model.BatchConflictResolution

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commit, err := api.service.ResolveBatchConflicts(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": commit})
}

// AutoResolveConflicts 自动解决冲突
func (api *VersionApi) AutoResolveConflicts(c *gin.Context) {
	var req struct {
		ProjectID string `json:"projectId" binding:"required"`
		NodeID    string `json:"nodeId" binding:"required"`
		AuthorID  string `json:"authorId" binding:"required"`
		Message   string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 首先检测冲突
	conflict, err := api.service.DetectConflicts(c.Request.Context(), req.ProjectID, req.NodeID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !conflict.HasConflict {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No conflicts to resolve"})
		return
	}

	// 尝试自动解决
	resolvedContent, err := api.service.AutoResolveConflicts(c.Request.Context(), req.ProjectID, req.NodeID, conflict.ConflictingRevisions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建解决方案并应用
	resolution := &model.BatchConflictResolution{
		ProjectID: req.ProjectID,
		AuthorID:  req.AuthorID,
		Message:   req.Message,
		Resolutions: map[string]*model.ConflictResolution{
			req.NodeID: {
				Strategy:      "auto",
				ResolvedBy:    req.AuthorID,
				Resolution:    "Auto-merged conflicting changes",
				MergedContent: resolvedContent,
			},
		},
	}

	commit, err := api.service.ResolveBatchConflicts(c.Request.Context(), resolution)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": commit})
}

// GET /api/v1/document/:nodeId/version/current
func (a *VersionApi) GetCurrentVersion(c *gin.Context) {
	projectId := c.Param("projectId")
	nodeId := c.Param("nodeId")
	
	version, err := a.service.GetCurrentVersion(c.Request.Context(), projectId, nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"data": version})
}

// rollback request body
