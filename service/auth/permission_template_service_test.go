package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authModel "Qingyu_backend/models/auth"
)

// MockPermissionTemplateRepository 模拟权限模板仓储
type MockPermissionTemplateRepository struct {
	mock.Mock
}

func (m *MockPermissionTemplateRepository) CreateTemplate(ctx context.Context, template *authModel.PermissionTemplate) error {
	args := m.Called(ctx, template)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockPermissionTemplateRepository) GetTemplateByID(ctx context.Context, templateID string) (*authModel.PermissionTemplate, error) {
	args := m.Called(ctx, templateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.PermissionTemplate), args.Error(1)
}

func (m *MockPermissionTemplateRepository) GetTemplateByCode(ctx context.Context, code string) (*authModel.PermissionTemplate, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.PermissionTemplate), args.Error(1)
}

func (m *MockPermissionTemplateRepository) UpdateTemplate(ctx context.Context, templateID string, updates map[string]interface{}) error {
	args := m.Called(ctx, templateID, updates)
	return args.Error(0)
}

func (m *MockPermissionTemplateRepository) DeleteTemplate(ctx context.Context, templateID string) error {
	args := m.Called(ctx, templateID)
	return args.Error(0)
}

func (m *MockPermissionTemplateRepository) ListTemplates(ctx context.Context) ([]*authModel.PermissionTemplate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.PermissionTemplate), args.Error(1)
}

func (m *MockPermissionTemplateRepository) ListTemplatesByCategory(ctx context.Context, category string) ([]*authModel.PermissionTemplate, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.PermissionTemplate), args.Error(1)
}

func (m *MockPermissionTemplateRepository) ApplyTemplateToRole(ctx context.Context, templateID, roleID string) error {
	args := m.Called(ctx, templateID, roleID)
	return args.Error(0)
}

func (m *MockPermissionTemplateRepository) GetSystemTemplates(ctx context.Context) ([]*authModel.PermissionTemplate, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.PermissionTemplate), args.Error(1)
}

func (m *MockPermissionTemplateRepository) InitializeSystemTemplates(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPermissionTemplateRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockRoleRepository 模拟角色仓储
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) CreateRole(ctx context.Context, role *authModel.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetRole(ctx context.Context, roleID string) (*authModel.Role, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleByName(ctx context.Context, name string) (*authModel.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModel.Role), args.Error(1)
}

func (m *MockRoleRepository) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}) error {
	args := m.Called(ctx, roleID, updates)
	return args.Error(0)
}

func (m *MockRoleRepository) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) ListRoles(ctx context.Context) ([]*authModel.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Role), args.Error(1)
}

func (m *MockRoleRepository) AssignUserRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]*authModel.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*authModel.Role), args.Error(1)
}

func (m *MockRoleRepository) HasUserRole(ctx context.Context, userID, roleID string) (bool, error) {
	args := m.Called(ctx, userID, roleID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRoleRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRoleRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ============ 测试用例 ============

// TestPermissionTemplateService_CreateTemplate_Success 测试成功创建模板
func TestPermissionTemplateService_CreateTemplate_Success(t *testing.T) {
	// 设置
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()
	req := &CreateTemplateRequest{
		Name:        "测试模板",
		Code:        "test_template",
		Description: "测试用模板",
		Permissions: []string{"user.read", "book.read"},
		Category:    authModel.CategoryCustom,
	}

	// mock期望
	mockTemplateRepo.On("GetTemplateByCode", ctx, req.Code).Return(nil, errors.New("not found"))
	mockTemplateRepo.On("CreateTemplate", ctx, mock.AnythingOfType("*auth.PermissionTemplate")).Return(nil)

	// 执行
	resp, err := service.CreateTemplate(ctx, req)

	// 验证
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, req.Code, resp.Code)
	assert.Equal(t, req.Permissions, resp.Permissions)
	assert.False(t, resp.IsSystem)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_CreateTemplate_EmptyName 测试名称为空
func TestPermissionTemplateService_CreateTemplate_EmptyName(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()
	req := &CreateTemplateRequest{
		Name:        "",
		Code:        "test",
		Permissions: []string{"user.read"},
	}

	_, err := service.CreateTemplate(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "模板名称不能为空")
}

// TestPermissionTemplateService_CreateTemplate_CodeExists 测试代码已存在
func TestPermissionTemplateService_CreateTemplate_CodeExists(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()
	req := &CreateTemplateRequest{
		Name:        "测试模板",
		Code:        "existing_code",
		Permissions: []string{"user.read"},
	}

	existingTemplate := &authModel.PermissionTemplate{
		ID:     "existing_id",
		Code:   "existing_code",
		Name:   "已存在的模板",
		IsSystem: false,
	}

	mockTemplateRepo.On("GetTemplateByCode", ctx, req.Code).Return(existingTemplate, nil)

	_, err := service.CreateTemplate(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "模板代码已存在")

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_UpdateSystemTemplate 测试不能更新系统模板
func TestPermissionTemplateService_UpdateSystemTemplate(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()
	systemTemplate := &authModel.PermissionTemplate{
		ID:       "system_id",
		Name:     "系统模板",
		Code:     "system_template",
		IsSystem: true,
	}

	mockTemplateRepo.On("GetTemplateByID", ctx, "system_id").Return(systemTemplate, nil)

	req := &UpdateTemplateRequest{
		Name: "新名称",
	}

	err := service.UpdateTemplate(ctx, "system_id", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不能修改系统模板")

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_DeleteSystemTemplate 测试不能删除系统模板
func TestPermissionTemplateService_DeleteSystemTemplate(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	// mock仓储会在删除时检查是否是系统模板
	mockTemplateRepo.On("DeleteTemplate", ctx, "system_id").Return(errors.New("cannot delete system template"))

	err := service.DeleteTemplate(ctx, "system_id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "删除模板失败")

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_ListTemplates_WithCategory 测试按分类列出模板
func TestPermissionTemplateService_ListTemplates_WithCategory(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	templates := []*authModel.PermissionTemplate{
		{
			ID:          "1",
			Name:        "读者模板1",
			Code:        "reader_1",
			Category:    authModel.CategoryReader,
			Permissions: []string{"book.read"},
			IsSystem:    false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "读者模板2",
			Code:        "reader_2",
			Category:    authModel.CategoryReader,
			Permissions: []string{"document.read"},
			IsSystem:    false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockTemplateRepo.On("ListTemplatesByCategory", ctx, authModel.CategoryReader).Return(templates, nil)

	result, err := service.ListTemplates(ctx, authModel.CategoryReader)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "读者模板1", result[0].Name)
	assert.Equal(t, "读者模板2", result[1].Name)
	assert.Equal(t, authModel.CategoryReader, result[0].Category)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_ListTemplates_All 测试列出所有模板
func TestPermissionTemplateService_ListTemplates_All(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	templates := []*authModel.PermissionTemplate{
		{
			ID:          "1",
			Name:        "模板1",
			Code:        "template_1",
			Category:    authModel.CategoryReader,
			Permissions: []string{"book.read"},
			IsSystem:    false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "模板2",
			Code:        "template_2",
			Category:    authModel.CategoryAuthor,
			Permissions: []string{"book.write"},
			IsSystem:    false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockTemplateRepo.On("ListTemplates", ctx).Return(templates, nil)

	result, err := service.ListTemplates(ctx, "")

	require.NoError(t, err)
	assert.Len(t, result, 2)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_ApplyTemplate 测试应用模板
func TestPermissionTemplateService_ApplyTemplate(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	template := &authModel.PermissionTemplate{
		ID:          "template_id",
		Name:        "测试模板",
		Code:        "test_template",
		Permissions: []string{"user.read", "book.write"},
		IsSystem:    false,
	}

	mockTemplateRepo.On("GetTemplateByID", ctx, "template_id").Return(template, nil)
	mockTemplateRepo.On("ApplyTemplateToRole", ctx, "template_id", "role_id").Return(nil)

	err := service.ApplyTemplate(ctx, "template_id", "role_id")

	require.NoError(t, err)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_InitializeSystemTemplates 测试初始化系统模板
func TestPermissionTemplateService_InitializeSystemTemplates(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	mockTemplateRepo.On("InitializeSystemTemplates", ctx).Return(nil)

	err := service.InitializeSystemTemplates(ctx)

	require.NoError(t, err)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_GetTemplate 测试获取模板
func TestPermissionTemplateService_GetTemplate(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	template := &authModel.PermissionTemplate{
		ID:          "template_id",
		Name:        "测试模板",
		Code:        "test_code",
		Description: "测试描述",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockTemplateRepo.On("GetTemplateByID", ctx, "template_id").Return(template, nil)

	result, err := service.GetTemplate(ctx, "template_id")

	require.NoError(t, err)
	assert.Equal(t, template.Name, result.Name)
	assert.Equal(t, template.Code, result.Code)
	assert.Equal(t, template.Description, result.Description)
	assert.Equal(t, template.Permissions, result.Permissions)
	assert.False(t, result.IsSystem)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_GetTemplateByCode 测试根据代码获取模板
func TestPermissionTemplateService_GetTemplateByCode(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	template := &authModel.PermissionTemplate{
		ID:          "template_id",
		Name:        "测试模板",
		Code:        "test_code",
		Description: "测试描述",
		Permissions: []string{"user.read"},
		IsSystem:    false,
		Category:    authModel.CategoryCustom,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockTemplateRepo.On("GetTemplateByCode", ctx, "test_code").Return(template, nil)

	result, err := service.GetTemplateByCode(ctx, "test_code")

	require.NoError(t, err)
	assert.Equal(t, template.Code, result.Code)
	assert.Equal(t, template.Name, result.Name)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_UpdateTemplate 测试更新模板
func TestPermissionTemplateService_UpdateTemplate(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	template := &authModel.PermissionTemplate{
		ID:       "template_id",
		Name:     "原始名称",
		Code:     "test_code",
		IsSystem: false,
	}

	mockTemplateRepo.On("GetTemplateByID", ctx, "template_id").Return(template, nil)
	mockTemplateRepo.On("UpdateTemplate", ctx, "template_id", mock.MatchedBy(func(updates map[string]interface{}) bool {
		return updates["name"] == "新名称"
	})).Return(nil)

	req := &UpdateTemplateRequest{
		Name: "新名称",
	}

	err := service.UpdateTemplate(ctx, "template_id", req)

	require.NoError(t, err)

	mockTemplateRepo.AssertExpectations(t)
}

// TestPermissionTemplateService_UpdateTemplate_EmptyUpdates 测试空更新
func TestPermissionTemplateService_UpdateTemplate_EmptyUpdates(t *testing.T) {
	mockTemplateRepo := new(MockPermissionTemplateRepository)
	mockRoleRepo := new(MockRoleRepository)
	service := NewPermissionTemplateService(mockTemplateRepo, mockRoleRepo)

	ctx := context.Background()

	template := &authModel.PermissionTemplate{
		ID:       "template_id",
		Name:     "原始名称",
		Code:     "test_code",
		IsSystem: false,
	}

	mockTemplateRepo.On("GetTemplateByID", ctx, "template_id").Return(template, nil)

	req := &UpdateTemplateRequest{}

	err := service.UpdateTemplate(ctx, "template_id", req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "没有要更新的内容")

	mockTemplateRepo.AssertExpectations(t)
}
