package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	AiRpc zrpc.RpcClientConf
	AI    struct {
		Enabled bool `json:",default=true"`
	}
}
