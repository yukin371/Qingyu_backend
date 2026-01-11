package shared

// TitledEntity 标题实体混入
type TitledEntity struct {
	Title string `bson:"title" json:"title" validate:"required,min=1,max=200"`
}

// NamedEntity 命名实体混入
type NamedEntity struct {
	Name string `bson:"name" json:"name" validate:"required,min=1,max=100"`
}

// DescriptedEntity 描述实体混入
type DescriptedEntity struct {
	Description string `bson:"description,omitempty" json:"description,omitempty" validate:"max=1000"`
}
