package bookstore

import (
	"strings"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// BookDetailAPI 书籍详情API处理器
type BookDetailAPI struct {
	bookstoreService bookstoreService.BookstoreService
}

// NewBookDetailAPI 创建书籍详情API实例
func NewBookDetailAPI(bookstoreService bookstoreService.BookstoreService) *BookDetailAPI {
	return &BookDetailAPI{
		bookstoreService: bookstoreService,
	}
}

// GetBookDetail 获取书籍详情
//
//	@Summary		获取书籍详情
//	@Description	根据书籍ID获取详细信息
//	@Tags			书籍详情
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"书籍ID"
//	@Success 200 {object} response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/detail [get]
func (api *BookDetailAPI) GetBookDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	book, err := api.bookstoreService.GetBookByID(c.Request.Context(), idStr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not available") {
			response.NotFound(c, "书籍不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, book)
}

// IncrementViewCount 增加书籍浏览量
//
//	@Summary		增加书籍浏览量
//	@Description	记录用户浏览书籍详情，增加浏览量统计
//	@Tags			书籍交互
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"书籍ID"
//	@Success		200		{object}	response.APIResponse
//	@Failure		400		{object}	response.APIResponse
//	@Failure		404		{object}	response.APIResponse
//	@Failure		500		{object}	response.APIResponse
//	@Router			/api/v1/bookstore/books/{id}/view [post]
func (api *BookDetailAPI) IncrementViewCount(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.BadRequest(c, "参数错误", "书籍ID不能为空")
		return
	}

	err := api.bookstoreService.IncrementBookView(c.Request.Context(), idStr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not available") {
			response.NotFound(c, "书籍不存在")
			return
		}
		response.InternalError(c, err)
		return
	}

	response.SuccessWithMessage(c, "浏览量增加成功", nil)
}
