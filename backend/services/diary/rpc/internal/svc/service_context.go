package svc

import (
	"github.com/me2/ai/rpc/ai_client"
	"github.com/me2/avatar/rpc/avatar_client"
	"github.com/me2/diary/rpc/internal/config"
	"github.com/me2/diary/rpc/internal/model"
	"github.com/me2/event/rpc/event_client"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	DiaryModel  model.DiariesModel
	EventRpc    event_client.Event
	AIRpc       ai_client.Ai
	AvatarRpc   avatar_client.Avatar
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// 打印 AI RPC 的超时配置
	logx.Infof("AI RPC Timeout配置: %d ms", c.AIRpc.Timeout)

	return &ServiceContext{
		Config:     c,
		DiaryModel: model.NewDiariesModel(conn),
		EventRpc:   event_client.NewEvent(zrpc.MustNewClient(c.EventRpc)),
		AIRpc:      ai_client.NewAi(zrpc.MustNewClient(c.AIRpc)),
		AvatarRpc:  avatar_client.NewAvatar(zrpc.MustNewClient(c.AvatarRpc)),
	}
}
