package writer

import (
	"errors"
	"testing"

	lockpkg "Qingyu_backend/pkg/lock"
)

func TestIsLockedError_RecognizesSentinel(t *testing.T) {
	err := errors.Join(lockpkg.ErrDocumentLocked, errors.New("document is locked by alice"))

	if !isLockedError(err) {
		t.Fatal("expected sentinel document locked error to be recognized")
	}
}

func TestIsPermissionError_RecognizesSentinel(t *testing.T) {
	err := errors.Join(lockpkg.ErrLockPermissionDenied, errors.New("lock owned by alice"))

	if !isPermissionError(err) {
		t.Fatal("expected sentinel permission error to be recognized")
	}
}

func TestIsLockedError_KeepsLegacyStringFallback(t *testing.T) {
	if !isLockedError(errors.New("document is locked by alice")) {
		t.Fatal("expected legacy lock message to remain recognized")
	}
}

func TestIsPermissionError_KeepsLegacyStringFallback(t *testing.T) {
	if !isPermissionError(errors.New("permission denied: lock owned by alice")) {
		t.Fatal("expected legacy permission message to remain recognized")
	}
}
