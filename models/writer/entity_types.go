package writer

// EntityType 实体类型枚举
type EntityType string

const (
	EntityTypeCharacter     EntityType = "character"
	EntityTypeItem          EntityType = "item"
	EntityTypeLocation      EntityType = "location"
	EntityTypeOrganization  EntityType = "organization"
	EntityTypeForeshadowing EntityType = "foreshadowing"
)

// IsValid 验证实体类型是否合法
func (et EntityType) IsValid() bool {
	switch et {
	case EntityTypeCharacter, EntityTypeItem, EntityTypeLocation,
		EntityTypeOrganization, EntityTypeForeshadowing:
		return true
	default:
		return false
	}
}

// StateValue 结构化状态值
type StateValue struct {
	Current     any         `bson:"current" json:"current"`
	Min         *float64    `bson:"min,omitempty" json:"min,omitempty"`
	Max         *float64    `bson:"max,omitempty" json:"max,omitempty"`
	Unit        string      `bson:"unit,omitempty" json:"unit,omitempty"`
	Description string      `bson:"description,omitempty" json:"description,omitempty"`
}
