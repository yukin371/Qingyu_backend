package service

import (
	"testing"

	"Qingyu_backend/service/document"

	"github.com/stretchr/testify/assert"
)

func TestWordCountService_CalculateWordCount(t *testing.T) {
	service := document.NewWordCountService()

	tests := []struct {
		name     string
		content  string
		expected *document.WordCountResult
	}{
		{
			name:    "空内容",
			content: "",
			expected: &document.WordCountResult{
				TotalCount:      0,
				ChineseCount:    0,
				EnglishCount:    0,
				NumberCount:     0,
				ParagraphCount:  0,
				SentenceCount:   0,
				ReadingTime:     0,
				ReadingTimeText: "",
			},
		},
		{
			name:    "纯中文",
			content: "这是一个测试。",
			expected: &document.WordCountResult{
				TotalCount:     6,
				ChineseCount:   6,
				EnglishCount:   0,
				NumberCount:    0,
				ParagraphCount: 1,
				SentenceCount:  1,
				ReadingTime:    1,
			},
		},
		{
			name:    "纯英文",
			content: "This is a test.",
			expected: &document.WordCountResult{
				TotalCount:     4,
				ChineseCount:   0,
				EnglishCount:   4,
				NumberCount:    0,
				ParagraphCount: 1,
				SentenceCount:  1,
				ReadingTime:    1,
			},
		},
		{
			name:    "中英文混合",
			content: "Hello世界！This is 测试123。",
			expected: &document.WordCountResult{
				TotalCount:     10, // Hello(1) + 世界(2) + This(1) + is(1) + 测试(2) + 123(3)
				ChineseCount:   4,  // 世界 + 测试
				EnglishCount:   3,  // Hello + This + is
				NumberCount:    3,  // 123
				ParagraphCount: 1,
				SentenceCount:  1,
			},
		},
		{
			name: "多段落",
			content: `第一段内容。

第二段内容。

第三段内容。`,
			expected: &document.WordCountResult{
				TotalCount:     15,
				ChineseCount:   15,
				EnglishCount:   0,
				NumberCount:    0,
				ParagraphCount: 3,
				SentenceCount:  3,
			},
		},
		{
			name:    "包含数字",
			content: "2024年有365天",
			expected: &document.WordCountResult{
				TotalCount:     10, // 2024(4) + 年(1) + 有(1) + 365(3) + 天(1)
				ChineseCount:   3,
				EnglishCount:   0,
				NumberCount:    7, // 2024 + 365
				ParagraphCount: 1,
				SentenceCount:  0, // 没有句号
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateWordCount(tt.content)

			assert.Equal(t, tt.expected.TotalCount, result.TotalCount, "总字数不匹配")
			assert.Equal(t, tt.expected.ChineseCount, result.ChineseCount, "中文字数不匹配")
			assert.Equal(t, tt.expected.EnglishCount, result.EnglishCount, "英文单词数不匹配")
			assert.Equal(t, tt.expected.NumberCount, result.NumberCount, "数字个数不匹配")
			assert.Equal(t, tt.expected.ParagraphCount, result.ParagraphCount, "段落数不匹配")
			assert.Equal(t, tt.expected.SentenceCount, result.SentenceCount, "句子数不匹配")

			// 验证阅读时长
			if tt.expected.ReadingTime > 0 {
				assert.GreaterOrEqual(t, result.ReadingTime, tt.expected.ReadingTime, "阅读时长应该至少为预期值")
			}

			// 验证阅读时长文本
			if result.ReadingTime > 0 {
				assert.NotEmpty(t, result.ReadingTimeText, "阅读时长文本不应为空")
			}
		})
	}
}

func TestWordCountService_CalculateWordCountWithMarkdown(t *testing.T) {
	service := document.NewWordCountService()

	tests := []struct {
		name            string
		content         string
		expectedMinimum int // 最小字数（Markdown过滤后应该减少）
	}{
		{
			name: "标题",
			content: `# 一级标题
## 二级标题
### 三级标题`,
			expectedMinimum: 12, // 一级标题 + 二级标题 + 三级标题
		},
		{
			name:            "粗体斜体",
			content:         "**粗体内容** *斜体内容* ~~删除线~~",
			expectedMinimum: 12, // 粗体内容 + 斜体内容 + 删除线
		},
		{
			name:            "链接",
			content:         "[链接文本](https://example.com)",
			expectedMinimum: 4, // 链接文本
		},
		{
			name:            "图片",
			content:         "![图片描述](https://example.com/image.jpg)",
			expectedMinimum: 4, // 图片描述
		},
		{
			name:            "代码块",
			content:         "```go\nfunc main() {\n}\n```\n这是正文内容",
			expectedMinimum: 6, // 这是正文内容
		},
		{
			name:            "行内代码",
			content:         "这是`代码`内容",
			expectedMinimum: 4, // 这是内容（代码被过滤）
		},
		{
			name: "列表",
			content: `- 列表项1
- 列表项2
* 列表项3`,
			expectedMinimum: 9, // 列表项1 + 列表项2 + 列表项3
		},
		{
			name: "有序列表",
			content: `1. 第一项
2. 第二项
3. 第三项`,
			expectedMinimum: 9, // 第一项 + 第二项 + 第三项
		},
		{
			name: "引用",
			content: `> 这是引用内容
> 第二行引用`,
			expectedMinimum: 12, // 这是引用内容 + 第二行引用
		},
		{
			name:            "综合测试",
			content:         "# 文档标题\n\n这是**正文内容**，包含[链接](url)和图片![alt](url)。\n\n```go\ncode here\n```\n\n- 列表项1\n- 列表项2",
			expectedMinimum: 20, // 预估最小字数
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateWordCountWithMarkdown(tt.content)

			// 验证过滤后的字数
			assert.GreaterOrEqual(t, result.TotalCount, tt.expectedMinimum,
				"过滤Markdown后的字数应该大于等于预期最小值")

			// 验证过滤确实生效（字数应该比原始内容少或相等）
			originalResult := service.CalculateWordCount(tt.content)
			assert.LessOrEqual(t, result.TotalCount, originalResult.TotalCount,
				"过滤后的字数应该小于等于原始字数")
		})
	}
}

func TestWordCountService_ReadingTime(t *testing.T) {
	service := document.NewWordCountService()

	tests := []struct {
		name                string
		content             string
		expectedMinReadTime int
		expectedMaxReadTime int
	}{
		{
			name:                "短文本（少于1分钟）",
			content:             "这是一个很短的测试文本",
			expectedMinReadTime: 1,
			expectedMaxReadTime: 1,
		},
		{
			name:                "500字中文（约1分钟）",
			content:             generateChineseText(500),
			expectedMinReadTime: 1,
			expectedMaxReadTime: 1,
		},
		{
			name:                "1000字中文（约2分钟）",
			content:             generateChineseText(1000),
			expectedMinReadTime: 2,
			expectedMaxReadTime: 2,
		},
		{
			name:                "200个英文单词（约1分钟）",
			content:             generateEnglishText(200),
			expectedMinReadTime: 1,
			expectedMaxReadTime: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateWordCount(tt.content)

			assert.GreaterOrEqual(t, result.ReadingTime, tt.expectedMinReadTime,
				"阅读时长应该大于等于预期最小值")
			assert.LessOrEqual(t, result.ReadingTime, tt.expectedMaxReadTime,
				"阅读时长应该小于等于预期最大值")

			// 验证阅读时长文本格式
			assert.NotEmpty(t, result.ReadingTimeText, "阅读时长文本不应为空")
			assert.Contains(t, result.ReadingTimeText, "分钟", "阅读时长文本应该包含'分钟'")
		})
	}
}

func TestWordCountService_LargeDocument(t *testing.T) {
	service := document.NewWordCountService()

	// 测试大文档（10000字）
	largeContent := generateChineseText(10000)

	result := service.CalculateWordCount(largeContent)

	// 验证统计准确性
	assert.Equal(t, 10000, result.ChineseCount, "大文档中文字数应该准确")
	assert.Equal(t, 10000, result.TotalCount, "大文档总字数应该准确")

	// 验证阅读时长合理（10000字约20分钟）
	assert.GreaterOrEqual(t, result.ReadingTime, 20, "大文档阅读时长应该合理")
	assert.LessOrEqual(t, result.ReadingTime, 20, "大文档阅读时长应该准确")
}

// 辅助函数：生成指定字数的中文文本
func generateChineseText(count int) string {
	text := ""
	baseText := "这是一个测试文本用于字数统计验证功能"
	for len([]rune(text)) < count {
		text += baseText
	}
	// 截取到指定字数
	runes := []rune(text)
	if len(runes) > count {
		runes = runes[:count]
	}
	return string(runes)
}

// 辅助函数：生成指定单词数的英文文本
func generateEnglishText(wordCount int) string {
	text := ""
	words := []string{"This", "is", "a", "test", "document", "for", "word", "count", "validation"}
	for i := 0; i < wordCount; i++ {
		text += words[i%len(words)] + " "
	}
	return text
}

func BenchmarkWordCountService_CalculateWordCount(b *testing.B) {
	service := document.NewWordCountService()
	content := generateChineseText(1000) + " " + generateEnglishText(200)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CalculateWordCount(content)
	}
}

func BenchmarkWordCountService_CalculateWordCountWithMarkdown(b *testing.B) {
	service := document.NewWordCountService()
	content := "# 标题\n\n这是**粗体**和*斜体*内容。\n\n```go\nfunc main() {\n}\n```\n\n- 列表项1\n- 列表项2\n\n[链接](url)和![图片](url)。"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CalculateWordCountWithMarkdown(content)
	}
}
