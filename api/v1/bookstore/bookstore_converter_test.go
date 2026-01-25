package bookstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/bookstore"
	dto "Qingyu_backend/models/dto"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/shared/types"
)

// TestToBookDTO 测试 Model → DTO 转换
func TestToBookDTO(t *testing.T) {
	// 准备测试数据
	authorID := primitive.NewObjectID()
	categoryID1 := primitive.NewObjectID()
	categoryID2 := primitive.NewObjectID()
	now := time.Now()
	publishedAt := now.Add(-24 * time.Hour)
	lastUpdateAt := now.Add(-1 * time.Hour)

	bookModel := &bookstore.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
		BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
		Title:            "测试书籍",
		Author:           "测试作者",
		AuthorID:         authorID,
		Introduction:     "这是一本测试书籍的简介",
		Cover:            "https://example.com/cover.jpg",
		CategoryIDs:      []primitive.ObjectID{categoryID1, categoryID2},
		Categories:       []string{"玄幻", "修仙"},
		Tags:             []string{"热血", "爽文"},
		Status:           bookstore.BookStatusOngoing,
		Rating:           4.5,
		RatingCount:      1000,
		ViewCount:        50000,
		WordCount:        1000000,
		ChapterCount:     500,
		Price:            9900, // 99.00元，以分为单位
		IsFree:           false,
		IsRecommended:    true,
		IsFeatured:       false,
		IsHot:            true,
		PublishedAt:      &publishedAt,
		LastUpdateAt:     &lastUpdateAt,
	}

	// 执行转换
	bookDTO := ToBookDTO(bookModel)

	// 验证基础字段
	require.NotNil(t, bookDTO)
	assert.NotEmpty(t, bookDTO.ID)
	assert.NotEmpty(t, bookDTO.CreatedAt)
	assert.NotEmpty(t, bookDTO.UpdatedAt)

	// 验证基本信息
	assert.Equal(t, "测试书籍", bookDTO.Title)
	assert.Equal(t, "测试作者", bookDTO.Author)
	assert.Equal(t, authorID.Hex(), bookDTO.AuthorID)
	assert.Equal(t, "这是一本测试书籍的简介", bookDTO.Introduction)
	assert.Equal(t, "https://example.com/cover.jpg", bookDTO.Cover)

	// 验证分类和标签
	assert.Equal(t, 2, len(bookDTO.CategoryIDs))
	assert.Contains(t, bookDTO.CategoryIDs, categoryID1.Hex())
	assert.Contains(t, bookDTO.CategoryIDs, categoryID2.Hex())
	assert.Equal(t, []string{"玄幻", "修仙"}, bookDTO.Categories)
	assert.Equal(t, []string{"热血", "爽文"}, bookDTO.Tags)

	// 验证状态和统计
	assert.Equal(t, "ongoing", bookDTO.Status)
	assert.Equal(t, float64(4.5), bookDTO.Rating)
	assert.Equal(t, 1000, bookDTO.RatingCount)
	assert.Equal(t, int64(50000), bookDTO.ViewCount)
	assert.Equal(t, int64(1000000), bookDTO.WordCount)
	assert.Equal(t, 500, bookDTO.ChapterCount)

	// 验证价格 - 应该是 "¥99.00"
	assert.Equal(t, "¥99.00", bookDTO.Price)

	// 验证标记
	assert.False(t, bookDTO.IsFree)
	assert.True(t, bookDTO.IsRecommended)
	assert.False(t, bookDTO.IsFeatured)
	assert.True(t, bookDTO.IsHot)

	// 验证发布信息
	assert.NotEmpty(t, bookDTO.PublishedAt)
	assert.NotEmpty(t, bookDTO.LastUpdateAt)
}

// TestToBookDTOWithNil 测试 nil 输入
func TestToBookDTOWithNil(t *testing.T) {
	bookDTO := ToBookDTO(nil)
	assert.Nil(t, bookDTO)
}

// TestToBookDTOs 测试批量转换
func TestToBookDTOs(t *testing.T) {
	books := []bookstore.Book{
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Title:            "书籍1",
			Author:           "作者1",
			Status:           bookstore.BookStatusCompleted,
			Price:            5000,
		},
		{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Title:            "书籍2",
			Author:           "作者2",
			Status:           bookstore.BookStatusDraft,
			Price:            0,
		},
	}

	dtos := ToBookDTOs(books)

	assert.Equal(t, 2, len(dtos))
	assert.Equal(t, "书籍1", dtos[0].Title)
	assert.Equal(t, "作者1", dtos[0].Author)
	assert.Equal(t, "completed", dtos[0].Status)
	assert.Equal(t, "¥50.00", dtos[0].Price)

	assert.Equal(t, "书籍2", dtos[1].Title)
	assert.Equal(t, "作者2", dtos[1].Author)
	assert.Equal(t, "draft", dtos[1].Status)
	assert.Equal(t, "¥0.00", dtos[1].Price)
}

// TestToBookModel 测试 DTO → Model 转换
func TestToBookModel(t *testing.T) {
	bookID := primitive.NewObjectID()
	authorID := primitive.NewObjectID()
	categoryID1 := primitive.NewObjectID()
	categoryID2 := primitive.NewObjectID()
	now := time.Now()
	publishedAt := now.Add(-24 * time.Hour)
	lastUpdateAt := now.Add(-1 * time.Hour)

	bookDTO := &dto.BookDTO{
		ID:           bookID.Hex(),
		CreatedAt:    now.Format(time.RFC3339),
		UpdatedAt:    now.Format(time.RFC3339),
		Title:        "测试书籍",
		Author:       "测试作者",
		AuthorID:     authorID.Hex(),
		Introduction: "这是一本测试书籍的简介",
		Cover:        "https://example.com/cover.jpg",
		CategoryIDs:  []string{categoryID1.Hex(), categoryID2.Hex()},
		Categories:   []string{"玄幻", "修仙"},
		Tags:         []string{"热血", "爽文"},
		Status:       "ongoing",
		Rating:       4.5,
		RatingCount:  1000,
		ViewCount:    50000,
		WordCount:    1000000,
		ChapterCount: 500,
		Price:        "¥99.00",
		IsFree:       false,
		IsRecommended: true,
		IsFeatured:   false,
		IsHot:        true,
		PublishedAt:  publishedAt.Format(time.RFC3339),
		LastUpdateAt: lastUpdateAt.Format(time.RFC3339),
	}

	// 执行转换
	bookModel, err := ToBookModel(bookDTO)

	// 验证转换成功
	require.NoError(t, err)
	require.NotNil(t, bookModel)

	// 验证 ID
	assert.Equal(t, bookID, bookModel.ID)

	// 验证基本信息
	assert.Equal(t, "测试书籍", bookModel.Title)
	assert.Equal(t, "测试作者", bookModel.Author)
	assert.Equal(t, authorID, bookModel.AuthorID)
	assert.Equal(t, "这是一本测试书籍的简介", bookModel.Introduction)
	assert.Equal(t, "https://example.com/cover.jpg", bookModel.Cover)

	// 验证分类和标签
	assert.Equal(t, 2, len(bookModel.CategoryIDs))
	assert.Contains(t, bookModel.CategoryIDs, categoryID1)
	assert.Contains(t, bookModel.CategoryIDs, categoryID2)
	assert.Equal(t, []string{"玄幻", "修仙"}, bookModel.Categories)
	assert.Equal(t, []string{"热血", "爽文"}, bookModel.Tags)

	// 验证状态
	assert.Equal(t, bookstore.BookStatusOngoing, bookModel.Status)

	// 验证评分
	assert.Equal(t, types.Rating(4.5), bookModel.Rating)
	assert.Equal(t, int64(1000), bookModel.RatingCount)

	// 验证统计
	assert.Equal(t, int64(50000), bookModel.ViewCount)
	assert.Equal(t, int64(1000000), bookModel.WordCount)
	assert.Equal(t, 500, bookModel.ChapterCount)

	// 验证价格 - 应该是 9900 分
	assert.Equal(t, int64(9900), bookModel.Price)

	// 验证标记
	assert.False(t, bookModel.IsFree)
	assert.True(t, bookModel.IsRecommended)
	assert.False(t, bookModel.IsFeatured)
	assert.True(t, bookModel.IsHot)

	// 验证发布时间
	require.NotNil(t, bookModel.PublishedAt)
	require.NotNil(t, bookModel.LastUpdateAt)
	assert.WithinDuration(t, publishedAt, *bookModel.PublishedAt, time.Second)
	assert.WithinDuration(t, lastUpdateAt, *bookModel.LastUpdateAt, time.Second)
}

// TestToBookModelWithNil 测试 nil 输入
func TestToBookModelWithNil(t *testing.T) {
	bookModel, err := ToBookModel(nil)
	assert.NoError(t, err)
	assert.Nil(t, bookModel)
}

// TestBookStatusConversion 测试书籍状态转换
func TestBookStatusConversion(t *testing.T) {
	testCases := []struct {
		name     string
		status   string
		expected bookstore.BookStatus
	}{
		{"草稿", "draft", bookstore.BookStatusDraft},
		{"连载中", "ongoing", bookstore.BookStatusOngoing},
		{"已完结", "completed", bookstore.BookStatusCompleted},
		{"暂停更新", "paused", bookstore.BookStatusPaused},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bookDTO := &dto.BookDTO{
				ID:        primitive.NewObjectID().Hex(),
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
				Title:     "测试",
				Author:    "作者",
				AuthorID:  primitive.NewObjectID().Hex(),
				Status:    tc.status,
				Price:     "¥0.00",
			}

			bookModel, err := ToBookModel(bookDTO)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, bookModel.Status)
		})
	}
}

// TestPriceConversion 测试价格转换
func TestPriceConversion(t *testing.T) {
	testCases := []struct {
		name     string
		price    string
		expected int64 // 以分为单位
	}{
		{"免费", "¥0.00", 0},
		{"1元", "¥1.00", 100},
		{"99元", "¥99.00", 9900},
		{"199.99元", "¥199.99", 19999},
		{"999.99元", "¥999.99", 99999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Model → DTO
			bookModel := &bookstore.Book{
				IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
				BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:            "测试",
				Author:           "作者",
				Status:           bookstore.BookStatusDraft,
				Price:            float64(tc.expected),
			}

			bookDTO := ToBookDTO(bookModel)
			assert.Equal(t, tc.price, bookDTO.Price)

			// DTO → Model
			bookDTO2 := &dto.BookDTO{
				ID:        primitive.NewObjectID().Hex(),
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
				Title:     "测试",
				Author:    "作者",
				AuthorID:  primitive.NewObjectID().Hex(),
				Status:    "draft",
				Price:     tc.price,
			}

			bookModel2, err := ToBookModel(bookDTO2)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, bookModel2.Price)
		})
	}
}

// TestRatingConversion 测试评分转换
func TestRatingConversion(t *testing.T) {
	testCases := []struct {
		name     string
		rating   float32
		expected float32
	}{
		{"零分", 0.0, 0.0},
		{"2.5分", 2.5, 2.5},
		{"满分", 5.0, 5.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Model → DTO
			bookModel := &bookstore.Book{
				IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
				BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Title:            "测试",
				Author:           "作者",
				Status:           bookstore.BookStatusDraft,
				Rating:           types.Rating(tc.expected),
			}

			bookDTO := ToBookDTO(bookModel)
			assert.Equal(t, tc.rating, bookDTO.Rating)

			// DTO → Model
			bookDTO2 := &dto.BookDTO{
				ID:        primitive.NewObjectID().Hex(),
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
				Title:     "测试",
				Author:    "作者",
				AuthorID:  primitive.NewObjectID().Hex(),
				Status:    "draft",
				Rating:    float64(tc.rating),
				Price:     "¥0.00",
			}

			bookModel2, err := ToBookModel(bookDTO2)
			require.NoError(t, err)
			assert.Equal(t, types.Rating(tc.expected), bookModel2.Rating)
		})
	}
}

// TestTimestampConversion 测试时间戳转换
func TestTimestampConversion(t *testing.T) {
	now := time.Now()

	// Model → DTO
	bookModel := &bookstore.Book{
		IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
		BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
		Title:            "测试",
		Author:           "作者",
		Status:           bookstore.BookStatusDraft,
	}

	bookDTO := ToBookDTO(bookModel)

	// 验证时间格式为 ISO8601
	parsedTime, err := time.Parse(time.RFC3339, bookDTO.CreatedAt)
	require.NoError(t, err)
	assert.WithinDuration(t, now, parsedTime, time.Second)

	parsedTime, err = time.Parse(time.RFC3339, bookDTO.UpdatedAt)
	require.NoError(t, err)
	assert.WithinDuration(t, now, parsedTime, time.Second)

	// DTO → Model
	bookDTO2 := &dto.BookDTO{
		ID:        primitive.NewObjectID().Hex(),
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
		Title:     "测试",
		Author:    "作者",
		AuthorID:  primitive.NewObjectID().Hex(),
		Status:    "draft",
		Price:     "¥0.00",
		Rating:    0.0,
	}

	bookModel2, err := ToBookModel(bookDTO2)
	require.NoError(t, err)
	assert.WithinDuration(t, now, bookModel2.CreatedAt, time.Second)
	assert.WithinDuration(t, now, bookModel2.UpdatedAt, time.Second)
}
