package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DTOConverter DTO 转换辅助
type DTOConverter struct{}

// ===== ID 转换 =====

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

// ===== 金额转换 =====

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

// ===== 评分转换 =====

// RatingToDTO Rating → 评分字符串 (Rating → "4.5")
func (DTOConverter) RatingToDTO(rating Rating) string {
	return rating.String()
}

// RatingToFloat Rating → float64
func (DTOConverter) RatingToFloat(rating Rating) float64 {
	return rating.ToFloat()
}

// DTORatingToModel 评分字符串 → Rating ("4.5" → Rating)
func (DTOConverter) DTORatingToModel(s string) (Rating, error) {
	var value float64
	_, err := fmt.Sscanf(s, "%f", &value)
	if err != nil {
		return RatingDefault, err
	}
	return NewRating(value)
}

// DTORatingToFloat float64 → Rating
func (DTOConverter) DTORatingToFloat(value float64) (Rating, error) {
	return NewRating(value)
}

// ===== 进度转换 =====

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

// ===== 时间戳转换 =====

// TimeToISO8601 time.Time → ISO8601 字符串 (RFC3339)
// 用于将 BaseEntity 的时间字段转换为 API 层的字符串格式
func (DTOConverter) TimeToISO8601(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// TimeToISO8601Ptr *time.Time → *ISO8601 字符串指针 (RFC3339)
// 用于将 BaseEntity 的时间字段指针转换为 API 层的字符串格式指针
func (DTOConverter) TimeToISO8601Ptr(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	result := t.Format(time.RFC3339)
	return &result
}

// TimeToISO8601PtrToString *time.Time → string (RFC3339)
// 用于将 *time.Time 转换为 string，如果为 nil 则返回空字符串
func (DTOConverter) TimeToISO8601PtrToString(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// StringToISO8601ToTimePtr string → *time.Time (RFC3339)
// 用于将 string 转换为 *time.Time，如果为空字符串则返回 nil
func (DTOConverter) StringToISO8601ToTimePtr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ISO8601ToTime ISO8601 字符串 → time.Time
// 用于将 API 层的时间字符串转换为 Model 层的 time.Time
func (DTOConverter) ISO8601ToTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}

// ISO8601ToTimePtr *ISO8601 字符串指针 → *time.Time
// 用于将 API 层的时间字符串指针转换为 Model 层的 time.Time 指针
func (DTOConverter) ISO8601ToTimePtr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// TimeToUnix time.Time → Unix 时间戳（秒）
// 用于返回 Unix 时间戳给前端
func (DTOConverter) TimeToUnix(t time.Time) int64 {
	return t.Unix()
}

// UnixToTime Unix 时间戳 → time.Time
// 用于将前端传来的 Unix 时间戳转换为 time.Time
func (DTOConverter) UnixToTime(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// TimesToISO8601 批量转换 time.Time → ISO8601 字符串
func (DTOConverter) TimesToISO8601(times []time.Time) []string {
	result := make([]string, len(times))
	for i, t := range times {
		result[i] = t.Format(time.RFC3339)
	}
	return result
}

// ISO8601sToTimes 批量转换 ISO8601 字符串 → time.Time
func (DTOConverter) ISO8601sToTimes(strs []string) ([]time.Time, error) {
	result := make([]time.Time, len(strs))
	for i, s := range strs {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time at index %d: %w", i, err)
		}
		result[i] = t
	}
	return result, nil
}

// ===== 枚举转换 =====

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

// ===== BaseEntity 转换辅助 =====

// CopyBaseFields 复制 BaseEntity 字段到 DTO 格式
// 返回: (id string, createdAt string, updatedAt string)
// 用法:
//   dto.ID, dto.CreatedAt, dto.UpdatedAt = converter.CopyBaseFields(model.ID, model.CreatedAt, model.UpdatedAt)
func (DTOConverter) CopyBaseFields(id primitive.ObjectID, createdAt, updatedAt time.Time) (string, string, string) {
	return ToHex(id),
		createdAt.Format(time.RFC3339),
		updatedAt.Format(time.RFC3339)
}

// ParseBaseFields 从 DTO 解析 BaseEntity 字段
// 返回: (id ObjectID, createdAt time.Time, updatedAt time.Time, error)
// 用法:
//   id, createdAt, updatedAt, err := converter.ParseBaseFields(dto.ID, dto.CreatedAt, dto.UpdatedAt)
func (DTOConverter) ParseBaseFields(idStr, createdAtStr, updatedAtStr string) (primitive.ObjectID, time.Time, time.Time, error) {
	id, err := ParseObjectID(idStr)
	if err != nil {
		return primitive.ObjectID{}, time.Time{}, time.Time{}, fmt.Errorf("invalid id: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return primitive.ObjectID{}, time.Time{}, time.Time{}, fmt.Errorf("invalid createdAt: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return primitive.ObjectID{}, time.Time{}, time.Time{}, fmt.Errorf("invalid updatedAt: %w", err)
	}

	return id, createdAt, updatedAt, nil
}

// DefaultConverter 默认转换器实例
var DefaultConverter = DTOConverter{}
