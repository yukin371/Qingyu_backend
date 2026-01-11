package bookstore

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/api/v1/shared"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// BookStatisticsAPI 图书统计API控制器
type BookStatisticsAPI struct {
	BookStatisticsService bookstoreService.BookStatisticsService
}

// NewBookStatisticsAPI 创建新的图书统计API实例
func NewBookStatisticsAPI(bookStatisticsService bookstoreService.BookStatisticsService) *BookStatisticsAPI {
	return &BookStatisticsAPI{
		BookStatisticsService: bookStatisticsService,
	}
}

// GetBookStatistics 获取图书统计信息
// @Summary 获取图书统计信息
// @Description 根据图书ID获取统计信息
// @Tags 图书统计
// @Accept json
// @Produce json
// @Param book_id path string true "图书ID"
// @Success 200 {object} APIResponse "成功"
// @Failure 400 {object} APIResponse "请求参数错误"
// @Failure 404 {object} APIResponse "统计信息不存在"
// @Failure 500 {object} APIResponse "服务器内部错误"
// @Router /api/v1/reading/books/{book_id}/statistics [get]
func (api *BookStatisticsAPI) GetBookStatistics(c *gin.Context) {
	bookIDStr := c.Param("book_id")
	if bookIDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "图书ID不能为空",
			Data:    nil,
		})
		return
	}

	bookID, err := primitive.ObjectIDFromHex(bookIDStr)
