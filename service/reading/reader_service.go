package reading

import (
	reader2 "Qingyu_backend/models/reader"
	"context"
	"fmt"
	"time"

	readingRepo "Qingyu_backend/repository/interfaces/reading"
	"Qingyu_backend/service/base"
)

// ReaderService 阅读器服务
type ReaderService struct {
	chapterRepo    readingRepo.ChapterRepository
	progressRepo   readingRepo.ReadingProgressRepository
	annotationRepo readingRepo.AnnotationRepository
	settingsRepo   readingRepo.ReadingSettingsRepository
	eventBus       base.EventBus
	cacheService   ReaderCacheService
	vipService     VIPPermissionService
	serviceName    string
	version        string
}

// NewReaderService 创建阅读器服务实例
func NewReaderService(
	chapterRepo readingRepo.ChapterRepository,
	progressRepo readingRepo.ReadingProgressRepository,
	annotationRepo readingRepo.AnnotationRepository,
	settingsRepo readingRepo.ReadingSettingsRepository,
	eventBus base.EventBus,
	cacheService ReaderCacheService,
	vipService VIPPermissionService,
) *ReaderService {
	return &ReaderService{
		chapterRepo:    chapterRepo,
		progressRepo:   progressRepo,
		annotationRepo: annotationRepo,
		settingsRepo:   settingsRepo,
		eventBus:       eventBus,
		cacheService:   cacheService,
		vipService:     vipService,
		serviceName:    "ReaderService",
		version:        "1.0.0",
	}
}

// =========================
// BaseService 接口实现
// =========================

// Initialize 初始化服务
func (s *ReaderService) Initialize(ctx context.Context) error {
	return nil
}

// Health 健康检查
func (s *ReaderService) Health(ctx context.Context) error {
	if err := s.chapterRepo.Health(ctx); err != nil {
		return fmt.Errorf("章节Repository健康检查失败: %w", err)
	}
	if err := s.progressRepo.Health(ctx); err != nil {
		return fmt.Errorf("进度Repository健康检查失败: %w", err)
	}
	if err := s.annotationRepo.Health(ctx); err != nil {
		return fmt.Errorf("标注Repository健康检查失败: %w", err)
	}
	return nil
}

// Close 关闭服务
func (s *ReaderService) Close(ctx context.Context) error {
	return nil
}

// GetServiceName 获取服务名称
func (s *ReaderService) GetServiceName() string {
	return s.serviceName
}

// GetVersion 获取服务版本
func (s *ReaderService) GetVersion() string {
	return s.version
}

// =========================
// 章节相关方法
// =========================

// GetChapterByID 根据ID获取章节
func (s *ReaderService) GetChapterByID(ctx context.Context, chapterID string) (*reader2.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节失败: %w", err)
	}
	return chapter, nil
}

// GetChapterByNum 根据章节号获取章节
func (s *ReaderService) GetChapterByNum(ctx context.Context, bookID string, chapterNum int) (*reader2.Chapter, error) {
	chapter, err := s.chapterRepo.GetByChapterNum(ctx, bookID, chapterNum)
	if err != nil {
		return nil, fmt.Errorf("获取章节失败: %w", err)
	}
	return chapter, nil
}

// GetBookChapters 获取书籍的所有章节
func (s *ReaderService) GetBookChapters(ctx context.Context, bookID string) ([]*reader2.Chapter, error) {
	chapters, err := s.chapterRepo.GetByBookID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}
	return chapters, nil
}

// GetBookChaptersWithPagination 分页获取书籍章节
func (s *ReaderService) GetBookChaptersWithPagination(ctx context.Context, bookID string, page, size int) ([]*reader2.Chapter, int64, error) {
	offset := int64((page - 1) * size)
	limit := int64(size)

	chapters, err := s.chapterRepo.GetByBookIDWithPagination(ctx, bookID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取章节列表失败: %w", err)
	}

	total, err := s.chapterRepo.CountByBookID(ctx, bookID)
	if err != nil {
		return nil, 0, fmt.Errorf("统计章节数失败: %w", err)
	}

	return chapters, total, nil
}

// GetPrevChapter 获取上一章
func (s *ReaderService) GetPrevChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader2.Chapter, error) {
	chapter, err := s.chapterRepo.GetPrevChapter(ctx, bookID, currentChapterNum)
	if err != nil {
		return nil, fmt.Errorf("获取上一章失败: %w", err)
	}
	return chapter, nil
}

// GetNextChapter 获取下一章
func (s *ReaderService) GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader2.Chapter, error) {
	chapter, err := s.chapterRepo.GetNextChapter(ctx, bookID, currentChapterNum)
	if err != nil {
		return nil, fmt.Errorf("获取下一章失败: %w", err)
	}
	return chapter, nil
}

// GetChapterContent 获取章节内容
func (s *ReaderService) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	// 1. 尝试从缓存获取章节内容
	if s.cacheService != nil {
		cachedContent, err := s.cacheService.GetChapterContent(ctx, chapterID)
		if err == nil && cachedContent != "" {
			// 缓存命中，仍需验证VIP权限
			isVIP, _ := s.chapterRepo.CheckVIPAccess(ctx, chapterID)
			if isVIP {
				hasAccess, err := s.vipService.CheckVIPAccess(ctx, userID, chapterID, true)
				if err != nil {
					return "", fmt.Errorf("检查VIP权限失败: %w", err)
				}
				if !hasAccess {
					return "", fmt.Errorf("该章节为VIP章节，需要VIP权限或购买后才能阅读")
				}
			}

			// 发布阅读事件
			s.publishReadingEvent(ctx, userID, chapterID)
			return cachedContent, nil
		}
	}

	// 2. 检查VIP权限
	isVIP, err := s.chapterRepo.CheckVIPAccess(ctx, chapterID)
	if err != nil {
		return "", fmt.Errorf("检查VIP权限失败: %w", err)
	}

	if isVIP {
		// 验证用户是否有VIP权限或已购买该章节
		if s.vipService != nil {
			hasAccess, err := s.vipService.CheckVIPAccess(ctx, userID, chapterID, true)
			if err != nil {
				return "", fmt.Errorf("检查VIP权限失败: %w", err)
			}
			if !hasAccess {
				return "", fmt.Errorf("该章节为VIP章节，需要VIP权限或购买后才能阅读")
			}
		}
	}

	// 3. 从数据库获取章节内容
	content, err := s.chapterRepo.GetChapterContent(ctx, chapterID)
	if err != nil {
		return "", fmt.Errorf("获取章节内容失败: %w", err)
	}

	// 4. 缓存章节内容（30分钟）
	if s.cacheService != nil {
		_ = s.cacheService.SetChapterContent(ctx, chapterID, content, 30*time.Minute)
	}

	// 5. 发布阅读事件
	s.publishReadingEvent(ctx, userID, chapterID)

	return content, nil
}

// GetFirstChapter 获取第一章
func (s *ReaderService) GetFirstChapter(ctx context.Context, bookID string) (*reader2.Chapter, error) {
	chapter, err := s.chapterRepo.GetFirstChapter(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取第一章失败: %w", err)
	}
	return chapter, nil
}

// GetLastChapter 获取最后一章
func (s *ReaderService) GetLastChapter(ctx context.Context, bookID string) (*reader2.Chapter, error) {
	chapter, err := s.chapterRepo.GetLastChapter(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取最后一章失败: %w", err)
	}
	return chapter, nil
}

// =========================
// 阅读进度相关方法
// =========================

// GetReadingProgress 获取阅读进度
func (s *ReaderService) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader2.ReadingProgress, error) {
	progress, err := s.progressRepo.GetByUserAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取阅读进度失败: %w", err)
	}

	// 如果没有阅读记录，返回空进度
	if progress == nil {
		progress = &reader2.ReadingProgress{
			UserID:      userID,
			BookID:      bookID,
			Progress:    0,
			ReadingTime: 0,
		}
	}

	return progress, nil
}

// SaveReadingProgress 保存阅读进度
func (s *ReaderService) SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	if progress < 0 || progress > 1 {
		return fmt.Errorf("进度值必须在0-1之间")
	}

	err := s.progressRepo.SaveProgress(ctx, userID, bookID, chapterID, progress)
	if err != nil {
		return fmt.Errorf("保存阅读进度失败: %w", err)
	}

	// 发布进度更新事件
	s.publishProgressEvent(ctx, userID, bookID, chapterID, progress)

	return nil
}

// UpdateReadingTime 更新阅读时长
func (s *ReaderService) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	if duration <= 0 {
		return fmt.Errorf("阅读时长必须大于0")
	}

	err := s.progressRepo.UpdateReadingTime(ctx, userID, bookID, duration)
	if err != nil {
		return fmt.Errorf("更新阅读时长失败: %w", err)
	}

	return nil
}

// GetRecentReading 获取最近阅读记录
func (s *ReaderService) GetRecentReading(ctx context.Context, userID string, limit int) ([]*reader2.ReadingProgress, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	progresses, err := s.progressRepo.GetRecentReadingByUser(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取最近阅读记录失败: %w", err)
	}

	return progresses, nil
}

// GetReadingHistory 获取阅读历史
func (s *ReaderService) GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*reader2.ReadingProgress, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	offset := (page - 1) * size
	progresses, err := s.progressRepo.GetReadingHistory(ctx, userID, size, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("获取阅读历史失败: %w", err)
	}

	total, err := s.progressRepo.CountReadingBooks(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("统计阅读书籍数失败: %w", err)
	}

	return progresses, total, nil
}

// GetTotalReadingTime 获取总阅读时长
func (s *ReaderService) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	total, err := s.progressRepo.GetTotalReadingTime(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("获取总阅读时长失败: %w", err)
	}
	return total, nil
}

// GetReadingTimeByPeriod 获取时间段内的阅读时长
func (s *ReaderService) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	total, err := s.progressRepo.GetReadingTimeByPeriod(ctx, userID, startTime, endTime)
	if err != nil {
		return 0, fmt.Errorf("获取时间段阅读时长失败: %w", err)
	}
	return total, nil
}

// GetUnfinishedBooks 获取未读完的书籍
func (s *ReaderService) GetUnfinishedBooks(ctx context.Context, userID string) ([]*reader2.ReadingProgress, error) {
	progresses, err := s.progressRepo.GetUnfinishedBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取未读完书籍失败: %w", err)
	}
	return progresses, nil
}

// GetFinishedBooks 获取已读完的书籍
func (s *ReaderService) GetFinishedBooks(ctx context.Context, userID string) ([]*reader2.ReadingProgress, error) {
	progresses, err := s.progressRepo.GetFinishedBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取已读完书籍失败: %w", err)
	}
	return progresses, nil
}

// =========================
// 标注相关方法
// =========================

// CreateAnnotation 创建标注
func (s *ReaderService) CreateAnnotation(ctx context.Context, annotation *reader2.Annotation) error {
	// 参数验证
	if err := s.validateAnnotation(annotation); err != nil {
		return fmt.Errorf("标注参数验证失败: %w", err)
	}

	err := s.annotationRepo.Create(ctx, annotation)
	if err != nil {
		return fmt.Errorf("创建标注失败: %w", err)
	}

	// 发布标注创建事件
	s.publishAnnotationEvent(ctx, "created", annotation)

	return nil
}

// UpdateAnnotation 更新标注
func (s *ReaderService) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	err := s.annotationRepo.Update(ctx, annotationID, updates)
	if err != nil {
		return fmt.Errorf("更新标注失败: %w", err)
	}

	return nil
}

// DeleteAnnotation 删除标注
func (s *ReaderService) DeleteAnnotation(ctx context.Context, annotationID string) error {
	err := s.annotationRepo.Delete(ctx, annotationID)
	if err != nil {
		return fmt.Errorf("删除标注失败: %w", err)
	}

	return nil
}

// GetAnnotationsByChapter 获取章节的标注
func (s *ReaderService) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader2.Annotation, error) {
	annotations, err := s.annotationRepo.GetByUserAndChapter(ctx, userID, bookID, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节标注失败: %w", err)
	}
	return annotations, nil
}

// GetAnnotationsByBook 获取书籍的所有标注
func (s *ReaderService) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	annotations, err := s.annotationRepo.GetByUserAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取书籍标注失败: %w", err)
	}
	return annotations, nil
}

// GetNotes 获取笔记
func (s *ReaderService) GetNotes(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	notes, err := s.annotationRepo.GetNotes(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取笔记失败: %w", err)
	}
	return notes, nil
}

// SearchNotes 搜索笔记
func (s *ReaderService) SearchNotes(ctx context.Context, userID, keyword string) ([]*reader2.Annotation, error) {
	if keyword == "" {
		return nil, fmt.Errorf("搜索关键词不能为空")
	}

	notes, err := s.annotationRepo.SearchNotes(ctx, userID, keyword)
	if err != nil {
		return nil, fmt.Errorf("搜索笔记失败: %w", err)
	}
	return notes, nil
}

// GetBookmarks 获取书签
func (s *ReaderService) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	bookmarks, err := s.annotationRepo.GetBookmarks(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取书签失败: %w", err)
	}
	return bookmarks, nil
}

// GetLatestBookmark 获取最新的书签
func (s *ReaderService) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader2.Annotation, error) {
	bookmark, err := s.annotationRepo.GetLatestBookmark(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取最新书签失败: %w", err)
	}
	return bookmark, nil
}

// GetHighlights 获取高亮
func (s *ReaderService) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader2.Annotation, error) {
	highlights, err := s.annotationRepo.GetHighlights(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取高亮失败: %w", err)
	}
	return highlights, nil
}

// GetRecentAnnotations 获取最近的标注
func (s *ReaderService) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader2.Annotation, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	annotations, err := s.annotationRepo.GetRecentAnnotations(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取最近标注失败: %w", err)
	}
	return annotations, nil
}

// GetPublicAnnotations 获取公开的标注
func (s *ReaderService) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader2.Annotation, error) {
	annotations, err := s.annotationRepo.GetPublicAnnotations(ctx, bookID, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取公开标注失败: %w", err)
	}
	return annotations, nil
}

// =========================
// 阅读设置相关方法
// =========================

// GetReadingSettings 获取阅读设置
func (s *ReaderService) GetReadingSettings(ctx context.Context, userID string) (*reader2.ReadingSettings, error) {
	// 1. 尝试从缓存获取
	if s.cacheService != nil {
		cachedSettings, err := s.cacheService.GetReadingSettings(ctx, userID)
		if err == nil && cachedSettings != nil {
			return cachedSettings, nil
		}
	}

	// 2. 从数据库获取
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取阅读设置失败: %w", err)
	}

	// 3. 如果没有设置，返回默认设置
	if settings == nil {
		settings = s.getDefaultSettings(userID)
	} else {
		// 4. 缓存设置（1小时）
		if s.cacheService != nil {
			_ = s.cacheService.SetReadingSettings(ctx, userID, settings, time.Hour)
		}
	}

	return settings, nil
}

// SaveReadingSettings 保存阅读设置
func (s *ReaderService) SaveReadingSettings(ctx context.Context, settings *reader2.ReadingSettings) error {
	if settings.UserID == "" {
		return fmt.Errorf("用户ID不能为空")
	}

	// 检查是否已存在
	exists, err := s.settingsRepo.ExistsByUserID(ctx, settings.UserID)
	if err != nil {
		return fmt.Errorf("检查设置是否存在失败: %w", err)
	}

	if exists {
		// 更新现有设置
		err = s.settingsRepo.UpdateByUserID(ctx, settings.UserID, settings)
		if err != nil {
			return fmt.Errorf("更新阅读设置失败: %w", err)
		}
	} else {
		// 创建新设置
		err = s.settingsRepo.Create(ctx, settings)
		if err != nil {
			return fmt.Errorf("创建阅读设置失败: %w", err)
		}
	}

	// 更新缓存
	if s.cacheService != nil {
		_ = s.cacheService.SetReadingSettings(ctx, settings.UserID, settings, time.Hour)
	}

	return nil
}

// UpdateReadingSettings 更新阅读设置
func (s *ReaderService) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	// 获取现有设置
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取阅读设置失败: %w", err)
	}

	if settings == nil {
		return fmt.Errorf("阅读设置不存在")
	}

	// 应用更新
	if fontSize, ok := updates["font_size"]; ok {
		settings.FontSize = fontSize.(int)
	}
	if fontFamily, ok := updates["font_family"]; ok {
		settings.FontFamily = fontFamily.(string)
	}
	if lineHeight, ok := updates["line_height"]; ok {
		settings.LineHeight = lineHeight.(float64)
	}
	if theme, ok := updates["theme"]; ok {
		settings.Theme = theme.(string)
	}
	if background, ok := updates["background"]; ok {
		settings.Background = background.(string)
	}
	if pageMode, ok := updates["page_mode"]; ok {
		settings.PageMode = pageMode.(int)
	}
	if autoScroll, ok := updates["auto_scroll"]; ok {
		settings.AutoScroll = autoScroll.(bool)
	}
	if scrollSpeed, ok := updates["scroll_speed"]; ok {
		settings.ScrollSpeed = scrollSpeed.(int)
	}

	settings.UpdatedAt = time.Now()

	err = s.settingsRepo.UpdateByUserID(ctx, userID, settings)
	if err != nil {
		return fmt.Errorf("更新阅读设置失败: %w", err)
	}

	// 更新缓存
	if s.cacheService != nil {
		_ = s.cacheService.SetReadingSettings(ctx, userID, settings, time.Hour)
	}

	return nil
}

// =========================
// 私有辅助方法
// =========================

// validateAnnotation 验证标注参数
func (s *ReaderService) validateAnnotation(annotation *reader2.Annotation) error {
	if annotation.UserID == "" {
		return fmt.Errorf("用户ID不能为空")
	}
	if annotation.BookID == "" {
		return fmt.Errorf("书籍ID不能为空")
	}
	if annotation.ChapterID == "" {
		return fmt.Errorf("章节ID不能为空")
	}
	if annotation.Type == "" {
		return fmt.Errorf("标注类型不能为空")
	}
	// 验证标注类型是否为有效值
	if annotation.Type != "bookmark" && annotation.Type != "highlight" && annotation.Type != "note" {
		return fmt.Errorf("标注类型必须是bookmark(书签)、highlight(高亮)或note(笔记)")
	}
	return nil
}

// getDefaultSettings 获取默认阅读设置
func (s *ReaderService) getDefaultSettings(userID string) *reader2.ReadingSettings {
	return &reader2.ReadingSettings{
		UserID:      userID,
		FontSize:    16,
		FontFamily:  "serif",
		LineHeight:  1.5,
		Theme:       "light",
		Background:  "#FFFFFF",
		PageMode:    1, // 1-滑动
		AutoScroll:  false,
		ScrollSpeed: 50,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// publishReadingEvent 发布阅读事件
func (s *ReaderService) publishReadingEvent(ctx context.Context, userID, chapterID string) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: "reader.chapter.read",
		EventData: map[string]interface{}{
			"user_id":    userID,
			"chapter_id": chapterID,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}

// publishProgressEvent 发布进度更新事件
func (s *ReaderService) publishProgressEvent(ctx context.Context, userID, bookID, chapterID string, progress float64) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: "reader.progress.updated",
		EventData: map[string]interface{}{
			"user_id":    userID,
			"book_id":    bookID,
			"chapter_id": chapterID,
			"progress":   progress,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}

// publishAnnotationEvent 发布标注事件
func (s *ReaderService) publishAnnotationEvent(ctx context.Context, action string, annotation *reader2.Annotation) {
	if s.eventBus == nil {
		return
	}

	event := &base.BaseEvent{
		EventType: fmt.Sprintf("reader.annotation.%s", action),
		EventData: map[string]interface{}{
			"user_id":    annotation.UserID,
			"book_id":    annotation.BookID,
			"chapter_id": annotation.ChapterID,
			"type":       annotation.Type,
		},
		Timestamp: time.Now(),
		Source:    s.serviceName,
	}

	s.eventBus.PublishAsync(ctx, event)
}

// =========================
// 批量操作方法
// =========================

// BatchCreateAnnotations 批量创建注记
func (s *ReaderService) BatchCreateAnnotations(ctx context.Context, annotations []*reader2.Annotation) error {
	if len(annotations) == 0 {
		return nil
	}

	if len(annotations) > 50 {
		return fmt.Errorf("批量创建注记数量不能超过50个")
	}

	// 批量创建
	for _, annotation := range annotations {
		if err := s.CreateAnnotation(ctx, annotation); err != nil {
			return fmt.Errorf("批量创建注记失败: %w", err)
		}
	}

	return nil
}

// BatchDeleteAnnotations 批量删除注记
func (s *ReaderService) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	if len(annotationIDs) == 0 {
		return nil
	}

	if len(annotationIDs) > 100 {
		return fmt.Errorf("批量删除注记数量不能超过100个")
	}

	// 批量删除
	for _, id := range annotationIDs {
		if err := s.DeleteAnnotation(ctx, id); err != nil {
			return fmt.Errorf("批量删除注记失败: %w", err)
		}
	}

	return nil
}

// GetAnnotationStats 获取注记统计
func (s *ReaderService) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	annotations, err := s.GetAnnotationsByBook(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取注记失败: %w", err)
	}

	stats := map[string]interface{}{
		"totalCount":     len(annotations),
		"bookmarkCount":  0,
		"highlightCount": 0,
		"noteCount":      0,
	}

	// 统计各类型注记数量
	for _, ann := range annotations {
		switch ann.Type {
		case "bookmark":
			stats["bookmarkCount"] = stats["bookmarkCount"].(int) + 1
		case "highlight":
			stats["highlightCount"] = stats["highlightCount"].(int) + 1
		case "note":
			stats["noteCount"] = stats["noteCount"].(int) + 1
		}
	}

	return stats, nil
}

// SyncAnnotationsRequest 同步注记请求（内部使用）
type SyncAnnotationsRequest struct {
	BookID           string
	LastSyncTime     int64
	LocalAnnotations []*reader2.Annotation
}

// SyncAnnotations 同步注记（多端同步）
func (s *ReaderService) SyncAnnotations(ctx context.Context, userID string, req interface{}) (map[string]interface{}, error) {
	// 类型断言
	syncReq, ok := req.(*SyncAnnotationsRequest)
	if !ok {
		return nil, fmt.Errorf("无效的同步请求类型")
	}

	// 1. 获取服务器端的注记
	serverAnnotations, err := s.GetAnnotationsByBook(ctx, userID, syncReq.BookID)
	if err != nil {
		return nil, fmt.Errorf("获取服务器注记失败: %w", err)
	}

	// 2. 过滤出需要下发的注记（比lastSyncTime更新的）
	newAnnotations := make([]*reader2.Annotation, 0)
	for _, ann := range serverAnnotations {
		if ann.CreatedAt.Unix() > syncReq.LastSyncTime {
			newAnnotations = append(newAnnotations, ann)
		}
	}

	// 3. 上传本地新增的注记
	uploadedCount := 0
	if len(syncReq.LocalAnnotations) > 0 {
		for _, ann := range syncReq.LocalAnnotations {
			ann.UserID = userID // 确保UserID正确
			if err := s.CreateAnnotation(ctx, ann); err != nil {
				// 记录错误但继续
				continue
			}
			uploadedCount++
		}
	}

	return map[string]interface{}{
		"newAnnotations":  newAnnotations,
		"syncTime":        time.Now().Unix(),
		"uploadedCount":   uploadedCount,
		"downloadedCount": len(newAnnotations),
	}, nil
}
