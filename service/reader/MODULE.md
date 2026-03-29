# Reader Service

> 最后更新：2026-03-29

## 职责

阅读体验层，管理章节内容获取、阅读进度、阅读历史、书签/高亮/批注、阅读设置、阅读时长统计。不管理书籍元数据（由 Bookstore 负责）。

## 数据流

```
API Handler → ReaderService → Repository → MongoDB
                ↓
         EventBus（reading.progress、reading.annotation、reading.event）
```

## 约定 & 陷阱

- **阅读进度自动保存**：前端通过 `SaveReadingProgress` 定期提交，需注意防抖避免频繁写入
- **批注同步**：`SyncAnnotations` 支持客户端离线后批量同步，使用时间戳解决冲突（last-write-wins）
- **阅读时长统计**：`GetReadingTimeByPeriod` 按天/周/月聚合，依赖 `UpdateReadingTime` 的准确上报
- **书籍状态管理**：`UpdateBookStatus`/`BatchUpdateBookStatus` 管理"在读/已读/想读"状态，变更会触发事件
- **批注可见性**：`GetPublicAnnotations` 返回公开批注，私有批注只能通过 `GetAnnotationsByChapter` 获取
