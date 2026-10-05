package testhierarchy

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var baseURL = "http://localhost:9090/api/v1"

// Run performs the hierarchy API integration flow when explicitly enabled.
func Run(args []string) error {
	apply, err := parseApply(args)
	if err != nil {
		return err
	}
	if !apply {
		fmt.Println("预览模式：test_hierarchy 将注册用户并创建项目/文档；使用 --apply 才执行")
		return nil
	}
	if configuredURL := os.Getenv("QINGYU_TEST_BASE_URL"); configuredURL != "" {
		baseURL = configuredURL
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 等待服务启动
	fmt.Println("等待服务启动...")
	time.Sleep(3 * time.Second)

	// 1. 注册并登录
	fmt.Println("\n=== 1. 注册并登录 ===")
	token, userID, err := registerAndLogin("testuser2", "test2@example.com", "Password123!")
	if err != nil {
		return fmt.Errorf("注册/登录失败: %w", err)
	}
	fmt.Printf("✓ 登录成功, userID: %s\n", userID)

	// 2. 创建项目
	fmt.Println("\n=== 2. 创建项目 ===")
	projectID, err := createProject(token, "层级测试项目", "测试大纲层级结构")
	if err != nil {
		return fmt.Errorf("创建项目失败: %w", err)
	}
	fmt.Printf("✓ 项目创建成功, ID: %s\n", projectID)

	// 3. 创建第一个卷
	fmt.Println("\n=== 3. 创建第一个卷 ===")
	volume1DocID, err := createVolume(token, projectID, "卷一", 1)
	if err != nil {
		return fmt.Errorf("创建卷一失败: %w", err)
	}
	fmt.Printf("✓ 卷一创建成功, DocumentID: %s\n", volume1DocID)

	// 4. 创建第二个卷
	fmt.Println("\n=== 4. 创建第二个卷 ===")
	volume2DocID, err := createVolume(token, projectID, "卷二", 2)
	if err != nil {
		return fmt.Errorf("创建卷二失败: %w", err)
	}
	fmt.Printf("✓ 卷二创建成功, DocumentID: %s\n", volume2DocID)

	// 5. 在卷一下创建章
	fmt.Println("\n=== 5. 在卷一下创建第一章 ===")
	chapter1DocID, err := createChapter(token, projectID, volume1DocID, "第一章 开始", 1)
	if err != nil {
		return fmt.Errorf("创建第一章失败: %w", err)
	}
	fmt.Printf("✓ 第一章创建成功, DocumentID: %s\n", chapter1DocID)

	// 6. 在卷一下创建第二章
	fmt.Println("\n=== 6. 在卷一下创建第二章 ===")
	chapter2DocID, err := createChapter(token, projectID, volume1DocID, "第二章 发展", 2)
	if err != nil {
		return fmt.Errorf("创建第二章失败: %w", err)
	}
	fmt.Printf("✓ 第二章创建成功, DocumentID: %s\n", chapter2DocID)

	// 7. 获取大纲树验证层级结构
	fmt.Println("\n=== 7. 验证大纲树结构 ===")
	tree, err := getOutlineTree(token, projectID)
	if err != nil {
		return fmt.Errorf("获取大纲树失败: %w", err)
	}
	// tree["data"] is the array of root nodes
	if data, ok := tree["data"].([]interface{}); ok {
		fmt.Printf("根节点数量: %d\n", len(data))
		for _, node := range data {
			printTree(node, 0)
		}
	} else {
		fmt.Printf("树数据: %+v\n", tree["data"])
	}

	// 8. 验证数据库中的结构
	fmt.Println("\n=== 8. 验证数据库结构 ===")
	verifyStructure(ctx)

	fmt.Println("\n=== 测试完成 ===")
	return nil
}

func parseApply(args []string) (bool, error) {
	fs := flag.NewFlagSet("test_hierarchy", flag.ContinueOnError)
	apply := fs.Bool("apply", false, "执行真实 API 集成流程")
	if err := fs.Parse(args); err != nil {
		return false, err
	}
	return *apply, nil
}

func registerAndLogin(username, email, password string) (string, string, error) {
	// 先尝试注册
	registerData := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}
	resp, err := doPost("/shared/auth/register", registerData, "")
	if err != nil {
		return "", "", err
	}
	fmt.Printf("注册响应: %s\n", string(resp))

	var regResult map[string]interface{}
	if err := json.Unmarshal(resp, &regResult); err != nil {
		return "", "", err
	}

	// 然后登录
	loginData := map[string]string{
		"username": username,
		"password": password,
	}
	resp, err = doPost("/shared/auth/login", loginData, "")
	if err != nil {
		return "", "", err
	}
	fmt.Printf("登录响应: %s\n", string(resp))

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", "", err
	}

	if err := responseCode(result); err != nil {
		return "", "", fmt.Errorf("登录失败: %w", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("登录响应缺少 data")
	}
	tokenStr, ok := data["token"].(string)
	if !ok || tokenStr == "" {
		return "", "", fmt.Errorf("登录响应缺少 token")
	}
	user, ok := data["user"].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("登录响应缺少 user")
	}
	userID, ok := user["id"].(string)
	if !ok || userID == "" {
		return "", "", fmt.Errorf("登录响应缺少 user.id")
	}
	return tokenStr, userID, nil
}

func createProject(token, title, description string) (string, error) {
	data := map[string]string{
		"title":       title,
		"description": description,
	}
	resp, err := doPost("/writer/projects", data, token)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	if err := responseCode(result); err != nil {
		return "", fmt.Errorf("创建项目失败: %w", err)
	}
	project, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("创建项目响应缺少 data")
	}
	projectID, ok := project["id"].(string)
	if !ok || projectID == "" {
		return "", fmt.Errorf("创建项目响应缺少 id")
	}
	return projectID, nil
}

func createVolume(token, projectID, title string, order int) (string, error) {
	data := map[string]interface{}{
		"title":    title,
		"type":     "volume",
		"level":    1,
		"parentId": "",
		"order":    order,
		"status":   "planned",
	}
	resp, err := doPost(fmt.Sprintf("/writer/project/%s/documents", projectID), data, token)
	if err != nil {
		return "", err
	}
	fmt.Printf("创建卷响应: %s\n", string(resp))

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	if err := responseCode(result); err != nil {
		return "", fmt.Errorf("创建卷失败: %w", err)
	}

	doc, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("创建卷响应缺少 data")
	}
	docID, ok := doc["documentId"].(string)
	if !ok || docID == "" {
		return "", fmt.Errorf("创建卷响应缺少 documentId")
	}

	return docID, nil
}

func createChapter(token, projectID, parentID, title string, order int) (string, error) {
	data := map[string]interface{}{
		"title":    title,
		"type":     "chapter",
		"level":    2,
		"parentId": parentID,
		"order":    order,
		"status":   "planned",
	}
	resp, err := doPost(fmt.Sprintf("/writer/project/%s/documents", projectID), data, token)
	if err != nil {
		return "", err
	}
	fmt.Printf("创建章响应: %s\n", string(resp))

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}

	if err := responseCode(result); err != nil {
		return "", fmt.Errorf("创建章失败: %w", err)
	}

	doc, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("创建章响应缺少 data")
	}
	docID, ok := doc["documentId"].(string)
	if !ok || docID == "" {
		return "", fmt.Errorf("创建章响应缺少 documentId")
	}

	return docID, nil
}

func getOutlineTree(token, projectID string) (map[string]interface{}, error) {
	resp, err := doGet(fmt.Sprintf("/writer/projects/%s/outlines/tree", projectID), token)
	if err != nil {
		return nil, err
	}
	fmt.Printf("大纲树响应: %s\n", string(resp))

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func printTree(node interface{}, level int) {
	switch v := node.(type) {
	case []interface{}:
		for _, item := range v {
			printTree(item, level)
		}
	case map[string]interface{}:
		indent := ""
		for i := 0; i < level; i++ {
			indent += "  "
		}
		title, _ := v["title"].(string)
		nodeType, _ := v["type"].(string)
		docID, _ := v["documentId"].(string)
		if docID == "" {
			docID, _ = v["document_id"].(string)
		}
		parentID, _ := v["parentId"].(string)
		if parentID == "" {
			parentID, _ = v["parent_id"].(string)
		}
		fmt.Printf("%s- [%s] %s (doc: %s, parent: %s)\n", indent, nodeType, title, docID, parentID)

		if v["children"] != nil {
			printTree(v["children"], level+1)
		}
	}
}

func doPost(path string, data interface{}, token string) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("编码请求失败: %w", err)
	}
	req, err := http.NewRequest("POST", baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}
	return responseBody, nil
}

func doGet(path string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}
	return responseBody, nil
}

func verifyStructure(ctx context.Context) {
	_ = ctx
	fmt.Println("数据库结构复核未实现，已跳过；详见 MODULE.md TODO")
}

func responseCode(result map[string]interface{}) error {
	value, exists := result["code"]
	if !exists || value == nil {
		return nil
	}
	code, ok := value.(float64)
	if !ok {
		return fmt.Errorf("响应 code 类型无效: %T", value)
	}
	if code != 0 {
		return fmt.Errorf("响应 code=%v", code)
	}
	return nil
}
