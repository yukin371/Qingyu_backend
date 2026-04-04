package writer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CharacterAppearance 角色在大纲节点中的登场信息
type CharacterAppearance struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DocumentID     primitive.ObjectID `bson:"document_id" json:"documentId"`       // 大纲节点ID
	CharacterID    primitive.ObjectID `bson:"character_id" json:"characterId"`     // 角色ID
	CharacterName  string             `bson:"character_name" json:"characterName"` // 冗余存储，方便查询
	RoleID         primitive.ObjectID `bson:"role_id" json:"roleId"`              // 引用CharacterRole
	RoleName       string             `bson:"role_name" json:"roleName"`          // 冗余存储角色类型名称
	FirstAppearance bool             `bson:"first_appearance" json:"firstAppearance"` // 是否首次登场
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`     // 如"第1卷男二号"
	CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updatedAt"`
}

// NewCharacterAppearance 创建登场信息
func NewCharacterAppearance(documentID, characterID, roleID primitive.ObjectID, characterName, roleName string, firstAppearance bool) CharacterAppearance {
	now := time.Now()
	return CharacterAppearance{
		ID:             primitive.NewObjectID(),
		DocumentID:     documentID,
		CharacterID:    characterID,
		CharacterName:  characterName,
		RoleID:         roleID,
		RoleName:       roleName,
		FirstAppearance: firstAppearance,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
