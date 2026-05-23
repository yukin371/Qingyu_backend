package events

import "go.uber.org/zap"

func logEventRuntime(level string, message string, fields map[string]interface{}) {
	logger := NewStructuredEventLogger(zap.L().Named("events"))
	logger.LogWithFields(level, message, fields)
}
