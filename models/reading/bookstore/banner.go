package bookstore

type Banner struct {
	ID     string `bson:"_id,omitempty" json:"id"`
	Title  string `bson:"title" json:"title"`   // 标题
	Image  string `bson:"image" json:"image"`   // 图片
	Target string `bson:"target" json:"target"` // 跳转目标
}
