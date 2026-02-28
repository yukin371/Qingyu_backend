package auth

import "testing"

func TestSanitizeRoleName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid-basic", input: "reader", want: "reader"},
		{name: "valid-trim", input: " admin ", want: "admin"},
		{name: "valid-with-dash", input: "super-admin", want: "super-admin"},
		{name: "valid-with-underscore", input: "content_editor", want: "content_editor"},
		{name: "invalid-empty", input: "", wantErr: true},
		{name: "invalid-space-only", input: "   ", wantErr: true},
		{name: "invalid-dollar", input: "$ne", wantErr: true},
		{name: "invalid-dot", input: "a.b", wantErr: true},
		{name: "invalid-unicode", input: "管理员", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sanitizeRoleName(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}
