package validation

import (
	"context"
	"errors"
	"fmt"
)

// UserRepository 用户仓库接口（定义最小需要的接口）
type UserRepository interface {
	Exists(ctx context.Context, userID string) (bool, error)
}

// BookRepository 书籍仓库接口（定义最小需要的接口）
type BookRepository interface {
	Exists(ctx context.Context, bookID string) (bool, error)
}

// ChapterRepository 章节仓库接口（预留）
type ChapterRepository interface {
	Exists(ctx context.Context, chapterID string) (bool, error)
}

// CommentRepository 评论仓库接口（定义最小需要的接口）
type CommentRepository interface {
	Exists(ctx context.Context, commentID string) (bool, error)
}

// ReferenceValidator 引用验证器 - 用于验证外键引用的有效性
type ReferenceValidator struct {
	userRepo     UserRepository
	bookRepo     BookRepository
	chapterRepo  ChapterRepository
	commentRepo  CommentRepository
}

// NewReferenceValidator 创建引用验证器
func NewReferenceValidator(
	userRepo UserRepository,
	bookRepo BookRepository,
	chapterRepo ChapterRepository,
	commentRepo CommentRepository,
) *ReferenceValidator {
	return &ReferenceValidator{
		userRepo:    userRepo,
		bookRepo:    bookRepo,
		chapterRepo: chapterRepo,
		commentRepo: commentRepo,
	}
}

// =======================
// 用户验证
// =======================

// ValidateUserExists 验证用户存在
func (v *ReferenceValidator) ValidateUserExists(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("用户ID不能为空")
	}

	exists, err := v.userRepo.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("检查用户失败: %w", err)
	}
	if !exists {
		return errors.New("用户不存在")
	}
	return nil
}

// =======================
// 书籍验证
// =======================

// ValidateBookExists 验证书籍存在
func (v *ReferenceValidator) ValidateBookExists(ctx context.Context, bookID string) error {
	if bookID == "" {
		return errors.New("书籍ID不能为空")
	}

	exists, err := v.bookRepo.Exists(ctx, bookID)
	if err != nil {
		return fmt.Errorf("检查书籍失败: %w", err)
	}
	if !exists {
		return errors.New("书籍不存在")
	}
	return nil
}

// =======================
// 章节验证
// =======================

// ValidateChapterExists 验证章节存在
func (v *ReferenceValidator) ValidateChapterExists(ctx context.Context, chapterID string) error {
	if chapterID == "" {
		return errors.New("章节ID不能为空")
	}

	if v.chapterRepo == nil {
		// 如果没有提供章节仓库，跳过验证
		return nil
	}

	exists, err := v.chapterRepo.Exists(ctx, chapterID)
	if err != nil {
		return fmt.Errorf("检查章节失败: %w", err)
	}
	if !exists {
		return errors.New("章节不存在")
	}
	return nil
}

// =======================
// 评论验证
// =======================

// ValidateCommentExists 验证评论存在
func (v *ReferenceValidator) ValidateCommentExists(ctx context.Context, commentID string) error {
	if commentID == "" {
		return errors.New("评论ID不能为空")
	}

	exists, err := v.commentRepo.Exists(ctx, commentID)
	if err != nil {
		return fmt.Errorf("检查评论失败: %w", err)
	}
	if !exists {
		return errors.New("评论不存在")
	}
	return nil
}

// =======================
// 点赞引用验证
// =======================

// ValidateLikeReference 验证点赞引用
// targetType: "book", "comment" 等
func (v *ReferenceValidator) ValidateLikeReference(ctx context.Context, userID, targetID, targetType string) error {
	// 验证用户存在
	if err := v.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	// 根据目标类型验证目标存在
	switch targetType {
	case "book":
		if err := v.ValidateBookExists(ctx, targetID); err != nil {
			return err
		}
	case "comment":
		if err := v.ValidateCommentExists(ctx, targetID); err != nil {
			return err
		}
	default:
		return errors.New("无效的目标类型")
	}

	return nil
}

// =======================
// 评论引用验证
// =======================

// ValidateCommentReference 验证评论引用（书籍评论）
func (v *ReferenceValidator) ValidateCommentReference(ctx context.Context, authorID, bookID string) error {
	// 验证作者存在
	if err := v.ValidateUserExists(ctx, authorID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	return nil
}

// ValidateCommentReferenceReply 验证评论引用（回复评论）
func (v *ReferenceValidator) ValidateCommentReferenceReply(ctx context.Context, authorID, parentCommentID string) error {
	// 验证作者存在
	if err := v.ValidateUserExists(ctx, authorID); err != nil {
		return err
	}

	// 验证父评论存在
	if err := v.ValidateCommentExists(ctx, parentCommentID); err != nil {
		return errors.New("父评论不存在")
	}

	return nil
}

// =======================
// 关注引用验证
// =======================

// ValidateFollowReference 验证关注引用
func (v *ReferenceValidator) ValidateFollowReference(ctx context.Context, followerID, followingID string) error {
	// 先检查是否尝试关注自己
	if followerID == followingID {
		return errors.New("不能关注自己")
	}

	// 验证关注者存在
	if err := v.ValidateUserExists(ctx, followerID); err != nil {
		return errors.New("关注者不存在")
	}

	// 验证被关注用户存在
	if err := v.ValidateUserExists(ctx, followingID); err != nil {
		return errors.New("被关注用户不存在")
	}

	return nil
}

// =======================
// 收益记录引用验证
// =======================

// ValidateRevenueReference 验证收益记录引用
func (v *ReferenceValidator) ValidateRevenueReference(ctx context.Context, authorID, bookID string) error {
	// 验证作者存在
	if err := v.ValidateUserExists(ctx, authorID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	return nil
}

// =======================
// 阅读进度引用验证
// =======================

// ValidateReadingProgressReference 验证阅读进度引用
func (v *ReferenceValidator) ValidateReadingProgressReference(ctx context.Context, userID, bookID string) error {
	// 验证用户存在
	if err := v.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	return nil
}

// =======================
// 阅读历史引用验证
// =======================

// ValidateReadingHistoryReference 验证阅读历史引用
func (v *ReferenceValidator) ValidateReadingHistoryReference(ctx context.Context, userID, bookID, chapterID string) error {
	// 验证用户存在
	if err := v.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	// 验证章节存在（可选）
	if chapterID != "" {
		if err := v.ValidateChapterExists(ctx, chapterID); err != nil {
			return err
		}
	}

	return nil
}

// =======================
// 收藏引用验证
// =======================

// ValidateCollectionReference 验证收藏引用
func (v *ReferenceValidator) ValidateCollectionReference(ctx context.Context, userID, bookID string) error {
	// 验证用户存在
	if err := v.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	return nil
}

// =======================
// 书签引用验证
// =======================

// ValidateBookmarkReference 验证书签引用
func (v *ReferenceValidator) ValidateBookmarkReference(ctx context.Context, userID, bookID, chapterID string) error {
	// 验证用户存在
	if err := v.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	// 验证章节存在
	if chapterID != "" {
		if err := v.ValidateChapterExists(ctx, chapterID); err != nil {
			return err
		}
	}

	return nil
}

// =======================
// 批注引用验证
// =======================

// ValidateAnnotationReference 验证批注引用
func (v *ReferenceValidator) ValidateAnnotationReference(ctx context.Context, userID, bookID, chapterID string) error {
	// 验证用户存在
	if err := v.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	// 验证章节存在
	if chapterID != "" {
		if err := v.ValidateChapterExists(ctx, chapterID); err != nil {
			return err
		}
	}

	return nil
}

// =======================
// 评分引用验证
// =======================

// ValidateRatingReference 验证评分引用
func (v *ReferenceValidator) ValidateRatingReference(ctx context.Context, userID, bookID string) error {
	// 验证用户存在
	if err := v.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	// 验证书籍存在
	if err := v.ValidateBookExists(ctx, bookID); err != nil {
		return err
	}

	return nil
}
