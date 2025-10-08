package reader

import "time"

type ReadingSettings struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	UserID      string    `bson:"user_id" json:"userId"`           // 用户ID
	FontFamily  string    `bson:"font_family" json:"fontFamily"`   // 字体
	FontSize    int       `bson:"font_size" json:"fontSize"`       // 字号
	LineHeight  float64   `bson:"line_height" json:"lineHeight"`   // 行高
	Theme       string    `bson:"theme" json:"theme"`              // 主题
	Background  string    `bson:"background" json:"background"`    // 背景色
	PageMode    int       `bson:"page_mode" json:"pageMode"`       // 翻页模式：1-滑动，2-仿真
	AutoScroll  bool      `bson:"auto_scroll" json:"autoScroll"`   // 自动滚动
	ScrollSpeed int       `bson:"scroll_speed" json:"scrollSpeed"` // 滚动速度
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}
