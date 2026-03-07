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

type reviewRepoState struct {
	reviews           map[string]*socialModel.Review
	likes             map[string]*socialModel.ReviewLike
	failCreateLike    error
	failDeleteLike    error
	failIncrementLike error
	failDecrementLike error
}

func newReviewRepoState() *reviewRepoState {
	return &reviewRepoState{
		reviews: make(map[string]*socialModel.Review),
		likes:   make(map[string]*socialModel.ReviewLike),
	}
}

func (m *reviewRepoState) CreateReview(ctx context.Context, review *socialModel.Review) error {
	return nil
}
func (m *reviewRepoState) GetReviewByID(ctx context.Context, reviewID string) (*socialModel.Review, error) {
	review, ok := m.reviews[reviewID]
	if !ok {
		return nil, nil
	}
	cloned := *review
	return &cloned, nil
}
func (m *reviewRepoState) GetReviewsByBook(ctx context.Context, bookID string, page, size int) ([]*socialModel.Review, int64, error) {
	return nil, 0, nil
}
func (m *reviewRepoState) GetReviewsByUser(ctx context.Context, userID string, page, size int) ([]*socialModel.Review, int64, error) {
	return nil, 0, nil
}
func (m *reviewRepoState) GetPublicReviews(ctx context.Context, page, size int) ([]*socialModel.Review, int64, error) {
	return nil, 0, nil
}
func (m *reviewRepoState) GetReviewsByRating(ctx context.Context, bookID string, rating int, page, size int) ([]*socialModel.Review, int64, error) {
	return nil, 0, nil
}
func (m *reviewRepoState) UpdateReview(ctx context.Context, reviewID string, updates map[string]interface{}) error {
	return nil
}
func (m *reviewRepoState) DeleteReview(ctx context.Context, reviewID string) error { return nil }
func (m *reviewRepoState) CreateReviewLike(ctx context.Context, reviewLike *socialModel.ReviewLike) error {
	if m.failCreateLike != nil {
		return m.failCreateLike
	}
	if reviewLike.ID.IsZero() {
		reviewLike.ID = primitive.NewObjectID()
	}
	cloned := *reviewLike
	m.likes[reviewLike.ReviewID+":"+reviewLike.UserID] = &cloned
	return nil
}
func (m *reviewRepoState) DeleteReviewLike(ctx context.Context, reviewID, userID string) error {
	if m.failDeleteLike != nil {
		return m.failDeleteLike
	}
	delete(m.likes, reviewID+":"+userID)
	return nil
}
func (m *reviewRepoState) GetReviewLike(ctx context.Context, reviewID, userID string) (*socialModel.ReviewLike, error) {
	like, ok := m.likes[reviewID+":"+userID]
	if !ok {
		return nil, nil
	}
	cloned := *like
	return &cloned, nil
}
func (m *reviewRepoState) IsReviewLiked(ctx context.Context, reviewID, userID string) (bool, error) {
	_, ok := m.likes[reviewID+":"+userID]
	return ok, nil
}
func (m *reviewRepoState) GetReviewLikes(ctx context.Context, reviewID string, page, size int) ([]*socialModel.ReviewLike, int64, error) {
	return nil, 0, nil
}
func (m *reviewRepoState) IncrementReviewLikeCount(ctx context.Context, reviewID string) error {
	if m.failIncrementLike != nil {
		return m.failIncrementLike
	}
	if review, ok := m.reviews[reviewID]; ok {
		review.LikeCount++
	}
	return nil
}
func (m *reviewRepoState) DecrementReviewLikeCount(ctx context.Context, reviewID string) error {
	if m.failDecrementLike != nil {
		return m.failDecrementLike
	}
	if review, ok := m.reviews[reviewID]; ok && review.LikeCount > 0 {
		review.LikeCount--
	}
	return nil
}
func (m *reviewRepoState) GetAverageRating(ctx context.Context, bookID string) (float64, error) {
	return 0, nil
}
func (m *reviewRepoState) GetRatingDistribution(ctx context.Context, bookID string) (map[int]int64, error) {
	return nil, nil
}
func (m *reviewRepoState) CountReviews(ctx context.Context, bookID string) (int64, error) {
	return 0, nil
}
func (m *reviewRepoState) CountUserReviews(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}
func (m *reviewRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	reviewSnapshot := cloneReviewMap(m.reviews)
	likeSnapshot := cloneReviewLikeMap(m.likes)
	if err := fn(ctx); err != nil {
		m.reviews = reviewSnapshot
		m.likes = likeSnapshot
		return err
	}
	return nil
}
func (m *reviewRepoState) Health(ctx context.Context) error { return nil }

func TestLikeReviewRollbackOnLikeCountFailure(t *testing.T) {
	repo := newReviewRepoState()
	repo.reviews["review-1"] = &socialModel.Review{
		ID:        primitive.NewObjectID(),
		BookID:    "book-1",
		UserID:    "owner-1",
		Title:     "review",
		Content:   "content",
		Rating:    5,
		LikeCount: 0,
	}
	repo.failIncrementLike = errors.New("mock increment failure")

	service := NewReviewService(repo, nil)

	err := service.LikeReview(context.Background(), "user-1", "review-1")
	assert.Error(t, err)
	assert.Empty(t, repo.likes)
	assert.EqualValues(t, 0, repo.reviews["review-1"].LikeCount)
}

func TestUnlikeReviewRollbackOnLikeCountFailure(t *testing.T) {
	repo := newReviewRepoState()
	repo.reviews["review-1"] = &socialModel.Review{
		ID:        primitive.NewObjectID(),
		BookID:    "book-1",
		UserID:    "owner-1",
		Title:     "review",
		Content:   "content",
		Rating:    5,
		LikeCount: 1,
	}
	repo.likes["review-1:user-1"] = &socialModel.ReviewLike{
		ID:        primitive.NewObjectID(),
		ReviewID:  "review-1",
		UserID:    "user-1",
		CreatedAt: time.Now(),
	}
	repo.failDecrementLike = errors.New("mock decrement failure")

	service := NewReviewService(repo, nil)

	err := service.UnlikeReview(context.Background(), "user-1", "review-1")
	assert.Error(t, err)
	assert.Len(t, repo.likes, 1)
	assert.EqualValues(t, 1, repo.reviews["review-1"].LikeCount)
}

func cloneReviewMap(source map[string]*socialModel.Review) map[string]*socialModel.Review {
	cloned := make(map[string]*socialModel.Review, len(source))
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

func cloneReviewLikeMap(source map[string]*socialModel.ReviewLike) map[string]*socialModel.ReviewLike {
	cloned := make(map[string]*socialModel.ReviewLike, len(source))
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
