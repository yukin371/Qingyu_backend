package writer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CharacterRole 项目级角色类型定义
type CharacterRole struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID   primitive.ObjectID `bson:"project_id" json:"projectId"`
	Name        string             `bson:"name" json:"name"`                 // 如"主角"、"配角"
	Color       string             `bson:"color,omitempty" json:"color,omitempty"`       // 可选的显示颜色
	Icon        string             `bson:"icon,omitempty" json:"icon,omitempty"`         // 可选的图标
	Order       int                `bson:"order" json:"order"`                   // 排序权重
	IsDefault   bool               `bson:"is_default" json:"isDefault"`         // 是否系统预设
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// 预设角色类型常量
const (
	RoleProtagonist = "主角"
	RoleSupporting  = "配角"
	RoleCameo       = "龙套"
)

// GetDefaultCharacterRoles 获取系统预设的角色类型列表
func GetDefaultCharacterRoles(projectID primitive.ObjectID) []CharacterRole {
	now := time.Now()
	return []CharacterRole{
		{
			ProjectID: projectID,
			Name:      RoleProtagonist,
			Color:     "#ff6b6b",
			Icon:      "star",
			Order:     1,
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ProjectID: projectID,
			Name:      RoleSupporting,
			Color:     "#4ecdc4",
			Icon:      "user",
			Order:     2,
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ProjectID: projectID,
			Name:      RoleCameo,
			Color:     "#95e1d3",
			Icon:      "users",
			Order:     3,
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// NewCharacterRole 创建自定义角色类型
func NewCharacterRole(projectID primitive.ObjectID, name, color, icon string, order int) CharacterRole {
	now := time.Now()
	return CharacterRole{
		ID:        primitive.NewObjectID(),
		ProjectID: projectID,
		Name:      name,
		Color:     color,
		Icon:      icon,
		Order:     order,
		IsDefault: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
