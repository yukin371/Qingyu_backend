package benchmarkapp

import "testing"

func TestParseFlagsUsesStableDefaults(t *testing.T) {
	config, err := parseFlags(nil)
	if err != nil {
		t.Fatalf("parseFlags() error = %v", err)
	}
	if config.BaseURL != "http://localhost:9090" || config.Requests != 1000 || config.Concurrent != 50 {
		t.Fatalf("unexpected defaults: %#v", config)
	}
}

func TestParseFlagsAcceptsOverrides(t *testing.T) {
	config, err := parseFlags([]string{"-base-url", "http://example.test", "-requests", "2", "-concurrent", "1", "-with-cache=false"})
	if err != nil {
		t.Fatalf("parseFlags() error = %v", err)
	}
	if config.BaseURL != "http://example.test" || config.Requests != 2 || config.Concurrent != 1 || config.WithCache {
		t.Fatalf("unexpected overrides: %#v", config)
	}
}
