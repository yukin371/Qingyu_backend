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

type likeRepoState struct {
	likes      map[string]*socialModel.Like
	failAdd    error
	failRemove error
}

func newLikeRepoState() *likeRepoState {
	return &likeRepoState{likes: make(map[string]*socialModel.Like)}
}

func (m *likeRepoState) AddLike(ctx context.Context, like *socialModel.Like) error {
	if m.failAdd != nil {
		return m.failAdd
	}
	if like.ID.IsZero() {
		like.ID = primitive.NewObjectID()
	}
	key := likeKey(like.UserID, like.TargetType, like.TargetID)
	if _, exists := m.likes[key]; exists {
		return errors.New("已经点赞过了")
	}
	cloned := *like
	m.likes[key] = &cloned
	return nil
}

func (m *likeRepoState) RemoveLike(ctx context.Context, userID, targetType, targetID string) error {
	if m.failRemove != nil {
		return m.failRemove
	}
	key := likeKey(userID, targetType, targetID)
	if _, exists := m.likes[key]; !exists {
		return errors.New("点赞记录不存在")
	}
	delete(m.likes, key)
	return nil
}

func (m *likeRepoState) IsLiked(ctx context.Context, userID, targetType, targetID string) (bool, error) {
	_, ok := m.likes[likeKey(userID, targetType, targetID)]
	return ok, nil
}

func (m *likeRepoState) GetByID(ctx context.Context, id string) (*socialModel.Like, error) {
	return nil, nil
}

func (m *likeRepoState) GetUserLikes(ctx context.Context, userID, targetType string, page, size int) ([]*socialModel.Like, int64, error) {
	return nil, 0, nil
}

func (m *likeRepoState) GetLikeCount(ctx context.Context, targetType, targetID string) (int64, error) {
	var count int64
	for _, like := range m.likes {
		if like.TargetType == targetType && like.TargetID == targetID {
			count++
		}
	}
	return count, nil
}

func (m *likeRepoState) GetLikesCountBatch(ctx context.Context, targetType string, targetIDs []string) (map[string]int64, error) {
	return nil, nil
}

func (m *likeRepoState) GetUserLikeStatusBatch(ctx context.Context, userID, targetType string, targetIDs []string) (map[string]bool, error) {
	return nil, nil
}

func (m *likeRepoState) CountUserLikes(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *likeRepoState) CountTargetLikes(ctx context.Context, targetType, targetID string) (int64, error) {
	return 0, nil
}

func (m *likeRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	snapshot := cloneLikeMap(m.likes)
	if err := fn(ctx); err != nil {
		m.likes = snapshot
		return err
	}
	return nil
}

func (m *likeRepoState) Health(ctx context.Context) error { return nil }

type likeCommentRepoState struct {
	comments      map[string]*socialModel.Comment
	failIncrement error
	failDecrement error
}

func newLikeCommentRepoState() *likeCommentRepoState {
	return &likeCommentRepoState{comments: make(map[string]*socialModel.Comment)}
}

func (m *likeCommentRepoState) Create(ctx context.Context, comment *socialModel.Comment) error {
	return nil
}
func (m *likeCommentRepoState) GetByID(ctx context.Context, id string) (*socialModel.Comment, error) {
	comment, ok := m.comments[id]
	if !ok {
		return nil, nil
	}
	cloned := *comment
	return &cloned, nil
}
func (m *likeCommentRepoState) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}
func (m *likeCommentRepoState) Delete(ctx context.Context, id string) error { return nil }
func (m *likeCommentRepoState) GetCommentsByBookID(ctx context.Context, bookID string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *likeCommentRepoState) GetCommentsByBookIDSorted(ctx context.Context, bookID, sortBy string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *likeCommentRepoState) GetCommentsByChapterID(ctx context.Context, chapterID string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *likeCommentRepoState) GetCommentsByIDs(ctx context.Context, ids []string) ([]*socialModel.Comment, error) {
	return nil, nil
}
func (m *likeCommentRepoState) GetCommentsByUserID(ctx context.Context, userID string, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *likeCommentRepoState) GetRepliesByCommentID(ctx context.Context, commentID string) ([]*socialModel.Comment, error) {
	return nil, nil
}
func (m *likeCommentRepoState) UpdateCommentStatus(ctx context.Context, id, status, reason string) error {
	return nil
}
func (m *likeCommentRepoState) GetPendingComments(ctx context.Context, page, size int) ([]*socialModel.Comment, int64, error) {
	return nil, 0, nil
}
func (m *likeCommentRepoState) IncrementLikeCount(ctx context.Context, id string) error {
	if m.failIncrement != nil {
		return m.failIncrement
	}
	if comment, ok := m.comments[id]; ok {
		comment.LikeCount++
		comment.UpdatedAt = time.Now()
	}
	return nil
}
func (m *likeCommentRepoState) DecrementLikeCount(ctx context.Context, id string) error {
	if m.failDecrement != nil {
		return m.failDecrement
	}
	if comment, ok := m.comments[id]; ok && comment.LikeCount > 0 {
		comment.LikeCount--
		comment.UpdatedAt = time.Now()
	}
	return nil
}
func (m *likeCommentRepoState) IncrementReplyCount(ctx context.Context, id string) error { return nil }
func (m *likeCommentRepoState) DecrementReplyCount(ctx context.Context, id string) error { return nil }
func (m *likeCommentRepoState) GetBookRatingStats(ctx context.Context, bookID string) (map[string]interface{}, error) {
	return nil, nil
}
func (m *likeCommentRepoState) DeleteCommentsByBookID(ctx context.Context, bookID string) error {
	return nil
}
func (m *likeCommentRepoState) GetCommentCount(ctx context.Context, bookID string) (int64, error) {
	return 0, nil
}
func (m *likeCommentRepoState) Health(ctx context.Context) error { return nil }
func (m *likeCommentRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
func (m *likeCommentRepoState) Exists(ctx context.Context, id string) (bool, error) {
	_, ok := m.comments[id]
	return ok, nil
}

func TestLikeCommentRollbackOnLikeCountFailure(t *testing.T) {
	likeRepo := newLikeRepoState()
	commentRepo := newLikeCommentRepoState()
	commentRepo.comments["comment-1"] = &socialModel.Comment{
		IdentifiedEntity: socialModel.IdentifiedEntity{ID: primitive.NewObjectID()},
		Likable:          socialModel.Likable{LikeCount: 0},
	}
	commentRepo.failIncrement = errors.New("mock increment failure")

	service := NewLikeService(likeRepo, commentRepo, nil)

	err := service.LikeComment(context.Background(), "user-1", "comment-1")
	assert.Error(t, err)
	assert.Empty(t, likeRepo.likes)
	assert.EqualValues(t, 0, commentRepo.comments["comment-1"].LikeCount)
}

func TestUnlikeCommentRollbackOnLikeCountFailure(t *testing.T) {
	likeRepo := newLikeRepoState()
	commentRepo := newLikeCommentRepoState()
	commentRepo.comments["comment-1"] = &socialModel.Comment{
		IdentifiedEntity: socialModel.IdentifiedEntity{ID: primitive.NewObjectID()},
		Likable:          socialModel.Likable{LikeCount: 1},
	}
	likeRepo.likes[likeKey("user-1", socialModel.LikeTargetTypeComment, "comment-1")] = &socialModel.Like{
		ID:         primitive.NewObjectID(),
		UserID:     "user-1",
		TargetType: socialModel.LikeTargetTypeComment,
		TargetID:   "comment-1",
		CreatedAt:  time.Now(),
	}
	commentRepo.failDecrement = errors.New("mock decrement failure")

	service := NewLikeService(likeRepo, commentRepo, nil)

	err := service.UnlikeComment(context.Background(), "user-1", "comment-1")
	assert.Error(t, err)
	assert.Len(t, likeRepo.likes, 1)
	assert.EqualValues(t, 1, commentRepo.comments["comment-1"].LikeCount)
}

func likeKey(userID, targetType, targetID string) string {
	return userID + ":" + targetType + ":" + targetID
}

func cloneLikeMap(source map[string]*socialModel.Like) map[string]*socialModel.Like {
	cloned := make(map[string]*socialModel.Like, len(source))
	for key, value := range source {
		if value == nil {
			cloned[key] = nil
			continue
		}
		copyValue := *value
		cloned[key] = &copyValue
	}
	return cloned
}
