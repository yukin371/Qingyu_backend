package writer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
