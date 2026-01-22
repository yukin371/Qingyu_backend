package stats

import "time"

// ReaderBehavior 读者行为数据模型
type ReaderBehavior struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	UserID    string `bson:"user_id" json:"$1$2"`       // 用户ID
	BookID    string `bson:"book_id" json:"$1$2"`       // 作品ID
	ChapterID string `bson:"chapter_id" json:"$1$2"` // 章节ID

	// 行为类型
	BehaviorType string `bson:"behavior_type" json:"$1$2"` // view/complete/drop_off/subscribe

	// 阅读进度
	StartPosition int     `bson:"start_position" json:"$1$2"` // 开始位置(字符数)
	EndPosition   int     `bson:"end_position" json:"$1$2"`     // 结束位置(字符数)
	Progress      float64 `bson:"progress" json:"progress"`             // 阅读进度(0-1)

	// 时间数据
	ReadDuration int       `bson:"read_duration" json:"$1$2"` // 阅读时长(秒)
	ReadAt       time.Time `bson:"read_at" json:"$1$2"`             // 阅读时间

	// 设备信息
	DeviceType string `bson:"device_type" json:"$1$2"` // mobile/desktop/tablet
	ClientIP   string `bson:"client_ip" json:"$1$2"`     // IP地址

	// 来源信息
	Source   string `bson:"source" json:"source"`     // 来源(推荐/搜索/书架)
	Referrer string `bson:"referrer" json:"referrer"` // 引荐页面

	CreatedAt time.Time `bson:"created_at" json:"$1$2"`
}

// ReadingSession 阅读会话
type ReadingSession struct {
	SessionID string `bson:"session_id" json:"$1$2"`
	UserID    string `bson:"user_id" json:"$1$2"`
	BookID    string `bson:"book_id" json:"$1$2"`

	StartChapter string `bson:"start_chapter" json:"$1$2"` // 开始章节
	EndChapter   string `bson:"end_chapter" json:"$1$2"`     // 结束章节
	ChaptersRead int    `bson:"chapters_read" json:"$1$2"` // 已读章节数

	TotalDuration int       `bson:"total_duration" json:"$1$2"` // 总时长(秒)
	StartTime     time.Time `bson:"start_time" json:"$1$2"`
	EndTime       time.Time `bson:"end_time" json:"$1$2"`
}

// ReaderRetention 读者留存数据
type ReaderRetention struct {
	BookID string `bson:"book_id" json:"$1$2"`

	Day1Retention  float64 `bson:"day1_retention" json:"$1$2"`   // 次日留存率
	Day3Retention  float64 `bson:"day3_retention" json:"$1$2"`   // 3日留存率
	Day7Retention  float64 `bson:"day7_retention" json:"$1$2"`   // 7日留存率
	Day30Retention float64 `bson:"day30_retention" json:"$1$2"` // 30日留存率

	NewReaders    int64 `bson:"new_readers" json:"$1$2"`       // 新增读者
	ActiveReaders int64 `bson:"active_readers" json:"$1$2"` // 活跃读者

	StatDate  time.Time `bson:"stat_date" json:"$1$2"`
	CreatedAt time.Time `bson:"created_at" json:"$1$2"`
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
