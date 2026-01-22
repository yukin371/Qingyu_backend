package types

import (
	"fmt"
	"math"
	"strconv"
)

// Money 金额类型（最小货币单位：分）
type Money int64

const (
	// MoneyZero 零金额
	MoneyZero Money = 0

	// CentsPerYuan 每元对应的分数
	CentsPerYuan Money = 100

	// MaxMoney 最大金额（约 922 万亿）
	MaxMoney = math.MaxInt64
)

// NewMoneyFromYuan 从元创建金额（浮点）
func NewMoneyFromYuan(yuan float64) Money {
	// 四舍五入到分
	return Money(int64(math.Round(yuan * float64(CentsPerYuan))))
}

// NewMoneyFromCents 从分创建金额
func NewMoneyFromCents(cents int64) Money {
	return Money(cents)
}

// ToYuan 转换为元（浮点，仅用于展示）
func (m Money) ToYuan() float64 {
	return float64(m) / float64(CentsPerYuan)
}

// ToCents 转换为分（int64）
func (m Money) ToCents() int64 {
	return int64(m)
}

// String 格式化为货币字符串（如 "¥12.99"）
func (m Money) String() string {
	yuan := m.ToYuan()
	return fmt.Sprintf("¥%.2f", yuan)
}

// Add 金额相加
func (m Money) Add(other Money) Money {
	return m + other
}

// Sub 金额相减
func (m Money) Sub(other Money) Money {
	return m - other
}

// Mul 金额乘法（乘以系数）
func (m Money) Mul(factor float64) Money {
	return NewMoneyFromYuan(m.ToYuan() * factor)
}

// Div 金额除法（除以系数）
func (m Money) Div(divisor float64) Money {
	return NewMoneyFromYuan(m.ToYuan() / divisor)
}

// IsZero 是否为零
func (m Money) IsZero() bool {
	return m == MoneyZero
}

// IsNegative 是否为负
func (m Money) IsNegative() bool {
	return m < MoneyZero
}

// IsPositive 是否为正
func (m Money) IsPositive() bool {
	return m > MoneyZero
}

// Compare 比较（-1: <, 0: =, 1: >）
func (m Money) Compare(other Money) int {
	switch {
	case m < other:
		return -1
	case m > other:
		return 1
	default:
		return 0
	}
}

// GreaterThan 是否大于
func (m Money) GreaterThan(other Money) bool {
	return m > other
}

// LessThan 是否小于
func (m Money) LessThan(other Money) bool {
	return m < other
}

// GreaterThanOrEqual 是否大于等于
func (m Money) GreaterThanOrEqual(other Money) bool {
	return m >= other
}

// LessThanOrEqual 是否小于等于
func (m Money) LessThanOrEqual(other Money) bool {
	return m <= other
}

// Abs 取绝对值
func (m Money) Abs() Money {
	if m < MoneyZero {
		return -m
	}
	return m
}

// Min 返回较小值
func MinMoney(a, b Money) Money {
	if a < b {
		return a
	}
	return b
}

// Max 返回较大值
func MaxMoneyValue(a, b Money) Money {
	if a > b {
		return a
	}
	return b
}

// ParseMoney 从字符串解析金额
// 支持格式："12.99", "¥12.99", "1299"（分）
func ParseMoney(s string) (Money, error) {
	// 移除货币符号
	if len(s) > 0 && s[0] == '¥' {
		s = s[1:]
	}

	// 尝试解析为浮点数（元）
	if len(s) > 0 && (s[0] == '-' || (s[0] >= '0' && s[0] <= '9')) {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return NewMoneyFromYuan(f), nil
		}
	}

	// 尝试解析为整数（分）
	cents, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return NewMoneyFromCents(cents), nil
	}

	return MoneyZero, fmt.Errorf("invalid money format: %s", s)
}

// ApplyDiscount 应用折扣（0-100）
func (m Money) ApplyDiscount(discountPercent int) Money {
	if discountPercent <= 0 {
		return m
	}
	if discountPercent >= 100 {
		return MoneyZero
	}
	factor := float64(100-discountPercent) / 100.0
	return m.Mul(factor)
}
