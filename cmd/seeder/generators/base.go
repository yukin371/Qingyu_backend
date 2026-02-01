// Package generators 提供数据生成器
package generators

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

// BaseGenerator 基础数据生成器
type BaseGenerator struct {
	faker *gofakeit.Faker
}

// NewBaseGenerator 创建基础生成器实例
func NewBaseGenerator() *BaseGenerator {
	faker := gofakeit.New(0)
	return &BaseGenerator{faker: faker}
}

var bookTemplates = map[string][]string{
	"仙侠": {"%s修仙传", "%s%s诀", "%s之%s"},
	"都市": {"都市之%s", "重生%s", "%s之王"},
	"科幻": {"星际%s", "%s文明", "%s纪元"},
	"历史": {"大明%s", "%s风云", "%s天下"},
}

var keywords = map[string][]string{
	"仙侠": {"青羽", "太古", "混沌", "阴阳", "五行", "苍穹", "九天"},
	"都市": {"兵王", "神医", "首富", "战神", "总裁", "神豪"},
	"科幻": {"穿越", "进化", "战争", "探索", "机械", "时空"},
	"历史": {"王朝", "英雄", "江山", "谋士", "将军", "太子"},
}

// BookName 根据分类生成随机书名
func (g *BaseGenerator) BookName(category string) string {
	templates, ok := bookTemplates[category]
	if !ok {
		templates = bookTemplates["仙侠"]
	}

	words, ok := keywords[category]
	if !ok {
		words = keywords["仙侠"]
	}

	template := g.faker.RandomString(templates)
	word1 := g.faker.RandomString(words)

	if strings.Count(template, "%s") == 2 {
		word2 := g.faker.RandomString(words)
		return fmt.Sprintf(template, word1, word2)
	}

	return fmt.Sprintf(template, word1)
}

// ChapterContent 生成随机章节内容
func (g *BaseGenerator) ChapterContent(minWords, maxWords int) string {
	// 生成足够多的段落以满足最小字数要求
	paragraphs := g.faker.Number(10, 20)
	var content strings.Builder

	for i := 0; i < paragraphs; i++ {
		content.WriteString(g.faker.Paragraph(15, 25, 80, " "))
		content.WriteString("\n\n")
	}

	return content.String()
}

// Username 根据用户类型生成用户名
func (g *BaseGenerator) Username(userType string) string {
	if userType == "author" {
		return g.faker.FirstName() + g.faker.LastName()
	}
	return "user_" + g.faker.Username()
}

// Email 生成随机邮箱地址
func (g *BaseGenerator) Email() string {
	return g.faker.Email()
}

// PhoneNumber 生成随机手机号码
func (g *BaseGenerator) PhoneNumber() string {
	return g.faker.Phone()
}

// ID 生成唯一ID
func (g *BaseGenerator) ID() string {
	return uuid.New().String()
}

// ChineseName 生成中文姓名
func (g *BaseGenerator) ChineseName() string {
	surnames := []string{"王", "李", "张", "刘", "陈", "杨", "赵", "黄", "周", "吴"}
	names := []string{"伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军"}
	return surnames[rand.Intn(len(surnames))] + names[rand.Intn(len(names))]
}

// Device 生成设备类型
func (g *BaseGenerator) Device() string {
	devices := []string{"mobile", "tablet", "desktop"}
	return devices[rand.Intn(len(devices))]
}
