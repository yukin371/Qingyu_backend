package writer

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetDefaultCharacterRoles(t *testing.T) {
	projectID := primitive.NewObjectID()
	roles := GetDefaultCharacterRoles(projectID)

	if len(roles) != 3 {
		t.Errorf("期望3个默认角色，得到%d个", len(roles))
	}

	// 验证主角
	if roles[0].Name != RoleProtagonist {
		t.Errorf("第1个角色应该是%q，得到%q", RoleProtagonist, roles[0].Name)
	}
	if !roles[0].IsDefault {
		t.Error("预设角色应该标记为IsDefault=true")
	}

	// 验证排序
	if roles[0].Order != 1 || roles[1].Order != 2 || roles[2].Order != 3 {
		t.Error("预设角色应该按1,2,3排序")
	}
}

func TestNewCharacterRole(t *testing.T) {
	projectID := primitive.NewObjectID()
	role := NewCharacterRole(projectID, "反派", "#ff0000", "villain", 4)

	if role.ProjectID != projectID {
		t.Error("角色ProjectID应该匹配")
	}
	if role.Name != "反派" {
		t.Errorf("角色名称应该是'反派'，得到%q", role.Name)
	}
	if role.IsDefault {
		t.Error("新建的自定义角色应该标记为IsDefault=false")
	}
	if role.Order != 4 {
		t.Errorf("角色Order应该是4，得到%d", role.Order)
	}
}

func TestProjectSettings_GetRoleByName(t *testing.T) {
	projectID := primitive.NewObjectID()
	settings := NewProjectSettings(projectID)

	// 测试查找存在的角色
	role := settings.GetRoleByName(RoleProtagonist)
	if role == nil {
		t.Error("应该找到主角角色")
	}
	if role.Name != RoleProtagonist {
		t.Errorf("找到的角色应该是%q", RoleProtagonist)
	}

	// 测试查找不存在的角色
	role = settings.GetRoleByName("不存在的角色")
	if role != nil {
		t.Error("不应该找到不存在的角色")
	}
}

func TestProjectSettings_GetRoleByID(t *testing.T) {
	projectID := primitive.NewObjectID()
	settings := NewProjectSettings(projectID)

	// 测试查找存在的角色
	targetRole := settings.CharacterRoles[0]
	role := settings.GetRoleByID(targetRole.ID.Hex())
	if role == nil {
		t.Error("应该找到角色")
	}
	if role.ID != targetRole.ID {
		t.Error("找到的角色ID应该匹配")
	}

	// 测试查找不存在的角色
	role = settings.GetRoleByID("ffffffffffffffffffffffff")
	if role != nil {
		t.Error("不应该找到不存在的角色ID")
	}
}
