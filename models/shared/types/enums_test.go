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

func TestParseRecommendationBehaviorType_NormalizesLegacyAliases(t *testing.T) {
	tests := map[string]RecommendationBehaviorType{
		"favorite": RecommendationBehaviorCollect,
		"collect":  RecommendationBehaviorCollect,
		"complete": RecommendationBehaviorFinish,
		"finish":   RecommendationBehaviorFinish,
	}

	for input, expected := range tests {
		actual, err := ParseRecommendationBehaviorType(input)
		if err != nil {
			t.Fatalf("expected %s to parse: %v", input, err)
		}
		if actual != expected {
			t.Fatalf("expected %s to normalize to %s, got %s", input, expected, actual)
		}
	}
}

func TestRecommendationBehaviorQueryValues_IncludeLegacyAliases(t *testing.T) {
	values := RecommendationBehaviorQueryValues(RecommendationBehaviorCollect)
	if len(values) != 2 || values[0] != RecommendationBehaviorCollect.String() || values[1] != "favorite" {
		t.Fatalf("expected collect query values to include favorite alias, got %#v", values)
	}
}

func TestParseReaderBehaviorType_NormalizesFinishToComplete(t *testing.T) {
	actual, err := ParseReaderBehaviorType("finish")
	if err != nil {
		t.Fatalf("expected finish to parse for reader behavior: %v", err)
	}
	if actual != ReaderBehaviorComplete {
		t.Fatalf("expected finish to normalize to complete, got %s", actual)
	}
}
