package reader

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	bookstoreService "Qingyu_backend/service/bookstore"
	socialService "Qingyu_backend/service/social"

	readerModels "Qingyu_backend/models/reader"
	socialModels "Qingyu_backend/models/social"

	"Qingyu_backend/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterCommentAPI 章节评论API
type ChapterCommentAPI struct {
	commentService *socialService.CommentService
	chapterService bookstoreService.ChapterService
}

// NewChapterCommentAPI 创建章节评论API实例
func NewChapterCommentAPI() *ChapterCommentAPI {
	return &ChapterCommentAPI{}
}

// BindServices 绑定章节评论所需服务
func (api *ChapterCommentAPI) BindServices(commentService *socialService.CommentService, chapterService bookstoreService.ChapterService) {
	api.commentService = commentService
	api.chapterService = chapterService
}

type paragraphRichContent struct {
	ParagraphID    string `json:"paragraph_id,omitempty"`
	ParagraphIndex int    `json:"paragraph_index,omitempty"`
	ParagraphText  string `json:"paragraph_text,omitempty"`
}

func (api *ChapterCommentAPI) requireCommentService(c *gin.Context) bool {
	if api.commentService == nil {
		response.InternalError(c, fmt.Errorf("评论服务未初始化"))
		return false
	}
	return true
}

func (api *ChapterCommentAPI) getUserIdentity(c *gin.Context) (string, string, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return "", "", false
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		response.Unauthorized(c, "请先登录")
		return "", "", false
	}

	username, _ := c.Get("username")
	usernameStr, _ := username.(string)
	return userID, usernameStr, true
}

func (api *ChapterCommentAPI) loadParagraphRefs(c *gin.Context, chapterID, userID string) ([]map[string]interface{}, bool) {
	if api.chapterService == nil {
		response.InternalError(c, fmt.Errorf("章节服务未初始化"))
		return nil, false
	}

	paragraphs, err := api.chapterService.GetChapterParagraphs(c.Request.Context(), chapterID, userID)
	if err != nil {
		c.Error(err)
		return nil, false
	}

	refs := make([]map[string]interface{}, 0, len(paragraphs))
	for index, paragraph := range paragraphs {
		if paragraph == nil {
			continue
		}
		refs = append(refs, map[string]interface{}{
			"id":      paragraph.ID.Hex(),
			"index":   index,
			"content": paragraph.Content,
		})
	}
	return refs, true
}

func parseParagraphRichContent(raw interface{}) paragraphRichContent {
	meta := paragraphRichContent{}
	switch typed := raw.(type) {
	case map[string]interface{}:
		if value, ok := typed["paragraph_id"].(string); ok {
			meta.ParagraphID = value
		}
		if value, ok := typed["paragraph_index"].(int); ok {
			meta.ParagraphIndex = value
		}
		if value, ok := typed["paragraph_index"].(float64); ok {
			meta.ParagraphIndex = int(value)
		}
		if value, ok := typed["paragraph_text"].(string); ok {
			meta.ParagraphText = value
		}
	case primitive.M:
		if value, ok := typed["paragraph_id"].(string); ok {
			meta.ParagraphID = value
		}
		if value, ok := typed["paragraph_index"].(int32); ok {
			meta.ParagraphIndex = int(value)
		}
		if value, ok := typed["paragraph_index"].(int64); ok {
			meta.ParagraphIndex = int(value)
		}
		if value, ok := typed["paragraph_index"].(float64); ok {
			meta.ParagraphIndex = int(value)
		}
		if value, ok := typed["paragraph_text"].(string); ok {
			meta.ParagraphText = value
		}
	case primitive.D:
		return parseParagraphRichContent(typed.Map())
	}
	return meta
}

func toChapterCommentDTO(comment *socialModels.Comment) *readerModels.ChapterComment {
	if comment == nil {
		return nil
	}

	meta := parseParagraphRichContent(comment.RichContent)
	dto := &readerModels.ChapterComment{
		ID:         comment.ID.Hex(),
		ChapterID:  comment.ChapterID,
		BookID:     comment.BookID,
		UserID:     comment.AuthorID,
		Content:    comment.Content,
		Rating:     comment.Rating,
		ParentID:   comment.ParentID,
		RootID:     comment.RootID,
		ReplyCount: int(comment.ReplyCount),
		LikeCount:  int(comment.LikeCount),
		IsVisible:  comment.State == socialModels.CommentStateNormal,
		IsDeleted:  comment.State == socialModels.CommentStateDeleted,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}

	if meta.ParagraphID != "" {
		dto.ParagraphID = &meta.ParagraphID
	}
	if meta.ParagraphText != "" {
		dto.ParagraphText = &meta.ParagraphText
	}
	if comment.RichContent != nil {
		dto.ParagraphIndex = &meta.ParagraphIndex
	}

	if comment.AuthorSnapshot != nil {
		dto.UserSnapshot = &readerModels.CommentUserSnapshot{
			ID:       comment.AuthorSnapshot.ID,
			Username: comment.AuthorSnapshot.Username,
			Avatar:   comment.AuthorSnapshot.Avatar,
		}
	}

	return dto
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
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/comments [get]
// @Summary GetChapterComments 操作
// @Description TODO: 补充详细描述
// @Tags reader
// @Accept json
// @Produce json
// @Security Bearer
// @Param chapterId path string true "ChapterId"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /reader/{chapterId}/comments [get]

func (api *ChapterCommentAPI) GetChapterComments(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		response.BadRequest(c, "章节ID不能为空", "章节ID不能为空")
		return
	}

	// 验证章节ID格式
	if _, err := primitive.ObjectIDFromHex(chapterID); err != nil {
		response.BadRequest(c, "章节ID格式无效", "章节ID格式无效")
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

	userID, _, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	if _, err := api.chapterService.GetChapterByID(c.Request.Context(), chapterID); err != nil {
		c.Error(err)
		return
	}
	if _, ok := api.loadParagraphRefs(c, chapterID, userID); !ok {
		return
	}

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
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	parentID := ""
	parentIDPtr := &parentID
	if parentIdParam, exists := c.GetQuery("parentId"); exists {
		if parentIdParam != "" {
			if _, err := primitive.ObjectIDFromHex(parentIdParam); err != nil {
				response.BadRequest(c, "父评论ID格式无效", "父评论ID格式无效")
				return
			}
		}
		parentID = parentIdParam
	}

	targetType := socialModels.CommentTargetTypeChapter
	state := socialModels.CommentStateNormal
	filter := &socialModels.CommentFilter{
		TargetType: &targetType,
		TargetID:   &chapterID,
		ChapterID:  &chapterID,
		State:      &state,
		ParentID:   parentIDPtr,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
		Limit:      pageSize,
		Offset:     (page - 1) * pageSize,
	}

	result, total, err := api.commentService.ListCommentsByFilter(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	comments := make([]*readerModels.ChapterComment, 0, len(result))
	ratingTotal := 0
	ratingCount := 0
	for _, item := range result {
		dto := toChapterCommentDTO(item)
		if dto == nil {
			continue
		}
		comments = append(comments, dto)
		if item.Rating > 0 {
			ratingTotal += item.Rating
			ratingCount++
		}
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	avgRating := 0.0
	if ratingCount > 0 {
		avgRating = float64(ratingTotal) / float64(ratingCount)
	}

	response.Success(c, readerModels.ChapterCommentListResponse{
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
//	@Param		request		body object	true	"评论内容"
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/comments [post]
func (api *ChapterCommentAPI) CreateChapterComment(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}
	userID, username, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		response.BadRequest(c, "章节ID不能为空", "章节ID不能为空")
		return
	}

	// 验证章节ID格式
	if _, err := primitive.ObjectIDFromHex(chapterID); err != nil {
		response.BadRequest(c, "章节ID格式无效", "章节ID格式无效")
		return
	}

	var req readerModels.CreateChapterCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 确保ChapterID匹配
	req.ChapterID = chapterID

	chapter, err := api.chapterService.GetChapterByID(c.Request.Context(), chapterID)
	if err != nil {
		c.Error(err)
		return
	}
	if req.BookID == "" {
		req.BookID = chapter.BookID
	}

	// 验证评分
	if req.Rating < 0 || req.Rating > 5 {
		response.BadRequest(c, "评分必须在0-5之间", "评分必须在0-5之间")
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if req.ParentID != nil && *req.ParentID != "" {
		// 验证父评论ID格式
		if _, err := primitive.ObjectIDFromHex(*req.ParentID); err != nil {
			response.BadRequest(c, "父评论ID格式无效", "父评论ID格式无效")
			return
		}
	}

	comment := &socialModels.Comment{
		TargetType: socialModels.CommentTargetTypeChapter,
		TargetID:   chapterID,
		BookID:     req.BookID,
		ChapterID:  chapterID,
		AuthorID:   userID,
		Content:    strings.TrimSpace(req.Content),
		Rating:     req.Rating,
		AuthorSnapshot: &socialModels.CommentAuthorSnapshot{
			ID:       userID,
			Username: username,
		},
	}
	comment.ParentID = req.ParentID

	created, err := api.commentService.CreateTargetedComment(c.Request.Context(), comment)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, gin.H{
		"comment": toChapterCommentDTO(created),
		"message": "评论发表成功",
	})
}

// GetChapterComment 获取单条评论详情
//
//	@Summary	获取单条评论详情
//	@Tags		阅读器-章节评论
//	@Param		commentId	path	string	true	"评论ID"
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{commentId} [get]
func (api *ChapterCommentAPI) GetChapterComment(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		response.BadRequest(c, "评论ID不能为空", "评论ID不能为空")
		return
	}

	// 验证评论ID格式
	if _, err := primitive.ObjectIDFromHex(commentID); err != nil {
		response.BadRequest(c, "评论ID格式无效", "评论ID格式无效")
		return
	}

	comment, err := api.commentService.GetCommentDetail(c.Request.Context(), commentID)
	if err != nil {
		response.NotFound(c, "评论不存在")
		return
	}
	response.Success(c, gin.H{"comment": toChapterCommentDTO(comment)})
}

// UpdateChapterComment 更新章节评论
//
//	@Summary	更新章节评论
//	@Tags		阅读器-章节评论
//	@Param		commentId	path	string								true	"评论ID"
//	@Param		request		body object	true	"更新内容"
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{commentId} [put]
func (api *ChapterCommentAPI) UpdateChapterComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		response.BadRequest(c, "评论ID不能为空", "评论ID不能为空")
		return
	}

	var req readerModels.UpdateChapterCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 验证评分
	if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 5) {
		response.BadRequest(c, "评分必须在0-5之间", "评分必须在0-5之间")
		return
	}

	// 实际应用中应该：
	// 1. 从数据库获取评论
	// 2. 验证用户是否为评论作者
	// 3. 检查是否在可编辑时间内（30分钟）
	// 4. 更新评论内容
	// 5. 如果评分改变，更新章节/书籍评分统计

	response.Success(c, gin.H{
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
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{commentId} [delete]
func (api *ChapterCommentAPI) DeleteChapterComment(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}
	userID, _, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		response.BadRequest(c, "评论ID不能为空", "评论ID不能为空")
		return
	}

	if err := api.commentService.DeleteComment(c.Request.Context(), userID, commentID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
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
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{commentId}/like [post]
func (api *ChapterCommentAPI) LikeChapterComment(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}
	userID, _, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		response.BadRequest(c, "评论ID不能为空", "评论ID不能为空")
		return
	}

	if err := api.commentService.LikeComment(c.Request.Context(), userID, commentID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
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
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/comments/{commentId}/like [delete]
func (api *ChapterCommentAPI) UnlikeChapterComment(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}
	userID, _, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	commentID := c.Param("commentId")
	if commentID == "" {
		response.BadRequest(c, "评论ID不能为空", "评论ID不能为空")
		return
	}

	if err := api.commentService.UnlikeComment(c.Request.Context(), userID, commentID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
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
//	@Success	200				{object}	response.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/paragraphs/{paragraphIndex}/comments [get]
func (api *ChapterCommentAPI) GetParagraphComments(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		response.BadRequest(c, "章节ID不能为空", "章节ID不能为空")
		return
	}

	paragraphIndexStr := c.Param("paragraphIndex")
	if paragraphIndexStr == "" {
		response.BadRequest(c, "段落索引不能为空", "段落索引不能为空")
		return
	}

	paragraphIndex, err := strconv.Atoi(paragraphIndexStr)
	if err != nil || paragraphIndex < 0 {
		response.BadRequest(c, "段落索引格式无效", "段落索引格式无效")
		return
	}

	userID, _, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	paragraphRefs, ok := api.loadParagraphRefs(c, chapterID, userID)
	if !ok {
		return
	}

	var selected paragraphRichContent
	found := false
	for _, ref := range paragraphRefs {
		refIndex, _ := ref["index"].(int)
		if refIndex != paragraphIndex {
			continue
		}
		selected.ParagraphIndex = paragraphIndex
		selected.ParagraphID, _ = ref["id"].(string)
		selected.ParagraphText, _ = ref["content"].(string)
		found = true
		break
	}
	if !found {
		response.NotFound(c, "段落不存在")
		return
	}

	targetType := socialModels.CommentTargetTypeChapter
	state := socialModels.CommentStateNormal
	parentID := ""
	filter := &socialModels.CommentFilter{
		TargetType:     &targetType,
		TargetID:       &chapterID,
		ChapterID:      &chapterID,
		ParagraphIndex: &paragraphIndex,
		State:          &state,
		ParentID:       &parentID,
		SortBy:         "created_at",
		SortOrder:      "asc",
		Limit:          200,
	}
	result, _, err := api.commentService.ListCommentsByFilter(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	comments := make([]*readerModels.ChapterComment, 0, len(result))
	for _, item := range result {
		dto := toChapterCommentDTO(item)
		if dto == nil {
			continue
		}
		comments = append(comments, dto)
		replies, _, err := api.commentService.GetCommentReplies(c.Request.Context(), item.ID.Hex(), 1, 200)
		if err != nil {
			c.Error(err)
			return
		}
		for _, reply := range replies {
			replyDTO := toChapterCommentDTO(reply)
			if replyDTO != nil {
				comments = append(comments, replyDTO)
			}
		}
	}

	response.Success(c, readerModels.ParagraphCommentResponse{
		ParagraphIndex: paragraphIndex,
		ParagraphID:    selected.ParagraphID,
		ParagraphText:  selected.ParagraphText,
		CommentCount:   len(comments),
		Comments:       comments,
	})
}

// CreateParagraphComment 发表段落级评论
//
//	@Summary	发表段落级评论
//	@Tags		阅读器-段落评论
//	@Param		chapterId	path	string								true	"章节ID"
//	@Param		request		body object	true	"评论内容"
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/paragraph-comments [post]
func (api *ChapterCommentAPI) CreateParagraphComment(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}
	userID, username, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		response.BadRequest(c, "章节ID不能为空", "章节ID不能为空")
		return
	}

	var req readerModels.CreateChapterCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	// 段落级评论必须指定段落索引
	if req.ParagraphIndex == nil {
		response.BadRequest(c, "段落索引不能为空", "段落索引不能为空")
		return
	}

	chapter, err := api.chapterService.GetChapterByID(c.Request.Context(), chapterID)
	if err != nil {
		c.Error(err)
		return
	}
	if req.BookID == "" {
		req.BookID = chapter.BookID
	}

	paragraphRefs, ok := api.loadParagraphRefs(c, chapterID, userID)
	if !ok {
		return
	}

	var selected paragraphRichContent
	found := false
	for _, ref := range paragraphRefs {
		refIndex, _ := ref["index"].(int)
		if refIndex != *req.ParagraphIndex {
			continue
		}
		selected.ParagraphIndex = *req.ParagraphIndex
		selected.ParagraphID, _ = ref["id"].(string)
		selected.ParagraphText, _ = ref["content"].(string)
		found = true
		break
	}
	if !found {
		response.NotFound(c, "段落不存在")
		return
	}
	if req.ParagraphID != nil && *req.ParagraphID != "" && *req.ParagraphID != selected.ParagraphID {
		response.BadRequest(c, "段落ID与索引不匹配", "段落ID与索引不匹配")
		return
	}
	if req.ParentID != nil && *req.ParentID != "" {
		if _, err := primitive.ObjectIDFromHex(*req.ParentID); err != nil {
			response.BadRequest(c, "父评论ID格式无效", "父评论ID格式无效")
			return
		}
	}

	richContent := map[string]interface{}{
		"paragraph_id":    selected.ParagraphID,
		"paragraph_index": selected.ParagraphIndex,
	}
	if selected.ParagraphText != "" {
		richContent["paragraph_text"] = selected.ParagraphText
	}

	comment := &socialModels.Comment{
		TargetType:  socialModels.CommentTargetTypeChapter,
		TargetID:    chapterID,
		BookID:      req.BookID,
		ChapterID:   chapterID,
		AuthorID:    userID,
		Content:     strings.TrimSpace(req.Content),
		Rating:      0,
		RichContent: richContent,
		AuthorSnapshot: &socialModels.CommentAuthorSnapshot{
			ID:       userID,
			Username: username,
		},
	}
	comment.ParentID = req.ParentID

	created, err := api.commentService.CreateTargetedComment(c.Request.Context(), comment)
	if err != nil {
		c.Error(err)
		return
	}
	response.Created(c, gin.H{
		"comment": toChapterCommentDTO(created),
		"message": fmt.Sprintf("段落 %d 评论发表成功", *req.ParagraphIndex),
	})
}

// GetChapterParagraphComments 获取章节所有段落评论概览
//
//	@Summary	获取章节所有段落评论概览
//	@Tags		阅读器-段落评论
//	@Param		chapterId	path	string	true	"章节ID"
//	@Success	200			{object}	response.APIResponse
//	@Router		/api/v1/reader/chapters/{chapterId}/paragraph-comments [get]
func (api *ChapterCommentAPI) GetChapterParagraphComments(c *gin.Context) {
	if !api.requireCommentService(c) {
		return
	}

	chapterID := c.Param("chapterId")
	if chapterID == "" {
		response.BadRequest(c, "章节ID不能为空", "章节ID不能为空")
		return
	}

	userID, _, ok := api.getUserIdentity(c)
	if !ok {
		return
	}

	paragraphRefs, ok := api.loadParagraphRefs(c, chapterID, userID)
	if !ok {
		return
	}

	targetType := socialModels.CommentTargetTypeChapter
	state := socialModels.CommentStateNormal
	filter := &socialModels.CommentFilter{
		TargetType: &targetType,
		TargetID:   &chapterID,
		ChapterID:  &chapterID,
		State:      &state,
		SortBy:     "created_at",
		SortOrder:  "desc",
		Limit:      500,
	}
	result, _, err := api.commentService.ListCommentsByFilter(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	type paragraphSummary struct {
		ParagraphID    string
		ParagraphIndex int
		CommentCount   int
		LatestComment  *struct {
			Content  string `json:"content"`
			Username string `json:"username"`
			Time     string `json:"time"`
		}
	}
	summaryByParagraph := make(map[string]*paragraphSummary)
	for _, ref := range paragraphRefs {
		refID, _ := ref["id"].(string)
		refIndex, _ := ref["index"].(int)
		summaryByParagraph[refID] = &paragraphSummary{
			ParagraphID:    refID,
			ParagraphIndex: refIndex,
		}
	}
	for _, item := range result {
		meta := parseParagraphRichContent(item.RichContent)
		if meta.ParagraphID == "" {
			continue
		}
		summary, exists := summaryByParagraph[meta.ParagraphID]
		if !exists {
			summary = &paragraphSummary{
				ParagraphID:    meta.ParagraphID,
				ParagraphIndex: meta.ParagraphIndex,
			}
			summaryByParagraph[meta.ParagraphID] = summary
		}
		summary.CommentCount++
		if summary.LatestComment == nil {
			username := ""
			if item.AuthorSnapshot != nil {
				username = item.AuthorSnapshot.Username
			}
			summary.LatestComment = &struct {
				Content  string `json:"content"`
				Username string `json:"username"`
				Time     string `json:"time"`
			}{
				Content:  item.Content,
				Username: username,
				Time:     item.CreatedAt.Format(time.RFC3339),
			}
		}
	}

	items := make([]readerModels.ParagraphCommentSummaryItem, 0, len(summaryByParagraph))
	for _, ref := range paragraphRefs {
		refID, _ := ref["id"].(string)
		summary := summaryByParagraph[refID]
		if summary == nil || summary.CommentCount == 0 {
			continue
		}
		items = append(items, readerModels.ParagraphCommentSummaryItem{
			ParagraphID:    summary.ParagraphID,
			ParagraphIndex: summary.ParagraphIndex,
			CommentCount:   summary.CommentCount,
			LatestComment:  summary.LatestComment,
		})
	}

	response.Success(c, gin.H{
		"chapterId":      chapterID,
		"paragraphStats": items,
	})
}
