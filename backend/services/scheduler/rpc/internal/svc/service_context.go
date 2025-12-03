package svc

import (
	"github.com/me2/action/rpc/action_client"
	"github.com/me2/scheduler/rpc/internal/config"
	"github.com/me2/scheduler/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	AvatarScheduleModel *model.AvatarScheduleModel
	ActionRpc           action_client.Action
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:              c,
		AvatarScheduleModel: model.NewAvatarScheduleModel(conn),
		ActionRpc:           action_client.NewAction(zrpc.MustNewClient(c.ActionRpc)),
	}
}
