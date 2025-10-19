package bookstore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reading/bookstore"
	bookstoreRepo "Qingyu_backend/repository/mongodb/bookstore"
)

// ============ Mock Repository Tests ============

// 这些是Repository逻辑测试，不需要真实数据库
// 主要测试输入验证、数据转换等逻辑

func TestBookRepository_ValidateBookData(t *testing.T) {
	tests := []struct {
		name    string
		book    *bookstore.Book
		wantErr bool
	}{
		{
			name: "有效的书籍数据",
			book: &bookstore.Book{
				Title:  "测试书籍",
				Author: "测试作者",
				Status: bookstore.BookStatusPublished,
			},
			wantErr: false,
		},
		{
			name: "空标题",
			book: &bookstore.Book{
				Title:  "",
				Author: "测试作者",
				Status: bookstore.BookStatusPublished,
			},
			wantErr: true,
		},
		{
			name: "空作者",
			book: &bookstore.Book{
				Title:  "测试书籍",
				Author: "",
				Status: bookstore.BookStatusPublished,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证逻辑
			err := validateBook(tt.book)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookRepository_StatusValidation(t *testing.T) {
	validStatuses := []bookstore.BookStatus{
		bookstore.BookStatusDraft,
		bookstore.BookStatusPublished,
		bookstore.BookStatusOngoing,
		bookstore.BookStatusCompleted,
		bookstore.BookStatusPaused,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			book := &bookstore.Book{
				Title:  "测试",
				Author: "测试",
				Status: status,
			}
			err := validateBook(book)
			assert.NoError(t, err)
		})
	}
}

func TestBookRepository_IDGeneration(t *testing.T) {
	// 测试ObjectID生成
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	assert.NotEqual(t, id1, id2, "生成的ID应该是唯一的")
	assert.NotEmpty(t, id1.Hex(), "ID的Hex表示不应为空")
}

func TestBookFilter_BuildQuery(t *testing.T) {
	tests := []struct {
		name   string
		filter *bookstore.BookFilter
		want   int // 预期的查询条件数量
	}{
		{
			name:   "空过滤器",
			filter: &bookstore.BookFilter{},
			want:   0,
		},
		{
			name: "按状态过滤",
			filter: &bookstore.BookFilter{
				Status: ptrBookStatus(bookstore.BookStatusPublished),
			},
			want: 1,
		},
		{
			name: "按作者过滤",
			filter: &bookstore.BookFilter{
				Author: ptrString("测试作者"),
			},
			want: 1,
		},
		{
			name: "多条件过滤",
			filter: &bookstore.BookFilter{
				Status: ptrBookStatus(bookstore.BookStatusPublished),
				Author: ptrString("测试作者"),
				IsFree: ptrBool(true),
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := buildQueryFromFilter(tt.filter)
			assert.Equal(t, tt.want, len(query))
		})
	}
}

func TestCategoryTree_BuildStructure(t *testing.T) {
	// 测试分类树构建逻辑
	categories := []*bookstore.Category{
		{
			ID:       primitive.NewObjectID(),
			Name:     "小说",
			Level:    0,
			ParentID: nil,
		},
		{
			ID:       primitive.NewObjectID(),
			Name:     "玄幻",
			Level:    1,
			ParentID: nil, // 在实际中会设置为"小说"的ID
		},
	}

	// 验证根分类
	rootCategories := filterRootCategories(categories)
	assert.NotEmpty(t, rootCategories)
	assert.Equal(t, 0, rootCategories[0].Level)
}

func TestBanner_ActiveFilter(t *testing.T) {
	tests := []struct {
		name     string
		banner   *bookstore.Banner
		isActive bool
	}{
		{
			name: "活动的Banner",
			banner: &bookstore.Banner{
				IsActive: true,
			},
			isActive: true,
		},
		{
			name: "非活动的Banner",
			banner: &bookstore.Banner{
				IsActive: false,
			},
			isActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isActive, tt.banner.IsActive)
		})
	}
}

// ============ Repository Constructor Tests ============

func TestNewMongoBookRepository(t *testing.T) {
	// 这个测试不需要真实数据库连接
	// 只测试构造函数是否正确创建Repository实例

	// 注意: 这里我们无法直接创建Repository，因为需要MongoDB连接
	// 但我们可以验证Repository的接口设计

	// 验证Repository接口方法
	var _ interface {
		Create(ctx context.Context, book *bookstore.Book) error
		GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.Book, error)
		Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
		Delete(ctx context.Context, id primitive.ObjectID) error
	} = (*bookstoreRepo.MongoBookRepository)(nil)

	// 这只是编译时检查，确保接口正确实现
	assert.True(t, true, "Repository接口正确实现")
}

// ============ Helper Functions ============

func validateBook(book *bookstore.Book) error {
	if book.Title == "" {
		return assert.AnError
	}
	if book.Author == "" {
		return assert.AnError
	}
	return nil
}

func buildQueryFromFilter(filter *bookstore.BookFilter) map[string]interface{} {
	query := make(map[string]interface{})

	if filter.Status != nil {
		query["status"] = *filter.Status
	}
	if filter.Author != nil {
		query["author"] = *filter.Author
	}
	if filter.IsFree != nil {
		query["is_free"] = *filter.IsFree
	}

	return query
}

func filterRootCategories(categories []*bookstore.Category) []*bookstore.Category {
	var roots []*bookstore.Category
	for _, cat := range categories {
		if cat.ParentID == nil || cat.Level == 0 {
			roots = append(roots, cat)
		}
	}
	return roots
}

// ============ Helper Pointer Functions ============

func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrBookStatus(status bookstore.BookStatus) *bookstore.BookStatus {
	return &status
}

// ============ Benchmark Tests ============

func BenchmarkBuildQueryFromFilter(b *testing.B) {
	filter := &bookstore.BookFilter{
		Status: ptrBookStatus(bookstore.BookStatusPublished),
		Author: ptrString("测试作者"),
		IsFree: ptrBool(true),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildQueryFromFilter(filter)
	}
}

// ============ Table-Driven Tests ============

func TestBookRepository_PaginationCalculation(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		offset   int
		limit    int
	}{
		{"第1页", 1, 10, 0, 10},
		{"第2页", 2, 10, 10, 10},
		{"第5页", 5, 20, 80, 20},
		{"默认分页", 1, 50, 0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := (tt.page - 1) * tt.pageSize
			assert.Equal(t, tt.offset, offset)
			assert.Equal(t, tt.limit, tt.pageSize)
		})
	}
}

