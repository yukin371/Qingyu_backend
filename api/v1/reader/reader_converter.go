package reader

import (
	"fmt"
	"time"

	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/shared/types"
	readerModel "Qingyu_backend/models/reader"
)

// ===========================
// ReadingProgress Model ↔ DTO 转换函数
// ===========================

// ToReadingProgressDTO Model → DTO 转换
// 将 ReadingProgress Model 转换为 ReadingProgressDTO 用于 API 层返回
func ToReadingProgressDTO(progress *readerModel.ReadingProgress) *ReadingProgressDTO {
	if progress == nil {
		return nil
	}

	var converter types.DTOConverter

	return &ReadingProgressDTO{
		ID:        converter.ModelIDToDTO(progress.ID),
		CreatedAt: converter.TimeToISO8601(progress.CreatedAt),
		UpdatedAt: converter.TimeToISO8601(progress.UpdatedAt),

		// 关联信息
		UserID:    converter.ModelIDToDTO(progress.UserID),
		BookID:    converter.ModelIDToDTO(progress.BookID),
		ChapterID: converter.ModelIDToDTO(progress.ChapterID),

		// 阅读进度
		Progress:    converter.ProgressToDTO(progress.Progress),
		ReadingTime: progress.ReadingTime,
		LastReadAt:  converter.TimeToISO8601(progress.LastReadAt),

		// 状态
		Status: progress.Status, // Status 已经是 string 类型
	}
}

// ToReadingProgressDTOs 批量转换 Model → DTO
func ToReadingProgressDTOs(progresses []readerModel.ReadingProgress) []*ReadingProgressDTO {
	result := make([]*ReadingProgressDTO, len(progresses))
	for i := range progresses {
		result[i] = ToReadingProgressDTO(&progresses[i])
	}
	return result
}

// ToReadingProgressDTOsFromPtrSlice 批量转换 Model 指针切片 → DTO
func ToReadingProgressDTOsFromPtrSlice(progresses []*readerModel.ReadingProgress) []*ReadingProgressDTO {
	result := make([]*ReadingProgressDTO, len(progresses))
	for i := range progresses {
		result[i] = ToReadingProgressDTO(progresses[i])
	}
	return result
}

// ToReadingProgressModel 从 DTO 创建 Model（用于更新）
func ToReadingProgressModel(dto *ReadingProgressDTO) (*readerModel.ReadingProgress, error) {
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

	userID, err := converter.DTOIDToModel(dto.UserID)
	if err != nil {
		return nil, err
	}

	bookID, err := converter.DTOIDToModel(dto.BookID)
	if err != nil {
		return nil, err
	}

	chapterID, err := converter.DTOIDToModel(dto.ChapterID)
	if err != nil {
		return nil, err
	}

	progress, err := converter.DTOProgressToModel(dto.Progress)
	if err != nil {
		return nil, err
	}

	lastReadAt, err := converter.ISO8601ToTime(dto.LastReadAt)
	if err != nil {
		return nil, err
	}

	status := dto.Status
	// 验证状态
	switch status {
	case "reading", "want_read", "finished":
		// 有效状态
	default:
		return nil, fmt.Errorf("invalid reading status: %s", dto.Status)
	}

	return &readerModel.ReadingProgress{
		IdentifiedEntity: shared.IdentifiedEntity{ID: id},
		BaseEntity:       shared.BaseEntity{CreatedAt: createdAt, UpdatedAt: updatedAt},

		UserID:    userID,
		BookID:    bookID,
		ChapterID: chapterID,

		Progress:    progress,
		ReadingTime: dto.ReadingTime,
		LastReadAt:  lastReadAt,

		Status: status,
	}, nil
}

// ToReadingProgressModelWithoutID 从 DTO 创建 Model（用于创建新进度）
// 不设置 ID，让数据库自动生成
func ToReadingProgressModelWithoutID(dto *ReadingProgressDTO) (*readerModel.ReadingProgress, error) {
	if dto == nil {
		return nil, nil
	}

	var converter types.DTOConverter

	userID, err := converter.DTOIDToModel(dto.UserID)
	if err != nil {
		return nil, err
	}

	bookID, err := converter.DTOIDToModel(dto.BookID)
	if err != nil {
		return nil, err
	}

	chapterID, err := converter.DTOIDToModel(dto.ChapterID)
	if err != nil {
		return nil, err
	}

	progress, err := converter.DTOProgressToModel(dto.Progress)
	if err != nil {
		return nil, err
	}

	lastReadAt, err := converter.ISO8601ToTime(dto.LastReadAt)
	if err != nil {
		// 如果 LastReadAt 为空，使用零值
		lastReadAt = time.Time{}
	}

	status := dto.Status
	// 验证状态
	switch status {
	case "reading", "want_read", "finished":
		// 有效状态
	default:
		return nil, fmt.Errorf("invalid reading status: %s", dto.Status)
	}

	return &readerModel.ReadingProgress{
		IdentifiedEntity: shared.IdentifiedEntity{}, // ID 将由数据库生成
		BaseEntity:       shared.BaseEntity{},        // 时间戳将由数据库设置

		UserID:    userID,
		BookID:    bookID,
		ChapterID: chapterID,

		Progress:    progress,
		ReadingTime: dto.ReadingTime,
		LastReadAt:  lastReadAt,

		Status: status,
	}, nil
}
