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

type followRepoState struct {
	follows          map[string]*socialModel.Follow
	failMutualUpdate error
	failDeleteMutual error
	failStats        error
}

func newFollowRepoState() *followRepoState {
	return &followRepoState{follows: make(map[string]*socialModel.Follow)}
}

func (m *followRepoState) CreateFollow(ctx context.Context, follow *socialModel.Follow) error {
	if follow.ID.IsZero() {
		follow.ID = primitive.NewObjectID()
	}
	m.follows[m.followKey(follow.FollowerID, follow.FollowingID, follow.FollowType)] = cloneFollow(follow)
	return nil
}

func (m *followRepoState) DeleteFollow(ctx context.Context, followerID, followingID, followType string) error {
	delete(m.follows, m.followKey(followerID, followingID, followType))
	return nil
}

func (m *followRepoState) GetFollow(ctx context.Context, followerID, followingID, followType string) (*socialModel.Follow, error) {
	follow, ok := m.follows[m.followKey(followerID, followingID, followType)]
	if !ok {
		return nil, nil
	}
	return cloneFollow(follow), nil
}

func (m *followRepoState) IsFollowing(ctx context.Context, followerID, followingID, followType string) (bool, error) {
	_, ok := m.follows[m.followKey(followerID, followingID, followType)]
	return ok, nil
}

func (m *followRepoState) GetFollowers(ctx context.Context, userID string, followType string, page, size int) ([]*socialModel.FollowInfo, int64, error) {
	return nil, 0, nil
}

func (m *followRepoState) GetFollowing(ctx context.Context, userID string, followType string, page, size int) ([]*socialModel.FollowingInfo, int64, error) {
	return nil, 0, nil
}

func (m *followRepoState) UpdateMutualStatus(ctx context.Context, followerID, followingID, followType string, isMutual bool) error {
	if isMutual && m.failMutualUpdate != nil {
		return m.failMutualUpdate
	}
	if !isMutual && m.failDeleteMutual != nil {
		return m.failDeleteMutual
	}
	follow, ok := m.follows[m.followKey(followerID, followingID, followType)]
	if !ok {
		return nil
	}
	follow.IsMutual = isMutual
	follow.UpdatedAt = time.Now()
	return nil
}

func (m *followRepoState) CreateAuthorFollow(ctx context.Context, authorFollow *socialModel.AuthorFollow) error {
	return nil
}
func (m *followRepoState) DeleteAuthorFollow(ctx context.Context, userID, authorID string) error {
	return nil
}
func (m *followRepoState) GetAuthorFollow(ctx context.Context, userID, authorID string) (*socialModel.AuthorFollow, error) {
	return nil, nil
}
func (m *followRepoState) GetAuthorFollowers(ctx context.Context, authorID string, page, size int) ([]*socialModel.FollowInfo, int64, error) {
	return nil, 0, nil
}
func (m *followRepoState) GetUserFollowingAuthors(ctx context.Context, userID string, page, size int) ([]*socialModel.AuthorFollow, int64, error) {
	return nil, 0, nil
}
func (m *followRepoState) GetFollowStats(ctx context.Context, userID string) (*socialModel.FollowStats, error) {
	return nil, nil
}
func (m *followRepoState) UpdateFollowStats(ctx context.Context, userID string, followerDelta, followingDelta int) error {
	return m.failStats
}
func (m *followRepoState) CountFollowers(ctx context.Context, userID, followType string) (int64, error) {
	return 0, nil
}
func (m *followRepoState) CountFollowing(ctx context.Context, userID, followType string) (int64, error) {
	return 0, nil
}
func (m *followRepoState) Health(ctx context.Context) error { return nil }
func (m *followRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	snapshot := cloneFollowMap(m.follows)
	if err := fn(ctx); err != nil {
		m.follows = snapshot
		return err
	}
	return nil
}

func (m *followRepoState) followKey(followerID, followingID, followType string) string {
	return followerID + ":" + followingID + ":" + followType
}

func TestFollowUserRollbackOnMutualUpdateFailure(t *testing.T) {
	repo := newFollowRepoState()
	repo.follows[repo.followKey("u2", "u1", "user")] = &socialModel.Follow{
		ID:          primitive.NewObjectID(),
		FollowerID:  "u2",
		FollowingID: "u1",
		FollowType:  "user",
		IsMutual:    false,
	}
	repo.failMutualUpdate = errors.New("mock mutual update failure")

	service := NewFollowService(repo, nil)

	err := service.FollowUser(context.Background(), "u1", "u2")
	assert.Error(t, err)
	assert.False(t, hasFollow(repo.follows, repo.followKey("u1", "u2", "user")))
	assert.False(t, repo.follows[repo.followKey("u2", "u1", "user")].IsMutual)
}

func TestUnfollowUserRollbackOnMutualClearFailure(t *testing.T) {
	repo := newFollowRepoState()
	repo.follows[repo.followKey("u1", "u2", "user")] = &socialModel.Follow{
		ID:          primitive.NewObjectID(),
		FollowerID:  "u1",
		FollowingID: "u2",
		FollowType:  "user",
		IsMutual:    true,
	}
	repo.follows[repo.followKey("u2", "u1", "user")] = &socialModel.Follow{
		ID:          primitive.NewObjectID(),
		FollowerID:  "u2",
		FollowingID: "u1",
		FollowType:  "user",
		IsMutual:    true,
	}
	repo.failDeleteMutual = errors.New("mock mutual clear failure")

	service := NewFollowService(repo, nil)

	err := service.UnfollowUser(context.Background(), "u1", "u2")
	assert.Error(t, err)
	assert.True(t, hasFollow(repo.follows, repo.followKey("u1", "u2", "user")))
	assert.True(t, repo.follows[repo.followKey("u2", "u1", "user")].IsMutual)
}

func hasFollow(source map[string]*socialModel.Follow, key string) bool {
	_, ok := source[key]
	return ok
}

func cloneFollow(follow *socialModel.Follow) *socialModel.Follow {
	if follow == nil {
		return nil
	}
	cloned := *follow
	return &cloned
}

func cloneFollowMap(source map[string]*socialModel.Follow) map[string]*socialModel.Follow {
	cloned := make(map[string]*socialModel.Follow, len(source))
	for key, value := range source {
		cloned[key] = cloneFollow(value)
	}
	return cloned
}
