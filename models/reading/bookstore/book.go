package bookstore

type Book struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	Title        string `bson:"title" json:"title"`               // 书名
	Author       string `bson:"author" json:"author"`             // 作者
	Introduction string `bson:"introduction" json:"introduction"` // 简介
	Cover        string `bson:"cover" json:"cover"`               // 封面
}
