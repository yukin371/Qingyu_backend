package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOBackend MinIO对象存储后端
type MinIOBackend struct {
	client     *minio.Client
	bucketName string
	useSSL     bool
	location   string // 存储桶位置 (例如: us-east-1)
}

// 编译期契约校验：MinIOBackend 必须实现 StorageBackend 端口。
var _ StorageBackend = (*MinIOBackend)(nil)

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Endpoint   string // MinIO服务器地址 (例如: localhost:9000)
	AccessKey  string // 访问密钥
	SecretKey  string // 密钥
	BucketName string // 存储桶名称
	UseSSL     bool   // 是否使用SSL
	Location   string // 存储桶位置 (默认: us-east-1)
}

// NewMinIOBackend 创建MinIO存储后端
func NewMinIOBackend(config *MinIOConfig) (*MinIOBackend, error) {
	// 1. 创建MinIO客户端
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// 2. 设置默认位置
	location := config.Location
	if location == "" {
		location = "us-east-1"
	}

	backend := &MinIOBackend{
		client:     client,
		bucketName: config.BucketName,
		useSSL:     config.UseSSL,
		location:   location,
	}

	// 3. 确保存储桶存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		// 创建存储桶
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{
			Region: location,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return backend, nil
}

// Save 保存文件到MinIO
func (b *MinIOBackend) Save(ctx context.Context, path string, reader io.Reader) error {
	// 上传文件到MinIO
	_, err := b.client.PutObject(
		ctx,
		b.bucketName,
		path,
		reader,
		-1, // 未知文件大小
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to save file to MinIO: %w", err)
	}

	return nil
}

// SaveWithSize 保存文件到MinIO（指定文件大小）
func (b *MinIOBackend) SaveWithSize(ctx context.Context, path string, reader io.Reader, size int64, contentType string) error {
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	} else {
		opts.ContentType = "application/octet-stream"
	}

	_, err := b.client.PutObject(
		ctx,
		b.bucketName,
		path,
		reader,
		size,
		opts,
	)
	if err != nil {
		return fmt.Errorf("failed to save file to MinIO: %w", err)
	}

	return nil
}

// Load 从MinIO加载文件
func (b *MinIOBackend) Load(ctx context.Context, path string) (io.ReadCloser, error) {
	object, err := b.client.GetObject(ctx, b.bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to load file from MinIO: %w", err)
	}

	// 检查对象是否存在
	_, err = object.Stat()
	if err != nil {
		object.Close()
		return nil, fmt.Errorf("file not found in MinIO: %w", err)
	}

	return object, nil
}

// Delete 从MinIO删除文件
func (b *MinIOBackend) Delete(ctx context.Context, path string) error {
	err := b.client.RemoveObject(ctx, b.bucketName, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (b *MinIOBackend) Exists(ctx context.Context, path string) (bool, error) {
	_, err := b.client.StatObject(ctx, b.bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		// 检查是否是"对象不存在"错误
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GetURL 生成预签名下载URL
func (b *MinIOBackend) GetURL(ctx context.Context, path string, expiresIn time.Duration) (string, error) {
	// 生成预签名URL
	url, err := b.client.PresignedGetObject(ctx, b.bucketName, path, expiresIn, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// GetUploadURL 生成预签名上传URL
func (b *MinIOBackend) GetUploadURL(ctx context.Context, path string, expiresIn time.Duration) (string, error) {
	url, err := b.client.PresignedPutObject(ctx, b.bucketName, path, expiresIn)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return url.String(), nil
}

// GetFileInfo 获取文件信息
func (b *MinIOBackend) GetFileInfo(ctx context.Context, path string) (*minio.ObjectInfo, error) {
	info, err := b.client.StatObject(ctx, b.bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &info, nil
}

// ListObjects 列出对象
func (b *MinIOBackend) ListObjects(ctx context.Context, prefix string, recursive bool) <-chan minio.ObjectInfo {
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}
	return b.client.ListObjects(ctx, b.bucketName, opts)
}

// CopyObject 复制对象
func (b *MinIOBackend) CopyObject(ctx context.Context, srcPath, destPath string) error {
	src := minio.CopySrcOptions{
		Bucket: b.bucketName,
		Object: srcPath,
	}

	dest := minio.CopyDestOptions{
		Bucket: b.bucketName,
		Object: destPath,
	}

	_, err := b.client.CopyObject(ctx, dest, src)
	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	return nil
}

// ============ 分片上传支持 ============

// InitiateMultipartUpload 初始化分片上传
func (b *MinIOBackend) InitiateMultipartUpload(ctx context.Context, path string, contentType string) (string, error) {
	// MinIO的分片上传是自动处理的，这里返回一个上传ID用于跟踪
	// 实际的分片上传逻辑在PutObject中自动处理
	return fmt.Sprintf("upload_%s_%d", path, time.Now().Unix()), nil
}

// UploadPart 上传文件分片
// Note: MinIO会自动处理分片上传，通常不需要手动分片
func (b *MinIOBackend) UploadPart(ctx context.Context, path string, partNumber int, reader io.Reader, size int64) error {
	// MinIO SDK会自动处理分片，这里使用标准的PutObject
	// 对于大文件，SDK会自动进行分片上传
	return b.SaveWithSize(ctx, path, reader, size, "")
}

// CompleteMultipartUpload 完成分片上传
func (b *MinIOBackend) CompleteMultipartUpload(ctx context.Context, uploadID string) error {
	// MinIO自动处理，这里只是一个占位符
	return nil
}

// AbortMultipartUpload 中止分片上传
func (b *MinIOBackend) AbortMultipartUpload(ctx context.Context, uploadID string) error {
	// MinIO自动处理，这里只是一个占位符
	return nil
}

// ============ 高级功能 ============

// SetBucketPolicy 设置存储桶策略
func (b *MinIOBackend) SetBucketPolicy(ctx context.Context, policy string) error {
	err := b.client.SetBucketPolicy(ctx, b.bucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}
	return nil
}

// GetBucketPolicy 获取存储桶策略
func (b *MinIOBackend) GetBucketPolicy(ctx context.Context) (string, error) {
	policy, err := b.client.GetBucketPolicy(ctx, b.bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket policy: %w", err)
	}
	return policy, nil
}

// SetObjectMetadata 设置对象元数据
func (b *MinIOBackend) SetObjectMetadata(ctx context.Context, path string, metadata map[string]string) error {
	// 通过复制对象来更新元数据
	src := minio.CopySrcOptions{
		Bucket: b.bucketName,
		Object: path,
	}

	dest := minio.CopyDestOptions{
		Bucket:          b.bucketName,
		Object:          path,
		UserMetadata:    metadata,
		ReplaceMetadata: true,
	}

	_, err := b.client.CopyObject(ctx, dest, src)
	if err != nil {
		return fmt.Errorf("failed to set object metadata: %w", err)
	}

	return nil
}

// GetClient 获取MinIO客户端（用于高级操作）
func (b *MinIOBackend) GetClient() *minio.Client {
	return b.client
}

// HealthCheck 健康检查
func (b *MinIOBackend) HealthCheck(ctx context.Context) error {
	// 检查存储桶是否存在
	exists, err := b.client.BucketExists(ctx, b.bucketName)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	if !exists {
		return fmt.Errorf("bucket %s does not exist", b.bucketName)
	}
	return nil
}
