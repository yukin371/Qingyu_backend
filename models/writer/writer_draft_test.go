package writer

import "testing"

func TestWriterDraft_UpdateContent_DoesNotIncrementVersion(t *testing.T) {
	draft := &WriterDraft{
		Version: 1,
	}

	draft.UpdateContent("hello world")
	if draft.Version != 1 {
		t.Fatalf("expected version to remain 1, got %d", draft.Version)
	}
}
