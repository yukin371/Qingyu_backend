package social

// 本文件已弃用，所有基础模型已迁移到 models/shared
// 为了向后兼容，这里重新导出 shared 包的类型
//
// 使用方式：
//   import "Qingyu_backend/models/shared"
//   或继续使用 social 包（通过下面的类型别名）

import (
	shared "Qingyu_backend/models/shared"
)

// 类型别名 - 向后兼容
type BaseEntity = shared.BaseEntity
type IdentifiedEntity = shared.IdentifiedEntity
type Timestamps = shared.BaseEntity
type Likable = shared.Likable
type ThreadedConversation = shared.ThreadedConversation
