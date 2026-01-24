package dto

// ===========================
// 书城 DTO（符合分层架构规范）
// ===========================

// BookDTO 书籍数据传输对象
// 用于：Service 层和 API 层数据传输，ID 和时间字段使用字符串类型
type BookDTO struct {
	ID        string `json:"id" validate:"required"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// 基本信息
	Title        string `json:"title" validate:"required,max=200"`
	Author       string `json:"author" validate:"required,max=100"`
	AuthorID     string `json:"authorId" validate:"required"`
	Introduction string `json:"introduction,omitempty"`
	Cover        string `json:"cover,omitempty"`

	// 分类和标签
	CategoryIDs []string `json:"categoryIds"` // 分类ID列表
	Categories  []string `json:"categories"`  // 分类名称列表
	Tags        []string `json:"tags"`        // 标签列表

	// 状态和统计
	Status      string  `json:"status" validate:"required,oneof=draft published rejected deleted"`
	Price       string  `json:"price" validate:"omitempty,numeric"` // 价格（字符串格式，避免浮点数精度问题）
	Rating      float64 `json:"rating" validate:"min=0,max=5"`       // 平均评分 0-5
	RatingCount int     `json:"ratingCount" validate:"min=0"`        // 评分人数
	ViewCount   int64   `json:"viewCount" validate:"min=0"`         // 浏览次数
	WordCount   int64   `json:"wordCount" validate:"min=0"`         // 字数
	ChapterCount int    `json:"chapterCount" validate:"min=0"`      // 章节数

	// 标记
	IsFree       bool `json:"isFree"`       // 是否免费
	IsRecommended bool `json:"isRecommended"` // 是否推荐
	IsFeatured    bool `json:"isFeatured"`    // 是否精选
	IsHot        bool `json:"isHot"`         // 是否热门

	// 发布信息
	PublishedAt  string `json:"publishedAt,omitempty"`  // 发布时间（ISO8601）
	LastUpdateAt string `json:"lastUpdateAt,omitempty"` // 最后更新时间（ISO8601）
}
