package checktypes

import "testing"

func TestRunFlagSetRejectsUnknownFlag(t *testing.T) {
	if _, err := parseArgs([]string{"--unknown"}); err == nil {
		t.Fatal("expected unknown flag error")
	}
}

func TestParseArgsDefaultsToReadOnly(t *testing.T) {
	apply, err := parseArgs(nil)
	if err != nil {
		t.Fatalf("parseArgs() error = %v", err)
	}
	if apply {
		t.Fatal("expected read-only default")
	}
}

func TestParseArgsAcceptsApply(t *testing.T) {
	apply, err := parseArgs([]string{"--apply"})
	if err != nil {
		t.Fatalf("parseArgs() error = %v", err)
	}
	if !apply {
		t.Fatal("expected apply=true")
	}
}
