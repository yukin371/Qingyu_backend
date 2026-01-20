//go:build e2e
// +build e2e

package data

import (
	"testing"
)

// TestScenarioBuilder_BuildReaderWithProgress 测试构建阅读进度场景
func TestScenarioBuilder_BuildReaderWithProgress(t *testing.T) {
	SetupTestEnvironment(t)

	builder := NewScenarioBuilder(t)
	scenario := builder.BuildReaderWithProgress()

	// 验证场景数据
	if scenario.User == nil {
		t.Error("expected user to be created")
	}
	if scenario.CurrentBook == nil {
		t.Error("expected current book to be created")
	}
	if scenario.CurrentChapter == nil {
		t.Error("expected current chapter to be created")
	}
	if len(scenario.Books) == 0 {
		t.Error("expected books to be created")
	}
	if scenario.Progress != 0.3 {
		t.Errorf("expected progress 0.3, got %f", scenario.Progress)
	}

	t.Logf("Created reader scenario: user=%s, book=%s, chapter=%d",
		scenario.User.Username, scenario.CurrentBook.Title, scenario.CurrentChapter.ChapterNum)
}

// TestScenarioBuilder_BuildSocialInteraction 测试构建社交互动场景
func TestScenarioBuilder_BuildSocialInteraction(t *testing.T) {
	SetupTestEnvironment(t)

	builder := NewScenarioBuilder(t)
	scenario := builder.BuildSocialInteraction(5)

	// 验证场景数据
	if len(scenario.Users) != 5 {
		t.Errorf("expected 5 users, got %d", len(scenario.Users))
	}
	if scenario.Book == nil {
		t.Error("expected book to be created")
	}
	if len(scenario.Comments) != 5 {
		t.Errorf("expected 5 comments, got %d", len(scenario.Comments))
	}
	if len(scenario.Collections) != 2 { // 5 / 2 = 2
		t.Errorf("expected 2 collections, got %d", len(scenario.Collections))
	}

	// 验证评论
	for i, comment := range scenario.Comments {
		if comment.ID == "" {
			t.Errorf("comment %d: expected ID to be set", i)
		}
		if comment.AuthorID != scenario.Users[i].ID {
			t.Errorf("comment %d: expected author ID to match", i)
		}
	}

	// 验证收藏
	for i, collection := range scenario.Collections {
		if collection.ID.IsZero() {
			t.Errorf("collection %d: expected ID to be set", i)
		}
	}

	t.Logf("Created social interaction scenario: %d users, %d comments, %d collections",
		len(scenario.Users), len(scenario.Comments), len(scenario.Collections))
}

// TestScenarioBuilder_BuildPaidContent 测试构建付费内容场景
func TestScenarioBuilder_BuildPaidContent(t *testing.T) {
	SetupTestEnvironment(t)

	builder := NewScenarioBuilder(t)
	scenario := builder.BuildPaidContent()

	// 验证场景数据
	if scenario.Author == nil {
		t.Error("expected author to be created")
	}
	if scenario.FreeUser == nil {
		t.Error("expected free user to be created")
	}
	if scenario.VIPUser == nil {
		t.Error("expected VIP user to be created")
	}
	if scenario.PaidBook == nil {
		t.Error("expected paid book to be created")
	}
	if len(scenario.FreeChapters) != 3 {
		t.Errorf("expected 3 free chapters, got %d", len(scenario.FreeChapters))
	}
	if len(scenario.PaidChapters) != 7 {
		t.Errorf("expected 7 paid chapters, got %d", len(scenario.PaidChapters))
	}

	// 验证免费章节
	for i, chapter := range scenario.FreeChapters {
		if !chapter.IsFree {
			t.Errorf("free chapter %d: expected to be free", i)
		}
	}

	// 验证付费章节
	for i, chapter := range scenario.PaidChapters {
		if chapter.IsFree {
			t.Errorf("paid chapter %d: expected to be paid", i)
		}
	}

	// 验证用户VIP等级
	if scenario.FreeUser.VIPLevel != 0 {
		t.Errorf("expected free user VIP level 0, got %d", scenario.FreeUser.VIPLevel)
	}
	if scenario.VIPUser.VIPLevel != 1 {
		t.Errorf("expected VIP user VIP level 1, got %d", scenario.VIPUser.VIPLevel)
	}

	t.Logf("Created paid content scenario: author=%s, %d free chapters, %d paid chapters",
		scenario.Author.Username, len(scenario.FreeChapters), len(scenario.PaidChapters))
}

// TestScenarioBuilder_Integration 测试场景构建器的集成使用
func TestScenarioBuilder_Integration(t *testing.T) {
	SetupTestEnvironment(t)

	builder := NewScenarioBuilder(t)

	// 构建阅读进度场景
	readerScenario := builder.BuildReaderWithProgress()

	// 构建社交互动场景
	socialScenario := builder.BuildSocialInteraction(3)

	// 构建付费内容场景
	paidScenario := builder.BuildPaidContent()

	// 验证所有场景都能正常工作
	if readerScenario.User.ID == "" {
		t.Error("reader scenario: expected user ID")
	}
	if len(socialScenario.Comments) != 3 {
		t.Error("social scenario: expected 3 comments")
	}
	if len(paidScenario.FreeChapters) != 3 {
		t.Error("paid scenario: expected 3 free chapters")
	}

	// 可以使用这些场景进行更复杂的E2E测试
	// 例如：模拟用户阅读付费章节、发表评论、收藏书籍等

	t.Log("Integration test passed: all scenarios work correctly")
}

