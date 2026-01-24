package reader

// ===========================
// 重新导出 models/dto 中的类型（向后兼容）
// ===========================

import "Qingyu_backend/models/dto"

// ReadingProgressDTO 阅读进度数据传输对象
// 重新导出 models/dto 中的类型
type ReadingProgressDTO = dto.ReadingProgressDTO
