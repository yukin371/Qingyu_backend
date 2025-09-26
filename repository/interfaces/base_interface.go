package interfaces

import (
	"context"
	"time"
)

// 核心CRUD接口
type CRUDRepository[T any, ID comparable] interface {
	Create(ctx context.Context, entity T) error
	GetByID(ctx context.Context, id ID) (T, error)
	Update(ctx context.Context, id ID, updates map[string]interface{}) error
	Delete(ctx context.Context, id ID) error
}

// Filter 定义查询过滤器接口
type Filter interface {
	// 获取查询条件
	GetConditions() map[string]interface{}

	// 获取排序条件
	GetSort() map[string]int

	// 获取字段选择
	GetFields() []string

	// 验证过滤器
	Validate() error
}

// BaseFilter 基础过滤器实现
type BaseFilter struct {
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Sort       map[string]int         `json:"sort,omitempty"`
	Fields     []string               `json:"fields,omitempty"`
	Search     string                 `json:"search,omitempty"`
	CreatedAt  *TimeRange             `json:"createdAt,omitempty"`
	UpdatedAt  *TimeRange             `json:"updatedAt,omitempty"`
}

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"pageSize" validate:"min=1,max=100"`
	Skip     int `json:"skip,omitempty"`
}

// PagedResult 分页结果
type PagedResult[T any] struct {
	Data       []*T  `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// TimeRange 时间范围
type TimeRange struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// GetConditions 实现Filter接口
func (f *BaseFilter) GetConditions() map[string]interface{} {
	if f.Conditions == nil {
		f.Conditions = make(map[string]interface{})
	}

	// 添加时间范围条件
	if f.CreatedAt != nil {
		if f.CreatedAt.Start != nil {
			f.Conditions["createdAt.$gte"] = f.CreatedAt.Start
		}
		if f.CreatedAt.End != nil {
			f.Conditions["createdAt.$lte"] = f.CreatedAt.End
		}
	}

	if f.UpdatedAt != nil {
		if f.UpdatedAt.Start != nil {
			f.Conditions["updatedAt.$gte"] = f.UpdatedAt.Start
		}
		if f.UpdatedAt.End != nil {
			f.Conditions["updatedAt.$lte"] = f.UpdatedAt.End
		}
	}

	return f.Conditions
}

// GetSort 实现Filter接口
func (f *BaseFilter) GetSort() map[string]int {
	if f.Sort == nil {
		f.Sort = map[string]int{"createdAt": -1} // 默认按创建时间倒序
	}
	return f.Sort
}

// GetFields 实现Filter接口
func (f *BaseFilter) GetFields() []string {
	return f.Fields
}

// Validate 实现Filter接口
func (f *BaseFilter) Validate() error {
	// 基础验证逻辑
	if f.CreatedAt != nil {
		if f.CreatedAt.Start != nil && f.CreatedAt.End != nil {
			if f.CreatedAt.Start.After(*f.CreatedAt.End) {
				return NewValidationError("创建时间范围无效：开始时间不能晚于结束时间")
			}
		}
	}

	if f.UpdatedAt != nil {
		if f.UpdatedAt.Start != nil && f.UpdatedAt.End != nil {
			if f.UpdatedAt.Start.After(*f.UpdatedAt.End) {
				return NewValidationError("更新时间范围无效：开始时间不能晚于结束时间")
			}
		}
	}

	return nil
}

// CalculatePagination 计算分页参数
func (p *Pagination) CalculatePagination() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	p.Skip = (p.Page - 1) * p.PageSize
}

// NewPagedResult 创建分页结果
func NewPagedResult[T any](data []*T, total int64, pagination Pagination) *PagedResult[T] {
	totalPages := int((total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))

	return &PagedResult[T]{
		Data:       data,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError 创建验证错误
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// NewFieldValidationError 创建字段验证错误
func NewFieldValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// 查询接口
type QueryRepository[T any] interface {
	List(ctx context.Context, filter Filter) ([]T, error)
	FindWithPagination(ctx context.Context, filter Filter, pagination Pagination) (*PagedResult[T], error)
}

// 批量操作接口
type BatchRepository[T any, ID comparable] interface {
	BatchCreate(ctx context.Context, entities []T) error
	BatchUpdate(ctx context.Context, ids []ID, updates map[string]interface{}) error
	BatchDelete(ctx context.Context, ids []ID) error
}

// 健康检查接口
type HealthRepository interface {
	Health(ctx context.Context) error
}
