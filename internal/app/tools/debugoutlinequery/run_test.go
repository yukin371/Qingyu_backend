package debugoutlinequery

import "testing"

func TestParseOptionsDefault(t *testing.T) {
	opts, err := parseOptions(nil)
	if err != nil || opts.projectID != defaultProjectID {
		t.Fatalf("parseOptions() = %#v, %v", opts, err)
	}
}
