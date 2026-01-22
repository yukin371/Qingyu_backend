package types

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Errors
var (
	ErrInvalidIDFormat = errors.New("invalid ID format: must be 24-character hex")
	ErrEmptyID         = errors.New("ID cannot be empty")
)

// ParseObjectID 将 hex 字符串解析为 ObjectID
// 输入：24字符的 hex 字符串
// 输出：primitive.ObjectID 或 error
func ParseObjectID(s string) (primitive.ObjectID, error) {
	if s == "" {
		return primitive.NilObjectID, ErrEmptyID
	}

	oid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%w: %s", ErrInvalidIDFormat, s)
	}

	return oid, nil
}

// MustParseObjectID 解析 ObjectID，panic on error
// 仅在测试或确定 ID 有效时使用
func MustParseObjectID(s string) primitive.ObjectID {
	oid, err := ParseObjectID(s)
	if err != nil {
		panic(err)
	}
	return oid
}

// ToHex 将 ObjectID 转换为 hex 字符串
// 输入：primitive.ObjectID
// 输出：24字符的 hex 字符串
func ToHex(id primitive.ObjectID) string {
	return id.Hex()
}

// IsValidObjectID 检查字符串是否为有效的 ObjectID hex 格式
func IsValidObjectID(s string) bool {
	_, err := primitive.ObjectIDFromHex(s)
	return err == nil
}

// ParseObjectIDSlice 批量解析 ID 字符串
// 返回：成功解析的 ObjectID 列表和失败索引的映射
func ParseObjectIDSlice(ss []string) ([]primitive.ObjectID, map[int]error) {
	oids := make([]primitive.ObjectID, 0, len(ss))
	errs := make(map[int]error)

	for i, s := range ss {
		oid, err := ParseObjectID(s)
		if err != nil {
			errs[i] = err
			continue
		}
		oids = append(oids, oid)
	}

	return oids, errs
}

// ToHexSlice 批量转换 ObjectID 为 hex 字符串
func ToHexSlice(ids []primitive.ObjectID) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = ToHex(id)
	}
	return result
}

// GenerateNewObjectID 生成新的 ObjectID 并返回 hex 字符串
func GenerateNewObjectID() string {
	return primitive.NewObjectID().Hex()
}

// IsNilObjectID 检查 ObjectID 是否为零值
func IsNilObjectID(id primitive.ObjectID) bool {
	return id.IsZero()
}
