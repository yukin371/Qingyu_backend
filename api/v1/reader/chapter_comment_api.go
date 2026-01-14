package reader

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"Qingyu_backend/api/v1/shared"
	readerModels "Qingyu_backend/models/reader"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterCommentAPI 章节评论API
type ChapterCommentAPI struct {
	// 可以注入ChapterCommentService
}

// NewChapterCommentAPI 创建章节评论API实例
func NewChapterCommentAPI() *ChapterCommentAPI {
	return &ChapterCommentAPI{}
}

// GetChapterComments 获取章节评论列表
//
//	@Summary	获取章节评论列表
//	@Tags		阅读器-章节评论
//	@Param		chapterId	path	string	true	"章节ID"
//	@Param		page		query	int		false	"页码"	default(1)
//	@Param		pageSize	query	int		false	"每页数量"	default(20)
//	@Param		sortBy		query	string	false	"排序字段：created_at/like_count/rating"	default(created_at)
//	@Param		sortOrder	query	string	false	"排序方向：asc/desc"	default(desc)
//	@Param		parentId	query	string	false	"父评论ID（空字符串表示顶级评论）"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/comments [get]
func (api *ChapterCommentAPI) GetChapterComments(c *gin.Context) {
	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.Error(c, http.StatusBadRequest, "章节ID不能为空", "章节ID不能为空")
		return
	}

	// 验证章节ID格式
	if _, err := primitive.ObjectIDFromHex(chapterID); err != nil {
		shared.Error(c, http.StatusBadRequest, "章节ID格式无效", "章节ID格式无效")
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 20
	if p, ok := c.GetQuery("page"); ok {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	if ps, ok := c.GetQuery("pageSize"); ok {
		if n, err := strconv.Atoi(ps); err == nil && n > 0 && n <= 100 {
			pageSize = n
		}
	}

	// 获取排序参数
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	// 验证排序字段
	validSortFields := map[string]bool{
		"created_at": true,
		"like_count": true,
		"rating":     true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at" //nolint:ineffassign // TODO: 实现排序功能
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// 获取父评论ID参数
	var _ *string
	if parentIdParam, ok := c.GetQuery("parentId"); ok {
		if parentIdParam == "" {
			// 空字符串表示查询顶级评论
			emptyStr := ""
			_ = &emptyStr
		} else {
			// 验证父评论ID格式
			if _, err := primitive.ObjectIDFromHex(parentIdParam); err == nil {
				_ = &parentIdParam
			}
		}
	}

	// 实际应用中应该从数据库查询
	// 这里返回模拟数据
	comments := make([]*readerModels.ChapterComment, 0)
	total := int64(0)

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// 计算平均评分和评分数量
	avgRating := 0.0
	ratingCount := 0

	shared.Success(c, http.StatusOK, "获取成功", readerModels.ChapterCommentListResponse{
		Comments:    comments,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		AvgRating:   avgRating,
		RatingCount: ratingCount,
	})
}

// CreateChapterComment 发表章节评论
//
//	@Summary	发表章节评论
//	@Tags		阅读器-章节评论
//	@Param		chapterId	path	string								true	"章节ID"
//	@Param		request		body	reader.CreateChapterCommentRequest	true	"评论内容"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/comments [post]
func (api *ChapterCommentAPI) CreateChapterComment(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.Error(c, http.StatusBadRequest, "章节ID不能为空", "章节ID不能为空")
		return
	}

	// 验证章节ID格式
	if _, err := primitive.ObjectIDFromHex(chapterID); err != nil {
		shared.Error(c, http.StatusBadRequest, "章节ID格式无效", "章节ID格式无效")
		return
	}

	var req readerModels.CreateChapterCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 确保ChapterID匹配
	req.ChapterID = chapterID

	// 验证评分
	if req.Rating < 0 || req.Rating > 5 {
		shared.Error(c, http.StatusBadRequest, "评分必须在0-5之间", "评分必须在0-5之间")
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if req.ParentID != nil && *req.ParentID != "" {
		// 验证父评论ID格式
		if _, err := primitive.ObjectIDFromHex(*req.ParentID); err != nil {
			shared.Error(c, http.StatusBadRequest, "父评论ID格式无效", "父评论ID格式无效")
			return
		}
		// 实际应用中应该查询数据库验证父评论存在
	}

	// 创建评论
	comment := &readerModels.ChapterComment{
		ID:         primitive.NewObjectID().Hex(),
		ChapterID:  req.ChapterID,
		BookID:     req.BookID,
		UserID:     userID.(string),
		Content:    req.Content,
		Rating:     req.Rating,
		ParentID:   req.ParentID,
		ReplyCount: 0,
		LikeCount:  0,
		IsVisible:  true,
		IsDeleted:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 段落级评论
	if req.ParagraphIndex != nil {
		comment.ParagraphIndex = req.ParagraphIndex
		comment.CharStart = req.CharStart
		comment.CharEnd = req.CharEnd
	}

	// 实际应用中应该：
	// 1. 保存到数据库
	// 2. 如果是回复，更新父评论的reply_count
	// 3. 如果是顶级评论且有评分，更新书籍/章节的评分统计
	// 4. 发布评论创建事件

	shared.Success(c, http.StatusCreated, "评论成功", gin.H{
		"comment": comment,
		"message": "评论发表成功",
	})
}

// GetChapterComment 获取单条评论详情
//
//	@Summary	获取单条评论详情
//	@Tags		阅读器-章节评论
//	@Param		commentId	path	string	true	"评论ID"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{commentId} [get]
func (api *ChapterCommentAPI) GetChapterComment(c *gin.Context) {
	commentID := c.Param("commentId")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "评论ID不能为空", "评论ID不能为空")
		return
	}

	// 验证评论ID格式
	if _, err := primitive.ObjectIDFromHex(commentID); err != nil {
		shared.Error(c, http.StatusBadRequest, "评论ID格式无效", "评论ID格式无效")
		return
	}

	// 实际应用中应该从数据库查询
	shared.Error(c, http.StatusNotFound, "评论不存在", "未找到指定评论")
}

// UpdateChapterComment 更新章节评论
//
//	@Summary	更新章节评论
//	@Tags		阅读器-章节评论
//	@Param		commentId	path	string								true	"评论ID"
//	@Param		request		body	reader.UpdateChapterCommentRequest	true	"更新内容"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{commentId} [put]
func (api *ChapterCommentAPI) UpdateChapterComment(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "评论ID不能为空", "评论ID不能为空")
		return
	}

	var req readerModels.UpdateChapterCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 验证评分
	if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 5) {
		shared.Error(c, http.StatusBadRequest, "评分必须在0-5之间", "评分必须在0-5之间")
		return
	}

	// 实际应用中应该：
	// 1. 从数据库获取评论
	// 2. 验证用户是否为评论作者
	// 3. 检查是否在可编辑时间内（30分钟）
	// 4. 更新评论内容
	// 5. 如果评分改变，更新章节/书籍评分统计

	shared.Success(c, http.StatusOK, "更新成功", gin.H{
		"message":   "评论更新成功",
		"commentId": commentID,
		"userId":    userID,
	})
}

// DeleteChapterComment 删除章节评论
//
//	@Summary	删除章节评论
//	@Tags		阅读器-章节评论
//	@Param		commentId	path	string	true	"评论ID"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{commentId} [delete]
func (api *ChapterCommentAPI) DeleteChapterComment(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "评论ID不能为空", "评论ID不能为空")
		return
	}

	// 实际应用中应该：
	// 1. 从数据库获取评论
	// 2. 验证用户权限（评论作者或管理员）
	// 3. 软删除评论（设置is_deleted=true）
	// 4. 如果有评分，更新章节/书籍评分统计
	// 5. 如果是回复，减少父评论的reply_count

	shared.Success(c, http.StatusOK, "删除成功", gin.H{
		"message":   "评论删除成功",
		"commentId": commentID,
		"userId":    userID,
	})
}

// LikeChapterComment 点赞章节评论
//
//	@Summary	点赞章节评论
//	@Tags		阅读器-章节评论
//	@Param		commentId	path	string	true	"评论ID"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{commentId}/like [post]
func (api *ChapterCommentAPI) LikeChapterComment(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "评论ID不能为空", "评论ID不能为空")
		return
	}

	// 实际应用中应该：
	// 1. 检查用户是否已点赞
	// 2. 如果未点赞，创建点赞记录并增加评论的like_count
	// 3. 如果已点赞，取消点赞并减少like_count

	shared.Success(c, http.StatusOK, "点赞成功", gin.H{
		"message":   "评论点赞成功",
		"commentId": commentID,
		"userId":    userID,
	})
}

// UnlikeChapterComment 取消点赞章节评论
//
//	@Summary	取消点赞章节评论
//	@Tags		阅读器-章节评论
//	@Param		commentId	path	string	true	"评论ID"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/comments/{commentId}/like [delete]
func (api *ChapterCommentAPI) UnlikeChapterComment(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "评论ID不能为空", "评论ID不能为空")
		return
	}

	shared.Success(c, http.StatusOK, "取消点赞成功", gin.H{
		"message":   "取消点赞成功",
		"commentId": commentID,
		"userId":    userID,
	})
}

// GetParagraphComments 获取段落级评论
//
//	@Summary	获取段落级评论
//	@Tags		阅读器-段落评论
//	@Param		chapterId		path	string	true	"章节ID"
//	@Param		paragraphIndex	query	int		true	"段落索引"
//	@Success	200				{object}	shared.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/paragraphs/{paragraphIndex}/comments [get]
func (api *ChapterCommentAPI) GetParagraphComments(c *gin.Context) {
	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.Error(c, http.StatusBadRequest, "章节ID不能为空", "章节ID不能为空")
		return
	}

	paragraphIndexStr := c.Param("paragraphIndex")
	if paragraphIndexStr == "" {
		shared.Error(c, http.StatusBadRequest, "段落索引不能为空", "段落索引不能为空")
		return
	}

	paragraphIndex, err := strconv.Atoi(paragraphIndexStr)
	if err != nil || paragraphIndex < 0 {
		shared.Error(c, http.StatusBadRequest, "段落索引格式无效", "段落索引格式无效")
		return
	}

	// 实际应用中应该从数据库查询
	comments := make([]*readerModels.ChapterComment, 0)

	shared.Success(c, http.StatusOK, "获取成功", readerModels.ParagraphCommentResponse{
		ParagraphIndex: paragraphIndex,
		ParagraphText:  "", // 应该从章节内容中获取
		CommentCount:   len(comments),
		Comments:       comments,
	})
}

// CreateParagraphComment 发表段落级评论
//
//	@Summary	发表段落级评论
//	@Tags		阅读器-段落评论
//	@Param		chapterId	path	string								true	"章节ID"
//	@Param		request		body	reader.CreateChapterCommentRequest	true	"评论内容"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/paragraph-comments [post]
func (api *ChapterCommentAPI) CreateParagraphComment(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.Error(c, http.StatusBadRequest, "章节ID不能为空", "章节ID不能为空")
		return
	}

	var req readerModels.CreateChapterCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ValidationError(c, err)
		return
	}

	// 段落级评论必须指定段落索引
	if req.ParagraphIndex == nil {
		shared.Error(c, http.StatusBadRequest, "段落索引不能为空", "段落索引不能为空")
		return
	}

	// 创建段落评论
	comment := &readerModels.ChapterComment{
		ID:             primitive.NewObjectID().Hex(),
		ChapterID:      chapterID,
		BookID:         req.BookID,
		UserID:         userID.(string),
		Content:        req.Content,
		Rating:         0, // 段落评论通常不包含评分
		ParagraphIndex: req.ParagraphIndex,
		CharStart:      req.CharStart,
		CharEnd:        req.CharEnd,
		ReplyCount:     0,
		LikeCount:      0,
		IsVisible:      true,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 实际应用中应该保存到数据库
	shared.Success(c, http.StatusCreated, "评论成功", gin.H{
		"comment": comment,
		"message": fmt.Sprintf("段落 %d 评论发表成功", *req.ParagraphIndex),
	})
}

// GetChapterParagraphComments 获取章节所有段落评论概览
//
//	@Summary	获取章节所有段落评论概览
//	@Tags		阅读器-段落评论
//	@Param		chapterId	path	string	true	"章节ID"
//	@Success	200			{object}	shared.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/paragraph-comments [get]
func (api *ChapterCommentAPI) GetChapterParagraphComments(c *gin.Context) {
	chapterID := c.Param("chapterId")
	if chapterID == "" {
		shared.Error(c, http.StatusBadRequest, "章节ID不能为空", "章节ID不能为空")
		return
	}

	// 实际应用中应该：
	// 1. 从数据库查询该章节所有段落评论
	// 2. 按段落索引分组
	// 3. 返回每个段落的评论数量统计

	result := make(map[int]int)

	shared.Success(c, http.StatusOK, "获取成功", gin.H{
		"chapterId":      chapterID,
		"paragraphStats": result,
	})
}
