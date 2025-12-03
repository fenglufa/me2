package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	RedisConf  redis.RedisConf
	Deepseek   DeepseekConfig
	PromptFile string // Prompt 模板配置文件路径
}

type DeepseekConfig struct {
	ApiKey  string
	BaseURL string
	Model   string
	Timeout int64
}
