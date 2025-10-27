package collection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/models/reader"
	"Qingyu_backend/service/base"
	"Qingyu_backend/service/reading"
)

// MockCollectionRepository Mock收藏Repository
type MockCollectionRepository struct {
	mock.Mock
}

func (m *MockCollectionRepository) Create(ctx context.Context, collection *reader.Collection) error {
	args := m.Called(ctx, collection)
	return args.Error(0)
}

func (m *MockCollectionRepository) GetByID(ctx context.Context, id string) (*reader.Collection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Collection), args.Error(1)
}

func (m *MockCollectionRepository) GetByUserAndTarget(ctx context.Context, userID, targetID string) (*reader.Collection, error) {
	args := m.Called(ctx, userID, targetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Collection), args.Error(1)
}

func (m *MockCollectionRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*reader.Collection, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.Collection), args.Error(1)
}

func (m *MockCollectionRepository) GetCollectionsByUser(ctx context.Context, userID string, folderID string, page, size int) ([]*reader.Collection, int64, error) {
	args := m.Called(ctx, userID, folderID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*reader.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockCollectionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCollectionRepository) GetCollectionsByTag(ctx context.Context, userID, tag string, page, size int) ([]*reader.Collection, int64, error) {
	args := m.Called(ctx, userID, tag, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*reader.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) CountUserCollections(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCollectionRepository) GetPublicCollections(ctx context.Context, page, size int) ([]*reader.Collection, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*reader.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) CreateFolder(ctx context.Context, folder *reader.CollectionFolder) error {
	args := m.Called(ctx, folder)
	return args.Error(0)
}

func (m *MockCollectionRepository) GetFolderByID(ctx context.Context, id string) (*reader.CollectionFolder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.CollectionFolder), args.Error(1)
}

func (m *MockCollectionRepository) GetFoldersByUser(ctx context.Context, userID string) ([]*reader.CollectionFolder, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader.CollectionFolder), args.Error(1)
}

func (m *MockCollectionRepository) UpdateFolder(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockCollectionRepository) DeleteFolder(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCollectionRepository) IncrementFolderBookCount(ctx context.Context, folderID string) error {
	args := m.Called(ctx, folderID)
	return args.Error(0)
}

func (m *MockCollectionRepository) DecrementFolderBookCount(ctx context.Context, folderID string) error {
	args := m.Called(ctx, folderID)
	return args.Error(0)
}

func (m *MockCollectionRepository) GetPublicFolders(ctx context.Context, page, size int) ([]*reader.CollectionFolder, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*reader.CollectionFolder), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	events []base.Event
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		events: make([]base.Event, 0),
	}
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	return nil
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	return nil
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	m.events = append(m.events, event)
	return nil
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	m.events = append(m.events, event)
	return nil
}

func (m *MockEventBus) GetServiceName() string {
	return "MockEventBus"
}

func (m *MockEventBus) GetVersion() string {
	return "1.0.0"
}

func (m *MockEventBus) Initialize(ctx context.Context) error {
	return nil
}

func (m *MockEventBus) Health(ctx context.Context) error {
	return nil
}

func (m *MockEventBus) Close(ctx context.Context) error {
	return nil
}

// TestCollectionService_AddCollection 添加收藏测试
func TestCollectionService_AddCollection(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("AddCollection_Success", func(t *testing.T) {
		// Mock检查是否已收藏
		mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).
			Return(nil, nil).Once()

		// Mock创建收藏
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Collection")).
			Return(nil).Once()

		// 添加收藏
		collection, err := service.AddToCollection(ctx, testUserID, testBookID, "", "很好的书", []string{"玄幻"}, false)

		assert.NoError(t, err)
		assert.NotNil(t, collection)
		assert.Equal(t, testUserID, collection.UserID)
		assert.Equal(t, testBookID, collection.BookID)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 添加收藏成功")
	})

	t.Run("AddCollection_AlreadyCollected", func(t *testing.T) {
		// Mock检查是否已收藏（已收藏）
		existingCollection := &reader.Collection{
			ID:     primitive.NewObjectID(),
			UserID: testUserID,
			BookID: testBookID,
		}
		mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).
			Return(existingCollection, nil).Once()

		// 添加收藏
		_, err := service.AddToCollection(ctx, testUserID, testBookID, "", "笔记", []string{}, false)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已经收藏")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 重复收藏检测通过")
	})

	t.Run("AddCollection_EmptyBookID", func(t *testing.T) {
		_, err := service.AddToCollection(ctx, testUserID, "", "", "", []string{}, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "书籍ID")

		t.Logf("✓ 空书籍ID验证通过")
	})
}

// TestCollectionService_AddCollectionToFolder 添加到收藏夹测试
func TestCollectionService_AddCollectionToFolder(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("AddCollectionToFolder_Success", func(t *testing.T) {
		// Mock检查是否已收藏
		mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).
			Return(nil, nil).Once()

		// Mock创建收藏
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reader.Collection")).
			Return(nil).Once()

		// Mock增加收藏夹计数
		mockRepo.On("IncrementFolderBookCount", ctx, testFolderID).
			Return(nil).Once()

		// 添加到收藏夹
		collection, err := service.AddToCollection(ctx, testUserID, testBookID, testFolderID, "", []string{}, false)

		assert.NoError(t, err)
		assert.NotNil(t, collection)
		assert.Equal(t, testFolderID, collection.FolderID)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 添加到收藏夹成功")
	})
}

// TestCollectionService_RemoveCollection 删除收藏测试
func TestCollectionService_RemoveCollection(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCollectionID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("RemoveCollection_Success", func(t *testing.T) {
		// Mock获取收藏
		collection := &reader.Collection{
			ID:       primitive.NewObjectID(),
			UserID:   testUserID,
			FolderID: testFolderID,
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// Mock减少收藏夹计数
		mockRepo.On("DecrementFolderBookCount", ctx, testFolderID).
			Return(nil).Once()

		// Mock删除收藏
		mockRepo.On("Delete", ctx, testCollectionID).
			Return(nil).Once()

		// 删除收藏
		err := service.RemoveFromCollection(ctx, testUserID, testCollectionID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除收藏成功")
	})

	t.Run("RemoveCollection_NotOwner", func(t *testing.T) {
		// Mock获取收藏（不是自己的）
		collection := &reader.Collection{
			ID:     primitive.NewObjectID(),
			UserID: primitive.NewObjectID().Hex(),
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// 删除收藏
		err := service.RemoveFromCollection(ctx, testUserID, testCollectionID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无权删除")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除权限检查通过")
	})
}

// TestCollectionService_UpdateCollection 更新收藏测试
func TestCollectionService_UpdateCollection(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCollectionID := primitive.NewObjectID().Hex()

	t.Run("UpdateCollection_Success", func(t *testing.T) {
		// Mock获取收藏
		collection := &reader.Collection{
			ID:     primitive.NewObjectID(),
			UserID: testUserID,
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// Mock更新
		mockRepo.On("Update", ctx, testCollectionID, mock.AnythingOfType("map[string]interface {}")).
			Return(nil).Once()

		// 更新收藏
		updates := map[string]interface{}{
			"notes": "新笔记",
			"tags":  []string{"新标签"},
		}
		err := service.UpdateCollection(ctx, testUserID, testCollectionID, updates)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新收藏成功")
	})

	t.Run("UpdateCollection_NotOwner", func(t *testing.T) {
		// Mock获取收藏（不是自己的）
		collection := &reader.Collection{
			ID:     primitive.NewObjectID(),
			UserID: primitive.NewObjectID().Hex(),
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// 更新收藏
		updates := map[string]interface{}{"notes": "新笔记"}
		err := service.UpdateCollection(ctx, testUserID, testCollectionID, updates)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无权更新")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新权限检查通过")
	})
}

// TestCollectionService_GetUserCollections 获取用户收藏列表测试
func TestCollectionService_GetUserCollections(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserCollections_Success", func(t *testing.T) {
		// Mock查询
		collections := []*reader.Collection{
			{ID: primitive.NewObjectID(), BookID: "book1"},
			{ID: primitive.NewObjectID(), BookID: "book2"},
		}
		mockRepo.On("GetCollectionsByUser", ctx, testUserID, "", 1, 20).
			Return(collections, int64(2), nil).Once()

		// 获取列表
		result, total, err := service.GetUserCollections(ctx, testUserID, "", 1, 20)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(2), total)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取收藏列表成功")
	})
}

// TestCollectionService_CreateFolder 创建收藏夹测试
func TestCollectionService_CreateFolder(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("CreateFolder_Success", func(t *testing.T) {
		// Mock创建
		mockRepo.On("CreateFolder", ctx, mock.AnythingOfType("*reader.CollectionFolder")).
			Return(nil).Once()

		// 创建收藏夹
		folder, err := service.CreateFolder(ctx, testUserID, "我的最爱", "经典作品", true)

		assert.NoError(t, err)
		assert.NotNil(t, folder)
		assert.Equal(t, "我的最爱", folder.Name)
		assert.True(t, folder.IsPublic)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 创建收藏夹成功")
	})

	t.Run("CreateFolder_EmptyName", func(t *testing.T) {
		_, err := service.CreateFolder(ctx, testUserID, "", "", false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "收藏夹名称")

		t.Logf("✓ 空名称验证通过")
	})
}

// TestCollectionService_UpdateFolder 更新收藏夹测试
func TestCollectionService_UpdateFolder(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("UpdateFolder_Success", func(t *testing.T) {
		// Mock获取收藏夹
		folder := &reader.CollectionFolder{
			ID:     primitive.NewObjectID(),
			UserID: testUserID,
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// Mock更新
		mockRepo.On("UpdateFolder", ctx, testFolderID, mock.AnythingOfType("map[string]interface {}")).
			Return(nil).Once()

		// 更新收藏夹
		updates := map[string]interface{}{
			"name":        "新名称",
			"description": "新描述",
		}
		err := service.UpdateFolder(ctx, testUserID, testFolderID, updates)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新收藏夹成功")
	})

	t.Run("UpdateFolder_NotOwner", func(t *testing.T) {
		// Mock获取收藏夹（不是自己的）
		folder := &reader.CollectionFolder{
			ID:     primitive.NewObjectID(),
			UserID: primitive.NewObjectID().Hex(),
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// 更新收藏夹
		updates := map[string]interface{}{"name": "新名称"}
		err := service.UpdateFolder(ctx, testUserID, testFolderID, updates)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无权更新")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新权限检查通过")
	})
}

// TestCollectionService_DeleteFolder 删除收藏夹测试
func TestCollectionService_DeleteFolder(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("DeleteFolder_Success", func(t *testing.T) {
		// Mock获取收藏夹
		folder := &reader.CollectionFolder{
			ID:        primitive.NewObjectID(),
			UserID:    testUserID,
			BookCount: 0,
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// Mock删除
		mockRepo.On("DeleteFolder", ctx, testFolderID).
			Return(nil).Once()

		// 删除收藏夹
		err := service.DeleteFolder(ctx, testUserID, testFolderID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除收藏夹成功")
	})

	t.Run("DeleteFolder_NotEmpty", func(t *testing.T) {
		// Mock获取收藏夹（包含收藏）
		folder := &reader.CollectionFolder{
			ID:        primitive.NewObjectID(),
			UserID:    testUserID,
			BookCount: 5,
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// 删除收藏夹
		err := service.DeleteFolder(ctx, testUserID, testFolderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "收藏夹不为空")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 非空收藏夹删除检测通过")
	})
}

// TestCollectionService_ShareCollection 分享收藏测试
func TestCollectionService_ShareCollection(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCollectionID := primitive.NewObjectID().Hex()

	t.Run("ShareCollection_Success", func(t *testing.T) {
		// Mock获取收藏
		collection := &reader.Collection{
			ID:       primitive.NewObjectID(),
			UserID:   testUserID,
			IsPublic: false,
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// Mock更新为公开
		mockRepo.On("Update", ctx, testCollectionID, mock.AnythingOfType("map[string]interface {}")).
			Return(nil).Once()

		// 分享收藏
		err := service.ShareCollection(ctx, testUserID, testCollectionID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 分享收藏成功")
	})
}

// TestCollectionService_GetCollectionStats 获取收藏统计测试
func TestCollectionService_GetCollectionStats(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := reading.NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetCollectionStats_Success", func(t *testing.T) {
		// Mock统计
		mockRepo.On("CountUserCollections", ctx, testUserID).
			Return(int64(50), nil).Once()

		// Mock收藏夹列表
		folders := []*reader.CollectionFolder{
			{
				ID:     primitive.NewObjectID(),
				UserID: testUserID,
				Name:   "默认收藏夹",
			},
		}
		mockRepo.On("GetFoldersByUser", ctx, testUserID).
			Return(folders, nil).Once()

		// 获取统计
		stats, err := service.GetUserCollectionStats(ctx, testUserID)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(50), stats["total_collections"])
		assert.Equal(t, 1, stats["total_folders"])

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取收藏统计成功")
	})
}
