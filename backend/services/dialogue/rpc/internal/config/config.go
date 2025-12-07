package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql     MysqlConfig
	AvatarRpc zrpc.RpcClientConf
	EventRpc  zrpc.RpcClientConf
	AiRpc     zrpc.RpcClientConf
	Context   ContextConfig
}

type MysqlConfig struct {
	DataSource string
}

type ContextConfig struct {
	MaxHistoryMessages int32
	MaxRecentEvents    int32
}
