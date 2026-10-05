package reader

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	readerModels "Qingyu_backend/models/reader"
)

type bookmarkRepositorySecurityStub struct {
	bookmark *readerModels.Bookmark
	updated  bool
	deleted  bool
}

func (r *bookmarkRepositorySecurityStub) Create(context.Context, *readerModels.Bookmark) error {
	return nil
}

func (r *bookmarkRepositorySecurityStub) GetByID(context.Context, primitive.ObjectID) (*readerModels.Bookmark, error) {
	if r.bookmark == nil {
		return nil, errors.New("not found")
	}
	return r.bookmark, nil
}

func (r *bookmarkRepositorySecurityStub) GetByUser(context.Context, primitive.ObjectID, *readerModels.BookmarkFilter, int, int) ([]*readerModels.Bookmark, int64, error) {
	return nil, 0, nil
}

func (r *bookmarkRepositorySecurityStub) GetByBook(context.Context, primitive.ObjectID, primitive.ObjectID, int, int) ([]*readerModels.Bookmark, int64, error) {
	return nil, 0, nil
}

func (r *bookmarkRepositorySecurityStub) GetByChapter(context.Context, primitive.ObjectID, primitive.ObjectID, int, int) ([]*readerModels.Bookmark, int64, error) {
	return nil, 0, nil
}

func (r *bookmarkRepositorySecurityStub) Update(context.Context, *readerModels.Bookmark) error {
	r.updated = true
	return nil
}

func (r *bookmarkRepositorySecurityStub) Delete(context.Context, primitive.ObjectID) error {
	r.deleted = true
	return nil
}

func (r *bookmarkRepositorySecurityStub) DeleteByBook(context.Context, primitive.ObjectID, primitive.ObjectID) error {
	return nil
}

func (r *bookmarkRepositorySecurityStub) Exists(context.Context, primitive.ObjectID, primitive.ObjectID, int) (bool, error) {
	return false, nil
}

func (r *bookmarkRepositorySecurityStub) GetStats(context.Context, primitive.ObjectID) (*readerModels.BookmarkStats, error) {
	return nil, nil
}

func (r *bookmarkRepositorySecurityStub) GetPublicBookmarks(context.Context, int, int) ([]*readerModels.Bookmark, int64, error) {
	return nil, 0, nil
}

func (r *bookmarkRepositorySecurityStub) Search(context.Context, primitive.ObjectID, string, int, int) ([]*readerModels.Bookmark, int64, error) {
	return nil, 0, nil
}

func TestBookmarkService_OwnerScopeRejectsOtherUser(t *testing.T) {
	ownerID := primitive.NewObjectID()
	requesterID := primitive.NewObjectID()
	bookmark := &readerModels.Bookmark{ID: primitive.NewObjectID(), UserID: ownerID}
	repo := &bookmarkRepositorySecurityStub{bookmark: bookmark}
	service := NewBookmarkService(repo, nil)

	_, err := service.GetBookmark(context.Background(), requesterID.Hex(), bookmark.ID.Hex())
	require.ErrorIs(t, err, ErrBookmarkNotFound)

	err = service.UpdateBookmark(context.Background(), requesterID.Hex(), bookmark.ID.Hex(), &readerModels.Bookmark{Note: "越权"})
	require.ErrorIs(t, err, ErrBookmarkNotFound)
	require.False(t, repo.updated)

	err = service.DeleteBookmark(context.Background(), requesterID.Hex(), bookmark.ID.Hex())
	require.ErrorIs(t, err, ErrBookmarkNotFound)
	require.False(t, repo.deleted)
}

func TestBookmarkService_OwnerScopeAllowsOwner(t *testing.T) {
	ownerID := primitive.NewObjectID()
	bookmark := &readerModels.Bookmark{ID: primitive.NewObjectID(), UserID: ownerID, Note: "旧"}
	repo := &bookmarkRepositorySecurityStub{bookmark: bookmark}
	service := NewBookmarkService(repo, nil)

	got, err := service.GetBookmark(context.Background(), ownerID.Hex(), bookmark.ID.Hex())
	require.NoError(t, err)
	require.Equal(t, bookmark.ID, got.ID)

	err = service.UpdateBookmark(context.Background(), ownerID.Hex(), bookmark.ID.Hex(), &readerModels.Bookmark{Note: "新"})
	require.NoError(t, err)
	require.True(t, repo.updated)

	err = service.DeleteBookmark(context.Background(), ownerID.Hex(), bookmark.ID.Hex())
	require.NoError(t, err)
	require.True(t, repo.deleted)
}
