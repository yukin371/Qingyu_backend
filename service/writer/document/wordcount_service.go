package document

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// WordCountService 字数统计服务
type WordCountService struct{}

// NewWordCountService 创建字数统计服务
func NewWordCountService() *WordCountService {
	return &WordCountService{}
}

// WordCountResult 字数统计结果
type WordCountResult struct {
	TotalCount      int    `json:"totalCount"`      // 总字数
	ChineseCount    int    `json:"chineseCount"`    // 中文字数
	EnglishCount    int    `json:"englishCount"`    // 英文单词数
	NumberCount     int    `json:"numberCount"`     // 数字个数
	ParagraphCount  int    `json:"paragraphCount"`  // 段落数
	SentenceCount   int    `json:"sentenceCount"`   // 句子数
	ReadingTime     int    `json:"readingTime"`     // 预计阅读时长（分钟）
	ReadingTimeText string `json:"readingTimeText"` // 阅读时长文本
}

// CalculateWordCount 计算字数
func (s *WordCountService) CalculateWordCount(content string) *WordCountResult {
	if content == "" {
		return &WordCountResult{}
	}

	result := &WordCountResult{}

	// 统计段落数（按换行符分割，忽略空行）
	paragraphs := strings.Split(content, "\n")
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			result.ParagraphCount++
		}
	}

	// 统计句子数（按句号、问号、感叹号分割）
	sentencePattern := regexp.MustCompile(`[。！？.!?]+`)
	sentences := sentencePattern.Split(content, -1)
	for _, s := range sentences {
		if strings.TrimSpace(s) != "" {
			result.SentenceCount++
		}
	}

	// 统计中文、英文、数字
	var englishWords []string
	var currentWord strings.Builder

	for _, r := range content {
		if isChineseChar(r) {
			// 中文字符
			result.ChineseCount++

			// 如果有累积的英文单词，先处理
			if currentWord.Len() > 0 {
				englishWords = append(englishWords, currentWord.String())
				currentWord.Reset()
			}
		} else if unicode.IsLetter(r) {
			// 英文字母
			currentWord.WriteRune(r)
		} else if unicode.IsDigit(r) {
			// 数字
			result.NumberCount++

			// 如果有累积的英文单词，先处理
			if currentWord.Len() > 0 {
				englishWords = append(englishWords, currentWord.String())
				currentWord.Reset()
			}
		} else {
			// 其他字符（空格、标点等）
			if currentWord.Len() > 0 {
				englishWords = append(englishWords, currentWord.String())
				currentWord.Reset()
			}
		}
	}

	// 处理最后的英文单词
	if currentWord.Len() > 0 {
		englishWords = append(englishWords, currentWord.String())
	}

	result.EnglishCount = len(englishWords)

	// 计算总字数（中文 + 英文单词 + 数字）
	result.TotalCount = result.ChineseCount + result.EnglishCount + result.NumberCount

	// 计算预计阅读时长
	// 假设：中文阅读速度 500字/分钟，英文阅读速度 200词/分钟
	readingMinutes := float64(result.ChineseCount)/500.0 + float64(result.EnglishCount)/200.0
	result.ReadingTime = int(readingMinutes)
	if result.ReadingTime < 1 {
		result.ReadingTime = 1
	}

	// 生成阅读时长文本
	result.ReadingTimeText = formatReadingTime(result.ReadingTime)

	return result
}

// CalculateWordCountWithMarkdown 计算字数（过滤Markdown）
func (s *WordCountService) CalculateWordCountWithMarkdown(content string) *WordCountResult {
	// 移除Markdown语法
	cleaned := removeMarkdownSyntax(content)
	return s.CalculateWordCount(cleaned)
}

// isChineseChar 判断是否为中文字符
func isChineseChar(r rune) bool {
	// 中文字符Unicode范围：
	// 常用汉字：\u4e00-\u9fa5
	// 扩展A：\u3400-\u4dbf
	// 扩展B：\u20000-\u2a6df
	return (r >= 0x4e00 && r <= 0x9fa5) ||
		(r >= 0x3400 && r <= 0x4dbf) ||
		(r >= 0x20000 && r <= 0x2a6df)
}

// formatReadingTime 格式化阅读时长
func formatReadingTime(minutes int) string {
	if minutes < 1 {
		return "少于1分钟"
	}
	if minutes < 60 {
		return fmt.Sprintf("%d分钟", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%d小时", hours)
	}
	return fmt.Sprintf("%d小时%d分钟", hours, mins)
}

// removeMarkdownSyntax 移除Markdown语法
func removeMarkdownSyntax(content string) string {
	// 移除代码块
	codeBlockPattern := regexp.MustCompile("```[\\s\\S]*?```")
	content = codeBlockPattern.ReplaceAllString(content, "")

	// 移除行内代码
	inlineCodePattern := regexp.MustCompile("`[^`]+`")
	content = inlineCodePattern.ReplaceAllString(content, "")

	// 移除链接 [text](url)
	linkPattern := regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
	content = linkPattern.ReplaceAllString(content, "$1")

	// 移除图片 ![alt](url)
	imagePattern := regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`)
	content = imagePattern.ReplaceAllString(content, "$1")

	// 移除标题标记
	headerPattern := regexp.MustCompile(`^#{1,6}\s+`)
	content = headerPattern.ReplaceAllString(content, "")

	// 移除粗体和斜体标记
	boldItalicPattern := regexp.MustCompile(`[*_]{1,3}([^*_]+)[*_]{1,3}`)
	content = boldItalicPattern.ReplaceAllString(content, "$1")

	// 移除删除线
	strikethroughPattern := regexp.MustCompile(`~~([^~]+)~~`)
	content = strikethroughPattern.ReplaceAllString(content, "$1")

	// 移除引用标记
	quotePattern := regexp.MustCompile(`^>\s+`)
	content = quotePattern.ReplaceAllString(content, "")

	// 移除列表标记
	listPattern := regexp.MustCompile(`^[\*\-\+]\s+`)
	content = listPattern.ReplaceAllString(content, "")

	// 移除有序列表标记
	orderedListPattern := regexp.MustCompile(`^\d+\.\s+`)
	content = orderedListPattern.ReplaceAllString(content, "")

	// 移除分隔线
	hrPattern := regexp.MustCompile(`^[\*\-_]{3,}$`)
	content = hrPattern.ReplaceAllString(content, "")

	return content
}
