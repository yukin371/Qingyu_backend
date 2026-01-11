package base

import (
	"time"

	shared "Qingyu_backend/models/shared"
)

// BaseEntity 基础实体接口
type BaseEntity interface {
	GetID() string
	SetID(string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	Touch()
}

// 类型别名 - 使用 shared 包的实现
type Timestamps = shared.BaseEntity
type NamedEntity = shared.NamedEntity
type DescriptedEntity = shared.DescriptedEntity
type TitledEntity = shared.TitledEntity
