package social

import (
	"context"
	"errors"
	"testing"
	"time"

	socialModel "Qingyu_backend/models/social"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentRepoState struct {
	comments           map[string]*socialModel.Comment
	failIncrementReply error
	failDecrementReply error
}

func newCommentRepoState() *commentRepoState {
	return &commentRepoState{comments: make(map[string]*socialModel.Comment)}
}

func (m *commentRepoState) Create(ctx context.Context, comment *socialModel.Comment) error {
	if comment.ID.IsZero() {
		comment.ID = primitive.NewObjectID()
	}
	m.comments[comment.ID.Hex()] = cloneComment(comment)
	return nil
}
func (m *commentRepoState) GetByID(ctx context.Context, id string) (*socialModel.Comment, error) {
	comment, ok := m.comments[id]
	if !ok {
		return nil, errors.New("comment not found")
	}
	return cloneComment(comment), nil
}
func (m *commentRepoState) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}
func (m *commentRepoState) Delete(ctx context.Context, id string) error {
	comment, ok := m.comments[id]
	if !ok {
		return errors.New("comment not found")
	}
	comment.State = socialModel.CommentStateDeleted
	comment.UpdatedAt = time.Now()
	return nil
}
func (m *commentRepoState) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.comments[id]
	return ok, nil
}
func (m *commentRepoState) GetCommentsByBookID(ctx context.Context, bookID string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *commentRepoState) GetCommentsByUserID(ctx context.Context, userID string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *commentRepoState) GetRepliesByCommentID(ctx context.Context, commentID string) ([]*socialModel.Comment, error) {
	return nil, nil
}
func (m *commentRepoState) GetCommentsByChapterID(ctx context.Context, chapterID string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *commentRepoState) ListByFilter(ctx context.Context, filter *socialModel.CommentFilter) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *commentRepoState) GetCommentsByBookIDSorted(ctx context.Context, bookID string, sortBy string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *commentRepoState) UpdateCommentStatus(ctx context.Context, id, status, reason string) error {
	return nil
}
func (m *commentRepoState) GetPendingComments(ctx context.Context, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *commentRepoState) IncrementLikeCount(ctx context.Context, id string) error { return nil }
func (m *commentRepoState) DecrementLikeCount(ctx context.Context, id string) error { return nil }
func (m *commentRepoState) IncrementReplyCount(ctx context.Context, id string) error {
	if m.failIncrementReply != nil {
		return m.failIncrementReply
	}
	comment, ok := m.comments[id]
	if !ok {
		return errors.New("comment not found")
	}
	comment.ReplyCount++
	return nil
}
func (m *commentRepoState) DecrementReplyCount(ctx context.Context, id string) error {
	if m.failDecrementReply != nil {
		return m.failDecrementReply
	}
	comment, ok := m.comments[id]
	if !ok {
		return errors.New("comment not found")
	}
	comment.ReplyCount--
	return nil
}
func (m *commentRepoState) GetBookRatingStats(ctx context.Context, bookID string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *commentRepoState) GetCommentCount(ctx context.Context, bookID string) (int64, error) {
	return 0, nil
}
func (m *commentRepoState) GetCommentsByIDs(ctx context.Context, ids []string) ([]*socialModel.Comment, error) {
	return nil, nil
}
func (m *commentRepoState) DeleteCommentsByBookID(ctx context.Context, bookID string) error {
	return nil
}
func (m *commentRepoState) Health(ctx context.Context) error { return nil }
func (m *commentRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	snapshot := cloneCommentMap(m.comments)
	if err := fn(ctx); err != nil {
		m.comments = snapshot
		return err
	}
	return nil
}

func TestReplyCommentRollbackOnReplyCountFailure(t *testing.T) {
	repo := newCommentRepoState()
	parentID := primitive.NewObjectID()
	repo.comments[parentID.Hex()] = &socialModel.Comment{
		IdentifiedEntity: socialModel.IdentifiedEntity{ID: parentID},
		AuthorID:         "parent-user",
		TargetID:         primitive.NewObjectID().Hex(),
		State:            socialModel.CommentStateNormal,
	}
	repo.failIncrementReply = errors.New("mock reply count failure")

	service := NewCommentService(repo, nil, nil)

	reply, err := service.ReplyComment(context.Background(), "reply-user", parentID.Hex(), "这是一个足够长的回复内容")
	assert.Error(t, err)
	assert.Nil(t, reply)
	assert.Len(t, repo.comments, 1)
	assert.Equal(t, int64(0), repo.comments[parentID.Hex()].ReplyCount)
}

func TestDeleteReplyRollbackOnReplyCountFailure(t *testing.T) {
	repo := newCommentRepoState()
	parentID := primitive.NewObjectID()
	replyID := primitive.NewObjectID()
	parentIDHex := parentID.Hex()
	repo.comments[parentIDHex] = &socialModel.Comment{
		IdentifiedEntity:     socialModel.IdentifiedEntity{ID: parentID},
		AuthorID:             "parent-user",
		TargetID:             primitive.NewObjectID().Hex(),
		State:                socialModel.CommentStateNormal,
		ThreadedConversation: socialModel.ThreadedConversation{ReplyCount: 1},
	}
	repo.comments[replyID.Hex()] = &socialModel.Comment{
		IdentifiedEntity: socialModel.IdentifiedEntity{ID: replyID},
		AuthorID:         "reply-user",
		TargetID:         primitive.NewObjectID().Hex(),
		State:            socialModel.CommentStateNormal,
		ThreadedConversation: socialModel.ThreadedConversation{
			ParentID: &parentIDHex,
		},
	}
	repo.failDecrementReply = errors.New("mock decrement failure")

	service := NewCommentService(repo, nil, nil)

	err := service.DeleteComment(context.Background(), "reply-user", replyID.Hex())
	assert.Error(t, err)
	assert.Equal(t, socialModel.CommentStateNormal, repo.comments[replyID.Hex()].State)
	assert.Equal(t, int64(1), repo.comments[parentIDHex].ReplyCount)
}

func cloneComment(comment *socialModel.Comment) *socialModel.Comment {
	if comment == nil {
		return nil
	}
	cloned := *comment
	if comment.ParentID != nil {
		parent := *comment.ParentID
		cloned.ParentID = &parent
	}
	if comment.RootID != nil {
		root := *comment.RootID
		cloned.RootID = &root
	}
	if comment.ReplyToUserID != nil {
		replyTo := *comment.ReplyToUserID
		cloned.ReplyToUserID = &replyTo
	}
	return &cloned
}

func cloneCommentMap(source map[string]*socialModel.Comment) map[string]*socialModel.Comment {
	cloned := make(map[string]*socialModel.Comment, len(source))
	for key, value := range source {
		cloned[key] = cloneComment(value)
	}
	return cloned
}
