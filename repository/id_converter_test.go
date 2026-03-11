package repository

import (
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()

	tests := []struct {
		name        string
		input       string
		wantErr     error
		wantNilOID  bool
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
			got, err := ParseID(tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ParseID() error = %v, want %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("ParseID() unexpected error = %v", err)
			}

			if tt.wantNilOID && got != primitive.NilObjectID {
				t.Errorf("ParseID() = %v, want NilObjectID", got)
			}

			if tt.wantErr == nil && got.Hex() != tt.input {
				t.Errorf("ParseID() = %v, want %v", got.Hex(), tt.input)
			}
		})
	}
}

func TestParseOptionalID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()

	tests := []struct {
		name      string
		input     string
		wantNil   bool
		wantErr   bool
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
			got, err := ParseOptionalID(tt.input)

			if tt.wantErr && err == nil {
				t.Error("ParseOptionalID() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ParseOptionalID() unexpected error = %v", err)
			}
			if tt.wantNil && got != nil {
				t.Errorf("ParseOptionalID() = %v, want nil", got)
			}
			if !tt.wantNil && tt.input != "" && got != nil && got.Hex() != tt.input {
				t.Errorf("ParseOptionalID() = %v, want %v", got.Hex(), tt.input)
			}
		})
	}
}

func TestParseIDs(t *testing.T) {
	validID1 := primitive.NewObjectID().Hex()
	validID2 := primitive.NewObjectID().Hex()

	tests := []struct {
		name       string
		input      []string
		wantLen    int
		wantErr    bool
		errContain string
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
			name:       "contains empty string",
			input:      []string{validID1, "", validID2},
			wantErr:    true,
			errContain: "ids[1]",
		},
		{
			name:       "contains invalid ID",
			input:      []string{validID1, "invalid"},
			wantErr:    true,
			errContain: "ids[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIDs(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseIDs() expected error, got nil")
				} else if tt.errContain != "" && !errors.Is(err, ErrEmptyID) && !errors.Is(err, ErrInvalidIDFormat) {
					t.Errorf("ParseIDs() error should contain ID error, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("ParseIDs() unexpected error = %v", err)
				}
				if tt.wantLen == 0 && got != nil {
					t.Errorf("ParseIDs() = %v, want nil", got)
				}
				if tt.wantLen > 0 && len(got) != tt.wantLen {
					t.Errorf("ParseIDs() len = %d, want %d", len(got), tt.wantLen)
				}
			}
		})
	}
}

func TestParseOptionalIDs(t *testing.T) {
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
			got, err := ParseOptionalIDs(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseOptionalIDs() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ParseOptionalIDs() unexpected error = %v", err)
				}
				if tt.wantLen == 0 && got != nil && len(got) != 0 {
					t.Errorf("ParseOptionalIDs() = %v, want empty", got)
				}
				if tt.wantLen > 0 && len(got) != tt.wantLen {
					t.Errorf("ParseOptionalIDs() len = %d, want %d", len(got), tt.wantLen)
				}
			}
		})
	}
}

func TestIsIDError(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		want  bool
	}{
		{
			name: "ErrEmptyID",
			err:  ErrEmptyID,
			want: true,
		},
		{
			name: "ErrInvalidIDFormat",
			err:  ErrInvalidIDFormat,
			want: true,
		},
		{
			name: "wrapped ErrEmptyID",
			err:  errors.Join(ErrEmptyID, errors.New("context")),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("some other error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIDError(tt.err); got != tt.want {
				t.Errorf("IsIDError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBackwardCompatibility 验证现有函数仍能正常工作
func TestBackwardCompatibility(t *testing.T) {
	validID := primitive.NewObjectID().Hex()

	t.Run("StringToObjectId - empty returns NilObjectID", func(t *testing.T) {
		got, err := StringToObjectId("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != primitive.NilObjectID {
			t.Errorf("expected NilObjectID, got %v", got)
		}
	})

	t.Run("StringToObjectId - valid ID", func(t *testing.T) {
		got, err := StringToObjectId(validID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got.Hex() != validID {
			t.Errorf("expected %s, got %s", validID, got.Hex())
		}
	})

	t.Run("StringSliceToObjectIDSlice - skips empty", func(t *testing.T) {
		got, err := StringSliceToObjectIDSlice([]string{validID, "", validID})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("expected 2 IDs, got %d", len(got))
		}
	})
}
