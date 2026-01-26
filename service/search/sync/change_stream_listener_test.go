package sync

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"Qingyu_backend/models/search"
)

// MockMongoClient 模拟 MongoDB 客户端
type MockMongoClient struct {
	events []bson.M
}

// NewMockMongoClient 创建模拟 MongoDB 客户端
func NewMockMongoClient() *MockMongoClient {
	return &MockMongoClient{
		events: []bson.M{
			{
				"_id": primitive.NewObjectID(),
				"operationType": "insert",
				"fullDocument": bson.M{
					"_id":        primitive.NewObjectID(),
					"title":      "Test Book",
					"author_id":  primitive.NewObjectID().Hex(),
					"created_at": time.Now(),
				},
				"ns": bson.M{
					"db":         "testdb",
					"collection": "books",
				},
				"documentKey": bson.M{
					"_id": primitive.NewObjectID(),
				},
				"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
			},
			{
				"_id": primitive.NewObjectID(),
				"operationType": "update",
				"updateDescription": bson.M{
					"updatedFields": bson.M{
						"title": "Updated Book",
					},
				},
				"ns": bson.M{
					"db":         "testdb",
					"collection": "books",
				},
				"documentKey": bson.M{
					"_id": primitive.NewObjectID(),
				},
				"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
			},
			{
				"_id": primitive.NewObjectID(),
				"operationType": "delete",
				"ns": bson.M{
					"db":         "testdb",
					"collection": "books",
				},
				"documentKey": bson.M{
					"_id": primitive.NewObjectID(),
				},
				"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
			},
		},
	}
}

// TestNewChangeStreamListener 测试创建监听器
func TestNewChangeStreamListener(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// 注意：由于 ChangeStreamListener 需要 *mongo.Client 和 *mongo.Database
	// 这里我们只测试基本结构，实际的集成测试需要真实的 MongoDB
	t.Run("valid structure", func(t *testing.T) {
		listener := &ChangeStreamListenerImpl{
			zapLogger:      logger,
			collections:    []string{"books", "projects", "documents", "users"},
			eventBuffer:    make([]search.SyncEvent, 0, 100),
			resumeTokenMap: make(map[string][]byte),
		}

		assert.NotNil(t, listener)
		assert.Equal(t, []string{"books", "projects", "documents", "users"}, listener.collections)
		assert.NotNil(t, listener.eventBuffer)
		assert.NotNil(t, listener.resumeTokenMap)
	})
}

// TestConvertEventToSyncEvent 测试事件转换
func TestConvertEventToSyncEvent(t *testing.T) {
	logger := zaptest.NewLogger(t)
	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	listener := &ChangeStreamListenerImpl{
		zapLogger:      logger,
		redisClient:    redisClient,
		eventBuffer:    make([]search.SyncEvent, 0, 100),
		resumeTokenMap: make(map[string][]byte),
	}

	t.Run("convert insert event", func(t *testing.T) {
		docID := primitive.NewObjectID()
		eventData := bson.M{
			"_id":           primitive.NewObjectID(),
			"operationType": "insert",
			"fullDocument": bson.M{
				"_id":        docID,
				"title":      "Test Book",
				"author_id":  primitive.NewObjectID().Hex(),
				"created_at": time.Now(),
			},
			"ns": bson.M{
				"db":         "testdb",
				"collection": "books",
			},
			"documentKey": bson.M{
				"_id": docID,
			},
			"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
		}

		eventBytes, _ := bson.Marshal(eventData)
		rawEvent := bson.Raw(eventBytes)

		syncEvent, err := listener.ConvertEventToSyncEvent("books", &rawEvent)
		require.NoError(t, err)
		assert.NotNil(t, syncEvent)
		assert.Equal(t, search.SyncEventInsert, syncEvent.Type)
		assert.Equal(t, "books", syncEvent.Index)
		assert.NotNil(t, syncEvent.FullDocument)
	})

	t.Run("convert update event", func(t *testing.T) {
		docID := primitive.NewObjectID()
		eventData := bson.M{
			"_id":           primitive.NewObjectID(),
			"operationType": "update",
			"updateDescription": bson.M{
				"updatedFields": bson.M{
					"title": "Updated Book",
				},
			},
			"ns": bson.M{
				"db":         "testdb",
				"collection": "books",
			},
			"documentKey": bson.M{
				"_id": docID,
			},
			"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
		}

		eventBytes, _ := bson.Marshal(eventData)
		rawEvent := bson.Raw(eventBytes)

		syncEvent, err := listener.ConvertEventToSyncEvent("books", &rawEvent)
		require.NoError(t, err)
		assert.NotNil(t, syncEvent)
		assert.Equal(t, search.SyncEventUpdate, syncEvent.Type)
		assert.Equal(t, "books", syncEvent.Index)
		assert.NotEmpty(t, syncEvent.ChangedFields)
	})

	t.Run("convert delete event", func(t *testing.T) {
		docID := primitive.NewObjectID()
		eventData := bson.M{
			"_id":           primitive.NewObjectID(),
			"operationType": "delete",
			"ns": bson.M{
				"db":         "testdb",
				"collection": "books",
			},
			"documentKey": bson.M{
				"_id": docID,
			},
			"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
		}

		eventBytes, _ := bson.Marshal(eventData)
		rawEvent := bson.Raw(eventBytes)

		syncEvent, err := listener.ConvertEventToSyncEvent("books", &rawEvent)
		require.NoError(t, err)
		assert.NotNil(t, syncEvent)
		assert.Equal(t, search.SyncEventDelete, syncEvent.Type)
		assert.Equal(t, "books", syncEvent.Index)
	})

	t.Run("convert replace event as update", func(t *testing.T) {
		docID := primitive.NewObjectID()
		eventData := bson.M{
			"_id":           primitive.NewObjectID(),
			"operationType": "replace",
			"fullDocument": bson.M{
				"_id":   docID,
				"title": "Replaced Book",
			},
			"ns": bson.M{
				"db":         "testdb",
				"collection": "books",
			},
			"documentKey": bson.M{
				"_id": docID,
			},
			"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
		}

		eventBytes, _ := bson.Marshal(eventData)
		rawEvent := bson.Raw(eventBytes)

		syncEvent, err := listener.ConvertEventToSyncEvent("books", &rawEvent)
		require.NoError(t, err)
		assert.NotNil(t, syncEvent)
		assert.Equal(t, search.SyncEventUpdate, syncEvent.Type)
	})

	t.Run("handle invalid event", func(t *testing.T) {
		invalidData := bson.M{
			"invalid": "data",
		}

		eventBytes, _ := bson.Marshal(invalidData)
		rawEvent := bson.Raw(eventBytes)

		_, err := listener.ConvertEventToSyncEvent("books", &rawEvent)
		assert.Error(t, err)
	})
}

// TestFlushEvents 测试刷新事件到 Redis
func TestFlushEvents(t *testing.T) {
	logger := zaptest.NewLogger(t)
	miniRedis := miniredis.RunT(t)

	t.Run("flush events to redis", func(t *testing.T) {
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})

		listener := &ChangeStreamListenerImpl{
			zapLogger:      logger,
			redisClient: redisClient,
			eventBuffer: []search.SyncEvent{
				{
					ID:     "event1",
					Type:   search.SyncEventInsert,
					Index:  "books",
					OpType: search.SyncEventInsert,
				},
				{
					ID:     "event2",
					Type:   search.SyncEventUpdate,
					Index:  "projects",
					OpType: search.SyncEventUpdate,
				},
			},
		}

		err := listener.FlushEvents()
		require.NoError(t, err)

		// 验证事件已发送到 Redis
		ctx := context.Background()
		count := redisClient.LLen(ctx, "search:sync:events").Val()
		assert.Equal(t, int64(2), count)

		// 验证缓冲区已清空
		assert.Empty(t, listener.eventBuffer)
	})

	t.Run("flush empty buffer", func(t *testing.T) {
		redisClient := redis.NewClient(&redis.Options{
			Addr: miniRedis.Addr(),
		})

		listener := &ChangeStreamListenerImpl{
			zapLogger:      logger,
			redisClient: redisClient,
			eventBuffer: []search.SyncEvent{},
		}

		err := listener.FlushEvents()
		require.NoError(t, err)
		assert.Empty(t, listener.eventBuffer)
	})
}

// TestSaveResumeToken 测试保存 resume token
func TestSaveResumeToken(t *testing.T) {
	logger := zaptest.NewLogger(t)

	t.Run("save resume token", func(t *testing.T) {
		listener := &ChangeStreamListenerImpl{
			zapLogger:      logger,
			resumeTokenMap: make(map[string][]byte),
		}

		token := []byte("test-token")
		listener.SaveResumeToken("books", token)

		savedToken, exists := listener.resumeTokenMap["books"]
		assert.True(t, exists)
		assert.Equal(t, token, savedToken)
	})

	t.Run("overwrite existing token", func(t *testing.T) {
		listener := &ChangeStreamListenerImpl{
			zapLogger:      logger,
			resumeTokenMap: make(map[string][]byte),
		}

		token1 := []byte("token1")
		token2 := []byte("token2")

		listener.SaveResumeToken("books", token1)
		listener.SaveResumeToken("books", token2)

		savedToken := listener.resumeTokenMap["books"]
		assert.Equal(t, token2, savedToken)
	})
}

// TestEventSerialization 测试事件序列化
func TestEventSerialization(t *testing.T) {
	t.Run("serialize sync event", func(t *testing.T) {
		event := search.SyncEvent{
			ID:     "test-event-1",
			Type:   search.SyncEventInsert,
			Index:  "books",
			OpType: search.SyncEventInsert,
			FullDocument: map[string]interface{}{
				"_id":       primitive.NewObjectID().Hex(),
				"title":     "Test Book",
				"author_id": "author123",
			},
			Timestamp: time.Now(),
		}

		data, err := json.Marshal(event)
		require.NoError(t, err)

		var decoded search.SyncEvent
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, event.ID, decoded.ID)
		assert.Equal(t, event.Type, decoded.Type)
		assert.Equal(t, event.Index, decoded.Index)
	})
}

// TestCollectionsConfig 测试集合配置
func TestCollectionsConfig(t *testing.T) {
	t.Run("validate required collections", func(t *testing.T) {
		requiredCollections := []string{"books", "projects", "documents", "users"}
		assert.Equal(t, 4, len(requiredCollections))
		assert.Contains(t, requiredCollections, "books")
		assert.Contains(t, requiredCollections, "projects")
		assert.Contains(t, requiredCollections, "documents")
		assert.Contains(t, requiredCollections, "users")
	})
}

// BenchmarkConvertEvent 性能测试
func BenchmarkConvertEvent(b *testing.B) {
	logger := zap.NewNop()
	miniRedis := miniredis.RunT(b)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	listener := &ChangeStreamListenerImpl{
		zapLogger:      logger,
		redisClient:    redisClient,
		eventBuffer:    make([]search.SyncEvent, 0, 100),
		resumeTokenMap: make(map[string][]byte),
	}

	docID := primitive.NewObjectID()
	eventData := bson.M{
		"_id":           primitive.NewObjectID(),
		"operationType": "insert",
		"fullDocument": bson.M{
			"_id":        docID,
			"title":      "Test Book",
			"author_id":  primitive.NewObjectID().Hex(),
			"created_at": time.Now(),
		},
		"ns": bson.M{
			"db":         "testdb",
			"collection": "books",
		},
		"documentKey": bson.M{
			"_id": docID,
		},
		"clusterTime": primitive.Timestamp{T: uint32(time.Now().Unix())},
	}

	eventBytes, _ := bson.Marshal(eventData)
	rawEvent := bson.Raw(eventBytes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = listener.ConvertEventToSyncEvent("books", &rawEvent)
	}
}
