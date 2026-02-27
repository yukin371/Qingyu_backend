package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	adminService "Qingyu_backend/service/admin"
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
// @Success 200 {object} map[string]interface{} "code,message,data"
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("查询审计日志失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"logs":      logs,
			"total":     total,
			"page":      page,
			"size":      size,
			"totalPages": (total + int64(size) - 1) / int64(size),
		},
	})
}

// GetResourceAuditTrail 获取资源审计追踪
// @Summary 获取资源审计追踪
// @Description 根据资源类型和资源ID获取相关的审计日志
// @Tags Admin-Audit
// @Accept json
// @Produce json
// @Param type path string true "资源类型"
// @Param id path string true "资源ID"
// @Success 200 {object} map[string]interface{} "code,message,data"
// @Router /api/v1/admin/audit/trail/resource/{type}/{id} [get]
func (api *AuditAPI) GetResourceAuditTrail(c *gin.Context) {
	resourceType := c.Param("type")
	resourceID := c.Param("id")

	if resourceType == "" || resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "资源类型和资源ID不能为空",
		})
		return
	}

	// 查询资源相关日志
	logs, err := api.auditLogService.GetLogsByResource(context.Background(), resourceType, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("查询资源审计日志失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"logs":  logs,
			"total": len(logs),
		},
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("查询审计日志失败: %v", err),
		})
		return
	}

	// 导出
	if format == "csv" {
		csvData, err := api.exportAuditLogsToCSV(logs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": fmt.Sprintf("导出CSV失败: %v", err),
			})
			return
		}

		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=audit_trail_%s.csv", time.Now().Format("20060102_150405")))
		c.Data(http.StatusOK, "text/csv", csvData)
	} else if format == "xlsx" {
		// TODO: 实现Excel导出
		c.JSON(http.StatusNotImplemented, gin.H{
			"code":    501,
			"message": "Excel导出功能待实现",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的导出格式",
		})
	}
}

// GetAuditStatistics 获取审计统计
// @Summary 获取审计统计
// @Description 获取审计日志的统计信息
// @Tags Admin-Audit
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "code,message,data"
// @Router /api/v1/admin/audit/statistics [get]
func (api *AuditAPI) GetAuditStatistics(c *gin.Context) {
	// 获取总日志数
	req := &adminService.QueryAuditLogsRequest{
		Page:     1,
		PageSize: 1,
	}

	_, total, err := api.auditLogService.QueryAuditLogs(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取统计信息失败: %v", err),
		})
		return
	}

	// TODO: 添加更多统计信息
	// - 按操作类型统计
	// - 按管理员统计
	// - 按日期统计
	// - 按资源类型统计

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取统计信息成功",
		"data": gin.H{
			"total_logs": total,
			// 更多统计信息待添加
		},
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
