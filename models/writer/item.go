package writer

import "time"

// Item 物品/道具
type Item struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	ProjectID   string `bson:"project_id" json:"projectId"`
	Name        string `bson:"name" json:"name"`
	Type        string `bson:"type" json:"type"` // 武器、消耗品、关键道具、货币
	Description string `bson:"description" json:"description"`

	// 归属
	OwnerID    string `bson:"owner_id,omitempty" json:"ownerId,omitempty"`       // 当前持有者(角色ID)
	LocationID string `bson:"location_id,omitempty" json:"locationId,omitempty"` // 当前所在地(若无持有者)

	// 属性
	Rarity   string `bson:"rarity,omitempty" json:"rarity,omitempty"`     // 稀有度
	Function string `bson:"function,omitempty" json:"function,omitempty"` // 功能/用途
	Origin   string `bson:"origin,omitempty" json:"origin,omitempty"`     // 来源/出处

	ImageURL  string    `bson:"image_url,omitempty" json:"imageUrl,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}
