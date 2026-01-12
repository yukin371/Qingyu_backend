package querybuilder

import (
	"Qingyu_backend/repository/interfaces/infrastructure"
	"fmt"
	"regexp"
	"strings"
)

// conditionGroup 条件组，支持嵌套逻辑
type conditionGroup struct {
	logic      string              // "and", "or"
	conditions []map[string]interface{}
}

// MongoQueryBuilder MongoDB查询构建器实现
type MongoQueryBuilder struct {
	groups        []conditionGroup     // 条件分组，支持嵌套逻辑
	currentGroup  *conditionGroup      // 当前激活的条件组
	sortFields    map[string]int
	selectFields  []string
	excludeFields []string
	limitValue    int
	skipValue     int
	aggregations  []map[string]interface{}
}

// NewMongoQueryBuilder 创建MongoDB查询构建器
func NewMongoQueryBuilder() infrastructure.QueryBuilder {
	groups := []conditionGroup{{
		logic:      "and",
		conditions: make([]map[string]interface{}, 0),
	}}
	return &MongoQueryBuilder{
		groups:        groups,
		currentGroup:  &groups[0],
		sortFields:    make(map[string]int),
		selectFields:  make([]string, 0),
		excludeFields: make([]string, 0),
		aggregations:  make([]map[string]interface{}, 0),
	}
}

// Where 添加条件
func (qb *MongoQueryBuilder) Where(field string, operator string, value interface{}) infrastructure.QueryBuilder {
	condition := qb.buildCondition(field, operator, value)
	if condition != nil {
		qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	}
	return qb
}

// WhereIn 添加IN条件
func (qb *MongoQueryBuilder) WhereIn(field string, values []interface{}) infrastructure.QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	condition := map[string]interface{}{
		field: map[string]interface{}{"$in": values},
	}
	qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	return qb
}

// WhereNotIn 添加NOT IN条件
func (qb *MongoQueryBuilder) WhereNotIn(field string, values []interface{}) infrastructure.QueryBuilder {
	if len(values) == 0 {
		return qb
	}
	condition := map[string]interface{}{
		field: map[string]interface{}{"$nin": values},
	}
	qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	return qb
}

// WhereBetween 添加范围条件
func (qb *MongoQueryBuilder) WhereBetween(field string, start, end interface{}) infrastructure.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{
			"$gte": start,
			"$lte": end,
		},
	}
	qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	return qb
}

// WhereNull 添加NULL条件
func (qb *MongoQueryBuilder) WhereNull(field string) infrastructure.QueryBuilder {
	condition := map[string]interface{}{
		field: nil,
	}
	qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	return qb
}

// WhereNotNull 添加NOT NULL条件
func (qb *MongoQueryBuilder) WhereNotNull(field string) infrastructure.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{"$ne": nil},
	}
	qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	return qb
}

// WhereLike 添加LIKE条件（模糊匹配）
// 支持的通配符：
//   % - 匹配任意数量的字符
//   _ - 匹配单个字符
// 转义字符：使用 \ 可以转义通配符，例如 \% 匹配字面量 %
func (qb *MongoQueryBuilder) WhereLike(field string, pattern string) infrastructure.QueryBuilder {
	regexPattern := qb.convertLikeToRegex(pattern)
	condition := map[string]interface{}{
		field: map[string]interface{}{
			"$regex":   regexPattern,
			"$options": "i", // 不区分大小写
		},
	}
	qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	return qb
}

// convertLikeToRegex 将 SQL LIKE 模式转换为 MongoDB 正则表达式
// 支持转义字符，例如:
//   "abc%"     -> "^abc.*$"
//   "abc_def"  -> "^abc.def.$"
//   "abc\%def" -> "^abc%def$"
func (qb *MongoQueryBuilder) convertLikeToRegex(pattern string) string {
	var sb strings.Builder
	sb.WriteByte('^')

	escaped := false
	for i, r := range pattern {
		if escaped {
			// 前一个字符是转义符，当前字符作为字面量
			sb.WriteString(regexp.QuoteMeta(string(r)))
			escaped = false
			continue
		}

		switch r {
		case '\\':
			// 转义符
			if i < len(pattern)-1 {
				escaped = true
			} else {
				// 末尾的反斜杠，作为字面量
				sb.WriteString(`\\`)
			}
		case '%':
			// 匹配任意数量的字符
			sb.WriteString(".*")
		case '_':
			// 匹配单个字符
			sb.WriteString(".")
		default:
			// 其他字符进行转义，避免正则表达式特殊字符问题
			sb.WriteString(regexp.QuoteMeta(string(r)))
		}
	}

	sb.WriteByte('$')
	return sb.String()
}

// WhereRegex 添加正则表达式条件
func (qb *MongoQueryBuilder) WhereRegex(field string, pattern string) infrastructure.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{
			"$regex":   pattern,
			"$options": "i",
		},
	}
	qb.currentGroup.conditions = append(qb.currentGroup.conditions, condition)
	return qb
}

// And 开始一个新的 AND 逻辑组
// 后续添加的条件将使用 AND 逻辑与之前的条件组合
func (qb *MongoQueryBuilder) And() infrastructure.QueryBuilder {
	// 如果当前组已经有条件，创建一个新的 OR 组
	if len(qb.currentGroup.conditions) > 0 {
		newGroup := conditionGroup{
			logic:      "and",
			conditions: make([]map[string]interface{}, 0),
		}
		qb.groups = append(qb.groups, newGroup)
		qb.currentGroup = &qb.groups[len(qb.groups)-1]
	}
	return qb
}

// Or 开始一个新的 OR 逻辑组
// 后续添加的条件将使用 OR 逻辑与之前的条件组合
func (qb *MongoQueryBuilder) Or() infrastructure.QueryBuilder {
	// 如果当前组已经有条件，创建一个新的 OR 组
	if len(qb.currentGroup.conditions) > 0 {
		newGroup := conditionGroup{
			logic:      "or",
			conditions: make([]map[string]interface{}, 0),
		}
		qb.groups = append(qb.groups, newGroup)
		qb.currentGroup = &qb.groups[len(qb.groups)-1]
	}
	return qb
}

// Not 对当前条件组应用 NOT 逻辑
// 注意：这会影响当前组，而不是后续添加的条件
func (qb *MongoQueryBuilder) Not() infrastructure.QueryBuilder {
	// 标记当前组为 NOT 逻辑
	if len(qb.groups) > 0 {
		qb.groups[len(qb.groups)-1].logic = "not"
	}
	return qb
}

// OrderBy 添加排序
func (qb *MongoQueryBuilder) OrderBy(field string, direction string) infrastructure.QueryBuilder {
	dir := 1
	if strings.ToLower(direction) == "desc" {
		dir = -1
	}
	qb.sortFields[field] = dir
	return qb
}

// OrderByAsc 升序排序
func (qb *MongoQueryBuilder) OrderByAsc(field string) infrastructure.QueryBuilder {
	qb.sortFields[field] = 1
	return qb
}

// OrderByDesc 降序排序
func (qb *MongoQueryBuilder) OrderByDesc(field string) infrastructure.QueryBuilder {
	qb.sortFields[field] = -1
	return qb
}

// Select 选择字段
func (qb *MongoQueryBuilder) Select(fields ...string) infrastructure.QueryBuilder {
	qb.selectFields = append(qb.selectFields, fields...)
	return qb
}

// Exclude 排除字段
func (qb *MongoQueryBuilder) Exclude(fields ...string) infrastructure.QueryBuilder {
	qb.excludeFields = append(qb.excludeFields, fields...)
	return qb
}

// Limit 限制数量
func (qb *MongoQueryBuilder) Limit(limit int) infrastructure.QueryBuilder {
	qb.limitValue = limit
	return qb
}

// Skip 跳过数量
func (qb *MongoQueryBuilder) Skip(skip int) infrastructure.QueryBuilder {
	qb.skipValue = skip
	return qb
}

// Count 计数聚合
func (qb *MongoQueryBuilder) Count() infrastructure.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$count": "total",
	})
	return qb
}

// Sum 求和聚合
func (qb *MongoQueryBuilder) Sum(field string) infrastructure.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":   nil,
			"total": map[string]interface{}{"$sum": "$" + field},
		},
	})
	return qb
}

// Avg 平均值聚合
func (qb *MongoQueryBuilder) Avg(field string) infrastructure.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":     nil,
			"average": map[string]interface{}{"$avg": "$" + field},
		},
	})
	return qb
}

// Max 最大值聚合
func (qb *MongoQueryBuilder) Max(field string) infrastructure.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":     nil,
			"maximum": map[string]interface{}{"$max": "$" + field},
		},
	})
	return qb
}

// Min 最小值聚合
func (qb *MongoQueryBuilder) Min(field string) infrastructure.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":     nil,
			"minimum": map[string]interface{}{"$min": "$" + field},
		},
	})
	return qb
}

// Build 构建查询
func (qb *MongoQueryBuilder) Build() (map[string]interface{}, error) {
	query := make(map[string]interface{})

	// 构建查询条件 - 使用条件分组
	if len(qb.groups) > 0 {
		// 收集所有有条件的组
		activeGroups := make([]conditionGroup, 0, len(qb.groups))
		totalConditions := 0
		for _, g := range qb.groups {
			if len(g.conditions) > 0 {
				activeGroups = append(activeGroups, g)
				totalConditions += len(g.conditions)
			}
		}

		if totalConditions == 0 {
			// 没有条件，返回空查询
		} else if len(activeGroups) == 1 && totalConditions == 1 {
			// 只有一个条件，直接添加到查询中
			for k, v := range activeGroups[0].conditions[0] {
				query[k] = v
			}
		} else if len(activeGroups) == 1 && activeGroups[0].logic == "and" {
			// 只有一个 AND 组，使用 $and
			query["$and"] = activeGroups[0].conditions
		} else if len(activeGroups) == 1 && activeGroups[0].logic == "or" {
			// 只有一个 OR 组，使用 $or
			query["$or"] = activeGroups[0].conditions
		} else if len(activeGroups) == 1 && activeGroups[0].logic == "not" {
			// 只有一个 NOT 组，使用 $nor
			query["$nor"] = activeGroups[0].conditions
		} else {
			// 多个组，需要组合
			// 策略：将每个组包装成子查询，然后使用 $and 组合
			groupQueries := make([]map[string]interface{}, 0, len(activeGroups))
			for _, g := range activeGroups {
				if len(g.conditions) == 1 {
					// 单个条件的组，直接添加
					groupQueries = append(groupQueries, g.conditions[0])
				} else if g.logic == "and" {
					groupQueries = append(groupQueries, map[string]interface{}{"$and": g.conditions})
				} else if g.logic == "or" {
					groupQueries = append(groupQueries, map[string]interface{}{"$or": g.conditions})
				} else if g.logic == "not" {
					groupQueries = append(groupQueries, map[string]interface{}{"$nor": g.conditions})
				}
			}
			if len(groupQueries) == 1 {
				// 只有一个组查询，直接添加
				for k, v := range groupQueries[0] {
					query[k] = v
				}
			} else {
				// 多个组查询，使用 $and 组合
				query["$and"] = groupQueries
			}
		}
	}

	// 构建选项
	options := make(map[string]interface{})

	// 排序
	if len(qb.sortFields) > 0 {
		options["sort"] = qb.sortFields
	}

	// 字段选择
	if len(qb.selectFields) > 0 {
		projection := make(map[string]interface{})
		for _, field := range qb.selectFields {
			projection[field] = 1
		}
		options["projection"] = projection
	} else if len(qb.excludeFields) > 0 {
		projection := make(map[string]interface{})
		for _, field := range qb.excludeFields {
			projection[field] = 0
		}
		options["projection"] = projection
	}

	// 分页
	if qb.limitValue > 0 {
		options["limit"] = qb.limitValue
	}
	if qb.skipValue > 0 {
		options["skip"] = qb.skipValue
	}

	result := map[string]interface{}{
		"filter":  query,
		"options": options,
	}

	// 聚合管道
	if len(qb.aggregations) > 0 {
		pipeline := make([]map[string]interface{}, 0)

		// 添加匹配阶段
		if len(query) > 0 {
			pipeline = append(pipeline, map[string]interface{}{"$match": query})
		}

		// 添加聚合阶段
		pipeline = append(pipeline, qb.aggregations...)

		result["pipeline"] = pipeline
	}

	return result, nil
}

// BuildFilter 构建Filter对象
func (qb *MongoQueryBuilder) BuildFilter() (infrastructure.Filter, error) {
	query, err := qb.Build()
	if err != nil {
		return nil, err
	}

	filter := &MongoFilter{
		Query:   query["filter"].(map[string]interface{}),
		Options: query["options"].(map[string]interface{}),
	}

	if pipeline, exists := query["pipeline"]; exists {
		filter.Pipeline = pipeline.([]map[string]interface{})
	}

	return filter, nil
}

// Reset 重置查询构建器
func (qb *MongoQueryBuilder) Reset() infrastructure.QueryBuilder {
	qb.groups = []conditionGroup{{
		logic:      "and",
		conditions: make([]map[string]interface{}, 0),
	}}
	qb.currentGroup = &qb.groups[0]
	qb.sortFields = make(map[string]int)
	qb.selectFields = make([]string, 0)
	qb.excludeFields = make([]string, 0)
	qb.limitValue = 0
	qb.skipValue = 0
	qb.aggregations = make([]map[string]interface{}, 0)
	return qb
}

// Clone 克隆查询构建器
func (qb *MongoQueryBuilder) Clone() infrastructure.QueryBuilder {
	clone := &MongoQueryBuilder{
		groups:        make([]conditionGroup, len(qb.groups)),
		sortFields:    make(map[string]int),
		selectFields:  make([]string, len(qb.selectFields)),
		excludeFields: make([]string, len(qb.excludeFields)),
		aggregations:  make([]map[string]interface{}, len(qb.aggregations)),
		limitValue:    qb.limitValue,
		skipValue:     qb.skipValue,
	}

	// 深拷贝 groups
	for i, g := range qb.groups {
		clone.groups[i] = conditionGroup{
			logic:      g.logic,
			conditions: make([]map[string]interface{}, len(g.conditions)),
		}
		copy(clone.groups[i].conditions, g.conditions)
	}
	// 设置 currentGroup 指向第一个组
	if len(clone.groups) > 0 {
		clone.currentGroup = &clone.groups[0]
	}

	copy(clone.selectFields, qb.selectFields)
	copy(clone.excludeFields, qb.excludeFields)
	copy(clone.aggregations, qb.aggregations)

	for k, v := range qb.sortFields {
		clone.sortFields[k] = v
	}

	return clone
}

// buildCondition 构建单个条件
// 如果操作符无效，返回 nil
func (qb *MongoQueryBuilder) buildCondition(field string, operator string, value interface{}) map[string]interface{} {
	if field == "" {
		return nil
	}

	condition := make(map[string]interface{})

	switch strings.ToLower(operator) {
	case "=", "eq":
		condition[field] = value
	case "!=", "ne":
		condition[field] = map[string]interface{}{"$ne": value}
	case ">", "gt":
		condition[field] = map[string]interface{}{"$gt": value}
	case ">=", "gte":
		condition[field] = map[string]interface{}{"$gte": value}
	case "<", "lt":
		condition[field] = map[string]interface{}{"$lt": value}
	case "<=", "lte":
		condition[field] = map[string]interface{}{"$lte": value}
	case "in":
		// 直接传递值，应该在 WhereIn 中处理
		return nil
	case "nin":
		// 直接传递值，应该在 WhereNotIn 中处理
		return nil
	case "between":
		// 直接传递值，应该在 WhereBetween 中处理
		return nil
	case "like":
		// WhereLike 中已经处理，这里不处理
		return nil
	case "regex":
		// WhereRegex 中已经处理，这里不处理
		return nil
	case "exists":
		condition[field] = map[string]interface{}{"$exists": value}
	case "type":
		condition[field] = map[string]interface{}{"$type": value}
	default:
		// 未知操作符，使用相等条件作为默认行为
		// 这样可以确保查询不会因为操作符错误而失败
		condition[field] = value
	}

	return condition
}

// MongoFilter MongoDB过滤器实现
type MongoFilter struct {
	Query    map[string]interface{}   `json:"query"`
	Options  map[string]interface{}   `json:"options"`
	Pipeline []map[string]interface{} `json:"pipeline,omitempty"`
}

// GetConditions 实现Filter接口
func (f *MongoFilter) GetConditions() map[string]interface{} {
	return f.Query
}

// GetSort 实现Filter接口
func (f *MongoFilter) GetSort() map[string]int {
	if sort, exists := f.Options["sort"]; exists {
		if sortMap, ok := sort.(map[string]int); ok {
			return sortMap
		}
	}
	return map[string]int{"createdAt": -1}
}

// GetFields 实现Filter接口
func (f *MongoFilter) GetFields() []string {
	if projection, exists := f.Options["projection"]; exists {
		if projMap, ok := projection.(map[string]interface{}); ok {
			fields := make([]string, 0, len(projMap))
			for field, include := range projMap {
				if include == 1 {
					fields = append(fields, field)
				}
			}
			return fields
		}
	}
	return nil
}

// Validate 实现Filter接口
func (f *MongoFilter) Validate() error {
	// 验证查询条件
	if f.Query != nil {
		if err := f.validateQuery(f.Query); err != nil {
			return err
		}
	}

	// 验证选项
	if f.Options != nil {
		if err := f.validateOptions(f.Options); err != nil {
			return err
		}
	}

	return nil
}

// validateQuery 验证查询条件
func (f *MongoFilter) validateQuery(query map[string]interface{}) error {
	for field, value := range query {
		if strings.HasPrefix(field, "$") {
			// 验证MongoDB操作符
			if !f.isValidOperator(field) {
				return infrastructure.NewValidationError(fmt.Sprintf("无效的查询操作符: %s", field))
			}
		}

		// 递归验证嵌套查询
		if nestedQuery, ok := value.(map[string]interface{}); ok {
			if err := f.validateQuery(nestedQuery); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateOptions 验证选项
func (f *MongoFilter) validateOptions(options map[string]interface{}) error {
	validOptions := map[string]bool{
		"sort":       true,
		"projection": true,
		"limit":      true,
		"skip":       true,
		"hint":       true,
		"collation":  true,
	}

	for option := range options {
		if !validOptions[option] {
			return infrastructure.NewValidationError(fmt.Sprintf("无效的查询选项: %s", option))
		}
	}

	// 验证limit和skip
	if limit, exists := options["limit"]; exists {
		if limitVal, ok := limit.(int); ok && limitVal < 0 {
			return infrastructure.NewValidationError("limit不能为负数")
		}
	}

	if skip, exists := options["skip"]; exists {
		if skipVal, ok := skip.(int); ok && skipVal < 0 {
			return infrastructure.NewValidationError("skip不能为负数")
		}
	}

	return nil
}

// isValidOperator 检查是否为有效的MongoDB操作符
func (f *MongoFilter) isValidOperator(operator string) bool {
	validOperators := map[string]bool{
		"$and": true, "$or": true, "$nor": true, "$not": true,
		"$eq": true, "$ne": true, "$gt": true, "$gte": true, "$lt": true, "$lte": true,
		"$in": true, "$nin": true, "$exists": true, "$type": true,
		"$regex": true, "$options": true, "$text": true, "$where": true,
		"$all": true, "$elemMatch": true, "$size": true,
		"$mod": true, "$geoWithin": true, "$geoIntersects": true, "$near": true,
		"$match": true, "$group": true, "$sort": true, "$limit": true, "$skip": true,
		"$project": true, "$unwind": true, "$lookup": true, "$count": true,
	}

	return validOperators[operator]
}
