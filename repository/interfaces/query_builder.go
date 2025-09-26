package interfaces

// QueryBuilder 查询构建器接口
type QueryBuilder interface {
	// 条件构建
	Where(field string, operator string, value interface{}) QueryBuilder
	WhereIn(field string, values []interface{}) QueryBuilder
	WhereNotIn(field string, values []interface{}) QueryBuilder
	WhereBetween(field string, start, end interface{}) QueryBuilder
	WhereNull(field string) QueryBuilder
	WhereNotNull(field string) QueryBuilder
	WhereLike(field string, pattern string) QueryBuilder
	WhereRegex(field string, pattern string) QueryBuilder

	// 逻辑操作
	And() QueryBuilder
	Or() QueryBuilder
	Not() QueryBuilder

	// 排序
	OrderBy(field string, direction string) QueryBuilder
	OrderByAsc(field string) QueryBuilder
	OrderByDesc(field string) QueryBuilder

	// 字段选择
	Select(fields ...string) QueryBuilder
	Exclude(fields ...string) QueryBuilder

	// 分页
	Limit(limit int) QueryBuilder
	Skip(skip int) QueryBuilder

	// 聚合
	Count() QueryBuilder
	Sum(field string) QueryBuilder
	Avg(field string) QueryBuilder
	Max(field string) QueryBuilder
	Min(field string) QueryBuilder

	// 构建查询
	Build() (map[string]interface{}, error)
	BuildFilter() (Filter, error)

	// 重置
	Reset() QueryBuilder
	Clone() QueryBuilder
}
