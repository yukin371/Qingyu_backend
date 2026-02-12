package impl

import (
	"context"

	"Qingyu_backend/models/writer"
	serviceWriter "Qingyu_backend/service/interfaces/writer"
	writerproject "Qingyu_backend/service/writer/project"
)

// ProjectManagementImpl 项目管理端口实现
type ProjectManagementImpl struct {
	projectService *writerproject.ProjectService
	serviceName    string
	version        string
}

// NewProjectManagementImpl 创建项目管理端口实现
func NewProjectManagementImpl(projectService *writerproject.ProjectService) serviceWriter.ProjectManagementPort {
	return &ProjectManagementImpl{
		projectService: projectService,
		serviceName:    "ProjectManagementPort",
		version:        "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (p *ProjectManagementImpl) Initialize(ctx context.Context) error {
	return p.projectService.Initialize(ctx)
}

func (p *ProjectManagementImpl) Health(ctx context.Context) error {
	return p.projectService.Health(ctx)
}

func (p *ProjectManagementImpl) Close(ctx context.Context) error {
	return p.projectService.Close(ctx)
}

func (p *ProjectManagementImpl) GetServiceName() string {
	return p.serviceName
}

func (p *ProjectManagementImpl) GetVersion() string {
	return p.version
}

// ============================================================================
// ProjectManagementPort 方法实现
// ============================================================================

// CreateProject 创建新项目
func (p *ProjectManagementImpl) CreateProject(ctx context.Context, req *serviceWriter.CreateProjectRequest) (*serviceWriter.CreateProjectResponse, error) {
	// 转换请求类型
	projectReq := &writerproject.CreateProjectRequest{
		Title:    req.Title,
		Summary:  req.Summary,
		CoverURL: req.CoverURL,
		Category: req.Category,
		Tags:     req.Tags,
	}
	projectResp, err := p.projectService.CreateProject(ctx, projectReq)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.CreateProjectResponse{
		ProjectID: projectResp.ProjectID,
		Title:     projectResp.Title,
		Status:    projectResp.Status,
		CreatedAt: projectResp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetProject 获取项目详情
func (p *ProjectManagementImpl) GetProject(ctx context.Context, projectID string) (*writer.Project, error) {
	return p.projectService.GetProject(ctx, projectID)
}

// GetProjectByID 根据ID获取项目（别名方法，兼容API层调用）
func (p *ProjectManagementImpl) GetProjectByID(ctx context.Context, projectID string) (*writer.Project, error) {
	return p.projectService.GetProjectByID(ctx, projectID)
}

// GetByIDWithoutAuth 获取项目详情（无权限检查，用于内部服务调用如AI）
func (p *ProjectManagementImpl) GetByIDWithoutAuth(ctx context.Context, projectID string) (*writer.Project, error) {
	return p.projectService.GetByIDWithoutAuth(ctx, projectID)
}

// ListMyProjects 获取我的项目列表
func (p *ProjectManagementImpl) ListMyProjects(ctx context.Context, req *serviceWriter.ListProjectsRequest) (*serviceWriter.ListProjectsResponse, error) {
	// 转换请求类型
	projectReq := &writerproject.ListProjectsRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
	}
	projectResp, err := p.projectService.ListMyProjects(ctx, projectReq)
	if err != nil {
		return nil, err
	}
	// 转换响应类型
	return &serviceWriter.ListProjectsResponse{
		Projects: projectResp.Projects,
		Total:    projectResp.Total,
		Page:     projectResp.Page,
		PageSize: projectResp.PageSize,
	}, nil
}

// GetProjectList 获取项目列表（兼容API层调用）
func (p *ProjectManagementImpl) GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*writer.Project, error) {
	return p.projectService.GetProjectList(ctx, userID, status, limit, offset)
}

// UpdateProject 更新项目
func (p *ProjectManagementImpl) UpdateProject(ctx context.Context, projectID string, req *serviceWriter.UpdateProjectRequest) error {
	// 转换请求类型
	projectReq := &writerproject.UpdateProjectRequest{
		Title:    req.Title,
		Summary:  req.Summary,
		CoverURL: req.CoverURL,
		Category: req.Category,
		Tags:     req.Tags,
		Status:   req.Status,
	}
	return p.projectService.UpdateProject(ctx, projectID, projectReq)
}

// UpdateProjectByID 更新项目（兼容API层调用）
func (p *ProjectManagementImpl) UpdateProjectByID(ctx context.Context, projectID, userID string, req *writer.Project) error {
	return p.projectService.UpdateProjectByID(ctx, projectID, userID, req)
}

// DeleteProject 删除项目
func (p *ProjectManagementImpl) DeleteProject(ctx context.Context, projectID string) error {
	return p.projectService.DeleteProject(ctx, projectID)
}

// DeleteProjectByID 删除项目（兼容API层调用）
func (p *ProjectManagementImpl) DeleteProjectByID(ctx context.Context, projectID, userID string) error {
	return p.projectService.DeleteProjectByID(ctx, projectID, userID)
}

// RestoreProjectByID 恢复已删除的项目
func (p *ProjectManagementImpl) RestoreProjectByID(ctx context.Context, projectID, userID string) error {
	return p.projectService.RestoreProjectByID(ctx, projectID, userID)
}

// DeleteHard 物理删除项目（永久删除）
func (p *ProjectManagementImpl) DeleteHard(ctx context.Context, projectID string) error {
	return p.projectService.DeleteHard(ctx, projectID)
}

// UpdateProjectStatistics 更新项目统计
func (p *ProjectManagementImpl) UpdateProjectStatistics(ctx context.Context, projectID string, stats *writer.ProjectStats) error {
	return p.projectService.UpdateProjectStatistics(ctx, projectID, stats)
}

// RecalculateProjectStatistics 重新计算项目统计信息
func (p *ProjectManagementImpl) RecalculateProjectStatistics(ctx context.Context, projectID string) error {
	return p.projectService.RecalculateProjectStatistics(ctx, projectID)
}
