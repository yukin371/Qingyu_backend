package writer

import (
	"errors"
	"time"

	"Qingyu_backend/models/writer/base"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Template 模板模型
type Template struct {
	base.IdentifiedEntity `bson:",inline"` // ID
	base.Timestamps       `bson:",inline"` // CreatedAt, UpdatedAt, DeletedAt

	// 项目关联（可选，为空时表示全局模板）
	ProjectID *primitive.ObjectID `bson:"project_id,omitempty" json:"projectId,omitempty"`

	// 基本信息
	Name        string       `bson:"name" json:"name" validate:"required,max=200"`
	Description string       `bson:"description" json:"description" validate:"max=1000"`
	Type        TemplateType `bson:"type" json:"type" validate:"required"`
	Category    string       `bson:"category" json:"category"`

	// 模板内容
	Content   string             `bson:"content" json:"content" validate:"required"`
	Variables []TemplateVariable `bson:"variables" json:"variables"`

	// 状态和版本
	Status  TemplateStatus `bson:"status" json:"status"`
	Version int            `bson:"version" json:"version"`
	IsSystem bool          `bson:"is_system" json:"isSystem"`

	// 创建者
	CreatedBy string `bson:"created_by" json:"createdBy" validate:"required"`
}

// TemplateType 模板类型
type TemplateType string

const (
	TemplateTypeChapter TemplateType = "chapter" // 章节模板
	TemplateTypeOutline TemplateType = "outline" // 大纲模板
	TemplateTypeSetting TemplateType = "setting" // 设定模板
)

// TemplateStatus 模板状态
type TemplateStatus string

const (
	TemplateStatusActive     TemplateStatus = "active"     // 活跃
	TemplateStatusDeprecated TemplateStatus = "deprecated" // 已弃用
	TemplateStatusDeleted    TemplateStatus = "deleted"    // 已删除
)

// TemplateVariable 模板变量定义
type TemplateVariable struct {
	Name         string   `bson:"name" json:"name" validate:"required,max=50"`
	Label        string   `bson:"label" json:"label" validate:"required,max=100"`
	Type         string   `bson:"type" json:"type" validate:"required,oneof=text textarea select number"`
	Placeholder  string   `bson:"placeholder" json:"placeholder"`
	DefaultValue string   `bson:"default_value,omitempty" json:"defaultValue,omitempty"`
	Required     bool     `bson:"required" json:"required"`
	Order        int      `bson:"order" json:"order"`
	Options      []Option `bson:"options,omitempty" json:"options,omitempty"`
}

// Option 选项定义（用于select类型变量）
type Option struct {
	Label string `bson:"label" json:"label"`
	Value string `bson:"value" json:"value"`
}

// TouchForCreate 创建时设置默认值
func (t *Template) TouchForCreate() {
	t.IdentifiedEntity.GenerateID()
	t.Timestamps.TouchForCreate()
	
	// 设置默认状态
	if t.Status == "" {
		t.Status = TemplateStatusActive
	}
	
	// 设置默认版本
	if t.Version == 0 {
		t.Version = 1
	}
}

// TouchForUpdate 更新时设置默认值
func (t *Template) TouchForUpdate() {
	t.Timestamps.Touch()
	t.Version++
}

// Validate 验证模板数据
func (t *Template) Validate() error {
	// 验证名称
	if err := base.ValidateTitle(t.Name, 200); err != nil {
		return err
	}

	// 验证模板类型
	if !isValidTemplateType(t.Type) {
		return base.ErrInvalidDocumentType // 使用现有的错误类型
	}

	// 验证模板状态
	if !isValidTemplateStatus(t.Status) {
		return base.ErrInvalidStatus // 使用现有的错误类型
	}

	return nil
}

// isValidTemplateType 验证模板类型是否有效
func isValidTemplateType(ttype TemplateType) bool {
	switch ttype {
	case TemplateTypeChapter, TemplateTypeOutline, TemplateTypeSetting:
		return true
	default:
		return false
	}
}

// isValidTemplateStatus 验证模板状态是否有效
func isValidTemplateStatus(status TemplateStatus) bool {
	switch status {
	case TemplateStatusActive, TemplateStatusDeprecated, TemplateStatusDeleted:
		return true
	default:
		return false
	}
}

// IsActive 判断模板是否活跃
func (t *Template) IsActive() bool {
	return t.Status == TemplateStatusActive
}

// IsSystemTemplate 判断是否为系统模板
func (t *Template) IsSystemTemplate() bool {
	return t.IsSystem
}

// HasProject 判断是否关联到项目
func (t *Template) HasProject() bool {
	return t.ProjectID != nil && !t.ProjectID.IsZero()
}

// IsGlobal 判断是否为全局模板（未关联到项目）
func (t *Template) IsGlobal() bool {
	return !t.HasProject()
}

// GetProjectID 获取项目ID（如果存在）
func (t *Template) GetProjectID() *primitive.ObjectID {
	return t.ProjectID
}

// SetProjectID 设置项目ID
func (t *Template) SetProjectID(projectID primitive.ObjectID) {
	t.ProjectID = &projectID
}

// ClearProjectID 清除项目ID（设为全局模板）
func (t *Template) ClearProjectID() {
	t.ProjectID = nil
}

// GetVariableByName 根据名称获取变量定义
func (t *Template) GetVariableByName(name string) *TemplateVariable {
	for i := range t.Variables {
		if t.Variables[i].Name == name {
			return &t.Variables[i]
		}
	}
	return nil
}

// AddVariable 添加变量
func (t *Template) AddVariable(variable TemplateVariable) {
	t.Variables = append(t.Variables, variable)
}

// RemoveVariable 移除变量
func (t *Template) RemoveVariable(name string) bool {
	for i, v := range t.Variables {
		if v.Name == name {
			t.Variables = append(t.Variables[:i], t.Variables[i+1:]...)
			return true
		}
	}
	return false
}

// UpdateVariable 更新变量
func (t *Template) UpdateVariable(name string, newVar TemplateVariable) bool {
	for i := range t.Variables {
		if t.Variables[i].Name == name {
			t.Variables[i] = newVar
			return true
		}
	}
	return false
}

// ValidateVariables 验证变量定义
func (t *Template) ValidateVariables() error {
	// 检查变量名重复
	varMap := make(map[string]bool)
	for _, v := range t.Variables {
		if varMap[v.Name] {
			return errors.New("变量名重复: " + v.Name)
		}
		varMap[v.Name] = true
		
		// 如果是select类型，必须有选项
		if v.Type == "select" && len(v.Options) == 0 {
			return errors.New("select类型变量必须有选项: " + v.Name)
		}
	}
	return nil
}

// Deprecated 标记模板为已弃用
func (t *Template) Deprecated() {
	t.Status = TemplateStatusDeprecated
	t.TouchForUpdate()
}

// Activate 激活模板
func (t *Template) Activate() {
	t.Status = TemplateStatusActive
	t.TouchForUpdate()
}

// SoftDelete 软删除模板
func (t *Template) SoftDelete() {
	t.Status = TemplateStatusDeleted
	t.Timestamps.SoftDelete()
}

// IsDeleted 判断是否已删除
func (t *Template) IsDeleted() bool {
	return t.Status == TemplateStatusDeleted || !t.Timestamps.DeletedAt.IsZero()
}

// GetCreatedAt 获取创建时间
func (t *Template) GetCreatedAt() time.Time {
	return t.Timestamps.CreatedAt
}

// GetUpdatedAt 获取更新时间
func (t *Template) GetUpdatedAt() time.Time {
	return t.Timestamps.UpdatedAt
}
