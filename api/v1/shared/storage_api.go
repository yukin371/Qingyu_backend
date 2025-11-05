package shared

import (
	"Qingyu_backend/service/shared/storage"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// StorageAPI 文件存储API处理器
type StorageAPI struct {
	storageService   *storage.StorageServiceImpl
	multipartService *storage.MultipartUploadService
	imageProcessor   *storage.ImageProcessor
}

// NewStorageAPI 创建存储API实例
func NewStorageAPI(
	storageService *storage.StorageServiceImpl,
	multipartService *storage.MultipartUploadService,
	imageProcessor *storage.ImageProcessor,
) *StorageAPI {
	return &StorageAPI{
		storageService:   storageService,
		multipartService: multipartService,
		imageProcessor:   imageProcessor,
	}
}

// ============ 基础文件操作 ============

// UploadFile 上传文件
//
//	@Summary		上传文件
//	@Description	上传单个文件（小文件 <50MB）
//	@Tags			文件存储
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file		formData	file	true	"文件"
//	@Param			category	formData	string	false	"分类"	default(attachment)
//	@Param			is_public	formData	boolean	false	"是否公开"	default(false)
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/files/upload [post]
func (api *StorageAPI) UploadFile(c *gin.Context) {
	// 1. 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未授权")
		return
	}

	// 2. 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		BadRequest(c, "参数错误", "文件上传失败: "+err.Error())
		return
	}
	defer file.Close()

	// 3. 获取其他参数
	category := c.DefaultPostForm("category", "attachment")
	isPublicStr := c.DefaultPostForm("is_public", "false")
	isPublic := isPublicStr == "true"

	// 4. 构建上传请求
	req := &storage.UploadRequest{
		File:        file,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		UserID:      userID.(string),
		IsPublic:    isPublic,
		Category:    category,
	}

	// 5. 上传文件
	fileInfo, err := api.storageService.Upload(c.Request.Context(), req)
	if err != nil {
		InternalError(c, "上传失败", err)
		return
	}

	Success(c, http.StatusOK, "上传成功", fileInfo)
}

// DownloadFile 下载文件
//
//	@Summary		下载文件
//	@Description	下载文件
//	@Tags			文件存储
//	@Produce		octet-stream
//	@Param			id	path	string	true	"文件ID"
//	@Success		200	{file}	binary	"文件内容"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/files/{id}/download [get]
func (api *StorageAPI) DownloadFile(c *gin.Context) {
	// 1. 获取文件ID
	fileID := c.Param("id")
	if fileID == "" {
		BadRequest(c, "参数错误", "文件ID不能为空")
		return
	}

	// 2. 获取用户ID（可选，用于权限检查）
	userID, _ := c.Get("userId")

	// 3. 获取文件信息
	fileInfo, err := api.storageService.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		NotFound(c, "文件不存在")
		return
	}

	// 4. 权限检查
	if !fileInfo.IsPublic && userID != nil {
		hasAccess, err := api.storageService.CheckAccess(c.Request.Context(), fileID, userID.(string))
		if err != nil || !hasAccess {
			Forbidden(c, "您没有访问该文件的权限")
			return
		}
	}

	// 5. 下载文件
	reader, err := api.storageService.Download(c.Request.Context(), fileID)
	if err != nil {
		InternalError(c, "下载失败", err)
		return
	}
	defer reader.Close()

	// 6. 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.OriginalName))
	c.Header("Content-Type", fileInfo.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size))

	// 7. 流式传输文件
	c.DataFromReader(http.StatusOK, fileInfo.Size, fileInfo.ContentType, reader, nil)
}

// GetFileInfo 获取文件信息
//
//	@Summary		获取文件信息
//	@Description	获取文件详细信息
//	@Tags			文件存储
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"文件ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/files/{id} [get]
func (api *StorageAPI) GetFileInfo(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		BadRequest(c, "参数错误", "文件ID不能为空")
		return
	}

	fileInfo, err := api.storageService.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		NotFound(c, "文件不存在")
		return
	}

	Success(c, http.StatusOK, "获取成功", fileInfo)
}

// DeleteFile 删除文件
//
//	@Summary		删除文件
//	@Description	删除文件
//	@Tags			文件存储
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"文件ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/files/{id} [delete]
func (api *StorageAPI) DeleteFile(c *gin.Context) {
	// 1. 获取文件ID和用户ID
	fileID := c.Param("id")
	userID, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未授权")
		return
	}

	// 2. 获取文件信息（验证权限）
	fileInfo, err := api.storageService.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		NotFound(c, "文件不存在")
		return
	}

	// 3. 权限检查（只有文件所有者可以删除）
	if fileInfo.UserID != userID.(string) {
		Forbidden(c, "只有文件所有者可以删除文件")
		return
	}

	// 4. 删除文件
	err = api.storageService.Delete(c.Request.Context(), fileID)
	if err != nil {
		InternalError(c, "删除失败", err)
		return
	}

	Success(c, http.StatusOK, "删除成功", nil)
}

// ListFiles 查询文件列表
//
//	@Summary		查询文件列表
//	@Description	查询用户的文件列表
//	@Tags			文件存储
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			category	query		string	false	"分类"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页大小"	default(20)
//	@Success		200			{object}	PaginatedResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/files [get]
func (api *StorageAPI) ListFiles(c *gin.Context) {
	// 1. 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未授权")
		return
	}

	// 2. 获取查询参数
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 3. 查询文件列表
	files, err := api.storageService.ListFiles(c.Request.Context(), &storage.ListFilesRequest{
		UserID:   userID.(string),
		Category: category,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		InternalError(c, "查询失败", err)
		return
	}

	// 4. 返回分页结果
	Paginated(c, files, int64(len(files)), page, pageSize, "查询成功")
}

// GetDownloadURL 获取下载链接
//
//	@Summary		获取下载链接
//	@Description	生成临时下载链接
//	@Tags			文件存储
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id			path		string	true	"文件ID"
//	@Param			expires_in	query		int		false	"过期时间(秒)"	default(3600)
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/files/{id}/url [get]
func (api *StorageAPI) GetDownloadURL(c *gin.Context) {
	fileID := c.Param("id")
	expiresIn, _ := strconv.Atoi(c.DefaultQuery("expires_in", "3600"))

	url, err := api.storageService.GetDownloadURL(
		c.Request.Context(),
		fileID,
		time.Duration(expiresIn)*time.Second,
	)
	if err != nil {
		InternalError(c, "生成链接失败", err)
		return
	}

	Success(c, http.StatusOK, "生成成功", map[string]interface{}{
		"url":        url,
		"expires_in": expiresIn,
	})
}

// ============ 分片上传 ============

// InitiateMultipartUpload 初始化分片上传
//
//	@Summary		初始化分片上传
//	@Description	初始化大文件分片上传
//	@Tags			文件存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		storage.InitiateMultipartUploadRequest	true	"初始化请求"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/files/multipart/init [post]
func (api *StorageAPI) InitiateMultipartUpload(c *gin.Context) {
	// 1. 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		Unauthorized(c, "未授权")
		return
	}

	// 2. 绑定请求参数
	var req storage.InitiateMultipartUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}
	req.UploadedBy = userID.(string)

	// 3. 初始化分片上传
	resp, err := api.multipartService.InitiateMultipartUpload(c.Request.Context(), &req)
	if err != nil {
		InternalError(c, "初始化失败", err)
		return
	}

	Success(c, http.StatusOK, "初始化成功", resp)
}

// UploadChunk 上传文件分片
//
//	@Summary		上传文件分片
//	@Description	上传单个文件分片
//	@Tags			文件存储
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			upload_id	formData	string	true	"上传ID"
//	@Param			chunk_index	formData	int		true	"分片索引"
//	@Param			chunk		formData	file	true	"分片文件"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/files/multipart/upload [post]
func (api *StorageAPI) UploadChunk(c *gin.Context) {
	// 1. 获取参数
	uploadID := c.PostForm("upload_id")
	chunkIndexStr := c.PostForm("chunk_index")
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		BadRequest(c, "参数错误", "分片索引无效")
		return
	}

	// 2. 获取分片文件
	file, header, err := c.Request.FormFile("chunk")
	if err != nil {
		BadRequest(c, "参数错误", "分片文件上传失败")
		return
	}
	defer file.Close()

	// 3. 上传分片
	req := &storage.UploadChunkRequest{
		UploadID:   uploadID,
		ChunkIndex: chunkIndex,
		ChunkData:  file,
		ChunkSize:  header.Size,
	}

	err = api.multipartService.UploadChunk(c.Request.Context(), req)
	if err != nil {
		InternalError(c, "上传分片失败", err)
		return
	}

	Success(c, http.StatusOK, "上传成功", nil)
}

// CompleteMultipartUpload 完成分片上传
//
//	@Summary		完成分片上传
//	@Description	完成所有分片上传，合并文件
//	@Tags			文件存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		storage.CompleteMultipartUploadRequest	true	"完成请求"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/files/multipart/complete [post]
func (api *StorageAPI) CompleteMultipartUpload(c *gin.Context) {
	var req storage.CompleteMultipartUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	fileMetadata, err := api.multipartService.CompleteMultipartUpload(c.Request.Context(), &req)
	if err != nil {
		InternalError(c, "完成上传失败", err)
		return
	}

	Success(c, http.StatusOK, "上传完成", fileMetadata)
}

// AbortMultipartUpload 中止分片上传
//
//	@Summary		中止分片上传
//	@Description	中止分片上传任务
//	@Tags			文件存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			upload_id	query		string	true	"上传ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/files/multipart/abort [post]
func (api *StorageAPI) AbortMultipartUpload(c *gin.Context) {
	uploadID := c.Query("upload_id")
	if uploadID == "" {
		BadRequest(c, "参数错误", "上传ID不能为空")
		return
	}

	err := api.multipartService.AbortMultipartUpload(c.Request.Context(), uploadID)
	if err != nil {
		InternalError(c, "中止上传失败", err)
		return
	}

	Success(c, http.StatusOK, "已中止上传", nil)
}

// GetUploadProgress 获取上传进度
//
//	@Summary		获取上传进度
//	@Description	获取分片上传进度
//	@Tags			文件存储
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			upload_id	query		string	true	"上传ID"
//	@Success		200			{object}	APIResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/files/multipart/progress [get]
func (api *StorageAPI) GetUploadProgress(c *gin.Context) {
	uploadID := c.Query("upload_id")
	if uploadID == "" {
		BadRequest(c, "参数错误", "上传ID不能为空")
		return
	}

	progress, err := api.multipartService.GetUploadProgress(c.Request.Context(), uploadID)
	if err != nil {
		InternalError(c, "获取进度失败", err)
		return
	}

	Success(c, http.StatusOK, "获取成功", map[string]interface{}{
		"upload_id": uploadID,
		"progress":  progress,
	})
}

// ============ 图片处理 ============

// GenerateThumbnail 生成缩略图
//
//	@Summary		生成缩略图
//	@Description	为图片生成缩略图
//	@Tags			文件存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file_id	query		string	true	"文件ID"
//	@Param			width	query		int		false	"宽度"	default(200)
//	@Param			height	query		int		false	"高度"	default(200)
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/files/thumbnail [post]
func (api *StorageAPI) GenerateThumbnail(c *gin.Context) {
	fileID := c.Query("file_id")
	width, _ := strconv.Atoi(c.DefaultQuery("width", "200"))
	height, _ := strconv.Atoi(c.DefaultQuery("height", "200"))

	// 获取文件信息
	fileInfo, err := api.storageService.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		NotFound(c, "文件不存在")
		return
	}

	// 生成缩略图
	thumbnailPath, err := api.imageProcessor.GenerateThumbnail(
		c.Request.Context(),
		fileInfo.Path,
		width,
		height,
		true, // 保持宽高比
	)
	if err != nil {
		InternalError(c, "生成缩略图失败", err)
		return
	}

	Success(c, http.StatusOK, "生成成功", map[string]interface{}{
		"thumbnail_path": thumbnailPath,
	})
}

// ============ 权限管理 ============

// GrantAccess 授予访问权限
//
//	@Summary		授予访问权限
//	@Description	授予用户文件访问权限
//	@Tags			文件存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file_id	path	string							true	"文件ID"
//	@Param			request	body	map[string]interface{}	true	"请求体"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/files/{file_id}/access [post]
func (api *StorageAPI) GrantAccess(c *gin.Context) {
	fileID := c.Param("file_id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationError(c, err)
		return
	}

	err := api.storageService.GrantAccess(c.Request.Context(), fileID, req.UserID)
	if err != nil {
		InternalError(c, "授权失败", err)
		return
	}

	Success(c, http.StatusOK, "授权成功", nil)
}

// RevokeAccess 撤销访问权限
//
//	@Summary		撤销访问权限
//	@Description	撤销用户文件访问权限
//	@Tags			文件存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file_id	path		string	true	"文件ID"
//	@Param			user_id	query		string	true	"用户ID"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/v1/files/{file_id}/access [delete]
func (api *StorageAPI) RevokeAccess(c *gin.Context) {
	fileID := c.Param("file_id")
	userIDToRevoke := c.Query("user_id")

	if userIDToRevoke == "" {
		BadRequest(c, "参数错误", "用户ID不能为空")
		return
	}

	err := api.storageService.RevokeAccess(c.Request.Context(), fileID, userIDToRevoke)
	if err != nil {
		InternalError(c, "撤销失败", err)
		return
	}

	Success(c, http.StatusOK, "撤销成功", nil)
}
