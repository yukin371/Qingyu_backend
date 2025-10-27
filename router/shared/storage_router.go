package shared

import (
	sharedApi "Qingyu_backend/api/v1/shared"
	"Qingyu_backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitStorageRoutes 初始化存储路由
func InitStorageRoutes(router *gin.RouterGroup, api *sharedApi.StorageAPI) {
	storage := router.Group("/files")
	{
		// ============ 公开路由（可选认证） ============

		// 下载文件（公开文件无需认证）
		storage.GET("/:id/download", api.DownloadFile)

		// ============ 需要认证的路由 ============
		authenticated := storage.Group("")
		authenticated.Use(middleware.JWTAuth())
		{
			// 基础文件操作
			authenticated.POST("/upload", api.UploadFile)     // 上传文件
			authenticated.GET("/:id", api.GetFileInfo)        // 获取文件信息
			authenticated.DELETE("/:id", api.DeleteFile)      // 删除文件
			authenticated.GET("", api.ListFiles)              // 查询文件列表
			authenticated.GET("/:id/url", api.GetDownloadURL) // 获取下载链接

			// 分片上传
			authenticated.POST("/multipart/init", api.InitiateMultipartUpload)     // 初始化分片上传
			authenticated.POST("/multipart/upload", api.UploadChunk)               // 上传分片
			authenticated.POST("/multipart/complete", api.CompleteMultipartUpload) // 完成分片上传
			authenticated.POST("/multipart/abort", api.AbortMultipartUpload)       // 中止分片上传
			authenticated.GET("/multipart/progress", api.GetUploadProgress)        // 获取上传进度

			// 图片处理
			authenticated.POST("/thumbnail", api.GenerateThumbnail) // 生成缩略图

			// 权限管理
			authenticated.POST("/:file_id/access", api.GrantAccess)    // 授予访问权限
			authenticated.DELETE("/:file_id/access", api.RevokeAccess) // 撤销访问权限
		}
	}
}
