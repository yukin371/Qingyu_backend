// Seeder 工具 - 读者视角数据关联审查工具
// 用于验证 reading_progress, reading_history, bookmarks, annotations, book_lists 等集合的数据完整性
package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditResult 审查结果
type AuditResult struct {
	Name        string
	Description string
	Count       int64
	Samples     []bson.M
	Error       error
}

// ReaderAuditor 读者视角数据审查器
type ReaderAuditor struct {
	db *mongo.Database
}

// NewReaderAuditor 创建审查器
func NewReaderAuditor(db *mongo.Database) *ReaderAuditor {
	return &ReaderAuditor{db: db}
}

// AuditAll 执行所有审查
func (a *ReaderAuditor) AuditAll(ctx context.Context) []*AuditResult {
	var results []*AuditResult

	fmt.Println("🔍 开始读者视角数据关联审查...")
	fmt.Println("====================================")

	// 外键关系完整性检查
	results = append(results, a.auditReadingProgressOrphanUsers(ctx)...)
	results = append(results, a.auditReadingProgressOrphanBooks(ctx)...)
	results = append(results, a.auditReadingProgressDuplicates(ctx)...)
	results = append(results, a.auditReadingHistoryOrphans(ctx)...)
	results = append(results, a.auditBookmarksOrphans(ctx)...)
	results = append(results, a.auditAnnotationsOrphans(ctx)...)
	results = append(results, a.auditBookListsOrphanUsers(ctx)...)

	// 业务规则一致性检查
	results = append(results, a.auditReadingProgressRange(ctx)...)
	results = append(results, a.auditReadingDuration(ctx)...)
	results = append(results, a.auditBookmarkPosition(ctx)...)

	return results
}

// ============ reading_progress 检查 ============

// 验证 user_id → users._id
func (a *ReaderAuditor) auditReadingProgressOrphanUsers(ctx context.Context) []*AuditResult {
	fmt.Println("\n📖 检查 reading_progress 孤儿用户...")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "user", Value: bson.D{{Key: "$size", Value: 0}}},
		}}},
	}

	return a.executeAudit(ctx, "reading_progress", "孤儿用户 (user_id 不存在)", pipeline, 5)
}

// 验证 book_id → books._id
func (a *ReaderAuditor) auditReadingProgressOrphanBooks(ctx context.Context) []*AuditResult {
	fmt.Println("📖 检查 reading_progress 孤儿书籍...")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "books"},
			{Key: "localField", Value: "book_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "book"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "book", Value: bson.D{{Key: "$size", Value: 0}}},
		}}},
	}

	return a.executeAudit(ctx, "reading_progress", "孤儿书籍 (book_id 不存在)", pipeline, 5)
}

// 验证唯一性（同一用户-书籍不应有多条记录）
func (a *ReaderAuditor) auditReadingProgressDuplicates(ctx context.Context) []*AuditResult {
	fmt.Println("📖 检查 reading_progress 重复记录...")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "user_id", Value: "$user_id"},
				{Key: "book_id", Value: "$book_id"},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "docs", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}},
		}}},
	}

	return a.executeAudit(ctx, "reading_progress", "重复记录 (同一用户-书籍组合)", pipeline, 5)
}

// 验证阅读进度范围 [0, 100]
func (a *ReaderAuditor) auditReadingProgressRange(ctx context.Context) []*AuditResult {
	fmt.Println("📖 检查 reading_progress 进度范围...")

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "progress", Value: bson.D{{Key: "$lt", Value: 0}}}},
			bson.D{{Key: "progress", Value: bson.D{{Key: "$gt", Value: 100}}}},
		}},
	}

	return a.executeFindAudit(ctx, "reading_progress", "阅读进度超出范围 [0,100]", filter, 5)
}

// ============ reading_history 检查 ============

// 验证 user_id, book_id, chapter_id
func (a *ReaderAuditor) auditReadingHistoryOrphans(ctx context.Context) []*AuditResult {
	fmt.Println("\n📚 检查 reading_history 孤儿记录...")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "books"},
			{Key: "localField", Value: "book_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "book"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "chapters"},
			{Key: "localField", Value: "chapter_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "chapter"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "user", Value: bson.D{{Key: "$size", Value: 0}}}},
				bson.D{{Key: "book", Value: bson.D{{Key: "$size", Value: 0}}}},
				bson.D{{Key: "chapter", Value: bson.D{{Key: "$size", Value: 0}}}},
			}},
		}}},
	}

	return a.executeAudit(ctx, "reading_history", "孤儿记录 (外键不存在)", pipeline, 5)
}

// 验证阅读时长合理性
func (a *ReaderAuditor) auditReadingDuration(ctx context.Context) []*AuditResult {
	fmt.Println("📚 检查 reading_history 阅读时长...")

	filter := bson.D{{Key: "duration", Value: bson.D{{Key: "$lt", Value: 0}}}}

	return a.executeFindAudit(ctx, "reading_history", "阅读时长为负数", filter, 5)
}

// ============ bookmarks 检查 ============

// 验证 user_id, book_id, chapter_id
func (a *ReaderAuditor) auditBookmarksOrphans(ctx context.Context) []*AuditResult {
	fmt.Println("\n🔖 检查 bookmarks 孤儿记录...")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "books"},
			{Key: "localField", Value: "book_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "book"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "chapters"},
			{Key: "localField", Value: "chapter_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "chapter"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "user", Value: bson.D{{Key: "$size", Value: 0}}}},
				bson.D{{Key: "book", Value: bson.D{{Key: "$size", Value: 0}}}},
				bson.D{{Key: "chapter", Value: bson.D{{Key: "$size", Value: 0}}}},
			}},
		}}},
	}

	return a.executeAudit(ctx, "bookmarks", "孤儿记录 (外键不存在)", pipeline, 5)
}

// 验证书签位置合理性
func (a *ReaderAuditor) auditBookmarkPosition(ctx context.Context) []*AuditResult {
	fmt.Println("🔖 检查 bookmarks 位置合理性...")

	filter := bson.D{{Key: "position", Value: bson.D{{Key: "$lt", Value: 0}}}}

	return a.executeFindAudit(ctx, "bookmarks", "书签位置为负数", filter, 5)
}

// ============ annotations 检查 ============

// 验证 user_id, book_id, chapter_id
func (a *ReaderAuditor) auditAnnotationsOrphans(ctx context.Context) []*AuditResult {
	fmt.Println("\n✏️ 检查 annotations 孤儿记录...")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "books"},
			{Key: "localField", Value: "book_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "book"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "chapters"},
			{Key: "localField", Value: "chapter_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "chapter"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "user", Value: bson.D{{Key: "$size", Value: 0}}}},
				bson.D{{Key: "book", Value: bson.D{{Key: "$size", Value: 0}}}},
				bson.D{{Key: "chapter", Value: bson.D{{Key: "$size", Value: 0}}}},
			}},
		}}},
	}

	return a.executeAudit(ctx, "annotations", "孤儿记录 (外键不存在)", pipeline, 5)
}

// ============ book_lists 检查 ============

// 验证 user_id → users._id
func (a *ReaderAuditor) auditBookListsOrphanUsers(ctx context.Context) []*AuditResult {
	fmt.Println("\n📋 检查 book_lists 孤儿用户...")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "user", Value: bson.D{{Key: "$size", Value: 0}}},
		}}},
	}

	return a.executeAudit(ctx, "book_lists", "孤儿用户 (user_id 不存在)", pipeline, 5)
}

// ============ 辅助方法 ============

// executeAudit 执行聚合管道审查
func (a *ReaderAuditor) executeAudit(ctx context.Context, collection, description string, pipeline mongo.Pipeline, sampleLimit int) []*AuditResult {
	coll := a.db.Collection(collection)

	// 先获取总数
	countPipeline := append(pipeline, bson.D{{Key: "$count", Value: "count"}})
	countCursor, err := coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return []*AuditResult{{
			Name:        collection,
			Description: description,
			Count:       -1,
			Error:       err,
		}}
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err = countCursor.All(ctx, &countResult); err != nil {
		return []*AuditResult{{
			Name:        collection,
			Description: description,
			Count:       -1,
			Error:       err,
		}}
	}

	count := int64(0)
	if len(countResult) > 0 {
		if c, ok := countResult[0]["count"].(int32); ok {
			count = int64(c)
		} else if c, ok := countResult[0]["count"].(int64); ok {
			count = c
		}
	}

	// 获取样本
	sampleCursor, err := coll.Aggregate(ctx, append(pipeline, bson.D{{Key: "$limit", Value: sampleLimit}}))
	if err != nil {
		return []*AuditResult{{
			Name:        collection,
			Description: description,
			Count:       count,
			Error:       err,
		}}
	}
	defer sampleCursor.Close(ctx)

	var samples []bson.M
	if err = sampleCursor.All(ctx, &samples); err != nil {
		samples = nil
	}

	result := &AuditResult{
		Name:        collection,
		Description: description,
		Count:       count,
		Samples:     samples,
	}

	return []*AuditResult{result}
}

// executeFindAudit 执行简单查询审查
func (a *ReaderAuditor) executeFindAudit(ctx context.Context, collection, description string, filter bson.D, sampleLimit int) []*AuditResult {
	coll := a.db.Collection(collection)

	// 获取总数
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return []*AuditResult{{
			Name:        collection,
			Description: description,
			Count:       -1,
			Error:       err,
		}}
	}

	// 获取样本
	cursor, err := coll.Find(ctx, filter, options.Find().SetLimit(int64(sampleLimit)))
	if err != nil {
		return []*AuditResult{{
			Name:        collection,
			Description: description,
			Count:       count,
			Error:       err,
		}}
	}
	defer cursor.Close(ctx)

	var samples []bson.M
	if err = cursor.All(ctx, &samples); err != nil {
		samples = nil
	}

	return []*AuditResult{{
		Name:        collection,
		Description: description,
		Count:       count,
		Samples:     samples,
	}}
}

// ============ 集合存在性检查 ============

// CheckCollectionExists 检查集合是否存在
func (a *ReaderAuditor) CheckCollectionExists(ctx context.Context, collectionName string) bool {
	collections, err := a.db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return false
	}
	return len(collections) > 0
}

// GetCollectionStats 获取集合统计信息
func (a *ReaderAuditor) GetCollectionStats(ctx context.Context, collectionName string) (string, int64) {
	if !a.CheckCollectionExists(ctx, collectionName) {
		return "不存在", 0
	}
	count, _ := a.db.Collection(collectionName).CountDocuments(ctx, bson.M{})
	return "存在", count
}

// IDFormatCheck 检查 ID 格式一致性
func (a *ReaderAuditor) IDFormatCheck(ctx context.Context) map[string][]string {
	fmt.Println("\n🔍 检查各集合的 ID 格式...")

	result := make(map[string][]string)

	collections := []string{"users", "books", "reading_progress", "reading_history", "bookmarks", "annotations", "book_lists"}
	idFields := map[string][]string{
		"users":             {"_id"},
		"books":             {"_id", "author_id"},
		"reading_progress":  {"_id", "user_id", "book_id"},
		"reading_history":   {"_id", "user_id", "book_id", "chapter_id"},
		"bookmarks":         {"_id", "user_id", "book_id", "chapter_id"},
		"annotations":       {"_id", "user_id", "book_id", "chapter_id"},
		"book_lists":        {"_id", "user_id"},
		"book_list_items":   {"_id", "book_id", "list_id"},
	}

	for _, coll := range collections {
		if !a.CheckCollectionExists(ctx, coll) {
			result[coll] = []string{"集合不存在"}
			continue
		}

		var formats []string
		for _, field := range idFields[coll] {
			// 获取样本 - 使用 bson.M 构造动态字段名
			matchStage := bson.M{field: bson.M{"$exists": true}}
			projectStage := bson.M{field: 1}
			cursor, _ := a.db.Collection(coll).Aggregate(ctx, mongo.Pipeline{
				bson.D{{Key: "$match", Value: matchStage}},
				bson.D{{Key: "$limit", Value: 10}},
				bson.D{{Key: "$project", Value: projectStage}},
			})

			var samples []bson.M
			cursor.All(ctx, &samples)
			cursor.Close(ctx)

			if len(samples) == 0 {
				formats = append(formats, fmt.Sprintf("%s: 无数据", field))
				continue
			}

			// 检查 ID 类型
			var hasObjectID, hasString bool
			for _, sample := range samples {
				if val, ok := sample[field]; ok {
					switch val.(type) {
					case primitive.ObjectID:
						hasObjectID = true
					case string:
						hasString = true
					}
				}
			}

			if hasObjectID && hasString {
				formats = append(formats, fmt.Sprintf("%s: 混合格式 (ObjectID + String)", field))
			} else if hasObjectID {
				formats = append(formats, fmt.Sprintf("%s: ObjectID", field))
			} else if hasString {
				formats = append(formats, fmt.Sprintf("%s: String", field))
			} else {
				formats = append(formats, fmt.Sprintf("%s: 未知格式", field))
			}
		}

		result[coll] = formats
	}

	return result
}

// GenerateReport 生成审查报告
func (a *ReaderAuditor) GenerateReport(results []*AuditResult, idFormats map[string][]string) string {
	report := "# 读者视角数据关联审查报告\n\n"
	report += fmt.Sprintf("**审查时间:** %s\n", time.Now().Format("2006-01-02 15:04:05"))
	report += "**审查人:** 数据审查专家女仆\n\n"

	// 集合统计概览
	report += "## 📊 集合统计概览\n\n"
	ctx := context.Background()
	collections := []string{
		"users", "books", "chapters", "reading_progress",
		"reading_history", "bookmarks", "annotations", "book_lists", "book_list_items",
	}

	report += "| 集合名称 | 状态 | 文档数量 |\n"
	report += "|---------|------|----------|\n"
	for _, coll := range collections {
		status, count := a.GetCollectionStats(ctx, coll)
		report += fmt.Sprintf("| %s | %s | %d |\n", coll, status, count)
	}
	report += "\n"

	// ID 格式检查结果
	report += "## 🔍 ID 格式检查结果\n\n"
	report += "这是本次审查的重点，因为 Seeder 工具之前发现了 ID 类型不匹配问题。\n\n"
	report += "| 集合 | ID 字段 | 格式 |\n"
	report += "|------|---------|------|\n"

	for coll, formats := range idFormats {
		for _, format := range formats {
			report += fmt.Sprintf("| %s | %s |\n", coll, format)
		}
	}
	report += "\n"

	// 外键关系完整性
	report += "## 🔗 外键关系完整性\n\n"

	// 分类结果
	foreignKeyResults := make(map[string][]*AuditResult)
	businessRuleResults := make(map[string][]*AuditResult)

	for _, result := range results {
		if result.Error != nil {
			report += fmt.Sprintf("### ❌ %s - %s\n", result.Name, result.Description)
			report += fmt.Sprintf("**错误:** %v\n\n", result.Error)
			continue
		}

		// 分类
		if contains(result.Description, "孤儿") || contains(result.Description, "重复") {
			foreignKeyResults[result.Name] = append(foreignKeyResults[result.Name], result)
		} else {
			businessRuleResults[result.Name] = append(businessRuleResults[result.Name], result)
		}
	}

	// 输出外键关系结果
	for _, coll := range []string{"reading_progress", "reading_history", "bookmarks", "annotations", "book_lists"} {
		if results, ok := foreignKeyResults[coll]; ok && len(results) > 0 {
			report += fmt.Sprintf("### %s\n", coll)
			for _, r := range results {
				if r.Count > 0 {
					report += fmt.Sprintf("- ⚠️ **%s:** %d 条异常\n", r.Description, r.Count)
					if len(r.Samples) > 0 {
						report += fmt.Sprintf("  - 示例: ` %+v`\n", r.Samples[0])
					}
				} else {
					report += fmt.Sprintf("- ✅ **%s:** 无异常\n", r.Description)
				}
			}
			report += "\n"
		}
	}

	// 业务规则一致性
	report += "## 📋 业务规则一致性\n\n"

	for _, coll := range []string{"reading_progress", "reading_history", "bookmarks", "annotations", "book_lists"} {
		if results, ok := businessRuleResults[coll]; ok && len(results) > 0 {
			report += fmt.Sprintf("### %s\n", coll)
			for _, r := range results {
				if r.Count > 0 {
					report += fmt.Sprintf("- ⚠️ **%s:** %d 条异常\n", r.Description, r.Count)
					if len(r.Samples) > 0 {
						report += fmt.Sprintf("  - 示例: ` %+v`\n", r.Samples[0])
					}
				} else {
					report += fmt.Sprintf("- ✅ **%s:** 无异常\n", r.Description)
				}
			}
			report += "\n"
		}
	}

	// 数据质量评估
	report += "## 📈 数据质量评估\n\n"

	totalIssues := 0
	criticalIssues := 0

	for _, result := range results {
		if result.Error == nil && result.Count > 0 {
			totalIssues += int(result.Count)
			if contains(result.Description, "孤儿") || contains(result.Description, "重复") {
				criticalIssues += int(result.Count)
			}
		}
	}

	score := "优秀"
	if criticalIssues > 0 {
		score = "差"
	} else if totalIssues > 100 {
		score = "一般"
	} else if totalIssues > 0 {
		score = "良好"
	}

	report += fmt.Sprintf("**整体评分:** %s\n\n", score)
	report += fmt.Sprintf("**总问题数:** %d 条\n", totalIssues)
	report += fmt.Sprintf("**严重问题:** %d 条 (孤儿记录/重复记录)\n\n", criticalIssues)

	if totalIssues > 0 {
		report += "### 主要问题汇总\n\n"

		// 按严重程度分类
		if criticalIssues > 0 {
			report += "#### 🔴 严重问题\n\n"
			for _, result := range results {
				if result.Error == nil && result.Count > 0 &&
					(contains(result.Description, "孤儿") || contains(result.Description, "重复")) {
					report += fmt.Sprintf("- **%s.%s**: %d 条\n", result.Name, result.Description, result.Count)
				}
			}
			report += "\n"
		}

		otherIssues := totalIssues - criticalIssues
		if otherIssues > 0 {
			report += "#### 🟡 一般问题\n\n"
			for _, result := range results {
				if result.Error == nil && result.Count > 0 &&
					!contains(result.Description, "孤儿") && !contains(result.Description, "重复") {
					report += fmt.Sprintf("- **%s.%s**: %d 条\n", result.Name, result.Description, result.Count)
				}
			}
			report += "\n"
		}
	} else {
		report += "🎉 **未发现数据问题！**\n\n"
	}

	// 修复建议
	if totalIssues > 0 {
		report += "## 🔧 修复建议\n\n"

		if criticalIssues > 0 {
			report += "### 严重问题修复\n\n"
			report += "1. **孤儿记录处理**\n"
			report += "   - 删除无效的孤儿记录\n"
			report += "   - 或重新建立关联关系\n\n"

			report += "2. **重复记录处理**\n"
			report += "   - 保留最新的记录\n"
			report += "   - 删除旧的重复记录\n\n"

			report += "3. **ID 格式统一** (如果发现问题)\n"
			report += "   - 统一使用 ObjectID 或 String 类型\n"
			report += "   - 更新相关代码和 Seeder 工具\n\n"
		}

		if totalIssues-criticalIssues > 0 {
			report += "### 一般问题修复\n\n"
			report += "1. **数据范围修正**\n"
			report += "   - 将超出范围的值调整到合理范围\n"
			report += "   - 添加数据验证规则\n\n"

			report += "2. **业务规则完善**\n"
			report += "   - 在应用层添加更严格的数据验证\n"
			report += "   - 在数据库层添加约束（如果支持）\n\n"
		}

		report += "### 预防措施\n\n"
		report += "1. **增强 Seeder 工具**\n"
		report += "   - 添加数据关联验证\n"
		report += "   - 统一 ID 类型使用\n"
		report += "   - 添加外键关系检查\n\n"

		report += "2. **添加监控**\n"
		report += "   - 定期执行数据完整性检查\n"
		report += "   - 添加数据质量监控告警\n\n"
	}

	report += "---\n\n"
	report += "*本报告由数据审查专家女仆自动生成*\n"

	return report
}

// contains 辅助函数：检查字符串包含
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============ 主函数入口 ============

// RunReaderAudit 执行读者视角数据关联审查
func RunReaderAudit(db *mongo.Database) error {
	auditor := NewReaderAuditor(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 执行 ID 格式检查
	idFormats := auditor.IDFormatCheck(ctx)

	// 执行所有审查
	results := auditor.AuditAll(ctx)

	// 生成报告
	report := auditor.GenerateReport(results, idFormats)

	// 输出到控制台
	fmt.Println("\n" + report)

	return nil
}
