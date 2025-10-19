package stats

import "time"

// ChapterStats 章节统计模型
type ChapterStats struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	BookID    string    `bson:"book_id" json:"book_id"`                       // 作品ID
	ChapterID string    `bson:"chapter_id" json:"chapter_id"`                 // 章节ID
	Title     string    `bson:"title" json:"title"`                           // 章节标题
	WordCount int       `bson:"word_count" json:"word_count"`                 // 字数
	
	// 阅读数据
	ViewCount      int64   `bson:"view_count" json:"view_count"`               // 阅读次数
	UniqueViewers  int64   `bson:"unique_viewers" json:"unique_viewers"`       // 独立读者数
	AvgReadTime    float64 `bson:"avg_read_time" json:"avg_read_time"`         // 平均阅读时长(秒)
	CompletionRate float64 `bson:"completion_rate" json:"completion_rate"`     // 完读率(0-1)
	
	// 跳出数据
	DropOffCount int64   `bson:"drop_off_count" json:"drop_off_count"`       // 跳出次数
	DropOffRate  float64 `bson:"drop_off_rate" json:"drop_off_rate"`         // 跳出率(0-1)
	
	// 互动数据
	CommentCount  int64 `bson:"comment_count" json:"comment_count"`         // 评论数
	LikeCount     int64 `bson:"like_count" json:"like_count"`               // 点赞数
	BookmarkCount int64 `bson:"bookmark_count" json:"bookmark_count"`       // 书签数
	
	// 订阅数据（付费章节）
	SubscribeCount int64   `bson:"subscribe_count" json:"subscribe_count"`     // 订阅数
	Revenue        float64 `bson:"revenue" json:"revenue"`                     // 收入（元）
	
	// 时间戳
	StatDate  time.Time `bson:"stat_date" json:"stat_date"`     // 统计日期
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// ChapterStatsAggregate 章节统计聚合结果
type ChapterStatsAggregate struct {
	ChapterID      string  `json:"chapter_id"`
	Title          string  `json:"title"`
	ViewCount      int64   `json:"view_count"`
	UniqueViewers  int64   `json:"unique_viewers"`
	CompletionRate float64 `json:"completion_rate"`
	DropOffRate    float64 `json:"drop_off_rate"`
	Revenue        float64 `json:"revenue"`
}

// HeatmapPoint 热力图数据点
type HeatmapPoint struct {
	ChapterNum     int     `json:"chapter_num"`      // 章节序号
	ChapterID      string  `json:"chapter_id"`       // 章节ID
	ViewCount      int64   `json:"view_count"`       // 阅读量
	CompletionRate float64 `json:"completion_rate"`  // 完读率
	DropOffRate    float64 `json:"drop_off_rate"`    // 跳出率
	HeatScore      float64 `json:"heat_score"`       // 热度分数(0-100)
}

// TimeRangeStats 时间范围统计
type TimeRangeStats struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	
	TotalViews         int64   `json:"total_views"`
	TotalUniqueViewers int64   `json:"total_unique_viewers"`
	AvgCompletionRate  float64 `json:"avg_completion_rate"`
	AvgDropOffRate     float64 `json:"avg_drop_off_rate"`
	TotalRevenue       float64 `json:"total_revenue"`
}

