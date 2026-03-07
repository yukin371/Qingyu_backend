// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	yamlPath := "docs/api/swagger.yaml"
	jsonPath := "docs/api/swagger.json"

	// Check if YAML file exists
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		log.Fatalf("YAML文件不存在: %s", yamlPath)
	}

	// Read YAML file
	data, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("读取YAML文件失败: %v", err)
	}

	// Parse YAML
	var m interface{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("解析YAML失败: %v", err)
	}

	// Convert to JSON with proper formatting
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatalf("转换为JSON失败: %v", err)
	}

	// Write JSON file
	err = ioutil.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("写入JSON文件失败: %v", err)
	}

	// Get file info
	yamlInfo, _ := os.Stat(yamlPath)
	jsonInfo, _ := os.Stat(jsonPath)

	fmt.Println("✅ Swagger转换完成!")
	fmt.Printf("   来源: %s (%.1f KB)\n", yamlPath, float64(yamlInfo.Size())/1024)
	fmt.Printf("   输出: %s (%.1f KB)\n", jsonPath, float64(jsonInfo.Size())/1024)
	fmt.Println("")
	fmt.Println("导入方式:")
	fmt.Println("  Postman: Import -> Upload Files -> 选择 docs/api/swagger.json")
	fmt.Println("  Apifox:  项目设置 -> 导入数据 -> OpenAPI/Swagger -> 选择 docs/api/swagger.yaml")
}
