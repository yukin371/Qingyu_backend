package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id"`
}

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("http://localhost:9090/api/v1/writer/projects/69c54227f79b1f9c6a17f3b9/outlines/tree")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}

	fmt.Printf("=== 原始响应 ===\n%s\n\n", string(body))

	// 格式化JSON
	var formatted interface{}
	if err := json.Unmarshal(body, &formatted); err == nil {
		pretty, _ := json.MarshalIndent(formatted, "", "  ")
		fmt.Printf("=== 格式化JSON ===\n%s\n\n", string(pretty))
	}

	// 解析为APIResponse
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Printf("解析APIResponse失败: %v\n", err)
		return
	}

	fmt.Printf("=== API响应结构 ===\n")
	fmt.Printf("Code: %d\n", apiResp.Code)
	fmt.Printf("Message: %s\n", apiResp.Message)
	fmt.Printf("RequestID: %s\n", apiResp.RequestID)

	// 检查Data字段的类型
	switch v := apiResp.Data.(type) {
	case []interface{}:
		fmt.Printf("Data类型: []interface{}, 长度: %d\n", len(v))
		if len(v) > 0 {
			fmt.Printf("第一个元素类型: %T\n", v[0])
			if m, ok := v[0].(map[string]interface{}); ok {
				fmt.Printf("第一个元素的字段:\n")
				for key, val := range m {
					switch v2 := val.(type) {
					case []interface{}:
						fmt.Printf("  %s: []interface{} (长度: %d)\n", key, len(v2))
					case string:
						fmt.Printf("  %s: %s\n", key, v2)
					default:
						fmt.Printf("  %s: %T\n", key, v2)
					}
				}
			}
		}
	case map[string]interface{}:
		fmt.Printf("Data类型: map[string]interface{}\n")
		for key, val := range v {
			fmt.Printf("  %s: %T\n", key, val)
		}
	default:
		fmt.Printf("Data类型: %T\n", v)
	}
}
