package admin

import (
	"context"
	"time"

	"Qingyu_backend/models/admin"
	base "Qingyu_backend/repository/interfaces/infrastructure"
)

// ExportHistoryRepository 导出历史仓储接口
type ExportHistoryRepository interface {
	// 继承基础CRUD接口
	base.CRUDRepository[*admin.ExportHistory, string]
	// 继承 HealthRepository 接口
	base.HealthRepository

	// ListByUser 按管理员ID查询导出历史
	ListByUser(ctx context.Context, adminID string, page, pageSize int) ([]*admin.ExportHistory, int64, error)

	// ListByDateRange 按日期范围查询导出历史
	ListByDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*admin.ExportHistory, int64, error)

	// CleanOldRecords 清理指定日期之前的旧记录
	CleanOldRecords(ctx context.Context, beforeDate time.Time) (int64, error)
}

// ExportHistoryFilter 导出历史过滤器
type ExportHistoryFilter struct {
	AdminID    string     `json:"admin_id,omitempty"`    // 管理员ID
	ExportType string     `json:"export_type,omitempty"` // 导出类型
	Format     string     `json:"format,omitempty"`      // 导出格式
	Status     string     `json:"status,omitempty"`      // 状态
	StartDate  *time.Time `json:"start_date,omitempty"`  // 开始日期
	EndDate    *time.Time `json:"end_date,omitempty"`    // 结束日期
	Sort       map[string]int `json:"sort,omitempty"`     // 排序
	Fields     []string   `json:"fields,omitempty"`      // 字段选择
}

// GetConditions 实现Filter接口
func (f *ExportHistoryFilter) GetConditions() map[string]interface{} {
	conditions := make(map[string]interface{})

	// 添加基础条件
	if f.AdminID != "" {
		conditions["admin_id"] = f.AdminID
	}
	if f.ExportType != "" {
		conditions["export_type"] = f.ExportType
	}
	if f.Format != "" {
		conditions["format"] = f.Format
	}
	if f.Status != "" {
		conditions["status"] = f.Status
	}

	// 添加日期范围条件
	if f.StartDate != nil {
		conditions["created_at.$gte"] = f.StartDate
	}
	if f.EndDate != nil {
		conditions["created_at.$lte"] = f.EndDate
	}

	return conditions
}

// GetSort 实现Filter接口
func (f *ExportHistoryFilter) GetSort() map[string]int {
	if f.Sort != nil && len(f.Sort) > 0 {
		return f.Sort
	}
	// 默认按创建时间倒序
	return map[string]int{"created_at": -1}
}

// GetFields 实现Filter接口
func (f *ExportHistoryFilter) GetFields() []string {
	return f.Fields
}

// Validate 实现Filter接口
func (f *ExportHistoryFilter) Validate() error {
	// 验证日期范围
	if f.StartDate != nil && f.EndDate != nil {
		if f.StartDate.After(*f.EndDate) {
			return base.NewValidationError("开始日期不能晚于结束日期")
		}
	}

	// 验证导出类型
	if f.ExportType != "" {
		validTypes := map[string]bool{
			admin.ExportTypeBooks:    true,
			admin.ExportTypeChapters: true,
			admin.ExportTypeUsers:    true,
		}
		if !validTypes[f.ExportType] {
			return base.NewFieldValidationError("export_type", "无效的导出类型")
		}
	}

	// 验证格式
	if f.Format != "" {
		validFormats := map[string]bool{
			admin.ExportFormatCSV:   true,
			admin.ExportFormatExcel: true,
		}
		if !validFormats[f.Format] {
			return base.NewFieldValidationError("format", "无效的导出格式")
		}
	}

	// 验证状态
	if f.Status != "" {
		validStatuses := map[string]bool{
			admin.ExportStatusPending:   true,
			admin.ExportStatusCompleted: true,
			admin.ExportStatusFailed:    true,
		}
		if !validStatuses[f.Status] {
			return base.NewFieldValidationError("status", "无效的导出状态")
		}
	}

	return nil
}
