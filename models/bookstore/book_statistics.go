package bookstore

import (
	"Qingyu_backend/models/shared/types"
	"strconv"
	"time"
)

// BookStatistics 书籍统计模型
type BookStatistics struct {
	ID                 string                   `bson:"_id,omitempty" json:"id"`
	BookID             string                   `bson:"book_id" json:"$1$2"`
	ViewCount          int64                    `bson:"view_count" json:"$1$2"`
	FavoriteCount      int64                    `bson:"favorite_count" json:"$1$2"`
	CommentCount       int64                    `bson:"comment_count" json:"$1$2"`
	ShareCount         int64                    `bson:"share_count" json:"$1$2"`
	AverageRating      types.Rating             `bson:"average_rating" json:"$1$2"`
	RatingCount        int64                    `bson:"rating_count" json:"$1$2"`
	RatingDistribution types.RatingDistribution `bson:"rating_distribution" json:"$1$2"` // 键为 "1"-"5" 的字符串
	HotScore           float64                  `bson:"hot_score" json:"$1$2"`
	UpdatedAt          time.Time                `bson:"updated_at" json:"$1$2"`
}

// BeforeUpdate 在更新前刷新更新时间戳
func (bs *BookStatistics) BeforeUpdate() {
	bs.UpdatedAt = time.Now()
}

// IncrementViewCount 增加浏览次数
func (bs *BookStatistics) IncrementViewCount() {
	bs.ViewCount++
	bs.BeforeUpdate()
}

// IncrementFavoriteCount 增加收藏次数
func (bs *BookStatistics) IncrementFavoriteCount() {
	bs.FavoriteCount++
	bs.BeforeUpdate()
}

// DecrementFavoriteCount 减少收藏次数
func (bs *BookStatistics) DecrementFavoriteCount() {
	if bs.FavoriteCount > 0 {
		bs.FavoriteCount--
		bs.BeforeUpdate()
	}
}

// IncrementCommentCount 增加评论次数
func (bs *BookStatistics) IncrementCommentCount() {
	bs.CommentCount++
	bs.BeforeUpdate()
}

// IncrementShareCount 增加分享次数
func (bs *BookStatistics) IncrementShareCount() {
	bs.ShareCount++
	bs.BeforeUpdate()
}

// UpdateRating 更新评分统计
func (bs *BookStatistics) UpdateRating(rating int) {
	if bs.RatingDistribution == nil {
		bs.RatingDistribution = make(map[string]int64)
	}

	// 增加对应星级的计数 (使用字符串键 "1"-"5")
	ratingKey := strconv.Itoa(rating)
	bs.RatingDistribution[ratingKey]++
	bs.RatingCount++

	// 重新计算平均评分
	bs.calculateAverageRating()
	bs.BeforeUpdate()
}

// RemoveRating 移除评分统计
func (bs *BookStatistics) RemoveRating(rating int) {
	if bs.RatingDistribution == nil {
		return
	}

	ratingKey := strconv.Itoa(rating)
	if bs.RatingDistribution[ratingKey] > 0 {
		bs.RatingDistribution[ratingKey]--
		if bs.RatingCount > 0 {
			bs.RatingCount--
		}

		// 重新计算平均评分
		bs.calculateAverageRating()
		bs.BeforeUpdate()
	}
}

// calculateAverageRating 计算平均评分
func (bs *BookStatistics) calculateAverageRating() {
	if bs.RatingCount == 0 {
		bs.AverageRating = 0
		return
	}

	totalScore := int64(0)
	for ratingKey, count := range bs.RatingDistribution {
		// 将字符串键转换为整数 ("1"-"5" -> 1-5)
		rating, _ := strconv.Atoi(ratingKey)
		totalScore += int64(rating) * count
	}

	bs.AverageRating = types.Rating(float64(totalScore) / float64(bs.RatingCount))
}

// CalculateHotScore 计算热度分数
func (bs *BookStatistics) CalculateHotScore() {
	// 热度分数计算公式：浏览数*0.1 + 收藏数*0.5 + 评论数*0.3 + 分享数*0.2 + 平均评分*10
	bs.HotScore = float64(bs.ViewCount)*0.1 +
		float64(bs.FavoriteCount)*0.5 +
		float64(bs.CommentCount)*0.3 +
		float64(bs.ShareCount)*0.2 +
		float64(bs.AverageRating)*10
	bs.BeforeUpdate()
}

// GetPopularityLevel 获取热门程度等级
func (bs *BookStatistics) GetPopularityLevel() string {
	if bs.HotScore >= 1000 {
		return "extremely_hot"
	} else if bs.HotScore >= 500 {
		return "very_hot"
	} else if bs.HotScore >= 100 {
		return "hot"
	} else if bs.HotScore >= 50 {
		return "warm"
	}
	return "cold"
}
