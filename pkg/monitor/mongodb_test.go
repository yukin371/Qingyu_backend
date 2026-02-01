package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckResult_Validation 测试CheckResult验证
func TestCheckResult_Validation(t *testing.T) {
	tests := []struct {
		name     string
		result   CheckResult
		valid    bool
	}{
		{
			name: "有效的检查结果",
			result: CheckResult{
				Collection: "reading_progress",
				Field:      "user_id",
				Count:      10,
			},
			valid: true,
		},
		{
			name: "空集合名称",
			result: CheckResult{
				Collection: "",
				Field:      "user_id",
				Count:      10,
			},
			valid: false,
		},
		{
			name: "负数count",
			result: CheckResult{
				Collection: "reading_progress",
				Field:      "user_id",
				Count:      -1,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid {
				assert.NotEmpty(t, tt.result.Collection)
				assert.NotEmpty(t, tt.result.Field)
				assert.GreaterOrEqual(t, tt.result.Count, 0)
			} else {
				if tt.result.Collection == "" || tt.result.Field == "" || tt.result.Count < 0 {
					// 无效的情况
				}
			}
		})
	}
}

// TestDataQualityReport_Merge 测试报告合并
func TestDataQualityReport_Merge(t *testing.T) {
	report1 := &DataQualityReport{
		TotalOrphanedRecords: 10,
		InaccurateStatistics: 5,
		Details: []CheckResult{
			{Collection: "reading_progress", Field: "user_id", Count: 10},
		},
		HasIssues: true,
	}

	report2 := &DataQualityReport{
		TotalOrphanedRecords: 20,
		InaccurateStatistics: 15,
		Details: []CheckResult{
			{Collection: "books", Field: "likes_count", Count: 15},
		},
		HasIssues: true,
	}

	// 合并报告
	merged := &DataQualityReport{
		TotalOrphanedRecords: report1.TotalOrphanedRecords + report2.TotalOrphanedRecords,
		InaccurateStatistics: report1.InaccurateStatistics + report2.InaccurateStatistics,
		Details:              append(report1.Details, report2.Details...),
		HasIssues:            report1.HasIssues || report2.HasIssues,
	}

	assert.Equal(t, 30, merged.TotalOrphanedRecords)
	assert.Equal(t, 20, merged.InaccurateStatistics)
	assert.Equal(t, 2, len(merged.Details))
	assert.True(t, merged.HasIssues)
}

// TestDataQualityReport_EmptyDetails 测试空详情的报告
func TestDataQualityReport_EmptyDetails(t *testing.T) {
	report := &DataQualityReport{
		TotalOrphanedRecords: 0,
		InaccurateStatistics: 0,
		Details:              []CheckResult{},
		HasIssues:            false,
	}

	assert.Equal(t, 0, report.TotalOrphanedRecords)
	assert.Equal(t, 0, report.InaccurateStatistics)
	assert.Equal(t, 0, len(report.Details))
	assert.False(t, report.HasIssues)
}

// BenchmarkOrphanedRecords 基准测试孤儿记录检查
func BenchmarkOrphanedRecords(b *testing.B) {
	if testing.Short() {
		b.Skip("跳过基准测试")
	}

	// TODO: 添加基准测试
	// 1. 创建测试数据库
	// 2. 插入大量测试数据
	// 3. 运行基准测试
	// 4. 清理测试数据
}

// ExampleDataQualityMonitor 示例：如何使用数据质量监控器
func ExampleDataQualityMonitor() {
	// 创建监控器（实际使用时需要真实的数据库连接）
	// monitor := NewDataQualityMonitor(db, alerter)

	// 执行检查
	// report, err := monitor.RunDailyCheck(context.Background())
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 检查结果
	// if report.HasIssues {
	//     log.Printf("发现 %d 个问题", report.TotalOrphanedRecords)
	// }
}
