package writer

import (
	"testing"

	writermodels "Qingyu_backend/models/writer"

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
