package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	response "Qingyu_backend/pkg/response"
	adminService "Qingyu_backend/service/admin"
	"github.com/gin-gonic/gin"
)

// AuditAPI 审计追踪API
type AuditAPI struct {
	auditLogService adminService.AuditLogService
	exportService   adminService.ExportService
}

// NewAuditAPI 创建审计API实例
func NewAuditAPI(auditLogService adminService.AuditLogService, exportService adminService.ExportService) *AuditAPI {
	return &AuditAPI{
		auditLogService: auditLogService,
		exportService:   exportService,
	}
}

// GetAuditTrail 获取审计追踪
// @Summary 获取审计追踪
// @Description 获取审计日志列表，支持多种过滤条件
// @Tags Admin-Audit
// @Accept json
// @Produce json
// @Param admin_id query string false "管理员ID"
// @Param operation query string false "操作类型"
// @Param resource_type query string false "资源类型"
// @Param resource_id query string false "资源ID"
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} response.APIResponse "code,message,data"
// @Router /api/v1/admin/audit/trail [get]
func (api *AuditAPI) GetAuditTrail(c *gin.Context) {
	// 解析查询参数
	adminID := c.Query("admin_id")
	operation := c.Query("operation")
	resourceType := c.Query("resource_type")
	resourceID := c.Query("resource_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 解析日期
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	// 构建查询请求
	req := &adminService.QueryAuditLogsRequest{
		AdminID:      adminID,
		Operation:    operation,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		StartDate:    startDate,
		EndDate:      endDate,
		Page:         page,
		PageSize:     size,
	}

	// 查询审计日志
	logs, total, err := api.auditLogService.QueryAuditLogs(context.Background(), req)
	if err != nil {
		response.InternalError(c, fmt.Errorf("查询审计日志失败: %w", err))
		return
	}

	response.Paginated(c, logs, total, page, size, "查询成功")
}

// GetResourceAuditTrail 获取资源审计追踪
// @Summary 获取资源审计追踪
// @Description 根据资源类型和资源ID获取相关的审计日志
// @Tags Admin-Audit
// @Accept json
// @Produce json
// @Param type path string true "资源类型"
// @Param id path string true "资源ID"
// @Success 200 {object} response.APIResponse "code,message,data"
// @Router /api/v1/admin/audit/trail/resource/{type}/{id} [get]
func (api *AuditAPI) GetResourceAuditTrail(c *gin.Context) {
	resourceType := c.Param("type")
	resourceID := c.Param("id")

	if resourceType == "" || resourceID == "" {
		response.BadRequest(c, "资源类型和资源ID不能为空", nil)
		return
	}

	// 查询资源相关日志
	logs, err := api.auditLogService.GetLogsByResource(context.Background(), resourceType, resourceID)
	if err != nil {
		response.InternalError(c, fmt.Errorf("查询资源审计日志失败: %w", err))
		return
	}

	response.SuccessWithMessage(c, "查询成功", gin.H{
		"logs":  logs,
		"total": len(logs),
	})
}

// ExportAuditTrail 导出审计追踪
// @Summary 导出审计追踪
// @Description 导出审计日志为CSV或Excel格式
// @Tags Admin-Audit
// @Accept json
// @Produce csv/xlsx
// @Param format query string false "导出格式 (csv/xlsx)" default(csv)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param admin_id query string false "管理员ID"
// @Param operation query string false "操作类型"
// @Success 200 {file} file
// @Router /api/v1/admin/audit/trail/export [get]
func (api *AuditAPI) ExportAuditTrail(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	adminID := c.Query("admin_id")
	operation := c.Query("operation")

	// 解析日期
	var startDate, endDate *time.Time
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	// 查询审计日志
	req := &adminService.QueryAuditLogsRequest{
		AdminID:   adminID,
		Operation: operation,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      1,
		PageSize:  10000, // 导出时设置较大的页面大小
	}

	logs, _, err := api.auditLogService.QueryAuditLogs(context.Background(), req)
	if err != nil {
		response.InternalError(c, fmt.Errorf("查询审计日志失败: %w", err))
		return
	}

	// 导出
	if format == "csv" {
		csvData, err := api.exportAuditLogsToCSV(logs)
		if err != nil {
			response.InternalError(c, fmt.Errorf("导出CSV失败: %w", err))
			return
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=audit_trail_%s.csv", time.Now().Format("20060102_150405")))
		c.Data(http.StatusOK, "text/csv", csvData)
	} else if format == "xlsx" {
		// TODO: 实现Excel导出
		response.JSON(c, http.StatusNotImplemented, response.APIResponse{
			Code:      response.CodeInternalError,
			Message:   "Excel导出功能待实现",
			Timestamp: time.Now().UnixMilli(),
		})
	} else {
		response.BadRequest(c, "不支持的导出格式", nil)
	}
}

// GetAuditStatistics 获取审计统计
// @Summary 获取审计统计
// @Description 获取审计日志的统计信息
// @Tags Admin-Audit
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse "code,message,data"
// @Router /api/v1/admin/audit/statistics [get]
func (api *AuditAPI) GetAuditStatistics(c *gin.Context) {
	// 获取总日志数
	req := &adminService.QueryAuditLogsRequest{
		Page:     1,
		PageSize: 1,
	}

	_, total, err := api.auditLogService.QueryAuditLogs(context.Background(), req)
	if err != nil {
		response.InternalError(c, fmt.Errorf("获取统计信息失败: %w", err))
		return
	}

	// TODO: 添加更多统计信息
	// - 按操作类型统计
	// - 按管理员统计
	// - 按日期统计
	// - 按资源类型统计

	response.SuccessWithMessage(c, "获取统计信息成功", gin.H{
		"total_logs": total,
	})
}

// exportAuditLogsToCSV 将审计日志导出为CSV
func (api *AuditAPI) exportAuditLogsToCSV(logs interface{}) ([]byte, error) {
	// 简化实现：直接返回CSV格式的字节数组
	// 实际实现中应该使用专门的CSV导出库或复用ExportService

	csvContent := "ID,管理员ID,管理员名称,操作,资源类型,资源ID,IP,创建时间\n"

	// TODO: 实现完整的CSV导出逻辑
	// 这里简化处理，实际应该遍历logs并格式化每条记录

	return []byte(csvContent), nil
}
