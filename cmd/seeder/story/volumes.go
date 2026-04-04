package story

// DemoVolume 卷数据结构
type DemoVolume struct {
	Title    string
	Summary  string
	Order    int
	Chapters []DemoChapter
}

// DemoChapter 章节数据结构
type DemoChapter struct {
	Title       string
	Summary     string
	Order       int
	WordCount   int
	HasOutline3 bool // 是否创建3级大纲
	Content     string
}

// DemoVolumes 所有4卷数据
var DemoVolumes = []DemoVolume{
	{
		Title:   "第一卷：火星觉醒",
		Summary: "考古学家林风在火星遗迹中发现神秘晶体，意外觉醒超能力。地球政府介入，强行带走晶体进行研究。林风与同伴苏文、陈博士被迫逃亡，寻找真相。",
		Order:   1,
		Chapters: []DemoChapter{
			{
				Title:       "火星遗迹",
				Summary:     "林风在火星考古中发现神秘晶体，意外觉醒念力天赋",
				Order:       1,
				WordCount:   3500,
				HasOutline3: true,
			},
			{
				Title:       "能力觉醒",
				Summary:     "晶体激发林风潜能，念力觉醒。苏文首次见证超能力",
				Order:       2,
				WordCount:   3200,
				HasOutline3: true,
			},
			{
				Title:       "政府介入",
				Summary:     "地球政府军方强行介入，没收晶体。陈博士意识到事态严重",
				Order:       3,
				WordCount:   3800,
				HasOutline3: true,
			},
			{
				Title:       "逃亡之路",
				Summary:     "林风等人逃离火星，在逃亡过程中逐渐发现晶体的秘密",
				Order:       4,
				WordCount:   4100,
				HasOutline3: true,
			},
			{
				Title:       "追捕与反追捕",
				Summary:     "军方持续追捕，林风利用念力多次化险为夷",
				Order:       5,
				WordCount:   3600,
				HasOutline3: false,
			},
			{
				Title:       "抉择时刻",
				Summary:     "陈博士发现晶体真相：外星文明遗产。林风决定返回地球寻找答案",
				Order:       6,
				WordCount:   3900,
				HasOutline3: true,
			},
		},
	},
	{
		Title:   "第二卷：独立战争",
		Summary: "火星殖民地宣布独立，地球与火星爆发内战。林风被迫在战争中寻找立场，最终促成双方停火。火星独立，但代价惨重。",
		Order:   2,
		Chapters: []DemoChapter{
			{
				Title:       "火星独立宣言",
				Summary:     "火星总督雷诺兹宣布独立，林风被卷入政治漩涡",
				Order:       1,
				WordCount:   3700,
				HasOutline3: true,
			},
			{
				Title:       "暗流涌动",
				Summary:     "战争前夕的紧张氛围，各方势力暗中角力",
				Order:       2,
				WordCount:   3400,
				HasOutline3: true,
			},
			{
				Title:       "内战爆发",
				Summary:     "地球军队进攻火星殖民地，全面战争爆发",
				Order:       3,
				WordCount:   4200,
				HasOutline3: true,
			},
			{
				Title:       "血腥谈判",
				Summary:     "林风尝试在双方之间调停，但谈判破裂",
				Order:       4,
				WordCount:   3800,
				HasOutline3: true,
			},
			{
				Title:       "真相浮现",
				Summary:     "晶体激活记忆投影：外星文明曾经造访地球",
				Order:       5,
				WordCount:   3500,
				HasOutline3: true,
			},
			{
				Title:       "停火协议",
				Summary:     "林风展示外星威胁证据，促成双方停火。火星获得独立",
				Order:       6,
				WordCount:   4000,
				HasOutline3: true,
			},
		},
	},
	{
		Title:   "第三卷：第一次接触",
		Summary: "外星舰队突然出现，鹰派主张开战，林风努力促成沟通。成功建立联盟，人类与外星文明首次正式接触。",
		Order:   3,
		Chapters: []DemoChapter{
			{
				Title:       "外星舰队",
				Summary:     "庞大外星舰队出现在太阳系边缘，人类陷入恐慌",
				Order:       1,
				WordCount:   3600,
				HasOutline3: true,
			},
			{
				Title:       "恐惧与敌意",
				Summary:     "鹰派将军主张先发制人，联合国安理会紧急会议",
				Order:       2,
				WordCount:   3300,
				HasOutline3: true,
			},
			{
				Title:       "跨越语言的沟通",
				Summary:     "林风通过晶体与外星使者建立心灵连接",
				Order:       3,
				WordCount:   3800,
				HasOutline3: true,
			},
			{
				Title:       "建立信任",
				Summary:     "外星使者说明来意：更强大的敌人正在逼近",
				Order:       4,
				WordCount:   3500,
				HasOutline3: true,
			},
			{
				Title:       "威胁逼近",
				Summary:     "侦察报告显示远方有一支毁灭性舰队",
				Order:       5,
				WordCount:   3400,
				HasOutline3: false,
			},
			{
				Title:       "危机化解",
				Summary:     "林风促成人类-外星联盟，共同备战",
				Order:       6,
				WordCount:   3900,
				HasOutline3: true,
			},
		},
	},
	{
		Title:   "第四卷：联盟纪元",
		Summary: "人类与外星文明结成联盟，共同应对宇宙威胁。刺杀阴谋被揭露，联盟巩固。人类文明进入新纪元。",
		Order:   4,
		Chapters: []DemoChapter{
			{
				Title:       "联盟谈判",
				Summary:     "人类与外星文明正式签署联盟条约",
				Order:       1,
				WordCount:   3700,
				HasOutline3: true,
			},
			{
				Title:       "最后的阻碍",
				Summary:     "鹰派策划刺杀外星使者，企图破坏联盟",
				Order:       2,
				WordCount:   3600,
				HasOutline3: true,
			},
			{
				Title:       "刺杀危机",
				Summary:     "林风识破刺杀阴谋，小琳舍身保护使者",
				Order:       3,
				WordCount:   4000,
				HasOutline3: true,
			},
			{
				Title:       "真相大白",
				Summary:     "记忆体揭示完整历史：人类是外星文明的后裔",
				Order:       4,
				WordCount:   3800,
				HasOutline3: true,
			},
			{
				Title:       "联盟成立",
				Summary:     "人类-外星联盟正式成立，林风担任第一任大使",
				Order:       5,
				WordCount:   3500,
				HasOutline3: false,
			},
			{
				Title:       "新的起点",
				Summary:     "人类文明进入星际时代，林风与苏文在星空下展望未来",
				Order:       6,
				WordCount:   3200,
				HasOutline3: true,
			},
		},
	},
}
