package types

// ArticleType 文章写作类型
type ArticleType struct{}

// NewArticleType 创建文章类型
func NewArticleType() *ArticleType {
	return &ArticleType{}
}

func (a *ArticleType) GetTypeCode() string {
	return "article"
}

func (a *ArticleType) GetTypeName() string {
	return "文章"
}

func (a *ArticleType) GetDescription() string {
	return "新闻报道、技术文章、博客文章等"
}

func (a *ArticleType) GetDocumentTypes() []DocumentTypeDefinition {
	return []DocumentTypeDefinition{
		{
			Code:            "article",
			Name:            "文章",
			Description:     "文章主体",
			Level:           0,
			CanHaveChildren: true,
			Icon:            "file-text",
			Color:           "#059669",
			RequiredFields:  []string{"title"},
		},
		{
			Code:            "section",
			Name:            "小节",
			Description:     "文章的小节",
			Level:           1,
			CanHaveChildren: true,
			Icon:            "hash",
			Color:           "#0D9488",
			RequiredFields:  []string{"title"},
		},
		{
			Code:            "paragraph",
			Name:            "段落",
			Description:     "文本段落",
			Level:           2,
			CanHaveChildren: false,
			Icon:            "align-left",
			Color:           "#6366F1",
		},
	}
}

func (a *ArticleType) GetDefaultHierarchy() []string {
	return []string{"article", "section", "paragraph"}
}

func (a *ArticleType) ValidateDocumentType(docType string) bool {
	validTypes := map[string]bool{
		"article":   true,
		"section":   true,
		"paragraph": true,
	}
	return validTypes[docType]
}

func (a *ArticleType) GetParentType(docType string) (string, bool) {
	parents := map[string]string{
		"section":   "article",
		"paragraph": "section",
	}
	parent, ok := parents[docType]
	return parent, ok
}

func (a *ArticleType) GetChildTypes(docType string) []string {
	children := map[string][]string{
		"article":   {"section"},
		"section":   {"section", "paragraph"},
		"paragraph": {},
	}
	return children[docType]
}

func (a *ArticleType) CanHaveChildren(docType string) bool {
	return docType != "paragraph"
}

func (a *ArticleType) GetMaxDepth() int {
	return 3
}

func (a *ArticleType) SupportsFeature(feature string) bool {
	features := map[string]bool{
		"hierarchy":       true,
		"outline":         false,
		"timeline":        false,
		"character_map":   false,
		"word_count_goal": true,
		"tags":            true,
	}
	return features[feature]
}
