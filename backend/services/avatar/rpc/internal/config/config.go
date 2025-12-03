package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	OssRpc       zrpc.RpcClientConf
	SchedulerRpc zrpc.RpcClientConf
	MachineID    int64
}
