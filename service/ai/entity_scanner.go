package ai

import (
	"context"
	"strings"

	writerRepo "Qingyu_backend/repository/interfaces/writer"
)

// EntityScanner 实体扫描器
// 基于已有实体库进行文本匹配，自动关联文档与实体
type EntityScanner struct {
	characterRepo writerRepo.CharacterRepository
	locationRepo  writerRepo.LocationRepository
}

func NewEntityScanner(
	characterRepo writerRepo.CharacterRepository,
	locationRepo writerRepo.LocationRepository,
) *EntityScanner {
	return &EntityScanner{
		characterRepo: characterRepo,
		locationRepo:  locationRepo,
	}
}

// ScanResult 扫描结果
type ScanResult struct {
	Characters []MatchedEntity `json:"characters"`
	Locations  []MatchedEntity `json:"locations"`
}

// MatchedEntity 匹配到的实体
type MatchedEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ScanContent 扫描文本内容，匹配已有实体
func (s *EntityScanner) ScanContent(ctx context.Context, projectID string, text string) (*ScanResult, error) {
	result := &ScanResult{}

	// 获取并匹配角色
	characters, err := s.characterRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, char := range characters {
		if strings.Contains(text, char.Name) {
			result.Characters = append(result.Characters, MatchedEntity{
				ID:   char.ID.Hex(),
				Name: char.Name,
			})
			continue
		}
		for _, alias := range char.Alias {
			if strings.Contains(text, alias) {
				result.Characters = append(result.Characters, MatchedEntity{
					ID:   char.ID.Hex(),
					Name: char.Name,
				})
				break
			}
		}
	}

	// 获取并匹配地点
	locations, err := s.locationRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, loc := range locations {
		if strings.Contains(text, loc.Name) {
			result.Locations = append(result.Locations, MatchedEntity{
				ID:   loc.ID.Hex(),
				Name: loc.Name,
			})
		}
	}

	return result, nil
}
