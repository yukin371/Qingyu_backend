package user

import (
	"testing"
	"time"

	"Qingyu_backend/models/dto"
	usersModel "Qingyu_backend/models/users"
)

// TestToUserDTO_Normal 测试正常Model转DTO
func TestToUserDTO_Normal(t *testing.T) {
	// 创建测试用的Model对象
	now := time.Now()
	userModel := &usersModel.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Roles:    []string{"user"},
		VIPLevel: 1,
		Status:   usersModel.UserStatusActive,
		Avatar:   "https://example.com/avatar.jpg",
		Nickname: "Test User",
		Bio:      "Test bio",
		Gender:   "male",
		Birthday: &now,
		Location: "Beijing",
		Website:  "https://example.com",
		EmailVerified: true,
		PhoneVerified: false,
		LastLoginAt:   now,
		LastLoginIP:   "192.168.1.1",
	}

	// 执行转换
	dto := ToUserDTO(userModel)

	// 验证转换结果
	if dto == nil {
		t.Fatal("ToUserDTO() returned nil")
	}

	if dto.Username != userModel.Username {
		t.Errorf("Username = %v, want %v", dto.Username, userModel.Username)
	}

	if dto.Email != userModel.Email {
		t.Errorf("Email = %v, want %v", dto.Email, userModel.Email)
	}

	if dto.Phone != userModel.Phone {
		t.Errorf("Phone = %v, want %v", dto.Phone, userModel.Phone)
	}

	if dto.Status != string(userModel.Status) {
		t.Errorf("Status = %v, want %v", dto.Status, string(userModel.Status))
	}

	if dto.EmailVerified != userModel.EmailVerified {
		t.Errorf("EmailVerified = %v, want %v", dto.EmailVerified, userModel.EmailVerified)
	}

	if dto.PhoneVerified != userModel.PhoneVerified {
		t.Errorf("PhoneVerified = %v, want %v", dto.PhoneVerified, userModel.PhoneVerified)
	}
}

// TestToUserDTO_Nil 测试nil输入
func TestToUserDTO_Nil(t *testing.T) {
	// 测试nil输入
	dto := ToUserDTO(nil)

	// 验证返回nil
	if dto != nil {
		t.Errorf("ToUserDTO(nil) = %v, want nil", dto)
	}
}

// TestToUserDTO_EmptyUser 测试空User对象
func TestToUserDTO_EmptyUser(t *testing.T) {
	// 创建空的User对象
	userModel := &usersModel.User{}

	// 执行转换
	dto := ToUserDTO(userModel)

	// 验证转换结果
	if dto == nil {
		t.Fatal("ToUserDTO() returned nil for empty user")
	}

	// 验证空值被正确转换
	if dto.Username != "" {
		t.Errorf("Username = %v, want empty string", dto.Username)
	}

	if dto.Email != "" {
		t.Errorf("Email = %v, want empty string", dto.Email)
	}
}

// TestToUserDTOs_Normal 测试批量转换
func TestToUserDTOs_Normal(t *testing.T) {
	// 创建测试用的Model对象列表
	user1 := &usersModel.User{
		Username: "user1",
		Email:    "user1@example.com",
		Status:   usersModel.UserStatusActive,
	}

	user2 := &usersModel.User{
		Username: "user2",
		Email:    "user2@example.com",
		Status:   usersModel.UserStatusInactive,
	}

	users := []*usersModel.User{user1, user2}

	// 执行批量转换
	dtos := ToUserDTOs(users)

	// 验证转换结果
	if len(dtos) != len(users) {
		t.Errorf("ToUserDTOs() returned %d items, want %d", len(dtos), len(users))
	}

	if dtos[0].Username != user1.Username {
		t.Errorf("dtos[0].Username = %v, want %v", dtos[0].Username, user1.Username)
	}

	if dtos[1].Username != user2.Username {
		t.Errorf("dtos[1].Username = %v, want %v", dtos[1].Username, user2.Username)
	}
}

// TestToUserDTOs_Empty 测试空列表转换
func TestToUserDTOs_Empty(t *testing.T) {
	// 测试空列表
	users := []*usersModel.User{}
	dtos := ToUserDTOs(users)

	// 验证返回空列表
	if len(dtos) != 0 {
		t.Errorf("ToUserDTOs(empty) returned %d items, want 0", len(dtos))
	}
}

// TestToUserDTOs_NilSlice 测试包含nil的列表转换
func TestToUserDTOs_NilSlice(t *testing.T) {
	// 测试包含nil的列表
	users := []*usersModel.User{nil, nil}
	dtos := ToUserDTOs(users)

	// 验证转换结果
	if len(dtos) != 2 {
		t.Errorf("ToUserDTOs() returned %d items, want 2", len(dtos))
	}

	if dtos[0] != nil {
		t.Errorf("dtos[0] = %v, want nil", dtos[0])
	}

	if dtos[1] != nil {
		t.Errorf("dtos[1] = %v, want nil", dtos[1])
	}
}

// TestToUserDTOsFromSlice_Normal 测试从切片转换
func TestToUserDTOsFromSlice_Normal(t *testing.T) {
	// 创建测试用的User切片
	users := []usersModel.User{
		{
			Username: "user1",
			Email:    "user1@example.com",
			Status:   usersModel.UserStatusActive,
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
			Status:   usersModel.UserStatusInactive,
		},
	}

	// 执行转换
	dtos := ToUserDTOsFromSlice(users)

	// 验证转换结果
	if len(dtos) != len(users) {
		t.Errorf("ToUserDTOsFromSlice() returned %d items, want %d", len(dtos), len(users))
	}

	if dtos[0].Username != users[0].Username {
		t.Errorf("dtos[0].Username = %v, want %v", dtos[0].Username, users[0].Username)
	}

	if dtos[1].Username != users[1].Username {
		t.Errorf("dtos[1].Username = %v, want %v", dtos[1].Username, users[1].Username)
	}
}

// TestToUserDTOsFromSlice_Empty 测试空切片转换
func TestToUserDTOsFromSlice_Empty(t *testing.T) {
	// 测试空切片
	users := []usersModel.User{}
	dtos := ToUserDTOsFromSlice(users)

	// 验证返回空列表
	if len(dtos) != 0 {
		t.Errorf("ToUserDTOsFromSlice(empty) returned %d items, want 0", len(dtos))
	}
}

// TestToUser_Normal 测试正常DTO转Model
func TestToUser_Normal(t *testing.T) {
	// 创建测试用的DTO对象
	userDTO := &dto.UserDTO{
		ID:      "507f1f77bcf86cd799439011",
		Username: "testuser",
		Email:   "test@example.com",
		Phone:   "13800138000",
		Roles:   []string{"user"},
		VIPLevel: 1,
		Status:  "active",
		Avatar:  "https://example.com/avatar.jpg",
		Nickname: "Test User",
		Bio:     "Test bio",
		Gender:  "male",
		Location: "Beijing",
		Website: "https://example.com",
		EmailVerified: true,
		PhoneVerified: false,
		LastLoginIP:   "192.168.1.1",
	}

	// 执行转换
	userModel, err := ToUser(userDTO)

	// 验证转换结果
	if err != nil {
		t.Fatalf("ToUser() returned error: %v", err)
	}

	if userModel == nil {
		t.Fatal("ToUser() returned nil model")
	}

	if userModel.Username != userDTO.Username {
		t.Errorf("Username = %v, want %v", userModel.Username, userDTO.Username)
	}

	if userModel.Email != userDTO.Email {
		t.Errorf("Email = %v, want %v", userModel.Email, userDTO.Email)
	}

	if userModel.Status != usersModel.UserStatus(userDTO.Status) {
		t.Errorf("Status = %v, want %v", userModel.Status, usersModel.UserStatus(userDTO.Status))
	}
}

// TestToUser_Nil 测试nil DTO输入
func TestToUser_Nil(t *testing.T) {
	// 测试nil输入
	userModel, err := ToUser(nil)

	// 验证返回nil且无错误
	if err != nil {
		t.Fatalf("ToUser(nil) returned error: %v", err)
	}

	if userModel != nil {
		t.Errorf("ToUser(nil) = %v, want nil", userModel)
	}
}

// TestToUser_InvalidStatus 测试无效状态
func TestToUser_InvalidStatus(t *testing.T) {
	// 创建包含无效状态的DTO
	userDTO := &dto.UserDTO{
		ID:       "507f1f77bcf86cd799439011",
		Username: "testuser",
		Status:   "invalid_status",
	}

	// 执行转换
	userModel, err := ToUser(userDTO)

	// 验证返回错误
	if err == nil {
		t.Error("ToUser() should return error for invalid status")
	}

	if userModel != nil {
		t.Errorf("ToUser() returned model despite error: %v", userModel)
	}
}

// TestToUser_AllValidStatuses 测试所有有效状态
func TestToUser_AllValidStatuses(t *testing.T) {
	validStatuses := []string{
		string(usersModel.UserStatusActive),
		string(usersModel.UserStatusInactive),
		string(usersModel.UserStatusBanned),
		string(usersModel.UserStatusDeleted),
	}

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			userDTO := &dto.UserDTO{
				ID:       "507f1f77bcf86cd799439011",
				Username: "testuser",
				Status:   status,
			}

			userModel, err := ToUser(userDTO)

			if err != nil {
				t.Errorf("ToUser() returned error for valid status %s: %v", status, err)
			}

			if userModel.Status != usersModel.UserStatus(status) {
				t.Errorf("Status = %v, want %v", userModel.Status, usersModel.UserStatus(status))
			}
		})
	}
}

// TestToUserWithoutID_Normal 测试无ID转换
func TestToUserWithoutID_Normal(t *testing.T) {
	// 创建测试用的DTO对象
	userDTO := &dto.UserDTO{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Roles:    []string{"user"},
		VIPLevel: 1,
		Status:   "active",
		Avatar:   "https://example.com/avatar.jpg",
		Nickname: "Test User",
		Bio:      "Test bio",
		Gender:   "male",
		Location: "Beijing",
		Website:  "https://example.com",
		EmailVerified: true,
		PhoneVerified: false,
		LastLoginIP:   "192.168.1.1",
	}

	// 执行转换
	userModel, err := ToUserWithoutID(userDTO)

	// 验证转换结果
	if err != nil {
		t.Fatalf("ToUserWithoutID() returned error: %v", err)
	}

	if userModel == nil {
		t.Fatal("ToUserWithoutID() returned nil model")
	}

	// 验证ID为零值（由数据库生成）
// 注意：Go中的零值ObjectId不是完全为空，而是"000000000000000000000000"
if userModel.ID.Hex() != "000000000000000000000000" {
	t.Errorf("ID should be zero value ObjectId for ToUserWithoutID, got %v", userModel.ID.Hex())
}

	if userModel.Username != userDTO.Username {
		t.Errorf("Username = %v, want %v", userModel.Username, userDTO.Username)
	}

	if userModel.Email != userDTO.Email {
		t.Errorf("Email = %v, want %v", userModel.Email, userDTO.Email)
	}

	if userModel.Status != usersModel.UserStatus(userDTO.Status) {
		t.Errorf("Status = %v, want %v", userModel.Status, usersModel.UserStatus(userDTO.Status))
	}
}

// TestToUserWithoutID_Nil 测试nil DTO输入
func TestToUserWithoutID_Nil(t *testing.T) {
	// 测试nil输入
	userModel, err := ToUserWithoutID(nil)

	// 验证返回nil且无错误
	if err != nil {
		t.Fatalf("ToUserWithoutID(nil) returned error: %v", err)
	}

	if userModel != nil {
		t.Errorf("ToUserWithoutID(nil) = %v, want nil", userModel)
	}
}

// TestToUserWithoutID_InvalidStatus 测试无效状态
func TestToUserWithoutID_InvalidStatus(t *testing.T) {
	// 创建包含无效状态的DTO
	userDTO := &dto.UserDTO{
		Username: "testuser",
		Status:   "invalid_status",
	}

	// 执行转换
	userModel, err := ToUserWithoutID(userDTO)

	// 验证返回错误
	if err == nil {
		t.Error("ToUserWithoutID() should return error for invalid status")
	}

	if userModel != nil {
		t.Errorf("ToUserWithoutID() returned model despite error: %v", userModel)
	}
}

// TestRoundTrip 测试往返转换（Model -> DTO -> Model）
func TestRoundTrip(t *testing.T) {
	// 创建原始Model
	originalUser := &usersModel.User{
		Username: "testuser",
		Email:    "test@example.com",
		Phone:    "13800138000",
		Roles:    []string{"user", "admin"},
		VIPLevel: 2,
		Status:   usersModel.UserStatusActive,
		Avatar:   "https://example.com/avatar.jpg",
		Nickname: "Test User",
		Bio:      "Test bio",
		Gender:   "male",
		Location: "Beijing",
		Website:  "https://example.com",
		EmailVerified: true,
		PhoneVerified: false,
		LastLoginIP:   "192.168.1.1",
	}

	// Model -> DTO
	userDTO := ToUserDTO(originalUser)

	// DTO -> Model (带ID)
	convertedUser, err := ToUser(userDTO)
	if err != nil {
		t.Fatalf("ToUser() returned error: %v", err)
	}

	// 验证关键字段保持一致
	if convertedUser.Username != originalUser.Username {
		t.Errorf("Username changed: %v -> %v", originalUser.Username, convertedUser.Username)
	}

	if convertedUser.Email != originalUser.Email {
		t.Errorf("Email changed: %v -> %v", originalUser.Email, convertedUser.Email)
	}

	if convertedUser.Status != originalUser.Status {
		t.Errorf("Status changed: %v -> %v", originalUser.Status, convertedUser.Status)
	}

	// 验证角色列表
	if len(convertedUser.Roles) != len(originalUser.Roles) {
		t.Errorf("Roles count changed: %d -> %d", len(originalUser.Roles), len(convertedUser.Roles))
	}

	for i, role := range originalUser.Roles {
		if convertedUser.Roles[i] != role {
			t.Errorf("Role[%d] changed: %v -> %v", i, role, convertedUser.Roles[i])
		}
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("EmptyRoles", func(t *testing.T) {
		userModel := &usersModel.User{
			Username: "testuser",
			Roles:    []string{},
		}

		dto := ToUserDTO(userModel)
		if dto.Roles == nil || len(dto.Roles) != 0 {
			t.Errorf("Roles should be empty slice, got %v", dto.Roles)
		}

		// 跳过反向转换测试，因为需要有效的ID格式
		// 这个测试主要验证Model->DTO的方向
	})

	t.Run("NilRoles", func(t *testing.T) {
		userModel := &usersModel.User{
			Username: "testuser",
			Roles:    nil,
		}

		dto := ToUserDTO(userModel)
		// nil应该被转换为nil或空切片
		if dto.Roles != nil && len(dto.Roles) != 0 {
			t.Errorf("Roles should be nil or empty, got %v", dto.Roles)
		}
	})

	t.Run("ZeroVIPLevel", func(t *testing.T) {
		userModel := &usersModel.User{
			Username: "testuser",
			VIPLevel: 0,
		}

		dto := ToUserDTO(userModel)
		if dto.VIPLevel != 0 {
			t.Errorf("VIPLevel = %v, want 0", dto.VIPLevel)
		}
	})

	t.Run("MaxVIPLevel", func(t *testing.T) {
		userModel := &usersModel.User{
			Username: "testuser",
			VIPLevel: 999,
		}

		dto := ToUserDTO(userModel)
		if dto.VIPLevel != 999 {
			t.Errorf("VIPLevel = %v, want 999", dto.VIPLevel)
		}
	})
}
