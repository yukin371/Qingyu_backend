package ai

import (
	"context"
	"strings"
	"time"

	aiModels "Qingyu_backend/models/ai"
	financeModel "Qingyu_backend/models/finance"
	usersModel "Qingyu_backend/models/users"
)

const (
	quotaCustomFieldManualOverride = "manualOverride"
	quotaCustomFieldProfileManaged = "profileManaged"
)

type quotaProfileUserReader interface {
	GetByID(ctx context.Context, id string) (*usersModel.User, error)
}

type quotaProfileMembershipReader interface {
	GetMembership(ctx context.Context, userID string) (*financeModel.UserMembership, error)
}

// QuotaProfile 表示当前用户的有效 AI 配额画像。
type QuotaProfile struct {
	UserID          string
	UserRole        string
	MembershipLevel string
}

// QuotaProfileResolver 负责从用户与会员真相解析 AI 配额画像。
type QuotaProfileResolver struct {
	userReader       quotaProfileUserReader
	membershipReader quotaProfileMembershipReader
}

func NewQuotaProfileResolver(userReader quotaProfileUserReader, membershipReader quotaProfileMembershipReader) *QuotaProfileResolver {
	return &QuotaProfileResolver{
		userReader:       userReader,
		membershipReader: membershipReader,
	}
}

func (r *QuotaProfileResolver) Resolve(ctx context.Context, userID string) (*QuotaProfile, error) {
	profile := &QuotaProfile{
		UserID:          userID,
		UserRole:        string(aiModels.UserRoleReader),
		MembershipLevel: financeModel.MembershipLevelNormal,
	}

	if r == nil {
		return profile, nil
	}

	if r.userReader != nil {
		if user, err := r.userReader.GetByID(ctx, userID); err == nil && user != nil {
			profile.UserRole = normalizeQuotaUserRole(user.Roles)
		}
	}

	if r.membershipReader != nil {
		if membership, err := r.membershipReader.GetMembership(ctx, userID); err == nil && membership != nil && isMembershipActive(membership) {
			profile.MembershipLevel = normalizeQuotaMembershipLevel(membership.Level)
		}
	}

	return profile, nil
}

func normalizeQuotaUserRole(roles []string) string {
	if len(roles) == 0 {
		return string(aiModels.UserRoleReader)
	}

	seen := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		seen[strings.ToLower(strings.TrimSpace(role))] = struct{}{}
	}

	if _, ok := seen["admin"]; ok {
		return string(aiModels.UserRoleAdmin)
	}
	if _, ok := seen["author"]; ok {
		return string(aiModels.UserRoleWriter)
	}
	if _, ok := seen["writer"]; ok {
		return string(aiModels.UserRoleWriter)
	}
	return string(aiModels.UserRoleReader)
}

func normalizeQuotaMembershipLevel(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "", "normal":
		return financeModel.MembershipLevelNormal
	case "vip":
		return financeModel.MembershipLevelVIPMonthly
	case financeModel.MembershipLevelVIPMonthly:
		return financeModel.MembershipLevelVIPMonthly
	case financeModel.MembershipLevelVIPYearly:
		return financeModel.MembershipLevelVIPYearly
	case financeModel.MembershipLevelSuperVIP:
		return financeModel.MembershipLevelSuperVIP
	default:
		return financeModel.MembershipLevelNormal
	}
}

func isMembershipActive(membership *financeModel.UserMembership) bool {
	if membership == nil {
		return false
	}
	if membership.Status != financeModel.MembershipStatusActive {
		return false
	}
	return !time.Now().After(membership.EndTime)
}

func ensureQuotaMetadata(quota *aiModels.UserQuota) *aiModels.QuotaMetadata {
	if quota.Metadata == nil {
		quota.Metadata = &aiModels.QuotaMetadata{}
	}
	if quota.Metadata.CustomFields == nil {
		quota.Metadata.CustomFields = make(map[string]interface{})
	}
	return quota.Metadata
}

func setQuotaManualOverride(quota *aiModels.UserQuota, enabled bool) {
	metadata := ensureQuotaMetadata(quota)
	metadata.CustomFields[quotaCustomFieldManualOverride] = enabled
}

func isQuotaManualOverride(quota *aiModels.UserQuota) bool {
	if quota == nil || quota.Metadata == nil || quota.Metadata.CustomFields == nil {
		return false
	}
	value, ok := quota.Metadata.CustomFields[quotaCustomFieldManualOverride]
	if !ok {
		return false
	}
	flag, ok := value.(bool)
	return ok && flag
}
