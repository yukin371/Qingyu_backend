package seeds

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Chapter 章节
type Chapter struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	BookID      string             `bson:"book_id"`
	ChapterNum  int                `bson:"chapter_num"`
	Title       string             `bson:"title"`
	WordCount   int                `bson:"word_count"`
	Price       float64            `bson:"price"`
	IsFree      bool               `bson:"is_free"`
	Status      string             `bson:"status"` // draft, published, deleted
	PublishedAt time.Time          `bson:"published_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// ChapterContent 章节内容
type ChapterContent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChapterID string             `bson:"chapter_id"`
	Content   string             `bson:"content"`
	WordCount int                `bson:"word_count"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// 示例章节内容
var sampleChapterContents = []string{
	`晨光熹微，透过破旧的窗棂，洒落在少年苍白的脸上。

李凡缓缓睁开眼睛，看着熟悉的茅草屋顶，心中五味杂陈。这是他来到这个世界的第三个月。

前世，他是一个普通的上班族，每天朝九晚五，为了生活奔波。一场意外，让他穿越到了这个名为"苍澜大陆"的修真世界。

"既然老天让我重活一次，我定要活出个模样来！"李凡握紧拳头，眼中闪过一丝坚定。

这是一个弱肉强食的世界，只有拥有力量，才能掌控自己的命运。

他想起昨天村长说的话，明日便是苍澜宗招收弟子的日子。这是整个青云城十年一次的大事，所有十二到十六岁的少年都可以参加。

李凡今年十五岁，正是最后一次机会。

"苍澜宗，我来了！"

他起身，整理了一下破旧的衣服，推门而出。`,
	`天空中飘着细雨，苍澜宗的山门笼罩在烟雨朦胧之中。

李凡站在山脚下，抬头望去。只见巍峨的山峰直插云霄，云雾缭绕间，隐约可见亭台楼阁，宛如仙境。

"这便是修真门派吗？"他喃喃自语，心中既向往又忐忑。

周围聚集了数千名少年，个个神情激动。有人穿着锦衣华服，显然出身富贵之家；也有人如李凡一般，衣衫朴素，眼中却透着坚毅。

"测试开始！"

随着一声高喝，一名身着青衫的中年人出现在众人面前。他便是负责此次招生的执事长老。

"第一关，测骨龄。十二到十六岁者，方可通过。"

一块巨大的测灵石被抬了上来，散发着淡淡的光芒。

李凡深吸一口气，走上前去。

将手放在测灵石上，一道光芒闪过。

"骨龄十五岁，合格！"

听到这句话，李凡心中一松，至少第一关过了。`,
	`"第二关，测灵根。"

执事长老的声音传来，让李凡从回忆中回过神来。

灵根，是修行的根本。没有灵根，便无法感应天地灵气，更别提修炼了。

传说中，灵根分为天地玄黄四等，每等又分上中下三品。

天灵根，万中无一，修炼速度一日千里。
地灵根，千里挑一，同样是各门派争抢的对象。
玄灵根，百中有一，足以在修真界立足。
黄灵根，最为常见，十人之中便有一人。

而李凡的灵根，会是哪一种呢？

他忐忑地将手放在测灵石上。

时间一分一秒过去，测灵石没有任何反应。

一息，两息，三息...

"下一个。"

执事长老淡漠的声音，如同晴天霹雳，让李凡如坠冰窟。

没有灵根！

他呆呆地站在那里，难以置信。

"怎么会这样...明明我感觉到天地灵气的存在..."`,
	`就在李凡即将离去的时候，异变突生。

他怀中那块从小佩戴的黑色玉佩，突然发出一缕微弱的光芒。

这光芒虽然微弱，但在场的修士都感应到了。

"嗯？"

执事长老眼睛一亮，身形一闪，便出现在李凡面前。

"小子，你怀中之物，拿出来看看。"

李凡不解，但还是将玉佩取出。

那是一块古朴的黑玉，上面刻着复杂的纹路，看起来平平无奇。

但在执事长老眼中，这块玉佩却非同小可。

"这是...混沌玉？！"

他的声音都有些颤抖，眼中满是不可置信。

周围的人群顿时一片哗然。

混沌玉，传说中的至宝，可吞噬万物，孕育混沌之气。拥有此物者，即便没有灵根，亦可踏上修真之路。

"你，被苍澜宗录取了。"

执事长老深吸一口气，郑重地说道。

李凡愣住了。

这就...录取了？

修真之路，就在这一刻，为他敞开了大门。`,
}

// SeedChapters 章节种子数据
func SeedChapters(ctx context.Context, db *mongo.Database) error {
	chapterCollection := db.Collection("chapters")
	contentCollection := db.Collection("chapter_contents")

	fmt.Println("========================================")
	fmt.Println("开始创建章节和内容测试数据...")
	fmt.Println("========================================")

	// 获取书籍列表
	bookCollection := db.Collection("books")
	cursor, err := bookCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("获取书籍列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var books []struct {
		ID     string `bson:"_id"`
		Title  string `bson:"title"`
		Status string `bson:"status"`
	}
	if err = cursor.All(ctx, &books); err != nil {
		return fmt.Errorf("解析书籍列表失败: %w", err)
	}

	if len(books) == 0 {
		fmt.Println("警告：没有找到书籍，跳过章节创建")
		return nil
	}

	// 检查是否已有章节数据
	existingCount, _ := chapterCollection.CountDocuments(ctx, bson.M{})
	if existingCount > 0 {
		fmt.Printf("章节数据已存在 (%d条)，跳过创建\n", existingCount)
		return nil
	}

	now := time.Now()
	totalChapters := 0
	totalContents := 0

	for _, book := range books {
		fmt.Printf("处理书籍: %s...\n", book.Title)

		// 根据书籍状态决定章节数量
		var chapterCount int
		switch book.Status {
		case "completed":
			chapterCount = 100 + rand.Intn(500) // 完结书 100-600章
		case "ongoing":
			chapterCount = 50 + rand.Intn(200) // 连载中 50-250章
		default:
			chapterCount = 10 + rand.Intn(50) // 其他 10-60章
		}

		for i := 1; i <= chapterCount; i++ {
			chapterID := primitive.NewObjectID()

			// 前10章免费
			isFree := i <= 10

			// 生成章节标题
			title := fmt.Sprintf("第%d章 %s", i, randomChapterTitle())

			// 随机字数 2000-5000
			wordCount := 2000 + rand.Intn(3001)

			// 创建章节
			publishedAt := now.Add(-time.Duration(chapterCount-i) * 24 * time.Hour)
			chapter := Chapter{
				ID:          chapterID,
				BookID:      book.ID,
				ChapterNum:  i,
				Title:       title,
				WordCount:   wordCount,
				Price:       0.05,
				IsFree:      isFree,
				Status:      "published",
				PublishedAt: publishedAt,
				CreatedAt:   publishedAt,
				UpdatedAt:   now,
			}

			_, err := chapterCollection.InsertOne(ctx, chapter)
			if err != nil {
				fmt.Printf("  创建第%d章失败: %v\n", i, err)
				continue
			}
			totalChapters++

			// 创建章节内容（前20章使用示例内容，后续章节使用简化内容）
			var content string
			if i <= 20 {
				// 使用示例内容
				content = sampleChapterContents[rand.Intn(len(sampleChapterContents))]
			} else {
				// 使用简化内容
				content = generateSimpleContent(i, title)
			}

			chapterContent := ChapterContent{
				ID:        primitive.NewObjectID(),
				ChapterID: chapterID.Hex(),
				Content:   content,
				WordCount: len([]rune(content)),
				CreatedAt: publishedAt,
				UpdatedAt: now,
			}

			_, err = contentCollection.InsertOne(ctx, chapterContent)
			if err == nil {
				totalContents++
			}
		}

		fmt.Printf("  创建了 %d 个章节\n", chapterCount)
	}

	fmt.Println("========================================")
	fmt.Println("章节创建完成")
	fmt.Println("========================================")
	fmt.Printf("总计创建: %d 个章节\n", totalChapters)
	fmt.Printf("总计创建: %d 个内容\n", totalContents)
	fmt.Println()

	return nil
}

func randomChapterTitle() string {
	titles := []string{
		"初入江湖",
		"机缘巧合",
		"奇遇",
		"突破",
		"危机",
		"转机",
		"历练",
		"成长",
		"对决",
		"胜利",
		"发现",
		"秘密",
		"传承",
		"蜕变",
		"征程",
		"风云",
		"变故",
		"抉择",
		"战斗",
		"突破",
		"奇遇",
		"机遇",
		"挑战",
		"生死",
		"领悟",
		"进阶",
		"震惊",
		"威名",
		"荣耀",
		"传奇",
		"巅峰",
		"征程",
		"开启",
		"新的开始",
		"风云再起",
		"逆天改命",
		"一鸣惊人",
		"惊天动地",
		"王者归来",
		"绝世无双",
		"傲视群雄",
		"独步天下",
		"登峰造极",
	}
	return titles[rand.Intn(len(titles))]
}

func generateSimpleContent(chapterNum int, title string) string {
	return fmt.Sprintf(`这是第%d章：%s

在这里，故事继续发展。

主角在前面的经历中获得了宝贵的经验，这一次，他面临着新的挑战。

（此处为测试数据，实际内容应该更加丰富精彩。）

随着情节的推进，人物关系逐渐清晰，世界观也在不断展开。

读者在阅读过程中，会感受到作者的用心和故事的魅力。

修真之路，步步艰难，但主角从未放弃。

每一次突破，都伴随着汗水和努力。

每一次成长，都意味着新的责任。

前方还有更长的路要走...

（本章完，字数：%d字）

---
作者有话说：
感谢大家的支持，求推荐票，求收藏！
如果觉得好看，请告诉你的朋友~
`, chapterNum, title, 2000+rand.Intn(3000))
}
