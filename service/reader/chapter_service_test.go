package reader

import (
	bookstoreModel "Qingyu_backend/models/bookstore"
	readerModel "Qingyu_backend/models/reader"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupChapterServiceTest(t *testing.T) (ChapterService, *ReaderService, *MockReadingProgressRepository, *MockChapterService, *MockVIPPermissionService) {
	t.Helper()

	readerService, mockProgressRepo, _, _, mockChapterService, _, _ := setupReaderService()
	mockVIPService := new(MockVIPPermissionService)

	chapterService := NewChapterService(mockChapterService, readerService, mockVIPService)

	return chapterService, readerService, mockProgressRepo, mockChapterService, mockVIPService
}

func TestChapterService_GetChapterContent_ProtectedChapterDeniedWithoutLogin(t *testing.T) {
	chapterService, _, _, mockChapterService, _ := setupChapterServiceTest(t)

	ctx := context.Background()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	chapterObjectID, _ := primitive.ObjectIDFromHex(chapterID)

	chapter := &bookstoreModel.Chapter{
		ID:          chapterObjectID,
		BookID:      bookID,
		Title:       "付费章节",
		ChapterNum:  3,
		WordCount:   1234,
		IsFree:      false,
		PublishTime: time.Now().Add(-time.Hour),
	}

	mockChapterService.On("GetChapterByID", mock.Anything, chapterID).Return(chapter, nil).Once()

	result, err := chapterService.GetChapterContent(ctx, "", bookID, chapterID)

	require.ErrorIs(t, err, ErrAccessDenied)
	require.NotNil(t, result)
	assert.Equal(t, chapterID, result.ChapterID)
	assert.Equal(t, bookID, result.BookID)
	assert.False(t, result.CanAccess)
	assert.Empty(t, result.Content)
	assert.Equal(t, "需要登录后阅读付费章节", result.AccessReason)
	mockChapterService.AssertNotCalled(t, "GetChapterContent", mock.Anything, chapterID, mock.Anything)
	mockChapterService.AssertNotCalled(t, "GetChapterParagraphs", mock.Anything, chapterID, mock.Anything)
	mockChapterService.AssertExpectations(t)
}

func TestChapterService_GetChapterContent_AllowsVIPAccess(t *testing.T) {
	chapterService, readerService, mockProgressRepo, mockChapterService, mockVIPService := setupChapterServiceTest(t)

	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	chapterObjectID, _ := primitive.ObjectIDFromHex(chapterID)
	readerService.eventBus = nil

	chapter := &bookstoreModel.Chapter{
		ID:          chapterObjectID,
		BookID:      bookID,
		Title:       "VIP章节",
		ChapterNum:  4,
		WordCount:   4321,
		IsFree:      false,
		PublishTime: time.Now().Add(-time.Hour),
	}

	paragraphs := []*bookstoreModel.ChapterContent{
		{
			ID:             primitive.NewObjectID(),
			ChapterID:      chapterObjectID,
			Content:        "第一段",
			Format:         bookstoreModel.ContentFormatText,
			ParagraphOrder: 1,
			WordCount:      4,
		},
	}

	mockChapterService.On("GetChapterByID", mock.Anything, chapterID).Return(chapter, nil).Once()
	mockVIPService.On("CheckVIPAccess", mock.Anything, userID, chapterID, true).Return(true, nil).Once()
	mockChapterService.On("GetChapterContent", mock.Anything, chapterID, userID).Return("VIP正文", nil).Once()
	mockChapterService.On("GetChapterParagraphs", mock.Anything, chapterID, userID).Return(paragraphs, nil).Once()
	mockChapterService.On("GetNextChapter", mock.Anything, bookID, chapter.ChapterNum).Return((*bookstoreModel.Chapter)(nil), nil).Once()
	mockChapterService.On("GetPreviousChapter", mock.Anything, bookID, chapter.ChapterNum).Return((*bookstoreModel.Chapter)(nil), nil).Once()
	mockProgressRepo.On("GetByUserAndBook", mock.Anything, userID, bookID).Return((*readerModel.ReadingProgress)(nil), nil).Once()
	mockProgressRepo.On("SaveProgressWithInitial", mock.Anything, userID, bookID, chapterID, 0.0, int64(0)).Return(nil).Once()

	result, err := chapterService.GetChapterContent(ctx, userID, bookID, chapterID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.CanAccess)
	assert.Equal(t, "VIP正文", result.Content)
	assert.Len(t, result.Paragraphs, 1)
	assert.Equal(t, chapterID, result.ChapterID)
	assert.Equal(t, bookID, result.BookID)
	assert.Equal(t, 0.0, result.Progress)
	assert.Equal(t, int64(0), result.ReadingTime)
	mockChapterService.AssertExpectations(t)
	mockVIPService.AssertExpectations(t)
	mockProgressRepo.AssertExpectations(t)
}
