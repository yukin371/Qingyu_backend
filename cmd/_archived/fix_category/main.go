// Package main 提供书籍分类关联修复功能
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var allBooks []bson.M // 全局变量，用于存储所有书籍

// CategoryFixReport 修复报告
type CategoryFixReport struct {
	TotalBooks         int64
	InvalidBooks       int64
	TotalCategories    int64
	BookCategories     map[string]int // 书籍使用的分类及数量
	CategoryNames      []string       // 分类表中的分类名称
	UnmatchedBooks     []string       // 未匹配的书籍
	FixesApplied       []string       // 执行的修复操作
	ProblemsFound      []string       // 发现的问题
	FixedCount         int64          // 修复的书籍数量
	VerificationPassed bool           // 验证是否通过
}

func main() {
	fmt.Println("========================================")
	fmt.Println("   书籍分类关联诊断与修复工具")
	fmt.Println("========================================")
	fmt.Println()

	// 连接数据库
	fmt.Println("正在连接数据库...")
	uri := "mongodb://localhost:27017"
	dbName := "qingyu"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("❌ 数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	booksCollection := db.Collection("books")
	categoriesCollection := db.Collection("categories")

	fmt.Println("✅ 数据库连接成功")
	fmt.Println()

	// 创建报告
	report := &CategoryFixReport{
		BookCategories: make(map[string]int),
	}

	// 第一步: 诊断问题
	fmt.Println("🔍 第一步: 诊断书籍分类问题")
	fmt.Println(strings.Repeat("-", 40))

	// 1.1 检查所有书籍的分类
	fmt.Println("\n1.1 正在检查所有书籍的分类...")
	bookCursor, err := booksCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("❌ 查询书籍失败: %v\n", err)
		os.Exit(1)
	}
	defer bookCursor.Close(ctx)

	if err = bookCursor.All(ctx, &allBooks); err != nil {
		fmt.Printf("❌ 解析书籍数据失败: %v\n", err)
		os.Exit(1)
	}

	report.TotalBooks = int64(len(allBooks))
	fmt.Printf("   共找到 %d 本书籍\n", report.TotalBooks)

	// 收集书籍使用的分类
	for _, book := range allBooks {
		if categories, ok := book["categories"].(bson.A); ok && len(categories) > 0 {
			for _, cat := range categories {
				if catStr, ok := cat.(string); ok {
					report.BookCategories[catStr]++
				}
			}
		}
	}

	fmt.Printf("\n   书籍使用的分类:\n")
	for cat, count := range report.BookCategories {
		fmt.Printf("     - %s: %d 本书\n", cat, count)
	}

	// 1.2 检查分类表
	fmt.Println("\n1.2 正在检查分类表...")
	categoryCursor, err := categoriesCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("❌ 查询分类失败: %v\n", err)
		os.Exit(1)
	}
	defer categoryCursor.Close(ctx)

	var allCategories []bson.M
	if err = categoryCursor.All(ctx, &allCategories); err != nil {
		fmt.Printf("❌ 解析分类数据失败: %v\n", err)
		os.Exit(1)
	}

	report.TotalCategories = int64(len(allCategories))
	report.CategoryNames = make([]string, 0, len(allCategories))
	for _, cat := range allCategories {
		if name, ok := cat["name"].(string); ok {
			report.CategoryNames = append(report.CategoryNames, name)
		}
	}

	fmt.Printf("   共找到 %d 个分类:\n", report.TotalCategories)
	for _, name := range report.CategoryNames {
		fmt.Printf("     - %s\n", name)
	}

	// 1.3 对比分析
	fmt.Println("\n1.3 正在对比分析...")
	categorySet := make(map[string]bool)
	for _, name := range report.CategoryNames {
		categorySet[name] = true
	}

	report.UnmatchedBooks = make([]string, 0)
	report.InvalidBooks = 0

	for _, book := range allBooks {
		if categories, ok := book["categories"].(bson.A); ok && len(categories) > 0 {
			hasInvalid := false
			for _, cat := range categories {
				if catStr, ok := cat.(string); ok {
					if !categorySet[catStr] {
						hasInvalid = true
						title := "未知"
						if t, ok := book["title"].(string); ok {
							title = t
						}
						report.UnmatchedBooks = append(report.UnmatchedBooks,
							fmt.Sprintf("%s (分类: %s)", title, catStr))
						break
					}
				}
			}
			if hasInvalid {
				report.InvalidBooks++
			}
		}
	}

	fmt.Printf("   无效分类书籍数量: %d\n", report.InvalidBooks)

	if report.InvalidBooks > 0 {
		fmt.Println("\n   ⚠️  发现问题:")
		report.ProblemsFound = append(report.ProblemsFound,
			fmt.Sprintf("有 %d 本书籍使用了不存在的分类", report.InvalidBooks))

		// 显示前10个问题书籍
		maxShow := 10
		if len(report.UnmatchedBooks) < maxShow {
			maxShow = len(report.UnmatchedBooks)
		}
		for i := 0; i < maxShow; i++ {
			fmt.Printf("     - %s\n", report.UnmatchedBooks[i])
		}
		if len(report.UnmatchedBooks) > maxShow {
			fmt.Printf("     ... 还有 %d 本\n", len(report.UnmatchedBooks)-maxShow)
		}

		// 分析原因
		fmt.Println("\n   分析问题原因:")
		usedCategories := make([]string, 0, len(report.BookCategories))
		for cat := range report.BookCategories {
			usedCategories = append(usedCategories, cat)
		}

		// 检查是否是大小写问题
		for _, used := range usedCategories {
			for _, valid := range report.CategoryNames {
				if strings.EqualFold(used, valid) {
					report.ProblemsFound = append(report.ProblemsFound,
						fmt.Sprintf("大小写不匹配: 书籍使用 '%s', 分类表是 '%s'", used, valid))
					fmt.Printf("     ⚠️  大小写不匹配: '%s' vs '%s'\n", used, valid)
				}
			}
		}

		// 检查是否是空格问题
		for _, used := range usedCategories {
			trimmed := strings.TrimSpace(used)
			if trimmed != used && categorySet[trimmed] {
				report.ProblemsFound = append(report.ProblemsFound,
					fmt.Sprintf("包含空格: 书籍使用 '%s', 分类表是 '%s'", used, trimmed))
				fmt.Printf("     ⚠️  包含空格: '%s' vs '%s'\n", used, trimmed)
			}
		}

		// 检查是否有完全不匹配的分类
		for _, used := range usedCategories {
			found := false
			for _, valid := range report.CategoryNames {
				if strings.EqualFold(used, valid) ||
					strings.EqualFold(strings.TrimSpace(used), valid) {
					found = true
					break
				}
			}
			if !found {
				report.ProblemsFound = append(report.ProblemsFound,
					fmt.Sprintf("完全未知分类: '%s' 不在分类表中", used))
				fmt.Printf("     ⚠️  未知分类: '%s'\n", used)
			}
		}
	}

	// 第二步: 确定修复方案
	fmt.Println("\n========================================")
	fmt.Println("🛠️  第二步: 确定修复方案")
	fmt.Println(strings.Repeat("-", 40))

	if report.InvalidBooks == 0 {
		fmt.Println("✅ 没有发现分类问题，无需修复!")
		report.VerificationPassed = true
	} else {
		fmt.Println("\n可选修复方案:")
		fmt.Println("1. 方案A: 更新书籍分类 - 将书籍的分类更新为分类表中存在的值")
		fmt.Println("2. 方案B: 更新分类表 - 将分类表更新为书籍实际使用的值")
		fmt.Println("3. 方案C: 重建分类数据 - 根据书籍实际使用的分类重建分类表")
		fmt.Println("4. 仅生成诊断报告 - 不执行任何修复操作")

		fmt.Print("\n请选择方案 (1-4): ")
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			executeFixSchemeA(ctx, booksCollection, categoriesCollection, report)
		case 2:
			executeFixSchemeB(ctx, booksCollection, categoriesCollection, report)
		case 3:
			executeFixSchemeC(ctx, booksCollection, categoriesCollection, report)
		case 4:
			fmt.Println("仅生成报告，不执行修复")
		default:
			fmt.Println("无效选择")
		}

		// 第三步: 验证结果
		if len(report.FixesApplied) > 0 {
			fmt.Println("\n========================================")
			fmt.Println("✅ 第三步: 验证修复结果")
			fmt.Println(strings.Repeat("-", 40))

			verifyFixResults(ctx, booksCollection, categoriesCollection, report)
		}
	}

	// 生成报告
	fmt.Println("\n========================================")
	fmt.Println("📄 生成修复报告")
	fmt.Println(strings.Repeat("-", 40))

	generateReport(report)
}

// executeFixSchemeA 执行方案A: 更新书籍分类
func executeFixSchemeA(ctx context.Context, booksColl, catColl *mongo.Collection, report *CategoryFixReport) {
	fmt.Println("\n🔄 执行方案A: 更新书籍分类")

	// 构建分类映射
	categorySet := make(map[string]bool)
	for _, name := range report.CategoryNames {
		categorySet[name] = true
	}

	// 构建模糊匹配映射
	normalizedMap := make(map[string]string) // 标准化后的名称 -> 正确的分类名
	for _, validName := range report.CategoryNames {
		normalizedMap[strings.ToLower(strings.TrimSpace(validName))] = validName
	}

	// 统计和更新
	fixCount := 0
	for _, book := range allBooks {
		if categories, ok := book["categories"].(bson.A); ok && len(categories) > 0 {
			needsUpdate := false
			newCategories := make(bson.A, 0, len(categories))

			for _, cat := range categories {
				if catStr, ok := cat.(string); ok {
					// 检查分类是否有效
					if categorySet[catStr] {
						newCategories = append(newCategories, catStr)
					} else {
						// 尝试模糊匹配
						normalized := strings.ToLower(strings.TrimSpace(catStr))
						if validName, exists := normalizedMap[normalized]; exists {
							newCategories = append(newCategories, validName)
							needsUpdate = true
							report.FixesApplied = append(report.FixesApplied,
								fmt.Sprintf("书籍 '%s' 的分类 '%s' -> '%s'",
									getBookTitle(book), catStr, validName))
						} else {
							// 无法匹配，保留原值
							newCategories = append(newCategories, catStr)
						}
					}
				}
			}

			if needsUpdate {
				id := book["_id"]
				_, err := booksColl.UpdateOne(ctx,
					bson.M{"_id": id},
					bson.M{"$set": bson.M{"categories": newCategories}})
				if err != nil {
					fmt.Printf("   ❌ 更新书籍 %s 失败: %v\n", getBookTitle(book), err)
				} else {
					fixCount++
				}
			}
		}
	}

	report.FixedCount = int64(fixCount)
	fmt.Printf("   ✅ 已更新 %d 本书籍的分类\n", fixCount)
}

// executeFixSchemeB 执行方案B: 更新分类表
func executeFixSchemeB(ctx context.Context, booksColl, catColl *mongo.Collection, report *CategoryFixReport) {
	fmt.Println("\n🔄 执行方案B: 更新分类表")

	// 收集书籍使用但分类表中没有的分类
	missingCategories := make(map[string]bool)
	for cat := range report.BookCategories {
		found := false
		for _, name := range report.CategoryNames {
			if name == cat {
				found = true
				break
			}
		}
		if !found {
			missingCategories[cat] = true
		}
	}

	if len(missingCategories) == 0 {
		fmt.Println("   ℹ️  所有书籍使用的分类都已在分类表中")
		return
	}

	// 添加缺失的分类
	now := time.Now()
	addedCount := 0
	for cat := range missingCategories {
		newCat := bson.M{
			"_id":         primitive.NewObjectID(),
			"name":        cat,
			"slug":        strings.ToLower(strings.ReplaceAll(cat, " ", "-")),
			"description": cat + "分类",
			"icon":        "/images/icons/default.png",
			"parent_id":   nil,
			"sort_order":  len(report.CategoryNames) + addedCount + 1,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		}

		_, err := catColl.InsertOne(ctx, newCat)
		if err != nil {
			fmt.Printf("   ❌ 添加分类 '%s' 失败: %v\n", cat, err)
		} else {
			report.FixesApplied = append(report.FixesApplied,
				fmt.Sprintf("添加新分类: '%s'", cat))
			report.CategoryNames = append(report.CategoryNames, cat)
			addedCount++
		}
	}

	fmt.Printf("   ✅ 已添加 %d 个新分类\n", addedCount)
}

// executeFixSchemeC 执行方案C: 重建分类数据
func executeFixSchemeC(ctx context.Context, booksColl, catColl *mongo.Collection, report *CategoryFixReport) {
	fmt.Println("\n🔄 执行方案C: 重建分类数据")

	// 1. 清空现有分类
	_, err := catColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		fmt.Printf("   ❌ 清空分类表失败: %v\n", err)
		return
	}
	fmt.Println("   ✓ 已清空现有分类表")

	// 2. 根据书籍使用的分类重建
	now := time.Now()
	sortOrder := 1
	addedCount := 0

	for cat := range report.BookCategories {
		newCat := bson.M{
			"_id":         primitive.NewObjectID(),
			"name":        cat,
			"slug":        strings.ToLower(strings.ReplaceAll(cat, " ", "-")),
			"description": cat + "分类",
			"icon":        "/images/icons/default.png",
			"parent_id":   nil,
			"sort_order":  sortOrder,
			"is_active":   true,
			"created_at":  now,
			"updated_at":  now,
		}

		_, err := catColl.InsertOne(ctx, newCat)
		if err != nil {
			fmt.Printf("   ❌ 添加分类 '%s' 失败: %v\n", cat, err)
		} else {
			report.FixesApplied = append(report.FixesApplied,
				fmt.Sprintf("重建分类: '%s'", cat))
			addedCount++
			sortOrder++
		}
	}

	// 更新报告中的分类名称
	report.CategoryNames = make([]string, 0, len(report.BookCategories))
	for cat := range report.BookCategories {
		report.CategoryNames = append(report.CategoryNames, cat)
	}

	fmt.Printf("   ✅ 已重建 %d 个分类\n", addedCount)
}

// verifyFixResults 验证修复结果
func verifyFixResults(ctx context.Context, booksColl, catColl *mongo.Collection, report *CategoryFixReport) {
	// 重新检查无效分类书籍数量
	categorySet := make(map[string]bool)
	for _, name := range report.CategoryNames {
		categorySet[name] = true
	}

	invalidCount := int64(0)
	bookCursor, err := booksColl.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("   ❌ 验证查询失败: %v\n", err)
		return
	}
	defer bookCursor.Close(ctx)

	var books []bson.M
	bookCursor.All(ctx, &books)

	for _, book := range books {
		if categories, ok := book["categories"].(bson.A); ok && len(categories) > 0 {
			hasInvalid := false
			for _, cat := range categories {
				if catStr, ok := cat.(string); ok {
					if !categorySet[catStr] {
						hasInvalid = true
						break
					}
				}
			}
			if hasInvalid {
				invalidCount++
			}
		}
	}

	report.VerificationPassed = (invalidCount == 0)

	fmt.Printf("   修复前无效分类书籍: %d\n", report.InvalidBooks)
	fmt.Printf("   修复后无效分类书籍: %d\n", invalidCount)

	if report.VerificationPassed {
		fmt.Println("   ✅ 验证通过! 所有书籍的分类都有效")
	} else {
		fmt.Println("   ⚠️  验证未完全通过，仍有书籍使用无效分类")
	}
}

// generateReport 生成报告
func generateReport(report *CategoryFixReport) {
	reportContent := fmt.Sprintf(`# 书籍分类关联修复报告

**生成时间**: %s
**修复状态**: %s

## 一、诊断结果

### 1.1 基本信息
- 总书籍数量: %d
- 无效分类书籍数量: %d
- 分类表分类数量: %d

### 1.2 书籍使用的分类
%s

### 1.3 分类表中的分类
%s

### 1.4 发现的问题
%s

## 二、修复方案

%s

## 三、执行的操作

%s

## 四、验证结果

- 修复书籍数量: %d
- 验证状态: %s

%s

---
*此报告由书籍分类关联修复工具自动生成*
`,
		time.Now().Format("2006-01-02 15:04:05"),
		getStatusText(report),
		report.TotalBooks,
		report.InvalidBooks,
		report.TotalCategories,
		formatMapList(report.BookCategories),
		formatStringList(report.CategoryNames),
		formatProblemList(report.ProblemsFound),
		getSchemeDescription(report),
		formatFixList(report.FixesApplied),
		report.FixedCount,
		getVerificationText(report),
		getConclusion(report),
	)

	// 保存报告
	reportPath := "docs/reports/2026-02-01-category-fix-report.md"
	_ = os.MkdirAll("docs/reports", 0755)

	err := os.WriteFile(reportPath, []byte(reportContent), 0644)
	if err != nil {
		fmt.Printf("⚠️  保存报告失败: %v\n", err)
		fmt.Println("\n报告内容:")
		fmt.Println(reportContent)
	} else {
		fmt.Printf("✅ 报告已保存到: %s\n", reportPath)
	}
}

// 辅助函数
func getBookTitle(book bson.M) string {
	if title, ok := book["title"].(string); ok {
		return title
	}
	return "未知"
}

func formatMapList(m map[string]int) string {
	if len(m) == 0 {
		return "(无)"
	}
	lines := make([]string, 0, len(m))
	for k, v := range m {
		lines = append(lines, fmt.Sprintf("- %s: %d 本书", k, v))
	}
	return "\n" + strings.Join(lines, "\n")
}

func formatStringList(list []string) string {
	if len(list) == 0 {
		return "(无)"
	}
	lines := make([]string, 0, len(list))
	for _, v := range list {
		lines = append(lines, fmt.Sprintf("- %s", v))
	}
	return "\n" + strings.Join(lines, "\n")
}

func formatProblemList(problems []string) string {
	if len(problems) == 0 {
		return "未发现问题"
	}
	lines := make([]string, 0, len(problems))
	for _, p := range problems {
		lines = append(lines, fmt.Sprintf("- %s", p))
	}
	return "\n" + strings.Join(lines, "\n")
}

func formatFixList(fixes []string) string {
	if len(fixes) == 0 {
		return "未执行任何修复操作"
	}
	lines := make([]string, 0, len(fixes))
	for _, f := range fixes {
		lines = append(lines, fmt.Sprintf("- %s", f))
	}
	return "\n" + strings.Join(lines, "\n")
}

func getStatusText(report *CategoryFixReport) string {
	if report.InvalidBooks == 0 {
		return "✅ 无需修复"
	}
	if report.VerificationPassed {
		return "✅ 修复成功"
	}
	return "⚠️ 修复未完全成功"
}

func getSchemeDescription(report *CategoryFixReport) string {
	if len(report.FixesApplied) == 0 {
		return "未执行修复方案"
	}
	return "根据选择的方案执行了相应的修复操作"
}

func getVerificationText(report *CategoryFixReport) string {
	if report.VerificationPassed {
		return "✅ 通过"
	}
	return "⚠️ 未完全通过"
}

func getConclusion(report *CategoryFixReport) string {
	if report.InvalidBooks == 0 {
		return "## 结论\n\n所有书籍的分类关联正常，无需修复。"
	}
	if report.VerificationPassed {
		return fmt.Sprintf("## 结论\n\n成功修复了 %d 本书籍的分类关联问题。", report.FixedCount)
	}
	return fmt.Sprintf("## 结论\n\n修复了 %d 本书籍，但仍有问题存在。建议进一步检查。", report.FixedCount)
}
