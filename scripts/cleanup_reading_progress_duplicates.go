package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DuplicateGroup 重复记录组
type DuplicateGroup struct {
	UserID     string `bson:"_id_user_id"`  // 使用 string 支持 UUID 和 ObjectID
	BookID     string `bson:"_id_book_id"`  // 使用 string 支持 UUID 和 ObjectID
	Count      int    `bson:"count"`
	Docs       []Doc  `bson:"docs"`
}

// Doc 文档信息
type Doc struct {
	ID        string    `bson:"id"`  // 使用 string 支持 UUID 和 ObjectID
	UpdatedAt time.Time `bson:"updated_at"`
	CreatedAt time.Time `bson:"created_at"`
}

// CleanupReport 清理报告
type CleanupReport struct {
	ExecutionTime    time.Time `json:"execution_time"`
	TotalGroups      int       `json:"total_groups"`
	TotalDuplicates  int       `json:"total_duplicates"`
	DeletedCount     int       `json:"deleted_count"`
	RemainingCount   int       `json:"remaining_count"`
	DuplicateGroups  []GroupDetail `json:"duplicate_groups"`
}

// GroupDetail 重复组详情
type GroupDetail struct {
	UserID        string `json:"user_id"`
	BookID        string `json:"book_id"`
	TotalCount    int    `json:"total_count"`
	DeletedCount  int    `json:"deleted_count"`
	RemainingID   string `json:"remaining_id"`
	DeletedIDs    []string `json:"deleted_ids"`
}

const (
	collectionName = "reading_progress"
	databaseName   = "qingyu"
	mongoURI       = "mongodb://localhost:27017"
)

func main() {
	log.Println("=== Reading Progress 重复记录清理工具 ===")
	log.Printf("执行时间: %s\n", time.Now().Format(time.RFC3339))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// 连接数据库
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(databaseName)
	coll := db.Collection(collectionName)

	log.Printf("数据库: %s, 集合: %s\n", databaseName, collectionName)

	// 第一步：查找重复记录
	log.Println("\n第1步: 查找重复记录...")
	duplicateGroups, err := findDuplicates(ctx, coll)
	if err != nil {
		log.Fatalf("查找重复记录失败: %v", err)
	}

	if len(duplicateGroups) == 0 {
		log.Println("未发现重复记录，数据库状态正常！")
		return
	}

	log.Printf("发现 %d 组重复记录\n", len(duplicateGroups))

	// 第二步：分析并确认要删除的记录
	log.Println("\n第2步: 分析重复记录...")
	idsToDelete := make(map[string]bool)
	report := CleanupReport{
		ExecutionTime: time.Now(),
		DuplicateGroups: []GroupDetail{},
	}

	totalDuplicates := 0
	for _, group := range duplicateGroups {
		totalDuplicates += group.Count

		// 按 updated_at 降序排序，保留最新的
		sortedDocs := sortDocsByUpdatedAt(group.Docs)

		detail := GroupDetail{
			UserID:       group.UserID,
			BookID:       group.BookID,
			TotalCount:   group.Count,
			RemainingID:  sortedDocs[0].ID,
			DeletedIDs:   []string{},
		}

		// 第一个保留，其余标记为删除
		for i := 1; i < len(sortedDocs); i++ {
			idsToDelete[sortedDocs[i].ID] = true
			detail.DeletedIDs = append(detail.DeletedIDs, sortedDocs[i].ID)
		}
		detail.DeletedCount = len(detail.DeletedIDs)

		report.DuplicateGroups = append(report.DuplicateGroups, detail)
	}

	report.TotalGroups = len(duplicateGroups)
	report.TotalDuplicates = totalDuplicates

	log.Printf("总共 %d 条重复记录需要删除\n", len(idsToDelete))

	// 第三步：执行删除
	log.Println("\n第3步: 删除重复记录...")
	deleteIDs := make([]string, 0, len(idsToDelete))
	for id := range idsToDelete {
		deleteIDs = append(deleteIDs, id)
	}

	// 批量删除（每次1000条）
	batchSize := 1000
	deletedTotal := 0
	if len(deleteIDs) == 0 {
		log.Println("警告: 没有需要删除的记录")
	} else {
		log.Printf("准备删除 %d 条记录\n", len(deleteIDs))
		// 调试: 显示前几个要删除的 ID
		for i := 0; i < len(deleteIDs) && i < 5; i++ {
			log.Printf("  要删除的ID[%d]: %s\n", i, deleteIDs[i])

			// 验证 ID 是否存在
			filter := bson.M{"_id": deleteIDs[i]}
			count, _ := coll.CountDocuments(ctx, filter)
			log.Printf("    验证: _id=%s 的记录数 = %d\n", deleteIDs[i], count)
		}

		for i := 0; i < len(deleteIDs); i += batchSize {
			end := i + batchSize
			if end > len(deleteIDs) {
				end = len(deleteIDs)
			}

			batch := deleteIDs[i:end]
			log.Printf("删除批次 [%d:%d], 共 %d 条记录\n", i, end, len(batch))
			result, err := coll.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": batch}})
			if err != nil {
				log.Printf("批量删除失败 [%d:%d]: %v\n", i, end, err)
				continue
			}
			deletedTotal += int(result.DeletedCount)
			log.Printf("已删除 [%d:%d]: %d 条\n", i, end, result.DeletedCount)
		}
	}

	report.DeletedCount = deletedTotal

	// 第四步：验证结果
	log.Println("\n第4步: 验证结果...")
	remainingDuplicates, err := findDuplicates(ctx, coll)
	if err != nil {
		log.Printf("验证失败: %v\n", err)
	} else {
		if len(remainingDuplicates) > 0 {
			log.Printf("警告: 仍有 %d 组重复记录\n", len(remainingDuplicates))
		} else {
			log.Println("验证成功: 无重复记录")
		}
	}

	// 获取当前集合总文档数
	totalDocs, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("获取文档总数失败: %v\n", err)
	} else {
		report.RemainingCount = int(totalDocs)
		log.Printf("当前集合文档总数: %d\n", totalDocs)
	}

	// 第五步：生成报告
	log.Println("\n第5步: 生成报告...")
	generateReport(report)

	log.Println("\n=== 清理完成 ===")
	log.Printf("总重复组数: %d\n", report.TotalGroups)
	log.Printf("总重复记录数: %d\n", report.TotalDuplicates)
	log.Printf("已删除记录数: %d\n", report.DeletedCount)
	log.Printf("剩余记录数: %d\n", report.RemainingCount)
}

// findDuplicates 查找重复记录
func findDuplicates(ctx context.Context, coll *mongo.Collection) ([]DuplicateGroup, error) {
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", bson.D{
				{"user_id", "$user_id"},
				{"book_id", "$book_id"},
			}},
			{"docs", bson.D{
				{"$push", bson.D{
					{"id", "$_id"},
					{"updated_at", "$updated_at"},
					{"created_at", "$created_at"},
				}},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}}},
		{{"$match", bson.D{
			{"count", bson.D{{"$gt", 1}}},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	groups := make([]DuplicateGroup, 0, len(results))
	for _, result := range results {
		id := result["_id"].(bson.M)
		var userID, bookID string

		// 处理 user_id，可能是 string 或 ObjectID
		if uid, ok := id["user_id"].(string); ok {
			userID = uid
		} else if uid, ok := id["user_id"].(primitive.ObjectID); ok {
			userID = uid.Hex()
		}

		// 处理 book_id，可能是 string 或 ObjectID
		if bid, ok := id["book_id"].(string); ok {
			bookID = bid
		} else if bid, ok := id["book_id"].(primitive.ObjectID); ok {
			bookID = bid.Hex()
		}

		docsRaw := result["docs"].(bson.A)
		docs := make([]Doc, 0, len(docsRaw))
		for _, docRaw := range docsRaw {
			doc := docRaw.(bson.M)
			// 处理 ID，可能是 string 或 ObjectID
			var docID string
			if idStr, ok := doc["id"].(string); ok {
				docID = idStr
			} else if oid, ok := doc["id"].(primitive.ObjectID); ok {
				docID = oid.Hex()
			}

			// 处理时间字段，可能是 time.Time 或 primitive.DateTime
			var updatedAt, createdAt time.Time
			if ut, ok := doc["updated_at"].(time.Time); ok {
				updatedAt = ut
			} else if pd, ok := doc["updated_at"].(primitive.DateTime); ok {
				updatedAt = pd.Time()
			}
			if ct, ok := doc["created_at"].(time.Time); ok {
				createdAt = ct
			} else if pd, ok := doc["created_at"].(primitive.DateTime); ok {
				createdAt = pd.Time()
			}

			docs = append(docs, Doc{
				ID:        docID,
				UpdatedAt: updatedAt,
				CreatedAt: createdAt,
			})
		}

		groups = append(groups, DuplicateGroup{
			UserID: userID,
			BookID: bookID,
			Count:  int(result["count"].(int32)),
			Docs:   docs,
		})
	}

	return groups, nil
}

// sortDocsByUpdatedAt 按 updated_at 降序排序
func sortDocsByUpdatedAt(docs []Doc) []Doc {
	sorted := make([]Doc, len(docs))
	copy(sorted, docs)

	// 使用简单的冒泡排序（文档数量通常很少）
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].UpdatedAt.After(sorted[i].UpdatedAt) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// generateReport 生成报告
func generateReport(report CleanupReport) {
	// 确保报告目录存在
	reportDir := "../../docs/reports"
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		log.Printf("创建报告目录失败: %v\n", err)
		return
	}

	// 生成 JSON 文件
	jsonPath := fmt.Sprintf("%s/2026-02-01-duplicate-cleanup-report.json", reportDir)
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Printf("生成 JSON 报告失败: %v\n", err)
	} else {
		if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
			log.Printf("写入 JSON 报告失败: %v\n", err)
		} else {
			log.Printf("JSON 报告已保存: %s\n", jsonPath)
		}
	}

	// 生成 Markdown 文件
	mdPath := fmt.Sprintf("%s/2026-02-01-duplicate-cleanup-report.md", reportDir)
	mdContent := formatMarkdownReport(report)
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		log.Printf("写入 Markdown 报告失败: %v\n", err)
	} else {
		log.Printf("Markdown 报告已保存: %s\n", mdPath)
	}
}

// formatMarkdownReport 格式化 Markdown 报告
func formatMarkdownReport(report CleanupReport) string {
	md := fmt.Sprintf(`# Reading Progress 重复记录清理报告

**执行时间**: %s

## 概览

| 指标 | 数值 |
|------|------|
| 重复组总数 | %d |
| 重复记录总数 | %d |
| 已删除记录数 | %d |
| 剩余记录数 | %d |

## 重复组详情

`,
		report.ExecutionTime.Format("2006-01-02 15:04:05"),
		report.TotalGroups,
		report.TotalDuplicates,
		report.DeletedCount,
		report.RemainingCount,
	)

	for i, group := range report.DuplicateGroups {
		md += fmt.Sprintf("### %d. User: %s, Book: %s\n\n", i+1, group.UserID, group.BookID)
		md += fmt.Sprintf("- 总记录数: %d\n", group.TotalCount)
		md += fmt.Sprintf("- 删除记录数: %d\n", group.DeletedCount)
		md += fmt.Sprintf("- 保留记录ID: `%s`\n", group.RemainingID)
		if len(group.DeletedIDs) > 0 {
			md += "- 已删除ID:\n"
			for _, id := range group.DeletedIDs {
				md += fmt.Sprintf("  - `%s`\n", id)
			}
		}
		md += "\n"
	}

	md += "## 验证结果\n\n"
	md += "请执行以下命令验证无重复记录：\n\n"
	md += "```javascript\n"
	md += "db.reading_progress.aggregate([\n"
	md += "  { $group: {\n"
	md += "    _id: { user_id: \"$user_id\", book_id: \"$book_id\" },\n"
	md += "    count: { $sum: 1 }\n"
	md += "  }},\n"
	md += "  { $match: { count: { $gt: 1 } } },\n"
	md += "  { $count: \"duplicates\" }\n"
	md += "])\n"
	md += "```\n\n"
	md += "## 建议\n\n"
	md += "1. 考虑在 `user_id` 和 `book_id` 上创建唯一索引，防止未来产生重复记录\n"
	md += "2. 检查应用代码中是否存在并发更新导致重复的逻辑\n\n"
	md += "---\n\n"
	md += "*此报告由 cleanup_reading_progress_duplicates.go 自动生成*\n"

	return md
}
