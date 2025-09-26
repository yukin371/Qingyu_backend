package mongodb

import (
	"fmt"
	"strings"

	interfaces "Qingyu_backend/repository/interfaces"
)

// MongoQueryBuilder MongoDB查询构建器实现
type MongoQueryBuilder struct {
	conditions    []map[string]interface{}
	sortFields    map[string]int
	selectFields  []string
	excludeFields []string
	limitValue    int
	skipValue     int
	aggregations  []map[string]interface{}
	currentLogic  string // "and", "or", "not"
}

// NewMongoQueryBuilder 创建MongoDB查询构建器
func NewMongoQueryBuilder() interfaces.QueryBuilder {
	return &MongoQueryBuilder{
		conditions:    make([]map[string]interface{}, 0),
		sortFields:    make(map[string]int),
		selectFields:  make([]string, 0),
		excludeFields: make([]string, 0),
		aggregations:  make([]map[string]interface{}, 0),
		currentLogic:  "and",
	}
}

// Where 添加条件
func (qb *MongoQueryBuilder) Where(field string, operator string, value interface{}) interfaces.QueryBuilder {
	condition := qb.buildCondition(field, operator, value)
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereIn 添加IN条件
func (qb *MongoQueryBuilder) WhereIn(field string, values []interface{}) interfaces.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{"$in": values},
	}
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereNotIn 添加NOT IN条件
func (qb *MongoQueryBuilder) WhereNotIn(field string, values []interface{}) interfaces.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{"$nin": values},
	}
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereBetween 添加范围条件
func (qb *MongoQueryBuilder) WhereBetween(field string, start, end interface{}) interfaces.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{
			"$gte": start,
			"$lte": end,
		},
	}
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereNull 添加NULL条件
func (qb *MongoQueryBuilder) WhereNull(field string) interfaces.QueryBuilder {
	condition := map[string]interface{}{
		field: nil,
	}
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereNotNull 添加NOT NULL条件
func (qb *MongoQueryBuilder) WhereNotNull(field string) interfaces.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{"$ne": nil},
	}
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereLike 添加LIKE条件（模糊匹配）
func (qb *MongoQueryBuilder) WhereLike(field string, pattern string) interfaces.QueryBuilder {
	// 转换SQL LIKE模式到MongoDB正则表达式
	regexPattern := strings.ReplaceAll(pattern, "%", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "_", ".")

	condition := map[string]interface{}{
		field: map[string]interface{}{
			"$regex":   regexPattern,
			"$options": "i", // 不区分大小写
		},
	}
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereRegex 添加正则表达式条件
func (qb *MongoQueryBuilder) WhereRegex(field string, pattern string) interfaces.QueryBuilder {
	condition := map[string]interface{}{
		field: map[string]interface{}{
			"$regex":   pattern,
			"$options": "i",
		},
	}
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// And 设置AND逻辑
func (qb *MongoQueryBuilder) And() interfaces.QueryBuilder {
	qb.currentLogic = "and"
	return qb
}

// Or 设置OR逻辑
func (qb *MongoQueryBuilder) Or() interfaces.QueryBuilder {
	qb.currentLogic = "or"
	return qb
}

// Not 设置NOT逻辑
func (qb *MongoQueryBuilder) Not() interfaces.QueryBuilder {
	qb.currentLogic = "not"
	return qb
}

// OrderBy 添加排序
func (qb *MongoQueryBuilder) OrderBy(field string, direction string) interfaces.QueryBuilder {
	dir := 1
	if strings.ToLower(direction) == "desc" {
		dir = -1
	}
	qb.sortFields[field] = dir
	return qb
}

// OrderByAsc 升序排序
func (qb *MongoQueryBuilder) OrderByAsc(field string) interfaces.QueryBuilder {
	qb.sortFields[field] = 1
	return qb
}

// OrderByDesc 降序排序
func (qb *MongoQueryBuilder) OrderByDesc(field string) interfaces.QueryBuilder {
	qb.sortFields[field] = -1
	return qb
}

// Select 选择字段
func (qb *MongoQueryBuilder) Select(fields ...string) interfaces.QueryBuilder {
	qb.selectFields = append(qb.selectFields, fields...)
	return qb
}

// Exclude 排除字段
func (qb *MongoQueryBuilder) Exclude(fields ...string) interfaces.QueryBuilder {
	qb.excludeFields = append(qb.excludeFields, fields...)
	return qb
}

// Limit 限制数量
func (qb *MongoQueryBuilder) Limit(limit int) interfaces.QueryBuilder {
	qb.limitValue = limit
	return qb
}

// Skip 跳过数量
func (qb *MongoQueryBuilder) Skip(skip int) interfaces.QueryBuilder {
	qb.skipValue = skip
	return qb
}

// Count 计数聚合
func (qb *MongoQueryBuilder) Count() interfaces.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$count": "total",
	})
	return qb
}

// Sum 求和聚合
func (qb *MongoQueryBuilder) Sum(field string) interfaces.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":   nil,
			"total": map[string]interface{}{"$sum": "$" + field},
		},
	})
	return qb
}

// Avg 平均值聚合
func (qb *MongoQueryBuilder) Avg(field string) interfaces.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":     nil,
			"average": map[string]interface{}{"$avg": "$" + field},
		},
	})
	return qb
}

// Max 最大值聚合
func (qb *MongoQueryBuilder) Max(field string) interfaces.QueryBuilder {
	qb.aggregations = append(qb.aggregations, map[string]interface{}{
		"$group": map[string]interface{}{
			"_id":     nil,
			"maximum": map[string]interface{}{"$max": "$" + field},
		},
	})
	return qb
}

// Min 最小值聚合
func (qb *MongoQueryBuilder) Min(field string) interfaces.QueryBuilder {
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

	// 构建查询条件
	if len(qb.conditions) > 0 {
		if len(qb.conditions) == 1 {
			for k, v := range qb.conditions[0] {
				query[k] = v
			}
		} else {
			switch qb.currentLogic {
			case "and":
				query["$and"] = qb.conditions
			case "or":
				query["$or"] = qb.conditions
			case "not":
				query["$nor"] = qb.conditions
			default:
				query["$and"] = qb.conditions
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
func (qb *MongoQueryBuilder) BuildFilter() (interfaces.Filter, error) {
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
func (qb *MongoQueryBuilder) Reset() interfaces.QueryBuilder {
	qb.conditions = make([]map[string]interface{}, 0)
	qb.sortFields = make(map[string]int)
	qb.selectFields = make([]string, 0)
	qb.excludeFields = make([]string, 0)
	qb.limitValue = 0
	qb.skipValue = 0
	qb.aggregations = make([]map[string]interface{}, 0)
	qb.currentLogic = "and"
	return qb
}

// Clone 克隆查询构建器
func (qb *MongoQueryBuilder) Clone() interfaces.QueryBuilder {
	clone := &MongoQueryBuilder{
		conditions:    make([]map[string]interface{}, len(qb.conditions)),
		sortFields:    make(map[string]int),
		selectFields:  make([]string, len(qb.selectFields)),
		excludeFields: make([]string, len(qb.excludeFields)),
		aggregations:  make([]map[string]interface{}, len(qb.aggregations)),
		limitValue:    qb.limitValue,
		skipValue:     qb.skipValue,
		currentLogic:  qb.currentLogic,
	}

	copy(clone.conditions, qb.conditions)
	copy(clone.selectFields, qb.selectFields)
	copy(clone.excludeFields, qb.excludeFields)
	copy(clone.aggregations, qb.aggregations)

	for k, v := range qb.sortFields {
		clone.sortFields[k] = v
	}

	return clone
}

// buildCondition 构建单个条件
func (qb *MongoQueryBuilder) buildCondition(field string, operator string, value interface{}) map[string]interface{} {
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
	case "like":
		pattern := fmt.Sprintf(".*%s.*", value)
		condition[field] = map[string]interface{}{
			"$regex":   pattern,
			"$options": "i",
		}
	case "exists":
		condition[field] = map[string]interface{}{"$exists": value}
	case "type":
		condition[field] = map[string]interface{}{"$type": value}
	default:
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
				return interfaces.NewValidationError(fmt.Sprintf("无效的查询操作符: %s", field))
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
			return interfaces.NewValidationError(fmt.Sprintf("无效的查询选项: %s", option))
		}
	}

	// 验证limit和skip
	if limit, exists := options["limit"]; exists {
		if limitVal, ok := limit.(int); ok && limitVal < 0 {
			return interfaces.NewValidationError("limit不能为负数")
		}
	}

	if skip, exists := options["skip"]; exists {
		if skipVal, ok := skip.(int); ok && skipVal < 0 {
			return interfaces.NewValidationError("skip不能为负数")
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
