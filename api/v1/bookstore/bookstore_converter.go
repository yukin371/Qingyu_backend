package bookstore

import (
	"fmt"
	"time"

	"Qingyu_backend/models/bookstore"
	dto "Qingyu_backend/models/dto"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/shared/types"
)

// ===========================
// Book Model ↔ DTO 转换函数
// ===========================

// ToBookDTO Model → DTO 转换
// 将 Book Model 转换为 BookDTO 用于 API 层返回
func ToBookDTO(book *bookstore.Book) *dto.BookDTO {
	if book == nil {
		return nil
	}

	var converter types.DTOConverter

	// 处理可空时间字段
	var publishedAt, lastUpdateAt string
	if book.PublishedAt != nil {
		publishedAt = converter.TimeToISO8601(*book.PublishedAt)
	}
	if book.LastUpdateAt != nil {
		lastUpdateAt = converter.TimeToISO8601(*book.LastUpdateAt)
	}

	// 转换价格：分 → 金额字符串（Price 现在是 float64，需要转为 int64）
	money := types.NewMoneyFromCents(int64(book.Price))

	return &dto.BookDTO{
		ID:        converter.ModelIDToDTO(book.ID),
		CreatedAt: converter.TimeToISO8601(book.CreatedAt),
		UpdatedAt: converter.TimeToISO8601(book.UpdatedAt),

		// 基本信息
		Title:        book.Title,
		Author:       book.Author,
		AuthorID:     book.AuthorID, // AuthorID 现在是 string 类型，直接赋值
		Introduction: book.Introduction,
		Cover:        book.Cover,

		// 分类和标签
		CategoryIDs: converter.ModelIDsToDTO(book.CategoryIDs),
		Categories:  book.Categories,
		Tags:        book.Tags,

		// 状态和统计
		Status:       string(book.Status),
		Price:        money.String(),
		Rating:       converter.RatingToFloat(book.Rating),
		RatingCount:  int(book.RatingCount),
		ViewCount:    book.ViewCount,
		WordCount:    book.WordCount,
		ChapterCount: book.ChapterCount,

		// 标记
		IsFree:        book.IsFree,
		IsRecommended: book.IsRecommended,
		IsFeatured:    book.IsFeatured,
		IsHot:         book.IsHot,

		// 发布信息
		PublishedAt:  publishedAt,
		LastUpdateAt: lastUpdateAt,
	}
}

// ToBookDTOs 批量转换 Model → DTO（值切片）
func ToBookDTOs(books []bookstore.Book) []*dto.BookDTO {
	result := make([]*BookDTO, len(books))
	for i := range books {
		result[i] = ToBookDTO(&books[i])
	}
	return result
}

// ToBookDTOsFromPtrSlice 批量转换 Model → DTO（指针切片）
func ToBookDTOsFromPtrSlice(books []*bookstore.Book) []*dto.BookDTO {
	result := make([]*BookDTO, len(books))
	for i := range books {
		result[i] = ToBookDTO(books[i])
	}
	return result
}

// ToBookModel 从 DTO 创建 Model（用于更新）
func ToBookModel(dto *dto.BookDTO) (*bookstore.Book, error) {
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

	// AuthorID 现在是 string 类型，直接使用
	authorID := dto.AuthorID

	// 转换分类 ID 列表
	categoryIDs, err := converter.DTOIDsToModel(dto.CategoryIDs)
	if err != nil {
		return nil, err
	}

	// 解析状态 - 直接解析为 bookstore.BookStatus
	status, err := parseBookStoreStatus(dto.Status)
	if err != nil {
		return nil, err
	}

	// 解析评分
	rating, err := converter.DTORatingToFloat(dto.Rating)
	if err != nil {
		return nil, err
	}

	// 解析价格
	price, err := types.ParseMoney(dto.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	// 处理可空时间字段
	var publishedAt, lastUpdateAt *time.Time
	if dto.PublishedAt != "" {
		t, err := converter.ISO8601ToTime(dto.PublishedAt)
		if err != nil {
			return nil, err
		}
		publishedAt = &t
	}
	if dto.LastUpdateAt != "" {
		t, err := converter.ISO8601ToTime(dto.LastUpdateAt)
		if err != nil {
			return nil, err
		}
		lastUpdateAt = &t
	}

	return &bookstore.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: id},
		BaseEntity:       shared.BaseEntity{CreatedAt: createdAt, UpdatedAt: updatedAt},

		Title:        dto.Title,
		Author:       dto.Author,
		AuthorID:     authorID,
		Introduction: dto.Introduction,
		Cover:        dto.Cover,

		CategoryIDs: categoryIDs,
		Categories:  dto.Categories,
		Tags:        dto.Tags,

		Status:       status,
		Rating:       rating,
		RatingCount:  int64(dto.RatingCount),
		ViewCount:    dto.ViewCount,
		WordCount:    dto.WordCount,
		ChapterCount: dto.ChapterCount,
		Price:        float64(price.ToCents()), // Money → Book.Price (int64 → float64)

		IsFree:        dto.IsFree,
		IsRecommended: dto.IsRecommended,
		IsFeatured:    dto.IsFeatured,
		IsHot:         dto.IsHot,

		PublishedAt:  publishedAt,
		LastUpdateAt: lastUpdateAt,
	}, nil
}

// parseBookStoreStatus 从字符串解析书籍状态 (bookstore.BookStatus)
func parseBookStoreStatus(s string) (bookstore.BookStatus, error) {
	status := bookstore.BookStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid book status: %s", s)
	}
	return status, nil
}

// ToBookModelWithoutID 从 DTO 创建 Model（用于创建新书籍）
// 不设置 ID，让数据库自动生成
func ToBookModelWithoutID(dto *dto.BookDTO) (*bookstore.Book, error) {
	if dto == nil {
		return nil, nil
	}

	var converter types.DTOConverter

	// AuthorID 现在是 string 类型，直接使用
	authorID := dto.AuthorID

	// 转换分类 ID 列表
	categoryIDs, err := converter.DTOIDsToModel(dto.CategoryIDs)
	if err != nil {
		return nil, err
	}

	// 解析状态 - 直接解析为 bookstore.BookStatus
	status, err := parseBookStoreStatus(dto.Status)
	if err != nil {
		return nil, err
	}

	// 解析评分
	rating, err := converter.DTORatingToFloat(dto.Rating)
	if err != nil {
		return nil, err
	}

	// 解析价格
	price, err := types.ParseMoney(dto.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	// 处理可空时间字段
	var publishedAt, lastUpdateAt *time.Time
	if dto.PublishedAt != "" {
		t, err := converter.ISO8601ToTime(dto.PublishedAt)
		if err != nil {
			return nil, err
		}
		publishedAt = &t
	}
	if dto.LastUpdateAt != "" {
		t, err := converter.ISO8601ToTime(dto.LastUpdateAt)
		if err != nil {
			return nil, err
		}
		lastUpdateAt = &t
	}

	return &bookstore.Book{
		IdentifiedEntity: shared.IdentifiedEntity{}, // ID 将由数据库生成
		BaseEntity:       shared.BaseEntity{},        // 时间戳将由数据库设置

		Title:        dto.Title,
		Author:       dto.Author,
		AuthorID:     authorID,
		Introduction: dto.Introduction,
		Cover:        dto.Cover,

		CategoryIDs: categoryIDs,
		Categories:  dto.Categories,
		Tags:        dto.Tags,

		Status:       status,
		Rating:       rating,
		RatingCount:  int64(dto.RatingCount),
		ViewCount:    dto.ViewCount,
		WordCount:    dto.WordCount,
		ChapterCount: dto.ChapterCount,
		Price:        float64(price.ToCents()), // Money → Book.Price (int64 → float64)

		IsFree:        dto.IsFree,
		IsRecommended: dto.IsRecommended,
		IsFeatured:    dto.IsFeatured,
		IsHot:         dto.IsHot,

		PublishedAt:  publishedAt,
		LastUpdateAt: lastUpdateAt,
	}, nil
}
