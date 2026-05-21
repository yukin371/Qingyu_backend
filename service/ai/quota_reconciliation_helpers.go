package ai

import (
	"fmt"
	"math"
	"strings"
	"time"

	"Qingyu_backend/config"
	aiModels "Qingyu_backend/models/ai"
)

const (
	quotaConsistencyScopeUser     = "user"
	quotaConsistencyScopeWorkflow = "workflow"
	quotaConsistencyScopeGlobal   = "global"
)

func resolveQuotaReconciliationWindow(timeRange string, now time.Time) (string, time.Time, time.Time, error) {
	normalized := strings.ToLower(strings.TrimSpace(timeRange))
	if normalized == "" {
		normalized = "day"
	}

	location := now.Location()
	switch normalized {
	case "day":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
		return normalized, start, now, nil
	case "week":
		start := now.AddDate(0, 0, -7)
		return normalized, start, now, nil
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
		return normalized, start, now, nil
	case "all":
		return normalized, time.Unix(0, 0), now, nil
	default:
		return "", time.Time{}, time.Time{}, fmt.Errorf("不支持的时间范围: %s", timeRange)
	}
}

func normalizeQuotaSummaryGroupBy(groupBy string) string {
	normalized := strings.ToLower(strings.TrimSpace(groupBy))
	if normalized == "workflow" {
		return "workflow"
	}
	return "user"
}

func normalizeQuotaSummaryPagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func determineConsistencyAlertLevel(backendTokens, aiTokens int) (aiModels.QuotaAlertLevel, bool) {
	return determineConsistencyAlertLevelForScope(quotaConsistencyScopeUser, backendTokens, aiTokens)
}

func determineConsistencyAlertLevelForScope(scope string, backendTokens, aiTokens int) (aiModels.QuotaAlertLevel, bool) {
	return determineConsistencyAlertLevelForScopeInt64(scope, int64(backendTokens), int64(aiTokens))
}

func determineConsistencyAlertLevelInt64(backendTokens, aiTokens int64) (aiModels.QuotaAlertLevel, bool) {
	return determineConsistencyAlertLevelForScopeInt64(quotaConsistencyScopeUser, backendTokens, aiTokens)
}

func determineConsistencyAlertLevelForScopeInt64(scope string, backendTokens, aiTokens int64) (aiModels.QuotaAlertLevel, bool) {
	diff := absInt64(backendTokens - aiTokens)
	maxTokens := backendTokens
	if aiTokens > maxTokens {
		maxTokens = aiTokens
	}
	if maxTokens <= 0 {
		return aiModels.QuotaAlertLevelInfo, false
	}

	diffRatio := float64(diff) / float64(maxTokens)
	threshold := getQuotaConsistencyThreshold(scope)
	if diff >= int64(threshold.CriticalTokens) || diffRatio > threshold.CriticalRatio {
		return aiModels.QuotaAlertLevelCritical, true
	}
	if diff >= int64(threshold.WarningTokens) || diffRatio > threshold.WarningRatio {
		return aiModels.QuotaAlertLevelWarning, true
	}
	return aiModels.QuotaAlertLevelInfo, false
}

func getQuotaConsistencyThreshold(scope string) *config.QuotaConsistencyThresholdConfig {
	if config.GlobalConfig == nil || config.GlobalConfig.AIQuota == nil {
		return (&config.AIQuotaConfig{}).GetConsistencyThreshold(scope)
	}
	return config.GlobalConfig.AIQuota.GetConsistencyThreshold(scope)
}

func calculateDifferenceRatioInt64(backendTokens, aiTokens int64) float64 {
	maxTokens := backendTokens
	if aiTokens > maxTokens {
		maxTokens = aiTokens
	}
	if maxTokens <= 0 {
		return 0
	}
	return float64(absInt64(backendTokens-aiTokens)) / float64(maxTokens)
}

func absInt(value int) int {
	return int(math.Abs(float64(value)))
}

func absInt64(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

func buildConsistencyAlertKey(userID string, data map[string]interface{}) string {
	scope := strings.ToLower(strings.TrimSpace(stringifyConsistencyAlertKeyPart(data["scope"])))
	timeRange := strings.ToLower(strings.TrimSpace(stringifyConsistencyAlertKeyPart(data["timeRange"])))
	if timeRange == "" {
		timeRange = "day"
	}

	switch scope {
	case quotaConsistencyScopeWorkflow:
		groupBy := normalizeQuotaSummaryGroupBy(stringifyConsistencyAlertKeyPart(data["groupBy"]))
		groupKey := strings.TrimSpace(stringifyConsistencyAlertKeyPart(data["groupKey"]))
		return fmt.Sprintf("%s|%s|%s|%s", scope, timeRange, groupBy, groupKey)
	case quotaConsistencyScopeGlobal:
		groupBy := normalizeQuotaSummaryGroupBy(stringifyConsistencyAlertKeyPart(data["groupBy"]))
		return fmt.Sprintf("%s|%s|%s", scope, timeRange, groupBy)
	default:
		if userID == "" {
			userID = strings.TrimSpace(stringifyConsistencyAlertKeyPart(data["userId"]))
		}
		return fmt.Sprintf("%s|%s|%s", quotaConsistencyScopeUser, timeRange, userID)
	}
}

func stringifyConsistencyAlertKeyPart(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	default:
		return fmt.Sprint(value)
	}
}
