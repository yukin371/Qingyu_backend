package service

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/service/base"
	"Qingyu_backend/service/events"

	"github.com/stretchr/testify/assert"
)

// TestSimpleEventBus_Subscribe 测试事件订阅
func TestSimpleEventBus_Subscribe(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()
	handler := events.NewWelcomeEmailHandler()

	// Act
	err := eventBus.Subscribe(events.EventTypeUserRegistered, handler)

	// Assert
	assert.NoError(t, err)
}

// TestSimpleEventBus_Publish 测试同步事件发布
func TestSimpleEventBus_Publish(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()
	handler := events.NewUserActivityLogHandler()

	err := eventBus.Subscribe(events.EventTypeUserRegistered, handler)
	assert.NoError(t, err)

	// Act
	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	err = eventBus.Publish(context.Background(), event)

	// Assert
	assert.NoError(t, err)
}

// TestSimpleEventBus_PublishAsync 测试异步事件发布
func TestSimpleEventBus_PublishAsync(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()
	handler := events.NewUserActivityLogHandler()

	err := eventBus.Subscribe(events.EventTypeUserRegistered, handler)
	assert.NoError(t, err)

	// Act
	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	err = eventBus.PublishAsync(context.Background(), event)

	// Assert
	assert.NoError(t, err)

	// 等待异步处理完成
	time.Sleep(100 * time.Millisecond)
}

// TestSimpleEventBus_MultipleHandlers 测试多个处理器
func TestSimpleEventBus_MultipleHandlers(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()

	// 注册多个处理器
	welcomeHandler := events.NewWelcomeEmailHandler()
	activityHandler := events.NewUserActivityLogHandler()
	statisticsHandler := events.NewUserStatisticsHandler()

	eventBus.Subscribe(events.EventTypeUserRegistered, welcomeHandler)
	eventBus.Subscribe(events.EventTypeUserRegistered, activityHandler)
	eventBus.Subscribe(events.EventTypeUserRegistered, statisticsHandler)

	// Act
	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	err := eventBus.Publish(context.Background(), event)

	// Assert
	assert.NoError(t, err)
}

// TestSimpleEventBus_Unsubscribe 测试取消订阅
func TestSimpleEventBus_Unsubscribe(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()
	handler := events.NewWelcomeEmailHandler()

	eventBus.Subscribe(events.EventTypeUserRegistered, handler)

	// Act
	err := eventBus.Unsubscribe(events.EventTypeUserRegistered, handler.GetHandlerName())

	// Assert
	assert.NoError(t, err)

	// 发布事件，不应该被处理
	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	err = eventBus.Publish(context.Background(), event)
	assert.NoError(t, err)
}

// TestSimpleEventBus_DifferentEventTypes 测试不同事件类型
func TestSimpleEventBus_DifferentEventTypes(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()
	activityHandler := events.NewUserActivityLogHandler()

	// 订阅多个事件类型
	eventBus.Subscribe(events.EventTypeUserRegistered, activityHandler)
	eventBus.Subscribe(events.EventTypeUserLoggedIn, activityHandler)

	// Act & Assert - 用户注册事件
	registerEvent := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	err := eventBus.Publish(context.Background(), registerEvent)
	assert.NoError(t, err)

	// Act & Assert - 用户登录事件
	loginEvent := events.NewUserLoggedInEvent("user123", "testuser")
	err = eventBus.Publish(context.Background(), loginEvent)
	assert.NoError(t, err)
}

// TestReadingEvents 测试阅读事件
func TestReadingEvents(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()
	statisticsHandler := events.NewReadingStatisticsHandler()
	recommendationHandler := events.NewRecommendationUpdateHandler()

	eventBus.Subscribe(events.EventTypeChapterRead, statisticsHandler)
	eventBus.Subscribe(events.EventTypeChapterRead, recommendationHandler)

	// Act
	event := events.NewChapterReadEvent("user123", "book456", "chapter789")
	err := eventBus.Publish(context.Background(), event)

	// Assert
	assert.NoError(t, err)
}

// TestReadingProgressEvents 测试阅读进度事件
func TestReadingProgressEvents(t *testing.T) {
	// Arrange
	eventBus := base.NewSimpleEventBus()
	statisticsHandler := events.NewReadingStatisticsHandler()

	eventBus.Subscribe(events.EventTypeReadingProgress, statisticsHandler)

	// Act
	event := events.NewReadingProgressEvent("user123", "book456", 75)
	err := eventBus.Publish(context.Background(), event)

	// Assert
	assert.NoError(t, err)
}

// TestEventData 测试事件数据
func TestEventData(t *testing.T) {
	// Act
	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")

	// Assert
	assert.Equal(t, events.EventTypeUserRegistered, event.GetEventType())
	assert.Equal(t, "UserService", event.GetSource())
	assert.NotNil(t, event.GetEventData())
	assert.NotZero(t, event.GetTimestamp())

	// 验证事件数据
	data, ok := event.GetEventData().(events.UserEventData)
	assert.True(t, ok)
	assert.Equal(t, "user123", data.UserID)
	assert.Equal(t, "testuser", data.Username)
	assert.Equal(t, "test@example.com", data.Email)
	assert.Equal(t, "registered", data.Action)
}

// BenchmarkEventPublish 性能测试 - 同步发布
func BenchmarkEventPublish(b *testing.B) {
	eventBus := base.NewSimpleEventBus()
	handler := events.NewUserActivityLogHandler()
	eventBus.Subscribe(events.EventTypeUserRegistered, handler)

	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventBus.Publish(ctx, event)
	}
}

// BenchmarkEventPublishAsync 性能测试 - 异步发布
func BenchmarkEventPublishAsync(b *testing.B) {
	eventBus := base.NewSimpleEventBus()
	handler := events.NewUserActivityLogHandler()
	eventBus.Subscribe(events.EventTypeUserRegistered, handler)

	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventBus.PublishAsync(ctx, event)
	}
}

// BenchmarkMultipleHandlers 性能测试 - 多个处理器
func BenchmarkMultipleHandlers(b *testing.B) {
	eventBus := base.NewSimpleEventBus()

	// 注册多个处理器
	eventBus.Subscribe(events.EventTypeUserRegistered, events.NewWelcomeEmailHandler())
	eventBus.Subscribe(events.EventTypeUserRegistered, events.NewUserActivityLogHandler())
	eventBus.Subscribe(events.EventTypeUserRegistered, events.NewUserStatisticsHandler())

	event := events.NewUserRegisteredEvent("user123", "testuser", "test@example.com")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventBus.Publish(ctx, event)
	}
}
