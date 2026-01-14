// Package generators 提供数据生成器
package generators

import (
	"testing"
)

func TestBaseGenerator_BookName(t *testing.T) {
	gen := NewBaseGenerator()

	name := gen.BookName("仙侠")
	if name == "" {
		t.Error("BookName 不应该返回空字符串")
	}

	name2 := gen.BookName("仙侠")
	if name == name2 {
		t.Error("BookName 应该生成不同的书名")
	}
}

func TestBaseGenerator_ChapterContent(t *testing.T) {
	gen := NewBaseGenerator()

	content := gen.ChapterContent(500, 3000)
	if len(content) < 500 {
		t.Errorf("内容长度应该至少500字符，实际获得 %d", len(content))
	}
}

func TestBaseGenerator_Username(t *testing.T) {
	gen := NewBaseGenerator()

	username := gen.Username("reader")
	if username == "" {
		t.Error("Username 不应该返回空字符串")
	}

	authorName := gen.Username("author")
	if authorName == "" {
		t.Error("作者用户名不应该返回空字符串")
	}
}
