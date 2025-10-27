package storage

import (
	"fmt"
)

// BackendType 存储后端类型
type BackendType string

const (
	BackendTypeLocal BackendType = "local"
	BackendTypeMinIO BackendType = "minio"
	BackendTypeOSS   BackendType = "oss" // 阿里云OSS (未实现)
	BackendTypeCOS   BackendType = "cos" // 腾讯云COS (未实现)
	BackendTypeS3    BackendType = "s3"  // AWS S3 (未实现)
)

// BackendConfig 存储后端配置接口
type BackendConfig interface {
	GetType() BackendType
}

// LocalBackendConfig 本地存储配置
type LocalBackendConfig struct {
	BasePath string `json:"base_path" yaml:"base_path"` // 存储根目录
	BaseURL  string `json:"base_url" yaml:"base_url"`   // 访问基础URL
}

func (c *LocalBackendConfig) GetType() BackendType {
	return BackendTypeLocal
}

// MinIOBackendConfig MinIO存储配置
type MinIOBackendConfig struct {
	Endpoint   string `json:"endpoint" yaml:"endpoint"`       // MinIO服务器地址
	AccessKey  string `json:"access_key" yaml:"access_key"`   // 访问密钥
	SecretKey  string `json:"secret_key" yaml:"secret_key"`   // 密钥
	BucketName string `json:"bucket_name" yaml:"bucket_name"` // 存储桶名称
	UseSSL     bool   `json:"use_ssl" yaml:"use_ssl"`         // 是否使用SSL
	Location   string `json:"location" yaml:"location"`       // 存储桶位置
}

func (c *MinIOBackendConfig) GetType() BackendType {
	return BackendTypeMinIO
}

// BackendFactory 存储后端工厂
type BackendFactory struct {
}

// NewBackendFactory 创建存储后端工厂
func NewBackendFactory() *BackendFactory {
	return &BackendFactory{}
}

// CreateBackend 根据配置创建存储后端
func (f *BackendFactory) CreateBackend(config BackendConfig) (StorageBackend, error) {
	if config == nil {
		return nil, fmt.Errorf("backend config is nil")
	}

	switch config.GetType() {
	case BackendTypeLocal:
		return f.createLocalBackend(config)
	case BackendTypeMinIO:
		return f.createMinIOBackend(config)
	case BackendTypeOSS:
		return nil, fmt.Errorf("OSS backend not implemented yet")
	case BackendTypeCOS:
		return nil, fmt.Errorf("COS backend not implemented yet")
	case BackendTypeS3:
		return nil, fmt.Errorf("S3 backend not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported backend type: %s", config.GetType())
	}
}

// createLocalBackend 创建本地存储后端
func (f *BackendFactory) createLocalBackend(config BackendConfig) (StorageBackend, error) {
	localConfig, ok := config.(*LocalBackendConfig)
	if !ok {
		return nil, fmt.Errorf("invalid local backend config")
	}

	if localConfig.BasePath == "" {
		return nil, fmt.Errorf("base_path is required for local backend")
	}

	return NewLocalBackend(localConfig.BasePath, localConfig.BaseURL), nil
}

// createMinIOBackend 创建MinIO存储后端
func (f *BackendFactory) createMinIOBackend(config BackendConfig) (StorageBackend, error) {
	minioConfig, ok := config.(*MinIOBackendConfig)
	if !ok {
		return nil, fmt.Errorf("invalid MinIO backend config")
	}

	if minioConfig.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required for MinIO backend")
	}
	if minioConfig.AccessKey == "" {
		return nil, fmt.Errorf("access_key is required for MinIO backend")
	}
	if minioConfig.SecretKey == "" {
		return nil, fmt.Errorf("secret_key is required for MinIO backend")
	}
	if minioConfig.BucketName == "" {
		return nil, fmt.Errorf("bucket_name is required for MinIO backend")
	}

	return NewMinIOBackend(&MinIOConfig{
		Endpoint:   minioConfig.Endpoint,
		AccessKey:  minioConfig.AccessKey,
		SecretKey:  minioConfig.SecretKey,
		BucketName: minioConfig.BucketName,
		UseSSL:     minioConfig.UseSSL,
		Location:   minioConfig.Location,
	})
}

// GetSupportedBackends 获取支持的存储后端列表
func (f *BackendFactory) GetSupportedBackends() []BackendType {
	return []BackendType{
		BackendTypeLocal,
		BackendTypeMinIO,
		// BackendTypeOSS,  // 待实现
		// BackendTypeCOS,  // 待实现
		// BackendTypeS3,   // 待实现
	}
}

// ValidateConfig 验证存储后端配置
func (f *BackendFactory) ValidateConfig(config BackendConfig) error {
	if config == nil {
		return fmt.Errorf("backend config is nil")
	}

	switch config.GetType() {
	case BackendTypeLocal:
		localConfig, ok := config.(*LocalBackendConfig)
		if !ok {
			return fmt.Errorf("invalid local backend config type")
		}
		if localConfig.BasePath == "" {
			return fmt.Errorf("base_path is required")
		}
		return nil

	case BackendTypeMinIO:
		minioConfig, ok := config.(*MinIOBackendConfig)
		if !ok {
			return fmt.Errorf("invalid MinIO backend config type")
		}
		if minioConfig.Endpoint == "" {
			return fmt.Errorf("endpoint is required")
		}
		if minioConfig.AccessKey == "" {
			return fmt.Errorf("access_key is required")
		}
		if minioConfig.SecretKey == "" {
			return fmt.Errorf("secret_key is required")
		}
		if minioConfig.BucketName == "" {
			return fmt.Errorf("bucket_name is required")
		}
		return nil

	default:
		return fmt.Errorf("unsupported backend type: %s", config.GetType())
	}
}
