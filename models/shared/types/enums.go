package types

import (
	"fmt"
)

// UserRole 用户角色
type UserRole string

const (
	RoleReader UserRole = "reader"
	RoleAuthor UserRole = "author"
	RoleAdmin  UserRole = "admin"
)

// AllUserRoles 所有有效角色
var AllUserRoles = []UserRole{RoleReader, RoleAuthor, RoleAdmin}

// IsValid 检查角色是否有效
func (r UserRole) IsValid() bool {
	switch r {
	case RoleReader, RoleAuthor, RoleAdmin:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (r UserRole) String() string {
	return string(r)
}

// ParseUserRole 从字符串解析角色
func ParseUserRole(s string) (UserRole, error) {
	role := UserRole(s)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid user role: %s", s)
	}
	return role, nil
}

// CanPublish 是否可以发布内容
func (r UserRole) CanPublish() bool {
	return r == RoleAuthor || r == RoleAdmin
}

// CanModerate 是否可以管理内容
func (r UserRole) CanModerate() bool {
	return r == RoleAdmin
}

// PageMode 阅读翻页模式
type PageMode string

const (
	PageModeScroll    PageMode = "scroll"
	PageModePaginate  PageMode = "paginate"
)

// AllPageModes 所有有效模式
var AllPageModes = []PageMode{PageModeScroll, PageModePaginate}

// IsValid 检查模式是否有效
func (m PageMode) IsValid() bool {
	switch m {
	case PageModeScroll, PageModePaginate:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (m PageMode) String() string {
	return string(m)
}

// ParsePageMode 从字符串解析模式
func ParsePageMode(s string) (PageMode, error) {
	mode := PageMode(s)
	if !mode.IsValid() {
		return "", fmt.Errorf("invalid page mode: %s", s)
	}
	return mode, nil
}

// DocumentStatus 写作文档状态
type DocumentStatus string

const (
	DocumentStatusDraft     DocumentStatus = "draft"
	DocumentStatusPublished DocumentStatus = "published"
	DocumentStatusArchived  DocumentStatus = "archived"
	DocumentStatusDeleted   DocumentStatus = "deleted"
)

// AllDocumentStatuses 所有有效状态
var AllDocumentStatuses = []DocumentStatus{
	DocumentStatusDraft,
	DocumentStatusPublished,
	DocumentStatusArchived,
	DocumentStatusDeleted,
}

// IsValid 检查状态是否有效
func (s DocumentStatus) IsValid() bool {
	switch s {
	case DocumentStatusDraft, DocumentStatusPublished,
		DocumentStatusArchived, DocumentStatusDeleted:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (s DocumentStatus) String() string {
	return string(s)
}

// ParseDocumentStatus 从字符串解析状态
func ParseDocumentStatus(s string) (DocumentStatus, error) {
	status := DocumentStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid document status: %s", s)
	}
	return status, nil
}

// IsPublic 是否公开状态
func (s DocumentStatus) IsPublic() bool {
	return s == DocumentStatusPublished
}

// CanEdit 是否可编辑
func (s DocumentStatus) CanEdit() bool {
	return s == DocumentStatusDraft || s == DocumentStatusPublished
}

// CanDelete 是否可删除
func (s DocumentStatus) CanDelete() bool {
	return s != DocumentStatusDeleted
}

// WithdrawalStatus 提现状态
type WithdrawalStatus string

const (
	WithdrawalStatusPending   WithdrawalStatus = "pending"
	WithdrawalStatusApproved  WithdrawalStatus = "approved"
	WithdrawalStatusRejected  WithdrawalStatus = "rejected"
	WithdrawalStatusCompleted WithdrawalStatus = "completed"
)

// AllWithdrawalStatuses 所有有效状态
var AllWithdrawalStatuses = []WithdrawalStatus{
	WithdrawalStatusPending,
	WithdrawalStatusApproved,
	WithdrawalStatusRejected,
	WithdrawalStatusCompleted,
}

// IsValid 检查状态是否有效
func (s WithdrawalStatus) IsValid() bool {
	switch s {
	case WithdrawalStatusPending, WithdrawalStatusApproved,
		WithdrawalStatusRejected, WithdrawalStatusCompleted:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (s WithdrawalStatus) String() string {
	return string(s)
}

// ParseWithdrawalStatus 从字符串解析状态
func ParseWithdrawalStatus(s string) (WithdrawalStatus, error) {
	status := WithdrawalStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid withdrawal status: %s", s)
	}
	return status, nil
}

// IsFinal 是否为最终状态
func (s WithdrawalStatus) IsFinal() bool {
	return s == WithdrawalStatusRejected || s == WithdrawalStatusCompleted
}

// CanApprove 是否可批准
func (s WithdrawalStatus) CanApprove() bool {
	return s == WithdrawalStatusPending
}

// CanReject 是否可拒绝
func (s WithdrawalStatus) CanReject() bool {
	return s == WithdrawalStatusPending
}

// CanComplete 是否可完成
func (s WithdrawalStatus) CanComplete() bool {
	return s == WithdrawalStatusApproved
}

// BookStatus 书籍状态
type BookStatus string

const (
	BookStatusDraft     BookStatus = "draft"
	BookStatusPublished BookStatus = "published"
	BookStatusCompleted BookStatus = "completed"
	BookStatusPaused    BookStatus = "paused"
	BookStatusDeleted   BookStatus = "deleted"
)

// AllBookStatuses 所有有效状态
var AllBookStatuses = []BookStatus{
	BookStatusDraft,
	BookStatusPublished,
	BookStatusCompleted,
	BookStatusPaused,
	BookStatusDeleted,
}

// IsValid 检查状态是否有效
func (s BookStatus) IsValid() bool {
	switch s {
	case BookStatusDraft, BookStatusPublished,
		BookStatusCompleted, BookStatusPaused, BookStatusDeleted:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (s BookStatus) String() string {
	return string(s)
}

// ParseBookStatus 从字符串解析状态
func ParseBookStatus(s string) (BookStatus, error) {
	status := BookStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid book status: %s", s)
	}
	return status, nil
}

// IsPublic 是否公开
func (s BookStatus) IsPublic() bool {
	return s == BookStatusPublished
}

// CanEdit 是否可编辑
func (s BookStatus) CanEdit() bool {
	return s == BookStatusDraft || s == BookStatusPaused
}

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
)

// AllOrderStatuses 所有有效状态
var AllOrderStatuses = []OrderStatus{
	OrderStatusPending,
	OrderStatusPaid,
	OrderStatusCompleted,
	OrderStatusCancelled,
	OrderStatusRefunded,
}

// IsValid 检查状态是否有效
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusPaid,
		OrderStatusCompleted, OrderStatusCancelled, OrderStatusRefunded:
		return true
	default:
		return false
	}
}

// String 转换为字符串
func (s OrderStatus) String() string {
	return string(s)
}

// ParseOrderStatus 从字符串解析状态
func ParseOrderStatus(s string) (OrderStatus, error) {
	status := OrderStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid order status: %s", s)
	}
	return status, nil
}

// IsPaid 是否已支付
func (s OrderStatus) IsPaid() bool {
	return s == OrderStatusPaid || s == OrderStatusCompleted
}

// CanCancel 是否可取消
func (s OrderStatus) CanCancel() bool {
	return s == OrderStatusPending
}

// CanRefund 是否可退款
func (s OrderStatus) CanRefund() bool {
	return s == OrderStatusPaid || s == OrderStatusCompleted
}
