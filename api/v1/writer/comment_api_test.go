package writer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	writermodels "Qingyu_backend/models/writer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateCommentRequest_ToComment_RequiresParagraphID(t *testing.T) {
	docID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	req := &CreateCommentRequest{
		Content: "test",
		Type:    "comment",
	}

	comment, err := req.ToComment(docID, userID, "tester")
	assert.Error(t, err)
	assert.Nil(t, comment)
}

func TestCreateCommentRequest_ToComment_RejectsDeprecatedPosition(t *testing.T) {
	docID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	paragraphID := primitive.NewObjectID().Hex()

	req := &CreateCommentRequest{
		Content:     "test",
		Type:        "comment",
		ParagraphID: paragraphID,
		Position: &writermodels.CommentPosition{
			Paragraph: 1,
			Offset:    0,
			Length:    2,
		},
	}

	comment, err := req.ToComment(docID, userID, "tester")
	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.Contains(t, err.Error(), "position")
}

func TestCreateCommentRequest_ToComment_SucceedsWithParagraphID(t *testing.T) {
	docID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	paragraphID := primitive.NewObjectID().Hex()

	req := &CreateCommentRequest{
		Content:     "test",
		Type:        "comment",
		ParagraphID: paragraphID,
	}

	comment, err := req.ToComment(docID, userID, "tester")
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "test", comment.Content)
	expectedParagraphID, convErr := primitive.ObjectIDFromHex(paragraphID)
	assert.NoError(t, convErr)
	assert.Equal(t, expectedParagraphID, comment.ParagraphID)
}

func TestParseResolvedQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		query          string
		expected       *bool
		expectedOK     bool
		expectedStatus int
	}{
		{
			name:           "missing query",
			query:          "",
			expected:       nil,
			expectedOK:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid true",
			query:          "?resolved=true",
			expected:       boolPtr(true),
			expectedOK:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid false",
			query:          "?resolved=false",
			expected:       boolPtr(false),
			expectedOK:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid boolean",
			query:          "?resolved=maybe",
			expected:       nil,
			expectedOK:     false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			req := httptest.NewRequest(http.MethodGet, "/comments"+tt.query, nil)
			c.Request = req

			resolved, ok := parseResolvedQuery(c)

			assert.Equal(t, tt.expectedOK, ok)
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.expected == nil {
				assert.Nil(t, resolved)
				return
			}

			if assert.NotNil(t, resolved) {
				assert.Equal(t, *tt.expected, *resolved)
			}
		})
	}
}

func boolPtr(value bool) *bool {
	return &value
}
