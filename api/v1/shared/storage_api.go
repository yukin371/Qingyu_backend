package shared

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/service/shared/storage"
)

// StorageAPI 存储服务API处理器
type StorageAPI struct {
	storageService storage.StorageService
}

// NewStorageAPI 创建存储API实例
func NewStorageAPI(storageService storage.StorageService) *StorageAPI {
	return &StorageAPI{
		storageService: storageService,
	}
}

// UploadFile 上传文件
//
//	@Summary		上传文件
//	@Description	上传文件到存储服务
//	@Tags			存储
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file	formData	file	true	"上传文件"
//	@Param			path	formData	string	false	"存储路径"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/storage/upload [post]
func (api *StorageAPI) UploadFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "获取文件失败: " + err.Error(),
		})
		return
	}

	// 打开文件
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "打开文件失败: " + err.Error(),
		})
		return
	}
	defer fileReader.Close()

	// 获取存储路径
	path := c.PostForm("path")
	if path == "" {
		path = "uploads/" + userID.(string)
	}

	req := &storage.UploadRequest{
		File:        fileReader,
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
		UserID:      userID.(string),
		Category:    path,
	}

	fileInfo, err := api.storageService.Upload(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "上传文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "上传文件成功",
		Data:    fileInfo,
	})
}

// DownloadFile 下载文件
//
//	@Summary		下载文件
//	@Description	下载指定的文件
//	@Tags			存储
//	@Accept			json
//	@Produce		application/octet-stream
//	@Security		ApiKeyAuth
//	@Param			file_id	path		string	true	"文件ID"
//	@Success		200		{file}		binary
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		404		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/storage/download/{file_id} [get]
func (api *StorageAPI) DownloadFile(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	fileID := c.Param("file_id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "缺少文件ID",
		})
		return
	}

	file, err := api.storageService.Download(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "下载文件失败: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+fileID)
	c.Header("Content-Type", "application/octet-stream")

	// 将文件内容复制到响应
	c.DataFromReader(http.StatusOK, -1, "application/octet-stream", file, map[string]string{
		"Content-Disposition": "attachment; filename=" + fileID,
	})
}

// DeleteFile 删除文件
//
//	@Summary		删除文件
//	@Description	删除指定的文件
//	@Tags			存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file_id	path		string	true	"文件ID"
//	@Success		200		{object}	APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		403		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/storage/files/{file_id} [delete]
func (api *StorageAPI) DeleteFile(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	fileID := c.Param("file_id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "缺少文件ID",
		})
		return
	}

	err := api.storageService.Delete(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "删除文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "删除文件成功",
	})
}

// GetFileInfo 获取文件信息
//
//	@Summary		获取文件信息
//	@Description	获取指定文件的详细信息
//	@Tags			存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file_id	path		string	true	"文件ID"
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		404		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/storage/files/{file_id} [get]
func (api *StorageAPI) GetFileInfo(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	fileID := c.Param("file_id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "缺少文件ID",
		})
		return
	}

	fileInfo, err := api.storageService.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取文件信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取文件信息成功",
		Data:    fileInfo,
	})
}

// ListFiles 列出文件
//
//	@Summary		列出文件
//	@Description	列出用户的文件列表
//	@Tags			存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			page_size	query		int		false	"每页数量"	default(20)
//	@Param			path		query		string	false	"路径前缀"
//	@Success 200 {object} APIResponse
//	@Failure		401			{object}	APIResponse
//	@Failure		500			{object}	APIResponse
//	@Router			/api/v1/shared/storage/files [get]
func (api *StorageAPI) ListFiles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	req := &storage.ListFilesRequest{
		UserID:   userID.(string),
		Category: category,
		Page:     page,
		PageSize: pageSize,
	}

	files, err := api.storageService.ListFiles(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "列出文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponseHelper(
		files,
		int64(len(files)),
		page,
		pageSize,
		"列出文件成功",
	))
}

// GetFileURL 获取文件访问URL
//
//	@Summary		获取文件访问URL
//	@Description	获取文件的临时访问URL
//	@Tags			存储
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file_id	path		string	true	"文件ID"
//	@Param			expire	query		int		false	"过期时间(秒)"	default(3600)
//	@Success 200 {object} APIResponse
//	@Failure		400		{object}	APIResponse
//	@Failure		401		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/shared/storage/files/{file_id}/url [get]
func (api *StorageAPI) GetFileURL(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Code:    401,
			Message: "未认证",
		})
		return
	}

	fileID := c.Param("file_id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "缺少文件ID",
		})
		return
	}

	expire, _ := strconv.Atoi(c.DefaultQuery("expire", "3600"))

	url, err := api.storageService.GetDownloadURL(c.Request.Context(), fileID, time.Duration(expire)*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "获取文件URL失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "获取文件URL成功",
		Data:    map[string]string{"url": url},
	})
}
