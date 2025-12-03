package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql          MysqlConfig
	AvatarRpc      zrpc.RpcClientConf
	WorldRpc       zrpc.RpcClientConf
	ActionSchedule ActionScheduleConfig
}

type MysqlConfig struct {
	DataSource string
}

type ActionScheduleConfig struct {
	MinIntervalHours int
	MaxIntervalHours int
}
