# API迁移常见问题FAQ

> **版本**: v1.0
> **更新日期**: 2026-01-29
> **维护者**: Backend Team

## 📋 目录

1. [迁移相关](#迁移相关)
2. [代码变更](#代码变更)
3. [测试相关](#测试相关)
4. [错误处理](#错误处理)
5. [特殊场景](#特殊场景)
6. [工具和流程](#工具和流程)

---

## 迁移相关

### Q1: 为什么要从shared包迁移到response包？

**A**: 统一响应格式，简化API调用，规范错误码。

**收益**:
- 代码更简洁（4参数→2参数）
- 错误码统一（6位→4位）
- 响应格式一致
- 依赖更少

### Q2: 迁移会影响现有功能吗？

**A**: 不会。迁移只是改变响应调用的方式，不改变业务逻辑。

**保证**:
- 响应数据结构兼容
- HTTP状态码一致
- 错误信息完整
- 测试全覆盖

### Q3: 迁移需要多长时间？

**A**: 取决于文件复杂度。

**参考**:
- 简单文件（10-20次调用）: 20-30分钟
- 中等文件（20-40次调用）: 30-45分钟
- 复杂文件（40+次调用）: 1小时+
- Writer模块总计: 预计12.5小时（1.5-2天）

### Q4: 可以部分迁移吗？

**A**: 可以，但不建议。

**原因**:
- 部分迁移导致代码不一致
- 增加维护成本
- 容易遗漏

**建议**: 按文件完整迁移，一次完成一个文件。

### Q5: 迁移后发现错误怎么办？

**A**: 立即回滚，分析问题，重新迁移。

**步骤**:
1. Git revert或回退到上一个commit
2. 分析错误原因
3. 修复问题
4. 重新迁移

---

## 代码变更

### Q6: 如何替换shared.Error调用？

**A**: 根据HTTP状态码选择对应的response函数。

```go
// 400 Bad Request
shared.Error(c, http.StatusBadRequest, "参数错误", err.Error())
→ response.BadRequest(c, "参数错误", err.Error())

// 401 Unauthorized
shared.Error(c, http.StatusUnauthorized, "未授权", "请先登录")
→ response.Unauthorized(c, "请先登录")

// 403 Forbidden
shared.Error(c, http.StatusForbidden, "禁止访问", "无权限")
→ response.Forbidden(c, "无权限")

// 404 Not Found
shared.Error(c, http.StatusNotFound, "未找到", "资源不存在")
→ response.NotFound(c, "资源不存在")

// 409 Conflict
shared.Error(c, http.StatusConflict, "版本冲突", "文档已被修改")
→ response.Conflict(c, "版本冲突", "文档已被修改")

// 500 Internal Error
shared.Error(c, http.StatusInternalServerError, "服务器错误", err.Error())
→ response.InternalError(c, err)
```

### Q7: 如何替换shared.Success调用？

**A**: 根据操作类型选择Success或Created。

```go
// 200 OK
shared.Success(c, http.StatusOK, "获取成功", data)
→ response.Success(c, data)

// 201 Created
shared.Success(c, http.StatusCreated, "创建成功", data)
→ response.Created(c, data)
```

### Q8: 如何处理shared.ValidationError？

**A**: 替换为response.BadRequest。

```go
// 旧代码
shared.ValidationError(c, err)

// 新代码
response.BadRequest(c, "参数错误", err.Error())
```

### Q9: 可以保留消息参数吗？

**A**: 可以，但通常不需要。

**说明**:
- response包会自动设置合适的消息
- 自定义消息可以通过参数传递
- 建议让response包自动处理

### Q10: 如何清理导入依赖？

**A**: 移除shared和net/http（除非WebSocket）。

```go
// 移除
import (
    "net/http"  // 移除（WebSocket除外）
    "Qingyu_backend/api/v1/shared"  // 移除
)

// 添加
import (
    "Qingyu_backend/pkg/response"  // 添加
)
```

### Q11: Swagger注释如何更新？

**A**: 替换shared.APIResponse为response.APIResponse。

```go
// 旧注释
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse

// 新注释
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
```

### Q12: 如何验证迁移完整性？

**A**: 使用grep搜索残留的shared调用。

```bash
# 搜索shared.Error
grep -r "shared\.Error" api/v1/writer

# 搜索shared.Success
grep -r "shared\.Success" api/v1/writer

# 搜索shared.ValidationError
grep -r "shared\.ValidationError" api/v1/writer

# 搜索shared包导入
grep -r "Qingyu_backend/api/v1/shared" api/v1/writer
```

---

## 测试相关

### Q13: 迁移后测试失败怎么办？

**A**: 检查响应格式和错误码是否匹配。

**常见问题**:
1. 响应结构变化
2. 错误码变化（6位→4位）
3. 时间戳格式变化（秒→毫秒）

**解决**:
```go
// 旧测试断言
assert.Equal(t, 100001, response.Code)

// 新测试断言
assert.Equal(t, 1001, response.Code)
```

### Q14: 如何编写新的单元测试？

**A**: 参考迁移指南中的示例代码。

**结构**:
```go
func TestAPI_GetXxx(t *testing.T) {
    // 1. 设置测试环境
    gin.SetMode(gin.TestMode)
    router := gin.New()
    api := NewXxxAPI()
    router.GET("/xxx", api.GetXxx)

    // 2. 创建请求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/xxx", nil)

    // 3. 执行请求
    router.ServeHTTP(w, req)

    // 4. 断言响应
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "\"code\":0")
}
```

### Q15: 如何处理Mock测试？

**A**: 使用testify/mock或类似的Mock框架。

**示例**:
```go
func TestAPI_GetXxx_Mock(t *testing.T) {
    // 创建Mock服务
    mockService := new(MockXxxService)
    mockService.On("GetXxx", mock.Anything, "id123").Return(&Xxx{}, nil)

    // 创建API并注入Mock
    api := NewXxxAPI(mockService)

    // 执行测试
    // ...
}
```

### Q16: 如何提高测试覆盖率？

**A**: 覆盖所有分支和错误场景。

**策略**:
1. 成功路径测试
2. 参数错误测试
3. 未授权测试
4. 资源不存在测试
5. 服务器错误测试

---

## 错误处理

### Q17: 如何处理版本冲突？

**A**: 使用response.Conflict。

```go
if err.Error() == "版本冲突" {
    response.Conflict(c, "版本冲突", "文档已被其他用户修改，请刷新后重试")
    return
}
```

### Q18: 如何处理自定义错误？

**A**: 包装错误并返回适当的响应。

```go
// 定义自定义错误
var ErrNotFound = errors.New("文档不存在")

// 使用
if err != nil {
    if errors.Is(err, ErrNotFound) {
        response.NotFound(c, "文档不存在")
        return
    }
    response.InternalError(c, err)
    return
}
```

### Q19: 如何处理第三方服务错误？

**A**: 转换为内部错误。

```go
resp, err := thirdPartyClient.Call()
if err != nil {
    response.InternalError(c, fmt.Errorf("第三方服务错误: %w", err))
    return
}
```

### Q20: 错误码如何选择？

**A**: 根据错误类型选择对应的错误码。

```go
// 参数错误 → 1001
response.BadRequest(c, "参数错误", details)  // Code: 1001

// 未授权 → 1002
response.Unauthorized(c, "请先登录")  // Code: 1002

// 禁止访问 → 1003
response.Forbidden(c, "无权限")  // Code: 1003

// 资源不存在 → 1004
response.NotFound(c, "资源不存在")  // Code: 1004

// 版本冲突 → 1006
response.Conflict(c, "版本冲突", details)  // Code: 1006

// 服务器错误 → 5000
response.InternalError(c, err)  // Code: 5000
```

---

## 特殊场景

### Q21: WebSocket如何处理？

**A**: 保留net/http导入，WebSocket部分不变。

```go
import (
    "net/http"  // 保留，WebSocket需要
    "Qingyu_backend/pkg/response"
)

// WebSocket升级不需要修改
upgrader := websocket.Upgrader{}
conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
if err != nil {
    response.InternalError(c, err)
    return
}
```

### Q22: 文件下载如何处理？

**A**: 文件下载部分不变，错误处理迁移。

```go
// 文件下载不需要修改
c.FileAttachment(filePath, fileName)

// 错误处理需要迁移
if err != nil {
    response.InternalError(c, err)
    return
}
```

### Q23: 批量操作如何处理？

**A**: 提交后立即返回，异步执行。

```go
// 提交批量操作
response.Success(c, gin.H{
    "batchId": batchOp.ID.Hex(),
    "status": "submitted",
})

// 异步执行
go func() {
    api.batchOpSvc.Execute(ctx, batchId)
}()
```

### Q24: 分页响应如何处理？

**A**: 使用response.Paginated或自定义结构。

```go
// 推荐：使用Paginated
response.Paginated(c, data, total, page, pageSize, "获取成功")

// 或者：自定义结构
response.Success(c, gin.H{
    "list": data,
    "total": total,
    "page": page,
    "pageSize": pageSize,
})
```

### Q25: 流式响应如何处理？

**A**: 流式响应不使用response包。

```go
// 流式响应直接使用gin.Context
c.Stream(func(w io.Writer) bool {
    // 写入流数据
    return true
})
```

---

## 工具和流程

### Q26: 有自动化工具可以辅助迁移吗？

**A**: 可以使用grep/sed批量替换，但要仔细检查。

**示例**:
```bash
# 批量替换shared.Error → response.BadRequest
sed -i 's/shared\.Error(c, http\.StatusBadRequest,/response.BadRequest(c,/g' xxx_api.go

# 注意：需要仔细检查每个替换
```

### Q27: 如何创建迁移分支？

**A**: 使用git checkout -b创建feature分支。

```bash
git checkout -b feature/block8-writer-migration
```

### Q28: 如何提交迁移代码？

**A**: 使用规范的commit信息。

```bash
git add api/v1/writer/xxx_api.go
git commit -m "feat(api): migrate xxx_api to new response package

- Replace all shared.Error calls with response functions
- Replace all shared.Success calls with response functions
- Remove HTTP status code parameters
- Update Swagger annotations
- Clean up imports"
```

### Q29: 如何创建PR？

**A**: 使用gh CLI或GitHub网页。

```bash
# 推送到远程
git push origin feature/block8-writer-migration

# 创建PR
gh pr create --title "[Block 8] API迁移 - xxx模块" --body "PR描述..."
```

### Q30: 如何验证PR？

**A**: 等待CI检查通过，代码审查通过。

**检查项**:
- ✅ CI测试全部通过
- ✅ 代码审查通过
- ✅ 无shared包残留
- ✅ Swagger文档完整

---

## 更多帮助

### 获取更多帮助

1. **查阅迁移指南**: [api-migration-guide.md](api-migration-guide.md)
2. **查看检查清单**: [api-migration-checklist.md](api-migration-checklist.md)
3. **参考Block 7**: [Block 7进展报告](../plans/2026-01-28-block7-api-standardization-progress.md)
4. **分析Writer模块**: [Writer模块预分析报告](../analysis/2026-01-29-writer-migration-analysis.md)

### 联系方式

- **问题反馈**: GitHub Issues
- **技术讨论**: 团队会议
- **紧急问题**: 联系Tech Lead

---

**FAQ版本**: v1.0
**最后更新**: 2026-01-29
**维护者**: Backend Team
