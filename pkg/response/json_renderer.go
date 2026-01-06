package response

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JsonWithNoEscape 自定义JSON渲染器，不转义非ASCII字符
// 解决中文字符被转义为 \uXXXX 格式的问题
func JsonWithNoEscape(c *gin.Context, code int, obj any) {
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	encoder.SetEscapeHTML(false) // 关键：禁用HTML转义，保持中文字符原样输出
	if err := encoder.Encode(obj); err != nil {
		// 如果编码失败，返回错误
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// JSONWithNoEscapeString 将对象编码为不转义的JSON字符串
func JSONWithNoEscapeString(obj any) (string, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(obj); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// ParseJSONUnescape 解析JSON并反转义HTML实体
func ParseJSONUnescape(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

// UnescapeHTML 反转义HTML实体
func UnescapeHTML(s string) string {
	return html.UnescapeString(s)
}

// ReadJSONWithNoEscape 从io.Reader读取JSON并反转义
func ReadJSONWithNoEscape(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	return JSONWithNoEscapeString(result)
}
