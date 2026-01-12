package reader

import "time"

// ReaderFont 阅读器字体配置
type ReaderFont struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Name        string `bson:"name" json:"name"`                // 字体名称（唯一标识）
	DisplayName string `bson:"display_name" json:"displayName"` // 显示名称
	FontFamily  string `bson:"font_family" json:"fontFamily"`   // CSS font-family 值
	Description string `bson:"description" json:"description"`  // 字体描述
	Category    string `bson:"category" json:"category"`        // 字体分类：serif/sans-serif/monospace

	// 字体文件信息（如果支持自定义字体上传）
	FontURL    string `bson:"font_url,omitempty" json:"fontUrl"`       // 字体文件URL
	FontFormat string `bson:"font_format,omitempty" json:"fontFormat"` // 字体格式：woff/woff2/ttf

	// 字体预览
	PreviewText string `bson:"preview_text" json:"previewText"`         // 预览文本
	PreviewURL  string `bson:"preview_url,omitempty" json:"previewUrl"` // 预览图片URL

	// 属性
	IsBuiltIn   bool  `bson:"is_built_in" json:"isBuiltIn"`    // 是否内置字体
	IsActive    bool  `bson:"is_active" json:"isActive"`       // 是否激活可用
	SupportSize []int `bson:"support_size" json:"supportSize"` // 支持的字号列表

	// 统计
	UseCount  int64     `bson:"use_count" json:"useCount"` // 使用次数
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// BuiltInFonts 内置字体列表
var BuiltInFonts = []*ReaderFont{
	{
		Name:        "system-serif",
		DisplayName: "宋体/衬线",
		FontFamily:  "SimSun, 'Songti SC', 'Noto Serif SC', serif",
		Description: "经典衬线字体，适合正文阅读",
		Category:    "serif",
		IsBuiltIn:   true,
		IsActive:    true,
		SupportSize: []int{12, 14, 16, 18, 20, 22, 24, 28, 32},
		PreviewText: "这是一段宋体预览文字 The quick brown fox jumps over the lazy dog.",
		UseCount:    0,
	},
	{
		Name:        "system-sans",
		DisplayName: "黑体/无衬线",
		FontFamily:  "'Microsoft YaHei', 'PingFang SC', 'Noto Sans SC', sans-serif",
		Description: "现代无衬线字体，清晰易读",
		Category:    "sans-serif",
		IsBuiltIn:   true,
		IsActive:    true,
		SupportSize: []int{12, 14, 16, 18, 20, 22, 24, 28, 32},
		PreviewText: "这是一段黑体预览文字 The quick brown fox jumps over the lazy dog.",
		UseCount:    0,
	},
	{
		Name:        "kai",
		DisplayName: "楷体",
		FontFamily:  "KaiTi, 'Kaiti SC', 'STKaiti', serif",
		Description: "传统楷体，书法风格",
		Category:    "serif",
		IsBuiltIn:   true,
		IsActive:    true,
		SupportSize: []int{14, 16, 18, 20, 22, 24, 28, 32},
		PreviewText: "这是一段楷体预览文字 The quick brown fox jumps over the lazy dog.",
		UseCount:    0,
	},
	{
		Name:        "fangsong",
		DisplayName: "仿宋",
		FontFamily:  "FangSong, 'STFangsong', serif",
		Description: "仿宋字体，正式文档常用",
		Category:    "serif",
		IsBuiltIn:   true,
		IsActive:    true,
		SupportSize: []int{14, 16, 18, 20, 22, 24},
		PreviewText: "这是一段仿宋预览文字 The quick brown fox jumps over the lazy dog.",
		UseCount:    0,
	},
	{
		Name:        "monospace",
		DisplayName: "等宽字体",
		FontFamily:  "'Courier New', Consolas, monospace",
		Description: "等宽字体，适合代码阅读",
		Category:    "monospace",
		IsBuiltIn:   true,
		IsActive:    true,
		SupportSize: []int{12, 14, 16, 18, 20},
		PreviewText: "这是一段等宽字体预览文字 The quick brown fox jumps over the lazy dog.",
		UseCount:    0,
	},
}

// FontPreference 字体偏好设置
type FontPreference struct {
	UserID        string  `bson:"user_id" json:"userId"`               // 用户ID
	FontName      string  `bson:"font_name" json:"fontName"`           // 字体名称
	FontSize      int     `bson:"font_size" json:"fontSize"`           // 字号
	LineHeight    float64 `bson:"line_height" json:"lineHeight"`       // 行高
	LetterSpacing float64 `bson:"letter_spacing" json:"letterSpacing"` // 字间距
}

// CreateCustomFontRequest 创建自定义字体请求
type CreateCustomFontRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	DisplayName string `json:"displayName" validate:"required,min=1,max=50"`
	FontFamily  string `json:"fontFamily" validate:"required"`
	Description string `json:"description" validate:"max=200"`
	Category    string `json:"category" validate:"required,oneof=serif sans-serif monospace"`
	FontURL     string `json:"fontUrl"`
	PreviewText string `json:"previewText"`
}

// UpdateFontRequest 更新字体请求
type UpdateFontRequest struct {
	DisplayName *string `json:"displayName" validate:"omitempty,min=1,max=50"`
	Description *string `json:"description" validate:"omitempty,max=200"`
	IsActive    *bool   `json:"isActive"`
}
