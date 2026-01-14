// Package generators 提供数据生成器
package generators

import (
	"time"

	"Qingyu_backend/cmd/seeder/models"
	"github.com/google/uuid"
)

// BookGenerator 书籍数据生成器
type BookGenerator struct {
	*BaseGenerator
}

// NewBookGenerator 创建书籍生成器
func NewBookGenerator() *BookGenerator {
	return &BookGenerator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// GenerateBook 根据分类生成单本书籍
func (g *BookGenerator) GenerateBook(category string) models.Book {
	now := time.Now()
	publishedAt := now.Add(-time.Duration(g.faker.IntRange(30, 365)) * 24 * time.Hour)

	// 生成热度值 (0-1)
	popularity := g.faker.Float64Range(0, 1)

	// 根据热度值确定评分和推荐状态
	var rating float64
	var ratingCount int64
	var isHot, isRecommended, isFeatured bool

	if popularity > 0.8 {
		// 高热度: 高评分
		rating = g.faker.Float64Range(8.5, 9.5)
		ratingCount = int64(g.faker.IntRange(1000, 50000))
		isHot = true
		isRecommended = true
		isFeatured = true
	} else if popularity > 0.5 {
		// 中热度: 中评分
		rating = g.faker.Float64Range(6.0, 8.5)
		ratingCount = int64(g.faker.IntRange(100, 1000))
		isRecommended = g.faker.Bool() // 50% 概率推荐
		isHot = false
		isFeatured = false
	} else {
		// 低热度: 低评分
		rating = g.faker.Float64Range(4.0, 6.0)
		ratingCount = int64(g.faker.IntRange(0, 100))
		isHot = false
		isRecommended = false
		isFeatured = false
	}

	// 生成章节数和字数
	chapterCount := g.faker.IntRange(10, 500)
	wordCount := int64(chapterCount * g.faker.IntRange(2000, 5000))

	// 生成浏览数 (通常远大于评分人数)
	viewCount := ratingCount * int64(g.faker.IntRange(10, 100))

	// 生成价格
	var price float64
	isFree := g.faker.Bool() // 50% 概率免费
	if isFree {
		price = 0
	} else {
		price = g.faker.Float64Range(9.9, 99.9)
	}

	// 生成状态 (连载/完结)
	statuses := []string{"ongoing", "completed"}
	status := g.faker.RandomString(statuses)

	// 使用 BaseGenerator 的 BookName 方法生成书名
	title := g.BookName(category)

	// 生成作者名
	author := g.Username("author")

	// 生成简介
	introduction := g.faker.Paragraph(2, 4, 30, " ")

	// 生成封面
	cover := "/images/covers/" + uuid.New().String() + ".jpg"

	// 生成分类和标签
	categories := []string{category}
	tags := g.generateTags(category)

	// 生成更新时间
	lastUpdateAt := publishedAt.Add(time.Duration(g.faker.IntRange(1, 30)) * 24 * time.Hour)
	if lastUpdateAt.After(now) {
		lastUpdateAt = now
	}

	return models.Book{
		ID:            uuid.New().String(),
		Title:         title,
		Author:        author,
		Introduction:  introduction,
		Cover:         cover,
		Categories:    categories,
		Tags:          tags,
		Status:        status,
		Rating:        rating,
		RatingCount:   ratingCount,
		ViewCount:     viewCount,
		WordCount:     wordCount,
		ChapterCount:  chapterCount,
		Price:         price,
		IsFree:        isFree,
		IsRecommended: isRecommended,
		IsFeatured:    isFeatured,
		IsHot:         isHot,
		PublishedAt:   publishedAt,
		LastUpdateAt:  lastUpdateAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// generateTags 根据分类生成标签
func (g *BookGenerator) generateTags(category string) []string {
	// 标签池
	tagPool := map[string][]string{
		"仙侠": {"修仙", "玄幻", "热血", "冒险", "升级", "逆天", "天才", "法宝", "丹药", "阵法"},
		"都市": {"爽文", "系统", "神豪", "异能", "医生", "兵王", "总裁", "重生", "逆袭", "都市"},
		"科幻": {"未来", "星际", "机甲", "末世", "赛博朋克", "人工智能", "时空", "进化", "战争", "探索"},
		"历史": {"穿越", "权谋", "争霸", "帝王", "名将", "智囊", "建国", "架空", "战争", "策略"},
	}

	// 获取对应分类的标签池
	tags, ok := tagPool[category]
	if !ok {
		tags = tagPool["仙侠"] // 默认使用仙侠标签
	}

	// 随机选择 3-6 个标签
	tagCount := g.faker.IntRange(3, 6)
	selectedTags := make([]string, 0, tagCount)

	// 使用 map 避免重复
	tagMap := make(map[string]bool)
	for len(selectedTags) < tagCount {
		tag := g.faker.RandomString(tags)
		if !tagMap[tag] {
			tagMap[tag] = true
			selectedTags = append(selectedTags, tag)
		}
	}

	return selectedTags
}

// GenerateBooks 批量生成书籍
func (g *BookGenerator) GenerateBooks(count int, category string) []models.Book {
	books := make([]models.Book, count)
	for i := 0; i < count; i++ {
		books[i] = g.GenerateBook(category)
	}
	return books
}
