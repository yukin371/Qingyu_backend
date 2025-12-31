package types

// NovelType 小说写作类型
type NovelType struct{}

// NewNovelType 创建小说类型
func NewNovelType() *NovelType {
	return &NovelType{}
}

func (n *NovelType) GetTypeCode() string {
	return "novel"
}

func (n *NovelType) GetTypeName() string {
	return "小说"
}

func (n *NovelType) GetDescription() string {
	return "长篇、中篇、短篇小说等文学作品"
}

func (n *NovelType) GetDocumentTypes() []DocumentTypeDefinition {
	return []DocumentTypeDefinition{
		{
			Code:            "volume",
			Name:            "卷",
			Description:     "最高层级，用于组织章节集合",
			Level:           0,
			CanHaveChildren: true,
			Icon:            "folder",
			Color:           "#8B5CF6",
			RequiredFields:  []string{"title", "order"},
		},
		{
			Code:            "chapter",
			Name:            "章",
			Description:     "主要的故事单元",
			Level:           1,
			CanHaveChildren: true,
			Icon:            "file",
			Color:           "#3B82F6",
			RequiredFields:  []string{"title", "order"},
		},
		{
			Code:            "section",
			Name:            "节",
			Description:     "章的细分（最小叙事单元）",
			Level:           2,
			CanHaveChildren: false, // 节是最小单元，不能有子节点
			Icon:            "file-text",
			Color:           "#10B981",
			RequiredFields:  []string{"title", "order"},
		},
	}
}

func (n *NovelType) GetDefaultHierarchy() []string {
	return []string{"volume", "chapter", "section"}
}

func (n *NovelType) ValidateDocumentType(docType string) bool {
	validTypes := map[string]bool{
		"volume":  true,
		"chapter": true,
		"section": true,
	}
	return validTypes[docType]
}

func (n *NovelType) GetParentType(docType string) (string, bool) {
	parents := map[string]string{
		"chapter": "volume",
		"section": "chapter",
	}
	parent, ok := parents[docType]
	return parent, ok
}

func (n *NovelType) GetChildTypes(docType string) []string {
	children := map[string][]string{
		"volume":  {"chapter"},
		"chapter": {"section"},
		"section": {}, // 节没有子节点
	}
	return children[docType]
}

func (n *NovelType) CanHaveChildren(docType string) bool {
	return docType != "section" // 只有节不能有子节点
}

func (n *NovelType) GetMaxDepth() int {
	return 3 // volume -> chapter -> section
}

func (n *NovelType) SupportsFeature(feature string) bool {
	features := map[string]bool{
		"hierarchy":         true,
		"outline":           true,
		"timeline":          true,
		"character_map":     true,
		"location_tracking": true,
		"word_count_goal":   true,
	}
	return features[feature]
}
