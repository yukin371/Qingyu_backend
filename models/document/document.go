package document

import "time"

// Document 表示系统中的文档数据模型
type Document struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	UserID    string    `bson:"user_id" json:"userId"`
	Title     string    `bson:"title" json:"title"`
	Content   string    `bson:"content" json:"content"`
	Tags      []string  `bson:"tags" json:"tags"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

// TouchForCreate 在创建前设置时间戳
func (d *Document) TouchForCreate() {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
}

// TouchForUpdate 在更新前刷新更新时间戳
func (d *Document) TouchForUpdate() {
	d.UpdatedAt = time.Now()
}
