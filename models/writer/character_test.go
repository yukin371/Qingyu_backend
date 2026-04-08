package writer

import (
	"Qingyu_backend/models/writer/base"
	"testing"
)

func TestCharacterRelation_IsValidAtChapter(t *testing.T) {
	chapterOrderMap := map[string]int{
		"chapter-1": 1,
		"chapter-2": 2,
		"chapter-3": 3,
		"chapter-4": 4,
		"chapter-5": 5,
	}

	tests := []struct {
		name              string
		validFromChapter  *string
		validUntilChapter *string
		testChapterOrder  int
		expected          bool
	}{
		{
			name:              "全局有效关系",
			validFromChapter:  nil,
			validUntilChapter: nil,
			testChapterOrder:  3,
			expected:          true,
		},
		{
			name:              "在第3章生效，测试第5章",
			validFromChapter:  strPtr("chapter-3"),
			validUntilChapter: nil,
			testChapterOrder:  5,
			expected:          true,
		},
		{
			name:              "在第3章生效，测试第2章",
			validFromChapter:  strPtr("chapter-3"),
			validUntilChapter: nil,
			testChapterOrder:  2,
			expected:          false,
		},
		{
			name:              "在第5章失效，测试第3章",
			validFromChapter:  nil,
			validUntilChapter: strPtr("chapter-5"),
			testChapterOrder:  3,
			expected:          true,
		},
		{
			name:              "在第5章失效，测试第5章",
			validFromChapter:  nil,
			validUntilChapter: strPtr("chapter-5"),
			testChapterOrder:  5,
			expected:          false,
		},
		{
			name:              "第3章生效，第5章失效，测试第4章",
			validFromChapter:  strPtr("chapter-3"),
			validUntilChapter: strPtr("chapter-5"),
			testChapterOrder:  4,
			expected:          true,
		},
		{
			name:              "第3章生效，第5章失效，测试第2章",
			validFromChapter:  strPtr("chapter-3"),
			validUntilChapter: strPtr("chapter-5"),
			testChapterOrder:  2,
			expected:          false,
		},
		{
			name:              "第3章生效，第5章失效，测试第5章",
			validFromChapter:  strPtr("chapter-3"),
			validUntilChapter: strPtr("chapter-5"),
			testChapterOrder:  5,
			expected:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relation := &CharacterRelation{
				ValidFromChapterID:  tt.validFromChapter,
				ValidUntilChapterID: tt.validUntilChapter,
			}

			result := relation.IsValidAtChapter(tt.testChapterOrder, chapterOrderMap)
			if result != tt.expected {
				t.Errorf("IsValidAtChapter() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCharacterRelation_IsValidAtChapter_NonExistentChapter(t *testing.T) {
	// 测试章节不存在的情况（保守处理：认为关系有效）
	chapterOrderMap := map[string]int{
		"chapter-1": 1,
		"chapter-2": 2,
	}

	relation := &CharacterRelation{
		ValidFromChapterID:  strPtr("chapter-nonexistent"),
		ValidUntilChapterID: strPtr("chapter-nonexistent"),
	}

	// 章节不存在时应该返回true（保守处理）
	if !relation.IsValidAtChapter(1, chapterOrderMap) {
		t.Error("IsValidAtChapter() should return true for non-existent chapters (conservative approach)")
	}
}

// 辅助函数：返回字符串指针
func strPtr(s string) *string {
	return &s
}

func TestCharacterEntityTypeField(t *testing.T) {
	c := Character{
		NamedEntity: base.NamedEntity{Name: "测试角色"},
		EntityType:  EntityTypeCharacter,
		StateFields: map[string]StateValue{
			"恐惧值": {Current: 30.0, Min: ptrFloat64(0), Max: ptrFloat64(100)},
		},
	}
	if c.EntityType != EntityTypeCharacter {
		t.Errorf("expected character, got %s", c.EntityType)
	}
	if c.StateFields["恐惧值"].Current != 30.0 {
		t.Errorf("expected 30, got %v", c.StateFields["恐惧值"].Current)
	}
	if !EntityTypeCharacter.IsValid() {
		t.Error("EntityTypeCharacter should be valid")
	}
	if EntityType("invalid").IsValid() {
		t.Error("invalid EntityType should not be valid")
	}
}

func ptrFloat64(v float64) *float64 { return &v }

func TestCharacterRelationCrossType(t *testing.T) {
	r := CharacterRelation{
		FromID:   "char-1",
		ToID:     "item-1",
		FromType: EntityTypeCharacter,
		ToType:   EntityTypeItem,
	}
	if r.FromType != EntityTypeCharacter {
		t.Errorf("expected character, got %s", r.FromType)
	}
	if r.ToType != EntityTypeItem {
		t.Errorf("expected item, got %s", r.ToType)
	}
}
