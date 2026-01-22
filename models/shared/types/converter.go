package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DTOConverter DTO 转换辅助
type DTOConverter struct{}

// ModelIDToDTO Model ID → DTO ID (ObjectID → string)
func (DTOConverter) ModelIDToDTO(id primitive.ObjectID) string {
	return ToHex(id)
}

// ModelIDsToDTO 批量转换 Model ID → DTO ID
func (DTOConverter) ModelIDsToDTO(ids []primitive.ObjectID) []string {
	return ToHexSlice(ids)
}

// DTOIDToModel DTO ID → Model ID (string → ObjectID)
func (DTOConverter) DTOIDToModel(id string) (primitive.ObjectID, error) {
	return ParseObjectID(id)
}

// DTOIDsToModel 批量转换 DTO ID → Model ID
func (DTOConverter) DTOIDsToModel(ids []string) ([]primitive.ObjectID, error) {
	oids, errs := ParseObjectIDSlice(ids)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to parse %d IDs", len(errs))
	}
	return oids, nil
}

// MoneyToDTO Money → 金额字符串 (Money → "¥12.99")
func (DTOConverter) MoneyToDTO(money Money) string {
	return money.String()
}

// MoneyToYuan Money → 元 (Money → float64)
func (DTOConverter) MoneyToYuan(money Money) float64 {
	return money.ToYuan()
}

// MoneyToCents Money → 分 (Money → int64)
func (DTOConverter) MoneyToCents(money Money) int64 {
	return money.ToCents()
}

// DTOMoneyToYuan 金额字符串 → Money ("12.99" → Money)
func (DTOConverter) DTOMoneyToYuan(yuan float64) Money {
	return NewMoneyFromYuan(yuan)
}

// DTOMoneyToCents 分 → Money (int64 → Money)
func (DTOConverter) DTOMoneyToCents(cents int64) Money {
	return NewMoneyFromCents(cents)
}

// RatingToDTO Rating → 评分字符串 (Rating → "4.5")
func (DTOConverter) RatingToDTO(rating Rating) string {
	return rating.String()
}

// RatingToFloat Rating → float32
func (DTOConverter) RatingToFloat(rating Rating) float32 {
	return rating.ToFloat()
}

// DTORatingToModel 评分字符串 → Rating ("4.5" → Rating)
func (DTOConverter) DTORatingToModel(s string) (Rating, error) {
	var value float32
	_, err := fmt.Sscanf(s, "%f", &value)
	if err != nil {
		return RatingDefault, err
	}
	return NewRating(value)
}

// DTORatingToFloat float32 → Rating
func (DTOConverter) DTORatingToFloat(value float32) (Rating, error) {
	return NewRating(value)
}

// ProgressToDTO Progress → 百分比 (Progress → 75)
func (DTOConverter) ProgressToDTO(progress Progress) int {
	return progress.ToPercent()
}

// ProgressToFloat Progress → float32 (0-1)
func (DTOConverter) ProgressToFloat(progress Progress) float32 {
	return progress.ToFloat()
}

// ProgressToString Progress → 百分比字符串 (Progress → "75%")
func (DTOConverter) ProgressToString(progress Progress) string {
	return progress.String()
}

// DTOProgressToModel 百分比 → Progress (75 → Progress)
func (DTOConverter) DTOProgressToModel(percent int) (Progress, error) {
	return NewProgressFromPercent(percent)
}

// DTOProgressFromFloat float32 → Progress (0.75 → Progress)
func (DTOConverter) DTOProgressFromFloat(value float32) (Progress, error) {
	return NewProgress(value)
}

// UserRoleToString UserRole → string
func (DTOConverter) UserRoleToString(role UserRole) string {
	return role.String()
}

// StringToUserRole string → UserRole
func (DTOConverter) StringToUserRole(s string) (UserRole, error) {
	return ParseUserRole(s)
}

// PageModeToString PageMode → string
func (DTOConverter) PageModeToString(mode PageMode) string {
	return mode.String()
}

// StringToPageMode string → PageMode
func (DTOConverter) StringToPageMode(s string) (PageMode, error) {
	return ParsePageMode(s)
}

// DocumentStatusToString DocumentStatus → string
func (DTOConverter) DocumentStatusToString(status DocumentStatus) string {
	return status.String()
}

// StringToDocumentStatus string → DocumentStatus
func (DTOConverter) StringToDocumentStatus(s string) (DocumentStatus, error) {
	return ParseDocumentStatus(s)
}

// BookStatusToString BookStatus → string
func (DTOConverter) BookStatusToString(status BookStatus) string {
	return status.String()
}

// StringToBookStatus string → BookStatus
func (DTOConverter) StringToBookStatus(s string) (BookStatus, error) {
	return ParseBookStatus(s)
}

// WithdrawalStatusToString WithdrawalStatus → string
func (DTOConverter) WithdrawalStatusToString(status WithdrawalStatus) string {
	return status.String()
}

// StringToWithdrawalStatus string → WithdrawalStatus
func (DTOConverter) StringToWithdrawalStatus(s string) (WithdrawalStatus, error) {
	return ParseWithdrawalStatus(s)
}

// OrderStatusToString OrderStatus → string
func (DTOConverter) OrderStatusToString(status OrderStatus) string {
	return status.String()
}

// StringToOrderStatus string → OrderStatus
func (DTOConverter) StringToOrderStatus(s string) (OrderStatus, error) {
	return ParseOrderStatus(s)
}

// DefaultConverter 默认转换器实例
var DefaultConverter = DTOConverter{}
