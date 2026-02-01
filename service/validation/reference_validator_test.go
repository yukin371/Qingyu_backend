package validation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo 模拟用户仓库
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Exists(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

// MockBookRepo 模拟书籍仓库
type MockBookRepo struct {
	mock.Mock
}

func (m *MockBookRepo) Exists(ctx context.Context, bookID string) (bool, error) {
	args := m.Called(ctx, bookID)
	return args.Bool(0), args.Error(1)
}

// MockCommentRepo 模拟评论仓库
type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) Exists(ctx context.Context, commentID string) (bool, error) {
	args := m.Called(ctx, commentID)
	return args.Bool(0), args.Error(1)
}

// =======================
// 用户验证测试
// =======================

func TestReferenceValidator_ValidateUserExists_ValidUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)

	// Act
	err := validator.ValidateUserExists(context.Background(), "valid-user-id")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateUserExists_InvalidUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "invalid-user-id").Return(false, nil)

	// Act
	err := validator.ValidateUserExists(context.Background(), "invalid-user-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "用户不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateUserExists_RepositoryError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "error-user-id").Return(false, errors.New("数据库连接失败"))

	// Act
	err := validator.ValidateUserExists(context.Background(), "error-user-id")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "检查用户失败")
	mockUserRepo.AssertExpectations(t)
}

// =======================
// 书籍验证测试
// =======================

func TestReferenceValidator_ValidateBookExists_ValidBook(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(nil, mockBookRepo, nil, nil)

	mockBookRepo.On("Exists", mock.Anything, "valid-book-id").Return(true, nil)

	// Act
	err := validator.ValidateBookExists(context.Background(), "valid-book-id")

	// Assert
	assert.NoError(t, err)
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateBookExists_InvalidBook(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(nil, mockBookRepo, nil, nil)

	mockBookRepo.On("Exists", mock.Anything, "invalid-book-id").Return(false, nil)

	// Act
	err := validator.ValidateBookExists(context.Background(), "invalid-book-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "书籍不存在", err.Error())
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateBookExists_RepositoryError(t *testing.T) {
	// Arrange
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(nil, mockBookRepo, nil, nil)

	mockBookRepo.On("Exists", mock.Anything, "error-book-id").Return(false, errors.New("数据库连接失败"))

	// Act
	err := validator.ValidateBookExists(context.Background(), "error-book-id")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "检查书籍失败")
	mockBookRepo.AssertExpectations(t)
}

// =======================
// 点赞引用验证测试
// =======================

func TestReferenceValidator_ValidateLikeReference_ValidBookLike(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "valid-book-id").Return(true, nil)

	// Act
	err := validator.ValidateLikeReference(context.Background(), "valid-user-id", "valid-book-id", "book")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateLikeReference_InvalidUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "invalid-user-id").Return(false, nil)

	// Act
	err := validator.ValidateLikeReference(context.Background(), "invalid-user-id", "valid-book-id", "book")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "用户不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateLikeReference_InvalidBook(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "invalid-book-id").Return(false, nil)

	// Act
	err := validator.ValidateLikeReference(context.Background(), "valid-user-id", "invalid-book-id", "book")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "书籍不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateLikeReference_InvalidTargetType(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)

	// Act
	err := validator.ValidateLikeReference(context.Background(), "valid-user-id", "target-id", "invalid-type")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "无效的目标类型", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateLikeReference_ValidCommentLike(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockCommentRepo := new(MockCommentRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, mockCommentRepo)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)
	mockCommentRepo.On("Exists", mock.Anything, "valid-comment-id").Return(true, nil)

	// Act
	err := validator.ValidateLikeReference(context.Background(), "valid-user-id", "valid-comment-id", "comment")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateLikeReference_InvalidComment(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockCommentRepo := new(MockCommentRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, mockCommentRepo)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)
	mockCommentRepo.On("Exists", mock.Anything, "invalid-comment-id").Return(false, nil)

	// Act
	err := validator.ValidateLikeReference(context.Background(), "valid-user-id", "invalid-comment-id", "comment")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "评论不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

// =======================
// 评论引用验证测试
// =======================

func TestReferenceValidator_ValidateCommentReference_Valid(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-author-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "valid-book-id").Return(true, nil)

	// Act
	err := validator.ValidateCommentReference(context.Background(), "valid-author-id", "valid-book-id")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateCommentReference_InvalidAuthor(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "invalid-author-id").Return(false, nil)

	// Act
	err := validator.ValidateCommentReference(context.Background(), "invalid-author-id", "valid-book-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "用户不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateCommentReference_InvalidBook(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-author-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "invalid-book-id").Return(false, nil)

	// Act
	err := validator.ValidateCommentReference(context.Background(), "valid-author-id", "invalid-book-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "书籍不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateCommentReference_ReplyToComment_Valid(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockCommentRepo := new(MockCommentRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, mockCommentRepo)

	mockUserRepo.On("Exists", mock.Anything, "valid-author-id").Return(true, nil)
	mockCommentRepo.On("Exists", mock.Anything, "valid-parent-comment-id").Return(true, nil)

	// Act
	err := validator.ValidateCommentReferenceReply(context.Background(), "valid-author-id", "valid-parent-comment-id")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateCommentReference_ReplyToComment_InvalidParentComment(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockCommentRepo := new(MockCommentRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, mockCommentRepo)

	mockUserRepo.On("Exists", mock.Anything, "valid-author-id").Return(true, nil)
	mockCommentRepo.On("Exists", mock.Anything, "invalid-parent-comment-id").Return(false, nil)

	// Act
	err := validator.ValidateCommentReferenceReply(context.Background(), "valid-author-id", "invalid-parent-comment-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "父评论不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

// =======================
// 关注引用验证测试
// =======================

func TestReferenceValidator_ValidateFollowReference_Valid(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	// 设置两次调用期望（验证follower和following）
	mockUserRepo.On("Exists", mock.Anything, "follower-id").Return(true, nil)
	mockUserRepo.On("Exists", mock.Anything, "following-id").Return(true, nil)

	// Act
	err := validator.ValidateFollowReference(context.Background(), "follower-id", "following-id")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateFollowReference_InvalidFollower(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "invalid-follower-id").Return(false, nil)

	// Act
	err := validator.ValidateFollowReference(context.Background(), "invalid-follower-id", "following-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "关注者不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateFollowReference_InvalidFollowing(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-follower-id").Return(true, nil)
	mockUserRepo.On("Exists", mock.Anything, "invalid-following-id").Return(false, nil)

	// Act
	err := validator.ValidateFollowReference(context.Background(), "valid-follower-id", "invalid-following-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "被关注用户不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateFollowReference_SelfFollow(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	validator := NewReferenceValidator(mockUserRepo, nil, nil, nil)

	// 自关注检查在验证用户之前，所以不会调用Exists

	// Act
	err := validator.ValidateFollowReference(context.Background(), "same-user-id", "same-user-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "不能关注自己", err.Error())
	// 验证没有调用Exists（因为先检查自关注）
	mockUserRepo.AssertNotCalled(t, "Exists", mock.Anything, mock.Anything)
}

// =======================
// 收益记录引用验证测试
// =======================

func TestReferenceValidator_ValidateRevenueReference_Valid(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-author-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "valid-book-id").Return(true, nil)

	// Act
	err := validator.ValidateRevenueReference(context.Background(), "valid-author-id", "valid-book-id")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateRevenueReference_InvalidAuthor(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "invalid-author-id").Return(false, nil)

	// Act
	err := validator.ValidateRevenueReference(context.Background(), "invalid-author-id", "valid-book-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "用户不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateRevenueReference_InvalidBook(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-author-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "invalid-book-id").Return(false, nil)

	// Act
	err := validator.ValidateRevenueReference(context.Background(), "valid-author-id", "invalid-book-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "书籍不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

// =======================
// 阅读进度引用验证测试
// =======================

func TestReferenceValidator_ValidateReadingProgressReference_Valid(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "valid-book-id").Return(true, nil)

	// Act
	err := validator.ValidateReadingProgressReference(context.Background(), "valid-user-id", "valid-book-id")

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateReadingProgressReference_InvalidUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "invalid-user-id").Return(false, nil)

	// Act
	err := validator.ValidateReadingProgressReference(context.Background(), "invalid-user-id", "valid-book-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "用户不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestReferenceValidator_ValidateReadingProgressReference_InvalidBook(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepo)
	mockBookRepo := new(MockBookRepo)
	validator := NewReferenceValidator(mockUserRepo, mockBookRepo, nil, nil)

	mockUserRepo.On("Exists", mock.Anything, "valid-user-id").Return(true, nil)
	mockBookRepo.On("Exists", mock.Anything, "invalid-book-id").Return(false, nil)

	// Act
	err := validator.ValidateReadingProgressReference(context.Background(), "valid-user-id", "invalid-book-id")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "书籍不存在", err.Error())
	mockUserRepo.AssertExpectations(t)
	mockBookRepo.AssertExpectations(t)
}
