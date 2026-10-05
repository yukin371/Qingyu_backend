package cleartestdata

import "testing"

func TestRunDefaultsToReadOnly(t *testing.T) {
	t.Setenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA", "")
	if err := Run(nil); err != nil {
		t.Fatalf("Run(nil) error = %v", err)
	}
}

func TestRunRequiresExplicitDestructiveOptIn(t *testing.T) {
	t.Setenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA", "")
	if err := Run([]string{"--apply"}); err == nil {
		t.Fatal("expected destructive opt-in error")
	}
}
