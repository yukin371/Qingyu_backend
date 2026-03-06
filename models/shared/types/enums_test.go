package types

import "testing"

func TestParseBookStatus_NormalizesLegacyPublished(t *testing.T) {
	status, err := ParseBookStatus("published")
	if err != nil {
		t.Fatalf("expected legacy published to parse, got error: %v", err)
	}
	if status != BookStatusOngoing {
		t.Fatalf("expected published to normalize to ongoing, got %s", status)
	}
}

func TestDocumentStatus_UsesWriterLifecycleValues(t *testing.T) {
	valid := []DocumentStatus{DocumentStatusPlanned, DocumentStatusWriting, DocumentStatusCompleted}
	for _, status := range valid {
		if !status.IsValid() {
			t.Fatalf("expected %s to be valid", status)
		}
	}

	if _, err := ParseDocumentStatus("published"); err == nil {
		t.Fatal("expected legacy published document status to be rejected")
	}
}
