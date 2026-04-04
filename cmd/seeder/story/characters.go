package story

// DemoCharacter 角色数据结构
type DemoCharacter struct {
	ID               string
	Name             string
	Alias            string
	Summary          string
	Traits           []string
	Background       string
	AvatarURL        string
	ShortDescription string
}

// DemoCharacters 12个角色数据
var DemoCharacters = []DemoCharacter{
	{
		ID:        "char-linfeng",
		Name:      "林风",
		Alias:     "小林",
		Summary:   "年轻的考古学家，在火星遗迹中意外获得外星晶体，觉醒念力能力。性格坚韧、善良，在危机中成长为文明的桥梁。",
		Traits:    []string{"勇敢", "善良", "聪慧", "领导力"},
		Background: "2055年出生于地球，2080年获得考古学博士学位，2082年前往火星参与考古项目。",
		AvatarURL: "/images/avatars/linfeng.png",
		ShortDescription: "主角，考古学家，念力觉醒者",
	},
	{
		ID:        "char-suwen",
		Name:      "苏文",
		Alias:     "苏苏",
		Summary:   "生物学家，林风的恋人。理智冷静，是团队中的智囊。在逃亡过程中始终支持林风，最终成为外星生命研究专家。",
		Traits:    []string{"理智", "忠诚", "敏锐", "温柔"},
		Background: "2056年出生于地球，2081年获得生物学博士学位，同年前往火星。",
		AvatarURL: "/images/avatars/suwen.png",
		ShortDescription: "女主角，生物学家，林风的恋人",
	},
	{
		ID:        "char-chen",
		Name:      "陈博士",
		Alias:     "老陈",
		Summary:   "资深考古学家，林风的导师。德高望重，在学术界享有崇高声誉。第一个意识到晶体重要性的人。",
		Traits:    []string{"博学", "谨慎", "正直", "慈祥"},
		Background: "2040年出生于地球，从事考古研究40年，火星考古项目总负责人。",
		AvatarURL: "/images/avatars/chen.png",
		ShortDescription: "导师，考古学家，德高望重",
	},
	{
		ID:        "char-reynold",
		Name:      "雷诺兹",
		Alias:     "总督",
		Summary:   "火星殖民地总督，政治家。务实、果断，在关键时刻做出独立决定。战后成为火星第一任总统。",
		Traits:    []string{"果断", "务实", "远见", "威严"},
		Background: "2048年出生于地球，2075年从政，2080年出任火星殖民地总督。",
		AvatarURL: "/images/avatars/reynold.png",
		ShortDescription: "火星总督，政治家，独立运动领袖",
	},
	{
		ID:        "char-ava",
		Name:      "艾娃",
		Alias:     "女将军",
		Summary:   "火星自卫军总司令，女将军。英勇善战，忠于火星，但对地球同胞怀有复杂情感。",
		Traits:    []string{"勇敢", "忠诚", "刚毅", "矛盾"},
		Background: "2052年出生于地球，2070年参军，2078年调任火星殖民地。",
		AvatarURL: "/images/avatars/ava.png",
		ShortDescription: "火星自卫军总司令，女将军",
	},
	{
		ID:        "char-general-li",
		Name:      "李将军",
		Alias:     "鹰派",
		Summary:   "地球军方鹰派领袖，强硬、排外。主张对外星文明采取敌对态度，多次策划破坏联盟。",
		Traits:    []string{"强硬", "排外", "野心", "狡诈"},
		Background: "2045年出生于地球，军旅生涯30年，地球武装力量总司令。",
		AvatarURL: "/images/avatars/general-li.png",
		ShortDescription: "地球鹰派将军，强硬排外",
	},
	{
		ID:        "char-zhanggong",
		Name:      "张公",
		Alias:     "外长",
		Summary:   "地球外交部长，温和派。主张和平解决争端，是林风在地球政府中的重要盟友。",
		Traits:    []string{"温和", "圆滑", "智慧", "耐心"},
		Background: "2042年出生于地球，外交生涯35年，以善于斡旋著称。",
		AvatarURL: "/images/avatars/zhanggong.png",
		ShortDescription: "地球外长，温和派政治家",
	},
	{
		ID:        "char-messenger",
		Name:      "信使",
		Alias:     "外星使者",
		Summary:   "外星文明使者，拥有高度发达的心灵感应能力。温和、智慧，致力于促进文明间理解。",
		Traits:    []string{"智慧", "温和", "神秘", "博爱"},
		Background: "来自半人马座阿尔星文明，拥有超过地球文明10000年的历史。",
		AvatarURL: "/images/avatars/messenger.png",
		ShortDescription: "外星文明使者，心灵感应者",
	},
	{
		ID:        "char-hawk",
		Name:      "鹰",
		Alias:     "鹰派军官",
		Summary:   "李将军的副官，极端鹰派。多次参与刺杀行动，最终被制止。",
		Traits:    []string{"极端", "忠诚", "冷酷", "狂热"},
		Background: "2058年出生于地球，从小接受军事化教育，李将军的忠实追随者。",
		AvatarURL: "/images/avatars/hawk.png",
		ShortDescription: "极端鹰派军官，刺杀行动执行者",
	},
	{
		ID:        "char-un-sg",
		Name:      "联合国秘书长",
		Alias:     "秘书长",
		Summary:   "联合国秘书长，政治家。在战争与和平之间艰难平衡，最终支持联盟。",
		Traits:    []string{"谨慎", "中立", "理性", "负责"},
		Background: "2046年出生于欧洲，外交官出身，2083年当选联合国秘书长。",
		AvatarURL: "/images/avatars/un-sg.png",
		ShortDescription: "联合国秘书长，中立派领袖",
	},
	{
		ID:        "char-xiaolin",
		Name:      "小琳",
		Alias:     "少女",
		Summary:   "年轻女孩，战争孤儿。被林风救下后成为追随者，在刺杀危机中舍身保护外星使者。",
		Traits:    []string{"纯真", "勇敢", "感恩", "无私"},
		Background: "2070年出生于火星，父母在独立战争中牺牲，后被林风救下。",
		AvatarURL: "/images/avatars/xiaolin.png",
		ShortDescription: "战争孤儿，勇敢的少女",
	},
	{
		ID:        "char-memory",
		Name:      "记忆体",
		Alias:     "AI",
		Summary:   "外星文明留下的AI记忆体，储存着完整的历史信息。帮助林风了解真相。",
		Traits:    []string{"智慧", "理性", "客观", "古老"},
		Background: "由外星文明在10000年前创造，一直在晶体中沉睡。",
		AvatarURL: "/images/avatars/memory.png",
		ShortDescription: "外星AI记忆体，历史记录者",
	},
}
