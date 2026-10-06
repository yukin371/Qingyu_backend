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

type catalogChapterRepositoryStub struct {
	BookstoreRepo.ChapterRepository
	countErr error
}

func (r *catalogChapterRepositoryStub) GetByBookID(context.Context, string, int, int) ([]*bookstoreModel.Chapter, error) {
	return []*bookstoreModel.Chapter{}, nil
}

func (r *catalogChapterRepositoryStub) CountByBookID(context.Context, string) (int64, error) {
	return 0, r.countErr
}

func (r *catalogChapterRepositoryStub) GetTotalWordCount(context.Context, string) (int64, error) {
	return 0, nil
}

type catalogBookRepositoryStub struct {
	BookstoreRepo.BookRepository
}

func (r *catalogBookRepositoryStub) GetByID(context.Context, string) (*bookstoreModel.Book, error) {
	return &bookstoreModel.Book{Title: "测试书籍"}, nil
}

type catalogPurchaseRepositoryStub struct {
	BookstoreRepo.ChapterPurchaseRepository
}

func (r *catalogPurchaseRepositoryStub) GetPurchasedChapterIDs(context.Context, string, string) ([]string, error) {
	return nil, nil
}

type batchPurchaseRepositoryStub struct {
	BookstoreRepo.ChapterPurchaseRepository
	checkErr error
}

func (r *batchPurchaseRepositoryStub) GetByUserAndChapter(context.Context, string, string) (*bookstoreModel.ChapterPurchase, error) {
	return nil, r.checkErr
}

type batchChapterRepositoryStub struct {
	BookstoreRepo.ChapterRepository
}

func (r *batchChapterRepositoryStub) GetByID(context.Context, string) (*bookstoreModel.Chapter, error) {
	return &bookstoreModel.Chapter{BookID: "book-1", Price: 10}, nil
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

func TestChapterPurchaseService_GetChapterCatalogReturnsCountError(t *testing.T) {
	countErr := errors.New("count unavailable")
	service := &ChapterPurchaseServiceImpl{
		chapterRepo:  &catalogChapterRepositoryStub{countErr: countErr},
		bookRepo:     &catalogBookRepositoryStub{},
		purchaseRepo: &catalogPurchaseRepositoryStub{},
	}

	catalog, err := service.GetChapterCatalog(context.Background(), "", "book-1")

	require.Error(t, err)
	require.ErrorIs(t, err, countErr)
	require.Nil(t, catalog)
}

func TestChapterPurchaseService_PurchaseChaptersReturnsExistingPurchaseCheckError(t *testing.T) {
	checkErr := errors.New("purchase check unavailable")
	service := &ChapterPurchaseServiceImpl{
		chapterRepo:  &batchChapterRepositoryStub{},
		purchaseRepo: &batchPurchaseRepositoryStub{checkErr: checkErr},
	}

	batch, err := service.PurchaseChapters(context.Background(), "user-1", []string{"chapter-1"})

	require.Error(t, err)
	require.ErrorIs(t, err, checkErr)
	require.Nil(t, batch)
}
