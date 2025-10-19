package audit

import (
	"strings"
	"sync"
)

// DFAFilter DFA敏感词过滤器
type DFAFilter struct {
	root *TrieNode
	mu   sync.RWMutex
}

// TrieNode Trie树节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	word     string // 完整的敏感词
	level    int    // 敏感词等级
	category string // 敏感词分类
}

// NewDFAFilter 创建DFA过滤器
func NewDFAFilter() *DFAFilter {
	return &DFAFilter{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

// AddWord 添加敏感词
func (f *DFAFilter) AddWord(word string, level int, category string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if word == "" {
		return
	}

	// 转换为小写并去除空格
	word = strings.ToLower(strings.TrimSpace(word))
	runes := []rune(word)

	node := f.root
	for _, r := range runes {
		if node.children[r] == nil {
			node.children[r] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[r]
	}

	node.isEnd = true
	node.word = word
	node.level = level
	node.category = category
}

// BatchAddWords 批量添加敏感词
func (f *DFAFilter) BatchAddWords(words []SensitiveWordInfo) {
	for _, w := range words {
		f.AddWord(w.Word, w.Level, w.Category)
	}
}

// RemoveWord 移除敏感词
func (f *DFAFilter) RemoveWord(word string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if word == "" {
		return
	}

	word = strings.ToLower(strings.TrimSpace(word))
	runes := []rune(word)

	node := f.root
	parents := make([]*TrieNode, 0, len(runes))
	parents = append(parents, node)

	// 找到目标节点
	for _, r := range runes {
		if node.children[r] == nil {
			return // 词不存在
		}
		node = node.children[r]
		parents = append(parents, node)
	}

	// 标记为非结束节点
	node.isEnd = false
	node.word = ""
	node.level = 0
	node.category = ""

	// 如果没有子节点，向上清理
	for i := len(parents) - 1; i > 0; i-- {
		if len(parents[i].children) == 0 && !parents[i].isEnd {
			r := runes[i-1]
			delete(parents[i-1].children, r)
		} else {
			break
		}
	}
}

// Check 检查文本是否包含敏感词
func (f *DFAFilter) Check(text string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	text = strings.ToLower(text)
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		node := f.root
		j := i

		for j < len(runes) {
			r := runes[j]
			if node.children[r] == nil {
				break
			}
			node = node.children[r]
			if node.isEnd {
				return true // 找到敏感词
			}
			j++
		}
	}

	return false
}

// FindAll 查找所有敏感词
func (f *DFAFilter) FindAll(text string) []MatchResult {
	f.mu.RLock()
	defer f.mu.RUnlock()

	text = strings.ToLower(text)
	runes := []rune(text)
	results := make([]MatchResult, 0)

	for i := 0; i < len(runes); i++ {
		node := f.root
		j := i
		lastMatch := -1
		var matchedNode *TrieNode

		for j < len(runes) {
			r := runes[j]
			if node.children[r] == nil {
				break
			}
			node = node.children[r]
			if node.isEnd {
				lastMatch = j
				matchedNode = node
			}
			j++
		}

		if lastMatch >= 0 && matchedNode != nil {
			// 找到敏感词
			result := MatchResult{
				Word:     matchedNode.word,
				Start:    i,
				End:      lastMatch + 1,
				Level:    matchedNode.level,
				Category: matchedNode.category,
				Context:  extractContext(runes, i, lastMatch+1, 10),
			}
			results = append(results, result)
			i = lastMatch // 跳过已匹配的部分
		}
	}

	return results
}

// Replace 替换敏感词
func (f *DFAFilter) Replace(text string, replacement string) string {
	matches := f.FindAll(text)
	if len(matches) == 0 {
		return text
	}

	runes := []rune(text)

	// 从后往前替换，避免索引变化
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		// 用replacement替换敏感词
		before := runes[:match.Start]
		after := runes[match.End:]
		runes = append(before, []rune(replacement)...)
		runes = append(runes, after...)
	}

	return string(runes)
}

// ReplaceWithMask 用掩码替换敏感词（如：***）
func (f *DFAFilter) ReplaceWithMask(text string, mask rune) string {
	matches := f.FindAll(text)
	if len(matches) == 0 {
		return text
	}

	runes := []rune(text)

	// 标记需要替换的位置
	for _, match := range matches {
		for i := match.Start; i < match.End; i++ {
			runes[i] = mask
		}
	}

	return string(runes)
}

// GetStatistics 获取过滤器统计信息
func (f *DFAFilter) GetStatistics() FilterStatistics {
	f.mu.RLock()
	defer f.mu.RUnlock()

	stats := FilterStatistics{
		TotalWords: 0,
		ByCategory: make(map[string]int),
		ByLevel:    make(map[int]int),
	}

	// 遍历Trie树统计
	var traverse func(*TrieNode)
	traverse = func(node *TrieNode) {
		if node.isEnd {
			stats.TotalWords++
			stats.ByCategory[node.category]++
			stats.ByLevel[node.level]++
		}
		for _, child := range node.children {
			traverse(child)
		}
	}

	traverse(f.root)
	return stats
}

// Clear 清空所有敏感词
func (f *DFAFilter) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.root = &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

// MatchResult 匹配结果
type MatchResult struct {
	Word     string `json:"word"`     // 匹配的敏感词
	Start    int    `json:"start"`    // 起始位置（字符索引）
	End      int    `json:"end"`      // 结束位置（字符索引）
	Level    int    `json:"level"`    // 敏感词等级
	Category string `json:"category"` // 敏感词分类
	Context  string `json:"context"`  // 上下文
}

// SensitiveWordInfo 敏感词信息
type SensitiveWordInfo struct {
	Word     string
	Level    int
	Category string
}

// FilterStatistics 过滤器统计信息
type FilterStatistics struct {
	TotalWords int            `json:"totalWords"` // 总敏感词数
	ByCategory map[string]int `json:"byCategory"` // 按分类统计
	ByLevel    map[int]int    `json:"byLevel"`    // 按等级统计
}

// extractContext 提取上下文
func extractContext(runes []rune, start, end, contextLen int) string {
	contextStart := start - contextLen
	if contextStart < 0 {
		contextStart = 0
	}

	contextEnd := end + contextLen
	if contextEnd > len(runes) {
		contextEnd = len(runes)
	}

	// 添加省略号标记
	prefix := ""
	suffix := ""
	if contextStart > 0 {
		prefix = "..."
	}
	if contextEnd < len(runes) {
		suffix = "..."
	}

	return prefix + string(runes[contextStart:contextEnd]) + suffix
}
