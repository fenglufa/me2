package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRpc   zrpc.RpcClientConf
	AvatarRpc zrpc.RpcClientConf
	WorldRpc  zrpc.RpcClientConf
	EventRpc  zrpc.RpcClientConf
	ActionRpc zrpc.RpcClientConf
	DiaryRpc  zrpc.RpcClientConf
}
