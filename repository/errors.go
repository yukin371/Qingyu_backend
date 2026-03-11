package repository

import "errors"

// ID相关错误定义
// 用于统一ID转换和验证的错误语义

var (
	// ErrEmptyID 表示ID为空字符串
	// 用于必需ID的场景，如GetByID、UpdateByID等
	ErrEmptyID = errors.New("ID cannot be empty")

	// ErrInvalidIDFormat 表示ID格式无效（不是有效的ObjectID）
	ErrInvalidIDFormat = errors.New("invalid ID format")
)
