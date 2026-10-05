package shared

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWalletAPI_GetBalance_InvalidUserIDReturnsUnauthorized(t *testing.T) {
	api := NewWalletAPI(nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Set("user_id", 123)

	assert.NotPanics(t, func() {
		api.GetBalance(c)
	})
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	resp := decodeSharedAPIResponse(t, recorder)
	assert.Equal(t, "用户信息格式错误", resp["message"])
}
