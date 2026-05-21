package discovery

import (
	"errors"
	"strconv"

	discoveryService "Qingyu_backend/service/discovery"

	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
)

var errDiscoveryServiceNotInitialized = errors.New("discovery service not initialized")

type DiscoveryBook struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Cover       string   `json:"cover"`
	Description string   `json:"description"`
	Rating      float64  `json:"rating"`
	Tags        []string `json:"tags,omitempty"`
	Category    string   `json:"category,omitempty"`
	PublishDate string   `json:"publishDate,omitempty"`
}

type DiscoveryBookListResponse struct {
	List  []DiscoveryBook `json:"list"`
	Total int64           `json:"total"`
}

type DiscoveryEditorPick struct {
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

type DiscoveryEditorPickListResponse struct {
	List  []DiscoveryEditorPick `json:"list"`
	Total int64                 `json:"total"`
}

type DiscoveryTrendingBook struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Cover       string   `json:"cover"`
	Description string   `json:"description,omitempty"`
	Rating      float64  `json:"rating"`
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

type DiscoveryRecommendationItem struct {
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
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Books       []DiscoveryBook `json:"books"`
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

type PersonalizedRecommendations struct {
	Books struct {
		ForYou   []DiscoveryBook `json:"forYou"`
		Similar  []DiscoveryBook `json:"similar"`
		Trending []DiscoveryBook `json:"trending"`
	} `json:"books"`
	Booklists struct {
		Recommended []map[string]any `json:"recommended"`
		Popular     []map[string]any `json:"popular"`
	} `json:"booklists"`
	Authors struct {
		Suggested []map[string]any `json:"suggested"`
	} `json:"authors"`
}

type DiscoveryAPI struct {
	service *discoveryService.DiscoveryService
}

func NewDiscoveryAPI(service *discoveryService.DiscoveryService) *DiscoveryAPI {
	return &DiscoveryAPI{service: service}
}

func (api *DiscoveryAPI) GetPersonalized(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	limit := normalizeLimit(c.DefaultQuery("limit", "10"), 10)
	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)

	result, err := api.service.GetPersonalized(c.Request.Context(), userIDStr, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	payload := PersonalizedRecommendations{}
	payload.Books.ForYou = toDiscoveryBooks(result.Books.ForYou)
	payload.Books.Similar = toDiscoveryBooks(result.Books.Similar)
	payload.Books.Trending = toDiscoveryBooks(result.Books.Trending)
	payload.Booklists.Recommended = []map[string]any{}
	payload.Booklists.Popular = []map[string]any{}
	payload.Authors.Suggested = []map[string]any{}

	response.SuccessWithMessage(c, "获取发现推荐成功", payload)
}

func (api *DiscoveryAPI) GetNewReleases(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	page := normalizePage(c.DefaultQuery("page", "1"))
	size := normalizeLimit(c.DefaultQuery("size", "20"), 20)
	categoryID := c.Query("categoryId")

	result, err := api.service.GetNewReleases(c.Request.Context(), page, size, categoryID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "获取新书上架成功", DiscoveryBookListResponse{
		List:  toDiscoveryBooks(result.List),
		Total: result.Total,
	})
}

func (api *DiscoveryAPI) GetEditorsPick(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	page := normalizePage(c.DefaultQuery("page", "1"))
	size := normalizeLimit(c.DefaultQuery("size", "20"), 20)
	tag := c.Query("tag")

	result, err := api.service.GetEditorsPick(c.Request.Context(), page, size, tag)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	list := make([]DiscoveryEditorPick, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, DiscoveryEditorPick{
			ID:           item.ID,
			Title:        item.Title,
			Author:       item.Author,
			Cover:        item.Cover,
			Description:  item.Description,
			Reason:       item.Reason,
			ViewCount:    item.ViewCount,
			CollectCount: item.CollectCount,
			Rating:       item.Rating,
		})
	}

	response.SuccessWithMessage(c, "获取编辑推荐成功", DiscoveryEditorPickListResponse{
		List:  list,
		Total: result.Total,
	})
}

func (api *DiscoveryAPI) GetTrending(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	trendingType := c.DefaultQuery("type", "daily")
	limit := normalizeLimit(c.DefaultQuery("limit", "10"), 10)

	result, err := api.service.GetTrending(c.Request.Context(), trendingType, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	list := make([]DiscoveryTrendingBook, 0, len(result))
	for _, item := range result {
		list = append(list, DiscoveryTrendingBook{
			ID:          item.ID,
			Title:       item.Title,
			Author:      item.Author,
			Cover:       item.Cover,
			Description: item.Description,
			Rating:      item.Rating,
			Tags:        append([]string{}, item.Tags...),
			Category:    item.Category,
			PublishDate: item.PublishDate,
			Rank:        item.Rank,
			RankingType: item.RankingType,
			Period:      item.Period,
			Score:       item.Score,
			ViewCount:   item.ViewCount,
			LikeCount:   item.LikeCount,
		})
	}

	response.SuccessWithMessage(c, "获取热门榜单成功", list)
}

func (api *DiscoveryAPI) TrackAction(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	var req struct {
		ItemID   string                 `json:"itemId" binding:"required"`
		ItemType string                 `json:"itemType"`
		Action   string                 `json:"action" binding:"required"`
		Metadata map[string]interface{} `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)

	err := api.service.TrackAction(c.Request.Context(), userIDStr, discoveryService.TrackActionRequest{
		ItemID:   req.ItemID,
		ItemType: req.ItemType,
		Action:   req.Action,
		Metadata: req.Metadata,
	})
	if err != nil {
		if errors.Is(err, discoveryService.ErrDiscoveryTrackItemIDRequired) ||
			errors.Is(err, discoveryService.ErrDiscoveryTrackActionRequired) ||
			errors.Is(err, discoveryService.ErrDiscoveryTrackActionUnsupported) {
			response.BadRequest(c, "参数错误", err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "记录发现页行为成功", nil)
}

func (api *DiscoveryAPI) GetRecommendations(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	slot := c.Query("slot")
	result, err := api.service.GetRecommendations(c.Request.Context(), slot)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	list := make([]DiscoveryRecommendationItem, 0, len(result))
	for _, item := range result {
		list = append(list, DiscoveryRecommendationItem{
			ID:          item.ID,
			Type:        item.Type,
			Slot:        item.Slot,
			Title:       item.Title,
			Description: item.Description,
			Cover:       item.Cover,
			Link:        item.Link,
			Priority:    item.Priority,
			StartTime:   item.StartTime,
			EndTime:     item.EndTime,
		})
	}

	response.SuccessWithMessage(c, "获取发现推荐位成功", list)
}

func (api *DiscoveryAPI) GetTopics(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	page := normalizePage(c.DefaultQuery("page", "1"))
	size := normalizeLimit(c.DefaultQuery("size", "10"), 10)

	result, err := api.service.GetTopics(c.Request.Context(), page, size)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	list := make([]DiscoveryTopic, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, DiscoveryTopic{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			Books:       toDiscoveryBooks(item.Books),
		})
	}

	response.SuccessWithMessage(c, "获取发现主题成功", DiscoveryTopicListResponse{
		List:  list,
		Total: result.Total,
	})
}

func (api *DiscoveryAPI) UpdatePreferences(c *gin.Context) {
	if api.service == nil {
		response.InternalError(c, errDiscoveryServiceNotInitialized)
		return
	}

	var req RecommendationConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)
	result, err := api.service.UpdatePreferences(c.Request.Context(), userIDStr, discoveryService.RecommendationConfig{
		UserID: req.UserID,
		Preferences: discoveryService.RecommendationPreferences{
			FavoriteGenres:  req.Preferences.FavoriteGenres,
			FavoriteAuthors: req.Preferences.FavoriteAuthors,
			FavoriteTags:    req.Preferences.FavoriteTags,
		},
		History: discoveryService.RecommendationHistory{
			ViewedBooks: req.History.ViewedBooks,
			ViewedLists: req.History.ViewedLists,
		},
	})
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "更新发现偏好成功", RecommendationConfig{
		UserID: result.UserID,
		Preferences: RecommendationPreferences{
			FavoriteGenres:  result.Preferences.FavoriteGenres,
			FavoriteAuthors: result.Preferences.FavoriteAuthors,
			FavoriteTags:    result.Preferences.FavoriteTags,
		},
		History: RecommendationHistory{
			ViewedBooks: result.History.ViewedBooks,
			ViewedLists: result.History.ViewedLists,
		},
	})
}

func normalizePage(pageStr string) int {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func normalizeLimit(limitStr string, fallback int) int {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return fallback
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func toDiscoveryBooks(books []discoveryService.BookBrief) []DiscoveryBook {
	result := make([]DiscoveryBook, 0, len(books))
	for _, book := range books {
		result = append(result, DiscoveryBook{
			ID:          book.ID,
			Title:       book.Title,
			Author:      book.Author,
			Cover:       book.Cover,
			Description: book.Description,
			Rating:      book.Rating,
			Tags:        append([]string{}, book.Tags...),
			Category:    book.Category,
			PublishDate: book.PublishDate,
		})
	}
	return result
}
