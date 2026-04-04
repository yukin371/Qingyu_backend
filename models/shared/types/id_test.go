package types

import (
	"errors"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseObjectID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()

	tests := []struct {
		name       string
		input      string
		wantErr    error
		wantNilOID bool
	}{
		{
			name:    "valid ID",
			input:   validID,
			wantErr: nil,
		},
		{
			name:       "empty string",
			input:      "",
			wantErr:    ErrEmptyID,
			wantNilOID: true,
		},
		{
			name:       "invalid format - too short",
			input:      "abc",
			wantErr:    ErrInvalidIDFormat,
			wantNilOID: true,
		},
		{
			name:       "invalid format - invalid chars",
			input:      "507f1f77bcf86cd79943901x",
			wantErr:    ErrInvalidIDFormat,
			wantNilOID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseObjectID(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ParseObjectID() error = %v, want %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("ParseObjectID() unexpected error = %v", err)
			}

			if tt.wantNilOID && got != primitive.NilObjectID {
				t.Errorf("ParseObjectID() = %v, want NilObjectID", got)
			}

			if tt.wantErr == nil && got.Hex() != tt.input {
				t.Errorf("ParseObjectID() = %v, want %v", got.Hex(), tt.input)
			}
		})
	}
}

func TestMustParseObjectID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()

	t.Run("valid ID returns ObjectID", func(t *testing.T) {
		got := MustParseObjectID(validID)
		if got.Hex() != validID {
			t.Errorf("MustParseObjectID() = %v, want %v", got.Hex(), validID)
		}
	})

	t.Run("invalid ID panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustParseObjectID() expected panic, got none")
			}
		}()
		MustParseObjectID("invalid")
	})
}

func TestParseOptionalObjectID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()

	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantErr bool
	}{
		{
			name:    "valid ID",
			input:   validID,
			wantNil: false,
			wantErr: false,
		},
		{
			name:    "empty string returns nil",
			input:   "",
			wantNil: true,
			wantErr: false,
		},
		{
			name:    "invalid format returns error",
			input:   "invalid",
			wantNil: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOptionalObjectID(tt.input)

			if tt.wantErr && err == nil {
				t.Error("ParseOptionalObjectID() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ParseOptionalObjectID() unexpected error = %v", err)
			}
			if tt.wantNil && got != nil {
				t.Errorf("ParseOptionalObjectID() = %v, want nil", got)
			}
			if !tt.wantNil && tt.input != "" && got != nil && got.Hex() != tt.input {
				t.Errorf("ParseOptionalObjectID() = %v, want %v", got.Hex(), tt.input)
			}
		})
	}
}

func TestToHex(t *testing.T) {
	oid := primitive.NewObjectID()
	got := ToHex(oid)
	if got != oid.Hex() {
		t.Errorf("ToHex() = %v, want %v", got, oid.Hex())
	}
}

func TestIsValidObjectID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()

	tests := []struct {
		input string
		want  bool
	}{
		{validID, true},
		{"", false},
		{"invalid", false},
		{"507f1f77bcf86cd79943901x", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidObjectID(tt.input); got != tt.want {
				t.Errorf("IsValidObjectID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseObjectIDSlice(t *testing.T) {
	validID1 := primitive.NewObjectID().Hex()
	validID2 := primitive.NewObjectID().Hex()

	tests := []struct {
		name       string
		input      []string
		wantLen    int
		wantErrLen int
	}{
		{
			name:       "empty slice",
			input:      []string{},
			wantLen:    0,
			wantErrLen: 0,
		},
		{
			name:       "all valid IDs",
			input:      []string{validID1, validID2},
			wantLen:    2,
			wantErrLen: 0,
		},
		{
			name:       "mixed valid and invalid",
			input:      []string{validID1, "invalid", validID2},
			wantLen:    2,
			wantErrLen: 1,
		},
		{
			name:       "all invalid",
			input:      []string{"invalid1", "invalid2"},
			wantLen:    0,
			wantErrLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, errs := ParseObjectIDSlice(tt.input)
			if len(got) != tt.wantLen {
				t.Errorf("ParseObjectIDSlice() len = %d, want %d", len(got), tt.wantLen)
			}
			if len(errs) != tt.wantErrLen {
				t.Errorf("ParseObjectIDSlice() errs len = %d, want %d", len(errs), tt.wantErrLen)
			}
		})
	}
}

func TestParseOptionalObjectIDSlice(t *testing.T) {
	validID1 := primitive.NewObjectID().Hex()
	validID2 := primitive.NewObjectID().Hex()

	tests := []struct {
		name    string
		input   []string
		wantLen int
		wantErr bool
	}{
		{
			name:    "empty slice returns nil",
			input:   []string{},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "nil slice returns nil",
			input:   nil,
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "valid IDs",
			input:   []string{validID1, validID2},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "skips empty strings",
			input:   []string{validID1, "", validID2},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "all empty returns empty",
			input:   []string{"", "", ""},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "invalid ID returns error",
			input:   []string{validID1, "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOptionalObjectIDSlice(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseOptionalObjectIDSlice() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ParseOptionalObjectIDSlice() unexpected error = %v", err)
				}
				if tt.wantLen == 0 && got != nil && len(got) != 0 {
					t.Errorf("ParseOptionalObjectIDSlice() = %v, want empty", got)
				}
				if tt.wantLen > 0 && len(got) != tt.wantLen {
					t.Errorf("ParseOptionalObjectIDSlice() len = %d, want %d", len(got), tt.wantLen)
				}
			}
		})
	}
}

func TestToHexSlice(t *testing.T) {
	oid1 := primitive.NewObjectID()
	oid2 := primitive.NewObjectID()

	t.Run("nil slice returns nil", func(t *testing.T) {
		got := ToHexSlice(nil)
		if got != nil {
			t.Errorf("ToHexSlice(nil) = %v, want nil", got)
		}
	})

	t.Run("valid slice", func(t *testing.T) {
		input := []primitive.ObjectID{oid1, oid2}
		got := ToHexSlice(input)
		if len(got) != 2 {
			t.Errorf("ToHexSlice() len = %d, want 2", len(got))
		}
		if got[0] != oid1.Hex() {
			t.Errorf("ToHexSlice()[0] = %v, want %v", got[0], oid1.Hex())
		}
	})
}

func TestGenerateNewObjectID(t *testing.T) {
	id1 := GenerateNewObjectID()
	id2 := GenerateNewObjectID()

	if id1 == "" || id2 == "" {
		t.Error("GenerateNewObjectID() returned empty string")
	}
	if id1 == id2 {
		t.Error("GenerateNewObjectID() returned duplicate IDs")
	}
	if len(id1) != 24 {
		t.Errorf("GenerateNewObjectID() len = %d, want 24", len(id1))
	}
}

func TestIsNilObjectID(t *testing.T) {
	tests := []struct {
		name  string
		input primitive.ObjectID
		want  bool
	}{
		{"nil ObjectID", primitive.NilObjectID, true},
		{"zero ObjectID", primitive.ObjectID{}, true},
		{"valid ObjectID", primitive.NewObjectID(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNilObjectID(tt.input); got != tt.want {
				t.Errorf("IsNilObjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIDError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrEmptyID", ErrEmptyID, true},
		{"ErrInvalidIDFormat", ErrInvalidIDFormat, true},
		{"wrapped ErrEmptyID", errors.Join(ErrEmptyID, errors.New("context")), true},
		{"wrapped ErrInvalidIDFormat", fmt.Errorf("wrapped: %w", ErrInvalidIDFormat), true},
		{"other error", errors.New("some other error"), false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIDError(tt.err); got != tt.want {
				t.Errorf("IsIDError() = %v, want %v", got, tt.want)
			}
		})
	}
}
