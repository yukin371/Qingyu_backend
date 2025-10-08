package response

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// gzipWriterPool gzip writer池，减少内存分配
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(io.Discard)
	},
}

// gzipResponseWriter 实现Gin的ResponseWriter接口
type gzipResponseWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipResponseWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

// GzipMiddleware Gzip压缩中间件
func GzipMiddleware(level int) gin.HandlerFunc {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	return func(c *gin.Context) {
		// 检查客户端是否支持gzip
		if !shouldCompress(c.Request) {
			c.Next()
			return
		}

		// 从池中获取writer
		gz := gzipWriterPool.Get().(*gzip.Writer)
		defer gzipWriterPool.Put(gz)

		gz.Reset(c.Writer)
		defer gz.Close()

		// 设置响应头
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// 替换writer
		c.Writer = &gzipResponseWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		c.Next()
	}
}

// shouldCompress 判断是否应该压缩响应
func shouldCompress(req *http.Request) bool {
	// 检查Accept-Encoding头
	acceptEncoding := req.Header.Get("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		return false
	}

	// 检查Content-Type（某些类型不需要压缩）
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, "image/") ||
		strings.Contains(contentType, "video/") ||
		strings.Contains(contentType, "audio/") {
		return false
	}

	return true
}

// MinSizeForCompression 最小压缩大小（小于此大小不压缩）
const MinSizeForCompression = 1024 // 1KB

// ShouldCompressResponse 判断响应是否应该压缩
func ShouldCompressResponse(contentLength int64, contentType string) bool {
	// 检查大小
	if contentLength > 0 && contentLength < MinSizeForCompression {
		return false
	}

	// 检查内容类型
	compressibleTypes := []string{
		"application/json",
		"application/xml",
		"text/html",
		"text/plain",
		"text/css",
		"text/javascript",
		"application/javascript",
	}

	for _, t := range compressibleTypes {
		if strings.Contains(contentType, t) {
			return true
		}
	}

	return false
}
