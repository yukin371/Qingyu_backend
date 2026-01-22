package stats

import "time"

// BookStats 作品统计模型
type BookStats struct {
	ID     string `bson:"_id,omitempty" json:"id"`
	BookID string `bson:"book_id" json:"$1$2"` // 作品ID

	// 基础信息
	Title        string `bson:"title" json:"title"`
	AuthorID     string `bson:"author_id" json:"$1$2"`
	TotalChapter int    `bson:"total_chapter" json:"$1$2"` // 总章节数
	TotalWords   int64  `bson:"total_words" json:"$1$2"`     // 总字数

	// 阅读数据
	TotalViews         int64   `bson:"total_views" json:"$1$2"`                   // 总阅读量
	UniqueReaders      int64   `bson:"unique_readers" json:"$1$2"`             // 独立读者数
	AvgChapterViews    float64 `bson:"avg_chapter_views" json:"$1$2"`       // 平均章节阅读量
	AvgCompletionRate  float64 `bson:"avg_completion_rate" json:"$1$2"`   // 平均完读率
	AvgReadingDuration float64 `bson:"avg_reading_duration" json:"$1$2"` // 平均阅读时长(秒)

	// 跳出数据
	TotalDropOffs  int64   `bson:"total_drop_offs" json:"$1$2"`     // 总跳出次数
	AvgDropOffRate float64 `bson:"avg_drop_off_rate" json:"$1$2"` // 平均跳出率
	DropOffChapter string  `bson:"drop_off_chapter" json:"$1$2"`   // 最高跳出章节

	// 互动数据
	TotalComments  int64 `bson:"total_comments" json:"$1$2"`   // 总评论数
	TotalLikes     int64 `bson:"total_likes" json:"$1$2"`         // 总点赞数
	TotalBookmarks int64 `bson:"total_bookmarks" json:"$1$2"` // 总书签数
	TotalShares    int64 `bson:"total_shares" json:"$1$2"`       // 总分享数

	// 订阅数据
	TotalSubscribers int64   `bson:"total_subscribers" json:"$1$2"`   // 总订阅数
	AvgSubscribeRate float64 `bson:"avg_subscribe_rate" json:"$1$2"` // 平均订阅率

	// 收入数据
	TotalRevenue      float64 `bson:"total_revenue" json:"$1$2"`               // 总收入
	ChapterRevenue    float64 `bson:"chapter_revenue" json:"$1$2"`           // 章节收入
	SubscribeRevenue  float64 `bson:"subscribe_revenue" json:"$1$2"`       // 订阅收入
	RewardRevenue     float64 `bson:"reward_revenue" json:"$1$2"`             // 打赏收入
	AvgRevenuePerUser float64 `bson:"avg_revenue_per_user" json:"$1$2"` // 平均用户贡献

	// 留存数据
	Day1Retention  float64 `bson:"day1_retention" json:"$1$2"`   // 次日留存
	Day7Retention  float64 `bson:"day7_retention" json:"$1$2"`   // 7日留存
	Day30Retention float64 `bson:"day30_retention" json:"$1$2"` // 30日留存

	// 趋势数据
	ViewTrend    string `bson:"view_trend" json:"$1$2"`       // 阅读量趋势(up/down/stable)
	RevenueTrend string `bson:"revenue_trend" json:"$1$2"` // 收入趋势

	// 时间戳
	StatDate  time.Time `bson:"stat_date" json:"$1$2"` // 统计日期
	CreatedAt time.Time `bson:"created_at" json:"$1$2"`
	UpdatedAt time.Time `bson:"updated_at" json:"$1$2"`
}

// BookStatsDaily 作品每日统计
type BookStatsDaily struct {
	ID     string    `bson:"_id,omitempty" json:"id"`
	BookID string    `bson:"book_id" json:"$1$2"`
	Date   time.Time `bson:"date" json:"date"` // 统计日期

	DailyViews       int64   `bson:"daily_views" json:"$1$2"`             // 当日阅读量
	DailyNewReaders  int64   `bson:"daily_new_readers" json:"$1$2"` // 当日新增读者
	DailyRevenue     float64 `bson:"daily_revenue" json:"$1$2"`         // 当日收入
	DailySubscribers int64   `bson:"daily_subscribers" json:"$1$2"` // 当日订阅数

	CreatedAt time.Time `bson:"created_at" json:"$1$2"`
	UpdatedAt time.Time `bson:"updated_at" json:"$1$2"`
}

// RevenueBreakdown 收入细分
type RevenueBreakdown struct {
	BookID string `json:"$1$2"`

	ChapterRevenue   float64 `json:"$1$2"`   // 章节付费
	SubscribeRevenue float64 `json:"$1$2"` // 订阅收入
	RewardRevenue    float64 `json:"$1$2"`    // 打赏收入
	AdRevenue        float64 `json:"$1$2"`        // 广告收入

	TotalRevenue float64 `json:"$1$2"`

	StartDate time.Time `json:"$1$2"`
	EndDate   time.Time `json:"$1$2"`
}

// TopChapters 热门章节
type TopChapters struct {
	BookID string `json:"$1$2"`

	MostViewed       []*ChapterStatsAggregate `json:"$1$2"`       // 阅读量最高
	HighestRevenue   []*ChapterStatsAggregate `json:"$1$2"`   // 收入最高
	LowestCompletion []*ChapterStatsAggregate `json:"$1$2"` // 完读率最低
	HighestDropOff   []*ChapterStatsAggregate `json:"$1$2"`  // 跳出率最高
}

// Trend 趋势常量
const (
	TrendUp     = "up"
	TrendDown   = "down"
	TrendStable = "stable"
)
