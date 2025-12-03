package svc

import (
	"github.com/me2/world/rpc/internal/config"
	"github.com/me2/world/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	// Models
	WorldMapModel    *model.WorldMapModel
	WorldRegionModel *model.WorldRegionModel
	WorldSceneModel  *model.WorldSceneModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:           c,
		WorldMapModel:    model.NewWorldMapModel(conn),
		WorldRegionModel: model.NewWorldRegionModel(conn),
		WorldSceneModel:  model.NewWorldSceneModel(conn),
	}
}
