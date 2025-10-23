package project

import (
	"Qingyu_backend/models/writer"
	"time"
)

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Title    string   `json:"title" validate:"required,min=1,max=100"`
	Summary  string   `json:"summary,omitempty"`
	CoverURL string   `json:"coverUrl,omitempty"`
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// CreateProjectResponse 创建项目响应
type CreateProjectResponse struct {
	ProjectID string    `json:"projectId"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Title    string   `json:"title,omitempty"`
	Summary  string   `json:"summary,omitempty"`
	CoverURL string   `json:"coverUrl,omitempty"`
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Status   string   `json:"status,omitempty"`
}

// ListProjectsRequest 项目列表请求
type ListProjectsRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Status   string `json:"status,omitempty"`
	Category string `json:"category,omitempty"`
}

// ListProjectsResponse 项目列表响应
type ListProjectsResponse struct {
	Projects []*writer.Project `json:"projects"`
	Total    int64             `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

// GetProjectResponse 获取项目详情响应
type GetProjectResponse struct {
	*writer.Project
}

// UpdateStatisticsRequest 更新统计请求
type UpdateStatisticsRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
}

// AddCollaboratorRequest 添加协作者请求
type AddCollaboratorRequest struct {
	ProjectID string                    `json:"projectId" validate:"required"`
	UserID string                  `json:"userId" validate:"required"`
	Role   writer.CollaboratorRole `json:"role" validate:"required"`
}

// RemoveCollaboratorRequest 移除协作者请求
type RemoveCollaboratorRequest struct {
	ProjectID string `json:"projectId" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
}
