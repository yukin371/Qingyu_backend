package stats

import "time"

// ChapterStats 章节统计模型
type ChapterStats struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	BookID    string `bson:"book_id" json:"$1$2"`       // 作品ID
	ChapterID string `bson:"chapter_id" json:"$1$2"` // 章节ID
	Title     string `bson:"title" json:"title"`           // 章节标题
	WordCount int    `bson:"word_count" json:"$1$2"` // 字数

	// 阅读数据
	ViewCount      int64   `bson:"view_count" json:"$1$2"`           // 阅读次数
	UniqueViewers  int64   `bson:"unique_viewers" json:"$1$2"`   // 独立读者数
	AvgReadTime    float64 `bson:"avg_read_time" json:"$1$2"`     // 平均阅读时长(秒)
	CompletionRate float64 `bson:"completion_rate" json:"$1$2"` // 完读率(0-1)

	// 跳出数据
	DropOffCount int64   `bson:"drop_off_count" json:"$1$2"` // 跳出次数
	DropOffRate  float64 `bson:"drop_off_rate" json:"$1$2"`   // 跳出率(0-1)

	// 互动数据
	CommentCount  int64 `bson:"comment_count" json:"$1$2"`   // 评论数
	LikeCount     int64 `bson:"like_count" json:"$1$2"`         // 点赞数
	BookmarkCount int64 `bson:"bookmark_count" json:"$1$2"` // 书签数

	// 订阅数据（付费章节）
	SubscribeCount int64   `bson:"subscribe_count" json:"$1$2"` // 订阅数
	Revenue        float64 `bson:"revenue" json:"revenue"`                 // 收入（元）

	// 时间戳
	StatDate  time.Time `bson:"stat_date" json:"$1$2"` // 统计日期
	CreatedAt time.Time `bson:"created_at" json:"$1$2"`
	UpdatedAt time.Time `bson:"updated_at" json:"$1$2"`
}

// ChapterStatsAggregate 章节统计聚合结果
type ChapterStatsAggregate struct {
	ChapterID      string  `json:"$1$2"`
	Title          string  `json:"title"`
	ViewCount      int64   `json:"$1$2"`
	UniqueViewers  int64   `json:"$1$2"`
	CompletionRate float64 `json:"$1$2"`
	DropOffRate    float64 `json:"$1$2"`
	Revenue        float64 `json:"revenue"`
}

// HeatmapPoint 热力图数据点
type HeatmapPoint struct {
	ChapterNum     int     `json:"$1$2"`     // 章节序号
	ChapterID      string  `json:"$1$2"`      // 章节ID
	ViewCount      int64   `json:"$1$2"`      // 阅读量
	CompletionRate float64 `json:"$1$2"` // 完读率
	DropOffRate    float64 `json:"$1$2"`   // 跳出率
	HeatScore      float64 `json:"$1$2"`      // 热度分数(0-100)
}

// TimeRangeStats 时间范围统计
type TimeRangeStats struct {
	StartDate time.Time `json:"$1$2"`
	EndDate   time.Time `json:"$1$2"`

	TotalViews         int64   `json:"$1$2"`
	TotalUniqueViewers int64   `json:"$1$2"`
	AvgCompletionRate  float64 `json:"$1$2"`
	AvgDropOffRate     float64 `json:"$1$2"`
	TotalRevenue       float64 `json:"$1$2"`
}
