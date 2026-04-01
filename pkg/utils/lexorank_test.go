package utils

import (
	"sort"
	"testing"
)

// TestGenerateOrderKey_FirstNode 测试生成第一个节点的排序键
func TestGenerateOrderKey_FirstNode(t *testing.T) {
	result := GenerateOrderKey("", "")

	// 新实现使用 base62，初始键是中间值 "UYYYYYY"（6个字符）
	// 只要保证非空且有序即可
	if result == "" {
		t.Errorf("GenerateOrderKey(\"\", \"\") should not be empty")
	}
	t.Logf("Initial key: %s (length: %d)", result, len(result))
}

// TestGenerateOrderKey_Sequential 测试顺序生成的排序键保持字典序
func TestGenerateOrderKey_Sequential(t *testing.T) {
	keys := []string{}
	prevKey := ""

	// 生成10个顺序键
	for i := 0; i < 10; i++ {
		newKey := GenerateOrderKey(prevKey, "")
		if newKey == "" {
			t.Errorf("GenerateOrderKey(%q, \"\") returned empty string", prevKey)
			break
		}
		if prevKey != "" && newKey <= prevKey {
			t.Errorf("GenerateOrderKey(%q, \"\") = %q; should be > prevKey", prevKey, newKey)
			break
		}
		keys = append(keys, newKey)
		prevKey = newKey
	}

	t.Logf("Generated keys: %v", keys)

	// 验证长度不会无限增长（base62的特性）
	maxLen := 0
	for _, k := range keys {
		if len(k) > maxLen {
			maxLen = len(k)
		}
	}
	// 10个顺序键，长度应该保持在合理范围（不超过20字符）
	if maxLen > 20 {
		t.Errorf("Key length grew too large: %d characters after 10 keys", maxLen)
	}
}

// TestGenerateOrderKey_Between 测试在两个键之间生成新键
func TestGenerateOrderKey_Between(t *testing.T) {
	// 先生成初始键
	key1 := GenerateOrderKey("", "")
	// 在末尾添加
	key3 := GenerateOrderKey(key1, "")
	// 在中间插入
	key2 := GenerateOrderKey(key1, key3)

	t.Logf("key1: %s, key2: %s, key3: %s", key1, key2, key3)

	// 验证顺序
	if !(key1 < key2 && key2 < key3) {
		t.Errorf("Keys not in order: key1=%q < key2=%q < key3=%q should be true", key1, key2, key3)
	}
}

// TestGenerateOrderKey_Ordering 测试生成的排序键保持字典序
func TestGenerateOrderKey_Ordering(t *testing.T) {
	keys := []string{}
	prevKey := ""

	for i := 0; i < 100; i++ {
		newKey := GenerateOrderKey(prevKey, "")
		keys = append(keys, newKey)
		prevKey = newKey
	}

	// 验证keys是有序的
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)

	for i := range keys {
		if keys[i] != sortedKeys[i] {
			t.Errorf("Keys not in order: got %v", keys[:10])
			break
		}
	}

	// 验证长度控制
	lastKey := keys[len(keys)-1]
	if len(lastKey) > 50 {
		t.Errorf("Key length too large after 100 iterations: %d", len(lastKey))
	}
	t.Logf("100 keys generated, last key length: %d", len(lastKey))
}

// TestGenerateSiblingOrderKey_Sequential 测试同级排序键的顺序生成
func TestGenerateSiblingOrderKey_Sequential(t *testing.T) {
	keys := []string{}
	prevKey := ""

	for i := 0; i < 10; i++ {
		newKey := GenerateSiblingOrderKey(prevKey)
		if newKey == "" {
			t.Errorf("GenerateSiblingOrderKey(%q) returned empty string", prevKey)
			break
		}
		if prevKey != "" && newKey <= prevKey {
			t.Errorf("GenerateSiblingOrderKey(%q) = %q; should be > prevKey", prevKey, newKey)
			break
		}
		keys = append(keys, newKey)
		prevKey = newKey
	}

	t.Logf("Generated sibling keys: %v", keys)
}

// TestGenerateSiblingOrderKey_EmptyKey 测试空参考键的情况
func TestGenerateSiblingOrderKey_EmptyKey(t *testing.T) {
	result := GenerateSiblingOrderKey("")

	if result == "" {
		t.Errorf("GenerateSiblingOrderKey(\"\") should not be empty")
	}
	t.Logf("Initial sibling key: %s", result)
}

// TestGenerateSiblingOrderKey_Ordering 测试同级排序键保持字典序
func TestGenerateSiblingOrderKey_Ordering(t *testing.T) {
	keys := []string{}
	prevKey := ""

	for i := 0; i < 50; i++ {
		newKey := GenerateSiblingOrderKey(prevKey)
		keys = append(keys, newKey)
		prevKey = newKey
	}

	// 验证keys是有序的
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)

	for i := range keys {
		if keys[i] != sortedKeys[i] {
			t.Errorf("Sibling keys not in order")
			break
		}
	}

	// 验证长度控制
	lastKey := keys[len(keys)-1]
	if len(lastKey) > 30 {
		t.Errorf("Key length too large after 50 sibling keys: %d", len(lastKey))
	}
	t.Logf("50 sibling keys generated, last key length: %d", len(lastKey))
}

// TestLexicographicOrdering 测试生成的排序键满足字典序要求
func TestLexicographicOrdering(t *testing.T) {
	// 模拟文档树中的多个节点
	orderKeys := []string{}
	prevKey := ""

	// 生成100个节点的排序键
	for i := 0; i < 100; i++ {
		newKey := GenerateOrderKey(prevKey, "")
		orderKeys = append(orderKeys, newKey)
		prevKey = newKey
	}

	// 验证排序键的字典序
	for i := 0; i < len(orderKeys)-1; i++ {
		if orderKeys[i] >= orderKeys[i+1] {
			t.Errorf("Order keys not in lexicographic order: %s >= %s", orderKeys[i], orderKeys[i+1])
			break
		}
	}
}

// TestKeyLengthGrowth 测试键长度增长
// 与旧的简化实现不同，新的 base62 实现应该保持键长度稳定
func TestKeyLengthGrowth(t *testing.T) {
	keys := []string{}
	prevKey := ""

	// 生成1000个键，验证长度不会线性增长
	for i := 0; i < 1000; i++ {
		newKey := GenerateOrderKey(prevKey, "")
		keys = append(keys, newKey)
		prevKey = newKey
	}

	// 检查最后一个键的长度
	lastKey := keys[len(keys)-1]
	t.Logf("After 1000 keys, last key: %q (length: %d)", lastKey, len(lastKey))

	// 新的 base62 实现应该保持键长度在合理范围
	// 旧实现会产生1001字符的键，新实现应该保持在10字符以内
	if len(lastKey) > 10 {
		t.Errorf("Key length too large: %d (expected < 10 for base62)", len(lastKey))
	}
}

// TestInsertInMiddle 测试在中间插入的能力
func TestInsertInMiddle(t *testing.T) {
	// 生成初始序列
	keys := []string{}
	prevKey := ""
	for i := 0; i < 10; i++ {
		newKey := GenerateOrderKey(prevKey, "")
		keys = append(keys, newKey)
		prevKey = newKey
	}

	// 在第5和第6个键之间插入新键
	insertKey := GenerateOrderKey(keys[4], keys[5])

	// 验证插入位置正确
	if !(keys[4] < insertKey && insertKey < keys[5]) {
		t.Errorf("Insert key not between neighbors: %q < %q < %q should be true", keys[4], insertKey, keys[5])
	}

	// 验证可以继续在插入的键之间再插入
	insertKey2 := GenerateOrderKey(keys[4], insertKey)
	if !(keys[4] < insertKey2 && insertKey2 < insertKey) {
		t.Errorf("Second insert key not between neighbors: %q < %q < %q should be true", keys[4], insertKey2, insertKey)
	}

	t.Logf("Original: %s...%s", keys[4], keys[5])
	t.Logf("Insert1: %s, Insert2: %s", insertKey, insertKey2)
}

// TestGenerator_Between 测试 Generator 的 Between 方法
func TestGenerator_Between(t *testing.T) {
	g := NewGenerator()

	// 测试初始键
	initial, err := g.Initial()
	if err != nil {
		t.Fatalf("Initial() failed: %v", err)
	}
	if initial == "" {
		t.Error("Initial key should not be empty")
	}

	// 测试在两个键之间生成
	key1 := initial
	key3, err := g.Next(key1)
	if err != nil {
		t.Fatalf("Next() failed: %v", err)
	}

	key2, err := g.Between(key1, key3)
	if err != nil {
		t.Fatalf("Between() failed: %v", err)
	}

	if !(string(key1) < string(key2) && string(key2) < string(key3)) {
		t.Errorf("Between key not in order: %s < %s < %s", key1, key2, key3)
	}
}

// TestGenerator_MultipleInsertions 测试在固定两点之间多次插入
func TestGenerator_MultipleInsertions(t *testing.T) {
	g := NewGenerator()

	// 生成初始键和末尾键
	startKey, _ := g.Initial()
	endKey, _ := g.Next(startKey)

	// 在 startKey 和 endKey 之间多次插入
	// 每次都在 startKey 和上一次插入的键之间插入
	prevKey := endKey
	for i := 0; i < 10; i++ {
		newKey, err := g.Between(startKey, prevKey)
		if err != nil {
			t.Fatalf("Between() failed at iteration %d: %v", i, err)
		}
		// 验证新键确实在两者之间
		if !(string(startKey) < string(newKey) && string(newKey) < string(prevKey)) {
			t.Errorf("Insert not between: %s < %s < %s", startKey, newKey, prevKey)
		}
		prevKey = newKey
	}

	t.Logf("After 10 insertions between %s and %s, final key: %s", startKey, endKey, prevKey)
}

// TestBackwardCompatibility 测试向后兼容性
// 确保旧代码调用仍然能正常工作
func TestBackwardCompatibility(t *testing.T) {
	// 测试空参数调用
	key1 := GenerateOrderKey("", "")
	if key1 == "" {
		t.Error("GenerateOrderKey with empty args should return non-empty key")
	}

	// 测试只有一个参数
	key2 := GenerateOrderKey(key1, "")
	if key2 == "" {
		t.Error("GenerateOrderKey with prevKey only should return non-empty key")
	}
	if key2 <= key1 {
		t.Errorf("GenerateOrderKey(%q, \"\") = %q should be > prevKey", key1, key2)
	}

	// 测试 GenerateSiblingOrderKey
	key3 := GenerateSiblingOrderKey("")
	if key3 == "" {
		t.Error("GenerateSiblingOrderKey with empty arg should return non-empty key")
	}

	key4 := GenerateSiblingOrderKey(key3)
	if key4 == "" {
		t.Error("GenerateSiblingOrderKey should return non-empty key")
	}
	if key4 <= key3 {
		t.Errorf("GenerateSiblingOrderKey(%q) = %q should be > afterKey", key3, key4)
	}
}
