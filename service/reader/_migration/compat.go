package reader

import (
	"context"
	"time"

	readerModel "Qingyu_backend/models/reader"
	readeriface "Qingyu_backend/service/interfaces/reader"
)

// ============================================================================
// 兼容层 - 向后兼容支持
// ============================================================================

// ReaderServiceAdapter 旧 ReaderService 接口的适配器
// 将旧的 ReaderService 方法调用委托给新的 Port 接口
type ReaderServiceAdapter struct {
	progressPort   readeriface.ReadingProgressPort
	annotationPort readeriface.AnnotationPort
	chapterPort    readeriface.ChapterContentPort
	settingsPort   readeriface.ReaderSettingsPort
	syncPort       readeriface.ReaderSyncPort
}

// ReaderChapterServiceAdapter 旧 ReaderChapterService 接口的适配器
// 将旧的 ReaderChapterService 方法调用委托给 ChapterContentPort
type ReaderChapterServiceAdapter struct {
	chapterPort readeriface.ChapterContentPort
}

// ============================================================================
// ReaderServiceAdapter - ReaderService 接口实现
// ============================================================================

// NewReaderServiceAdapter 创建新的 ReaderService 适配器
func NewReaderServiceAdapter(
	progressPort readeriface.ReadingProgressPort,
	annotationPort readeriface.AnnotationPort,
	chapterPort readeriface.ChapterContentPort,
	settingsPort readeriface.ReaderSettingsPort,
	syncPort readeriface.ReaderSyncPort,
) *ReaderServiceAdapter {
	return &ReaderServiceAdapter{
		progressPort:   progressPort,
		annotationPort: annotationPort,
		chapterPort:    chapterPort,
		settingsPort:   settingsPort,
		syncPort:       syncPort,
	}
}

// ============================================================================
// BaseService 接口实现 - 委托给 progressPort
// ============================================================================

// Initialize 初始化服务
func (a *ReaderServiceAdapter) Initialize(ctx context.Context) error {
	return a.progressPort.Initialize(ctx)
}

// Health 健康检查
func (a *ReaderServiceAdapter) Health(ctx context.Context) error {
	return a.progressPort.Health(ctx)
}

// Close 关闭服务
func (a *ReaderServiceAdapter) Close(ctx context.Context) error {
	// 关闭所有 Port
	if err := a.progressPort.Close(ctx); err != nil {
		return err
	}
	if err := a.annotationPort.Close(ctx); err != nil {
		return err
	}
	if err := a.chapterPort.Close(ctx); err != nil {
		return err
	}
	if err := a.settingsPort.Close(ctx); err != nil {
		return err
	}
	if err := a.syncPort.Close(ctx); err != nil {
		return err
	}
	return nil
}

// GetServiceName 获取服务名称
func (a *ReaderServiceAdapter) GetServiceName() string {
	return a.progressPort.GetServiceName()
}

// GetVersion 获取服务版本
func (a *ReaderServiceAdapter) GetVersion() string {
	return a.progressPort.GetVersion()
}

// ============================================================================
// 章节相关方法 - 委托给 ChapterContentPort
// ============================================================================

// GetChapterContent 委托给 ChapterContentPort
func (a *ReaderServiceAdapter) GetChapterContent(ctx context.Context, userID, chapterID string) (string, error) {
	return a.chapterPort.GetChapterContent(ctx, userID, chapterID)
}

// GetChapterByID 委托给 ChapterContentPort
func (a *ReaderServiceAdapter) GetChapterByID(ctx context.Context, chapterID string) (interface{}, error) {
	return a.chapterPort.GetChapterByID(ctx, chapterID)
}

// GetBookChapters 委托给 ChapterContentPort
func (a *ReaderServiceAdapter) GetBookChapters(ctx context.Context, bookID string, page, size int) (interface{}, int64, error) {
	return a.chapterPort.GetBookChapters(ctx, bookID, page, size)
}

// ============================================================================
// 阅读进度相关方法 - 委托给 ReadingProgressPort
// ============================================================================

// GetReadingProgress 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) GetReadingProgress(ctx context.Context, userID, bookID string) (*readerModel.ReadingProgress, error) {
	return a.progressPort.GetReadingProgress(ctx, userID, bookID)
}

// SaveReadingProgress 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) SaveReadingProgress(ctx context.Context, userID, bookID, chapterID string, progress float64) error {
	req := &readeriface.SaveReadingProgressRequest{
		UserID:    userID,
		BookID:    bookID,
		ChapterID: chapterID,
		Progress:  progress,
	}
	return a.progressPort.SaveReadingProgress(ctx, req)
}

// UpdateReadingTime 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) UpdateReadingTime(ctx context.Context, userID, bookID string, duration int64) error {
	req := &readeriface.UpdateReadingTimeRequest{
		UserID:   userID,
		BookID:   bookID,
		Duration: duration,
	}
	return a.progressPort.UpdateReadingTime(ctx, req)
}

// GetRecentReading 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) GetRecentReading(ctx context.Context, userID string, limit int) ([]*readerModel.ReadingProgress, error) {
	return a.progressPort.GetRecentReading(ctx, userID, limit)
}

// GetReadingHistory 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) GetReadingHistory(ctx context.Context, userID string, page, size int) ([]*readerModel.ReadingProgress, int64, error) {
	req := &readeriface.GetReadingHistoryRequest{
		UserID: userID,
		Page:   page,
		Size:   size,
	}
	resp, err := a.progressPort.GetReadingHistory(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Progresses, resp.Total, nil
}

// GetTotalReadingTime 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) GetTotalReadingTime(ctx context.Context, userID string) (int64, error) {
	return a.progressPort.GetTotalReadingTime(ctx, userID)
}

// GetReadingTimeByPeriod 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) GetReadingTimeByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (int64, error) {
	req := &readeriface.GetReadingTimeByPeriodRequest{
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
	}
	return a.progressPort.GetReadingTimeByPeriod(ctx, req)
}

// GetUnfinishedBooks 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) GetUnfinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	return a.progressPort.GetUnfinishedBooks(ctx, userID)
}

// GetFinishedBooks 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) GetFinishedBooks(ctx context.Context, userID string) ([]*readerModel.ReadingProgress, error) {
	return a.progressPort.GetFinishedBooks(ctx, userID)
}

// DeleteReadingProgress 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) DeleteReadingProgress(ctx context.Context, userID, bookID string) error {
	return a.progressPort.DeleteReadingProgress(ctx, userID, bookID)
}

// UpdateBookStatus 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) UpdateBookStatus(ctx context.Context, userID, bookID, status string) error {
	req := &readeriface.UpdateBookStatusRequest{
		UserID: userID,
		BookID: bookID,
		Status: status,
	}
	return a.progressPort.UpdateBookStatus(ctx, req)
}

// BatchUpdateBookStatus 委托给 ReadingProgressPort
func (a *ReaderServiceAdapter) BatchUpdateBookStatus(ctx context.Context, userID string, bookIDs []string, status string) error {
	req := &readeriface.BatchUpdateBookStatusRequest{
		UserID:  userID,
		BookIDs: bookIDs,
		Status:  status,
	}
	return a.progressPort.BatchUpdateBookStatus(ctx, req)
}

// ============================================================================
// 标注相关方法 - 委托给 AnnotationPort
// ============================================================================

// CreateAnnotation 委托给 AnnotationPort
func (a *ReaderServiceAdapter) CreateAnnotation(ctx context.Context, annotation *readerModel.Annotation) error {
	return a.annotationPort.CreateAnnotation(ctx, annotation)
}

// UpdateAnnotation 委托给 AnnotationPort
func (a *ReaderServiceAdapter) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	return a.annotationPort.UpdateAnnotation(ctx, annotationID, updates)
}

// DeleteAnnotation 委托给 AnnotationPort
func (a *ReaderServiceAdapter) DeleteAnnotation(ctx context.Context, annotationID string) error {
	return a.annotationPort.DeleteAnnotation(ctx, annotationID)
}

// GetAnnotationsByChapter 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*readerModel.Annotation, error) {
	return a.annotationPort.GetAnnotationsByChapter(ctx, userID, bookID, chapterID)
}

// GetAnnotationsByBook 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return a.annotationPort.GetAnnotationsByBook(ctx, userID, bookID)
}

// GetNotes 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetNotes(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return a.annotationPort.GetNotes(ctx, userID, bookID)
}

// SearchNotes 委托给 AnnotationPort
func (a *ReaderServiceAdapter) SearchNotes(ctx context.Context, userID, keyword string) ([]*readerModel.Annotation, error) {
	return a.annotationPort.SearchNotes(ctx, userID, keyword)
}

// GetBookmarks 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetBookmarks(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return a.annotationPort.GetBookmarks(ctx, userID, bookID)
}

// GetLatestBookmark 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetLatestBookmark(ctx context.Context, userID, bookID string) (*readerModel.Annotation, error) {
	return a.annotationPort.GetLatestBookmark(ctx, userID, bookID)
}

// GetHighlights 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetHighlights(ctx context.Context, userID, bookID string) ([]*readerModel.Annotation, error) {
	return a.annotationPort.GetHighlights(ctx, userID, bookID)
}

// GetRecentAnnotations 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*readerModel.Annotation, error) {
	return a.annotationPort.GetRecentAnnotations(ctx, userID, limit)
}

// GetPublicAnnotations 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*readerModel.Annotation, error) {
	return a.annotationPort.GetPublicAnnotations(ctx, bookID, chapterID)
}

// GetAnnotationStats 委托给 AnnotationPort
func (a *ReaderServiceAdapter) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	return a.annotationPort.GetAnnotationStats(ctx, userID, bookID)
}

// BatchCreateAnnotations 委托给 AnnotationPort
func (a *ReaderServiceAdapter) BatchCreateAnnotations(ctx context.Context, annotations []*readerModel.Annotation) error {
	return a.annotationPort.BatchCreateAnnotations(ctx, annotations)
}

// BatchDeleteAnnotations 委托给 AnnotationPort
func (a *ReaderServiceAdapter) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	return a.annotationPort.BatchDeleteAnnotations(ctx, annotationIDs)
}

// ============================================================================
// 阅读设置相关方法 - 委托给 ReaderSettingsPort
// ============================================================================

// GetReadingSettings 委托给 ReaderSettingsPort
func (a *ReaderServiceAdapter) GetReadingSettings(ctx context.Context, userID string) (*readerModel.ReadingSettings, error) {
	return a.settingsPort.GetReadingSettings(ctx, userID)
}

// SaveReadingSettings 委托给 ReaderSettingsPort
func (a *ReaderServiceAdapter) SaveReadingSettings(ctx context.Context, settings *readerModel.ReadingSettings) error {
	return a.settingsPort.SaveReadingSettings(ctx, settings)
}

// UpdateReadingSettings 委托给 ReaderSettingsPort
func (a *ReaderServiceAdapter) UpdateReadingSettings(ctx context.Context, userID string, updates map[string]interface{}) error {
	return a.settingsPort.UpdateReadingSettings(ctx, userID, updates)
}

// ============================================================================
// 同步方法 - 委托给 ReaderSyncPort
// ============================================================================

// SyncAnnotations 委托给 ReaderSyncPort
func (a *ReaderServiceAdapter) SyncAnnotations(ctx context.Context, userID string, req interface{}) (map[string]interface{}, error) {
	// 从旧的 req 结构体中提取数据
	type legacySyncReq struct {
		BookID           string
		LastSyncTime     int64
		LocalAnnotations []*readerModel.Annotation
	}
	legacyReq := req.(legacySyncReq)

	syncReq := &readeriface.SyncAnnotationsRequest{
		UserID:           userID,
		BookID:           legacyReq.BookID,
		LastSyncTime:     legacyReq.LastSyncTime,
		LocalAnnotations: legacyReq.LocalAnnotations,
	}
	resp, err := a.syncPort.SyncAnnotations(ctx, syncReq)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"newAnnotations":  resp.NewAnnotations,
		"syncTime":        resp.SyncTime,
		"uploadedCount":   resp.UploadedCount,
		"downloadedCount": resp.DownloadedCount,
	}, nil
}

// ============================================================================
// ReaderChapterServiceAdapter - ReaderChapterService 接口实现
// ============================================================================

// NewReaderChapterServiceAdapter 创建新的 ReaderChapterService 适配器
func NewReaderChapterServiceAdapter(chapterPort readeriface.ChapterContentPort) *ReaderChapterServiceAdapter {
	return &ReaderChapterServiceAdapter{
		chapterPort: chapterPort,
	}
}

// GetChapterContent 委托给 ChapterContentPort
func (a *ReaderChapterServiceAdapter) GetChapterContent(ctx context.Context, userID, bookID, chapterID string) (interface{}, error) {
	req := &readeriface.GetChapterContentRequest{
		UserID:    userID,
		BookID:    bookID,
		ChapterID: chapterID,
	}
	return a.chapterPort.GetChapterContentWithProgress(ctx, req)
}

// GetChapterByNumber 委托给 ChapterContentPort
func (a *ReaderChapterServiceAdapter) GetChapterByNumber(ctx context.Context, userID, bookID string, chapterNum int) (interface{}, error) {
	return a.chapterPort.GetChapterByNumber(ctx, userID, bookID, chapterNum)
}

// GetNextChapter 委托给 ChapterContentPort
func (a *ReaderChapterServiceAdapter) GetNextChapter(ctx context.Context, userID, bookID, chapterID string) (interface{}, error) {
	return a.chapterPort.GetNextChapter(ctx, userID, bookID, chapterID)
}

// GetPreviousChapter 委托给 ChapterContentPort
func (a *ReaderChapterServiceAdapter) GetPreviousChapter(ctx context.Context, userID, bookID, chapterID string) (interface{}, error) {
	return a.chapterPort.GetPreviousChapter(ctx, userID, bookID, chapterID)
}

// GetChapterList 委托给 ChapterContentPort
func (a *ReaderChapterServiceAdapter) GetChapterList(ctx context.Context, userID, bookID string, page, size int) (interface{}, error) {
	return a.chapterPort.GetChapterList(ctx, userID, bookID, page, size)
}

// GetChapterInfo 委托给 ChapterContentPort
func (a *ReaderChapterServiceAdapter) GetChapterInfo(ctx context.Context, userID, chapterID string) (interface{}, error) {
	return a.chapterPort.GetChapterInfo(ctx, userID, chapterID)
}

// ============================================================================
// 编译时检查
// ============================================================================

// 确保 ReaderServiceAdapter 实现了旧 ReaderService 接口
// var _ interfaces.ReaderService = (*ReaderServiceAdapter)(nil)

// 确保 ReaderChapterServiceAdapter 实现了旧 ReaderChapterService 接口
// var _ interfaces.ReaderChapterService = (*ReaderChapterServiceAdapter)(nil)
