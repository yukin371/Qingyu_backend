package writer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Qingyu_backend/api/v1/shared"
	"Qingyu_backend/service/interfaces"
)

// LocationApi 地点API处理器
type LocationApi struct {
	locationService interfaces.LocationService
}

// NewLocationApi 创建LocationApi实例
func NewLocationApi(locationService interfaces.LocationService) *LocationApi {
	return &LocationApi{
		locationService: locationService,
	}
}

// CreateLocation 创建地点
func (api *LocationApi) CreateLocation(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	userID := ""
	if uid, exists := c.Get("userId"); exists {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	location, err := api.locationService.Create(c.Request.Context(), projectID, userID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建地点失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", location)
}

// GetLocation 获取地点详情
func (api *LocationApi) GetLocation(c *gin.Context) {
	locationID := c.Param("locationId")
	projectID := c.Query("projectId")

	if locationID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "locationId和projectId不能为空")
		return
	}

	location, err := api.locationService.GetByID(c.Request.Context(), locationID, projectID)
	if err != nil {
		shared.Error(c, http.StatusNotFound, "地点不存在", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", location)
}

// ListLocations 获取项目地点列表
func (api *LocationApi) ListLocations(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	locations, err := api.locationService.List(c.Request.Context(), projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取地点列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", locations)
}

// GetLocationTree 获取地点层级树
func (api *LocationApi) GetLocationTree(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	tree, err := api.locationService.GetLocationTree(c.Request.Context(), projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取地点树失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", tree)
}

// UpdateLocation 更新地点
func (api *LocationApi) UpdateLocation(c *gin.Context) {
	locationID := c.Param("locationId")
	projectID := c.Query("projectId")

	if locationID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "locationId和projectId不能为空")
		return
	}

	var req interfaces.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	location, err := api.locationService.Update(c.Request.Context(), locationID, projectID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "更新地点失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "更新成功", location)
}

// DeleteLocation 删除地点
func (api *LocationApi) DeleteLocation(c *gin.Context) {
	locationID := c.Param("locationId")
	projectID := c.Query("projectId")

	if locationID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "locationId和projectId不能为空")
		return
	}

	err := api.locationService.Delete(c.Request.Context(), locationID, projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除地点失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}

// CreateLocationRelation 创建地点关系
func (api *LocationApi) CreateLocationRelation(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	var req interfaces.CreateLocationRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	relation, err := api.locationService.CreateRelation(c.Request.Context(), projectID, &req)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "创建关系失败", err.Error())
		return
	}

	shared.Success(c, http.StatusCreated, "创建成功", relation)
}

// ListLocationRelations 获取地点关系列表
func (api *LocationApi) ListLocationRelations(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		shared.Error(c, http.StatusBadRequest, "项目ID不能为空", "")
		return
	}

	locationID := c.Query("locationId")
	var locIDPtr *string
	if locationID != "" {
		locIDPtr = &locationID
	}

	relations, err := api.locationService.ListRelations(c.Request.Context(), projectID, locIDPtr)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "获取关系列表失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "获取成功", relations)
}

// DeleteLocationRelation 删除地点关系
func (api *LocationApi) DeleteLocationRelation(c *gin.Context) {
	relationID := c.Param("relationId")
	projectID := c.Query("projectId")

	if relationID == "" || projectID == "" {
		shared.Error(c, http.StatusBadRequest, "参数错误", "relationId和projectId不能为空")
		return
	}

	err := api.locationService.DeleteRelation(c.Request.Context(), relationID, projectID)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "删除关系失败", err.Error())
		return
	}

	shared.Success(c, http.StatusOK, "删除成功", nil)
}
