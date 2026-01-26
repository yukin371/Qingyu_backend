package search

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// SearchIndicesConfig 搜索索引配置
type SearchIndicesConfig struct {
	Settings map[string]interface{} `yaml:"settings"`
	Indices  map[string]IndexConfig `yaml:"indices"`
}

// IndexConfig 单个索引配置
type IndexConfig struct {
	Alias        string                 `yaml:"alias"`
	Shards       int                    `yaml:"shards"`
	Replicas     int                    `yaml:"replicas"`
	Mappings     map[string]interface{} `yaml:"mappings"`
	Settings     map[string]interface{} `yaml:"settings"`
}

// BuildIndexMapping 构建索引映射
func (ic *IndexConfig) BuildIndexMapping(globalSettings map[string]interface{}) map[string]interface{} {
	mapping := make(map[string]interface{})

	// 合并全局设置
	if globalSettings != nil {
		for k, v := range globalSettings {
			mapping[k] = v
		}
	}

	// 合并索引级别的设置
	if ic.Settings != nil {
		for k, v := range ic.Settings {
			mapping[k] = v
		}
	}

	// 添加 mappings
	if ic.Mappings != nil {
		mapping["mappings"] = ic.Mappings
	}

	// 添加 settings
	settings := make(map[string]interface{})
	if shards, ok := mapping["number_of_shards"]; ok {
		settings["number_of_shards"] = shards
	} else {
		settings["number_of_shards"] = ic.Shards
	}

	if replicas, ok := mapping["number_of_replicas"]; ok {
		settings["number_of_replicas"] = replicas
	} else {
		settings["number_of_replicas"] = ic.Replicas
	}

	mapping["settings"] = settings

	return mapping
}

// LoadSearchIndicesConfig 从文件加载搜索索引配置
func LoadSearchIndicesConfig(configPath string) (*SearchIndicesConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &SearchIndicesConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetDefaultIndicesConfig 获取默认索引配置
func GetDefaultIndicesConfig() *SearchIndicesConfig {
	return &SearchIndicesConfig{
		Settings: map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 1,
		},
		Indices: map[string]IndexConfig{
			"books": {
				Alias:    "books",
				Shards:   1,
				Replicas: 1,
				Mappings: map[string]interface{}{
					"properties": map[string]interface{}{
						"title":       map[string]string{"type": "text"},
						"author":      map[string]string{"type": "text"},
						"description": map[string]string{"type": "text"},
						"tags":        map[string]string{"type": "keyword"},
					},
				},
			},
		},
	}
}
