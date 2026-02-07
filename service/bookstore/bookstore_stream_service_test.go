package bookstore

import (
	"context"
	"errors"
	"testing"

	"Qingyu_backend/models/bookstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// =========================
// Mock Stream Repository
// =========================

// MockBookCursor Mock游标接口
type MockBookCursor struct {
	mock.Mock
	books      []*bookstore.Book
	currentIdx int
	closed     bool
	err        error
}

func (m *MockBookCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	if !args.Bool(0) {
		return false
	}
	if m.currentIdx < len(m.books) {
		m.currentIdx++
		return true
	}
	return false
}

func (m *MockBookCursor) Decode(val interface{}) error {
	args := m.Called(val)
	if m.currentIdx <= 0 || m.currentIdx > len(m.books) {
		return errors.New("no data")
	}
	book := m.books[m.currentIdx-1]
	if b, ok := val.(*bookstore.Book); ok {
		*b = *book
		return args.Error(0)
	}
	return args.Error(0)
}

func (m *MockBookCursor) Close(ctx context.Context) error {
	m.closed = true
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBookCursor) Err() error {
	if m.err != nil {
		return m.err
	}
	args := m.Called()
	return args.Error(0)
}

// MockCursorManager Mock游标管理器
type MockCursorManager struct {
	mock.Mock
}

func (m *MockCursorManager) EncodeCursor(cursorType bookstore.CursorType, value interface{}) (string, error) {
	args := m.Called(cursorType, value)
	return args.String(0), args.Error(1)
}

func (m *MockCursorManager) DecodeCursor(encoded string) (*bookstore.StreamCursor, error) {
	args := m.Called(encoded)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.StreamCursor), args.Error(1)
}

func (m *MockCursorManager) GenerateNextCursor(book *bookstore.Book, cursorType bookstore.CursorType, sortField string) (string, error) {
	args := m.Called(book, cursorType, sortField)
	return args.String(0), args.Error(1)
}

// MockBookStreamRepository Mock流式仓储
type MockBookStreamRepository struct {
	mock.Mock
	cursorMgr *MockCursorManager
}

func (m *MockBookStreamRepository) StreamSearch(ctx context.Context, filter *bookstore.BookFilter) (BookCursor, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(BookCursor), args.Error(1)
}

func (m *MockBookStreamRepository) StreamByCursor(ctx context.Context, filter *bookstore.BookFilter) (BookCursor, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(BookCursor), args.Error(1)
}

func (m *MockBookStreamRepository) CountWithFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookStreamRepository) GetCursorManager() CursorManagerInterface {
	return m.cursorMgr
}


func createMockCursor(books []*bookstore.Book, err error) *MockBookCursor {
	cursor := &MockBookCursor{
		books: books,
		err:   err,
	}
	if len(books) > 0 && err == nil {
		cursor.On("Next", mock.Anything).Return(true).Times(len(books))
		cursor.On("Next", mock.Anything).Return(false)
		for _, book := range books {
			cursor.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
				if b, ok := args.Get(0).(*bookstore.Book); ok {
					*b = *book
				}
			}).Return(nil)
		}
	} else {
		cursor.On("Next", mock.Anything).Return(false)
	}
	cursor.On("Close", mock.Anything).Return(nil)
	cursor.On("Err").Return(err)
	return cursor
}

// =========================
// StreamSearch 测试
// =========================

func TestBookstoreStreamService_StreamSearch(t *testing.T) {
	tests := []struct {
		name        string
		filter      *bookstore.BookFilter
		setupMock   func(*MockBookStreamRepository, *MockCursorManager)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *StreamSearchResult)
	}{
		{
			name: "成功流式搜索-有数据且有更多",
			filter: &bookstore.BookFilter{
				Keyword: stringPtr("测试"),
				Limit:   2, // 设置为2，返回2本书表示还有更多
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("测试书籍1", "", bookstore.BookStatusOngoing),
					newTestBook("测试书籍2", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(cursor, nil)
				cm.On("GenerateNextCursor", mock.Anything, mock.Anything, mock.Anything).Return("next-cursor", nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 2)
				assert.Equal(t, "next-cursor", result.NextCursor)
				assert.True(t, result.HasMore)
			},
		},
		{
			name: "成功流式搜索-有数据但无更多",
			filter: &bookstore.BookFilter{
				Keyword: stringPtr("测试"),
				Limit:   20, // 设置为20，但只返回1本书，表示没有更多
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("测试书籍1", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(cursor, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 1)
				assert.Empty(t, result.NextCursor) // 没有更多数据，不生成游标
				assert.False(t, result.HasMore)
			},
		},
		{
			name: "成功流式搜索-无数据",
			filter: &bookstore.BookFilter{
				Limit: 20,
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				cursor := createMockCursor([]*bookstore.Book{}, nil)
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(cursor, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 0)
				assert.Empty(t, result.NextCursor)
				assert.False(t, result.HasMore)
			},
		},
		{
			name:   "filter为nil",
			filter: nil,
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
			},
			wantErr:     true,
			errContains: "filter cannot be nil",
		},
		{
			name: "使用默认Limit",
			filter: &bookstore.BookFilter{
				Limit: 0,
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("测试书籍", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamSearch", mock.Anything, mock.MatchedBy(func(f *bookstore.BookFilter) bool {
					return f.Limit == 20 // 应该使用默认值
				})).Return(cursor, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 1)
				assert.False(t, result.HasMore) // 只有1条数据，Limit是20，所以没有更多
			},
		},
		{
			name: "Limit超过最大值",
			filter: &bookstore.BookFilter{
				Limit: 1000,
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("测试书籍", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamSearch", mock.Anything, mock.MatchedBy(func(f *bookstore.BookFilter) bool {
					return f.Limit == 100 // 应该限制到100
				})).Return(cursor, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 1)
				assert.False(t, result.HasMore)
			},
		},
		{
			name: "仓储错误",
			filter: &bookstore.BookFilter{
				Limit: 20,
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "stream search failed",
		},
		{
			name: "游标错误",
			filter: &bookstore.BookFilter{
				Limit: 20,
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				cursor := createMockCursor([]*bookstore.Book{}, errors.New("cursor error"))
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(cursor, nil)
			},
			wantErr:     true,
			errContains: "cursor error",
		},
		{
			name: "生成游标失败",
			filter: &bookstore.BookFilter{
				Limit: 1, // 设置为1，让hasMore为true
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("测试书籍", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(cursor, nil)
				cm.On("GenerateNextCursor", mock.Anything, mock.Anything, mock.Anything).Return("", errors.New("cursor encoding failed"))
			},
			wantErr:     true,
			errContains: "generate cursor failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cm := new(MockCursorManager)
			repo := &MockBookStreamRepository{cursorMgr: cm}
			service := NewBookstoreStreamService(repo)
			ctx := context.Background()
			tt.setupMock(repo, cm)

			// Act
			result, err := service.StreamSearch(ctx, tt.filter)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
			repo.AssertExpectations(t)
			cm.AssertExpectations(t)
		})
	}
}

// =========================
// StreamSearchWithCursor 测试
// =========================

func TestBookstoreStreamService_StreamSearchWithCursor(t *testing.T) {
	tests := []struct {
		name        string
		filter      *bookstore.BookFilter
		setupMock   func(*MockBookStreamRepository, *MockCursorManager)
		wantErr     bool
		errContains string
		validate    func(*testing.T, *StreamSearchResult)
	}{
		{
			name: "使用游标继续搜索",
			filter: &bookstore.BookFilter{
				Cursor: stringPtr("valid-cursor"),
				Limit:  1, // 设置为1，让hasMore为true
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("书籍1", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamByCursor", mock.Anything, mock.Anything).Return(cursor, nil)
				cm.On("GenerateNextCursor", mock.Anything, mock.Anything, mock.Anything).Return("next-cursor", nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 1)
				assert.Equal(t, "next-cursor", result.NextCursor)
				assert.True(t, result.HasMore)
			},
		},
		{
			name: "空游标-回退到普通搜索",
			filter: &bookstore.BookFilter{
				Cursor: stringPtr(""),
				Limit:  20,
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("书籍1", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(cursor, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 1)
				assert.False(t, result.HasMore)
			},
		},
		{
			name: "nil游标-回退到普通搜索",
			filter: &bookstore.BookFilter{
				Limit: 20,
			},
			setupMock: func(repo *MockBookStreamRepository, cm *MockCursorManager) {
				books := []*bookstore.Book{
					newTestBook("书籍1", "", bookstore.BookStatusOngoing),
				}
				cursor := createMockCursor(books, nil)
				repo.On("StreamSearch", mock.Anything, mock.Anything).Return(cursor, nil)
			},
			wantErr: false,
			validate: func(t *testing.T, result *StreamSearchResult) {
				assert.Len(t, result.Books, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cm := new(MockCursorManager)
			repo := &MockBookStreamRepository{cursorMgr: cm}
			service := NewBookstoreStreamService(repo)
			ctx := context.Background()
			tt.setupMock(repo, cm)

			// Act
			result, err := service.StreamSearchWithCursor(ctx, tt.filter)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
			repo.AssertExpectations(t)
		})
	}
}

// =========================
// CountWithFilter 测试
// =========================

func TestBookstoreStreamService_CountWithFilter(t *testing.T) {
	tests := []struct {
		name      string
		filter    *bookstore.BookFilter
		setupMock func(*MockBookStreamRepository)
		wantCount int64
		wantErr   bool
	}{
		{
			name: "成功统计",
			filter: &bookstore.BookFilter{
				Keyword: stringPtr("测试"),
			},
			setupMock: func(repo *MockBookStreamRepository) {
				repo.On("CountWithFilter", mock.Anything, mock.Anything).Return(int64(100), nil)
			},
			wantCount: 100,
			wantErr:   false,
		},
		{
			name:      "nil filter",
			filter:    nil,
			setupMock: func(repo *MockBookStreamRepository) {
				repo.On("CountWithFilter", mock.Anything, mock.MatchedBy(func(f *bookstore.BookFilter) bool {
					return f != nil
				})).Return(int64(0), nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "仓储错误",
			filter: &bookstore.BookFilter{
				Limit: 20,
			},
			setupMock: func(repo *MockBookStreamRepository) {
				repo.On("CountWithFilter", mock.Anything, mock.Anything).Return(int64(0), errors.New("count failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cm := new(MockCursorManager)
			repo := &MockBookStreamRepository{cursorMgr: cm}
			service := NewBookstoreStreamService(repo)
			ctx := context.Background()
			tt.setupMock(repo)

			// Act
			count, err := service.CountWithFilter(ctx, tt.filter)

			// Assert
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
			}
			repo.AssertExpectations(t)
		})
	}
}

// =========================
// getCursorType 测试
// =========================

func TestBookstoreStreamService_getCursorType(t *testing.T) {
	service := NewBookstoreStreamService(nil)

	tests := []struct {
		name     string
		sortBy   string
		expected bookstore.CursorType
	}{
		{
			name:     "按created_at排序-时间戳游标",
			sortBy:   "created_at",
			expected: bookstore.CursorTypeTimestamp,
		},
		{
			name:     "按updated_at排序-时间戳游标",
			sortBy:   "updated_at",
			expected: bookstore.CursorTypeTimestamp,
		},
		{
			name:     "按published_at排序-时间戳游标",
			sortBy:   "published_at",
			expected: bookstore.CursorTypeTimestamp,
		},
		{
			name:     "按_id排序-ID游标",
			sortBy:   "_id",
			expected: bookstore.CursorTypeID,
		},
		{
			name:     "其他字段-默认时间戳游标",
			sortBy:   "view_count",
			expected: bookstore.CursorTypeTimestamp,
		},
		{
			name:     "空字段-默认时间戳游标",
			sortBy:   "",
			expected: bookstore.CursorTypeTimestamp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getCursorType(tt.sortBy)
			assert.Equal(t, tt.expected, result)
		})
	}
}

