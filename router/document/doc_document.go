package document

import (
	api "Qingyu_backend/api/v1/document"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文档相关路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 文档路由组
	docRouter := r.Group("/document")
	{
		// 文档相关路由
		d := api.NewDocumentApi()
		// 创建文档
		docRouter.POST("/", d.Create)

		// 获取文档列表
		docRouter.GET("/", d.List)

		// 获取单个文档
		docRouter.GET("/:id", d.Get)

		// 更新文档
		docRouter.PUT("/:id", d.Update)

		// 删除文档
		docRouter.DELETE("/:id", d.Delete)

		// 版本相关路由
		v := api.NewVersionApi()

		docRouter.POST("/version/create", v.CreateVersion)     // 创建版本
		docRouter.POST("/version/rollback", v.Rollback)        // 回滚版本
		docRouter.GET("/version/list", v.ListVersions)         // 获取版本列表
		docRouter.POST("/version/patch/create", v.CreatePatch) // 创建补丁
		docRouter.POST("/version/patch/apply", v.ApplyPatch)   // 应用补丁

		// 批量提交相关路由
		docRouter.POST("/version/commit", v.CreateCommit)              // 创建提交
		docRouter.GET("/version/commits", v.ListCommits)               // 获取提交列表
		docRouter.GET("/version/commit/:commitId", v.GetCommitDetails) // 获取提交详情
		docRouter.POST("/version/conflicts/detect", v.DetectConflicts) // 检测冲突

		// 单文件版本管理API
		docRouter.POST("/version/create", v.CreateVersion)  // 创建版本
		docRouter.PUT("/version/update", v.UpdateVersion)   // 更新版本
		docRouter.GET("/version/revisions", v.GetRevisions) // 获取版本修订列表

		// 冲突处理API
		docRouter.POST("/version/conflicts/resolve", v.ResolveBatchConflicts)     // 批量处理冲突
		docRouter.POST("/version/conflicts/auto-resolve", v.AutoResolveConflicts) // 自动处理冲突
	}
}
