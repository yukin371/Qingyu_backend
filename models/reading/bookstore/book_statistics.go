package bookstore

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookStatistics 书籍统计模型
type BookStatistics struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID             primitive.ObjectID `bson:"book_id" json:"book_id"`
	ViewCount          int64              `bson:"view_count" json:"view_count"`
	FavoriteCount      int64              `bson:"favorite_count" json:"favorite_count"`
	CommentCount       int64              `bson:"comment_count" json:"comment_count"`
	ShareCount         int64              `bson:"share_count" json:"share_count"`
	AverageRating      float64            `bson:"average_rating" json:"average_rating"`
	RatingCount        int64              `bson:"rating_count" json:"rating_count"`
	RatingDistribution map[int]int64      `bson:"rating_distribution" json:"rating_distribution"`
	HotScore           float64            `bson:"hot_score" json:"hot_score"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
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
		bs.RatingDistribution = make(map[int]int64)
	}

	// 增加对应星级的计数
	bs.RatingDistribution[rating]++
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

	if bs.RatingDistribution[rating] > 0 {
		bs.RatingDistribution[rating]--
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
	for rating, count := range bs.RatingDistribution {
		totalScore += int64(rating) * count
	}

	bs.AverageRating = float64(totalScore) / float64(bs.RatingCount)
}

// CalculateHotScore 计算热度分数
func (bs *BookStatistics) CalculateHotScore() {
	// 热度分数计算公式：浏览数*0.1 + 收藏数*0.5 + 评论数*0.3 + 分享数*0.2 + 平均评分*10
	bs.HotScore = float64(bs.ViewCount)*0.1 +
		float64(bs.FavoriteCount)*0.5 +
		float64(bs.CommentCount)*0.3 +
		float64(bs.ShareCount)*0.2 +
		bs.AverageRating*10
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
