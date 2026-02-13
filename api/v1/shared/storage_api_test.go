package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	storageModel "Qingyu_backend/models/storage"
	storageRepo "Qingyu_backend/repository/interfaces/storage"
	"Qingyu_backend/service/shared/storage"
	storageMock "Qingyu_backend/service/shared/storage/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type apiStorageRepoAdapter struct {
	repo storageRepo.StorageRepository
}

func (a *apiStorageRepoAdapter) Create(ctx context.Context, file *storageModel.FileInfo) error {
	return a.repo.CreateFile(ctx, file)
}

func (a *apiStorageRepoAdapter) Get(ctx context.Context, fileID string) (*storageModel.FileInfo, error) {
	return a.repo.GetFile(ctx, fileID)
}

func (a *apiStorageRepoAdapter) Update(ctx context.Context, fileID string, updates map[string]interface{}) error {
	return a.repo.UpdateFile(ctx, fileID, updates)
}

func (a *apiStorageRepoAdapter) Delete(ctx context.Context, fileID string) error {
	return a.repo.DeleteFile(ctx, fileID)
}

func (a *apiStorageRepoAdapter) List(ctx context.Context, userID, category string, page, pageSize int) ([]*storageModel.FileInfo, error) {
	filter := &storageRepo.FileFilter{
		UserID:   userID,
		Category: category,
		Page:     page,
		PageSize: pageSize,
	}
	files, _, err := a.repo.ListFiles(ctx, filter)
	return files, err
}

func (a *apiStorageRepoAdapter) GrantAccess(ctx context.Context, fileID, userID string) error {
	return a.repo.GrantAccess(ctx, fileID, userID)
}

func (a *apiStorageRepoAdapter) RevokeAccess(ctx context.Context, fileID, userID string) error {
	return a.repo.RevokeAccess(ctx, fileID, userID)
}

func (a *apiStorageRepoAdapter) CheckAccess(ctx context.Context, fileID, userID string) (bool, error) {
	return a.repo.CheckAccess(ctx, fileID, userID)
}

func setupStorageMultipartRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	backend := storageMock.NewMockStorageBackend()
	repo := storageMock.NewMockStorageRepository()
	storageSvc := storage.NewStorageService(backend, &apiStorageRepoAdapter{repo: repo})
	multipartSvc := storage.NewMultipartUploadService(backend, repo)
	api := NewStorageAPI(storageSvc, multipartSvc, nil)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", "user123")
		c.Next()
	})

	files := router.Group("/api/v1/files")
	files.POST("/multipart/init", api.InitiateMultipartUpload)
	files.POST("/multipart/upload", api.UploadChunk)
	files.POST("/multipart/complete", api.CompleteMultipartUpload)
	files.POST("/multipart/abort", api.AbortMultipartUpload)
	files.GET("/multipart/progress", api.GetUploadProgress)

	return router
}

func decodeAPIResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

func TestStorageAPI_MultipartFlow_Success(t *testing.T) {
	router := setupStorageMultipartRouter()

	// 1) init (不传 uploaded_by，验证由服务端注入)
	initBody := `{"file_name":"demo.txt","file_size":11,"category":"document"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/multipart/init", strings.NewReader(initBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	initResp := decodeAPIResponse(t, w)
	data, ok := initResp["data"].(map[string]interface{})
	require.True(t, ok)
	uploadID, _ := data["upload_id"].(string)
	require.NotEmpty(t, uploadID)

	// 2) upload chunk
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("upload_id", uploadID))
	require.NoError(t, writer.WriteField("chunk_index", "0"))
	part, err := writer.CreateFormFile("chunk", "chunk0.bin")
	require.NoError(t, err)
	_, err = part.Write([]byte("hello world"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req = httptest.NewRequest(http.MethodPost, "/api/v1/files/multipart/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// 3) progress
	req = httptest.NewRequest(http.MethodGet, "/api/v1/files/multipart/progress?upload_id="+uploadID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	progressResp := decodeAPIResponse(t, w)
	progressData, ok := progressResp["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, float64(100), progressData["progress"])

	// 4) complete
	completeBody := `{"upload_id":"` + uploadID + `"}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/files/multipart/complete", strings.NewReader(completeBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestStorageAPI_MultipartAbort_Success(t *testing.T) {
	router := setupStorageMultipartRouter()

	initBody := `{"file_name":"to-abort.txt","file_size":5}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/multipart/init", strings.NewReader(initBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	initResp := decodeAPIResponse(t, w)
	data := initResp["data"].(map[string]interface{})
	uploadID := data["upload_id"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/files/multipart/abort?upload_id="+uploadID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestStorageAPI_MultipartInit_ValidationError(t *testing.T) {
	router := setupStorageMultipartRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/multipart/init", strings.NewReader(`{"file_name":"","file_size":0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStorageAPI_UploadChunk_MissingUploadID(t *testing.T) {
	router := setupStorageMultipartRouter()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("chunk_index", "0"))
	part, err := writer.CreateFormFile("chunk", "chunk0.bin")
	require.NoError(t, err)
	_, err = part.Write([]byte("hello"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/multipart/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeAPIResponse(t, w)
	assert.Equal(t, "参数错误", resp["message"])
}
