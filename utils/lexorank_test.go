package utils

import (
	"sort"
	"testing"
)

// TestGenerateOrderKey_FirstNode 测试生成第一个节点的排序键
func TestGenerateOrderKey_FirstNode(t *testing.T) {
	result := GenerateOrderKey("", "")

	if result != "a0" {
		t.Errorf("GenerateOrderKey(\"\", \"\") = %s; want \"a0\"", result)
	}
}

// TestGenerateOrderKey_SecondNode 测试生成第二个节点的排序键
func TestGenerateOrderKey_SecondNode(t *testing.T) {
	result := GenerateOrderKey("a0", "")

	expected := "a00"
	if result != expected {
		t.Errorf("GenerateOrderKey(\"a0\", \"\") = %s; want \"%s\"", result, expected)
	}
}

// TestGenerateOrderKey_ThirdNode 测试生成第三个节点的排序键
func TestGenerateOrderKey_ThirdNode(t *testing.T) {
	result := GenerateOrderKey("a00", "")

	expected := "a000"
	if result != expected {
		t.Errorf("GenerateOrderKey(\"a00\", \"\") = %s; want \"%s\"", result, expected)
	}
}

// TestGenerateOrderKey_MultipleLevels 测试生成多层级的排序键
func TestGenerateOrderKey_MultipleLevels(t *testing.T) {
	testCases := []struct {
		prevKey    string
		wantPrefix string
	}{
		{"a0", "a00"},
		{"a00", "a000"},
		{"a000", "a0000"},
		{"b0", "b00"},
		{"z0", "z00"},
	}

	for _, tc := range testCases {
		result := GenerateOrderKey(tc.prevKey, "")
		if result != tc.wantPrefix {
			t.Errorf("GenerateOrderKey(%q, \"\") = %s; want %s", tc.prevKey, result, tc.wantPrefix)
		}
	}
}

// TestGenerateOrderKey_Ordering 测试生成的排序键保持字典序
func TestGenerateOrderKey_Ordering(t *testing.T) {
	keys := []string{
		GenerateOrderKey("", ""),      // a0
		GenerateOrderKey("a0", ""),    // a00
		GenerateOrderKey("a00", ""),   // a000
		GenerateOrderKey("a000", ""),  // a0000
	}

	// 验证keys是有序的
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)

	for i := range keys {
		if keys[i] != sortedKeys[i] {
			t.Errorf("Keys not in order: got %v, want %v", keys, sortedKeys)
			break
		}
	}
}

// TestGenerateSiblingOrderKey_AfterA0 测试在a0后生成同级排序键
func TestGenerateSiblingOrderKey_AfterA0(t *testing.T) {
	result := GenerateSiblingOrderKey("a0")

	expected := "a00"
	if result != expected {
		t.Errorf("GenerateSiblingOrderKey(\"a0\") = %s; want \"%s\"", result, expected)
	}
}

// TestGenerateSiblingOrderKey_AfterA00 测试在a00后生成同级排序键
func TestGenerateSiblingOrderKey_AfterA00(t *testing.T) {
	result := GenerateSiblingOrderKey("a00")

	expected := "a000"
	if result != expected {
		t.Errorf("GenerateSiblingOrderKey(\"a00\") = %s; want \"%s\"", result, expected)
	}
}

// TestGenerateSiblingOrderKey_EmptyKey 测试空参考键的情况
func TestGenerateSiblingOrderKey_EmptyKey(t *testing.T) {
	result := GenerateSiblingOrderKey("")

	expected := "a0"
	if result != expected {
		t.Errorf("GenerateSiblingOrderKey(\"\") = %s; want \"%s\"", result, expected)
	}
}

// TestGenerateSiblingOrderKey_MultipleKeys 测试多个同级排序键的生成
func TestGenerateSiblingOrderKey_MultipleKeys(t *testing.T) {
	testCases := []struct {
		afterKey    string
		wantResult  string
	}{
		{"a0", "a00"},
		{"a00", "a000"},
		{"a000", "a0000"},
		{"b0", "b00"},
		{"z0", "z00"},
	}

	for _, tc := range testCases {
		result := GenerateSiblingOrderKey(tc.afterKey)
		if result != tc.wantResult {
			t.Errorf("GenerateSiblingOrderKey(%q) = %s; want %s", tc.afterKey, result, tc.wantResult)
		}
	}
}

// TestGenerateSiblingOrderKey_Ordering 测试同级排序键保持字典序
func TestGenerateSiblingOrderKey_Ordering(t *testing.T) {
	keys := []string{
		GenerateSiblingOrderKey(""),      // a0
		GenerateSiblingOrderKey("a0"),    // a00
		GenerateSiblingOrderKey("a00"),   // a000
		GenerateSiblingOrderKey("a000"),  // a0000
	}

	// 验证keys是有序的
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)

	for i := range keys {
		if keys[i] != sortedKeys[i] {
			t.Errorf("Sibling keys not in order: got %v, want %v", keys, sortedKeys)
			break
		}
	}
}

// TestLexicographicOrdering 测试生成的排序键满足字典序要求
func TestLexicographicOrdering(t *testing.T) {
	// 模拟文档树中的多个节点
	orderKeys := []string{}
	prevKey := ""

	// 生成10个节点的排序键
	for i := 0; i < 10; i++ {
		newKey := GenerateOrderKey(prevKey, "")
		orderKeys = append(orderKeys, newKey)
		prevKey = newKey
	}

	// 验证排序键的字典序
	for i := 0; i < len(orderKeys)-1; i++ {
		if orderKeys[i] >= orderKeys[i+1] {
			t.Errorf("Order keys not in lexicographic order: %s >= %s", orderKeys[i], orderKeys[i+1])
		}
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		prevKey  string
		nextKey  string
		expected string
	}{
		{"Empty prev and next", "", "", "a0"},
		{"Prev with single digit", "a0", "", "a00"},
		{"Prev with multiple digits", "a000", "", "a0000"},
		{"Letter only prev", "abc", "", "abc0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateOrderKey(tc.prevKey, tc.nextKey)
			if result != tc.expected {
				t.Errorf("GenerateOrderKey(%q, %q) = %s; want %s", tc.prevKey, tc.nextKey, result, tc.expected)
			}
		})
	}
}
