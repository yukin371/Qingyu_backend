package bookstore

import (
	"context"
	"errors"
	"testing"

	bookstoreModel "Qingyu_backend/models/bookstore"
	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"

	"github.com/stretchr/testify/require"
)

type trialChapterRepositoryStub struct {
	BookstoreRepo.ChapterRepository
	freeErr   error
	byBookErr error
}

func (r *trialChapterRepositoryStub) GetFreeChapters(context.Context, string, int, int) ([]*bookstoreModel.Chapter, error) {
	if r.freeErr != nil {
		return nil, r.freeErr
	}
	return []*bookstoreModel.Chapter{{IsFree: true}}, nil
}

func (r *trialChapterRepositoryStub) GetByBookID(context.Context, string, int, int) ([]*bookstoreModel.Chapter, error) {
	return nil, r.byBookErr
}

func TestChapterPurchaseService_GetTrialChaptersReturnsSupplementError(t *testing.T) {
	supplementErr := errors.New("chapter repository unavailable")
	service := &ChapterPurchaseServiceImpl{
		chapterRepo: &trialChapterRepositoryStub{byBookErr: supplementErr},
	}

	chapters, err := service.GetTrialChapters(context.Background(), "book-1", 2)

	require.Error(t, err)
	require.ErrorIs(t, err, supplementErr)
	require.Nil(t, chapters)
}
