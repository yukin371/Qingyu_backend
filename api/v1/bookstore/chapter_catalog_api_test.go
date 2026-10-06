package bookstore

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRequiredUserObjectIDRejectsMissingContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	_, ok := requiredUserObjectID(c)

	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, c.Writer.Status())
}

func TestRequiredUserObjectIDRejectsMalformedContextType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", 123)

	_, ok := requiredUserObjectID(c)

	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, c.Writer.Status())
}

func TestRequiredUserObjectIDRejectsMalformedID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "not-an-object-id")

	_, ok := requiredUserObjectID(c)

	require.False(t, ok)
	require.Equal(t, http.StatusBadRequest, c.Writer.Status())
}

func TestRequiredUserObjectIDReturnsParsedID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expected := primitive.NewObjectID()
	c.Set("user_id", expected.Hex())

	actual, ok := requiredUserObjectID(c)

	require.True(t, ok)
	require.Equal(t, expected, actual)
}

func TestOptionalUserObjectIDTreatsMissingAndMalformedValuesAsAnonymous(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for name, value := range map[string]interface{}{
		"missing":    nil,
		"wrong type": 123,
		"malformed":  "not-an-object-id",
	} {
		t.Run(name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			if value != nil {
				c.Set("user_id", value)
			}

			require.Equal(t, primitive.NilObjectID, optionalUserObjectID(c))
		})
	}
}

func TestOptionalUserObjectIDReturnsParsedID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expected := primitive.NewObjectID()
	c.Set("user_id", expected.Hex())

	require.Equal(t, expected, optionalUserObjectID(c))
}
