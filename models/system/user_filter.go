package system

// UserFilter 用户查询过滤条件
type UserFilter struct {
	ID       string `bson:"id,omitempty" json:"id,omitempty"`
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	Email    string `bson:"email,omitempty" json:"email,omitempty"`
}
