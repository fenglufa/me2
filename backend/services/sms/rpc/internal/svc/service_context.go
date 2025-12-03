package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/me2/sms/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis: redis.MustNewRedis(c.RedisConf),
	}
}
