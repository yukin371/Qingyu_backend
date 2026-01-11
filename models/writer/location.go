package writer

import (
	"Qingyu_backend/models/writer/base"
)

// Location 地点/场景
type Location struct {
	base.IdentifiedEntity    `bson:",inline"` // ID
	base.Timestamps          `bson:",inline"` // 时间戳
	base.ProjectScopedEntity `bson:",inline"` // ProjectID
	base.NamedEntity         `bson:",inline"` // Name

	// 详细信息（保持BSON字段名不变，确保数据库兼容）
	Description string `bson:"description,omitempty" json:"description,omitempty" validate:"max=1000"`
	Climate     string `bson:"climate,omitempty" json:"climate,omitempty" validate:"max=100"`
	Culture     string `bson:"culture,omitempty" json:"culture,omitempty" validate:"max=200"`
	Geography   string `bson:"geography,omitempty" json:"geography,omitempty" validate:"max=200"`
	Atmosphere  string `bson:"atmosphere,omitempty" json:"atmosphere,omitempty" validate:"max=200"`
	ParentID    string `bson:"parent_id,omitempty" json:"parentId,omitempty"` // 父级地点ID，支持层级结构
	ImageURL    string `bson:"image_url,omitempty" json:"imageUrl,omitempty" validate:"omitempty,url"`
}

// TouchForCreate 创建时设置默认值
func (l *Location) TouchForCreate() {
	l.IdentifiedEntity.GenerateID()
	l.Timestamps.TouchForCreate()
}

// Validate 验证地点数据
func (l *Location) Validate() error {
	if err := base.ValidateName(l.Name, 100); err != nil {
		return err
	}
	if l.ProjectID == "" {
		return base.ErrProjectIDRequired
	}
	if err := base.ValidateURL(l.ImageURL); err != nil {
		return err
	}
	return nil
}

// LocationRelation 地点关系（如：相邻、包含等）
type LocationRelation struct {
	base.IdentifiedEntity    `bson:",inline"`
	base.Timestamps          `bson:",inline"`
	base.ProjectScopedEntity `bson:",inline"`

	FromID   string                 `bson:"from_id" json:"fromId" validate:"required"`
	ToID     string                 `bson:"to_id" json:"toId" validate:"required"`
	Type     LocationRelationType   `bson:"type" json:"type" validate:"required"`
	Distance string                 `bson:"distance,omitempty" json:"distance,omitempty" validate:"max=100"` // 距离描述
	Notes    string                 `bson:"notes,omitempty" json:"notes,omitempty" validate:"max=500"`
}

// TouchForCreate 创建时设置默认值
func (lr *LocationRelation) TouchForCreate() {
	lr.IdentifiedEntity.GenerateID()
	lr.Timestamps.TouchForCreate()
}

// Validate 验证地点关系
func (lr *LocationRelation) Validate() error {
	if lr.ProjectID == "" {
		return base.ErrProjectIDRequired
	}
	if !lr.Type.IsValid() {
		return base.ErrInvalidEnum
	}
	return nil
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

// IsValid 验证地点关系类型是否合法
func (lrt LocationRelationType) IsValid() bool {
	switch lrt {
	case LocationRelationAdjacent, LocationRelationContains, LocationRelationNear, LocationRelationFar, LocationRelationConnected:
		return true
	default:
		return false
	}
}

// IsValidLocationRelationType 验证地点关系类型是否合法（保持向后兼容）
func IsValidLocationRelationType(t string) bool {
	return LocationRelationType(t).IsValid()
}
