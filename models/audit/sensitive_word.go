package audit

import "time"

// SensitiveWord 敏感词模型
type SensitiveWord struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Word        string    `bson:"word" json:"word" validate:"required"`               // 敏感词
	Category    string    `bson:"category" json:"category" validate:"required"`       // 分类（政治、色情、暴力等）
	Level       int       `bson:"level" json:"level" validate:"required,min=1,max=5"` // 严重等级 1-5
	Replacement string    `bson:"replacement" json:"replacement"`                     // 替换词（如：***）
	IsEnabled   bool      `bson:"isEnabled" json:"isEnabled"`                         // 是否启用
	Source      string    `bson:"source" json:"source"`                               // 来源（系统/用户添加）
	Description string    `bson:"description" json:"description"`                     // 描述说明
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
	CreatedBy   string    `bson:"createdBy" json:"createdBy"` // 创建人
}

// SensitiveWordCategory 敏感词分类常量
const (
	CategoryPolitics = "politics" // 政治敏感
	CategoryPorn     = "porn"     // 色情内容
	CategoryViolence = "violence" // 暴力内容
	CategoryGambling = "gambling" // 赌博相关
	CategoryDrugs    = "drugs"    // 毒品相关
	CategoryCult     = "cult"     // 邪教相关
	CategoryInsult   = "insult"   // 侮辱谩骂
	CategoryAd       = "ad"       // 广告推广
	CategoryOther    = "other"    // 其他
)

// SensitiveWordLevel 敏感词等级
const (
	LevelLow      = 1 // 低风险（警告）
	LevelMedium   = 2 // 中风险（需要复核）
	LevelHigh     = 3 // 高风险（自动拒绝）
	LevelCritical = 4 // 严重（自动拒绝+封号警告）
	LevelBanned   = 5 // 禁用（自动封号）
)

// GetCategoryName 获取分类中文名
func GetCategoryName(category string) string {
	names := map[string]string{
		CategoryPolitics: "政治敏感",
		CategoryPorn:     "色情内容",
		CategoryViolence: "暴力内容",
		CategoryGambling: "赌博相关",
		CategoryDrugs:    "毒品相关",
		CategoryCult:     "邪教相关",
		CategoryInsult:   "侮辱谩骂",
		CategoryAd:       "广告推广",
		CategoryOther:    "其他",
	}
	if name, ok := names[category]; ok {
		return name
	}
	return "未知分类"
}

// GetLevelName 获取等级中文名
func GetLevelName(level int) string {
	names := map[int]string{
		LevelLow:      "低风险",
		LevelMedium:   "中风险",
		LevelHigh:     "高风险",
		LevelCritical: "严重",
		LevelBanned:   "禁用",
	}
	if name, ok := names[level]; ok {
		return name
	}
	return "未知等级"
}

// IsHighRisk 是否高风险
func (s *SensitiveWord) IsHighRisk() bool {
	return s.Level >= LevelHigh
}

// ShouldBan 是否应该封号
func (s *SensitiveWord) ShouldBan() bool {
	return s.Level >= LevelBanned
}
