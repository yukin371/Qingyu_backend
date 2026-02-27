//go:build e2e
// +build e2e

package data

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"Qingyu_backend/global"
	"Qingyu_backend/pkg/logger"
)

// ConsistencyIssue 一致性问题
type ConsistencyIssue struct {
	Type        string                 `json:"type"`        // 问题类型
	Description string                 `json:"description"` // 问题描述
	Severity    string                 `json:"severity"`    // 严重级别: error, warning
	Details     map[string]interface{} `json:"details"`     // 详细信息
}

// ConsistencyValidator 数据一致性验证器
type ConsistencyValidator struct {
	t *testing.T
}

// NewConsistencyValidator 创建一致性验证器
func NewConsistencyValidator(t *testing.T) *ConsistencyValidator {
	return &ConsistencyValidator{t: t}
}

// ValidateUserData 验证用户数据一致性
func (v *ConsistencyValidator) ValidateUserData(ctx context.Context, userID string) []ConsistencyIssue {
	var issues []ConsistencyIssue

	// 1. 检查用户基础信息存在
	user := v.getUser(ctx, userID)
	if user == nil {
		issues = append(issues, ConsistencyIssue{
			Type:        "user_not_found",
			Description: fmt.Sprintf("用户不存在: %s", userID),
			Severity:    "error",
			Details:     map[string]interface{}{"user_id": userID},
		})
		return issues
	}

	// 2. 验证阅读进度数据一致性（书籍存在）
	v.validateReadingProgress(ctx, userID, &issues)

	// 3. 验证社交数据一致性（评论/收藏的target存在）
	v.validateComments(ctx, userID, &issues)
	v.validateCollections(ctx, userID, &issues)

	// 4. 验证财务数据一致性（购买记录的书籍存在）
	v.validatePurchases(ctx, userID, &issues)

	return issues
}

// ValidateBookData 验证书籍数据一致性
func (v *ConsistencyValidator) ValidateBookData(ctx context.Context, bookID string) []ConsistencyIssue {
	var issues []ConsistencyIssue

	// 1. 检查书籍存在
	bookObjID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		issues = append(issues, ConsistencyIssue{
			Type:        "invalid_book_id",
			Description: fmt.Sprintf("无效的书籍ID: %s", bookID),
			Severity:    "error",
			Details:     map[string]interface{}{"book_id": bookID, "error": err.Error()},
		})
		return issues
	}

	book := v.getBook(ctx, bookObjID)
	if book == nil {
		issues = append(issues, ConsistencyIssue{
			Type:        "book_not_found",
			Description: fmt.Sprintf("书籍不存在: %s", bookID),
			Severity:    "error",
			Details:     map[string]interface{}{"book_id": bookID},
		})
		return issues
	}

	// 2. 验证章节数量与book.chapter_count一致
	chapters := v.getChaptersByBook(ctx, bookObjID)
	actualCount := len(chapters)
	chapterCount, ok := book["chapter_count"].(int32)
	if !ok {
		issues = append(issues, ConsistencyIssue{
			Type:        "invalid_chapter_count_type",
			Description: fmt.Sprintf("书籍chapter_count字段类型错误: %s", bookID),
			Severity:    "error",
			Details:     map[string]interface{}{"book_id": bookID, "actual_type": fmt.Sprintf("%T", book["chapter_count"])},
		})
		return issues
	}
	expectedCount := int(chapterCount)

	if actualCount != expectedCount {
		issues = append(issues, ConsistencyIssue{
			Type:        "chapter_count_mismatch",
			Description: fmt.Sprintf("章节数量不一致: 期望 %d, 实际 %d", expectedCount, actualCount),
			Severity:    "error",
			Details: map[string]interface{}{
				"book_id":        bookID,
				"expected_count": expectedCount,
				"actual_count":   actualCount,
			},
		})
	}

	// 3. 验证所有章节都有内容
	for _, chapter := range chapters {
		chapterID, ok := chapter["_id"].(primitive.ObjectID)
		if !ok {
			issues = append(issues, ConsistencyIssue{
				Type:        "invalid_chapter_id_type",
				Description: "章节_id字段类型错误",
				Severity:    "error",
				Details:     map[string]interface{}{"book_id": bookID, "chapter": chapter},
			})
			continue
		}

		hasContent := v.hasChapterContent(ctx, chapterID)
		if !hasContent {
			chapterNum, okNum := chapter["chapter_num"].(int32)
			chapterTitle, okTitle := chapter["title"].(string)

			if !okNum || !okTitle {
				issues = append(issues, ConsistencyIssue{
					Type:        "invalid_chapter_fields",
					Description: fmt.Sprintf("章节字段类型错误: %s", chapterID.Hex()),
					Severity:    "error",
					Details:     map[string]interface{}{"book_id": bookID, "chapter_id": chapterID.Hex(), "chapter": chapter},
				})
				continue
			}

			issues = append(issues, ConsistencyIssue{
				Type:        "missing_chapter_content",
				Description: fmt.Sprintf("章节缺少内容: %s (第%d章)", chapterID.Hex(), chapterNum),
				Severity:    "error",
				Details: map[string]interface{}{
					"book_id":       bookID,
					"chapter_id":    chapterID.Hex(),
					"chapter_num":   chapterNum,
					"chapter_title": chapterTitle,
				},
			})
		}
	}

	return issues
}

// AssertNoConsistencyIssues 断言没有一致性错误
func (v *ConsistencyValidator) AssertNoConsistencyIssues(t *testing.T, issues []ConsistencyIssue) {
	errorCount := 0
	warningCount := 0

	for _, issue := range issues {
		if issue.Severity == "error" {
			errorCount++
			t.Errorf("【一致性错误】%s: %s", issue.Type, issue.Description)
			if len(issue.Details) > 0 {
				for key, value := range issue.Details {
					t.Logf("  - %s: %v", key, value)
				}
			}
		} else if issue.Severity == "warning" {
			warningCount++
			t.Logf("【一致性警告】%s: %s", issue.Type, issue.Description)
		}
	}

	if errorCount == 0 && warningCount == 0 {
		t.Logf("✓ 数据一致性验证通过: 未发现问题")
	} else if errorCount == 0 && warningCount > 0 {
		t.Logf("⚠ 数据一致性验证完成: 发现 %d 个警告", warningCount)
	}
}

// ============ 私有辅助方法 ============

// getUser 获取用户
func (v *ConsistencyValidator) getUser(ctx context.Context, userID string) map[string]interface{} {
	var user map[string]interface{}

	// 兼容用户ID在不同集合中可能的存储类型（ObjectID 或 string）。
	filter := bson.M{"_id": userID}
	if objID, err := primitive.ObjectIDFromHex(userID); err == nil {
		filter = bson.M{
			"$or": []bson.M{
				{"_id": objID},
				{"_id": userID},
			},
		}
	}

	err := global.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil
	}
	return user
}

// getBook 获取书籍
func (v *ConsistencyValidator) getBook(ctx context.Context, bookID primitive.ObjectID) map[string]interface{} {
	var book map[string]interface{}
	err := global.DB.Collection("books").FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
	if err != nil {
		return nil
	}
	return book
}

// getChaptersByBook 获取书籍的所有章节
func (v *ConsistencyValidator) getChaptersByBook(ctx context.Context, bookID primitive.ObjectID) []map[string]interface{} {
	cursor, err := global.DB.Collection("chapters").Find(ctx, bson.M{"book_id": bookID})
	if err != nil {
		return []map[string]interface{}{}
	}
	defer cursor.Close(ctx)

	var chapters []map[string]interface{}
	cursor.All(ctx, &chapters)
	return chapters
}

// hasChapterContent 检查章节是否有内容
func (v *ConsistencyValidator) hasChapterContent(ctx context.Context, chapterID primitive.ObjectID) bool {
	count, err := global.DB.Collection("chapter_contents").CountDocuments(ctx, bson.M{"chapter_id": chapterID})
	if err != nil {
		return false
	}
	return count > 0
}

// validateReadingProgress 验证阅读进度
func (v *ConsistencyValidator) validateReadingProgress(ctx context.Context, userID string, issues *[]ConsistencyIssue) {
	cursor, err := global.DB.Collection("reading_progress").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Debug("Failed to query reading_progress",
			zap.String("user_id", userID),
			zap.String("error", err.Error()))
		return
	}
	defer cursor.Close(ctx)

	var progressList []map[string]interface{}
	cursor.All(ctx, &progressList)

	for _, progress := range progressList {
		bookID, ok := progress["book_id"].(string)
		if !ok {
			// book_id可能是ObjectID类型
			if bookObjID, ok := progress["book_id"].(primitive.ObjectID); ok {
				bookID = bookObjID.Hex()
			} else {
				*issues = append(*issues, ConsistencyIssue{
					Type:        "invalid_reading_progress_book_id_type",
					Description: "阅读进度book_id字段类型错误",
					Severity:    "warning",
					Details:     map[string]interface{}{"progress": progress},
				})
				continue
			}
		}

		// 检查书籍是否存在
		bookObjID, err := primitive.ObjectIDFromHex(bookID)
		if err != nil {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "invalid_reading_progress_book_id",
				Description: fmt.Sprintf("阅读进度引用的书籍ID无效: %s", bookID),
				Severity:    "warning",
				Details:     map[string]interface{}{"progress": progress},
			})
			continue
		}

		book := v.getBook(ctx, bookObjID)
		if book == nil {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "orphaned_reading_progress_book",
				Description: fmt.Sprintf("阅读进度引用的书籍不存在: %s", bookID),
				Severity:    "warning",
				Details:     map[string]interface{}{"progress": progress},
			})
		}
	}
}

// validateComments 验证评论
func (v *ConsistencyValidator) validateComments(ctx context.Context, userID string, issues *[]ConsistencyIssue) {
	cursor, err := global.DB.Collection("comments").Find(ctx, bson.M{"author_id": userID})
	if err != nil {
		logger.Debug("Failed to query comments",
			zap.String("user_id", userID),
			zap.String("error", err.Error()))
		return
	}
	defer cursor.Close(ctx)

	var comments []map[string]interface{}
	cursor.All(ctx, &comments)

	for _, comment := range comments {
		targetID, okID := comment["target_id"].(string)
		targetType, okType := comment["target_type"].(string)

		if !okID || !okType {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "invalid_comment_fields",
				Description: "评论target_id或target_type字段类型错误",
				Severity:    "warning",
				Details:     map[string]interface{}{"comment": comment},
			})
			continue
		}

		// 根据目标类型检查目标是否存在
		var exists bool
		if targetType == "book" {
			bookObjID, err := primitive.ObjectIDFromHex(targetID)
			if err != nil {
				exists = false
			} else {
				book := v.getBook(ctx, bookObjID)
				exists = book != nil
			}
		} else if targetType == "chapter" {
			chapterObjID, err := primitive.ObjectIDFromHex(targetID)
			if err != nil {
				exists = false
			} else {
				count, _ := global.DB.Collection("chapters").CountDocuments(ctx, bson.M{"_id": chapterObjID})
				exists = count > 0
			}
		} else {
			// 未知类型的评论，记录警告
			*issues = append(*issues, ConsistencyIssue{
				Type:        "unknown_comment_target_type",
				Description: fmt.Sprintf("评论引用了未知的目标类型: %s (类型: %s)", targetID, targetType),
				Severity:    "warning",
				Details:     map[string]interface{}{"comment": comment},
			})
			continue
		}

		if !exists {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "orphaned_comment_target",
				Description: fmt.Sprintf("评论引用的目标不存在: %s (类型: %s)", targetID, targetType),
				Severity:    "warning",
				Details:     map[string]interface{}{"comment": comment},
			})
		}
	}
}

// validateCollections 验证收藏
func (v *ConsistencyValidator) validateCollections(ctx context.Context, userID string, issues *[]ConsistencyIssue) {
	cursor, err := global.DB.Collection("collections").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Debug("Failed to query collections",
			zap.String("user_id", userID),
			zap.String("error", err.Error()))
		return
	}
	defer cursor.Close(ctx)

	var collections []map[string]interface{}
	cursor.All(ctx, &collections)

	for _, collection := range collections {
		bookID, ok := collection["book_id"].(string)
		if !ok {
			// book_id可能是ObjectID类型
			if bookObjID, ok := collection["book_id"].(primitive.ObjectID); ok {
				bookID = bookObjID.Hex()
			} else {
				*issues = append(*issues, ConsistencyIssue{
					Type:        "invalid_collection_book_id_type",
					Description: "收藏book_id字段类型错误",
					Severity:    "warning",
					Details:     map[string]interface{}{"collection": collection},
				})
				continue
			}
		}

		// 检查书籍是否存在
		bookObjID, err := primitive.ObjectIDFromHex(bookID)
		if err != nil {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "invalid_collection_book_id",
				Description: fmt.Sprintf("收藏引用的书籍ID无效: %s", bookID),
				Severity:    "warning",
				Details:     map[string]interface{}{"collection": collection},
			})
			continue
		}

		book := v.getBook(ctx, bookObjID)
		if book == nil {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "orphaned_collection_book",
				Description: fmt.Sprintf("收藏引用的书籍不存在: %s", bookID),
				Severity:    "warning",
				Details:     map[string]interface{}{"collection": collection},
			})
		}
	}
}

// validatePurchases 验证购买记录
func (v *ConsistencyValidator) validatePurchases(ctx context.Context, userID string, issues *[]ConsistencyIssue) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return
	}

	cursor, err := global.DB.Collection("chapter_purchases").Find(ctx, bson.M{"user_id": userObjID})
	if err != nil {
		logger.Debug("Failed to query chapter_purchases",
			zap.String("user_id", userID),
			zap.String("error", err.Error()))
		return
	}
	defer cursor.Close(ctx)

	var purchases []map[string]interface{}
	cursor.All(ctx, &purchases)

	for _, purchase := range purchases {
		bookID, ok := purchase["book_id"].(string)
		if !ok {
			// book_id可能是ObjectID
			if bookObjID, ok := purchase["book_id"].(primitive.ObjectID); ok {
				bookID = bookObjID.Hex()
			} else {
				continue
			}
		}

		// 检查书籍是否存在
		bookObjID, err := primitive.ObjectIDFromHex(bookID)
		if err != nil {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "invalid_purchase_book_id",
				Description: fmt.Sprintf("购买记录引用的书籍ID无效: %s", bookID),
				Severity:    "warning",
				Details:     map[string]interface{}{"purchase": purchase},
			})
			continue
		}

		book := v.getBook(ctx, bookObjID)
		if book == nil {
			*issues = append(*issues, ConsistencyIssue{
				Type:        "orphaned_purchase_book",
				Description: fmt.Sprintf("购买记录引用的书籍不存在: %s", bookID),
				Severity:    "warning",
				Details:     map[string]interface{}{"purchase": purchase},
			})
		}
	}
}
