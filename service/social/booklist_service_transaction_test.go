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

type bookListRepoState struct {
	bookLists         map[string]*socialModel.BookList
	likes             map[string]*socialModel.BookListLike
	failCreateLike    error
	failIncrementLike error
	failFork          error
	failIncrementFork error
}

func newBookListRepoState() *bookListRepoState {
	return &bookListRepoState{
		bookLists: make(map[string]*socialModel.BookList),
		likes:     make(map[string]*socialModel.BookListLike),
	}
}

func (m *bookListRepoState) CreateBookList(ctx context.Context, bookList *socialModel.BookList) error {
	if bookList.ID.IsZero() {
		bookList.ID = primitive.NewObjectID()
	}
	cloned := *bookList
	m.bookLists[bookList.ID.Hex()] = &cloned
	return nil
}

func (m *bookListRepoState) GetBookListByID(ctx context.Context, bookListID string) (*socialModel.BookList, error) {
	bookList, ok := m.bookLists[bookListID]
	if !ok {
		return nil, nil
	}
	cloned := *bookList
	return &cloned, nil
}

func (m *bookListRepoState) GetBookListsByUser(ctx context.Context, userID string, page, size int) ([]*socialModel.BookList, int64, error) {
	return nil, 0, nil
}

func (m *bookListRepoState) GetPublicBookLists(ctx context.Context, page, size int) ([]*socialModel.BookList, int64, error) {
	return nil, 0, nil
}

func (m *bookListRepoState) GetBookListsByCategory(ctx context.Context, category string, page, size int) ([]*socialModel.BookList, int64, error) {
	return nil, 0, nil
}

func (m *bookListRepoState) GetBookListsByTag(ctx context.Context, tag string, page, size int) ([]*socialModel.BookList, int64, error) {
	return nil, 0, nil
}

func (m *bookListRepoState) SearchBookLists(ctx context.Context, keyword string, page, size int) ([]*socialModel.BookList, int64, error) {
	return nil, 0, nil
}

func (m *bookListRepoState) UpdateBookList(ctx context.Context, bookListID string, updates map[string]interface{}) error {
	return nil
}

func (m *bookListRepoState) DeleteBookList(ctx context.Context, bookListID string) error {
	return nil
}

func (m *bookListRepoState) AddBookToList(ctx context.Context, bookListID string, bookItem *socialModel.BookListItem) error {
	return nil
}

func (m *bookListRepoState) RemoveBookFromList(ctx context.Context, bookListID, bookID string) error {
	return nil
}

func (m *bookListRepoState) UpdateBookInList(ctx context.Context, bookListID, bookID string, updates map[string]interface{}) error {
	return nil
}

func (m *bookListRepoState) ReorderBooks(ctx context.Context, bookListID string, bookOrders map[string]int) error {
	return nil
}

func (m *bookListRepoState) GetBooksInList(ctx context.Context, bookListID string) ([]*socialModel.BookListItem, error) {
	return nil, nil
}

func (m *bookListRepoState) CreateBookListLike(ctx context.Context, bookListLike *socialModel.BookListLike) error {
	if m.failCreateLike != nil {
		return m.failCreateLike
	}
	if bookListLike.ID.IsZero() {
		bookListLike.ID = primitive.NewObjectID()
	}
	cloned := *bookListLike
	m.likes[bookListLike.BookListID+":"+bookListLike.UserID] = &cloned
	return nil
}

func (m *bookListRepoState) DeleteBookListLike(ctx context.Context, bookListID, userID string) error {
	delete(m.likes, bookListID+":"+userID)
	return nil
}

func (m *bookListRepoState) GetBookListLike(ctx context.Context, bookListID, userID string) (*socialModel.BookListLike, error) {
	like, ok := m.likes[bookListID+":"+userID]
	if !ok {
		return nil, nil
	}
	cloned := *like
	return &cloned, nil
}

func (m *bookListRepoState) IsBookListLiked(ctx context.Context, bookListID, userID string) (bool, error) {
	_, ok := m.likes[bookListID+":"+userID]
	return ok, nil
}

func (m *bookListRepoState) GetBookListLikes(ctx context.Context, bookListID string, page, size int) ([]*socialModel.BookListLike, int64, error) {
	return nil, 0, nil
}

func (m *bookListRepoState) IncrementBookListLikeCount(ctx context.Context, bookListID string) error {
	if m.failIncrementLike != nil {
		return m.failIncrementLike
	}
	if list, ok := m.bookLists[bookListID]; ok {
		list.LikeCount++
	}
	return nil
}

func (m *bookListRepoState) DecrementBookListLikeCount(ctx context.Context, bookListID string) error {
	if list, ok := m.bookLists[bookListID]; ok && list.LikeCount > 0 {
		list.LikeCount--
	}
	return nil
}

func (m *bookListRepoState) ForkBookList(ctx context.Context, originalID, userID string) (*socialModel.BookList, error) {
	if m.failFork != nil {
		return nil, m.failFork
	}
	original, ok := m.bookLists[originalID]
	if !ok {
		return nil, errors.New("原始书单不存在")
	}
	forked := &socialModel.BookList{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Title:       original.Title,
		Description: original.Description,
		Cover:       original.Cover,
		Books:       append([]socialModel.BookListItem(nil), original.Books...),
		Category:    original.Category,
		Tags:        append([]string(nil), original.Tags...),
		IsPublic:    original.IsPublic,
		OriginalID:  &original.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.bookLists[forked.ID.Hex()] = forked
	return forked, nil
}

func (m *bookListRepoState) IncrementForkCount(ctx context.Context, bookListID string) error {
	if m.failIncrementFork != nil {
		return m.failIncrementFork
	}
	if list, ok := m.bookLists[bookListID]; ok {
		list.ForkCount++
	}
	return nil
}

func (m *bookListRepoState) GetForkedBookLists(ctx context.Context, originalID string, page, size int) ([]*socialModel.BookList, int64, error) {
	return nil, 0, nil
}

func (m *bookListRepoState) IncrementViewCount(ctx context.Context, bookListID string) error {
	return nil
}

func (m *bookListRepoState) CountUserBookLists(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *bookListRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	bookListsSnapshot := cloneBookListMap(m.bookLists)
	likesSnapshot := cloneBookListLikeMap(m.likes)
	if err := fn(ctx); err != nil {
		m.bookLists = bookListsSnapshot
		m.likes = likesSnapshot
		return err
	}
	return nil
}

func (m *bookListRepoState) Health(ctx context.Context) error { return nil }

func TestLikeBookListRollbackOnLikeCountFailure(t *testing.T) {
	repo := newBookListRepoState()
	repo.bookLists["booklist-1"] = &socialModel.BookList{
		ID:        primitive.NewObjectID(),
		UserID:    "owner-1",
		Title:     "list",
		LikeCount: 0,
	}
	repo.failIncrementLike = errors.New("mock like count failure")

	service := NewBookListService(repo, nil)

	err := service.LikeBookList(context.Background(), "user-1", "booklist-1")
	assert.Error(t, err)
	assert.Empty(t, repo.likes)
	assert.EqualValues(t, 0, repo.bookLists["booklist-1"].LikeCount)
}

func TestForkBookListRollbackOnForkCountFailure(t *testing.T) {
	repo := newBookListRepoState()
	repo.bookLists["booklist-1"] = &socialModel.BookList{
		ID:        primitive.NewObjectID(),
		UserID:    "owner-1",
		Title:     "list",
		ForkCount: 0,
	}
	repo.failIncrementFork = errors.New("mock fork count failure")

	service := NewBookListService(repo, nil)

	forked, err := service.ForkBookList(context.Background(), "user-2", "booklist-1")
	assert.Error(t, err)
	assert.Nil(t, forked)
	assert.Len(t, repo.bookLists, 1)
	assert.EqualValues(t, 0, repo.bookLists["booklist-1"].ForkCount)
}

func cloneBookListMap(source map[string]*socialModel.BookList) map[string]*socialModel.BookList {
	cloned := make(map[string]*socialModel.BookList, len(source))
	for key, value := range source {
		if value == nil {
			cloned[key] = nil
			continue
		}
		copyValue := *value
		if value.OriginalID != nil {
			originalID := *value.OriginalID
			copyValue.OriginalID = &originalID
		}
		copyValue.Books = append([]socialModel.BookListItem(nil), value.Books...)
		copyValue.Tags = append([]string(nil), value.Tags...)
		cloned[key] = &copyValue
	}
	return cloned
}

func cloneBookListLikeMap(source map[string]*socialModel.BookListLike) map[string]*socialModel.BookListLike {
	cloned := make(map[string]*socialModel.BookListLike, len(source))
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
