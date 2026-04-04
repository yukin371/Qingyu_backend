package story

// DemoRelation 关系数据结构
type DemoRelation struct {
	FromID   string // 发起者角色ID
	ToID     string // 接收者角色ID
	Type     string // 关系类型
	Strength int    // 关系强度 (1-10)
	Notes    string // 关系说明
}

// DemoRelations 30+条关系数据
var DemoRelations = []DemoRelation{
	// 林风的亲密关系
	{FromID: "char-linfeng", ToID: "char-suwen", Type: "恋人", Strength: 10, Notes: "彼此深爱，相互支持"},
	{FromID: "char-linfeng", ToID: "char-chen", Type: "师徒", Strength: 9, Notes: "导师与学生，深厚的情谊"},
	{FromID: "char-linfeng", ToID: "char-reynold", Type: "盟友", Strength: 7, Notes: "政治盟友，共同推动火星独立"},
	{FromID: "char-linfeng", ToID: "char-ava", Type: "朋友", Strength: 6, Notes: "战友关系，彼此尊重"},
	{FromID: "char-linfeng", ToID: "char-messenger", Type: "朋友", Strength: 9, Notes: "跨越文明的友谊"},
	{FromID: "char-linfeng", ToID: "char-xiaolin", Type: "保护者", Strength: 8, Notes: "林风救下的战争孤儿"},
	
	// 苏文的关系
	{FromID: "char-suwen", ToID: "char-chen", Type: "朋友", Strength: 7, Notes: "学术伙伴"},
	{FromID: "char-suwen", ToID: "char-reynold", Type: "朋友", Strength: 6, Notes: "相互尊重"},
	{FromID: "char-suwen", ToID: "char-messenger", Type: "研究伙伴", Strength: 8, Notes: "共同研究外星生命"},
	
	// 陈博士的关系
	{FromID: "char-chen", ToID: "char-reynold", Type: "朋友", Strength: 6, Notes: "老朋友"},
	{FromID: "char-chen", ToID: "char-zhanggong", Type: "朋友", Strength: 7, Notes: "学术界旧识"},
	
	// 雷诺兹的政治关系
	{FromID: "char-reynold", ToID: "char-ava", Type: "下属", Strength: 8, Notes: "总督与总司令"},
	{FromID: "char-reynold", ToID: "char-general-li", Type: "敌对", Strength: 9, Notes: "战争中的对手"},
	{FromID: "char-reynold", ToID: "char-zhanggong", Type: "谈判对手", Strength: 5, Notes: "停火谈判"},
	{FromID: "char-reynold", ToID: "char-un-sg", Type: "盟友", Strength: 7, Notes: "共同推动和平"},
	
	// 艾娃的关系
	{FromID: "char-ava", ToID: "char-general-li", Type: "敌对", Strength: 10, Notes: "战场上的死敌"},
	{FromID: "char-ava", ToID: "char-linfeng", Type: "战友", Strength: 7, Notes: "并肩作战"},
	{FromID: "char-ava", ToID: "char-xiaolin", Type: "保护者", Strength: 7, Notes: "收养战争孤儿"},
	
	// 李将军的政治关系
	{FromID: "char-general-li", ToID: "char-hawk", Type: "上司", Strength: 9, Notes: "副官"},
	{FromID: "char-general-li", ToID: "char-zhanggong", Type: "政敌", Strength: 7, Notes: "鹰派vs温和派"},
	{FromID: "char-general-li", ToID: "char-un-sg", Type: "竞争", Strength: 6, Notes: "争夺决策权"},
	{FromID: "char-general-li", ToID: "char-messenger", Type: "敌对", Strength: 10, Notes: "极端排外"},
	
	// 张公的外交关系
	{FromID: "char-zhanggong", ToID: "char-un-sg", Type: "盟友", Strength: 8, Notes: "温和派联盟"},
	{FromID: "char-zhanggong", ToID: "char-reynold", Type: "谈判伙伴", Strength: 6, Notes: "停火谈判"},
	
	// 联合国秘书长的关系
	{FromID: "char-un-sg", ToID: "char-general-li", Type: "竞争", Strength: 6, Notes: "文武之争"},
	{FromID: "char-un-sg", ToID: "char-zhanggong", Type: "盟友", Strength: 8, Notes: "温和派联盟"},
	{FromID: "char-un-sg", ToID: "char-messenger", Type: "盟友", Strength: 8, Notes: "支持联盟"},
	{FromID: "char-un-sg", ToID: "char-reynold", Type: "盟友", Strength: 7, Notes: "推动和平"},
	
	// 鹰的关系
	{FromID: "char-hawk", ToID: "char-general-li", Type: "上司", Strength: 10, Notes: "忠诚追随"},
	{FromID: "char-hawk", ToID: "char-messenger", Type: "敌对", Strength: 10, Notes: "刺杀目标"},
	{FromID: "char-hawk", ToID: "char-linfeng", Type: "敌对", Strength: 9, Notes: "多次交手"},
	
	// 小琳的关系
	{FromID: "char-xiaolin", ToID: "char-linfeng", Type: "救命恩人", Strength: 10, Notes: "林风救下的"},
	{FromID: "char-xiaolin", ToID: "char-ava", Type: "被监护人", Strength: 8, Notes: "艾娃收养"},
	{FromID: "char-xiaolin", ToID: "char-messenger", Type: "保护", Strength: 9, Notes: "舍身保护"},
	
	// 外星使者的关系
	{FromID: "char-messenger", ToID: "char-linfeng", Type: "朋友", Strength: 9, Notes: "心灵连接"},
	{FromID: "char-messenger", ToID: "char-suwen", Type: "研究伙伴", Strength: 8, Notes: "生命研究"},
	{FromID: "char-messenger", ToID: "char-un-sg", Type: "盟友", Strength: 8, Notes: "联盟伙伴"},
	
	// 记忆体的关系
	{FromID: "char-memory", ToID: "char-linfeng", Type: "引导者", Strength: 10, Notes: "揭示历史真相"},
	{FromID: "char-memory", ToID: "char-messenger", Type: "创造者与创造物", Strength: 10, Notes: "外星AI"},
}
