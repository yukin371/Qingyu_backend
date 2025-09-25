package bookstore

type Category struct {
	ID   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
}
