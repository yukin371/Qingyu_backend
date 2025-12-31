package types

// DocumentTypeDefinition 文档类型定义
type DocumentTypeDefinition struct {
	Code            string            `json:"code"`                      // 类型代码
	Name            string            `json:"name"`                      // 类型名称
	Description     string            `json:"description"`               // 描述
	Level           int               `json:"level"`                     // 层级
	CanHaveChildren bool              `json:"canHaveChildren"`           // 是否可以有子节点
	Icon            string            `json:"icon,omitempty"`            // 图标
	Color           string            `json:"color,omitempty"`           // 颜色
	RequiredFields  []string          `json:"requiredFields,omitempty"`  // 必填字段
	OptionalFields  []string          `json:"optionalFields,omitempty"`  // 可选字段
	CustomFields    map[string]string `json:"customFields,omitempty"`    // 自定义字段
}

// WritingType 写作类型接口
// 每种写作类型（小说、文章、剧本等）实现此接口以定义自己的结构
type WritingType interface {
	// GetTypeCode 返回类型代码
	GetTypeCode() string

	// GetTypeName 返回类型名称
	GetTypeName() string

	// GetDescription 返回类型描述
	GetDescription() string

	// GetDocumentTypes 返回该类型支持的文档类型列表
	GetDocumentTypes() []DocumentTypeDefinition

	// GetDefaultHierarchy 返回默认层级结构
	// 例如：小说返回 [卷, 章, 节, 场景]
	GetDefaultHierarchy() []string

	// ValidateDocumentType 验证文档类型是否属于该写作类型
	ValidateDocumentType(docType string) bool

	// GetParentType 获取给定类型的父类型
	GetParentType(docType string) (string, bool)

	// GetChildTypes 获取给定类型的子类型列表
	GetChildTypes(docType string) []string

	// CanHaveChildren 判断给定类型是否可以有子节点
	CanHaveChildren(docType string) bool

	// GetMaxDepth 获取最大层级深度
	GetMaxDepth() int

	// SupportsFeature 判断是否支持某个特性
	SupportsFeature(feature string) bool
}

// WritingTypeRegistry 写作类型注册表
type WritingTypeRegistry struct {
	types map[string]WritingType
}

// NewWritingTypeRegistry 创建写作类型注册表
func NewWritingTypeRegistry() *WritingTypeRegistry {
	registry := &WritingTypeRegistry{
		types: make(map[string]WritingType),
	}

	// 注册内置类型
	registry.Register(NewNovelType())
	registry.Register(NewArticleType())
	registry.Register(NewScriptType())

	return registry
}

// Register 注册写作类型
func (r *WritingTypeRegistry) Register(wt WritingType) {
	r.types[wt.GetTypeCode()] = wt
}

// Get 获取写作类型
func (r *WritingTypeRegistry) Get(code string) (WritingType, bool) {
	wt, ok := r.types[code]
	return wt, ok
}

// MustGet 获取写作类型，如果不存在则panic
func (r *WritingTypeRegistry) MustGet(code string) WritingType {
	wt, ok := r.types[code]
	if !ok {
		panic("writing type not found: " + code)
	}
	return wt
}

// GetAll 获取所有写作类型
func (r *WritingTypeRegistry) GetAll() []WritingType {
	result := make([]WritingType, 0, len(r.types))
	for _, wt := range r.types {
		result = append(result, wt)
	}
	return result
}

// ValidateWritingType 验证写作类型是否有效
func (r *WritingTypeRegistry) ValidateWritingType(code string) bool {
	_, ok := r.types[code]
	return ok
}

// ValidateDocumentType 验证文档类型是否有效
func (r *WritingTypeRegistry) ValidateDocumentType(writingTypeCode, docType string) bool {
	wt, ok := r.types[writingTypeCode]
	if !ok {
		return false
	}
	return wt.ValidateDocumentType(docType)
}

// GetDocumentTypeDefinition 获取文档类型定义
func (r *WritingTypeRegistry) GetDocumentTypeDefinition(writingTypeCode, docType string) (DocumentTypeDefinition, bool) {
	wt, ok := r.types[writingTypeCode]
	if !ok {
		return DocumentTypeDefinition{}, false
	}

	for _, def := range wt.GetDocumentTypes() {
		if def.Code == docType {
			return def, true
		}
	}

	return DocumentTypeDefinition{}, false
}

// 全局注册表实例
var GlobalRegistry = NewWritingTypeRegistry()
