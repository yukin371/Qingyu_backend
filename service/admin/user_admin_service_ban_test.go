package admin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	adminModel "Qingyu_backend/models/admin"
	"Qingyu_backend/models/users"
	adminrepo "Qingyu_backend/repository/interfaces/admin"
	base "Qingyu_backend/repository/interfaces/infrastructure"
)

type stubUserAdminRepo struct {
	getByIDFn func(ctx context.Context, userID string) (*users.User, error)
	updateFn  func(ctx context.Context, userID string, updates map[string]interface{}) error
}

// === BaseUserRepository 接口方法 (使用 string ID) ===

func (s *stubUserAdminRepo) Create(ctx context.Context, user *users.User) error { return nil }
func (s *stubUserAdminRepo) GetByID(ctx context.Context, id string) (*users.User, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, ErrUserNotFound
}
func (s *stubUserAdminRepo) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, id, updates)
	}
	return nil
}
func (s *stubUserAdminRepo) Delete(ctx context.Context, id string) error { return nil }
func (s *stubUserAdminRepo) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	return nil, nil
}
func (s *stubUserAdminRepo) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (s *stubUserAdminRepo) UpdateStatus(ctx context.Context, id string, status users.UserStatus) error {
	return nil
}
func (s *stubUserAdminRepo) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	return nil
}
func (s *stubUserAdminRepo) SetEmailVerified(ctx context.Context, id string, verified bool) error {
	return nil
}
func (s *stubUserAdminRepo) BatchUpdateStatus(ctx context.Context, ids []string, status users.UserStatus) error {
	return nil
}
func (s *stubUserAdminRepo) BatchDelete(ctx context.Context, ids []string) error {
	return nil
}
func (s *stubUserAdminRepo) CountByStatus(ctx context.Context, status users.UserStatus) (int64, error) {
	return 0, nil
}
func (s *stubUserAdminRepo) CountByRole(ctx context.Context, role string) (int64, error) {
	return 0, nil
}
func (s *stubUserAdminRepo) List(ctx context.Context, filter base.Filter) ([]*users.User, error) {
	return nil, nil
}
func (s *stubUserAdminRepo) Count(ctx context.Context, filter base.Filter) (int64, error) {
	return 0, nil
}
func (s *stubUserAdminRepo) Exists(ctx context.Context, id string) (bool, error) {
	return false, nil
}

// === admin 特有方法 ===

func (s *stubUserAdminRepo) ListWithPagination(ctx context.Context, filter *adminrepo.UserFilter, page, pageSize int) ([]*users.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserAdminRepo) BatchCreate(ctx context.Context, usersList []*users.User) error {
	return nil
}
func (s *stubUserAdminRepo) HardDelete(ctx context.Context, userID primitive.ObjectID) error {
	return nil
}
func (s *stubUserAdminRepo) GetActivities(ctx context.Context, userID primitive.ObjectID, page, pageSize int) ([]*users.UserActivity, int64, error) {
	return nil, 0, nil
}
func (s *stubUserAdminRepo) GetStatistics(ctx context.Context, userID primitive.ObjectID) (*users.UserStatistics, error) {
	return nil, nil
}
func (s *stubUserAdminRepo) ResetPassword(ctx context.Context, userID primitive.ObjectID, newPassword string) error {
	return nil
}
func (s *stubUserAdminRepo) UpdateRoles(ctx context.Context, userID primitive.ObjectID, role string) error {
	return nil
}
func (s *stubUserAdminRepo) SearchUsers(ctx context.Context, keyword string, page, pageSize int) ([]*users.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserAdminRepo) GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]*users.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserAdminRepo) CountByStatusMap(ctx context.Context) (map[string]int64, error) {
	return nil, nil
}
func (s *stubUserAdminRepo) GetRecentUsers(ctx context.Context, limit int) ([]*users.User, error) {
	return nil, nil
}
func (s *stubUserAdminRepo) GetActiveUsers(ctx context.Context, days int, limit int) ([]*users.User, error) {
	return nil, nil
}

type stubBanRecordRepo struct {
	createFn func(ctx context.Context, record *adminModel.BanRecord) error
}

func (s *stubBanRecordRepo) Create(ctx context.Context, record *adminModel.BanRecord) error {
	if s.createFn != nil {
		return s.createFn(ctx, record)
	}
	return nil
}
func (s *stubBanRecordRepo) GetByUserID(ctx context.Context, userID string, page, pageSize int) ([]*adminModel.BanRecord, int64, error) {
	return nil, 0, nil
}
func (s *stubBanRecordRepo) GetActiveBan(ctx context.Context, userID string) (*adminModel.BanRecord, error) {
	return nil, nil
}

func TestUpdateUserStatusWithReason_BanRequiresReason(t *testing.T) {
	testID := primitive.NewObjectID().Hex()
	repo := &stubUserAdminRepo{
		getByIDFn: func(ctx context.Context, userID string) (*users.User, error) {
			return &users.User{
				Username: "u1",
				Roles:    []string{"reader"},
				Status:   users.UserStatusActive,
			}, nil
		},
	}
	svc := &UserAdminServiceImpl{userRepo: repo}

	err := svc.UpdateUserStatusWithReason(context.Background(), testID, users.UserStatusBanned, "admin1", nil)
	assert.ErrorIs(t, err, ErrBanReasonRequired)
}

func TestUpdateUserStatusWithReason_BanSetsFieldsAndCreatesRecord(t *testing.T) {
	var updatedUpdates map[string]interface{}
	var record *adminModel.BanRecord
	reason := "恶意刷接口"

	repo := &stubUserAdminRepo{
		getByIDFn: func(ctx context.Context, userID string) (*users.User, error) {
			return &users.User{
				Username: "u1",
				Roles:    []string{"reader"},
				Status:   users.UserStatusActive,
			}, nil
		},
		updateFn: func(ctx context.Context, userID string, updates map[string]interface{}) error {
			updatedUpdates = updates
			return nil
		},
	}
	banRepo := &stubBanRecordRepo{
		createFn: func(ctx context.Context, r *adminModel.BanRecord) error {
			record = r
			return nil
		},
	}
	svc := &UserAdminServiceImpl{userRepo: repo, banRecordRepo: banRepo}

	testID := primitive.NewObjectID().Hex()
	err := svc.UpdateUserStatusWithReason(context.Background(), testID, users.UserStatusBanned, "admin1", &reason)
	assert.NoError(t, err)
	if assert.NotNil(t, updatedUpdates) {
		assert.NotNil(t, updatedUpdates["banned_at"])
		assert.Equal(t, "admin1", updatedUpdates["banned_by"])
		assert.Equal(t, reason, updatedUpdates["ban_reason"])
	}
	if assert.NotNil(t, record) {
		assert.Equal(t, "ban", record.Action)
		assert.Equal(t, reason, record.Reason)
		assert.Equal(t, "admin1", record.OperatorID)
	}
}

func TestUpdateUserStatusWithReason_UnbanClearsFieldsAndCreatesRecord(t *testing.T) {
	var updatedUpdates map[string]interface{}
	var record *adminModel.BanRecord
	now := time.Now()

	repo := &stubUserAdminRepo{
		getByIDFn: func(ctx context.Context, userID string) (*users.User, error) {
			return &users.User{
				Username:  "u1",
				Roles:     []string{"reader"},
				Status:    users.UserStatusBanned,
				BannedAt:  &now,
				BannedBy:  "admin0",
				BanReason: "旧原因",
			}, nil
		},
		updateFn: func(ctx context.Context, userID string, updates map[string]interface{}) error {
			updatedUpdates = updates
			return nil
		},
	}
	banRepo := &stubBanRecordRepo{
		createFn: func(ctx context.Context, r *adminModel.BanRecord) error {
			record = r
			return nil
		},
	}
	svc := &UserAdminServiceImpl{userRepo: repo, banRecordRepo: banRepo}

	testID := primitive.NewObjectID().Hex()
	err := svc.UpdateUserStatusWithReason(context.Background(), testID, users.UserStatusActive, "admin2", nil)
	assert.NoError(t, err)
	if assert.NotNil(t, updatedUpdates) {
		assert.Nil(t, updatedUpdates["banned_at"])
		assert.Equal(t, "", updatedUpdates["banned_by"])
		assert.Equal(t, "", updatedUpdates["ban_reason"])
	}
	if assert.NotNil(t, record) {
		assert.Equal(t, "unban", record.Action)
		assert.Equal(t, "解除封禁", record.Reason)
		assert.Equal(t, "admin2", record.OperatorID)
	}
}
