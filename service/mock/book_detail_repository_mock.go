// Code generated for testing purposes
// Package mock provides mock implementations for testing
package mock

import (
	"context"

	bookstoreModel "Qingyu_backend/models/bookstore"
	"Qingyu_backend/repository/interfaces/bookstore"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockBookDetailRepository is a mock implementation of BookDetailRepository interface
type MockBookDetailRepository struct {
	mock.Mock
}

// Create mocks base method
func (m *MockBookDetailRepository) Create(ctx context.Context, bookDetail *bookstoreModel.BookDetail) error {
	args := m.Called(ctx, bookDetail)
	return args.Error(0)
}

// GetByID mocks base method
func (m *MockBookDetailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookDetail), args.Error(1)
}

// GetByTitle mocks base method
func (m *MockBookDetailRepository) GetByTitle(ctx context.Context, title string) (*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookDetail), args.Error(1)
}

// Update mocks base method
func (m *MockBookDetailRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

// Delete mocks base method
func (m *MockBookDetailRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetByAuthor mocks base method
func (m *MockBookDetailRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, author, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// GetByAuthorID mocks base method
func (m *MockBookDetailRepository) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// GetByCategory mocks base method
func (m *MockBookDetailRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, category, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// GetByStatus mocks base method
func (m *MockBookDetailRepository) GetByStatus(ctx context.Context, status bookstoreModel.BookStatus, limit, offset int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// GetByTags mocks base method
func (m *MockBookDetailRepository) GetByTags(ctx context.Context, tags []string, limit, offset int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, tags, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// Search mocks base method
func (m *MockBookDetailRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// SearchByFilter mocks base method
func (m *MockBookDetailRepository) SearchByFilter(ctx context.Context, filter *bookstore.BookDetailFilter) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// Count mocks base method
func (m *MockBookDetailRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// CountByAuthor mocks base method
func (m *MockBookDetailRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	args := m.Called(ctx, author)
	return args.Get(0).(int64), args.Error(1)
}

// CountByAuthorID mocks base method
func (m *MockBookDetailRepository) CountByAuthorID(ctx context.Context, authorID primitive.ObjectID) (int64, error) {
	args := m.Called(ctx, authorID)
	return args.Get(0).(int64), args.Error(1)
}

// CountByCategory mocks base method
func (m *MockBookDetailRepository) CountByCategory(ctx context.Context, category string) (int64, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(int64), args.Error(1)
}

// CountByStatus mocks base method
func (m *MockBookDetailRepository) CountByStatus(ctx context.Context, status bookstoreModel.BookStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

// CountByTags mocks base method
func (m *MockBookDetailRepository) CountByTags(ctx context.Context, tags []string) (int64, error) {
	args := m.Called(ctx, tags)
	return args.Get(0).(int64), args.Error(1)
}

// CountByFilter mocks base method
func (m *MockBookDetailRepository) CountByFilter(ctx context.Context, filter *bookstore.BookDetailFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// IncrementViewCount mocks base method
func (m *MockBookDetailRepository) IncrementViewCount(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// IncrementLikeCount mocks base method
func (m *MockBookDetailRepository) IncrementLikeCount(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DecrementLikeCount mocks base method
func (m *MockBookDetailRepository) DecrementLikeCount(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// IncrementCommentCount mocks base method
func (m *MockBookDetailRepository) IncrementCommentCount(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DecrementCommentCount mocks base method
func (m *MockBookDetailRepository) DecrementCommentCount(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// IncrementShareCount mocks base method
func (m *MockBookDetailRepository) IncrementShareCount(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// UpdateRating mocks base method
func (m *MockBookDetailRepository) UpdateRating(ctx context.Context, id primitive.ObjectID, rating float64, ratingCount int64) error {
	args := m.Called(ctx, id, rating, ratingCount)
	return args.Error(0)
}

// UpdateLastChapter mocks base method
func (m *MockBookDetailRepository) UpdateLastChapter(ctx context.Context, id primitive.ObjectID, chapterTitle string) error {
	args := m.Called(ctx, id, chapterTitle)
	return args.Error(0)
}

// BatchUpdateStatus mocks base method
func (m *MockBookDetailRepository) BatchUpdateStatus(ctx context.Context, ids []primitive.ObjectID, status bookstoreModel.BookStatus) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

// BatchUpdateCategories mocks base method
func (m *MockBookDetailRepository) BatchUpdateCategories(ctx context.Context, ids []primitive.ObjectID, categoryIDs []string) error {
	args := m.Called(ctx, ids, categoryIDs)
	return args.Error(0)
}

// BatchUpdateTags mocks base method
func (m *MockBookDetailRepository) BatchUpdateTags(ctx context.Context, ids []primitive.ObjectID, tags []string) error {
	args := m.Called(ctx, ids, tags)
	return args.Error(0)
}

// GetByISBN mocks base method
func (m *MockBookDetailRepository) GetByISBN(ctx context.Context, isbn string) (*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookDetail), args.Error(1)
}

// GetByPublisher mocks base method
func (m *MockBookDetailRepository) GetByPublisher(ctx context.Context, publisher string, limit, offset int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, publisher, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// GetByBookID mocks base method
func (m *MockBookDetailRepository) GetByBookID(ctx context.Context, bookID primitive.ObjectID) (*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreModel.BookDetail), args.Error(1)
}

// GetByBookIDs mocks base method
func (m *MockBookDetailRepository) GetByBookIDs(ctx context.Context, bookIDs []primitive.ObjectID) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, bookIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// UpdateAuthor mocks base method
func (m *MockBookDetailRepository) UpdateAuthor(ctx context.Context, bookID primitive.ObjectID, authorID primitive.ObjectID, authorName string) error {
	args := m.Called(ctx, bookID, authorID, authorName)
	return args.Error(0)
}

// GetSimilarBooks mocks base method
func (m *MockBookDetailRepository) GetSimilarBooks(ctx context.Context, bookID primitive.ObjectID, limit int) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, bookID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// CountByPublisher mocks base method
func (m *MockBookDetailRepository) CountByPublisher(ctx context.Context, publisher string) (int64, error) {
	args := m.Called(ctx, publisher)
	return args.Get(0).(int64), args.Error(1)
}

// BatchUpdatePublisher mocks base method
func (m *MockBookDetailRepository) BatchUpdatePublisher(ctx context.Context, bookIDs []primitive.ObjectID, publisher string) error {
	args := m.Called(ctx, bookIDs, publisher)
	return args.Error(0)
}

// List mocks base method
func (m *MockBookDetailRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*bookstoreModel.BookDetail, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstoreModel.BookDetail), args.Error(1)
}

// Exists mocks base method
func (m *MockBookDetailRepository) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// Health mocks base method
func (m *MockBookDetailRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Transaction mocks base method
func (m *MockBookDetailRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
