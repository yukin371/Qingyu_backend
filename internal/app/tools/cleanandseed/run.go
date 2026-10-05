package cleanandseed

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"Qingyu_backend/internal/app/tools/mongodbconfig"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Run previews or explicitly recreates the complete demonstration dataset.
func Run(args []string) error {
	apply, err := parseApply(args)
	if err != nil {
		return err
	}
	if !apply {
		fmt.Println("预览模式：clean_and_seed 将清空演示集合并重建数据；使用 --apply 并设置安全环境变量后执行")
		return nil
	}
	if os.Getenv("QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA") != "true" {
		return fmt.Errorf("重建演示数据需要设置 QINGYU_ALLOW_DESTRUCTIVE_TEST_DATA=true")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbconfig.URI()))
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping 失败: %w", err)
	}
	db := client.Database(mongodbconfig.Database())

	// Clean all collections
	cols := []string{
		"character_relations", "characters", "timelines", "timeline_events",
		"locations", "location_relations", "items", "outlines",
		"document_contents", "documents", "projects",
	}
	for _, col := range cols {
		r, err := db.Collection(col).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("清理集合 %s 失败: %w", col, err)
		}
		fmt.Printf("Cleaned %s: %d\n", col, r.DeletedCount)
	}

	// Get admin user
	var existingAdmin bson.M
	if err := db.Collection("users").FindOne(ctx, bson.M{"username": "admin"}).Decode(&existingAdmin); err != nil {
		return fmt.Errorf("读取 admin 用户失败: %w", err)
	}
	adminID, ok := existingAdmin["_id"].(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("admin 用户 _id 类型无效: %T", existingAdmin["_id"])
	}
	fmt.Printf("Admin ID: %s\n", adminID.Hex())

	now := time.Now()

	// Create project
	projectID := primitive.NewObjectID()
	project := bson.M{
		"_id":          projectID,
		"author_id":    adminID,
		"title":        "星际觉醒",
		"summary":      "在人类迈向星际的时代，一支考古队在火星发现了远古先民遗留的神秘遗物，由此揭开了一段横跨银河系文明的惊天秘密。",
		"cover_url":    "",
		"category":     "科幻",
		"writing_type": "novel",
		"status":       "serializing",
		"visibility":   "public",
		"tags":         []string{"科幻", "冒险", "星际文明"},
		"statistics": bson.M{
			"total_words":    0,
			"chapter_count":  24,
			"document_count": 28,
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
	if err := insertDocument(ctx, db, "projects", project); err != nil {
		return err
	}
	fmt.Printf("Project: %s\n", projectID.Hex())

	// Global outline (root)
	outline1ID := primitive.NewObjectID()
	if err := insertDocument(ctx, db, "outlines", bson.M{
		"_id": outline1ID, "project_id": projectID, "title": "星际觉醒：故事总纲",
		"parent_id": "", "order": 0, "type": "global", "tension": 5,
		"summary": project["summary"], "document_id": "", "characters": bson.A{},
		"items": bson.A{}, "tags": bson.A{},
		"created_at": now, "updated_at": now,
	}); err != nil {
		return err
	}

	volTitles := []string{"第一卷：火星遗迹", "第二卷：觉醒之路", "第三卷：星际冲突", "第四卷：命运抉择"}
	volSummaries := []string{
		"考古队发现先民遗迹，林风获得能力",
		"林风能力觉醒，火星独立运动",
		"外星文明接触，星际战争危机",
		"联盟成立，人类新纪元",
	}
	chapterTitles := [][]string{
		{"火星遗迹", "能力觉醒", "政府介入", "逃亡之路", "追捕与反追捕", "抉择时刻"},
		{"火星独立宣言", "暗流涌动", "内战爆发", "血腥谈判", "真相浮现", "停火协议"},
		{"外星舰队", "恐惧与敌意", "跨越语言的沟通", "建立信任", "威胁逼近", "危机化解"},
		{"联盟谈判", "最后的阻碍", "刺杀危机", "真相大白", "联盟成立", "新的起点"},
	}

	totalDocs := 0
	for vi, volTitle := range volTitles {
		volID := primitive.NewObjectID()
		volOutlineID := primitive.NewObjectID()
		volOrder := vi + 1

		// Volume document
		if err := insertDocument(ctx, db, "documents", bson.M{
			"_id": volID, "project_id": projectID, "title": volTitle,
			"type": "volume", "level": 1, "order": volOrder,
			"parent_id": primitive.NilObjectID, "status": "completed",
			"word_count": 0, "stable_ref": primitive.NewObjectID().Hex(),
			"order_key":       fmt.Sprintf("%04d", volOrder*10000),
			"outline_node_id": volOutlineID.Hex(),
			"author_id":       adminID,
			"character_ids":   bson.A{}, "location_ids": bson.A{},
			"timeline_ids": bson.A{}, "tags": bson.A{}, "notes": "",
			"created_at": now, "updated_at": now,
		}); err != nil {
			return err
		}
		totalDocs++

		// Volume outline
		if err := insertDocument(ctx, db, "outlines", bson.M{
			"_id": volOutlineID, "project_id": projectID, "title": volTitle + "大纲",
			"parent_id": outline1ID.Hex(), "order": vi, "type": "arc",
			"tension": 5 + vi, "summary": volSummaries[vi],
			"document_id": volID.Hex(), "characters": bson.A{},
			"items": bson.A{}, "tags": bson.A{},
			"created_at": now, "updated_at": now,
		}); err != nil {
			return err
		}

		for ci, chTitle := range chapterTitles[vi] {
			chID := primitive.NewObjectID()
			chOutlineID := primitive.NewObjectID()
			chOrder := ci + 1

			// Chapter document
			if err := insertDocument(ctx, db, "documents", bson.M{
				"_id": chID, "project_id": projectID, "title": chTitle,
				"type": "chapter", "level": 2, "order": chOrder,
				"parent_id": volID, "status": "planned",
				"word_count": 0, "stable_ref": primitive.NewObjectID().Hex(),
				"order_key":       fmt.Sprintf("%04d%02d%04d", volOrder*10000, 0, chOrder*1000),
				"outline_node_id": chOutlineID.Hex(),
				"author_id":       adminID,
				"character_ids":   bson.A{}, "location_ids": bson.A{},
				"timeline_ids": bson.A{}, "tags": bson.A{}, "notes": "",
				"created_at": now, "updated_at": now,
			}); err != nil {
				return err
			}
			totalDocs++

			// Chapter outline
			if err := insertDocument(ctx, db, "outlines", bson.M{
				"_id": chOutlineID, "project_id": projectID, "title": chTitle + "场景大纲",
				"parent_id": volOutlineID.Hex(), "order": ci, "type": "scene",
				"tension": 3 + (chOrder % 5), "summary": chTitle + "的精彩故事",
				"document_id": chID.Hex(), "characters": bson.A{},
				"items": bson.A{}, "tags": bson.A{},
				"created_at": now, "updated_at": now,
			}); err != nil {
				return err
			}

			// Chapter content
			content := fmt.Sprintf("# %s\n\n## 章节概述\n\n%s\n\n---\n\n*本章为演示数据自动生成的内容*", chTitle, chTitle+"的精彩故事")
			if err := insertDocument(ctx, db, "document_contents", bson.M{
				"_id": primitive.NewObjectID(), "document_id": chID,
				"content": content, "content_type": "markdown",
				"word_count": len([]rune(content)), "char_count": len(content),
				"paragraph_order": 0, "version": 1,
				"last_saved_at": now, "last_edited_by": "system",
				"created_at": now, "updated_at": now,
			}); err != nil {
				return err
			}
		}
		fmt.Printf("Volume %d: %s (%d chapters)\n", vi+1, volTitle, len(chapterTitles[vi]))
	}

	// Characters (8)
	charData := []struct{ Name, Summary string }{
		{"林风", "主角，考古学家，后发现先民遗物获得超能力"},
		{"苏文", "林风的挚友和搭档，语言学家"},
		{"陈博士", "资深科学家，研究先民文明的权威"},
		{"雷诺将军", "地球联邦军方鹰派代表"},
		{"艾娃", "火星AI系统，逐渐获得自我意识"},
		{"李将军", "地球联邦军方鸽派代表"},
		{"信使", "外星文明的和平使者"},
		{"记忆者", "保存先民记忆的神秘存在"},
	}
	for _, c := range charData {
		if err := insertDocument(ctx, db, "characters", bson.M{
			"_id": primitive.NewObjectID(), "project_id": projectID,
			"name": c.Name, "alias": bson.A{}, "summary": c.Summary,
			"traits": bson.A{"勇敢", "聪明"}, "background": c.Summary,
			"avatar_url": "", "short_description": c.Summary,
			"personality_prompt": "", "speech_pattern": "", "current_state": "",
			"author_id":  adminID,
			"created_at": now, "updated_at": now,
		}); err != nil {
			return err
		}
	}

	// Locations (6)
	locData := []struct{ Name, Desc string }{
		{"火星考古基地", "位于火星赤道附近的考古挖掘基地"},
		{"火星首都星城", "火星最大的城市，政治经济中心"},
		{"地球联邦总部", "地球联邦政府所在地"},
		{"先民遗迹", "先民遗留的神秘建筑群"},
		{"星际空间站", "人类在太阳系边缘建造的最大空间站"},
		{"外星母舰", "外星文明的主力飞船"},
	}
	for _, l := range locData {
		if err := insertDocument(ctx, db, "locations", bson.M{
			"_id": primitive.NewObjectID(), "project_id": projectID,
			"name": l.Name, "description": l.Desc,
			"climate": "", "culture": "", "geography": "",
			"atmosphere": "", "parent_id": "", "image_url": "",
			"author_id":  adminID,
			"created_at": now, "updated_at": now,
		}); err != nil {
			return err
		}
	}

	// Timelines (3) with events
	tlData := []struct {
		Name, Desc string
		Events     []struct {
			Title, Desc           string
			Year, Month, Day, Imp int
		}
	}{
		{"主线时间线", "星际觉醒主线故事时间线", []struct {
			Title, Desc           string
			Year, Month, Day, Imp int
		}{
			{"发现遗迹", "考古队在火星发现先民遗迹", 2187, 3, 15, 9},
			{"能力觉醒", "林风接触遗物获得超能力", 2187, 4, 1, 8},
			{"政府介入", "地球联邦派出军队控制遗迹", 2187, 5, 20, 7},
			{"联盟成立", "人类-先民联盟正式成立", 2189, 1, 1, 10},
		}},
		{"火星独立运动", "火星从殖民地到独立的过程", []struct {
			Title, Desc           string
			Year, Month, Day, Imp int
		}{
			{"火星独立宣言", "林风宣布火星独立", 2188, 1, 1, 9},
			{"内战爆发", "地球联邦对火星发起军事行动", 2188, 3, 15, 8},
			{"停火协议", "双方签订停火协议", 2188, 6, 1, 7},
		}},
		{"星际接触", "与外星文明接触的关键事件", []struct {
			Title, Desc           string
			Year, Month, Day, Imp int
		}{
			{"外星舰队抵达", "外星舰队进入太阳系", 2188, 8, 1, 10},
			{"首次沟通", "林风与外星信使首次对话", 2188, 9, 1, 9},
			{"危机化解", "成功化解星际战争危机", 2188, 11, 1, 8},
		}},
	}
	for _, tl := range tlData {
		tlID := primitive.NewObjectID()
		if err := insertDocument(ctx, db, "timelines", bson.M{
			"_id": tlID, "project_id": projectID,
			"name": tl.Name, "description": tl.Desc,
			"start_time": nil, "end_time": nil,
			"author_id":  adminID,
			"created_at": now, "updated_at": now,
		}); err != nil {
			return err
		}
		for _, ev := range tl.Events {
			if err := insertDocument(ctx, db, "timeline_events", bson.M{
				"_id": primitive.NewObjectID(), "project_id": projectID,
				"timeline_id": tlID.Hex(), "title": ev.Title,
				"description": ev.Desc,
				"story_time":  bson.M{"year": ev.Year, "month": ev.Month, "day": ev.Day},
				"importance":  ev.Imp, "event_type": "plot",
				"author_id":  adminID,
				"created_at": now, "updated_at": now,
			}); err != nil {
				return err
			}
		}
	}

	// Items (8)
	itemData := []struct{ Name, Type, Desc, Rarity, Origin string }{
		{"先民遗物", "key_item", "先民遗留的神秘晶体，蕴含巨大能量", "legendary", "先民遗物"},
		{"星际通讯器", "tool", "可以跨越星系进行通讯的设备", "rare", "先民技术"},
		{"能量水晶", "material", "先民文明的能量存储装置", "rare", "先民制造"},
		{"先民数据库", "key_item", "记录先民全部知识的数据库", "legendary", "先民遗物"},
		{"隐身装置", "tool", "先民技术的隐身设备", "epic", "先民技术"},
		{"时空门钥匙", "key_item", "开启时空门的钥匙", "legendary", "先民遗物"},
		{"防护盾", "equipment", "先民防护技术装置", "epic", "先民技术"},
		{"翻译矩阵", "tool", "自动翻译各种语言的AI矩阵", "rare", "先民技术"},
	}
	for _, it := range itemData {
		if err := insertDocument(ctx, db, "items", bson.M{
			"_id": primitive.NewObjectID().Hex(), "project_id": projectID,
			"name": it.Name, "type": it.Type, "description": it.Desc,
			"owner_id": "", "location_id": "",
			"rarity": it.Rarity, "function": "", "origin": it.Origin,
			"author_id":  adminID,
			"created_at": now, "updated_at": now,
		}); err != nil {
			return err
		}
	}

	fmt.Println("\n=== Demo data seeded ===")
	fmt.Printf("Project: %s\n", projectID.Hex())
	fmt.Println("Login: admin / qingyu2024")

	// Verify
	fmt.Println("\nVerification:")
	for _, col := range []string{"documents", "outlines", "characters", "locations", "timelines", "items", "document_contents"} {
		count, err := db.Collection(col).CountDocuments(ctx, bson.M{"project_id": projectID})
		if err != nil {
			return fmt.Errorf("验证集合 %s 失败: %w", col, err)
		}
		fmt.Printf("  %s: %d\n", col, count)
	}
	return nil
}

func parseApply(args []string) (bool, error) {
	fs := flag.NewFlagSet("clean_and_seed", flag.ContinueOnError)
	apply := fs.Bool("apply", false, "清空并重建演示数据")
	if err := fs.Parse(args); err != nil {
		return false, err
	}
	return *apply, nil
}

func insertDocument(ctx context.Context, db *mongo.Database, collection string, document interface{}) error {
	if _, err := db.Collection(collection).InsertOne(ctx, document); err != nil {
		return fmt.Errorf("写入集合 %s 失败: %w", collection, err)
	}
	return nil
}
