package main

import (
	"fmt"
	"strings"
	"time"

	"Qingyu_backend/cmd/seeder/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type showcaseChapterPlan struct {
	Number     int
	Title      string
	Paragraphs []string
}

type showcaseBannerPlan struct {
	Title       string
	Description string
	Image       string
}

// TODO(showcase): 后续新增第二本、第三本精选书时，只需要补这三块数据：
// 1. showcaseBookSpecs       - 书籍元数据（showcase_books.go）
// 2. showcaseBannerPlans     - 首页/运营位 banner
// 3. showcaseChapterPlans    - 前几章手工正文模板
//
// 推荐补充顺序：
// - 先补 1 本书的 banner + 前 3 章标题 + 前 3 章正文
// - 跑 `go run ./cmd/seeder showcase`
// - 再从首页、榜单、详情页、阅读页回归验证
//
// 当前已完整模板化的作品：
// - 云海问剑录
//
// TODO(showcase-next):
// - 长安雪尽时：补 banner + 前 2-3 章正文
// - 夜航星尘档案：补 banner + 前 2-3 章正文
// - 霓虹停机坪：补前 2 章正文
// - 旧城游戏策展人：补前 2 章正文

var showcaseChapterPlans = map[string][]showcaseChapterPlan{
	"云海问剑录": {
		{Number: 1, Title: "第一章 残卷入海", Paragraphs: []string{
			"海雾压着废弃栈桥向前滚，顾照临把最后一捆旧绳拖进仓棚时，听见木板下传来了一声极轻的金石摩擦。那声音不像潮水拍桩，更像有什么东西在黑水里缓慢翻身。",
			"他掀开被海盐浸透的破板，看见半截乌木匣卡在横梁之间。匣面没有锁孔，只有被岁月磨平的剑纹。指尖刚碰上去，掌心便像被细针扎了一下，凉意顺着经脉一路窜进肩背，让他几乎当场失手。",
			"匣中只有一卷残书，纸页薄如鱼腹，开篇只写着七个字，海上剑墟，照夜归潮。顾照临原本以为那是骗子写给散修看的旧戏本，可当他照着第一页的吐纳图运行一遍，停滞三年的灵息竟第一次有了回应。",
			"他把残卷重新卷起，藏进贴身衣襟里，耳边已传来同伴催工的叫喊。仓棚外的海风裹着潮腥，吹得人睁不开眼，只有胸口那卷书像一枚刚从炉火里夹出来的铁片，灼得人发慌。顾照临不敢再试第二遍，只能一边系绳一边强迫自己把呼吸放稳，免得旁人看出异样。",
			"傍晚收工时，年长的船工老常拎着酒壶在栈桥边坐下，眯眼看他许久，忽然说，小顾，你今天脚步轻了。顾照临只当没听见，俯身去搬最后一只鱼篓。老常咂了口劣酒，语气却比海风更淡，港里近来不太平，夜里听见什么、看见什么，都当不知道，能活得久些。",
			"这句话本是老码头人的自保之道，可落到顾照临耳里，却像专门说给自己听。他回头望了一眼乌沉沉的外海，远处几座废弃灯塔正被雾吞得只剩轮廓。残卷在胸前微微发热，仿佛也在催促他离开这里，去看那片雾后更远的海。",
			"夜深后，他躺在漏雨的旧屋里，把门闩横了两道，才敢再度摊开残卷。纸页上的剑纹在油灯下泛起极淡的银光，每一道呼吸轨迹都写得苛刻而古怪，像是专门给经脉残破的人留的活路。顾照临照着运转到第三轮时，胸口多年未散的闷痛忽然松开一线，仿佛有人在黑暗中推开了一扇尘封很久的门。",
			"他强压住心头狂跳，继续把那一线灵息往下引，直到指尖渗出细密汗珠。窗外潮声一阵高过一阵，屋脊木梁被海风吹得咯吱作响，可那声音之下，另有一道更隐秘的回响，像是极远处有人隔着几百年，轻轻敲了一下剑鞘。顾照临在那一瞬间明白，自己捡到的绝不是一件能悄悄藏起来自保的小机缘。",
		}},
		{Number: 2, Title: "第二章 夜潮剑鸣", Paragraphs: []string{
			"当夜，外海起了逆潮。渔港所有系船木桩都在发颤，像是被同一根看不见的线牵住。顾照临抱着残卷站在堤上，耳边却没有风声，只有一声比一声清晰的剑鸣。",
			"那剑鸣并不从海面传来，而是从他体内震出。残卷上的每一道墨痕，都在识海里化成一道极细的银线，逼着他去看更远的地方。海雾深处仿佛有一整片沉没的山门正缓缓抬头，屋脊上插满了断裂古剑。",
			"他明白，从自己翻开残卷的那一刻起，渔港小工这个身份就已经走到尽头。真正逼近他的，不是机缘，而是一场比海潮更慢、更冷、更久远的旧债。",
			"顾照临原本只是想找个没人的地方把胸口那股躁动压下去，没料到越靠近堤岸，残卷越发滚烫。海面上没有月，远处巡夜船的灯火却被潮雾拉成一线细碎白痕，像是浮在黑水上的一串骨节。逆潮拍岸时，他甚至能听清每一道浪脊里都夹着细细的金属颤音，仿佛海底埋着一整座尚未出鞘的兵库。",
			"他试着收束心神，可识海里的银线却像有了自己的意志，一寸寸往外海延伸。那不是寻常修士感知中的探路灵识，而更像一条被旧誓驱使的归途。银线穿过海雾、掠过沉没礁群，最终落在一片黑沉海域，那里本该什么都没有，却隐约显出山门轮廓，牌坊半塌、石阶尽没，屋脊上密密插着断裂古剑，连风都是冷的。",
			"剑鸣在那片幻象里骤然高了一截，顾照临胸口一震，嘴角溢出一点血丝。可他没有退。三年来他守着渔港讨生活，日日被人说资质断绝、此生难再入道，如今终于有一道门缝愿意朝他打开，即便门后站着的是债，他也想看清那债到底从何而来。",
			"潮头忽然在堤前折开，一道被浪托起的黑影从海里翻出，重重砸在礁石上。那是一把只剩半截的古剑，剑身尽是盐蚀斑驳，仍能看出与残卷上同源的剑纹。顾照临伸手去碰，掌心瞬间被一道寒意割开，血珠落在剑纹上，整片海面随即安静了一瞬，像是某个沉睡已久的名字终于被重新唤醒。",
			"下一刻，巡海司的号角声在港口另一侧陡然吹响。渔港犬吠四起，火把在巷道间迅速亮起。顾照临把半截古剑裹进外衣，转身沿着堤岸向旧仓区奔去。他知道，今晚之后，再也不会有人允许自己像从前那样在海边默默打杂了。",
		}},
		{Number: 3, Title: "第三章 山门旧债", Paragraphs: []string{
			"第二天清晨，巡海司的人就到了码头。他们翻查旧仓、盘问船夫，只为找昨夜海面出现的那一道银色潮线。顾照临低着头搬盐袋，掌心却一直攥着那卷残书，仿佛稍一松手，就会被人听见其中的剑意。",
			"午后，一个背剑老人坐到他常去修补鱼篓的石阶上，开口便问，孩子，你知道照夜两个字怎么写吗。顾照临抬头，对方衣衫褴褛，鞋底还沾着外海黑沙，像是从很远的地方一路走来。",
			"老人没逼他回答，只把一枚生锈铜牌放在石阶上。牌上刻着同样的剑纹，背面只有一句话，山门未灭，弟子当归。顾照临望着那行字，忽然意识到残卷要他偿还的旧债，或许根本不是一卷书那么简单。",
			"码头上的气氛从天亮起就像绷紧的缆绳。巡海司挨船查验，连靠岸补网的老渔户都被赶到一旁盘问。顾照临照旧跟着苦工搬盐、点货、修漏板，动作比平时更稳，心里却一直压着昨夜那半截古剑带来的寒意。只要稍一分神，他就能感觉到怀里的残卷和屋里藏着的断剑彼此呼应，像两枚埋在不同地方的火种，被同一阵风悄悄吹亮。",
			"晌午时分，巡海司统领亲自走上栈桥，命人把昨夜在外海值守的几名船夫押来问话。顾照临从人群缝隙里看见那统领腰间悬着一枚青铜令牌，牌面纹路和残卷上的剑纹只有细微差别，心头不由一沉。若巡海司也在找同一件东西，那这卷残书显然不只是散修秘法，而是牵扯到旧山门和海司秘档的禁物。",
			"他正想着如何把藏在旧屋里的断剑转移出去，午后的太阳忽然被海雾盖住，码头边阴影一寸寸蔓延。背剑老人就是在这时出现的。老人坐下时没有半点声响，像是本就一直坐在石阶上，只是现在才被人看见。顾照临本能地想退开，对方却先笑了一声，说自己饿了半日，能不能分半个冷馒头。",
			"顾照临把身上唯一剩下的干粮递过去。老人接过馒头，只咬了一口，便随口问起照夜两个字怎么写。那语气轻得近乎闲聊，却让顾照临全身寒毛都竖了起来。昨夜残卷第一页的七个字里，正好就有照夜二字。老人见他沉默，也不逼问，只慢慢讲起二十年前外海失火、旧宗门覆灭的传闻，像在说一段无关紧要的码头旧闻。",
			"等到人群散去，老人把那枚生锈铜牌搁在石阶上，指节轻轻敲了两下。铜牌背面那句山门未灭，弟子当归在日光下显得格外刺目。顾照临盯着铜牌良久，终于开口问了一句，若山门还在，为何至今无人归去。老人抬眼看向外海，只说了一句，因为当年该回去的人，很多都死在回家的路上。",
			"这句话像一块石头压进胸口。顾照临忽然明白，自己捡到的残卷不是单纯的传承，而更像一封迟到了很多年的召回令。它选中的不是天资最好的人，而是仍活在海边、仍与旧债有缘的人。老人起身离开前，没有再多说半句，只留下一句今晚子时，若你还想知道自己为什么会被残卷选中，就带着铜牌来旧灯塔。",
		}},
	},
}

var showcaseBannerPlans = map[string]showcaseBannerPlan{
	"云海问剑录": {Title: "海上剑墟开启", Description: "仙侠头部作品，适合首页推荐与阅读演示。", Image: "/images/banners/showcase-yunhai.jpg"},
}

// TODO(showcase-template): 复制下面模板给下一本精选书补数据。
//
// showcaseBannerPlans["作品名"] = showcaseBannerPlan{
// 	Title:       "运营位标题",
// 	Description: "一句话卖点",
// 	Image:       "/images/banners/showcase-your-book.jpg",
// }
//
// showcaseChapterPlans["作品名"] = []showcaseChapterPlan{
// 	{
// 		Number: 1,
// 		Title:  "第一章 标题",
// 		Paragraphs: []string{
// 			"第一段。",
// 			"第二段。",
// 			"第三段。",
// 		},
// 	},
// }

func getShowcaseChapterPlan(bookTitle string, chapterNum int) *showcaseChapterPlan {
	for _, plan := range showcaseChapterPlans[bookTitle] {
		if plan.Number == chapterNum {
			copyPlan := plan
			return &copyPlan
		}
	}
	return nil
}

func formatShowcaseContent(paragraphs []string) string {
	return strings.Join(paragraphs, "\n\n")
}

func showcaseWordCount(paragraphs []string) int {
	total := 0
	for _, paragraph := range paragraphs {
		total += len([]rune(strings.ReplaceAll(paragraph, " ", "")))
	}
	return total
}

func buildShowcaseBanners(books []models.Book) []interface{} {
	now := time.Now()
	banners := make([]interface{}, 0, len(showcaseBannerPlans))
	for _, book := range books {
		plan, ok := showcaseBannerPlans[book.Title]
		if !ok {
			continue
		}
		start := now
		end := now.Add(45 * 24 * time.Hour)
		banners = append(banners, map[string]interface{}{
			"_id":         primitive.NewObjectID(),
			"title":       plan.Title,
			"description": plan.Description,
			"image":       plan.Image,
			"target":      fmt.Sprintf("/bookstore/books/%s", book.ID.Hex()),
			"target_type": "url",
			"sort_order":  10 + len(banners),
			"is_active":   true,
			"start_time":  &start,
			"end_time":    &end,
			"click_count": int64(120 + len(banners)*35),
			"created_at":  now,
			"updated_at":  now,
		})
	}
	return banners
}
