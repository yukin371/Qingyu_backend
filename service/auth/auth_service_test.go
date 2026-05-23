package auth

import (
	"context"
	"errors"
	"testing"

	authModel "Qingyu_backend/models/auth"
	usersModel "Qingyu_backend/models/users"
	authRepo "Qingyu_backend/repository/interfaces/auth"
	userRepo "Qingyu_backend/repository/interfaces/user"
	usermocks "Qingyu_backend/service/user/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

type jwtServiceStub struct {
	generateTokenFunc func(ctx context.Context, userID string, roles []string) (string, error)
	validateTokenFunc func(ctx context.Context, token string) (*TokenClaims, error)
	refreshTokenFunc  func(ctx context.Context, token string) (string, error)
	revokeTokenFunc   func(ctx context.Context, token string) error
}

func (s *jwtServiceStub) GenerateToken(ctx context.Context, userID string, roles []string) (string, error) {
	if s.generateTokenFunc != nil {
		return s.generateTokenFunc(ctx, userID, roles)
	}
	return "", nil
}

func (s *jwtServiceStub) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	if s.validateTokenFunc != nil {
		return s.validateTokenFunc(ctx, token)
	}
	return nil, nil
}

func (s *jwtServiceStub) RefreshToken(ctx context.Context, token string) (string, error) {
	if s.refreshTokenFunc != nil {
		return s.refreshTokenFunc(ctx, token)
	}
	return "", nil
}

func (s *jwtServiceStub) RevokeToken(ctx context.Context, token string) error {
	if s.revokeTokenFunc != nil {
		return s.revokeTokenFunc(ctx, token)
	}
	return nil
}

func (s *jwtServiceStub) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	return false, nil
}

type sessionServiceStub struct {
	createSessionFunc      func(ctx context.Context, userID string) (*Session, error)
	enforceDeviceLimitFunc func(ctx context.Context, userID string, maxDevices int) error
}

func (s *sessionServiceStub) CreateSession(ctx context.Context, userID string) (*Session, error) {
	if s.createSessionFunc != nil {
		return s.createSessionFunc(ctx, userID)
	}
	return nil, nil
}

func (s *sessionServiceStub) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	return nil, nil
}

func (s *sessionServiceStub) UpdateSession(ctx context.Context, sessionID string, data map[string]interface{}) error {
	return nil
}

func (s *sessionServiceStub) DestroySession(ctx context.Context, sessionID string) error {
	return nil
}

func (s *sessionServiceStub) RefreshSession(ctx context.Context, sessionID string) error {
	return nil
}

func (s *sessionServiceStub) CheckDeviceLimit(ctx context.Context, userID string, maxDevices int) error {
	return nil
}

func (s *sessionServiceStub) EnforceDeviceLimit(ctx context.Context, userID string, maxDevices int) error {
	if s.enforceDeviceLimitFunc != nil {
		return s.enforceDeviceLimitFunc(ctx, userID, maxDevices)
	}
	return nil
}

func (s *sessionServiceStub) GetUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	return nil, nil
}

func (s *sessionServiceStub) DestroyUserSessions(ctx context.Context, userID string) error {
	return nil
}

type oauthRepositoryStub struct {
	findByProviderAndProviderIDFunc func(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error)
	createFunc                      func(ctx context.Context, account *authModel.OAuthAccount) error
	updateLastLoginFunc             func(ctx context.Context, id string) error
}

func (s *oauthRepositoryStub) FindByProviderAndProviderID(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
	if s.findByProviderAndProviderIDFunc != nil {
		return s.findByProviderAndProviderIDFunc(ctx, provider, providerUserID)
	}
	return nil, nil
}

func (s *oauthRepositoryStub) FindByUserID(ctx context.Context, userID string) ([]*authModel.OAuthAccount, error) {
	return nil, nil
}

func (s *oauthRepositoryStub) FindByID(ctx context.Context, id string) (*authModel.OAuthAccount, error) {
	return nil, nil
}

func (s *oauthRepositoryStub) Create(ctx context.Context, account *authModel.OAuthAccount) error {
	if s.createFunc != nil {
		return s.createFunc(ctx, account)
	}
	return nil
}

func (s *oauthRepositoryStub) Update(ctx context.Context, account *authModel.OAuthAccount) error {
	return nil
}

func (s *oauthRepositoryStub) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *oauthRepositoryStub) UpdateLastLogin(ctx context.Context, id string) error {
	if s.updateLastLoginFunc != nil {
		return s.updateLastLoginFunc(ctx, id)
	}
	return nil
}

func (s *oauthRepositoryStub) UpdateTokens(ctx context.Context, id string, accessToken, refreshToken string, expiresAt primitive.DateTime) error {
	return nil
}

func (s *oauthRepositoryStub) SetPrimaryAccount(ctx context.Context, userID string, accountID string) error {
	return nil
}

func (s *oauthRepositoryStub) GetPrimaryAccount(ctx context.Context, userID string) (*authModel.OAuthAccount, error) {
	return nil, nil
}

func (s *oauthRepositoryStub) CountByUserID(ctx context.Context, userID string) (int64, error) {
	return 0, nil
}

func (s *oauthRepositoryStub) CreateSession(ctx context.Context, session *authModel.OAuthSession) error {
	return nil
}

func (s *oauthRepositoryStub) FindSessionByID(ctx context.Context, id string) (*authModel.OAuthSession, error) {
	return nil, nil
}

func (s *oauthRepositoryStub) FindSessionByState(ctx context.Context, state string) (*authModel.OAuthSession, error) {
	return nil, nil
}

func (s *oauthRepositoryStub) DeleteSession(ctx context.Context, id string) error {
	return nil
}

func (s *oauthRepositoryStub) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	return 0, nil
}

var _ authRepo.OAuthRepository = (*oauthRepositoryStub)(nil)
var _ userRepo.UserRepository = (*usermocks.MockUserRepository)(nil)

func TestAuthService_Logout_WrapsRevokeError(t *testing.T) {
	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			revokeTokenFunc: func(ctx context.Context, token string) error {
				assert.Equal(t, "token-123", token)
				return errors.New("blacklist unavailable")
			},
		},
	}

	err := service.Logout(context.Background(), "token-123")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "登出失败")
	assert.Contains(t, err.Error(), "blacklist unavailable")
}

func TestAuthService_RefreshToken_WrapsJWTError(t *testing.T) {
	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			refreshTokenFunc: func(ctx context.Context, token string) (string, error) {
				assert.Equal(t, "refresh-123", token)
				return "", errors.New("token expired")
			},
		},
	}

	newToken, err := service.RefreshToken(context.Background(), "refresh-123")

	require.Error(t, err)
	assert.Empty(t, newToken)
	assert.Contains(t, err.Error(), "刷新Token失败")
	assert.Contains(t, err.Error(), "token expired")
}

func TestAuthService_ValidateToken_ReturnsClaims(t *testing.T) {
	expectedClaims := &TokenClaims{
		UserID: "user-123",
		Roles:  []string{"reader", "author"},
		Exp:    1234567890,
		Iat:    1234567000,
	}

	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			validateTokenFunc: func(ctx context.Context, token string) (*TokenClaims, error) {
				assert.Equal(t, "access-token-123", token)
				return expectedClaims, nil
			},
		},
	}

	actualClaims, err := service.ValidateToken(context.Background(), "access-token-123")

	require.NoError(t, err)
	assert.Equal(t, expectedClaims, actualClaims)
}

func TestAuthService_Logout_SucceedsWhenTokenRevoked(t *testing.T) {
	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			revokeTokenFunc: func(ctx context.Context, token string) error {
				assert.Equal(t, "token-success", token)
				return nil
			},
		},
	}

	err := service.Logout(context.Background(), "token-success")

	require.NoError(t, err)
}

func TestAuthService_Login_SessionFailuresDoNotBlockLogin(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()

	user := &usersModel.User{
		Username: "reader-login",
		Email:    "reader@example.com",
		Roles:    []string{"reader"},
		Status:   usersModel.UserStatusActive,
	}
	user.ID = userID
	require.NoError(t, user.SetPassword("Password123!"))

	mockUserRepo := &usermocks.MockUserRepository{}
	mockAuthRepo := usermocks.NewMockAuthRepository()
	mockUserRepo.On("GetByUsername", mock.Anything, "reader-login").Return(user, nil).Once()
	mockUserRepo.On("UpdateLastLogin", mock.Anything, userID.Hex(), "unknown").Return(nil).Once()
	mockAuthRepo.On("GetUserRoles", mock.Anything, userID.Hex()).Return([]*authModel.Role{
		{Name: "reader"},
	}, nil).Once()

	deviceLimitCalled := false
	createSessionCalled := false
	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			generateTokenFunc: func(ctx context.Context, actualUserID string, roles []string) (string, error) {
				assert.Equal(t, userID.Hex(), actualUserID)
				assert.Equal(t, []string{"reader"}, roles)
				return "login-token", nil
			},
		},
		authRepo: mockAuthRepo,
		userRepo: mockUserRepo,
		sessionService: &sessionServiceStub{
			enforceDeviceLimitFunc: func(ctx context.Context, actualUserID string, maxDevices int) error {
				deviceLimitCalled = true
				assert.Equal(t, userID.Hex(), actualUserID)
				assert.Equal(t, 5, maxDevices)
				return errors.New("device tracker unavailable")
			},
			createSessionFunc: func(ctx context.Context, actualUserID string) (*Session, error) {
				createSessionCalled = true
				assert.Equal(t, userID.Hex(), actualUserID)
				return nil, errors.New("session store unavailable")
			},
		},
	}

	resp, err := service.Login(ctx, &LoginRequest{
		Username: "reader-login",
		Password: "Password123!",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, deviceLimitCalled)
	assert.True(t, createSessionCalled)
	assert.Equal(t, "login-token", resp.Token)
	assert.Equal(t, userID.Hex(), resp.User.ID)
	assert.Equal(t, []string{"reader"}, resp.User.Roles)
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_OAuthLogin_ExistingAccountUsesResolvedRoles(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	accountID := primitive.NewObjectID()

	mockUserRepo := &usermocks.MockUserRepository{}
	mockAuthRepo := usermocks.NewMockAuthRepository()
	oauthRepo := &oauthRepositoryStub{
		findByProviderAndProviderIDFunc: func(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
			assert.Equal(t, authModel.OAuthProviderGoogle, provider)
			assert.Equal(t, "google-user-1", providerUserID)
			return &authModel.OAuthAccount{
				ID:             accountID,
				UserID:         userID.Hex(),
				Provider:       provider,
				ProviderUserID: providerUserID,
			}, nil
		},
		updateLastLoginFunc: func(ctx context.Context, id string) error {
			assert.Equal(t, accountID.Hex(), id)
			return nil
		},
	}

	user := &usersModel.User{
		Username: "oauth-user",
		Email:    "oauth@example.com",
		Roles:    []string{"reader"},
		Status:   usersModel.UserStatusActive,
	}
	user.ID = userID

	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(user, nil).Once()
	mockAuthRepo.On("GetUserRoles", mock.Anything, userID.Hex()).Return([]*authModel.Role{
		{Name: "author"},
		{Name: "reader"},
	}, nil).Once()

	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			generateTokenFunc: func(ctx context.Context, actualUserID string, roles []string) (string, error) {
				assert.Equal(t, userID.Hex(), actualUserID)
				assert.Equal(t, []string{"author", "reader"}, roles)
				return "jwt-token", nil
			},
		},
		authRepo:       mockAuthRepo,
		oauthRepo:      oauthRepo,
		userRepo:       mockUserRepo,
		sessionService: &sessionServiceStub{},
	}

	resp, err := service.OAuthLogin(ctx, &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google-user-1",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "jwt-token", resp.Token)
	assert.Equal(t, userID.Hex(), resp.User.ID)
	assert.Equal(t, []string{"author", "reader"}, resp.User.Roles)
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_OAuthLogin_NewAccountCreatesReaderUser(t *testing.T) {
	ctx := context.Background()
	newUserID := primitive.NewObjectID()

	mockUserRepo := &usermocks.MockUserRepository{}
	mockAuthRepo := usermocks.NewMockAuthRepository()
	oauthRepo := &oauthRepositoryStub{
		findByProviderAndProviderIDFunc: func(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
			return nil, nil
		},
		createFunc: func(ctx context.Context, account *authModel.OAuthAccount) error {
			assert.Equal(t, newUserID.Hex(), account.UserID)
			assert.Equal(t, authModel.OAuthProviderGitHub, account.Provider)
			assert.Equal(t, "provider-user-99", account.ProviderUserID)
			assert.Equal(t, "oauth@example.com", account.Email)
			assert.True(t, account.IsPrimary)
			return nil
		},
	}

	mockUserRepo.On("ExistsByUsername", mock.Anything, mock.MatchedBy(func(username string) bool {
		return len(username) > 0 && username[:7] == "github_"
	})).Return(false, nil).Once()
	mockUserRepo.On("ExistsByEmail", mock.Anything, "oauth@example.com").Return(false, nil).Once()
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *usersModel.User) bool {
		user.ID = newUserID
		return user.Email == "oauth@example.com" &&
			len(user.Username) > 0 &&
			user.Username[:7] == "github_" &&
			user.Password != "" &&
			user.Status == usersModel.UserStatusActive &&
			assert.ObjectsAreEqual([]string{"reader"}, user.Roles)
	})).Return(nil).Once()

	roleID := primitive.NewObjectID()
	mockAuthRepo.On("GetRoleByName", mock.Anything, "reader").Return(&authModel.Role{ID: roleID, Name: "reader"}, nil).Once()
	mockAuthRepo.On("AssignUserRole", mock.Anything, newUserID.Hex(), roleID.Hex()).Return(nil).Once()

	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			generateTokenFunc: func(ctx context.Context, actualUserID string, roles []string) (string, error) {
				assert.Equal(t, newUserID.Hex(), actualUserID)
				assert.Equal(t, []string{"reader"}, roles)
				return "new-user-token", nil
			},
		},
		authRepo:       mockAuthRepo,
		oauthRepo:      oauthRepo,
		userRepo:       mockUserRepo,
		sessionService: &sessionServiceStub{},
	}

	resp, err := service.OAuthLogin(ctx, &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGitHub,
		ProviderID: "provider-user-99",
		Email:      "oauth@example.com",
		Avatar:     "https://example.com/avatar.png",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "new-user-token", resp.Token)
	assert.Equal(t, newUserID.Hex(), resp.User.ID)
	assert.Equal(t, []string{"reader"}, resp.User.Roles)
	assert.Equal(t, "oauth@example.com", resp.User.Email)
	assert.Equal(t, "github_provider", resp.User.Username)
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_OAuthLogin_ExistingAccountGenerateTokenFailureReturnsWrappedError(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	accountID := primitive.NewObjectID()

	mockUserRepo := &usermocks.MockUserRepository{}
	mockAuthRepo := usermocks.NewMockAuthRepository()
	oauthRepo := &oauthRepositoryStub{
		findByProviderAndProviderIDFunc: func(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
			return &authModel.OAuthAccount{
				ID:             accountID,
				UserID:         userID.Hex(),
				Provider:       provider,
				ProviderUserID: providerUserID,
			}, nil
		},
		updateLastLoginFunc: func(ctx context.Context, id string) error {
			assert.Equal(t, accountID.Hex(), id)
			return nil
		},
	}

	user := &usersModel.User{
		Username: "oauth-user",
		Email:    "oauth@example.com",
		Roles:    []string{"reader"},
		Status:   usersModel.UserStatusActive,
	}
	user.ID = userID

	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(user, nil).Once()
	mockAuthRepo.On("GetUserRoles", mock.Anything, userID.Hex()).Return([]*authModel.Role{
		{Name: "author"},
		{Name: "reader"},
	}, nil).Once()

	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			generateTokenFunc: func(ctx context.Context, actualUserID string, roles []string) (string, error) {
				assert.Equal(t, userID.Hex(), actualUserID)
				assert.Equal(t, []string{"author", "reader"}, roles)
				return "", errors.New("jwt signer unavailable")
			},
		},
		authRepo:       mockAuthRepo,
		oauthRepo:      oauthRepo,
		userRepo:       mockUserRepo,
		sessionService: &sessionServiceStub{},
	}

	resp, err := service.OAuthLogin(ctx, &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google-user-token-fail",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "生成Token失败")
	assert.Contains(t, err.Error(), "jwt signer unavailable")
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_OAuthLogin_NewAccountCreateOAuthRecordFailureReturnsWrappedError(t *testing.T) {
	ctx := context.Background()
	newUserID := primitive.NewObjectID()

	mockUserRepo := &usermocks.MockUserRepository{}
	mockAuthRepo := usermocks.NewMockAuthRepository()
	oauthRepo := &oauthRepositoryStub{
		findByProviderAndProviderIDFunc: func(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
			return nil, nil
		},
		createFunc: func(ctx context.Context, account *authModel.OAuthAccount) error {
			assert.Equal(t, newUserID.Hex(), account.UserID)
			assert.Equal(t, authModel.OAuthProviderGitHub, account.Provider)
			assert.Equal(t, "provider-user-create-fail", account.ProviderUserID)
			return errors.New("oauth repo down")
		},
	}

	mockUserRepo.On("ExistsByUsername", mock.Anything, mock.MatchedBy(func(username string) bool {
		return len(username) > 0 && username[:7] == "github_"
	})).Return(false, nil).Once()
	mockUserRepo.On("ExistsByEmail", mock.Anything, "oauth@example.com").Return(false, nil).Once()
	mockUserRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *usersModel.User) bool {
		user.ID = newUserID
		return user.Email == "oauth@example.com" &&
			len(user.Username) > 0 &&
			user.Username[:7] == "github_" &&
			user.Password != "" &&
			user.Status == usersModel.UserStatusActive &&
			assert.ObjectsAreEqual([]string{"reader"}, user.Roles)
	})).Return(nil).Once()

	roleID := primitive.NewObjectID()
	mockAuthRepo.On("GetRoleByName", mock.Anything, "reader").Return(&authModel.Role{ID: roleID, Name: "reader"}, nil).Once()
	mockAuthRepo.On("AssignUserRole", mock.Anything, newUserID.Hex(), roleID.Hex()).Return(nil).Once()

	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			generateTokenFunc: func(ctx context.Context, actualUserID string, roles []string) (string, error) {
				t.Fatalf("GenerateToken should not run when oauthRepo.Create fails")
				return "", nil
			},
		},
		authRepo:       mockAuthRepo,
		oauthRepo:      oauthRepo,
		userRepo:       mockUserRepo,
		sessionService: &sessionServiceStub{},
	}

	resp, err := service.OAuthLogin(ctx, &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGitHub,
		ProviderID: "provider-user-create-fail",
		Email:      "oauth@example.com",
		Avatar:     "https://example.com/avatar.png",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "创建OAuth账号失败")
	assert.Contains(t, err.Error(), "oauth repo down")
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_OAuthLogin_ExistingAccountResolveRolesFailureReturnsWrappedError(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	accountID := primitive.NewObjectID()

	mockUserRepo := &usermocks.MockUserRepository{}
	mockAuthRepo := usermocks.NewMockAuthRepository()
	oauthRepo := &oauthRepositoryStub{
		findByProviderAndProviderIDFunc: func(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
			return &authModel.OAuthAccount{
				ID:             accountID,
				UserID:         userID.Hex(),
				Provider:       provider,
				ProviderUserID: providerUserID,
			}, nil
		},
		updateLastLoginFunc: func(ctx context.Context, id string) error {
			assert.Equal(t, accountID.Hex(), id)
			return nil
		},
	}

	user := &usersModel.User{
		Username: "oauth-user",
		Email:    "oauth@example.com",
		Roles:    []string{"reader"},
		Status:   usersModel.UserStatusActive,
	}
	user.ID = userID

	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(user, nil).Once()
	mockAuthRepo.On("GetUserRoles", mock.Anything, userID.Hex()).Return(nil, errors.New("role repo timeout")).Once()

	service := &AuthServiceImpl{
		jwtService: &jwtServiceStub{
			generateTokenFunc: func(ctx context.Context, actualUserID string, roles []string) (string, error) {
				t.Fatalf("GenerateToken should not run when resolveUserRoles fails")
				return "", nil
			},
		},
		authRepo:       mockAuthRepo,
		oauthRepo:      oauthRepo,
		userRepo:       mockUserRepo,
		sessionService: &sessionServiceStub{},
	}

	resp, err := service.OAuthLogin(ctx, &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google-user-role-fail",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "获取用户角色失败")
	assert.Contains(t, err.Error(), "role repo timeout")
	mockUserRepo.AssertExpectations(t)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_OAuthLogin_ExistingAccountGetUserFailureReturnsWrappedError(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	accountID := primitive.NewObjectID()

	mockUserRepo := &usermocks.MockUserRepository{}
	oauthRepo := &oauthRepositoryStub{
		findByProviderAndProviderIDFunc: func(ctx context.Context, provider authModel.OAuthProvider, providerUserID string) (*authModel.OAuthAccount, error) {
			assert.Equal(t, authModel.OAuthProviderGoogle, provider)
			assert.Equal(t, "google-user-fail", providerUserID)
			return &authModel.OAuthAccount{
				ID:             accountID,
				UserID:         userID.Hex(),
				Provider:       provider,
				ProviderUserID: providerUserID,
			}, nil
		},
		updateLastLoginFunc: func(ctx context.Context, id string) error {
			t.Fatalf("UpdateLastLogin should not run when GetByID fails")
			return nil
		},
	}

	mockUserRepo.On("GetByID", mock.Anything, userID.Hex()).Return(nil, errors.New("user repo offline")).Once()

	service := &AuthServiceImpl{
		jwtService:      &jwtServiceStub{},
		oauthRepo:       oauthRepo,
		userRepo:        mockUserRepo,
		sessionService:  &sessionServiceStub{},
	}

	resp, err := service.OAuthLogin(ctx, &OAuthLoginRequest{
		Provider:   authModel.OAuthProviderGoogle,
		ProviderID: "google-user-fail",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "获取用户信息失败")
	assert.Contains(t, err.Error(), "user repo offline")
	mockUserRepo.AssertExpectations(t)
}

func TestGenerateUsernameFromProvider(t *testing.T) {
	assert.Equal(t, "github_provider", generateUsernameFromProvider(authModel.OAuthProviderGitHub, "provider-user-99"))
	assert.Equal(t, "oauth_abcdef", generateUsernameFromProvider("unknown", "abcdef"))
	assert.Equal(t, "qq_short", generateUsernameFromProvider(authModel.OAuthProviderQQ, "short"))
}

var _ oauth2.Token
