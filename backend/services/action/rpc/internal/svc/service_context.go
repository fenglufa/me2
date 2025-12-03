package svc

import (
	"github.com/me2/action/rpc/internal/config"
	"github.com/me2/action/rpc/internal/model"
	"github.com/me2/avatar/rpc/avatar_client"
	"github.com/me2/world/rpc/world_client"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	ActionLogModel *model.ActionLogModel
	AvatarRpc      avatar_client.Avatar
	WorldRpc       world_client.World
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		ActionLogModel: model.NewActionLogModel(conn),
		AvatarRpc:      avatar_client.NewAvatar(zrpc.MustNewClient(c.AvatarRpc)),
		WorldRpc:       world_client.NewWorld(zrpc.MustNewClient(c.WorldRpc)),
	}
}
