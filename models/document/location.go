package document

import "time"

// Location 地点/场景
type Location struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	ProjectID   string    `bson:"project_id" json:"projectId"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Climate     string    `bson:"climate,omitempty" json:"climate,omitempty"`
	Culture     string    `bson:"culture,omitempty" json:"culture,omitempty"`
	Geography   string    `bson:"geography,omitempty" json:"geography,omitempty"`
	Atmosphere  string    `bson:"atmosphere,omitempty" json:"atmosphere,omitempty"`
	ParentID    string    `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父级地点ID，支持层级结构
	ImageURL    string    `bson:"image_url,omitempty" json:"imageUrl,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt"`
}

// LocationRelation 地点关系（如：相邻、包含等）
type LocationRelation struct {
	ID        string               `bson:"_id,omitempty" json:"id"`
	ProjectID string               `bson:"project_id" json:"projectId"`
	FromID    string               `bson:"from_id" json:"fromId"`
	ToID      string               `bson:"to_id" json:"toId"`
	Type      LocationRelationType `bson:"type" json:"type"`
	Distance  string               `bson:"distance,omitempty" json:"distance,omitempty"` // 距离描述
	Notes     string               `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updatedAt"`
}

// LocationRelationType 地点关系类型
type LocationRelationType string

const (
	LocationRelationAdjacent  LocationRelationType = "adjacent"  // 相邻
	LocationRelationContains  LocationRelationType = "contains"  // 包含
	LocationRelationNear      LocationRelationType = "near"      // 附近
	LocationRelationFar       LocationRelationType = "far"       // 远离
	LocationRelationConnected LocationRelationType = "connected" // 连通
)

// IsValidLocationRelationType 验证地点关系类型是否合法
func IsValidLocationRelationType(t string) bool {
	switch LocationRelationType(t) {
	case LocationRelationAdjacent, LocationRelationContains, LocationRelationNear, LocationRelationFar, LocationRelationConnected:
		return true
	default:
		return false
	}
}
