package svc

import (
	"github.com/me2/dialogue/rpc/internal/config"
	"github.com/me2/dialogue/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/me2/avatar/rpc/avatar_client"
	"github.com/me2/event/rpc/event_client"
	"github.com/me2/ai/rpc/ai_client"
)

type ServiceContext struct {
	Config        config.Config
	SessionModel  model.DialogueSessionModel
	MessageModel  model.DialogueMessageModel
	AvatarRpc     avatar_client.Avatar
	EventRpc      event_client.Event
	AiRpc         ai_client.Ai
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:       c,
		SessionModel: model.NewDialogueSessionModel(conn),
		MessageModel: model.NewDialogueMessageModel(conn),
		AvatarRpc:    avatar_client.NewAvatar(zrpc.MustNewClient(c.AvatarRpc)),
		EventRpc:     event_client.NewEvent(zrpc.MustNewClient(c.EventRpc)),
		AiRpc:        ai_client.NewAi(zrpc.MustNewClient(c.AiRpc)),
	}
}
