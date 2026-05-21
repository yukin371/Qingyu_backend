package main

import (
	"context"
	"fmt"
	"time"

	"Qingyu_backend/cmd/seeder/story"
	"Qingyu_backend/cmd/seeder/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HarnessDataSeeder Harness测试数据填充器
type HarnessDataSeeder struct {
	db *utils.Database
}

// NewHarnessDataSeeder 创建Harness数据填充器
func NewHarnessDataSeeder(db *utils.Database) *HarnessDataSeeder {
	return &HarnessDataSeeder{db: db}
}

// CreateHarnessData 创建Harness测试数据
func (s *HarnessDataSeeder) CreateHarnessData() error {

	// 1. 检查项目是否已存在
	if s.projectExists(story.EmploymentGuideProject.Title) {
		fmt.Println("✓ Harness项目「就业指南：我在异世界当劝退专员」已存在，跳过创建")
		return nil
	}

	fmt.Println("🚀 开始创建Harness数据...")

	// 2. 获取/创建 Harness 专用用户
	harnessUserID, err := s.getOrCreateHarnessUser()
	if err != nil {
		return fmt.Errorf("获取 Harness 账号失败: %w", err)
	}
	fmt.Printf("✓ Harness 账号: %s\n", harnessUserID)

	// 3. 创建项目
	projectID, err := s.createProject(harnessUserID)
	if err != nil {
		return fmt.Errorf("创建项目失败: %w", err)
	}
	fmt.Printf("✓ 创建项目: %s\n", projectID)

	// 4. 创建文档树（卷 -> 章）
	docIDMap, err := s.createDocuments(projectID)
	if err != nil {
		return fmt.Errorf("创建Harness文档失败: %w", err)
	}
	fmt.Printf("✓ 创建文档: %d个\n", len(docIDMap))

	// 5. 创建章节内容
	if err := s.createChapterContents(docIDMap); err != nil {
		return fmt.Errorf("创建Harness正文失败: %w", err)
	}
	fmt.Println("✓ 创建章节正文: 1章")

	// 6. 创建角色
	characterIDs, err := s.createCharacters(projectID)
	if err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}
	fmt.Printf("✓ 创建角色: %d个\n", len(characterIDs))

	// 7. 创建最小时间线
	if err := s.createTimelines(projectID, docIDMap, characterIDs); err != nil {
		return fmt.Errorf("创建Harness时间线失败: %w", err)
	}
	fmt.Println("✓ 创建时间线: 1条")

	// 8. 更新项目统计
	if err := s.updateProjectStatistics(projectID, len(docIDMap)); err != nil {
		return fmt.Errorf("更新项目统计失败: %w", err)
	}

	// 9. 验证数据完整性
	if err := s.validateHarnessData(projectID, docIDMap); err != nil {
		return fmt.Errorf("数据验证失败: %w", err)
	}

	fmt.Println("\n✅ Harness数据创建完成！")
	fmt.Println("📖 项目名称: 就业指南：我在异世界当劝退专员")
	fmt.Println("🔐 登录账号: harness_writer / qingyu2024")

	return nil
}

// Clean 清空Harness数据
func (s *HarnessDataSeeder) Clean() error {
	ctx := context.Background()

	// 获取项目ID
	var project struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := s.db.Collection("projects").FindOne(ctx, bson.M{"title": story.EmploymentGuideProject.Title}).Decode(&project)
	if err != nil {
		fmt.Println("✓ Harness项目不存在，无需清理")
		return nil
	}

	projectID := project.ID
	documentCursor, err := s.db.Collection("documents").Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return fmt.Errorf("查询Harness文档失败: %w", err)
	}
	defer documentCursor.Close(ctx)

	documentIDs := make([]primitive.ObjectID, 0)
	for documentCursor.Next(ctx) {
		var item struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := documentCursor.Decode(&item); err != nil {
			return fmt.Errorf("解析Harness文档失败: %w", err)
		}
		documentIDs = append(documentIDs, item.ID)
	}
	if err := documentCursor.Err(); err != nil {
		return fmt.Errorf("遍历Harness文档失败: %w", err)
	}

	// 按依赖顺序删除
	collections := []string{
		"story_harness_projections",
		"change_requests",
		"change_request_batches",
		"timeline_events",
		"timelines",
		"document_contents",
		"documents",
		"character_relations",
		"characters",
		"projects",
	}

	for _, colName := range collections {
		filter := bson.M{"project_id": projectID}
		if colName == "document_contents" {
			if len(documentIDs) == 0 {
				continue
			}
			filter = bson.M{"document_id": bson.M{"$in": documentIDs}}
		}
		if colName == "projects" {
			filter = bson.M{"_id": projectID}
		}
		result, err := s.db.Collection(colName).DeleteMany(ctx, filter)
		if err != nil {
			return fmt.Errorf("清空 %s 失败: %w", colName, err)
		}
		if result.DeletedCount > 0 {
			fmt.Printf("✓ 清空 %s: %d条\n", colName, result.DeletedCount)
		}
	}

	fmt.Println("✓ Harness数据清理完成")
	return nil
}

// projectExists 检查项目是否存在
func (s *HarnessDataSeeder) projectExists(title string) bool {
	ctx := context.Background()
	count, _ := s.db.Collection("projects").CountDocuments(ctx, bson.M{"title": title})
	return count > 0
}

// getOrCreateHarnessUser 获取或创建Harness专用用户
// 独立于admin，避免被 seeder users --clean 误删
func (s *HarnessDataSeeder) getOrCreateHarnessUser() (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("users")

	// 查找专用的 harness 用户
	harnessUsername := "harness_writer"
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := collection.FindOne(ctx, bson.M{"username": harnessUsername}).Decode(&user)
	if err == nil {
		fmt.Printf("✓ 使用已存在Harness用户: %s\n", harnessUsername)
		return user.ID.Hex(), nil
	}

	// 不存在则创建专用用户
	return s.createHarnessUser()
}

// createHarnessUser 创建Harness专用用户
func (s *HarnessDataSeeder) createHarnessUser() (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("users")

	now := time.Now()
	id := primitive.NewObjectID()

	// 使用 bcrypt 哈希密码
	hashedPassword, err := utils.HashPassword("qingyu2024")
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}

	user := bson.M{
		"_id":        id,
		"username":   "harness_writer",
		"email":      "harness@qingyu.com",
		"password":   hashedPassword,
		"roles":      []string{"author"},
		"status":     "active",
		"nickname":   "Harness作家",
		"avatar":     "/images/avatars/default.png",
		"bio":        "Harness测试专用账号",
		"created_at": now,
		"updated_at": now,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	fmt.Println("✓ 创建 Harness 账号: harness_writer")
	return id.Hex(), nil
}

// createProject 创建项目
func (s *HarnessDataSeeder) createProject(authorID string) (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("projects")

	now := time.Now()
	id := primitive.NewObjectID()

	project := bson.M{
		"_id":          id,
		"author_id":    func() primitive.ObjectID { oid, _ := primitive.ObjectIDFromHex(authorID); return oid }(),
		"title":        story.EmploymentGuideProject.Title,
		"summary":      story.EmploymentGuideProject.Summary,
		"cover_url":    story.EmploymentGuideProject.CoverURL,
		"category":     story.EmploymentGuideProject.Category,
		"writing_type": story.EmploymentGuideProject.WritingType,
		"status":       story.EmploymentGuideProject.Status,
		"visibility":   story.EmploymentGuideProject.Visibility,
		"tags":         story.EmploymentGuideProject.Tags,
		"statistics": bson.M{
			"total_words":    0,
			"chapter_count":  0,
			"document_count": 0,
			"last_update_at": now,
		},
		"settings": bson.M{
			"auto_backup":     true,
			"backup_interval": 24,
		},
		"collaborators": []interface{}{},
		"created_at":    now,
		"updated_at":    now,
	}

	_, err := collection.InsertOne(ctx, project)
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

// createCharacters 创建角色
func (s *HarnessDataSeeder) createCharacters(projectID string) (map[string]string, error) {
	ctx := context.Background()
	collection := s.db.Collection("characters")

	now := time.Now()
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	characterIDs := make(map[string]string)

	for _, char := range story.NewCharacters {
		charID := primitive.NewObjectID()
		character := bson.M{
			"_id":                charID,
			"project_id":         projectOID,
			"name":               char.Name,
			"alias":              char.Alias,
			"summary":            char.Summary,
			"traits":             char.Traits,
			"background":         char.Background,
			"avatar_url":         char.AvatarURL,
			"short_description":  char.Summary,
			"personality_prompt": char.PersonalityPrompt,
			"speech_pattern":     char.SpeechPattern,
			"current_state":      "",
			"entity_type":        "character",
			"created_at":         now,
			"updated_at":         now,
		}
		_, err := collection.InsertOne(ctx, character)
		if err != nil {
			return nil, err
		}
		characterIDs[char.ID] = charID.Hex()
	}

	return characterIDs, nil
}

func (s *HarnessDataSeeder) createDocuments(projectID string) (map[string]string, error) {
	ctx := context.Background()
	collection := s.db.Collection("documents")
	now := time.Now()
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	docIDs := make(map[string]string, len(story.EmploymentGuideDocuments))

	for _, item := range story.EmploymentGuideDocuments {
		docID := primitive.NewObjectID()
		parentID := primitive.NilObjectID
		if item.ParentKey != "" {
			parentHex, ok := docIDs[item.ParentKey]
			if !ok {
				return nil, fmt.Errorf("父文档未创建: %s -> %s", item.Key, item.ParentKey)
			}
			parentID, _ = primitive.ObjectIDFromHex(parentHex)
		}
		document := bson.M{
			"_id":          docID,
			"project_id":   projectOID,
			"title":        item.Title,
			"type":         item.Type,
			"level":        item.Level,
			"order":        item.Order,
			"parent_id":    parentID,
			"status":       item.Status,
			"word_count":   0,
			"stable_ref":   primitive.NewObjectID().Hex(),
			"order_key":    item.OrderKey,
			"plot_threads": item.PlotThreads,
			"key_points":   item.KeyPoints,
			"notes":        item.Summary,
			"scene_goal":   item.SceneGoal,
			"created_at":   now,
			"updated_at":   now,
		}
		_, err := collection.InsertOne(ctx, document)
		if err != nil {
			return nil, err
		}
		docIDs[item.Key] = docID.Hex()
	}

	return docIDs, nil
}

func (s *HarnessDataSeeder) createChapterContents(docIDMap map[string]string) error {
	ctx := context.Background()
	collection := s.db.Collection("document_contents")
	chapterIDHex, ok := docIDMap[story.EmploymentGuideChapterOneKey]
	if !ok {
		return fmt.Errorf("缺少章节文档ID: %s", story.EmploymentGuideChapterOneKey)
	}
	chapterID, err := primitive.ObjectIDFromHex(chapterIDHex)
	if err != nil {
		return err
	}

	content := story.EmploymentGuideChapterOneContent
	now := time.Now()
	documentContent := bson.M{
		"_id":             primitive.NewObjectID(),
		"document_id":     chapterID,
		"content":         content,
		"content_type":    "markdown",
		"word_count":      len([]rune(content)),
		"char_count":      len(content),
		"paragraph_order": 0,
		"version":         1,
		"last_saved_at":   now,
		"last_edited_by":  "harness_writer",
		"created_at":      now,
		"updated_at":      now,
	}

	_, err = collection.InsertOne(ctx, documentContent)
	return err
}

func (s *HarnessDataSeeder) createTimelines(projectID string, docIDMap map[string]string, characterIDs map[string]string) error {
	ctx := context.Background()
	timelineCol := s.db.Collection("timelines")
	eventCol := s.db.Collection("timeline_events")
	now := time.Now()
	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	for _, timeline := range story.EmploymentGuideTimelines {
		timelineID := primitive.NewObjectID()
		timelineDoc := bson.M{
			"_id":         timelineID,
			"project_id":  projectOID,
			"name":        timeline.Name,
			"description": timeline.Description,
			"start_time":  nil,
			"end_time":    nil,
			"created_at":  now,
			"updated_at":  now,
		}
		if _, err := timelineCol.InsertOne(ctx, timelineDoc); err != nil {
			return err
		}

		for _, event := range timeline.Events {
			participants := make([]string, 0, len(event.ParticipantKeys))
			for _, key := range event.ParticipantKeys {
				if id, ok := characterIDs[key]; ok {
					participants = append(participants, id)
				}
			}
			chapterIDs := make([]string, 0, 1)
			if event.ChapterKey != "" {
				if chapterID, ok := docIDMap[event.ChapterKey]; ok {
					chapterIDs = append(chapterIDs, chapterID)
				}
			}

			eventDoc := bson.M{
				"_id":          primitive.NewObjectID(),
				"project_id":   projectOID,
				"timeline_id":  timelineID.Hex(),
				"title":        event.Title,
				"description":  event.Description,
				"story_time":   bson.M{"description": event.TimeLabel},
				"participants": participants,
				"chapter_ids":  chapterIDs,
				"event_type":   event.EventType,
				"importance":   event.Importance,
				"created_at":   now,
				"updated_at":   now,
			}
			if _, err := eventCol.InsertOne(ctx, eventDoc); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *HarnessDataSeeder) updateProjectStatistics(projectID string, documentCount int) error {
	ctx := context.Background()
	projectOID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}
	_, err = s.db.Collection("projects").UpdateByID(ctx, projectOID, bson.M{
		"$set": bson.M{
			"statistics.chapter_count":  1,
			"statistics.document_count": documentCount,
			"statistics.last_update_at": time.Now(),
			"updated_at":                time.Now(),
		},
	})
	return err
}

// validateHarnessData 验证数据完整性
func (s *HarnessDataSeeder) validateHarnessData(projectID string, docIDMap map[string]string) error {
	ctx := context.Background()
	projectOID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}

	chapterIDHex, ok := docIDMap[story.EmploymentGuideChapterOneKey]
	if !ok {
		return fmt.Errorf("缺少章节文档ID: %s", story.EmploymentGuideChapterOneKey)
	}
	chapterOID, err := primitive.ObjectIDFromHex(chapterIDHex)
	if err != nil {
		return err
	}

	// 检查各集合数量
	collections := []struct {
		name     string
		expected int64
		filter   bson.M
	}{
		{"characters", 3, bson.M{"project_id": projectOID}},
		{"documents", 2, bson.M{"project_id": projectOID}},
		{"document_contents", 1, bson.M{"document_id": chapterOID}},
		{"timelines", 1, bson.M{"project_id": projectOID}},
		{"timeline_events", 2, bson.M{"project_id": projectOID}},
	}

	fmt.Println("\n📊 数据完整性验证:")
	for _, col := range collections {
		count, err := s.db.Collection(col.name).CountDocuments(ctx, col.filter)
		if err != nil {
			return err
		}
		status := "✓"
		if count != col.expected {
			status = "⚠"
		}
		fmt.Printf("  %s %s: %d/%d\n", status, col.name, count, col.expected)
	}

	return nil
}
