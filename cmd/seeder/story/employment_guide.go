package story

// EmploymentGuideProject 就业指南项目元数据
var EmploymentGuideProject = struct {
	Title        string
	Summary      string
	CoverURL     string
	Category     string
	WritingType  string
	Status       string
	Visibility   string
	Tags         []string
}{
	Title:       "就业指南：我在异世界当劝退专员",
	Summary:     "一个关于异世界中间件服务商的奇幻故事。亚伯是一个战力低下的劝退专员，擅长用话术解决问题；诺艾尔是贫穷的圣女，每天要喝10杯限定版奶茶；伊莎贝拉是对亚伯有妄想型恋爱滤镜的监控者。",
	CoverURL:    "/images/covers/employment-guide.jpg",
	Category:    "奇幻",
	WritingType: "novel",
	Status:      "draft",
	Visibility:  "private",
	Tags:        []string{"奇幻", "异世界", "职场", "搞笑"},
}

// NewCharacter 新角色数据结构
type NewCharacter struct {
	ID               string
	Name             string
	Alias            []string
	Summary          string
	Traits           []string
	Background       string
	AvatarURL        string
	PersonalityPrompt string
	SpeechPattern    string
}

// NewCharacters 三个新角色数据
var NewCharacters = []NewCharacter{
	{
		ID:        "char-abel",
		Name:      "亚伯",
		Alias:     []string{"Abel", "劝退专员"},
		Summary:   "极度理性的利己主义者，战力低下但话术无敌，内心社恐但外表冷静。",
		Traits:    []string{"理性", "话术", "利己", "社恐"},
		Background: "来自异世界的中间件服务商，负责'维持世界平衡'。表面是审计员，实际是劝退专员。",
		AvatarURL: "/images/avatars/abel.png",
		PersonalityPrompt: "极度理性，话术无敌，内心社恐怕死，但外表表现冷静。被误解时外界会脑补为深谋远虑。",
		SpeechPattern: "冷静、理性、一针见血，常用财务和商业术语",
	},
	{
		ID:        "char-noelle",
		Name:      "诺艾尔",
		Alias:     []string{"Noelle", "圣女"},
		Summary:   "高输出低逻辑的笨蛋圣女，金发圣光但极其贫穷，每天要喝10杯限定版奶茶。",
		Traits:    []string{"笨蛋", "高输出", "贫穷", "吃货"},
		Background: "永恒之光勇者小队的圣女，智商波动大，容易被诱导背负债务。",
		AvatarURL: "/images/avatars/noelle.png",
		PersonalityPrompt: "智商波动大，高输出但零逻辑，贫穷吃货，容易被骗但很忠诚。",
		SpeechPattern: "活泼、直接、经常说蠢话",
	},
	{
		ID:        "char-isabella",
		Name:      "伊莎贝拉",
		Alias:     []string{"Isabella", "妄想系"},
		Summary:   "对亚伯有妄想型恋爱滤镜，负责监控亚伯行动并自动脑补完美剧情。",
		Traits:    []string{"妄想", "恋爱脑", "观察力"},
		Background: "亚伯的监控者/追求者，每时每刻都在脑补亚伯行为的深层含义。",
		AvatarURL: "/images/avatars/isabella.png",
		PersonalityPrompt: "妄想型恋爱引擎，自动对亚伯的所有行为进行过度解读，永远脑补出完美剧情。",
		SpeechPattern: "甜蜜、主观、充满浪漫幻想",
	},
}
