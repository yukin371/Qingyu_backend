//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("qingyu")

	// 获取测试书籍ID（修仙世界）
	bookID, _ := primitive.ObjectIDFromHex("696f35c4cee9d6ed15e66935")

	// 检查是否已有章节数据
	chapterCount, _ := db.Collection("chapters").CountDocuments(ctx, bson.M{"book_id": bookID})
	if chapterCount > 0 {
		fmt.Printf("书籍已有 %d 个章节，跳过创建\n", chapterCount)
		return
	}

	fmt.Println("开始创建测试章节数据...")
	fmt.Println("==========================================")

	now := time.Now()
	contentCount := len(sampleContents)

	// 创建5个测试章节
	for i := 1; i <= 5; i++ {
		chapterID := primitive.NewObjectID()
		chapterNum := i
		title := fmt.Sprintf("第%d章 %s", i, getChapterTitle(i))
		isFree := i <= 3 // 前3章免费
		wordCount := 2000 + i*200

		// 选择内容
		contentIndex := (i - 1) % contentCount
		content := sampleContents[contentIndex]

		// 创建章节元数据
		chapter := bson.M{
			"_id":          chapterID,
			"book_id":      bookID,
			"chapter_num":  chapterNum,
			"title":        title,
			"word_count":   wordCount,
			"price":        0.05,
			"is_free":      isFree,
			"publish_time": now.Add(-time.Duration(i) * 24 * time.Hour),
			"created_at":   now,
			"updated_at":   now,
		}

		_, err := db.Collection("chapters").InsertOne(ctx, chapter)
		if err != nil {
			log.Printf("创建第%d章失败: %v", i, err)
			continue
		}
		fmt.Printf("✓ 创建第%d章: %s\n", i, title)

		// 创建章节内容
		chapterContent := bson.M{
			"_id":        primitive.NewObjectID(),
			"chapter_id": chapterID,
			"content":    content,
			"format":     "markdown",
			"version":    1,
			"word_count": len([]rune(content)),
			"created_at": now,
			"updated_at": now,
		}

		_, err = db.Collection("chapter_contents").InsertOne(ctx, chapterContent)
		if err != nil {
			log.Printf("创建第%d章内容失败: %v", i, err)
		} else {
			fmt.Printf("  ✓ 添加内容 (%d字)\n", wordCount)
		}
	}

	fmt.Println("==========================================")
	fmt.Println("✓ 测试章节数据创建完成！")
	fmt.Printf("总计: 5 个章节\n")
}

func getChapterTitle(i int) string {
	titles := []string{
		"初入修真界",
		"灵根测试",
		"拜入宗门",
		"基础功法",
		"首次突破",
	}
	if i-1 < len(titles) {
		return titles[i-1]
	}
	return fmt.Sprintf("新的开始（第%d章）", i)
}

// 示例章节内容
var sampleContents = []string{
	`# 第一章 初入修真界

晨光熹微，透过破旧的窗棂，洒落在少年苍白的脸上。

李凡缓缓睁开眼睛，看着熟悉的茅草屋顶，心中五味杂陈。这是他来到这个世界的第三个月。

前世，他是一个普通的上班族，每天朝九晚五，为了生活奔波。一场意外，让他穿越到了这个名为"苍澜大陆"的修真世界。

"既然老天让我重活一次，我定要活出个模样来！"李凡握紧拳头，眼中闪过一丝坚定。

这是一个弱肉强食的世界，只有拥有力量，才能掌控自己的命运。

---

## 测试内容

这是一段测试章节内容，用于验证E2E测试的章节阅读功能。

*主要功能点：*
- 章节内容显示
- 字体大小调整
- 阅读主题切换
- 章节导航（上一章/下一章）

测试数据准备完成，可以开始进行章节阅读功能的E2E测试验证。`,

	`# 第二章 灵根测试

天空中飘着细雨，苍澜宗的山门笼罩在烟雨朦胧之中。

李凡站在山脚下，抬头望去。只见巍峨的山峰直插云霄，云雾缭绕间，隐约可见亭台楼阁，宛如仙境。

"这便是修真门派吗？"他喃喃自语，心中既向往又忐忑。

周围聚集了数千名少年，个个神情激动。有人穿着锦衣华服，显然出身富贵之家；也有人如李凡一般，衣衫朴素，眼中却透着坚毅。

---

## 测试灵根

**测试开始！**

随着一声高喝，一名身着青衫的中年人出现在众人面前。他便是负责此次招生的执事长老。

"第一关，测骨龄。十二到十六岁者，方可通过。"

一块巨大的测灵石被抬了上来，散发着淡淡的光芒。

李凡深吸一口气，走上前去。将手放在测灵石上，一道光芒闪过。

"骨龄十五岁，合格！"

听到这句话，李凡心中一松，至少第一关过了。`,

	`# 第三章 拜入宗门

"第二关，测灵根。"

执事长老的声音传来，让李凡从回忆中回过神来。

**灵根**，是修行的根本。没有灵根，便无法感应天地灵气，更别提修炼了。

传说中，灵根分为天地玄黄四等：
- 天灵根，万中无一
- 地灵根，千里挑一
- 玄灵根，百中有一
- 黄灵根，最为常见

---

## 测试进行中

李凡忐忑地将手放在测灵石上。

时间一分一秒过去，测灵石没有任何反应。

一息，两息，三息...

"下一个。"

执事长老淡漠的声音，如同晴天霹雳，让李凡如坠冰窟。

**没有灵根！**`,

	`# 第四章 基础功法

就在李凡即将离去的时候，异变突生。

他怀中那块从小佩戴的黑色玉佩，突然发出一缕微弱的光芒。

这光芒虽然微弱，但在场的修士都感应到了。

"嗯？"

执事长老眼睛一亮，身形一闪，便出现在李凡面前。

"小子，你怀中之物，拿出来看看。"

李凡不解，但还是将玉佩取出。那是一块古朴的黑玉，上面刻着复杂的纹路。

---

## 混沌玉

"这是...**混沌玉**？！"

他的声音都有些颤抖，眼中满是不可置信。

**混沌玉**，传说中的至宝，可吞噬万物，孕育混沌之气。拥有此物者，即便没有灵根，亦可踏上修真之路。

"你，被苍澜宗录取了。"

执事长老深吸一口气，郑重地说道。

修真之路，就在这一刻，为他敞开了大门。`,

	`# 第五章 首次突破

进入苍澜宗已经三个月了。

李凡盘坐在洞府中，感受着体内缓缓流动的灵气。虽然他没有灵根，但凭借混沌玉，他也能像普通修士一样修炼。

**修真境界**：
- 炼气期
- 筑基期
- 金丹期
- 元婴期
- 化神期

目前，他还处于炼气期一层。

---

## 功法运行

李凡按照师父传授的《引气诀》，开始运转周天。

混沌玉在他胸前微微发光，将周围的天地灵气吸收进来，经过提纯，送入他的经脉。

一个周天，两个周天，三个周天...

随着灵气在体内不断循环，李凡感觉到自己的修为在缓慢提升。

"这就是修真吗..."

他睁开眼睛，眼中满是震撼。

这就是修真的世界，一个充满奇迹和挑战的世界。而他的故事，才刚刚开始。`,
}
