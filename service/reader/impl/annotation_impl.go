package impl

import (
	"context"

	"Qingyu_backend/models/reader"
	serviceReader "Qingyu_backend/service/interfaces/reader"
	readerService "Qingyu_backend/service/reader"
)

// AnnotationImpl 标注与笔记管理端口实现
type AnnotationImpl struct {
	readerService *readerService.ReaderService
	cacheService  readerService.AnnotationCacheService
	serviceName   string
	version       string
}

// NewAnnotationImpl 创建标注与笔记管理端口实现
func NewAnnotationImpl(
	readerService *readerService.ReaderService,
	cacheService readerService.AnnotationCacheService,
) serviceReader.AnnotationPort {
	return &AnnotationImpl{
		readerService: readerService,
		cacheService:  cacheService,
		serviceName:   "AnnotationPort",
		version:       "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (a *AnnotationImpl) Initialize(ctx context.Context) error {
	return a.readerService.Initialize(ctx)
}

func (a *AnnotationImpl) Health(ctx context.Context) error {
	return a.readerService.Health(ctx)
}

func (a *AnnotationImpl) Close(ctx context.Context) error {
	return a.readerService.Close(ctx)
}

func (a *AnnotationImpl) GetServiceName() string {
	return a.serviceName
}

func (a *AnnotationImpl) GetVersion() string {
	return a.version
}

// ============================================================================
// AnnotationPort 方法实现
// ============================================================================

// CreateAnnotation 创建标注
func (a *AnnotationImpl) CreateAnnotation(ctx context.Context, annotation *reader.Annotation) error {
	return a.readerService.CreateAnnotation(ctx, annotation)
}

// UpdateAnnotation 更新标注
func (a *AnnotationImpl) UpdateAnnotation(ctx context.Context, annotationID string, updates map[string]interface{}) error {
	return a.readerService.UpdateAnnotation(ctx, annotationID, updates)
}

// DeleteAnnotation 删除标注
func (a *AnnotationImpl) DeleteAnnotation(ctx context.Context, annotationID string) error {
	return a.readerService.DeleteAnnotation(ctx, annotationID)
}

// GetAnnotationsByChapter 获取章节的标注
func (a *AnnotationImpl) GetAnnotationsByChapter(ctx context.Context, userID, bookID, chapterID string) ([]*reader.Annotation, error) {
	return a.readerService.GetAnnotationsByChapter(ctx, userID, bookID, chapterID)
}

// GetAnnotationsByBook 获取书籍的所有标注
func (a *AnnotationImpl) GetAnnotationsByBook(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return a.readerService.GetAnnotationsByBook(ctx, userID, bookID)
}

// GetNotes 获取笔记
func (a *AnnotationImpl) GetNotes(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return a.readerService.GetNotes(ctx, userID, bookID)
}

// SearchNotes 搜索笔记
func (a *AnnotationImpl) SearchNotes(ctx context.Context, userID, keyword string) ([]*reader.Annotation, error) {
	return a.readerService.SearchNotes(ctx, userID, keyword)
}

// GetBookmarks 获取书签
func (a *AnnotationImpl) GetBookmarks(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return a.readerService.GetBookmarks(ctx, userID, bookID)
}

// GetLatestBookmark 获取最新的书签
func (a *AnnotationImpl) GetLatestBookmark(ctx context.Context, userID, bookID string) (*reader.Annotation, error) {
	return a.readerService.GetLatestBookmark(ctx, userID, bookID)
}

// GetHighlights 获取高亮
func (a *AnnotationImpl) GetHighlights(ctx context.Context, userID, bookID string) ([]*reader.Annotation, error) {
	return a.readerService.GetHighlights(ctx, userID, bookID)
}

// GetRecentAnnotations 获取最近的标注
func (a *AnnotationImpl) GetRecentAnnotations(ctx context.Context, userID string, limit int) ([]*reader.Annotation, error) {
	return a.readerService.GetRecentAnnotations(ctx, userID, limit)
}

// GetPublicAnnotations 获取公开的标注
func (a *AnnotationImpl) GetPublicAnnotations(ctx context.Context, bookID, chapterID string) ([]*reader.Annotation, error) {
	return a.readerService.GetPublicAnnotations(ctx, bookID, chapterID)
}

// GetAnnotationStats 获取标注统计
func (a *AnnotationImpl) GetAnnotationStats(ctx context.Context, userID, bookID string) (map[string]interface{}, error) {
	return a.readerService.GetAnnotationStats(ctx, userID, bookID)
}

// BatchCreateAnnotations 批量创建标注
func (a *AnnotationImpl) BatchCreateAnnotations(ctx context.Context, annotations []*reader.Annotation) error {
	return a.readerService.BatchCreateAnnotations(ctx, annotations)
}

// BatchDeleteAnnotations 批量删除标注
func (a *AnnotationImpl) BatchDeleteAnnotations(ctx context.Context, annotationIDs []string) error {
	return a.readerService.BatchDeleteAnnotations(ctx, annotationIDs)
}
