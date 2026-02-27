package reader

import "testing"

func TestValidateMongoQueryValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid value", value: "user_123", wantErr: false},
		{name: "empty value", value: "   ", wantErr: true},
		{name: "dollar sign", value: "abc$def", wantErr: true},
		{name: "null byte", value: "abc\x00def", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMongoQueryValue("field", tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateMongoQueryValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

