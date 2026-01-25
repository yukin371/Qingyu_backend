package document

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"Qingyu_backend/models/writer"
	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// TemplateService 模板服务
type TemplateService struct {
	repo   writerRepo.TemplateRepository
	logger *zap.Logger
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(repo writerRepo.TemplateRepository, logger *zap.Logger) *TemplateService {
	return &TemplateService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTemplateRequest 创建模板请求
type CreateTemplateRequest struct {
	WorkspaceID string
	ProjectID   *primitive.ObjectID
	Name        string
	Description string
	Type        writer.TemplateType
	Category    string
	Content     string
	Variables   []writer.TemplateVariable
	CreatedBy   string
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name        *string
	Description *string
	Category    *string
	Content     *string
	Variables   *[]writer.TemplateVariable
}

// ListTemplatesRequest 列出模板请求
type ListTemplatesRequest struct {
	WorkspaceID *string
	ProjectID   *primitive.ObjectID
	Type        *writer.TemplateType
	Status      *writer.TemplateStatus
	Category    string
	Keyword     string
	IsSystem    *bool
	Page        int
	PageSize    int
	SortBy      string
	SortOrder   int
}

// CreateTemplate 创建模板
func (s *TemplateService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*writer.Template, error) {
	// 验证请求参数
	if req == nil {
		s.logger.Error("创建模板失败：请求参数为空")
		return nil, fmt.Errorf("请求参数不能为空")
	}

	if req.Name == "" {
		s.logger.Error("创建模板失败：模板名称为空")
		return nil, fmt.Errorf("模板名称不能为空")
	}

	if req.Content == "" {
		s.logger.Error("创建模板失败：模板内容为空")
		return nil, fmt.Errorf("模板内容不能为空")
	}

	if req.CreatedBy == "" {
		s.logger.Error("创建模板失败：创建者为空")
		return nil, fmt.Errorf("创建者不能为空")
	}

	// 检查名称唯一性（同项目下）
	exists, err := s.repo.ExistsByName(ctx, req.ProjectID, req.Name, nil)
	if err != nil {
		s.logger.Error("检查模板名称唯一性失败", zap.Error(err))
		return nil, fmt.Errorf("检查模板名称失败: %w", err)
	}

	if exists {
		s.logger.Warn("模板名称已存在", zap.String("name", req.Name))
		return nil, fmt.Errorf("模板名称 '%s' 已存在", req.Name)
	}

	// 验证模板语法
	if err := s.ValidateTemplateContent(req.Content); err != nil {
		s.logger.Error("模板内容验证失败", zap.Error(err))
		return nil, fmt.Errorf("模板内容验证失败: %w", err)
	}

	// 创建模板对象
	template := &writer.Template{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Content:     req.Content,
		Variables:   req.Variables,
		Status:      writer.TemplateStatusActive,
		Version:     1,
		IsSystem:    false,
		CreatedBy:   req.CreatedBy,
	}

	// 验证模板变量
	if err := template.ValidateVariables(); err != nil {
		s.logger.Error("模板变量验证失败", zap.Error(err))
		return nil, fmt.Errorf("模板变量验证失败: %w", err)
	}

	// 调用TouchForCreate初始化时间戳和ID
	template.TouchForCreate()

	// 保存到数据库
	if err := s.repo.Create(ctx, template); err != nil {
		s.logger.Error("创建模板失败", zap.Error(err))
		return nil, fmt.Errorf("创建模板失败: %w", err)
	}

	s.logger.Info("创建模板成功",
		zap.String("template_id", template.ID.Hex()),
		zap.String("name", template.Name),
		zap.String("type", string(template.Type)))

	return template, nil
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(ctx context.Context, id primitive.ObjectID, req *UpdateTemplateRequest) (*writer.Template, error) {
	// 验证请求参数
	if req == nil {
		s.logger.Error("更新模板失败：请求参数为空")
		return nil, fmt.Errorf("请求参数不能为空")
	}

	// 获取现有模板
	template, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取模板失败", zap.String("id", id.Hex()), zap.Error(err))
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	if template == nil {
		s.logger.Warn("模板不存在", zap.String("id", id.Hex()))
		return nil, fmt.Errorf("模板不存在")
	}

	// 检查是否为系统模板
	if template.IsSystem {
		s.logger.Warn("系统模板不允许修改", zap.String("id", id.Hex()))
		return nil, fmt.Errorf("系统模板不允许修改")
	}

	// 更新字段
	if req.Name != nil {
		// 检查名称唯一性（排除自己）
		exists, err := s.repo.ExistsByName(ctx, template.ProjectID, *req.Name, &id)
		if err != nil {
			s.logger.Error("检查模板名称唯一性失败", zap.Error(err))
			return nil, fmt.Errorf("检查模板名称失败: %w", err)
		}

		if exists {
			s.logger.Warn("模板名称已存在", zap.String("name", *req.Name))
			return nil, fmt.Errorf("模板名称 '%s' 已存在", *req.Name)
		}

		template.Name = *req.Name
	}

	if req.Description != nil {
		template.Description = *req.Description
	}

	if req.Category != nil {
		template.Category = *req.Category
	}

	if req.Content != nil {
		// 验证模板语法
		if err := s.ValidateTemplateContent(*req.Content); err != nil {
			s.logger.Error("模板内容验证失败", zap.Error(err))
			return nil, fmt.Errorf("模板内容验证失败: %w", err)
		}

		template.Content = *req.Content
	}

	if req.Variables != nil {
		template.Variables = *req.Variables

		// 验证模板变量
		if err := template.ValidateVariables(); err != nil {
			s.logger.Error("模板变量验证失败", zap.Error(err))
			return nil, fmt.Errorf("模板变量验证失败: %w", err)
		}
	}

	// 调用TouchForUpdate更新时间戳和版本号
	template.TouchForUpdate()

	// 保存更新
	if err := s.repo.Update(ctx, template); err != nil {
		s.logger.Error("更新模板失败", zap.Error(err))
		return nil, fmt.Errorf("更新模板失败: %w", err)
	}

	s.logger.Info("更新模板成功",
		zap.String("template_id", template.ID.Hex()),
		zap.String("name", template.Name),
		zap.Int("version", template.Version))

	return template, nil
}

// DeleteTemplate 删除模板（软删除）
func (s *TemplateService) DeleteTemplate(ctx context.Context, id primitive.ObjectID, userID string) error {
	// 获取现有模板
	template, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取模板失败", zap.String("id", id.Hex()), zap.Error(err))
		return fmt.Errorf("获取模板失败: %w", err)
	}

	if template == nil {
		s.logger.Warn("模板不存在", zap.String("id", id.Hex()))
		return fmt.Errorf("模板不存在")
	}

	// 检查是否为系统模板
	if template.IsSystem {
		s.logger.Warn("系统模板不允许删除", zap.String("id", id.Hex()))
		return fmt.Errorf("系统模板不允许删除")
	}

	// 检查权限（createdBy匹配）
	if template.CreatedBy != userID {
		s.logger.Warn("无权限删除模板",
			zap.String("template_id", id.Hex()),
			zap.String("created_by", template.CreatedBy),
			zap.String("user_id", userID))
		return fmt.Errorf("无权限删除此模板")
	}

	// 软删除
	template.SoftDelete()

	// 保存更新
	if err := s.repo.Update(ctx, template); err != nil {
		s.logger.Error("删除模板失败", zap.Error(err))
		return fmt.Errorf("删除模板失败: %w", err)
	}

	s.logger.Info("删除模板成功",
		zap.String("template_id", id.Hex()),
		zap.String("name", template.Name))

	return nil
}

// GetTemplate 获取模板详情
func (s *TemplateService) GetTemplate(ctx context.Context, id primitive.ObjectID) (*writer.Template, error) {
	template, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取模板失败", zap.String("id", id.Hex()), zap.Error(err))
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}

	if template == nil {
		s.logger.Warn("模板不存在", zap.String("id", id.Hex()))
		return nil, fmt.Errorf("模板不存在")
	}

	return template, nil
}

// ListTemplates 列出模板（支持分页、过滤）
func (s *TemplateService) ListTemplates(ctx context.Context, req *ListTemplatesRequest) ([]*writer.Template, int64, error) {
	// 构建过滤器
	filter := &writerRepo.TemplateFilter{
		Type:      req.Type,
		Status:    req.Status,
		Category:  req.Category,
		Keyword:   req.Keyword,
		IsSystem:  req.IsSystem,
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	var templates []*writer.Template
	var err error

	// 根据查询条件选择不同的查询方法
	if req.ProjectID != nil && !req.ProjectID.IsZero() {
		// 查询项目模板
		templates, err = s.repo.ListByProject(ctx, *req.ProjectID, filter)
	} else {
		// 查询全局模板
		templates, err = s.repo.ListGlobal(ctx, filter)
	}

	if err != nil {
		s.logger.Error("查询模板列表失败", zap.Error(err))
		return nil, 0, fmt.Errorf("查询模板列表失败: %w", err)
	}

	// 获取总数（这里简化处理，实际可能需要单独的count方法）
	total := int64(len(templates))

	s.logger.Info("查询模板列表成功",
		zap.Int64("total", total),
		zap.Int("count", len(templates)))

	return templates, total, nil
}

// ValidateTemplateContent 验证模板语法
func (s *TemplateService) ValidateTemplateContent(content string) error {
	if content == "" {
		return fmt.Errorf("模板内容不能为空")
	}

	// 定义变量正则表达式：{{var.xxx}}
	varRegex := regexp.MustCompile(`\{\{var\.([a-zA-Z][a-zA-Z0-9_]*)\}\}`)

	// 查找所有变量引用
	matches := varRegex.FindAllStringSubmatch(content, -1)

	// 验证变量名格式
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		varName := match[1]

		// 验证变量名是否符合规范
		if !isValidVariableName(varName) {
			return fmt.Errorf("变量名 '%s' 不符合规范：必须以字母开头，只能包含字母、数字和下划线", varName)
		}

		s.logger.Debug("找到模板变量", zap.String("variable", varName))
	}

	// 检查是否有未闭合的标签
	openBraces := strings.Count(content, "{{")
	closeBraces := strings.Count(content, "}}")

	if openBraces != closeBraces {
		return fmt.Errorf("模板语法错误：未闭合的标签（{{ }}不匹配）")
	}

	s.logger.Debug("模板内容验证通过",
		zap.Int("variable_count", len(matches)),
		zap.Int("open_braces", openBraces),
		zap.Int("close_braces", closeBraces))

	return nil
}

// isValidVariableName 验证变量名是否有效
func isValidVariableName(name string) bool {
	if name == "" {
		return false
	}

	// 变量名必须以字母开头
	if name[0] < 'a' || name[0] > 'z' {
		if name[0] < 'A' || name[0] > 'Z' {
			return false
		}
	}

	// 变量名只能包含字母、数字和下划线
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}

	return true
}
