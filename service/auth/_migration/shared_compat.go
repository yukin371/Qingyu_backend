//go:build ignore
// +build ignore

// Package auth 保留迁移历史说明。
//
// 本文件不会参与构建，仅用于记录：
// - 旧路径 Qingyu_backend/service/shared/auth 已删除；
// - 运行时代码必须使用 Qingyu_backend/service/auth；
// - 密码验证器归属 service/user；
// - 鉴权/权限中间件归属 internal/middleware/auth。
package auth
