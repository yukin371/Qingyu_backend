package querybuilder

import (
	"testing"
)

// TestNewMongoQueryBuilder 测试创建 QueryBuilder
func TestNewMongoQueryBuilder(t *testing.T) {
	qb := NewMongoQueryBuilder()
	if qb == nil {
		t.Fatal("NewMongoQueryBuilder 返回 nil")
	}

	mqb, ok := qb.(*MongoQueryBuilder)
	if !ok {
		t.Fatal("NewMongoQueryBuilder 未返回 *MongoQueryBuilder 类型")
	}

	if len(mqb.groups) != 1 {
		t.Errorf("预期有 1 个组，实际有 %d", len(mqb.groups))
	}

	if mqb.groups[0].logic != "and" {
		t.Errorf("预期初始逻辑为 'and'，实际为 '%s'", mqb.groups[0].logic)
	}
}

// TestWhere 测试条件查询
func TestWhere(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.Where("name", "=", "John").
		Where("age", ">", 30).
		Where("status", "!=", "deleted")

	if len(qb.currentGroup.conditions) != 3 {
		t.Errorf("预期有 3 个条件，实际有 %d", len(qb.currentGroup.conditions))
	}
}

// TestWhereIn 测试 IN 条件
func TestWhereIn(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.WhereIn("status", []interface{}{"active", "pending"})

	if len(qb.currentGroup.conditions) != 1 {
		t.Errorf("预期有 1 个条件，实际有 %d", len(qb.currentGroup.conditions))
	}

	// 测试空数组
	qb.WhereIn("empty", []interface{}{})
	// 空数组应该不添加条件
}

// TestWhereNotIn 测试 NOT IN 条件
func TestWhereNotIn(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.WhereNotIn("status", []interface{}{"deleted", "banned"})

	if len(qb.currentGroup.conditions) != 1 {
		t.Errorf("预期有 1 个条件，实际有 %d", len(qb.currentGroup.conditions))
	}
}

// TestWhereBetween 测试范围条件
func TestWhereBetween(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.WhereBetween("price", 10, 100)

	if len(qb.currentGroup.conditions) != 1 {
		t.Errorf("预期有 1 个条件，实际有 %d", len(qb.currentGroup.conditions))
	}
}

// TestWhereNull 测试 NULL 条件
func TestWhereNull(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.WhereNull("deletedAt")

	if len(qb.currentGroup.conditions) != 1 {
		t.Errorf("预期有 1 个条件，实际有 %d", len(qb.currentGroup.conditions))
	}
}

// TestWhereNotNull 测试 NOT NULL 条件
func TestWhereNotNull(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.WhereNotNull("email")

	if len(qb.currentGroup.conditions) != 1 {
		t.Errorf("预期有 1 个条件，实际有 %d", len(qb.currentGroup.conditions))
	}
}

// TestWhereLike 测试 LIKE 条件
func TestWhereLike(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		expected string
	}{
		{"前缀匹配", "abc%", "^abc.*$"},
		{"后缀匹配", "%abc", "^.*abc$"},
		{"包含匹配", "%abc%", "^.*abc.*$"},
		{"单字符匹配", "abc_def", "^abc.def$"},
		{"转义字符", "abc\\%def", "^abc%def$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := NewMongoQueryBuilder().(*MongoQueryBuilder)
			qb.WhereLike("name", tt.pattern)

			result, err := qb.Build()
			if err != nil {
				t.Fatalf("Build 失败: %v", err)
			}

			filter := result["filter"].(map[string]interface{})
			nameCondition := filter["name"].(map[string]interface{})
			regex := nameCondition["$regex"].(string)

			if regex != tt.expected {
				t.Errorf("模式 '%s' 转换错误，预期 '%s'，实际 '%s'", tt.pattern, tt.expected, regex)
			}
		})
	}
}

// TestConvertLikeToRegex 测试 LIKE 模式转换
func TestConvertLikeToRegex(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	tests := []struct {
		pattern  string
		expected string
	}{
		{"%", "^.*$"},
		{"_", "^.$"},
		{"abc", "^abc$"},
		{"abc%", "^abc.*$"},
		{"%abc", "^.*abc$"},
		{"%abc%", "^.*abc.*$"},
		{"abc_def", "^abc.def$"},
		{"abc\\%def", "^abc%def$"},
		{"abc\\_def", "^abc_def$"},
		{"abc\\\\def", "^abc\\\\def$"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			result := qb.convertLikeToRegex(tt.pattern)
			if result != tt.expected {
				t.Errorf("convertLikeToRegex('%s') = '%s'，预期 '%s'", tt.pattern, result, tt.expected)
			}
		})
	}
}

// TestOrderBy 测试排序
func TestOrderBy(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.OrderBy("name", "asc").
		OrderBy("age", "desc").
		OrderByAsc("createdAt").
		OrderByDesc("updatedAt")

	if len(qb.sortFields) != 4 {
		t.Errorf("预期有 4 个排序字段，实际有 %d", len(qb.sortFields))
	}

	if qb.sortFields["name"] != 1 {
		t.Errorf("name 字段排序应为 1，实际为 %d", qb.sortFields["name"])
	}

	if qb.sortFields["age"] != -1 {
		t.Errorf("age 字段排序应为 -1，实际为 %d", qb.sortFields["age"])
	}
}

// TestSelectExclude 测试字段选择
func TestSelectExclude(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.Select("name", "email").
		Exclude("password", "secret")

	if len(qb.selectFields) != 2 {
		t.Errorf("预期有 2 个选择字段，实际有 %d", len(qb.selectFields))
	}

	if len(qb.excludeFields) != 2 {
		t.Errorf("预期有 2 个排除字段，实际有 %d", len(qb.excludeFields))
	}
}

// TestLimitSkip 测试分页
func TestLimitSkip(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.Limit(10).Skip(20)

	if qb.limitValue != 10 {
		t.Errorf("limit 应为 10，实际为 %d", qb.limitValue)
	}

	if qb.skipValue != 20 {
		t.Errorf("skip 应为 20，实际为 %d", qb.skipValue)
	}
}

// TestAggregations 测试聚合
func TestAggregations(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.Count().
		Sum("price").
		Avg("rating").
		Max("views").
		Min("likes")

	if len(qb.aggregations) != 5 {
		t.Errorf("预期有 5 个聚合操作，实际有 %d", len(qb.aggregations))
	}
}

// TestBuildSingleCondition 测试构建单个条件
func TestBuildSingleCondition(t *testing.T) {
	qb := NewMongoQueryBuilder()
	qb.Where("name", "=", "John")

	result, err := qb.Build()
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	filter := result["filter"].(map[string]interface{})
	if name, ok := filter["name"]; !ok || name != "John" {
		t.Errorf("单个条件构建错误，预期 {name: 'John'}，实际 %v", filter)
	}
}

// TestBuildMultipleConditionsSameGroup 测试构建多个条件（同一组）
func TestBuildMultipleConditionsSameGroup(t *testing.T) {
	qb := NewMongoQueryBuilder()
	qb.Where("name", "=", "John").
		Where("age", ">", 30)

	result, err := qb.Build()
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	filter := result["filter"].(map[string]interface{})
	if _, ok := filter["$and"]; !ok {
		t.Errorf("多个条件应使用 $and，实际 %v", filter)
	}
}

// TestBuildOrConditions 测试构建 OR 条件
func TestBuildOrConditions(t *testing.T) {
	qb := NewMongoQueryBuilder()
	qb.Where("status", "=", "active").
		Or().
		Where("status", "=", "pending")

	result, err := qb.Build()
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	filter := result["filter"].(map[string]interface{})
	if _, ok := filter["$and"]; !ok {
		t.Errorf("OR 条件应包装在 $and 中，实际 %v", filter)
	}

	// 验证内部结构
	andConditions := filter["$and"].([]map[string]interface{})
	if len(andConditions) != 2 {
		t.Errorf("预期有 2 个组，实际有 %d", len(andConditions))
	}
}

// TestBuildComplexConditions 测试复杂条件
func TestBuildComplexConditions(t *testing.T) {
	qb := NewMongoQueryBuilder()
	qb.Where("status", "=", "active").
		Where("role", "=", "user").
		Or().
		Where("status", "=", "pending").
		Where("role", "=", "admin")

	result, err := qb.Build()
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	filter := result["filter"].(map[string]interface{})
	if _, ok := filter["$and"]; !ok {
		t.Errorf("复杂条件应使用 $and，实际 %v", filter)
	}
}

// TestBuildWithSortAndPagination 测试构建带排序和分页的查询
func TestBuildWithSortAndPagination(t *testing.T) {
	qb := NewMongoQueryBuilder()
	qb.Where("status", "=", "active").
		OrderBy("createdAt", "desc").
		Limit(10).
		Skip(20)

	result, err := qb.Build()
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	options := result["options"].(map[string]interface{})

	// 验证排序
	sort := options["sort"].(map[string]int)
	if sort["createdAt"] != -1 {
		t.Errorf("createdAt 应降序排序，实际为 %d", sort["createdAt"])
	}

	// 验证分页
	if options["limit"] != 10 {
		t.Errorf("limit 应为 10，实际为 %v", options["limit"])
	}

	if options["skip"] != 20 {
		t.Errorf("skip 应为 20，实际为 %v", options["skip"])
	}
}

// TestBuildWithProjection 测试构建带字段投影的查询
func TestBuildWithProjection(t *testing.T) {
	qb := NewMongoQueryBuilder()
	qb.Select("name", "email").
		Where("status", "=", "active")

	result, err := qb.Build()
	if err != nil {
		t.Fatalf("Build 失败: %v", err)
	}

	options := result["options"].(map[string]interface{})
	projection := options["projection"].(map[string]interface{})

	if projection["name"] != 1 {
		t.Errorf("name 字段应被选择")
	}

	if projection["email"] != 1 {
		t.Errorf("email 字段应被选择")
	}
}

// TestReset 测试重置
func TestReset(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.Where("name", "=", "John").
		OrderBy("age", "desc").
		Limit(10).
		Select("name")

	qb.Reset()

	if len(qb.groups) != 1 || len(qb.groups[0].conditions) != 0 {
		t.Errorf("Reset 后条件应被清空")
	}

	if len(qb.sortFields) != 0 {
		t.Errorf("Reset 后排序字段应被清空")
	}

	if qb.limitValue != 0 {
		t.Errorf("Reset 后 limit 应为 0，实际为 %d", qb.limitValue)
	}
}

// TestClone 测试克隆
func TestClone(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	qb.Where("name", "=", "John").
		OrderBy("age", "desc").
		Limit(10)

	clone := qb.Clone().(*MongoQueryBuilder)

	// 验证克隆的数据
	if len(clone.groups) != len(qb.groups) {
		t.Errorf("克隆的组数量应与原对象相同")
	}

	if len(clone.sortFields) != len(qb.sortFields) {
		t.Errorf("克隆的排序字段数量应与原对象相同")
	}

	if clone.limitValue != qb.limitValue {
		t.Errorf("克隆的 limit 应与原对象相同")
	}

	// 验证深拷贝 - 修改克隆不应影响原对象
	clone.Where("status", "=", "active")

	if len(qb.currentGroup.conditions) != 1 {
		t.Errorf("修改克隆对象不应影响原对象")
	}

	if len(clone.currentGroup.conditions) != 2 {
		t.Errorf("克隆对象应包含原始条件和新添加的条件")
	}
}

// TestBuildCondition 测试构建条件
func TestBuildCondition(t *testing.T) {
	qb := NewMongoQueryBuilder().(*MongoQueryBuilder)

	tests := []struct {
		name     string
		field    string
		operator string
		value    interface{}
		wantNil  bool
	}{
		{"等于", "field", "=", "test", false},
		{"不等于", "field", "!=", "test", false},
		{"大于", "field", ">", 10, false},
		{"大于等于", "field", ">=", 10, false},
		{"小于", "field", "<", 10, false},
		{"小于等于", "field", "<=", 10, false},
		{"存在", "field", "exists", true, false},
		{"类型", "field", "type", "string", false},
		{"未知操作符", "field", "unknown", "test", false},
		{"空字段名", "", "=", "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := qb.buildCondition(tt.field, tt.operator, tt.value)
			if tt.wantNil {
				if result != nil {
					t.Errorf("预期返回 nil，实际返回 %v", result)
				}
				return
			}
			if result == nil {
				t.Fatal("buildCondition 返回 nil")
			}
			// 验证条件包含正确的字段
			if _, ok := result[tt.field]; !ok && tt.field != "" {
				t.Errorf("buildCondition 未包含字段 %s", tt.field)
			}
		})
	}
}

// TestMongoFilter 测试 MongoFilter
func TestMongoFilter(t *testing.T) {
	qb := NewMongoQueryBuilder()
	qb.Where("status", "=", "active").
		OrderBy("createdAt", "desc").
		Limit(10)

	filter, err := qb.BuildFilter()
	if err != nil {
		t.Fatalf("BuildFilter 失败: %v", err)
	}

	conditions := filter.GetConditions()
	if conditions["status"] != "active" {
		t.Errorf("GetConditions 返回错误，预期 status=active，实际 %v", conditions)
	}

	sort := filter.GetSort()
	if sort["createdAt"] != -1 {
		t.Errorf("GetSort 返回错误，预期 createdAt=-1，实际 %v", sort)
	}

	fields := filter.GetFields()
	if fields != nil {
		t.Errorf("GetFields 应返回 nil（未设置投影），实际 %v", fields)
	}

	// 测试验证
	if err := filter.Validate(); err != nil {
		t.Errorf("Validate 失败: %v", err)
	}
}
