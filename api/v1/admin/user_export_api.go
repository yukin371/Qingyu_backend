package admin

import (
	"encoding/csv"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"

	"Qingyu_backend/models/users"
	"Qingyu_backend/pkg/response"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
	adminservice "Qingyu_backend/service/admin"
)

// UserExportAPI 用户导出API处理器
type UserExportAPI struct {
	userAdminService adminservice.UserAdminService
}

// NewUserExportAPI 创建用户导出API实例
func NewUserExportAPI(userAdminService adminservice.UserAdminService) *UserExportAPI {
	return &UserExportAPI{
		userAdminService: userAdminService,
	}
}

// ExportUsers 导出用户数据
//
//	@Summary		导出用户数据
//	@Description	管理员导出用户数据，支持CSV和Excel格式
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		csv,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			format	query		string	false	"导出格式"	Enums(csv,xlsx)
//	@Param			role		query		string	false	"角色筛选"
//	@Param			status		query		string	false	"状态筛选"
//	@Param			keyword		query		string	false	"关键词搜索"
//	@Success		200			{file}		file
//	@Failure		400			{object}	shared.APIResponse
//	@Failure		500			{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/export [get]
func (api *UserExportAPI) ExportUsers(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")

	// 验证格式
	if format != "csv" && format != "xlsx" {
		response.BadRequest(c, "参数错误", "不支持的导出格式，仅支持csv和xlsx")
		return
	}

	// 构建筛选条件
	filter := &adminrepo.UserFilter{
		Keyword: c.Query("keyword"),
		Role:    c.Query("role"),
	}

	if status := c.Query("status"); status != "" {
		filter.Status = users.UserStatus(status)
	}

	// 获取用户列表（分页大小设为10000以支持大批量导出）
	userList, _, err := api.userAdminService.GetUserList(c.Request.Context(), filter, 1, 10000)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	switch format {
	case "csv":
		api.exportToCSV(c, userList)
	case "xlsx":
		api.exportToExcel(c, userList)
	}
}

// GetUserExportTemplate 获取导出模板
//
//	@Summary		获取导出模板
//	@Description	获取用户导出的字段模板
//	@Tags			Admin-User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shared.APIResponse
//	@Router			/api/v1/admin/users/export/template [get]
func (api *UserExportAPI) GetUserExportTemplate(c *gin.Context) {
	template := []map[string]string{
		{"field": "username", "label": "用户名", "required": "true"},
		{"field": "email", "label": "邮箱", "required": "false"},
		{"field": "phone", "label": "手机号", "required": "false"},
		{"field": "nickname", "label": "昵称", "required": "false"},
		{"field": "bio", "label": "简介", "required": "false"},
		{"field": "role", "label": "角色", "required": "false"},
		{"field": "status", "label": "状态", "required": "false"},
	}

	response.Success(c, template)
}

// exportToCSV 导出为CSV格式
func (api *UserExportAPI) exportToCSV(c *gin.Context, userList []*users.User) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=users.csv")

	// 创建CSV写入器
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// 写入表头
	headers := []string{"ID", "用户名", "邮箱", "手机号", "昵称", "角色", "状态", "注册时间"}
	writer.Write(headers)

	// 写入数据
	for _, user := range userList {
		role := getPrimaryRole(user.Roles)
		record := []string{
			user.ID.Hex(),
			user.Username,
			user.Email,
			user.Phone,
			user.Nickname,
			role,
			string(user.Status),
			user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}
}

// exportToExcel 导出为Excel格式
func (api *UserExportAPI) exportToExcel(c *gin.Context, userList []*users.User) {
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=users.xlsx")

	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := "用户列表"

	// 设置表头
	headers := []string{"ID", "用户名", "邮箱", "手机号", "昵称", "角色", "状态", "注册时间"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for i, user := range userList {
		row := i + 2
		role := getPrimaryRole(user.Roles)

		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), user.ID.Hex())
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), user.Username)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), user.Email)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), user.Phone)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), user.Nickname)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), role)
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), string(user.Status))
		f.SetCellValue(sheetName, "H"+strconv.Itoa(row), user.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// 写入响应
	f.Write(c.Writer)
}

// getPrimaryRole 获取主要角色（优先级：admin > author > reader）
func getPrimaryRole(roles []string) string {
	if len(roles) == 0 {
		return ""
	}

	priority := map[string]int{
		"admin":  3,
		"author": 2,
		"reader": 1,
	}

	primaryRole := roles[0]
	maxPriority := 0

	for _, role := range roles {
		if p, ok := priority[role]; ok && p > maxPriority {
			maxPriority = p
			primaryRole = role
		}
	}

	return primaryRole
}
