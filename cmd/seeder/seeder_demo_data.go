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

// DemoDataSeeder 演示数据填充器
type DemoDataSeeder struct {
	db *utils.Database
}

// NewDemoDataSeeder 创建演示数据填充器
func NewDemoDataSeeder(db *utils.Database) *DemoDataSeeder {
	return &DemoDataSeeder{db: db}
}

// CreateDemoData 创建演示数据
func (s *DemoDataSeeder) CreateDemoData() error {

	// 1. 检查项目是否已存在
	if s.projectExists("星际觉醒") {
		fmt.Println("✓ 演示项目「星际觉醒」已存在，跳过创建")
		return nil
	}

	fmt.Println("🚀 开始创建演示数据...")
	fmt.Printf("📊 项目：星际觉醒\n")

	// 2. 获取/创建管理员用户
	adminID, err := s.getOrCreateAdminUser()
	if err != nil {
		return fmt.Errorf("获取管理员失败: %w", err)
	}
	fmt.Printf("✓ 管理员账号: %s\n", adminID)

	// 3. 创建项目
	projectID, err := s.createProject(adminID)
	if err != nil {
		return fmt.Errorf("创建项目失败: %w", err)
	}
	fmt.Printf("✓ 创建项目: %s\n", projectID)

	// 4. 创建文档树（卷 → 章）
	docIDMap, err := s.createDocuments(projectID)
	if err != nil {
		return fmt.Errorf("创建文档失败: %w", err)
	}
	fmt.Printf("✓ 创建文档: 4卷 + 24章\n")

	// 5. 创建章节内容
	if err := s.createChapterContents(projectID, docIDMap); err != nil {
		return fmt.Errorf("创建章节内容失败: %w", err)
	}
	fmt.Println("✓ 创建章节内容")

	// 6. 创建大纲树
	if err := s.createOutlines(projectID, docIDMap); err != nil {
		return fmt.Errorf("创建大纲失败: %w", err)
	}
	fmt.Println("✓ 创建大纲树: 1级 → 4级2级 → 16个3级")

	// 7. 创建角色
	characterIDs, err := s.createCharacters(projectID)
	if err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}
	fmt.Printf("✓ 创建角色: %d个\n", len(characterIDs))

	// 8. 创建关系网络
	if err := s.createRelations(projectID, characterIDs); err != nil {
		return fmt.Errorf("创建关系失败: %w", err)
	}
	fmt.Println("✓ 创建关系网络")

	// 9. 创建其他资产
	if err := s.createAssets(projectID, characterIDs); err != nil {
		return fmt.Errorf("创建资产失败: %w", err)
	}
	fmt.Println("✓ 创建其他资产: 8道具 + 6地点 + 3时间线")

	// 10. 验证数据完整性
	if err := s.validateDemoData(projectID); err != nil {
		return fmt.Errorf("数据验证失败: %w", err)
	}

	fmt.Println("\n✅ 演示数据创建完成！")
	fmt.Println("📖 项目名称: 星际觉醒")
	fmt.Println("🔐 登录账号: admin / qingyu2024")

	return nil
}

// Clean 清空演示数据
func (s *DemoDataSeeder) Clean() error {
	ctx := context.Background()

	// 获取项目ID
	var project struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := s.db.Collection("projects").FindOne(ctx, bson.M{"title": "星际觉醒"}).Decode(&project)
	if err != nil {
		fmt.Println("✓ 演示项目不存在，无需清理")
		return nil
	}

	projectID := project.ID.Hex()

	// 按依赖顺序删除
	collections := []string{
		"character_relations",
		"characters",
		"timelines",
		"locations",
		"items",
		"outlines",
		"document_contents",
		"documents",
		"projects",
	}

	for _, colName := range collections {
		filter := bson.M{"project_id": projectID}
		if colName == "projects" {
			filter = bson.M{"_id": project.ID}
		}
		result, err := s.db.Collection(colName).DeleteMany(ctx, filter)
		if err != nil {
			return fmt.Errorf("清空 %s 失败: %w", colName, err)
		}
		if result.DeletedCount > 0 {
			fmt.Printf("✓ 清空 %s: %d条\n", colName, result.DeletedCount)
		}
	}

	fmt.Println("✓ 演示数据清理完成")
	return nil
}

// projectExists 检查项目是否存在
func (s *DemoDataSeeder) projectExists(title string) bool {
	ctx := context.Background()
	count, _ := s.db.Collection("projects").CountDocuments(ctx, bson.M{"title": title})
	return count > 0
}

// getOrCreateAdminUser 获取或创建管理员
func (s *DemoDataSeeder) getOrCreateAdminUser() (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("users")

	// 按优先级查找
	usernames := []string{"admin", "testadmin001"}

	for _, username := range usernames {
		var user struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err == nil {
			fmt.Printf("✓ 使用已存在账号: %s\n", username)
			// 确保 password 字段使用 bcrypt 哈希（修复历史明文密码）
			hashedPassword, hashErr := utils.HashPassword("qingyu2024")
			if hashErr == nil {
				collection.UpdateByID(ctx, user.ID, bson.M{
					"$set": bson.M{"password": hashedPassword},
				})
			}
			return user.ID.Hex(), nil
		}
	}

	// 不存在则创建
	return s.createAdminUser()
}

// createAdminUser 创建管理员用户
func (s *DemoDataSeeder) createAdminUser() (string, error) {
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
		"_id":         id,
		"username":    "admin",
		"email":       "admin@qingyu.com",
		"password":    hashedPassword,
		"roles":       []string{"admin"},
		"status":      "active",
		"nickname":    "管理员",
		"avatar":      "/images/avatars/default.png",
		"bio":         "演示账号管理员",
		"created_at":  now,
		"updated_at":  now,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	fmt.Println("✓ 创建管理员账号: admin")
	return id.Hex(), nil
}

// createProject 创建项目
func (s *DemoDataSeeder) createProject(authorID string) (string, error) {
	ctx := context.Background()
	collection := s.db.Collection("projects")

	now := time.Now()
	id := primitive.NewObjectID()

	project := bson.M{
		"_id":            id,
		"author_id":      func() primitive.ObjectID { oid, _ := primitive.ObjectIDFromHex(authorID); return oid }(),
		"title":          story.DemoProject.Title,
		"summary":        story.DemoProject.Summary,
		"cover_url":      story.DemoProject.CoverURL,
		"category":       story.DemoProject.Category,
		"writing_type":   story.DemoProject.WritingType,
		"status":         story.DemoProject.Status,
		"visibility":     story.DemoProject.Visibility,
		"tags":           story.DemoProject.Tags,
		"statistics": bson.M{
			"total_words":    0,
			"chapter_count":  24,
			"document_count":  28,
			"last_update_at": now,
		},
		"settings": bson.M{
			"auto_backup":     true,
			"backup_interval": 24,
		},
		"collaborators": []interface{}{},
		"created_at":     now,
		"updated_at":     now,
	}

	_, err := collection.InsertOne(ctx, project)
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

// createDocuments 创建文档树
func (s *DemoDataSeeder) createDocuments(projectID string) (map[string]string, error) {
	ctx := context.Background()
	collection := s.db.Collection("documents")

	now := time.Now()
	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	volumeIDs := make(map[string]string)
	chapterIDs := make(map[string]string)

	// 创建卷
	for _, vol := range story.DemoVolumes {
		volID := primitive.NewObjectID()
		volumeDoc := bson.M{
			"_id":          volID,
			"project_id":    projectOID,
			"title":        vol.Title,
			"type":         "volume",
			"level":        1,
			"order":        vol.Order,
			"parent_id":    primitive.NilObjectID,
			"status":       "completed",
			"word_count":   0,
			"stable_ref":    primitive.NewObjectID().Hex(),
			"order_key":    fmt.Sprintf("%04d", vol.Order*10000),
			"created_at":   now,
			"updated_at":   now,
		}
		_, err := collection.InsertOne(ctx, volumeDoc)
		if err != nil {
			return nil, err
		}
		volumeIDs[vol.Title] = volID.Hex()

		// 创建该卷下的章节
		for _, ch := range vol.Chapters {
			chID := primitive.NewObjectID()
			chapterDoc := bson.M{
				"_id":          chID,
				"project_id":    projectOID,
				"title":        ch.Title,
				"type":         "chapter",
				"level":        2,
				"order":        ch.Order,
				"parent_id":    volID,
				"status":       "planned",
				"word_count":   0,
				"stable_ref":   primitive.NewObjectID().Hex(),
				"order_key":    fmt.Sprintf("%04d%02d%04d", vol.Order*10000, 0, ch.Order*1000),
				"created_at":   now,
				"updated_at":   now,
			}
			_, err = collection.InsertOne(ctx, chapterDoc)
			if err != nil {
				return nil, err
			}
			chapterIDs[ch.Title] = chID.Hex()
		}
	}

	// 合并两个 map
	result := make(map[string]string)
	for k, v := range volumeIDs {
		result[k] = v
	}
	for k, v := range chapterIDs {
		result[k] = v
	}

	return result, nil
}

// createChapterContents 创建章节内容
func (s *DemoDataSeeder) createChapterContents(projectID string, docIDMap map[string]string) error {
	ctx := context.Background()
	collection := s.db.Collection("document_contents")

	now := time.Now()

	for _, vol := range story.DemoVolumes {
		for _, ch := range vol.Chapters {
			docIDStr, ok := docIDMap[ch.Title]
			if !ok {
				continue
			}

			docOID, err := primitive.ObjectIDFromHex(docIDStr)
			if err != nil {
				continue
			}

			// 生成章节内容
			content := generateDemoChapterContent(ch.Title, vol.Title, ch.Summary)

			docContent := bson.M{
				"_id":            primitive.NewObjectID(),
				"document_id":    docOID,
				"content":        content,
				"content_type":   "markdown",
				"word_count":     len([]rune(content)),
				"char_count":     len(content),
				"paragraph_order": 0,
				"version":        1,
				"last_saved_at":  now,
				"last_edited_by": "system",
				"created_at":     now,
				"updated_at":     now,
			}

			_, err = collection.InsertOne(ctx, docContent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// createOutlines 创建大纲树
func (s *DemoDataSeeder) createOutlines(projectID string, docIDMap map[string]string) error {
	ctx := context.Background()
	collection := s.db.Collection("outlines")

	now := time.Now()

	// 1. 创建1级大纲（全局总纲）
	outline1ID := primitive.NewObjectID()
	outline1 := bson.M{
		"_id":         outline1ID,
		"project_id":  projectID,
		"title":       "星际觉醒：故事总纲",
		"parent_id":   "",
		"order":       0,
		"type":        "global",
		"tension":     5,
		"summary":     story.DemoProject.Summary,
		"document_id": "",
		"characters":  []string{},
		"items":       []string{},
		"created_at":  now,
		"updated_at":  now,
	}
	_, err := collection.InsertOne(ctx, outline1)
	if err != nil {
		return err
	}

	// 2. 创建2级大纲（卷级）
	outline2IDs := make([]primitive.ObjectID, len(story.DemoVolumes))
	for i, vol := range story.DemoVolumes {
		volID, ok := docIDMap[vol.Title]
		if !ok {
			continue
		}

		outline2ID := primitive.NewObjectID()
		outline2 := bson.M{
			"_id":         outline2ID,
			"project_id":  projectID,
			"title":       vol.Title + "大纲",
			"parent_id":   outline1ID.Hex(),
			"order":       i,
			"type":        "arc",
			"tension":     5 + i,
			"summary":     vol.Summary,
			"document_id": volID,
			"characters":  []string{},
			"items":       []string{},
			"created_at":  now,
			"updated_at":  now,
		}
		_, err := collection.InsertOne(ctx, outline2)
		if err != nil {
			return err
		}
		outline2IDs[i] = outline2ID

		// 反向更新文档的 outline_node_id（双向映射）
		docCol := s.db.Collection("documents")
		volDocOID, _ := primitive.ObjectIDFromHex(volID)
		docCol.UpdateByID(ctx, volDocOID, bson.M{
			"$set": bson.M{"outline_node_id": outline2ID.Hex()},
		})
	}

	// 3. 创建3级大纲（场景级）
	for volIdx, vol := range story.DemoVolumes {
		if volIdx >= len(outline2IDs) {
			continue
		}
		parent2ID := outline2IDs[volIdx].Hex()
		outline3Count := 0

		for _, ch := range vol.Chapters {
			if !ch.HasOutline3 {
				continue
			}

			// 获取章节绑定的文档ID
			chDocID, ok := docIDMap[ch.Title]
			if !ok {
				chDocID = ""
			}

			// 获取章节相关角色
			characters := getChapterCharacters(ch.Title)

			outline3 := bson.M{
				"_id":         primitive.NewObjectID(),
				"project_id":  projectID,
				"title":       ch.Title + "场景大纲",
				"parent_id":   parent2ID,
				"order":       outline3Count,
				"type":        "scene",
				"tension":     3 + (ch.Order % 5),
				"summary":     ch.Summary,
				"document_id": chDocID,
				"characters":  characters,
				"items":       []string{},
				"created_at":  now,
				"updated_at":  now,
			}
			outline3Doc := outline3["_id"].(primitive.ObjectID)
			_, err := collection.InsertOne(ctx, outline3)
			if err != nil {
				return err
			}

			// 反向更新章节文档的 outline_node_id（双向映射）
			if chDocID != "" {
				docCol := s.db.Collection("documents")
				chDocOID, _ := primitive.ObjectIDFromHex(chDocID)
				docCol.UpdateByID(ctx, chDocOID, bson.M{
					"$set": bson.M{"outline_node_id": outline3Doc.Hex()},
				})
			}

			outline3Count++
		}
	}

	return nil
}

// createCharacters 创建角色
func (s *DemoDataSeeder) createCharacters(projectID string) (map[string]string, error) {
	ctx := context.Background()
	collection := s.db.Collection("characters")

	now := time.Now()
	projectOID, _ := primitive.ObjectIDFromHex(projectID)
	characterIDs := make(map[string]string)

	for _, char := range story.DemoCharacters {
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
			"short_description":  char.ShortDescription,
			"personality_prompt": "",
			"speech_pattern":     "",
			"current_state":      "",
			"created_at":        now,
			"updated_at":        now,
		}
		_, err := collection.InsertOne(ctx, character)
		if err != nil {
			return nil, err
		}
		characterIDs[char.ID] = charID.Hex()
	}

	return characterIDs, nil
}

// createRelations 创建关系网络
func (s *DemoDataSeeder) createRelations(projectID string, characterIDs map[string]string) error {
	ctx := context.Background()
	collection := s.db.Collection("character_relations")

	now := time.Now()
	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	for _, rel := range story.DemoRelations {
		fromID, ok := characterIDs[rel.FromID]
		if !ok {
			continue
		}
		toID, ok := characterIDs[rel.ToID]
		if !ok {
			continue
		}

		relation := bson.M{
			"_id":         primitive.NewObjectID(),
			"project_id":  projectOID,
			"from_id":     fromID,
			"to_id":       toID,
			"type":        rel.Type,
			"strength":    rel.Strength,
			"notes":       rel.Notes,
			"created_at":  now,
			"updated_at":  now,
		}
		_, err := collection.InsertOne(ctx, relation)
		if err != nil {
			return err
		}

		// 双向关系：如果是朋友/盟友/恋人，也创建反向关系
		if rel.Type == "朋友" || rel.Type == "盟友" || rel.Type == "恋人" || rel.Type == "家庭" {
			reverseRelation := bson.M{
				"_id":         primitive.NewObjectID(),
				"project_id":  projectOID,
				"from_id":     toID,
				"to_id":       fromID,
				"type":        rel.Type,
				"strength":    rel.Strength,
				"notes":       rel.Notes,
				"created_at":  now,
				"updated_at":  now,
			}
			collection.InsertOne(ctx, reverseRelation)
		}
	}

	return nil
}

// createAssets 创建其他资产
func (s *DemoDataSeeder) createAssets(projectID string, characterIDs map[string]string) error {
	now := time.Now()

	// 创建道具
	if err := s.createItems(projectID, characterIDs, now); err != nil {
		return err
	}

	// 创建地点
	if err := s.createLocations(projectID, now); err != nil {
		return err
	}

	// 创建时间线
	if err := s.createTimelines(projectID, now); err != nil {
		return err
	}

	return nil
}

// createItems 创建道具
func (s *DemoDataSeeder) createItems(projectID string, characterIDs map[string]string, now time.Time) error {
	ctx := context.Background()
	collection := s.db.Collection("items")

	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	for _, item := range story.DemoItems {
		ownerID := ""
		if item.OwnerID != "" {
			ownerID = characterIDs[item.OwnerID]
		}

		itemDoc := bson.M{
			"_id":          primitive.NewObjectID().Hex(),
			"project_id":   projectOID,
			"name":         item.Name,
			"type":         item.Type,
			"description":  item.Description,
			"owner_id":     ownerID,
			"location_id":  item.LocationID,
			"rarity":       item.Rarity,
			"function":     item.Function,
			"origin":       item.Origin,
			"created_at":   now,
			"updated_at":   now,
		}
		_, err := collection.InsertOne(ctx, itemDoc)
		if err != nil {
			return err
		}
	}

	return nil
}

// createLocations 创建地点
func (s *DemoDataSeeder) createLocations(projectID string, now time.Time) error {
	ctx := context.Background()
	collection := s.db.Collection("locations")

	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	for _, loc := range story.DemoLocations {
		locDoc := bson.M{
			"_id":          primitive.NewObjectID(),
			"project_id":   projectOID,
			"name":         loc.Name,
			"description":  loc.Description,
			"climate":      loc.Climate,
			"culture":      loc.Culture,
			"geography":    loc.Geography,
			"atmosphere":   loc.Atmosphere,
			"parent_id":    "",
			"image_url":    "",
			"created_at":   now,
			"updated_at":   now,
		}
		_, err := collection.InsertOne(ctx, locDoc)
		if err != nil {
			return err
		}
	}

	return nil
}

// createTimelines 创建时间线
func (s *DemoDataSeeder) createTimelines(projectID string, now time.Time) error {
	ctx := context.Background()
	timelineCol := s.db.Collection("timelines")
	eventCol := s.db.Collection("timeline_events")

	projectOID, _ := primitive.ObjectIDFromHex(projectID)

	for _, tl := range story.DemoTimelines {
		timelineID := primitive.NewObjectID()
		timelineDoc := bson.M{
			"_id":          timelineID,
			"project_id":   projectOID,
			"name":         tl.Name,
			"description":  tl.Description,
			"start_time":   nil,
			"end_time":     nil,
			"created_at":   now,
			"updated_at":   now,
		}
		_, err := timelineCol.InsertOne(ctx, timelineDoc)
		if err != nil {
			return err
		}

		// 创建时间线事件
		for _, event := range tl.Events {
			eventDoc := bson.M{
				"_id":           primitive.NewObjectID(),
				"project_id":    projectOID,
				"timeline_id":   timelineID.Hex(),
				"title":         event.Title,
				"description":   event.Description,
				"story_time": bson.M{
					"year":  event.Year,
					"month": event.Month,
					"day":   event.Day,
				},
				"importance":    event.Importance,
				"event_type":    "plot",
				"created_at":    now,
				"updated_at":    now,
			}
			_, err := eventCol.InsertOne(ctx, eventDoc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// validateDemoData 验证数据完整性
func (s *DemoDataSeeder) validateDemoData(projectID string) error {
	ctx := context.Background()

	// 检查各集合数量
	collections := []struct {
		name     string
		expected int64
	}{
		{"documents", 28},
		{"characters", 12},
		{"outlines", 21},
		{"items", 8},
		{"locations", 6},
	}

	fmt.Println("\n📊 数据完整性验证:")
	for _, col := range collections {
		count, err := s.db.Collection(col.name).CountDocuments(ctx, bson.M{"project_id": projectID})
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

// generateDemoChapterContent 生成演示章节内容
func generateDemoChapterContent(chapterTitle, volumeTitle, summary string) string {
	return fmt.Sprintf(`# %s

## 章节概述

%s

---

## 正文

*本章为演示数据自动生成的内容，展示了《星际觉醒》故事中的精彩片段。*

### 场景一：命运的转折

夜幕降临，星光洒落在火星表面的红色沙丘上。林风站在考古遗址的入口处，手中紧握着刚刚发现的先民遗物。

这枚古老的晶体散发着微弱的光芒，仿佛在诉说着一段被遗忘的历史。

"林风，你发现了什么？"苏文的声音从身后传来。

林风转过头，眼中闪烁着兴奋与忧虑交织的光芒。"苏文，我觉得...我们可能发现了一些不得了的东西。"

### 场景二：真相浮现

遗迹内部的空气干燥而寒冷，墙壁上刻满了古老的符文。陈博士正在用专业设备进行扫描，神情严肃。

"这些符号..."陈博士的声音微微颤抖，"如果我的判断没错，这可能是外星文明留下的遗迹。"

林风和苏文对视一眼，都从对方眼中看到了震惊。

### 场景三：危机降临

突然，遗迹深处传来一阵低沉的轰鸣声。整个地面开始震动，碎石从顶部落下。

"快走！"陈博士大喊。

三人拼命向外跑去，但通道已经开始崩塌。林风感觉到手中的遗物突然变得滚烫，一道耀眼的光芒从晶体中迸发而出...

---

*欲知后事如何，请阅读下一章。*

---
**字数统计**: 约1200字
`, chapterTitle, summary)
}

// getChapterCharacters 获取章节相关角色
func getChapterCharacters(chapterTitle string) []string {
	charactersMap := map[string][]string{
		"火星遗迹":     {"char-linfeng", "char-suwen", "char-chen"},
		"能力觉醒":     {"char-linfeng", "char-suwen"},
		"政府介入":     {"char-linfeng", "char-chen", "char-general-li"},
		"逃亡之路":     {"char-linfeng", "char-suwen", "char-chen"},
		"追捕与反追捕":   {"char-linfeng", "char-suwen"},
		"抉择时刻":     {"char-linfeng", "char-suwen"},
		"火星独立宣言":   {"char-linfeng", "char-reynold", "char-ava"},
		"暗流涌动":     {"char-linfeng", "char-reynold", "char-ava", "char-zhanggong"},
		"内战爆发":     {"char-linfeng", "char-ava", "char-general-li"},
		"血腥谈判":     {"char-linfeng", "char-ava", "char-reynold"},
		"真相浮现":     {"char-linfeng", "char-memory"},
		"停火协议":     {"char-linfeng", "char-reynold", "char-ava", "char-general-li"},
		"外星舰队":     {"char-linfeng", "char-messenger"},
		"恐惧与敌意":    {"char-linfeng", "char-hawk", "char-un-sg"},
		"跨越语言的沟通":  {"char-linfeng", "char-messenger", "char-suwen"},
		"建立信任":     {"char-linfeng", "char-messenger", "char-un-sg"},
		"威胁逼近":     {"char-linfeng", "char-hawk", "char-messenger"},
		"危机化解":     {"char-linfeng", "char-messenger"},
		"联盟谈判":     {"char-linfeng", "char-reynold", "char-ava", "char-un-sg", "char-messenger"},
		"最后的阻碍":    {"char-linfeng", "char-hawk", "char-un-sg"},
		"刺杀危机":     {"char-linfeng", "char-messenger", "char-xiaolin"},
		"真相大白":     {"char-linfeng", "char-memory"},
		"联盟成立":     {"char-linfeng", "char-suwen", "char-messenger", "char-un-sg"},
		"新的起点":     {"char-linfeng", "char-suwen", "char-memory"},
	}

	if chars, ok := charactersMap[chapterTitle]; ok {
		return chars
	}
	return []string{"char-linfeng"}
}
