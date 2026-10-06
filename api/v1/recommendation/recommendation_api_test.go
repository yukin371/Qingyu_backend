package recommendation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthenticatedUserIDRejectsMissingContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, ok := authenticatedUserID(c)

	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, c.Writer.Status())
}

func TestAuthenticatedUserIDRejectsMalformedContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", 123)

	_, ok := authenticatedUserID(c)

	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, c.Writer.Status())
}

func TestAuthenticatedUserIDReturnsValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "user-a")

	userID, ok := authenticatedUserID(c)

	require.True(t, ok)
	require.Equal(t, "user-a", userID)
}
