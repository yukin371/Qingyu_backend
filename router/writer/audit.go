package writer

import (
	"Qingyu_backend/service/interfaces/audit"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/writer"
)

// InitAuditRouter 初始化审核路由
func InitAuditRouter(r *gin.RouterGroup, auditService audit.ContentAuditService) {
	auditApi := writer.NewAuditApi(auditService)

	// 公开的审核接口（需要认证）
	auditGroup := r.Group("/audit")
	{
		// 实时检测
		auditGroup.POST("/check", auditApi.CheckContent)

		// 申诉
		auditGroup.POST("/:id/appeal", auditApi.SubmitAppeal)
	}

	// 文档审核接口
	documentGroup := r.Group("/documents")
	{
		// 全文审核
		documentGroup.POST("/:id/audit", auditApi.AuditDocument)

		// 获取审核结果
		documentGroup.GET("/:id/audit-result", auditApi.GetAuditResult)
	}

	// 用户违规查询
	userGroup := r.Group("/users")
	{
		userGroup.GET("/:userId/violations", auditApi.GetUserViolations)
		userGroup.GET("/:userId/violation-summary", auditApi.GetUserViolationSummary)
	}

	// 管理员审核接口（需要管理员权限）
	adminGroup := r.Group("/admin/audit")
	// adminGroup.Use(middleware.AdminPermission()) // TODO: 添加管理员权限中间件
	{
		// 待复核列表
		adminGroup.GET("/pending", auditApi.GetPendingReviews)

		// 高风险记录
		adminGroup.GET("/high-risk", auditApi.GetHighRiskAudits)

		// 复核审核结果
		adminGroup.POST("/:id/review", auditApi.ReviewAudit)

		// 复核申诉
		adminGroup.POST("/:id/appeal/review", auditApi.ReviewAppeal)
	}
}
