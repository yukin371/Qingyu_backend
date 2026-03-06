package bookstore

import "testing"

func TestNormalizeBookStatus(t *testing.T) {
	if got := NormalizeBookStatus(BookStatus("published")); got != BookStatusOngoing {
		t.Fatalf("expected published to normalize to ongoing, got %s", got)
	}
}

func TestPublicBookStatusQueryValues_ContainsLegacyCompatibilityValue(t *testing.T) {
	values := PublicBookStatusQueryValues()
	want := map[string]bool{"ongoing": false, "completed": false, "published": false}
	for _, value := range values {
		if _, ok := want[value]; ok {
			want[value] = true
		}
	}
	for value, found := range want {
		if !found {
			t.Fatalf("expected query values to include %s", value)
		}
	}
}
