package discovery

import (
	"context"
	"errors"
	"testing"

	bookstoreModel "Qingyu_backend/models/bookstore"
	discoveryModel "Qingyu_backend/models/discovery"
	sharedModel "Qingyu_backend/models/shared"
	sharedtypes "Qingyu_backend/models/shared/types"
	socialModel "Qingyu_backend/models/social"
	bookstoreService "Qingyu_backend/service/bookstore"
	recommendationService "Qingyu_backend/service/recommendation"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stubDiscoveryBookstoreService struct {
	recommendedBooks []*bookstoreModel.Book
	featuredBooks    []*bookstoreModel.Book
	hotBooks         []*bookstoreModel.Book
	newReleaseBooks  []*bookstoreModel.Book
	realtimeRanking  []*bookstoreModel.RankingItem
	weeklyRanking    []*bookstoreModel.RankingItem
	monthlyRanking   []*bookstoreModel.RankingItem
	booksByID        map[string]*bookstoreModel.Book
}

func (s *stubDiscoveryBookstoreService) GetAllBooks(context.Context, int, int) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) GetBookByID(_ context.Context, id string) (*bookstoreModel.Book, error) {
	return s.booksByID[id], nil
}

func (s *stubDiscoveryBookstoreService) GetBooksByCategory(context.Context, string, int, int) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) GetBooksByAuthorID(context.Context, string, int, int) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) GetRecommendedBooks(context.Context, int, int) ([]*bookstoreModel.Book, int64, error) {
	return s.recommendedBooks, int64(len(s.recommendedBooks)), nil
}

func (s *stubDiscoveryBookstoreService) GetFeaturedBooks(context.Context, int, int) ([]*bookstoreModel.Book, int64, error) {
	return s.featuredBooks, int64(len(s.featuredBooks)), nil
}

func (s *stubDiscoveryBookstoreService) GetHotBooks(context.Context, int, int) ([]*bookstoreModel.Book, int64, error) {
	return s.hotBooks, int64(len(s.hotBooks)), nil
}

func (s *stubDiscoveryBookstoreService) GetNewReleases(context.Context, int, int) ([]*bookstoreModel.Book, int64, error) {
	return s.newReleaseBooks, int64(len(s.newReleaseBooks)), nil
}

func (s *stubDiscoveryBookstoreService) GetFreeBooks(context.Context, int, int) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) SearchBooks(context.Context, string, int, int) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) SearchBooksWithFilter(context.Context, *bookstoreModel.BookFilter) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) GetCategoryTree(context.Context) ([]*bookstoreModel.CategoryTree, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) GetCategoryByID(context.Context, string) (*bookstoreModel.Category, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) GetRootCategories(context.Context) ([]*bookstoreModel.Category, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) GetActiveBanners(context.Context, int) ([]*bookstoreModel.Banner, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) IncrementBannerClick(context.Context, string) error {
	return nil
}

func (s *stubDiscoveryBookstoreService) GetRealtimeRanking(context.Context, int) ([]*bookstoreModel.RankingItem, error) {
	return s.realtimeRanking, nil
}

func (s *stubDiscoveryBookstoreService) GetWeeklyRanking(context.Context, string, int) ([]*bookstoreModel.RankingItem, error) {
	return s.weeklyRanking, nil
}

func (s *stubDiscoveryBookstoreService) GetMonthlyRanking(context.Context, string, int) ([]*bookstoreModel.RankingItem, error) {
	return s.monthlyRanking, nil
}

func (s *stubDiscoveryBookstoreService) GetNewbieRanking(context.Context, string, int) ([]*bookstoreModel.RankingItem, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) GetRankingByType(context.Context, bookstoreModel.RankingType, string, int) ([]*bookstoreModel.RankingItem, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) UpdateRankings(context.Context, bookstoreModel.RankingType, string) error {
	return nil
}

func (s *stubDiscoveryBookstoreService) GetHomepageData(context.Context) (*bookstoreService.HomepageData, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) GetBookStats(context.Context) (*bookstoreModel.BookStats, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) IncrementBookView(context.Context, string) error {
	return nil
}

func (s *stubDiscoveryBookstoreService) GetYears(context.Context) ([]int, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) GetTags(context.Context, *string) ([]string, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) SearchByTitle(context.Context, string, int, int) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) SearchByAuthor(context.Context, string, int, int) ([]*bookstoreModel.Book, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryBookstoreService) GetSimilarBooks(context.Context, string, int) ([]*bookstoreModel.Book, error) {
	return nil, nil
}

func (s *stubDiscoveryBookstoreService) SetSearchService(interface{}) {}

type stubDiscoveryRecommendationService struct {
	recorded []*recommendationService.RecordBehaviorRequest
}

type stubDiscoveryPostService struct {
	posts []*socialModel.PostInfo
}

type stubDiscoveryPreferenceRepository struct {
	profile      *discoveryModel.PreferenceProfile
	upsertedUser string
	upserted     *discoveryModel.PreferenceProfile
	getErr       error
	upsertErr    error
}

func (s *stubDiscoveryRecommendationService) GetPersonalizedRecommendations(context.Context, string, int) ([]*recommendationService.RecommendedItem, error) {
	return nil, nil
}

func (s *stubDiscoveryRecommendationService) GetSimilarItems(context.Context, string, int) ([]*recommendationService.RecommendedItem, error) {
	return nil, nil
}

func (s *stubDiscoveryRecommendationService) GetHotItems(context.Context, string, int) ([]*recommendationService.RecommendedItem, error) {
	return nil, nil
}

func (s *stubDiscoveryRecommendationService) RecordUserBehavior(_ context.Context, req *recommendationService.RecordBehaviorRequest) error {
	s.recorded = append(s.recorded, req)
	return nil
}

func (s *stubDiscoveryRecommendationService) GetUserBehaviors(context.Context, string, int) ([]*recommendationService.UserBehavior, error) {
	return nil, nil
}

func (s *stubDiscoveryRecommendationService) RefreshRecommendations(context.Context, string) error {
	return nil
}

func (s *stubDiscoveryRecommendationService) RefreshHotItems(context.Context, string) error {
	return nil
}

func (s *stubDiscoveryRecommendationService) Health(context.Context) error {
	return nil
}

func (s *stubDiscoveryPostService) CreatePost(context.Context, string, string, string, int, socialModel.PostType, string, []string, string, string, string, string, string, string, int, []string) (*socialModel.Post, error) {
	return nil, nil
}

func (s *stubDiscoveryPostService) GetPosts(context.Context, string, int, int, string, string) ([]*socialModel.PostInfo, int64, error) {
	return s.posts, int64(len(s.posts)), nil
}

func (s *stubDiscoveryPostService) GetPostByID(context.Context, string, string) (*socialModel.PostInfo, error) {
	return nil, nil
}

func (s *stubDiscoveryPostService) UpdatePost(context.Context, string, string, string, []string) error {
	return nil
}

func (s *stubDiscoveryPostService) DeletePost(context.Context, string, string) error {
	return nil
}

func (s *stubDiscoveryPostService) GetUserPosts(context.Context, string, int, int) ([]*socialModel.PostInfo, int64, error) {
	return nil, 0, nil
}

func (s *stubDiscoveryPostService) ToggleLike(context.Context, string, string) (bool, error) {
	return false, nil
}

func (s *stubDiscoveryPostService) Like(context.Context, string, string) error {
	return nil
}

func (s *stubDiscoveryPostService) Unlike(context.Context, string, string) error {
	return nil
}

func (s *stubDiscoveryPreferenceRepository) GetByUserID(context.Context, string) (*discoveryModel.PreferenceProfile, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.profile, nil
}

func (s *stubDiscoveryPreferenceRepository) UpsertByUserID(_ context.Context, userID string, profile *discoveryModel.PreferenceProfile) error {
	if s.upsertErr != nil {
		return s.upsertErr
	}
	s.upsertedUser = userID
	s.upserted = profile
	return nil
}

func buildDiscoveryTestBook(id primitive.ObjectID, title string) *bookstoreModel.Book {
	return &bookstoreModel.Book{
		IdentifiedEntity: sharedModel.IdentifiedEntity{ID: id},
		Title:            title,
		Author:           "测试作者",
		Introduction:     "测试简介",
		Cover:            "/images/covers/test.jpg",
		Categories:       []string{"玄幻"},
		Tags:             []string{"热血", "冒险"},
		Rating:           sharedtypes.Rating(4.8),
	}
}

func TestGetTrending_DailyMapsToRealtimeRanking(t *testing.T) {
	bookID := primitive.NewObjectID()
	book := buildDiscoveryTestBook(bookID, "榜单图书")
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{
		realtimeRanking: []*bookstoreModel.RankingItem{{
			BookID:    bookID,
			Book:      book,
			Rank:      1,
			Score:     98.6,
			ViewCount: 3200,
			LikeCount: 240,
			Period:    "2026-04-20",
		}},
	}, nil)

	result, err := service.GetTrending(context.Background(), "daily", 8)
	if err != nil {
		t.Fatalf("expected trending request to succeed: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 trending item, got %d", len(result))
	}
	if result[0].RankingType != "daily" {
		t.Fatalf("expected ranking type to normalize to daily, got %s", result[0].RankingType)
	}
	if result[0].Rank != 1 || result[0].Title != "榜单图书" {
		t.Fatalf("expected ranked book payload, got %+v", result[0])
	}
}

func TestTrackAction_RecordsSupportedBehavior(t *testing.T) {
	recommendation := &stubDiscoveryRecommendationService{}
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{}, recommendation)

	err := service.TrackAction(context.Background(), "user-1", TrackActionRequest{
		ItemID: "book-1",
		Action: "like",
	})
	if err != nil {
		t.Fatalf("expected supported action to be recorded: %v", err)
	}
	if len(recommendation.recorded) != 1 {
		t.Fatalf("expected 1 behavior record, got %d", len(recommendation.recorded))
	}
	if recommendation.recorded[0].ItemType != "book" {
		t.Fatalf("expected blank item type to default to book, got %+v", recommendation.recorded[0])
	}
	if recommendation.recorded[0].ActionType != "like" {
		t.Fatalf("expected like action to be preserved, got %+v", recommendation.recorded[0])
	}
	if recommendation.recorded[0].Metadata["source"] != "discovery" {
		t.Fatalf("expected source metadata to be injected, got %+v", recommendation.recorded[0].Metadata)
	}
}

func TestTrackAction_DislikeIsRecordedAndAnonymousIsIgnored(t *testing.T) {
	recommendation := &stubDiscoveryRecommendationService{}
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{}, recommendation)

	if err := service.TrackAction(context.Background(), "", TrackActionRequest{ItemID: "book-1", Action: "view"}); err != nil {
		t.Fatalf("expected anonymous action to be ignored without error: %v", err)
	}
	if err := service.TrackAction(context.Background(), "user-1", TrackActionRequest{ItemID: "book-1", Action: "dislike"}); err != nil {
		t.Fatalf("expected dislike action to be recorded without error: %v", err)
	}
	if len(recommendation.recorded) != 1 || recommendation.recorded[0].ActionType != "dislike" {
		t.Fatalf("expected dislike action to be persisted, got %+v", recommendation.recorded)
	}
}

func TestGetRecommendations_FiltersBySlot(t *testing.T) {
	bookID := primitive.NewObjectID()
	book := buildDiscoveryTestBook(bookID, "推荐图书")
	customBookstore := &stubDiscoveryBookstoreService{
		featuredBooks:    []*bookstoreModel.Book{buildDiscoveryTestBook(primitive.NewObjectID(), "精选图书")},
		recommendedBooks: []*bookstoreModel.Book{buildDiscoveryTestBook(primitive.NewObjectID(), "推荐位图书")},
		newReleaseBooks:  []*bookstoreModel.Book{book},
		booksByID:        map[string]*bookstoreModel.Book{bookID.Hex(): book},
	}
	service := NewDiscoveryService(customBookstore, nil)

	result, err := service.GetRecommendations(context.Background(), slotNewReleases)
	if err != nil {
		t.Fatalf("expected recommendation slot request to succeed: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 new release recommendation, got %+v", result)
	}
	if result[0].Slot != slotNewReleases || result[0].Title != "推荐图书" {
		t.Fatalf("expected slot-filtered recommendation payload, got %+v", result[0])
	}
}

func TestUpdatePreferences_NormalizesValues(t *testing.T) {
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{}, nil)

	result, err := service.UpdatePreferences(context.Background(), "user-1", RecommendationConfig{
		Preferences: RecommendationPreferences{
			FavoriteGenres:  []string{" 玄幻 ", "玄幻", ""},
			FavoriteAuthors: []string{"作者A", "作者A"},
			FavoriteTags:    []string{"热血", " 热血 "},
		},
		History: RecommendationHistory{
			ViewedBooks: []string{" book-1 ", "book-1"},
			ViewedLists: []string{"list-1", ""},
		},
	})
	if err != nil {
		t.Fatalf("expected preferences update to succeed: %v", err)
	}

	if result.UserID != "user-1" {
		t.Fatalf("expected user id to be preserved, got %+v", result)
	}
	if len(result.Preferences.FavoriteGenres) != 1 || result.Preferences.FavoriteGenres[0] != "玄幻" {
		t.Fatalf("expected genres to be normalized, got %+v", result.Preferences.FavoriteGenres)
	}
	if len(result.History.ViewedBooks) != 1 || result.History.ViewedBooks[0] != "book-1" {
		t.Fatalf("expected viewed books to be normalized, got %+v", result.History.ViewedBooks)
	}
}

func TestUpdatePreferences_PersistsAuthenticatedUserAndIgnoresRequestUserID(t *testing.T) {
	repo := &stubDiscoveryPreferenceRepository{
		profile: &discoveryModel.PreferenceProfile{UserID: "user-1"},
	}
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{}, nil).WithPreferenceRepository(repo)

	result, err := service.UpdatePreferences(context.Background(), "user-1", RecommendationConfig{
		UserID: "spoofed-user",
		Preferences: RecommendationPreferences{
			FavoriteGenres:  []string{"科幻"},
			FavoriteAuthors: []string{"作者A"},
			FavoriteTags:    []string{"烧脑"},
		},
		History: RecommendationHistory{
			ViewedBooks: []string{"book-9"},
			ViewedLists: []string{"list-2"},
		},
	})
	if err != nil {
		t.Fatalf("expected authenticated preferences update to succeed: %v", err)
	}

	if result.UserID != "user-1" {
		t.Fatalf("expected result user id to come from auth context, got %+v", result)
	}
	if repo.upsertedUser != "user-1" || repo.upserted == nil {
		t.Fatalf("expected repo upsert to use authenticated user, got user=%s profile=%+v", repo.upsertedUser, repo.upserted)
	}
	if repo.upserted.UserID != "user-1" {
		t.Fatalf("expected persisted profile to use authenticated user, got %+v", repo.upserted)
	}
	if len(repo.upserted.FavoriteGenres) != 1 || repo.upserted.FavoriteGenres[0] != "科幻" {
		t.Fatalf("expected normalized genres to be persisted, got %+v", repo.upserted.FavoriteGenres)
	}
}

func TestUpdatePreferences_AnonymousOnlyNormalizesWithoutPersisting(t *testing.T) {
	repo := &stubDiscoveryPreferenceRepository{}
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{}, nil).WithPreferenceRepository(repo)

	result, err := service.UpdatePreferences(context.Background(), "", RecommendationConfig{
		UserID: "spoofed-user",
		Preferences: RecommendationPreferences{
			FavoriteTags: []string{" 热门 ", "热门"},
		},
	})
	if err != nil {
		t.Fatalf("expected anonymous preferences update to succeed: %v", err)
	}

	if result.UserID != "" {
		t.Fatalf("expected anonymous request to return empty user id, got %+v", result)
	}
	if repo.upserted != nil {
		t.Fatalf("expected anonymous request to skip persistence, got %+v", repo.upserted)
	}
}

func TestUpdatePreferences_ReturnsPersistenceError(t *testing.T) {
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{}, nil).WithPreferenceRepository(&stubDiscoveryPreferenceRepository{
		upsertErr: errors.New("boom"),
	})

	_, err := service.UpdatePreferences(context.Background(), "user-1", RecommendationConfig{})
	if err == nil {
		t.Fatal("expected persistence error to be returned")
	}
}

func TestGetTopics_PrefersCommunityPosts(t *testing.T) {
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{}, nil).WithPostService(&stubDiscoveryPostService{
		posts: []*socialModel.PostInfo{
			{
				Topics:       []string{"科幻", "推荐"},
				LikeCount:    12,
				CommentCount: 5,
				ShareCount:   2,
				Book: &socialModel.BookInfo{
					BookID: "book-1",
					Title:  "星海回声",
					Cover:  "/images/covers/scifi.jpg",
					Author: "测试作者A",
				},
			},
			{
				Topics:       []string{"科幻"},
				LikeCount:    6,
				CommentCount: 3,
				ShareCount:   1,
				Book: &socialModel.BookInfo{
					BookID: "book-2",
					Title:  "深空档案",
					Cover:  "/images/covers/scifi-2.jpg",
					Author: "测试作者B",
				},
			},
		},
	})

	result, err := service.GetTopics(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("expected community topics to load: %v", err)
	}
	if len(result.List) == 0 {
		t.Fatal("expected community topics to produce discovery topics")
	}
	if result.List[0].Title != "科幻" {
		t.Fatalf("expected hottest community topic to be 科幻, got %+v", result.List[0])
	}
	if len(result.List[0].Books) != 2 {
		t.Fatalf("expected topic books to be aggregated from community posts, got %+v", result.List[0].Books)
	}
}

func TestGetTopics_UsesCommunityTopicNamesEvenWithoutBookBinding(t *testing.T) {
	service := NewDiscoveryService(&stubDiscoveryBookstoreService{
		recommendedBooks: []*bookstoreModel.Book{buildDiscoveryTestBook(primitive.NewObjectID(), "兜底图书")},
	}, nil).WithPostService(&stubDiscoveryPostService{
		posts: []*socialModel.PostInfo{
			{
				Topics:       []string{"现实主义"},
				LikeCount:    8,
				CommentCount: 4,
				ShareCount:   1,
			},
		},
	})

	result, err := service.GetTopics(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("expected community-only topic to load: %v", err)
	}
	if len(result.List) != 1 || result.List[0].Title != "现实主义" {
		t.Fatalf("expected community topic name to be preserved, got %+v", result.List)
	}
	if len(result.List[0].Books) != 1 || result.List[0].Books[0].Title != "兜底图书" {
		t.Fatalf("expected fallback books to decorate community topic, got %+v", result.List[0].Books)
	}
}
