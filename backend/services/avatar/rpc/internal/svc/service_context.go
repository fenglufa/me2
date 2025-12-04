package svc

import (
	"github.com/me2/action/rpc/action_client"
	"github.com/me2/avatar/rpc/internal/config"
	"github.com/me2/avatar/rpc/internal/idgen"
	"github.com/me2/avatar/rpc/internal/model"
	"github.com/me2/oss/rpc/oss_client"
	"github.com/me2/scheduler/rpc/scheduler_client"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	DB           sqlx.SqlConn
	OssRpc       oss_client.Oss
	SchedulerRpc scheduler_client.Scheduler
	ActionRpc    action_client.Action
	IDGen        *idgen.Snowflake
	AvatarModel  *model.AvatarModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := sqlx.NewMysql(c.Mysql.DataSource)
	ossRpc := oss_client.NewOss(zrpc.MustNewClient(c.OssRpc))
	schedulerRpc := scheduler_client.NewScheduler(zrpc.MustNewClient(c.SchedulerRpc))
	actionRpc := action_client.NewAction(zrpc.MustNewClient(c.ActionRpc))

	idGen, err := idgen.NewSnowflake(c.MachineID)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:       c,
		DB:           db,
		OssRpc:       ossRpc,
		SchedulerRpc: schedulerRpc,
		ActionRpc:    actionRpc,
		IDGen:        idGen,
		AvatarModel:  model.NewAvatarModel(db),
	}
}
