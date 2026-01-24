package types

import (
	"fmt"
	"math"
)

// Rating 评分类型（0.0-5.0）
type Rating float64

const (
	// RatingMin 最小评分
	RatingMin Rating = 0.0
	// RatingMax 最大评分
	RatingMax Rating = 5.0
	// RatingDefault 默认评分
	RatingDefault Rating = 0.0
)

var (
	ErrInvalidRating = fmt.Errorf("rating must be between %.1f and %.1f", RatingMin, RatingMax)
)

// NewRating 创建评分
func NewRating(value float64) (Rating, error) {
	r := Rating(value)
	if !r.IsValid() {
		return RatingDefault, ErrInvalidRating
	}
	return r, nil
}

// MustRating 创建评分（panic on invalid）
func MustRating(value float64) Rating {
	r, err := NewRating(value)
	if err != nil {
		panic(err)
	}
	return r
}

// IsValid 检查评分是否有效
func (r Rating) IsValid() bool {
	return r >= RatingMin && r <= RatingMax
}

// ToFloat 转换为 float64
func (r Rating) ToFloat() float64 {
	return float64(r)
}

// String 格式化为字符串（保留 1 位小数）
func (r Rating) String() string {
	return fmt.Sprintf("%.1f", r.ToFloat())
}

// IsZero 是否为零分
func (r Rating) IsZero() bool {
	return r == RatingMin
}

// IsFull 是否满分
func (r Rating) IsFull() bool {
	return r == RatingMax
}

// Round 四舍五入到 1 位小数
func (r Rating) Round() Rating {
	return Rating(math.Round(float64(r)*10) / 10)
}

// RatingDistribution 评分分布
type RatingDistribution map[string]int64 // key: "1", "2", "3", "4", "5"

// NewRatingDistribution 创建空分布
func NewRatingDistribution() RatingDistribution {
	return RatingDistribution{
		"1": 0,
		"2": 0,
		"3": 0,
		"4": 0,
		"5": 0,
	}
}

// Add 添加评分
func (rd RatingDistribution) Add(rating Rating) error {
	if !rating.IsValid() {
		return ErrInvalidRating
	}

	// 将评分转换为星级（1-5）
	star := int(math.Ceil(float64(rating)))
	if star < 1 {
		star = 1
	}
	if star > 5 {
		star = 5
	}

	key := fmt.Sprintf("%d", star)
	rd[key]++
	return nil
}

// GetCount 获取某分数的个数
func (rd RatingDistribution) GetCount(star int) int64 {
	if star < 1 || star > 5 {
		return 0
	}
	key := fmt.Sprintf("%d", star)
	return rd[key]
}

// GetTotal 获取总评分数
func (rd RatingDistribution) GetTotal() int64 {
	var total int64
	for _, count := range rd {
		total += count
	}
	return total
}

// GetAverage 计算平均分
func (rd RatingDistribution) GetAverage() Rating {
	total := rd.GetTotal()
	if total == 0 {
		return RatingDefault
	}

	var sum float64
	for star := 1; star <= 5; star++ {
		key := fmt.Sprintf("%d", star)
		sum += float64(star) * float64(rd[key])
	}

	return Rating(sum / float64(total))
}

// GetPercentage 获取某星级的百分比
func (rd RatingDistribution) GetPercentage(star int) float64 {
	if star < 1 || star > 5 {
		return 0
	}

	total := rd.GetTotal()
	if total == 0 {
		return 0
	}

	return float64(rd.GetCount(star)) / float64(total) * 100
}

// ToBSON 转换为 BSON 兼容格式
func (rd RatingDistribution) ToBSON() map[string]int64 {
	return rd
}

// FromBSON 从 BSON 创建分布
func FromBSON(data map[string]int64) RatingDistribution {
	if data == nil {
		return NewRatingDistribution()
	}

	// 确保所有星级都存在
	result := NewRatingDistribution()
	for k, v := range data {
		result[k] = v
	}
	return result
}

// Clone 克隆分布
func (rd RatingDistribution) Clone() RatingDistribution {
	result := make(RatingDistribution, len(rd))
	for k, v := range rd {
		result[k] = v
	}
	return result
}

// Reset 重置所有计数
func (rd RatingDistribution) Reset() {
	for star := 1; star <= 5; star++ {
		key := fmt.Sprintf("%d", star)
		rd[key] = 0
	}
}

// GetMostCommonStar 获取最常见的星级
func (rd RatingDistribution) GetMostCommonStar() int {
	var maxStar int
	var maxCount int64

	for star := 1; star <= 5; star++ {
		count := rd.GetCount(star)
		if count > maxCount {
			maxCount = count
			maxStar = star
		}
	}

	return maxStar
}

// GetRatingSummary 获取评分摘要
func (rd RatingDistribution) GetRatingSummary() map[string]interface{} {
	return map[string]interface{}{
		"average":    rd.GetAverage().ToFloat(),
		"total":      rd.GetTotal(),
		"distribution": map[string]int64(rd),
		"most_common": rd.GetMostCommonStar(),
	}
}
