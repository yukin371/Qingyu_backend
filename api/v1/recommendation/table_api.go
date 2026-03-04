package recommendation

import (
	reco "Qingyu_backend/models/recommendation"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/pkg/response"
)

type upsertAutoTableRequest struct {
	Items []reco.TableItem `json:"items"`
}

type createManualTableRequest struct {
	Name   string           `json:"name" binding:"required"`
	Period string           `json:"period"`
	Items  []reco.TableItem `json:"items"`
}

type updateManualTableRequest struct {
	Name   string           `json:"name"`
	Period string           `json:"period"`
	Status reco.TableStatus `json:"status"`
	Items  []reco.TableItem `json:"items"`
}

// ListTables 获取推荐榜列表
func (api *RecommendationAPI) ListTables(c *gin.Context) {
	if api.tableService == nil {
		response.InternalError(c, errors.New("table service not initialized"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var tableType *reco.TableType
	if t := c.Query("type"); t != "" {
		if !reco.IsValidTableType(t) {
			response.BadRequest(c, "参数错误", "invalid table type")
			return
		}
		tt := reco.TableType(t)
		tableType = &tt
	}

	var source *reco.TableSource
	if s := c.Query("source"); s != "" {
		if !reco.IsValidTableSource(s) {
			response.BadRequest(c, "参数错误", "invalid source")
			return
		}
		ss := reco.TableSource(s)
		source = &ss
	}

	tables, total, err := api.tableService.ListTables(c.Request.Context(), tableType, source, page, size)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{
		"tables": tables,
		"total":  total,
		"page":   page,
		"size":   size,
	})
}

// GetTable 获取单个推荐榜
func (api *RecommendationAPI) GetTable(c *gin.Context) {
	if api.tableService == nil {
		response.InternalError(c, errors.New("table service not initialized"))
		return
	}

	id := c.Param("id")
	table, err := api.tableService.GetTable(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	if table == nil {
		response.NotFound(c, "table not found")
		return
	}
	response.Success(c, table)
}

// UpsertAutoTable 覆盖自动榜（周榜/月榜/月票榜）
func (api *RecommendationAPI) UpsertAutoTable(c *gin.Context) {
	if api.tableService == nil {
		response.InternalError(c, errors.New("table service not initialized"))
		return
	}

	tableType := c.Param("tableType")
	period := c.Param("period")
	if !reco.IsValidTableType(tableType) {
		response.BadRequest(c, "参数错误", "invalid table type")
		return
	}

	var req upsertAutoTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	updatedBy, _ := c.Get("user_id")
	updatedByStr, _ := updatedBy.(string)

	if err := api.tableService.UpsertAutoTable(c.Request.Context(), reco.TableType(tableType), period, req.Items, updatedByStr); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// CreateManualTable 创建手动推荐榜
func (api *RecommendationAPI) CreateManualTable(c *gin.Context) {
	if api.tableService == nil {
		response.InternalError(c, errors.New("table service not initialized"))
		return
	}

	var req createManualTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	updatedBy, _ := c.Get("user_id")
	updatedByStr, _ := updatedBy.(string)
	if err := api.tableService.CreateManualTable(c.Request.Context(), req.Name, req.Period, req.Items, updatedByStr); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// UpdateManualTable 更新手动推荐榜
func (api *RecommendationAPI) UpdateManualTable(c *gin.Context) {
	if api.tableService == nil {
		response.InternalError(c, errors.New("table service not initialized"))
		return
	}

	id := c.Param("id")
	var req updateManualTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	updatedBy, _ := c.Get("user_id")
	updatedByStr, _ := updatedBy.(string)
	if err := api.tableService.UpdateManualTable(c.Request.Context(), id, req.Name, req.Period, req.Items, req.Status, updatedByStr); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}

// DeleteTable 删除推荐榜
func (api *RecommendationAPI) DeleteTable(c *gin.Context) {
	if api.tableService == nil {
		response.InternalError(c, errors.New("table service not initialized"))
		return
	}

	id := c.Param("id")
	if err := api.tableService.DeleteTable(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	response.Success(c, nil)
}
