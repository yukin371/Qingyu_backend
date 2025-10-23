package bookstore

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RankingType 榜单类型
type RankingType string

const (
	RankingTypeRealtime RankingType = "realtime" // 实时榜
	RankingTypeWeekly   RankingType = "weekly"   // 周榜
	RankingTypeMonthly  RankingType = "monthly"  // 月榜
	RankingTypeNewbie   RankingType = "newbie"   // 新人榜
)

// RankingItem 榜单项目
type RankingItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookID    primitive.ObjectID `bson:"book_id" json:"bookId" validate:"required"`  // 书籍ID
	Book      *Book              `bson:"book,omitempty" json:"book,omitempty"`       // 书籍信息（查询时填充）
	Type      RankingType        `bson:"type" json:"type" validate:"required"`       // 榜单类型
	Rank      int                `bson:"rank" json:"rank" validate:"required,min=1"` // 排名
	Score     float64            `bson:"score" json:"score"`                         // 评分
	ViewCount int64              `bson:"view_count" json:"viewCount"`                // 浏览量
	LikeCount int64              `bson:"like_count" json:"likeCount"`                // 点赞数
	Period    string             `bson:"period" json:"period" validate:"required"`   // 统计周期 (格式: 2024-01-01 或 2024-W01)
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`                // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`                // 更新时间
}

// RankingFilter 榜单查询过滤器
type RankingFilter struct {
	Type     *RankingType `json:"type,omitempty"`
	Period   *string      `json:"period,omitempty"`
	Limit    int          `json:"limit,omitempty"`
	Offset   int          `json:"offset,omitempty"`
	WithBook bool         `json:"withBook,omitempty"` // 是否包含书籍详情
}

// RankingStats 榜单统计信息
type RankingStats struct {
	Type   RankingType `json:"type"`
	Period string      `json:"period"`
	TotalBooks    int64       `json:"totalBooks"`
	TotalViews    int64       `json:"totalViews"`
	TotalLikes    int64       `json:"totalLikes"`
	AverageScore  float64     `json:"averageScore"`
	LastUpdatedAt time.Time   `json:"lastUpdatedAt"`
}

// RankingResponse 榜单响应结构
type RankingResponse struct {
	Type      RankingType    `json:"type"`
	Period    string         `json:"period"`
	Items     []*RankingItem `json:"items"`
	Stats     *RankingStats  `json:"stats,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

// GetPeriodString 获取周期字符串
func GetPeriodString(t RankingType, date time.Time) string {
	switch t {
	case RankingTypeRealtime:
		return date.Format("2006-01-02")
	case RankingTypeWeekly:
		year, week := date.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case RankingTypeMonthly:
		return date.Format("2006-01")
	case RankingTypeNewbie:
		return date.Format("2006-01")
	default:
		return date.Format("2006-01-02")
	}
}

// IsValidRankingType 验证榜单类型是否有效
func IsValidRankingType(t string) bool {
	switch RankingType(t) {
	case RankingTypeRealtime, RankingTypeWeekly, RankingTypeMonthly, RankingTypeNewbie:
		return true
	default:
		return false
	}
}
