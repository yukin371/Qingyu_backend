package document

import "time"

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int      `json:"maxRetries" validate:"min=0,max=10"`      // 最大重试次数
	RetryDelay      int      `json:"retryDelay" validate:"min=100,max=60000"` // 重试延迟（毫秒）
	RetryableErrors []string `json:"retryableErrors"`                         // 可重试的错误码列表
}

// DefaultRetryConfig 返回默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 3,
		RetryDelay: 1000, // 1秒
		RetryableErrors: []string{
			"VERSION_CONFLICT",
			"NETWORK_ERROR",
			"TIMEOUT",
			"SERVICE_UNAVAILABLE",
		},
	}
}

// GetMaxDelay 获取最大延迟时间（60秒）
func (c *RetryConfig) GetMaxDelay() time.Duration {
	return 60 * time.Second
}

// RetryService 重试服务
type RetryService struct {
	// 可以添加logger字段
}

// NewRetryService 创建重试服务
func NewRetryService() *RetryService {
	return &RetryService{}
}

// ShouldRetry 判断错误是否可重试
func (s *RetryService) ShouldRetry(err error, config *RetryConfig) bool {
	if err == nil || config == nil {
		return false
	}

	// 从错误中提取错误码
	errorCode := s.extractErrorCode(err)

	// 检查是否在可重试列表中
	for _, retryableErr := range config.RetryableErrors {
		if errorCode == retryableErr {
			return true
		}
	}

	return false
}

// extractErrorCode 从错误中提取错误码
func (s *RetryService) extractErrorCode(err error) string {
	// TODO: 从具体错误类型中提取错误码
	// 当前简化实现：返回错误消息的某些部分
	// 后续可以根据实际错误类型实现
	return err.Error()
}

// GetRetryDelay 计算重试延迟时间（指数退避）
func (s *RetryService) GetRetryDelay(attempt int, config *RetryConfig) time.Duration {
	if config == nil || attempt < 0 {
		return 0
	}

	// 指数退避：delay = baseDelay * 2^attempt
	// baseDelay是毫秒，需要转换为time.Duration
	baseDelay := time.Duration(config.RetryDelay) * time.Millisecond
	delay := baseDelay * time.Duration(1<<uint(attempt))

	// 不超过最大延迟
	maxDelay := config.GetMaxDelay()
	if delay > maxDelay {
		return maxDelay
	}

	return delay
}

// CanRetry 判断是否还可以重试
func (s *RetryService) CanRetry(attempt int, config *RetryConfig) bool {
	if config == nil {
		return false
	}
	return attempt < config.MaxRetries
}
