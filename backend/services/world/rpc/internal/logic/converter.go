package logic

import (
	"github.com/me2/world/rpc/internal/model"
	"github.com/me2/world/rpc/world"
)

// ModelToProtoMap 将 Model 的 WorldMap 转换为 Proto 的 WorldMap
func ModelToProtoMap(m *model.WorldMap) *world.WorldMap {
	return &world.WorldMap{
		Id:          m.Id,
		Name:        m.Name,
		Description: m.Description,
		CoverImage:  m.CoverImage,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt.Unix(),
		UpdatedAt:   m.UpdatedAt.Unix(),
	}
}

// ModelToProtoRegion 将 Model 的 WorldRegion 转换为 Proto 的 WorldRegion
func ModelToProtoRegion(m *model.WorldRegion) *world.WorldRegion {
	return &world.WorldRegion{
		Id:          m.Id,
		MapId:       m.MapId,
		Name:        m.Name,
		Description: m.Description,
		CoverImage:  m.CoverImage,
		Atmosphere:  m.Atmosphere,
		Tags:        m.Tags,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt.Unix(),
		UpdatedAt:   m.UpdatedAt.Unix(),
	}
}

// ModelToProtoScene 将 Model 的 WorldScene 转换为 Proto 的 WorldScene
func ModelToProtoScene(m *model.WorldScene) *world.WorldScene {
	return &world.WorldScene{
		Id:              m.Id,
		RegionId:        m.RegionId,
		Name:            m.Name,
		Description:     m.Description,
		CoverImage:      m.CoverImage,
		Atmosphere:      m.Atmosphere,
		Tags:            m.Tags,
		SuitableActions: m.SuitableActions,
		Capacity:        m.Capacity,
		IsActive:        m.IsActive,
		Features: &world.SceneFeatures{
			HasWifi:      m.HasWifi,
			HasFood:      m.HasFood,
			HasSeating:   m.HasSeating,
			IsIndoor:     m.IsIndoor,
			IsQuiet:      m.IsQuiet,
			ComfortLevel: m.ComfortLevel,
			SocialLevel:  m.SocialLevel,
		},
		CreatedAt: m.CreatedAt.Unix(),
		UpdatedAt: m.UpdatedAt.Unix(),
	}
}
