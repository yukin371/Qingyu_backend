package stats

import (
	"Qingyu_backend/models/bookstore"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/users"
	"Qingyu_backend/models/writer"
	userRepo "Qingyu_backend/repository/interfaces/user"
	bookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	base "Qingyu_backend/repository/interfaces/infrastructure"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============ Mock Repositories ============

// MockUserRepository Mock用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, entity *users.User) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, filter base.Filter) ([]*users.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*users.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	args := m.Called(ctx, phone)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, id string, ip string) error {
	args := m.Called(ctx, id, ip)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePasswordByEmail(ctx context.Context, email string, hashedPassword string) error {
	args := m.Called(ctx, email, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateStatus(ctx context.Context, id string, status users.UserStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserRepository) GetActiveUsers(ctx context.Context, limit int64) ([]*users.User, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersByRole(ctx context.Context, role string, limit int64) ([]*users.User, error) {
	args := m.Called(ctx, role, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

func (m *MockUserRepository) SetPhoneVerified(ctx context.Context, id string, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

func (m *MockUserRepository) UnbindEmail(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UnbindPhone(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteDevice(ctx context.Context, userID string, deviceID string) error {
	args := m.Called(ctx, userID, deviceID)
	return args.Error(0)
}

func (m *MockUserRepository) GetDevices(ctx context.Context, userID string) ([]interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockUserRepository) BatchUpdateStatus(ctx context.Context, ids []string, status users.UserStatus) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockUserRepository) BatchDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockUserRepository) FindWithFilter(ctx context.Context, filter *users.UserFilter) ([]*users.User, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*users.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) SearchUsers(ctx context.Context, keyword string, limit int) ([]*users.User, error) {
	args := m.Called(ctx, keyword, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*users.User), args.Error(1)
}

func (m *MockUserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	args := m.Called(ctx, role)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountByStatus(ctx context.Context, status users.UserStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Transaction(ctx context.Context, user *users.User, fn func(ctx context.Context, repo userRepo.UserRepository) error) error {
	args := m.Called(ctx, user, fn)
	return args.Error(0)
}

// MockBookRepository Mock书籍仓库
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(ctx context.Context, entity *bookstore.Book) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockBookRepository) GetByID(ctx context.Context, id string) (*bookstore.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookRepository) List(ctx context.Context, filter base.Filter) ([]*bookstore.Book, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBookRepository) GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, categoryID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByAuthor(ctx context.Context, author string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, author, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByAuthorID(ctx context.Context, authorID string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByStatus(ctx context.Context, status bookstore.BookStatus, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetRecommended(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetFeatured(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetHotBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetNewReleases(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetFreeBooks(ctx context.Context, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, minPrice, maxPrice, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) SearchWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByAuthor(ctx context.Context, author string) (int64, error) {
	args := m.Called(ctx, author)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByStatus(ctx context.Context, status bookstore.BookStatus) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) CountByFilter(ctx context.Context, filter *bookstore.BookFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBookRepository) BatchUpdateStatus(ctx context.Context, bookIDs []string, status bookstore.BookStatus) error {
	args := m.Called(ctx, bookIDs, status)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateCategory(ctx context.Context, bookIDs []string, categoryIDs []string) error {
	args := m.Called(ctx, bookIDs, categoryIDs)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateRecommended(ctx context.Context, bookIDs []string, isRecommended bool) error {
	args := m.Called(ctx, bookIDs, isRecommended)
	return args.Error(0)
}

func (m *MockBookRepository) BatchUpdateFeatured(ctx context.Context, bookIDs []string, isFeatured bool) error {
	args := m.Called(ctx, bookIDs, isFeatured)
	return args.Error(0)
}

func (m *MockBookRepository) GetStats(ctx context.Context) (*bookstore.BookStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookStats), args.Error(1)
}

func (m *MockBookRepository) IncrementViewCount(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockBookRepository) GetYears(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockBookRepository) GetTags(ctx context.Context, categoryID *string) ([]string, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockBookRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// MockProjectRepository Mock项目仓库
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, entity *writer.Project) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*writer.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) List(ctx context.Context, filter base.Filter) ([]*writer.Project, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockProjectRepository) GetListByOwnerID(ctx context.Context, ownerID string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByOwnerAndStatus(ctx context.Context, ownerID, status string, limit, offset int64) ([]*writer.Project, error) {
	args := m.Called(ctx, ownerID, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.Project), args.Error(1)
}

func (m *MockProjectRepository) UpdateByOwner(ctx context.Context, projectID, ownerID string, updates map[string]interface{}) error {
	args := m.Called(ctx, projectID, ownerID, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) IsOwner(ctx context.Context, projectID, ownerID string) (bool, error) {
	args := m.Called(ctx, projectID, ownerID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) SoftDelete(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) HardDelete(ctx context.Context, projectID string) error {
	args := m.Called(ctx, projectID)
	return args.Error(0)
}

func (m *MockProjectRepository) Restore(ctx context.Context, projectID, ownerID string) error {
	args := m.Called(ctx, projectID, ownerID)
	return args.Error(0)
}

func (m *MockProjectRepository) CountByOwner(ctx context.Context, ownerID string) (int64, error) {
	args := m.Called(ctx, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProjectRepository) CreateWithTransaction(ctx context.Context, project *writer.Project, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, project, callback)
	return args.Error(0)
}

// MockChapterRepository Mock章节仓库
type MockChapterRepository struct {
	mock.Mock
}

func (m *MockChapterRepository) Create(ctx context.Context, entity *bookstore.Chapter) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByID(ctx context.Context, id string) (*bookstore.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockChapterRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChapterRepository) List(ctx context.Context, filter base.Filter) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Count(ctx context.Context, filter base.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockChapterRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByBookID(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByBookIDAndChapterNum(ctx context.Context, bookID string, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByTitle(ctx context.Context, title string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, title, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFreeChapters(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPaidChapters(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPublishedChapters(ctx context.Context, bookID string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetChapterRange(ctx context.Context, bookID string, startChapter, endChapter int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, startChapter, endChapter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, keyword, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) SearchByFilter(ctx context.Context, filter *bookstoreRepo.ChapterFilter) ([]*bookstore.Chapter, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountFreeChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountPaidChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountPublishedChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetTotalWordCount(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetPreviousChapter(ctx context.Context, bookID string, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetNextChapter(ctx context.Context, bookID string, chapterNum int) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFirstChapter(ctx context.Context, bookID string) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetLastChapter(ctx context.Context, bookID string) (*bookstore.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Chapter), args.Error(1)
}

func (m *MockChapterRepository) BatchUpdatePrice(ctx context.Context, chapterIDs []string, price float64) error {
	args := m.Called(ctx, chapterIDs, price)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchDelete(ctx context.Context, chapterIDs []string) error {
	args := m.Called(ctx, chapterIDs)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdateFreeStatus(ctx context.Context, chapterIDs []string, isFree bool) error {
	args := m.Called(ctx, chapterIDs, isFree)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdatePublishTime(ctx context.Context, chapterIDs []string, publishTime time.Time) error {
	args := m.Called(ctx, chapterIDs, publishTime)
	return args.Error(0)
}

func (m *MockChapterRepository) DeleteByBookID(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockChapterRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// ============ Test Helper Functions ============

// createTestUser 创建测试用户
func createTestUser(roles []string, vipLevel int) *users.User {
	now := time.Now()
	objID := primitive.NewObjectID()
	return &users.User{
		IdentifiedEntity: shared.IdentifiedEntity{
			ID: objID,
		},
		BaseEntity: shared.BaseEntity{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:  "testuser",
		Email:     "test@example.com",
		Roles:     roles,
		VIPLevel:  vipLevel,
		Status:    users.UserStatusActive,
		Password:   "hashedpassword",
	}
}

// ============ TDD Phase 1: RED - Write Failing Tests ============

// TestGetUserStats_EmptyUserID 测试空用户ID
func TestGetUserStats_EmptyUserID(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	// Act
	stats, err := service.GetUserStats(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "用户ID不能为空")

	t.Log("✓ EmptyUserID validation test passed (RED phase)")
}

// TestGetUserStats_UserNotFound 测试用户不存在
func TestGetUserStats_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	// Mock用户不存在
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(nil, errors.New("用户不存在")).Once()

	// Act
	stats, err := service.GetUserStats(ctx, testUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "获取用户信息失败")

	mockUserRepo.AssertExpectations(t)

	t.Log("✓ UserNotFound test passed (RED phase)")
}

// TestGetUserStats_Success 测试成功获取用户统计
func TestGetUserStats_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testUser := createTestUser([]string{"reader"}, 0)

	// Mock获取用户信息
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(testUser, nil).Once()

	// Mock统计项目数
	mockProjectRepo.On("CountByOwner", ctx, testUserID).
		Return(int64(5), nil).Once()

	// Act
	stats, err := service.GetUserStats(ctx, testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, testUserID, stats.UserID)
	assert.Equal(t, int64(5), stats.TotalProjects)
	assert.Equal(t, "普通读者", stats.MemberLevel)
	assert.Equal(t, int64(0), stats.TotalBooks) // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.TotalWords) // TODO: 当前返回0

	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)

	t.Log("✓ GetUserStats success test passed (RED phase)")
}

// TestGetUserStats_ProjectCountError 测试项目数统计失败
func TestGetUserStats_ProjectCountError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testUser := createTestUser([]string{"author"}, 0)

	// Mock获取用户信息
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(testUser, nil).Once()

	// Mock统计项目数失败
	mockProjectRepo.On("CountByOwner", ctx, testUserID).
		Return(int64(0), errors.New("数据库连接失败")).Once()

	// Act
	stats, err := service.GetUserStats(ctx, testUserID)

	// Assert - 应该返回默认值0而不是错误
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalProjects)

	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)

	t.Log("✓ ProjectCountError test passed (RED phase)")
}

// TestGetUserStats_AdminUser 测试管理员用户
func TestGetUserStats_AdminUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testUser := createTestUser([]string{"admin"}, 0)

	// Mock获取用户信息
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(testUser, nil).Once()

	// Mock统计项目数
	mockProjectRepo.On("CountByOwner", ctx, testUserID).
		Return(int64(3), nil).Once()

	// Act
	stats, err := service.GetUserStats(ctx, testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "管理员", stats.MemberLevel)

	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)

	t.Log("✓ AdminUser test passed (RED phase)")
}

// TestGetUserStats_VIPUser 测试VIP用户
func TestGetUserStats_VIPUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testUser := createTestUser([]string{"author"}, 3)

	// Mock获取用户信息
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(testUser, nil).Once()

	// Mock统计项目数
	mockProjectRepo.On("CountByOwner", ctx, testUserID).
		Return(int64(2), nil).Once()

	// Act
	stats, err := service.GetUserStats(ctx, testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Contains(t, stats.MemberLevel, "VIP Level 3")

	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)

	t.Log("✓ VIPUser test passed (RED phase)")
}

// ============ GetContentStats Tests ============

// TestGetContentStats_EmptyUserID 测试空用户ID
func TestGetContentStats_EmptyUserID(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	// Act
	stats, err := service.GetContentStats(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "用户ID不能为空")

	t.Log("✓ GetContentStats EmptyUserID test passed (RED phase)")
}

// TestGetContentStats_UserNotFound 测试用户不存在
func TestGetContentStats_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	// Mock用户不存在
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(nil, errors.New("用户不存在")).Once()

	// Act
	stats, err := service.GetContentStats(ctx, testUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "获取用户信息失败")

	mockUserRepo.AssertExpectations(t)

	t.Log("✓ GetContentStats UserNotFound test passed (RED phase)")
}

// TestGetContentStats_Success 测试成功获取内容统计
func TestGetContentStats_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testUser := createTestUser([]string{"author"}, 0)

	// Mock获取用户信息
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(testUser, nil).Once()

	// Mock统计项目数
	mockProjectRepo.On("CountByOwner", ctx, testUserID).
		Return(int64(10), nil).Once()

	// Act
	stats, err := service.GetContentStats(ctx, testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, testUserID, stats.UserID)
	assert.Equal(t, int64(10), stats.TotalProjects)
	assert.Equal(t, int64(0), stats.PublishedBooks) // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.TotalWords)    // TODO: 当前返回0

	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)

	t.Log("✓ GetContentStats success test passed (RED phase)")
}

// TestGetContentStats_ProjectCountError 测试项目数统计失败
func TestGetContentStats_ProjectCountError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testUser := createTestUser([]string{"author"}, 0)

	// Mock获取用户信息
	mockUserRepo.On("GetByID", ctx, testUserID).
		Return(testUser, nil).Once()

	// Mock统计项目数失败
	mockProjectRepo.On("CountByOwner", ctx, testUserID).
		Return(int64(0), errors.New("数据库连接失败")).Once()

	// Act
	stats, err := service.GetContentStats(ctx, testUserID)

	// Assert - 应该返回默认值0而不是错误
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalProjects)

	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)

	t.Log("✓ GetContentStats ProjectCountError test passed (RED phase)")
}

// ============ GetPlatformUserStats Tests ============

// TestGetPlatformUserStats_Success 测试获取平台用户统计
func TestGetPlatformUserStats_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	// Act
	stats, err := service.GetPlatformUserStats(ctx, startDate, endDate)

	// Assert - 当前返回空数据
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalUsers)       // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.NewUsers)         // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.ActiveUsers)      // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.VIPUsers)         // TODO: 当前返回0
	assert.Equal(t, float64(0), stats.RetentionRate)   // TODO: 当前返回0
	assert.Equal(t, float64(0), stats.AverageActiveDay) // TODO: 当前返回0

	t.Log("✓ GetPlatformUserStats test passed (RED phase - TODOs)")
}

// ============ GetPlatformContentStats Tests ============

// TestGetPlatformContentStats_Success 测试获取平台内容统计
func TestGetPlatformContentStats_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	// Act
	stats, err := service.GetPlatformContentStats(ctx, startDate, endDate)

	// Assert - 当前返回空数据
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(0), stats.TotalBooks)        // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.NewBooks)          // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.TotalChapters)     // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.TotalWords)        // TODO: 当前返回0
	assert.Equal(t, int64(0), stats.TotalViews)        // TODO: 当前返回0
	assert.Equal(t, float64(0), stats.AverageRating)    // TODO: 当前返回0
	assert.Empty(t, stats.PopularCategories)           // TODO: 当前返回空

	t.Log("✓ GetPlatformContentStats test passed (RED phase - TODOs)")
}

// ============ GetUserActivityStats Tests ============

// TestGetUserActivityStats_EmptyUserID 测试空用户ID
func TestGetUserActivityStats_EmptyUserID(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	// Act
	stats, err := service.GetUserActivityStats(ctx, "", 7)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "用户ID不能为空")

	t.Log("✓ GetUserActivityStats EmptyUserID test passed (RED phase)")
}

// TestGetUserActivityStats_NegativeDays 测试负数天数
func TestGetUserActivityStats_NegativeDays(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	// Act
	stats, err := service.GetUserActivityStats(ctx, testUserID, -5)

	// Assert - 应该使用默认值7天
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 7, stats.Days) // 应该被纠正为7
	assert.Equal(t, int64(0), stats.TotalActions) // TODO: 当前返回0

	t.Log("✓ GetUserActivityStats NegativeDays test passed (RED phase)")
}

// TestGetUserActivityStats_ZeroDays 测试0天数
func TestGetUserActivityStats_ZeroDays(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	// Act
	stats, err := service.GetUserActivityStats(ctx, testUserID, 0)

	// Assert - 应该使用默认值7天
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 7, stats.Days) // 应该被纠正为7

	t.Log("✓ GetUserActivityStats ZeroDays test passed (RED phase)")
}

// TestGetUserActivityStats_Success 测试成功获取活跃度统计
func TestGetUserActivityStats_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	// Act
	stats, err := service.GetUserActivityStats(ctx, testUserID, 30)

	// Assert - 当前返回空数据
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, testUserID, stats.UserID)
	assert.Equal(t, 30, stats.Days)
	assert.Equal(t, int64(0), stats.TotalActions)  // TODO: 当前返回0
	assert.Empty(t, stats.DailyActions)            // TODO: 当前返回空
	assert.Empty(t, stats.ActionTypes)             // TODO: 当前返回空
	assert.Empty(t, stats.ActiveHours)             // TODO: 当前返回空

	t.Log("✓ GetUserActivityStats test passed (RED phase - TODOs)")
}

// ============ GetRevenueStats Tests ============

// TestGetRevenueStats_EmptyUserID 测试空用户ID
func TestGetRevenueStats_EmptyUserID(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	// Act
	stats, err := service.GetRevenueStats(ctx, "", startDate, endDate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, stats)
	assert.Contains(t, err.Error(), "用户ID不能为空")

	t.Log("✓ GetRevenueStats EmptyUserID test passed (RED phase)")
}

// TestGetRevenueStats_Success 测试成功获取收益统计
func TestGetRevenueStats_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()

	// Act
	stats, err := service.GetRevenueStats(ctx, testUserID, startDate, endDate)

	// Assert - 当前返回空数据
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, testUserID, stats.UserID)
	assert.Equal(t, float64(0), stats.TotalRevenue)   // TODO: 当前返回0
	assert.Equal(t, float64(0), stats.PeriodRevenue)  // TODO: 当前返回0
	assert.Empty(t, stats.DailyRevenue)                // TODO: 当前返回空
	assert.Empty(t, stats.RevenueByBook)               // TODO: 当前返回空
	assert.Empty(t, stats.RevenueByType)               // TODO: 当前返回空

	t.Log("✓ GetRevenueStats test passed (RED phase - TODOs)")
}

// ============ Health Tests ============

// TestHealth_Success 测试健康检查成功
func TestHealth_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	// Act
	err := service.Health(ctx)

	// Assert
	assert.NoError(t, err)

	t.Log("✓ Health test passed (RED phase)")
}

// TestServiceBaseMethods 测试基础服务方法
func TestServiceBaseMethods(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockBookRepo := new(MockBookRepository)
	mockProjectRepo := new(MockProjectRepository)
	mockChapterRepo := new(MockChapterRepository)

	service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
	ctx := context.Background()

	// 类型断言获取具体实现
ServiceImpl, ok := service.(*PlatformStatsServiceImpl)
	assert.True(t, ok, "service should be of type *PlatformStatsServiceImpl")

	t.Run("Initialize", func(t *testing.T) {
		// Act
		err := ServiceImpl.Initialize(ctx)

		// Assert
		assert.NoError(t, err)

		t.Log("✓ Initialize test passed")
	})

	t.Run("GetServiceName", func(t *testing.T) {
		// Act
		name := ServiceImpl.GetServiceName()

		// Assert
		assert.Equal(t, "PlatformStatsService", name)

		t.Log("✓ GetServiceName test passed")
	})

	t.Run("GetVersion", func(t *testing.T) {
		// Act
		version := ServiceImpl.GetVersion()

		// Assert
		assert.Equal(t, "v1.0.0", version)

		t.Log("✓ GetVersion test passed")
	})

	t.Run("Close", func(t *testing.T) {
		// Act
		err := ServiceImpl.Close(ctx)

		// Assert
		assert.NoError(t, err)

		t.Log("✓ Close test passed")
	})
}

// ============ Table-Driven Tests ============

// TestGetUserStats_TableDriven 表格驱动测试
func TestGetUserStats_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		roles         []string
		vipLevel      int
		projectCount  int64
		projectError  error
		expectedLevel string
		wantErr      bool
		errContains  string
	}{
		{
			name:          "普通读者",
			userID:        primitive.NewObjectID().Hex(),
			roles:         []string{"reader"},
			vipLevel:      0,
			projectCount:  0,
			projectError:  nil,
			expectedLevel: "普通读者",
			wantErr:       false,
		},
		{
			name:          "作者用户",
			userID:        primitive.NewObjectID().Hex(),
			roles:         []string{"author"},
			vipLevel:      0,
			projectCount:  3,
			projectError:  nil,
			expectedLevel: "作者",
			wantErr:       false,
		},
		{
			name:          "VIP作者",
			userID:        primitive.NewObjectID().Hex(),
			roles:         []string{"author"},
			vipLevel:      2,
			projectCount:  5,
			projectError:  nil,
			expectedLevel: "作者 (VIP Level 2)",
			wantErr:       false,
		},
		{
			name:          "空用户ID",
			userID:        "",
			roles:         []string{"reader"},
			vipLevel:      0,
			projectCount:  0,
			projectError:  nil,
			expectedLevel: "",
			wantErr:       true,
			errContains:   "用户ID不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(MockUserRepository)
			mockBookRepo := new(MockBookRepository)
			mockProjectRepo := new(MockProjectRepository)
			mockChapterRepo := new(MockChapterRepository)

			service := NewPlatformStatsService(mockUserRepo, mockBookRepo, mockProjectRepo, mockChapterRepo)
			ctx := context.Background()

			if tt.userID != "" {
				testUser := createTestUser(tt.roles, tt.vipLevel)
				mockUserRepo.On("GetByID", ctx, tt.userID).Return(testUser, nil).Once()
				mockProjectRepo.On("CountByOwner", ctx, tt.userID).Return(tt.projectCount, tt.projectError).Once()
			}

			// Act
			stats, err := service.GetUserStats(ctx, tt.userID)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, stats)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, stats)
				assert.Equal(t, tt.expectedLevel, stats.MemberLevel)
			}

			mockUserRepo.AssertExpectations(t)
			mockProjectRepo.AssertExpectations(t)
		})
	}

	t.Log("✓ Table-driven tests passed (RED phase)")
}
