package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql          MysqlConfig
	ActionRpc      zrpc.RpcClientConf
	Schedule       ScheduleConfig
	ActionSchedule ActionScheduleConfig
}

type MysqlConfig struct {
	DataSource string
}

type ScheduleConfig struct {
	ScanInterval int // 扫描间隔（秒）
	MaxWorkers   int // 最大并发数
}

type ActionScheduleConfig struct {
	MinIntervalHours int // 最小间隔小时数
	MaxIntervalHours int // 最大间隔小时数
}
