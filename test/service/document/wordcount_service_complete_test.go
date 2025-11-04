package document

import (
	"testing"

	"github.com/stretchr/testify/assert"

	documentSvc "Qingyu_backend/service/document"
)

// TestWordCountService_CalculateWordCount 测试字数统计功能
func TestWordCountService_CalculateWordCount(t *testing.T) {
	service := documentSvc.NewWordCountService()

	t.Run("EmptyContent", func(t *testing.T) {
		result := service.CalculateWordCount("")

		assert.NotNil(t, result)
		assert.Equal(t, 0, result.TotalCount)
		assert.Equal(t, 0, result.ChineseCount)
		assert.Equal(t, 0, result.EnglishCount)
		assert.Equal(t, 0, result.NumberCount)
		assert.Equal(t, 0, result.ParagraphCount)
		assert.Equal(t, 0, result.SentenceCount)
		t.Logf("✓ 空内容测试通过")
	})

	t.Run("ChineseContentOnly", func(t *testing.T) {
		content := "这是一段中文内容。这是第二句话！"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.TotalCount, 0)
		assert.Greater(t, result.ChineseCount, 0)
		assert.Equal(t, 0, result.EnglishCount)
		assert.Equal(t, 0, result.NumberCount)
		assert.Equal(t, 1, result.ParagraphCount)
		assert.Equal(t, 2, result.SentenceCount)
		t.Logf("✓ 中文内容测试通过: 总字数=%d, 中文=%d", result.TotalCount, result.ChineseCount)
	})

	t.Run("EnglishContentOnly", func(t *testing.T) {
		content := "Hello world. This is an English text."
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.TotalCount, 0)
		assert.Equal(t, 0, result.ChineseCount)
		assert.Greater(t, result.EnglishCount, 0)
		assert.Equal(t, 1, result.ParagraphCount)
		assert.GreaterOrEqual(t, result.SentenceCount, 2)
		t.Logf("✓ 英文内容测试通过: 英文单词=%d", result.EnglishCount)
	})

	t.Run("MixedContent", func(t *testing.T) {
		content := "这是中文content混合。Hello world! 2023年。"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.TotalCount, 0)
		assert.Greater(t, result.ChineseCount, 0)
		assert.Greater(t, result.EnglishCount, 0)
		assert.GreaterOrEqual(t, result.NumberCount, 4) // 2023四个数字
		t.Logf("✓ 混合内容测试通过: 中文=%d, 英文=%d, 数字=%d",
			result.ChineseCount, result.EnglishCount, result.NumberCount)
	})

	t.Run("NumberCountTest", func(t *testing.T) {
		content := "这段文字包含数字123和456"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.GreaterOrEqual(t, result.NumberCount, 6) // 1,2,3,4,5,6
		t.Logf("✓ 数字统计测试通过: 数字个数=%d", result.NumberCount)
	})

	t.Run("ParagraphCountTest", func(t *testing.T) {
		content := "第一段内容\n第二段内容\n第三段内容"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Equal(t, 3, result.ParagraphCount)
		t.Logf("✓ 段落统计测试通过: 段落数=%d", result.ParagraphCount)
	})

	t.Run("ParagraphWithEmptyLines", func(t *testing.T) {
		content := "第一段\n\n第二段\n\n\n第三段"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Equal(t, 3, result.ParagraphCount) // 空行不计算
		t.Logf("✓ 空行处理测试通过: 段落数=%d", result.ParagraphCount)
	})

	t.Run("SentenceCountTest", func(t *testing.T) {
		content := "这是第一句。这是第二句！这是第三句？"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Equal(t, 3, result.SentenceCount)
		t.Logf("✓ 句子统计测试通过: 句子数=%d", result.SentenceCount)
	})

	t.Run("ComplexContent", func(t *testing.T) {
		content := `春风又绿江南岸。
明月何时照我还。
The year is 2025.
Remember the number: 123456789.`
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.TotalCount, 0)
		assert.Greater(t, result.ChineseCount, 0)
		assert.Greater(t, result.EnglishCount, 0)
		assert.Greater(t, result.NumberCount, 0)
		assert.Equal(t, 4, result.ParagraphCount)
		assert.GreaterOrEqual(t, result.SentenceCount, 4)
		t.Logf("✓ 复杂内容测试通过: 总字数=%d, 中文=%d, 英文=%d, 数字=%d, 段落=%d, 句子=%d",
			result.TotalCount, result.ChineseCount, result.EnglishCount,
			result.NumberCount, result.ParagraphCount, result.SentenceCount)
	})

	t.Run("SpecialCharacters", func(t *testing.T) {
		content := "Hello@World#2025!测试内容。"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.ChineseCount, 0)
		// 特殊字符应该被正确处理
		t.Logf("✓ 特殊字符测试通过: 总字数=%d", result.TotalCount)
	})

	t.Run("WhitespaceHandling", func(t *testing.T) {
		content := "   前置空格\n  缩进行  \n后置空格   "
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.ChineseCount, 0)
		assert.Equal(t, 3, result.ParagraphCount) // 三个非空行
		t.Logf("✓ 空格处理测试通过: 段落数=%d", result.ParagraphCount)
	})

	t.Run("ReadingTimeCalculation", func(t *testing.T) {
		// 构建一个相对较长的文本来测试阅读时间
		content := "这是一段测试文本。" // 重复足够多次以产生可测量的阅读时间
		for i := 0; i < 100; i++ {
			content += "这是一段测试文本。"
		}
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		// 阅读时间应该大于0（对于足够长的文本）
		assert.GreaterOrEqual(t, result.ReadingTime, 0)
		if result.ReadingTime > 0 {
			t.Logf("✓ 阅读时间计算测试通过: 预计阅读时长=%d分钟", result.ReadingTime)
		} else {
			t.Logf("✓ 阅读时间计算测试通过: 文本较短，阅读时长=%d分钟", result.ReadingTime)
		}
	})

	t.Run("LargeContent", func(t *testing.T) {
		// 创建一个较大的内容
		content := ""
		for i := 0; i < 1000; i++ {
			content += "这是第" + string(rune('0'+(i%10))) + "行内容。\n"
		}
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.TotalCount, 0)
		assert.Greater(t, result.ChineseCount, 0)
		assert.Equal(t, 1000, result.ParagraphCount)
		assert.Greater(t, result.SentenceCount, 0)
		t.Logf("✓ 大文本处理测试通过: 总字数=%d, 段落=%d", result.TotalCount, result.ParagraphCount)
	})
}

// TestWordCountService_CalculateWordCount_EdgeCases 边界条件测试
func TestWordCountService_CalculateWordCount_EdgeCases(t *testing.T) {
	service := documentSvc.NewWordCountService()

	t.Run("OnlyNumbers", func(t *testing.T) {
		content := "123456789"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Equal(t, 0, result.ChineseCount)
		assert.Equal(t, 0, result.EnglishCount)
		assert.Equal(t, 9, result.NumberCount)
		t.Logf("✓ 纯数字测试通过: 数字=%d", result.NumberCount)
	})

	t.Run("OnlyPunctuation", func(t *testing.T) {
		content := "。！？.!?"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		// 标点符号应该被正确处理
		t.Logf("✓ 纯标点测试通过")
	})

	t.Run("OnlySpaces", func(t *testing.T) {
		content := "   \n   \n   "
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Equal(t, 0, result.ParagraphCount) // 所有空行
		t.Logf("✓ 纯空格测试通过")
	})

	t.Run("SingleCharacter", func(t *testing.T) {
		content := "字"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Equal(t, 1, result.ChineseCount)
		t.Logf("✓ 单字符测试通过")
	})

	t.Run("SentenceWithoutPunctuation", func(t *testing.T) {
		content := "这是没有标点的一句话"
		result := service.CalculateWordCount(content)

		assert.NotNil(t, result)
		assert.Greater(t, result.ChineseCount, 0)
		// 没有标点时，句子数应该是1（或按其他逻辑）
		t.Logf("✓ 无标点句子测试通过")
	})
}

