package interfaces

import (
	"context"

	socialModel "Qingyu_backend/models/social"
)

// PostService 动态服务接口
type PostService interface {
	// =========================
	// 动态管理
	// =========================

	// CreatePost 创建动态
	// userID: 用户ID
	// userName: 用户昵称
	// userAvatar: 用户头像
	// userLevel: 用户等级
	// postType: 动态类型(text/image/book_recommendation/reading_progress/poll)
	// content: 内容
	// images: 图片列表
	// bookID, bookTitle, bookCover, bookAuthor: 书籍信息(可选)
	// chapterID, chapterTitle, progress: 阅读进度信息(可选)
	// topics: 话题标签
	CreatePost(ctx context.Context, userID, userName, userAvatar string, userLevel int,
		postType socialModel.PostType, content string, images []string,
		bookID, bookTitle, bookCover, bookAuthor string,
		chapterID, chapterTitle string, progress int,
		topics []string) (*socialModel.Post, error)

	// GetPosts 获取动态列表
	// userID: 当前用户ID(用于获取点赞状态，可为空)
	// topic: 话题筛选(可选)
	// sort: 排序方式("latest"/"hottest")
	GetPosts(ctx context.Context, userID string, page, size int, topic, sort string) ([]*socialModel.PostInfo, int64, error)

	// GetPostByID 获取动态详情
	// userID: 当前用户ID(用于获取点赞状态，可为空)
	GetPostByID(ctx context.Context, userID, postID string) (*socialModel.PostInfo, error)

	// UpdatePost 更新动态
	UpdatePost(ctx context.Context, userID, postID string, content string, topics []string) error

	// DeletePost 删除动态
	DeletePost(ctx context.Context, userID, postID string) error

	// GetUserPosts 获取用户发布的动态列表
	GetUserPosts(ctx context.Context, userID string, page, size int) ([]*socialModel.PostInfo, int64, error)

	// =========================
	// 动态点赞
	// =========================

	// ToggleLike 切换点赞状态
	// 返回当前点赞状态(true=已点赞, false=未点赞)
	ToggleLike(ctx context.Context, userID, postID string) (bool, error)

	// Like 点赞动态
	Like(ctx context.Context, userID, postID string) error

	// Unlike 取消点赞动态
	Unlike(ctx context.Context, userID, postID string) error
}
