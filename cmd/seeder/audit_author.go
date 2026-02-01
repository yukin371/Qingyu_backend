// Package main provides author perspective data relation audit tool
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// AuthorAuditReport 作者视角数据审查报告
type AuthorAuditReport struct {
	AuditTime     time.Time
	TotalScore    string
	Summary       string
	Findings      []Finding
}

// Finding 发现的问题
type Finding struct {
	Category    string
	Count       int64
	Description string
	Examples    []bson.M
}

// RunAuthorAudit 执行作者视角数据关联审查
func RunAuthorAudit(db *utils.Database) (*AuthorAuditReport, error) {
	report := &AuthorAuditReport{
		AuditTime: time.Now(),
		Findings:  []Finding{},
	}

	ctx := context.Background()

	// 1. 外键关系完整性验证
	fmt.Println("🔍 验证外键关系完整性...")

	// 1.1 书籍-章节关系验证
	fmt.Println("  1.1 验证书籍-章节关系...")
	finding, err := auditBookChapterRelation(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("书籍-章节关系验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 1.2 书籍-作者关系验证
	fmt.Println("  1.2 验证书籍-作者关系...")
	finding, err = auditBookAuthorRelation(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("书籍-作者关系验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 1.3 章节内容完整性验证
	fmt.Println("  1.3 验证章节内容完整性...")
	finding, err = auditChapterContentRelation(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("章节内容完整性验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 1.4 收益记录关联验证
	fmt.Println("  1.4 验证收益记录关联...")
	finding, err = auditRevenueRelation(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("收益记录关联验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 2. 统计数据准确性验证
	fmt.Println("🔍 验证统计数据准确性...")

	// 2.1 章节数量统计验证
	fmt.Println("  2.1 验证章节数量统计...")
	finding, err = auditChapterCountAccuracy(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("章节数量统计验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 2.2 字数统计验证
	fmt.Println("  2.2 验证字数统计...")
	finding, err = auditWordCountAccuracy(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("字数统计验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 3. 业务规则一致性验证
	fmt.Println("🔍 验证业务规则一致性...")

	// 3.1 章节序号唯一性验证
	fmt.Println("  3.1 验证章节序号唯一性...")
	finding, err = auditChapterNumberUniqueness(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("章节序号唯一性验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 3.2 免费章节定价规则验证
	fmt.Println("  3.2 验证免费章节定价规则...")
	finding, err = auditFreeChapterPricing(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("免费章节定价规则验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 3.3 状态一致性验证
	fmt.Println("  3.3 验证状态一致性...")
	finding, err = auditStatusConsistency(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("状态一致性验证失败: %w", err)
	}
	report.Findings = append(report.Findings, finding)

	// 计算总体评分
	report.TotalScore = calculateTotalScore(report.Findings)
	report.Summary = generateSummary(report.Findings)

	return report, nil
}

// auditBookChapterRelation 验证书籍-章节关系
func auditBookChapterRelation(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "书籍-章节关系",
		Description: "外键关系完整性",
		Examples:    []bson.M{},
	}

	// 查询无章节的书籍
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "chapters",
				"localField":   "_id",
				"foreignField": "book_id",
				"as":           "chapters",
			},
		},
		{
			"$match": bson.M{
				"chapters": bson.M{"$size": 0},
			},
		},
		{
			"$project": bson.M{
				"_id":   1,
				"title": 1,
			},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err := db.Collection("books").Aggregate(ctx, pipeline)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		return finding, err
	}

	finding.Count = int64(len(books))
	if len(books) > 0 {
		finding.Examples = books
	}

	return finding, nil
}

// auditBookAuthorRelation 验证书籍-作者关系
func auditBookAuthorRelation(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "书籍-作者关系",
		Description: "孤儿作者记录（作者不存在）",
		Examples:    []bson.M{},
	}

	// 查询作者不存在的书籍
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "author_id",
				"foreignField": "_id",
				"as":           "author",
			},
		},
		{
			"$match": bson.M{
				"author": bson.M{"$size": 0},
			},
		},
		{
			"$project": bson.M{
				"_id":        1,
				"title":      1,
				"author_id":  1,
			},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err := db.Collection("books").Aggregate(ctx, pipeline)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		return finding, err
	}

	finding.Count = int64(len(books))
	if len(books) > 0 {
		finding.Examples = books
	}

	return finding, nil
}

// auditChapterContentRelation 验证章节内容完整性
func auditChapterContentRelation(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "章节-内容关系",
		Description: "无内容的章节",
		Examples:    []bson.M{},
	}

	// 查询无内容的章节
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "chapter_contents",
				"localField":   "_id",
				"foreignField": "chapter_id",
				"as":           "content",
			},
		},
		{
			"$match": bson.M{
				"content": bson.M{"$size": 0},
			},
		},
		{
			"$project": bson.M{
				"_id":        1,
				"title":      1,
				"book_id":    1,
				"chapter_num": 1,
			},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err := db.Collection("chapters").Aggregate(ctx, pipeline)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var chapters []bson.M
	if err = cursor.All(ctx, &chapters); err != nil {
		return finding, err
	}

	finding.Count = int64(len(chapters))
	if len(chapters) > 0 {
		finding.Examples = chapters
	}

	return finding, nil
}

// auditRevenueRelation 验证收益记录关联
func auditRevenueRelation(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "收益记录关联",
		Description: "孤儿收益记录（用户或书籍不存在）",
		Examples:    []bson.M{},
	}

	// 查询孤儿收益记录
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "books",
				"localField":   "book_id",
				"foreignField": "_id",
				"as":           "book",
			},
		},
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"user": bson.M{"$size": 0}},
					{"book": bson.M{"$size": 0}},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":     1,
				"user_id": 1,
				"book_id": 1,
				"amount":  1,
			},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err := db.Collection("author_revenue").Aggregate(ctx, pipeline)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var revenues []bson.M
	if err = cursor.All(ctx, &revenues); err != nil {
		return finding, err
	}

	finding.Count = int64(len(revenues))
	if len(revenues) > 0 {
		finding.Examples = revenues
	}

	return finding, nil
}

// auditChapterCountAccuracy 验证章节数量统计准确性
func auditChapterCountAccuracy(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "章节数量统计",
		Description: "章节计数不一致的书籍",
		Examples:    []bson.M{},
	}

	// 查询章节数量不一致的书籍
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "chapters",
				"localField":   "_id",
				"foreignField": "book_id",
				"as":           "chapters",
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"title":        1,
				"stored_count": "$chapter_count",
				"actual_count": bson.M{"$size": "$chapters"},
				"diff": bson.M{
					"$subtract": []interface{}{
						"$chapter_count",
						bson.M{"$size": "$chapters"},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"diff": bson.M{"$ne": 0},
			},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err := db.Collection("books").Aggregate(ctx, pipeline)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		return finding, err
	}

	// 先获取总数
	countPipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "chapters",
				"localField":   "_id",
				"foreignField": "book_id",
				"as":           "chapters",
			},
		},
		{
			"$project": bson.M{
				"diff": bson.M{
					"$subtract": []interface{}{
						"$chapter_count",
						bson.M{"$size": "$chapters"},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"diff": bson.M{"$ne": 0},
			},
		},
		{
			"$count": "total",
		},
	}

	countCursor, err := db.Collection("books").Aggregate(ctx, countPipeline)
	if err != nil {
		return finding, err
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err = countCursor.All(ctx, &countResult); err == nil && len(countResult) > 0 {
		finding.Count = countResult[0]["total"].(int64)
	}

	finding.Examples = books

	return finding, nil
}

// auditWordCountAccuracy 验证字数统计准确性
func auditWordCountAccuracy(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "字数统计",
		Description: "字数统计异常的书籍（为0或负数）",
		Examples:    []bson.M{},
	}

	// 查询字数统计异常的书籍
	filter := bson.M{
		"$or": []bson.M{
			{"word_count": bson.M{"$lt": 0}},
			{"word_count": bson.M{"$exists": false}},
		},
	}

	cursor, err := db.Collection("books").Find(ctx, filter)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		return finding, err
	}

	// 获取总数
	count, err := db.Collection("books").CountDocuments(ctx, filter)
	if err != nil {
		return finding, err
	}

	finding.Count = count
	if len(books) > 10 {
		finding.Examples = books[:10]
	} else {
		finding.Examples = books
	}

	return finding, nil
}

// auditChapterNumberUniqueness 验证章节序号唯一性
func auditChapterNumberUniqueness(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "章节序号唯一性",
		Description: "重复章节号（同一本书内）",
		Examples:    []bson.M{},
	}

	// 查询重复的章节号
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"book_id":     "$book_id",
					"chapter_num": "$chapter_num",
				},
				"count": bson.M{"$sum": 1},
				"chapter_ids": bson.M{"$push": "$_id"},
			},
		},
		{
			"$match": bson.M{
				"count": bson.M{"$gt": 1},
			},
		},
		{
			"$limit": 10,
		},
	}

	cursor, err := db.Collection("chapters").Aggregate(ctx, pipeline)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var duplicates []bson.M
	if err = cursor.All(ctx, &duplicates); err != nil {
		return finding, err
	}

	// 获取总数
	countPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"book_id":     "$book_id",
					"chapter_num": "$chapter_num",
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$match": bson.M{
				"count": bson.M{"$gt": 1},
			},
		},
		{
			"$count": "total",
		},
	}

	countCursor, err := db.Collection("chapters").Aggregate(ctx, countPipeline)
	if err != nil {
		return finding, err
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	if err = countCursor.All(ctx, &countResult); err == nil && len(countResult) > 0 {
		finding.Count = countResult[0]["total"].(int64)
	}

	finding.Examples = duplicates

	return finding, nil
}

// auditFreeChapterPricing 验证免费章节定价规则
func auditFreeChapterPricing(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "免费章节定价规则",
		Description: "免费但价格>0的章节",
		Examples:    []bson.M{},
	}

	// 查询免费但价格>0的章节
	filter := bson.M{
		"is_free": true,
		"price":   bson.M{"$gt": 0},
	}

	cursor, err := db.Collection("chapters").Find(ctx, filter)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var chapters []bson.M
	if err = cursor.All(ctx, &chapters); err != nil {
		return finding, err
	}

	// 获取总数
	count, err := db.Collection("chapters").CountDocuments(ctx, filter)
	if err != nil {
		return finding, err
	}

	finding.Count = count
	if len(chapters) > 10 {
		finding.Examples = chapters[:10]
	} else {
		finding.Examples = chapters
	}

	return finding, nil
}

// auditStatusConsistency 验证状态一致性
func auditStatusConsistency(ctx context.Context, db *utils.Database) (Finding, error) {
	finding := Finding{
		Category:    "状态一致性",
		Description: "书籍状态异常（已完结但更新时间很近）",
		Examples:    []bson.M{},
	}

	// 查询状态不一致的书籍
	filter := bson.M{
		"status": "completed",
		"updated_at": bson.M{
			"$gte": time.Now().AddDate(0, 0, -7), // 最近7天更新
		},
	}

	cursor, err := db.Collection("books").Find(ctx, filter)
	if err != nil {
		return finding, err
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		return finding, err
	}

	// 获取总数
	count, err := db.Collection("books").CountDocuments(ctx, filter)
	if err != nil {
		return finding, err
	}

	finding.Count = count
	if len(books) > 10 {
		finding.Examples = books[:10]
	} else {
		finding.Examples = books
	}

	return finding, nil
}

// calculateTotalScore 计算总体评分
func calculateTotalScore(findings []Finding) string {
	totalIssues := int64(0)
	for _, f := range findings {
		totalIssues += f.Count
	}

	if totalIssues == 0 {
		return "优秀 (A)"
	} else if totalIssues < 10 {
		return "良好 (B)"
	} else if totalIssues < 50 {
		return "一般 (C)"
	}
	return "差 (D)"
}

// generateSummary 生成总结
func generateSummary(findings []Finding) string {
	criticalIssues := 0
	totalIssues := int64(0)

	for _, f := range findings {
		totalIssues += f.Count
		if f.Count > 0 {
			criticalIssues++
		}
	}

	if totalIssues == 0 {
		return "✅ 所有数据关联完整，未发现异常"
	}

	return fmt.Sprintf("⚠️ 发现 %d 类问题，共 %d 条异常数据需要处理", criticalIssues, totalIssues)
}

// PrintReport 打印审查报告
func (r *AuthorAuditReport) PrintReport() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 作者视角数据关联审查报告")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("审查时间: %s\n", r.AuditTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("总体评分: %s\n\n", r.TotalScore)

	// 外键关系完整性
	fmt.Println("📌 外键关系完整性")
	fmt.Println(strings.Repeat("-", 80))
	for _, f := range r.Findings {
		if f.Category == "书籍-章节关系" || f.Category == "书籍-作者关系" ||
			f.Category == "章节-内容关系" || f.Category == "收益记录关联" {
			fmt.Printf("  %s: %d 条\n", f.Category, f.Count)
			if f.Count > 0 && len(f.Examples) > 0 {
				fmt.Printf("    描述: %s\n", f.Description)
				fmt.Printf("    示例: %v\n", f.Examples[0])
			}
		}
	}

	// 统计数据准确性
	fmt.Println("\n📌 统计数据准确性")
	fmt.Println(strings.Repeat("-", 80))
	for _, f := range r.Findings {
		if f.Category == "章节数量统计" || f.Category == "字数统计" {
			fmt.Printf("  %s: %d 条\n", f.Category, f.Count)
			if f.Count > 0 && len(f.Examples) > 0 {
				fmt.Printf("    描述: %s\n", f.Description)
				fmt.Printf("    示例: %v\n", f.Examples[0])
			}
		}
	}

	// 业务规则一致性
	fmt.Println("\n📌 业务规则一致性")
	fmt.Println(strings.Repeat("-", 80))
	for _, f := range r.Findings {
		if f.Category == "章节序号唯一性" || f.Category == "免费章节定价规则" ||
			f.Category == "状态一致性" {
			fmt.Printf("  %s: %d 条\n", f.Category, f.Count)
			if f.Count > 0 && len(f.Examples) > 0 {
				fmt.Printf("    描述: %s\n", f.Description)
				fmt.Printf("    示例: %v\n", f.Examples[0])
			}
		}
	}

	// 数据质量评估
	fmt.Println("\n📌 数据质量评估")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("  整体评分: %s\n", r.TotalScore)
	fmt.Printf("  总结: %s\n", r.Summary)

	fmt.Println("\n" + strings.Repeat("=", 80))
}
