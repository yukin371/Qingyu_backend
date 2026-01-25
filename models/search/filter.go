package search

// Filter 过滤器
type Filter struct {
	Terms   map[string][]string // 精确匹配（或关系）
	Range   map[string]Range    // 范围过滤
	Exists  []string            // 字段存在
	Not     map[string][]string // 不匹配
}

// Range 范围
type Range struct {
	GTE interface{} // 大于等于
	GT  interface{} // 大于
	LTE interface{} // 小于等于
	LT  interface{} // 小于
}

// BoolFilter 布尔过滤器
type BoolFilter struct {
	Must     []Filter // 必须匹配
	Should   []Filter // 应该匹配
	MustNot  []Filter // 必须不匹配
	Filter   []Filter // 过滤（不计分）
}

// BookFilter 书籍过滤器
type BookFilter struct {
	CategoryID   string   // 分类 ID
	Author       string   // 作者
	Tags         []string // 标签
	Status       []string // 状态
	WordCountMin int      // 最小字数
	WordCountMax int      // 最大字数
	RatingMin    float64  // 最小评分
	IsPrivate    bool     // 是否私密
}

// ProjectFilter 项目过滤器
type ProjectFilter struct {
	AuthorID string   // 作者 ID（强制过滤）
	Status   []string // 状态
	Genre    string   // 类型
	Tags     []string // 标签
}

// DocumentFilter 文档过滤器
type DocumentFilter struct {
	UserID    string // 用户 ID（强制过滤）
	ProjectID string // 项目 ID
	Type      string // 文档类型
	Status    string // 状态
}

// UserFilter 用户过滤器
type UserFilter struct {
	Role       string // 角色
	IsVerified bool   // 是否认证
	Status     string // 状态
}
