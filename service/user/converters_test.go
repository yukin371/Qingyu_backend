package user

import (
	"testing"
	"time"

	"Qingyu_backend/models/dto"
	"Qingyu_backend/models/shared"
	usersModel "Qingyu_backend/models/users"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestUserToDTO 测试 User -> UserDTO 转换
func TestUserToDTO(t *testing.T) {
	t.Run("正常转换", func(t *testing.T) {
		now := time.Now()
		birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		userID := primitive.NewObjectID()

		userModel := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},

			Username: "testuser",
			Email:    "test@example.com",
			Phone:    "13800138000",
			Password: "hashed_password_here", // 密码不应出现在 DTO 中

			Roles:    []string{"reader", "author"},
			VIPLevel: 2,

			Status:   usersModel.UserStatusActive,
			Avatar:   "https://example.com/avatar.jpg",
			Nickname: "Test User",
			Bio:      "This is a test bio",

			Gender:   "male",
			Birthday: &birthday,
			Location: "Beijing, China",
			Website:  "https://example.com",

			EmailVerified: true,
			PhoneVerified: false,
			LastLoginAt:   now,
			LastLoginIP:   "192.168.1.1",
		}

		dto := ToUserDTO(userModel)

		assert.NotNil(t, dto)
		assert.Equal(t, userID.Hex(), dto.ID)
		assert.Equal(t, "testuser", dto.Username)
		assert.Equal(t, "test@example.com", dto.Email)
		assert.Equal(t, "13800138000", dto.Phone)
		assert.Equal(t, []string{"reader", "author"}, dto.Roles)
		assert.Equal(t, 2, dto.VIPLevel)
		assert.Equal(t, "active", dto.Status)
		assert.Equal(t, "https://example.com/avatar.jpg", dto.Avatar)
		assert.Equal(t, "Test User", dto.Nickname)
		assert.Equal(t, "This is a test bio", dto.Bio)
		assert.Equal(t, "male", dto.Gender)
		assert.Equal(t, "Beijing, China", dto.Location)
		assert.Equal(t, "https://example.com", dto.Website)
		assert.True(t, dto.EmailVerified)
		assert.False(t, dto.PhoneVerified)
		assert.Equal(t, "192.168.1.1", dto.LastLoginIP)
	})

	t.Run("nil 输入", func(t *testing.T) {
		dto := ToUserDTO(nil)
		assert.Nil(t, dto)
	})
}

// TestDTOToUser 测试 UserDTO -> User 转换
func TestDTOToUser(t *testing.T) {
	t.Run("正常转换", func(t *testing.T) {
		userID := primitive.NewObjectID()
		now := time.Now().Format(time.RFC3339)
		birthday := "1990-01-01T00:00:00Z"

		userDTO := &dto.UserDTO{
			ID:        userID.Hex(),
			CreatedAt: now,
			UpdatedAt: now,

			Username: "testuser",
			Email:    "test@example.com",
			Phone:    "13800138000",

			Roles:    []string{"reader", "author"},
			VIPLevel: 2,

			Status:   "active",
			Avatar:   "https://example.com/avatar.jpg",
			Nickname: "Test User",
			Bio:      "This is a test bio",

			Gender:   "male",
			Birthday: birthday,
			Location: "Beijing, China",
			Website:  "https://example.com",

			EmailVerified: true,
			PhoneVerified: false,
			LastLoginAt:   now,
			LastLoginIP:   "192.168.1.1",
		}

		userModel, err := ToUser(userDTO)

		assert.NoError(t, err)
		assert.NotNil(t, userModel)
		assert.Equal(t, userID, userModel.ID)
		assert.Equal(t, "testuser", userModel.Username)
		assert.Equal(t, "test@example.com", userModel.Email)
		assert.Equal(t, "13800138000", userModel.Phone)
		assert.Equal(t, []string{"reader", "author"}, userModel.Roles)
		assert.Equal(t, 2, userModel.VIPLevel)
		assert.Equal(t, usersModel.UserStatusActive, userModel.Status)
		assert.Equal(t, "https://example.com/avatar.jpg", userModel.Avatar)
		assert.Equal(t, "Test User", userModel.Nickname)
		assert.Equal(t, "This is a test bio", userModel.Bio)
		assert.Equal(t, "male", userModel.Gender)
		assert.Equal(t, "Beijing, China", userModel.Location)
		assert.Equal(t, "https://example.com", userModel.Website)
		assert.True(t, userModel.EmailVerified)
		assert.False(t, userModel.PhoneVerified)
		assert.Equal(t, "192.168.1.1", userModel.LastLoginIP)
	})

	t.Run("nil 输入", func(t *testing.T) {
		userModel, err := ToUser(nil)
		assert.NoError(t, err)
		assert.Nil(t, userModel)
	})

	t.Run("无效状态", func(t *testing.T) {
		userDTO := &dto.UserDTO{
			ID:       primitive.NewObjectID().Hex(),
			Username: "testuser",
			Status:   "invalid_status",
		}

		userModel, err := ToUser(userDTO)
		assert.Error(t, err)
		assert.Nil(t, userModel)
	})

	t.Run("无效 ID", func(t *testing.T) {
		userDTO := &dto.UserDTO{
			ID:       "invalid-id-format",
			Username: "testuser",
			Status:   "active",
		}

		userModel, err := ToUser(userDTO)
		assert.Error(t, err)
		assert.Nil(t, userModel)
	})
}

// TestUserToDTO_Security 测试密码字段不会被转换
func TestUserToDTO_Security(t *testing.T) {
	t.Run("密码字段不在 DTO 中", func(t *testing.T) {
		userID := primitive.NewObjectID()

		userModel := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},

			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password_12345",
			Status:   usersModel.UserStatusActive,
		}

		dto := ToUserDTO(userModel)

		assert.NotNil(t, dto)
		// DTO 中不应该有 Password 字段
		// 这是编译时检查，UserDTO 结构体中没有 Password 字段
		assert.Equal(t, "testuser", dto.Username)
		assert.Equal(t, "test@example.com", dto.Email)

		// 确保原始模型中的密码仍然存在
		assert.Equal(t, "hashed_password_12345", userModel.Password)
	})

	t.Run("往返转换不保留密码", func(t *testing.T) {
		userID := primitive.NewObjectID()

		originalUser := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},

			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password_12345",
			Roles:    []string{"reader"},
			Status:   usersModel.UserStatusActive,
		}

		// Model -> DTO
		userDTO := ToUserDTO(originalUser)

		// DTO -> Model
		convertedUser, err := ToUser(userDTO)
		assert.NoError(t, err)

		// 转换后的模型密码应该为空（因为 DTO 中没有密码字段）
		assert.Equal(t, "", convertedUser.Password)
		// 其他字段应该保持一致
		assert.Equal(t, originalUser.Username, convertedUser.Username)
		assert.Equal(t, originalUser.Email, convertedUser.Email)
	})
}

// TestUserToDTO_TimeFormat 测试时间格式转换（ISO8601）
func TestUserToDTO_TimeFormat(t *testing.T) {
	t.Run("标准时间格式", func(t *testing.T) {
		userID := primitive.NewObjectID()
		createdAt := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 2, 20, 15, 45, 30, 0, time.UTC)
		lastLoginAt := time.Date(2024, 3, 1, 8, 0, 0, 0, time.Local)

		userModel := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: createdAt, UpdatedAt: updatedAt},

			Username:   "testuser",
			Status:    usersModel.UserStatusActive,
			Roles:     []string{"reader"},
			LastLoginAt: lastLoginAt,
		}

		dto := ToUserDTO(userModel)

		assert.NotNil(t, dto)

		// 验证时间格式为 ISO8601 (RFC3339)
		expectedCreatedAt := "2024-01-15T10:30:00Z"
		expectedUpdatedAt := "2024-02-20T15:45:30Z"

		assert.Equal(t, expectedCreatedAt, dto.CreatedAt)
		assert.Equal(t, expectedUpdatedAt, dto.UpdatedAt)
		assert.NotEmpty(t, dto.LastLoginAt)

		// 验证可以解析回 time.Time
		parsedCreatedAt, err := time.Parse(time.RFC3339, dto.CreatedAt)
		assert.NoError(t, err)
		assert.Equal(t, createdAt, parsedCreatedAt)

		parsedUpdatedAt, err := time.Parse(time.RFC3339, dto.UpdatedAt)
		assert.NoError(t, err)
		assert.Equal(t, updatedAt, parsedUpdatedAt)
	})

	t.Run("零值时间转换", func(t *testing.T) {
		userID := primitive.NewObjectID()

		userModel := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Time{}, UpdatedAt: time.Time{}},

			Username: "testuser",
			Status:  usersModel.UserStatusActive,
			Roles:   []string{"reader"},
		}

		dto := ToUserDTO(userModel)

		assert.NotNil(t, dto)
		// 零值时间应该被转换为空字符串
		assert.Equal(t, "", dto.CreatedAt)
		assert.Equal(t, "", dto.UpdatedAt)
		assert.Equal(t, "", dto.LastLoginAt)
	})

	t.Run("生日时间指针转换", func(t *testing.T) {
		userID := primitive.NewObjectID()
		birthday := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

		userModel := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},

			Username: "testuser",
			Status:  usersModel.UserStatusActive,
			Roles:   []string{"reader"},
			Birthday: &birthday,
		}

		dto := ToUserDTO(userModel)

		assert.NotNil(t, dto)
		expectedBirthday := "1990-05-15T00:00:00Z"
		assert.Equal(t, expectedBirthday, dto.Birthday)

		// nil 生日应该转换为空字符串
		userModel.Birthday = nil
		dto = ToUserDTO(userModel)
		assert.Equal(t, "", dto.Birthday)
	})
}

// TestUserToDTO_IDConversion 测试 ObjectID -> string 转换
func TestUserToDTO_IDConversion(t *testing.T) {
	t.Run("正常 ObjectID 转换", func(t *testing.T) {
		userID := primitive.NewObjectID()

		userModel := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: userID},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},

			Username: "testuser",
			Status:  usersModel.UserStatusActive,
			Roles:   []string{"reader"},
		}

		dto := ToUserDTO(userModel)

		assert.NotNil(t, dto)
		assert.Equal(t, userID.Hex(), dto.ID)
		assert.NotEmpty(t, dto.ID)
		assert.Len(t, dto.ID, 24) // ObjectID Hex 长度为 24
	})

	t.Run("零值 ObjectID 转换", func(t *testing.T) {
		userModel := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.ObjectID{}},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},

			Username: "testuser",
			Status:  usersModel.UserStatusActive,
			Roles:   []string{"reader"},
		}

		dto := ToUserDTO(userModel)

		assert.NotNil(t, dto)
		assert.Equal(t, "000000000000000000000000", dto.ID)
	})

	t.Run("批量转换 ID", func(t *testing.T) {
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()
		id3 := primitive.NewObjectID()

		users := []*usersModel.User{
			{
				IdentifiedEntity: shared.IdentifiedEntity{ID: id1},
				BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username:         "user1",
				Status:          usersModel.UserStatusActive,
				Roles:           []string{"reader"},
			},
			{
				IdentifiedEntity: shared.IdentifiedEntity{ID: id2},
				BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username:         "user2",
				Status:          usersModel.UserStatusActive,
				Roles:           []string{"reader"},
			},
			{
				IdentifiedEntity: shared.IdentifiedEntity{ID: id3},
				BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Username:         "user3",
				Status:          usersModel.UserStatusActive,
				Roles:           []string{"reader"},
			},
		}

		dtos := ToUserDTOs(users)

		assert.Len(t, dtos, 3)
		assert.Equal(t, id1.Hex(), dtos[0].ID)
		assert.Equal(t, id2.Hex(), dtos[1].ID)
		assert.Equal(t, id3.Hex(), dtos[2].ID)
	})

	t.Run("往返 ID 转换", func(t *testing.T) {
		originalID := primitive.NewObjectID()

		originalUser := &usersModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: originalID},
			BaseEntity:       shared.BaseEntity{CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Username:         "testuser",
			Status:          usersModel.UserStatusActive,
			Roles:           []string{"reader"},
		}

		// Model -> DTO
		userDTO := ToUserDTO(originalUser)

		// DTO -> Model
		convertedUser, err := ToUser(userDTO)

		assert.NoError(t, err)
		assert.Equal(t, originalID, convertedUser.ID)
	})
}
