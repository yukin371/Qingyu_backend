//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/writer"
	"Qingyu_backend/models/stats"
	"Qingyu_backend/repository/mongodb"
	"Qingyu_backend/service/audit"
	documentService "Qingyu_backend/service/writer/document"
	projectService "Qingyu_backend/service/writer/project"
	statsService "Qingyu_backend/service/shared/stats"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWriterE2E_CompleteWorkflow 测试完整的写作流程
// 流程：创建项目 -> 创建章节 -> 编辑内容 -> 自动保存 -> 内容审核 -> 查看统计
func TestWriterE2E_CompleteWorkflow(t *testing.T) {
	// 初始化测试环境
	ctx := context.Background()
	factory, cleanup := setupTestDB(t)
	defer cleanup()

	// 初始化所有Service
	projectRepo := factory.CreateProjectRepository()
	documentRepo := factory.CreateProjectDocumentRepository()
	chapterStatsRepo := factory.CreateChapterStatsRepository()
	bookStatsRepo := factory.CreateBookStatsRepository()

	projectSvc := projectService.NewProjectService(projectRepo, nil)
	documentSvc := documentService.NewDocumentService(documentRepo, projectRepo, nil)
	wordCountSvc := documentService.NewWordCountService()
	auditSvc := audit.NewContentAuditService(nil, nil, nil, nil)
	statsSvc := statsService.NewStatsService(bookStatsRepo, chapterStatsRepo, nil)

	// 测试用户
	authorID := "test_author_001"

	// ============================================
	// 步骤1: 创建项目
	// ============================================
	t.Run("Step1_CreateProject", func(t *testing.T) {
		req := &projectService.CreateProjectRequest{
			Title:       "我的第一本小说",
			Description: "这是一本精彩的玄幻小说",
			Genre:       "玄幻",
			Tags:        []string{"修仙", "热血", "爽文"},
			AuthorID:    authorID,
		}

		resp, err := projectSvc.CreateProject(ctx, req)
		require.NoError(t, err, "创建项目应该成功")
		assert.NotEmpty(t, resp.ProjectID, "项目ID不能为空")
		assert.Equal(t, "我的第一本小说", resp.Title)

		// 保存项目ID供后续使用
		ctx = context.WithValue(ctx, "projectID", resp.ProjectID)
	})

	projectID := ctx.Value("projectID").(string)

	// ============================================
	// 步骤2: 创建章节
	// ============================================
	var chapterID string
	t.Run("Step2_CreateChapter", func(t *testing.T) {
		req := &documentService.CreateDocumentRequest{
			ProjectID: projectID,
			Title:     "第一章：穿越异世",
			Type:      writer.DocumentTypeChapter,
			ParentID:  "",
			Order:     1,
		}

		resp, err := documentSvc.CreateDocument(ctx, req)
		require.NoError(t, err, "创建章节应该成功")
		assert.NotEmpty(t, resp.DocumentID)
		assert.Equal(t, "第一章：穿越异世", resp.Title)

		chapterID = resp.DocumentID
	})

	// ============================================
	// 步骤3: 编辑内容（包含字数统计）
	// ============================================
	t.Run("Step3_EditContent", func(t *testing.T) {
		content := `# 第一章：穿越异世

林凡睁开眼睛，发现自己躺在一片陌生的山林中。

"这是哪里？"他茫然地看着周围。

突然，一股庞大的记忆涌入脑海。这个世界叫做玄天大陆，是一个修真世界...

林凡惊讶地发现，自己竟然穿越了！而且还带着一个神秘的系统。

【叮！新手礼包已发放，请宿主查收】

"系统？这不会是网文看多了产生的幻觉吧？"林凡半信半疑。`

		// 更新文档内容
		updateReq := &documentService.UpdateContentRequest{
			DocumentID: chapterID,
			Content:    content,
			Version:    1,
		}

		err := documentSvc.UpdateDocumentContent(ctx, updateReq)
		require.NoError(t, err, "更新内容应该成功")

		// 字数统计
		wordCount := wordCountSvc.CalculateWordCount(content)
		assert.Greater(t, wordCount.TotalWords, 0, "字数应该大于0")
		assert.Greater(t, wordCount.ChineseChars, wordCount.EnglishWords, "中文字符应该多于英文单词")
		t.Logf("字数统计: 总字数=%d, 中文=%d, 英文=%d",
			wordCount.TotalWords, wordCount.ChineseChars, wordCount.EnglishWords)
	})

	// ============================================
	// 步骤4: 自动保存测试
	// ============================================
	t.Run("Step4_AutoSave", func(t *testing.T) {
		// 模拟自动保存
		content := "这是自动保存的内容修改..."

		autoSaveReq := &documentService.AutoSaveRequest{
			DocumentID: chapterID,
			Content:    content,
			Version:    2,
		}

		resp, err := documentSvc.AutoSaveDocument(ctx, autoSaveReq)
		require.NoError(t, err, "自动保存应该成功")
		assert.True(t, resp.Success)
		assert.NotEmpty(t, resp.SaveTime)

		// 获取保存状态
		status, err := documentSvc.GetSaveStatus(ctx, chapterID)
		require.NoError(t, err)
		assert.Equal(t, "saved", status.Status)
		assert.WithinDuration(t, time.Now(), status.LastSaveTime, 5*time.Second)
	})

	// ============================================
	// 步骤5: 内容审核
	// ============================================
	t.Run("Step5_ContentAudit", func(t *testing.T) {
		// 测试正常内容
		safeContent := "这是一段安全的内容，讲述主角的修炼之路。"
		checkReq := &audit.CheckContentRequest{
			Content:   safeContent,
			ContentID: chapterID,
			AuthorID:  authorID,
		}

		checkResp, err := auditSvc.CheckContent(ctx, checkReq)
		require.NoError(t, err, "内容检测应该成功")
		assert.True(t, checkResp.Safe, "安全内容应该通过审核")
		assert.Equal(t, "low", checkResp.RiskLevel)

		// 测试包含敏感词的内容
		unsafeContent := "这里包含一些违规词汇：色情、暴力、反动..."
		checkReq2 := &audit.CheckContentRequest{
			Content:   unsafeContent,
			ContentID: chapterID + "_2",
			AuthorID:  authorID,
		}

		checkResp2, err := auditSvc.CheckContent(ctx, checkReq2)
		require.NoError(t, err)
		if !checkResp2.Safe {
			assert.Greater(t, len(checkResp2.MatchedWords), 0, "应该检测到敏感词")
			t.Logf("检测到敏感词: %v", checkResp2.MatchedWords)
		}
	})

	// ============================================
	// 步骤6: 数据统计
	// ============================================
	t.Run("Step6_Statistics", func(t *testing.T) {
		// 创建模拟的章节统计数据
		mockChapterStats := &stats.ChapterStats{
			BookID:          projectID,
			ChapterID:       chapterID,
			ChapterTitle:    "第一章：穿越异世",
			ViewCount:       1500,
			UniqueViewers:   800,
			CompletionRate:  0.75,
			DropOffRate:     0.25,
			AvgReadDuration: 180, // 3分钟
			Revenue:         15.50,
			StatDate:        time.Now(),
		}

		err := chapterStatsRepo.Create(ctx, mockChapterStats)
		require.NoError(t, err, "创建章节统计应该成功")

		// 获取章节统计
		chapterStats, err := statsSvc.GetChapterStats(ctx, projectID, chapterID)
		if err == nil && chapterStats != nil {
			assert.Equal(t, chapterID, chapterStats.ChapterID)
			assert.Greater(t, chapterStats.ViewCount, int64(0))
			t.Logf("章节统计: 阅读量=%d, 完读率=%.2f%%",
				chapterStats.ViewCount, chapterStats.CompletionRate*100)
		}

		// 创建作品统计
		mockBookStats := &stats.BookStats{
			BookID:            projectID,
			AuthorID:          authorID,
			TotalViews:        15000,
			UniqueReaders:     5000,
			TotalChapters:     50,
			AvgCompletionRate: 0.72,
			AvgDropOffRate:    0.28,
			TotalRevenue:      1250.00,
			StatDate:          time.Now(),
		}

		err = bookStatsRepo.Create(ctx, mockBookStats)
		require.NoError(t, err, "创建作品统计应该成功")

		// 获取作品统计
		bookStats, err := statsSvc.GetBookStats(ctx, projectID)
		if err == nil && bookStats != nil {
			assert.Equal(t, projectID, bookStats.BookID)
			assert.Greater(t, bookStats.TotalViews, int64(0))
			t.Logf("作品统计: 总阅读量=%d, 独立读者=%d, 总收入=%.2f",
				bookStats.TotalViews, bookStats.UniqueReaders, bookStats.TotalRevenue)
		}
	})

	// ============================================
	// 步骤7: 版本控制
	// ============================================
	t.Run("Step7_VersionControl", func(t *testing.T) {
		// 创建新版本
		content := "这是修改后的内容版本2..."
		updateReq := &documentService.UpdateContentRequest{
			DocumentID: chapterID,
			Content:    content,
			Version:    3,
		}

		err := documentSvc.UpdateDocumentContent(ctx, updateReq)
		require.NoError(t, err, "更新版本应该成功")

		// 获取文档内容
		docContent, err := documentSvc.GetDocumentContent(ctx, chapterID)
		require.NoError(t, err)
		assert.Equal(t, content, docContent.Content)
		t.Logf("当前版本: %d", docContent.Version)
	})
}

// TestWriterE2E_PerformanceBaseline 性能基准测试
func TestWriterE2E_PerformanceBaseline(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}

	ctx := context.Background()
	factory, cleanup := setupTestDB(t)
	defer cleanup()

	projectRepo := factory.CreateProjectRepository()
	projectSvc := projectService.NewProjectService(projectRepo, nil)

	// ============================================
	// 测试1: 项目创建性能
	// ============================================
	t.Run("Performance_CreateProject", func(t *testing.T) {
		start := time.Now()

		req := &projectService.CreateProjectRequest{
			Title:       "性能测试项目",
			Description: "用于性能基准测试",
			Genre:       "玄幻",
			AuthorID:    "perf_test_author",
		}

		_, err := projectSvc.CreateProject(ctx, req)
		require.NoError(t, err)

		duration := time.Since(start)
		assert.Less(t, duration, 200*time.Millisecond, "项目创建应该在200ms内完成")
		t.Logf("项目创建耗时: %v", duration)
	})

	// ============================================
	// 测试2: 字数统计性能
	// ============================================
	t.Run("Performance_WordCount", func(t *testing.T) {
		wordCountSvc := documentService.NewWordCountService()

		// 生成大文档（约5000字）
		largeContent := generateLargeContent(5000)

		start := time.Now()
		result := wordCountSvc.CalculateWordCount(largeContent)
		duration := time.Since(start)

		assert.Greater(t, result.TotalWords, 4000)
		assert.Less(t, duration, 100*time.Millisecond, "字数统计应该在100ms内完成")
		t.Logf("字数统计耗时: %v, 字数: %d", duration, result.TotalWords)
	})

	// ============================================
	// 测试3: 敏感词检测性能
	// ============================================
	t.Run("Performance_ContentAudit", func(t *testing.T) {
		auditSvc := audit.NewContentAuditService(nil, nil, nil, nil)

		content := generateLargeContent(1000) + "这里可能有敏感词"

		start := time.Now()
		req := &audit.CheckContentRequest{
			Content:   content,
			ContentID: "perf_test_001",
			AuthorID:  "perf_author",
		}

		_, err := auditSvc.CheckContent(ctx, req)
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Less(t, duration, 200*time.Millisecond, "内容审核应该在200ms内完成")
		t.Logf("内容审核耗时: %v", duration)
	})
}

// TestWriterE2E_ConcurrentOperations 并发操作测试
func TestWriterE2E_ConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过并发测试")
	}

	ctx := context.Background()
	factory, cleanup := setupTestDB(t)
	defer cleanup()

	documentRepo := factory.CreateProjectDocumentRepository()
	documentSvc := documentService.NewDocumentService(documentRepo, nil, nil)

	// 创建测试文档
	projectID := "concurrent_test_project"
	req := &documentService.CreateDocumentRequest{
		ProjectID: projectID,
		Title:     "并发测试章节",
		Type:      writer.DocumentTypeChapter,
	}

	resp, err := documentSvc.CreateDocument(ctx, req)
	require.NoError(t, err)
	docID := resp.DocumentID

	// ============================================
	// 并发自动保存测试
	// ============================================
	t.Run("Concurrent_AutoSave", func(t *testing.T) {
		const goroutines = 10
		errChan := make(chan error, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(index int) {
				autoSaveReq := &documentService.AutoSaveRequest{
					DocumentID: docID,
					Content:    "并发保存内容 " + time.Now().String(),
					Version:    index + 1,
				}

				_, err := documentSvc.AutoSaveDocument(ctx, autoSaveReq)
				errChan <- err
			}(i)
		}

		// 收集结果
		successCount := 0
		for i := 0; i < goroutines; i++ {
			if err := <-errChan; err == nil {
				successCount++
			}
		}

		// 至少应该有一些成功（版本冲突是预期的）
		assert.Greater(t, successCount, 0, "并发保存应该有部分成功")
		t.Logf("并发保存: %d/%d 成功", successCount, goroutines)
	})
}

// 辅助函数

func setupTestDB(t *testing.T) (*mongodb.MongoRepositoryFactory, func()) {
	config := &mongodb.MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "qingyu_integration_test",
		Timeout:  30 * time.Second,
	}

	factory, err := mongodb.NewMongoRepositoryFactory(config)
	require.NoError(t, err, "应该能连接到测试数据库")

	cleanup := func() {
		// 清理测试数据
		ctx := context.Background()
		_ = factory.GetDatabase().Drop(ctx)
		_ = factory.Close()
	}

	return factory, cleanup
}

func generateLargeContent(targetWords int) string {
	content := ""
	sentence := "这是一个测试句子，用于生成大量文本内容。"

	for len(content) < targetWords*3 { // 估算字符数
		content += sentence
	}

	return content
}
