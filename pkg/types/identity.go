package types

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EntityID 统一的实体ID类型，使用字符串存储
// 格式为 MongoDB ObjectID 的 24 位十六进制字符串
type EntityID string

// 常用空的 EntityID
var NilEntityID = EntityID("")

// String 实现 Stringer 接口
func (id EntityID) String() string {
	return string(id)
}

// IsEmpty 检查 ID 是否为空
func (id EntityID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid 验证 ID 是否为有效的 MongoDB ObjectID 格式
func (id EntityID) IsValid() bool {
	if id.IsEmpty() {
		return false
	}
	_, err := primitive.ObjectIDFromHex(string(id))
	return err == nil
}

// MustToObjectID 转换为 ObjectID，如果无效则 panic
func (id EntityID) MustToObjectID() primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		panic(fmt.Sprintf("invalid EntityID format: %s", id))
	}
	return oid
}

// ToObjectID 转换为 ObjectID，如果无效返回错误
func (id EntityID) ToObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(string(id))
}

// NewEntityID 从字符串创建 EntityID（自动验证）
func NewEntityID(id string) (EntityID, error) {
	entityID := EntityID(id)
	if !entityID.IsValid() {
		return NilEntityID, errors.New("invalid EntityID format: must be a valid MongoDB ObjectID hex string")
	}
	return entityID, nil
}

// MustEntityID 从字符串创建 EntityID，如果无效则 panic
func MustEntityID(id string) EntityID {
	entityID := EntityID(id)
	if !entityID.IsValid() {
		panic(fmt.Sprintf("invalid EntityID format: %s", id))
	}
	return entityID
}

// NewEntityIDFromObjectID 从 ObjectID 创建 EntityID
func NewEntityIDFromObjectID(oid primitive.ObjectID) EntityID {
	return EntityID(oid.Hex())
}

// GenerateNewEntityID 生成新的 EntityID
func GenerateNewEntityID() EntityID {
	return NewEntityIDFromObjectID(primitive.NewObjectID())
}

// ObjectIDToEntityID 转换工具函数：ObjectID -> EntityID
func ObjectIDToEntityID(oid primitive.ObjectID) EntityID {
	return NewEntityIDFromObjectID(oid)
}

// ObjectIDPtrToEntityIDPtr 转换工具函数：*ObjectID -> *EntityID
func ObjectIDPtrToEntityIDPtr(oidPtr *primitive.ObjectID) *EntityID {
	if oidPtr == nil {
		return nil
	}
	id := NewEntityIDFromObjectID(*oidPtr)
	return &id
}

// EntityIDPtrToObjectIDPtr 转换工具函数：*EntityID -> *ObjectID
func EntityIDPtrToObjectIDPtr(idPtr *EntityID) (*primitive.ObjectID, error) {
	if idPtr == nil {
		return nil, nil
	}
	oid, err := idPtr.ToObjectID()
	return &oid, err
}

// MustEntityIDPtrToObjectIDPtr 转换工具函数：*EntityID -> *ObjectID（panic on error）
func MustEntityIDPtrToObjectIDPtr(idPtr *EntityID) *primitive.ObjectID {
	if idPtr == nil {
		return nil
	}
	oid := idPtr.MustToObjectID()
	return &oid
}

// EntityIDSliceToObjectIdSlice 转换工具函数：[]EntityID -> []ObjectID
func EntityIDSliceToObjectIdSlice(ids []EntityID) ([]primitive.ObjectID, error) {
	result := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := id.ToObjectID()
		if err != nil {
			return nil, fmt.Errorf("failed to convert EntityID %s: %w", id, err)
		}
		result = append(result, oid)
	}
	return result, nil
}

// ObjectIdSliceToEntityIdSlice 转换工具函数：[]ObjectID -> []EntityID
func ObjectIdSliceToEntityIdSlice(oids []primitive.ObjectID) []EntityID {
	result := make([]EntityID, 0, len(oids))
	for _, oid := range oids {
		result = append(result, NewEntityIDFromObjectID(oid))
	}
	return result
}

// StringToEntityIDSlice 安全地将字符串切片转换为 EntityID 切片
func StringToEntityIDSlice(strs []string) ([]EntityID, error) {
	result := make([]EntityID, 0, len(strs))
	for _, s := range strs {
		id, err := NewEntityID(s)
		if err != nil {
			return nil, fmt.Errorf("invalid string ID %s: %w", s, err)
		}
		result = append(result, id)
	}
	return result, nil
}
