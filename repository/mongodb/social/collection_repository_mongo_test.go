package social

import "testing"

func TestValidateCollectionTag(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid value", value: "标签A_1", wantErr: false},
		{name: "empty value", value: "   ", wantErr: true},
		{name: "dollar sign", value: "abc$def", wantErr: true},
		{name: "null byte", value: "abc\x00def", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCollectionTag(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateCollectionTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeObjectIDHex(t *testing.T) {
	validID := "507f1f77bcf86cd799439011"
	if _, err := normalizeObjectIDHex("user_id", validID); err != nil {
		t.Fatalf("expected valid object id, got err=%v", err)
	}

	if _, err := normalizeObjectIDHex("user_id", "not-object-id"); err == nil {
		t.Fatal("expected error for invalid object id")
	}
}
