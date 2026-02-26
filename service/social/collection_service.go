package social

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"Qingyu_backend/models/reader"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	"Qingyu_backend/service/base"
)

// CollectionService 收藏服务
type CollectionService struct {
	collectionRepo socialRepo.CollectionRepository
	eventBus       base.EventBus
	serviceName    string
	version        string
}

// NewCollectionService 创建收藏服务实例
func NewCollectionService(
	collectionRepo socialRepo.CollectionRepository,
	eventBus base.EventBus,
) *CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
		eventBus:       eventBus,
		serviceName:    "CollectionService",
		version:        "1.0.0",
	}
}

// =========================
// BaseService 接口实现
// =========================

// Initialize 初始化服务
func (s *CollectionService) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *CollectionService) Health(ctx context.Context) error {
	if err := s.collectionRepo.Health(ctx); err != nil {
		return fmt.Errorf("收藏Repository健康检查失败: %w", err)
	}
	return nil
}

// Close 关闭服务
func (s *CollectionService) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *CollectionService) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *CollectionService) GetVersion() string {
	return s.version
}

// =========================
// 收藏管理
// =========================

// AddToCollection 添加收藏
func (s *CollectionService) AddToCollection(ctx context.Context, userID, bookID, folderID, note string, tags []string, isPublic bool) (*reader.Collection, error) {
	// 参数验证
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if bookID == "" {
		return nil, fmt.Errorf("书籍ID不能为空")
	}

	// 验证笔记长度
	if len(note) > 500 {
		return nil, fmt.Errorf("笔记最多500字")
	}

	// 验证标签
	if len(tags) > 10 {
		return nil, fmt.Errorf("最多10个标签")
	}
	for _, tag := range tags {
		if len(tag) > 20 {
			return nil, fmt.Errorf("标签最多20字")
		}
	}

	// 检查是否已收藏
	existing, err := s.collectionRepo.GetByUserAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("检查收藏失败: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("已经收藏过该书籍")
	}

	// 创建收藏
	collection := &reader.Collection{
		UserID:    userID,
		BookID:    bookID,
		FolderID:  folderID,
		Note:      note,
		Tags:      tags,
		IsPublic:  isPublic,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.collectionRepo.Create(ctx, collection); err != nil {
		return nil, fmt.Errorf("添加收藏失败: %w", err)
	}

	// 更新收藏夹计数
	if folderID != "" {
		if err := s.collectionRepo.IncrementFolderBookCount(ctx, folderID); err != nil {
			fmt.Printf("Warning: Failed to increment folder book count: %v\n", err)
		}
	}

	// 发布事件
	s.publishCollectionEvent(ctx, "collection.added", userID, bookID, collection.ID.Hex())

	return collection, nil
}

// RemoveFromCollection 取消收藏
func (s *CollectionService) RemoveFromCollection(ctx context.Context, userID, collectionID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if collectionID == "" {
		return fmt.Errorf("收藏ID不能为空")
	}

	// 获取收藏记录
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("获取收藏失败: %w", err)
	}

	// 权限检查
	if collection.UserID != userID {
		return fmt.Errorf("无权删除该收藏")
	}

	// 删除收藏
	if err := s.collectionRepo.Delete(ctx, collectionID); err != nil {
		return fmt.Errorf("删除收藏失败: %w", err)
	}

	// 更新收藏夹计数
	if collection.FolderID != "" {
		if err := s.collectionRepo.DecrementFolderBookCount(ctx, collection.FolderID); err != nil {
			fmt.Printf("Warning: Failed to decrement folder book count: %v\n", err)
		}
	}

	// 发布事件
	s.publishCollectionEvent(ctx, "collection.removed", userID, collection.BookID, collectionID)

	return nil
}

// UpdateCollection 更新收藏
func (s *CollectionService) UpdateCollection(ctx context.Context, userID, collectionID string, updates map[string]interface{}) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if collectionID == "" {
		return fmt.Errorf("收藏ID不能为空")
	}

	// 获取收藏记录
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("获取收藏失败: %w", err)
	}

	// 权限检查
	if collection.UserID != userID {
		return fmt.Errorf("无权更新该收藏")
	}

	// 参数验证
	if note, ok := updates["note"].(string); ok {
		if len(note) > 500 {
			return fmt.Errorf("笔记最多500字")
		}
	}
	if tags, ok := updates["tags"].([]string); ok {
		if len(tags) > 10 {
			return fmt.Errorf("最多10个标签")
		}
		for _, tag := range tags {
			if len(tag) > 20 {
				return fmt.Errorf("标签最多20字")
			}
		}
	}

	// 处理收藏夹变更
	if newFolderID, ok := updates["folder_id"].(string); ok {
		oldFolderID := collection.FolderID
		if oldFolderID != newFolderID {
			// 旧收藏夹计数-1
			if oldFolderID != "" {
				if err := s.collectionRepo.DecrementFolderBookCount(ctx, oldFolderID); err != nil {
					fmt.Printf("Warning: Failed to decrement old folder count: %v\n", err)
				}
			}
			// 新收藏夹计数+1
			if newFolderID != "" {
				if err := s.collectionRepo.IncrementFolderBookCount(ctx, newFolderID); err != nil {
					fmt.Printf("Warning: Failed to increment new folder count: %v\n", err)
				}
			}
		}
	}

	// 更新收藏
	if err := s.collectionRepo.Update(ctx, collectionID, updates); err != nil {
		return fmt.Errorf("更新收藏失败: %w", err)
	}

	return nil
}

// GetUserCollections 获取用户收藏列表
func (s *CollectionService) GetUserCollections(ctx context.Context, userID, folderID string, page, size int) ([]*reader.Collection, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}

	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	collections, total, err := s.collectionRepo.GetCollectionsByUser(ctx, userID, folderID, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取收藏列表失败: %w", err)
	}

	return collections, total, nil
}

// GetCollectionsByTag 根据标签获取收藏
func (s *CollectionService) GetCollectionsByTag(ctx context.Context, userID, tag string, page, size int) ([]*reader.Collection, int64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("用户ID不能为空")
	}
	if tag == "" {
		return nil, 0, fmt.Errorf("标签不能为空")
	}

	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	collections, total, err := s.collectionRepo.GetCollectionsByTag(ctx, userID, tag, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("根据标签获取收藏失败: %w", err)
	}

	return collections, total, nil
}

// IsCollected 检查是否已收藏
func (s *CollectionService) IsCollected(ctx context.Context, userID, bookID string) (bool, error) {
	if userID == "" || bookID == "" {
		return false, fmt.Errorf("用户ID和书籍ID不能为空")
	}

	collection, err := s.collectionRepo.GetByUserAndBook(ctx, userID, bookID)
	if err != nil {
		return false, fmt.Errorf("检查收藏失败: %w", err)
	}

	return collection != nil, nil
}

// =========================
// 收藏夹管理
// =========================

// CreateFolder 创建收藏夹
func (s *CollectionService) CreateFolder(ctx context.Context, userID, name, description string, isPublic bool) (*reader.CollectionFolder, error) {
	// 参数验证
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if name == "" {
		return nil, fmt.Errorf("收藏夹名称不能为空")
	}
	if len(name) > 50 {
		return nil, fmt.Errorf("收藏夹名称最多50字")
	}
	if len(description) > 200 {
		return nil, fmt.Errorf("收藏夹描述最多200字")
	}

	// 创建收藏夹
	folder := &reader.CollectionFolder{
		UserID:      userID,
		Name:        name,
		Description: description,
		IsPublic:    isPublic,
		BookCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.collectionRepo.CreateFolder(ctx, folder); err != nil {
		return nil, fmt.Errorf("创建收藏夹失败: %w", err)
	}

	// 发布事件
	s.publishFolderEvent(ctx, "folder.created", userID, folder.ID.Hex())

	return folder, nil
}

// GetUserFolders 获取用户收藏夹列表
func (s *CollectionService) GetUserFolders(ctx context.Context, userID string) ([]*reader.CollectionFolder, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	folders, err := s.collectionRepo.GetFoldersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取收藏夹列表失败: %w", err)
	}

	return folders, nil
}

// UpdateFolder 更新收藏夹
func (s *CollectionService) UpdateFolder(ctx context.Context, userID, folderID string, updates map[string]interface{}) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if folderID == "" {
		return fmt.Errorf("收藏夹ID不能为空")
	}

	// 获取收藏夹
	folder, err := s.collectionRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		return fmt.Errorf("获取收藏夹失败: %w", err)
	}

	// 权限检查
	if folder.UserID != userID {
		return fmt.Errorf("无权更新该收藏夹")
	}

	// 参数验证
	if name, ok := updates["name"].(string); ok {
		if name == "" {
			return fmt.Errorf("收藏夹名称不能为空")
		}
		if len(name) > 50 {
			return fmt.Errorf("收藏夹名称最多50字")
		}
	}
	if desc, ok := updates["description"].(string); ok {
		if len(desc) > 200 {
			return fmt.Errorf("收藏夹描述最多200字")
		}
	}

	// 更新收藏夹
	if err := s.collectionRepo.UpdateFolder(ctx, folderID, updates); err != nil {
		return fmt.Errorf("更新收藏夹失败: %w", err)
	}

	return nil
}

// DeleteFolder 删除收藏夹
func (s *CollectionService) DeleteFolder(ctx context.Context, userID, folderID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if folderID == "" {
		return fmt.Errorf("收藏夹ID不能为空")
	}

	// 获取收藏夹
	folder, err := s.collectionRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		return fmt.Errorf("获取收藏夹失败: %w", err)
	}

	// 权限检查
	if folder.UserID != userID {
		return fmt.Errorf("无权删除该收藏夹")
	}

	// 检查是否为空
	if folder.BookCount > 0 {
		return fmt.Errorf("收藏夹不为空，请先移除收藏")
	}

	// 删除收藏夹
	if err := s.collectionRepo.DeleteFolder(ctx, folderID); err != nil {
		return fmt.Errorf("删除收藏夹失败: %w", err)
	}

	// 发布事件
	s.publishFolderEvent(ctx, "folder.deleted", userID, folderID)

	return nil
}

// =========================
// 收藏分享
// =========================

// ShareCollection 分享收藏
func (s *CollectionService) ShareCollection(ctx context.Context, userID, collectionID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if collectionID == "" {
		return fmt.Errorf("收藏ID不能为空")
	}

	// 获取收藏
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("获取收藏失败: %w", err)
	}

	// 权限检查
	if collection.UserID != userID {
		return fmt.Errorf("无权分享该收藏")
	}

	// 更新为公开
	updates := map[string]interface{}{
		"is_public": true,
	}

	if err := s.collectionRepo.Update(ctx, collectionID, updates); err != nil {
		return fmt.Errorf("分享收藏失败: %w", err)
	}

	return nil
}

// UnshareCollection 取消分享
func (s *CollectionService) UnshareCollection(ctx context.Context, userID, collectionID string) error {
	if userID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if collectionID == "" {
		return fmt.Errorf("收藏ID不能为空")
	}

	// 获取收藏
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("获取收藏失败: %w", err)
	}

	// 权限检查
	if collection.UserID != userID {
		return fmt.Errorf("无权取消分享该收藏")
	}

	// 更新为私有
	updates := map[string]interface{}{
		"is_public": false,
	}

	if err := s.collectionRepo.Update(ctx, collectionID, updates); err != nil {
		return fmt.Errorf("取消分享失败: %w", err)
	}

	return nil
}

// GetPublicCollections 获取公开收藏列表
func (s *CollectionService) GetPublicCollections(ctx context.Context, page, size int) ([]*reader.Collection, int64, error) {
	// 参数验证和默认值
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	collections, total, err := s.collectionRepo.GetPublicCollections(ctx, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("获取公开收藏列表失败: %w", err)
	}

	return collections, total, nil
}

// =========================
// 统计
// =========================

// GetUserCollectionStats 获取用户收藏统计
func (s *CollectionService) GetUserCollectionStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 统计收藏数
	totalCount, err := s.collectionRepo.CountUserCollections(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取收藏统计失败: %w", err)
	}

	// 获取收藏夹列表
	folders, err := s.collectionRepo.GetFoldersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取收藏夹列表失败: %w", err)
	}

	stats := map[string]interface{}{
		"total_collections": totalCount,
		"total_folders":     len(folders),
	}

	return stats, nil
}

// =========================
// 私有辅助方法
// =========================

// publishCollectionEvent 发布收藏事件
func (s *CollectionService) publishCollectionEvent(ctx context.Context, eventType string, userID, bookID, collectionID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"user_id":       userID,
			"book_id":       bookID,
			"collection_id": collectionID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}

// publishFolderEvent 发布收藏夹事件
func (s *CollectionService) publishFolderEvent(ctx context.Context, eventType string, userID, folderID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: eventType,
		EventData: map[string]interface{}{
			"user_id":   userID,
			"folder_id": folderID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}

// =========================
// 分享链接生成（新增）
// =========================

// generateShareID 生成唯一的分享ID
func (s *CollectionService) generateShareID() (string, error) {
	// 生成9位随机字符串（小写字母+数字）
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 9)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("生成随机数失败: %w", err)
		}
		result[i] = chars[num.Int64()]
	}
	return string(result), nil
}

// ShareCollectionWithURL 分享收藏并返回分享链接
func (s *CollectionService) ShareCollectionWithURL(ctx context.Context, userID, collectionID string) (map[string]interface{}, error) {
	if userID == "" {
		return nil, fmt.Errorf("用户ID不能为空")
	}
	if collectionID == "" {
		return nil, fmt.Errorf("收藏ID不能为空")
	}

	// 获取收藏
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("获取收藏失败: %w", err)
	}

	// 权限检查
	if collection.UserID != userID {
		return nil, fmt.Errorf("无权分享该收藏")
	}

	// 如果已有分享ID，直接返回
	if collection.ShareID != "" {
		return map[string]interface{}{
			"share_id":   collection.ShareID,
			"share_url":  "/api/v1/reader/collections/shared/" + collection.ShareID,
			"expires_at": nil,
		}, nil
	}

	// 生成新的分享ID
	shareID, err := s.generateShareID()
	if err != nil {
		return nil, fmt.Errorf("生成分享ID失败: %w", err)
	}

	// 更新收藏
	updates := map[string]interface{}{
		"is_public": true,
		"share_id":  shareID,
	}

	if err := s.collectionRepo.Update(ctx, collectionID, updates); err != nil {
		return nil, fmt.Errorf("分享收藏失败: %w", err)
	}

	return map[string]interface{}{
		"share_id":   shareID,
		"share_url":  "/api/v1/reader/collections/shared/" + shareID,
		"expires_at": nil,
	}, nil
}

// GetSharedCollection 根据分享ID获取收藏详情
func (s *CollectionService) GetSharedCollection(ctx context.Context, shareID string) (*reader.Collection, error) {
	if shareID == "" {
		return nil, fmt.Errorf("分享ID不能为空")
	}

	collection, err := s.collectionRepo.GetByShareID(ctx, shareID)
	if err != nil {
		return nil, fmt.Errorf("获取分享收藏失败: %w", err)
	}

	// 检查是否为公开收藏
	if !collection.IsPublic {
		return nil, fmt.Errorf("该收藏未公开")
	}

	return collection, nil
}
