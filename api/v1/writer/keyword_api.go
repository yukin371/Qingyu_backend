package writer

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/pkg/response"
	"Qingyu_backend/service/interfaces"
)

// KeywordApi 关键词检索API
type KeywordApi struct {
	characterService interfaces.CharacterService
	locationService  interfaces.LocationService
}

// NewKeywordApi 创建关键词检索API
func NewKeywordApi(characterService interfaces.CharacterService, locationService interfaces.LocationService) *KeywordApi {
	return &KeywordApi{
		characterService: characterService,
		locationService:  locationService,
	}
}

type keywordSuggestion struct {
	Type      string `json:"type"` // character | location
	ID        string `json:"id"`
	Name      string `json:"name"`
	MatchMode string `json:"matchMode"` // exact | prefix | contains | alias
}

type keywordSearchResponse struct {
	Query       string              `json:"query"`
	Suggestions []keywordSuggestion `json:"suggestions"`
}

// SearchKeywords 按项目搜索角色/地点关键词（支持前缀补全）
func (api *KeywordApi) SearchKeywords(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		projectID = c.Param("id")
	}
	if projectID == "" {
		response.BadRequest(c, "项目ID不能为空", "")
		return
	}

	// 验证projectId格式
	if _, err := primitive.ObjectIDFromHex(projectID); err != nil {
		response.BadRequest(c, "项目ID格式无效", "")
		return
	}

	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		response.BadRequest(c, "q 不能为空", "")
		return
	}

	// 限制查询长度，防止DoS攻击
	const maxQueryLength = 50
	if len(query) > maxQueryLength {
		response.BadRequest(c, "查询字符串过长（最大50字符）", "")
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	suggestions := make([]keywordSuggestion, 0, limit)

	if api.characterService != nil {
		characters, err := api.characterService.List(c.Request.Context(), projectID)
		if err != nil {
			c.Error(err)
			return
		}
		for _, ch := range characters {
			matchMode, ok := matchKeyword(query, ch.Name, ch.Alias...)
			if !ok {
				continue
			}
			suggestions = append(suggestions, keywordSuggestion{
				Type:      "character",
				ID:        ch.ID.Hex(),
				Name:      ch.Name,
				MatchMode: matchMode,
			})
			if len(suggestions) >= limit {
				response.Success(c, &keywordSearchResponse{Query: query, Suggestions: suggestions})
				return
			}
		}
	}

	if api.locationService != nil {
		locations, err := api.locationService.List(c.Request.Context(), projectID)
		if err != nil {
			c.Error(err)
			return
		}
		for _, loc := range locations {
			matchMode, ok := matchKeyword(query, loc.Name, loc.Description, loc.Culture, loc.Climate, loc.Atmosphere, loc.Geography)
			if !ok {
				continue
			}
			suggestions = append(suggestions, keywordSuggestion{
				Type:      "location",
				ID:        loc.ID.Hex(),
				Name:      loc.Name,
				MatchMode: matchMode,
			})
			if len(suggestions) >= limit {
				break
			}
		}
	}

	response.Success(c, &keywordSearchResponse{
		Query:       query,
		Suggestions: suggestions,
	})
}

func matchKeyword(query string, primary string, candidates ...string) (string, bool) {
	q := normalizeForSearch(query)
	if q == "" {
		return "", false
	}

	for _, token := range buildSearchTokens(primary) {
		if token == q {
			return "exact", true
		}
		if strings.HasPrefix(token, q) {
			return "prefix", true
		}
		if strings.Contains(token, q) {
			return "contains", true
		}
	}

	for _, candidate := range candidates {
		for _, token := range buildSearchTokens(candidate) {
			if token == "" {
				continue
			}
			if token == q || strings.HasPrefix(token, q) || strings.Contains(token, q) {
				return "alias", true
			}
		}
	}

	return "", false
}

func normalizeForSearch(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func buildSearchTokens(v string) []string {
	n := normalizeForSearch(v)
	if n == "" {
		return nil
	}

	tokenSet := map[string]struct{}{
		n: {},
	}

	args := pinyin.NewArgs()
	py := pinyin.LazyPinyin(v, args)
	if len(py) > 0 {
		fullPinyin := normalizeForSearch(strings.Join(py, ""))
		if fullPinyin != "" {
			tokenSet[fullPinyin] = struct{}{}
		}

		var initials strings.Builder
		for _, syllable := range py {
			rs := []rune(syllable)
			if len(rs) == 0 {
				continue
			}
			initials.WriteRune(rs[0])
		}
		initialToken := normalizeForSearch(initials.String())
		if initialToken != "" {
			tokenSet[initialToken] = struct{}{}
		}
	}

	result := make([]string, 0, len(tokenSet))
	for token := range tokenSet {
		result = append(result, token)
	}
	return result
}
