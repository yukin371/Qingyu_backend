package audit

// ===========================
// 重新导出 models/dto 中的类型（向后兼容）
// ===========================

import "Qingyu_backend/models/dto"

// AuditRecordDTO 审计记录数据传输对象
// 重新导出 models/dto 中的类型
type AuditRecordDTO = dto.AuditRecordDTO

// ViolationDetailDTO 违规详情数据传输对象
// 重新导出 models/dto 中的类型
type ViolationDetailDTO = dto.ViolationDetailDTO
