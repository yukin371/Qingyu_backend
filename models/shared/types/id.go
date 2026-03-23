/*
* 该文件已被弃用
    请使用 repository.ParseID 替代
    相关迁移指南请参考: docs/guides/id-error-handling-guide.md
*/

package types

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"Qingyu_backend/repository"
)

// Errors - 与 repository/errors.go 保持一致
// Deprecated: 使用 repository.ErrEmptyID 和 repository.ErrInvalidIDFormat
var (
	ErrInvalidIDFormat = errors.New("invalid ID format")
	ErrEmptyID         = errors.New("ID cannot be empty")
)

// ParseObjectID 将 hex 字符串解析为 ObjectID
// Deprecated: 使用 repository.ParseID 替代
func ParseObjectID(s string) (primitive.ObjectID, error) {
    if s == "" {
        return primitive.NilObjectID, ErrEmptyID
    }
    // 统一使用 repository.ParseID
    return repository.ParseID(s)
}

// MustParseObjectID 解析 ObjectID， panic on error
// Deprecated: 使用 repository.ParseID 替代
func MustParseObjectID(s string) primitive.ObjectID {
    oid, err := ParseObjectID(s)
    if err != nil {
        panic(err)
    }
    return oid
}

// ToHex 将 ObjectID 转换为 hex 字符串
// Deprecated: 使用 oid.Hex() 替代
func ToHex(id primitive.ObjectID) string {
    return id.Hex()
}
// IsValidObjectID 检查字符串是否为有效的 ObjectID hex 格式
// Deprecated: 使用 repository.IsIDError(repository.ParseID(id))
func IsValidObjectID(s string) bool {
    _, err := primitive.ObjectIDFromHex(s)
    return err == nil
}
// ParseObjectIDSlice 批量解析 ID 字符串
// Deprecated: 使用 repository.ParseIDs 替代
func ParseObjectIDSlice(ss []string) ([]primitive.ObjectID, map[int]error) {
    oids := make([]primitive.ObjectID, 0, len(ss))
    errMap := make(map[int]error)
    for i, s := range ss {
        oid, err := ParseObjectID(s)
        if err != nil {
            errMap[i] = err
            continue
        }
        oids = append(oids, oid)
    }
    return oids, errMap
}
// ToHexSlice 批量转换 ObjectID 为 hex 字符串
// Deprecated: 使用 repository 包的 ToHexSlice 替代
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
