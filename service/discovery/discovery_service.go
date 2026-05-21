package discovery

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	bookstoreModel "Qingyu_backend/models/bookstore"
	discoveryModel "Qingyu_backend/models/discovery"
	sharedtypes "Qingyu_backend/models/shared/types"
	socialModel "Qingyu_backend/models/social"
	discoveryRepo "Qingyu_backend/repository/interfaces/discovery"
	bookstoreService "Qingyu_backend/service/bookstore"
	postInterfaces "Qingyu_backend/service/interfaces"
	recommendationService "Qingyu_backend/service/recommendation"
)

const (
	defaultBookDescription = "这本书暂时还没有补充简介。"
	defaultBookAuthor      = "佚名"
	defaultBookTitle       = "未命名书籍"
	timeLayoutRFC3339      = "2006-01-02T15:04:05Z07:00"
	defaultTrendingType    = "daily"
	defaultDiscoveryItem   = "book"
	slotBanner             = "banner"
	slotFeatured           = "featured"
	slotNewReleases        = "new_releases"
	slotEditorsPick        = "editors_pick"
)

var (
	ErrDiscoveryTrackItemIDRequired    = errors.New("itemId is required")
	ErrDiscoveryTrackActionRequired    = errors.New("action is required")
	ErrDiscoveryTrackActionUnsupported = errors.New("unsupported track action")
)

type DiscoveryService struct {
	bookstore      bookstoreService.BookstoreService
	recommendation recommendationService.RecommendationService
	postService    postInterfaces.PostService
	preferenceRepo discoveryRepo.PreferenceRepository
}

type PersonalizedResponse struct {
	Books     PersonalizedBooks     `json:"books"`
	Booklists PersonalizedBooklists `json:"booklists"`
	Authors   PersonalizedAuthors   `json:"authors"`
}

type PersonalizedBooks struct {
	ForYou   []BookBrief `json:"forYou"`
	Similar  []BookBrief `json:"similar"`
	Trending []BookBrief `json:"trending"`
}

type PersonalizedBooklists struct {
	Recommended []BooklistBrief `json:"recommended"`
	Popular     []BooklistBrief `json:"popular"`
}

type PersonalizedAuthors struct {
	Suggested []AuthorBrief `json:"suggested"`
}

type BookBrief struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Author      string   `json:"author"`
	Rating      float64  `json:"rating"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Category    string   `json:"category,omitempty"`
	PublishDate string   `json:"publishDate,omitempty"`
}

type BooklistBrief struct {
	ID string `json:"id"`
}

type AuthorBrief struct {
	ID string `json:"id"`
}

type BookListResponse struct {
	List  []BookBrief `json:"list"`
	Total int64       `json:"total"`
}

type EditorPick struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Cover        string  `json:"cover"`
	Description  string  `json:"description,omitempty"`
	Reason       string  `json:"reason"`
	ViewCount    int64   `json:"viewCount"`
	CollectCount int64   `json:"collectCount"`
	Rating       float64 `json:"rating"`
}

type EditorPickListResponse struct {
	List  []EditorPick `json:"list"`
	Total int64        `json:"total"`
}

type TrendingBook struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Cover       string   `json:"cover"`
	Author      string   `json:"author"`
	Rating      float64  `json:"rating"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
	Category    string   `json:"category,omitempty"`
	PublishDate string   `json:"publishDate,omitempty"`
	Rank        int      `json:"rank"`
	RankingType string   `json:"rankingType"`
	Period      string   `json:"period,omitempty"`
	Score       float64  `json:"score"`
	ViewCount   int64    `json:"viewCount"`
	LikeCount   int64    `json:"likeCount"`
}

type TrackActionRequest struct {
	ItemID   string                 `json:"itemId"`
	ItemType string                 `json:"itemType"`
	Action   string                 `json:"action"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type RecommendationItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Slot        string `json:"slot"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Link        string `json:"link"`
	Priority    int    `json:"priority"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}

type DiscoveryTopic struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Books       []BookBrief `json:"books"`
}

type DiscoveryTopicListResponse struct {
	List  []DiscoveryTopic `json:"list"`
	Total int64            `json:"total"`
}

type RecommendationPreferences struct {
	FavoriteGenres  []string `json:"favoriteGenres"`
	FavoriteAuthors []string `json:"favoriteAuthors"`
	FavoriteTags    []string `json:"favoriteTags"`
}

type RecommendationHistory struct {
	ViewedBooks []string `json:"viewedBooks"`
	ViewedLists []string `json:"viewedLists"`
}

type RecommendationConfig struct {
	UserID      string                    `json:"userId"`
	Preferences RecommendationPreferences `json:"preferences"`
	History     RecommendationHistory     `json:"history"`
}

func NewDiscoveryService(bookstore bookstoreService.BookstoreService, recommendation recommendationService.RecommendationService) *DiscoveryService {
	return &DiscoveryService{
		bookstore:      bookstore,
		recommendation: recommendation,
	}
}

func (s *DiscoveryService) WithPostService(postService postInterfaces.PostService) *DiscoveryService {
	s.postService = postService
	return s
}

func (s *DiscoveryService) WithPreferenceRepository(preferenceRepo discoveryRepo.PreferenceRepository) *DiscoveryService {
	s.preferenceRepo = preferenceRepo
	return s
}

func (s *DiscoveryService) GetPersonalized(ctx context.Context, userID string, limit int) (*PersonalizedResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	books := make([]BookBrief, 0, limit)
	if userID != "" && s.recommendation != nil {
		items, err := s.recommendation.GetPersonalizedRecommendations(ctx, userID, limit)
		if err == nil {
			books = s.hydrateRecommendedItems(ctx, items, limit)
		}
	}

	if len(books) == 0 {
		fallback, err := s.getFallbackBooks(ctx, limit)
		if err != nil {
			return nil, err
		}
		books = fallback
	}

	return &PersonalizedResponse{
		Books: PersonalizedBooks{
			ForYou:   books,
			Similar:  []BookBrief{},
			Trending: []BookBrief{},
		},
		Booklists: PersonalizedBooklists{
			Recommended: []BooklistBrief{},
			Popular:     []BooklistBrief{},
		},
		Authors: PersonalizedAuthors{
			Suggested: []AuthorBrief{},
		},
	}, nil
}

func (s *DiscoveryService) GetNewReleases(ctx context.Context, page, size int, categoryID string) (*BookListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var (
		books []*bookstoreModel.Book
		total int64
		err   error
	)

	if categoryID == "" {
		books, total, err = s.bookstore.GetNewReleases(ctx, page, size)
	} else {
		books, total, err = s.bookstore.GetBooksByCategory(ctx, categoryID, page, size)
	}
	if err != nil {
		return nil, err
	}

	return &BookListResponse{
		List:  toBookBriefs(books),
		Total: total,
	}, nil
}

func (s *DiscoveryService) GetEditorsPick(ctx context.Context, page, size int, tag string) (*EditorPickListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	books, total, err := s.bookstore.GetFeaturedBooks(ctx, page, size)
	if err != nil {
		return nil, err
	}

	list := make([]EditorPick, 0, len(books))
	for _, book := range books {
		if book == nil {
			continue
		}
		if tag != "" && !containsTag(book.Tags, tag) {
			continue
		}

		list = append(list, EditorPick{
			ID:           book.ID.Hex(),
			Title:        normalizeTitle(book.Title),
			Author:       normalizeAuthor(book.Author),
			Cover:        book.Cover,
			Description:  normalizeDescription(book.Introduction),
			Reason:       buildEditorReason(book),
			ViewCount:    book.ViewCount,
			CollectCount: book.CollectCount,
			Rating:       book.Rating.ToFloat(),
		})
	}

	return &EditorPickListResponse{
		List:  list,
		Total: total,
	}, nil
}

func (s *DiscoveryService) GetTrending(ctx context.Context, rankingType string, limit int) ([]TrendingBook, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if s.bookstore == nil {
		return []TrendingBook{}, nil
	}

	normalizedType := normalizeTrendingType(rankingType)

	var (
		items []*bookstoreModel.RankingItem
		err   error
	)

	switch normalizedType {
	case "weekly":
		items, err = s.bookstore.GetWeeklyRanking(ctx, "", limit)
	case "monthly":
		items, err = s.bookstore.GetMonthlyRanking(ctx, "", limit)
	default:
		items, err = s.bookstore.GetRealtimeRanking(ctx, limit)
	}
	if err != nil {
		return nil, err
	}

	result := make([]TrendingBook, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		book := item.Book
		if book == nil && !item.BookID.IsZero() {
			hydratedBook, hydrateErr := s.bookstore.GetBookByID(ctx, item.BookID.Hex())
			if hydrateErr == nil {
				book = hydratedBook
			}
		}
		if book == nil {
			continue
		}

		brief := toBookBrief(book)
		result = append(result, TrendingBook{
			ID:          brief.ID,
			Title:       brief.Title,
			Cover:       brief.Cover,
			Author:      brief.Author,
			Rating:      brief.Rating,
			Description: brief.Description,
			Tags:        brief.Tags,
			Category:    brief.Category,
			PublishDate: brief.PublishDate,
			Rank:        item.Rank,
			RankingType: normalizedType,
			Period:      item.Period,
			Score:       item.Score,
			ViewCount:   item.ViewCount,
			LikeCount:   item.LikeCount,
		})
	}

	return result, nil
}

func (s *DiscoveryService) TrackAction(ctx context.Context, userID string, req TrackActionRequest) error {
	itemID := strings.TrimSpace(req.ItemID)
	if itemID == "" {
		return ErrDiscoveryTrackItemIDRequired
	}

	action := strings.ToLower(strings.TrimSpace(req.Action))
	if action == "" {
		return ErrDiscoveryTrackActionRequired
	}
	if _, err := sharedtypes.ParseRecommendationBehaviorType(action); err != nil {
		return ErrDiscoveryTrackActionUnsupported
	}
	if strings.TrimSpace(userID) == "" || s.recommendation == nil {
		return nil
	}

	itemType := strings.TrimSpace(req.ItemType)
	if itemType == "" {
		itemType = defaultDiscoveryItem
	}

	metadata := make(map[string]interface{}, len(req.Metadata)+1)
	for key, value := range req.Metadata {
		metadata[key] = value
	}
	if _, exists := metadata["source"]; !exists {
		metadata["source"] = "discovery"
	}

	return s.recommendation.RecordUserBehavior(ctx, &recommendationService.RecordBehaviorRequest{
		UserID:     userID,
		ItemID:     itemID,
		ItemType:   itemType,
		ActionType: action,
		Metadata:   metadata,
	})
}

func (s *DiscoveryService) GetRecommendations(ctx context.Context, slot string) ([]RecommendationItem, error) {
	if s.bookstore == nil {
		return []RecommendationItem{}, nil
	}

	now := time.Now()
	startTime := now.Add(-24 * time.Hour).Format(timeLayoutRFC3339)
	endTime := now.Add(14 * 24 * time.Hour).Format(timeLayoutRFC3339)

	items := make([]RecommendationItem, 0, 12)
	priority := 1

	appendBooks := func(slotName string, books []*bookstoreModel.Book, limit int) {
		for _, book := range books {
			if book == nil {
				continue
			}
			items = append(items, RecommendationItem{
				ID:          "discovery-" + slotName + "-" + book.ID.Hex(),
				Type:        defaultDiscoveryItem,
				Slot:        slotName,
				Title:       normalizeTitle(book.Title),
				Description: normalizeDescription(book.Introduction),
				Cover:       book.Cover,
				Link:        "/bookstore/book/" + book.ID.Hex(),
				Priority:    priority,
				StartTime:   startTime,
				EndTime:     endTime,
			})
			priority++
			if limit > 0 && priority > limit {
				return
			}
		}
	}

	addSlot := func(slotName string, loader func(context.Context, int, int) ([]*bookstoreModel.Book, int64, error), size int) error {
		if slot != "" && slot != slotName {
			return nil
		}
		books, _, err := loader(ctx, 1, size)
		if err != nil {
			return err
		}
		appendBooks(slotName, books, 0)
		return nil
	}

	if err := addSlot(slotBanner, s.bookstore.GetFeaturedBooks, 2); err != nil {
		return nil, err
	}
	if err := addSlot(slotFeatured, s.bookstore.GetRecommendedBooks, 4); err != nil {
		return nil, err
	}
	if err := addSlot(slotNewReleases, s.bookstore.GetNewReleases, 4); err != nil {
		return nil, err
	}
	if err := addSlot(slotEditorsPick, s.bookstore.GetFeaturedBooks, 4); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *DiscoveryService) GetTopics(ctx context.Context, page, size int) (*DiscoveryTopicListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	if s.bookstore == nil {
		return &DiscoveryTopicListResponse{List: []DiscoveryTopic{}, Total: 0}, nil
	}
	if s.postService != nil {
		topics, err := s.getTopicsFromCommunity(ctx, page, size)
		if err == nil && len(topics.List) > 0 {
			return topics, nil
		}
	}

	featured, _, featuredErr := s.bookstore.GetFeaturedBooks(ctx, 1, 4)
	if featuredErr != nil {
		featured = []*bookstoreModel.Book{}
	}
	newReleases, _, newErr := s.bookstore.GetNewReleases(ctx, 1, 4)
	if newErr != nil {
		newReleases = []*bookstoreModel.Book{}
	}
	recommended, _, recommendedErr := s.bookstore.GetRecommendedBooks(ctx, 1, 4)
	if recommendedErr != nil {
		recommended = []*bookstoreModel.Book{}
	}

	topics := []DiscoveryTopic{
		{
			ID:          "topic-discovery-1",
			Title:       "编辑部热推",
			Description: "适合想先看站内主推内容的读者，优先收录近期精选与高完成度作品。",
			Books:       toBookBriefs(featured),
		},
		{
			ID:          "topic-discovery-2",
			Title:       "新读者友好",
			Description: "优先展示近期上架、上手门槛低的新书，适合快速建立阅读兴趣。",
			Books:       toBookBriefs(newReleases),
		},
		{
			ID:          "topic-discovery-3",
			Title:       "高分延读",
			Description: "优先挑选评价稳定、适合持续阅读的作品，减少挑书成本。",
			Books:       toBookBriefs(recommended),
		},
	}

	filtered := make([]DiscoveryTopic, 0, len(topics))
	for _, topic := range topics {
		if len(topic.Books) == 0 {
			continue
		}
		filtered = append(filtered, topic)
	}

	start := (page - 1) * size
	if start >= len(filtered) {
		return &DiscoveryTopicListResponse{List: []DiscoveryTopic{}, Total: int64(len(filtered))}, nil
	}
	end := start + size
	if end > len(filtered) {
		end = len(filtered)
	}

	return &DiscoveryTopicListResponse{
		List:  filtered[start:end],
		Total: int64(len(filtered)),
	}, nil
}

func (s *DiscoveryService) getTopicsFromCommunity(ctx context.Context, page, size int) (*DiscoveryTopicListResponse, error) {
	posts, _, err := s.postService.GetPosts(ctx, "", 1, 50, "", "hottest")
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return &DiscoveryTopicListResponse{List: []DiscoveryTopic{}, Total: 0}, nil
	}

	type topicAggregate struct {
		name      string
		score     int
		books     []BookBrief
		seenBooks map[string]struct{}
	}

	aggregates := make(map[string]*topicAggregate)
	for _, post := range posts {
		if post == nil || len(post.Topics) == 0 {
			continue
		}
		for _, topic := range post.Topics {
			normalizedTopic := strings.TrimSpace(topic)
			if normalizedTopic == "" {
				continue
			}

			aggregate, exists := aggregates[normalizedTopic]
			if !exists {
				aggregate = &topicAggregate{
					name:      normalizedTopic,
					books:     make([]BookBrief, 0, 4),
					seenBooks: make(map[string]struct{}),
				}
				aggregates[normalizedTopic] = aggregate
			}

			aggregate.score += post.LikeCount + post.CommentCount + post.ShareCount + 1

			if post.Book == nil || post.Book.BookID == "" {
				continue
			}
			if _, seen := aggregate.seenBooks[post.Book.BookID]; seen {
				continue
			}

			bookBrief := s.toBookBriefFromPostBook(ctx, post.Book)
			aggregate.books = append(aggregate.books, bookBrief)
			aggregate.seenBooks[post.Book.BookID] = struct{}{}
		}
	}

	if len(aggregates) == 0 {
		return &DiscoveryTopicListResponse{List: []DiscoveryTopic{}, Total: 0}, nil
	}

	items := make([]*topicAggregate, 0, len(aggregates))
	for _, aggregate := range aggregates {
		items = append(items, aggregate)
	}
	if len(items) == 0 {
		return &DiscoveryTopicListResponse{List: []DiscoveryTopic{}, Total: 0}, nil
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].score == items[j].score {
			return items[i].name < items[j].name
		}
		return items[i].score > items[j].score
	})

	start := (page - 1) * size
	if start >= len(items) {
		return &DiscoveryTopicListResponse{List: []DiscoveryTopic{}, Total: int64(len(items))}, nil
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	fallbackBooks, _ := s.getFallbackBooks(ctx, 4)

	result := make([]DiscoveryTopic, 0, end-start)
	for _, item := range items[start:end] {
		books := item.books
		if len(books) == 0 {
			books = append([]BookBrief{}, fallbackBooks...)
		}
		result = append(result, DiscoveryTopic{
			ID:          "topic-community-" + strings.ToLower(strings.ReplaceAll(item.name, " ", "-")),
			Title:       item.name,
			Description: buildCommunityTopicDescription(item.name, item.score),
			Books:       books,
		})
	}

	return &DiscoveryTopicListResponse{
		List:  result,
		Total: int64(len(items)),
	}, nil
}

func (s *DiscoveryService) UpdatePreferences(ctx context.Context, userID string, config RecommendationConfig) (RecommendationConfig, error) {
	config.UserID = ""
	config.Preferences.FavoriteGenres = normalizeStringSlice(config.Preferences.FavoriteGenres)
	config.Preferences.FavoriteAuthors = normalizeStringSlice(config.Preferences.FavoriteAuthors)
	config.Preferences.FavoriteTags = normalizeStringSlice(config.Preferences.FavoriteTags)
	config.History.ViewedBooks = normalizeStringSlice(config.History.ViewedBooks)
	config.History.ViewedLists = normalizeStringSlice(config.History.ViewedLists)

	authenticatedUserID := strings.TrimSpace(userID)
	if authenticatedUserID == "" {
		return config, nil
	}

	config.UserID = authenticatedUserID
	if s.preferenceRepo == nil {
		return config, nil
	}

	profile := &discoveryModel.PreferenceProfile{
		UserID:          authenticatedUserID,
		FavoriteGenres:  append([]string{}, config.Preferences.FavoriteGenres...),
		FavoriteAuthors: append([]string{}, config.Preferences.FavoriteAuthors...),
		FavoriteTags:    append([]string{}, config.Preferences.FavoriteTags...),
		ViewedBooks:     append([]string{}, config.History.ViewedBooks...),
		ViewedLists:     append([]string{}, config.History.ViewedLists...),
	}

	existing, err := s.preferenceRepo.GetByUserID(ctx, authenticatedUserID)
	if err != nil {
		return RecommendationConfig{}, err
	}
	if existing != nil {
		profile.CreatedAt = existing.CreatedAt
	}

	if err := s.preferenceRepo.UpsertByUserID(ctx, authenticatedUserID, profile); err != nil {
		return RecommendationConfig{}, err
	}

	return config, nil
}

func (s *DiscoveryService) hydrateRecommendedItems(ctx context.Context, items []*recommendationService.RecommendedItem, limit int) []BookBrief {
	if s.bookstore == nil || len(items) == 0 {
		return []BookBrief{}
	}

	seen := make(map[string]struct{}, len(items))
	result := make([]BookBrief, 0, limit)

	for _, item := range items {
		if item == nil || item.ItemID == "" {
			continue
		}
		if _, exists := seen[item.ItemID]; exists {
			continue
		}

		book, err := s.bookstore.GetBookByID(ctx, item.ItemID)
		if err != nil || book == nil {
			continue
		}

		result = append(result, toBookBrief(book))
		seen[item.ItemID] = struct{}{}
		if len(result) >= limit {
			break
		}
	}

	return result
}

func (s *DiscoveryService) getFallbackBooks(ctx context.Context, limit int) ([]BookBrief, error) {
	if s.bookstore == nil {
		return []BookBrief{}, nil
	}

	fallbackSources := []func(context.Context, int, int) ([]*bookstoreModel.Book, int64, error){
		s.bookstore.GetRecommendedBooks,
		s.bookstore.GetFeaturedBooks,
		s.bookstore.GetHotBooks,
		s.bookstore.GetNewReleases,
	}

	for _, source := range fallbackSources {
		books, _, err := source(ctx, 1, limit)
		if err != nil {
			continue
		}
		if len(books) > 0 {
			return toBookBriefs(books), nil
		}
	}

	return []BookBrief{}, nil
}

func toBookBriefs(books []*bookstoreModel.Book) []BookBrief {
	result := make([]BookBrief, 0, len(books))
	for _, book := range books {
		if book == nil {
			continue
		}
		result = append(result, toBookBrief(book))
	}
	return result
}

func toBookBrief(book *bookstoreModel.Book) BookBrief {
	return BookBrief{
		ID:          book.ID.Hex(),
		Title:       normalizeTitle(book.Title),
		Cover:       book.Cover,
		Author:      normalizeAuthor(book.Author),
		Rating:      book.Rating.ToFloat(),
		Description: normalizeDescription(book.Introduction),
		Tags:        append([]string{}, book.Tags...),
		Category:    firstNonEmpty(book.Categories...),
		PublishDate: resolvePublishDate(book),
	}
}

func normalizeTitle(title string) string {
	if strings.TrimSpace(title) == "" {
		return defaultBookTitle
	}
	return title
}

func normalizeAuthor(author string) string {
	if strings.TrimSpace(author) == "" {
		return defaultBookAuthor
	}
	return author
}

func normalizeDescription(description string) string {
	if strings.TrimSpace(description) == "" {
		return defaultBookDescription
	}
	return description
}

func resolvePublishDate(book *bookstoreModel.Book) string {
	if book == nil {
		return ""
	}
	if book.PublishedAt != nil {
		return book.PublishedAt.Format(timeLayoutRFC3339)
	}
	if book.LastUpdateAt != nil {
		return book.LastUpdateAt.Format(timeLayoutRFC3339)
	}
	return book.UpdatedAt.Format(timeLayoutRFC3339)
}

func buildEditorReason(book *bookstoreModel.Book) string {
	title := normalizeTitle(book.Title)
	switch {
	case len(book.Tags) > 0:
		return "编辑部推荐《" + title + "》，因为它在“" + book.Tags[0] + "”方向的读者反馈表现稳定。"
	case len(book.Categories) > 0:
		return "编辑部推荐《" + title + "》，因为它在" + book.Categories[0] + "题材中有不错的完成度。"
	default:
		return "编辑部推荐《" + title + "》，因为它在最近的读者反馈中表现稳定。"
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func containsTag(tags []string, target string) bool {
	normalizedTarget := strings.TrimSpace(strings.ToLower(target))
	if normalizedTarget == "" {
		return true
	}
	for _, tag := range tags {
		if strings.TrimSpace(strings.ToLower(tag)) == normalizedTarget {
			return true
		}
	}
	return false
}

func normalizeTrendingType(value string) string {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", defaultTrendingType, "realtime":
		return defaultTrendingType
	case "weekly":
		return "weekly"
	case "monthly":
		return "monthly"
	default:
		return defaultTrendingType
	}
}

func normalizeStringSlice(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func (s *DiscoveryService) toBookBriefFromPostBook(ctx context.Context, book *socialModel.BookInfo) BookBrief {
	if book == nil {
		return BookBrief{}
	}

	if s.bookstore != nil && strings.TrimSpace(book.BookID) != "" {
		if hydratedBook, err := s.bookstore.GetBookByID(ctx, book.BookID); err == nil && hydratedBook != nil {
			return toBookBrief(hydratedBook)
		}
	}

	return BookBrief{
		ID:          strings.TrimSpace(book.BookID),
		Title:       normalizeTitle(book.Title),
		Cover:       book.Cover,
		Author:      normalizeAuthor(book.Author),
		Rating:      0,
		Description: defaultBookDescription,
		Tags:        []string{},
	}
}

func buildCommunityTopicDescription(topic string, score int) string {
	return fmt.Sprintf("社区里最近关于“%s”的讨论明显升温，当前聚合热度分值为 %d。", strings.TrimSpace(topic), score)
}
