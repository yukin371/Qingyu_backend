package impl

import (
	"context"

	"Qingyu_backend/models/reader"
	serviceReader "Qingyu_backend/service/interfaces/reader"
	readerService "Qingyu_backend/service/reader"
)

// ReaderSyncImpl 阅读数据同步端口实现
type ReaderSyncImpl struct {
	readerService *readerService.ReaderService
	serviceName   string
	version       string
}

// NewReaderSyncImpl 创建阅读数据同步端口实现
func NewReaderSyncImpl(readerService *readerService.ReaderService) serviceReader.ReaderSyncPort {
	return &ReaderSyncImpl{
		readerService: readerService,
		serviceName:   "ReaderSyncPort",
		version:       "1.0.0",
	}
}

// ============================================================================
// BaseService 生命周期方法实现
// ============================================================================

func (r *ReaderSyncImpl) Initialize(ctx context.Context) error {
	return r.readerService.Initialize(ctx)
}

func (r *ReaderSyncImpl) Health(ctx context.Context) error {
	return r.readerService.Health(ctx)
}

func (r *ReaderSyncImpl) Close(ctx context.Context) error {
	return r.readerService.Close(ctx)
}

func (r *ReaderSyncImpl) GetServiceName() string {
	return r.serviceName
}

func (r *ReaderSyncImpl) GetVersion() string {
	return r.version
}

// ============================================================================
// ReaderSyncPort 方法实现
// ============================================================================

// SyncAnnotations 同步标注（多端同步）
func (r *ReaderSyncImpl) SyncAnnotations(ctx context.Context, req *serviceReader.SyncAnnotationsRequest) (*serviceReader.SyncAnnotationsResponse, error) {
	// 构建内部同步请求
	internalReq := &readerService.SyncAnnotationsRequest{
		BookID:           req.BookID,
		LastSyncTime:     req.LastSyncTime,
		LocalAnnotations: req.LocalAnnotations,
	}

	// 委托给现有 Service
	result, err := r.readerService.SyncAnnotations(ctx, req.UserID, internalReq)
	if err != nil {
		return nil, err
	}

	// 转换响应类型
	newAnnotations := make([]*reader.Annotation, 0)
	if anns, ok := result["newAnnotations"].([]*reader.Annotation); ok {
		newAnnotations = anns
	}

	syncTime := int64(0)
	if st, ok := result["syncTime"].(int64); ok {
		syncTime = st
	} else if stFloat, ok := result["syncTime"].(float64); ok {
		syncTime = int64(stFloat)
	}

	uploadedCount := 0
	if uc, ok := result["uploadedCount"].(int); ok {
		uploadedCount = uc
	} else if ucFloat, ok := result["uploadedCount"].(float64); ok {
		uploadedCount = int(ucFloat)
	}

	downloadedCount := 0
	if dc, ok := result["downloadedCount"].(int); ok {
		downloadedCount = dc
	} else if dcFloat, ok := result["downloadedCount"].(float64); ok {
		downloadedCount = int(dcFloat)
	}

	return &serviceReader.SyncAnnotationsResponse{
		NewAnnotations:  newAnnotations,
		SyncTime:        syncTime,
		UploadedCount:   uploadedCount,
		DownloadedCount: downloadedCount,
	}, nil
}
