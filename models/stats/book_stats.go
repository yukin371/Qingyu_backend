package stats

import "time"

// BookStats 作品统计模型
type BookStats struct {
	ID     string `bson:"_id,omitempty" json:"id"`
	BookID string `bson:"book_id" json:"book_id"` // 作品ID
	
	// 基础信息
	Title        string `bson:"title" json:"title"`
	AuthorID     string `bson:"author_id" json:"author_id"`
	TotalChapter int    `bson:"total_chapter" json:"total_chapter"`     // 总章节数
	TotalWords   int64  `bson:"total_words" json:"total_words"`         // 总字数
	
	// 阅读数据
	TotalViews         int64   `bson:"total_views" json:"total_views"`                   // 总阅读量
	UniqueReaders      int64   `bson:"unique_readers" json:"unique_readers"`             // 独立读者数
	AvgChapterViews    float64 `bson:"avg_chapter_views" json:"avg_chapter_views"`       // 平均章节阅读量
	AvgCompletionRate  float64 `bson:"avg_completion_rate" json:"avg_completion_rate"`   // 平均完读率
	AvgReadingDuration float64 `bson:"avg_reading_duration" json:"avg_reading_duration"` // 平均阅读时长(秒)
	
	// 跳出数据
	TotalDropOffs  int64   `bson:"total_drop_offs" json:"total_drop_offs"`   // 总跳出次数
	AvgDropOffRate float64 `bson:"avg_drop_off_rate" json:"avg_drop_off_rate"` // 平均跳出率
	DropOffChapter string  `bson:"drop_off_chapter" json:"drop_off_chapter"` // 最高跳出章节
	
	// 互动数据
	TotalComments  int64 `bson:"total_comments" json:"total_comments"`   // 总评论数
	TotalLikes     int64 `bson:"total_likes" json:"total_likes"`         // 总点赞数
	TotalBookmarks int64 `bson:"total_bookmarks" json:"total_bookmarks"` // 总书签数
	TotalShares    int64 `bson:"total_shares" json:"total_shares"`       // 总分享数
	
	// 订阅数据
	TotalSubscribers int64   `bson:"total_subscribers" json:"total_subscribers"` // 总订阅数
	AvgSubscribeRate float64 `bson:"avg_subscribe_rate" json:"avg_subscribe_rate"` // 平均订阅率
	
	// 收入数据
	TotalRevenue      float64 `bson:"total_revenue" json:"total_revenue"`             // 总收入
	ChapterRevenue    float64 `bson:"chapter_revenue" json:"chapter_revenue"`         // 章节收入
	SubscribeRevenue  float64 `bson:"subscribe_revenue" json:"subscribe_revenue"`     // 订阅收入
	RewardRevenue     float64 `bson:"reward_revenue" json:"reward_revenue"`           // 打赏收入
	AvgRevenuePerUser float64 `bson:"avg_revenue_per_user" json:"avg_revenue_per_user"` // 平均用户贡献
	
	// 留存数据
	Day1Retention  float64 `bson:"day1_retention" json:"day1_retention"`   // 次日留存
	Day7Retention  float64 `bson:"day7_retention" json:"day7_retention"`   // 7日留存
	Day30Retention float64 `bson:"day30_retention" json:"day30_retention"` // 30日留存
	
	// 趋势数据
	ViewTrend     string `bson:"view_trend" json:"view_trend"`         // 阅读量趋势(up/down/stable)
	RevenueTrend  string `bson:"revenue_trend" json:"revenue_trend"`   // 收入趋势
	
	// 时间戳
	StatDate  time.Time `bson:"stat_date" json:"stat_date"`     // 统计日期
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// BookStatsDaily 作品每日统计
type BookStatsDaily struct {
	ID     string    `bson:"_id,omitempty" json:"id"`
	BookID string    `bson:"book_id" json:"book_id"`
	Date   time.Time `bson:"date" json:"date"` // 统计日期
	
	DailyViews       int64   `bson:"daily_views" json:"daily_views"`             // 当日阅读量
	DailyNewReaders  int64   `bson:"daily_new_readers" json:"daily_new_readers"` // 当日新增读者
	DailyRevenue     float64 `bson:"daily_revenue" json:"daily_revenue"`         // 当日收入
	DailySubscribers int64   `bson:"daily_subscribers" json:"daily_subscribers"` // 当日订阅数
	
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// RevenueBreakdown 收入细分
type RevenueBreakdown struct {
	BookID string `json:"book_id"`
	
	ChapterRevenue   float64 `json:"chapter_revenue"`    // 章节付费
	SubscribeRevenue float64 `json:"subscribe_revenue"`  // 订阅收入
	RewardRevenue    float64 `json:"reward_revenue"`     // 打赏收入
	AdRevenue        float64 `json:"ad_revenue"`         // 广告收入
	
	TotalRevenue float64 `json:"total_revenue"`
	
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// TopChapters 热门章节
type TopChapters struct {
	BookID string `json:"book_id"`
	
	MostViewed      []ChapterStatsAggregate `json:"most_viewed"`       // 阅读量最高
	HighestRevenue  []ChapterStatsAggregate `json:"highest_revenue"`   // 收入最高
	LowestCompletion []ChapterStatsAggregate `json:"lowest_completion"` // 完读率最低
	HighestDropOff  []ChapterStatsAggregate `json:"highest_drop_off"`  // 跳出率最高
}

// Trend 趋势常量
const (
	TrendUp     = "up"
	TrendDown   = "down"
	TrendStable = "stable"
)

