package events

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/service/base"
)

// ============ 辅助函数 - 简单的Mock EventBus ============

// NewSimpleEventBus 创建简单的事件总线用于测试
func NewSimpleEventBus() base.EventBus {
	return &simpleEventBus{
		handlers: make(map[string][]base.EventHandler),
	}
}

// simpleEventBus 简单事件总线实现
type simpleEventBus struct {
	handlers map[string][]base.EventHandler
}

func (b *simpleEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

func (b *simpleEventBus) Unsubscribe(eventType string, handlerName string) error {
	handlers := b.handlers[eventType]
	for i, h := range handlers {
		if h.GetHandlerName() == handlerName {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	return nil
}

func (b *simpleEventBus) Publish(ctx context.Context, event base.Event) error {
	for _, h := range b.handlers[event.GetEventType()] {
		if err := h.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (b *simpleEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	go func() {
		for _, h := range b.handlers[event.GetEventType()] {
			_ = h.Handle(ctx, event)
		}
	}()
	return nil
}

// ============ Mock Store ============

// MockEventStore 用于测试的Mock事件存储
type MockEventStore struct {
	mock.Mock
}

func (m *MockEventStore) Store(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventStore) StoreBatch(ctx context.Context, events []base.Event) error {
	args := m.Called(ctx, events)
	return args.Error(0)
}

func (m *MockEventStore) GetByID(ctx context.Context, eventID string) (*StoredEvent, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*StoredEvent), args.Error(1)
}

func (m *MockEventStore) GetByType(ctx context.Context, eventType string, limit, offset int64) ([]*StoredEvent, error) {
	args := m.Called(ctx, eventType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*StoredEvent), args.Error(1)
}

func (m *MockEventStore) GetBySource(ctx context.Context, source string, limit, offset int64) ([]*StoredEvent, error) {
	args := m.Called(ctx, source, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*StoredEvent), args.Error(1)
}

func (m *MockEventStore) GetByTimeRange(ctx context.Context, start, end time.Time, limit, offset int64) ([]*StoredEvent, error) {
	args := m.Called(ctx, start, end, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*StoredEvent), args.Error(1)
}

func (m *MockEventStore) GetByTypeAndTimeRange(ctx context.Context, eventType string, start, end time.Time, limit, offset int64) ([]*StoredEvent, error) {
	args := m.Called(ctx, eventType, start, end, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*StoredEvent), args.Error(1)
}

func (m *MockEventStore) Replay(ctx context.Context, handler base.EventHandler, filter EventFilter) (*ReplayResult, error) {
	args := m.Called(ctx, handler, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ReplayResult), args.Error(1)
}

func (m *MockEventStore) Cleanup(ctx context.Context, before time.Time) (int64, error) {
	args := m.Called(ctx, before)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEventStore) Count(ctx context.Context, filter EventFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEventStore) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ============ Mock EventBus ============

// MockEventBus 用于测试的Mock事件总线
type MockEventBus struct {
	mock.Mock
	handlers map[string][]base.EventHandler
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		handlers: make(map[string][]base.EventHandler),
	}
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// ============ 测试辅助函数 ============

// newTestPersistedEventBus 创建测试用的PersistedEventBus
func newTestPersistedEventBus(store EventStore, asyncStore bool) *PersistedEventBus {
	eventBus := NewSimpleEventBus()
	return NewPersistedEventBus(eventBus, store, asyncStore).(*PersistedEventBus)
}

// newTestEvent 创建测试事件
func newTestPersistedEvent(eventType string, data interface{}) base.Event {
	return &base.BaseEvent{
		EventType: eventType,
		EventData: data,
		Timestamp: time.Now(),
		Source:    "test",
	}
}

// ============ PersistedEventBus Replay 方法测试 ============

// TestPersistedEventBus_Replay_Success 测试成功重放事件
func TestPersistedEventBus_Replay_Success(t *testing.T) {
	// Arrange - 准备测试数据
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockHandler := new(MockEventHandler)

	// 创建测试事件
	testEvents := []*StoredEvent{
		{
			ID:        "evt1",
			EventType: "test.event",
			EventData: map[string]interface{}{"key": "value1"},
			Timestamp: time.Now(),
			Source:    "test",
			Processed: false,
		},
		{
			ID:        "evt2",
			EventType: "test.event",
			EventData: map[string]interface{}{"key": "value2"},
			Timestamp: time.Now(),
			Source:    "test",
			Processed: false,
		},
	}

	// 设置mock期望
	mockStore.On("Replay", ctx, mock.AnythingOfType("*events.MockEventHandler"), mock.AnythingOfType("events.EventFilter")).Return(nil).Run(func(args mock.Arguments) {
		// 模拟replay过程中调用handler
		handler := args.Get(1).(base.EventHandler)
		for _, evt := range testEvents {
			event := &base.BaseEvent{
				EventType: evt.EventType,
				EventData: evt.EventData,
				Timestamp: evt.Timestamp,
				Source:    evt.Source,
			}
			_ = handler.Handle(ctx, event)
		}
	})

	mockHandler.On("Handle", ctx, mock.AnythingOfType("*base.BaseEvent")).Return(nil).Times(len(testEvents))
	mockHandler.On("GetHandlerName").Return("test-handler")
	mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

	_ = newTestPersistedEventBus(mockStore, false)

	// Act - 执行被测试的方法
	// 注意：PersistedEventBus需要实现Replay方法
	// 这里假设方法签名为: Replay(ctx context.Context, handler base.EventHandler, filter EventFilter) error

	// Assert - 验证结果
	// 验证handler被正确调用了预期次数
	t.Skip("PersistedEventBus需要实现Replay方法")
}

// TestPersistedEventBus_Replay_DryRun 测试dry-run模式
func TestPersistedEventBus_Replay_DryRun(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockHandler := new(MockEventHandler)

	// dry-run模式下，应该只统计不调用handler
	mockStore.On("Replay", ctx, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
		// 验证dry-run标志被正确传递
		return true
	})).Return(nil)

	mockHandler.On("Handle", ctx, mock.Anything).Return(nil).Times(0) // dry-run不应该调用handler
	mockHandler.On("GetHandlerName").Return("test-handler")
	mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

	_ = newTestPersistedEventBus(mockStore, false)

	// Act
	// 执行dry-run replay
	_ = ctx

	// Assert
	// 验证handler未被调用，但返回了正确的统计信息
	t.Skip("需要实现ReplayOptions支持dry-run模式")
}

// TestPersistedEventBus_Replay_WithFilter 测试带过滤器的重放
func TestPersistedEventBus_Replay_WithFilter(t *testing.T) {
	tests := []struct {
		name      string
		filter    EventFilter
		wantCall  int
		setupMock func(*MockEventStore, *MockEventHandler)
	}{
		{
			name: "按类型过滤",
			filter: EventFilter{
				EventType: "specific.event",
			},
			wantCall: 3,
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.EventType == "specific.event"
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(3)
			},
		},
		{
			name: "按来源过滤",
			filter: EventFilter{
				Source: "specific.source",
			},
			wantCall: 2,
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.Source == "specific.source"
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(2)
			},
		},
		{
			name: "按时间范围过滤",
			filter: EventFilter{
				StartTime: timePtr(time.Now().Add(-1 * time.Hour)),
				EndTime:   timePtr(time.Now()),
			},
			wantCall: 5,
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.StartTime != nil && f.EndTime != nil
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(5)
			},
		},
		{
			name: "组合过滤",
			filter: EventFilter{
				EventType: "test.event",
				Source:    "test.source",
				Limit:     10,
			},
			wantCall: 4,
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.EventType == "test.event" &&
						f.Source == "test.source" &&
						f.Limit == 10
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(4)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockStore := new(MockEventStore)
			mockHandler := new(MockEventHandler)

			if tt.setupMock != nil {
				tt.setupMock(mockStore, mockHandler)
			}

			mockHandler.On("GetHandlerName").Return("test-handler")
			mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

			// Act
			// 执行带过滤器的replay
			_ = mockStore

			// Assert
			// 验证过滤器被正确传递给store
			// 验证handler被调用正确次数
			t.Skip("PersistedEventBus需要实现Replay方法")
		})
	}
}

// TestPersistedEventBus_Replay_HandlerError 测试handler返回错误
func TestPersistedEventBus_Replay_HandlerError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockHandler := new(MockEventHandler)
	expectedErr := errors.New("handler processing error")

	// 设置mock - handler返回错误
	mockStore.On("Replay", ctx, mock.Anything, mock.Anything).Return(expectedErr).Run(func(args mock.Arguments) {
		handler := args.Get(1).(base.EventHandler)
		event := newTestPersistedEvent("test.event", nil)
		_ = handler.Handle(ctx, event) // 这个调用会返回错误
	})

	mockHandler.On("Handle", ctx, mock.Anything).Return(expectedErr).Once()
	mockHandler.On("GetHandlerName").Return("test-handler")
	mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

	// Act
	// 执行replay
	_ = mockStore

	// Assert
	// 验证错误被正确返回
	t.Skip("PersistedEventBus需要实现Replay方法")
}

// TestPersistedEventBus_Replay_EmptyResult 测试空结果
func TestPersistedEventBus_Replay_EmptyResult(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockHandler := new(MockEventHandler)

	mockStore.On("Replay", ctx, mock.Anything, mock.Anything).Return(nil)

	// handler不应该被调用
	mockHandler.On("Handle", ctx, mock.Anything).Return(nil).Times(0)
	mockHandler.On("GetHandlerName").Return("test-handler")
	mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

	// Act
	// 执行replay，没有匹配的事件
	_ = mockStore

	// Assert
	// 验证handler未被调用
	t.Skip("PersistedEventBus需要实现Replay方法")
}

// TestPersistedEventBus_Replay_Limit 测试限制数量
func TestPersistedEventBus_Replay_Limit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int64
		total    int
		wantCall int
	}{
		{
			name:     "Limit小于总数",
			limit:    5,
			total:    100,
			wantCall: 5,
		},
		{
			name:     "Limit大于总数",
			limit:    1000,
			total:    100,
			wantCall: 100,
		},
		{
			name:     "Limit为0(无限制)",
			limit:    0,
			total:    100,
			wantCall: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockStore := new(MockEventStore)
			mockHandler := new(MockEventHandler)

			mockStore.On("Replay", ctx, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
				return f.Limit == tt.limit
			})).Return(nil)

			mockHandler.On("Handle", ctx, mock.Anything).Return(nil).Times(tt.wantCall)
			mockHandler.On("GetHandlerName").Return("test-handler")
			mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

			// Act
			// 执行replay
			_ = ctx

			// Assert
			// 验证只处理了limit数量的事件
			t.Skip("PersistedEventBus需要实现Replay方法")
		})
	}
}

// TestPersistedEventBus_Replay_ContextCancellation 测试上下文取消
func TestPersistedEventBus_Replay_ContextCancellation(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	mockStore := new(MockEventStore)
	mockHandler := new(MockEventHandler)

	mockStore.On("Replay", ctx, mock.Anything, mock.Anything).Return(context.Canceled)

	mockHandler.On("Handle", ctx, mock.Anything).Return(nil)
	mockHandler.On("GetHandlerName").Return("test-handler")
	mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

	// Act
	// 在处理过程中取消上下文
	_ = mockStore
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	// Assert
	// 验证上下文取消时replay正确返回错误
	t.Skip("PersistedEventBus需要实现Replay方法")
}

// ============ PersistedEventBus Publish 方法测试 ============

// TestPersistedEventBus_Publish_StoresEvent 测试发布事件时存储
func TestPersistedEventBus_Publish_StoresEvent(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockEventBus := NewMockEventBus()

	event := newTestPersistedEvent("test.event", map[string]interface{}{"key": "value"})

	mockStore.On("Store", ctx, mock.AnythingOfType("*base.BaseEvent")).Return(nil)
	mockEventBus.On("Publish", ctx, mock.AnythingOfType("*base.BaseEvent")).Return(nil)

	bus := &PersistedEventBus{
		eventBus:   mockEventBus,
		store:      mockStore,
		asyncStore: false,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000),
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	assert.NoError(t, err, "发布事件应该成功")
	mockStore.AssertCalled(t, "Store", ctx, mock.AnythingOfType("*base.BaseEvent"))
	mockEventBus.AssertCalled(t, "Publish", ctx, mock.AnythingOfType("*base.BaseEvent"))
}

// TestPersistedEventBus_Publish_StoreErrorDoesNotAffectPublish 测试存储失败不影响发布
func TestPersistedEventBus_Publish_StoreErrorDoesNotAffectPublish(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockEventBus := NewMockEventBus()

	event := newTestPersistedEvent("test.event", nil)
	storeErr := errors.New("storage error")

	mockStore.On("Store", ctx, mock.Anything).Return(storeErr)
	mockEventBus.On("Publish", ctx, mock.Anything).Return(nil)

	bus := &PersistedEventBus{
		eventBus:   mockEventBus,
		store:      mockStore,
		asyncStore: false,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000),
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	// 存储失败不应该影响事件发布
	assert.NoError(t, err, "存储失败时事件发布仍应成功")
	mockStore.AssertCalled(t, "Store", ctx, mock.Anything)
	mockEventBus.AssertCalled(t, "Publish", ctx, mock.Anything)
}

// TestPersistedEventBus_PublishAsync_StoresEvent 测试异步发布事件时存储
func TestPersistedEventBus_PublishAsync_StoresEvent(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockEventBus := NewMockEventBus()

	event := newTestPersistedEvent("test.event", nil)

	mockStore.On("Store", ctx, mock.Anything).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	bus := &PersistedEventBus{
		eventBus:   mockEventBus,
		store:      mockStore,
		asyncStore: false,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000),
	}

	// Act
	err := bus.PublishAsync(ctx, event)

	// Assert
	assert.NoError(t, err, "异步发布事件应该成功")
	mockStore.AssertCalled(t, "Store", ctx, mock.Anything)
	mockEventBus.AssertCalled(t, "PublishAsync", ctx, mock.Anything)
}

// TestPersistedEventBus_PublishAsync_WithAsyncStore 测试异步存储模式
func TestPersistedEventBus_PublishAsync_WithAsyncStore(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockEventBus := NewMockEventBus()

	event := newTestPersistedEvent("test.event", nil)

	// 异步模式下，事件先放入队列
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	bus := &PersistedEventBus{
		eventBus:   mockEventBus,
		store:      mockStore,
		asyncStore: true,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000),
	}

	// Act
	err := bus.PublishAsync(ctx, event)

	// Assert
	assert.NoError(t, err, "异步发布应该成功")

	// 验证事件被放入队列
	select {
	case <-bus.eventCh:
		// 事件成功放入队列
	case <-time.After(100 * time.Millisecond):
		t.Error("事件应该被放入异步队列")
	}
}

// ============ PersistedEventBus Subscribe/Unsubscribe 测试 ============

// TestPersistedEventBus_Subscribe 测试订阅事件
func TestPersistedEventBus_Subscribe(t *testing.T) {
	// Arrange
	mockStore := new(MockEventStore)
	mockEventBus := NewMockEventBus()
	mockHandler := new(MockEventHandler)

	mockEventBus.On("Subscribe", "test.event", mock.AnythingOfType("*events.MockEventHandler")).Return(nil)
	mockHandler.On("GetHandlerName").Return("test-handler")
	mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

	bus := &PersistedEventBus{
		eventBus:   mockEventBus,
		store:      mockStore,
		asyncStore: false,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000),
	}

	// Act
	err := bus.Subscribe("test.event", mockHandler)

	// Assert
	assert.NoError(t, err, "订阅应该成功")
	mockEventBus.AssertCalled(t, "Subscribe", "test.event", mockHandler)
}

// TestPersistedEventBus_Unsubscribe 测试取消订阅
func TestPersistedEventBus_Unsubscribe(t *testing.T) {
	// Arrange
	mockStore := new(MockEventStore)
	mockEventBus := NewMockEventBus()

	mockEventBus.On("Unsubscribe", "test.event", "test-handler").Return(nil)

	bus := &PersistedEventBus{
		eventBus:   mockEventBus,
		store:      mockStore,
		asyncStore: false,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000),
	}

	// Act
	err := bus.Unsubscribe("test.event", "test-handler")

	// Assert
	assert.NoError(t, err, "取消订阅应该成功")
	mockEventBus.AssertCalled(t, "Unsubscribe", "test.event", "test-handler")
}

// ============ PersistedEventBus Close 测试 ============

// TestPersistedEventBus_Close 测试关闭事件总线
func TestPersistedEventBus_Close(t *testing.T) {
	// Arrange
	mockStore := new(MockEventStore)
	mockEventBus := NewMockEventBus()

	bus := &PersistedEventBus{
		eventBus:   mockEventBus,
		store:      mockStore,
		asyncStore: true,
		stopCh:     make(chan struct{}),
		eventCh:    make(chan *base.Event, 1000),
		wg:         sync.WaitGroup{}, // 需要添加sync.WaitGroup初始化
	}

	// 启动异步处理器
	bus.startAsyncProcessor()

	// Act
	err := bus.Close()

	// Assert
	assert.NoError(t, err, "关闭应该成功")

	// 验证stopCh被关闭
	select {
	case <-bus.stopCh:
		// stopCh已关闭
	default:
		t.Error("stopCh应该被关闭")
	}
}

// ============ 辅助函数 ============

// timePtr 返回time.Time的指针
func timePtr(t time.Time) *time.Time {
	return &t
}

// ============ 表格驱动测试 ============

// TestPersistedEventBus_Replay_Combinations 综合测试各种组合
func TestPersistedEventBus_Replay_Combinations(t *testing.T) {
	tests := []struct {
		name      string
		filter    EventFilter
		setupMock func(*MockEventStore, *MockEventHandler)
		wantError bool
		errorType error
	}{
		{
			name:   "无过滤条件",
			filter: EventFilter{},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(10)
			},
			wantError: false,
		},
		{
			name: "类型过滤",
			filter: EventFilter{
				EventType: "test.event",
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.EventType == "test.event"
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(5)
			},
			wantError: false,
		},
		{
			name: "来源过滤",
			filter: EventFilter{
				Source: "test.source",
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.Source == "test.source"
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(3)
			},
			wantError: false,
		},
		{
			name: "时间范围过滤",
			filter: EventFilter{
				StartTime: timePtr(time.Now().Add(-24 * time.Hour)),
				EndTime:   timePtr(time.Now()),
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.StartTime != nil && f.EndTime != nil
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(7)
			},
			wantError: false,
		},
		{
			name: "限制数量",
			filter: EventFilter{
				Limit: 10,
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.Limit == 10
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(10)
			},
			wantError: false,
		},
		{
			name: "偏移量",
			filter: EventFilter{
				Offset: 5,
				Limit:  10,
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.Offset == 5 && f.Limit == 10
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(10)
			},
			wantError: false,
		},
		{
			name: "组合过滤",
			filter: EventFilter{
				EventType: "test.event",
				Source:    "test.source",
				Limit:     10,
				Offset:    0,
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.MatchedBy(func(f EventFilter) bool {
					return f.EventType == "test.event" &&
						f.Source == "test.source" &&
						f.Limit == 10 &&
						f.Offset == 0
				})).Return(nil)
				handler.On("Handle", mock.Anything, mock.Anything).Return(nil).Times(10)
			},
			wantError: false,
		},
		{
			name: "Store返回错误",
			filter: EventFilter{
				EventType: "test.event",
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				store.On("Replay", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("store error"))
			},
			wantError: true,
			errorType: errors.New("store error"),
		},
		{
			name: "Handler返回错误",
			filter: EventFilter{
				EventType: "test.event",
			},
			setupMock: func(store *MockEventStore, handler *MockEventHandler) {
				handlerErr := errors.New("handler error")
				store.On("Replay", mock.Anything, mock.Anything, mock.Anything).Return(handlerErr).Run(func(args mock.Arguments) {
					h := args.Get(1).(base.EventHandler)
					event := newTestPersistedEvent("test.event", nil)
					_ = h.Handle(context.Background(), event)
				})
				handler.On("Handle", mock.Anything, mock.Anything).Return(handlerErr).Once()
			},
			wantError: true,
			errorType: errors.New("handler error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockStore := new(MockEventStore)
			mockHandler := new(MockEventHandler)

			if tt.setupMock != nil {
				tt.setupMock(mockStore, mockHandler)
			}

			mockHandler.On("GetHandlerName").Return("test-handler")
			mockHandler.On("GetSupportedEventTypes").Return([]string{"test.event"})

			_ = newTestPersistedEventBus(mockStore, false)

			// Act
			// err := bus.Replay(ctx, mockHandler, tt.filter)

			// Assert
			if tt.wantError {
				// assert.Error(t, err)
				// assert.Contains(t, err.Error(), tt.errorType.Error())
			} else {
				// assert.NoError(t, err)
			}
			t.Skip("PersistedEventBus需要实现Replay方法")
		})
	}
}
