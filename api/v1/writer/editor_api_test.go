package writer

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	writerModel "Qingyu_backend/models/writer"
	writerBase "Qingyu_backend/models/writer/base"
	servicemock "Qingyu_backend/service/mock"
	documentSvc "Qingyu_backend/service/writer/document"
	writerRepo "Qingyu_backend/repository/mongodb/writer"
	"Qingyu_backend/test/testutil"
)

func TestEditorApiCalculateWordCount_UsesMarkdownFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewEditorApi(nil)
	c, w := newWriterTestContext(http.MethodPost, "/api/v1/writer/editor/word-count", `{
		"content":"**你好** [Go](https://go.dev) 123",
		"filterMarkdown":true
	}`, nil)

	api.CalculateWordCount(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"totalCount":6`)
	assert.Contains(t, w.Body.String(), `"chineseCount":2`)
	assert.Contains(t, w.Body.String(), `"englishCount":1`)
	assert.Contains(t, w.Body.String(), `"numberCount":3`)
}

func TestEditorApiCalculateWordCount_UsesPlainTextPathWithoutMarkdownFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewEditorApi(nil)
	c, w := newWriterTestContext(http.MethodPost, "/api/v1/writer/editor/word-count", `{
		"content":"你好 abc 123"
	}`, nil)

	api.CalculateWordCount(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"totalCount":6`)
	assert.Contains(t, w.Body.String(), `"chineseCount":2`)
	assert.Contains(t, w.Body.String(), `"englishCount":1`)
	assert.Contains(t, w.Body.String(), `"numberCount":3`)
}

func TestEditorApiGetDocumentContent_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, cleanup := testutil.SetupTestDB(t)
	t.Cleanup(cleanup)

	ctx := context.Background()
	projectRepo := writerRepo.NewMongoProjectRepository(db)
	documentRepo := writerRepo.NewMongoDocumentRepository(db)
	contentRepo := new(servicemock.MockDocumentContentRepository)

	userID := primitive.NewObjectID()
	projectID := primitive.NewObjectID()
	documentID := primitive.NewObjectID()
	savedAt := time.Date(2026, 5, 21, 9, 30, 0, 0, time.UTC)

	must := func(err error) {
		require.NoError(t, err)
	}

	must(projectRepo.Create(ctx, &writerModel.Project{
		IdentifiedEntity: writerBase.IdentifiedEntity{ID: projectID},
		OwnedEntity:  writerBase.OwnedEntity{AuthorID: userID},
		TitledEntity: writerBase.TitledEntity{Title: "测试项目"},
		WritingType:  "novel",
		CoverURL:     "https://example.com/cover.jpg",
		Status:       writerModel.StatusDraft,
		Visibility:   writerModel.VisibilityPrivate,
	}))

	must(documentRepo.Create(ctx, &writerModel.Document{
		IdentifiedEntity: writerBase.IdentifiedEntity{ID: documentID},
		ProjectID:        projectID,
		Title:            "第一章",
		StableRef:        primitive.NewObjectID().Hex(),
		OrderKey:         writerModel.DefaultOrderKey,
		Type:             writerModel.TypeChapter,
		Level:            0,
		WordCount:        18,
	}))

	contentRepo.On("GetByDocumentID", mock.Anything, documentID.Hex()).Return(&writerModel.DocumentContent{
		DocumentID:  documentID,
		Content:     "正文内容",
		ContentType: "markdown",
		Version:     3,
		WordCount:   5,
		LastSavedAt: savedAt,
	}, nil).Once()

	api := NewEditorApi(documentSvc.NewDocumentService(documentRepo, contentRepo, projectRepo, nil))
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/documents/"+documentID.Hex()+"/content", "", gin.Params{{Key: "id", Value: documentID.Hex()}})
	c.Set("user_id", userID.Hex())

	api.GetDocumentContent(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("X-API-Deprecated"))
	assert.Equal(t, "/api/v1/writer/documents/{id}/contents", w.Header().Get("X-API-Replacement"))
	assert.Contains(t, w.Body.String(), `"documentId":"`+documentID.Hex()+`"`)
	assert.Contains(t, w.Body.String(), `"content":"正文内容"`)
	assert.Contains(t, w.Body.String(), `"contentType":"markdown"`)
	assert.Contains(t, w.Body.String(), `"version":3`)
	assert.Contains(t, w.Body.String(), `"wordCount":18`)
}
