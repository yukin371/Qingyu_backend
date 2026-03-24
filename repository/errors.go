package repository

import (
	"errors"

	"Qingyu_backend/models/shared/types"
)

// ID相关错误定义
// 用于统一ID转换和验证的错误语义
// 统一使用 models/shared/types 中定义的错误，避免重复

var (
	// ErrEmptyID 表示ID为空字符串
	// 用于必需ID的场景，如GetByID、UpdateByID等
	ErrEmptyID = types.ErrEmptyID

	// ErrInvalidIDFormat 表示ID格式无效（不是有效的ObjectID）
	ErrInvalidIDFormat = types.ErrInvalidIDFormat
)

// IsIDError 判断是否为ID相关错误（兼容旧代码）
func IsIDError(err error) bool {
	return errors.Is(err, types.ErrEmptyID) || errors.Is(err, types.ErrInvalidIDFormat)
}
