package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	storageModel "Qingyu_backend/models/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStorageBackend struct {
	savedData   map[string][]byte
	lastSaveKey string
	getURLPath  string
	getURLTTL   time.Duration
	deleteErr   error
}

func newTestStorageBackend() *testStorageBackend {
	return &testStorageBackend{savedData: map[string][]byte{}}
}

func (b *testStorageBackend) Save(ctx context.Context, path string, reader io.Reader) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	b.savedData[path] = data
	b.lastSaveKey = path
	return nil
}

func (b *testStorageBackend) Load(ctx context.Context, path string) (io.ReadCloser, error) {
	data, ok := b.savedData[path]
	if !ok {
		return nil, errors.New("not found")
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (b *testStorageBackend) Delete(ctx context.Context, path string) error {
	delete(b.savedData, path)
	return b.deleteErr
}

func (b *testStorageBackend) Exists(ctx context.Context, path string) (bool, error) {
	_, ok := b.savedData[path]
	return ok, nil
}

func (b *testStorageBackend) GetURL(ctx context.Context, path string, expiresIn time.Duration) (string, error) {
	b.getURLPath = path
	b.getURLTTL = expiresIn
	return "http://example.com/" + path, nil
}

type testFileRepo struct {
	files     map[string]*storageModel.FileInfo
	lastFile  *storageModel.FileInfo
	deletedID string
}

func newTestFileRepo() *testFileRepo {
	return &testFileRepo{files: map[string]*storageModel.FileInfo{}}
}

func (r *testFileRepo) Create(ctx context.Context, file *FileInfo) error {
	copied := *file
	r.files[file.ID] = &copied
	r.lastFile = &copied
	return nil
}

func (r *testFileRepo) Get(ctx context.Context, fileID string) (*FileInfo, error) {
	f, ok := r.files[fileID]
	if !ok {
		return nil, errors.New("not found")
	}
	c := *f
	return &c, nil
}

func (r *testFileRepo) Update(ctx context.Context, fileID string, updates map[string]interface{}) error {
	return nil
}

func (r *testFileRepo) Delete(ctx context.Context, fileID string) error {
	delete(r.files, fileID)
	r.deletedID = fileID
	return nil
}

func (r *testFileRepo) List(ctx context.Context, userID, category string, page, pageSize int) ([]*FileInfo, error) {
	var out []*FileInfo
	for _, f := range r.files {
		c := *f
		out = append(out, &c)
	}
	return out, nil
}

func (r *testFileRepo) GrantAccess(ctx context.Context, fileID, userID string) error {
	return nil
}

func (r *testFileRepo) RevokeAccess(ctx context.Context, fileID, userID string) error {
	return nil
}

func (r *testFileRepo) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	return false, nil
}

func TestStorageService_Upload_DefaultCategoryAndSize(t *testing.T) {
	backend := newTestStorageBackend()
	repo := newTestFileRepo()
	svc := NewStorageService(backend, repo)

	data := []byte("hello-storage")
	req := &UploadRequest{
		File:     bytes.NewReader(data),
		Filename: "demo.txt",
		UserID:   "user-1",
	}

	fileInfo, err := svc.Upload(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, fileInfo)

	assert.Equal(t, defaultCategory, fileInfo.Category)
	assert.Equal(t, int64(len(data)), fileInfo.Size)
	assert.Equal(t, calculateBytesMD5(data), fileInfo.MD5)
	require.NotEmpty(t, backend.lastSaveKey)
	assert.Equal(t, data, backend.savedData[backend.lastSaveKey])
}

func TestStorageService_Upload_ContextCanceled(t *testing.T) {
	backend := newTestStorageBackend()
	repo := newTestFileRepo()
	svc := NewStorageService(backend, repo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fileInfo, err := svc.Upload(ctx, &UploadRequest{
		File:     bytes.NewReader([]byte("x")),
		Filename: "x.txt",
		UserID:   "u1",
	})

	assert.Nil(t, fileInfo)
	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestStorageService_Delete_BackendFailStillDeletesMetadata(t *testing.T) {
	backend := newTestStorageBackend()
	backend.deleteErr = errors.New("backend delete failed")
	repo := newTestFileRepo()
	svc := NewStorageService(backend, repo)

	repo.files["f1"] = &storageModel.FileInfo{
		ID:   "f1",
		Path: "general/2026/02/13/f1.txt",
	}

	err := svc.Delete(context.Background(), "f1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "删除文件失败")
	assert.Equal(t, "f1", repo.deletedID)
}

func TestStorageService_GetDownloadURL_DefaultTTL(t *testing.T) {
	backend := newTestStorageBackend()
	repo := newTestFileRepo()
	svc := NewStorageService(backend, repo)

	repo.files["f2"] = &storageModel.FileInfo{
		ID:   "f2",
		Path: "general/2026/02/13/f2.txt",
	}

	url, err := svc.GetDownloadURL(context.Background(), "f2", 0)
	require.NoError(t, err)
	assert.Equal(t, "http://example.com/general/2026/02/13/f2.txt", url)
	assert.Equal(t, defaultDownloadTTL, backend.getURLTTL)
}
