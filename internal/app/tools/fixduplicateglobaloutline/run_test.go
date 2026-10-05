package fixduplicateglobaloutline

import "testing"

func TestParseOptionsDefaultsToPreview(t *testing.T) {
	opts, err := parseOptions(nil)
	if err != nil || opts.apply {
		t.Fatalf("parseOptions() = %#v, %v", opts, err)
	}
}
