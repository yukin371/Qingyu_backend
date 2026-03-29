package story

// DemoProject 演示项目元数据
var DemoProject = struct {
	Title        string
	Summary      string
	CoverURL     string
	Category     string
	WritingType  string
	Status       string
	Visibility   string
	Tags         []string
}{
	Title: "星际觉醒",
	Summary: "2085年，人类已经在火星建立了殖民地。年轻的考古学家林风在火星遗迹中发现了一枚神秘的晶体，这枚晶体不仅是外星文明留下的遗产，更是开启人类文明新纪元的钥匙。随着真相逐渐浮现，地球与火星的关系迅速恶化，独立战争爆发。在战争最黑暗的时刻，更强大的外星舰队突然出现，人类面临着前所未有的考验。林风和他的同伴们必须在恐惧与敌意中寻找沟通的桥梁，促成人类与外星文明的联盟，共同应对来自宇宙深处的更大威胁。",
	CoverURL:     "/images/covers/interstellar-awakening.jpg",
	Category:     "科幻",
	WritingType:  "novel",
	Status:       "draft",
	Visibility:   "private",
	Tags:         []string{"科幻", "星际", "战争", "外星", "冒险"},
}
