package reco

import "time"

// ItemFeature 物品特征（书籍特征）
// 用于内容推荐与相似度计算
type ItemFeature struct {
	ID         string             `bson:"_id,omitempty" json:"id"`
	ItemID     string             `bson:"item_id" json:"itemId"` // 对应书籍ID
	Tags       map[string]float64 `bson:"tags" json:"tags"`
	Authors    []string           `bson:"authors" json:"authors"`
	Categories []string           `bson:"categories" json:"categories"`
	Vector     []float64          `bson:"vector,omitempty" json:"vector,omitempty"` // 预留向量字段（可选）
	UpdatedAt  time.Time          `bson:"updated_at" json:"updatedAt"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
}
