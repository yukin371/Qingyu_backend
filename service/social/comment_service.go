package social

import (
	"Qingyu_backend/models/social"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	auditRepo "Qingyu_backend/repository/interfaces/audit"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/service/base"
)

// CommentService 评论服务
type CommentService struct {
	commentRepo       socialRepo.CommentRepository
	sensitiveWordRepo auditRepo.SensitiveWordRepository
	eventBus          base.EventBus
	serviceName       string
	version           string
}

// NewCommentService 创建评论服务实例
func NewCommentService(
	commentRepo socialRepo.CommentRepository,
	sensitiveWordRepo auditRepo.SensitiveWordRepository,
	eventBus base.EventBus,
) *CommentService {
	return &CommentService{
		commentRepo:       commentRepo,
		sensitiveWordRepo: sensitiveWordRepo,
		eventBus:          eventBus,
		serviceName:       "CommentService",
		version:           "1.0.0",
	}
}

// =========================
// BaseService 接口实现
// =========================

// Initialize 初始化服务
func (s *CommentService) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *CommentService) Health(ctx context.Context) error {
	if err := s.commentRepo.Health(ctx); err != nil {
		return fmt.Errorf("评论Repository健康检查失败: %w", err)
	}
	return nil
}

// Close 关闭服务
func (s *CommentService) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *CommentService) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *CommentService) GetVersion() string {
	return s.version
}

// =========================
// 评论管理方法
// =========================

// PublishComment 发表评论
func (s *CommentService) PublishComment(ctx context.Context, userID, bookID, chapterID, content string, rating int) (*social.Comment, error) {
	// 参数验证
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if bookID == "" {
		return nil, fmt.Errorf("书籍ID不能为空")
	}
	if len(content) < 10 || len(content) > 500 {
		return nil, fmt.Errorf("评论内容长度必须在10-500字之间")
	}
	if rating < 0 || rating > 5 {
		return nil, fmt.Errorf("评分必须在0-5之间")
	}

	// 创建评论对象
	comment := &social.Comment{
		TargetType: social.CommentTargetTypeBook,
		TargetID:   bookID,
		AuthorID:   userID,
		Content:    strings.TrimSpace(content),
		Rating:     rating,
		State:      social.CommentStateNormal,
	}

	// 兼容旧字段（可选）
	if chapterID != "" {
		comment.ChapterID = chapterID
	}

	// 自动审核
	state, reason, err := s.AutoReviewComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("评论审核失败: %w", err)
	}

	comment.State = state
	if reason != "" {
		comment.RejectReason = reason
	}

	// 保存评论
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("保存评论失败: %w", err)
	}

	// 发布事件
	s.publishCommentEvent(ctx, "comment.created", comment)

	return comment, nil
}

// ReplyComment 回复评论
func (s *CommentService) ReplyComment(ctx context.Context, userID, parentCommentID, content string) (*social.Comment, error) {
	// 参数验证
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if parentCommentID == "" {
		return nil, fmt.Errorf("父评论ID不能为空")
	}
	if len(content) < 10 || len(content) > 500 {
		return nil, fmt.Errorf("回复内容长度必须在10-500字之间")
	}

	// 获取父评论
	parentComment, err := s.commentRepo.GetByID(ctx, parentCommentID)
	if err != nil {
		return nil, fmt.Errorf("父评论不存在: %w", err)
	}

	// 检查父评论是否已删除
	if parentComment.State == social.CommentStateDeleted {
		return nil, fmt.Errorf("无法回复已删除的评论")
	}

	// 创建回复评论
	parentID := parentCommentID
	rootID := parentComment.RootID
	if rootID == nil {
		rootID = &parentCommentID
	}
	replyToUserID := parentComment.AuthorID

	comment := &social.Comment{
		TargetType: parentComment.TargetType,
		TargetID:   parentComment.TargetID,
		AuthorID:   userID,
		Content:    strings.TrimSpace(content),
		Rating:     0, // 回复不包含评分
		State:      social.CommentStateNormal,
	}
	// 设置嵌入字段的字段
	comment.ParentID = &parentID
	comment.RootID = rootID
	comment.ReplyToUserID = &replyToUserID

	// 自动审核
	state, reason, err := s.AutoReviewComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("回复审核失败: %w", err)
	}

	comment.State = state
	if reason != "" {
		comment.RejectReason = reason
	}

	// 保存回复
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("保存回复失败: %w", err)
	}

	// 增加父评论的回复数
	if comment.State == social.CommentStateNormal {
		if err := s.commentRepo.IncrementReplyCount(ctx, parentCommentID); err != nil {
			// 非致命错误，只记录日志
			log.Printf("Warning: failed to increment reply count")
		}
	}

	// 发布事件
	s.publishCommentEvent(ctx, "comment.replied", comment)

	return comment, nil
}

// GetCommentList 获取评论列表
func (s *CommentService) GetCommentList(ctx context.Context, bookID string, sortBy string, page, size int) ([]*social.Comment, int64, error) {
	if bookID == "" {
		return nil, 0, fmt.Errorf("书籍ID不能为空")
	}

	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 验证排序方式
	if sortBy != social.CommentSortByLatest && sortBy != social.CommentSortByHot {
		sortBy = social.CommentSortByLatest
	}

	// 获取评论列表
	comments, total, err := s.commentRepo.GetCommentsByBookIDSorted(ctx, bookID, sortBy, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取评论列表失败: %w", err)
	}

	return comments, total, nil
}

// GetCommentDetail 获取评论详情
func (s *CommentService) GetCommentDetail(ctx context.Context, commentID string) (*social.Comment, error) {
	if commentID == "" {
		return nil, fmt.Errorf("评论ID不能为空")
	}

	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("获取评论详情失败: %w", err)
	}

	return comment, nil
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(ctx context.Context, userID, commentID, content string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if commentID == "" {
		return fmt.Errorf("评论ID不能为空")
	}
	if len(content) < 10 || len(content) > 500 {
		return fmt.Errorf("评论内容长度必须在10-500字之间")
	}

	// 获取原评论
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("评论不存在: %w", err)
	}

	// 权限检查
	if comment.AuthorID != userID {
		return fmt.Errorf("没有权限修改此评论")
	}

	// 检查是否可编辑（30分钟内）
	if !comment.IsEditable() {
		return fmt.Errorf("评论发表超过30分钟，无法编辑")
	}

	// 敏感词检测
	if s.sensitiveWordRepo != nil {
		hasSensitive, _, err := s.checkSensitiveWords(ctx, content)
		if err != nil {
			return fmt.Errorf("敏感词检测失败: %w", err)
		}
		if hasSensitive {
			return fmt.Errorf("评论包含敏感词，无法发布")
		}
	}

	// 更新评论
	updates := map[string]interface{}{
		"content": strings.TrimSpace(content),
	}

	if err := s.commentRepo.Update(ctx, commentID, updates); err != nil {
		return fmt.Errorf("更新评论失败: %w", err)
	}

	return nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(ctx context.Context, userID, commentID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if commentID == "" {
		return fmt.Errorf("评论ID不能为空")
	}

	// 获取评论
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("评论不存在: %w", err)
	}

	// 权限检查
	if comment.AuthorID != userID {
		return fmt.Errorf("没有权限删除此评论")
	}

	// 删除评论
	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("删除评论失败: %w", err)
	}

	// 如果是回复，减少父评论的回复数
	if comment.ParentID != nil && *comment.ParentID != "" {
		if err := s.commentRepo.DecrementReplyCount(ctx, *comment.ParentID); err != nil {
			fmt.Printf("Warning: Failed to decrement reply count: %v\n", err)
		}
	}

	// 发布事件
	s.publishCommentEvent(ctx, "comment.deleted", comment)

	return nil
}

// =========================
// 点赞方法（与LikeService集成）
// =========================

// LikeComment 点赞评论（委托给LikeService）
func (s *CommentService) LikeComment(ctx context.Context, userID, commentID string) error {
	// 此方法将由LikeService调用Repository的IncrementLikeCount
	// 这里仅作为占位符，实际实现在LikeService中
	if err := s.commentRepo.IncrementLikeCount(ctx, commentID); err != nil {
		return fmt.Errorf("点赞评论失败: %w", err)
	}
	return nil
}

// UnlikeComment 取消点赞评论（委托给LikeService）
func (s *CommentService) UnlikeComment(ctx context.Context, userID, commentID string) error {
	// 此方法将由LikeService调用Repository的DecrementLikeCount
	// 这里仅作为占位符，实际实现在LikeService中
	if err := s.commentRepo.DecrementLikeCount(ctx, commentID); err != nil {
		return fmt.Errorf("取消点赞评论失败: %w", err)
	}
	return nil
}

// =========================
// 审核方法
// =========================

// AutoReviewComment 自动审核评论
func (s *CommentService) AutoReviewComment(ctx context.Context, comment *social.Comment) (social.CommentState, string, error) {
	// 如果敏感词库未配置，默认通过
	if s.sensitiveWordRepo == nil {
		return social.CommentStateNormal, "", nil
	}

	// 检查敏感词
	hasSensitive, sensitiveWords, err := s.checkSensitiveWords(ctx, comment.Content)
	if err != nil {
		return social.CommentStateNormal, "", fmt.Errorf("敏感词检测失败: %w", err)
	}

	if hasSensitive {
		reason := fmt.Sprintf("评论包含敏感词: %s", strings.Join(sensitiveWords, ", "))
		return social.CommentStateRejected, reason, nil
	}

	// 通过审核
	return social.CommentStateNormal, "", nil
}

// =========================
// 统计方法
// =========================

// GetBookCommentStats 获取书籍评论统计
func (s *CommentService) GetBookCommentStats(ctx context.Context, bookID string) (map[string]interface{}, error) {
	if bookID == "" {
		return nil, fmt.Errorf("书籍ID不能为空")
	}

	// 获取评分统计
	ratingStats, err := s.commentRepo.GetBookRatingStats(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取评分统计失败: %w", err)
	}

	// 获取评论总数
	commentCount, err := s.commentRepo.GetCommentCount(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取评论总数失败: %w", err)
	}

	stats := map[string]interface{}{
		"comment_count": commentCount,
		"rating_stats":  ratingStats,
	}

	return stats, nil
}

// GetUserComments 获取用户评论列表
func (s *CommentService) GetUserComments(ctx context.Context, userID string, page, size int) ([]*social.Comment, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	comments, total, err := s.commentRepo.GetCommentsByUserID(ctx, userID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户评论列表失败: %w", err)
	}

	return comments, total, nil
}

// =========================
// 私有辅助方法
// =========================

// checkSensitiveWords 检查敏感词
func (s *CommentService) checkSensitiveWords(ctx context.Context, content string) (bool, []string, error) {
	if s.sensitiveWordRepo == nil {
		return false, nil, nil
	}

	// 获取启用的敏感词列表
	enabledWords, err := s.sensitiveWordRepo.GetEnabledWords(ctx)
	if err != nil {
		return false, nil, fmt.Errorf("获取敏感词列表失败: %w", err)
	}

	// 检查内容中是否包含敏感词
	var foundWords []string
	contentLower := strings.ToLower(content)

	for _, word := range enabledWords {
		wordLower := strings.ToLower(word.Word)
		if strings.Contains(contentLower, wordLower) {
			foundWords = append(foundWords, word.Word)
		}
	}

	return len(foundWords) > 0, foundWords, nil
}

// getRootID 获取根评论ID
func (s *CommentService) getRootID(comment *social.Comment) *string {
	if comment.RootID != nil {
		return comment.RootID
	}
	if comment.ParentID != nil {
		return comment.ParentID
	}
	idHex := comment.ID.Hex()
	return &idHex
}

// publishCommentEvent 发布评论事件
func (s *CommentService) publishCommentEvent(ctx context.Context, eventType string, comment *social.Comment) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"comment_id":  comment.ID,
			"target_id":   comment.TargetID,
			"target_type": comment.TargetType,
			"author_id":   comment.AuthorID,
			"state":       comment.State,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}

// =========================
// 高级查询方法
// =========================

// GetCommentThread 获取评论线程（树状结构）
func (s *CommentService) GetCommentThread(ctx context.Context, commentID string) (*social.CommentThread, error) {
	if commentID == "" {
		return nil, fmt.Errorf("评论ID不能为空")
	}

	// 获取主评论
	comment, err := s.GetCommentDetail(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("获取评论详情失败: %w", err)
	}

	// 获取所有回复（不分页，获取全部）
	replies, total, err := s.GetCommentReplies(ctx, commentID, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("获取评论回复失败: %w", err)
	}

	// 构建回复线程（简化版，只支持一级回复）
	replyThreads := make([]*social.CommentThread, len(replies))
	for i, reply := range replies {
		replyThreads[i] = &social.CommentThread{
			Comment: reply,
			Replies: nil, // 简化版不递归获取嵌套回复
			Total:   0,
			HasMore: false,
		}
	}

	thread := &social.CommentThread{
		Comment: comment,
		Replies: replyThreads,
		Total:   total,
		HasMore: total > 100, // 如果超过100条回复，表示还有更多
	}

	return thread, nil
}

// GetTopComments 获取热门评论（按点赞数排序）
func (s *CommentService) GetTopComments(ctx context.Context, bookID string, limit int) ([]*social.Comment, error) {
	if bookID == "" {
		return nil, fmt.Errorf("书籍ID不能为空")
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// 使用热度排序获取评论列表
	comments, _, err := s.GetCommentList(ctx, bookID, social.CommentSortByHot, 1, limit)
	if err != nil {
		return nil, fmt.Errorf("获取热门评论失败: %w", err)
	}

	return comments, nil
}

// GetCommentReplies 获取评论的所有回复
func (s *CommentService) GetCommentReplies(ctx context.Context, commentID string, page, size int) ([]*social.Comment, int64, error) {
	if commentID == "" {
		return nil, 0, fmt.Errorf("评论ID不能为空")
	}

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 获取该评论的所有回复（Repository 层暂不支持分页，返回所有回复）
	allReplies, err := s.commentRepo.GetRepliesByCommentID(ctx, commentID)
	if err != nil {
		return nil, 0, fmt.Errorf("获取评论回复失败: %w", err)
	}

	total := int64(len(allReplies))

	// 在 Service 层实现分页
	start := (page - 1) * size
	end := start + size

	if start >= len(allReplies) {
		return []*social.Comment{}, total, nil
	}

	if end > len(allReplies) {
		end = len(allReplies)
	}

	replies := allReplies[start:end]

	return replies, total, nil
}
