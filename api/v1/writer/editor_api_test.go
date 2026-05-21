package writer

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
