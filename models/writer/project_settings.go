package writer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewProjectSettings 创建新项目设置（含默认角色类型）
// 注意：此函数返回的 ProjectSettings 结构体应嵌入到 Project 中使用
func NewProjectSettings(projectID primitive.ObjectID) ProjectSettings {
	return ProjectSettings{
		AutoBackup:     false,
		BackupInterval: 24,
		WordCountGoal:  0,
		CharacterRoles: GetDefaultCharacterRoles(projectID),
	}
}

// GetRoleByName 根据名称获取角色类型
func (ps *ProjectSettings) GetRoleByName(name string) *CharacterRole {
	for i := range ps.CharacterRoles {
		if ps.CharacterRoles[i].Name == name {
			return &ps.CharacterRoles[i]
		}
	}
	return nil
}

// GetRoleByID 根据ID获取角色类型
func (ps *ProjectSettings) GetRoleByID(id string) *CharacterRole {
	for i := range ps.CharacterRoles {
		if ps.CharacterRoles[i].ID.Hex() == id {
			return &ps.CharacterRoles[i]
		}
	}
	return nil
}
