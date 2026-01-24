package user

import (
	"fmt"
	"time"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/shared"
	"Qingyu_backend/models/shared/types"
	usersModel "Qingyu_backend/models/users"
)

// ToUserDTO Model → DTO 转换
// 将 User Model 转换为 UserDTO 用于 API 层返回
func ToUserDTO(user *usersModel.User) *dto.UserDTO {
	if user == nil {
		return nil
	}

	var converter types.DTOConverter

	return &dto.UserDTO{
		ID:        converter.ModelIDToDTO(user.ID),
		CreatedAt: converter.TimeToISO8601(user.CreatedAt),
		UpdatedAt: converter.TimeToISO8601(user.UpdatedAt),

		// 基本信息
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,

		// 角色和权限
		Roles:    user.Roles,
		VIPLevel: user.VIPLevel,

		// 状态和资料
		Status:   string(user.Status), // 直接转换为字符串
		Avatar:   user.Avatar,
		Nickname: user.Nickname,
		Bio:      user.Bio,

		// 认证相关
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		LastLoginAt:   converter.TimeToISO8601(user.LastLoginAt),
		LastLoginIP:   user.LastLoginIP,
	}
}

// ToUserDTOs 批量转换 Model → DTO
func ToUserDTOs(users []*usersModel.User) []*dto.UserDTO {
	result := make([]*dto.UserDTO, len(users))
	for i := range users {
		result[i] = ToUserDTO(users[i])
	}
	return result
}

// ToUserDTOsFromSlice 从切片转换 Model → DTO
func ToUserDTOsFromSlice(users []usersModel.User) []*dto.UserDTO {
	result := make([]*dto.UserDTO, len(users))
	for i := range users {
		result[i] = ToUserDTO(&users[i])
	}
	return result
}

// ToUser 从 DTO 创建 Model（用于更新）
// 注意：不包含密码字段，密码更新应该使用专门的方法
func ToUser(dto *dto.UserDTO) (*usersModel.User, error) {
	if dto == nil {
		return nil, nil
	}

	var converter types.DTOConverter

	id, err := converter.DTOIDToModel(dto.ID)
	if err != nil {
		return nil, err
	}

	createdAt, err := converter.ISO8601ToTime(dto.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := converter.ISO8601ToTime(dto.UpdatedAt)
	if err != nil {
		return nil, err
	}

	status := usersModel.UserStatus(dto.Status)
	// 简单验证：确保状态是有效的枚举值
	switch status {
	case usersModel.UserStatusActive, usersModel.UserStatusInactive, usersModel.UserStatusBanned, usersModel.UserStatusDeleted:
		// 有效状态
	default:
		return nil, fmt.Errorf("invalid user status: %s", dto.Status)
	}

	lastLoginAt, err := converter.ISO8601ToTime(dto.LastLoginAt)
	if err != nil {
		return nil, err
	}

	return &usersModel.User{
		IdentifiedEntity: shared.IdentifiedEntity{ID: id},
		BaseEntity:       shared.BaseEntity{CreatedAt: createdAt, UpdatedAt: updatedAt},

		Username: dto.Username,
		Email:    dto.Email,
		Phone:    dto.Phone,

		Roles:    dto.Roles,
		VIPLevel: dto.VIPLevel,

		Status:   status,
		Avatar:   dto.Avatar,
		Nickname: dto.Nickname,
		Bio:      dto.Bio,

		EmailVerified: dto.EmailVerified,
		PhoneVerified: dto.PhoneVerified,
		LastLoginAt:   lastLoginAt,
		LastLoginIP:   dto.LastLoginIP,
	}, nil
}

// ToUserWithoutID 从 DTO 创建 Model（用于创建新用户）
// 不设置 ID，让数据库自动生成
func ToUserWithoutID(dto *dto.UserDTO) (*usersModel.User, error) {
	if dto == nil {
		return nil, nil
	}

	var converter types.DTOConverter

	status := usersModel.UserStatus(dto.Status)
	// 简单验证：确保状态是有效的枚举值
	switch status {
	case usersModel.UserStatusActive, usersModel.UserStatusInactive, usersModel.UserStatusBanned, usersModel.UserStatusDeleted:
		// 有效状态
	default:
		return nil, fmt.Errorf("invalid user status: %s", dto.Status)
	}

	lastLoginAt, err := converter.ISO8601ToTime(dto.LastLoginAt)
	if err != nil {
		// 如果 LastLoginAt 为空，使用零值
		lastLoginAt = time.Time{}
	}

	return &usersModel.User{
		IdentifiedEntity: shared.IdentifiedEntity{}, // ID 将由数据库生成
		BaseEntity:       shared.BaseEntity{},        // 时间戳将由数据库设置

		Username: dto.Username,
		Email:    dto.Email,
		Phone:    dto.Phone,

		Roles:    dto.Roles,
		VIPLevel: dto.VIPLevel,

		Status:   status,
		Avatar:   dto.Avatar,
		Nickname: dto.Nickname,
		Bio:      dto.Bio,

		EmailVerified: dto.EmailVerified,
		PhoneVerified: dto.PhoneVerified,
		LastLoginAt:   lastLoginAt,
		LastLoginIP:   dto.LastLoginIP,
	}, nil
}
