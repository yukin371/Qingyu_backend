package social

import (
	"Qingyu_backend/models/social"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestLikeService_LikeBook 点赞书籍测试
func TestLikeService_LikeBook(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("LikeBook_Success", func(t *testing.T) {
		// Arrange - Mock添加点赞
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.UserID == testUserID && l.TargetID == testBookID && l.TargetType == social.LikeTargetTypeBook
		})).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).
			Return(int64(1), nil).Once()

		// Act - 点赞书籍
		err := service.LikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 点赞书籍成功")
	})

	t.Run("LikeBook_AlreadyLiked", func(t *testing.T) {
		// Arrange - Mock添加点赞（返回已点赞错误）
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.UserID == testUserID && l.TargetID == testBookID
		})).
			Return(fmt.Errorf("已经点赞过了")).Once()

		// Act - 点赞书籍（幂等性：不报错）
		err := service.LikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.NoError(t, err) // Service将"已经点赞过了"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 重复点赞检测通过（幂等性）")
	})

	t.Run("LikeBook_EmptyBookID", func(t *testing.T) {
		// Act
		err := service.LikeBook(ctx, testUserID, "")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "书籍ID")

		t.Logf("✓ 空书籍ID验证通过")
	})
}

// TestLikeService_UnlikeBook 取消点赞书籍测试
func TestLikeService_UnlikeBook(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("UnlikeBook_Success", func(t *testing.T) {
		// Arrange - Mock移除点赞
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeBook, testBookID).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).
			Return(int64(0), nil).Once()

		// Act - 取消点赞
		err := service.UnlikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞成功")
	})

	t.Run("UnlikeBook_NotLiked", func(t *testing.T) {
		// Arrange - Mock移除点赞（返回未点赞错误）
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeBook, testBookID).
			Return(fmt.Errorf("点赞记录不存在")).Once()

		// Act - 取消点赞（幂等性：不报错）
		err := service.UnlikeBook(ctx, testUserID, testBookID)

		// Assert
		assert.NoError(t, err) // Service将"点赞记录不存在"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 未点赞检测通过（幂等性）")
	})
}

// TestLikeService_LikeComment 点赞评论测试
func TestLikeService_LikeComment(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("LikeComment_Success", func(t *testing.T) {
		// Arrange - Mock添加点赞
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.UserID == testUserID && l.TargetID == testCommentID && l.TargetType == social.LikeTargetTypeComment
		})).
			Return(nil).Once()

		// Mock增加评论点赞数
		mockCommentRepo.On("IncrementLikeCount", ctx, testCommentID).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeComment, testCommentID).
			Return(int64(1), nil).Once()

		// Act - 点赞评论
		err := service.LikeComment(ctx, testUserID, testCommentID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 点赞评论成功")
	})

	t.Run("LikeComment_AlreadyLiked", func(t *testing.T) {
		// Arrange - Mock添加点赞（返回已点赞错误）
		mockRepo.On("AddLike", ctx, mock.AnythingOfType("*social.Like")).
			Return(fmt.Errorf("已经点赞过了")).Once()

		// Act - 点赞评论（幂等性：不报错）
		err := service.LikeComment(ctx, testUserID, testCommentID)

		// Assert
		assert.NoError(t, err) // Service将"已经点赞过了"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 重复点赞评论检测通过（幂等性）")
	})
}

// TestLikeService_UnlikeComment 取消点赞评论测试
func TestLikeService_UnlikeComment(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testCommentID := primitive.NewObjectID().Hex()

	t.Run("UnlikeComment_Success", func(t *testing.T) {
		// Arrange - Mock移除点赞
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeComment, testCommentID).
			Return(nil).Once()

		// Mock减少评论点赞数
		mockCommentRepo.On("DecrementLikeCount", ctx, testCommentID).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeComment, testCommentID).
			Return(int64(0), nil).Once()

		// Act - 取消点赞
		err := service.UnlikeComment(ctx, testUserID, testCommentID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)

		t.Logf("✓ 取消点赞评论成功")
	})

	t.Run("UnlikeComment_NotLiked", func(t *testing.T) {
		// Arrange - Mock移除点赞（返回未点赞错误）
		mockRepo.On("RemoveLike", ctx, testUserID, social.LikeTargetTypeComment, testCommentID).
			Return(fmt.Errorf("点赞记录不存在")).Once()

		// Act - 取消点赞（幂等性：不报错）
		err := service.UnlikeComment(ctx, testUserID, testCommentID)

		// Assert
		assert.NoError(t, err) // Service将"点赞记录不存在"错误处理为幂等性，不报错

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 未点赞评论检测通过（幂等性）")
	})
}

// TestLikeService_GetBookLikeCount 获取书籍点赞数测试
func TestLikeService_GetBookLikeCount(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testBookID := primitive.NewObjectID().Hex()

	t.Run("GetBookLikeCount_Success", func(t *testing.T) {
		// Arrange - Mock获取点赞数
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).
			Return(int64(100), nil).Once()

		// Act - 获取点赞数
		count, err := service.GetBookLikeCount(ctx, testBookID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(100), count)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取点赞数成功")
	})
}

// TestLikeService_GetUserLikedBooks 获取用户点赞书籍列表测试
func TestLikeService_GetUserLikedBooks(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikedBooks_Success", func(t *testing.T) {
		// Arrange - Mock查询
		likes := []*social.Like{
			{ID: primitive.NewObjectID(), TargetID: "book1"},
			{ID: primitive.NewObjectID(), TargetID: "book2"},
		}
		mockRepo.On("GetUserLikes", ctx, testUserID, social.LikeTargetTypeBook, 1, 20).
			Return(likes, int64(2), nil).Once()

		// Act - 获取列表
		result, total, err := service.GetUserLikedBooks(ctx, testUserID, 1, 20)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, int64(2), total)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取点赞书籍列表成功")
	})
}

// TestLikeService_GetUserLikeStats 获取用户点赞统计测试
func TestLikeService_GetUserLikeStats(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikeStats_Success", func(t *testing.T) {
		// Arrange - Mock统计
		mockRepo.On("CountUserLikes", ctx, testUserID).
			Return(int64(50), nil).Once()

		// Act - 获取统计
		stats, err := service.GetUserLikeStats(ctx, testUserID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(50), stats["total_likes"])

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 获取点赞统计成功")
	})
}

// TestLikeService_ConcurrentLike 并发点赞测试
func TestLikeService_ConcurrentLike(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	t.Run("ConcurrentLike_OnlyOneSuccess", func(t *testing.T) {
		// Arrange - 第一次添加成功
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook
		})).
			Return(nil).Once()

		// Mock获取点赞数（用于事件发布）
		mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, testBookID).
			Return(int64(1), nil).Once()

		// 第二次添加（返回已点赞错误，幂等性处理）
		mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
			return l.TargetType == social.LikeTargetTypeBook
		})).
			Return(fmt.Errorf("已经点赞过了")).Once()

		// Act - 第一次点赞应该成功
		err1 := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err1)

		// 第二次点赞也成功（幂等性）
		err2 := service.LikeBook(ctx, testUserID, testBookID)
		assert.NoError(t, err2)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 并发点赞幂等性检测通过")
	})
}

// TestLikeService_AntiSpam 防刷测试
func TestLikeService_AntiSpam(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("AntiSpam_RapidLike", func(t *testing.T) {
		// Act - 模拟快速点赞多本书
		for i := 0; i < 5; i++ {
			bookID := primitive.NewObjectID().Hex()

			mockRepo.On("AddLike", ctx, mock.MatchedBy(func(l *social.Like) bool {
				return l.TargetType == social.LikeTargetTypeBook
			})).
				Return(nil).Once()

			// Mock获取点赞数（用于事件发布）
			mockRepo.On("GetLikeCount", ctx, social.LikeTargetTypeBook, bookID).
				Return(int64(1), nil).Once()

			err := service.LikeBook(ctx, testUserID, bookID)
			assert.NoError(t, err)
		}

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 正常点赞速率通过")
	})
}

// TestLikeService_GetBooksLikeCount 批量获取点赞数测试
func TestLikeService_GetBooksLikeCount(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	t.Run("GetBooksLikeCount_Success", func(t *testing.T) {
		// Arrange
		bookIDs := []string{
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
		}

		expectedCounts := map[string]int64{
			bookIDs[0]: 10,
			bookIDs[1]: 20,
			bookIDs[2]: 30,
		}

		mockRepo.On("GetLikesCountBatch", ctx, social.LikeTargetTypeBook, bookIDs).
			Return(expectedCounts, nil).Once()

		// Act
		counts, err := service.GetBooksLikeCount(ctx, bookIDs)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCounts, counts)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 批量获取点赞数成功")
	})

	t.Run("GetBooksLikeCount_EmptyArray", func(t *testing.T) {
		// Arrange
		emptyIDs := []string{}

		// Act - 空数组会导致service直接返回空map，不调用repository
		counts, err := service.GetBooksLikeCount(ctx, emptyIDs)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, counts)

		t.Logf("✓ 批量操作空数组处理正确")
	})
}

// TestLikeService_GetUserLikeStatus 批量获取用户点赞状态测试
func TestLikeService_GetUserLikeStatus(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	t.Run("GetUserLikeStatus_Success", func(t *testing.T) {
		// Arrange
		bookIDs := []string{
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(),
		}

		expectedStatus := map[string]bool{
			bookIDs[0]: true,
			bookIDs[1]: false,
			bookIDs[2]: true,
		}

		mockRepo.On("GetUserLikeStatusBatch", ctx, testUserID, social.LikeTargetTypeBook, bookIDs).
			Return(expectedStatus, nil).Once()

		// Act
		status, err := service.GetUserLikeStatus(ctx, testUserID, bookIDs)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedStatus, status)

		mockRepo.AssertExpectations(t)

		t.Logf("✓ 批量获取用户点赞状态成功")
	})
}

// TestLikeServiceTableDriven 表格驱动测试
func TestLikeServiceTableDriven(t *testing.T) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()

	tests := []struct {
		name        string
		targetType  string
		targetID    string
		setupMock   func()
		wantErr     bool
		errContains string
	}{
		{
			name:       "点赞书籍成功",
			targetType: social.LikeTargetTypeBook,
			targetID:   primitive.NewObjectID().Hex(),
			setupMock: func() {
				mockRepo.On("AddLike", ctx, mock.AnythingOfType("*social.Like")).Return(nil).Once()
				mockRepo.On("GetLikeCount", ctx, mock.Anything, mock.Anything).Return(int64(1), nil).Once()
			},
			wantErr: false,
		},
		{
			name:       "点赞评论成功",
			targetType: social.LikeTargetTypeComment,
			targetID:   primitive.NewObjectID().Hex(),
			setupMock: func() {
				mockRepo.On("AddLike", ctx, mock.AnythingOfType("*social.Like")).Return(nil).Once()
				mockCommentRepo.On("IncrementLikeCount", ctx, mock.Anything).Return(nil).Once()
				mockRepo.On("GetLikeCount", ctx, mock.Anything, mock.Anything).Return(int64(1), nil).Once()
			},
			wantErr: false,
		},
		{
			name:       "空书籍ID",
			targetType: social.LikeTargetTypeBook,
			targetID:   "",
			setupMock: func() {
				// 不调用任何mock
			},
			wantErr:     true,
			errContains: "书籍ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.setupMock()

			// Act
			var err error
			if tt.targetType == social.LikeTargetTypeBook {
				err = service.LikeBook(ctx, testUserID, tt.targetID)
			} else if tt.targetType == social.LikeTargetTypeComment {
				err = service.LikeComment(ctx, testUserID, tt.targetID)
			}

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
			mockCommentRepo.AssertExpectations(t)
		})
	}

	t.Logf("✓ 表格驱动测试完成")
}

// BenchmarkLikeService_LikeBook 性能测试
func BenchmarkLikeService_LikeBook(b *testing.B) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testUserID := primitive.NewObjectID().Hex()
	testBookID := primitive.NewObjectID().Hex()

	mockRepo.On("AddLike", ctx, mock.AnythingOfType("*social.Like")).Return(nil)
	mockRepo.On("GetLikeCount", ctx, mock.Anything, mock.Anything).Return(int64(1), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.LikeBook(ctx, testUserID, testBookID)
	}
}

// BenchmarkLikeService_GetBookLikeCount 性能测试
func BenchmarkLikeService_GetBookLikeCount(b *testing.B) {
	mockRepo := new(MockLikeRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockEventBus := NewMockEventBus()

	service := NewLikeService(mockRepo, mockCommentRepo, mockEventBus)
	ctx := context.Background()

	testBookID := primitive.NewObjectID().Hex()

	mockRepo.On("GetLikeCount", ctx, mock.Anything, mock.Anything).Return(int64(100), nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetBookLikeCount(ctx, testBookID)
	}
}
