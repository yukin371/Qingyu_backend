package reader

import "time"

// ReaderTheme 阅读器主题
type ReaderTheme struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Name        string `bson:"name" json:"name"`                      // 主题名称
	DisplayName string `bson:"display_name" json:"displayName"`       // 显示名称
	Description string `bson:"description" json:"description"`        // 主题描述
	IsBuiltIn   bool   `bson:"is_built_in" json:"isBuiltIn"`          // 是否内置主题
	IsPublic    bool   `bson:"is_public" json:"isPublic"`             // 是否公开（其他用户可见）
	CreatorID   string `bson:"creator_id,omitempty" json:"creatorId"` // 创建者ID（自定义主题）

	// 颜色设置
	Colors ThemeColors `bson:"colors" json:"colors"` // 主题颜色配置

	// 状态
	IsActive bool `bson:"is_active" json:"isActive"` // 是否激活

	// 使用统计
	UseCount  int64     `bson:"use_count" json:"useCount"` // 使用次数
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// ThemeColors 主题颜色配置
type ThemeColors struct {
	// 背景颜色
	Background          string `bson:"background" json:"background"`                    // 主背景色
	SecondaryBackground string `bson:"secondary_background" json:"secondaryBackground"` // 次背景色

	// 文字颜色
	TextPrimary   string `bson:"text_primary" json:"textPrimary"`     // 主要文字颜色
	TextSecondary string `bson:"text_secondary" json:"textSecondary"` // 次要文字颜色
	TextDisabled  string `bson:"text_disabled" json:"textDisabled"`   // 禁用文字颜色
	LinkColor     string `bson:"link_color" json:"linkColor"`         // 链接颜色

	// 强调色
	AccentColor string `bson:"accent_color" json:"accentColor"` // 强调色
	AccentHover string `bson:"accent_hover" json:"accentHover"` // 强调色悬停

	// 边框和分隔线
	BorderColor  string `bson:"border_color" json:"borderColor"`   // 边框颜色
	DividerColor string `bson:"divider_color" json:"dividerColor"` // 分隔线颜色

	// 特殊元素
	HighlightColor  string `bson:"highlight_color" json:"highlightColor"`   // 高亮颜色
	BookmarkColor   string `bson:"bookmark_color" json:"bookmarkColor"`     // 书签颜色
	AnnotationColor string `bson:"annotation_color" json:"annotationColor"` // 标注颜色

	// 阴影
	ShadowColor string `bson:"shadow_color" json:"shadowColor"` // 阴影颜色
}

// BuiltInThemes 内置主题列表
var BuiltInThemes = []*ReaderTheme{
	{
		Name:        "light",
		DisplayName: "明亮模式",
		Description: "默认明亮主题，适合白天阅读",
		IsBuiltIn:   true,
		IsPublic:    true,
		Colors: ThemeColors{
			Background:          "#FFFFFF",
			SecondaryBackground: "#F5F5F5",
			TextPrimary:         "#212121",
			TextSecondary:       "#757575",
			TextDisabled:        "#BDBDBD",
			LinkColor:           "#1976D2",
			AccentColor:         "#1976D2",
			AccentHover:         "#1565C0",
			BorderColor:         "#E0E0E0",
			DividerColor:        "#EEEEEE",
			HighlightColor:      "#FFEB3B",
			BookmarkColor:       "#FF9800",
			AnnotationColor:     "#4CAF50",
			ShadowColor:         "rgba(0, 0, 0, 0.1)",
		},
		IsActive: true,
		UseCount: 0,
	},
	{
		Name:        "dark",
		DisplayName: "暗黑模式",
		Description: "护眼暗色主题，适合夜间阅读",
		IsBuiltIn:   true,
		IsPublic:    true,
		Colors: ThemeColors{
			Background:          "#121212",
			SecondaryBackground: "#1E1E1E",
			TextPrimary:         "#FFFFFF",
			TextSecondary:       "#B0B0B0",
			TextDisabled:        "#666666",
			LinkColor:           "#64B5F6",
			AccentColor:         "#64B5F6",
			AccentHover:         "#42A5F5",
			BorderColor:         "#333333",
			DividerColor:        "#2C2C2C",
			HighlightColor:      "#FFD54F",
			BookmarkColor:       "#FFB74D",
			AnnotationColor:     "#81C784",
			ShadowColor:         "rgba(0, 0, 0, 0.3)",
		},
		IsActive: false,
		UseCount: 0,
	},
	{
		Name:        "sepia",
		DisplayName: "羊皮纸模式",
		Description: "复古羊皮纸主题，温暖舒适",
		IsBuiltIn:   true,
		IsPublic:    true,
		Colors: ThemeColors{
			Background:          "#F4ECD8",
			SecondaryBackground: "#E8DFC9",
			TextPrimary:         "#5C4B37",
			TextSecondary:       "#8B7355",
			TextDisabled:        "#B8A88A",
			LinkColor:           "#8B4513",
			AccentColor:         "#8B4513",
			AccentHover:         "#A0522D",
			BorderColor:         "#D4C5A8",
			DividerColor:        "#DDD0B8",
			HighlightColor:      "#FFE082",
			BookmarkColor:       "#FFB74D",
			AnnotationColor:     "#A1887F",
			ShadowColor:         "rgba(92, 75, 55, 0.1)",
		},
		IsActive: false,
		UseCount: 0,
	},
	{
		Name:        "eye-care",
		DisplayName: "护眼模式",
		Description: "绿色护眼主题，减轻视觉疲劳",
		IsBuiltIn:   true,
		IsPublic:    true,
		Colors: ThemeColors{
			Background:          "#C7EDCC",
			SecondaryBackground: "#B8DFC0",
			TextPrimary:         "#2E4F33",
			TextSecondary:       "#4A6B4F",
			TextDisabled:        "#7A8B7A",
			LinkColor:           "#2E7D32",
			AccentColor:         "#2E7D32",
			AccentHover:         "#1B5E20",
			BorderColor:         "#A5D4AA",
			DividerColor:        "#B8DFC0",
			HighlightColor:      "#C5E1A5",
			BookmarkColor:       "#81C784",
			AnnotationColor:     "#4CAF50",
			ShadowColor:         "rgba(46, 79, 51, 0.1)",
		},
		IsActive: false,
		UseCount: 0,
	},
}

// CreateCustomThemeRequest 创建自定义主题请求
type CreateCustomThemeRequest struct {
	Name        string      `json:"name" validate:"required,min=1,max=50"`
	DisplayName string      `json:"displayName" validate:"required,min=1,max=50"`
	Description string      `json:"description" validate:"max=200"`
	IsPublic    bool        `json:"isPublic"`
	Colors      ThemeColors `json:"colors" validate:"required"`
}

// UpdateThemeRequest 更新主题请求
type UpdateThemeRequest struct {
	DisplayName *string      `json:"displayName" validate:"omitempty,min=1,max=50"`
	Description *string      `json:"description" validate:"omitempty,max=200"`
	IsPublic    *bool        `json:"isPublic"`
	Colors      *ThemeColors `json:"colors"`
}
