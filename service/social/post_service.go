package social

import (
	"context"
	"fmt"
	"time"

	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/service/base"
	"Qingyu_backend/service/interfaces"

	socialModel "Qingyu_backend/models/social"
)

// PostService 动态服务
type PostService struct {
	postRepo    socialRepo.PostRepository
	eventBus    base.EventBus
	serviceName string
	version     string
}

// NewPostService 创建动态服务实例
func NewPostService(
	postRepo socialRepo.PostRepository,
	eventBus base.EventBus,
) *PostService {
	if eventBus == nil {
		eventBus = base.NewSimpleEventBus()
	}
	return &PostService{
		postRepo:    postRepo,
		eventBus:    eventBus,
		serviceName: "PostService",
		version:     "1.0.0",
	}
}

// 确保 PostService 实现了 interfaces.PostService 接口
var _ interfaces.PostService = (*PostService)(nil)

// BaseService 接口实现
func (s *PostService) Initialize(ctx context.Context) error { return nil }
func (s *PostService) Health(ctx context.Context) error {
	if err := s.postRepo.Health(ctx); err != nil {
		return fmt.Errorf("动态Repository健康检查失败: %w", err)
	}
	return nil
}
func (s *PostService) Close(ctx context.Context) error { return nil }
func (s *PostService) GetServiceName() string          { return s.serviceName }
func (s *PostService) GetVersion() string              { return s.version }

// CreatePost 创建动态
func (s *PostService) CreatePost(ctx context.Context, userID, userName, userAvatar string, userLevel int,
	postType socialModel.PostType, content string, images []string,
	bookID, bookTitle, bookCover, bookAuthor string,
	chapterID, chapterTitle string, progress int,
	topics []string) (*socialModel.Post, error) {

	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if content == "" {
		return nil, fmt.Errorf("内容不能为空")
	}
	if len(content) > 5000 {
		return nil, fmt.Errorf("内容最多5000字")
	}

	post := &socialModel.Post{
		UserID:       userID,
		UserName:     userName,
		UserAvatar:   userAvatar,
		UserLevel:    userLevel,
		Type:         postType,
		Content:      content,
		Images:       images,
		BookID:       bookID,
		BookTitle:    bookTitle,
		BookCover:    bookCover,
		BookAuthor:   bookAuthor,
		ChapterID:    chapterID,
		ChapterTitle: chapterTitle,
		Progress:     progress,
		Topics:       topics,
		LikeCount:    0,
		CommentCount: 0,
		ShareCount:   0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("创建动态失败: %w", err)
	}

	s.publishPostEvent(ctx, "post.created", userID, post.ID.Hex())
	return post, nil
}

// GetPosts 获取动态列表
func (s *PostService) GetPosts(ctx context.Context, userID string, page, size int, topic, sort string) ([]*socialModel.PostInfo, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	if sort == "" {
		sort = "latest"
	}

	posts, total, err := s.postRepo.List(ctx, page, size, topic, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("获取动态列表失败: %w", err)
	}

	// 批量获取点赞状态
	postIDs := make([]string, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID.Hex()
	}

	likedMap := make(map[string]bool)
	if userID != "" && len(postIDs) > 0 {
		likedMap, err = s.postRepo.GetLikedPostIDs(ctx, userID, postIDs)
		if err != nil {
			// 如果获取点赞状态失败，继续返回列表但不带点赞状态
			likedMap = make(map[string]bool)
		}
	}

	// 转换为 PostInfo
	postInfos := make([]*socialModel.PostInfo, len(posts))
	for i, post := range posts {
		postInfos[i] = post.ToPostInfo(likedMap[post.ID.Hex()], false)
	}

	return postInfos, total, nil
}

// GetPostByID 获取动态详情
func (s *PostService) GetPostByID(ctx context.Context, userID, postID string) (*socialModel.PostInfo, error) {
	if postID == "" {
		return nil, fmt.Errorf("动态ID不能为空")
	}

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("获取动态详情失败: %w", err)
	}

	// 获取点赞状态
	isLiked := false
	if userID != "" {
		isLiked, err = s.postRepo.IsLiked(ctx, postID, userID)
		if err != nil {
			isLiked = false
		}
	}

	return post.ToPostInfo(isLiked, false), nil
}

// UpdatePost 更新动态
func (s *PostService) UpdatePost(ctx context.Context, userID, postID string, content string, topics []string) error {
	if userID == "" || postID == "" {
		return fmt.Errorf("用户ID和动态ID不能为空")
	}
	if content == "" {
		return fmt.Errorf("内容不能为空")
	}
	if len(content) > 5000 {
		return fmt.Errorf("内容最多5000字")
	}

	// 获取动态
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("获取动态失败: %w", err)
	}

	// 权限检查
	if post.UserID != userID {
		return fmt.Errorf("无权更新该动态")
	}

	// 构建更新
	updates := map[string]interface{}{
		"content": content,
		"topics":  topics,
	}

	if err := s.postRepo.Update(ctx, postID, updates); err != nil {
		return fmt.Errorf("更新动态失败: %w", err)
	}

	return nil
}

// DeletePost 删除动态
func (s *PostService) DeletePost(ctx context.Context, userID, postID string) error {
	if userID == "" || postID == "" {
		return fmt.Errorf("用户ID和动态ID不能为空")
	}

	// 获取动态
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("获取动态失败: %w", err)
	}

	// 权限检查
	if post.UserID != userID {
		return fmt.Errorf("无权删除该动态")
	}

	if err := s.postRepo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("删除动态失败: %w", err)
	}

	s.publishPostEvent(ctx, "post.deleted", userID, postID)
	return nil
}

// GetUserPosts 获取用户发布的动态列表
func (s *PostService) GetUserPosts(ctx context.Context, userID string, page, size int) ([]*socialModel.PostInfo, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	posts, total, err := s.postRepo.GetByUser(ctx, userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户动态列表失败: %w", err)
	}

	// 转换为 PostInfo（不包含点赞状态）
	postInfos := make([]*socialModel.PostInfo, len(posts))
	for i, post := range posts {
		postInfos[i] = post.ToPostInfo(false, false)
	}

	return postInfos, total, nil
}

// ToggleLike 切换点赞状态
func (s *PostService) ToggleLike(ctx context.Context, userID, postID string) (bool, error) {
	if userID == "" || postID == "" {
		return false, fmt.Errorf("用户ID和动态ID不能为空")
	}

	// 检查是否已点赞
	isLiked, err := s.postRepo.IsLiked(ctx, postID, userID)
	if err != nil {
		return false, fmt.Errorf("检查点赞状态失败: %w", err)
	}

	if isLiked {
		// 已点赞，取消点赞
		if err := s.Unlike(ctx, userID, postID); err != nil {
			return false, err
		}
		return false, nil
	} else {
		// 未点赞，执行点赞
		if err := s.Like(ctx, userID, postID); err != nil {
			return false, err
		}
		return true, nil
	}
}

// Like 点赞动态
func (s *PostService) Like(ctx context.Context, userID, postID string) error {
	if userID == "" || postID == "" {
		return fmt.Errorf("用户ID和动态ID不能为空")
	}

	// 检查是否已点赞
	isLiked, err := s.postRepo.IsLiked(ctx, postID, userID)
	if err != nil {
		return fmt.Errorf("检查点赞状态失败: %w", err)
	}
	if isLiked {
		return fmt.Errorf("已经点赞过该动态")
	}

	// 事务：创建点赞记录 + 增加点赞数
	if err := s.postRepo.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.Like(txCtx, postID, userID); err != nil {
			return fmt.Errorf("点赞失败: %w", err)
		}
		if err := s.postRepo.IncrementLikeCount(txCtx, postID); err != nil {
			return fmt.Errorf("增加点赞数失败: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Unlike 取消点赞动态
func (s *PostService) Unlike(ctx context.Context, userID, postID string) error {
	if userID == "" || postID == "" {
		return fmt.Errorf("用户ID和动态ID不能为空")
	}

	// 检查是否已点赞
	isLiked, err := s.postRepo.IsLiked(ctx, postID, userID)
	if err != nil {
		return fmt.Errorf("检查点赞状态失败: %w", err)
	}
	if !isLiked {
		return fmt.Errorf("未点赞该动态")
	}

	// 事务：删除点赞记录 + 减少点赞数
	if err := s.postRepo.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.Unlike(txCtx, postID, userID); err != nil {
			return fmt.Errorf("取消点赞失败: %w", err)
		}
		if err := s.postRepo.DecrementLikeCount(txCtx, postID); err != nil {
			return fmt.Errorf("减少点赞数失败: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *PostService) publishPostEvent(ctx context.Context, eventType, userID, postID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"user_id": userID,
			"post_id": postID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}
