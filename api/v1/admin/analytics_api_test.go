package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	adminService "Qingyu_backend/service/admin"
)

// TestAnalyticsAPI_GetUserGrowthTrend_Success 测试获取用户增长趋势成功
func TestAnalyticsAPI_GetUserGrowthTrend_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建模拟服务
	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	// 创建测试路由
	router := gin.New()
	router.GET("/analytics/user-growth", api.GetUserGrowthTrend)

	// 创建测试请求
	startDate := "2026-01-01"
	endDate := "2026-01-07"
	req, _ := http.NewRequest("GET", "/analytics/user-growth?start_date="+startDate+"&end_date="+endDate+"&interval=daily", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "获取成功", response["message"])
	assert.NotNil(t, response["data"])
}

// TestAnalyticsAPI_GetUserGrowthTrend_MissingDates 测试缺少日期参数
func TestAnalyticsAPI_GetUserGrowthTrend_MissingDates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/user-growth", api.GetUserGrowthTrend)

	req, _ := http.NewRequest("GET", "/analytics/user-growth?interval=daily", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "开始日期和结束日期不能为空")
}

// TestAnalyticsAPI_GetUserGrowthTrend_InvalidDateFormat 测试无效日期格式
func TestAnalyticsAPI_GetUserGrowthTrend_InvalidDateFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/user-growth", api.GetUserGrowthTrend)

	req, _ := http.NewRequest("GET", "/analytics/user-growth?start_date=invalid&end_date=2026-01-07&interval=daily", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "开始日期格式错误")
}

// TestAnalyticsAPI_GetUserGrowthTrend_InvalidInterval 测试无效间隔参数
func TestAnalyticsAPI_GetUserGrowthTrend_InvalidInterval(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/user-growth", api.GetUserGrowthTrend)

	req, _ := http.NewRequest("GET", "/analytics/user-growth?start_date=2026-01-01&end_date=2026-01-07&interval=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "间隔参数必须是")
}

// TestAnalyticsAPI_GetContentStatistics_Success 测试获取内容统计成功
func TestAnalyticsAPI_GetContentStatistics_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/content-statistics", api.GetContentStatistics)

	req, _ := http.NewRequest("GET", "/analytics/content-statistics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
}

// TestAnalyticsAPI_GetContentStatistics_WithDateRange 测试带日期范围
func TestAnalyticsAPI_GetContentStatistics_WithDateRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/content-statistics", api.GetContentStatistics)

	req, _ := http.NewRequest("GET", "/analytics/content-statistics?start_date=2026-01-01&end_date=2026-01-31", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAnalyticsAPI_GetContentStatistics_InvalidDateFormat 测试无效日期格式
func TestAnalyticsAPI_GetContentStatistics_InvalidDateFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/content-statistics", api.GetContentStatistics)

	req, _ := http.NewRequest("GET", "/analytics/content-statistics?start_date=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "开始日期格式错误")
}

// TestAnalyticsAPI_GetRevenueReport_Success 测试获取收入报告成功
func TestAnalyticsAPI_GetRevenueReport_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/revenue-report", api.GetRevenueReport)

	req, _ := http.NewRequest("GET", "/analytics/revenue-report?start_date=2026-01-01&end_date=2026-01-31&interval=daily", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
}

// TestAnalyticsAPI_GetRevenueReport_MissingParameters 测试缺少参数
func TestAnalyticsAPI_GetRevenueReport_MissingParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/revenue-report", api.GetRevenueReport)

	req, _ := http.NewRequest("GET", "/analytics/revenue-report", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "开始日期和结束日期不能为空")
}

// TestAnalyticsAPI_GetActiveUsersReport_Success 测试获取活跃用户报告成功
func TestAnalyticsAPI_GetActiveUsersReport_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/active-users", api.GetActiveUsersReport)

	req, _ := http.NewRequest("GET", "/analytics/active-users?start_date=2026-01-01&end_date=2026-01-07&type=dau", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
}

// TestAnalyticsAPI_GetActiveUsersReport_MissingType 测试缺少类型参数
func TestAnalyticsAPI_GetActiveUsersReport_MissingType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/active-users", api.GetActiveUsersReport)

	req, _ := http.NewRequest("GET", "/analytics/active-users?start_date=2026-01-01&end_date=2026-01-07", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "报告类型不能为空")
}

// TestAnalyticsAPI_GetActiveUsersReport_InvalidType 测试无效类型参数
func TestAnalyticsAPI_GetActiveUsersReport_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/active-users", api.GetActiveUsersReport)

	req, _ := http.NewRequest("GET", "/analytics/active-users?start_date=2026-01-01&end_date=2026-01-07&type=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "类型参数必须是")
}

// TestAnalyticsAPI_GetSystemOverview_Success 测试获取系统概览成功
func TestAnalyticsAPI_GetSystemOverview_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/system-overview", api.GetSystemOverview)

	req, _ := http.NewRequest("GET", "/analytics/system-overview", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])

	// 验证返回的数据结构
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "total_users")
	assert.Contains(t, data, "total_books")
	assert.Contains(t, data, "system_status")
}

// TestAnalyticsAPI_ExportAnalyticsReport_CSV 测试导出CSV报告
func TestAnalyticsAPI_ExportAnalyticsReport_CSV(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/export", api.ExportAnalyticsReport)

	req, _ := http.NewRequest("GET", "/analytics/export?report_type=user-growth&format=csv", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
}

// TestAnalyticsAPI_ExportAnalyticsReport_MissingReportType 测试缺少报告类型
func TestAnalyticsAPI_ExportAnalyticsReport_MissingReportType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/export", api.ExportAnalyticsReport)

	req, _ := http.NewRequest("GET", "/analytics/export", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "报告类型不能为空")
}

// TestAnalyticsAPI_ExportAnalyticsReport_InvalidReportType 测试无效报告类型
func TestAnalyticsAPI_ExportAnalyticsReport_InvalidReportType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/export", api.ExportAnalyticsReport)

	req, _ := http.NewRequest("GET", "/analytics/export?report_type=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "无效的报告类型")
}

// TestAnalyticsAPI_GetAnalyticsDashboard_Success 测试获取仪表板成功
func TestAnalyticsAPI_GetAnalyticsDashboard_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/dashboard", api.GetAnalyticsDashboard)

	req, _ := http.NewRequest("GET", "/analytics/dashboard", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])

	// Go JSON 序列化使用 PascalCase，所以字段名是 SystemOverview, UserGrowth 等
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "SystemOverview")
	assert.Contains(t, data, "UserGrowth")
	assert.Contains(t, data, "ContentStats")
	assert.Contains(t, data, "RevenueReport")
	assert.Contains(t, data, "ActiveUsers")
}

// TestAnalyticsAPI_GetCustomAnalyticsQuery_Success 测试自定义查询成功
func TestAnalyticsAPI_GetCustomAnalyticsQuery_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.POST("/analytics/custom", api.GetCustomAnalyticsQuery)

	query := CustomAnalyticsQuery{
		Metrics:   []string{"users", "revenue"},
		Dimensions: []string{"date"},
		StartDate:  "2026-01-01",
		EndDate:    "2026-01-31",
		Limit:      100,
	}

	jsonData, _ := json.Marshal(query)
	req, _ := http.NewRequest("POST", "/analytics/custom", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
}

// TestAnalyticsAPI_GetCustomAnalyticsQuery_InvalidJSON 测试无效JSON
func TestAnalyticsAPI_GetCustomAnalyticsQuery_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.POST("/analytics/custom", api.GetCustomAnalyticsQuery)

	req, _ := http.NewRequest("POST", "/analytics/custom", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAnalyticsAPI_CompareAnalyticsPeriods_Success 测试时期对比成功
func TestAnalyticsAPI_CompareAnalyticsPeriods_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/compare", api.CompareAnalyticsPeriods)

	req, _ := http.NewRequest("GET", "/analytics/compare?period1_start=2026-01-01&period1_end=2026-01-07&period2_start=2026-01-08&period2_end=2026-01-14&metrics=users,revenue", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "period1")
	assert.Contains(t, data, "period2")
	assert.Contains(t, data, "comparison")
}

// TestAnalyticsAPI_CompareAnalyticsPeriods_MissingParameters 测试缺少参数
func TestAnalyticsAPI_CompareAnalyticsPeriods_MissingParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/compare", api.CompareAnalyticsPeriods)

	req, _ := http.NewRequest("GET", "/analytics/compare?period1_start=2026-01-01", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "缺少必要参数")
}

// TestAnalyticsAPI_GetRealTimeStats_Success 测试获取实时统计成功
func TestAnalyticsAPI_GetRealTimeStats_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/realtime", api.GetRealTimeStats)

	req, _ := http.NewRequest("GET", "/analytics/realtime", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "online_users")
	assert.Contains(t, data, "active_sessions")
	assert.Contains(t, data, "timestamp")
}

// TestAnalyticsAPI_GetAnalyticsPredictions_Success 测试获取预测成功
func TestAnalyticsAPI_GetAnalyticsPredictions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/predict", api.GetAnalyticsPredictions)

	req, _ := http.NewRequest("GET", "/analytics/predict?type=users&days=30", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "users", data["type"])
	assert.Equal(t, float64(30), data["days"])
}

// TestAnalyticsAPI_GetAnalyticsPredictions_InvalidDays 测试无效天数
func TestAnalyticsAPI_GetAnalyticsPredictions_InvalidDays(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/predict", api.GetAnalyticsPredictions)

	req, _ := http.NewRequest("GET", "/analytics/predict?type=users&days=400", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "预测天数必须在1-365之间")
}

// TestAnalyticsAPI_GetAnalyticsPredictions_MissingType 测试缺少类型参数
func TestAnalyticsAPI_GetAnalyticsPredictions_MissingType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &mockAnalyticsService{}
	api := NewAnalyticsAPI(mockService)

	router := gin.New()
	router.GET("/analytics/predict", api.GetAnalyticsPredictions)

	req, _ := http.NewRequest("GET", "/analytics/predict", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "预测类型不能为空")
}

// =========================== Mock Service ===========================

// mockAnalyticsService 模拟统计分析服务
type mockAnalyticsService struct{}

func (m *mockAnalyticsService) GetUserGrowthTrend(ctx context.Context, req *adminService.UserGrowthTrendRequest) (*adminService.UserGrowthTrendResponse, error) {
	data := []adminService.UserGrowthDataPoint{
		{Date: "2026-01-01", Count: 100},
		{Date: "2026-01-02", Count: 120},
	}
	return &adminService.UserGrowthTrendResponse{
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Interval:      req.Interval,
		TotalNewUsers: 220,
		Data:          data,
		GrowthRate:    12.5,
	}, nil
}

func (m *mockAnalyticsService) GetContentStatistics(ctx context.Context, req *adminService.ContentStatisticsRequest) (*adminService.ContentStatisticsResponse, error) {
	return &adminService.ContentStatisticsResponse{
		TotalBooks:     1234,
		TotalChapters:  56789,
		TotalComments:  89012,
		TotalWords:     1234567890,
		PendingReviews: 45,
		PublishedToday: 12,
		CategoryStats: []adminService.CategoryStat{
			{CategoryName: "玄幻", BookCount: 234, ChapterCount: 12345},
		},
		TrendingBooks: []adminService.TrendingBook{
			{BookID: "book1", Title: "热门书籍1", AuthorID: "author1", ViewCount: 123456, ChapterCount: 120},
		},
	}, nil
}

func (m *mockAnalyticsService) GetRevenueReport(ctx context.Context, req *adminService.RevenueReportRequest) (*adminService.RevenueReportResponse, error) {
	data := []adminService.RevenueDataPoint{
		{Date: "2026-01-01", Amount: 1000.0, Orders: 50},
	}
	return &adminService.RevenueReportResponse{
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Interval:     req.Interval,
		TotalRevenue: 1000.0,
		Data:         data,
		GrowthRate:   8.3,
	}, nil
}

func (m *mockAnalyticsService) GetActiveUsersReport(ctx context.Context, req *adminService.ActiveUsersReportRequest) (*adminService.ActiveUsersReportResponse, error) {
	data := []adminService.ActiveUserDataPoint{
		{Date: "2026-01-01", Count: 500},
	}
	return &adminService.ActiveUsersReportResponse{
		StartDate:           req.StartDate,
		EndDate:             req.EndDate,
		Type:                req.Type,
		AverageActiveUsers:  500,
		PeakActiveUsers:     500,
		PeakDate:            "2026-01-01",
		Data:                data,
	}, nil
}

func (m *mockAnalyticsService) GetSystemOverview(ctx context.Context) (*adminService.SystemOverviewResponse, error) {
	return &adminService.SystemOverviewResponse{
		TotalUsers:     12345,
		NewUsersToday:  67,
		ActiveUsers:    890,
		TotalBooks:     1234,
		TotalChapters:  56789,
		TotalComments:  89012,
		PendingReviews: 45,
		TotalRevenue:   123456.78,
		RevenueToday:   1234.56,
		SystemStatus:   "healthy",
		LastUpdated:    time.Now(),
	}, nil
}
