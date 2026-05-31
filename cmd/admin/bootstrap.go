package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"Qingyu_backend/config"
	authModel "Qingyu_backend/models/auth"
	"Qingyu_backend/models/shared"
	userModel "Qingyu_backend/models/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultBootstrapPasswordEnv = "QINGYU_BOOTSTRAP_ADMIN_PASSWORD"
	defaultBootstrapOperation   = "bootstrap_admin"
)

type bootstrapOptions struct {
	Username           string
	Email              string
	ConfigPath         string
	PasswordEnv        string
	ForceResetPassword bool
	DryRun             bool
}

type bootstrapExecutor struct {
	now func() time.Time
}

func shouldCreateBootstrapAdmin(findErr error, adminCount int64) (bool, error) {
	switch {
	case errors.Is(findErr, mongo.ErrNoDocuments):
		if adminCount > 0 {
			return false, fmt.Errorf("系统中已存在 %d 个管理员，拒绝继续初始化；如需重置现有账号，请使用 --force-reset-password", adminCount)
		}
		return true, nil
	case findErr != nil:
		return false, fmt.Errorf("查询管理员账号失败: %w", findErr)
	default:
		return false, nil
	}
}

func runBootstrap(args []string) error {
	fs := newBootstrapFlagSet()

	opts := bootstrapOptions{}
	fs.StringVar(&opts.Username, "username", "", "管理员用户名")
	fs.StringVar(&opts.Email, "email", "", "管理员邮箱")
	fs.StringVar(&opts.ConfigPath, "config", "", "配置文件路径，默认自动发现 configs/config.yaml")
	fs.StringVar(&opts.PasswordEnv, "password-env", defaultBootstrapPasswordEnv, "读取密码的环境变量名")
	fs.BoolVar(&opts.ForceResetPassword, "force-reset-password", false, "若管理员已存在，则重置同名管理员密码并补齐 admin 角色")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "仅校验参数与配置发现逻辑，不执行数据库写入")

	if err := fs.Parse(args); err != nil {
		return err
	}

	executor := &bootstrapExecutor{now: time.Now}
	return executor.execute(context.Background(), opts)
}

func (e *bootstrapExecutor) execute(ctx context.Context, opts bootstrapOptions) error {
	if err := validateBootstrapOptions(opts); err != nil {
		return err
	}

	password, err := resolveBootstrapPassword(opts.PasswordEnv, os.Getenv)
	if err != nil {
		return err
	}

	configPath, err := detectBootstrapConfigPath(opts.ConfigPath)
	if err != nil {
		return err
	}

	if opts.DryRun {
		fmt.Printf("Dry run 成功。\n配置文件: %s\n管理员用户名: %s\n管理员邮箱: %s\n密码环境变量: %s\n",
			configPath, opts.Username, opts.Email, opts.PasswordEnv)
		return nil
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	if cfg.Database == nil {
		return errors.New("数据库配置缺失")
	}

	mongoCfg, err := cfg.Database.GetMongoConfig()
	if err != nil {
		return fmt.Errorf("解析 MongoDB 配置失败: %w", err)
	}

	clientOpts := options.Client().
		ApplyURI(mongoCfg.URI).
		SetConnectTimeout(mongoCfg.ConnectTimeout).
		SetServerSelectionTimeout(mongoCfg.ServerTimeout)

	connectCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, clientOpts)
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %w", err)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
	}()

	if err := client.Ping(connectCtx, nil); err != nil {
		return fmt.Errorf("MongoDB ping 失败: %w", err)
	}

	db := client.Database(mongoCfg.Database)
	return e.bootstrapAdmin(ctx, db, opts, password)
}

func (e *bootstrapExecutor) bootstrapAdmin(ctx context.Context, db *mongo.Database, opts bootstrapOptions, password string) error {
	usersCollection := db.Collection("users")
	adminLogsCollection := db.Collection("admin_logs")

	var existing userModel.User
	findErr := usersCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": opts.Username},
			{"email": opts.Email},
		},
	}).Decode(&existing)

	adminCount, err := usersCollection.CountDocuments(ctx, bson.M{"roles": authModel.RoleAdmin})
	if err != nil {
		return fmt.Errorf("统计现有管理员失败: %w", err)
	}

	now := e.now()
	shouldCreate, err := shouldCreateBootstrapAdmin(findErr, adminCount)
	if err != nil {
		return err
	}

	switch {
	case shouldCreate:
		adminUser := userModel.User{
			IdentifiedEntity: shared.IdentifiedEntity{ID: primitive.NewObjectID()},
			BaseEntity:       shared.BaseEntity{CreatedAt: now, UpdatedAt: now},
			Username:         opts.Username,
			Email:            opts.Email,
			Roles:            []string{authModel.RoleAdmin},
			VIPLevel:         0,
			Status:           userModel.UserStatusActive,
			Nickname:         opts.Username,
			EmailVerified:    true,
		}
		if err := adminUser.SetPassword(password); err != nil {
			return fmt.Errorf("设置管理员密码失败: %w", err)
		}

		if _, err := usersCollection.InsertOne(ctx, adminUser); err != nil {
			return fmt.Errorf("创建管理员账号失败: %w", err)
		}

		_ = writeBootstrapAdminLog(ctx, adminLogsCollection, adminUser.ID.Hex(), adminUser.Username, "create", now)

		fmt.Printf("管理员初始化成功。\n用户名: %s\n邮箱: %s\n用户ID: %s\n",
			adminUser.Username, adminUser.Email, adminUser.ID.Hex())
		return nil

	default:
		if existing.Username != opts.Username || existing.Email != opts.Email {
			return fmt.Errorf("已存在冲突账号（username=%s, email=%s），请确认后再执行初始化", existing.Username, existing.Email)
		}
		if !opts.ForceResetPassword {
			return fmt.Errorf("管理员账号已存在（userID=%s），如需重置密码请追加 --force-reset-password", existing.ID.Hex())
		}

		hashedUser := existing
		if err := hashedUser.SetPassword(password); err != nil {
			return fmt.Errorf("重置管理员密码失败: %w", err)
		}

		newRoles := mergeAdminRoles(existing.Roles)
		update := bson.M{
			"$set": bson.M{
				"password":       hashedUser.Password,
				"roles":          newRoles,
				"status":         userModel.UserStatusActive,
				"email_verified": true,
				"updated_at":     now,
			},
		}
		if _, err := usersCollection.UpdateByID(ctx, existing.ID, update); err != nil {
			return fmt.Errorf("更新管理员账号失败: %w", err)
		}

		_ = writeBootstrapAdminLog(ctx, adminLogsCollection, existing.ID.Hex(), existing.Username, "reset_password", now)

		fmt.Printf("管理员账号已更新。\n用户名: %s\n邮箱: %s\n用户ID: %s\n",
			existing.Username, existing.Email, existing.ID.Hex())
		return nil
	}
}

func validateBootstrapOptions(opts bootstrapOptions) error {
	if strings.TrimSpace(opts.Username) == "" {
		return errors.New("username 不能为空")
	}
	if strings.TrimSpace(opts.Email) == "" {
		return errors.New("email 不能为空")
	}
	if !strings.Contains(opts.Email, "@") {
		return errors.New("email 格式无效")
	}
	if strings.TrimSpace(opts.PasswordEnv) == "" {
		return errors.New("password-env 不能为空")
	}
	return nil
}

func resolveBootstrapPassword(envName string, getenv func(string) string) (string, error) {
	password := strings.TrimSpace(getenv(envName))
	if password == "" {
		return "", fmt.Errorf("环境变量 %s 未设置，请先注入初始管理员密码", envName)
	}
	if len(password) < 8 {
		return "", fmt.Errorf("环境变量 %s 中的密码长度不能小于 8 位", envName)
	}
	return password, nil
}

func detectBootstrapConfigPath(provided string) (string, error) {
	if provided != "" {
		path := normalizeExistingPath(provided)
		if fileExists(path) {
			return path, nil
		}
		return "", fmt.Errorf("指定的配置文件不存在: %s", path)
	}

	for _, envKey := range []string{"CONFIG_FILE", "CONFIG_PATH"} {
		if candidate := strings.TrimSpace(os.Getenv(envKey)); candidate != "" && fileExists(candidate) {
			return normalizeExistingPath(candidate), nil
		}
	}

	candidates := []string{
		"/app/configs/config.yaml",
		"/app/config/config.yaml",
		"./configs/config.yaml",
		"./config/config.yaml",
		"../../configs/config.yaml",
		"../../config/config.yaml",
	}
	for _, candidate := range candidates {
		if fileExists(candidate) {
			return normalizeExistingPath(candidate), nil
		}
	}

	return "", errors.New("未找到可用配置文件，请通过 --config 显式指定")
}

func normalizeExistingPath(path string) string {
	cleaned := filepath.Clean(path)
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return cleaned
	}
	return abs
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func mergeAdminRoles(roles []string) []string {
	seen := make(map[string]struct{}, len(roles)+1)
	result := make([]string, 0, len(roles)+1)

	for _, role := range roles {
		role = strings.TrimSpace(role)
		if role == "" {
			continue
		}
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		result = append(result, role)
	}

	if _, ok := seen[authModel.RoleAdmin]; !ok {
		result = append(result, authModel.RoleAdmin)
	}

	return result
}

func writeBootstrapAdminLog(ctx context.Context, collection *mongo.Collection, adminID, adminName, action string, now time.Time) error {
	_, err := collection.InsertOne(ctx, bson.M{
		"admin_id":      adminID,
		"admin_name":    adminName,
		"operation":     defaultBootstrapOperation,
		"target":        adminID,
		"target_type":   "user",
		"resource_type": "user",
		"resource_id":   adminID,
		"details": bson.M{
			"action": action,
			"source": "admin-cli",
		},
		"ip":         "bootstrap-cli",
		"user_agent": "cmd/admin",
		"created_at": now,
	})
	return err
}
