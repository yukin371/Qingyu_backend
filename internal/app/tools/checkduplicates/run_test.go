package checkduplicates

import "testing"

func TestParseOptionsDefaultsToHistoricalProject(t *testing.T) {
	opts, err := parseOptions(nil)
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}
	if opts.projectID != defaultProjectID {
		t.Fatalf("projectID = %q, want %q", opts.projectID, defaultProjectID)
	}
}

func TestParseOptionsAcceptsProjectID(t *testing.T) {
	opts, err := parseOptions([]string{"--project-id", "507f1f77bcf86cd799439011"})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}
	if opts.projectID != "507f1f77bcf86cd799439011" {
		t.Fatalf("projectID = %q", opts.projectID)
	}
}
