package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OutlineNode struct {
	ID         string        `json:"id"`
	Title      string        `json:"title"`
	DocumentID string        `json:"documentId,omitempty"`
	Tags       []string      `json:"tags,omitempty"`
	Children   []OutlineNode `json:"children,omitempty"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    []OutlineNode `json:"data,omitempty"`
}

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	// 注意：这个API需要认证，这里仅作演示
	// 实际使用时需要添加认证token
	resp, err := client.Get("http://localhost:9090/api/v1/writer/projects/69c54227f79b1f9c6a17f3b9/outlines/tree")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// 先尝试直接解析为outline数组
	var outlines []OutlineNode
	if err := json.Unmarshal(body, &outlines); err == nil {
		checkTags(outlines, 0)
		return
	}

	// 尝试解析为APIResponse
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil {
		if apiResp.Code == 0 && len(apiResp.Data) > 0 {
			checkTags(apiResp.Data, 0)
			return
		}
	}

	fmt.Printf("响应: %s\n", string(body))
}

func checkTags(nodes []OutlineNode, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	for _, node := range nodes {
		tagsStatus := "❌ 无tags"
		if len(node.Tags) > 0 {
			tagsStatus = fmt.Sprintf("✓ tags: %v", node.Tags)
		}
		fmt.Printf("%s[%s] %s - %s\n", indent, node.ID, node.Title, tagsStatus)

		if len(node.Children) > 0 {
			checkTags(node.Children, depth+1)
		}
	}
}
