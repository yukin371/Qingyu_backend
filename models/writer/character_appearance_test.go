package writer

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewCharacterAppearance(t *testing.T) {
	documentID := primitive.NewObjectID()
	characterID := primitive.NewObjectID()
	roleID := primitive.NewObjectID()

	appearance := NewCharacterAppearance(
		documentID,
		characterID,
		roleID,
		"李明",
		"主角",
		true, // 首次登场
	)

	// 验证基本信息
	if appearance.DocumentID != documentID {
		t.Error("DocumentID应该匹配")
	}
	if appearance.CharacterID != characterID {
		t.Error("CharacterID应该匹配")
	}
	if appearance.CharacterName != "李明" {
		t.Errorf("CharacterName应该是'李明'，得到%q", appearance.CharacterName)
	}

	// 验证角色信息
	if appearance.RoleID != roleID {
		t.Error("RoleID应该匹配")
	}
	if appearance.RoleName != "主角" {
		t.Errorf("RoleName应该是'主角'，得到%q", appearance.RoleName)
	}

	// 验证首次登场标记
	if !appearance.FirstAppearance {
		t.Error("FirstAppearance应该是true")
	}

	// 验证时间戳
	if appearance.CreatedAt.IsZero() {
		t.Error("CreatedAt应该被设置")
	}
	if appearance.UpdatedAt.IsZero() {
		t.Error("UpdatedAt应该被设置")
	}
}
