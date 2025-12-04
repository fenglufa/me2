package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	AvatarRpc           zrpc.RpcClientConf
	WorldRpc            zrpc.RpcClientConf
	AIRpc               zrpc.RpcClientConf
	AIGenerateTimeout   int64 // AI 生成超时时间（毫秒），默认 30000
}
