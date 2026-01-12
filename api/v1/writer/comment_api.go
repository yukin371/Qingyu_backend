package writer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"Qingyu_backend/api/v1/shared"
	writermodels "Qingyu_backend/models/writer"
	writerservice "Qingyu_backend/service/writer"
)

// CommentAPI 批注API
type CommentAPI struct {
	commentService writerservice.CommentService
}

// NewCommentAPI 创建批注API实例
func NewCommentAPI(commentService writerservice.CommentService) *CommentAPI {
	return &CommentAPI{
		commentService: commentService,
	}
}

// CreateComment 创建批注
//
//	@Summary		创建批注
//	@Description	在文档中添加批注
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"文档ID"
//	@Param			request	body		CreateCommentRequest	true	"批注信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/writer/documents/{id}/comments [post]
func (api *CommentAPI) CreateComment(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	userName := ""
	if name, exists := c.Get("userName"); exists {
		userName = name.(string)
	}

	// 构建批注
	comment, err := req.ToComment(documentID, userID.(string), userName)
	if err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	// 创建批注
	created, err := api.commentService.CreateComment(c.Request.Context(), comment)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "创建成功", created)
}

// GetComments 获取批注列表
//
//	@Summary		获取批注列表
//	@Description	获取文档的所有批注
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string	true	"文档ID"
//	@Param			resolved		query		bool	false	"是否包含已解决"
//	@Param			type			query		string	false	"批注类型"
//	@Param			page			query		int		false	"页码"	default(1)
//	@Param			size			query		int		false	"每页数量"	default(20)
//	@Success		200				{object}	shared.APIResponse
//	@Failure		400				{object}	shared.APIResponse
//	@Failure		401				{object}	shared.APIResponse
//	@Router			/api/v1/writer/documents/{id}/comments [get]
func (api *CommentAPI) GetComments(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 构建筛选条件
	filter := &writermodels.CommentFilter{
		Resolved: nil,
	}

	if resolved := c.Query("resolved"); resolved != "" {
		if resolved == "true" {
			trueVal := true
			filter.Resolved = &trueVal
		} else {
			falseVal := false
			filter.Resolved = &falseVal
		}
	}

	if commentType := c.Query("type"); commentType != "" {
		filter.Type = writermodels.CommentType(commentType)
	}

	comments, total, err := api.commentService.ListComments(c.Request.Context(), filter, page, size)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	result := gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"size":     size,
	}

	shared.Success(c, http.StatusOK, "获取成功", result)
}

// GetComment 获取批注详情
//
//	@Summary		获取批注详情
//	@Description	获取单个批注的详细信息
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"批注ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/{id} [get]
func (api *CommentAPI) GetComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "批注ID不能为空")
		return
	}

	comment, err := api.commentService.GetComment(c.Request.Context(), commentID)
	if err != nil {
		if err == writerservice.ErrCommentNotFound {
			shared.Error(c, http.StatusNotFound, "批注不存在", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", comment)
}

// UpdateComment 更新批注
//
//	@Summary		更新批注
//	@Description	更新批注内容
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"批注ID"
//	@Param			request	body		UpdateCommentRequest	true	"更新信息"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/{id} [put]
func (api *CommentAPI) UpdateComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "批注ID不能为空")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	comment := &writermodels.DocumentComment{
		Content:  req.Content,
		Type:     writermodels.CommentType(req.Type),
		Metadata: req.Metadata,
	}

	if err := api.commentService.UpdateComment(c.Request.Context(), commentID, comment); err != nil {
		if err == writerservice.ErrCommentNotFound {
			shared.Error(c, http.StatusNotFound, "批注不存在", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", nil)
}

// DeleteComment 删除批注
//
//	@Summary		删除批注
//	@Description	删除指定批注
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"批注ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/{id} [delete]
func (api *CommentAPI) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "批注ID不能为空")
		return
	}

	if err := api.commentService.DeleteComment(c.Request.Context(), commentID); err != nil {
		if err == writerservice.ErrCommentNotFound {
			shared.Error(c, http.StatusNotFound, "批注不存在", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "删除失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// ResolveComment 标记批注为已解决
//
//	@Summary		标记批注为已解决
//	@Description	标记批注为已解决
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"批注ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/{id}/resolve [post]
func (api *CommentAPI) ResolveComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "批注ID不能为空")
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	if err := api.commentService.ResolveComment(c.Request.Context(), commentID, userID.(string)); err != nil {
		if err == writerservice.ErrCommentNotFound {
			shared.Error(c, http.StatusNotFound, "批注不存在", err.Error())
			return
		}
		if err == writerservice.ErrCommentAlreadyResolved {
			shared.Error(c, http.StatusBadRequest, "批注已解决", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "操作失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "标记成功", nil)
}

// UnresolveComment 标记批注为未解决
//
//	@Summary		标记批注为未解决
//	@Description	标记批注为未解决
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"批注ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/{id}/unresolve [post]
func (api *CommentAPI) UnresolveComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "批注ID不能为空")
		return
	}

	if err := api.commentService.UnresolveComment(c.Request.Context(), commentID); err != nil {
		if err == writerservice.ErrCommentNotFound {
			shared.Error(c, http.StatusNotFound, "批注不存在", err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "操作失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "标记成功", nil)
}

// ReplyComment 回复批注
//
//	@Summary		回复批注
//	@Description	回复批注
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"批注ID"
//	@Param			request	body		ReplyCommentRequest	true	"回复内容"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Failure		404		{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/{id}/reply [post]
func (api *CommentAPI) ReplyComment(c *gin.Context) {
	parentID := c.Param("id")
	if parentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "批注ID不能为空")
		return
	}

	var req ReplyCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "未授权", "需要登录")
		return
	}

	userName := ""
	if name, exists := c.Get("userName"); exists {
		userName = name.(string)
	}

	reply, err := api.commentService.ReplyComment(c.Request.Context(), parentID, req.Content, userID.(string), userName)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "回复失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "回复成功", reply)
}

// GetCommentThread 获取批注线程
//
//	@Summary		获取批注线程
//	@Description	获取批注线程及其所有回复
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			threadId	path		string	true	"线程ID"
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Failure		404			{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/threads/{threadId} [get]
func (api *CommentAPI) GetCommentThread(c *gin.Context) {
	threadID := c.Param("threadId")
	if threadID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "线程ID不能为空")
		return
	}

	thread, err := api.commentService.GetCommentThread(c.Request.Context(), threadID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", thread)
}

// GetCommentStats 获取批注统计
//
//	@Summary		批注统计
//	@Description	获取文档的批注统计信息
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"文档ID"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/writer/documents/{id}/comments/stats [get]
func (api *CommentAPI) GetCommentStats(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	stats, err := api.commentService.GetCommentStats(c.Request.Context(), documentID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", stats)
}

// SearchComments 搜索批注
//
//	@Summary		搜索批注
//	@Description	在文档中搜索批注
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"文档ID"
//	@Param			keyword	query		string	true	"搜索关键词"
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			size		query		int		false	"每页数量"	default(20)
//	@Success		200			{object}	shared.APIResponse
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		401			{object}	shared.APIResponse
//	@Router			/api/v1/writer/documents/{id}/comments/search [get]
func (api *CommentAPI) SearchComments(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "文档ID不能为空")
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "关键词不能为空")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	comments, total, err := api.commentService.SearchComments(c.Request.Context(), keyword, documentID, page, size)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "搜索失败", err.Error())
		return
	}

	result := gin.H{
		"comments": comments,
		"total":    total,
		"page":     page,
		"size":     size,
		"keyword":  keyword,
	}

	shared.Success(c, http.StatusOK, "搜索成功", result)
}

// BatchDeleteComments 批量删除批注
//
//	@Summary		批量删除批注
//	@Description	批量删除批注
//	@Tags			Writer-Comment
//	@Accept			json
//	@Produce		json
//	@Param			request	body		BatchDeleteRequest	true	"批量删除请求"
//	@Success		200		{object}	shared.APIResponse
//	@Failure		400		{object}	shared.APIResponse
//	@Failure		401		{object}	shared.APIResponse
//	@Router			/api/v1/writer/comments/batch-delete [post]
func (api *CommentAPI) BatchDeleteComments(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	if len(req.CommentIDs) == 0 {
		shared.Error(c, http.StatusBadRequest, "参数错误", "批注ID列表不能为空")
		return
	}

	if err := api.commentService.BatchDeleteComments(c.Request.Context(), req.CommentIDs); err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// ============ 请求体结构 ============

// CreateCommentRequest 创建批注请求
type CreateCommentRequest struct {
	Content   string                       `json:"content" binding:"required"`
	Type      string                       `json:"type"`
	ChapterID string                       `json:"chapterId"`
	Position  writermodels.CommentPosition `json:"position" binding:"required"`
	ParentID  *string                      `json:"parentId"`
	Metadata  writermodels.CommentMetadata `json:"metadata"`
}

// ToComment 转换为批注模型
func (r *CreateCommentRequest) ToComment(documentID, userID, userName string) (*writermodels.DocumentComment, error) {
	docID, err := primitive.ObjectIDFromHex(documentID)
	if err != nil {
		return nil, err
	}

	var chapterID primitive.ObjectID
	if r.ChapterID != "" {
		chapterID, err = primitive.ObjectIDFromHex(r.ChapterID)
		if err != nil {
			return nil, err
		}
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var parentID *primitive.ObjectID
	if r.ParentID != nil {
		pid, err := primitive.ObjectIDFromHex(*r.ParentID)
		if err != nil {
			return nil, err
		}
		parentID = &pid
	}

	comment := &writermodels.DocumentComment{
		DocumentID: docID,
		ChapterID:  chapterID,
		UserID:     userObjID,
		UserName:   userName,
		Content:    r.Content,
		Type:       writermodels.CommentType(r.Type),
		Position:   r.Position,
		ParentID:   parentID,
		Metadata:   r.Metadata,
		Resolved:   false,
	}

	return comment, nil
}

// UpdateCommentRequest 更新批注请求
type UpdateCommentRequest struct {
	Content  string                       `json:"content"`
	Type     string                       `json:"type"`
	Metadata writermodels.CommentMetadata `json:"metadata"`
}

// ReplyCommentRequest 回复批注请求
type ReplyCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	CommentIDs []string `json:"commentIds" binding:"required"`
}
