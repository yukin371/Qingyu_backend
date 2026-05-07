package project

import "testing"

func TestCalculateSnapshotWordCount(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "empty",
			content:  "",
			expected: 0,
		},
		{
			name:     "ascii",
			content:  "hello world",
			expected: 11,
		},
		{
			name:     "unicode",
			content:  "你好，世界",
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateSnapshotWordCount(tt.content); got != tt.expected {
				t.Fatalf("calculateSnapshotWordCount(%q) = %d, want %d", tt.content, got, tt.expected)
			}
		})
	}
}
