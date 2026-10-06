package bookstore

import (
	"context"
	"testing"

	bookstoreModel "Qingyu_backend/models/bookstore"
	BookstoreRepo "Qingyu_backend/repository/interfaces/bookstore"
	repoInfra "Qingyu_backend/repository/interfaces/infrastructure"
)

type searchChapterRepositoryStub struct {
	BookstoreRepo.ChapterRepository
	chapters []*bookstoreModel.Chapter
	total    int64
}

func (r *searchChapterRepositoryStub) Search(context.Context, string, int, int) ([]*bookstoreModel.Chapter, error) {
	return r.chapters, nil
}

func (r *searchChapterRepositoryStub) Count(context.Context, repoInfra.Filter) (int64, error) {
	return r.total, nil
}

func TestChapterServiceSearchChaptersUsesRepositoryCount(t *testing.T) {
	repo := &searchChapterRepositoryStub{
		chapters: []*bookstoreModel.Chapter{{Title: "第一章"}},
		total:    37,
	}
	service := NewChapterService(repo, nil, nil)

	chapters, total, err := service.SearchChapters(context.Background(), "第一", 2, 10)
	if err != nil {
		t.Fatalf("SearchChapters() error = %v", err)
	}
	if len(chapters) != 1 {
		t.Fatalf("SearchChapters() returned %d chapters, want 1", len(chapters))
	}
	if total != 37 {
		t.Fatalf("SearchChapters() total = %d, want repository count 37", total)
	}
}
