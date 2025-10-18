package project

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/models/document"
	pkgErrors "Qingyu_backend/pkg/errors"
	writingRepo "Qingyu_backend/repository/interfaces/writing"
	"Qingyu_backend/service/base"
)

// ProjectService 项目服务
type ProjectService struct {
	projectRepo writingRepo.ProjectRepository
	eventBus    base.EventBus
	serviceName string
	version     string
}

// NewProjectService 创建项目服务
func NewProjectService(
	projectRepo writingRepo.ProjectRepository,
	eventBus base.EventBus,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		eventBus:    eventBus,
		serviceName: "ProjectService",
		version:     "1.0.0",
	}
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {
	// 1. 参数验证
	if err := s.validateCreateProjectRequest(req); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "参数验证失败", err.Error(), err)
	}

	// 2. 获取用户ID
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	// 3. 创建项目对象
	project := &document.Project{
		AuthorID:   userID,
		Title:      req.Title,
		Summary:    req.Summary,
		CoverURL:   req.CoverURL,
		Category:   req.Category,
		Tags:       req.Tags,
		Status:     document.StatusDraft,       // 默认状态为草稿
		Visibility: document.VisibilityPrivate, // 默认为私密
	}

	// 4. 保存到数据库
	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "创建项目失败", "", err)
	}

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "project.created",
			EventData: map[string]interface{}{
				"project_id": project.ID,
				"author_id":  project.AuthorID,
				"title":      project.Title,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	// 6. 返回响应
	return &CreateProjectResponse{
		ProjectID: project.ID,
		Title:     project.Title,
		Status:    string(project.Status),
		CreatedAt: project.CreatedAt,
	}, nil
}

// GetProject 获取项目详情
func (s *ProjectService) GetProject(ctx context.Context, projectID string) (*document.Project, error) {
	// 1. 查询项目
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	// 2. 权限检查
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanView(userID) {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限查看该项目", "", nil)
	}

	return project, nil
}

// ListMyProjects 获取我的项目列表
func (s *ProjectService) ListMyProjects(ctx context.Context, req *ListProjectsRequest) (*ListProjectsResponse, error) {
	// 1. 获取用户ID
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	// 2. 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 3. 计算偏移量
	offset := int64((req.Page - 1) * req.PageSize)
	limit := int64(req.PageSize)

	// 4. 查询项目列表
	var projects []*document.Project
	var err error

	if req.Status != "" {
		// 按状态查询
		projects, err = s.projectRepo.GetByOwnerAndStatus(ctx, userID, req.Status, limit, offset)
	} else {
		// 查询所有项目
		projects, err = s.projectRepo.GetListByOwnerID(ctx, userID, limit, offset)
	}

	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目列表失败", "", err)
	}

	// 5. 统计总数
	total, err := s.projectRepo.CountByOwner(ctx, userID)
	if err != nil {
		return nil, pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "统计项目数失败", "", err)
	}

	// 6. 返回响应
	return &ListProjectsResponse{
		Projects: projects,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(ctx context.Context, projectID string, req *UpdateProjectRequest) error {
	// 1. 查询项目
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	// 2. 权限检查
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.CanEdit(userID) {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "无权限编辑该项目", "", nil)
	}

	// 3. 构建更新数据
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.CoverURL != "" {
		updates["cover_url"] = req.CoverURL
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Status != "" {
		// 验证状态值
		status := document.ProjectStatus(req.Status)
		if !status.IsValid() {
			return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorValidation, "无效的项目状态", "", nil)
		}
		updates["status"] = req.Status
	}

	// 4. 更新项目
	if err := s.projectRepo.Update(ctx, projectID, updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新项目失败", "", err)
	}

	// 5. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "project.updated",
			EventData: map[string]interface{}{
				"project_id": projectID,
				"updates":    updates,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// DeleteProject 删除项目
func (s *ProjectService) DeleteProject(ctx context.Context, projectID string) error {
	// 1. 查询项目
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	// 2. 权限检查（只有所有者可以删除）
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorUnauthorized, "用户未登录", "", nil)
	}

	if !project.IsOwner(userID) {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "只有项目所有者可以删除项目", "", nil)
	}

	// 3. 软删除项目
	if err := s.projectRepo.SoftDelete(ctx, projectID, userID); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "删除项目失败", "", err)
	}

	// 4. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "project.deleted",
			EventData: map[string]interface{}{
				"project_id": projectID,
				"author_id":  project.AuthorID,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// UpdateProjectStatistics 更新项目统计
func (s *ProjectService) UpdateProjectStatistics(ctx context.Context, projectID string, stats *document.ProjectStats) error {
	// 1. 验证项目存在
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	// 2. 更新统计
	updates := map[string]interface{}{
		"statistics": stats,
	}

	if err := s.projectRepo.Update(ctx, projectID, updates); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "更新统计失败", "", err)
	}

	return nil
}

// GetProjectByID 根据ID获取项目（别名方法，兼容API层调用）
func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*document.Project, error) {
	return s.GetProject(ctx, projectID)
}

// GetProjectList 获取项目列表（兼容API层调用）
func (s *ProjectService) GetProjectList(ctx context.Context, userID, status string, limit, offset int64) ([]*document.Project, error) {
	// 构建请求
	req := &ListProjectsRequest{
		Page:     int(offset/limit) + 1,
		PageSize: int(limit),
		Status:   status,
	}

	// 使用上下文注入用户ID
	ctxWithUser := context.WithValue(ctx, "userID", userID)

	// 调用ListMyProjects
	resp, err := s.ListMyProjects(ctxWithUser, req)
	if err != nil {
		return nil, err
	}

	return resp.Projects, nil
}

// UpdateProjectByID 更新项目（兼容API层调用）
func (s *ProjectService) UpdateProjectByID(ctx context.Context, projectID, userID string, req *document.Project) error {
	// 构建更新请求
	updateReq := &UpdateProjectRequest{
		Title:    req.Title,
		Summary:  req.Summary,
		CoverURL: req.CoverURL,
		Category: req.Category,
		Tags:     req.Tags,
		Status:   string(req.Status),
	}

	// 使用上下文注入用户ID
	ctxWithUser := context.WithValue(ctx, "userID", userID)

	return s.UpdateProject(ctxWithUser, projectID, updateReq)
}

// DeleteProjectByID 删除项目（兼容API层调用）
func (s *ProjectService) DeleteProjectByID(ctx context.Context, projectID, userID string) error {
	// 使用上下文注入用户ID
	ctxWithUser := context.WithValue(ctx, "userID", userID)

	return s.DeleteProject(ctxWithUser, projectID)
}

// RestoreProjectByID 恢复已删除的项目
func (s *ProjectService) RestoreProjectByID(ctx context.Context, projectID, userID string) error {
	// 1. 查询项目（包括已删除的）
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "查询项目失败", "", err)
	}

	if project == nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorNotFound, "项目不存在", "", nil)
	}

	// 2. 权限检查（只有所有者可以恢复）
	if !project.IsOwner(userID) {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorForbidden, "只有项目所有者可以恢复项目", "", nil)
	}

	// 3. 恢复项目
	if err := s.projectRepo.Restore(ctx, projectID, userID); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "恢复项目失败", "", err)
	}

	// 4. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "project.restored",
			EventData: map[string]interface{}{
				"project_id": projectID,
				"author_id":  project.AuthorID,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// DeleteHard 物理删除项目（永久删除）
func (s *ProjectService) DeleteHard(ctx context.Context, projectID string) error {
	// 注意：硬删除不需要权限检查，通常由管理员或系统调用
	// 如果需要权限检查，应该在调用方（API层）进行

	// 1. 物理删除项目
	if err := s.projectRepo.HardDelete(ctx, projectID); err != nil {
		return pkgErrors.NewServiceError(s.serviceName, pkgErrors.ServiceErrorInternal, "物理删除项目失败", "", err)
	}

	// 2. 发布事件
	if s.eventBus != nil {
		s.eventBus.PublishAsync(ctx, &base.BaseEvent{
			EventType: "project.hard_deleted",
			EventData: map[string]interface{}{
				"project_id": projectID,
			},
			Timestamp: time.Now(),
			Source:    s.serviceName,
		})
	}

	return nil
}

// 私有方法
func (s *ProjectService) validateCreateProjectRequest(req *CreateProjectRequest) error {
	if req.Title == "" {
		return fmt.Errorf("项目标题不能为空")
	}
	if len(req.Title) > 100 {
		return fmt.Errorf("项目标题不能超过100字符")
	}
	return nil
}

// BaseService接口实现
func (s *ProjectService) Initialize(ctx context.Context) error {
	return nil
}

func (s *ProjectService) Health(ctx context.Context) error {
	return s.projectRepo.Health(ctx)
}

func (s *ProjectService) Close(ctx context.Context) error {
	return nil
}

func (s *ProjectService) GetServiceName() string {
	return s.serviceName
}

func (s *ProjectService) GetVersion() string {
	return s.version
}
