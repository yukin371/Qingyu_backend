package story

// DemoItem 道具数据结构
type DemoItem struct {
	Name        string
	Type        string
	Description string
	OwnerID     string // 拥有者角色ID
	LocationID  string // 所在地点ID
	Rarity      string // 稀有度
	Function    string // 功能
	Origin      string // 来源
}

// DemoLocation 地点数据结构
type DemoLocation struct {
	Name        string
	Description string
	Climate     string
	Culture     string
	Geography   string
	Atmosphere  string
}

// DemoTimeline 时间线数据结构
type DemoTimeline struct {
	Name        string
	Description string
	Events      []DemoTimelineEvent
}

// DemoTimelineEvent 时间线事件
type DemoTimelineEvent struct {
	Title       string
	Description string
	Year        int
	Month       int
	Day         int
	Importance  int // 重要性 1-10
}

// DemoItems 8个道具数据
var DemoItems = []DemoItem{
	{
		Name:        "先民晶体",
		Type:        "遗物",
		Description: "外星文明留下的神秘晶体，能够激发人类潜能，储存着重要信息",
		OwnerID:     "char-linfeng",
		LocationID:  "",
		Rarity:      "传说",
		Function:    "激发潜能，储存记忆，心灵连接",
		Origin:      "火星遗迹",
	},
	{
		Name:        "记忆体",
		Type:        "AI设备",
		Description: "外星文明创造的AI，储存着完整的历史信息",
		OwnerID:     "",
		LocationID:  "loc-mars-ruins",
		Rarity:      "传说",
		Function:    "储存信息，回答问题",
		Origin:      "外星文明",
	},
	{
		Name:        "念力增幅器",
		Type:        "科技装备",
		Description: "地球军方研发的装置，可以增幅念力能力",
		OwnerID:     "char-general-li",
		LocationID:  "",
		Rarity:      "稀有",
		Function:    "增幅念力",
		Origin:      "地球科技",
	},
	{
		Name:        "心灵翻译器",
		Type:        "科技装备",
		Description: "苏文研发的装置，辅助跨物种沟通",
		OwnerID:     "char-suwen",
		LocationID:  "",
		Rarity:      "稀有",
		Function:    "辅助沟通",
		Origin:      "地球科技",
	},
	{
		Name:        "能量护盾",
		Type:        "防御装备",
		Description: "外星科技，可以抵挡能量攻击",
		OwnerID:     "char-ava",
		LocationID:  "",
		Rarity:      "史诗",
		Function:    "能量防护",
		Origin:      "外星科技",
	},
	{
		Name:        "考古扫描仪",
		Type:        "科研设备",
		Description: "陈博士的专业设备，可以分析古代遗迹",
		OwnerID:     "char-chen",
		LocationID:  "",
		Rarity:      "普通",
		Function:    "考古分析",
		Origin:      "地球科技",
	},
	{
		Name:        "星际通讯器",
		Type:        "通讯设备",
		Description: "外星文明的高级通讯设备",
		OwnerID:     "char-messenger",
		LocationID:  "",
		Rarity:      "史诗",
		Function:    "跨星际通讯",
		Origin:      "外星科技",
	},
	{
		Name:        "纳米医疗包",
		Type:        "医疗设备",
		Description: "可以快速治疗伤势的纳米机器人",
		OwnerID:     "char-xiaolin",
		LocationID:  "",
		Rarity:      "稀有",
		Function:    "快速治疗",
		Origin:      "火星科技",
	},
}

// DemoLocations 6个地点数据
var DemoLocations = []DemoLocation{
	{
		Name:        "火星遗迹",
		Description: "火星地表的古代遗迹，外星文明曾经造访的地方",
		Climate:     "干燥寒冷",
		Culture:     "考古遗址",
		Geography:   "火星地表沙漠",
		Atmosphere:  "神秘、庄严",
	},
	{
		Name:        "火星殖民地",
		Description: "人类在火星建立的第一个大型殖民地",
		Climate:     "人造气候",
		Culture:     "地球-火星混合文化",
		Geography:   "火星地下城市",
		Atmosphere:  "繁忙、希望",
	},
	{
		Name:        "地球联合国总部",
		Description: "人类政治中心，位于纽约",
		Climate:     "温带",
		Culture:     "地球文化",
		Geography:   "地球城市",
		Atmosphere:  "庄严、紧张",
	},
	{
		Name:        "太空舰队驻地",
		Description: "地球军方太空舰队基地",
		Climate:     "真空",
		Culture:     "军事文化",
		Geography:   "太空站",
		Atmosphere:  "肃杀、紧张",
	},
	{
		Name:        "外星母舰",
		Description: "外星文明的巨型母舰，停泊在太阳系边缘",
		Climate:     "人造气候",
		Culture:     "外星文明",
		Geography:   "太空",
		Atmosphere:  "神秘、先进",
	},
	{
		Name:        "谈判大厅",
		Description: "人类与外星文明首次正式接触的地点",
		Climate:     "恒温",
		Culture:     "中立区",
		Geography:   "太空站",
		Atmosphere:  "历史性、紧张",
	},
}

// DemoTimelines 3条时间线数据
var DemoTimelines = []DemoTimeline{
	{
		Name:        "主线故事",
		Description: "《星际觉醒》主要剧情时间线",
		Events: []DemoTimelineEvent{
			{
				Title:       "火星遗迹发现",
				Description: "林风在火星遗迹中发现先民晶体",
				Year:        2085,
				Month:      3,
				Day:        15,
				Importance:  10,
			},
			{
				Title:       "能力觉醒",
				Description: "晶体激发林风念力能力",
				Year:        2085,
				Month:      3,
				Day:        16,
				Importance:  9,
			},
			{
				Title:       "政府介入",
				Description: "地球政府军方强行介入",
				Year:        2085,
				Month:      3,
				Day:        20,
				Importance:  8,
			},
			{
				Title:       "火星独立宣言",
				Description: "火星总督雷诺兹宣布独立",
				Year:        2085,
				Month:      5,
				Day:        1,
				Importance:  10,
			},
			{
				Title:       "内战爆发",
				Description: "地球与火星爆发全面战争",
				Year:        2085,
				Month:      5,
				Day:        15,
				Importance:  9,
			},
			{
				Title:       "停火协议",
				Description: "林风促成双方停火",
				Year:        2085,
				Month:      8,
				Day:        20,
				Importance:  10,
			},
			{
				Title:       "外星舰队出现",
				Description: "外星舰队出现在太阳系边缘",
				Year:        2085,
				Month:      9,
				Day:        1,
				Importance:  10,
			},
			{
				Title:       "首次接触",
				Description: "林风与外星使者建立心灵连接",
				Year:        2085,
				Month:      9,
				Day:        10,
				Importance:  10,
			},
			{
				Title:       "联盟成立",
				Description: "人类与外星文明正式结盟",
				Year:        2085,
				Month:      10,
				Day:        1,
				Importance:  10,
			},
		},
	},
	{
		Name:        "外星文明历史",
		Description: "外星文明从兴起到现在的历史",
		Events: []DemoTimelineEvent{
			{
				Title:       "文明诞生",
				Description: "外星文明在半人马座阿尔星诞生",
				Year:        -8000,
				Month:      1,
				Day:        1,
				Importance:  10,
			},
			{
				Title:       "星际扩张",
				Description: "开始向其他星系扩张",
				Year:        -5000,
				Month:      1,
				Day:        1,
				Importance:  8,
			},
			{
				Title:       "造访地球",
				Description: "第一次造访地球，留下人类祖先",
				Year:        -10000,
				Month:      1,
				Day:        1,
				Importance:  10,
			},
			{
				Title:       "再次造访",
				Description: "第二次造访地球，留下晶体和记忆体",
				Year:        -3000,
				Month:      1,
				Day:        1,
				Importance:  9,
			},
			{
				Title:       "威胁出现",
				Description: "发现更强大的敌对文明",
				Year:        -500,
				Month:      1,
				Day:        1,
				Importance:  9,
			},
			{
				Title:       "寻找盟友",
				Description: "向各星系派遣使者，寻找盟友",
				Year:        2080,
				Month:      1,
				Day:        1,
				Importance:  8,
			},
		},
	},
	{
		Name:        "人类文明发展",
		Description: "人类从古至今的关键节点",
		Events: []DemoTimelineEvent{
			{
				Title:       "火星殖民开始",
				Description: "人类开始在火星建立殖民地",
				Year:        2050,
				Month:      1,
				Day:        1,
				Importance:  8,
			},
			{
				Title:       "第一次火星考古发现",
				Description: "在火星发现古代遗迹",
				Year:        2075,
				Month:      1,
				Day:        1,
				Importance:  7,
			},
			{
				Title:       "星际觉醒",
				Description: "林风发现晶体，人类进入新纪元",
				Year:        2085,
				Month:      3,
				Day:        15,
				Importance:  10,
			},
			{
				Title:       "火星独立",
				Description: "火星获得独立地位",
				Year:        2085,
				Month:      8,
				Day:        20,
				Importance:  9,
			},
			{
				Title:       "第一次接触",
				Description: "人类与外星文明首次正式接触",
				Year:        2085,
				Month:      9,
				Day:        10,
				Importance:  10,
			},
			{
				Title:       "星际联盟成立",
				Description: "人类加入星际联盟",
				Year:        2085,
				Month:      10,
				Day:        1,
				Importance:  10,
			},
			{
				Title:       "星际时代开启",
				Description: "人类文明进入星际时代",
				Year:        2090,
				Month:      1,
				Day:        1,
				Importance:  10,
			},
		},
	},
}
