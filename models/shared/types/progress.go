package types

import (
	"fmt"
	"math"
)

// Progress 进度类型（0.0-1.0）
type Progress float32

const (
	// ProgressMin 最小进度
	ProgressMin Progress = 0.0
	// ProgressMax 最大进度
	ProgressMax Progress = 1.0
	// ProgressZero 零进度
	ProgressZero Progress = 0.0
	// ProgressFull 完整进度
	ProgressFull Progress = 1.0
)

var (
	ErrInvalidProgress = fmt.Errorf("progress must be between %.1f and %.1f", ProgressMin, ProgressMax)
)

// NewProgress 创建进度（0-1）
func NewProgress(value float32) (Progress, error) {
	p := Progress(value)
	if !p.IsValid() {
		return ProgressZero, ErrInvalidProgress
	}
	return p, nil
}

// NewProgressFromPercent 从百分比创建进度（0-100）
func NewProgressFromPercent(percent int) (Progress, error) {
	if percent < 0 || percent > 100 {
		return ProgressZero, ErrInvalidProgress
	}
	return Progress(float32(percent) / 100.0), nil
}

// NewProgressFromRatio 从比例创建进度
func NewProgressFromRatio(numerator, denominator int) (Progress, error) {
	if denominator <= 0 {
		return ProgressZero, fmt.Errorf("denominator must be positive")
	}
	if numerator < 0 {
		return ProgressZero, fmt.Errorf("numerator cannot be negative")
	}
	if numerator > denominator {
		return ProgressFull, nil // 超过 100% 时视为完成
	}
	return Progress(float32(numerator) / float32(denominator)), nil
}

// MustProgress 创建进度（panic on invalid）
func MustProgress(value float32) Progress {
	p, err := NewProgress(value)
	if err != nil {
		panic(err)
	}
	return p
}

// MustProgressFromPercent 从百分比创建进度（panic on invalid）
func MustProgressFromPercent(percent int) Progress {
	p, err := NewProgressFromPercent(percent)
	if err != nil {
		panic(err)
	}
	return p
}

// IsValid 检查进度是否有效
func (p Progress) IsValid() bool {
	return p >= ProgressMin && p <= ProgressMax
}

// ToPercent 转换为百分比（0-100）
func (p Progress) ToPercent() int {
	return int(math.Round(float64(p) * 100))
}

// ToFloat 转换为 float32
func (p Progress) ToFloat() float32 {
	return float32(p)
}

// String 格式化为百分比字符串（如 "75%"）
func (p Progress) String() string {
	return fmt.Sprintf("%d%%", p.ToPercent())
}

// IsComplete 是否完成（100%）
func (p Progress) IsComplete() bool {
	return p >= ProgressFull
}

// IsStarted 是否已开始（> 0%）
func (p Progress) IsStarted() bool {
	return p > ProgressZero
}

// IsZero 是否为零进度
func (p Progress) IsZero() bool {
	return p == ProgressZero
}

// Add 累加进度
func (p Progress) Add(other Progress) Progress {
	result := float32(p) + float32(other)
	if result > float32(ProgressMax) {
		return ProgressFull
	}
	return Progress(result)
}

// Sub 减去进度
func (p Progress) Sub(other Progress) Progress {
	result := float32(p) - float32(other)
	if result < float32(ProgressMin) {
		return ProgressZero
	}
	return Progress(result)
}

// Percentage 计算占比（相对于另一个进度）
func (p Progress) Percentage(other Progress) int {
	if other == ProgressZero {
		return 0
	}
	ratio := float32(p) / float32(other)
	return int(math.Round(float64(ratio) * 100))
}

// Round 四舍五入到指定百分比精度
func (p Progress) Round(precision int) Progress {
	multiplier := math.Pow(10, float64(precision))
	rounded := math.Round(float64(p)*100*multiplier) / multiplier
	result := Progress(rounded / 100)
	if result > ProgressFull {
		return ProgressFull
	}
	return result
}

// Clamp 限制进度在有效范围内
func (p Progress) Clamp() Progress {
	if p < ProgressMin {
		return ProgressMin
	}
	if p > ProgressMax {
		return ProgressMax
	}
	return p
}

// Compare 比较进度（-1: <, 0: =, 1: >）
func (p Progress) Compare(other Progress) int {
	switch {
	case p < other:
		return -1
	case p > other:
		return 1
	default:
		return 0
	}
}

// GetStage 获取进度阶段
// 0-25%: 起始阶段
// 26-50%: 早期阶段
// 51-75%: 中期阶段
// 76-99%: 后期阶段
// 100%: 完成
func (p Progress) GetStage() string {
	percent := p.ToPercent()
	switch {
	case percent == 0:
		return "未开始"
	case percent <= 25:
		return "起始阶段"
	case percent <= 50:
		return "早期阶段"
	case percent <= 75:
		return "中期阶段"
	case percent < 100:
		return "后期阶段"
	default:
		return "已完成"
	}
}

// ProgressSlice 进度切片类型
type ProgressSlice []Progress

// Average 计算平均进度
func (ps ProgressSlice) Average() Progress {
	if len(ps) == 0 {
		return ProgressZero
	}

	var sum float32
	for _, p := range ps {
		sum += float32(p)
	}

	return Progress(sum / float32(len(ps)))
}

// Max 获取最大进度
func (ps ProgressSlice) Max() Progress {
	if len(ps) == 0 {
		return ProgressZero
	}

	max := ps[0]
	for _, p := range ps {
		if p > max {
			max = p
		}
	}
	return max
}

// Min 获取最小进度
func (ps ProgressSlice) Min() Progress {
	if len(ps) == 0 {
		return ProgressZero
	}

	min := ps[0]
	for _, p := range ps {
		if p < min {
			min = p
		}
	}
	return min
}
