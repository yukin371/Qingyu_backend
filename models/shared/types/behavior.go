package types

import (
	"fmt"
	"strings"
)

// RecommendationBehaviorType 推荐行为类型。
// 这一组用于推荐/曝光/交互事件，不与阅读统计行为强行合并。
type RecommendationBehaviorType string

const (
	RecommendationBehaviorView     RecommendationBehaviorType = "view"
	RecommendationBehaviorClick    RecommendationBehaviorType = "click"
	RecommendationBehaviorCollect  RecommendationBehaviorType = "collect"
	RecommendationBehaviorRead     RecommendationBehaviorType = "read"
	RecommendationBehaviorFinish   RecommendationBehaviorType = "finish"
	RecommendationBehaviorLike     RecommendationBehaviorType = "like"
	RecommendationBehaviorShare    RecommendationBehaviorType = "share"
	RecommendationBehaviorPurchase RecommendationBehaviorType = "purchase"
	RecommendationBehaviorComment  RecommendationBehaviorType = "comment"
	RecommendationBehaviorRate     RecommendationBehaviorType = "rate"
)

const (
	legacyRecommendationBehaviorFavorite = "favorite"
	legacyRecommendationBehaviorComplete = "complete"
)

var AllRecommendationBehaviorTypes = []RecommendationBehaviorType{
	RecommendationBehaviorView,
	RecommendationBehaviorClick,
	RecommendationBehaviorCollect,
	RecommendationBehaviorRead,
	RecommendationBehaviorFinish,
	RecommendationBehaviorLike,
	RecommendationBehaviorShare,
	RecommendationBehaviorPurchase,
	RecommendationBehaviorComment,
	RecommendationBehaviorRate,
}

func normalizeRecommendationBehaviorType(value string) RecommendationBehaviorType {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case legacyRecommendationBehaviorFavorite:
		return RecommendationBehaviorCollect
	case legacyRecommendationBehaviorComplete:
		return RecommendationBehaviorFinish
	default:
		return RecommendationBehaviorType(strings.TrimSpace(strings.ToLower(value)))
	}
}

func (t RecommendationBehaviorType) IsValid() bool {
	switch normalizeRecommendationBehaviorType(string(t)) {
	case RecommendationBehaviorView,
		RecommendationBehaviorClick,
		RecommendationBehaviorCollect,
		RecommendationBehaviorRead,
		RecommendationBehaviorFinish,
		RecommendationBehaviorLike,
		RecommendationBehaviorShare,
		RecommendationBehaviorPurchase,
		RecommendationBehaviorComment,
		RecommendationBehaviorRate:
		return true
	default:
		return false
	}
}

func (t RecommendationBehaviorType) String() string {
	return string(t)
}

func ParseRecommendationBehaviorType(value string) (RecommendationBehaviorType, error) {
	behaviorType := normalizeRecommendationBehaviorType(value)
	if !behaviorType.IsValid() {
		return "", fmt.Errorf("invalid recommendation behavior type: %s", value)
	}
	return behaviorType, nil
}

func RecommendationBehaviorQueryValues(behaviorType RecommendationBehaviorType) []string {
	switch normalizeRecommendationBehaviorType(string(behaviorType)) {
	case RecommendationBehaviorCollect:
		return []string{RecommendationBehaviorCollect.String(), legacyRecommendationBehaviorFavorite}
	case RecommendationBehaviorFinish:
		return []string{RecommendationBehaviorFinish.String(), legacyRecommendationBehaviorComplete}
	default:
		return []string{normalizeRecommendationBehaviorType(string(behaviorType)).String()}
	}
}

// ReaderBehaviorType 阅读统计行为类型。
// 这一组用于阅读漏斗/留存统计，与 recommendation 行为分开维护。
type ReaderBehaviorType string

const (
	ReaderBehaviorView      ReaderBehaviorType = "view"
	ReaderBehaviorComplete  ReaderBehaviorType = "complete"
	ReaderBehaviorDropOff   ReaderBehaviorType = "drop_off"
	ReaderBehaviorSubscribe ReaderBehaviorType = "subscribe"
	ReaderBehaviorBookmark  ReaderBehaviorType = "bookmark"
	ReaderBehaviorComment   ReaderBehaviorType = "comment"
	ReaderBehaviorLike      ReaderBehaviorType = "like"
)

var AllReaderBehaviorTypes = []ReaderBehaviorType{
	ReaderBehaviorView,
	ReaderBehaviorComplete,
	ReaderBehaviorDropOff,
	ReaderBehaviorSubscribe,
	ReaderBehaviorBookmark,
	ReaderBehaviorComment,
	ReaderBehaviorLike,
}

func normalizeReaderBehaviorType(value string) ReaderBehaviorType {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case string(RecommendationBehaviorFinish):
		return ReaderBehaviorComplete
	default:
		return ReaderBehaviorType(strings.TrimSpace(strings.ToLower(value)))
	}
}

func (t ReaderBehaviorType) IsValid() bool {
	switch normalizeReaderBehaviorType(string(t)) {
	case ReaderBehaviorView,
		ReaderBehaviorComplete,
		ReaderBehaviorDropOff,
		ReaderBehaviorSubscribe,
		ReaderBehaviorBookmark,
		ReaderBehaviorComment,
		ReaderBehaviorLike:
		return true
	default:
		return false
	}
}

func (t ReaderBehaviorType) String() string {
	return string(t)
}

func ParseReaderBehaviorType(value string) (ReaderBehaviorType, error) {
	behaviorType := normalizeReaderBehaviorType(value)
	if !behaviorType.IsValid() {
		return "", fmt.Errorf("invalid reader behavior type: %s", value)
	}
	return behaviorType, nil
}
