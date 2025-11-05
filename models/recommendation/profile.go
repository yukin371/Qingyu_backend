package recommendation

import "time"

// UserProfile 用户画像
// 聚合用户偏好（题材、标签、作者、阅读强度等）
type UserProfile struct {
	ID         string             `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"user_id" json:"userId"`
	Tags       map[string]float64 `bson:"tags" json:"tags"`             // 标签偏好权重
	Authors    map[string]float64 `bson:"authors" json:"authors"`       // 作者偏好权重
	Categories map[string]float64 `bson:"categories" json:"categories"` // 分类偏好权重
	UpdatedAt  time.Time          `bson:"updated_at" json:"updatedAt"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
}
