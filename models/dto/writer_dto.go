package dto

import "time"

// ===========================
// Writer DTO（符合分层架构规范）
// ===========================
//
// 本文件包含 Writer 模块的数据传输对象（DTO）
//
// 命名和标签规范：
// - DTO 结构体使用驼峰命名（PascalCase）
// - JSON 字段标签使用驼峰命名（camelCase）
// - 对应的 MongoDB 模型（位于 models/writer/）使用蛇形命名（snake_case）的 BSON 标签
//
// 用途：
// - 用于 Service 层和 API 层之间的数据传输
// - ID 和时间字段统一使用字符串类型
// - 避免直接暴露 MongoDB 模型到 API 层

// WriterDTO 作家数据传输对象
// 用于：Service 层和 API 层数据传输，ID 和时间字段使用字符串类型
type WriterDTO struct {
	ID        string `json:"id" validate:"required"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// 基本信息
	UserID   string `json:"userId" validate:"required"`          // 关联的用户ID
	Nickname string `json:"nickname" validate:"required,max=50"` // 笔名/昵称
	Bio      string `json:"bio,omitempty"`                       // 个人简介
	Avatar   string `json:"avatar,omitempty"`                    // 头像URL

	// 状态和统计
	Status        string `json:"status" validate:"required,oneof=active suspended deleted"` // 状态
	FollowerCount int    `json:"followerCount" validate:"min=0"`                            // 粉丝数
	BookCount     int    `json:"bookCount" validate:"min=0"`                                // 作品数
	WordCount     int64  `json:"wordCount" validate:"min=0"`                                // 总字数

	// 标记
	IsVerified bool `json:"isVerified"` // 是否认证作者
	IsOfficial bool `json:"isOfficial"` // 是否官方账号
}

// WriterStatsDTO 作家统计数据传输对象
// 用于：展示作家统计信息
type WriterStatsDTO struct {
	WriterID      string  `json:"writerId" validate:"required"`
	BookCount     int     `json:"bookCount" validate:"min=0"`
	WordCount     int64   `json:"wordCount" validate:"min=0"`
	FollowerCount int     `json:"followerCount" validate:"min=0"`
	ViewCount     int64   `json:"viewCount" validate:"min=0"`
	Rating        float64 `json:"rating" validate:"min=0,max=5"`

	// 统计时间
	LastCalculatedAt string `json:"lastCalculatedAt,omitempty"` // 最后计算时间（ISO8601）
}

// CreateWriterRequestDTO 创建作家请求数据传输对象
// 用于：API 层接收创建作家的请求
type CreateWriterRequestDTO struct {
	UserID   string `json:"userId" validate:"required"`
	Nickname string `json:"nickname" validate:"required,max=50"`
	Bio      string `json:"bio,omitempty" validate:"max=500"`
	Avatar   string `json:"avatar,omitempty" validate:"url"`
}

// UpdateWriterRequestDTO 更新作家请求数据传输对象
// 用于：API 层接收更新作家的请求
type UpdateWriterRequestDTO struct {
	Nickname string `json:"nickname,omitempty" validate:"omitempty,max=50"`
	Bio      string `json:"bio,omitempty" validate:"omitempty,max=500"`
	Avatar   string `json:"avatar,omitempty" validate:"omitempty,url"`
	Status   string `json:"status,omitempty" validate:"omitempty,oneof=active suspended deleted"`
}
