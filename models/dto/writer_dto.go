package dto

import "time"

// ===========================
// Writer DTO（符合分层架构规范）
// ===========================
//
// 本文件包含 Writer 模块的数据传输对象（DTO）
//
// 命名和标签规范：
// - DTO 结构体使用驼峰命名（PascalCase）
// - JSON 字段标签使用驼峰命名（camelCase）
// - 对应的 MongoDB 模型（位于 models/writer/）使用蛇形命名（snake_case）的 BSON 标签
//
// 用途：
// - 用于 Service 层和 API 层之间的数据传输
// - ID 和时间字段统一使用字符串类型
// - 避免直接暴露 MongoDB 模型到 API 层

// ===========================
// Project DTOs
// ===========================

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=100"`
	Summary  string   `json:"summary,omitempty" validate:"max=500"`
	CoverURL string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
	Tags     []string `json:"tags,omitempty" validate:"max=10,dive,min=1,max=50"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Title    *string   `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Summary  *string   `json:"summary,omitempty" validate:"omitempty,max=500"`
	CoverURL *string   `json:"coverUrl,omitempty" validate:"omitempty,url,max=500"`
	Tags     *[]string `json:"tags,omitempty" validate:"omitempty,max=10,dive,min=1,max=50"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	CoverURL  string    `json:"coverUrl"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ListProjectsRequest 查询参数用于列出项目
type ListProjectsRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Status   string `form:"status" validate:"omitempty,oneof=draft published archived"`
	Sort     string `form:"sort" validate:"omitempty,oneof=created_at updated_at title"`
	Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
}

// ProjectListResponse 分页项目列表响应
type ProjectListResponse struct {
	Items    []ProjectResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}
