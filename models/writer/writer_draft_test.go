package writer

import "testing"

func TestWriterDraft_UpdateContent_RefreshesVersionAndWordCount(t *testing.T) {
	draft := &WriterDraft{
		Version: 1,
	}

	draft.UpdateContent("hello world")

	if draft.Content != "hello world" {
		t.Fatalf("expected content to be updated, got %q", draft.Content)
	}
	if draft.WordCount != 11 {
		t.Fatalf("expected word count to be 11, got %d", draft.WordCount)
	}
	if draft.Version != 2 {
		t.Fatalf("expected version to increment to 2, got %d", draft.Version)
	}
}
