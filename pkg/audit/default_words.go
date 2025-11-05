package audit

// GetDefaultSensitiveWords 获取默认敏感词库
// MVP版本包含基础敏感词，实际生产环境应该从数据库或配置文件加载
func GetDefaultSensitiveWords() []SensitiveWordInfo {
	return []SensitiveWordInfo{
		// 政治敏感（高风险）
		{Word: "政治敏感词1", Level: 5, Category: "politics"},
		{Word: "政治敏感词2", Level: 5, Category: "politics"},
		{Word: "政治敏感词3", Level: 4, Category: "politics"},

		// 色情内容（高风险）
		{Word: "色情", Level: 4, Category: "porn"},
		{Word: "黄色", Level: 3, Category: "porn"},
		{Word: "成人", Level: 2, Category: "porn"},

		// 暴力内容（中高风险）
		{Word: "暴力", Level: 3, Category: "violence"},
		{Word: "血腥", Level: 3, Category: "violence"},
		{Word: "杀人", Level: 4, Category: "violence"},
		{Word: "自杀", Level: 4, Category: "violence"},

		// 赌博相关（中风险）
		{Word: "赌博", Level: 3, Category: "gambling"},
		{Word: "赌场", Level: 3, Category: "gambling"},
		{Word: "百家乐", Level: 3, Category: "gambling"},
		{Word: "老虎机", Level: 2, Category: "gambling"},

		// 毒品相关（高风险）
		{Word: "毒品", Level: 5, Category: "drugs"},
		{Word: "大麻", Level: 4, Category: "drugs"},
		{Word: "海洛因", Level: 5, Category: "drugs"},
		{Word: "冰毒", Level: 5, Category: "drugs"},

		// 邪教相关（高风险）
		{Word: "邪教", Level: 5, Category: "cult"},
		{Word: "邪教组织", Level: 5, Category: "cult"},

		// 侮辱谩骂（低中风险）
		{Word: "傻逼", Level: 2, Category: "insult"},
		{Word: "妈的", Level: 2, Category: "insult"},
		{Word: "草泥马", Level: 2, Category: "insult"},
		{Word: "操", Level: 2, Category: "insult"},
		{Word: "fuck", Level: 2, Category: "insult"},

		// 广告推广（低风险）
		{Word: "微信号", Level: 1, Category: "ad"},
		{Word: "加q", Level: 1, Category: "ad"},
		{Word: "加微信", Level: 1, Category: "ad"},
		{Word: "扫码", Level: 1, Category: "ad"},
		{Word: "点击链接", Level: 2, Category: "ad"},

		// 其他常见敏感词
		{Word: "枪支", Level: 3, Category: "violence"},
		{Word: "炸药", Level: 4, Category: "violence"},
		{Word: "恐怖袭击", Level: 5, Category: "violence"},
		{Word: "洗钱", Level: 3, Category: "other"},
		{Word: "诈骗", Level: 3, Category: "other"},
		{Word: "传销", Level: 3, Category: "other"},
	}
}

// GetTestSensitiveWords 获取测试用敏感词库（用于单元测试）
func GetTestSensitiveWords() []SensitiveWordInfo {
	return []SensitiveWordInfo{
		{Word: "测试敏感词", Level: 3, Category: "test"},
		{Word: "test", Level: 2, Category: "test"},
		{Word: "敏感", Level: 3, Category: "test"},
		{Word: "违规", Level: 4, Category: "test"},
		{Word: "色情暴力", Level: 5, Category: "test"},
	}
}

// LoadDefaultWords 加载默认敏感词到DFA过滤器
func LoadDefaultWords(filter *DFAFilter) {
	words := GetDefaultSensitiveWords()
	filter.BatchAddWords(words)
}

// LoadTestWords 加载测试敏感词到DFA过滤器
func LoadTestWords(filter *DFAFilter) {
	words := GetTestSensitiveWords()
	filter.BatchAddWords(words)
}
