package content

import (
	"context"
	"fmt"

	"Qingyu_backend/models/dto"
	readerService "Qingyu_backend/service/reader"
)

// ChapterAdapter 章节适配器
// 将现有ReaderService的章节功能适配到ChapterService接口
type ChapterAdapter struct {
	readerService *readerService.ReaderService
}

// NewChapterAdapter 创建章节适配器
func NewChapterAdapter(readerService *readerService.ReaderService) *ChapterAdapter {
	return &ChapterAdapter{
		readerService: readerService,
	}
}

// =========================
// 章节内容获取
// =========================

// GetChapter 获取章节内容
func (a *ChapterAdapter) GetChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error) {
	// 调用ReaderService获取章节内容
	content, err := a.readerService.GetChapterContent(ctx, "", chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节内容失败: %w", err)
	}

	// 获取章节信息
	_, err = a.readerService.GetChapterByID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节信息失败: %w", err)
	}

	// 转换为DTO
	// 注意：章节信息返回的是interface{}，需要类型断言
	// 这里简化处理，实际需要根据具体的chapter类型来转换
	response := &dto.ChapterResponse{
		ChapterID:   chapterID,
		BookID:      bookID,
		ChapterNum:  0, // 需要从章节信息获取
		Title:       "", // 需要从章节信息获取
		Content:     content,
		WordCount:   len(content), // 简单计算字数
		IsVIP:       false,        // 需要从章节信息获取
		PublishedAt: 0,            // 需要从章节信息获取
	}

	return response, nil
}

// GetChapterByNumber 根据章节号获取章节
func (a *ChapterAdapter) GetChapterByNumber(ctx context.Context, bookID string, chapterNum int) (*dto.ChapterResponse, error) {
	// 现有ReaderService没有直接按章节号获取的方法
	// 需要通过获取章节列表后查找
	// 这里返回错误，提示需要扩展功能
	return nil, fmt.Errorf("暂不支持按章节号获取章节，需要扩展ReaderService或直接使用Bookstore的ChapterService")
}

// GetChapterInfo 获取章节信息（不含内容）
func (a *ChapterAdapter) GetChapterInfo(ctx context.Context, chapterID string) (*dto.ChapterInfo, error) {
	// 获取章节信息
	_, err := a.readerService.GetChapterByID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节信息失败: %w", err)
	}

	// 转换为DTO
	// 注意：这里需要根据实际的章节信息类型进行转换
	info := &dto.ChapterInfo{
		ChapterID:   chapterID,
		BookID:      "", // 需要从章节信息获取
		ChapterNum:  0,  // 需要从章节信息获取
		Title:       "", // 需要从章节信息获取
		WordCount:   0,  // 需要从章节信息获取
		IsVIP:       false,
		PublishedAt: 0,
	}

	return info, nil
}

// =========================
// 章节导航
// =========================

// GetNextChapter 获取下一章
func (a *ChapterAdapter) GetNextChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error) {
	// 现有ReaderService没有直接获取下一章的方法
	// 需要获取章节列表后查找
	return nil, fmt.Errorf("暂不支持获取下一章，需要扩展功能")
}

// GetPreviousChapter 获取上一章
func (a *ChapterAdapter) GetPreviousChapter(ctx context.Context, bookID, chapterID string) (*dto.ChapterResponse, error) {
	// 现有ReaderService没有直接获取上一章的方法
	// 需要获取章节列表后查找
	return nil, fmt.Errorf("暂不支持获取上一章，需要扩展功能")
}

// ListChapters 获取章节列表
func (a *ChapterAdapter) ListChapters(ctx context.Context, bookID string) (*dto.ChapterListResponse, error) {
	// 调用ReaderService获取章节列表
	_, total, err := a.readerService.GetBookChapters(ctx, bookID, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// 转换为DTO
	// 注意：chapters是interface{}类型，需要根据实际类型转换
	chapterInfos := make([]*dto.ChapterInfo, 0)

	// TODO: 这里需要根据实际的chapters类型进行转换
	// 例如，如果chapters是[]*bookstore.Chapter：
	// for _, chapter := range chapters {
	//     chapterInfos = append(chapterInfos, &dto.ChapterInfo{
	//         ChapterID:   chapter.ID.Hex(),
	//         BookID:      chapter.BookID.Hex(),
	//         ChapterNum:  chapter.ChapterNum,
	//         Title:       chapter.Title,
	//         WordCount:   chapter.WordCount,
	//         IsVIP:       chapter.IsVIP,
	//         PublishedAt: chapter.PublishedAt.Unix(),
	//     })
	// }

	return &dto.ChapterListResponse{
		Chapters: chapterInfos,
		Total:    total,
		Page:     1,
		PageSize: 1000,
	}, nil
}

// =========================
// 章节搜索与筛选
// =========================

// SearchChapters 搜索章节
func (a *ChapterAdapter) SearchChapters(ctx context.Context, bookID string, keyword string, page, pageSize int) (*dto.ChapterListResponse, error) {
	// 现有ReaderService没有搜索功能
	// 需要扩展或使用Bookstore的ChapterService
	return nil, fmt.Errorf("暂不支持章节搜索，需要扩展功能")
}

// GetChaptersByType 根据类型获取章节
func (a *ChapterAdapter) GetChaptersByType(ctx context.Context, bookID, chapterType string) (*dto.ChapterListResponse, error) {
	// 现有ReaderService没有按类型筛选的功能
	return nil, fmt.Errorf("暂不支持按类型获取章节，需要扩展功能")
}

// =========================
// 章节状态管理
// =========================

// GetChapterPublishStatus 获取章节发布状态
func (a *ChapterAdapter) GetChapterPublishStatus(ctx context.Context, chapterID string) (*dto.ChapterPublishStatus, error) {
	// 现有ReaderService没有获取发布状态的方法
	// 需要扩展或使用Bookstore的ChapterService
	return nil, fmt.Errorf("暂不支持获取章节发布状态，需要扩展功能")
}

// UpdateChapterPublishStatus 更新章节发布状态
func (a *ChapterAdapter) UpdateChapterPublishStatus(ctx context.Context, chapterID string, req *dto.UpdateChapterPublishStatusRequest) error {
	// 现有ReaderService没有更新发布状态的方法
	// 需要扩展或使用Bookstore的ChapterService
	return fmt.Errorf("暂不支持更新章节发布状态，需要扩展功能")
}

// =========================
// 批量操作
// =========================

// BatchGetChapters 批量获取章节
func (a *ChapterAdapter) BatchGetChapters(ctx context.Context, chapterIDs []string) ([]*dto.ChapterResponse, error) {
	responses := make([]*dto.ChapterResponse, 0, len(chapterIDs))

	// 逐个获取章节（后续可以优化为批量查询）
	for _, chapterID := range chapterIDs {
		// 这里需要bookID，暂时使用空字符串
		chapter, err := a.GetChapter(ctx, "", chapterID)
		if err != nil {
			// 记录错误但继续处理
			continue
		}
		responses = append(responses, chapter)
	}

	return responses, nil
}

// GetChapterRange 获取章节范围
func (a *ChapterAdapter) GetChapterRange(ctx context.Context, bookID string, startNum, endNum int) (*dto.ChapterListResponse, error) {
	// 现有ReaderService没有按范围获取章节的方法
	// 需要扩展或使用Bookstore的ChapterService
	return nil, fmt.Errorf("暂不支持按范围获取章节，需要扩展功能")
}

// =========================
// 扩展功能
// =========================

// GetChapterContentWithProgress 获取章节内容并更新阅读进度
// 这是一个常用的组合操作，避免多次API调用
func (a *ChapterAdapter) GetChapterContentWithProgress(ctx context.Context, userID, bookID, chapterID string) (*dto.ChapterResponse, error) {
	// 1. 获取章节内容
	chapter, err := a.GetChapter(ctx, bookID, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节内容失败: %w", err)
	}

	// TODO: 2. 更新阅读进度（需要注入ProgressAdapter）
	// progressAdapter := NewProgressAdapter(a.readerService)
	// progressAdapter.SaveProgress(ctx, &dto.SaveProgressRequest{
	//     UserID:    userID,
	//     BookID:    bookID,
	//     ChapterID: chapterID,
	//     Progress:  calculateProgress(chapterID),
	// })

	return chapter, nil
}

// GetChapterNavigation 获取章节导航信息
// 返回上一章、当前章、下一章的信息
func (a *ChapterAdapter) GetChapterNavigation(ctx context.Context, bookID, chapterID string) (*ChapterNavigationInfo, error) {
	// 获取章节列表
	allChapters, err := a.ListChapters(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// 查找当前章节的索引
	currentIndex := -1
	for i, chapter := range allChapters.Chapters {
		if chapter.ChapterID == chapterID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return nil, fmt.Errorf("未找到指定章节")
	}

	navInfo := &ChapterNavigationInfo{
		CurrentChapter: allChapters.Chapters[currentIndex],
	}

	// 获取上一章
	if currentIndex > 0 {
		navInfo.PreviousChapter = allChapters.Chapters[currentIndex-1]
	}

	// 获取下一章
	if currentIndex < len(allChapters.Chapters)-1 {
		navInfo.NextChapter = allChapters.Chapters[currentIndex+1]
	}

	return navInfo, nil
}

// ChapterNavigationInfo 章节导航信息
type ChapterNavigationInfo struct {
	PreviousChapter *dto.ChapterInfo `json:"previousChapter,omitempty"`
	CurrentChapter  *dto.ChapterInfo `json:"currentChapter"`
	NextChapter     *dto.ChapterInfo `json:"nextChapter,omitempty"`
}

// GetChapterByTitle 根据标题获取章节
func (a *ChapterAdapter) GetChapterByTitle(ctx context.Context, bookID, title string) (*dto.ChapterResponse, error) {
	// 获取所有章节
	allChapters, err := a.ListChapters(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// 查找匹配标题的章节
	for _, chapter := range allChapters.Chapters {
		if chapter.Title == title {
			return a.GetChapter(ctx, bookID, chapter.ChapterID)
		}
	}

	return nil, fmt.Errorf("未找到标题为'%s'的章节", title)
}

// CalculateReadingTime 估算章节阅读时间
// 根据章节字数和平均阅读速度估算
func (a *ChapterAdapter) CalculateReadingTime(ctx context.Context, chapterID string) (int64, error) {
	// 获取章节信息
	chapter, err := a.GetChapterInfo(ctx, chapterID)
	if err != nil {
		return 0, fmt.Errorf("获取章节信息失败: %w", err)
	}

	// 假设平均阅读速度为500字/分钟
	wordsPerMinute := 500
	readingMinutes := chapter.WordCount / wordsPerMinute

	return int64(readingMinutes * 60), nil // 返回秒数
}

// ValidateChapterAccess 验证章节访问权限
// 检查用户是否有权限访问该章节（如VIP章节）
func (a *ChapterAdapter) ValidateChapterAccess(ctx context.Context, userID, chapterID string) (bool, error) {
	// 获取章节信息
	chapter, err := a.GetChapterInfo(ctx, chapterID)
	if err != nil {
		return false, fmt.Errorf("获取章节信息失败: %w", err)
	}

	// 如果不是VIP章节，直接允许访问
	if !chapter.IsVIP {
		return true, nil
	}

	// TODO: 检查用户VIP状态
	// vipService := ...
	// return vipService.HasAccessToVIPContent(ctx, userID)

	// 暂时返回true
	return true, nil
}

// GetChapterCount 获取书籍章节数量
func (a *ChapterAdapter) GetChapterCount(ctx context.Context, bookID string) (int, error) {
	// 获取章节列表
	allChapters, err := a.ListChapters(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("获取章节列表失败: %w", err)
	}

	return len(allChapters.Chapters), nil
}

// GetFirstChapter 获取第一章
func (a *ChapterAdapter) GetFirstChapter(ctx context.Context, bookID string) (*dto.ChapterResponse, error) {
	// 获取章节列表
	allChapters, err := a.ListChapters(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// 返回第一章
	if len(allChapters.Chapters) > 0 {
		return a.GetChapter(ctx, bookID, allChapters.Chapters[0].ChapterID)
	}

	return nil, fmt.Errorf("书籍没有章节")
}

// GetLastChapter 获取最后一章
func (a *ChapterAdapter) GetLastChapter(ctx context.Context, bookID string) (*dto.ChapterResponse, error) {
	// 获取章节列表
	allChapters, err := a.ListChapters(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// 返回最后一章
	if len(allChapters.Chapters) > 0 {
		lastIndex := len(allChapters.Chapters) - 1
		return a.GetChapter(ctx, bookID, allChapters.Chapters[lastIndex].ChapterID)
	}

	return nil, fmt.Errorf("书籍没有章节")
}

// GetChaptersBatch 分批获取章节内容
// 用于离线下载等场景
func (a *ChapterAdapter) GetChaptersBatch(ctx context.Context, bookID string, startNum, batchSize int) ([]*dto.ChapterResponse, error) {
	// 获取章节列表
	allChapters, err := a.ListChapters(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// 计算批次范围
	endNum := startNum + batchSize
	if endNum > len(allChapters.Chapters) {
		endNum = len(allChapters.Chapters)
	}

	// 获取批次内的章节
	chapters := make([]*dto.ChapterResponse, 0)
	for i := startNum; i < endNum; i++ {
		chapter, err := a.GetChapter(ctx, bookID, allChapters.Chapters[i].ChapterID)
		if err != nil {
			// 记录错误但继续处理
			continue
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

// GetChapterWordCount 获取章节字数统计
func (a *ChapterAdapter) GetChapterWordCount(ctx context.Context, chapterID string) (int, error) {
	chapter, err := a.GetChapterInfo(ctx, chapterID)
	if err != nil {
		return 0, fmt.Errorf("获取章节信息失败: %w", err)
	}

	return chapter.WordCount, nil
}

// GetBookWordCount 获取书籍总字数
func (a *ChapterAdapter) GetBookWordCount(ctx context.Context, bookID string) (int, error) {
	// 获取章节列表
	allChapters, err := a.ListChapters(ctx, bookID)
	if err != nil {
		return 0, fmt.Errorf("获取章节列表失败: %w", err)
	}

	// 累加字数
	totalWords := 0
	for _, chapter := range allChapters.Chapters {
		totalWords += chapter.WordCount
	}

	return totalWords, nil
}

// UpdateChapterReadingProgress 更新章节阅读进度
// 记录用户阅读到该章节的哪个位置
func (a *ChapterAdapter) UpdateChapterReadingProgress(ctx context.Context, userID, chapterID string, position int) error {
	// TODO: 实现章节内阅读进度跟踪
	// 这需要扩展Annotation模型来支持章节内的位置标记
	return fmt.Errorf("暂不支持章节内进度跟踪，需要扩展功能")
}

// GetChapterReadingProgress 获取章节阅读进度
func (a *ChapterAdapter) GetChapterReadingProgress(ctx context.Context, userID, chapterID string) (int, error) {
	// TODO: 实现章节内阅读进度查询
	return 0, fmt.Errorf("暂不支持章节内进度查询，需要扩展功能")
}

// MarkChapterAsRead 标记章节为已读
func (a *ChapterAdapter) MarkChapterAsRead(ctx context.Context, userID, chapterID string) error {
	// TODO: 实现章节已读标记
	return fmt.Errorf("暂不支持章节已读标记，需要扩展功能")
}

// GetUnreadChapters 获取未读章节列表
func (a *ChapterAdapter) GetUnreadChapters(ctx context.Context, userID, bookID string) ([]*dto.ChapterInfo, error) {
	// TODO: 实现未读章节查询
	return nil, fmt.Errorf("暂不支持未读章节查询，需要扩展功能")
}

// IsChapterRead 检查章节是否已读
func (a *ChapterAdapter) IsChapterRead(ctx context.Context, userID, chapterID string) (bool, error) {
	// TODO: 实现章节已读状态检查
	return false, fmt.Errorf("暂不支持章节已读状态检查，需要扩展功能")
}

// GetChapterComments 获取章节评论
func (a *ChapterAdapter) GetChapterComments(ctx context.Context, chapterID string, page, pageSize int) ([]interface{}, error) {
	// TODO: 集成评论服务
	return nil, fmt.Errorf("暂不支持获取章节评论，需要集成评论服务")
}

// GetChapterAnnotations 获取章节标注
func (a *ChapterAdapter) GetChapterAnnotations(ctx context.Context, userID, bookID, chapterID string) ([]interface{}, error) {
	// 调用ReaderService的标注功能
	annotations, err := a.readerService.GetAnnotationsByChapter(ctx, userID, bookID, chapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节标注失败: %w", err)
	}

	// 转换为通用格式（这里简化处理）
	result := make([]interface{}, 0, len(annotations))
	for _, ann := range annotations {
		result = append(result, ann)
	}

	return result, nil
}

// GetChapterBookmarks 获取章节书签
func (a *ChapterAdapter) GetChapterBookmarks(ctx context.Context, userID, bookID, chapterID string) ([]interface{}, error) {
	// 调用ReaderService的书签功能
	bookmarks, err := a.readerService.GetBookmarks(ctx, userID, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取章节书签失败: %w", err)
	}

	// 筛选属于该章节的书签
	result := make([]interface{}, 0)
	for _, bookmark := range bookmarks {
		// bookmark是*reader.Annotation类型
		if bookmark.ChapterID.Hex() == chapterID {
			result = append(result, bookmark)
		}
	}

	return result, nil
}
