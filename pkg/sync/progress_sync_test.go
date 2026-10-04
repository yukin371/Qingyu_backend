package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	readerModel "Qingyu_backend/models/reader"
)

type progressPortStub struct {
	getProgressErr error
	saveErr        error
	saveCalls      int
}

func (s *progressPortStub) GetReadingProgress(context.Context, string, string) (*readerModel.ReadingProgress, error) {
	if s.getProgressErr != nil {
		return nil, s.getProgressErr
	}
	return &readerModel.ReadingProgress{LastReadAt: time.Now().Add(time.Hour)}, nil
}

func (s *progressPortStub) SaveReadingProgress(context.Context, string, string, string, float64) error {
	s.saveCalls++
	return s.saveErr
}

func TestMergeOfflineProgressesReturnsSaveError(t *testing.T) {
	saveErr := errors.New("storage unavailable")
	port := &progressPortStub{
		getProgressErr: errors.New("progress not found"),
		saveErr:        saveErr,
	}
	service := NewProgressSyncService(port)

	err := service.MergeOfflineProgresses(context.Background(), "507f1f77bcf86cd799439011", []OfflineProgress{{
		UserID:    "507f1f77bcf86cd799439011",
		BookID:    "book-1",
		ChapterID: "chapter-1",
		Progress:  0.5,
		Timestamp: time.Now(),
	}})

	if !errors.Is(err, saveErr) {
		t.Fatalf("expected save error to be returned, got %v", err)
	}
	if port.saveCalls != 1 {
		t.Fatalf("expected one save attempt, got %d", port.saveCalls)
	}
}
