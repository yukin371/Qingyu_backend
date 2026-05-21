package writer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"Qingyu_backend/service/writer/document"
)

func newWriterTestContext(method, target string, body string, params gin.Params) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	c.Request = req
	c.Params = params
	return c, w
}

func TestCharacterApiGetCharacter_RequiresProjectIDQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewCharacterApi(nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/characters/char-1", "", gin.Params{
		{Key: "characterId", Value: "char-1"},
	})

	api.GetCharacter(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestLocationApiGetLocation_RequiresProjectIDQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewLocationApi(nil)
	c, w := newWriterTestContext(http.MethodGet, "/api/v1/writer/locations/location-1", "", gin.Params{
		{Key: "locationId", Value: "location-1"},
	})

	api.GetLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestExportApiExportDocument_RequiresProjectIDQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewExportApi(nil)
	c, w := newWriterTestContext(http.MethodPost, "/api/v1/writer/documents/doc-1/export", "", gin.Params{
		{Key: "id", Value: "doc-1"},
	})

	api.ExportDocument(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestChangeRequestApiProcessChangeRequest_InvalidJSONReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := NewChangeRequestApi(nil)
	c, w := newWriterTestContext(http.MethodPut, "/api/v1/writer/change-requests/cr-1/status", "{\"status\":", gin.Params{
		{Key: "requestId", Value: "cr-1"},
	})

	api.ProcessChangeRequest(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "参数错误")
}

func TestDocumentApiCreateDocument_RequiresLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{documentService: &document.DocumentService{}}
	c, w := newWriterTestContext(http.MethodPost, "/api/v1/projects/project-1/documents", `{
		"title":"新文档",
		"type":"chapter"
	}`, gin.Params{
		{Key: "projectId", Value: "project-1"},
	})

	api.CreateDocument(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "请先登录")
}

func TestDocumentApiCreateDocument_RequiresProjectIDPathParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{documentService: &document.DocumentService{}}
	c, w := newWriterTestContext(http.MethodPost, "/api/v1/projects//documents", `{
		"title":"新文档",
		"type":"chapter"
	}`, nil)
	c.Set("user_id", "user-1")

	api.CreateDocument(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "项目ID不能为空")
}

func TestDocumentApiCreateDocument_RejectsMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	api := &DocumentApi{documentService: &document.DocumentService{}}
	c, w := newWriterTestContext(http.MethodPost, "/api/v1/projects/project-1/documents", `{"title":`, gin.Params{
		{Key: "projectId", Value: "project-1"},
	})
	c.Set("user_id", "user-1")

	api.CreateDocument(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "参数错误")
}
