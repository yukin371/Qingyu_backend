package utils

import (
	"strings"
)

// GenerateOrderKey 在两个key之间生成新的排序键
// P1简化实现：使用简单的字符串追加策略
// 完整的base62算法实现将在P2阶段完成
//
// 参数:
//   - prevKey: 前一个节点的排序键，为空时表示这是第一个节点
//   - nextKey: 后一个节点的排序键（P1阶段暂不使用）
//
// 返回:
//   - 新生成的排序键
func GenerateOrderKey(prevKey, nextKey string) string {
	if prevKey == "" {
		// 第一个元素，使用默认起始键
		return "a0"
	}

	// P1简化实现：在prevKey后追加"0"
	// 例如: "a0" -> "a00", "a00" -> "a000"
	// 这种简化方式能够保持字典序，满足基本排序需求
	base := strings.TrimRight(prevKey, "0123456789")
	numStr := strings.TrimLeft(prevKey, base)

	// 如果prevKey不以数字结尾，直接追加"0"
	if numStr == "" {
		return prevKey + "0"
	}

	// 在prevKey后追加"0"
	return prevKey + "0"
}

// GenerateSiblingOrderKey 生成同级排序键
// 用于在同一父节点下追加新节点
//
// P1简化实现：在afterKey后追加"0"
//
// 参数:
//   - afterKey: 参考节点的排序键，新节点将插入到该节点之后
//
// 返回:
//   - 新生成的排序键
func GenerateSiblingOrderKey(afterKey string) string {
	if afterKey == "" {
		// 如果没有参考节点，返回默认起始键
		return "a0"
	}

	// P1简化实现：在afterKey后追加"0"
	// 例如: "a0" -> "a00", "a00" -> "a000"
	return afterKey + "0"
}
