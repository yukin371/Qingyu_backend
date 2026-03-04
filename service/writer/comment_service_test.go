package writer

import (
	"context"
	"testing"

	models "Qingyu_backend/models/writer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockCommentRepository struct {
	mock.Mock
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *models.DocumentComment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}
func (m *mockCommentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.DocumentComment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DocumentComment), args.Error(1)
}
func (m *mockCommentRepository) Update(ctx context.Context, id primitive.ObjectID, comment *models.DocumentComment) error {
	args := m.Called(ctx, id, comment)
	return args.Error(0)
}
func (m *mockCommentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockCommentRepository) HardDelete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockCommentRepository) List(ctx context.Context, filter *models.CommentFilter, page, pageSize int) ([]*models.DocumentComment, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.DocumentComment), args.Get(1).(int64), args.Error(2)
}
func (m *mockCommentRepository) GetByDocument(ctx context.Context, documentID primitive.ObjectID, includeResolved bool) ([]*models.DocumentComment, error) {
	args := m.Called(ctx, documentID, includeResolved)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DocumentComment), args.Error(1)
}
func (m *mockCommentRepository) GetByChapter(ctx context.Context, chapterID primitive.ObjectID, includeResolved bool) ([]*models.DocumentComment, error) {
	args := m.Called(ctx, chapterID, includeResolved)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DocumentComment), args.Error(1)
}
func (m *mockCommentRepository) GetThread(ctx context.Context, threadID primitive.ObjectID) (*models.CommentThread, error) {
	args := m.Called(ctx, threadID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CommentThread), args.Error(1)
}
func (m *mockCommentRepository) GetReplies(ctx context.Context, parentID primitive.ObjectID) ([]*models.DocumentComment, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DocumentComment), args.Error(1)
}
func (m *mockCommentRepository) MarkAsResolved(ctx context.Context, id primitive.ObjectID, resolvedBy primitive.ObjectID) error {
	args := m.Called(ctx, id, resolvedBy)
	return args.Error(0)
}
func (m *mockCommentRepository) MarkAsUnresolved(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockCommentRepository) GetStats(ctx context.Context, documentID primitive.ObjectID) (*models.CommentStats, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CommentStats), args.Error(1)
}
func (m *mockCommentRepository) GetUserComments(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*models.DocumentComment, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.DocumentComment), args.Get(1).(int64), args.Error(2)
}
func (m *mockCommentRepository) BatchDelete(ctx context.Context, ids []primitive.ObjectID) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}
func (m *mockCommentRepository) Search(ctx context.Context, keyword string, documentID primitive.ObjectID, page, pageSize int) ([]*models.DocumentComment, int64, error) {
	args := m.Called(ctx, keyword, documentID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.DocumentComment), args.Get(1).(int64), args.Error(2)
}

func TestCommentService_CreateComment_RequiresParagraphID(t *testing.T) {
	repo := new(mockCommentRepository)
	svc := NewCommentService(repo)

	comment := &models.DocumentComment{
		DocumentID: primitive.NewObjectID(),
		UserID:     primitive.NewObjectID(),
		Content:    "test",
		Type:       models.CommentTypeComment,
	}

	created, err := svc.CreateComment(context.Background(), comment)
	assert.Error(t, err)
	assert.Nil(t, created)
	assert.True(t, IsErrorCode(err, ErrInvalidInput))
}

func TestCommentService_ReplyComment_InheritsParagraphID(t *testing.T) {
	repo := new(mockCommentRepository)
	svc := NewCommentService(repo)

	parentID := primitive.NewObjectID()
	paragraphID := primitive.NewObjectID()
	threadID := primitive.NewObjectID()
	docID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	parent := &models.DocumentComment{
		ID:          parentID,
		DocumentID:  docID,
		ParagraphID: paragraphID,
		ThreadID:    &threadID,
	}

	repo.On("GetByID", mock.Anything, parentID).Return(parent, nil).Once()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(c *models.DocumentComment) bool {
		return c != nil && c.ParagraphID == paragraphID && c.DocumentID == docID
	})).Return(nil).Once()

	reply, err := svc.ReplyComment(context.Background(), parentID.Hex(), "reply", userID.Hex(), "tester")
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Equal(t, paragraphID, reply.ParagraphID)

	repo.AssertExpectations(t)
}
