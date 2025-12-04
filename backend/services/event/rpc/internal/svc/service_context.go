package svc

import (
	"github.com/me2/ai/rpc/ai_client"
	"github.com/me2/avatar/rpc/avatar_client"
	"github.com/me2/event/rpc/internal/config"
	"github.com/me2/event/rpc/internal/model"
	"github.com/me2/world/rpc/world_client"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	DB                 sqlx.SqlConn
	EventTemplateModel *model.EventTemplateModel
	EventHistoryModel  *model.EventHistoryModel
	AvatarRpc          avatar_client.Avatar
	WorldRpc           world_client.World
	AIRpc              ai_client.Ai
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := sqlx.NewMysql(c.Mysql.DataSource)
	avatarRpc := avatar_client.NewAvatar(zrpc.MustNewClient(c.AvatarRpc))
	worldRpc := world_client.NewWorld(zrpc.MustNewClient(c.WorldRpc))
	aiRpc := ai_client.NewAi(zrpc.MustNewClient(c.AIRpc))

	return &ServiceContext{
		Config:             c,
		DB:                 db,
		EventTemplateModel: model.NewEventTemplateModel(db),
		EventHistoryModel:  model.NewEventHistoryModel(db),
		AvatarRpc:          avatarRpc,
		WorldRpc:           worldRpc,
		AIRpc:              aiRpc,
	}
}
