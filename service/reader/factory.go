package reader

import (
	"Qingyu_backend/service/interfaces/reader"
	readermigration "Qingyu_backend/service/reader/_migration"
)

// ReaderServiceFactory Reader服务工厂
// 提供创建和组装Reader服务的方法
type ReaderServiceFactory struct{}

// NewReaderServiceFactory 创建工厂实例
func NewReaderServiceFactory() *ReaderServiceFactory {
	return &ReaderServiceFactory{}
}

// CreateWithPorts 使用 Port 接口创建服务（推荐方式）
//
// 新架构推荐的使用方式：
// 1. 实现 5 个 Port 接口的具体实现
// 2. 使用 ReaderServiceAdapter 组装它们
// 3. 返回兼容的 ReaderService 接口供 API 层使用
func (f *ReaderServiceFactory) CreateWithPorts(
	progressPort reader.ReadingProgressPort,
	annotationPort reader.AnnotationPort,
	chapterPort reader.ChapterContentPort,
	settingsPort reader.ReaderSettingsPort,
	syncPort reader.ReaderSyncPort,
) interface{} {
	return readermigration.NewReaderServiceAdapter(
		progressPort,
		annotationPort,
		chapterPort,
		settingsPort,
		syncPort,
	)
}

// CreateChapterServiceWithPorts 使用 Port 接口创建章节服务（推荐方式）
func (f *ReaderServiceFactory) CreateChapterServiceWithPorts(
	chapterPort reader.ChapterContentPort,
) interface{} {
	return readermigration.NewReaderChapterServiceAdapter(chapterPort)
}

// PortImplementations Port 接口实现集合
type PortImplementations struct {
	ProgressPort   reader.ReadingProgressPort
	AnnotationPort reader.AnnotationPort
	ChapterPort    reader.ChapterContentPort
	SettingsPort   reader.ReaderSettingsPort
	SyncPort       reader.ReaderSyncPort
}

// CreateFromImplementations 从结构体创建服务
func (f *ReaderServiceFactory) CreateFromImplementations(ports PortImplementations) interface{} {
	return readermigration.NewReaderServiceAdapter(
		ports.ProgressPort,
		ports.AnnotationPort,
		ports.ChapterPort,
		ports.SettingsPort,
		ports.SyncPort,
	)
}

// ChapterPortImplementations 章节 Port 接口实现集合
type ChapterPortImplementations struct {
	ChapterPort reader.ChapterContentPort
}

// CreateChapterServiceFromImplementations 从结构体创建章节服务
func (f *ReaderServiceFactory) CreateChapterServiceFromImplementations(ports ChapterPortImplementations) interface{} {
	return readermigration.NewReaderChapterServiceAdapter(ports.ChapterPort)
}
