package syncexistingvolumes

import "testing"

func TestParseOptionsDefaultsToPreview(t *testing.T) {
	opts, err := parseOptions(nil)
	if err != nil || opts.apply || opts.projectID != defaultProjectID {
		t.Fatalf("parseOptions() = %#v, %v", opts, err)
	}
}

func TestNumericValue(t *testing.T) {
	if value, err := numericValue(int32(3)); err != nil || value != 3 {
		t.Fatalf("numericValue() = %d, %v", value, err)
	}
}
