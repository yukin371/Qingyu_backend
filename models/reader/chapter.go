package reader

import "time"

type Chapter struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	BookID      string    `bson:"book_id" json:"bookId"`           // 书籍ID
	Title       string    `bson:"title" json:"title"`              // 章节标题
	Content     string    `bson:"content" json:"content"`          // 章节内容
	WordCount   int       `bson:"word_count" json:"wordCount"`     // 字数
	ChapterNum  int       `bson:"chapter_num" json:"chapterNum"`   // 章节序号
	IsVIP       bool      `bson:"is_vip" json:"isVip"`             // 是否VIP章节
	Price       int64     `bson:"price" json:"price"`              // 价格（分）
	Status      int       `bson:"status" json:"status"`            // 状态：1-正常，2-删除
	PublishTime time.Time `bson:"publish_time" json:"publishTime"` // 发布时间
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`     // 创建时间
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`     // 更新时间
}
