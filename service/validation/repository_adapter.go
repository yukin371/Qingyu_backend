package validation

import (
	"context"

	"Qingyu_backend/repository/interfaces/bookstore"
	"Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/repository/interfaces/user"
)

// UserRepoAdapter 用户仓储适配器
type UserRepoAdapter struct {
	repo user.UserRepository
}

func NewUserRepoAdapter(repo user.UserRepository) *UserRepoAdapter {
	return &UserRepoAdapter{repo: repo}
}

func (a *UserRepoAdapter) Exists(ctx context.Context, userID string) (bool, error) {
	// 使用GetByID来检查存在性
	_, err := a.repo.GetByID(ctx, userID)
	if err != nil {
		// 检查是否是"未找到"错误
		if user.IsNotFoundError(err) || err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// BookRepoAdapter 书籍仓储适配器
type BookRepoAdapter struct {
	repo bookstore.BookDetailRepository
}

func NewBookRepoAdapter(repo bookstore.BookDetailRepository) *BookRepoAdapter {
	return &BookRepoAdapter{repo: repo}
}

func (a *BookRepoAdapter) Exists(ctx context.Context, bookID string) (bool, error) {
	// 使用GetByID来检查存在性
	_, err := a.repo.GetByID(ctx, bookID)
	if err != nil {
		// 检查是否是"未找到"错误
		if err.Error() == "record not found" || err.Error() == "book detail not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CommentRepoAdapter 评论仓储适配器
type CommentRepoAdapter struct {
	repo social.CommentRepository
}

func NewCommentRepoAdapter(repo social.CommentRepository) *CommentRepoAdapter {
	return &CommentRepoAdapter{repo: repo}
}

func (a *CommentRepoAdapter) Exists(ctx context.Context, commentID string) (bool, error) {
	// 使用Exists方法（如果已实现）
	return a.repo.Exists(ctx, commentID)
}

// ChapterRepoAdapter 章节仓储适配器（预留）
type ChapterRepoAdapter struct {
	// 章节仓储接口
}

func NewChapterRepoAdapter() *ChapterRepoAdapter {
	return &ChapterRepoAdapter{}
}

func (a *ChapterRepoAdapter) Exists(ctx context.Context, chapterID string) (bool, error) {
	// 暂不实现，返回true跳过验证
	return true, nil
}
