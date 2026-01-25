package social

import (
	"Qingyu_backend/models/social"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCollectionRepository Mock收藏Repository
type MockCollectionRepository struct {
	mock.Mock
}

func (m *MockCollectionRepository) Create(ctx context.Context, collection *social.Collection) error {
	args := m.Called(ctx, collection)
	return args.Error(0)
}

func (m *MockCollectionRepository) GetByID(ctx context.Context, id string) (*social.Collection, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Collection), args.Error(1)
}

func (m *MockCollectionRepository) GetByUserAndBook(ctx context.Context, userID, bookID string) (*social.Collection, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Collection), args.Error(1)
}

func (m *MockCollectionRepository) GetCollectionsByUser(ctx context.Context, userID string, folderID string, page, size int) ([]*social.Collection, int64, error) {
	args := m.Called(ctx, userID, folderID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockCollectionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCollectionRepository) GetCollectionsByTag(ctx context.Context, userID string, tag string, page, size int) ([]*social.Collection, int64, error) {
	args := m.Called(ctx, userID, tag, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) CountUserCollections(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCollectionRepository) GetPublicCollections(ctx context.Context, page, size int) ([]*social.Collection, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) CreateFolder(ctx context.Context, folder *social.CollectionFolder) error {
	args := m.Called(ctx, folder)
	return args.Error(0)
}

func (m *MockCollectionRepository) GetFolderByID(ctx context.Context, id string) (*social.CollectionFolder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.CollectionFolder), args.Error(1)
}

func (m *MockCollectionRepository) GetFoldersByUser(ctx context.Context, userID string) ([]*social.CollectionFolder, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*social.CollectionFolder), args.Error(1)
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

func (m *MockCollectionRepository) GetPublicFolders(ctx context.Context, page, size int) ([]*social.CollectionFolder, int64, error) {
	args := m.Called(ctx, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.CollectionFolder), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestCollectionService_AddCollection 添加收藏测试
func TestCollectionService_AddCollection(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("AddCollection_Success", func(t *testing.T) {
		// Arrange - Mock检查是否已收藏
		mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).
			Return(nil, nil).Once()

		// Mock创建收藏
		mockRepo.On("Create", ctx, mock.MatchedBy(func(c *social.Collection) bool {
			return c.UserID == testUserID && c.BookID == testBookID
		})).
			Return(nil).Once()

		// Act - 添加收藏
		collection, err := service.AddToCollection(ctx, testUserID, testBookID, "", "很好的书", []string{"玄幻"}, false)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, collection)
		assert.Equal(t, testUserID, collection.UserID)
		assert.Equal(t, testBookID, collection.BookID)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 添加收藏成功")
	})

	t.Run("AddCollection_AlreadyCollected", func(t *testing.T) {
		// Arrange - Mock检查是否已收藏（已收藏）
		existingCollection := &social.Collection{
			ID:     primitive.NewObjectID(),
			UserID: testUserID,
			BookID: testBookID,
		}
		mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).
			Return(existingCollection, nil).Once()

		// Act - 添加收藏
		_, err := service.AddToCollection(ctx, testUserID, testBookID, "", "笔记", []string{}, false)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已经收藏")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 重复收藏检测通过")
	})

	t.Run("AddCollection_EmptyBookID", func(t *testing.T) {
		// Act
		_, err := service.AddToCollection(ctx, testUserID, "", "", "", []string{}, false)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "书籍ID")

		t.Logf("✓ 空书籍ID验证通过")
	})
}

// TestCollectionService_AddCollectionToFolder 添加到收藏夹测试
func TestCollectionService_AddCollectionToFolder(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("AddCollectionToFolder_Success", func(t *testing.T) {
		// Arrange - Mock检查是否已收藏
		mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).
			Return(nil, nil).Once()

		// Mock创建收藏
		mockRepo.On("Create", ctx, mock.MatchedBy(func(c *social.Collection) bool {
			return c.FolderID == testFolderID
		})).
			Return(nil).Once()

		// Mock增加收藏夹计数
		mockRepo.On("IncrementFolderBookCount", ctx, testFolderID).
			Return(nil).Once()

		// Act - 添加到收藏夹
		collection, err := service.AddToCollection(ctx, testUserID, testBookID, testFolderID, "", []string{}, false)

		// Assert
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

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCollectionID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("RemoveCollection_Success", func(t *testing.T) {
		// Arrange - Mock获取收藏
		collection := &social.Collection{
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

		// Act - 删除收藏
		err := service.RemoveFromCollection(ctx, testUserID, testCollectionID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除收藏成功")
	})

	t.Run("RemoveCollection_NotOwner", func(t *testing.T) {
		// Arrange - Mock获取收藏（不是自己的）
		collection := &social.Collection{
			ID:     primitive.NewObjectID(),
			UserID: primitive.NewObjectID(),
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// Act - 删除收藏
		err := service.RemoveFromCollection(ctx, testUserID, testCollectionID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无权")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除权限检查通过")
	})
}

// TestCollectionService_UpdateCollection 更新收藏测试
func TestCollectionService_UpdateCollection(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCollectionID := primitive.NewObjectID().Hex()

	t.Run("UpdateCollection_Success", func(t *testing.T) {
		// Arrange - Mock获取收藏
		collection := &social.Collection{
			ID:     primitive.NewObjectID(),
			UserID: testUserID,
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// Mock更新
		mockRepo.On("Update", ctx, testCollectionID, mock.AnythingOfType("map[string]interface {}")).
			Return(nil).Once()

		// Act - 更新收藏
		updates := map[string]interface{}{
			"note": "新笔记",
			"tags": []string{"新标签"},
		}
		err := service.UpdateCollection(ctx, testUserID, testCollectionID, updates)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新收藏成功")
	})

	t.Run("UpdateCollection_NotOwner", func(t *testing.T) {
		// Arrange - Mock获取收藏（不是自己的）
		collection := &social.Collection{
			ID:     primitive.NewObjectID(),
			UserID: primitive.NewObjectID(),
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// Act - 更新收藏
		updates := map[string]interface{}{"note": "新笔记"}
		err := service.UpdateCollection(ctx, testUserID, testCollectionID, updates)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无权")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新权限检查通过")
	})
}

// TestCollectionService_GetUserCollections 获取用户收藏列表测试
func TestCollectionService_GetUserCollections(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserCollections_Success", func(t *testing.T) {
		// Arrange - Mock查询
		collections := []*social.Collection{
			{ID: primitive.NewObjectID(), BookID: "book1"},
			{ID: primitive.NewObjectID(), BookID: "book2"},
		}
		mockRepo.On("GetCollectionsByUser", ctx, testUserID, "", 1, 20).
			Return(collections, int64(2), nil).Once()

		// Act - 获取列表
		result, total, err := service.GetUserCollections(ctx, testUserID, "", 1, 20)

		// Assert
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

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("CreateFolder_Success", func(t *testing.T) {
		// Arrange - Mock创建
		mockRepo.On("CreateFolder", ctx, mock.MatchedBy(func(f *social.CollectionFolder) bool {
			return f.UserID == testUserID && f.Name == "我的最爱"
		})).
			Return(nil).Once()

		// Act - 创建收藏夹
		folder, err := service.CreateFolder(ctx, testUserID, "我的最爱", "经典作品", true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, folder)
		assert.Equal(t, "我的最爱", folder.Name)
		assert.True(t, folder.IsPublic)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 创建收藏夹成功")
	})

	t.Run("CreateFolder_EmptyName", func(t *testing.T) {
		// Act
		_, err := service.CreateFolder(ctx, testUserID, "", "", false)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "收藏夹名称")

		t.Logf("✓ 空名称验证通过")
	})
}

// TestCollectionService_UpdateFolder 更新收藏夹测试
func TestCollectionService_UpdateFolder(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("UpdateFolder_Success", func(t *testing.T) {
		// Arrange - Mock获取收藏夹
		folder := &social.CollectionFolder{
			ID:     primitive.NewObjectID(),
			UserID: testUserID,
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// Mock更新
		mockRepo.On("UpdateFolder", ctx, testFolderID, mock.AnythingOfType("map[string]interface {}")).
			Return(nil).Once()

		// Act - 更新收藏夹
		updates := map[string]interface{}{
			"name":        "新名称",
			"description": "新描述",
		}
		err := service.UpdateFolder(ctx, testUserID, testFolderID, updates)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新收藏夹成功")
	})

	t.Run("UpdateFolder_NotOwner", func(t *testing.T) {
		// Arrange - Mock获取收藏夹（不是自己的）
		folder := &social.CollectionFolder{
			ID:     primitive.NewObjectID(),
			UserID: primitive.NewObjectID(),
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// Act - 更新收藏夹
		updates := map[string]interface{}{"name": "新名称"}
		err := service.UpdateFolder(ctx, testUserID, testFolderID, updates)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无权")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 更新权限检查通过")
	})
}

// TestCollectionService_DeleteFolder 删除收藏夹测试
func TestCollectionService_DeleteFolder(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testFolderID := primitive.NewObjectID().Hex()

	t.Run("DeleteFolder_Success", func(t *testing.T) {
		// Arrange - Mock获取收藏夹
		folder := &social.CollectionFolder{
			ID:        primitive.NewObjectID(),
			UserID:    testUserID,
			BookCount: 0,
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// Mock删除
		mockRepo.On("DeleteFolder", ctx, testFolderID).
			Return(nil).Once()

		// Act - 删除收藏夹
		err := service.DeleteFolder(ctx, testUserID, testFolderID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 删除收藏夹成功")
	})

	t.Run("DeleteFolder_NotEmpty", func(t *testing.T) {
		// Arrange - Mock获取收藏夹（包含收藏）
		folder := &social.CollectionFolder{
			ID:        primitive.NewObjectID(),
			UserID:    testUserID,
			BookCount: 5,
		}
		mockRepo.On("GetFolderByID", ctx, testFolderID).
			Return(folder, nil).Once()

		// Act - 删除收藏夹
		err := service.DeleteFolder(ctx, testUserID, testFolderID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "收藏夹")

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 非空收藏夹删除检测通过")
	})
}

// TestCollectionService_ShareCollection 分享收藏测试
func TestCollectionService_ShareCollection(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCollectionID := primitive.NewObjectID().Hex()

	t.Run("ShareCollection_Success", func(t *testing.T) {
		// Arrange - Mock获取收藏
		collection := &social.Collection{
			ID:       primitive.NewObjectID(),
			UserID:   testUserID,
			IsPublic: false,
		}
		mockRepo.On("GetByID", ctx, testCollectionID).
			Return(collection, nil).Once()

		// Mock更新为公开
		mockRepo.On("Update", ctx, testCollectionID, mock.AnythingOfType("map[string]interface {}")).
			Return(nil).Once()

		// Act - 分享收藏
		err := service.ShareCollection(ctx, testUserID, testCollectionID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 分享收藏成功")
	})
}

// TestCollectionService_GetCollectionStats 获取收藏统计测试
func TestCollectionService_GetCollectionStats(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetCollectionStats_Success", func(t *testing.T) {
		// Arrange - Mock统计
		mockRepo.On("CountUserCollections", ctx, testUserID).
			Return(int64(50), nil).Once()

		// Mock收藏夹列表
		folders := []*social.CollectionFolder{
			{
				ID:     primitive.NewObjectID(),
				UserID: testUserID,
				Name:   "默认收藏夹",
			},
		}
		mockRepo.On("GetFoldersByUser", ctx, testUserID).
			Return(folders, nil).Once()

		// Act - 获取统计
		stats, err := service.GetUserCollectionStats(ctx, testUserID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(50), stats["total_collections"])
		assert.Equal(t, 1, stats["total_folders"])

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取收藏统计成功")
	})
}

// TestCollectionServiceTableDriven 表格驱动测试
func TestCollectionServiceTableDriven(t *testing.T) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	tests := []struct {
		name        string
		setupMock   func()
		action      func() error
		wantErr     bool
		errContains string
	}{
		{
			name: "成功添加收藏",
			setupMock: func() {
				mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).Return(nil, nil).Once()
				mockRepo.On("Create", ctx, mock.AnythingOfType("*social.Collection")).Return(nil).Once()
			},
			action: func() error {
				_, err := service.AddToCollection(ctx, testUserID, testBookID, "", "笔记", []string{}, false)
				return err
			},
			wantErr: false,
		},
		{
			name: "书籍ID为空",
			setupMock: func() {
				// 不调用mock
			},
			action: func() error {
				_, err := service.AddToCollection(ctx, testUserID, "", "", "笔记", []string{}, false)
				return err
			},
			wantErr:     true,
			errContains: "书籍ID",
		},
		{
			name: "重复收藏",
			setupMock: func() {
				existing := &social.Collection{ID: primitive.NewObjectID(), UserID: testUserID, BookID: testBookID}
				mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).Return(existing, nil).Once()
			},
			action: func() error {
				_, err := service.AddToCollection(ctx, testUserID, testBookID, "", "笔记", []string{}, false)
				return err
			},
			wantErr:     true,
			errContains: "已经收藏",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMock()

			// Act
			err := tt.action()

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)

			t.Logf("✓ %s", tt.name)
		})
	}
}

// BenchmarkCollectionService_AddToCollection 性能测试
func BenchmarkCollectionService_AddToCollection(b *testing.B) {
	mockRepo := new(MockCollectionRepository)
	mockEventBus := NewMockEventBus()

	service := NewCollectionService(mockRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	mockRepo.On("GetByUserAndBook", ctx, testUserID, testBookID).Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*social.Collection")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.AddToCollection(ctx, testUserID, testBookID, "", "笔记", []string{}, false)
	}
}
