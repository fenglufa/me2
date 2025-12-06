package svc

import (
	"net/http"

	"github.com/me2/action/rpc/action_client"
	"github.com/me2/avatar/rpc/avatar_client"
	"github.com/me2/diary/rpc/diary_client"
	"github.com/me2/event/rpc/event_client"
	"github.com/me2/gateway/api/internal/config"
	"github.com/me2/gateway/api/internal/middleware"
	"github.com/me2/user/rpc/user_client"
	"github.com/me2/world/rpc/world_client"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Auth      func(http.HandlerFunc) http.HandlerFunc
	UserRpc   user_client.User
	AvatarRpc avatar_client.Avatar
	WorldRpc  world_client.World
	EventRpc  event_client.Event
	ActionRpc action_client.Action
	DiaryRpc  diary_client.Diary
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Auth:      middleware.NewAuthMiddleware(c.Auth.AccessSecret),
		UserRpc:   user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
		AvatarRpc: avatar_client.NewAvatar(zrpc.MustNewClient(c.AvatarRpc)),
		WorldRpc:  world_client.NewWorld(zrpc.MustNewClient(c.WorldRpc)),
		EventRpc:  event_client.NewEvent(zrpc.MustNewClient(c.EventRpc)),
		ActionRpc: action_client.NewAction(zrpc.MustNewClient(c.ActionRpc)),
		DiaryRpc:  diary_client.NewDiary(zrpc.MustNewClient(c.DiaryRpc)),
	}
}
