package cleanandseed

import "testing"

func TestParseApplyDefaultsToPreview(t *testing.T) {
	apply, err := parseApply(nil)
	if err != nil || apply {
		t.Fatalf("parseApply() = %v, %v", apply, err)
	}
}
