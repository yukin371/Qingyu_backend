package stats

import "time"

// ReaderBehavior 读者行为数据模型
type ReaderBehavior struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	UserID    string `bson:"user_id" json:"user_id"`       // 用户ID
	BookID    string `bson:"book_id" json:"book_id"`       // 作品ID
	ChapterID string `bson:"chapter_id" json:"chapter_id"` // 章节ID

	// 行为类型
	BehaviorType string `bson:"behavior_type" json:"behavior_type"` // view/complete/drop_off/subscribe

	// 阅读进度
	StartPosition int     `bson:"start_position" json:"start_position"` // 开始位置(字符数)
	EndPosition   int     `bson:"end_position" json:"end_position"`     // 结束位置(字符数)
	Progress      float64 `bson:"progress" json:"progress"`             // 阅读进度(0-1)

	// 时间数据
	ReadDuration int       `bson:"read_duration" json:"read_duration"` // 阅读时长(秒)
	ReadAt       time.Time `bson:"read_at" json:"read_at"`             // 阅读时间

	// 设备信息
	DeviceType string `bson:"device_type" json:"device_type"` // mobile/desktop/tablet
	ClientIP   string `bson:"client_ip" json:"client_ip"`     // IP地址

	// 来源信息
	Source   string `bson:"source" json:"source"`     // 来源(推荐/搜索/书架)
	Referrer string `bson:"referrer" json:"referrer"` // 引荐页面

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// ReadingSession 阅读会话
type ReadingSession struct {
	SessionID string `bson:"session_id" json:"session_id"`
	UserID    string `bson:"user_id" json:"user_id"`
	BookID    string `bson:"book_id" json:"book_id"`

	StartChapter string `bson:"start_chapter" json:"start_chapter"` // 开始章节
	EndChapter   string `bson:"end_chapter" json:"end_chapter"`     // 结束章节
	ChaptersRead int    `bson:"chapters_read" json:"chapters_read"` // 已读章节数

	TotalDuration int       `bson:"total_duration" json:"total_duration"` // 总时长(秒)
	StartTime     time.Time `bson:"start_time" json:"start_time"`
	EndTime       time.Time `bson:"end_time" json:"end_time"`
}

// ReaderRetention 读者留存数据
type ReaderRetention struct {
	BookID string `bson:"book_id" json:"book_id"`

	Day1Retention  float64 `bson:"day1_retention" json:"day1_retention"`   // 次日留存率
	Day3Retention  float64 `bson:"day3_retention" json:"day3_retention"`   // 3日留存率
	Day7Retention  float64 `bson:"day7_retention" json:"day7_retention"`   // 7日留存率
	Day30Retention float64 `bson:"day30_retention" json:"day30_retention"` // 30日留存率

	NewReaders    int64 `bson:"new_readers" json:"new_readers"`       // 新增读者
	ActiveReaders int64 `bson:"active_readers" json:"active_readers"` // 活跃读者

	StatDate  time.Time `bson:"stat_date" json:"stat_date"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// BehaviorType 行为类型常量
const (
	BehaviorTypeView      = "view"      // 浏览
	BehaviorTypeComplete  = "complete"  // 完读
	BehaviorTypeDropOff   = "drop_off"  // 跳出
	BehaviorTypeSubscribe = "subscribe" // 订阅
	BehaviorTypeBookmark  = "bookmark"  // 书签
	BehaviorTypeComment   = "comment"   // 评论
	BehaviorTypeLike      = "like"      // 点赞
)

// DeviceType 设备类型常量
const (
	DeviceTypeMobile  = "mobile"
	DeviceTypeDesktop = "desktop"
	DeviceTypeTablet  = "tablet"
)

// Source 来源常量
const (
	SourceRecommendation = "recommendation" // 推荐
	SourceSearch         = "search"         // 搜索
	SourceBookshelf      = "bookshelf"      // 书架
	SourceRanking        = "ranking"        // 榜单
	SourceCategory       = "category"       // 分类
)
