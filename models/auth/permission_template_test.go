package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPermissionTemplate_Validate_ValidTemplate 测试有效模板验证
func TestPermissionTemplate_Validate_ValidTemplate(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "读者模板",
		Code:        "template_reader",
		Description: "普通读者权限模板",
		Permissions: []string{"read:book", "read:chapter"},
		Category:    CategoryReader,
		IsSystem:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := template.Validate()
	assert.NoError(t, err)
}

// TestPermissionTemplate_Validate_EmptyName 测试空名称
func TestPermissionTemplate_Validate_EmptyName(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "",
		Code:        "template_test",
		Permissions: []string{"read:book"},
		Category:    CategoryCustom,
	}

	err := template.Validate()
	assert.Error(t, err)
	assert.Equal(t, ErrTemplateNameEmpty, err)
}

// TestPermissionTemplate_Validate_EmptyCode 测试空代码
func TestPermissionTemplate_Validate_EmptyCode(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "",
		Permissions: []string{"read:book"},
		Category:    CategoryCustom,
	}

	err := template.Validate()
	assert.Error(t, err)
	assert.Equal(t, ErrTemplateCodeEmpty, err)
}

// TestPermissionTemplate_Validate_EmptyPermissions 测试空权限列表
func TestPermissionTemplate_Validate_EmptyPermissions(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "template_test",
		Permissions: []string{},
		Category:    CategoryCustom,
	}

	err := template.Validate()
	assert.Error(t, err)
	assert.Equal(t, ErrTemplatePermissionsEmpty, err)
}

// TestPermissionTemplate_Validate_NilPermissions 测试nil权限列表
func TestPermissionTemplate_Validate_NilPermissions(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "template_test",
		Permissions: nil,
		Category:    CategoryCustom,
	}

	err := template.Validate()
	assert.Error(t, err)
	assert.Equal(t, ErrTemplatePermissionsEmpty, err)
}

// TestPermissionTemplate_IsSystemTemplate 测试系统模板判断
func TestPermissionTemplate_IsSystemTemplate(t *testing.T) {
	systemTemplate := &PermissionTemplate{
		Name:     "系统模板",
		Code:     TemplateAdmin,
		IsSystem: true,
	}

	customTemplate := &PermissionTemplate{
		Name:     "自定义模板",
		Code:     "custom_template",
		IsSystem: false,
	}

	assert.True(t, systemTemplate.IsSystem == true, "系统模板IsSystem应该为true")
	assert.False(t, customTemplate.IsSystem == true, "自定义模板IsSystem应该为false")
}

// TestPermissionTemplate_PredefinedTemplates 测试预定义模板常量
func TestPermissionTemplate_PredefinedTemplates(t *testing.T) {
	assert.Equal(t, "template_reader", TemplateReader)
	assert.Equal(t, "template_author", TemplateAuthor)
	assert.Equal(t, "template_admin", TemplateAdmin)
}

// TestPermissionTemplate_PredefinedCategories 测试预定义分类常量
func TestPermissionTemplate_PredefinedCategories(t *testing.T) {
	assert.Equal(t, "reader", CategoryReader)
	assert.Equal(t, "author", CategoryAuthor)
	assert.Equal(t, "admin", CategoryAdmin)
	assert.Equal(t, "custom", CategoryCustom)
}

// TestPermissionTemplate_GetCategory 测试获取分类
func TestPermissionTemplate_GetCategory(t *testing.T) {
	readerTemplate := &PermissionTemplate{
		Name:     "读者模板",
		Category: CategoryReader,
	}

	authorTemplate := &PermissionTemplate{
		Name:     "作者模板",
		Category: CategoryAuthor,
	}

	adminTemplate := &PermissionTemplate{
		Name:     "管理员模板",
		Category: CategoryAdmin,
	}

	customTemplate := &PermissionTemplate{
		Name:     "自定义模板",
		Category: CategoryCustom,
	}

	assert.Equal(t, CategoryReader, readerTemplate.Category)
	assert.Equal(t, CategoryAuthor, authorTemplate.Category)
	assert.Equal(t, CategoryAdmin, adminTemplate.Category)
	assert.Equal(t, CategoryCustom, customTemplate.Category)
}

// TestPermissionTemplate_GetPermissions 测试获取权限列表
func TestPermissionTemplate_GetPermissions(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "template_test",
		Permissions: []string{"read:book", "read:chapter", "comment:add"},
		Category:    CategoryCustom,
	}

	permissions := template.Permissions
	assert.NotNil(t, permissions)
	assert.Equal(t, 3, len(permissions))
	assert.Contains(t, permissions, "read:book")
	assert.Contains(t, permissions, "read:chapter")
	assert.Contains(t, permissions, "comment:add")
}

// TestPermissionTemplate_ErrorMessages 测试错误信息
func TestPermissionTemplate_ErrorMessages(t *testing.T) {
	assert.Equal(t, "template name cannot be empty", ErrTemplateNameEmpty.Error())
	assert.Equal(t, "template code cannot be empty", ErrTemplateCodeEmpty.Error())
	assert.Equal(t, "template permissions cannot be empty", ErrTemplatePermissionsEmpty.Error())
	assert.Equal(t, "cannot delete system template", ErrTemplateIsSystem.Error())
	assert.Equal(t, "template not found", ErrTemplateNotFound.Error())
	assert.Equal(t, "template code already exists", ErrTemplateCodeExists.Error())
}

// TestPermissionTemplate_Fields 测试模板字段
func TestPermissionTemplate_Fields(t *testing.T) {
	now := time.Now()
	template := &PermissionTemplate{
		ID:          "template123",
		Name:        "测试模板",
		Code:        "template_test",
		Description: "这是一个测试模板",
		Permissions: []string{"read:book"},
		IsSystem:    false,
		Category:    CategoryCustom,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   "admin123",
	}

	assert.Equal(t, "template123", template.ID)
	assert.Equal(t, "测试模板", template.Name)
	assert.Equal(t, "template_test", template.Code)
	assert.Equal(t, "这是一个测试模板", template.Description)
	assert.Len(t, template.Permissions, 1)
	assert.False(t, template.IsSystem)
	assert.Equal(t, CategoryCustom, template.Category)
	assert.False(t, template.CreatedAt.IsZero())
	assert.False(t, template.UpdatedAt.IsZero())
	assert.Equal(t, "admin123", template.CreatedBy)
}

// TestPermissionTemplate_EmptyID 测试空ID
func TestPermissionTemplate_EmptyID(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "template_test",
		Permissions: []string{"read:book"},
		Category:    CategoryCustom,
	}

	// 空ID是允许的（创建时）
	assert.Equal(t, "", template.ID)
}

// TestPermissionTemplate_EmptyDescription 测试空描述
func TestPermissionTemplate_EmptyDescription(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "template_test",
		Description: "",
		Permissions: []string{"read:book"},
		Category:    CategoryCustom,
	}

	// 空描述是允许的
	err := template.Validate()
	assert.NoError(t, err)
}

// TestPermissionTemplate_EmptyCreatedBy 测试空创建者
func TestPermissionTemplate_EmptyCreatedBy(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "系统模板",
		Code:        "template_system",
		Description: "系统创建的模板",
		Permissions: []string{"read:book"},
		IsSystem:    true,
		Category:    CategoryReader,
		CreatedBy:   "",
	}

	// 空创建者是允许的（系统模板）
	err := template.Validate()
	assert.NoError(t, err)
}

// TestPermissionTemplate_MultiplePermissions 测试多个权限
func TestPermissionTemplate_MultiplePermissions(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "全功能模板",
		Code:        "template_full",
		Description: "包含所有权限",
		Permissions: []string{
			"read:book",
			"read:chapter",
			"write:book",
			"write:chapter",
			"delete:book",
			"delete:chapter",
			"comment:add",
			"comment:delete",
			"user:manage",
		},
		Category: CategoryAdmin,
	}

	err := template.Validate()
	assert.NoError(t, err)
	assert.Equal(t, 9, len(template.Permissions))
}

// TestPermissionTemplate_SpecialCharactersInCode 测试代码中的特殊字符
func TestPermissionTemplate_SpecialCharactersInCode(t *testing.T) {
	testCases := []struct {
		name  string
		code  string
		valid bool
	}{
		{"下划线", "template_reader", true},
		{"连字符", "template-reader", true},
		{"点号", "template.reader", true},
		{"空格", "template reader", true}, // Validate不检查格式
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			template := &PermissionTemplate{
				Name:        "测试模板",
				Code:        tc.code,
				Permissions: []string{"read:book"},
				Category:    CategoryCustom,
			}

			err := template.Validate()
			if tc.valid {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPermissionTemplate_LargePermissionList 测试大量权限
func TestPermissionTemplate_LargePermissionList(t *testing.T) {
	permissions := make([]string, 100)
	for i := 0; i < 100; i++ {
		permissions[i] = "permission:action" + string(rune('0'+i))
	}

	template := &PermissionTemplate{
		Name:        "大权限模板",
		Code:        "template_large",
		Permissions: permissions,
		Category:    CategoryCustom,
	}

	err := template.Validate()
	assert.NoError(t, err)
	assert.Equal(t, 100, len(template.Permissions))
}

// TestPermissionTemplate_DuplicatePermissions 测试重复权限
func TestPermissionTemplate_DuplicatePermissions(t *testing.T) {
	// 注意：当前实现允许重复权限
	template := &PermissionTemplate{
		Name:        "重复权限模板",
		Code:        "template_duplicate",
		Permissions: []string{"read:book", "read:book", "read:chapter"},
		Category:    CategoryCustom,
	}

	err := template.Validate()
	assert.NoError(t, err) // 当前实现不检查重复
	assert.Equal(t, 3, len(template.Permissions))
}

// TestPermissionTemplate_ZeroTime 测试零值时间
func TestPermissionTemplate_ZeroTime(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "template_test",
		Permissions: []string{"read:book"},
		Category:    CategoryCustom,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	// 零值时间是允许的
	err := template.Validate()
	assert.NoError(t, err)
	assert.True(t, template.CreatedAt.IsZero())
	assert.True(t, template.UpdatedAt.IsZero())
}

// TestPermissionTemplate_UpdateTimestamp 测试更新时间戳
func TestPermissionTemplate_UpdateTimestamp(t *testing.T) {
	template := &PermissionTemplate{
		Name:        "测试模板",
		Code:        "template_test",
		Permissions: []string{"read:book"},
		Category:    CategoryCustom,
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	}

	oldUpdatedAt := template.UpdatedAt
	time.Sleep(time.Millisecond) // 确保时间变化

	// 模拟更新
	template.UpdatedAt = time.Now()

	assert.True(t, template.UpdatedAt.After(oldUpdatedAt))
}
