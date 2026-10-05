package cleanduplicates

import "testing"

func TestParseOptionsDefaultsToPreview(t *testing.T) {
	opts, err := parseOptions(nil)
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}
	if opts.apply || opts.projectID != defaultProjectID {
		t.Fatalf("unexpected defaults: %#v", opts)
	}
}

func TestParseOptionsAcceptsProjectAndApply(t *testing.T) {
	opts, err := parseOptions([]string{"--project-id", "507f1f77bcf86cd799439011", "--apply"})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}
	if !opts.apply || opts.projectID != "507f1f77bcf86cd799439011" {
		t.Fatalf("unexpected options: %#v", opts)
	}
}
