package writer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchKeyword_PinyinFullAndInitials(t *testing.T) {
	mode, ok := matchKeyword("zhang", "张三")
	assert.True(t, ok)
	assert.Equal(t, "prefix", mode)

	mode, ok = matchKeyword("zha", "张三")
	assert.True(t, ok)
	assert.Equal(t, "prefix", mode)

	mode, ok = matchKeyword("zs", "张三")
	assert.True(t, ok)
	assert.Contains(t, []string{"exact", "prefix"}, mode)
}

func TestMatchKeyword_ChineseDirectAndAlias(t *testing.T) {
	mode, ok := matchKeyword("长安", "长安城")
	assert.True(t, ok)
	assert.Equal(t, "prefix", mode)

	mode, ok = matchKeyword("chuanshuo", "A角色", "传说中的人")
	assert.True(t, ok)
	assert.Equal(t, "alias", mode)
}

func TestBuildSearchTokens_ContainsRawAndPinyin(t *testing.T) {
	tokens := buildSearchTokens("张三")
	assert.NotEmpty(t, tokens)

	hasRaw := false
	hasFull := false
	hasInitials := false

	for _, token := range tokens {
		switch token {
		case "张三":
			hasRaw = true
		case "zhangsan":
			hasFull = true
		case "zs":
			hasInitials = true
		}
	}

	assert.True(t, hasRaw)
	assert.True(t, hasFull)
	assert.True(t, hasInitials)
}
