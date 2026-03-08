# Swagger API 文档导出说明

## 概述

本项目的 Swagger API 文档用于导出到 Postman、Apifox 等 API 测试工具。

## 文件位置

- **YAML格式**: `docs/api/swagger.yaml` (494 KB, 18,433 行)
- **JSON格式**: `docs/api/swagger.json` (696 KB, 28,124 行)

## 快速使用

### 生成 JSON 格式（推荐）

```bash
make swagger-convert
```

此命令会将现有的 `swagger.yaml` 转换为 `swagger.json`。

### 导入到 API 测试工具

#### Postman

1. 打开 Postman
2. 点击 `Import` 按钮
3. 选择 `Upload Files`
4. 选择 `docs/api/swagger.json` 文件
5. 点击导入

#### Apifox

1. 打开 Apifox 项目
2. 进入 `项目设置` -> `导入数据`
3. 选择 `OpenAPI/Swagger`
4. 选择 `docs/api/swagger.yaml` 或 `docs/api/swagger.json` 文件
5. 点击导入

#### Insomnia

1. 打开 Insomnia
2. 点击 `Import/Export`
3. 选择 `From File`
4. 选择 `docs/api/swagger.json` 文件
5. 点击导入

## 技术说明

### 为什么使用转换而不是重新生成？

由于 `swaggo/swag` 工具在处理复杂类型时存在已知的栈溢出 bug，当解析包含递归类型引用或复杂的嵌套结构时会崩溃。

```bash
# 以下命令会因栈溢出而失败：
swag init --generalInfo cmd/server/main.go --output docs --parseInternal

# 错误信息：
# runtime: goroutine stack exceeds 1000000000-byte limit
# fatal error: stack overflow
```

### 当前的解决方案

我们使用预先生成的 `swagger.yaml` 文件（历史遗留），并通过转换脚本生成 JSON 格式：

```bash
go run scripts/swagger_convert.go
```

### 更新 API 文档

由于无法使用 `swag init` 重新生成完整文档，当添加或修改 API 时，需要手动维护 swagger 注释：

1. **添加新的 API 端点**：
   ```go
   // @Summary 新增用户
   // @Description 创建新的用户账号
   // @Tags 用户管理
   // @Accept json
   // @Produce json
   // @Security BearerAuth
   // @Param request body dto.CreateUserRequest true "创建用户请求"
   // @Success 200 {object} response.APIResponse
   // @Failure 400 {object} response.APIResponse
   // @Router /api/v1/users [post]
   func (h *UserHandler) CreateUser(c *gin.Context) {
       // ...
   }
   ```

2. **添加新的 DTO 类型**：
   ```go
   // CreateUserRequest 创建用户请求（用于swagger文档）
   type CreateUserRequest struct {
       Username string `json:"username" binding:"required" example:"testuser"`
       Email    string `json:"email" binding:"required,email" example:"test@example.com"`
       Password string `json:"password" binding:"required,min=6" example:"password123"`
   }
   ```

3. **更新文档**：
   - 可以尝试使用 `swag init`（可能失败）
   - 或手动编辑 `docs/api/swagger.yaml` 文件
   - 然后运行 `make swagger-convert` 生成 JSON

## 相关文件

- `Makefile` - 包含 swagger 相关命令
- `scripts/swagger_convert.go` - YAML 到 JSON 转换脚本
- `docs/api/swagger.yaml` - Swagger YAML 格式源文件
- `docs/api/swagger.json` - Swagger JSON 格式输出文件

## Makefile 命令

```bash
# 查看所有可用命令
make help

# 转换 YAML 到 JSON（推荐）
make swagger-convert

# 尝试生成完整文档（可能因 bug 失败）
make swagger

# 仅生成 YAML 格式
make swagger-yaml

# 清理文档
make swagger-clean
```

## 已知问题

1. **swaggo 栈溢出**：无法使用 `swag init` 重新生成完整文档
   - 问题：GitHub issue: https://github.com/swaggo/swag/issues/xxx
   - 临时解决方案：使用现有 YAML + 转换脚本

2. **重复路由警告**：许多路由被声明多次，需要在路由注册时检查

## 替代方案

如果需要完整的 OpenAPI 文档生成，可以考虑：

1. **使用 go-swagger**：
   ```bash
   go install github.com/go-swagger/go-swagger/cmd/swagger@latest
   swagger generate spec -o docs/api/swagger.json --scan-models
   ```
   注意：go-swagger 使用不同的注解格式，与现有代码不兼容。

2. **手动维护 OpenAPI 规范**：
   直接编辑 `docs/api/swagger.yaml` 或 `docs/api/swagger.json`

3. **使用其他工具**：
   - oapi-codegen
   - kin-openapi
   - 手动编写 OpenAPI 3.0 规范

## 验证文档

可以使用在线工具验证生成的 Swagger 文档：

- https://editor.swagger.io/
- https://apifox.com/api-debugger/

## 联系方式

如有问题，请联系开发团队。
