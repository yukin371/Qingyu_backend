package types

// ScriptType 剧本写作类型
type ScriptType struct{}

// NewScriptType 创建剧本类型
func NewScriptType() *ScriptType {
	return &ScriptType{}
}

func (s *ScriptType) GetTypeCode() string {
	return "script"
}

func (s *ScriptType) GetTypeName() string {
	return "剧本"
}

func (s *ScriptType) GetDescription() string {
	return "游戏剧本、电影剧本、舞台剧、电视剧等"
}

func (s *ScriptType) GetDocumentTypes() []DocumentTypeDefinition {
	return []DocumentTypeDefinition{
		{
			Code:            "act",
			Name:            "幕",
			Description:     "剧本的最大划分单位",
			Level:           0,
			CanHaveChildren: true,
			Icon:            "film",
			Color:           "#DC2626",
			RequiredFields:  []string{"title"},
		},
		{
			Code:            "scene",
			Name:            "场",
			Description:     "场景",
			Level:           1,
			CanHaveChildren: true,
			Icon:            "map-pin",
			Color:           "#EA580C",
			RequiredFields:  []string{"title", "location"},
			CustomFields: map[string]string{
				"int_ext":   "string", // 内景/外景
				"day_night": "string", // 日/夜
				"time":      "string",
			},
		},
		{
			Code:            "beat",
			Name:            "节拍",
			Description:     "最小的叙事单元",
			Level:           2,
			CanHaveChildren: false,
			Icon:            "music",
			Color:           "#7C3AED",
			RequiredFields:  []string{"content"},
		},
	}
}

func (s *ScriptType) GetDefaultHierarchy() []string {
	return []string{"act", "scene", "beat"}
}

func (s *ScriptType) ValidateDocumentType(docType string) bool {
	validTypes := map[string]bool{
		"act":   true,
		"scene": true,
		"beat":  true,
	}
	return validTypes[docType]
}

func (s *ScriptType) GetParentType(docType string) (string, bool) {
	parents := map[string]string{
		"scene": "act",
		"beat":  "scene",
	}
	parent, ok := parents[docType]
	return parent, ok
}

func (s *ScriptType) GetChildTypes(docType string) []string {
	children := map[string][]string{
		"act":   {"scene"},
		"scene": {"beat"},
		"beat":  {},
	}
	return children[docType]
}

func (s *ScriptType) CanHaveChildren(docType string) bool {
	return docType != "beat"
}

func (s *ScriptType) GetMaxDepth() int {
	return 3
}

func (s *ScriptType) SupportsFeature(feature string) bool {
	features := map[string]bool{
		"hierarchy":         true,
		"outline":           true,
		"timeline":          true,
		"character_map":     true,
		"location_tracking": true,
		"script_format":     true,
	}
	return features[feature]
}
