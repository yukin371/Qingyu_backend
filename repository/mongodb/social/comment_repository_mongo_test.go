package social

import (
	"context"
	"testing"
)

func TestGetRepliesByCommentID_InvalidID(t *testing.T) {
	repo := &MongoCommentRepository{}

	_, err := repo.GetRepliesByCommentID(context.Background(), "invalid-id")
	if err == nil {
		t.Fatal("expected error for invalid comment id")
	}
}

