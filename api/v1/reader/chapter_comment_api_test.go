package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
	readerModels "Qingyu_backend/models/reader"
	socialModels "Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	bookstoreService "Qingyu_backend/service/bookstore"
	socialService "Qingyu_backend/service/social"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type chapterCommentRepoStub struct {
	comments map[string]*socialModels.Comment
}

func newChapterCommentRepoStub(seed ...*socialModels.Comment) *chapterCommentRepoStub {
	repo := &chapterCommentRepoStub{comments: make(map[string]*socialModels.Comment)}
	for _, item := range seed {
		repo.comments[item.ID.Hex()] = item
	}
	return repo
}

func (m *chapterCommentRepoStub) Create(_ context.Context, comment *socialModels.Comment) error {
	if comment.ID.IsZero() {
		comment.ID = primitive.NewObjectID()
	}
	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}
	comment.UpdatedAt = time.Now()
	m.comments[comment.ID.Hex()] = comment
	return nil
}

func (m *chapterCommentRepoStub) GetByID(_ context.Context, id string) (*socialModels.Comment, error) {
	if comment, ok := m.comments[id]; ok {
		return comment, nil
	}
	return nil, errors.New("comment not found")
}

func (m *chapterCommentRepoStub) Update(_ context.Context, id string, updates map[string]interface{}) error {
	comment, ok := m.comments[id]
	if !ok {
		return errors.New("comment not found")
	}
	if content, ok := updates["content"].(string); ok {
		comment.Content = content
	}
	comment.UpdatedAt = time.Now()
	return nil
}

func (m *chapterCommentRepoStub) Delete(_ context.Context, id string) error {
	delete(m.comments, id)
	return nil
}

func (m *chapterCommentRepoStub) Exists(_ context.Context, id string) (bool, error) {
	_, ok := m.comments[id]
	return ok, nil
}

func (m *chapterCommentRepoStub) GetCommentsByBookID(_ context.Context, _ string, _ int, _ int) ([]*socialModels.Comment, int64, error) {
	return nil, 0, nil
}

func (m *chapterCommentRepoStub) GetCommentsByUserID(_ context.Context, _ string, _ int, _ int) ([]*socialModels.Comment, int64, error) {
	return nil, 0, nil
}

func (m *chapterCommentRepoStub) GetRepliesByCommentID(_ context.Context, commentID string) ([]*socialModels.Comment, error) {
	var replies []*socialModels.Comment
	for _, item := range m.comments {
		if item.ParentID != nil && *item.ParentID == commentID {
			replies = append(replies, item)
		}
	}
	return replies, nil
}

func (m *chapterCommentRepoStub) GetCommentsByChapterID(_ context.Context, _ string, _ int, _ int) ([]*socialModels.Comment, int64, error) {
	return nil, 0, nil
}

func (m *chapterCommentRepoStub) ListByFilter(_ context.Context, filter *socialModels.CommentFilter) ([]*socialModels.Comment, int64, error) {
	var result []*socialModels.Comment
	for _, item := range m.comments {
		if filter.TargetID != nil && item.TargetID != *filter.TargetID {
			continue
		}
		if filter.ChapterID != nil && item.ChapterID != *filter.ChapterID {
			continue
		}
		if filter.State != nil && item.State != *filter.State {
			continue
		}
		if filter.ParentID != nil {
			parentID := ""
			if item.ParentID != nil {
				parentID = *item.ParentID
			}
			if parentID != *filter.ParentID {
				continue
			}
		}
		if filter.ParagraphIndex != nil {
			meta := parseParagraphRichContent(item.RichContent)
			if meta.ParagraphIndex != *filter.ParagraphIndex {
				continue
			}
		}
		result = append(result, item)
	}
	return result, int64(len(result)), nil
}

func (m *chapterCommentRepoStub) GetCommentsByBookIDSorted(_ context.Context, _ string, _ string, _ int, _ int) ([]*socialModels.Comment, int64, error) {
	return nil, 0, nil
}

func (m *chapterCommentRepoStub) UpdateCommentStatus(_ context.Context, _ string, _ string, _ string) error {
	return nil
}

func (m *chapterCommentRepoStub) GetPendingComments(_ context.Context, _ int, _ int) ([]*socialModels.Comment, int64, error) {
	return nil, 0, nil
}

func (m *chapterCommentRepoStub) IncrementLikeCount(_ context.Context, id string) error {
	comment, err := m.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	comment.LikeCount++
	return nil
}

func (m *chapterCommentRepoStub) DecrementLikeCount(_ context.Context, id string) error {
	comment, err := m.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	if comment.LikeCount > 0 {
		comment.LikeCount--
	}
	return nil
}

func (m *chapterCommentRepoStub) IncrementReplyCount(_ context.Context, id string) error {
	comment, err := m.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	comment.ReplyCount++
	return nil
}

func (m *chapterCommentRepoStub) DecrementReplyCount(_ context.Context, id string) error {
	comment, err := m.GetByID(context.Background(), id)
	if err != nil {
		return err
	}
	if comment.ReplyCount > 0 {
		comment.ReplyCount--
	}
	return nil
}

func (m *chapterCommentRepoStub) GetBookRatingStats(_ context.Context, _ string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (m *chapterCommentRepoStub) GetCommentCount(_ context.Context, _ string) (int64, error) {
	return int64(len(m.comments)), nil
}

func (m *chapterCommentRepoStub) GetCommentsByIDs(_ context.Context, ids []string) ([]*socialModels.Comment, error) {
	result := make([]*socialModels.Comment, 0, len(ids))
	for _, id := range ids {
		if item, ok := m.comments[id]; ok {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *chapterCommentRepoStub) DeleteCommentsByBookID(_ context.Context, _ string) error {
	return nil
}

func (m *chapterCommentRepoStub) Health(_ context.Context) error {
	return nil
}

func (m *chapterCommentRepoStub) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

var _ socialRepo.CommentRepository = (*chapterCommentRepoStub)(nil)

type chapterServiceStub struct {
	chapter    *bookstoreModel.Chapter
	paragraphs []*bookstoreModel.ChapterContent
}

func (m *chapterServiceStub) CreateChapter(context.Context, *bookstoreModel.Chapter) error { return nil }
func (m *chapterServiceStub) GetChapterByID(context.Context, string) (*bookstoreModel.Chapter, error) {
	if m.chapter == nil {
		return nil, errors.New("chapter not found")
	}
	return m.chapter, nil
}
func (m *chapterServiceStub) UpdateChapter(context.Context, *bookstoreModel.Chapter) error { return nil }
func (m *chapterServiceStub) DeleteChapter(context.Context, string) error                  { return nil }
func (m *chapterServiceStub) GetChaptersByBookID(context.Context, string, int, int) ([]*bookstoreModel.Chapter, int64, error) {
	return nil, 0, nil
}
func (m *chapterServiceStub) GetChapterByBookIDAndNum(context.Context, string, int) (*bookstoreModel.Chapter, error) {
	return nil, nil
}
func (m *chapterServiceStub) GetChaptersByTitle(context.Context, string, int, int) ([]*bookstoreModel.Chapter, int64, error) {
	return nil, 0, nil
}
func (m *chapterServiceStub) GetFreeChaptersByBookID(context.Context, string, int, int) ([]*bookstoreModel.Chapter, int64, error) {
	return nil, 0, nil
}
func (m *chapterServiceStub) GetPaidChaptersByBookID(context.Context, string, int, int) ([]*bookstoreModel.Chapter, int64, error) {
	return nil, 0, nil
}
func (m *chapterServiceStub) GetPublishedChaptersByBookID(context.Context, string, int, int) ([]*bookstoreModel.Chapter, int64, error) {
	return nil, 0, nil
}
func (m *chapterServiceStub) GetPreviousChapter(context.Context, string, int) (*bookstoreModel.Chapter, error) {
	return nil, nil
}
func (m *chapterServiceStub) GetNextChapter(context.Context, string, int) (*bookstoreModel.Chapter, error) {
	return nil, nil
}
func (m *chapterServiceStub) GetFirstChapter(context.Context, string) (*bookstoreModel.Chapter, error) {
	return nil, nil
}
func (m *chapterServiceStub) GetLastChapter(context.Context, string) (*bookstoreModel.Chapter, error) {
	return nil, nil
}
func (m *chapterServiceStub) GetChapterCountByBookID(context.Context, string) (int64, error) {
	return 0, nil
}
func (m *chapterServiceStub) GetFreeChapterCountByBookID(context.Context, string) (int64, error) {
	return 0, nil
}
func (m *chapterServiceStub) GetPaidChapterCountByBookID(context.Context, string) (int64, error) {
	return 0, nil
}
func (m *chapterServiceStub) GetTotalWordCountByBookID(context.Context, string) (int64, error) {
	return 0, nil
}
func (m *chapterServiceStub) GetChapterStats(context.Context, string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
func (m *chapterServiceStub) GetChapterContent(context.Context, string, string) (string, error) {
	return "", nil
}
func (m *chapterServiceStub) GetChapterParagraphs(context.Context, string, string) ([]*bookstoreModel.ChapterContent, error) {
	return m.paragraphs, nil
}
func (m *chapterServiceStub) UpdateChapterContent(context.Context, string, string) error { return nil }
func (m *chapterServiceStub) PublishChapter(context.Context, string) error                { return nil }
func (m *chapterServiceStub) UnpublishChapter(context.Context, string) error              { return nil }
func (m *chapterServiceStub) BatchUpdateChapterPrice(context.Context, []string, float64) error {
	return nil
}
func (m *chapterServiceStub) BatchPublishChapters(context.Context, []string) error      { return nil }
func (m *chapterServiceStub) BatchDeleteChapters(context.Context, []string) error       { return nil }
func (m *chapterServiceStub) BatchDeleteChaptersByBookID(context.Context, string) error { return nil }
func (m *chapterServiceStub) SearchChapters(context.Context, string, int, int) ([]*bookstoreModel.Chapter, int64, error) {
	return nil, 0, nil
}

var _ bookstoreService.ChapterService = (*chapterServiceStub)(nil)

func mustObjectID(hex string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(hex)
	return id
}

func makeChapterFixture(chapterID, bookID string) (*bookstoreModel.Chapter, []*bookstoreModel.ChapterContent) {
	return &bookstoreModel.Chapter{
		ID:     mustObjectID(chapterID),
		BookID: bookID,
		Title:  "测试章节",
	}, []*bookstoreModel.ChapterContent{
		{
			ID:        primitive.NewObjectID(),
			ChapterID: mustObjectID(chapterID),
			Content:   "第一段内容",
		},
		{
			ID:        primitive.NewObjectID(),
			ChapterID: mustObjectID(chapterID),
			Content:   "第二段内容",
		},
	}
}

func setupChapterCommentTestRouter(userID string, repo *chapterCommentRepoStub, chapterSvc *chapterServiceStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("username", "tester")
		}
		c.Next()
	})

	api := NewChapterCommentAPI()
	api.BindServices(socialService.NewCommentService(repo, nil, nil), chapterSvc)

	v1 := r.Group("/api/v1/reader")
	v1.GET("/chapters/:chapterId/comments", api.GetChapterComments)
	v1.POST("/chapters/:chapterId/comments", api.CreateChapterComment)
	v1.GET("/comments/:commentId", api.GetChapterComment)
	v1.PUT("/comments/:commentId", api.UpdateChapterComment)
	v1.DELETE("/comments/:commentId", api.DeleteChapterComment)
	v1.POST("/comments/:commentId/like", api.LikeChapterComment)
	v1.DELETE("/comments/:commentId/like", api.UnlikeChapterComment)
	v1.GET("/chapters/:chapterId/paragraphs/:paragraphIndex/comments", api.GetParagraphComments)
	v1.POST("/chapters/:chapterId/paragraph-comments", api.CreateParagraphComment)
	v1.GET("/chapters/:chapterId/paragraph-comments", api.GetChapterParagraphComments)
	return r
}

func TestChapterCommentAPI_GetChapterComments_Success(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	repo := newChapterCommentRepoStub(&socialModels.Comment{
		IdentifiedEntity: socialModels.IdentifiedEntity{ID: primitive.NewObjectID()},
		Timestamps:       socialModels.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TargetType: socialModels.CommentTargetTypeChapter,
		TargetID:   chapterID,
		BookID:     bookID,
		ChapterID:  chapterID,
		AuthorID:   userID,
		Content:    "第一条章节评论",
		Rating:     5,
		State:      socialModels.CommentStateNormal,
		AuthorSnapshot: &socialModels.CommentAuthorSnapshot{
			ID:       userID,
			Username: "tester",
		},
	})
	router := setupChapterCommentTestRouter(userID, repo, &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID+"/comments?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"total":1`)
}

func TestChapterCommentAPI_GetChapterComments_Unauthorized(t *testing.T) {
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	router := setupChapterCommentTestRouter("", newChapterCommentRepoStub(), &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID+"/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestChapterCommentAPI_GetChapterComments_InvalidChapterID(t *testing.T) {
	router := setupChapterCommentTestRouter("user-1", newChapterCommentRepoStub(), &chapterServiceStub{})
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reader/chapters/invalid-id/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_CreateChapterComment_Success(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	repo := newChapterCommentRepoStub()
	router := setupChapterCommentTestRouter(userID, repo, &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	reqBody := readerModels.CreateChapterCommentRequest{
		Content: "这是一条很棒的评论",
		Rating:  5,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/reader/chapters/"+chapterID+"/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Len(t, repo.comments, 1)
	assert.Contains(t, w.Body.String(), "评论发表成功")
}

func TestChapterCommentAPI_CreateChapterComment_InvalidRating(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	router := setupChapterCommentTestRouter(userID, newChapterCommentRepoStub(), &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	reqBody := readerModels.CreateChapterCommentRequest{
		Content: "这是一条评论",
		Rating:  6,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/reader/chapters/"+chapterID+"/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_GetChapterComment_NotFound(t *testing.T) {
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	router := setupChapterCommentTestRouter("user-1", newChapterCommentRepoStub(), &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})
	commentID := primitive.NewObjectID().Hex()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reader/comments/"+commentID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestChapterCommentAPI_DeleteChapterComment_Success(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	comment := &socialModels.Comment{
		IdentifiedEntity: socialModels.IdentifiedEntity{ID: primitive.NewObjectID()},
		Timestamps:       socialModels.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TargetType: socialModels.CommentTargetTypeChapter,
		TargetID:   chapterID,
		BookID:     bookID,
		ChapterID:  chapterID,
		AuthorID:   userID,
		Content:    "待删除评论",
		State:      socialModels.CommentStateNormal,
	}
	repo := newChapterCommentRepoStub(comment)
	router := setupChapterCommentTestRouter(userID, repo, &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/reader/comments/"+comment.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	_, exists := repo.comments[comment.ID.Hex()]
	assert.False(t, exists)
}

func TestChapterCommentAPI_LikeChapterComment_Success(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	comment := &socialModels.Comment{
		IdentifiedEntity: socialModels.IdentifiedEntity{ID: primitive.NewObjectID()},
		Timestamps:       socialModels.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TargetType: socialModels.CommentTargetTypeChapter,
		TargetID:   chapterID,
		BookID:     bookID,
		ChapterID:  chapterID,
		AuthorID:   "another-user",
		Content:    "待点赞评论",
		State:      socialModels.CommentStateNormal,
	}
	repo := newChapterCommentRepoStub(comment)
	router := setupChapterCommentTestRouter(userID, repo, &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/reader/comments/"+comment.ID.Hex()+"/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), comment.LikeCount)
}

func TestChapterCommentAPI_GetParagraphComments_Success(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	paragraphID := paragraphs[0].ID.Hex()
	parentID := primitive.NewObjectID().Hex()

	topLevel := &socialModels.Comment{
		IdentifiedEntity: socialModels.IdentifiedEntity{ID: mustObjectID(parentID)},
		Timestamps:       socialModels.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TargetType: socialModels.CommentTargetTypeChapter,
		TargetID:   chapterID,
		BookID:     bookID,
		ChapterID:  chapterID,
		AuthorID:   userID,
		Content:    "段落主评论",
		State:      socialModels.CommentStateNormal,
		RichContent: map[string]interface{}{
			"paragraph_id":    paragraphID,
			"paragraph_index": 0,
			"paragraph_text":  paragraphs[0].Content,
		},
		AuthorSnapshot: &socialModels.CommentAuthorSnapshot{ID: userID, Username: "tester"},
	}
	replyParentID := parentID
	reply := &socialModels.Comment{
		IdentifiedEntity: socialModels.IdentifiedEntity{ID: primitive.NewObjectID()},
		Timestamps:       socialModels.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		ThreadedConversation: socialModels.ThreadedConversation{
			ParentID: &replyParentID,
		},
		TargetType: socialModels.CommentTargetTypeChapter,
		TargetID:   chapterID,
		BookID:     bookID,
		ChapterID:  chapterID,
		AuthorID:   "reply-user",
		Content:    "段落回复",
		State:      socialModels.CommentStateNormal,
		RichContent: map[string]interface{}{
			"paragraph_id":    paragraphID,
			"paragraph_index": 0,
			"paragraph_text":  paragraphs[0].Content,
		},
	}
	repo := newChapterCommentRepoStub(topLevel, reply)
	router := setupChapterCommentTestRouter(userID, repo, &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID+"/paragraphs/0/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"commentCount":2`)
}

func TestChapterCommentAPI_GetParagraphComments_InvalidIndex(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	router := setupChapterCommentTestRouter(userID, newChapterCommentRepoStub(), &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID+"/paragraphs/invalid/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_CreateParagraphComment_MissingParagraphIndex(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	router := setupChapterCommentTestRouter(userID, newChapterCommentRepoStub(), &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	reqBody := readerModels.CreateChapterCommentRequest{
		BookID:         bookID,
		Content:        "段落评论",
		ParagraphIndex: nil,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/reader/chapters/"+chapterID+"/paragraph-comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChapterCommentAPI_GetChapterParagraphComments_Success(t *testing.T) {
	userID := primitive.NewObjectID().Hex()
	chapterID := primitive.NewObjectID().Hex()
	bookID := primitive.NewObjectID().Hex()
	chapter, paragraphs := makeChapterFixture(chapterID, bookID)
	comment := &socialModels.Comment{
		IdentifiedEntity: socialModels.IdentifiedEntity{ID: primitive.NewObjectID()},
		Timestamps:       socialModels.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		TargetType: socialModels.CommentTargetTypeChapter,
		TargetID:   chapterID,
		BookID:     bookID,
		ChapterID:  chapterID,
		AuthorID:   userID,
		Content:    "段落概览评论",
		State:      socialModels.CommentStateNormal,
		RichContent: map[string]interface{}{
			"paragraph_id":    paragraphs[1].ID.Hex(),
			"paragraph_index": 1,
			"paragraph_text":  paragraphs[1].Content,
		},
		AuthorSnapshot: &socialModels.CommentAuthorSnapshot{ID: userID, Username: "tester"},
	}
	repo := newChapterCommentRepoStub(comment)
	router := setupChapterCommentTestRouter(userID, repo, &chapterServiceStub{chapter: chapter, paragraphs: paragraphs})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/reader/chapters/"+chapterID+"/paragraph-comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), fmt.Sprintf(`"paragraphId":"%s"`, paragraphs[1].ID.Hex()))
}
