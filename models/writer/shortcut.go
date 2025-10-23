package writer

import "time"

// ShortcutConfig 快捷键配置
type ShortcutConfig struct {
	ID        string              `bson:"_id,omitempty" json:"id"`
	UserID    string              `bson:"userId" json:"userId"`       // 用户ID
	Shortcuts map[string]Shortcut `bson:"shortcuts" json:"shortcuts"` // 快捷键映射
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt"`
}

// Shortcut 单个快捷键
type Shortcut struct {
	Action      string `bson:"action" json:"action"`           // 操作名称
	Key         string `bson:"key" json:"key"`                 // 按键组合 (e.g., "Ctrl+S")
	Description string `bson:"description" json:"description"` // 描述
	Category    string `bson:"category" json:"category"`       // 分类
	IsCustom    bool   `bson:"isCustom" json:"isCustom"`       // 是否自定义
}

// DefaultShortcuts 默认快捷键配置
var DefaultShortcuts = map[string]Shortcut{
	// 文件操作
	"save": {
		Action:      "save",
		Key:         "Ctrl+S",
		Description: "保存文档",
		Category:    "文件",
		IsCustom:    false,
	},
	"save_all": {
		Action:      "save_all",
		Key:         "Ctrl+Shift+S",
		Description: "保存所有",
		Category:    "文件",
		IsCustom:    false,
	},
	"new_document": {
		Action:      "new_document",
		Key:         "Ctrl+N",
		Description: "新建文档",
		Category:    "文件",
		IsCustom:    false,
	},
	"close_document": {
		Action:      "close_document",
		Key:         "Ctrl+W",
		Description: "关闭文档",
		Category:    "文件",
		IsCustom:    false,
	},

	// 编辑操作
	"undo": {
		Action:      "undo",
		Key:         "Ctrl+Z",
		Description: "撤销",
		Category:    "编辑",
		IsCustom:    false,
	},
	"redo": {
		Action:      "redo",
		Key:         "Ctrl+Y",
		Description: "重做",
		Category:    "编辑",
		IsCustom:    false,
	},
	"cut": {
		Action:      "cut",
		Key:         "Ctrl+X",
		Description: "剪切",
		Category:    "编辑",
		IsCustom:    false,
	},
	"copy": {
		Action:      "copy",
		Key:         "Ctrl+C",
		Description: "复制",
		Category:    "编辑",
		IsCustom:    false,
	},
	"paste": {
		Action:      "paste",
		Key:         "Ctrl+V",
		Description: "粘贴",
		Category:    "编辑",
		IsCustom:    false,
	},
	"select_all": {
		Action:      "select_all",
		Key:         "Ctrl+A",
		Description: "全选",
		Category:    "编辑",
		IsCustom:    false,
	},
	"find": {
		Action:      "find",
		Key:         "Ctrl+F",
		Description: "查找",
		Category:    "编辑",
		IsCustom:    false,
	},
	"replace": {
		Action:      "replace",
		Key:         "Ctrl+H",
		Description: "替换",
		Category:    "编辑",
		IsCustom:    false,
	},

	// 格式化
	"bold": {
		Action:      "bold",
		Key:         "Ctrl+B",
		Description: "加粗",
		Category:    "格式",
		IsCustom:    false,
	},
	"italic": {
		Action:      "italic",
		Key:         "Ctrl+I",
		Description: "斜体",
		Category:    "格式",
		IsCustom:    false,
	},
	"underline": {
		Action:      "underline",
		Key:         "Ctrl+U",
		Description: "下划线",
		Category:    "格式",
		IsCustom:    false,
	},
	"strikethrough": {
		Action:      "strikethrough",
		Key:         "Ctrl+Shift+X",
		Description: "删除线",
		Category:    "格式",
		IsCustom:    false,
	},
	"heading1": {
		Action:      "heading1",
		Key:         "Ctrl+Alt+1",
		Description: "一级标题",
		Category:    "格式",
		IsCustom:    false,
	},
	"heading2": {
		Action:      "heading2",
		Key:         "Ctrl+Alt+2",
		Description: "二级标题",
		Category:    "格式",
		IsCustom:    false,
	},
	"heading3": {
		Action:      "heading3",
		Key:         "Ctrl+Alt+3",
		Description: "三级标题",
		Category:    "格式",
		IsCustom:    false,
	},

	// 段落
	"indent": {
		Action:      "indent",
		Key:         "Tab",
		Description: "缩进",
		Category:    "段落",
		IsCustom:    false,
	},
	"outdent": {
		Action:      "outdent",
		Key:         "Shift+Tab",
		Description: "取消缩进",
		Category:    "段落",
		IsCustom:    false,
	},
	"bullet_list": {
		Action:      "bullet_list",
		Key:         "Ctrl+Shift+8",
		Description: "无序列表",
		Category:    "段落",
		IsCustom:    false,
	},
	"ordered_list": {
		Action:      "ordered_list",
		Key:         "Ctrl+Shift+7",
		Description: "有序列表",
		Category:    "段落",
		IsCustom:    false,
	},
	"blockquote": {
		Action:      "blockquote",
		Key:         "Ctrl+Shift+9",
		Description: "引用",
		Category:    "段落",
		IsCustom:    false,
	},

	// 插入
	"insert_link": {
		Action:      "insert_link",
		Key:         "Ctrl+K",
		Description: "插入链接",
		Category:    "插入",
		IsCustom:    false,
	},
	"insert_image": {
		Action:      "insert_image",
		Key:         "Ctrl+Shift+I",
		Description: "插入图片",
		Category:    "插入",
		IsCustom:    false,
	},
	"insert_table": {
		Action:      "insert_table",
		Key:         "Ctrl+Shift+T",
		Description: "插入表格",
		Category:    "插入",
		IsCustom:    false,
	},
	"insert_code": {
		Action:      "insert_code",
		Key:         "Ctrl+Shift+C",
		Description: "插入代码",
		Category:    "插入",
		IsCustom:    false,
	},

	// 视图
	"toggle_sidebar": {
		Action:      "toggle_sidebar",
		Key:         "Ctrl+\\",
		Description: "切换侧边栏",
		Category:    "视图",
		IsCustom:    false,
	},
	"toggle_preview": {
		Action:      "toggle_preview",
		Key:         "Ctrl+Shift+P",
		Description: "切换预览",
		Category:    "视图",
		IsCustom:    false,
	},
	"toggle_fullscreen": {
		Action:      "toggle_fullscreen",
		Key:         "F11",
		Description: "全屏模式",
		Category:    "视图",
		IsCustom:    false,
	},
	"zoom_in": {
		Action:      "zoom_in",
		Key:         "Ctrl+=",
		Description: "放大",
		Category:    "视图",
		IsCustom:    false,
	},
	"zoom_out": {
		Action:      "zoom_out",
		Key:         "Ctrl+-",
		Description: "缩小",
		Category:    "视图",
		IsCustom:    false,
	},
	"reset_zoom": {
		Action:      "reset_zoom",
		Key:         "Ctrl+0",
		Description: "重置缩放",
		Category:    "视图",
		IsCustom:    false,
	},
}

// GetDefaultShortcuts 获取默认快捷键配置
func GetDefaultShortcuts() map[string]Shortcut {
	// 返回副本，避免被修改
	shortcuts := make(map[string]Shortcut)
	for k, v := range DefaultShortcuts {
		shortcuts[k] = v
	}
	return shortcuts
}

// ShortcutCategory 快捷键分类
type ShortcutCategory struct {
	Name      string     `json:"name"`
	Shortcuts []Shortcut `json:"shortcuts"`
}

// GetShortcutsByCategory 按分类获取快捷键
func GetShortcutsByCategory(shortcuts map[string]Shortcut) []ShortcutCategory {
	categories := make(map[string][]Shortcut)

	for _, shortcut := range shortcuts {
		categories[shortcut.Category] = append(categories[shortcut.Category], shortcut)
	}

	result := make([]ShortcutCategory, 0, len(categories))
	for name, items := range categories {
		result = append(result, ShortcutCategory{
			Name:      name,
			Shortcuts: items,
		})
	}

	return result
}
