package auth

import (
	"context"
	"fmt"

	authModel "Qingyu_backend/models/auth"
	repoAuth "Qingyu_backend/repository/interfaces/auth"
)

// PermissionTemplateService 权限模板服务接口
type PermissionTemplateService interface {
	// 模板管理
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*TemplateResponse, error)
	GetTemplate(ctx context.Context, templateID string) (*TemplateResponse, error)
	GetTemplateByCode(ctx context.Context, code string) (*TemplateResponse, error)
	UpdateTemplate(ctx context.Context, templateID string, req *UpdateTemplateRequest) error
	DeleteTemplate(ctx context.Context, templateID string) error
	ListTemplates(ctx context.Context, category string) ([]*TemplateResponse, error)

	// 模板应用
	ApplyTemplate(ctx context.Context, templateID, roleID string) error
	InitializeSystemTemplates(ctx context.Context) error
}

// PermissionTemplateServiceImpl 权限模板服务实现
type PermissionTemplateServiceImpl struct {
	templateRepo repoAuth.PermissionTemplateRepository
	roleRepo     repoAuth.RoleRepository
}

// NewPermissionTemplateService 创建权限模板服务
func NewPermissionTemplateService(
	templateRepo repoAuth.PermissionTemplateRepository,
	roleRepo repoAuth.RoleRepository,
) PermissionTemplateService {
	return &PermissionTemplateServiceImpl{
		templateRepo: templateRepo,
		roleRepo:     roleRepo,
	}
}

// ============ 请求/响应结构 ============

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Code        string   `json:"code" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions" binding:"required"`
	Category    string   `json:"category"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	Category    string   `json:"category"`
}

// TemplateResponse 模板响应
type TemplateResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	IsSystem    bool     `json:"is_system"`
	Category    string   `json:"category"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// ============ 模板管理实现 ============

// CreateTemplate 创建模板
func (s *PermissionTemplateServiceImpl) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*TemplateResponse, error) {
	// 1. 验证请求
	if req.Name == "" {
		return nil, fmt.Errorf("模板名称不能为空")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("模板代码不能为空")
	}
	if len(req.Permissions) == 0 {
		return nil, fmt.Errorf("权限列表不能为空")
	}

	// 2. 检查代码是否已存在
	existing, _ := s.templateRepo.GetTemplateByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("模板代码已存在: %s", req.Code)
	}

	// 3. 创建模板
	template := &authModel.PermissionTemplate{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Permissions: req.Permissions,
		IsSystem:    false, // 用户创建的不是系统模板
		Category:    req.Category,
	}

	if err := s.templateRepo.CreateTemplate(ctx, template); err != nil {
		return nil, fmt.Errorf("创建模板失败: %w", err)
	}

	// 4. 转换为响应
	return convertToTemplateResponse(template), nil
}

// GetTemplate 获取模板
func (s *PermissionTemplateServiceImpl) GetTemplate(ctx context.Context, templateID string) (*TemplateResponse, error) {
	template, err := s.templateRepo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	return convertToTemplateResponse(template), nil
}

// GetTemplateByCode 根据代码获取模板
func (s *PermissionTemplateServiceImpl) GetTemplateByCode(ctx context.Context, code string) (*TemplateResponse, error) {
	template, err := s.templateRepo.GetTemplateByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	return convertToTemplateResponse(template), nil
}

// UpdateTemplate 更新模板
func (s *PermissionTemplateServiceImpl) UpdateTemplate(ctx context.Context, templateID string, req *UpdateTemplateRequest) error {
	// 1. 获取模板
	template, err := s.templateRepo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return fmt.Errorf("模板不存在: %w", err)
	}

	// 2. 检查是否是系统模板
	if template.IsSystem {
		return fmt.Errorf("不能修改系统模板: %s", template.Name)
	}

	// 3. 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Permissions != nil {
		updates["permissions"] = req.Permissions
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}

	if len(updates) == 0 {
		return fmt.Errorf("没有要更新的内容")
	}

	// 4. 更新模板
	if err := s.templateRepo.UpdateTemplate(ctx, templateID, updates); err != nil {
		return fmt.Errorf("更新模板失败: %w", err)
	}

	return nil
}

// DeleteTemplate 删除模板
func (s *PermissionTemplateServiceImpl) DeleteTemplate(ctx context.Context, templateID string) error {
	// 删除模板（Repository会检查是否是系统模板）
	if err := s.templateRepo.DeleteTemplate(ctx, templateID); err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}

	return nil
}

// ListTemplates 列出模板
func (s *PermissionTemplateServiceImpl) ListTemplates(ctx context.Context, category string) ([]*TemplateResponse, error) {
	var templates []*authModel.PermissionTemplate
	var err error

	if category != "" {
		templates, err = s.templateRepo.ListTemplatesByCategory(ctx, category)
	} else {
		templates, err = s.templateRepo.ListTemplates(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("获取模板列表失败: %w", err)
	}

	// 转换为响应格式
	result := make([]*TemplateResponse, len(templates))
	for i, template := range templates {
		result[i] = convertToTemplateResponse(template)
	}

	return result, nil
}

// ============ 模板应用实现 ============

// ApplyTemplate 应用模板到角色
func (s *PermissionTemplateServiceImpl) ApplyTemplate(ctx context.Context, templateID, roleID string) error {
	// 1. 获取模板
	_, err := s.templateRepo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return fmt.Errorf("模板不存在: %w", err)
	}

	// 2. 应用模板到角色
	if err := s.templateRepo.ApplyTemplateToRole(ctx, templateID, roleID); err != nil {
		return fmt.Errorf("应用模板失败: %w", err)
	}

	return nil
}

// InitializeSystemTemplates 初始化系统预设模板
func (s *PermissionTemplateServiceImpl) InitializeSystemTemplates(ctx context.Context) error {
	if err := s.templateRepo.InitializeSystemTemplates(ctx); err != nil {
		return fmt.Errorf("初始化系统模板失败: %w", err)
	}

	return nil
}

// ============ 辅助函数 ============

// convertToTemplateResponse 转换为响应格式
func convertToTemplateResponse(template *authModel.PermissionTemplate) *TemplateResponse {
	return &TemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Code:        template.Code,
		Description: template.Description,
		Permissions: template.Permissions,
		IsSystem:    template.IsSystem,
		Category:    template.Category,
		CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
