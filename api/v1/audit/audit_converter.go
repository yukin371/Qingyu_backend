package audit

import (
	"time"

	"Qingyu_backend/models/audit"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/shared/types"
)

// ===========================
// AuditRecord Model ↔ DTO 转换函数
// ===========================

// ToViolationDetailDTO Model → DTO 转换（ViolationDetail）
func ToViolationDetailDTO(violation audit.ViolationDetail) ViolationDetailDTO {
	return ViolationDetailDTO{
		Type:        violation.Type,
		Category:    violation.Category,
		Level:       violation.Level,
		Description: violation.Description,
		Position:    violation.Position,
		Context:     violation.Context,
		Keywords:    violation.Keywords,
	}
}

// ToViolationDetailDTOs 批量转换 ViolationDetail
func ToViolationDetailDTOs(violations []audit.ViolationDetail) []ViolationDetailDTO {
	result := make([]ViolationDetailDTO, len(violations))
	for i := range violations {
		result[i] = ToViolationDetailDTO(violations[i])
	}
	return result
}

// ToViolationDetailModel DTO → Model 转换（ViolationDetail）
func ToViolationDetailModel(dto ViolationDetailDTO) audit.ViolationDetail {
	return audit.ViolationDetail{
		Type:        dto.Type,
		Category:    dto.Category,
		Level:       dto.Level,
		Description: dto.Description,
		Position:    dto.Position,
		Context:     dto.Context,
		Keywords:    dto.Keywords,
	}
}

// ToViolationDetailModels 批量转换 ViolationDetail
func ToViolationDetailModels(dtos []ViolationDetailDTO) []audit.ViolationDetail {
	result := make([]audit.ViolationDetail, len(dtos))
	for i := range dtos {
		result[i] = ToViolationDetailModel(dtos[i])
	}
	return result
}

// ToAuditRecordDTO Model → DTO 转换
// 将 AuditRecord Model 转换为 AuditRecordDTO 用于 API 层返回
func ToAuditRecordDTO(record *audit.AuditRecord) *AuditRecordDTO {
	if record == nil {
		return nil
	}

	var converter types.DTOConverter

	// 处理可空时间字段
	var reviewedAt string
	if record.ReviewedAt != nil {
		reviewedAt = converter.TimeToISO8601(*record.ReviewedAt)
	}

	// 转换违规详情
	violationsDTOs := ToViolationDetailDTOs(record.Violations)

	return &AuditRecordDTO{
		ID:        converter.ModelIDToDTO(record.ID),
		CreatedAt: converter.TimeToISO8601(record.CreatedAt),
		UpdatedAt: converter.TimeToISO8601(record.UpdatedAt),
		ReviewedAt: reviewedAt,

		// 审核对象
		TargetType: record.TargetType,
		TargetID:   converter.ModelIDToDTO(record.TargetID),
		AuthorID:   converter.ModelIDToDTO(record.AuthorID),

		// 审核内容
		Content: record.Content,
		Status:  record.Status,
		Result:  record.Result,

		// 风险评估
		RiskLevel: record.RiskLevel,
		RiskScore: float32(record.RiskScore), // float64 → float32

		// 违规详情
		Violations: violationsDTOs,

		// 复核信息
		ReviewerID: converter.ModelIDToDTO(record.ReviewerID),
		ReviewNote: record.ReviewNote,

		// 申诉
		AppealStatus: record.AppealStatus,

		// 元数据
		Metadata: record.Metadata,
	}
}

// ToAuditRecordDTOs 批量转换 Model → DTO
func ToAuditRecordDTOs(records []audit.AuditRecord) []*AuditRecordDTO {
	result := make([]*AuditRecordDTO, len(records))
	for i := range records {
		result[i] = ToAuditRecordDTO(&records[i])
	}
	return result
}

// ToAuditRecordModel 从 DTO 创建 Model（用于更新）
func ToAuditRecordModel(dto *AuditRecordDTO) (*audit.AuditRecord, error) {
	if dto == nil {
		return nil, nil
	}

	var converter types.DTOConverter

	id, err := converter.DTOIDToModel(dto.ID)
	if err != nil {
		return nil, err
	}

	createdAt, err := converter.ISO8601ToTime(dto.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := converter.ISO8601ToTime(dto.UpdatedAt)
	if err != nil {
		return nil, err
	}

	targetID, err := converter.DTOIDToModel(dto.TargetID)
	if err != nil {
		return nil, err
	}

	authorID, err := converter.DTOIDToModel(dto.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewerID, err := converter.DTOIDToModel(dto.ReviewerID)
	if err != nil {
		return nil, err
	}

	// 处理可空时间字段
	var reviewedAt *time.Time
	if dto.ReviewedAt != "" {
		t, err := converter.ISO8601ToTime(dto.ReviewedAt)
		if err != nil {
			return nil, err
		}
		reviewedAt = &t
	}

	// 转换违规详情
	violations := ToViolationDetailModels(dto.Violations)

	return &audit.AuditRecord{
		IdentifiedEntity: shared.IdentifiedEntity{ID: id},
		BaseEntity:       shared.BaseEntity{CreatedAt: createdAt, UpdatedAt: updatedAt},

		TargetType: dto.TargetType,
		TargetID:   targetID,
		AuthorID:   authorID,

		Content: dto.Content,
		Status:  dto.Status,
		Result:  dto.Result,

		RiskLevel: dto.RiskLevel,
		RiskScore: float64(dto.RiskScore), // float32 → float64

		Violations: violations,

		ReviewerID: reviewerID,
		ReviewNote: dto.ReviewNote,

		AppealStatus: dto.AppealStatus,

		Metadata:   dto.Metadata,
		ReviewedAt: reviewedAt,
	}, nil
}

// ToAuditRecordModelWithoutID 从 DTO 创建 Model（用于创建新审核记录）
// 不设置 ID，让数据库自动生成
func ToAuditRecordModelWithoutID(dto *AuditRecordDTO) (*audit.AuditRecord, error) {
	if dto == nil {
		return nil, nil
	}

	var converter types.DTOConverter

	targetID, err := converter.DTOIDToModel(dto.TargetID)
	if err != nil {
		return nil, err
	}

	authorID, err := converter.DTOIDToModel(dto.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewerID, err := converter.DTOIDToModel(dto.ReviewerID)
	if err != nil {
		return nil, err
	}

	// 处理可空时间字段
	var reviewedAt *time.Time
	if dto.ReviewedAt != "" {
		t, err := converter.ISO8601ToTime(dto.ReviewedAt)
		if err != nil {
			return nil, err
		}
		reviewedAt = &t
	}

	// 转换违规详情
	violations := ToViolationDetailModels(dto.Violations)

	return &audit.AuditRecord{
		IdentifiedEntity: shared.IdentifiedEntity{}, // ID 将由数据库生成
		BaseEntity:       shared.BaseEntity{},        // 时间戳将由数据库设置

		TargetType: dto.TargetType,
		TargetID:   targetID,
		AuthorID:   authorID,

		Content: dto.Content,
		Status:  dto.Status,
		Result:  dto.Result,

		RiskLevel: dto.RiskLevel,
		RiskScore: float64(dto.RiskScore), // float32 → float64

		Violations: violations,

		ReviewerID: reviewerID,
		ReviewNote: dto.ReviewNote,

		AppealStatus: dto.AppealStatus,

		Metadata:   dto.Metadata,
		ReviewedAt: reviewedAt,
	}, nil
}
