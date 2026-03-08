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

type collectionRepoState struct {
	collections         map[string]*socialModel.Collection
	folders             map[string]*socialModel.CollectionFolder
	failIncrementFolder error
	failDecrementFolder error
}

func newCollectionRepoState() *collectionRepoState {
	return &collectionRepoState{
		collections: make(map[string]*socialModel.Collection),
		folders:     make(map[string]*socialModel.CollectionFolder),
	}
}

func (m *collectionRepoState) Create(ctx context.Context, collection *socialModel.Collection) error {
	if collection.ID.IsZero() {
		collection.ID = primitive.NewObjectID()
	}
	cloned := *collection
	m.collections[collection.ID.Hex()] = &cloned
	return nil
}

func (m *collectionRepoState) GetByID(ctx context.Context, id string) (*socialModel.Collection, error) {
	collection, ok := m.collections[id]
	if !ok {
		return nil, errors.New("collection not found")
	}
	cloned := *collection
	return &cloned, nil
}

func (m *collectionRepoState) GetByUserAndBook(ctx context.Context, userID, bookID string) (*socialModel.Collection, error) {
	for _, collection := range m.collections {
		if collection.UserID == userID && collection.BookID == bookID {
			cloned := *collection
			return &cloned, nil
		}
	}
	return nil, nil
}

func (m *collectionRepoState) GetByShareID(ctx context.Context, shareID string) (*socialModel.Collection, error) {
	return nil, nil
}

func (m *collectionRepoState) GetCollectionsByUser(ctx context.Context, userID string, folderID string, page, size int) ([]*socialModel.Collection, int64, error) {
	return nil, 0, nil
}

func (m *collectionRepoState) GetCollectionsByTag(ctx context.Context, userID string, tag string, page, size int) ([]*socialModel.Collection, int64, error) {
	return nil, 0, nil
}

func (m *collectionRepoState) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	collection, ok := m.collections[id]
	if !ok {
		return errors.New("collection not found")
	}
	if folderID, ok := updates["folder_id"].(string); ok {
		collection.FolderID = folderID
	}
	if note, ok := updates["note"].(string); ok {
		collection.Note = note
	}
	collection.UpdatedAt = time.Now()
	return nil
}

func (m *collectionRepoState) Delete(ctx context.Context, id string) error {
	if _, ok := m.collections[id]; !ok {
		return errors.New("collection not found")
	}
	delete(m.collections, id)
	return nil
}

func (m *collectionRepoState) CreateFolder(ctx context.Context, folder *socialModel.CollectionFolder) error {
	return nil
}

func (m *collectionRepoState) GetFolderByID(ctx context.Context, id string) (*socialModel.CollectionFolder, error) {
	folder, ok := m.folders[id]
	if !ok {
		return nil, errors.New("folder not found")
	}
	cloned := *folder
	return &cloned, nil
}

func (m *collectionRepoState) GetFoldersByUser(ctx context.Context, userID string) ([]*socialModel.CollectionFolder, error) {
	return nil, nil
}

func (m *collectionRepoState) UpdateFolder(ctx context.Context, id string, updates map[string]interface{}) error {
	return nil
}

func (m *collectionRepoState) DeleteFolder(ctx context.Context, id string) error {
	return nil
}

func (m *collectionRepoState) IncrementFolderBookCount(ctx context.Context, folderID string) error {
	if m.failIncrementFolder != nil {
		return m.failIncrementFolder
	}
	if folder, ok := m.folders[folderID]; ok {
		folder.BookCount++
		folder.UpdatedAt = time.Now()
	}
	return nil
}

func (m *collectionRepoState) DecrementFolderBookCount(ctx context.Context, folderID string) error {
	if m.failDecrementFolder != nil {
		return m.failDecrementFolder
	}
	if folder, ok := m.folders[folderID]; ok && folder.BookCount > 0 {
		folder.BookCount--
		folder.UpdatedAt = time.Now()
	}
	return nil
}

func (m *collectionRepoState) GetPublicCollections(ctx context.Context, page, size int) ([]*socialModel.Collection, int64, error) {
	return nil, 0, nil
}

func (m *collectionRepoState) GetPublicFolders(ctx context.Context, page, size int) ([]*socialModel.CollectionFolder, int64, error) {
	return nil, 0, nil
}

func (m *collectionRepoState) CountUserCollections(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (m *collectionRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	collectionsSnapshot := cloneCollectionMap(m.collections)
	foldersSnapshot := cloneFolderMap(m.folders)
	if err := fn(ctx); err != nil {
		m.collections = collectionsSnapshot
		m.folders = foldersSnapshot
		return err
	}
	return nil
}

func (m *collectionRepoState) Health(ctx context.Context) error { return nil }

func TestAddToCollectionRollbackOnFolderCountFailure(t *testing.T) {
	repo := newCollectionRepoState()
	repo.folders["folder-1"] = &socialModel.CollectionFolder{
		ID:        primitive.NewObjectID(),
		UserID:    "user-1",
		Name:      "folder",
		BookCount: 0,
	}
	repo.failIncrementFolder = errors.New("mock increment failure")

	service := NewCollectionService(repo, nil)

	collection, err := service.AddToCollection(context.Background(), "user-1", "book-1", "folder-1", "", nil, false)
	assert.Error(t, err)
	assert.Nil(t, collection)
	assert.Empty(t, repo.collections)
	assert.EqualValues(t, 0, repo.folders["folder-1"].BookCount)
}

func TestRemoveFromCollectionRollbackOnFolderCountFailure(t *testing.T) {
	repo := newCollectionRepoState()
	repo.collections["collection-1"] = &socialModel.Collection{
		ID:        primitive.NewObjectID(),
		UserID:    "user-1",
		BookID:    "book-1",
		FolderID:  "folder-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.folders["folder-1"] = &socialModel.CollectionFolder{
		ID:        primitive.NewObjectID(),
		UserID:    "user-1",
		Name:      "folder",
		BookCount: 1,
	}
	repo.failDecrementFolder = errors.New("mock decrement failure")

	service := NewCollectionService(repo, nil)

	err := service.RemoveFromCollection(context.Background(), "user-1", "collection-1")
	assert.Error(t, err)
	assert.Contains(t, repo.collections, "collection-1")
	assert.EqualValues(t, 1, repo.folders["folder-1"].BookCount)
}

func TestUpdateCollectionRollbackOnTargetFolderCountFailure(t *testing.T) {
	repo := newCollectionRepoState()
	repo.collections["collection-1"] = &socialModel.Collection{
		ID:        primitive.NewObjectID(),
		UserID:    "user-1",
		BookID:    "book-1",
		FolderID:  "folder-old",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.folders["folder-old"] = &socialModel.CollectionFolder{
		ID:        primitive.NewObjectID(),
		UserID:    "user-1",
		Name:      "old",
		BookCount: 1,
	}
	repo.folders["folder-new"] = &socialModel.CollectionFolder{
		ID:        primitive.NewObjectID(),
		UserID:    "user-1",
		Name:      "new",
		BookCount: 0,
	}
	repo.failIncrementFolder = errors.New("mock increment failure")

	service := NewCollectionService(repo, nil)

	err := service.UpdateCollection(context.Background(), "user-1", "collection-1", map[string]interface{}{"folder_id": "folder-new"})
	assert.Error(t, err)
	assert.Equal(t, "folder-old", repo.collections["collection-1"].FolderID)
	assert.EqualValues(t, 1, repo.folders["folder-old"].BookCount)
	assert.EqualValues(t, 0, repo.folders["folder-new"].BookCount)
}

func cloneCollectionMap(source map[string]*socialModel.Collection) map[string]*socialModel.Collection {
	cloned := make(map[string]*socialModel.Collection, len(source))
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

func cloneFolderMap(source map[string]*socialModel.CollectionFolder) map[string]*socialModel.CollectionFolder {
	cloned := make(map[string]*socialModel.CollectionFolder, len(source))
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
